package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/galaxylabs/gogal-studio/internal/db"
	"github.com/spf13/cobra"
)

type doctorCommonConfig struct {
	DefaultSite   string `json:"default_site"`
	ServerPort    int    `json:"server_port"`
	DeveloperMode bool   `json:"developer_mode"`
}

type doctorSiteConfig struct {
	SiteName         string   `json:"site_name"`
	DBType           string   `json:"db_type"`
	DBHost           string   `json:"db_host"`
	DBPort           int      `json:"db_port"`
	DBName           string   `json:"db_name"`
	DBUser           string   `json:"db_user"`
	DBPassword       string   `json:"db_password"`
	InstalledApps    []string `json:"installed_apps"`
	InstalledModules []string `json:"installed_modules"`
}

func NewDoctorCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check Gogal Studio installation health",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Gogal Studio doctor")
			fmt.Println("-------------------")

			commonCfg, err := readDoctorCommonConfig()
			if err != nil {
				return err
			}

			siteCfg, err := readDoctorSiteConfig(commonCfg.DefaultSite)
			if err != nil {
				return err
			}

			fmt.Println("Default site:", commonCfg.DefaultSite)
			fmt.Println("Database:", siteCfg.DBType, siteCfg.DBHost, siteCfg.DBPort, siteCfg.DBName)
			fmt.Println("DB user:", siteCfg.DBUser)

			database, err := db.Connect(siteCfg.DatabaseURL())
			if err != nil {
				return err
			}
			defer database.Close()

			fmt.Println("PostgreSQL connection: OK")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			checks := []struct {
				Label string
				SQL   string
			}{
				{"Installed apps", `SELECT COUNT(*) FROM "tabInstalled App"`},
				{"Installed modules", `SELECT COUNT(*) FROM "tabModule Def"`},
				{"Users", `SELECT COUNT(*) FROM "tabUser"`},
				{"Roles", `SELECT COUNT(*) FROM "tabRole"`},
				{"DocTypes", `SELECT COUNT(*) FROM "tabDocType"`},
				{"DocFields", `SELECT COUNT(*) FROM "tabDocField"`},
				{"DocPerms", `SELECT COUNT(*) FROM "tabDocPerm"`},
			}

			for _, check := range checks {
				var count int

				if err := database.QueryRow(ctx, check.SQL).Scan(&count); err != nil {
					return fmt.Errorf("%s check failed: %w", check.Label, err)
				}

				fmt.Printf("%s: %d\n", check.Label, count)
			}

			fmt.Println("-------------------")
			fmt.Println("Gogal Studio installation: OK")

			return nil
		},
	}
}

func readDoctorCommonConfig() (doctorCommonConfig, error) {
	var cfg doctorCommonConfig

	raw, err := os.ReadFile("sites/common_site_config.json")
	if err != nil {
		return cfg, fmt.Errorf("read common site config: %w", err)
	}

	if err := json.Unmarshal(raw, &cfg); err != nil {
		return cfg, fmt.Errorf("parse common site config: %w", err)
	}

	if cfg.DefaultSite == "" {
		return cfg, fmt.Errorf("default_site is missing in common_site_config.json")
	}

	return cfg, nil
}

func readDoctorSiteConfig(siteName string) (doctorSiteConfig, error) {
	var cfg doctorSiteConfig

	path := "sites/" + siteName + "/site_config.json"

	raw, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("read site config: %w", err)
	}

	if err := json.Unmarshal(raw, &cfg); err != nil {
		return cfg, fmt.Errorf("parse site config: %w", err)
	}

	return cfg, nil
}

func (cfg doctorSiteConfig) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		url.QueryEscape(cfg.DBUser),
		url.QueryEscape(cfg.DBPassword),
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
}
