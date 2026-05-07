package module

import (
	"fmt"
	"regexp"
	"strings"
)

var moduleNamePattern = regexp.MustCompile(`^[A-Z][A-Za-z0-9]*$`)

func ValidateModuleName(name string) error {
	raw := name
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("module name is required")
	}

	if raw != name {
		return fmt.Errorf("module name %q has leading or trailing spaces", raw)
	}

	if !moduleNamePattern.MatchString(name) {
		return fmt.Errorf("invalid module name %q: use one Title Case word like Core, Setup, Security, Desk", name)
	}

	return nil
}
