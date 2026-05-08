package doctype

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/galaxylabs/gogal-studio/internal/core/slug"
)

type ReadResult struct {
	JSONPath string      `json:"json_path"`
	DocType  JSONDocType `json:"doctype"`
}

func ReadDocTypeJSON(jsonPath string) (ReadResult, error) {
	if jsonPath == "" {
		return ReadResult{}, fmt.Errorf("json path is required")
	}

	payload, err := os.ReadFile(jsonPath)
	if err != nil {
		return ReadResult{}, err
	}

	var doc JSONDocType
	if err := json.Unmarshal(payload, &doc); err != nil {
		return ReadResult{}, fmt.Errorf("invalid doctype json %s: %w", jsonPath, err)
	}

	doc = normalizeJSONDocType(doc)

	if err := ValidateJSONDocType(doc); err != nil {
		return ReadResult{}, fmt.Errorf("invalid DocType JSON %s: %w", jsonPath, err)
	}

	return ReadResult{
		JSONPath: jsonPath,
		DocType:  doc,
	}, nil
}

func ReadDocTypeJSONByName(rootPath string, moduleName string, doctypeName string) (ReadResult, error) {
	if rootPath == "" {
		rootPath = "."
	}

	jsonPath := filepath.Join(
		rootPath,
		slug.DocTypeJSONPath(moduleName, doctypeName),
	)

	return ReadDocTypeJSON(jsonPath)
}

func FindDocTypeJSONFiles(rootPath string) ([]string, error) {
	if rootPath == "" {
		rootPath = "."
	}

	basePath := filepath.Join(rootPath, "modules")

	files := []string{}

	err := filepath.WalkDir(basePath, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".json" {
			return nil
		}

		if filepath.Base(path) == "module.json" {
			return nil
		}

		files = append(files, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func ReadAllDocTypeJSON(rootPath string) ([]ReadResult, error) {
	files, err := FindDocTypeJSONFiles(rootPath)
	if err != nil {
		return nil, err
	}

	results := make([]ReadResult, 0, len(files))

	for _, file := range files {
		result, err := ReadDocTypeJSON(file)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}
