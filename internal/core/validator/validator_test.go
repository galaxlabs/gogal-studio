package validator

import (
	"testing"

	"github.com/galaxylabs/gogal-studio/internal/core/meta"
)

func TestIsStoredFieldtype(t *testing.T) {
	stored := []string{
		"Data", "Int", "Float", "Currency", "Check", "Text", "Long Text",
		"Small Text", "Code", "Date", "Datetime", "Time", "Select", "Link",
		"Dynamic Link", "Attach", "Attach Image", "Password", "Phone",
		"JSON", "Percent", "Rating", "Duration",
	}
	for _, ft := range stored {
		if !IsStoredFieldtype(ft) {
			t.Errorf("expected %q to be a stored fieldtype", ft)
		}
	}

	layout := []string{"Section Break", "Column Break", "Tab Break", "HTML", "Button", "Table", "Heading", "Fold"}
	for _, ft := range layout {
		if IsStoredFieldtype(ft) {
			t.Errorf("expected %q to NOT be a stored fieldtype", ft)
		}
	}
}

func TestValidateRequiredFieldMissing(t *testing.T) {
	m := meta.DocTypeMeta{
		Name: "Item",
		Fields: []meta.FieldMeta{
			{Fieldname: "item_code", Label: "Item Code", Fieldtype: "Data", Reqd: true},
		},
	}
	doc := map[string]any{}
	err := Validate(m, doc)
	if err == nil {
		t.Fatal("expected error for missing required field, got nil")
	}
}

func TestValidateRequiredFieldEmpty(t *testing.T) {
	m := meta.DocTypeMeta{
		Name: "Item",
		Fields: []meta.FieldMeta{
			{Fieldname: "item_code", Label: "Item Code", Fieldtype: "Data", Reqd: true},
		},
	}
	doc := map[string]any{"item_code": ""}
	err := Validate(m, doc)
	if err == nil {
		t.Fatal("expected error for empty required field, got nil")
	}
}

func TestValidateRequiredFieldPresent(t *testing.T) {
	m := meta.DocTypeMeta{
		Name: "Item",
		Fields: []meta.FieldMeta{
			{Fieldname: "item_code", Label: "Item Code", Fieldtype: "Data", Reqd: true},
		},
	}
	doc := map[string]any{"item_code": "ITEM-001"}
	if err := Validate(m, doc); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateOptionalFieldMissing(t *testing.T) {
	m := meta.DocTypeMeta{
		Name: "Item",
		Fields: []meta.FieldMeta{
			{Fieldname: "description", Fieldtype: "Text", Reqd: false},
		},
	}
	doc := map[string]any{}
	if err := Validate(m, doc); err != nil {
		t.Fatalf("optional field missing should not error, got %v", err)
	}
}

func TestValidateSelectValidOption(t *testing.T) {
	m := meta.DocTypeMeta{
		Name: "Item",
		Fields: []meta.FieldMeta{
			{Fieldname: "status", Fieldtype: "Select", Options: "Draft\nSubmitted\nCancelled"},
		},
	}
	doc := map[string]any{"status": "Draft"}
	if err := Validate(m, doc); err != nil {
		t.Fatalf("expected no error for valid select option, got %v", err)
	}
}

func TestValidateSelectInvalidOption(t *testing.T) {
	m := meta.DocTypeMeta{
		Name: "Item",
		Fields: []meta.FieldMeta{
			{Fieldname: "status", Fieldtype: "Select", Options: "Draft\nSubmitted\nCancelled"},
		},
	}
	doc := map[string]any{"status": "Unknown"}
	err := Validate(m, doc)
	if err == nil {
		t.Fatal("expected error for invalid select option, got nil")
	}
}

func TestValidateSkipsLayoutFields(t *testing.T) {
	m := meta.DocTypeMeta{
		Name: "Item",
		Fields: []meta.FieldMeta{
			{Fieldname: "section_1", Fieldtype: "Section Break", Reqd: true},
			{Fieldname: "col_1", Fieldtype: "Column Break", Reqd: true},
		},
	}
	doc := map[string]any{}
	if err := Validate(m, doc); err != nil {
		t.Fatalf("layout fields should be skipped, got %v", err)
	}
}

func TestValidateSkipsHiddenRequiredFields(t *testing.T) {
	m := meta.DocTypeMeta{
		Name: "Item",
		Fields: []meta.FieldMeta{
			{Fieldname: "internal_code", Fieldtype: "Data", Reqd: true, Hidden: true},
		},
	}
	doc := map[string]any{}
	if err := Validate(m, doc); err != nil {
		t.Fatalf("hidden required fields should be skipped, got %v", err)
	}
}

func TestValidateSkipsReadOnlyRequiredFields(t *testing.T) {
	m := meta.DocTypeMeta{
		Name: "Item",
		Fields: []meta.FieldMeta{
			{Fieldname: "computed_field", Fieldtype: "Data", Reqd: true, ReadOnly: true},
		},
	}
	doc := map[string]any{}
	if err := Validate(m, doc); err != nil {
		t.Fatalf("read-only required fields should be skipped, got %v", err)
	}
}
