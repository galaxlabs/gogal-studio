package cli

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/galaxylabs/gogal-studio/internal/bootstrap"
	"github.com/galaxylabs/gogal-studio/internal/db"
	"github.com/galaxylabs/gogal-studio/internal/site"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Gogal Studio infrastructure",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Gogal Studio installer")
		fmt.Println("----------------------")

		if err := createDefaultSiteFiles(); err != nil {
			return err
		}

		commonCfg, err := site.LoadCommonConfig()
		if err != nil {
			return err
		}

		siteCfg, err := site.LoadSiteConfig(commonCfg.DefaultSite)
		if err != nil {
			return err
		}

		if siteCfg.DBPassword == "" || siteCfg.DBPassword == "change-this-password" {
			siteCfg.DBPassword = generateSecurePassword()
			if err := site.SaveSiteConfig(siteCfg); err != nil {
				return err
			}
		}

		fmt.Println("Loaded config:")
		fmt.Println("Default site:", commonCfg.DefaultSite)
		fmt.Println("Database:", siteCfg.DBType, siteCfg.DBHost, siteCfg.DBPort, siteCfg.DBName)
		fmt.Println("DB user:", siteCfg.DBUser)
		fmt.Println("Installed app:", siteCfg.InstalledApps)
		fmt.Println("Installed modules:", siteCfg.InstalledModules)
		fmt.Println()

		fmt.Println("Preparing PostgreSQL database and user...")

		if err := setupPostgresInfrastructure(siteCfg); err != nil {
			return err
		}

		fmt.Println("PostgreSQL database/user: OK")
		fmt.Println("Testing site database connection...")

		database, err := db.Connect(siteCfg.DatabaseURL())
		if err != nil {
			return fmt.Errorf("postgres connection failed: %w", err)
		}
		defer database.Close()

		fmt.Println("PostgreSQL connection: OK")
		fmt.Println("Creating Gogal core metadata bootstrap tables...")

		if err := bootstrap.CreateCoreTables(database); err != nil {
			return err
		}

		fmt.Println("Gogal core metadata tables: OK")
		fmt.Println("Seeding default app, modules, role, and Administrator...")

		adminPasswordHash, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash Administrator password: %w", err)
		}

		if err := bootstrap.SeedCoreData(
			database,
			"gogal_studio",
			siteCfg.InstalledModules,
			string(adminPasswordHash),
		); err != nil {
			return err
		}

		fmt.Println("Core seed data: OK")
		fmt.Println("Seeding core DocType metadata records...")
		if err := bootstrap.SeedCoreDocTypeMetadata(database); err != nil {
			return err
		}
		fmt.Println("Core DocType metadata seed: OK")
		fmt.Println("Expanding Core metadata columns...")
		if err := bootstrap.EnsureCoreMetaColumns(database); err != nil {
			return err
		}
		fmt.Println("Core metadata columns: OK")
		fmt.Println("Adding Core system columns...")
		if err := bootstrap.EnsureCoreSystemColumns(database); err != nil {
			return err
		}
		fmt.Println("Core system columns: OK")
		fmt.Println("Seeding default DocPerm records...")
		if err := bootstrap.SeedCoreDocPerms(database); err != nil {
			return err
		}
		fmt.Println("Default DocPerm seed: OK")
		fmt.Println("Seeding default Naming Series records...")
		if err := bootstrap.SeedDefaultNamingSeries(database); err != nil {
			return err
		}
		fmt.Println("Default Naming Series seed: OK")
		fmt.Println()
		fmt.Println("Administrator username: Administrator")
		fmt.Println("Administrator password: admin")
		fmt.Println()
		fmt.Println("Next step: verify seeded core DocType metadata records")
		return nil
	},
}

func createDefaultSiteFiles() error {
	sitesDir := "sites"
	defaultSite := "gogal.dev"
	defaultSiteDir := filepath.Join(sitesDir, defaultSite)

	if err := os.MkdirAll(defaultSiteDir, 0755); err != nil {
		return fmt.Errorf("failed to create site directory: %w", err)
	}

	commonPath := filepath.Join(sitesDir, "common_site_config.json")
	sitePath := filepath.Join(defaultSiteDir, "site_config.json")

	if _, err := os.Stat(commonPath); os.IsNotExist(err) {
		commonConfig := site.CommonSiteConfig{
			DefaultSite:   defaultSite,
			ServerPort:    8080,
			DeveloperMode: true,
		}

		if err := writeJSON(commonPath, commonConfig); err != nil {
			return err
		}

		fmt.Println("Created sites/common_site_config.json")
	}

	if _, err := os.Stat(sitePath); os.IsNotExist(err) {
		siteConfig := site.SiteConfig{
			SiteName:   defaultSite,
			DBType:     "postgres",
			DBHost:     "127.0.0.1",
			DBPort:     5432,
			DBName:     "gogal_dev",
			DBUser:     "gogal_dev_user",
			DBPassword: generateSecurePassword(),
			InstalledApps: []string{
				"gogal_studio",
			},
			InstalledModules: []string{
				"Core",
				"Setup",
				"Security",
				"Desk",
				"Workspace",
				"Navigation",
			},
		}

		if err := writeJSON(sitePath, siteConfig); err != nil {
			return err
		}

		fmt.Println("Created sites/gogal.dev/site_config.json")
	}

	return nil
}

