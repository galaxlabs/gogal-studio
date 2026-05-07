package slug

import "testing"

func TestFromAppName(t *testing.T) {
	got := FromAppName("gogal_studio")
	want := "gogal_studio"

	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestFromModuleName(t *testing.T) {
	got := FromModuleName("Core")
	want := "core"

	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestFromDocTypeName(t *testing.T) {
	got := FromDocTypeName("Naming Series")
	want := "naming_series"

	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestDocTypeFolderPath(t *testing.T) {
	got := DocTypeFolderPath("Core", "Naming Series")
	want := "modules/core/doctype/naming_series"

	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestDocTypeJSONPath(t *testing.T) {
	got := DocTypeJSONPath("Core", "Naming Series")
	want := "modules/core/doctype/naming_series/naming_series.json"

	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}
