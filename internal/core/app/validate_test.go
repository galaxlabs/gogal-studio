package app

import "testing"

func TestValidateAppNameValid(t *testing.T) {
	valid := []string{
		"gogal_studio",
		"erp",
		"crm",
		"transport",
		"door_app",
		"goods_transport",
		"app1",
		"gogal2_studio",
	}

	for _, name := range valid {
		if err := ValidateAppName(name); err != nil {
			t.Fatalf("expected valid app name %q: %v", name, err)
		}
	}
}

func TestValidateAppNameInvalid(t *testing.T) {
	invalid := []string{
		"",
		"Gogal Studio",
		"gogal studio",
		"gogal-studio",
		"gogal.studio",
		"1gogal",
		"_gogal",
		"gogal_",
		"gogal__studio",
		" gogal",
		"gogal ",
	}

	for _, name := range invalid {
		if err := ValidateAppName(name); err == nil {
			t.Fatalf("expected invalid app name: %q", name)
		}
	}
}
