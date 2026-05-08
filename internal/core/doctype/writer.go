package doctype

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func WriteDocTypeFile(basePath string, doc DocTypeJSON) (string, string, error) {
	if doc.Name == "" {
		return "", "", fmt.Errorf("doctype name is required")
	}

	if doc.Module == "" {
		return "", "", fmt.Errorf("module is required for doctype %s", doc.Name)
	}

	if doc.Label == "" {
		doc.Label = doc.Name
	}

	if doc.TableName == "" {
		doc.TableName = defaultTableName(doc.Name)
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

	moduleSlug := slugify(doc.Module)
	doctypeSlug := slugify(doc.Name)

	moduleDir := filepath.Join(basePath, "modules", moduleSlug)
	doctypeDir := filepath.Join(moduleDir, "doctype", doctypeSlug)
	filePath := filepath.Join(doctypeDir, doctypeSlug+".json")

	if err := os.MkdirAll(doctypeDir, 0o755); err != nil {
		return "", "", err
	}

	if err := ensureModuleFile(basePath, doc.Module, moduleSlug, doc.IsCore); err != nil {
		return "", "", err
	}

	raw, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", "", err
	}

	raw = append(raw, '\n')

	if err := os.WriteFile(filePath, raw, 0o644); err != nil {
		return "", "", err
	}

	relativePath := filepath.ToSlash(filepath.Join("modules", moduleSlug, "doctype", doctypeSlug, doctypeSlug+".json"))

	return relativePath, hashBytes(raw), nil
}

func ensureModuleFile(basePath string, moduleName string, moduleSlug string, isCore bool) error {
	modulePath := filepath.Join(basePath, "modules", moduleSlug, "module.json")

	if _, err := os.Stat(modulePath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(modulePath), 0o755); err != nil {
		return err
	}

	module := ModuleJSON{
		Name:        moduleName,
		Label:       moduleName,
		Description: fmt.Sprintf("%s module", moduleName),
		Version:     "0.0.1",
		IsCore:      FlexBool(isCore),
		Enabled:     true,
		DoctypePath: filepath.ToSlash(filepath.Join("modules", moduleSlug, "doctype")),
	}

	raw, err := json.MarshalIndent(module, "", "  ")
	if err != nil {
		return err
	}

	raw = append(raw, '\n')

	return os.WriteFile(modulePath, raw, 0o644)
}

func defaultTableName(value string) string {
	normalized := strings.Join(strings.Fields(value), " ")
	if normalized == "" {
		return ""
	}

	return "tab" + normalized
}

// JSONDocType is the canonical struct for DocType JSON files written to
// modules/{module_slug}/doctype/{doctype_slug}/{doctype_slug}.json
type JSONDocType struct {
	Name          string         `json:"name"`
	Label         string         `json:"label"`
	Module        string         `json:"module"`
	AppName       string         `json:"app_name"`
	TableName     string         `json:"table_name"`
	Autoname      string         `json:"autoname"`
	NamingRule    string         `json:"naming_rule"`
	TitleField    string         `json:"title_field"`
	SortField     string         `json:"sort_field"`
	SortOrder     string         `json:"sort_order"`
	DocumentType  string         `json:"document_type"`
	IsSingle      bool           `json:"is_single"`
	IsSubmittable bool           `json:"is_submittable"`
	IsChildTable  bool           `json:"is_child_table"`
	IsTree        bool           `json:"is_tree"`
	AllowImport   bool           `json:"allow_import"`
	AllowExport   bool           `json:"allow_export"`
	AllowRename   bool           `json:"allow_rename"`
	TrackChanges  bool           `json:"track_changes"`
	QuickEntry    bool           `json:"quick_entry"`
	EditableGrid  bool           `json:"editable_grid"`
	Fields        []JSONDocField `json:"fields"`
	Permissions   []JSONDocPerm  `json:"permissions"`
}

