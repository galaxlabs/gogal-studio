package lifecycle

import "testing"

func TestDocStatusString(t *testing.T) {
	tests := []struct {
		value DocStatus
		want  string
	}{
		{Draft, "Draft"},
		{Submitted, "Submitted"},
		{Cancelled, "Cancelled"},
		{DocStatus(99), "Unknown"},
	}

	for _, test := range tests {
		if got := test.value.String(); got != test.want {
			t.Fatalf("expected %s, got %s", test.want, got)
		}
	}
}

func TestIsValid(t *testing.T) {
	if !IsValid(0) || !IsValid(1) || !IsValid(2) {
		t.Fatal("expected 0, 1, 2 to be valid")
	}

	if IsValid(3) {
		t.Fatal("expected 3 to be invalid")
	}
}

func TestCanEdit(t *testing.T) {
	if !CanEdit(0, true) {
		t.Fatal("submittable draft should be editable")
	}

	if CanEdit(1, true) {
		t.Fatal("submittable submitted document should not be editable")
	}

	if CanEdit(2, true) {
		t.Fatal("submittable cancelled document should not be editable")
	}

	if !CanEdit(1, false) {
		t.Fatal("non-submittable submitted-style value should be editable unless cancelled")
	}

	if CanEdit(2, false) {
		t.Fatal("cancelled document should not be editable")
	}
}
