package migration

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/galaxylabs/gogal-studio/internal/core/meta"
)

type ValidationIssue struct {
	Fieldname string `json:"fieldname,omitempty"`
	Message   string `json:"message"`
}

type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Issues []ValidationIssue `json:"issues"`
}

func ValidateDocTypeMeta(doc meta.DocTypeMeta) ValidationResult {
	result := ValidationResult{
		Valid:  true,
		Issues: []ValidationIssue{},
	}

	if !isValidTableName(doc.TableName) {
		result.add("", fmt.Sprintf("invalid table name: %s", doc.TableName))
	}

	seenFields := map[string]bool{}

	for _, field := range doc.Fields {
		// System fields (hidden=true, read_only=true) are managed internally — skip validation
		if field.Hidden && field.ReadOnly {
			continue
		}

		fieldname := strings.TrimSpace(field.Fieldname)

		if fieldname == "" {
			result.add("", "fieldname is required")
			continue
		}

		if !isValidFieldname(fieldname) {
			result.add(fieldname, "invalid fieldname. Use lowercase letters, numbers, and underscore only.")
		}

		if isReservedFieldname(fieldname) {
			result.add(fieldname, "reserved fieldname cannot be used as custom field")
		}

		if seenFields[fieldname] {
			result.add(fieldname, "duplicate fieldname")
		}

		seenFields[fieldname] = true

		if _, ok := postgresType(field.Fieldtype); !ok && !isLayoutOnlyFieldtype(field.Fieldtype) {
			result.add(fieldname, fmt.Sprintf("unsupported fieldtype: %s", field.Fieldtype))
		}

		if field.Fieldtype == "Link" && strings.TrimSpace(field.Options) == "" {
			result.add(fieldname, "Link field requires options target DocType")
		}

		if field.Fieldtype == "Table" && strings.TrimSpace(field.Options) == "" {
			result.add(fieldname, "Table field requires options child DocType")
		}
	}

	return result
}

func (r *ValidationResult) add(fieldname string, message string) {
	r.Issues = append(r.Issues, ValidationIssue{
		Fieldname: fieldname,
		Message:   message,
	})
	r.Valid = false
}

func isValidTableName(tableName string) bool {
	tableName = strings.TrimSpace(tableName)

	if tableName == "" {
		return false
	}

	if len(tableName) > 63 {
		return false
	}

	re := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_ ]*$`)
	return re.MatchString(tableName)
}

func isValidFieldname(fieldname string) bool {
	fieldname = strings.TrimSpace(fieldname)

	if fieldname == "" {
		return false
	}

	if len(fieldname) > 63 {
		return false
	}

	re := regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
	return re.MatchString(fieldname)
}

func isReservedFieldname(fieldname string) bool {
	reserved := map[string]bool{
		"id":          true,
		"name":        true,
		"owner":       true,
		"creation":    true,
		"modified":    true,
		"modified_by": true,
		"docstatus":   true,
		"idx":         true,
	}

	return reserved[fieldname]
}

func isLayoutOnlyFieldtype(fieldtype string) bool {
	layoutTypes := map[string]bool{
		"Section Break": true,
		"Column Break":  true,
		"Tab Break":     true,
		"Button":        true,
		"HTML":          true,
		"Heading":       true,
		"Fold":          true,
		"Image":         true,
		"Read Only":     true,
	}

	return layoutTypes[fieldtype]
}
