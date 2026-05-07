package main

import (
	"log"

	"github.com/galaxylabs/gogal-studio/internal/config"
	"github.com/galaxylabs/gogal-studio/internal/db"
	apphttp "github.com/galaxylabs/gogal-studio/internal/http"
)

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	server := apphttp.NewServer(cfg, database)

	log.Println(cfg.AppName + " running on http://127.0.0.1:" + cfg.AppPort)

	if err := server.App.Listen(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
