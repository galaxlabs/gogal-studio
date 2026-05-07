package site

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

type CommonSiteConfig struct {
	DefaultSite   string `json:"default_site"`
	ServerPort    int    `json:"server_port"`
	DeveloperMode bool   `json:"developer_mode"`
}

type SiteConfig struct {
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

func CommonConfigPath() string {
	return filepath.Join("sites", "common_site_config.json")
}

func SiteConfigPath(siteName string) string {
	return filepath.Join("sites", siteName, "site_config.json")
}

func LoadCommonConfig() (CommonSiteConfig, error) {
	var cfg CommonSiteConfig

	raw, err := os.ReadFile(CommonConfigPath())
	if err != nil {
		return cfg, fmt.Errorf("failed to read common site config: %w", err)
	}

	if err := json.Unmarshal(raw, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse common site config: %w", err)
	}

	return cfg, nil
}

func LoadSiteConfig(siteName string) (SiteConfig, error) {
	var cfg SiteConfig

	raw, err := os.ReadFile(SiteConfigPath(siteName))
	if err != nil {
		return cfg, fmt.Errorf("failed to read site config for %s: %w", siteName, err)
	}

	if err := json.Unmarshal(raw, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse site config for %s: %w", siteName, err)
	}

	return cfg, nil
}
func (cfg SiteConfig) DatabaseURL() string {
	password := url.QueryEscape(cfg.DBPassword)

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser,
		password,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
}
func SaveSiteConfig(cfg SiteConfig) error {
	path := SiteConfigPath(cfg.SiteName)

	raw, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode site config: %w", err)
	}

	raw = append(raw, '\n')

	if err := os.WriteFile(path, raw, 0644); err != nil {
		return fmt.Errorf("failed to write site config: %w", err)
	}

	return nil
}
