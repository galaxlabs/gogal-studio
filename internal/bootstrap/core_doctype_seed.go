package bootstrap

import (
	"context"
	"fmt"
	"time"

	coreapp "github.com/galaxylabs/gogal-studio/internal/core/app"
	coredoctype "github.com/galaxylabs/gogal-studio/internal/core/doctype"
	"github.com/galaxylabs/gogal-studio/internal/core/fieldtype"
	coremodule "github.com/galaxylabs/gogal-studio/internal/core/module"
	"github.com/galaxylabs/gogal-studio/internal/core/system"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	systemOwner      = "Administrator"
	systemModifiedBy = "Administrator"
	systemDocStatus  = 0
)

type coreDocTypeSeed struct {
	Name          string
	Module        string
	AppName       string
	TableName     string
	Label         string
	Autoname      string
	NamingRule    string
	TitleField    string
	SortField     string
	SortOrder     string
	DocumentType  string
	IsSingle      bool
	IsSubmittable bool
	IsChildTable  bool
	IsTree        bool
	AllowImport   bool
	AllowExport   bool
	AllowRename   bool
	TrackChanges  bool
	QuickEntry    bool
	EditableGrid  bool
	Idx           int
	Fields        []coreDocFieldSeed
}

func normalizeCoreDocTypeSeed(dt coreDocTypeSeed) coreDocTypeSeed {
	if dt.Label == "" {
		dt.Label = dt.Name
	}

	if dt.Autoname == "" {
		dt.Autoname = "field:name"
	}

	if dt.NamingRule == "" {
		dt.NamingRule = "By fieldname"
	}

	if dt.TitleField == "" {
		dt.TitleField = "name"
	}

	if dt.SortField == "" {
		dt.SortField = "idx"
	}

	if dt.SortOrder == "" {
		dt.SortOrder = "ASC"
	}

	if dt.DocumentType == "" {
		dt.DocumentType = dt.Module
	}

	dt.AllowImport = true
	dt.AllowExport = true
	dt.AllowRename = true
	dt.TrackChanges = true
	dt.EditableGrid = true

	return dt
}

type coreDocFieldSeed struct {
	Fieldname string
	Label     string
	Fieldtype string
	Options   string

	Reqd       bool
	Hidden     bool
	ReadOnly   bool
	InListView bool

	Parentfield string
	Parenttype  string

	DefaultValue       string
	Description        string
	DependsOn          string
	MandatoryDependsOn string
	ReadOnlyDependsOn  string
	InStandardFilter   bool
	InFilter           bool
	InGlobalSearch     bool
	InPreview          bool
	InImportTemplate   bool
	SearchIndex        bool
	UniqueField        bool
	NoCopy             bool
	SetOnlyOnce        bool
	AllowOnSubmit      bool
	AllowBulkEdit      bool
	IgnoreUserPerms    bool
	IgnoreXSSFilter    bool
	Translatable       bool
	Bold               bool
	PrintHide          bool
	ReportHide         bool
	Permlevel          int
	Columns            int
	Length             int
	PrecisionValue     int
	Width              string
	PrintWidth         string
	MaxHeight          string
	Placeholder        string
	FetchFrom          string
	FetchIfEmpty       bool
	Oldfieldname       string
	Oldfieldtype       string
	ValidationRule     string
	Idx                int
}

func normalizeCoreDocFieldSeed(field coreDocFieldSeed) coreDocFieldSeed {
	if field.Label == "" {
		field.Label = field.Fieldname
	}

	if field.Fieldtype == "" {
		field.Fieldtype = "Data"
	}

	if field.Parentfield == "" {
		field.Parentfield = "fields"
	}

	if field.Parenttype == "" {
		field.Parenttype = "DocType"
	}

	if field.Length == 0 {
		field.Length = fieldtype.DefaultLength(field.Fieldtype)
	}

	if field.Columns == 0 {
		field.Columns = fieldtype.DefaultColumns(field.Fieldtype)
	}

	return field
}

