package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/galaxylabs/gogal-studio/internal/core/lifecycle"
	"github.com/galaxylabs/gogal-studio/internal/db"
	"github.com/spf13/cobra"
)

func NewShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show Gogal Studio core records",
	}

	cmd.AddCommand(newShowDocTypeCommand())
	cmd.AddCommand(newShowNamingSeriesCommand())

	return cmd
}

func newShowDocTypeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "doctype [name]",
		Short: "Show one DocType with permissions",
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

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			fmt.Printf("DocType: %s\n\n", args[0])

			permRows, err := database.Query(ctx, `
				SELECT
					role,
					permlevel,
					"read",
					"write",
					create_perm,
					delete_perm,
					submit_perm,
					cancel_perm,
					amend_perm,
					idx
				FROM "tabDocPerm"
				WHERE parent = $1
				ORDER BY idx, name
			`, args[0])
			if err != nil {
				return err
			}
			defer permRows.Close()

			fmt.Println("Permissions")
			fmt.Println("-----------")

			for permRows.Next() {
				var (
					role       string
					permlevel  int
					read       bool
					write      bool
					createPerm bool
					deletePerm bool
					submitPerm bool
					cancelPerm bool
					amendPerm  bool
					idx        int
				)

				if err := permRows.Scan(
					&role,
					&permlevel,
					&read,
					&write,
					&createPerm,
					&deletePerm,
					&submitPerm,
					&cancelPerm,
					&amendPerm,
					&idx,
				); err != nil {
					return err
				}

				fmt.Printf(
					"%d | %s | level=%d | read=%v write=%v create=%v delete=%v submit=%v cancel=%v amend=%v\n",
					idx,
					role,
					permlevel,
					read,
					write,
					createPerm,
					deletePerm,
					submitPerm,
					cancelPerm,
					amendPerm,
				)
			}

			return permRows.Err()
		},
	}
}

func newShowNamingSeriesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "naming-series [series_key]",
		Short: "Show one naming series",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			seriesKey := args[0]

			cfg, err := loadListSiteConfig()
			if err != nil {
				return err
			}

			database, err := db.Connect(cfg.DatabaseURL())
			if err != nil {
				return err
			}
			defer database.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			var (
				name         string
				prefix       string
				currentValue int64
				digits       int
				description  string
				owner        string
				modifiedBy   string
				docstatus    int
				idx          int
			)

			err = database.QueryRow(ctx, `
				SELECT
					name,
					prefix,
					current_value,
					digits,
					description,
					owner,
					modified_by,
					docstatus,
					idx
				FROM "tabNaming Series"
				WHERE series_key = $1
			`, seriesKey).Scan(
				&name,
				&prefix,
				&currentValue,
				&digits,
				&description,
				&owner,
				&modifiedBy,
				&docstatus,
				&idx,
			)
			if err != nil {
				return err
			}

			fmt.Println("Naming Series")
			fmt.Println("-------------")
			fmt.Println("Name:         ", name)
			fmt.Println("Series Key:   ", seriesKey)
			fmt.Println("Prefix:       ", prefix)
			fmt.Println("Current Value:", currentValue)
			fmt.Println("Digits:       ", digits)
			fmt.Println("Description:  ", description)
			fmt.Println("Owner:        ", owner)
			fmt.Println("Modified By:  ", modifiedBy)
			fmt.Println("Docstatus:    ", docstatus)
			fmt.Println("Status Label: ", lifecycle.DocStatus(docstatus).String())
			fmt.Println("Idx:          ", idx)

			return nil
		},
	}
}
