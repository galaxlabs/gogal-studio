package naming

import "testing"

func TestTableNameFromDocType(t *testing.T) {
	got := TableNameFromDocType("Module Def")
	want := "tabModule Def"

	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestGenerateNameFromNameField(t *testing.T) {
	name, err := GenerateName("User", "field:name", Document{
		"name": "Administrator",
	}, nil)

	if err != nil {
		t.Fatal(err)
	}

	if name != "Administrator" {
		t.Fatalf("expected Administrator, got %s", name)
	}
}

func TestGenerateNameFromCustomField(t *testing.T) {
	name, err := GenerateName("User", "field:username", Document{
		"username": "admin",
	}, nil)

	if err != nil {
		t.Fatal(err)
	}

	if name != "admin" {
		t.Fatalf("expected admin, got %s", name)
	}
}

func TestGenerateNameFromSeries(t *testing.T) {
	name, err := GenerateName("Invoice", "series:INV-.#####", Document{}, func(prefix string, digits int) (string, error) {
		if prefix != "INV-." {
			t.Fatalf("expected prefix INV-., got %s", prefix)
		}

		if digits != 5 {
			t.Fatalf("expected 5 digits, got %d", digits)
		}

		return "INV-.00001", nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if name != "INV-.00001" {
		t.Fatalf("expected INV-.00001, got %s", name)
	}
}

func TestGenerateNameManualNeedsProvidedName(t *testing.T) {
	_, err := GenerateName("User", "manual", Document{}, nil)

	if err == nil {
		t.Fatal("expected error")
	}
}
