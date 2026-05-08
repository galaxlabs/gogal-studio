package doctype

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/galaxylabs/gogal-studio/internal/core/system"
	"github.com/jackc/pgx/v5"
)

const (
	importOwner      = "Administrator"
	importModifiedBy = "Administrator"
	importDocStatus  = 0
)

func (im *Importer) upsertDocType(ctx context.Context, tx pgx.Tx, doc JSONDocType) error {
	now := time.Now().UTC()

	_, err := tx.Exec(ctx, `
		INSERT INTO "tabDocType" (
			name, creation, modified, modified_by, owner,
			docstatus, idx,
			module, app_name, label, table_name,
			autoname, naming_rule, title_field, sort_field, sort_order, document_type,
			is_single, is_submittable, is_child_table, is_tree,
			allow_import, allow_export, allow_rename,
			track_changes, quick_entry, editable_grid
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7,
			$8, $9, $10, $11,
			$12, $13, $14, $15, $16, $17,
			$18, $19, $20, $21,
			$22, $23, $24,
			$25, $26, $27
		)
		ON CONFLICT (name) DO UPDATE SET
			modified      = NOW(),
			modified_by   = EXCLUDED.modified_by,
			module        = EXCLUDED.module,
			app_name      = EXCLUDED.app_name,
			label         = EXCLUDED.label,
			table_name    = EXCLUDED.table_name,
			autoname      = EXCLUDED.autoname,
			naming_rule   = EXCLUDED.naming_rule,
			title_field   = EXCLUDED.title_field,
			sort_field    = EXCLUDED.sort_field,
			sort_order    = EXCLUDED.sort_order,
			document_type = EXCLUDED.document_type,
			is_single      = EXCLUDED.is_single,
			is_submittable = EXCLUDED.is_submittable,
			is_child_table = EXCLUDED.is_child_table,
			is_tree        = EXCLUDED.is_tree,
			allow_import   = EXCLUDED.allow_import,
			allow_export   = EXCLUDED.allow_export,
			allow_rename   = EXCLUDED.allow_rename,
			track_changes  = EXCLUDED.track_changes,
			quick_entry    = EXCLUDED.quick_entry,
			editable_grid  = EXCLUDED.editable_grid
	`,
		doc.Name, now, now, importModifiedBy, importOwner,
		importDocStatus, 0,
		doc.Module, nullStr(doc.AppName), nullStr(doc.Label), nullStr(doc.TableName),
		nullStr(doc.Autoname), nullStr(doc.NamingRule), nullStr(doc.TitleField),
		nullStr(doc.SortField), nullStr(doc.SortOrder), nullStr(doc.DocumentType),
		doc.IsSingle, doc.IsSubmittable, doc.IsChildTable, doc.IsTree,
		doc.AllowImport, doc.AllowExport, doc.AllowRename,
		doc.TrackChanges, doc.QuickEntry, doc.EditableGrid,
	)
	if err != nil {
		return fmt.Errorf("upsert DocType %s: %w", doc.Name, err)
	}

	return nil
}

