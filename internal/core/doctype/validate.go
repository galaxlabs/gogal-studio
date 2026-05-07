package doctype

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/galaxylabs/gogal-studio/internal/core/naming"
)

var doctypeNamePattern = regexp.MustCompile(`^[A-Z][A-Za-z0-9]*( [A-Z][A-Za-z0-9]*)*$`)

func ValidateDocTypeName(name string) error {
	raw := name
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("doctype name is required")
	}

	if raw != name {
		return fmt.Errorf("doctype name %q has leading or trailing spaces", raw)
	}

	if strings.Contains(name, "  ") {
		return fmt.Errorf("doctype name %q must not contain double spaces", name)
	}

	if !doctypeNamePattern.MatchString(name) {
		return fmt.Errorf("invalid doctype name %q: use Title Case words like Sales Invoice", name)
	}

	return nil
}

func TableName(name string) (string, error) {
	if err := ValidateDocTypeName(name); err != nil {
		return "", err
	}

	return naming.TableNameFromDocType(name), nil
}
