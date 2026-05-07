package doctype

import "testing"

func TestValidateDocTypeNameValid(t *testing.T) {
	valid := []string{
		"User",
		"Role",
		"Module Def",
		"Naming Series",
		"Sales Invoice",
		"Purchase Invoice Item",
		"UOM",
		"OAuth Client",
	}

	for _, name := range valid {
		if err := ValidateDocTypeName(name); err != nil {
			t.Fatalf("expected valid doctype name %q: %v", name, err)
		}
	}
}

func TestValidateDocTypeNameInvalid(t *testing.T) {
	invalid := []string{
		"",
		"sales invoice",
		"Sales invoice",
		"1 Sales Invoice",
		"Sales-Invoice",
		"Sales.Invoice",
		"Sales  Invoice",
		" Sales Invoice",
		"Sales Invoice ",
	}

	for _, name := range invalid {
		if err := ValidateDocTypeName(name); err == nil {
			t.Fatalf("expected invalid doctype name: %q", name)
		}
	}
}

func TestTableName(t *testing.T) {
	got, err := TableName("Sales Invoice")
	if err != nil {
		t.Fatal(err)
	}

	want := "tabSales Invoice"
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}
