package cli

import (
	"encoding/json"
	"fmt"

	"github.com/galaxylabs/gogal-studio/internal/site"
	"github.com/spf13/cobra"
)

func NewConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Show and manage Gogal Studio configuration",
	}

	cmd.AddCommand(newConfigShowCommand())

	return cmd
}

func newConfigShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show common and current site configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			commonCfg, err := site.LoadCommonConfig()
			if err != nil {
				return err
			}

			siteCfg, err := site.LoadSiteConfig(commonCfg.DefaultSite)
			if err != nil {
				return err
			}

			fmt.Println("Gogal Studio config")
			fmt.Println("-------------------")

			fmt.Println("Common config:")
			commonJSON, err := json.MarshalIndent(commonCfg, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(commonJSON))

			fmt.Println()
			fmt.Println("Site config:")
			siteJSON, err := json.MarshalIndent(maskSiteConfig(siteCfg), "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(siteJSON))

			return nil
		},
	}
}

func maskSiteConfig(siteCfg site.SiteConfig) site.SiteConfig {
	if siteCfg.DBPassword != "" {
		siteCfg.DBPassword = maskPassword(siteCfg.DBPassword)
	}

	return siteCfg
}