func validateCoreDocFieldSeed(parent string, field coreDocFieldSeed) error {
	validator := fieldtype.ValidateFieldSpec

	if field.Hidden && field.ReadOnly && fieldtype.IsReservedFieldname(field.Fieldname) {
		validator = fieldtype.ValidateSystemFieldSpec
	}

	err := validator(fieldtype.FieldSpec{
		Fieldname: field.Fieldname,
		Fieldtype: field.Fieldtype,
		Options:   field.Options,
	})
	if err != nil {
		return fmt.Errorf("invalid field %s.%s: %w", parent, field.Fieldname, err)
	}

	return nil
}

func buildSeedDocTypeMap(doctypes []coreDocTypeSeed) map[string]bool {
	known := map[string]bool{}

	for _, dt := range doctypes {
		if dt.Name != "" {
			known[dt.Name] = true
		}
	}

	return known
}

func docTypeExistsInDB(ctx context.Context, tx pgx.Tx, doctypeName string) (bool, error) {
	var exists bool

	err := tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM "tabDocType"
			WHERE name = $1
		)
	`, doctypeName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func validateCoreDocFieldLinks(ctx context.Context, tx pgx.Tx, seedDocTypes map[string]bool, parent string, field coreDocFieldSeed) error {
	if field.Options == "" {
		return nil
	}

	switch field.Fieldtype {
	case "Link", "Table":
		target := field.Options

		if seedDocTypes[target] {
			return nil
		}

		exists, err := docTypeExistsInDB(ctx, tx, target)
		if err != nil {
			return fmt.Errorf("check linked doctype for %s.%s: %w", parent, field.Fieldname, err)
		}

		if !exists {
			return fmt.Errorf("invalid field %s.%s: target DocType does not exist: %s", parent, field.Fieldname, target)
		}

	case "Dynamic Link":
		// For Dynamic Link, options points to another field that stores the DocType name.
		return nil
	}

	return nil
}

func standardSystemFields() []coreDocFieldSeed {
	fields := system.StandardFields()
	seeds := make([]coreDocFieldSeed, 0, len(fields))

	for _, f := range fields {
		seeds = append(seeds, coreDocFieldSeed{
			Fieldname:          f.Fieldname,
			Label:              f.Label,
			Fieldtype:          f.Fieldtype,
			Options:            f.Options,
			Reqd:               f.Reqd,
			Hidden:             f.Hidden,
			ReadOnly:           f.ReadOnly,
			InListView:         f.InListView,
			InStandardFilter:   f.InStandardFilter,
			InFilter:           f.InFilter,
			InGlobalSearch:     f.InGlobalSearch,
			SearchIndex:        f.SearchIndex,
			UniqueField:        f.UniqueField,
			NoCopy:             f.NoCopy,
			SetOnlyOnce:        f.SetOnlyOnce,
			AllowOnSubmit:      f.AllowOnSubmit,
			Permlevel:          f.Permlevel,
			Columns:            f.Columns,
			Length:             f.Length,
			PrecisionValue:     f.PrecisionValue,
			DefaultValue:       f.DefaultValue,
			Description:        f.Description,
			DependsOn:          f.DependsOn,
			MandatoryDependsOn: f.MandatoryDependsOn,
			ReadOnlyDependsOn:  f.ReadOnlyDependsOn,
			Placeholder:        f.Placeholder,
			FetchFrom:          f.FetchFrom,
			ValidationRule:     f.ValidationRule,
			Idx:                f.Idx,
			Parentfield:        "fields",
			Parenttype:         "DocType",
		})
	}

	return seeds
}

func withStandardSystemFields(fields []coreDocFieldSeed) []coreDocFieldSeed {
	systemFields := standardSystemFields()

	merged := make([]coreDocFieldSeed, 0, len(systemFields)+len(fields))
	merged = append(merged, systemFields...)

	for _, field := range fields {
		field.Idx = field.Idx + len(systemFields)
		merged = append(merged, field)
	}

	return merged
}

func SeedCoreDocTypeMetadata(database *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tx, err := database.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	doctypes := []coreDocTypeSeed{
		{
			Name:      "Installed App",
			Module:    "Core",
			AppName:   "gogal_studio",
			TableName: "tabInstalled App",
			Idx:       1,
			Fields: []coreDocFieldSeed{
				{Fieldname: "app_name", Label: "App Name", Fieldtype: "Data", Reqd: true, InListView: true, Idx: 1},
				{Fieldname: "app_version", Label: "App Version", Fieldtype: "Data", InListView: true, Idx: 2},
			},
		},
		{
			Name:      "Installed Module",
			Module:    "Core",
			AppName:   "gogal_studio",
			TableName: "tabInstalled Module",
			Idx:       2,
			Fields: []coreDocFieldSeed{
				{Fieldname: "module_name", Label: "Module Name", Fieldtype: "Data", Reqd: true, InListView: true, Idx: 1},
				{Fieldname: "app_name", Label: "App Name", Fieldtype: "Link", Options: "Installed App", Reqd: true, InListView: true, Idx: 2},
			},
		},
		{
			Name:      "Module Def",
			Module:    "Core",
			AppName:   "gogal_studio",
			TableName: "tabModule Def",
			Idx:       3,
			Fields: []coreDocFieldSeed{
				{Fieldname: "module_name", Label: "Module Name", Fieldtype: "Data", Reqd: true, InListView: true, Idx: 1},
				{Fieldname: "app_name", Label: "App Name", Fieldtype: "Link", Options: "Installed App", Reqd: true, InListView: true, Idx: 2},
			},
		},
		{
			Name:      "DocType",
			Module:    "Core",
			AppName:   "gogal_studio",
			TableName: "tabDocType",
			Idx:       4,
			Fields: []coreDocFieldSeed{
				{Fieldname: "module", Label: "Module", Fieldtype: "Link", Options: "Module Def", Reqd: true, InListView: true, Idx: 1},
				{Fieldname: "app_name", Label: "App Name", Fieldtype: "Link", Options: "Installed App", Reqd: true, InListView: true, Idx: 2},
				{Fieldname: "table_name", Label: "Table Name", Fieldtype: "Data", Hidden: true, ReadOnly: true, Idx: 3},
				{Fieldname: "is_single", Label: "Is Single", Fieldtype: "Check", Idx: 4},
				{Fieldname: "is_child_table", Label: "Is Child Table", Fieldtype: "Check", Idx: 5},
				{Fieldname: "is_submittable", Label: "Is Submittable", Fieldtype: "Check", Idx: 6},
			},
		},
		{
			Name:         "DocField",
			Module:       "Core",
			AppName:      "gogal_studio",
			TableName:    "tabDocField",
			IsChildTable: true,
			Idx:          5,
			Fields: []coreDocFieldSeed{
				{Fieldname: "parent", Label: "Parent DocType", Fieldtype: "Link", Options: "DocType", Reqd: true, Hidden: true, ReadOnly: true, InListView: true, Idx: 1},
				{Fieldname: "fieldname", Label: "Field Name", Fieldtype: "Data", Reqd: true, InListView: true, Idx: 2},
				{Fieldname: "label", Label: "Label", Fieldtype: "Data", Reqd: true, InListView: true, Idx: 3},
				{Fieldname: "fieldtype", Label: "Field Type", Fieldtype: "Select", Options: "Data\nText\nSmall Text\nLong Text\nInt\nFloat\nCurrency\nCheck\nDate\nDatetime\nSelect\nLink\nTable\nAttach\nJSON\nCode", Reqd: true, InListView: true, Idx: 4},
				{Fieldname: "options", Label: "Options", Fieldtype: "Small Text", Idx: 5},
				{Fieldname: "reqd", Label: "Required", Fieldtype: "Check", Idx: 6},
				{Fieldname: "hidden", Label: "Hidden", Fieldtype: "Check", Idx: 7},
				{Fieldname: "read_only", Label: "Read Only", Fieldtype: "Check", Idx: 8},
				{Fieldname: "in_list_view", Label: "In List View", Fieldtype: "Check", Idx: 9},
			},
		},
		{
			Name:      "User",
			Module:    "Security",
			AppName:   "gogal_studio",
			TableName: "tabUser",
			Idx:       6,
			Fields: []coreDocFieldSeed{
				{Fieldname: "username", Label: "Username", Fieldtype: "Data", Reqd: true, InListView: true, Idx: 1},
				{Fieldname: "email", Label: "Email", Fieldtype: "Data", Reqd: true, InListView: true, Idx: 2},
				{Fieldname: "password_hash", Label: "Password Hash", Fieldtype: "Password", Reqd: true, Hidden: true, ReadOnly: true, Idx: 3},
			},
		},
		{
			Name:      "Role",
			Module:    "Security",
			AppName:   "gogal_studio",
			TableName: "tabRole",
			Idx:       7,
			Fields: []coreDocFieldSeed{
				{Fieldname: "role_name", Label: "Role Name", Fieldtype: "Data", Reqd: true, InListView: true, Idx: 1},
			},
		},
		{
			Name:         "Has Role",
			Module:       "Security",
			AppName:      "gogal_studio",
			TableName:    "tabHas Role",
			IsChildTable: true,
			Idx:          8,
			Fields: []coreDocFieldSeed{
				{Fieldname: "parent", Label: "User", Fieldtype: "Link", Options: "User", Reqd: true, InListView: true, Idx: 1},
				{Fieldname: "role", Label: "Role", Fieldtype: "Link", Options: "Role", Reqd: true, InListView: true, Idx: 2},
			},
		},
		{
			Name:         "DocPerm",
			Module:       "Security",
			AppName:      "gogal_studio",
			TableName:    "tabDocPerm",
			IsChildTable: true,
			Idx:          9,
			Fields: []coreDocFieldSeed{
				{Fieldname: "parent", Label: "DocType", Fieldtype: "Link", Options: "DocType", Reqd: true, InListView: true, Idx: 1},
				{Fieldname: "role", Label: "Role", Fieldtype: "Link", Options: "Role", Reqd: true, InListView: true, Idx: 2},
				{Fieldname: "read", Label: "Read", Fieldtype: "Check", Idx: 3},
				{Fieldname: "write", Label: "Write", Fieldtype: "Check", Idx: 4},
				{Fieldname: "create", Label: "Create", Fieldtype: "Check", Idx: 5},
				{Fieldname: "delete", Label: "Delete", Fieldtype: "Check", Idx: 6},
			},
		},
		{
			Name:         "Naming Series",
			Module:       "Core",
			AppName:      "gogal_studio",
			TableName:    "tabNaming Series",
			Label:        "Naming Series",
			Autoname:     "field:series_key",
			NamingRule:   "By fieldname",
			TitleField:   "series_key",
			SortField:    "idx",
			SortOrder:    "ASC",
			DocumentType: "Core",
			Idx:          10,
			Fields: []coreDocFieldSeed{
				{
					Fieldname:   "series_key",
					Label:       "Series Key",
					Fieldtype:   "Data",
					Reqd:        true,
					InListView:  true,
					UniqueField: true,
					SearchIndex: true,
					Length:      140,
					Columns:     6,
					Description: "Unique key for this naming series.",
					Idx:         1,
				},
				{
					Fieldname:   "prefix",
					Label:       "Prefix",
					Fieldtype:   "Data",
					InListView:  true,
					Length:      140,
					Columns:     6,
					Description: "Prefix used before numeric counter.",
					Idx:         2,
				},
				{
					Fieldname:   "current_value",
					Label:       "Current Value",
					Fieldtype:   "Int",
					InListView:  true,
					Columns:     4,
					Description: "Current numeric counter value.",
					Idx:         3,
				},
				{
					Fieldname:    "digits",
					Label:        "Digits",
					Fieldtype:    "Int",
					InListView:   true,
					Columns:      4,
					DefaultValue: "5",
					Description:  "Number of padded digits.",
					Idx:          4,
				},
				{
					Fieldname: "description",
					Label:     "Description",
					Fieldtype: "Small Text",
					Columns:   12,
					Idx:       5,
				},
			},
		},
	}

	seedDocTypes := buildSeedDocTypeMap(doctypes)

	for _, dt := range doctypes {
		dt = normalizeCoreDocTypeSeed(dt)

		if err := coredoctype.ValidateDocTypeName(dt.Name); err != nil {
			return fmt.Errorf("invalid seeded DocType %q: %w", dt.Name, err)
		}

		if err := coremodule.ValidateModuleName(dt.Module); err != nil {
			return fmt.Errorf("invalid module for seeded DocType %q: %w", dt.Name, err)
		}

		if err := coreapp.ValidateAppName(dt.AppName); err != nil {
			return fmt.Errorf("invalid app for seeded DocType %q: %w", dt.Name, err)
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO "tabDocType" (
				name,
				module,
				app_name,
				table_name,
				label,
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
				editable_grid,
				owner,
				modified_by,
				docstatus,
				idx
			)
			VALUES (
				$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,
				$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,
				$21,$22,$23,$24,$25
			)
			ON CONFLICT (name)
			DO UPDATE SET
				module = EXCLUDED.module,
				app_name = EXCLUDED.app_name,
				table_name = EXCLUDED.table_name,
				label = EXCLUDED.label,
				autoname = EXCLUDED.autoname,
				naming_rule = EXCLUDED.naming_rule,
				title_field = EXCLUDED.title_field,
				sort_field = EXCLUDED.sort_field,
				sort_order = EXCLUDED.sort_order,
				document_type = EXCLUDED.document_type,
				is_single = EXCLUDED.is_single,
				is_submittable = EXCLUDED.is_submittable,
				is_child_table = EXCLUDED.is_child_table,
				is_tree = EXCLUDED.is_tree,
				allow_import = EXCLUDED.allow_import,
				allow_export = EXCLUDED.allow_export,
				allow_rename = EXCLUDED.allow_rename,
				track_changes = EXCLUDED.track_changes,
				quick_entry = EXCLUDED.quick_entry,
				editable_grid = EXCLUDED.editable_grid,
				owner = EXCLUDED.owner,
				modified_by = EXCLUDED.modified_by,
				docstatus = EXCLUDED.docstatus,
				idx = EXCLUDED.idx,
				modified = NOW()
		`,
			dt.Name,
			dt.Module,
			dt.AppName,
			DocTypeTableName(dt.Name),
			dt.Label,
			dt.Autoname,
			dt.NamingRule,
			dt.TitleField,
			dt.SortField,
			dt.SortOrder,
			dt.DocumentType,
			dt.IsSingle,
			dt.IsSubmittable,
			dt.IsChildTable,
			dt.IsTree,
			dt.AllowImport,
			dt.AllowExport,
			dt.AllowRename,
			dt.TrackChanges,
			dt.QuickEntry,
			dt.EditableGrid,
			systemOwner,
			systemModifiedBy,
			systemDocStatus,
			dt.Idx,
		)
		if err != nil {
			return fmt.Errorf("seed doctype %s: %w", dt.Name, err)
		}

		_, err = tx.Exec(ctx, `DELETE FROM "tabDocField" WHERE parent = $1`, dt.Name)
		if err != nil {
			return fmt.Errorf("clear fields for %s: %w", dt.Name, err)
		}

		for _, field := range withStandardSystemFields(dt.Fields) {
			field = normalizeCoreDocFieldSeed(field)

			if err := validateCoreDocFieldSeed(dt.Name, field); err != nil {
				return err
			}

			if err := validateCoreDocFieldLinks(ctx, tx, seedDocTypes, dt.Name, field); err != nil {
				return err
			}

			_, err := tx.Exec(ctx, `
				INSERT INTO "tabDocField" (
					name,
					parent,
					fieldname,
					label,
					fieldtype,
					options,
					reqd,
					hidden,
					read_only,
					in_list_view,
					parentfield,
					parenttype,
					default_value,
					description,
					depends_on,
					mandatory_depends_on,
					read_only_depends_on,
					in_standard_filter,
					in_filter,
					in_global_search,
					in_preview,
					in_import_template,
					search_index,
					unique_field,
					no_copy,
					set_only_once,
					allow_on_submit,
					allow_bulk_edit,
					ignore_user_permissions,
					ignore_xss_filter,
					translatable,
					bold,
					print_hide,
					report_hide,
					permlevel,
					columns,
					length,
					precision_value,
					width,
					print_width,
					max_height,
					placeholder,
					fetch_from,
					fetch_if_empty,
					oldfieldname,
					oldfieldtype,
					validation_rule,
					owner,
					modified_by,
					docstatus,
					idx
				)
				VALUES (
					$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,
					$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,
					$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,
					$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,
					$41,$42,$43,$44,$45,$46,$47,$48,$49,$50,
					$51
				)
			`,
				dt.Name+"."+field.Fieldname,
				dt.Name,
				field.Fieldname,
				field.Label,
				field.Fieldtype,
				field.Options,
				field.Reqd,
				field.Hidden,
				field.ReadOnly,
				field.InListView,
				field.Parentfield,
				field.Parenttype,
				field.DefaultValue,
				field.Description,
				field.DependsOn,
				field.MandatoryDependsOn,
				field.ReadOnlyDependsOn,
				field.InStandardFilter,
				field.InFilter,
				field.InGlobalSearch,
				field.InPreview,
				field.InImportTemplate,
				field.SearchIndex,
				field.UniqueField,
				field.NoCopy,
				field.SetOnlyOnce,
				field.AllowOnSubmit,
				field.AllowBulkEdit,
				field.IgnoreUserPerms,
				field.IgnoreXSSFilter,
				field.Translatable,
				field.Bold,
				field.PrintHide,
				field.ReportHide,
				field.Permlevel,
				field.Columns,
				field.Length,
				field.PrecisionValue,
				field.Width,
				field.PrintWidth,
				field.MaxHeight,
				field.Placeholder,
				field.FetchFrom,
				field.FetchIfEmpty,
				field.Oldfieldname,
				field.Oldfieldtype,
				field.ValidationRule,
				systemOwner,
				systemModifiedBy,
				systemDocStatus,
				field.Idx,
			)
			if err != nil {
				return fmt.Errorf("seed field %s.%s: %w", dt.Name, field.Fieldname, err)
			}
		}
	}

	return tx.Commit(ctx)
}
func SeedCoreDocPerms(database *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tx, err := database.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	doctypes := []string{
		"Installed App",
		"Installed Module",
		"Module Def",
		"DocType",
		"DocField",
		"User",
		"Role",
		"Has Role",
		"DocPerm",
		"Naming Series",
	}

	for idx, doctypeName := range doctypes {
		permName := fmt.Sprintf("%s-System Manager-0", doctypeName)

		_, err := tx.Exec(ctx, `
			INSERT INTO "tabDocPerm" (
				name,
				parent,
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
				owner,
				modified_by,
				docstatus,
				idx
			)
			VALUES (
				$1,$2,$3,0,
				TRUE,TRUE,TRUE,TRUE,
				TRUE,TRUE,TRUE,
				TRUE,TRUE,TRUE,TRUE,TRUE,TRUE,
				$4,$5,$6,$7
			)
			ON CONFLICT ON CONSTRAINT "tabDocPerm_parent_role_permlevel_key"
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
			permName,
			doctypeName,
			"System Manager",
			systemOwner,
			systemModifiedBy,
			systemDocStatus,
			idx+1,
		)
		if err != nil {
			return fmt.Errorf("seed permission for %s: %w", doctypeName, err)
		}
	}

	return tx.Commit(ctx)
}

func SeedDefaultNamingSeries(database *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tx, err := database.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	series := []struct {
		Name         string
		SeriesKey    string
		Prefix       string
		CurrentValue int
		Digits       int
		Description  string
		Idx          int
	}{
		{
			Name:         "CORE-.#####",
			SeriesKey:    "CORE-.#####",
			Prefix:       "CORE-.",
			CurrentValue: 0,
			Digits:       5,
			Description:  "Default Core document series.",
			Idx:          1,
		},
	}

	for _, row := range series {
		_, err := tx.Exec(ctx, `
			INSERT INTO "tabNaming Series" (
				name,
				series_key,
				prefix,
				current_value,
				digits,
				description,
				owner,
				modified_by,
				docstatus,
				idx
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
			ON CONFLICT (name)
			DO UPDATE SET
				series_key = EXCLUDED.series_key,
				prefix = EXCLUDED.prefix,
				current_value = EXCLUDED.current_value,
				digits = EXCLUDED.digits,
				description = EXCLUDED.description,
				owner = EXCLUDED.owner,
				modified_by = EXCLUDED.modified_by,
				docstatus = EXCLUDED.docstatus,
				idx = EXCLUDED.idx,
				modified = NOW()
		`,
			row.Name,
			row.SeriesKey,
			row.Prefix,
			row.CurrentValue,
			row.Digits,
			row.Description,
			systemOwner,
			systemModifiedBy,
			systemDocStatus,
			row.Idx,
		)
		if err != nil {
			return fmt.Errorf("seed naming series %s: %w", row.Name, err)
		}
	}

	return tx.Commit(ctx)
}
