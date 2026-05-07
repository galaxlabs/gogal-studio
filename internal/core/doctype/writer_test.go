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
