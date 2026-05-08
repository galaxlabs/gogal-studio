package migration

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/galaxylabs/gogal-studio/internal/core/meta"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const ConfirmApplyToken = "APPLY_SCHEMA_CHANGES"

const (
	OperationStatusPending = "pending"
	OperationStatusApplied = "applied"
	OperationStatusSkipped = "skipped"
	OperationStatusFailed  = "failed"
)

type ColumnInfo struct {
	Name       string
	DataType   string
	UDTName    string
	IsNullable string
	MaxLength  *int
	Precision  *int
	Scale      *int
}

type Operation struct {
	Action    string `json:"action"`
	Table     string `json:"table"`
	Column    string `json:"column,omitempty"`
	Fieldname string `json:"fieldname,omitempty"`
	Fieldtype string `json:"fieldtype,omitempty"`
	SQL       string `json:"sql,omitempty"`
	Message   string `json:"message,omitempty"`
	Dangerous bool   `json:"dangerous"`
	Status    string `json:"status"`
}

type Plan struct {
	DocType    string           `json:"doctype"`
	TableName  string           `json:"table_name"`
	TableFound bool             `json:"table_found"`
	Operations []Operation      `json:"operations"`
	Summary    MigrationSummary `json:"summary"`
}

type MigrationSummary struct {
	AppliedCount           int    `json:"applied_count"`
	SkippedCount           int    `json:"skipped_count"`
	FailedCount            int    `json:"failed_count"`
	HasDangerousOperations bool   `json:"has_dangerous_operations"`
	Status                 string `json:"status"`
}

type Planner struct {
	DB     *pgxpool.Pool
	Loader *meta.Loader
}

func NewPlanner(db *pgxpool.Pool) *Planner {
	return &Planner{
		DB:     db,
		Loader: meta.NewLoader(db),
	}
}

