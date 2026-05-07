package module

import "testing"

func TestValidateModuleNameValid(t *testing.T) {
	valid := []string{
		"Core",
		"Setup",
		"Security",
		"Desk",
		"Workspace",
		"Navigation",
		"Stock",
		"Accounts",
		"Selling",
		"Buying",
		"HR",
		"CRM",
		"POS",
	}

	for _, name := range valid {
		if err := ValidateModuleName(name); err != nil {
			t.Fatalf("expected valid module name %q: %v", name, err)
		}
	}
}

func TestValidateModuleNameInvalid(t *testing.T) {
	invalid := []string{
		"",
		"core",
		"setup",
		"Setup Module",
		"setup module",
		"1Core",
		"Core-Setup",
		"Core.Setup",
		"Core_Setup",
		" Core",
		"Core ",
	}

	for _, name := range invalid {
		if err := ValidateModuleName(name); err == nil {
			t.Fatalf("expected invalid module name: %q", name)
		}
	}
}