func setupPostgresInfrastructure(siteCfg site.SiteConfig) error {
	if siteCfg.DBType != "postgres" {
		return fmt.Errorf("only postgres is supported in this milestone")
	}

	adminURL := os.Getenv("GOGAL_PG_ADMIN_URL")
	if adminURL == "" {
		adminURL = fmt.Sprintf(
			"postgres://door:door@%s:%d/postgres?sslmode=disable",
			siteCfg.DBHost,
			siteCfg.DBPort,
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	adminDB, err := pgxpool.New(ctx, adminURL)
	if err != nil {
		return fmt.Errorf("failed to create postgres admin connection: %w", err)
	}
	defer adminDB.Close()

	if err := adminDB.Ping(ctx); err != nil {
		return fmt.Errorf("failed to connect postgres admin database: %w", err)
	}

	if err := createOrUpdatePostgresRole(ctx, adminDB, siteCfg.DBUser, siteCfg.DBPassword); err != nil {
		return err
	}

	if err := createPostgresDatabaseIfMissing(ctx, adminDB, siteCfg.DBName, siteCfg.DBUser); err != nil {
		return err
	}

	siteAdminURL := fmt.Sprintf(
		"postgres://door:door@%s:%d/%s?sslmode=disable",
		siteCfg.DBHost,
		siteCfg.DBPort,
		siteCfg.DBName,
	)

	siteAdminDB, err := pgxpool.New(ctx, siteAdminURL)
	if err != nil {
		return fmt.Errorf("failed to connect site database as admin: %w", err)
	}
	defer siteAdminDB.Close()

	if err := grantPostgresPermissions(ctx, siteAdminDB, siteCfg.DBName, siteCfg.DBUser); err != nil {
		return err
	}

	return nil
}

func createOrUpdatePostgresRole(ctx context.Context, pool *pgxpool.Pool, roleName string, password string) error {
	if !isSafeIdentifier(roleName) {
		return fmt.Errorf("unsafe postgres role name: %s", roleName)
	}

	var exists bool

	if err := pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM pg_roles WHERE rolname = $1
		)
	`, roleName).Scan(&exists); err != nil {
		return err
	}

	quotedRole := quoteIdentifier(roleName)
	quotedPassword := quoteLiteral(password)

	if exists {
		_, err := pool.Exec(ctx, fmt.Sprintf(
			"ALTER ROLE %s WITH LOGIN PASSWORD %s",
			quotedRole,
			quotedPassword,
		))
		return err
	}

	_, err := pool.Exec(ctx, fmt.Sprintf(
		"CREATE ROLE %s WITH LOGIN PASSWORD %s",
		quotedRole,
		quotedPassword,
	))

	return err
}

func createPostgresDatabaseIfMissing(ctx context.Context, pool *pgxpool.Pool, dbName string, owner string) error {
	if !isSafeIdentifier(dbName) {
		return fmt.Errorf("unsafe postgres database name: %s", dbName)
	}

	if !isSafeIdentifier(owner) {
		return fmt.Errorf("unsafe postgres owner name: %s", owner)
	}

	var exists bool

	if err := pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM pg_database WHERE datname = $1
		)
	`, dbName).Scan(&exists); err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err := pool.Exec(ctx, fmt.Sprintf(
		"CREATE DATABASE %s OWNER %s",
		quoteIdentifier(dbName),
		quoteIdentifier(owner),
	))

	return err
}

func grantPostgresPermissions(ctx context.Context, pool *pgxpool.Pool, dbName string, userName string) error {
	if !isSafeIdentifier(dbName) {
		return fmt.Errorf("unsafe postgres database name: %s", dbName)
	}

	if !isSafeIdentifier(userName) {
		return fmt.Errorf("unsafe postgres user name: %s", userName)
	}

	quotedDB := quoteIdentifier(dbName)
	quotedUser := quoteIdentifier(userName)

	statements := []string{
		fmt.Sprintf("GRANT CONNECT ON DATABASE %s TO %s", quotedDB, quotedUser),
		fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", quotedDB, quotedUser),
		fmt.Sprintf("GRANT USAGE, CREATE ON SCHEMA public TO %s", quotedUser),
		fmt.Sprintf("ALTER SCHEMA public OWNER TO %s", quotedUser),
		fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO %s", quotedUser),
		fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO %s", quotedUser),
	}

	for _, statement := range statements {
		if _, err := pool.Exec(ctx, statement); err != nil {
			return fmt.Errorf("failed grant statement %q: %w", statement, err)
		}
	}

	return nil
}

func generateSecurePassword() string {
	raw := make([]byte, 32)

	if _, err := rand.Read(raw); err != nil {
		return "_RbmVJTd84fHV_fallback_change_me"
	}

	return "_" + base64.RawURLEncoding.EncodeToString(raw)
}

func isSafeIdentifier(value string) bool {
	if value == "" {
		return false
	}

	for _, ch := range value {
		if (ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '_' {
			continue
		}

		return false
	}

	return true
}

func quoteIdentifier(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
}

func quoteLiteral(value string) string {
	return `'` + strings.ReplaceAll(value, `'`, `''`) + `'`
}

func writeJSON(path string, data any) error {
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode json for %s: %w", path, err)
	}

	raw = append(raw, '\n')

	if err := os.WriteFile(path, raw, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}

	return nil
}
