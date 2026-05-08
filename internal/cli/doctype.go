package cli

import (
	"fmt"

	coredoctype "github.com/galaxylabs/gogal-studio/internal/core/doctype"
	"github.com/galaxylabs/gogal-studio/internal/db"
	"github.com/galaxylabs/gogal-studio/internal/site"
	"github.com/spf13/cobra"
)

func NewDocTypeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctype",
		Short: "DocType utilities",
	}

	cmd.AddCommand(newWriteSampleDocTypeCommand())
	cmd.AddCommand(newReadDocTypeCommand())
	cmd.AddCommand(newReadAllDocTypesCommand())
	cmd.AddCommand(newExportDocTypeCommand())
	cmd.AddCommand(newExportAllDocTypesCommand())
	cmd.AddCommand(newImportDocTypeCommand())
	cmd.AddCommand(newImportAllDocTypesCommand())

	return cmd
}

func newWriteSampleDocTypeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "write-sample",
		Short: "Write a sample DocType JSON file",
		RunE: func(cmd *cobra.Command, args []string) error {
			doc := coredoctype.JSONDocType{
				Name:       "Sample DocType",
				Label:      "Sample DocType",
				Module:     "Core",
				AppName:    "gogal_studio",
				Autoname:   "field:name",
				NamingRule: "By fieldname",
				TitleField: "name",
				SortField:  "idx",
				SortOrder:  "ASC",
				Fields: []coredoctype.JSONDocField{
					{
						Fieldname:  "sample_title",
						Label:      "Sample Title",
						Fieldtype:  "Data",
						Reqd:       true,
						InListView: true,
						Columns:    6,
						Length:     140,
						Idx:        1,
					},
					{
						Fieldname: "description",
						Label:     "Description",
						Fieldtype: "Small Text",
						Columns:   12,
						Idx:       2,
					},
				},
				Permissions: []coredoctype.JSONDocPerm{
					{
						Role:      "System Manager",
						Permlevel: 0,
						Read:      true,
						Write:     true,
						Create:    true,
						Delete:    true,
						Idx:       1,
					},
				},
			}

			result, err := coredoctype.WriteDocTypeJSON(".", doc)
			if err != nil {
				return err
			}

			fmt.Println("DocType JSON written")
			fmt.Println("--------------------")
			fmt.Println("Folder:", result.FolderPath)
			fmt.Println("JSON:  ", result.JSONPath)

			return nil
		},
	}
}

func newReadDocTypeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "read [module] [doctype]",
		Short: "Read one DocType JSON file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := coredoctype.ReadDocTypeJSONByName(".", args[0], args[1])
			if err != nil {
				return err
			}

			fmt.Println("DocType JSON")
			fmt.Println("------------")
			fmt.Println("Path:      ", result.JSONPath)
			fmt.Println("Name:      ", result.DocType.Name)
			fmt.Println("Label:     ", result.DocType.Label)
			fmt.Println("Module:    ", result.DocType.Module)
			fmt.Println("App:       ", result.DocType.AppName)
			fmt.Println("Table:     ", result.DocType.TableName)
			fmt.Println("Fields:    ", len(result.DocType.Fields))
			fmt.Println("Perms:     ", len(result.DocType.Permissions))

			return nil
		},
	}
}

func newReadAllDocTypesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "read-all",
		Short: "Read all DocType JSON files",
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := coredoctype.ReadAllDocTypeJSON(".")
			if err != nil {
				return err
			}

			fmt.Println("DocType JSON Files")
			fmt.Println("------------------")

			for _, result := range results {
				fmt.Printf(
					"%s | %s | %s | fields=%d | perms=%d\n",
					result.DocType.Name,
					result.DocType.Module,
					result.JSONPath,
					len(result.DocType.Fields),
					len(result.DocType.Permissions),
				)
			}

			return nil
		},
	}
}

func newExportDocTypeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "export [doctype]",
		Short: "Export one DocType from database metadata to JSON file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			commonCfg, err := site.LoadCommonConfig()
			if err != nil {
				return err
			}

			siteName := commonCfg.DefaultSite
			if siteName == "" {
				siteName = "gogal.dev"
			}

			siteCfg, err := site.LoadSiteConfig(siteName)
			if err != nil {
				return err
			}

			database, err := db.Connect(siteCfg.DatabaseURL())
			if err != nil {
				return err
			}
			defer database.Close()

			exporter := coredoctype.NewExporter(database)

			result, err := exporter.ExportOneToFile(".", args[0])
			if err != nil {
				return err
			}

			fmt.Println("DocType exported")
			fmt.Println("----------------")
			fmt.Println("Folder:", result.FolderPath)
			fmt.Println("JSON:  ", result.JSONPath)

			return nil
		},
	}
}

func newExportAllDocTypesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "export-all",
		Short: "Export all DocTypes from database metadata to JSON files",
		RunE: func(cmd *cobra.Command, args []string) error {
			commonCfg, err := site.LoadCommonConfig()
			if err != nil {
				return err
			}

			siteName := commonCfg.DefaultSite
			if siteName == "" {
				siteName = "gogal.dev"
			}

			siteCfg, err := site.LoadSiteConfig(siteName)
			if err != nil {
				return err
			}

			database, err := db.Connect(siteCfg.DatabaseURL())
			if err != nil {
				return err
			}
			defer database.Close()

			exporter := coredoctype.NewExporter(database)

			results, err := exporter.ExportAllToFiles(".")
			if err != nil {
				return err
			}

			fmt.Println("DocTypes exported")
			fmt.Println("-----------------")

			for _, result := range results {
				fmt.Println(result.JSONPath)
			}

			fmt.Println()
			fmt.Println("Total:", len(results))

			return nil
		},
	}
}

func newImportDocTypeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "import [module] [doctype]",
		Short: "Import one DocType from JSON file into the database",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			commonCfg, err := site.LoadCommonConfig()
			if err != nil {
				return err
			}

			siteName := commonCfg.DefaultSite
			if siteName == "" {
				siteName = "gogal.dev"
			}

			siteCfg, err := site.LoadSiteConfig(siteName)
			if err != nil {
				return err
			}

			database, err := db.Connect(siteCfg.DatabaseURL())
			if err != nil {
				return err
			}
			defer database.Close()

			importer := coredoctype.NewImporter(database)

			result, err := importer.ImportOneByName(".", args[0], args[1])
			if err != nil {
				return err
			}

			fmt.Println("DocType imported")
			fmt.Println("----------------")
			fmt.Println("DocType:", result.DocType)
			fmt.Println("JSON:   ", result.JSONPath)
			fmt.Println("Fields: ", result.Fields)
			fmt.Println("Perms:  ", result.Perms)

			return nil
		},
	}
}

func newImportAllDocTypesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "import-all",
		Short: "Import all DocType JSON files into the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			commonCfg, err := site.LoadCommonConfig()
			if err != nil {
				return err
			}

			siteName := commonCfg.DefaultSite
			if siteName == "" {
				siteName = "gogal.dev"
			}

			siteCfg, err := site.LoadSiteConfig(siteName)
			if err != nil {
				return err
			}

			database, err := db.Connect(siteCfg.DatabaseURL())
			if err != nil {
				return err
			}
			defer database.Close()

			importer := coredoctype.NewImporter(database)

			results, err := importer.ImportAll(".")
			if err != nil {
				return err
			}

			fmt.Println("DocTypes imported")
			fmt.Println("-----------------")

			for _, result := range results {
				fmt.Printf("%s | fields=%d | perms=%d | %s\n",
					result.DocType,
					result.Fields,
					result.Perms,
					result.JSONPath,
				)
			}

			fmt.Println()
			fmt.Println("Total:", len(results))

			return nil
		},
	}
}
