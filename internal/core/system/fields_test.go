package system

import "testing"

func TestStandardFields(t *testing.T) {
	fields := StandardFields()

	if len(fields) != 7 {
		t.Fatalf("expected 7 standard fields, got %d", len(fields))
	}

	expected := []string{
		"name",
		"owner",
		"creation",
		"modified",
		"modified_by",
		"docstatus",
		"idx",
	}

	for i, fieldname := range expected {
		if fields[i].Fieldname != fieldname {
			t.Fatalf("expected field %s at index %d, got %s", fieldname, i, fields[i].Fieldname)
		}

		if !fields[i].Hidden {
			t.Fatalf("expected %s to be hidden", fieldname)
		}

		if !fields[i].ReadOnly {
			t.Fatalf("expected %s to be read_only", fieldname)
		}

		if fields[i].Idx != i+1 {
			t.Fatalf("expected %s to have idx=%d, got %d", fieldname, i+1, fields[i].Idx)
		}
	}
}

func TestIsSystemField(t *testing.T) {
	systemNames := []string{"name", "owner", "creation", "modified", "modified_by", "docstatus", "idx"}

	for _, n := range systemNames {
		if !IsSystemField(n) {
			t.Fatalf("expected %s to be a system field", n)
		}
	}

	nonSystem := []string{"customer_name", "series_key", "email", "title"}

	for _, n := range nonSystem {
		if IsSystemField(n) {
			t.Fatalf("did not expect %s to be a system field", n)
		}
	}
}

func TestSystemFieldCount(t *testing.T) {
	if SystemFieldCount() != 7 {
		t.Fatalf("expected SystemFieldCount()=7, got %d", SystemFieldCount())
	}
}
