package fieldtype

import "testing"

func TestFieldTypeRegistryHasCoreTypes(t *testing.T) {
	required := []string{
		"Data",
		"Small Text",
		"Long Text",
		"Text Editor",
		"Int",
		"Float",
		"Currency",
		"Check",
		"Date",
		"Datetime",
		"Select",
		"Link",
		"Table",
		"Attach",
		"Attach Image",
		"JSON",
		"Code",
		"Section",
		"Column",
	}

	for _, name := range required {
		if !IsValid(name) {
			t.Fatalf("expected field type %s to exist", name)
		}
	}
}

func TestLinkRequiresOptions(t *testing.T) {
	err := ValidateFieldSpec(FieldSpec{
		Fieldname: "customer",
		Fieldtype: "Link",
		Options:   "",
	})

	if err == nil {
		t.Fatal("expected Link without options to fail")
	}
}

func TestLinkWithOptionsPasses(t *testing.T) {
	err := ValidateFieldSpec(FieldSpec{
		Fieldname: "customer",
		Fieldtype: "Link",
		Options:   "Customer",
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestTableRequiresChildDocType(t *testing.T) {
	err := ValidateFieldSpec(FieldSpec{
		Fieldname: "items",
		Fieldtype: "Table",
		Options:   "",
	})

	if err == nil {
		t.Fatal("expected Table without options to fail")
	}
}

func TestSQLType(t *testing.T) {
	if SQLType("JSON") != "JSONB" {
		t.Fatalf("expected JSONB, got %s", SQLType("JSON"))
	}

	if SQLType("Data") != "TEXT" {
		t.Fatalf("expected TEXT, got %s", SQLType("Data"))
	}
}

func TestValidFieldnames(t *testing.T) {
	valid := []string{
		"name",
		"customer_name",
		"posting_date",
		"is_active",
		"total_amount_2",
	}

	for _, fieldname := range valid {
		if !IsValidFieldname(fieldname) {
			t.Fatalf("expected valid fieldname: %s", fieldname)
		}
	}
}

func TestInvalidFieldnames(t *testing.T) {
	invalid := []string{
		"",
		"Customer Name",
		"customer name",
		"customer-name",
		"customer.name",
		"1customer",
		"_customer",
		"customer$name",
	}

	for _, fieldname := range invalid {
		if IsValidFieldname(fieldname) {
			t.Fatalf("expected invalid fieldname: %s", fieldname)
		}
	}
}

func TestValidateFieldSpecRejectsBadFieldname(t *testing.T) {
	err := ValidateFieldSpec(FieldSpec{
		Fieldname: "Customer Name",
		Fieldtype: "Data",
	})

	if err == nil {
		t.Fatal("expected invalid fieldname error")
	}
}

func TestReservedFieldnames(t *testing.T) {
	reserved := []string{
		"name",
		"owner",
		"creation",
		"modified",
		"modified_by",
		"docstatus",
		"idx",
	}

	for _, fieldname := range reserved {
		if !IsReservedFieldname(fieldname) {
			t.Fatalf("expected reserved fieldname: %s", fieldname)
		}
	}
}

func TestValidateFieldSpecRejectsReservedFieldname(t *testing.T) {
	err := ValidateFieldSpec(FieldSpec{
		Fieldname: "docstatus",
		Fieldtype: "Int",
	})

	if err == nil {
		t.Fatal("expected reserved fieldname error")
	}
}

func TestValidateSystemFieldSpecAllowsReservedFieldname(t *testing.T) {
	err := ValidateSystemFieldSpec(FieldSpec{
		Fieldname: "docstatus",
		Fieldtype: "Int",
	})

	if err != nil {
		t.Fatal(err)
	}
}
