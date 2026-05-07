package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/galaxylabs/gogal-studio/internal/config"
	"github.com/galaxylabs/gogal-studio/internal/db"
	apphttp "github.com/galaxylabs/gogal-studio/internal/http"
	"github.com/spf13/cobra"
)

type startCommonSiteConfig struct {
	DefaultSite string `json:"default_site"`
	ServerPort  int    `json:"server_port"`
}

type startSiteConfig struct {
	SiteName   string `json:"site_name"`
	DBType     string `json:"db_type"`
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DBName     string `json:"db_name"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
}

func (c startSiteConfig) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		url.QueryEscape(c.DBUser),
		url.QueryEscape(c.DBPassword),
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

func NewStartCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start Gogal Studio server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Load()

			commonCfg, siteCfg, err := loadStartSiteConfig()
			if err != nil {
				return err
			}

			database, err := db.Connect(siteCfg.DatabaseURL())
			if err != nil {
				return err
			}
			defer database.Close()

			server := apphttp.NewServer(cfg, database)

			log.Println("Default site:", commonCfg.DefaultSite)
			log.Println("Database:", siteCfg.DBName)
			log.Println(cfg.AppName + " running on http://127.0.0.1:" + cfg.AppPort)

			return server.App.Listen(":" + cfg.AppPort)
		},
	}
}

func loadStartSiteConfig() (startCommonSiteConfig, startSiteConfig, error) {
	var commonCfg startCommonSiteConfig
	var siteCfg startSiteConfig

	commonPath := filepath.Join("sites", "common_site_config.json")

	commonRaw, err := os.ReadFile(commonPath)
	if err != nil {
		return commonCfg, siteCfg, fmt.Errorf("read common site config: %w", err)
	}

	if err := json.Unmarshal(commonRaw, &commonCfg); err != nil {
		return commonCfg, siteCfg, fmt.Errorf("parse common site config: %w", err)
	}

	if commonCfg.DefaultSite == "" {
		commonCfg.DefaultSite = "gogal.dev"
	}

	sitePath := filepath.Join("sites", commonCfg.DefaultSite, "site_config.json")

	siteRaw, err := os.ReadFile(sitePath)
	if err != nil {
		return commonCfg, siteCfg, fmt.Errorf("read site config: %w", err)
	}

	if err := json.Unmarshal(siteRaw, &siteCfg); err != nil {
		return commonCfg, siteCfg, fmt.Errorf("parse site config: %w", err)
	}

	return commonCfg, siteCfg, nil
}
