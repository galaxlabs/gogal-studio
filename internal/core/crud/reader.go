package crud

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/galaxylabs/gogal-studio/internal/core/meta"
	"github.com/galaxylabs/gogal-studio/internal/core/permission"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Reader is a read-only document accessor gated by permission checks.
type Reader struct {
	DB         *pgxpool.Pool
	MetaLoader *meta.Loader
	Perms      *permission.Checker
}

// NewReader creates a Reader backed by the given connection pool.
func NewReader(db *pgxpool.Pool) *Reader {
	return &Reader{
		DB:         db,
		MetaLoader: meta.NewLoader(db),
		Perms:      permission.NewChecker(db),
	}
}

// ListOptions controls how List behaves.
type ListOptions struct {
	User    string
	Limit   int
	Offset  int
	Fields  []string
	Filters map[string]any
}

// ReadResult is returned by List.
type ReadResult struct {
	Data                 []map[string]any `json:"data"`
	Columns              []string         `json:"columns"`
	MissingColumns       []string         `json:"missing_columns"`
	MissingFilterColumns []string         `json:"missing_filter_columns"`
}

// GetResult is returned by Get.
type GetResult struct {
	Data           map[string]any `json:"data"`
	Columns        []string       `json:"columns"`
	MissingColumns []string       `json:"missing_columns"`
}

// List returns up to Limit rows from the doctype table.
func (r *Reader) List(ctx context.Context, doctype string, opts ListOptions) (ReadResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	doctype = strings.TrimSpace(doctype)
	user := strings.TrimSpace(opts.User)

	if doctype == "" {
		return ReadResult{}, fmt.Errorf("doctype is required")
	}
	if user == "" {
		user = "Administrator"
	}

	allowed, err := r.Perms.CanUserRead(ctx, user, doctype)
	if err != nil {
		return ReadResult{}, err
	}
	if !allowed {
		return ReadResult{}, fmt.Errorf("permission denied: user %s cannot read %s", user, doctype)
	}

	m, err := r.MetaLoader.GetDocTypeMeta(ctx, doctype)
	if err != nil {
		return ReadResult{}, err
	}

	requested := readableColumns(m)
	if len(opts.Fields) > 0 {
		requested = opts.Fields
	}

	columns, missingColumns, err := r.safeColumns(ctx, m.TableName, requested)
	if err != nil {
		return ReadResult{}, err
	}
	if len(columns) == 0 {
		return ReadResult{}, fmt.Errorf("no readable database columns found for %s", doctype)
	}

	safeFilters, missingFilterColumns, err := r.safeFilters(ctx, m.TableName, opts.Filters)
	if err != nil {
		return ReadResult{}, err
	}

	limit := opts.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := opts.Offset
	if offset < 0 {
		offset = 0
	}

	whereSQL, args := buildWhereClause(safeFilters)
	limitPh := len(args) + 1
	offsetPh := len(args) + 2
	args = append(args, limit, offset)

	query := fmt.Sprintf(
		`SELECT %s FROM %s %s ORDER BY "name" ASC LIMIT $%d OFFSET $%d`,
		joinQuotedColumns(columns),
		quoteIdent(m.TableName),
		whereSQL,
		limitPh,
		offsetPh,
	)

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return ReadResult{}, err
	}
	defer rows.Close()

	data, err := scanRows(rows, columns)
	if err != nil {
		return ReadResult{}, err
	}

	return ReadResult{
		Data:                 data,
		Columns:              columns,
		MissingColumns:       missingColumns,
		MissingFilterColumns: missingFilterColumns,
	}, nil
}

// Get returns a single document by name.
func (r *Reader) Get(ctx context.Context, doctype, name, user string) (GetResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	doctype = strings.TrimSpace(doctype)
	name = strings.TrimSpace(name)
	user = strings.TrimSpace(user)

	if doctype == "" {
		return GetResult{}, fmt.Errorf("doctype is required")
	}
	if name == "" {
		return GetResult{}, fmt.Errorf("name is required")
	}
	if user == "" {
		user = "Administrator"
	}

	allowed, err := r.Perms.CanUserRead(ctx, user, doctype)
	if err != nil {
		return GetResult{}, err
	}
	if !allowed {
		return GetResult{}, fmt.Errorf("permission denied: user %s cannot read %s", user, doctype)
	}

	m, err := r.MetaLoader.GetDocTypeMeta(ctx, doctype)
	if err != nil {
		return GetResult{}, err
	}

	requested := readableColumns(m)
	columns, missing, err := r.safeColumns(ctx, m.TableName, requested)
	if err != nil {
		return GetResult{}, err
	}
	if len(columns) == 0 {
		return GetResult{}, fmt.Errorf("no readable database columns found for %s", doctype)
	}

	query := fmt.Sprintf(
		`SELECT %s FROM %s WHERE name = $1 LIMIT 1`,
		joinQuotedColumns(columns),
		quoteIdent(m.TableName),
	)

	rows, err := r.DB.Query(ctx, query, name)
	if err != nil {
		return GetResult{}, err
	}
	defer rows.Close()

	data, err := scanRows(rows, columns)
	if err != nil {
		return GetResult{}, err
	}
	if len(data) == 0 {
		return GetResult{}, fmt.Errorf("document not found: %s/%s", doctype, name)
	}

	return GetResult{Data: data[0], Columns: columns, MissingColumns: missing}, nil
}

