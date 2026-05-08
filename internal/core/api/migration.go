package api

import (
	"context"
	"strings"
	"time"

	"github.com/galaxylabs/gogal-studio/internal/core/migration"
	"github.com/gofiber/fiber/v2"
)

type MigrationPreviewRequest struct {
	DocType string `json:"doctype"`
}

type MigrationApplyRequest struct {
	DocType string `json:"doctype"`
	Confirm string `json:"confirm"`
}

func (h *Handler) PreviewMigration(c *fiber.Ctx) error {
	var req MigrationPreviewRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	req.DocType = strings.TrimSpace(req.DocType)

	if req.DocType == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "doctype is required",
		})
	}

	planner := migration.NewPlanner(h.DB)

	plan, err := planner.Preview(context.Background(), req.DocType)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Migration preview failed",
			"doctype": req.DocType,
			"detail":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": plan,
	})
}

func (h *Handler) MigrationLogs(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 20)

	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	rows, err := h.DB.Query(context.Background(), `
		SELECT
			doctype,
			table_name,
			operation,
			COALESCE(column_name, ''),
			COALESCE(fieldname, ''),
			COALESCE(fieldtype, ''),
			COALESCE(sql, ''),
			status,
			COALESCE(message, ''),
			applied_by,
			creation
		FROM "tabSchema Migration Log"
		ORDER BY creation DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":  "Failed to load migration logs",
			"detail": err.Error(),
		})
	}
	defer rows.Close()

	data := []fiber.Map{}

	for rows.Next() {
		var (
			doctype    string
			tableName  string
			operation  string
			columnName string
			fieldname  string
			fieldtype  string
			sqlText    string
			status     string
			message    string
			appliedBy  string
			creation   time.Time
		)

		if err := rows.Scan(
			&doctype,
			&tableName,
			&operation,
			&columnName,
			&fieldname,
			&fieldtype,
			&sqlText,
			&status,
			&message,
			&appliedBy,
			&creation,
		); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":  "Failed to scan migration log",
				"detail": err.Error(),
			})
		}

		data = append(data, fiber.Map{
			"doctype":     doctype,
			"table_name":  tableName,
			"operation":   operation,
			"column_name": columnName,
			"fieldname":   fieldname,
			"fieldtype":   fieldtype,
			"sql":         sqlText,
			"status":      status,
			"message":     message,
			"applied_by":  appliedBy,
			"creation":    creation.Format(time.RFC3339),
		})
	}

	if err := rows.Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":  "Failed to read migration logs",
			"detail": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":  data,
		"limit": limit,
	})
}

func (h *Handler) ApplyMigration(c *fiber.Ctx) error {
	var req MigrationApplyRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	req.DocType = strings.TrimSpace(req.DocType)
	req.Confirm = strings.TrimSpace(req.Confirm)

	if req.DocType == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "doctype is required",
		})
	}

	if req.Confirm != migration.ConfirmApplyToken {
		return c.Status(400).JSON(fiber.Map{
			"error":    "schema apply blocked",
			"required": migration.ConfirmApplyToken,
		})
	}

	planner := migration.NewPlanner(h.DB)

	plan, err := planner.Apply(context.Background(), req.DocType, req.Confirm)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Migration apply failed",
			"doctype": req.DocType,
			"detail":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": plan,
	})
}
