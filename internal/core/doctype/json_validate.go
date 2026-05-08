package doctype

import (
	"fmt"
	"strings"

	coreapp "github.com/galaxylabs/gogal-studio/internal/core/app"
	"github.com/galaxylabs/gogal-studio/internal/core/fieldtype"
	coremodule "github.com/galaxylabs/gogal-studio/internal/core/module"
	"github.com/galaxylabs/gogal-studio/internal/core/naming"
)

func ValidateJSONDocType(doc JSONDocType) error {
	doc.Name = strings.TrimSpace(doc.Name)
	doc.Module = strings.TrimSpace(doc.Module)
	doc.AppName = strings.TrimSpace(doc.AppName)

	if err := ValidateDocTypeName(doc.Name); err != nil {
		return err
	}

	if err := coremodule.ValidateModuleName(doc.Module); err != nil {
		return fmt.Errorf("invalid module for DocType %q: %w", doc.Name, err)
	}

	if err := coreapp.ValidateAppName(doc.AppName); err != nil {
		return fmt.Errorf("invalid app for DocType %q: %w", doc.Name, err)
	}

	expectedTableName := naming.TableNameFromDocType(doc.Name)

	if strings.TrimSpace(doc.TableName) != "" && doc.TableName != expectedTableName {
		return fmt.Errorf(
			"invalid table_name for DocType %q: expected %q, got %q",
			doc.Name,
			expectedTableName,
			doc.TableName,
		)
	}

	if err := validateJSONDocFields(doc); err != nil {
		return err
	}

	if err := validateJSONDocPerms(doc); err != nil {
		return err
	}

	return nil
}

func validateJSONDocFields(doc JSONDocType) error {
	seen := map[string]bool{}

	for _, field := range doc.Fields {
		field.Fieldname = strings.TrimSpace(field.Fieldname)
		field.Fieldtype = strings.TrimSpace(field.Fieldtype)
		field.Options = strings.TrimSpace(field.Options)

		if seen[field.Fieldname] {
			return fmt.Errorf("duplicate fieldname in DocType %q: %s", doc.Name, field.Fieldname)
		}

		seen[field.Fieldname] = true

		err := fieldtype.ValidateSystemFieldSpec(fieldtype.FieldSpec{
			Fieldname: field.Fieldname,
			Fieldtype: field.Fieldtype,
			Options:   field.Options,
		})
		if err != nil {
			return fmt.Errorf("invalid field %s.%s: %w", doc.Name, field.Fieldname, err)
		}
	}

	return nil
}

func validateJSONDocPerms(doc JSONDocType) error {
	seen := map[string]bool{}

	for _, perm := range doc.Permissions {
		role := strings.TrimSpace(perm.Role)
		if role == "" {
			return fmt.Errorf("permission role is required in DocType %q", doc.Name)
		}

		key := fmt.Sprintf("%s:%d", role, perm.Permlevel)
		if seen[key] {
			return fmt.Errorf("duplicate permission in DocType %q for role %s at permlevel %d", doc.Name, role, perm.Permlevel)
		}

		seen[key] = true
	}

	return nil
}
