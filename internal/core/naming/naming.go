package naming

import (
	"fmt"
	"regexp"
	"strings"
)

type Document map[string]any

func TableNameFromDocType(doctypeName string) string {
	doctypeName = strings.TrimSpace(doctypeName)
	if doctypeName == "" {
		return ""
	}

	return "tab" + doctypeName
}

func GenerateName(doctype string, rule string, doc Document, nextSeries func(prefix string, digits int) (string, error)) (string, error) {
	rule = strings.TrimSpace(rule)

	if rule == "" || rule == "autoname" {
		return "", fmt.Errorf("autoname requires database-backed naming series later")
	}

	if rule == "manual" {
		return "", fmt.Errorf("manual naming requires user-provided name")
	}

	if rule == "field:name" {
		return valueFromField(doc, "name")
	}

	if strings.HasPrefix(rule, "field:") {
		fieldname := strings.TrimSpace(strings.TrimPrefix(rule, "field:"))
		if fieldname == "" {
			return "", fmt.Errorf("field naming rule missing fieldname")
		}

		return valueFromField(doc, fieldname)
	}

	if strings.HasPrefix(rule, "series:") {
		if nextSeries == nil {
			return "", fmt.Errorf("series naming requires nextSeries function")
		}

		pattern := strings.TrimSpace(strings.TrimPrefix(rule, "series:"))
		prefix, digits, err := parseSeriesPattern(pattern)
		if err != nil {
			return "", err
		}

		return nextSeries(prefix, digits)
	}

	return "", fmt.Errorf("unsupported naming rule: %s", rule)
}

func valueFromField(doc Document, fieldname string) (string, error) {
	value, ok := doc[fieldname]
	if !ok || value == nil {
		return "", fmt.Errorf("field %s is required for naming", fieldname)
	}

	name := strings.TrimSpace(fmt.Sprint(value))
	if name == "" {
		return "", fmt.Errorf("field %s is empty", fieldname)
	}

	return name, nil
}

func parseSeriesPattern(pattern string) (string, int, error) {
	if pattern == "" {
		return "", 0, fmt.Errorf("series pattern is required")
	}

	re := regexp.MustCompile(`^(.*?)(#+)$`)
	matches := re.FindStringSubmatch(pattern)
	if len(matches) != 3 {
		return "", 0, fmt.Errorf("series pattern must end with # characters")
	}

	prefix := matches[1]
	digits := len(matches[2])

	if digits <= 0 {
		return "", 0, fmt.Errorf("series digits are required")
	}

	return prefix, digits, nil
}
