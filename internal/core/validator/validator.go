package validator

import (
	"fmt"
	"strings"

	"github.com/galaxylabs/gogal-studio/internal/core/meta"
)

// storedFieldtypes is the set of fieldtypes that map to actual DB columns.
// Layout-only types (Section Break, Column Break, HTML, Table, etc.) are excluded.
var storedFieldtypes = map[string]bool{
	"Autocomplete":    true,
	"Attach":          true,
	"Attach Image":    true,
	"Check":           true,
	"Code":            true,
	"Color":           true,
	"Currency":        true,
	"Data":            true,
	"Date":            true,
	"Datetime":        true,
	"Duration":        true,
	"Dynamic Link":    true,
	"Float":           true,
	"Geolocation":     true,
	"HTML Editor":     true,
	"Int":             true,
	"JSON":            true,
	"Link":            true,
	"Long Text":       true,
	"Markdown Editor": true,
	"Name":            true,
	"Password":        true,
	"Percent":         true,
	"Phone":           true,
	"Rating":          true,
	"Read Only":       true,
	"Select":          true,
	"Signature":       true,
	"Small Text":      true,
	"Text":            true,
	"Time":            true,
}

// IsStoredFieldtype returns true when this fieldtype is stored as a DB column.
func IsStoredFieldtype(ft string) bool {
	return storedFieldtypes[ft]
}

// Validate checks required fields and Select option constraints against meta.
// doc should be the user-supplied document map (without system fields yet).
func Validate(m meta.DocTypeMeta, doc map[string]any) error {
	for _, field := range m.Fields {
		if !IsStoredFieldtype(field.Fieldtype) {
			continue
		}
		// Skip layout-only or output-only fields
		if field.Hidden || field.ReadOnly {
			continue
		}

		val, ok := doc[field.Fieldname]

		if field.Reqd && (!ok || isEmpty(val)) {
			label := field.Label
			if label == "" {
				label = field.Fieldname
			}
			return fmt.Errorf("field %q is required", label)
		}

		if ok && !isEmpty(val) && field.Fieldtype == "Select" && field.Options != "" {
			if err := validateSelect(field, val); err != nil {
				return err
			}
		}
	}
	return nil
}

func isEmpty(v any) bool {
	if v == nil {
		return true
	}
	s, ok := v.(string)
	return ok && strings.TrimSpace(s) == ""
}

func validateSelect(field meta.FieldMeta, val any) error {
	raw, ok := val.(string)
	if !ok {
		label := field.Label
		if label == "" {
			label = field.Fieldname
		}
		return fmt.Errorf("field %q must be a string", label)
	}
	for _, opt := range strings.Split(field.Options, "\n") {
		if strings.TrimSpace(opt) == raw {
			return nil
		}
	}
	label := field.Label
	if label == "" {
		label = field.Fieldname
	}
	return fmt.Errorf("field %q value %q is not a valid option", label, raw)
}