func (p *Planner) Preview(ctx context.Context, doctypeName string) (Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	doc, err := p.Loader.GetDocTypeMeta(ctx, doctypeName)
	if err != nil {
		return Plan{}, err
	}

	validation := ValidateDocTypeMeta(doc)
	if !validation.Valid {
		return Plan{}, fmt.Errorf("schema validation failed: %v", validation.Issues)
	}

	plan := Plan{
		DocType:   doc.Name,
		TableName: doc.TableName,
	}

	tableFound, err := p.tableExists(ctx, doc.TableName)
	if err != nil {
		return Plan{}, err
	}

	plan.TableFound = tableFound

	if !tableFound {
		plan.Operations = append(plan.Operations, Operation{
			Action: "create_table",
			Table:  doc.TableName,
			SQL:    buildCreateTableSQL(doc.TableName),
		})
	}

	existingColumns := map[string]ColumnInfo{}
	if tableFound {
		existingColumns, err = p.existingColumns(ctx, doc.TableName)
		if err != nil {
			return Plan{}, err
		}
	}

	for _, field := range doc.Fields {
		if shouldSkipField(field) {
			continue
		}

		columnName := field.Fieldname

		if tableFound {
			if existingCol, exists := existingColumns[columnName]; exists {
				typeOps := compareColumnType(doc.TableName, field, existingCol)
				plan.Operations = append(plan.Operations, typeOps...)

				nullOps := compareColumnNullability(doc.TableName, field, existingCol)
				plan.Operations = append(plan.Operations, nullOps...)

				if field.UniqueField {
					indexName := buildIndexName(doc.TableName, columnName, "uniq")

					idxExists := false
					dangerous := false
					message := "Unique index can be created safely."

					indexAlreadyExists, err := p.indexExists(ctx, indexName)
					if err != nil {
						dangerous = true
						message = "Could not check unique index existence. Unique index blocked for safety: " + err.Error()
					} else {
						idxExists = indexAlreadyExists
					}

					if !idxExists && !dangerous {
						hasDuplicates, err := p.hasDuplicateValues(ctx, doc.TableName, columnName)
						if err != nil {
							dangerous = true
							message = "Could not verify duplicate values. Unique index blocked for safety: " + err.Error()
						} else if hasDuplicates {
							dangerous = true
							message = "Duplicate values exist. Unique index blocked until duplicates are cleaned."
						} else {
							dangerous = false
							message = "No duplicate values found. Unique index can be applied."
						}
					}

					if idxExists {
						plan.Operations = append(plan.Operations, Operation{
							Action:    "no_change",
							Table:     doc.TableName,
							Column:    columnName,
							Fieldname: field.Fieldname,
							Fieldtype: field.Fieldtype,
							Message:   "Unique index already exists: " + indexName,
							Dangerous: false,
						})
					} else {
						plan.Operations = append(plan.Operations, Operation{
							Action:    "create_unique_index",
							Table:     doc.TableName,
							Column:    columnName,
							Fieldname: field.Fieldname,
							Fieldtype: field.Fieldtype,
							SQL:       fmt.Sprintf(`CREATE UNIQUE INDEX IF NOT EXISTS %s ON %s (%s);`, quoteIdent(indexName), quoteIdent(doc.TableName), quoteIdent(columnName)),
							Message:   message,
							Dangerous: dangerous,
						})
					}
				}

				if field.Fieldtype == "Link" {
					indexName := buildIndexName(doc.TableName, columnName, "idx")

					idxExists := false

					indexAlreadyExists, err := p.indexExists(ctx, indexName)
					if err != nil {
						plan.Operations = append(plan.Operations, Operation{
							Action:    "warning",
							Table:     doc.TableName,
							Column:    columnName,
							Fieldname: field.Fieldname,
							Fieldtype: field.Fieldtype,
							Message:   "Could not check index existence: " + err.Error(),
							Dangerous: false,
						})
					} else {
						idxExists = indexAlreadyExists
					}

					if idxExists {
						plan.Operations = append(plan.Operations, Operation{
							Action:    "no_change",
							Table:     doc.TableName,
							Column:    columnName,
							Fieldname: field.Fieldname,
							Fieldtype: field.Fieldtype,
							Message:   "Index already exists: " + indexName,
							Dangerous: false,
						})
					} else {
						plan.Operations = append(plan.Operations, Operation{
							Action:    "create_index",
							Table:     doc.TableName,
							Column:    columnName,
							Fieldname: field.Fieldname,
							Fieldtype: field.Fieldtype,
							SQL:       fmt.Sprintf(`CREATE INDEX IF NOT EXISTS %s ON %s (%s);`, quoteIdent(indexName), quoteIdent(doc.TableName), quoteIdent(columnName)),
							Message:   "Index for Link field lookup performance.",
							Dangerous: false,
						})
					}
				}

				continue
			}
		}

		sqlType, ok := postgresType(field.Fieldtype)
		if !ok {
			plan.Operations = append(plan.Operations, Operation{
				Action:    "warning",
				Table:     doc.TableName,
				Column:    columnName,
				Fieldname: field.Fieldname,
				Fieldtype: field.Fieldtype,
				Message:   "Unsupported field type. Column will not be created.",
				Dangerous: false,
			})
			continue
		}

		plan.Operations = append(plan.Operations, Operation{
			Action:    "add_column",
			Table:     doc.TableName,
			Column:    columnName,
			Fieldname: field.Fieldname,
			Fieldtype: field.Fieldtype,
			SQL:       fmt.Sprintf(`ALTER TABLE %s ADD COLUMN IF NOT EXISTS %s %s;`, quoteIdent(doc.TableName), quoteIdent(columnName), sqlType),
			Dangerous: false,
		})

		if field.Reqd {
			if plan.TableFound {
				if hasFieldDefault(field) {
					defaultSQL := sqlDefaultLiteral(field)

					plan.Operations = append(plan.Operations, Operation{
						Action:    "set_default",
						Table:     doc.TableName,
						Column:    columnName,
						Fieldname: field.Fieldname,
						Fieldtype: field.Fieldtype,
						SQL:       fmt.Sprintf(`ALTER TABLE %s ALTER COLUMN %s SET DEFAULT %s;`, quoteIdent(doc.TableName), quoteIdent(columnName), defaultSQL),
						Message:   "Default added before backfill.",
						Dangerous: false,
					})

					plan.Operations = append(plan.Operations, Operation{
						Action:    "backfill_default",
						Table:     doc.TableName,
						Column:    columnName,
						Fieldname: field.Fieldname,
						Fieldtype: field.Fieldtype,
						SQL:       fmt.Sprintf(`UPDATE %s SET %s = %s WHERE %s IS NULL;`, quoteIdent(doc.TableName), quoteIdent(columnName), defaultSQL, quoteIdent(columnName)),
						Message:   "Existing NULL rows will be backfilled using default value.",
						Dangerous: false,
					})

					plan.Operations = append(plan.Operations, Operation{
						Action:    "set_not_null",
						Table:     doc.TableName,
						Column:    columnName,
						Fieldname: field.Fieldname,
						Fieldtype: field.Fieldtype,
						SQL:       fmt.Sprintf(`ALTER TABLE %s ALTER COLUMN %s SET NOT NULL;`, quoteIdent(doc.TableName), quoteIdent(columnName)),
						Message:   "Safe because default/backfill plan exists.",
						Dangerous: false,
					})
				} else {
					plan.Operations = append(plan.Operations, Operation{
						Action:    "set_not_null",
						Table:     doc.TableName,
						Column:    columnName,
						Fieldname: field.Fieldname,
						Fieldtype: field.Fieldtype,
						SQL:       fmt.Sprintf(`ALTER TABLE %s ALTER COLUMN %s SET NOT NULL;`, quoteIdent(doc.TableName), quoteIdent(columnName)),
						Message:   "Dangerous on existing table because no default/backfill value exists.",
						Dangerous: true,
					})
				}
			} else {
				plan.Operations = append(plan.Operations, Operation{
					Action:    "set_not_null",
					Table:     doc.TableName,
					Column:    columnName,
					Fieldname: field.Fieldname,
					Fieldtype: field.Fieldtype,
					SQL:       fmt.Sprintf(`ALTER TABLE %s ALTER COLUMN %s SET NOT NULL;`, quoteIdent(doc.TableName), quoteIdent(columnName)),
					Message:   "Safe on new table after column creation.",
					Dangerous: false,
				})
			}
		}

		if field.UniqueField {
			indexName := buildIndexName(doc.TableName, columnName, "uniq")

			exists := false
			dangerous := false
			message := "Unique index can be created safely."

			if plan.TableFound {
				indexAlreadyExists, err := p.indexExists(ctx, indexName)
				if err != nil {
					dangerous = true
					message = "Could not check unique index existence. Unique index blocked for safety: " + err.Error()
				} else {
					exists = indexAlreadyExists
				}

				if !exists && !dangerous {
					hasDuplicates, err := p.hasDuplicateValues(ctx, doc.TableName, columnName)
					if err != nil {
						dangerous = true
						message = "Could not verify duplicate values. Unique index blocked for safety: " + err.Error()
					} else if hasDuplicates {
						dangerous = true
						message = "Duplicate values exist. Unique index blocked until duplicates are cleaned."
					} else {
						dangerous = false
						message = "No duplicate values found. Unique index can be applied."
					}
				}
			}

			if exists {
				plan.Operations = append(plan.Operations, Operation{
					Action:    "no_change",
					Table:     doc.TableName,
					Column:    columnName,
					Fieldname: field.Fieldname,
					Fieldtype: field.Fieldtype,
					Message:   "Unique index already exists: " + indexName,
					Dangerous: false,
				})
			} else {
				plan.Operations = append(plan.Operations, Operation{
					Action:    "create_unique_index",
					Table:     doc.TableName,
					Column:    columnName,
					Fieldname: field.Fieldname,
					Fieldtype: field.Fieldtype,
					SQL:       fmt.Sprintf(`CREATE UNIQUE INDEX IF NOT EXISTS %s ON %s (%s);`, quoteIdent(indexName), quoteIdent(doc.TableName), quoteIdent(columnName)),
					Message:   message,
					Dangerous: dangerous,
				})
			}
		}

		if field.Fieldtype == "Link" {
			indexName := buildIndexName(doc.TableName, columnName, "idx")

			exists := false

			if plan.TableFound {
				indexAlreadyExists, err := p.indexExists(ctx, indexName)
				if err != nil {
					plan.Operations = append(plan.Operations, Operation{
						Action:    "warning",
						Table:     doc.TableName,
						Column:    columnName,
						Fieldname: field.Fieldname,
						Fieldtype: field.Fieldtype,
						Message:   "Could not check index existence: " + err.Error(),
						Dangerous: false,
					})
				} else {
					exists = indexAlreadyExists
				}
			}

			if exists {
				plan.Operations = append(plan.Operations, Operation{
					Action:    "no_change",
					Table:     doc.TableName,
					Column:    columnName,
					Fieldname: field.Fieldname,
					Fieldtype: field.Fieldtype,
					Message:   "Index already exists: " + indexName,
					Dangerous: false,
				})
			} else {
				plan.Operations = append(plan.Operations, Operation{
					Action:    "create_index",
					Table:     doc.TableName,
					Column:    columnName,
					Fieldname: field.Fieldname,
					Fieldtype: field.Fieldtype,
					SQL:       fmt.Sprintf(`CREATE INDEX IF NOT EXISTS %s ON %s (%s);`, quoteIdent(indexName), quoteIdent(doc.TableName), quoteIdent(columnName)),
					Message:   "Index for Link field lookup performance.",
					Dangerous: false,
				})
			}
		}
	}

	if len(plan.Operations) == 0 {
		plan.Operations = append(plan.Operations, Operation{
			Action:  "no_change",
			Table:   doc.TableName,
			Message: "No schema changes needed.",
		})
	}

	cleanupPlanOperations(&plan)
	normalizeOperationStatuses(&plan)
	updatePlanSummary(&plan)

	return plan, nil
}

