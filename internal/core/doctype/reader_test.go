package doctype

import (
	"path/filepath"
	"testing"
)

func TestReadDocTypeJSON(t *testing.T) {
	root := t.TempDir()

	doc := JSONDocType{
		Name:    "Naming Series",
		Module:  "Core",
		AppName: "gogal_studio",
		Fields: []JSONDocField{
			{
				Fieldname: "series_key",
				Label:     "Series Key",
				Fieldtype: "Data",
				Reqd:      true,
				Idx:       1,
			},
		},
	}

	writeResult, err := WriteDocTypeJSON(root, doc)
	if err != nil {
		t.Fatal(err)
	}

	readResult, err := ReadDocTypeJSON(writeResult.JSONPath)
	if err != nil {
		t.Fatal(err)
	}

	if readResult.DocType.Name != "Naming Series" {
		t.Fatalf("expected Naming Series, got %s", readResult.DocType.Name)
	}

	if readResult.DocType.TableName != "tabNaming Series" {
		t.Fatalf("expected tabNaming Series, got %s", readResult.DocType.TableName)
	}

	if len(readResult.DocType.Fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(readResult.DocType.Fields))
	}
}

func TestReadDocTypeJSONByName(t *testing.T) {
	root := t.TempDir()

	doc := JSONDocType{
		Name:    "Naming Series",
		Module:  "Core",
		AppName: "gogal_studio",
	}

	_, err := WriteDocTypeJSON(root, doc)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ReadDocTypeJSONByName(root, "Core", "Naming Series")
	if err != nil {
		t.Fatal(err)
	}

	if result.DocType.Name != "Naming Series" {
		t.Fatalf("expected Naming Series, got %s", result.DocType.Name)
	}
}

func TestFindDocTypeJSONFiles(t *testing.T) {
	root := t.TempDir()

	_, err := WriteDocTypeJSON(root, JSONDocType{
		Name:    "Naming Series",
		Module:  "Core",
		AppName: "gogal_studio",
	})
	if err != nil {
		t.Fatal(err)
	}

	files, err := FindDocTypeJSONFiles(root)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	want := filepath.Join(root, "modules", "core", "doctype", "naming_series", "naming_series.json")
	if files[0] != want {
		t.Fatalf("expected %s, got %s", want, files[0])
	}
}

func TestReadAllDocTypeJSON(t *testing.T) {
	root := t.TempDir()

	_, err := WriteDocTypeJSON(root, JSONDocType{
		Name:    "Naming Series",
		Module:  "Core",
		AppName: "gogal_studio",
	})
	if err != nil {
		t.Fatal(err)
	}

	results, err := ReadAllDocTypeJSON(root)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].DocType.Name != "Naming Series" {
		t.Fatalf("expected Naming Series, got %s", results[0].DocType.Name)
	}
}
