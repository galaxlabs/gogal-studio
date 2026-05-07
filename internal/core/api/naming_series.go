package api

import (
	"context"
	"net/url"
	"time"

	"github.com/galaxylabs/gogal-studio/internal/core/lifecycle"
	"github.com/galaxylabs/gogal-studio/internal/core/naming"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NamingSeriesHandler struct {
	DB *pgxpool.Pool
}

func NewNamingSeriesHandler(db *pgxpool.Pool) *NamingSeriesHandler {
	return &NamingSeriesHandler{DB: db}
}

func decodeSeriesKey(raw string) string {
	decoded, err := url.PathUnescape(raw)
	if err != nil {
		return raw
	}
	return decoded
}

func (h *NamingSeriesHandler) ListNamingSeries(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := h.DB.Query(ctx, `
		SELECT
			name,
			series_key,
			prefix,
			current_value,
			digits,
			description,
			owner,
			creation,
			modified,
			modified_by,
			docstatus,
			idx
		FROM "tabNaming Series"
		ORDER BY idx, name
	`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	data := []fiber.Map{}

	for rows.Next() {
		var (
			name         string
			seriesKey    string
			prefix       string
			currentValue int64
			digits       int
			description  string
			owner        string
			creation     time.Time
			modified     time.Time
			modifiedBy   string
			docstatus    int
			idx          int
		)

		if err := rows.Scan(
			&name,
			&seriesKey,
			&prefix,
			&currentValue,
			&digits,
			&description,
			&owner,
			&creation,
			&modified,
			&modifiedBy,
			&docstatus,
			&idx,
		); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		data = append(data, fiber.Map{
			"name":            name,
			"series_key":      seriesKey,
			"prefix":          prefix,
			"current_value":   currentValue,
			"digits":          digits,
			"description":     description,
			"owner":           owner,
			"creation":        creation,
			"modified":        modified,
			"modified_by":     modifiedBy,
			"docstatus":       docstatus,
			"docstatus_label": lifecycle.DocStatus(docstatus).String(),
			"idx":             idx,
		})
	}

	if err := rows.Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": data})
}

func (h *NamingSeriesHandler) GetNamingSeries(c *fiber.Ctx) error {
	seriesKey := decodeSeriesKey(c.Params("seriesKey"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		name         string
		prefix       string
		currentValue int64
		digits       int
		description  string
		owner        string
		creation     time.Time
		modified     time.Time
		modifiedBy   string
		docstatus    int
		idx          int
	)

	err := h.DB.QueryRow(ctx, `
		SELECT
			name,
			prefix,
			current_value,
			digits,
			description,
			owner,
			creation,
			modified,
			modified_by,
			docstatus,
			idx
		FROM "tabNaming Series"
		WHERE series_key = $1
	`, seriesKey).Scan(
		&name,
		&prefix,
		&currentValue,
		&digits,
		&description,
		&owner,
		&creation,
		&modified,
		&modifiedBy,
		&docstatus,
		&idx,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Naming Series not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"name":            name,
			"series_key":      seriesKey,
			"prefix":          prefix,
			"current_value":   currentValue,
			"digits":          digits,
			"description":     description,
			"owner":           owner,
			"creation":        creation,
			"modified":        modified,
			"modified_by":     modifiedBy,
			"docstatus":       docstatus,
			"docstatus_label": lifecycle.DocStatus(docstatus).String(),
			"idx":             idx,
		},
	})
}

func (h *NamingSeriesHandler) NextNamingSeries(c *fiber.Ctx) error {
	seriesKey := decodeSeriesKey(c.Params("seriesKey"))

	service := naming.NewSeriesService(h.DB)

	nextName, err := service.NextSeries(seriesKey)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"series_key": seriesKey,
			"name":       nextName,
		},
	})
}