// ─── helpers ─────────────────────────────────────────────────────────────────

// readableColumns builds the SELECT column list: system fields first, then
// non-hidden, non-layout meta fields.
func readableColumns(m meta.DocTypeMeta) []string {
	columns := []string{
		"name", "owner", "creation", "modified", "modified_by", "docstatus", "idx",
	}
	seen := make(map[string]bool, len(columns)+len(m.Fields))
	for _, c := range columns {
		seen[c] = true
	}
	for _, field := range m.Fields {
		fn := strings.TrimSpace(field.Fieldname)
		if fn == "" || seen[fn] || field.Hidden || isLayoutFieldtype(field.Fieldtype) {
			continue
		}
		columns = append(columns, fn)
		seen[fn] = true
	}
	return columns
}

var layoutFieldtypes = map[string]bool{
	"Section Break": true,
	"Column Break":  true,
	"Tab Break":     true,
	"HTML":          true,
	"Button":        true,
	"Heading":       true,
	"Fold":          true,
	"Table":         true,
}

func isLayoutFieldtype(ft string) bool { return layoutFieldtypes[ft] }

func joinQuotedColumns(columns []string) string {
	out := make([]string, len(columns))
	for i, c := range columns {
		out[i] = quoteIdent(c)
	}
	return strings.Join(out, ", ")
}

func quoteIdent(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

// safeColumns filters requested column names against what actually exists in the table.
func (r *Reader) safeColumns(ctx context.Context, tableName string, requested []string) (safe []string, missing []string, err error) {
	actual, err := r.actualColumns(ctx, tableName)
	if err != nil {
		return nil, nil, err
	}
	for _, col := range requested {
		if actual[col] {
			safe = append(safe, col)
		} else {
			missing = append(missing, col)
		}
	}
	return safe, missing, nil
}

// actualColumns queries information_schema for the real columns of a table.
func (r *Reader) actualColumns(ctx context.Context, tableName string) (map[string]bool, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT column_name
		FROM information_schema.columns
		WHERE table_schema = 'public'
		  AND table_name = $1
	`, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols := make(map[string]bool)
	for rows.Next() {
		var col string
		if err := rows.Scan(&col); err != nil {
			return nil, err
		}
		cols[col] = true
	}
	return cols, rows.Err()
}

// safeFilters intersects the requested filter keys against actual table columns.
func (r *Reader) safeFilters(ctx context.Context, tableName string, filters map[string]any) (map[string]any, []string, error) {
	if len(filters) == 0 {
		return map[string]any{}, []string{}, nil
	}
	actual, err := r.actualColumns(ctx, tableName)
	if err != nil {
		return nil, nil, err
	}
	safe := map[string]any{}
	missing := []string{}
	for fieldname, value := range filters {
		fieldname = strings.TrimSpace(fieldname)
		if fieldname == "" {
			continue
		}
		if actual[fieldname] {
			safe[fieldname] = value
		} else {
			missing = append(missing, fieldname)
		}
	}
	return safe, missing, nil
}

// buildWhereClause produces a parameterised WHERE clause from a safe filter map.
func buildWhereClause(filters map[string]any) (string, []any) {
	if len(filters) == 0 {
		return "", []any{}
	}
	clauses := []string{}
	args := []any{}
	for fieldname, value := range filters {
		args = append(args, value)
		clauses = append(clauses, fmt.Sprintf(`%s = $%d`, quoteIdent(fieldname), len(args)))
	}
	return "WHERE " + strings.Join(clauses, " AND "), args
}

func scanRows(rows pgx.Rows, columns []string) ([]map[string]any, error) {
	data := []map[string]any{}
	for rows.Next() {
		values := make([]any, len(columns))
		ptrs := make([]any, len(columns))
		for i := range values {
			ptrs[i] = &values[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		row := make(map[string]any, len(columns))
		for i, col := range columns {
			row[col] = values[i]
		}
		data = append(data, row)
	}
	return data, rows.Err()
}
