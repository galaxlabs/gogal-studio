package doctype

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func WriteDocTypeFile(basePath string, doc DocTypeJSON) (string, string, error) {
	if doc.Name == "" {
		return "", "", fmt.Errorf("doctype name is required")
	}

	if doc.Module == "" {
		return "", "", fmt.Errorf("module is required for doctype %s", doc.Name)
	}

	if doc.Label == "" {
		doc.Label = doc.Name
	}

	if doc.TableName == "" {
		doc.TableName = defaultTableName(doc.Name)
	}

	if doc.NamingRule == "" {
		doc.NamingRule = "autoname"
	}

	if doc.SortField == "" {
		doc.SortField = "created_at"
	}

	if doc.SortOrder == "" {
		doc.SortOrder = "DESC"
	}

	moduleSlug := slugify(doc.Module)
	doctypeSlug := slugify(doc.Name)

	moduleDir := filepath.Join(basePath, "modules", moduleSlug)
	doctypeDir := filepath.Join(moduleDir, "doctype", doctypeSlug)
	filePath := filepath.Join(doctypeDir, doctypeSlug+".json")

	if err := os.MkdirAll(doctypeDir, 0o755); err != nil {
		return "", "", err
	}

	if err := ensureModuleFile(basePath, doc.Module, moduleSlug, doc.IsCore); err != nil {
		return "", "", err
	}

	raw, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", "", err
	}

	raw = append(raw, '\n')

	if err := os.WriteFile(filePath, raw, 0o644); err != nil {
		return "", "", err
	}

	relativePath := filepath.ToSlash(filepath.Join("modules", moduleSlug, "doctype", doctypeSlug, doctypeSlug+".json"))

	return relativePath, hashBytes(raw), nil
}

func ensureModuleFile(basePath string, moduleName string, moduleSlug string, isCore bool) error {
	modulePath := filepath.Join(basePath, "modules", moduleSlug, "module.json")

	if _, err := os.Stat(modulePath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(modulePath), 0o755); err != nil {
		return err
	}

	module := ModuleJSON{
		Name:        moduleName,
		Label:       moduleName,
		Description: fmt.Sprintf("%s module", moduleName),
		Version:     "0.0.1",
		IsCore:      FlexBool(isCore),
		Enabled:     true,
		DoctypePath: filepath.ToSlash(filepath.Join("modules", moduleSlug, "doctype")),
	}

	raw, err := json.MarshalIndent(module, "", "  ")
	if err != nil {
		return err
	}

	raw = append(raw, '\n')

	return os.WriteFile(modulePath, raw, 0o644)
}

func defaultTableName(value string) string {
	normalized := strings.Join(strings.Fields(value), " ")
	if normalized == "" {
		return ""
	}

	return "tab" + normalized
}
