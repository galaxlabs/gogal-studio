package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/galaxylabs/gogal-studio/internal/core/meta"
	"github.com/galaxylabs/gogal-studio/internal/core/naming"
	"github.com/galaxylabs/gogal-studio/internal/core/permission"
	"github.com/galaxylabs/gogal-studio/internal/core/sysfields"
	"github.com/galaxylabs/gogal-studio/internal/core/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ListDocuments handles GET /api/core/documents/:doctype
//
// Query params:
//   - user      (required) — used for permission check
//   - limit     (optional, default 20, max 500)
//   - offset    (optional, default 0)
//   - order_by  (optional, default "name")
//   - order_dir (optional, "asc" or "desc", default "asc")
func (h *Handler) ListDocuments(c *fiber.Ctx) error {
	doctypeName := decodeCRUDParam(c.Params("doctype"))
	username := c.Query("user")

	if username == "" {
		return c.Status(400).JSON(fiber.Map{"error": "user query parameter is required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Permission gate.
	checker := permission.NewChecker(h.DB)

	canRead, err := checker.CanUserRead(ctx, username, doctypeName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if !canRead {
		return c.Status(403).JSON(fiber.Map{
			"error":   "permission denied",
			"user":    username,
			"doctype": doctypeName,
			"action":  "read",
		})
	}

	// Load meta to get table name.
	loader := meta.NewLoader(h.DB)

	m, err := loader.GetDocTypeMeta(ctx, doctypeName)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "doctype not found", "doctype": doctypeName})
	}

	if m.TableName == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "doctype has no table name configured",
			"doctype": doctypeName,
		})
	}

	// Parse pagination.
	limit := clampInt(c.QueryInt("limit", 20), 1, 500)
	offset := maxInt(c.QueryInt("offset", 0), 0)

	// Parse ordering — only allow safeIdentifier values.
	orderBy := "name"

	if raw := c.Query("order_by"); raw != "" && isSafeIdentifier(raw) {
		orderBy = raw
	}

	orderDir := "ASC"

	if strings.EqualFold(c.Query("order_dir"), "desc") {
		orderDir = "DESC"
	}

	query := fmt.Sprintf(
		`SELECT * FROM %s ORDER BY %s %s LIMIT $1 OFFSET $2`,
		quoteTableIdent(m.TableName),
		quoteTableIdent(orderBy),
		orderDir,
	)

	rows, err := h.DB.Query(ctx, query, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	results, err := rowsToMaps(rows)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Count total.
	var total int

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, quoteTableIdent(m.TableName))

	if err := h.DB.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		total = -1
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"doctype": doctypeName,
			"total":   total,
			"limit":   limit,
			"offset":  offset,
			"results": results,
		},
	})
}

// GetDocument handles GET /api/core/documents/:doctype/:name
//
// Query params:
//   - user (required) — used for permission check
func (h *Handler) GetDocument(c *fiber.Ctx) error {
	doctypeName := decodeCRUDParam(c.Params("doctype"))
	docName := decodeCRUDParam(c.Params("name"))
	username := c.Query("user")

	if username == "" {
		return c.Status(400).JSON(fiber.Map{"error": "user query parameter is required"})
	}

	if docName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "document name is required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Permission gate.
	checker := permission.NewChecker(h.DB)

	canRead, err := checker.CanUserRead(ctx, username, doctypeName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if !canRead {
		return c.Status(403).JSON(fiber.Map{
			"error":   "permission denied",
			"user":    username,
			"doctype": doctypeName,
			"action":  "read",
		})
	}

	// Load meta.
	loader := meta.NewLoader(h.DB)

	m, err := loader.GetDocTypeMeta(ctx, doctypeName)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "doctype not found", "doctype": doctypeName})
	}

	if m.TableName == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "doctype has no table name configured",
			"doctype": doctypeName,
		})
	}

	query := fmt.Sprintf(
		`SELECT * FROM %s WHERE name = $1 LIMIT 1`,
		quoteTableIdent(m.TableName),
	)

	rows, err := h.DB.Query(ctx, query, docName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	results, err := rowsToMaps(rows)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if len(results) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error":   "document not found",
			"doctype": doctypeName,
			"name":    docName,
		})
	}

	return c.JSON(fiber.Map{
		"data": results[0],
	})
}

// rowsToMaps converts pgx rows to a slice of map[string]any.
func rowsToMaps(rows pgx.Rows) ([]map[string]any, error) {
	descriptions := rows.FieldDescriptions()
	results := []map[string]any{}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		row := make(map[string]any, len(descriptions))

		for i, col := range descriptions {
			row[col.Name] = values[i]
		}

		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// decodeCRUDParam URL-decodes a Fiber route parameter.
func decodeCRUDParam(raw string) string {
	decoded, err := url.PathUnescape(raw)
	if err != nil {
		return raw
	}

	return decoded
}

// isSafeIdentifier allows only alphanumeric + underscore — safe for ORDER BY column names.
func isSafeIdentifier(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}

	return true
}

