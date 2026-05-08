package api

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/galaxylabs/gogal-studio/internal/core/crud"
	"github.com/gofiber/fiber/v2"
)

// ListResources handles GET /api/resource/:doctype
//
// Query params:
//   - user  (optional, default "Administrator")
//   - limit (optional, default 20, max 500)
func (h *Handler) ListResources(c *fiber.Ctx) error {
	doctype := decodeCRUDParam(c.Params("doctype"))
	user := strings.TrimSpace(c.Query("user", "Administrator"))
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)
	fields := parseFields(c.Query("fields"))
	filters, err := parseFilters(c.Query("filters"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":  "Invalid filters JSON",
			"detail": err.Error(),
		})
	}

	reader := crud.NewReader(h.DB)

	result, err := reader.List(context.Background(), doctype, crud.ListOptions{
		User:    user,
		Limit:   limit,
		Offset:  offset,
		Fields:  fields,
		Filters: filters,
	})
	if err != nil {
		if isMissingResourceTableError(err) {
			return missingResourceTable(c, doctype, err.Error())
		}

		status := 500
		if strings.Contains(err.Error(), "permission denied") {
			status = 403
		}
		if strings.Contains(err.Error(), "doctype not found") {
			status = 404
		}
		return c.Status(status).JSON(fiber.Map{
			"error":   err.Error(),
			"doctype": doctype,
		})
	}

	return c.JSON(fiber.Map{
		"doctype":                doctype,
		"limit":                  limit,
		"offset":                 offset,
		"data":                   result.Data,
		"columns":                result.Columns,
		"missing_columns":        result.MissingColumns,
		"missing_filter_columns": result.MissingFilterColumns,
	})
}

func parseFields(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{}
	}
	fields := []string{}
	for _, part := range strings.Split(raw, ",") {
		if f := strings.TrimSpace(part); f != "" {
			fields = append(fields, f)
		}
	}
	return fields
}

func parseFilters(raw string) (map[string]any, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return map[string]any{}, nil
	}
	if decoded, err := url.QueryUnescape(raw); err == nil {
		raw = decoded
	}
	filters := map[string]any{}
	if err := json.Unmarshal([]byte(raw), &filters); err != nil {
		return nil, err
	}
	return filters, nil
}

// GetResource handles GET /api/resource/:doctype/:name
//
// Query params:
//   - user (optional, default "Administrator")
func (h *Handler) GetResource(c *fiber.Ctx) error {
	doctype := decodeCRUDParam(c.Params("doctype"))
	name := decodeCRUDParam(c.Params("name"))
	user := strings.TrimSpace(c.Query("user", "Administrator"))

	reader := crud.NewReader(h.DB)

	result, err := reader.Get(context.Background(), doctype, name, user)
	if err != nil {
		if isMissingResourceTableError(err) {
			return missingResourceTable(c, doctype, err.Error())
		}

		status := 500
		if strings.Contains(err.Error(), "document not found") {
			status = 404
		}
		if strings.Contains(err.Error(), "permission denied") {
			status = 403
		}
		return c.Status(status).JSON(fiber.Map{
			"error":   err.Error(),
			"doctype": doctype,
			"name":    name,
		})
	}

	return c.JSON(fiber.Map{
		"doctype":         doctype,
		"name":            name,
		"data":            result.Data,
		"columns":         result.Columns,
		"missing_columns": result.MissingColumns,
	})
}
