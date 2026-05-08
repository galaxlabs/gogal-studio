package cli

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/galaxylabs/gogal-studio/internal/db"
	"github.com/spf13/cobra"
)

func NewDBCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "Database operations",
	}

	cmd.AddCommand(newDBInfoCommand())
	cmd.AddCommand(newDBPingCommand())
	cmd.AddCommand(newDBPSQLCommand())

	return cmd
}

func newDBInfoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Show current site database info",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadListSiteConfig()
			if err != nil {
				return err
			}

			fmt.Println("Gogal Studio database info")
			fmt.Println("--------------------------")
			fmt.Println("DB type:     ", cfg.DBType)
			fmt.Println("DB host:     ", cfg.DBHost)
			fmt.Println("DB port:     ", cfg.DBPort)
			fmt.Println("DB name:     ", cfg.DBName)
			fmt.Println("DB user:     ", cfg.DBUser)
			fmt.Println("DB password: ", maskPassword(cfg.DBPassword))

			return nil
		},
	}
}

func newDBPingCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ping",
		Short: "Test current site database connection",
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

			if err := database.Ping(ctx); err != nil {
				return err
			}

			fmt.Println("Database connection: OK")
			fmt.Println("Database:", cfg.DBName)
			fmt.Println("User:    ", cfg.DBUser)

			return nil
		},
	}
}

func newDBPSQLCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "psql",
		Short: "Print psql command for current site database",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadListSiteConfig()
			if err != nil {
				return err
			}

			fmt.Println("Use this command to open PostgreSQL shell:")
			fmt.Println()

			if runtime.GOOS == "windows" {
				fmt.Printf("docker exec -it door-app-postgres-1 psql -U %s -d %s\n", cfg.DBUser, cfg.DBName)
				fmt.Println()
				fmt.Println("Or if psql is installed locally:")
			}

			fmt.Printf("psql \"postgres://%s:%s@%s:%d/%s?sslmode=disable\"\n",
				cfg.DBUser,
				cfg.DBPassword,
				cfg.DBHost,
				cfg.DBPort,
				cfg.DBName,
			)

			return nil
		},
	}
}

func maskPassword(password string) string {
	if password == "" {
		return ""
	}

	if len(password) <= 4 {
		return "****"
	}

	return password[:2] + "********" + password[len(password)-2:]
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