type JSONDocField struct {
	Fieldname          string `json:"fieldname"`
	Label              string `json:"label"`
	Fieldtype          string `json:"fieldtype"`
	Options            string `json:"options,omitempty"`
	Reqd               bool   `json:"reqd"`
	Hidden             bool   `json:"hidden"`
	ReadOnly           bool   `json:"read_only"`
	InListView         bool   `json:"in_list_view"`
	InStandardFilter   bool   `json:"in_standard_filter"`
	SearchIndex        bool   `json:"search_index"`
	UniqueField        bool   `json:"unique_field"`
	NoCopy             bool   `json:"no_copy"`
	SetOnlyOnce        bool   `json:"set_only_once"`
	AllowOnSubmit      bool   `json:"allow_on_submit"`
	Permlevel          int    `json:"permlevel"`
	Columns            int    `json:"columns"`
	Length             int    `json:"length"`
	PrecisionValue     int    `json:"precision_value"`
	DefaultValue       string `json:"default_value,omitempty"`
	Description        string `json:"description,omitempty"`
	DependsOn          string `json:"depends_on,omitempty"`
	MandatoryDependsOn string `json:"mandatory_depends_on,omitempty"`
	ReadOnlyDependsOn  string `json:"read_only_depends_on,omitempty"`
	Placeholder        string `json:"placeholder,omitempty"`
	FetchFrom          string `json:"fetch_from,omitempty"`
	ValidationRule     string `json:"validation_rule,omitempty"`
	Idx                int    `json:"idx"`
}

type JSONDocPerm struct {
	Role      string `json:"role"`
	Permlevel int    `json:"permlevel"`
	Read      bool   `json:"read"`
	Write     bool   `json:"write"`
	Create    bool   `json:"create"`
	Delete    bool   `json:"delete"`
	Submit    bool   `json:"submit"`
	Cancel    bool   `json:"cancel"`
	Amend     bool   `json:"amend"`
	Print     bool   `json:"print"`
	Email     bool   `json:"email"`
	Export    bool   `json:"export"`
	Import    bool   `json:"import"`
	Share     bool   `json:"share"`
	Report    bool   `json:"report"`
	Idx       int    `json:"idx"`
}

type WriteResult struct {
	FolderPath string `json:"folder_path"`
	JSONPath   string `json:"json_path"`
}

// WriteDocTypeJSON writes a JSONDocType to
// {rootPath}/modules/{module_slug}/doctype/{doctype_slug}/{doctype_slug}.json
func WriteDocTypeJSON(rootPath string, doc JSONDocType) (WriteResult, error) {
	if strings.TrimSpace(rootPath) == "" {
		rootPath = "."
	}

	doc = normalizeJSONDocType(doc)

	if err := ValidateJSONDocType(doc); err != nil {
		return WriteResult{}, err
	}

	folderPath := filepath.Join(
		rootPath,
		"modules",
		slugify(doc.Module),
		"doctype",
		slugify(doc.Name),
	)

	jsonPath := filepath.Join(
		folderPath,
		slugify(doc.Name)+".json",
	)

	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return WriteResult{}, err
	}

	payload, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return WriteResult{}, err
	}

	payload = append(payload, '\n')

	if err := os.WriteFile(jsonPath, payload, 0644); err != nil {
		return WriteResult{}, err
	}

	return WriteResult{
		FolderPath: folderPath,
		JSONPath:   jsonPath,
	}, nil
}

func normalizeJSONDocType(doc JSONDocType) JSONDocType {
	if doc.Label == "" {
		doc.Label = doc.Name
	}

	if doc.TableName == "" {
		doc.TableName = "tab" + doc.Name
	}

	if doc.Autoname == "" {
		doc.Autoname = "field:name"
	}

	if doc.NamingRule == "" {
		doc.NamingRule = "By fieldname"
	}

	if doc.TitleField == "" {
		doc.TitleField = "name"
	}

	if doc.SortField == "" {
		doc.SortField = "idx"
	}

	if doc.SortOrder == "" {
		doc.SortOrder = "ASC"
	}

	if doc.DocumentType == "" {
		doc.DocumentType = doc.Module
	}

	return doc
}
