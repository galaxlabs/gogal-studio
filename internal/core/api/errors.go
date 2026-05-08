package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func apiError(c *fiber.Ctx, status int, code string, message string, detail string) error {
	body := fiber.Map{
		"error": fiber.Map{
			"code":    code,
			"message": message,
		},
	}

	if detail != "" {
		body["error"].(fiber.Map)["detail"] = detail
	}

	return c.Status(status).JSON(body)
}

func badRequest(c *fiber.Ctx, message string, detail string) error {
	return apiError(c, fiber.StatusBadRequest, "bad_request", message, detail)
}

func forbidden(c *fiber.Ctx, message string, detail string) error {
	return apiError(c, fiber.StatusForbidden, "forbidden", message, detail)
}

func notFound(c *fiber.Ctx, message string, detail string) error {
	return apiError(c, fiber.StatusNotFound, "not_found", message, detail)
}

func serverError(c *fiber.Ctx, message string, detail string) error {
	return apiError(c, fiber.StatusInternalServerError, "server_error", message, detail)
}

func isMissingResourceTableError(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())

	return strings.Contains(msg, "no readable database columns found") ||
		(strings.Contains(msg, "relation ") && strings.Contains(msg, " does not exist"))
}

func missingResourceTable(c *fiber.Ctx, doctype string, detail string) error {
	return c.Status(fiber.StatusConflict).JSON(fiber.Map{
		"error":   "record table is not available",
		"code":    "table_missing",
		"doctype": doctype,
		"detail":  detail,
	})
}
