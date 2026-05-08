package migration

import "testing"

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

func TestBuildIndexName(t *testing.T) {
	got := buildIndexName("tabNaming Series", "series_key", "uniq")
	want := "tabnaming_series_series_key_uniq"
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestBuildIndexNameTruncates(t *testing.T) {
	long := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" // 61 chars
	got := buildIndexName(long, "col", "idx")
	if len(got) > 60 {
		t.Fatalf("expected max 60 chars, got %d: %s", len(got), got)
	}
}

func TestIsSafeApplyOperation(t *testing.T) {
	tests := []struct {
		name string
		op   Operation
		want bool
	}{
		{
			name: "safe add column",
			op:   Operation{Action: "add_column"},
			want: true,
		},
		{
			name: "dangerous add column blocked",
			op:   Operation{Action: "add_column", Dangerous: true},
			want: false,
		},
		{
			name: "unknown action blocked",
			op:   Operation{Action: "drop_table"},
			want: false,
		},
		{
			name: "create_unique_index is safe",
			op:   Operation{Action: "create_unique_index"},
			want: true,
		},
		{
			name: "create_table is safe",
			op:   Operation{Action: "create_table"},
			want: true,
		},
		{
			name: "set_not_null is safe",
			op:   Operation{Action: "set_not_null"},
			want: true,
		},
		{
			name: "alter_column_type is safe when not dangerous",
			op:   Operation{Action: "alter_column_type"},
			want: true,
		},
		{
			name: "alter_column_type with dangerous flag is blocked",
			op:   Operation{Action: "alter_column_type", Dangerous: true},
			want: false,
		},
	}

	for _, tc := range tests {
		got := isSafeApplyOperation(tc.op)
		if got != tc.want {
			t.Errorf("%s: expected %v, got %v", tc.name, tc.want, got)
		}
	}
}

func TestIsSafeWidening(t *testing.T) {
	tests := []struct {
		actual   string
		expected string
		want     bool
	}{
		{"integer", "bigint", true},
		{"integer", "numeric", true},
		{"varchar", "text", true},
		{"bigint", "numeric", true},
		{"bigint", "double", true},
		{"text", "bigint", false},
		{"numeric", "integer", false},
		{"text", "text", true}, // same type is always safe
	}

	for _, tc := range tests {
		got := isSafeWidening(tc.actual, tc.expected)
		if got != tc.want {
			t.Errorf("%s -> %s: expected %v, got %v", tc.actual, tc.expected, tc.want, got)
		}
	}
}

func TestPostgresType(t *testing.T) {
	tests := []struct {
		fieldtype string
		wantType  string
		wantOK    bool
	}{
		{"Data", "TEXT", true},
		{"Int", "BIGINT", true},
		{"Currency", "NUMERIC(18,6)", true},
		{"Check", "BOOLEAN NOT NULL DEFAULT FALSE", true},
		{"Datetime", "TIMESTAMPTZ", true},
		{"JSON", "JSONB", true},
		{"Float", "DOUBLE PRECISION", true},
		{"Unknown Field", "", false},
		{"", "", false},
	}

	for _, tc := range tests {
		gotType, gotOK := postgresType(tc.fieldtype)
		if gotType != tc.wantType || gotOK != tc.wantOK {
			t.Errorf("%q: expected (%q, %v), got (%q, %v)",
				tc.fieldtype, tc.wantType, tc.wantOK, gotType, gotOK)
		}
	}
}
