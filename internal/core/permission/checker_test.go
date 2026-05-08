package permission

import "testing"

func TestActionColumnKnownActions(t *testing.T) {
	tests := []struct {
		action     Action
		wantColumn string
	}{
		{ActionRead, `"read"`},
		{ActionWrite, `"write"`},
		{ActionCreate, "create_perm"},
		{ActionDelete, "delete_perm"},
		{ActionSubmit, "submit_perm"},
		{ActionCancel, "cancel_perm"},
		{ActionPrint, "print_perm"},
		{ActionEmail, "email_perm"},
		{ActionExport, "export_perm"},
		{ActionImport, "import_perm"},
		{ActionShare, "share_perm"},
		{ActionReport, "report_perm"},
	}

	for _, tc := range tests {
		got, err := actionColumn(tc.action)
		if err != nil {
			t.Errorf("action %q: unexpected error: %v", tc.action, err)
			continue
		}
		if got != tc.wantColumn {
			t.Errorf("action %q: expected column %q, got %q", tc.action, tc.wantColumn, got)
		}
	}
}

func TestActionColumnUnknownReturnsError(t *testing.T) {
	_, err := actionColumn("drop_table")
	if err == nil {
		t.Fatal("expected error for unknown action, got nil")
	}
}

func TestActionColumnEmptyReturnsError(t *testing.T) {
	_, err := actionColumn("")
	if err == nil {
		t.Fatal("expected error for empty action, got nil")
	}
}
