package http

import (
	"context"
	"time"

	"github.com/galaxylabs/gogal-studio/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(app *fiber.App, cfg config.Config, database *pgxpool.Pool) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"app":     "gogal-studio",
			"name":    cfg.AppName,
			"message": "API server is running",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		dbStatus := "ok"

		if err := database.Ping(ctx); err != nil {
			dbStatus = "error"
		}

		return c.JSON(fiber.Map{
			"status":   "ok",
			"app":      "gogal-studio",
			"env":      cfg.AppEnv,
			"database": dbStatus,
		})
	})

	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Gogal Studio API root",
		})
	})

	api.Get("/version", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name":    "gogal-studio",
			"version": "0.0.1",
			"stage":   "foundation",
			"engine":  "Gogal Engine",
			"company": "Galaxy Labs",
		})
	})
}
