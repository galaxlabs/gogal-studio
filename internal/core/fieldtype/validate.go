package fieldtype

import (
	"fmt"
	"regexp"
	"strings"
)

var fieldnamePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

var reservedFieldnames = map[string]bool{
	"name":        true,
	"owner":       true,
	"creation":    true,
	"modified":    true,
	"modified_by": true,
	"docstatus":   true,
	"idx":         true,
}

type FieldSpec struct {
	Fieldname string
	Fieldtype string
	Options   string
}

func ValidateFieldSpec(field FieldSpec) error {
	return validateFieldSpec(field, false)
}

func ValidateSystemFieldSpec(field FieldSpec) error {
	return validateFieldSpec(field, true)
}

func validateFieldSpec(field FieldSpec, allowReserved bool) error {
	field.Fieldname = strings.TrimSpace(field.Fieldname)
	field.Fieldtype = strings.TrimSpace(field.Fieldtype)
	field.Options = strings.TrimSpace(field.Options)

	if field.Fieldname == "" {
		return fmt.Errorf("fieldname is required")
	}

	if !IsValidFieldname(field.Fieldname) {
		return fmt.Errorf("invalid fieldname %q: use lowercase letters, numbers, and underscore; must start with a letter", field.Fieldname)
	}

	if IsReservedFieldname(field.Fieldname) && !allowReserved {
		return fmt.Errorf("reserved fieldname %q cannot be created manually", field.Fieldname)
	}

	def, ok := Get(field.Fieldtype)
	if !ok {
		return fmt.Errorf("unknown field type: %s", field.Fieldtype)
	}

	if def.RequiresOptions && field.Options == "" {
		return fmt.Errorf("field %s with type %s requires options", field.Fieldname, field.Fieldtype)
	}

	return nil
}

func IsValidFieldname(fieldname string) bool {
	fieldname = strings.TrimSpace(fieldname)

	if fieldname == "" {
		return false
	}

	return fieldnamePattern.MatchString(fieldname)
}

func IsReservedFieldname(fieldname string) bool {
	fieldname = strings.TrimSpace(strings.ToLower(fieldname))
	return reservedFieldnames[fieldname]
}

func ReservedFieldnames() []string {
	names := make([]string, 0, len(reservedFieldnames))

	for name := range reservedFieldnames {
		names = append(names, name)
	}

	return names
}
