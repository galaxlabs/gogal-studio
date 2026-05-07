package api

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) ListInstalledApps(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := h.DB.Query(ctx, `
		SELECT name, app_name, app_version
		FROM "tabInstalled App"
		ORDER BY idx, name
	`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	data := []fiber.Map{}

	for rows.Next() {
		var name, appName, appVersion string

		if err := rows.Scan(&name, &appName, &appVersion); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		data = append(data, fiber.Map{
			"name":        name,
			"app_name":    appName,
			"app_version": appVersion,
		})
	}

	return c.JSON(fiber.Map{"data": data})
}

func (h *Handler) ListModules(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := h.DB.Query(ctx, `
		SELECT name, module_name, app_name
		FROM "tabModule Def"
		ORDER BY idx, name
	`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	data := []fiber.Map{}

	for rows.Next() {
		var name, moduleName, appName string

		if err := rows.Scan(&name, &moduleName, &appName); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		data = append(data, fiber.Map{
			"name":        name,
			"module_name": moduleName,
			"app_name":    appName,
		})
	}

	return c.JSON(fiber.Map{"data": data})
}
