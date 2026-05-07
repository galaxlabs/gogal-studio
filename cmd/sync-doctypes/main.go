package main

import (
	"fmt"
	"log"

	"github.com/galaxylabs/gogal-studio/internal/config"
	"github.com/galaxylabs/gogal-studio/internal/core/doctype"
	"github.com/galaxylabs/gogal-studio/internal/db"
)

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	syncer := doctype.NewSyncer(database)

	results, err := syncer.SyncAll("modules")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DocType Sync Results")
	fmt.Println("--------------------")

	successCount := 0
	failedCount := 0

	for _, result := range results {
		fmt.Printf("[%s] %s - %s\n", result.Status, result.DocType, result.FilePath)

		if result.Message != "" {
			fmt.Printf("      %s\n", result.Message)
		}

		if result.Status == "success" {
			successCount++
		} else {
			failedCount++
		}
	}

	fmt.Println("--------------------")
	fmt.Printf("Success: %d\n", successCount)
	fmt.Printf("Failed:  %d\n", failedCount)

	if failedCount > 0 {
		log.Fatal("some doctypes failed to sync")
	}
}
