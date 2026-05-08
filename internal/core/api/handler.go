package api

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/galaxylabs/gogal-studio/internal/core/lifecycle"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	DB *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{DB: db}
}

func decodeName(raw string) string {
	decoded, err := url.PathUnescape(raw)
	if err != nil {
		return raw
	}

	return decoded
}

func (h *Handler) ListDocTypes(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := h.DB.Query(ctx, `
		SELECT
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
			creation,
			modified,
			modified_by,
			docstatus,
			idx
		FROM "tabDocType"
		ORDER BY idx, name
	`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	data := []fiber.Map{}

	for rows.Next() {
		var (
			name          string
			module        string
			appName       string
			tableName     string
			label         string
			autoname      string
			namingRule    string
			titleField    string
			sortField     string
			sortOrder     string
			documentType  string
			isSingle      bool
			isSubmittable bool
			isChildTable  bool
			isTree        bool
			allowImport   bool
			allowExport   bool
			allowRename   bool
			trackChanges  bool
			quickEntry    bool
			editableGrid  bool
			owner         string
			creation      *time.Time
			modified      *time.Time
			modifiedBy    string
			docstatus     int
			idx           int
		)

		if err := rows.Scan(
			&name,
			&module,
			&appName,
			&tableName,
			&label,
			&autoname,
			&namingRule,
			&titleField,
			&sortField,
			&sortOrder,
			&documentType,
			&isSingle,
			&isSubmittable,
			&isChildTable,
			&isTree,
			&allowImport,
			&allowExport,
			&allowRename,
			&trackChanges,
			&quickEntry,
			&editableGrid,
			&owner,
			&creation,
			&modified,
			&modifiedBy,
			&docstatus,
			&idx,
		); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		data = append(data, fiber.Map{
			"name":            name,
			"module":          module,
			"app_name":        appName,
			"table_name":      tableName,
			"label":           label,
			"autoname":        autoname,
			"naming_rule":     namingRule,
			"title_field":     titleField,
			"sort_field":      sortField,
			"sort_order":      sortOrder,
			"document_type":   documentType,
			"is_single":       isSingle,
			"is_submittable":  isSubmittable,
			"is_child_table":  isChildTable,
			"is_tree":         isTree,
			"allow_import":    allowImport,
			"allow_export":    allowExport,
			"allow_rename":    allowRename,
			"track_changes":   trackChanges,
			"quick_entry":     quickEntry,
			"editable_grid":   editableGrid,
			"owner":           owner,
			"creation":        creation,
			"modified":        modified,
			"modified_by":     modifiedBy,
			"docstatus":       docstatus,
			"docstatus_label": lifecycle.DocStatus(docstatus).String(),
			"idx":             idx,
		})
	}

	return c.JSON(fiber.Map{"data": data})
}

func (h *Handler) GetDocType(c *fiber.Ctx) error {
	name := decodeName(c.Params("name"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		module        string
		appName       string
		tableName     string
		label         string
		autoname      string
		namingRule    string
		titleField    string
		sortField     string
		sortOrder     string
		documentType  string
		isSingle      bool
		isSubmittable bool
		isChildTable  bool
		isTree        bool
		allowImport   bool
		allowExport   bool
		allowRename   bool
		trackChanges  bool
		quickEntry    bool
		editableGrid  bool
		owner         string
		creation      *time.Time
		modified      *time.Time
		modifiedBy    string
		docstatus     int
		idx           int
	)

	err := h.DB.QueryRow(ctx, `
		SELECT
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
			creation,
			modified,
			modified_by,
			docstatus,
			idx
		FROM "tabDocType"
		WHERE name = $1
	`, name).Scan(
		&module,
		&appName,
		&tableName,
		&label,
		&autoname,
		&namingRule,
		&titleField,
		&sortField,
		&sortOrder,
		&documentType,
		&isSingle,
		&isSubmittable,
		&isChildTable,
		&isTree,
		&allowImport,
		&allowExport,
		&allowRename,
		&trackChanges,
		&quickEntry,
		&editableGrid,
		&owner,
		&creation,
		&modified,
		&modifiedBy,
		&docstatus,
		&idx,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "DocType not found"})
		}

		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"name":            name,
			"module":          module,
			"app_name":        appName,
			"table_name":      tableName,
			"label":           label,
			"autoname":        autoname,
			"naming_rule":     namingRule,
			"title_field":     titleField,
			"sort_field":      sortField,
			"sort_order":      sortOrder,
			"document_type":   documentType,
			"is_single":       isSingle,
			"is_submittable":  isSubmittable,
			"is_child_table":  isChildTable,
			"is_tree":         isTree,
			"allow_import":    allowImport,
			"allow_export":    allowExport,
			"allow_rename":    allowRename,
			"track_changes":   trackChanges,
			"quick_entry":     quickEntry,
			"editable_grid":   editableGrid,
			"owner":           owner,
			"creation":        creation,
			"modified":        modified,
			"modified_by":     modifiedBy,
			"docstatus":       docstatus,
			"docstatus_label": lifecycle.DocStatus(docstatus).String(),
			"idx":             idx,
		},
	})
}