func (p *Planner) tableExists(ctx context.Context, tableName string) (bool, error) {
	var exists bool

	err := p.DB.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = $1
		)
	`, tableName).Scan(&exists)

	return exists, err
}

func (p *Planner) existingColumns(ctx context.Context, tableName string) (map[string]ColumnInfo, error) {
	rows, err := p.DB.Query(ctx, `
		SELECT
			column_name,
			data_type,
			udt_name,
			is_nullable,
			character_maximum_length,
			numeric_precision,
			numeric_scale
		FROM information_schema.columns
		WHERE table_schema = 'public'
		AND table_name = $1
	`, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := map[string]ColumnInfo{}

	for rows.Next() {
		var col ColumnInfo

		if err := rows.Scan(
			&col.Name,
			&col.DataType,
			&col.UDTName,
			&col.IsNullable,
			&col.MaxLength,
			&col.Precision,
			&col.Scale,
		); err != nil {
			return nil, err
		}

		columns[col.Name] = col
	}

	return columns, rows.Err()
}

func buildCreateTableSQL(tableName string) string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	id BIGSERIAL PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	owner TEXT,
	creation TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	modified TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	modified_by TEXT,
	docstatus INT NOT NULL DEFAULT 0,
	idx INT NOT NULL DEFAULT 0
);`, quoteIdent(tableName))
}

