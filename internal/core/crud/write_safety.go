package crud

import (
	"context"
	"fmt"
	"strings"

	"github.com/galaxylabs/gogal-studio/internal/core/meta"
)

// SafePayloadResult is returned by SafeWritablePayload.
type SafePayloadResult struct {
	Values         map[string]any `json:"values"`
	SkippedFields  []string       `json:"skipped_fields"`
	MissingColumns []string       `json:"missing_columns"`
}

// SafeWritablePayload filters the user-supplied payload down to fields that are:
//   - not system-protected
//   - present in the DocType metadata and writable (not hidden, not read-only, not layout)
//   - present as an actual column in the database table
//
// Fields that exist in metadata but not yet in the DB are tracked in MissingColumns.
// Fields that are protected or not in metadata are tracked in SkippedFields.
func (r *Reader) SafeWritablePayload(ctx context.Context, doc meta.DocTypeMeta, payload map[string]any) (SafePayloadResult, error) {
	actualCols, err := r.actualColumns(ctx, doc.TableName)
	if err != nil {
		return SafePayloadResult{}, err
	}

	writableFields := writableFieldMap(doc)

	result := SafePayloadResult{
		Values:         map[string]any{},
		SkippedFields:  []string{},
		MissingColumns: []string{},
	}

	for key, value := range payload {
		fieldname := strings.TrimSpace(key)

		if fieldname == "" {
			continue
		}

		if isProtectedWriteField(fieldname) {
			result.SkippedFields = append(result.SkippedFields, fieldname)
			continue
		}

		if !writableFields[fieldname] {
			result.SkippedFields = append(result.SkippedFields, fieldname)
			continue
		}

		if !actualCols[fieldname] {
			result.MissingColumns = append(result.MissingColumns, fieldname)
			continue
		}

		result.Values[fieldname] = value
	}

	return result, nil
}

// writableFieldMap returns a set of fieldnames that are allowed in write payloads:
// stored, non-hidden, non-read-only, non-layout fields.
func writableFieldMap(doc meta.DocTypeMeta) map[string]bool {
	fields := map[string]bool{}

	for _, field := range doc.Fields {
		fieldname := strings.TrimSpace(field.Fieldname)

		if fieldname == "" {
			continue
		}

		if field.Hidden || field.ReadOnly {
			continue
		}

		if isLayoutFieldtype(field.Fieldtype) {
			continue
		}

		fields[fieldname] = true
	}

	return fields
}

// isProtectedWriteField returns true for fields that must never be set by user payload.
// System fields (modified, modified_by, idx) are injected automatically by sysfields.
func isProtectedWriteField(fieldname string) bool {
	protected := map[string]bool{
		"id":          true,
		"name":        true,
		"owner":       true,
		"creation":    true,
		"modified":    true,
		"modified_by": true,
		"docstatus":   true,
		"idx":         true,
	}

	return protected[fieldname]
}

// requireWritableValues returns an error if the safe payload has no writable fields.
func requireWritableValues(result SafePayloadResult) error {
	if len(result.Values) == 0 {
		return fmt.Errorf("no writable fields found in request payload")
	}

	return nil
}
