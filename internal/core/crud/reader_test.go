package crud

import (
	"reflect"
	"testing"

	"github.com/galaxylabs/gogal-studio/internal/core/meta"
)

func TestQuoteIdent(t *testing.T) {
	got := quoteIdent(`tabDocType`)
	want := `"tabDocType"`
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestQuoteIdentEscapesQuotes(t *testing.T) {
	got := quoteIdent(`bad"name`)
	want := `"bad""name"`
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestJoinQuotedColumns(t *testing.T) {
	got := joinQuotedColumns([]string{"name", "module"})
	want := `"name", "module"`
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestJoinQuotedColumnsSingle(t *testing.T) {
	got := joinQuotedColumns([]string{"name"})
	want := `"name"`
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestIsLayoutFieldtype(t *testing.T) {
	tests := []struct {
		fieldtype string
		want      bool
	}{
		{"Section Break", true},
		{"Column Break", true},
		{"Tab Break", true},
		{"HTML", true},
		{"Button", true},
		{"Heading", true},
		{"Fold", true},
		{"Table", true},
		{"Data", false},
		{"Link", false},
		{"Int", false},
		{"", false},
	}

	for _, tc := range tests {
		got := isLayoutFieldtype(tc.fieldtype)
		if got != tc.want {
			t.Errorf("%q: expected %v, got %v", tc.fieldtype, tc.want, got)
		}
	}
}

func TestReadableColumns(t *testing.T) {
	doc := meta.DocTypeMeta{
		Name:      "DocType",
		TableName: "tabDocType",
		Fields: []meta.FieldMeta{
			{Fieldname: "module", Fieldtype: "Link", Hidden: false},
			{Fieldname: "table_name", Fieldtype: "Data", Hidden: false},
			{Fieldname: "hidden_field", Fieldtype: "Data", Hidden: true},
			{Fieldname: "section_1", Fieldtype: "Section Break", Hidden: false},
			// duplicate — should appear only once
			{Fieldname: "module", Fieldtype: "Link", Hidden: false},
		},
	}

	got := readableColumns(doc)

	want := []string{
		"name", "owner", "creation", "modified", "modified_by", "docstatus", "idx",
		"module", "table_name",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %#v\ngot     %#v", want, got)
	}
}

func TestReadableColumnsNoExtraFields(t *testing.T) {
	doc := meta.DocTypeMeta{
		Name:      "Simple",
		TableName: "tabSimple",
		Fields:    []meta.FieldMeta{},
	}

	got := readableColumns(doc)
	// Should still contain the 7 system fields
	if len(got) != 7 {
		t.Fatalf("expected 7 system fields, got %d: %#v", len(got), got)
	}
}

func TestBuildWhereClauseEmpty(t *testing.T) {
	where, args := buildWhereClause(map[string]any{})
	if where != "" {
		t.Fatalf("expected empty where clause, got %q", where)
	}
	if len(args) != 0 {
		t.Fatalf("expected 0 args, got %d", len(args))
	}
}

func TestBuildWhereClauseSingleFilter(t *testing.T) {
	where, args := buildWhereClause(map[string]any{"module": "Core"})
	wantWhere := `WHERE "module" = $1`
	if where != wantWhere {
		t.Fatalf("expected %q, got %q", wantWhere, where)
	}
	if len(args) != 1 || args[0] != "Core" {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildWhereClauseNilFilters(t *testing.T) {
	where, args := buildWhereClause(nil)
	if where != "" {
		t.Fatalf("expected empty where clause for nil filters, got %q", where)
	}
	if len(args) != 0 {
		t.Fatalf("expected 0 args for nil filters, got %d", len(args))
	}
}