func shouldSkipField(field meta.FieldMeta) bool {
	if field.Fieldname == "" {
		return true
	}

	systemFields := map[string]bool{
		"id":          true,
		"name":        true,
		"owner":       true,
		"creation":    true,
		"modified":    true,
		"modified_by": true,
		"docstatus":   true,
		"idx":         true,
	}

	if systemFields[field.Fieldname] {
		return true
	}

	skipTypes := map[string]bool{
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

	return skipTypes[field.Fieldtype]
}

func postgresType(fieldtype string) (string, bool) {
	types := map[string]string{
		"Data":         "TEXT",
		"Text":         "TEXT",
		"Small Text":   "TEXT",
		"Long Text":    "TEXT",
		"Code":         "TEXT",
		"JSON":         "JSONB",
		"Int":          "BIGINT",
		"Float":        "DOUBLE PRECISION",
		"Currency":     "NUMERIC(18,6)",
		"Percent":      "NUMERIC(18,6)",
		"Check":        "BOOLEAN NOT NULL DEFAULT FALSE",
		"Date":         "DATE",
		"Datetime":     "TIMESTAMPTZ",
		"Time":         "TIME",
		"Select":       "TEXT",
		"Link":         "TEXT",
		"Dynamic Link": "TEXT",
		"Table":        "JSONB",
		"Attach":       "TEXT",
		"Attach Image": "TEXT",
		"Password":     "TEXT",
		"Phone":        "TEXT",
		"Email":        "TEXT",
		"URL":          "TEXT",
		"Color":        "TEXT",
		"Rating":       "INT",
		"Duration":     "BIGINT",
	}

	value, ok := types[fieldtype]
	return value, ok
}

func (p *Planner) Apply(ctx context.Context, doctypeName string, confirm string) (Plan, error) {
	if confirm != ConfirmApplyToken {
		return Plan{}, fmt.Errorf("schema apply blocked: invalid confirmation token")
	}

	plan, err := p.Preview(ctx, doctypeName)
	if err != nil {
		return Plan{}, err
	}

	tx, err := p.DB.Begin(ctx)
	if err != nil {
		return Plan{}, err
	}
	defer tx.Rollback(ctx)

	if err := ensureMigrationLogTable(ctx, tx); err != nil {
		return Plan{}, err
	}

	appliedCount := 0
	skippedCount := 0

	for index, op := range plan.Operations {
		if op.SQL == "" {
			continue
		}

		if op.Action == "warning" || op.Action == "no_change" {
			continue
		}

		if !isSafeApplyOperation(op) {
			op.Message = buildSkippedMessage(op)
			op.Status = OperationStatusSkipped

			if err := insertMigrationLog(ctx, tx, plan, op, "skipped", op.Message); err != nil {
				return Plan{}, err
			}

			plan.Operations[index].Message = op.Message
			plan.Operations[index].Dangerous = true
			plan.Operations[index].Status = OperationStatusSkipped
			skippedCount++

			continue
		}

		if _, err := tx.Exec(ctx, op.SQL); err != nil {
			plan.Operations[index].Status = OperationStatusFailed
			plan.Operations[index].Message = err.Error()
			op.Status = OperationStatusFailed
			_ = insertMigrationLog(ctx, tx, plan, op, "failed", err.Error())
			return Plan{}, fmt.Errorf("apply operation %s failed: %w", op.Action, err)
		}

		plan.Operations[index].Status = OperationStatusApplied

		if err := insertMigrationLog(ctx, tx, plan, op, "success", "applied"); err != nil {
			return Plan{}, err
		}

		appliedCount++
	}

	if appliedCount == 0 && skippedCount == 0 {
		noChangeOp := Operation{
			Action:  "no_change",
			Table:   plan.TableName,
			Message: "No schema changes applied.",
			Status:  OperationStatusApplied,
		}

		if err := insertMigrationLog(ctx, tx, plan, noChangeOp, "success", "no changes"); err != nil {
			return Plan{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Plan{}, err
	}

	appliedPlan, err := p.Preview(ctx, doctypeName)
	if err != nil {
		cleanupPlanOperations(&plan)
		keepMeaningfulApplyOperations(&plan)
		normalizeOperationStatuses(&plan)
		updatePlanSummary(&plan)
		return plan, nil
	}

	for _, op := range plan.Operations {
		if op.Dangerous {
			appliedPlan.Operations = append(appliedPlan.Operations, Operation{
				Action:    "skipped_" + op.Action,
				Table:     op.Table,
				Column:    op.Column,
				Fieldname: op.Fieldname,
				Fieldtype: op.Fieldtype,
				SQL:       op.SQL,
				Message:   buildSkippedMessage(op),
				Dangerous: true,
				Status:    OperationStatusSkipped,
			})
		}
	}

	cleanupPlanOperations(&appliedPlan)
	keepMeaningfulApplyOperations(&appliedPlan)
	normalizeOperationStatuses(&appliedPlan)
	updatePlanSummary(&appliedPlan)

	return appliedPlan, nil
}

func isSafeApplyOperation(op Operation) bool {
	if op.Dangerous {
		return false
	}

	safeActions := map[string]bool{
		"create_table":        true,
		"add_column":          true,
		"create_index":        true,
		"create_unique_index": true,
		"set_default":         true,
		"backfill_default":    true,
		"set_not_null":        true,
		"drop_not_null":       true,
		"alter_column_type":   true,
	}

	return safeActions[op.Action]
}

func normalizeOperationStatuses(plan *Plan) {
	for i := range plan.Operations {
		if plan.Operations[i].Status == "" {
			plan.Operations[i].Status = OperationStatusPending
		}
	}
}

func updatePlanSummary(plan *Plan) {
	summary := MigrationSummary{}

	for _, op := range plan.Operations {
		if op.Dangerous {
			summary.HasDangerousOperations = true
		}

		switch op.Status {
		case OperationStatusApplied:
			if op.SQL != "" {
				summary.AppliedCount++
			}
		case OperationStatusSkipped:
			summary.SkippedCount++
		case OperationStatusFailed:
			summary.FailedCount++
		}
	}

	switch {
	case summary.FailedCount > 0:
		summary.Status = "failed"
	case summary.AppliedCount > 0 && summary.SkippedCount > 0:
		summary.Status = "partial_success"
	case summary.AppliedCount > 0:
		summary.Status = "success"
	case summary.SkippedCount > 0:
		summary.Status = "skipped"
	default:
		summary.Status = "no_change"
	}

	plan.Summary = summary
}

func buildSkippedMessage(op Operation) string {
	if op.Message != "" {
		return "Skipped dangerous operation: " + op.Message
	}

	switch op.Action {
	case "set_not_null":
		return "Skipped dangerous operation: NOT NULL requires backfill/default safety check."
	case "create_unique_index":
		return "Skipped dangerous operation: UNIQUE index may fail if duplicate values already exist."
	default:
		return "Skipped dangerous operation: requires explicit safe migration support."
	}
}

func cleanupPlanOperations(plan *Plan) {
	cleaned := make([]Operation, 0, len(plan.Operations))
	seen := map[string]bool{}

	for _, op := range plan.Operations {
		key := operationDedupKey(op)

		// Keep all real change operations.
		if op.SQL != "" || op.Dangerous || op.Action == "warning" {
			if !seen[key] {
				cleaned = append(cleaned, op)
				seen[key] = true
			}
			continue
		}

		// Deduplicate no_change noise.
		if op.Action == "no_change" {
			if !seen[key] {
				cleaned = append(cleaned, op)
				seen[key] = true
			}
			continue
		}

		if !seen[key] {
			cleaned = append(cleaned, op)
			seen[key] = true
		}
	}

	plan.Operations = cleaned
}

func operationDedupKey(op Operation) string {
	return strings.Join([]string{
		op.Action,
		op.Table,
		op.Column,
		op.Fieldname,
		op.Fieldtype,
		op.SQL,
		op.Message,
	}, "|")
}

func keepMeaningfulApplyOperations(plan *Plan) {
	filtered := make([]Operation, 0, len(plan.Operations))

	for _, op := range plan.Operations {
		// Always keep applied/skipped/failed operations.
		if op.Status == OperationStatusApplied ||
			op.Status == OperationStatusSkipped ||
			op.Status == OperationStatusFailed {
			filtered = append(filtered, op)
			continue
		}

		// Keep dangerous operations because UI must show them.
		if op.Dangerous {
			filtered = append(filtered, op)
			continue
		}

		// Keep warnings.
		if op.Action == "warning" {
			filtered = append(filtered, op)
			continue
		}

		// Hide normal no_change rows from apply response.
		if op.Action == "no_change" {
			continue
		}

		filtered = append(filtered, op)
	}

	if len(filtered) == 0 {
		filtered = append(filtered, Operation{
			Action:  "no_change",
			Table:   plan.TableName,
			Message: "No meaningful schema changes.",
			Status:  OperationStatusApplied,
		})
	}

	plan.Operations = filtered
}

func quoteIdent(value string) string {
	safe := strings.ReplaceAll(value, `"`, `""`)
	return `"` + safe + `"`
}

func SafeTableName(doctypeName string) string {
	value := strings.ToLower(strings.TrimSpace(doctypeName))

	re := regexp.MustCompile(`[^a-z0-9]+`)
	value = re.ReplaceAllString(value, "_")
	value = strings.Trim(value, "_")

	if value == "" {
		value = "document"
	}

	if value[0] >= '0' && value[0] <= '9' {
		value = "dt_" + value
	}

	return "tab_" + value
}

func ensureMigrationLogTable(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS "tabSchema Migration Log" (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			doctype TEXT NOT NULL,
			table_name TEXT NOT NULL,
			operation TEXT NOT NULL,
			column_name TEXT,
			fieldname TEXT,
			fieldtype TEXT,
			sql TEXT,
			status TEXT NOT NULL,
			message TEXT,
			applied_by TEXT NOT NULL DEFAULT 'Administrator',
			creation TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			modified TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			idx INT NOT NULL DEFAULT 0
		)
	`)
	return err
}

func insertMigrationLog(ctx context.Context, tx pgx.Tx, plan Plan, op Operation, status string, message string) error {
	logName := fmt.Sprintf(
		"%s-%s-%d",
		plan.DocType,
		op.Action,
		time.Now().UnixNano(),
	)

	_, err := tx.Exec(ctx, `
		INSERT INTO "tabSchema Migration Log" (
			name,
			doctype,
			table_name,
			operation,
			column_name,
			fieldname,
			fieldtype,
			sql,
			status,
			message
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`,
		logName,
		plan.DocType,
		plan.TableName,
		op.Action,
		nullableString(op.Column),
		nullableString(op.Fieldname),
		nullableString(op.Fieldtype),
		nullableString(op.SQL),
		status,
		message,
	)

	return err
}

func nullableString(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	return value
}

func compareColumnNullability(tableName string, field meta.FieldMeta, existing ColumnInfo) []Operation {
	isNullable := strings.EqualFold(existing.IsNullable, "YES")

	if field.Reqd && isNullable {
		if hasFieldDefault(field) {
			defaultSQL := sqlDefaultLiteral(field)

			return []Operation{
				{
					Action:    "set_default",
					Table:     tableName,
					Column:    field.Fieldname,
					Fieldname: field.Fieldname,
					Fieldtype: field.Fieldtype,
					SQL:       fmt.Sprintf(`ALTER TABLE %s ALTER COLUMN %s SET DEFAULT %s;`, quoteIdent(tableName), quoteIdent(field.Fieldname), defaultSQL),
					Message:   "Default added before NOT NULL backfill.",
					Dangerous: false,
				},
				{
					Action:    "backfill_default",
					Table:     tableName,
					Column:    field.Fieldname,
					Fieldname: field.Fieldname,
					Fieldtype: field.Fieldtype,
					SQL:       fmt.Sprintf(`UPDATE %s SET %s = %s WHERE %s IS NULL;`, quoteIdent(tableName), quoteIdent(field.Fieldname), defaultSQL, quoteIdent(field.Fieldname)),
					Message:   "Existing NULL rows will be backfilled using default value.",
					Dangerous: false,
				},
				{
					Action:    "set_not_null",
					Table:     tableName,
					Column:    field.Fieldname,
					Fieldname: field.Fieldname,
					Fieldtype: field.Fieldtype,
					SQL:       fmt.Sprintf(`ALTER TABLE %s ALTER COLUMN %s SET NOT NULL;`, quoteIdent(tableName), quoteIdent(field.Fieldname)),
					Message:   "Safe because default/backfill plan exists.",
					Dangerous: false,
				},
			}
		}

		return []Operation{
			{
				Action:    "set_not_null",
				Table:     tableName,
				Column:    field.Fieldname,
				Fieldname: field.Fieldname,
				Fieldtype: field.Fieldtype,
				SQL:       fmt.Sprintf(`ALTER TABLE %s ALTER COLUMN %s SET NOT NULL;`, quoteIdent(tableName), quoteIdent(field.Fieldname)),
				Message:   "Dangerous: existing column is nullable and required field has no default/backfill value.",
				Dangerous: true,
			},
		}
	}

	if !field.Reqd && !isNullable {
		return []Operation{
			{
				Action:    "drop_not_null",
				Table:     tableName,
				Column:    field.Fieldname,
				Fieldname: field.Fieldname,
				Fieldtype: field.Fieldtype,
				SQL:       fmt.Sprintf(`ALTER TABLE %s ALTER COLUMN %s DROP NOT NULL;`, quoteIdent(tableName), quoteIdent(field.Fieldname)),
				Message:   "Optional relaxation: metadata says field is not required, but database column is NOT NULL.",
				Dangerous: false,
			},
		}
	}

	return []Operation{
		{
			Action:    "no_change",
			Table:     tableName,
			Column:    field.Fieldname,
			Fieldname: field.Fieldname,
			Fieldtype: field.Fieldtype,
			Message:   "Column nullability matches metadata.",
			Dangerous: false,
		},
	}
}

func compareColumnType(tableName string, field meta.FieldMeta, existing ColumnInfo) []Operation {
	expectedType, ok := postgresType(field.Fieldtype)
	if !ok {
		return []Operation{
			{
				Action:    "warning",
				Table:     tableName,
				Column:    field.Fieldname,
				Fieldname: field.Fieldname,
				Fieldtype: field.Fieldtype,
				Message:   "Unsupported fieldtype. Cannot compare column type.",
				Dangerous: false,
			},
		}
	}

	expectedKind := normalizeSQLType(expectedType)
	actualKind := normalizePostgresColumnType(existing)

	if expectedKind == actualKind {
		return []Operation{
			{
				Action:    "no_change",
				Table:     tableName,
				Column:    field.Fieldname,
				Fieldname: field.Fieldname,
				Fieldtype: field.Fieldtype,
				Message:   "Column type matches metadata.",
				Dangerous: false,
			},
		}
	}

	if isSafeWidening(actualKind, expectedKind) {
		return []Operation{
			{
				Action:    "alter_column_type",
				Table:     tableName,
				Column:    field.Fieldname,
				Fieldname: field.Fieldname,
				Fieldtype: field.Fieldtype,
				SQL: fmt.Sprintf(
					`ALTER TABLE %s ALTER COLUMN %s TYPE %s USING %s::%s;`,
					quoteIdent(tableName),
					quoteIdent(field.Fieldname),
					expectedType,
					quoteIdent(field.Fieldname),
					expectedType,
				),
				Message:   "Safe widening type change.",
				Dangerous: false,
			},
		}
	}

	return []Operation{
		{
			Action:    "alter_column_type",
			Table:     tableName,
			Column:    field.Fieldname,
			Fieldname: field.Fieldname,
			Fieldtype: field.Fieldtype,
			SQL: fmt.Sprintf(
				`ALTER TABLE %s ALTER COLUMN %s TYPE %s USING %s::%s;`,
				quoteIdent(tableName),
				quoteIdent(field.Fieldname),
				expectedType,
				quoteIdent(field.Fieldname),
				expectedType,
			),
			Message:   fmt.Sprintf("Risky type change blocked. Existing: %s, Expected: %s", actualKind, expectedKind),
			Dangerous: true,
		},
	}
}

func normalizeSQLType(sqlType string) string {
	value := strings.ToUpper(strings.TrimSpace(sqlType))

	switch {
	case value == "TEXT":
		return "text"
	case value == "JSONB":
		return "jsonb"
	case value == "BIGINT":
		return "bigint"
	case value == "INT" || value == "INTEGER":
		return "integer"
	case value == "DOUBLE PRECISION":
		return "double"
	case strings.HasPrefix(value, "NUMERIC"):
		return "numeric"
	case strings.HasPrefix(value, "BOOLEAN"):
		return "boolean"
	case value == "DATE":
		return "date"
	case value == "TIMESTAMPTZ" || value == "TIMESTAMP WITH TIME ZONE":
		return "timestamptz"
	case value == "TIME":
		return "time"
	default:
		return strings.ToLower(value)
	}
}

func normalizePostgresColumnType(col ColumnInfo) string {
	dataType := strings.ToLower(strings.TrimSpace(col.DataType))
	udtName := strings.ToLower(strings.TrimSpace(col.UDTName))

	switch {
	case dataType == "text":
		return "text"
	case dataType == "character varying":
		return "varchar"
	case dataType == "jsonb":
		return "jsonb"
	case dataType == "bigint" || udtName == "int8":
		return "bigint"
	case dataType == "integer" || udtName == "int4":
		return "integer"
	case dataType == "double precision" || udtName == "float8":
		return "double"
	case dataType == "numeric":
		return "numeric"
	case dataType == "boolean" || udtName == "bool":
		return "boolean"
	case dataType == "date":
		return "date"
	case dataType == "timestamp with time zone" || udtName == "timestamptz":
		return "timestamptz"
	case dataType == "time without time zone" || udtName == "time":
		return "time"
	default:
		if udtName != "" {
			return udtName
		}
		return dataType
	}
}

func isSafeWidening(actual string, expected string) bool {
	if actual == expected {
		return true
	}

	safe := map[string]map[string]bool{
		"varchar": {
			"text": true,
		},
		"integer": {
			"bigint":  true,
			"numeric": true,
			"double":  true,
		},
		"bigint": {
			"numeric": true,
			"double":  true,
		},
		"numeric": {
			"double": true,
		},
	}

	return safe[actual][expected]
}

func (p *Planner) indexExists(ctx context.Context, indexName string) (bool, error) {
	var exists bool

	err := p.DB.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM pg_indexes
			WHERE schemaname = 'public'
			AND indexname = $1
		)
	`, indexName).Scan(&exists)

	return exists, err
}

func (p *Planner) hasDuplicateValues(ctx context.Context, tableName string, columnName string) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM %s
			WHERE %s IS NOT NULL
			GROUP BY %s
			HAVING COUNT(*) > 1
			LIMIT 1
		)
	`, quoteIdent(tableName), quoteIdent(columnName), quoteIdent(columnName))

	var exists bool

	if err := p.DB.QueryRow(ctx, query).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func fieldDefaultValue(field meta.FieldMeta) string {
	return strings.TrimSpace(field.DefaultValue)
}

func hasFieldDefault(field meta.FieldMeta) bool {
	return fieldDefaultValue(field) != ""
}

func sqlDefaultLiteral(field meta.FieldMeta) string {
	value := fieldDefaultValue(field)

	switch field.Fieldtype {
	case "Int", "Rating", "Duration":
		return value
	case "Float", "Currency", "Percent":
		return value
	case "Check":
		if value == "1" || strings.EqualFold(value, "true") || strings.EqualFold(value, "yes") {
			return "TRUE"
		}
		return "FALSE"
	default:
		return "'" + strings.ReplaceAll(value, "'", "''") + "'"
	}
}

func buildIndexName(tableName string, columnName string, suffix string) string {
	clean := strings.ToLower(tableName + "_" + columnName + "_" + suffix)

	re := regexp.MustCompile(`[^a-z0-9_]+`)
	clean = re.ReplaceAllString(clean, "_")
	clean = strings.Trim(clean, "_")

	if len(clean) > 60 {
		clean = clean[:60]
	}

	return clean
}