func (im *Importer) replaceFields(ctx context.Context, tx pgx.Tx, doc JSONDocType) error {
	if _, err := tx.Exec(ctx, `DELETE FROM "tabDocField" WHERE parent = $1`, doc.Name); err != nil {
		return fmt.Errorf("delete fields for %s: %w", doc.Name, err)
	}

	// Build the full ordered field list: system fields first, then JSON fields.
	systemFields := system.StandardFields()
	offset := len(systemFields)

	allFields := make([]JSONDocField, 0, offset+len(doc.Fields))

	for _, sf := range systemFields {
		allFields = append(allFields, JSONDocField{
			Fieldname:          sf.Fieldname,
			Label:              sf.Label,
			Fieldtype:          sf.Fieldtype,
			Options:            sf.Options,
			Reqd:               sf.Reqd,
			Hidden:             sf.Hidden,
			ReadOnly:           sf.ReadOnly,
			InListView:         sf.InListView,
			InStandardFilter:   sf.InStandardFilter,
			SearchIndex:        sf.SearchIndex,
			UniqueField:        sf.UniqueField,
			NoCopy:             sf.NoCopy,
			SetOnlyOnce:        sf.SetOnlyOnce,
			AllowOnSubmit:      sf.AllowOnSubmit,
			Permlevel:          sf.Permlevel,
			Columns:            sf.Columns,
			Length:             sf.Length,
			PrecisionValue:     sf.PrecisionValue,
			DefaultValue:       sf.DefaultValue,
			Description:        sf.Description,
			DependsOn:          sf.DependsOn,
			MandatoryDependsOn: sf.MandatoryDependsOn,
			ReadOnlyDependsOn:  sf.ReadOnlyDependsOn,
			Placeholder:        sf.Placeholder,
			FetchFrom:          sf.FetchFrom,
			ValidationRule:     sf.ValidationRule,
			Idx:                sf.Idx,
		})
	}

	for i, field := range doc.Fields {
		if field.Idx == 0 {
			field.Idx = i + 1
		}
		field.Idx += offset
		allFields = append(allFields, field)
	}

	now := time.Now().UTC()

	for _, field := range allFields {
		rowName := doc.Name + "." + field.Fieldname

		_, err := tx.Exec(ctx, `
			INSERT INTO "tabDocField" (
				name, creation, modified, modified_by, owner,
				parent, parentfield, parenttype,
				docstatus, idx,
				fieldname, label, fieldtype, options,
				reqd, hidden, read_only,
				in_list_view, in_standard_filter, search_index,
				unique_field, no_copy, set_only_once, allow_on_submit,
				permlevel, columns, length, precision_value,
				default_value, description,
				depends_on, mandatory_depends_on, read_only_depends_on,
				placeholder, fetch_from, validation_rule
			) VALUES (
				$1, $2, $3, $4, $5,
				$6, 'fields', 'DocType',
				$7, $8,
				$9, $10, $11, $12,
				$13, $14, $15,
				$16, $17, $18,
				$19, $20, $21, $22,
				$23, $24, $25, $26,
				$27, $28,
				$29, $30, $31,
				$32, $33, $34
			)
		`,
			rowName, now, now, importModifiedBy, importOwner,
			doc.Name,
			importDocStatus, field.Idx,
			field.Fieldname, nullStr(field.Label), field.Fieldtype, nullStr(field.Options),
			field.Reqd, field.Hidden, field.ReadOnly,
			field.InListView, field.InStandardFilter, field.SearchIndex,
			field.UniqueField, field.NoCopy, field.SetOnlyOnce, field.AllowOnSubmit,
			field.Permlevel, field.Columns, field.Length, field.PrecisionValue,
			nullStr(field.DefaultValue), nullStr(field.Description),
			nullStr(field.DependsOn), nullStr(field.MandatoryDependsOn), nullStr(field.ReadOnlyDependsOn),
			nullStr(field.Placeholder), nullStr(field.FetchFrom), nullStr(field.ValidationRule),
		)
		if err != nil {
			return fmt.Errorf("insert field %s.%s: %w", doc.Name, field.Fieldname, err)
		}
	}

	return nil
}

func (im *Importer) replacePerms(ctx context.Context, tx pgx.Tx, doc JSONDocType) error {
	if _, err := tx.Exec(ctx, `DELETE FROM "tabDocPerm" WHERE parent = $1`, doc.Name); err != nil {
		return fmt.Errorf("delete permissions for %s: %w", doc.Name, err)
	}

	now := time.Now().UTC()

	for i, perm := range doc.Permissions {
		idx := perm.Idx
		if idx == 0 {
			idx = i + 1
		}

		permName := fmt.Sprintf("%s-%s-%d", doc.Name, perm.Role, perm.Permlevel)

		_, err := tx.Exec(ctx, `
			INSERT INTO "tabDocPerm" (
				name, creation, modified, modified_by, owner,
				parent, parentfield, parenttype,
				docstatus, idx,
				role, permlevel,
				"read", "write", create_perm, delete_perm,
				submit_perm, cancel_perm, amend_perm,
				print_perm, email_perm, export_perm,
				import_perm, share_perm, report_perm
			) VALUES (
				$1, $2, $3, $4, $5,
				$6, 'permissions', 'DocType',
				$7, $8,
				$9, $10,
				$11, $12, $13, $14,
				$15, $16, $17,
				$18, $19, $20,
				$21, $22, $23
			)
			ON CONFLICT (name)
			DO UPDATE SET
				parent = EXCLUDED.parent,
				role = EXCLUDED.role,
				permlevel = EXCLUDED.permlevel,
				"read" = EXCLUDED."read",
				"write" = EXCLUDED."write",
				create_perm = EXCLUDED.create_perm,
				delete_perm = EXCLUDED.delete_perm,
				submit_perm = EXCLUDED.submit_perm,
				cancel_perm = EXCLUDED.cancel_perm,
				amend_perm = EXCLUDED.amend_perm,
				print_perm = EXCLUDED.print_perm,
				email_perm = EXCLUDED.email_perm,
				export_perm = EXCLUDED.export_perm,
				import_perm = EXCLUDED.import_perm,
				share_perm = EXCLUDED.share_perm,
				report_perm = EXCLUDED.report_perm,
				owner = EXCLUDED.owner,
				modified_by = EXCLUDED.modified_by,
				docstatus = EXCLUDED.docstatus,
				idx = EXCLUDED.idx,
				modified = NOW()
		`,
			permName, now, now, importModifiedBy, importOwner,
			doc.Name,
			importDocStatus, idx,
			perm.Role, perm.Permlevel,
			perm.Read, perm.Write, perm.Create, perm.Delete,
			perm.Submit, perm.Cancel, perm.Amend,
			perm.Print, perm.Email, perm.Export,
			perm.Import, perm.Share, perm.Report,
		)
		if err != nil {
			return fmt.Errorf("insert permission %s.%s: %w", doc.Name, perm.Role, err)
		}
	}

	return nil
}

// nullStr converts an empty string to nil so the DB stores NULL instead of "".
func nullStr(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}

	return s
}
