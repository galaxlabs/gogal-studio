package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/galaxylabs/gogal-studio/internal/site"
	"github.com/spf13/cobra"
)

func NewSiteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site",
		Short: "Site operations",
	}

	cmd.AddCommand(newSiteInfoCommand())

	return cmd
}

func newSiteInfoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Show current default site information",
		RunE: func(cmd *cobra.Command, args []string) error {
			commonCfg, err := site.LoadCommonConfig()
			if err != nil {
				return err
			}

			siteCfg, err := site.LoadSiteConfig(commonCfg.DefaultSite)
			if err != nil {
				return err
			}

			commonPath := filepath.Join("sites", "common_site_config.json")
			sitePath := filepath.Join("sites", commonCfg.DefaultSite, "site_config.json")

			fmt.Println("Gogal Studio site info")
			fmt.Println("----------------------")
			fmt.Println("Default site:       ", commonCfg.DefaultSite)
			fmt.Println("Developer mode:     ", commonCfg.DeveloperMode)
			fmt.Println("Server port:        ", commonCfg.ServerPort)
			fmt.Println("Common config path: ", commonPath)
			fmt.Println("Site config path:   ", sitePath)
			fmt.Println()

			fmt.Println("Database")
			fmt.Println("--------")
			fmt.Println("DB type:            ", siteCfg.DBType)
			fmt.Println("DB host:            ", siteCfg.DBHost)
			fmt.Println("DB port:            ", siteCfg.DBPort)
			fmt.Println("DB name:            ", siteCfg.DBName)
			fmt.Println("DB user:            ", siteCfg.DBUser)
			fmt.Println("DB password:        ", maskPassword(siteCfg.DBPassword))
			fmt.Println()

			fmt.Println("Installed")
			fmt.Println("---------")
			fmt.Println("Installed apps:    ", siteCfg.InstalledApps)
			fmt.Println("Installed modules: ", siteCfg.InstalledModules)
			fmt.Println()

			fmt.Println("Files")
			fmt.Println("-----")
			printFileStatus(commonPath)
			printFileStatus(sitePath)

			return nil
		},
	}
}

func printFileStatus(path string) {
	if _, err := os.Stat(path); err != nil {
		fmt.Println(path + ": missing")
		return
	}

	fmt.Println(path + ": exists")
}
