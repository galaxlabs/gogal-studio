package cli

import (
	"context"
	"fmt"

	"github.com/galaxylabs/gogal-studio/internal/core/meta"
	"github.com/galaxylabs/gogal-studio/internal/db"
	"github.com/spf13/cobra"
)

func NewMetaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "meta",
		Short: "DocType metadata utilities",
	}

	cmd.AddCommand(newMetaShowCommand())

	return cmd
}

func newMetaShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show [doctype]",
		Short: "Show loaded DocType metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadListSiteConfig()
			if err != nil {
				return err
			}

			database, err := db.Connect(cfg.DatabaseURL())
			if err != nil {
				return err
			}
			defer database.Close()

			loader := meta.NewLoader(database)

			result, err := loader.GetDocTypeMeta(context.Background(), args[0])
			if err != nil {
				return err
			}

			fmt.Println("DocType Meta")
			fmt.Println("------------")
			fmt.Println("Name:        ", result.Name)
			fmt.Println("Module:      ", result.Module)
			fmt.Println("App:         ", result.AppName)
			fmt.Println("Table:       ", result.TableName)
			fmt.Println("Title Field: ", result.TitleField)
			fmt.Println("Autoname:    ", result.Autoname)
			fmt.Println("Submittable: ", result.IsSubmittable)
			fmt.Println("Child Table: ", result.IsChildTable)
			fmt.Println("Fields:      ", len(result.Fields))
			fmt.Println("Permissions: ", len(result.Permissions))

			fmt.Println()
			fmt.Println("Fields")
			fmt.Println("------")
			for _, field := range result.Fields {
				fmt.Printf("%d | %s | %s | %s | hidden=%v\n",
					field.Idx,
					field.Fieldname,
					field.Label,
					field.Fieldtype,
					field.Hidden,
				)
			}

			return nil
		},
	}
}