// quoteTableIdent quotes a table or column identifier safely.
func quoteTableIdent(s string) string {
	safe := strings.ReplaceAll(s, `"`, `""`)
	return `"` + safe + `"`
}

func clampInt(v, min, max int) int {
	if v < min {
		return min
	}

	if v > max {
		return max
	}

	return v
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// Prevent "declared and not used" if strconv is not otherwise referenced.
var _ = strconv.Itoa

// ─── CRUD Write ───────────────────────────────────────────────────────────────

// CreateDocument handles POST /api/core/documents/:doctype
//
// Body: JSON object with the fields to set.
// Query params:
//   - user (required) — owner + permission check
func (h *Handler) CreateDocument(c *fiber.Ctx) error {
	doctypeName := decodeCRUDParam(c.Params("doctype"))
	username := c.Query("user")

	if username == "" {
		return c.Status(400).JSON(fiber.Map{"error": "user query parameter is required"})
	}

	var doc map[string]any
	if err := c.BodyParser(&doc); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body: " + err.Error()})
	}
	if doc == nil {
		doc = map[string]any{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Permission gate.
	checker := permission.NewChecker(h.DB)
	canCreate, err := checker.CanUserCreate(ctx, username, doctypeName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if !canCreate {
		return c.Status(403).JSON(fiber.Map{
			"error":   "permission denied",
			"user":    username,
			"doctype": doctypeName,
			"action":  "create",
		})
	}

	// Load meta.
	loader := meta.NewLoader(h.DB)
	m, err := loader.GetDocTypeMeta(ctx, doctypeName)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "doctype not found", "doctype": doctypeName})
	}
	if m.TableName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "doctype has no table name configured", "doctype": doctypeName})
	}

	// Validate user-supplied fields.
	if err := validator.Validate(m, doc); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Resolve document name.
	docName, err := resolveName(ctx, h.DB, m, doc)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "naming error: " + err.Error()})
	}

	// Inject system fields.
	sysfields.InjectCreate(doc, docName, username)

	// Build allowed column set and execute INSERT.
	allowed := buildAllowedColumns(m)
	cols := []string{}
	placeholders := []string{}
	vals := []any{}
	i := 1
	for col, val := range doc {
		if !allowed[col] {
			continue
		}
		cols = append(cols, quoteTableIdent(col))
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		vals = append(vals, val)
		i++
	}

	if len(cols) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "no valid fields to insert"})
	}

	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s) RETURNING *`,
		quoteTableIdent(m.TableName),
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
	)

	rows, err := h.DB.Query(ctx, query, vals...)
	if err != nil {
		if isMissingResourceTableError(err) {
			return missingResourceTable(c, doctypeName, err.Error())
		}

		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	results, err := rowsToMaps(rows)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if len(results) == 0 {
		return c.Status(500).JSON(fiber.Map{"error": "insert returned no rows"})
	}

	return c.Status(201).JSON(fiber.Map{"data": results[0]})
}

// UpdateDocument handles PUT /api/core/documents/:doctype/:name
//
// Body: JSON object with the fields to update.
// Query params:
//   - user (required) — permission check
func (h *Handler) UpdateDocument(c *fiber.Ctx) error {
	doctypeName := decodeCRUDParam(c.Params("doctype"))
	docName := decodeCRUDParam(c.Params("name"))
	username := c.Query("user")

	if username == "" {
		return c.Status(400).JSON(fiber.Map{"error": "user query parameter is required"})
	}
	if docName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "document name is required"})
	}

	var updates map[string]any
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body: " + err.Error()})
	}
	if updates == nil {
		updates = map[string]any{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Permission gate.
	checker := permission.NewChecker(h.DB)
	canWrite, err := checker.CanUserWrite(ctx, username, doctypeName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if !canWrite {
		return c.Status(403).JSON(fiber.Map{
			"error":   "permission denied",
			"user":    username,
			"doctype": doctypeName,
			"action":  "write",
		})
	}

	// Load meta.
	loader := meta.NewLoader(h.DB)
	m, err := loader.GetDocTypeMeta(ctx, doctypeName)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "doctype not found", "doctype": doctypeName})
	}
	if m.TableName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "doctype has no table name configured", "doctype": doctypeName})
	}

	// Remove protected fields from user input — these must not change.
	for _, f := range sysfields.ProtectedFields {
		delete(updates, f)
	}

	// Validate user-supplied fields.
	if err := validator.Validate(m, updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Inject modified/modified_by.
	sysfields.InjectUpdate(updates, username)

	// Build allowed column set and execute UPDATE.
	allowed := buildAllowedColumns(m)
	setClauses := []string{}
	vals := []any{}
	i := 1
	for col, val := range updates {
		if !allowed[col] {
			continue
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", quoteTableIdent(col), i))
		vals = append(vals, val)
		i++
	}

	if len(setClauses) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "no valid fields to update"})
	}

	vals = append(vals, docName) // WHERE name = $i
	query := fmt.Sprintf(
		`UPDATE %s SET %s WHERE name = $%d RETURNING *`,
		quoteTableIdent(m.TableName),
		strings.Join(setClauses, ", "),
		i,
	)

	rows, err := h.DB.Query(ctx, query, vals...)
	if err != nil {
		if isMissingResourceTableError(err) {
			return missingResourceTable(c, doctypeName, err.Error())
		}

		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	results, err := rowsToMaps(rows)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if len(results) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error":   "document not found",
			"doctype": doctypeName,
			"name":    docName,
		})
	}

	return c.JSON(fiber.Map{"data": results[0]})
}

// DeleteDocument handles DELETE /api/core/documents/:doctype/:name
//
// Query params:
//   - user (required) — permission check
func (h *Handler) DeleteDocument(c *fiber.Ctx) error {
	doctypeName := decodeCRUDParam(c.Params("doctype"))
	docName := decodeCRUDParam(c.Params("name"))
	username := c.Query("user")

	if username == "" {
		return c.Status(400).JSON(fiber.Map{"error": "user query parameter is required"})
	}
	if docName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "document name is required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Permission gate.
	checker := permission.NewChecker(h.DB)
	canDelete, err := checker.CanUserDelete(ctx, username, doctypeName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if !canDelete {
		return c.Status(403).JSON(fiber.Map{
			"error":   "permission denied",
			"user":    username,
			"doctype": doctypeName,
			"action":  "delete",
		})
	}

	// Load meta.
	loader := meta.NewLoader(h.DB)
	m, err := loader.GetDocTypeMeta(ctx, doctypeName)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "doctype not found", "doctype": doctypeName})
	}
	if m.TableName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "doctype has no table name configured", "doctype": doctypeName})
	}

	query := fmt.Sprintf(`DELETE FROM %s WHERE name = $1`, quoteTableIdent(m.TableName))

	tag, err := h.DB.Exec(ctx, query, docName)
	if err != nil {
		if isMissingResourceTableError(err) {
			return missingResourceTable(c, doctypeName, err.Error())
		}

		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if tag.RowsAffected() == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error":   "document not found",
			"doctype": doctypeName,
			"name":    docName,
		})
	}

	return c.JSON(fiber.Map{
		"message": "deleted",
		"doctype": doctypeName,
		"name":    docName,
	})
}

// ─── Naming helpers ───────────────────────────────────────────────────────────

// resolveName determines the document name for a new record.
// Priority: (1) user-provided "name" in body, (2) naming_rule, (3) autoname pattern, (4) auto-generated.
func resolveName(ctx context.Context, db *pgxpool.Pool, m meta.DocTypeMeta, doc map[string]any) (string, error) {
	// 1. User-provided name.
	if nameVal, ok := doc["name"]; ok && nameVal != nil {
		if s, ok := nameVal.(string); ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s), nil
		}
	}

	// 2. naming_rule: series: or field: prefix.
	rule := strings.TrimSpace(m.NamingRule)
	if strings.HasPrefix(rule, "series:") || strings.HasPrefix(rule, "field:") {
		svc := naming.NewSeriesService(db)
		return svc.GenerateName(m.Name, rule, naming.Document(doc))
	}

	// 3. autoname column holds a series pattern (e.g. "CUST-.####").
	autoname := strings.TrimSpace(m.Autoname)
	if autoname != "" && autoname != "autoname" && autoname != "hash" && autoname != "manual" {
		svc := naming.NewSeriesService(db)
		seriesRule := autoname
		if !strings.HasPrefix(autoname, "series:") && !strings.HasPrefix(autoname, "field:") {
			seriesRule = "series:" + autoname
		}
		return svc.GenerateName(m.Name, seriesRule, naming.Document(doc))
	}

	// 4. Fallback: abbreviation + nanosecond timestamp hex.
	return abbreviate(m.Name) + fmt.Sprintf("-%X", time.Now().UnixNano()), nil
}

// abbreviate returns an uppercase abbreviation from the first letter of each word,
// falling back to the first 4 characters of the full name.
func abbreviate(s string) string {
	parts := strings.Fields(s)
	result := ""
	for _, p := range parts {
		if len(p) > 0 {
			result += strings.ToUpper(string(p[0]))
		}
	}
	if result == "" && len(s) > 0 {
		end := len(s)
		if end > 4 {
			end = 4
		}
		result = strings.ToUpper(s[:end])
	}
	return result
}

// buildAllowedColumns returns the set of column names that may appear in INSERT/UPDATE.
// It includes all stored meta fields plus the standard system fields.
func buildAllowedColumns(m meta.DocTypeMeta) map[string]bool {
	allowed := map[string]bool{
		"name": true, "owner": true, "creation": true,
		"modified": true, "modified_by": true, "docstatus": true, "idx": true,
		"parent": true, "parentfield": true, "parenttype": true,
	}
	for _, f := range m.Fields {
		if validator.IsStoredFieldtype(f.Fieldtype) {
			allowed[f.Fieldname] = true
		}
	}
	return allowed
}
