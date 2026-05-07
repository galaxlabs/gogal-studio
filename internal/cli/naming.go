package cli

import (
	"fmt"

	"github.com/galaxylabs/gogal-studio/internal/core/naming"
	"github.com/galaxylabs/gogal-studio/internal/db"
	"github.com/galaxylabs/gogal-studio/internal/site"
	"github.com/spf13/cobra"
)

func NewNamingCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "naming",
		Short: "Naming series utilities",
	}

	cmd.AddCommand(newNextSeriesCommand())
	cmd.AddCommand(newGenerateNameCommand())

	return cmd
}

func newNextSeriesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "next [series_key]",
		Short: "Generate next name from a naming series",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			commonCfg, err := site.LoadCommonConfig()
			if err != nil {
				return err
			}

			siteCfg, err := site.LoadSiteConfig(commonCfg.DefaultSite)
			if err != nil {
				return err
			}

			database, err := db.Connect(siteCfg.DatabaseURL())
			if err != nil {
				return err
			}
			defer database.Close()

			service := naming.NewSeriesService(database)

			nextName, err := service.NextSeries(args[0])
			if err != nil {
				return err
			}

			fmt.Println(nextName)
			return nil
		},
	}
}

func newGenerateNameCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "generate [doctype] [rule]",
		Short: "Generate name using a naming rule",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			doctype := args[0]
			rule := args[1]

			commonCfg, err := site.LoadCommonConfig()
			if err != nil {
				return err
			}

			siteCfg, err := site.LoadSiteConfig(commonCfg.DefaultSite)
			if err != nil {
				return err
			}

			database, err := db.Connect(siteCfg.DatabaseURL())
			if err != nil {
				return err
			}
			defer database.Close()

			service := naming.NewSeriesService(database)

			name, err := service.GenerateName(doctype, rule, naming.Document{
				"name": doctype,
			})
			if err != nil {
				return err
			}

			fmt.Println(name)
			return nil
		},
	}
}
