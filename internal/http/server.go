package http

import (
	"github.com/galaxylabs/gogal-studio/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	App *fiber.App
	Cfg config.Config
	DB  *pgxpool.Pool
}

func NewServer(cfg config.Config, database *pgxpool.Pool) *Server {
	app := fiber.New(fiber.Config{
		AppName: cfg.AppName,
	})

	server := &Server{
		App: app,
		Cfg: cfg,
		DB:  database,
	}

	RegisterRoutes(app, cfg, database)

	return server
}
