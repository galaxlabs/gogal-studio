package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/galaxylabs/gogal-studio/internal/core/migration"
	"github.com/galaxylabs/gogal-studio/internal/db"
	"github.com/spf13/cobra"
)

func NewMigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Schema migration tools",
	}

	cmd.AddCommand(newMigratePreviewCommand())
	cmd.AddCommand(newMigrateApplyCommand())
	cmd.AddCommand(newMigrateLogsCommand())

	return cmd
}

func newMigrateApplyCommand() *cobra.Command {
	var confirm string

	cmd := &cobra.Command{
		Use:   "apply [doctype]",
		Short: "Apply schema changes for a DocType",
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

			planner := migration.NewPlanner(database)

			plan, err := planner.Apply(context.Background(), args[0], confirm)
			if err != nil {
				return err
			}

			fmt.Println("Migration Applied")
			fmt.Println("-----------------")
			fmt.Println("DocType:    ", plan.DocType)
			fmt.Println("Table:      ", plan.TableName)
			fmt.Println("Table Found:", plan.TableFound)
			fmt.Println()
			fmt.Println("Summary")
			fmt.Println("-------")
			fmt.Println("Applied:              ", plan.Summary.AppliedCount)
			fmt.Println("Skipped:              ", plan.Summary.SkippedCount)
			fmt.Println("Failed:               ", plan.Summary.FailedCount)
			fmt.Println("Dangerous Operations: ", plan.Summary.HasDangerousOperations)
			fmt.Println("Status:               ", plan.Summary.Status)
			fmt.Println()

			for _, op := range plan.Operations {
				if op.Dangerous {
					fmt.Println("Skipped:", op.Action)
				} else {
					fmt.Println("Action:", op.Action)
				}

				if op.Status != "" {
					fmt.Println("Status:", op.Status)
				}

				if op.Column != "" {
					fmt.Println("Column:", op.Column)
				}

				if op.Fieldtype != "" {
					fmt.Println("Field Type:", op.Fieldtype)
				}

				if op.Message != "" {
					fmt.Println("Message:", op.Message)
				}

				if op.SQL != "" {
					fmt.Println("SQL:")
					fmt.Println(op.SQL)
				}

				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&confirm, "confirm", "", "Required confirmation token: APPLY_SCHEMA_CHANGES")

	return cmd
}

func newMigratePreviewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "preview [doctype]",
		Short: "Preview schema changes for a DocType",
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

			planner := migration.NewPlanner(database)

			plan, err := planner.Preview(context.Background(), args[0])
			if err != nil {
				return err
			}

			fmt.Println("Migration Preview")
			fmt.Println("-----------------")
			fmt.Println("DocType:    ", plan.DocType)
			fmt.Println("Table:      ", plan.TableName)
			fmt.Println("Table Found:", plan.TableFound)
			fmt.Println()
			fmt.Println("Summary")
			fmt.Println("-------")
			fmt.Println("Applied:              ", plan.Summary.AppliedCount)
			fmt.Println("Skipped:              ", plan.Summary.SkippedCount)
			fmt.Println("Failed:               ", plan.Summary.FailedCount)
			fmt.Println("Dangerous Operations: ", plan.Summary.HasDangerousOperations)
			fmt.Println("Status:               ", plan.Summary.Status)
			fmt.Println()

			for _, op := range plan.Operations {
				if op.Dangerous {
					fmt.Println("Skipped:", op.Action)
				} else {
					fmt.Println("Action:", op.Action)
				}

				if op.Status != "" {
					fmt.Println("Status:", op.Status)
				}

				if op.Column != "" {
					fmt.Println("Column:", op.Column)
				}

				if op.Fieldtype != "" {
					fmt.Println("Field Type:", op.Fieldtype)
				}

				if op.Message != "" {
					fmt.Println("Message:", op.Message)
				}

				if op.SQL != "" {
					fmt.Println("SQL:")
					fmt.Println(op.SQL)
				}

				fmt.Println()
			}

			return nil
		},
	}
}

func newMigrateLogsCommand() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Show schema migration audit logs",
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

			ctx := context.Background()

			rows, err := database.Query(ctx, `
				SELECT
					doctype,
					table_name,
					operation,
					COALESCE(column_name, ''),
					COALESCE(fieldname, ''),
					COALESCE(fieldtype, ''),
					status,
					COALESCE(message, ''),
					creation
				FROM "tabSchema Migration Log"
				ORDER BY creation DESC
				LIMIT $1
			`, limit)
			if err != nil {
				return err
			}
			defer rows.Close()

			fmt.Println("Schema Migration Logs")
			fmt.Println("---------------------")

			count := 0

			for rows.Next() {
				var (
					doctype    string
					tableName  string
					operation  string
					columnName string
					fieldname  string
					fieldtype  string
					status     string
					message    string
					creation   time.Time
				)

				if err := rows.Scan(
					&doctype,
					&tableName,
					&operation,
					&columnName,
					&fieldname,
					&fieldtype,
					&status,
					&message,
					&creation,
				); err != nil {
					return err
				}

				count++

				fmt.Printf("[%s] %s | %s | %s\n", status, doctype, tableName, operation)

				if columnName != "" {
					fmt.Println("  Column:   ", columnName)
				}

				if fieldname != "" {
					fmt.Println("  Field:    ", fieldname)
				}

				if fieldtype != "" {
					fmt.Println("  Type:     ", fieldtype)
				}

				if message != "" {
					fmt.Println("  Message:  ", message)
				}

				fmt.Println("  Created:  ", creation.Format(time.RFC3339))
				fmt.Println()
			}

			if err := rows.Err(); err != nil {
				return err
			}

			if count == 0 {
				fmt.Println("No migration logs found.")
			}

			return nil
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 20, "Number of migration logs to show")

	return cmd
}
