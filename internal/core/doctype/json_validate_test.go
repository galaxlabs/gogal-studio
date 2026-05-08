package doctype

import "testing"

func TestValidateJSONDocTypeValid(t *testing.T) {
	doc := JSONDocType{
		Name:      "Naming Series",
		Module:    "Core",
		AppName:   "gogal_studio",
		TableName: "tabNaming Series",
		Fields: []JSONDocField{
			{
				Fieldname: "series_key",
				Label:     "Series Key",
				Fieldtype: "Data",
				Idx:       1,
			},
			{
				Fieldname: "prefix",
				Label:     "Prefix",
				Fieldtype: "Data",
				Idx:       2,
			},
		},
		Permissions: []JSONDocPerm{
			{
				Role:      "System Manager",
				Permlevel: 0,
				Read:      true,
				Write:     true,
				Create:    true,
				Delete:    true,
				Idx:       1,
			},
		},
	}

	if err := ValidateJSONDocType(doc); err != nil {
		t.Fatal(err)
	}
}

func TestValidateJSONDocTypeRejectsBadDocTypeName(t *testing.T) {
	doc := JSONDocType{
		Name:    "naming series",
		Module:  "Core",
		AppName: "gogal_studio",
	}

	if err := ValidateJSONDocType(doc); err == nil {
		t.Fatal("expected invalid DocType name error")
	}
}

func TestValidateJSONDocTypeRejectsBadModuleName(t *testing.T) {
	doc := JSONDocType{
		Name:    "Naming Series",
		Module:  "core",
		AppName: "gogal_studio",
	}

	if err := ValidateJSONDocType(doc); err == nil {
		t.Fatal("expected invalid module name error")
	}
}

func TestValidateJSONDocTypeRejectsBadAppName(t *testing.T) {
	doc := JSONDocType{
		Name:    "Naming Series",
		Module:  "Core",
		AppName: "Gogal Studio",
	}

	if err := ValidateJSONDocType(doc); err == nil {
		t.Fatal("expected invalid app name error")
	}
}

func TestValidateJSONDocTypeRejectsBadTableName(t *testing.T) {
	doc := JSONDocType{
		Name:      "Naming Series",
		Module:    "Core",
		AppName:   "gogal_studio",
		TableName: "wrong_table",
	}

	if err := ValidateJSONDocType(doc); err == nil {
		t.Fatal("expected invalid table name error")
	}
}

func TestValidateJSONDocTypeRejectsDuplicateFieldname(t *testing.T) {
	doc := JSONDocType{
		Name:    "Naming Series",
		Module:  "Core",
		AppName: "gogal_studio",
		Fields: []JSONDocField{
			{
				Fieldname: "series_key",
				Fieldtype: "Data",
			},
			{
				Fieldname: "series_key",
				Fieldtype: "Data",
			},
		},
	}

	if err := ValidateJSONDocType(doc); err == nil {
		t.Fatal("expected duplicate fieldname error")
	}
}

func TestValidateJSONDocTypeRejectsLinkWithoutOptions(t *testing.T) {
	doc := JSONDocType{
		Name:    "Test DocType",
		Module:  "Core",
		AppName: "gogal_studio",
		Fields: []JSONDocField{
			{
				Fieldname: "user",
				Fieldtype: "Link",
			},
		},
	}

	if err := ValidateJSONDocType(doc); err == nil {
		t.Fatal("expected Link without options error")
	}
}

func TestValidateJSONDocTypeRejectsPermissionWithoutRole(t *testing.T) {
	doc := JSONDocType{
		Name:    "Naming Series",
		Module:  "Core",
		AppName: "gogal_studio",
		Permissions: []JSONDocPerm{
			{
				Role: "",
				Read: true,
			},
		},
	}

	if err := ValidateJSONDocType(doc); err == nil {
		t.Fatal("expected permission role error")
	}
}
