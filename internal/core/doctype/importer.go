package doctype

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Importer writes JSONDocType metadata (read from JSON files) into the database.
// Each import is atomic: the DocType row is upserted, child fields and permissions
// are replaced in full within a single transaction.
type Importer struct {
	DB *pgxpool.Pool
}

// ImportResult describes the outcome of a single import.
type ImportResult struct {
	DocType  string `json:"doctype"`
	JSONPath string `json:"json_path"`
	Fields   int    `json:"fields"`
	Perms    int    `json:"perms"`
}

func NewImporter(db *pgxpool.Pool) *Importer {
	return &Importer{DB: db}
}

// ImportOneFromFile reads jsonPath and imports it into the database.
func (im *Importer) ImportOneFromFile(jsonPath string) (ImportResult, error) {
	readResult, err := ReadDocTypeJSON(jsonPath)
	if err != nil {
		return ImportResult{}, err
	}

	return im.ImportOne(readResult.JSONPath, readResult.DocType)
}

// ImportOneByName reads the JSON for the given module+doctype and imports it.
func (im *Importer) ImportOneByName(rootPath, moduleName, doctypeName string) (ImportResult, error) {
	readResult, err := ReadDocTypeJSONByName(rootPath, moduleName, doctypeName)
	if err != nil {
		return ImportResult{}, err
	}

	return im.ImportOne(readResult.JSONPath, readResult.DocType)
}

// ImportAll reads every DocType JSON file under rootPath and imports each one.
func (im *Importer) ImportAll(rootPath string) ([]ImportResult, error) {
	readResults, err := ReadAllDocTypeJSON(rootPath)
	if err != nil {
		return nil, err
	}

	results := make([]ImportResult, 0, len(readResults))

	for _, rr := range readResults {
		result, err := im.ImportOne(rr.JSONPath, rr.DocType)
		if err != nil {
			return nil, fmt.Errorf("import %s: %w", rr.JSONPath, err)
		}

		results = append(results, result)
	}

	return results, nil
}

// ImportOne validates doc and upserts it into the database inside a transaction.
func (im *Importer) ImportOne(jsonPath string, doc JSONDocType) (ImportResult, error) {
	doc = normalizeJSONDocType(doc)

	if err := ValidateJSONDocType(doc); err != nil {
		return ImportResult{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := im.DB.Begin(ctx)
	if err != nil {
		return ImportResult{}, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if err := im.upsertDocType(ctx, tx, doc); err != nil {
		return ImportResult{}, err
	}

	if err := im.replaceFields(ctx, tx, doc); err != nil {
		return ImportResult{}, err
	}

	if err := im.replacePerms(ctx, tx, doc); err != nil {
		return ImportResult{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return ImportResult{}, err
	}

	return ImportResult{
		DocType:  doc.Name,
		JSONPath: jsonPath,
		Fields:   len(doc.Fields),
		Perms:    len(doc.Permissions),
	}, nil
}
