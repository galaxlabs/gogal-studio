package doctype

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SyncResult struct {
	FilePath string `json:"file_path"`
	DocType  string `json:"doctype"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

type Syncer struct {
	DB *pgxpool.Pool
}

func NewSyncer(db *pgxpool.Pool) *Syncer {
	return &Syncer{DB: db}
}

func (s *Syncer) SyncAll(rootPath string) ([]SyncResult, error) {
	files, err := findDocTypeJSONFiles(rootPath)
	if err != nil {
		return nil, err
	}

	results := make([]SyncResult, 0, len(files))

	for _, filePath := range files {
		result := SyncResult{
			FilePath: filePath,
			Status:   "pending",
		}

		doc, raw, err := readDocTypeJSON(filePath)
		if err != nil {
			result.Status = "failed"
			result.Message = err.Error()
			results = append(results, result)
			continue
		}

		result.DocType = doc.Name

		hash := hashBytes(raw)

		if err := s.SyncOne(filePath, hash, doc); err != nil {
			result.Status = "failed"
			result.Message = err.Error()
			results = append(results, result)
			continue
		}

		result.Status = "success"
		result.Message = "synced"
		results = append(results, result)
	}

	return results, nil
}

func (s *Syncer) SyncOne(jsonPath string, jsonHash string, doc DocTypeJSON) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if doc.Name == "" {
		return fmt.Errorf("doctype name is required")
	}

	if doc.Module == "" {
		return fmt.Errorf("module is required for doctype %s", doc.Name)
	}

	if doc.Label == "" {
		doc.Label = doc.Name
	}

	if doc.TableName == "" {
		return fmt.Errorf("table_name is required for doctype %s", doc.Name)
	}

	if doc.NamingRule == "" {
		doc.NamingRule = "autoname"
	}

	if doc.SortField == "" {
		doc.SortField = "created_at"
	}

	if doc.SortOrder == "" {
		doc.SortOrder = "DESC"
	}

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO core_module_def (
			name,
			label,
			is_core,
			enabled,
			updated_at
		)
		VALUES ($1, $2, $3, TRUE, NOW())
		ON CONFLICT (name)
		DO UPDATE SET
			label = EXCLUDED.label,
			is_core = EXCLUDED.is_core,
			updated_at = NOW()
	`, doc.Module, doc.Module, doc.IsCore)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO core_doctype (
			name,
			module,
			label,
			table_name,
			is_core,
			is_single,
			is_submittable,
			is_tree,
			is_child_table,
			allow_import,
			allow_export,
			track_changes,
			editable_grid,
			quick_entry,
			controller,
			route,
			naming_rule,
			title_field,
			image_field,
			sort_field,
			sort_order,
			json_hash,
			source_path,
			updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20,
			$21, $22, $23, NOW()
		)
		ON CONFLICT (name)
		DO UPDATE SET
			module = EXCLUDED.module,
			label = EXCLUDED.label,
			table_name = EXCLUDED.table_name,
			is_core = EXCLUDED.is_core,
			is_single = EXCLUDED.is_single,
			is_submittable = EXCLUDED.is_submittable,
			is_tree = EXCLUDED.is_tree,
			is_child_table = EXCLUDED.is_child_table,
			allow_import = EXCLUDED.allow_import,
			allow_export = EXCLUDED.allow_export,
			track_changes = EXCLUDED.track_changes,
			editable_grid = EXCLUDED.editable_grid,
			quick_entry = EXCLUDED.quick_entry,
			controller = EXCLUDED.controller,
			route = EXCLUDED.route,
			naming_rule = EXCLUDED.naming_rule,
			title_field = EXCLUDED.title_field,
			image_field = EXCLUDED.image_field,
			sort_field = EXCLUDED.sort_field,
			sort_order = EXCLUDED.sort_order,
			json_hash = EXCLUDED.json_hash,
			source_path = EXCLUDED.source_path,
			updated_at = NOW()
	`, doc.Name, doc.Module, doc.Label, doc.TableName,
		doc.IsCore,
		doc.IsSingle,
		doc.IsSubmittable,
		doc.IsTree,
		doc.IsChildTable,
		doc.AllowImport,
		doc.AllowExport,
		doc.TrackChanges,
		doc.EditableGrid,
		doc.QuickEntry,
		doc.Controller,
		doc.Route,
		doc.NamingRule,
		doc.TitleField,
		doc.ImageField,
		doc.SortField,
		doc.SortOrder,
		jsonHash,
		jsonPath,
	)
	if err != nil {
		return err
	}

	if err := syncFields(ctx, tx, doc); err != nil {
		return err
	}

	if err := syncActions(ctx, tx, doc); err != nil {
		return err
	}

	if err := syncLinks(ctx, tx, doc); err != nil {
		return err
	}

	if err := syncPermissions(ctx, tx, doc); err != nil {
		return err
	}

	if err := syncStates(ctx, tx, doc); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO core_migration_log (
			module_name,
			doctype,
			migration_name,
			migration_type,
			status,
			message,
			json_hash,
			source_path
		)
		VALUES ($1, $2, $3, 'metadata_sync', 'success', 'DocType synced from JSON', $4, $5)
	`, doc.Module, doc.Name, "sync_"+slugify(doc.Name), jsonHash, jsonPath)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func syncFields(ctx context.Context, tx pgx.Tx, doc DocTypeJSON) error {
	_, err := tx.Exec(ctx, `DELETE FROM core_docfield WHERE doctype = $1`, doc.Name)
	if err != nil {
		return err
	}

	for _, field := range doc.Fields {
		if field.Fieldname == "" {
			return fmt.Errorf("fieldname is required in doctype %s", doc.Name)
		}

		if field.Label == "" {
			field.Label = field.Fieldname
		}

		if field.Fieldtype == "" {
			field.Fieldtype = "Data"
		}

		required := field.Required || field.Reqd
		uniqueField := field.Unique || field.UniqueField

		linkFilters, err := normalizeLinkFilters(doc.Name, field.Fieldname, field.LinkFilters)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO core_docfield (
				doctype,
				fieldname,
				label,
				fieldtype,
				options,
				required,
				reqd,
				unique_field,
				read_only,
				hidden,
				in_list_view,
				in_standard_filter,
				in_filter,
				in_global_search,
				in_preview,
				in_import_template,
				search_index,
				sticky,
				idx,
				columns,
				length,
				precision_value,
				non_negative,
				default_value,
				description,
				depends_on,
				mandatory_depends_on,
				read_only_depends_on,
				collapsible,
				collapsible_depends_on,
				hide_border,
				fetch_from,
				fetch_if_empty,
				bold,
				translatable,
				allow_in_quick_entry,
				show_on_timeline,
				print_hide,
				print_hide_if_no_value,
				report_hide,
				print_width,
				width,
				max_height,
				permlevel,
				ignore_user_permissions,
				allow_on_submit,
				allow_bulk_edit,
				no_copy,
				set_only_once,
				remember_last_selected_value,
				ignore_xss_filter,
				alignment,
				placeholder,
				documentation_url,
				oldfieldname,
				oldfieldtype,
				hide_days,
				hide_seconds,
				sort_options,
				link_filters,
				make_attachment_public,
				mask,
				button_color,
				show_description_on_click,
				is_virtual,
				not_nullable,
				validation_rule
			)
			VALUES (
				$1, $2, $3, $4, $5,
				$6, $7, $8, $9, $10,
				$11, $12, $13, $14, $15,
				$16, $17, $18, $19, $20,
				$21, $22, $23, $24, $25,
				$26, $27, $28, $29, $30,
				$31, $32, $33, $34, $35,
				$36, $37, $38, $39, $40,
				$41, $42, $43, $44, $45,
				$46, $47, $48, $49, $50,
				$51, $52, $53, $54, $55,
				$56, $57, $58, $59, $60,
				$61, $62, $63, $64, $65,
				$66, $67
			)
		`,
			doc.Name,
			field.Fieldname,
			field.Label,
			field.Fieldtype,
			field.Options,
			required,
			required,
			uniqueField,
			field.ReadOnly,
			field.Hidden,
			field.InListView,
			field.InStandardFilter,
			field.InFilter,
			field.InGlobalSearch,
			field.InPreview,
			field.InImportTemplate,
			field.SearchIndex,
			field.Sticky,
			field.Idx,
			field.Columns,
			field.Length,
			field.Precision,
			field.NonNegative,
			field.DefaultValue,
			field.Description,
			field.DependsOn,
			field.MandatoryDependsOn,
			field.ReadOnlyDependsOn,
			field.Collapsible,
			field.CollapsibleDependsOn,
			field.HideBorder,
			field.FetchFrom,
			field.FetchIfEmpty,
			field.Bold,
			field.Translatable,
			field.AllowInQuickEntry,
			field.ShowOnTimeline,
			field.PrintHide,
			field.PrintHideIfNoValue,
			field.ReportHide,
			field.PrintWidth,
			field.Width,
			field.MaxHeight,
			field.PermLevel,
			field.IgnoreUserPermissions,
			field.AllowOnSubmit,
			field.AllowBulkEdit,
			field.NoCopy,
			field.SetOnlyOnce,
			field.RememberLastSelectedValue,
			field.IgnoreXSSFilter,
			field.Alignment,
			field.Placeholder,
			field.DocumentationURL,
			field.OldFieldname,
			field.OldFieldtype,
			field.HideDays,
			field.HideSeconds,
			field.SortOptions,
			linkFilters,
			field.MakeAttachmentPublic,
			field.Mask,
			field.ButtonColor,
			field.ShowDescriptionOnClick,
			field.IsVirtual,
			field.NotNullable,
			field.ValidationRule,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func syncActions(ctx context.Context, tx pgx.Tx, doc DocTypeJSON) error {
	_, err := tx.Exec(ctx, `DELETE FROM core_doctype_action WHERE doctype = $1`, doc.Name)
	if err != nil {
		return err
	}

	for _, action := range doc.Actions {
		if action.ActionName == "" {
			return fmt.Errorf("action_name is required in doctype %s", doc.Name)
		}

		if action.Label == "" {
			action.Label = action.ActionName
		}

		if action.ActionType == "" {
			action.ActionType = "server"
		}

		if action.Method == "" {
			action.Method = "POST"
		}

		groupName := action.GroupName
		if groupName == "" {
			groupName = action.Group
		}

		actionValue := action.Action
		if actionValue == "" {
			actionValue = action.Handler
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO core_doctype_action (
				doctype,
				action_name,
				label,
				group_name,
				action_type,
				action,
				handler,
				route,
				method,
				permission,
				visible_when,
				hidden,
				custom,
				enabled,
				idx
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		`,
			doc.Name,
			action.ActionName,
			action.Label,
			groupName,
			action.ActionType,
			actionValue,
			action.Handler,
			action.Route,
			action.Method,
			action.Permission,
			action.VisibleWhen,
			action.Hidden,
			action.Custom,
			action.Enabled,
			action.Idx,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func syncLinks(ctx context.Context, tx pgx.Tx, doc DocTypeJSON) error {
	_, err := tx.Exec(ctx, `DELETE FROM core_doctype_link WHERE doctype = $1`, doc.Name)
	if err != nil {
		return err
	}

	for _, link := range doc.Links {
		if link.LinkDoctype == "" {
			return fmt.Errorf("link_doctype is required in doctype %s", doc.Name)
		}

		if link.LinkFieldname == "" {
			return fmt.Errorf("link_fieldname is required in doctype %s", doc.Name)
		}

		groupName := link.GroupName
		if groupName == "" {
			groupName = link.Group
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO core_doctype_link (
				doctype,
				link_doctype,
				link_fieldname,
				parent_doctype,
				table_fieldname,
				group_name,
				hidden,
				is_child_table,
				custom,
				idx
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		`,
			doc.Name,
			link.LinkDoctype,
			link.LinkFieldname,
			link.ParentDoctype,
			link.TableFieldname,
			groupName,
			link.Hidden,
			link.IsChildTable,
			link.Custom,
			link.Idx,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func syncPermissions(ctx context.Context, tx pgx.Tx, doc DocTypeJSON) error {
	_, err := tx.Exec(ctx, `DELETE FROM core_docperm WHERE doctype = $1`, doc.Name)
	if err != nil {
		return err
	}

	for _, perm := range doc.Permissions {
		if perm.Role == "" {
			return fmt.Errorf("role is required in permissions for doctype %s", doc.Name)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO core_docperm (
				doctype,
				role,
				permlevel,
				if_owner,
				read,
				write,
				create_perm,
				delete_perm,
				submit,
				cancel,
				amend,
				report,
				export,
				import,
				share,
				print,
				email,
				select_perm,
				mask,
				idx
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)
		`,
			doc.Name,
			perm.Role,
			perm.PermLevel,
			perm.IfOwner,
			perm.Read,
			perm.Write,
			perm.Create,
			perm.Delete,
			perm.Submit,
			perm.Cancel,
			perm.Amend,
			perm.Report,
			perm.Export,
			perm.Import,
			perm.Share,
			perm.Print,
			perm.Email,
			perm.Select,
			perm.Mask,
			perm.Idx,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func syncStates(ctx context.Context, tx pgx.Tx, doc DocTypeJSON) error {
	_, err := tx.Exec(ctx, `DELETE FROM core_doctype_state WHERE doctype = $1`, doc.Name)
	if err != nil {
		return err
	}

	for _, state := range doc.States {
		if state.Title == "" {
			return fmt.Errorf("state title is required in doctype %s", doc.Name)
		}

		if state.Color == "" {
			state.Color = "Blue"
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO core_doctype_state (
				doctype,
				title,
				color,
				custom,
				idx
			)
			VALUES ($1,$2,$3,$4,$5)
		`,
			doc.Name,
			state.Title,
			state.Color,
			state.Custom,
			state.Idx,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func findDocTypeJSONFiles(rootPath string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".json" {
			return nil
		}

		cleanPath := filepath.Clean(path)
		parts := strings.Split(cleanPath, string(os.PathSeparator))

		if len(parts) < 5 {
			return nil
		}

		fileName := filepath.Base(cleanPath)
		dirName := filepath.Base(filepath.Dir(cleanPath))
		parentDir := filepath.Base(filepath.Dir(filepath.Dir(cleanPath)))

		if parentDir != "doctype" {
			return nil
		}

		expectedFile := dirName + ".json"
		if fileName != expectedFile {
			return nil
		}

		files = append(files, cleanPath)
		return nil
	})

	return files, err
}

func readDocTypeJSON(filePath string) (DocTypeJSON, []byte, error) {
	var doc DocTypeJSON

	raw, err := os.ReadFile(filePath)
	if err != nil {
		return doc, nil, err
	}

	if err := json.Unmarshal(raw, &doc); err != nil {
		return doc, raw, err
	}

	return doc, raw, nil
}

func hashBytes(raw []byte) string {
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func normalizeLinkFilters(doctype string, fieldname string, value string) (any, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return nil, nil
	}

	var parsed any
	if err := json.Unmarshal([]byte(value), &parsed); err != nil {
		return nil, fmt.Errorf("invalid link_filters JSON for DocType %s, field %s: %w", doctype, fieldname, err)
	}

	return []byte(value), nil
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))

	re := regexp.MustCompile(`[^a-z0-9]+`)
	value = re.ReplaceAllString(value, "_")

	value = strings.Trim(value, "_")

	if value == "" {
		return "unknown"
	}

	return value
}
