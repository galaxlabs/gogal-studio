package sysfields

import (
	"testing"
	"time"
)

func TestSystemFieldsContainsAll(t *testing.T) {
	expected := []string{"name", "owner", "creation", "modified", "modified_by", "docstatus", "idx"}
	for _, want := range expected {
		found := false
		for _, got := range SystemFields {
			if got == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("SystemFields missing %q", want)
		}
	}
}

func TestIsSystemField(t *testing.T) {
	tests := []struct {
		fieldname string
		want      bool
	}{
		{"name", true},
		{"owner", true},
		{"creation", true},
		{"modified", true},
		{"modified_by", true},
		{"docstatus", true},
		{"idx", true},
		{"customer_name", false},
		{"module", false},
		{"", false},
	}
	for _, tc := range tests {
		got := IsSystemField(tc.fieldname)
		if got != tc.want {
			t.Errorf("IsSystemField(%q): expected %v, got %v", tc.fieldname, tc.want, got)
		}
	}
}

func TestIsProtectedField(t *testing.T) {
	tests := []struct {
		fieldname string
		want      bool
	}{
		{"name", true},
		{"creation", true},
		{"owner", true},
		{"docstatus", true},
		{"modified", false},    // not protected — can be updated
		{"modified_by", false}, // not protected
		{"idx", false},
		{"module", false},
	}
	for _, tc := range tests {
		got := IsProtectedField(tc.fieldname)
		if got != tc.want {
			t.Errorf("IsProtectedField(%q): expected %v, got %v", tc.fieldname, tc.want, got)
		}
	}
}

func TestInjectCreate(t *testing.T) {
	doc := map[string]any{"module": "Core"}
	InjectCreate(doc, "DOC-0001", "Administrator")

	if doc["name"] != "DOC-0001" {
		t.Errorf("expected name=DOC-0001, got %v", doc["name"])
	}
	if doc["owner"] != "Administrator" {
		t.Errorf("expected owner=Administrator, got %v", doc["owner"])
	}
	if doc["modified_by"] != "Administrator" {
		t.Errorf("expected modified_by=Administrator, got %v", doc["modified_by"])
	}
	if doc["docstatus"] != 0 {
		t.Errorf("expected docstatus=0, got %v", doc["docstatus"])
	}
	if doc["idx"] != 0 {
		t.Errorf("expected idx=0, got %v", doc["idx"])
	}
	if _, ok := doc["creation"]; !ok {
		t.Error("expected creation to be set")
	}
	if _, ok := doc["modified"]; !ok {
		t.Error("expected modified to be set")
	}
}

func TestInjectCreateDoesNotOverwriteDocstatus(t *testing.T) {
	doc := map[string]any{"docstatus": 1}
	InjectCreate(doc, "X", "user")
	if doc["docstatus"] != 1 {
		t.Errorf("InjectCreate should not overwrite existing docstatus, got %v", doc["docstatus"])
	}
}

func TestInjectUpdate(t *testing.T) {
	doc := map[string]any{"name": "DOC-0001", "modified_by": "old"}
	before := time.Now().UTC()
	InjectUpdate(doc, "newuser")

	if doc["modified_by"] != "newuser" {
		t.Errorf("expected modified_by=newuser, got %v", doc["modified_by"])
	}
	ts, ok := doc["modified"].(time.Time)
	if !ok {
		t.Fatalf("expected modified to be time.Time, got %T", doc["modified"])
	}
	if ts.Before(before) {
		t.Error("expected modified to be updated to a recent time")
	}
}
