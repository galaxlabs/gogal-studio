package doctype

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteDocTypeFileCreatesJSONAndModuleFile(t *testing.T) {
	tempDir := t.TempDir()

	doc := DocTypeJSON{
		Name:      "Customer",
		Module:    "Selling",
		Label:     "Customer",
		TableName: "tab_customer",
		Fields: []DocFieldJSON{
			{
				Fieldname: "customer_name",
				Label:     "Customer Name",
				Fieldtype: "Data",
			},
		},
	}

	relativePath, jsonHash, err := WriteDocTypeFile(tempDir, doc)
	if err != nil {
		t.Fatalf("WriteDocTypeFile returned error: %v", err)
	}

	if relativePath != "modules/selling/doctype/customer/customer.json" {
		t.Fatalf("unexpected relative path: %s", relativePath)
	}

	if jsonHash == "" {
		t.Fatal("expected json hash to be set")
	}

	docPath := filepath.Join(tempDir, filepath.FromSlash(relativePath))
	if _, err := os.Stat(docPath); err != nil {
		t.Fatalf("doctype json was not written: %v", err)
	}

	modulePath := filepath.Join(tempDir, "modules", "selling", "module.json")
	rawModule, err := os.ReadFile(modulePath)
	if err != nil {
		t.Fatalf("module json was not written: %v", err)
	}

	var module ModuleJSON
	if err := json.Unmarshal(rawModule, &module); err != nil {
		t.Fatalf("module json is invalid: %v", err)
	}

	if module.Name != "Selling" {
		t.Fatalf("unexpected module name: %s", module.Name)
	}
}

func TestWriteDocTypeFileDefaultsTableNameToTabPrefixFormat(t *testing.T) {
	tempDir := t.TempDir()

	doc := DocTypeJSON{
		Name:   "Customer Group",
		Module: "Selling",
	}

	relativePath, _, err := WriteDocTypeFile(tempDir, doc)
	if err != nil {
		t.Fatalf("WriteDocTypeFile returned error: %v", err)
	}

	docPath := filepath.Join(tempDir, filepath.FromSlash(relativePath))
	raw, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("doctype json was not written: %v", err)
	}

	var written DocTypeJSON
	if err := json.Unmarshal(raw, &written); err != nil {
		t.Fatalf("doctype json is invalid: %v", err)
	}

	if written.TableName != "tabCustomer Group" {
		t.Fatalf("unexpected table name: %s", written.TableName)
	}
}

func TestWriteDocTypeFileRequiresNameAndModule(t *testing.T) {
	tempDir := t.TempDir()

	if _, _, err := WriteDocTypeFile(tempDir, DocTypeJSON{Module: "Selling"}); err == nil {
		t.Fatal("expected missing doctype name error")
	}

	if _, _, err := WriteDocTypeFile(tempDir, DocTypeJSON{Name: "Customer"}); err == nil {
		t.Fatal("expected missing module error")
	}
}

func TestWriteDocTypeJSON(t *testing.T) {
	root := t.TempDir()

	doc := JSONDocType{
		Name:    "Naming Series",
		Module:  "Core",
		AppName: "gogal_studio",
		Fields: []JSONDocField{
			{
				Fieldname:  "series_key",
				Label:      "Series Key",
				Fieldtype:  "Data",
				Reqd:       true,
				InListView: true,
				Idx:        1,
			},
		},
		Permissions: []JSONDocPerm{
			{
				Role:      "System Manager",
				Permlevel: 0,
				Read:      true,
				Write:     true,
				Create:    true,
				Delete:    true,
				Idx:       1,
			},
		},
	}

	result, err := WriteDocTypeJSON(root, doc)
	if err != nil {
		t.Fatal(err)
	}

	expectedFolder := filepath.Join(root, "modules", "core", "doctype", "naming_series")
	expectedJSON := filepath.Join(expectedFolder, "naming_series.json")

	if result.FolderPath != expectedFolder {
		t.Fatalf("expected folder %s, got %s", expectedFolder, result.FolderPath)
	}

	if result.JSONPath != expectedJSON {
		t.Fatalf("expected json path %s, got %s", expectedJSON, result.JSONPath)
	}

	if _, err := os.Stat(expectedJSON); err != nil {
		t.Fatalf("json file not found: %v", err)
	}
}

func TestWriteDocTypeJSONRejectsInvalidDocTypeName(t *testing.T) {
	root := t.TempDir()

	_, err := WriteDocTypeJSON(root, JSONDocType{
		Name:    "naming series",
		Module:  "Core",
		AppName: "gogal_studio",
	})

	if err == nil {
		t.Fatal("expected invalid DocType name error")
	}
}

func TestWriteDocTypeJSONRequiresModuleAndApp(t *testing.T) {
	root := t.TempDir()

	if _, err := WriteDocTypeJSON(root, JSONDocType{Name: "Customer", AppName: "gogal_studio"}); err == nil {
		t.Fatal("expected missing module error")
	}

	if _, err := WriteDocTypeJSON(root, JSONDocType{Name: "Customer", Module: "Selling"}); err == nil {
		t.Fatal("expected missing app_name error")
	}
}

func TestWriteDocTypeJSONSetsDefaults(t *testing.T) {
	root := t.TempDir()

	result, err := WriteDocTypeJSON(root, JSONDocType{
		Name:    "Sales Invoice",
		Module:  "Accounts",
		AppName: "erp",
	})
	if err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(result.JSONPath)
	if err != nil {
		t.Fatal(err)
	}

	var written JSONDocType
	if err := json.Unmarshal(raw, &written); err != nil {
		t.Fatal(err)
	}

	if written.TableName != "tabSales Invoice" {
		t.Fatalf("unexpected table_name: %s", written.TableName)
	}

	if written.SortOrder != "ASC" {
		t.Fatalf("unexpected sort_order: %s", written.SortOrder)
	}

	if written.NamingRule != "By fieldname" {
		t.Fatalf("unexpected naming_rule: %s", written.NamingRule)
	}
}
