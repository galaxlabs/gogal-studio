package app

import (
	"fmt"
	"regexp"
	"strings"
)

var appNamePattern = regexp.MustCompile(`^[a-z][a-z0-9]*(?:_[a-z0-9]+)*$`)

func ValidateAppName(name string) error {
	raw := name
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("app name is required")
	}

	if raw != name {
		return fmt.Errorf("app name %q has leading or trailing spaces", raw)
	}

	if !appNamePattern.MatchString(name) {
		return fmt.Errorf("invalid app name %q: use lowercase letters, numbers, and single underscores like gogal_studio", name)
	}

	return nil
}
