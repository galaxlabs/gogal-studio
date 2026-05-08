package api

import (
	"context"

	"github.com/galaxylabs/gogal-studio/internal/core/permission"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PermissionHandler struct {
	DB *pgxpool.Pool
}

func NewPermissionHandler(db *pgxpool.Pool) *PermissionHandler {
	return &PermissionHandler{DB: db}
}

// CheckPermission handles GET /api/core/permissions/check
//
// Query params:
//   - user     (required)
//   - doctype  (required)
//   - action   (optional) — if omitted, returns read/create/write/delete summary
func (h *PermissionHandler) CheckPermission(c *fiber.Ctx) error {
	username := c.Query("user")
	doctypeName := c.Query("doctype")
	action := permission.Action(c.Query("action"))

	if username == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "user query parameter is required",
		})
	}

	if doctypeName == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "doctype query parameter is required",
		})
	}

	checker := permission.NewChecker(h.DB)
	ctx := context.Background()

	roles, err := permission.GetRolesForUser(ctx, h.DB, username)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Single-action mode.
	if action != "" {
		allowed, err := checker.Can(ctx, doctypeName, roles, action)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"data": fiber.Map{
				"user":    username,
				"roles":   roles,
				"doctype": doctypeName,
				"action":  action,
				"allowed": allowed,
			},
		})
	}

	// Summary mode — return all four base permissions.
	canRead, err := checker.CanUserRead(ctx, username, doctypeName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	canCreate, err := checker.CanUserCreate(ctx, username, doctypeName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	canWrite, err := checker.CanUserWrite(ctx, username, doctypeName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	canDelete, err := checker.CanUserDelete(ctx, username, doctypeName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"user":    username,
			"roles":   roles,
			"doctype": doctypeName,
			"read":    canRead,
			"create":  canCreate,
			"write":   canWrite,
			"delete":  canDelete,
		},
	})
}
