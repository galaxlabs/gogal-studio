package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/galaxylabs/gogal-studio/internal/core/lifecycle"
	"github.com/galaxylabs/gogal-studio/internal/db"
	"github.com/spf13/cobra"
)

type listSiteConfig struct {
	DBType     string `json:"db_type"`
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DBName     string `json:"db_name"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
}

func (cfg listSiteConfig) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
}

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Gogal core records",
	}

	cmd.AddCommand(newListAppsCommand())
	cmd.AddCommand(newListModulesCommand())
	cmd.AddCommand(newListDocTypesCommand())
	cmd.AddCommand(newListDocTypeCommand())
	cmd.AddCommand(newListUsersCommand())
	cmd.AddCommand(newListRolesCommand())
	cmd.AddCommand(newListNamingSeriesCommand())

	return cmd
}

func loadListSiteConfig() (listSiteConfig, error) {
	var cfg listSiteConfig

	raw, err := os.ReadFile("sites/gogal.dev/site_config.json")
	if err != nil {
		return cfg, err
	}

	if err := json.Unmarshal(raw, &cfg); err != nil {
		return cfg, err
	}

	if cfg.DBType == "" {
		cfg.DBType = "postgres"
	}

	if cfg.DBHost == "" {
		cfg.DBHost = "127.0.0.1"
	}

	if cfg.DBPort == 0 {
		cfg.DBPort = 5432
	}

	return cfg, nil
}

func newListAppsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "apps",
		Short: "List installed apps",
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

			rows, err := database.Query(ctx, `
				SELECT name, app_name, app_version
				FROM "tabInstalled App"
				ORDER BY idx, name
			`)
			if err != nil {
				return err
			}
			defer rows.Close()

			fmt.Println("Installed Apps")
			fmt.Println("--------------")

			for rows.Next() {
				var name, appName, appVersion string

				if err := rows.Scan(&name, &appName, &appVersion); err != nil {
					return err
				}

				fmt.Printf("%s | %s | %s\n", name, appName, appVersion)
			}

			return rows.Err()
		},
	}
}

func newListModulesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "modules",
		Short: "List modules",
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

			rows, err := database.Query(ctx, `
				SELECT name, module_name, app_name
				FROM "tabModule Def"
				ORDER BY idx, name
			`)
			if err != nil {
				return err
			}
			defer rows.Close()

			fmt.Println("Modules")
			fmt.Println("-------")

			for rows.Next() {
				var name, moduleName, appName string

				if err := rows.Scan(&name, &moduleName, &appName); err != nil {
					return err
				}

				fmt.Printf("%s | %s | %s\n", name, moduleName, appName)
			}

			return rows.Err()
		},
	}
}

func newListDocTypesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "doctypes",
		Short: "List core DocTypes",
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

			rows, err := database.Query(ctx, `
				SELECT name, module, app_name, table_name, docstatus
				FROM "tabDocType"
				ORDER BY idx, name
			`)
			if err != nil {
				return err
			}
			defer rows.Close()

			fmt.Println("DocTypes")
			fmt.Println("--------")

			for rows.Next() {
				var (
					name      string
					module    string
					appName   string
					tableName string
					docstatus int
				)

				if err := rows.Scan(&name, &module, &appName, &tableName, &docstatus); err != nil {
					return err
				}

				fmt.Printf(
					"%s | %s | %s | %s | status=%s\n",
					name,
					module,
					appName,
					tableName,
					lifecycle.DocStatus(docstatus).String(),
				)
			}

			return rows.Err()
		},
	}
}

func newListDocTypeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "doctype [name]",
		Short: "Show one DocType",
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

			var (
				name          string
				module        string
				appName       string
				tableName     string
				isSingle      bool
				isSubmittable bool
				isChildTable  bool
				isTree        bool
				docstatus     int
			)

			err = database.QueryRow(ctx, `
				SELECT
					name,
					module,
					app_name,
					table_name,
					is_single,
					is_submittable,
					is_child_table,
					is_tree,
					docstatus
				FROM "tabDocType"
				WHERE name = $1
			`, args[0]).Scan(
				&name,
				&module,
				&appName,
				&tableName,
				&isSingle,
				&isSubmittable,
				&isChildTable,
				&isTree,
				&docstatus,
			)
			if err != nil {
				return err
			}

			fmt.Println("DocType")
			fmt.Println("-------")
			fmt.Println("Name:          ", name)
			fmt.Println("Module:        ", module)
			fmt.Println("App:           ", appName)
			fmt.Println("Table:         ", tableName)
			fmt.Println("Is Single:     ", isSingle)
			fmt.Println("Is Submittable:", isSubmittable)
			fmt.Println("Is Child Table:", isChildTable)
			fmt.Println("Is Tree:       ", isTree)
			fmt.Println("Docstatus:     ", docstatus)
			fmt.Println("Status Label:  ", lifecycle.DocStatus(docstatus).String())

			return nil
		},
	}
}

func newListUsersCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "users",
		Short: "List users",
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

			rows, err := database.Query(ctx, `
				SELECT name, username, email
				FROM "tabUser"
				ORDER BY idx, name
			`)
			if err != nil {
				return err
			}
			defer rows.Close()

			fmt.Println("Users")
			fmt.Println("-----")

			for rows.Next() {
				var name, username, email string

				if err := rows.Scan(&name, &username, &email); err != nil {
					return err
				}

				fmt.Printf("%s | %s | %s\n", name, username, email)
			}

			return rows.Err()
		},
	}
}

func newListRolesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "roles",
		Short: "List roles",
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

			rows, err := database.Query(ctx, `
				SELECT name, role_name
				FROM "tabRole"
				ORDER BY idx, name
			`)
			if err != nil {
				return err
			}
			defer rows.Close()

			fmt.Println("Roles")
			fmt.Println("-----")

			for rows.Next() {
				var name, roleName string

				if err := rows.Scan(&name, &roleName); err != nil {
					return err
				}

				fmt.Printf("%s | %s\n", name, roleName)
			}

			return rows.Err()
		},
	}
}

func newListNamingSeriesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "naming-series",
		Short: "List naming series",
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

			rows, err := database.Query(ctx, `
				SELECT
					name,
					series_key,
					prefix,
					current_value,
					digits,
					description,
					docstatus
				FROM "tabNaming Series"
				ORDER BY idx, name
			`)
			if err != nil {
				return err
			}
			defer rows.Close()

			fmt.Println("Naming Series")
			fmt.Println("-------------")

			for rows.Next() {
				var (
					name         string
					seriesKey    string
					prefix       string
					currentValue int64
					digits       int
					description  string
					docstatus    int
				)

				if err := rows.Scan(
					&name,
					&seriesKey,
					&prefix,
					&currentValue,
					&digits,
					&description,
					&docstatus,
				); err != nil {
					return err
				}

				fmt.Printf(
					"%s | key=%s | prefix=%s | current=%d | digits=%d | status=%s | %s\n",
					name,
					seriesKey,
					prefix,
					currentValue,
					digits,
					lifecycle.DocStatus(docstatus).String(),
					description,
				)
			}

			return rows.Err()
		},
	}
}