func (h *Handler) GetDocTypeFields(c *fiber.Ctx) error {
	name := decodeName(c.Params("name"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := h.DB.Query(ctx, `
		SELECT
			name,
			parent,
			fieldname,
			label,
			fieldtype,
			COALESCE(options, '') AS options,
			COALESCE(reqd, false) AS reqd,
			COALESCE(hidden, false) AS hidden,
			COALESCE(read_only, false) AS read_only,
			COALESCE(in_list_view, false) AS in_list_view,
			COALESCE(parentfield, '') AS parentfield,
			COALESCE(parenttype, '') AS parenttype,
			COALESCE(default_value, '') AS default_value,
			COALESCE(description, '') AS description,
			COALESCE(depends_on, '') AS depends_on,
			COALESCE(mandatory_depends_on, '') AS mandatory_depends_on,
			COALESCE(read_only_depends_on, '') AS read_only_depends_on,
			COALESCE(in_standard_filter, false) AS in_standard_filter,
			COALESCE(in_filter, false) AS in_filter,
			COALESCE(in_global_search, false) AS in_global_search,
			COALESCE(search_index, false) AS search_index,
			COALESCE(unique_field, false) AS unique_field,
			COALESCE(no_copy, false) AS no_copy,
			COALESCE(set_only_once, false) AS set_only_once,
			COALESCE(allow_on_submit, false) AS allow_on_submit,
			COALESCE(allow_bulk_edit, false) AS allow_bulk_edit,
			COALESCE(permlevel, 0) AS permlevel,
			COALESCE(columns, 0) AS columns,
			COALESCE(length, 0) AS length,
			COALESCE(precision_value, 0) AS precision_value,
			COALESCE(placeholder, '') AS placeholder,
			COALESCE(fetch_from, '') AS fetch_from,
			COALESCE(validation_rule, '') AS validation_rule,
			owner,
			creation,
			modified,
			modified_by,
			docstatus,
			idx
		FROM "tabDocField"
		WHERE parent = $1
		ORDER BY idx, name
	`, name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	data := []fiber.Map{}

	for rows.Next() {
		var (
			rowName            string
			parent             string
			fieldname          string
			label              string
			fieldtype          string
			options            string
			reqd               bool
			hidden             bool
			readOnly           bool
			inListView         bool
			parentfield        string
			parenttype         string
			defaultValue       string
			description        string
			dependsOn          string
			mandatoryDependsOn string
			readOnlyDependsOn  string
			inStandardFilter   bool
			inFilter           bool
			inGlobalSearch     bool
			searchIndex        bool
			uniqueField        bool
			noCopy             bool
			setOnlyOnce        bool
			allowOnSubmit      bool
			allowBulkEdit      bool
			permlevel          int
			columns            int
			length             int
			precisionValue     int
			placeholder        string
			fetchFrom          string
			validationRule     string
			owner              string
			creation           *time.Time
			modified           *time.Time
			modifiedBy         string
			docstatus          int
			idx                int
		)

		if err := rows.Scan(
			&rowName,
			&parent,
			&fieldname,
			&label,
			&fieldtype,
			&options,
			&reqd,
			&hidden,
			&readOnly,
			&inListView,
			&parentfield,
			&parenttype,
			&defaultValue,
			&description,
			&dependsOn,
			&mandatoryDependsOn,
			&readOnlyDependsOn,
			&inStandardFilter,
			&inFilter,
			&inGlobalSearch,
			&searchIndex,
			&uniqueField,
			&noCopy,
			&setOnlyOnce,
			&allowOnSubmit,
			&allowBulkEdit,
			&permlevel,
			&columns,
			&length,
			&precisionValue,
			&placeholder,
			&fetchFrom,
			&validationRule,
			&owner,
			&creation,
			&modified,
			&modifiedBy,
			&docstatus,
			&idx,
		); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		data = append(data, fiber.Map{
			"name":                 rowName,
			"parent":               parent,
			"fieldname":            fieldname,
			"label":                label,
			"fieldtype":            fieldtype,
			"options":              options,
			"reqd":                 reqd,
			"hidden":               hidden,
			"read_only":            readOnly,
			"in_list_view":         inListView,
			"parentfield":          parentfield,
			"parenttype":           parenttype,
			"default_value":        defaultValue,
			"description":          description,
			"depends_on":           dependsOn,
			"mandatory_depends_on": mandatoryDependsOn,
			"read_only_depends_on": readOnlyDependsOn,
			"in_standard_filter":   inStandardFilter,
			"in_filter":            inFilter,
			"in_global_search":     inGlobalSearch,
			"search_index":         searchIndex,
			"unique_field":         uniqueField,
			"no_copy":              noCopy,
			"set_only_once":        setOnlyOnce,
			"allow_on_submit":      allowOnSubmit,
			"allow_bulk_edit":      allowBulkEdit,
			"permlevel":            permlevel,
			"columns":              columns,
			"length":               length,
			"precision_value":      precisionValue,
			"placeholder":          placeholder,
			"fetch_from":           fetchFrom,
			"validation_rule":      validationRule,
			"owner":                owner,
			"creation":             creation,
			"modified":             modified,
			"modified_by":          modifiedBy,
			"docstatus":            docstatus,
			"docstatus_label":      lifecycle.DocStatus(docstatus).String(),
			"idx":                  idx,
		})
	}

	return c.JSON(fiber.Map{
		"doctype": name,
		"data":    data,
	})
}

func (h *Handler) GetDocTypePermissions(c *fiber.Ctx) error {
	name := decodeName(c.Params("name"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := h.DB.Query(ctx, `
		SELECT
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
			creation,
			modified,
			modified_by,
			docstatus,
			idx
		FROM "tabDocPerm"
		WHERE parent = $1
		ORDER BY idx, name
	`, name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	data := []fiber.Map{}

	for rows.Next() {
		var (
			rowName    string
			parent     string
			role       string
			permlevel  int
			read       bool
			write      bool
			createPerm bool
			deletePerm bool
			submitPerm bool
			cancelPerm bool
			amendPerm  bool
			printPerm  bool
			emailPerm  bool
			exportPerm bool
			importPerm bool
			sharePerm  bool
			reportPerm bool
			owner      string
			creation   *time.Time
			modified   *time.Time
			modifiedBy string
			docstatus  int
			idx        int
		)

		if err := rows.Scan(
			&rowName,
			&parent,
			&role,
			&permlevel,
			&read,
			&write,
			&createPerm,
			&deletePerm,
			&submitPerm,
			&cancelPerm,
			&amendPerm,
			&printPerm,
			&emailPerm,
			&exportPerm,
			&importPerm,
			&sharePerm,
			&reportPerm,
			&owner,
			&creation,
			&modified,
			&modifiedBy,
			&docstatus,
			&idx,
		); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		data = append(data, fiber.Map{
			"name":            rowName,
			"parent":          parent,
			"role":            role,
			"permlevel":       permlevel,
			"read":            read,
			"write":           write,
			"create":          createPerm,
			"delete":          deletePerm,
			"submit":          submitPerm,
			"cancel":          cancelPerm,
			"amend":           amendPerm,
			"print":           printPerm,
			"email":           emailPerm,
			"export":          exportPerm,
			"import":          importPerm,
			"share":           sharePerm,
			"report":          reportPerm,
			"owner":           owner,
			"creation":        creation,
			"modified":        modified,
			"modified_by":     modifiedBy,
			"docstatus":       docstatus,
			"docstatus_label": lifecycle.DocStatus(docstatus).String(),
			"idx":             idx,
		})
	}

	if err := rows.Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"doctype": name,
		"data":    data,
	})
}

type SaveDocTypeRequest struct {
	Name          string              `json:"name"`
	Module        string              `json:"module"`
	AppName       string              `json:"app_name"`
	TableName     string              `json:"table_name"`
	IsSingle      bool                `json:"is_single"`
	IsSubmittable bool                `json:"is_submittable"`
	IsChildTable  bool                `json:"is_child_table"`
	IsTree        bool                `json:"is_tree"`
	Fields        []SaveDocFieldInput `json:"fields"`
	Permissions   []SaveDocPermInput  `json:"permissions"`
}

type SaveDocFieldInput struct {
	Fieldname  string `json:"fieldname"`
	Label      string `json:"label"`
	Fieldtype  string `json:"fieldtype"`
	Options    string `json:"options"`
	Reqd       bool   `json:"reqd"`
	Hidden     bool   `json:"hidden"`
	ReadOnly   bool   `json:"read_only"`
	InListView bool   `json:"in_list_view"`
	Idx        int    `json:"idx"`
}

type SaveDocPermInput struct {
	Role      string `json:"role"`
	Permlevel int    `json:"permlevel"`
	Read      bool   `json:"read"`
	Write     bool   `json:"write"`
	Create    bool   `json:"create"`
	Delete    bool   `json:"delete"`
	Idx       int    `json:"idx"`
}

func (h *Handler) SaveDocType(c *fiber.Ctx) error {
	var req SaveDocTypeRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	normalizeDocTypeRequest(&req)

	if err := validateDocTypeRequest(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tx, err := h.DB.Begin(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO "tabDocType" (
			name,
			module,
			app_name,
			table_name,
			is_single,
			is_submittable,
			is_child_table,
			is_tree,
			idx
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,
			COALESCE((SELECT MAX(idx) + 1 FROM "tabDocType"), 1)
		)
		ON CONFLICT (name)
		DO UPDATE SET
			module = EXCLUDED.module,
			app_name = EXCLUDED.app_name,
			table_name = EXCLUDED.table_name,
			is_single = EXCLUDED.is_single,
			is_submittable = EXCLUDED.is_submittable,
			is_child_table = EXCLUDED.is_child_table,
			is_tree = EXCLUDED.is_tree
	`,
		req.Name,
		req.Module,
		req.AppName,
		req.TableName,
		req.IsSingle,
		req.IsSubmittable,
		req.IsChildTable,
		req.IsTree,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = tx.Exec(ctx, `DELETE FROM "tabDocField" WHERE parent = $1`, req.Name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	for idx, field := range req.Fields {
		if field.Idx == 0 {
			field.Idx = idx + 1
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO "tabDocField" (
				parent,
				fieldname,
				label,
				fieldtype,
				options,
				reqd,
				hidden,
				read_only,
				in_list_view,
				idx
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		`,
			req.Name,
			field.Fieldname,
			field.Label,
			field.Fieldtype,
			field.Options,
			field.Reqd,
			field.Hidden,
			field.ReadOnly,
			field.InListView,
			field.Idx,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}

	_, err = tx.Exec(ctx, `DELETE FROM "tabDocPerm" WHERE parent = $1`, req.Name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	for idx, perm := range req.Permissions {
		if perm.Idx == 0 {
			perm.Idx = idx + 1
		}

		permName := req.Name + "-" + perm.Role + "-" + fmt.Sprint(perm.Permlevel)

		_, err = tx.Exec(ctx, `
			INSERT INTO "tabDocPerm" (
				name,
				parent,
				role,
				permlevel,
				"read",
				"write",
				create_perm,
				delete_perm,
				idx
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		`,
			permName,
			req.Name,
			perm.Role,
			perm.Permlevel,
			perm.Read,
			perm.Write,
			perm.Create,
			perm.Delete,
			perm.Idx,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "DocType metadata saved",
		"data": fiber.Map{
			"name":        req.Name,
			"table_name":  req.TableName,
			"fields":      len(req.Fields),
			"permissions": len(req.Permissions),
		},
	})
}

func normalizeDocTypeRequest(req *SaveDocTypeRequest) {
	req.Name = strings.TrimSpace(req.Name)
	req.Module = strings.TrimSpace(req.Module)
	req.AppName = strings.TrimSpace(req.AppName)
	req.TableName = strings.TrimSpace(req.TableName)

	if req.AppName == "" {
		req.AppName = "gogal_studio"
	}

	if req.Module == "" {
		req.Module = "Core"
	}

	if req.TableName == "" && req.Name != "" {
		req.TableName = "tab" + req.Name
	}

	if req.IsChildTable {
		req.IsSingle = false
		req.IsSubmittable = false
		req.IsTree = false
	}

	if req.IsSingle {
		req.IsChildTable = false
		req.IsSubmittable = false
	}

	if req.IsSubmittable {
		req.IsChildTable = false
		req.IsSingle = false
	}

	for i := range req.Fields {
		req.Fields[i].Label = strings.TrimSpace(req.Fields[i].Label)
		req.Fields[i].Fieldname = strings.TrimSpace(req.Fields[i].Fieldname)

		if req.Fields[i].Fieldname == "" && req.Fields[i].Label != "" {
			req.Fields[i].Fieldname = slugFieldname(req.Fields[i].Label)
		}

		if req.Fields[i].Label == "" && req.Fields[i].Fieldname != "" {
			req.Fields[i].Label = titleFromFieldname(req.Fields[i].Fieldname)
		}

		if req.Fields[i].Fieldtype == "" {
			req.Fields[i].Fieldtype = "Data"
		}

		req.Fields[i].Fieldname = slugFieldname(req.Fields[i].Fieldname)
	}

	if len(req.Permissions) == 0 && !req.IsChildTable {
		req.Permissions = []SaveDocPermInput{
			{
				Role:      "System Manager",
				Permlevel: 0,
				Read:      true,
				Write:     true,
				Create:    true,
				Delete:    true,
				Idx:       1,
			},
		}
	}
}

func validateDocTypeRequest(req SaveDocTypeRequest) error {
	if req.Name == "" {
		return fmt.Errorf("doctype name is required")
	}

	if len(req.Name) > 61 {
		return fmt.Errorf("doctype name cannot be more than 61 characters")
	}

	if req.Module == "" {
		return fmt.Errorf("module is required")
	}

	if req.AppName == "" {
		return fmt.Errorf("app name is required")
	}

	if req.TableName == "" {
		return fmt.Errorf("table name is required")
	}

	seen := map[string]bool{}

	for idx, field := range req.Fields {
		if field.Fieldname == "" {
			return fmt.Errorf("row %d: fieldname is required", idx+1)
		}

		if seen[field.Fieldname] {
			return fmt.Errorf("duplicate fieldname: %s", field.Fieldname)
		}

		seen[field.Fieldname] = true

		if field.Fieldtype == "Link" || field.Fieldtype == "Table" {
			if strings.TrimSpace(field.Options) == "" {
				return fmt.Errorf("row %d: options is required for %s field", idx+1, field.Fieldtype)
			}
		}

		if field.Hidden && field.Reqd {
			return fmt.Errorf("row %d: hidden field cannot be required", idx+1)
		}
	}

	return nil
}

func slugFieldname(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, " ", "_")
	value = strings.ReplaceAll(value, "-", "_")

	re := regexp.MustCompile(`[^a-z0-9_]+`)
	value = re.ReplaceAllString(value, "")

	reUnderscore := regexp.MustCompile(`_+`)
	value = reUnderscore.ReplaceAllString(value, "_")

	value = strings.Trim(value, "_")

	reserved := map[string]bool{
		"name":        true,
		"parent":      true,
		"owner":       true,
		"creation":    true,
		"modified":    true,
		"modified_by": true,
		"docstatus":   true,
		"idx":         true,
	}

	if reserved[value] {
		value = value + "_field"
	}

	return value
}

func titleFromFieldname(value string) string {
	value = strings.ReplaceAll(value, "_", " ")
	parts := strings.Fields(value)

	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}

	return strings.Join(parts, " ")
}
func (h *Handler) ListUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := h.DB.Query(ctx, `
		SELECT name, username, email, idx
		FROM "tabUser"
		ORDER BY idx, name
	`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	data := []fiber.Map{}

	for rows.Next() {
		var (
			name     string
			username string
			email    string
			idx      int
		)

		if err := rows.Scan(&name, &username, &email, &idx); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		data = append(data, fiber.Map{
			"name":     name,
			"username": username,
			"email":    email,
			"idx":      idx,
		})
	}

	if err := rows.Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": data})
}

func (h *Handler) GetUser(c *fiber.Ctx) error {
	name := decodeName(c.Params("name"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		username string
		email    string
		idx      int
	)

	err := h.DB.QueryRow(ctx, `
		SELECT username, email, idx
		FROM "tabUser"
		WHERE name = $1
	`, name).Scan(&username, &email, &idx)

	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "User not found"})
		}

		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"name":     name,
			"username": username,
			"email":    email,
			"idx":      idx,
		},
	})
}

func (h *Handler) GetUserRoles(c *fiber.Ctx) error {
	name := decodeName(c.Params("name"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := h.DB.Query(ctx, `
		SELECT name, role, idx
		FROM "tabHas Role"
		WHERE parent = $1
		ORDER BY idx, name
	`, name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	data := []fiber.Map{}

	for rows.Next() {
		var (
			rowName string
			role    string
			idx     int
		)

		if err := rows.Scan(&rowName, &role, &idx); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		data = append(data, fiber.Map{
			"name": rowName,
			"role": role,
			"idx":  idx,
		})
	}

	if err := rows.Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"user": name,
		"data": data,
	})
}

func (h *Handler) ListRoles(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := h.DB.Query(ctx, `
		SELECT name, role_name, idx
		FROM "tabRole"
		ORDER BY idx, name
	`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	data := []fiber.Map{}

	for rows.Next() {
		var (
			name     string
			roleName string
			idx      int
		)

		if err := rows.Scan(&name, &roleName, &idx); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		data = append(data, fiber.Map{
			"name":      name,
			"role_name": roleName,
			"idx":       idx,
		})
	}

	if err := rows.Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": data})
}
