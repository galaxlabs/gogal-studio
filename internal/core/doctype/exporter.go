package doctype

import (
	"context"
	"fmt"
	"time"

	"github.com/galaxylabs/gogal-studio/internal/core/system"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Exporter struct {
	DB *pgxpool.Pool
}

func NewExporter(db *pgxpool.Pool) *Exporter {
	return &Exporter{DB: db}
}

func (e *Exporter) ExportOne(ctx context.Context, doctypeName string) (JSONDocType, error) {
	doc, err := e.loadDocType(ctx, doctypeName)
	if err != nil {
		return JSONDocType{}, err
	}

	fields, err := e.loadFields(ctx, doctypeName)
	if err != nil {
		return JSONDocType{}, err
	}

	perms, err := e.loadPermissions(ctx, doctypeName)
	if err != nil {
		return JSONDocType{}, err
	}

	doc.Fields = fields
	doc.Permissions = perms

	if err := ValidateJSONDocType(doc); err != nil {
		return JSONDocType{}, err
	}

	return doc, nil
}

func (e *Exporter) ExportOneToFile(rootPath string, doctypeName string) (WriteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	doc, err := e.ExportOne(ctx, doctypeName)
	if err != nil {
		return WriteResult{}, err
	}

	return WriteDocTypeJSON(rootPath, doc)
}

func (e *Exporter) loadDocType(ctx context.Context, doctypeName string) (JSONDocType, error) {
	var doc JSONDocType

	// Nullable string columns from the DB
	var (
		label        *string
		appName      *string
		tableName    *string
		autoname     *string
		namingRule   *string
		titleField   *string
		sortField    *string
		sortOrder    *string
		documentType *string
	)

	err := e.DB.QueryRow(ctx, `
		SELECT
			name,
			label,
			module,
			app_name,
			table_name,
			autoname,
			naming_rule,
			title_field,
			sort_field,
			sort_order,
			document_type,
			is_single,
			is_submittable,
			is_child_table,
			is_tree,
			allow_import,
			allow_export,
			allow_rename,
			track_changes,
			quick_entry,
			editable_grid
		FROM "tabDocType"
		WHERE name = $1
	`, doctypeName).Scan(
		&doc.Name,
		&label,
		&doc.Module,
		&appName,
		&tableName,
		&autoname,
		&namingRule,
		&titleField,
		&sortField,
		&sortOrder,
		&documentType,
		&doc.IsSingle,
		&doc.IsSubmittable,
		&doc.IsChildTable,
		&doc.IsTree,
		&doc.AllowImport,
		&doc.AllowExport,
		&doc.AllowRename,
		&doc.TrackChanges,
		&doc.QuickEntry,
		&doc.EditableGrid,
	)
	if err != nil {
		return JSONDocType{}, fmt.Errorf("load DocType %s: %w", doctypeName, err)
	}

	doc.Label = derefStr(label)
	doc.AppName = derefStr(appName)
	doc.TableName = derefStr(tableName)
	doc.Autoname = derefStr(autoname)
	doc.NamingRule = derefStr(namingRule)
	doc.TitleField = derefStr(titleField)
	doc.SortField = derefStr(sortField)
	doc.SortOrder = derefStr(sortOrder)
	doc.DocumentType = derefStr(documentType)

	return doc, nil
}

func (e *Exporter) loadFields(ctx context.Context, doctypeName string) ([]JSONDocField, error) {
	rows, err := e.DB.Query(ctx, `
		SELECT
			fieldname,
			label,
			fieldtype,
			options,
			reqd,
			hidden,
			read_only,
			in_list_view,
			in_standard_filter,
			search_index,
			unique_field,
			no_copy,
			set_only_once,
			allow_on_submit,
			permlevel,
			columns,
			length,
			precision_value,
			default_value,
			description,
			depends_on,
			mandatory_depends_on,
			read_only_depends_on,
			placeholder,
			fetch_from,
			validation_rule,
			idx
		FROM "tabDocField"
		WHERE parent = $1
		ORDER BY idx, name
	`, doctypeName)
	if err != nil {
		return nil, fmt.Errorf("load fields for %s: %w", doctypeName, err)
	}
	defer rows.Close()

	fields := []JSONDocField{}

	for rows.Next() {
		var field JSONDocField

		var (
			label              *string
			options            *string
			defaultValue       *string
			description        *string
			dependsOn          *string
			mandatoryDependsOn *string
			readOnlyDependsOn  *string
			placeholder        *string
			fetchFrom          *string
			validationRule     *string
		)

		if err := rows.Scan(
			&field.Fieldname,
			&label,
			&field.Fieldtype,
			&options,
			&field.Reqd,
			&field.Hidden,
			&field.ReadOnly,
			&field.InListView,
			&field.InStandardFilter,
			&field.SearchIndex,
			&field.UniqueField,
			&field.NoCopy,
			&field.SetOnlyOnce,
			&field.AllowOnSubmit,
			&field.Permlevel,
			&field.Columns,
			&field.Length,
			&field.PrecisionValue,
			&defaultValue,
			&description,
			&dependsOn,
			&mandatoryDependsOn,
			&readOnlyDependsOn,
			&placeholder,
			&fetchFrom,
			&validationRule,
			&field.Idx,
		); err != nil {
			return nil, err
		}

		field.Label = derefStr(label)
		field.Options = derefStr(options)
		field.DefaultValue = derefStr(defaultValue)
		field.Description = derefStr(description)
		field.DependsOn = derefStr(dependsOn)
		field.MandatoryDependsOn = derefStr(mandatoryDependsOn)
		field.ReadOnlyDependsOn = derefStr(readOnlyDependsOn)
		field.Placeholder = derefStr(placeholder)
		field.FetchFrom = derefStr(fetchFrom)
		field.ValidationRule = derefStr(validationRule)

		if system.IsSystemField(field.Fieldname) {
			continue
		}

		fields = append(fields, field)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fields, nil
}

func (e *Exporter) loadPermissions(ctx context.Context, doctypeName string) ([]JSONDocPerm, error) {
	rows, err := e.DB.Query(ctx, `
		SELECT
			role,
			permlevel,
			"read",
			"write",
			create_perm,
			delete_perm,
			submit_perm,
			cancel_perm,
			amend_perm,
			print_perm,
			email_perm,
			export_perm,
			import_perm,
			share_perm,
			report_perm,
			idx
		FROM "tabDocPerm"
		WHERE parent = $1
		ORDER BY idx, name
	`, doctypeName)
	if err != nil {
		return nil, fmt.Errorf("load permissions for %s: %w", doctypeName, err)
	}
	defer rows.Close()

	perms := []JSONDocPerm{}

	for rows.Next() {
		var perm JSONDocPerm

		if err := rows.Scan(
			&perm.Role,
			&perm.Permlevel,
			&perm.Read,
			&perm.Write,
			&perm.Create,
			&perm.Delete,
			&perm.Submit,
			&perm.Cancel,
			&perm.Amend,
			&perm.Print,
			&perm.Email,
			&perm.Export,
			&perm.Import,
			&perm.Share,
			&perm.Report,
			&perm.Idx,
		); err != nil {
			return nil, err
		}

		perms = append(perms, perm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return perms, nil
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (e *Exporter) ExportAllToFiles(rootPath string) ([]WriteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	rows, err := e.DB.Query(ctx, `
		SELECT name
		FROM "tabDocType"
		ORDER BY module, idx, name
	`)
	if err != nil {
		return nil, fmt.Errorf("load DocType names: %w", err)
	}
	defer rows.Close()

	results := []WriteResult{}

	for rows.Next() {
		var doctypeName string

		if err := rows.Scan(&doctypeName); err != nil {
			return nil, err
		}

		doc, err := e.ExportOne(ctx, doctypeName)
		if err != nil {
			return nil, fmt.Errorf("export DocType %s: %w", doctypeName, err)
		}

		result, err := WriteDocTypeJSON(rootPath, doc)
		if err != nil {
			return nil, fmt.Errorf("write DocType %s: %w", doctypeName, err)
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
