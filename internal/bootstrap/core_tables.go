package bootstrap

import (
	"context"
	"fmt"
	"time"

	coreapp "github.com/galaxylabs/gogal-studio/internal/core/app"
	coremodule "github.com/galaxylabs/gogal-studio/internal/core/module"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateCoreTables(database *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	statements := []string{
		`CREATE TABLE IF NOT EXISTS "tabApp" (
			name TEXT PRIMARY KEY,
			creation TIMESTAMPTZ,
			modified TIMESTAMPTZ,
			modified_by TEXT,
			owner TEXT,
			docstatus INT NOT NULL DEFAULT 0,
			idx INT NOT NULL DEFAULT 0,

			app_name TEXT UNIQUE,
			app_title TEXT,
			app_version TEXT,
			is_system BOOLEAN NOT NULL DEFAULT FALSE,
			enabled BOOLEAN NOT NULL DEFAULT TRUE
		)`,

		`CREATE TABLE IF NOT EXISTS "tabModule Def" (
			name TEXT PRIMARY KEY,
			creation TIMESTAMPTZ,
			modified TIMESTAMPTZ,
			modified_by TEXT,
			owner TEXT,
			docstatus INT NOT NULL DEFAULT 0,
			idx INT NOT NULL DEFAULT 0,

			module_name TEXT UNIQUE,
			app_name TEXT,
			custom BOOLEAN NOT NULL DEFAULT FALSE,
			package TEXT,
			enabled BOOLEAN NOT NULL DEFAULT TRUE
		)`,

		`CREATE TABLE IF NOT EXISTS "tabDocType" (
			name TEXT PRIMARY KEY,
			creation TIMESTAMPTZ,
			modified TIMESTAMPTZ,
			modified_by TEXT,
			owner TEXT,
			docstatus INT NOT NULL DEFAULT 0,
			idx INT NOT NULL DEFAULT 0,

			module TEXT,
			app TEXT,
			label TEXT,
			table_name TEXT UNIQUE,

			issingle BOOLEAN NOT NULL DEFAULT FALSE,
			istable BOOLEAN NOT NULL DEFAULT FALSE,
			is_tree BOOLEAN NOT NULL DEFAULT FALSE,
			is_submittable BOOLEAN NOT NULL DEFAULT FALSE,
			editable_grid BOOLEAN NOT NULL DEFAULT TRUE,
			quick_entry BOOLEAN NOT NULL DEFAULT FALSE,
			track_changes BOOLEAN NOT NULL DEFAULT TRUE,

			naming_rule TEXT,
			title_field TEXT,
			sort_field TEXT DEFAULT 'modified',
			sort_order TEXT DEFAULT 'DESC',

			custom BOOLEAN NOT NULL DEFAULT FALSE,
			is_system BOOLEAN NOT NULL DEFAULT FALSE,
			migration_hash TEXT
		)`,

		`CREATE TABLE IF NOT EXISTS "tabDocField" (
			name TEXT PRIMARY KEY,
			creation TIMESTAMPTZ,
			modified TIMESTAMPTZ,
			modified_by TEXT,
			owner TEXT,
			docstatus INT NOT NULL DEFAULT 0,

			parent TEXT,
			parentfield TEXT,
			parenttype TEXT,
			idx INT NOT NULL DEFAULT 0,

			fieldname TEXT,
			label TEXT,
			fieldtype TEXT DEFAULT 'Data',
			options TEXT,

			reqd BOOLEAN NOT NULL DEFAULT FALSE,
			unique_field BOOLEAN NOT NULL DEFAULT FALSE,
			read_only BOOLEAN NOT NULL DEFAULT FALSE,
			hidden BOOLEAN NOT NULL DEFAULT FALSE,
			in_list_view BOOLEAN NOT NULL DEFAULT FALSE,
			in_standard_filter BOOLEAN NOT NULL DEFAULT FALSE,
			in_filter BOOLEAN NOT NULL DEFAULT FALSE,

			default_value TEXT,
			description TEXT,
			depends_on TEXT,
			mandatory_depends_on TEXT,
			read_only_depends_on TEXT,

			columns INT NOT NULL DEFAULT 0,
			length INT NOT NULL DEFAULT 0,
			permlevel INT NOT NULL DEFAULT 0,

			UNIQUE (parent, fieldname)
		)`,

		`CREATE TABLE IF NOT EXISTS "tabRole" (
			name TEXT PRIMARY KEY,
			creation TIMESTAMPTZ,
			modified TIMESTAMPTZ,
			modified_by TEXT,
			owner TEXT,
			docstatus INT NOT NULL DEFAULT 0,
			idx INT NOT NULL DEFAULT 0,

			role_name TEXT UNIQUE,
			desk_access BOOLEAN NOT NULL DEFAULT TRUE,
			is_system BOOLEAN NOT NULL DEFAULT FALSE,
			enabled BOOLEAN NOT NULL DEFAULT TRUE
		)`,

		`CREATE TABLE IF NOT EXISTS "tabUser" (
			name TEXT PRIMARY KEY,
			creation TIMESTAMPTZ,
			modified TIMESTAMPTZ,
			modified_by TEXT,
			owner TEXT,
			docstatus INT NOT NULL DEFAULT 0,
			idx INT NOT NULL DEFAULT 0,

			enabled BOOLEAN NOT NULL DEFAULT TRUE,
			email TEXT UNIQUE,
			username TEXT UNIQUE,
			first_name TEXT,
			last_name TEXT,
			full_name TEXT,
			password_hash TEXT NOT NULL,
			user_type TEXT DEFAULT 'System User'
		)`,

		`CREATE TABLE IF NOT EXISTS "tabHas Role" (
			name TEXT PRIMARY KEY,
			creation TIMESTAMPTZ,
			modified TIMESTAMPTZ,
			modified_by TEXT,
			owner TEXT,
			docstatus INT NOT NULL DEFAULT 0,
			idx INT NOT NULL DEFAULT 0,

			role TEXT,
			parent TEXT,
			parentfield TEXT,
			parenttype TEXT,

			UNIQUE (parent, role)
		)`,

		`CREATE TABLE IF NOT EXISTS "tabDocPerm" (
			name TEXT PRIMARY KEY,
			creation TIMESTAMPTZ,
			modified TIMESTAMPTZ,
			modified_by TEXT,
			owner TEXT,
			docstatus INT NOT NULL DEFAULT 0,

			parent TEXT,
			parentfield TEXT,
			parenttype TEXT,
			idx INT NOT NULL DEFAULT 0,

			permlevel INT NOT NULL DEFAULT 0,
			role TEXT,

			read BOOLEAN NOT NULL DEFAULT TRUE,
			write BOOLEAN NOT NULL DEFAULT TRUE,
			create_perm BOOLEAN NOT NULL DEFAULT TRUE,
			delete_perm BOOLEAN NOT NULL DEFAULT TRUE,
			submit BOOLEAN NOT NULL DEFAULT FALSE,
			cancel BOOLEAN NOT NULL DEFAULT FALSE,
			amend BOOLEAN NOT NULL DEFAULT FALSE,
			report BOOLEAN NOT NULL DEFAULT TRUE,
			export BOOLEAN NOT NULL DEFAULT TRUE,
			import_perm BOOLEAN NOT NULL DEFAULT FALSE,
			share BOOLEAN NOT NULL DEFAULT TRUE,
			print BOOLEAN NOT NULL DEFAULT TRUE,
			email BOOLEAN NOT NULL DEFAULT TRUE,
			if_owner BOOLEAN NOT NULL DEFAULT FALSE,
			select_perm BOOLEAN NOT NULL DEFAULT FALSE,

			UNIQUE (parent, role, permlevel)
		)`,

		`CREATE TABLE IF NOT EXISTS "tabNaming Series" (
			name TEXT PRIMARY KEY,
			series_key TEXT NOT NULL UNIQUE,
			prefix TEXT NOT NULL DEFAULT '',
			current_value BIGINT NOT NULL DEFAULT 0,
			digits INT NOT NULL DEFAULT 5,
			description TEXT NOT NULL DEFAULT '',
			owner TEXT NOT NULL DEFAULT 'Administrator',
			creation TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			modified TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			modified_by TEXT NOT NULL DEFAULT 'Administrator',
			docstatus INT NOT NULL DEFAULT 0,
			idx INT NOT NULL DEFAULT 0
		)`,

		`CREATE TABLE IF NOT EXISTS "tabInstalled App" (
			name TEXT PRIMARY KEY,
			creation TIMESTAMPTZ,
			modified TIMESTAMPTZ,
			modified_by TEXT,
			owner TEXT,
			docstatus INT NOT NULL DEFAULT 0,
			idx INT NOT NULL DEFAULT 0,

			app_name TEXT,
			app_version TEXT,
			git_branch TEXT,
			has_setup_wizard BOOLEAN NOT NULL DEFAULT FALSE,
			is_setup_complete BOOLEAN NOT NULL DEFAULT FALSE
		)`,

		`CREATE TABLE IF NOT EXISTS "tabInstalled Module" (
			name TEXT PRIMARY KEY,
			creation TIMESTAMPTZ,
			modified TIMESTAMPTZ,
			modified_by TEXT,
			owner TEXT,
			docstatus INT NOT NULL DEFAULT 0,
			idx INT NOT NULL DEFAULT 0,

			app_name TEXT,
			module_name TEXT,
			enabled BOOLEAN NOT NULL DEFAULT TRUE,
			sort_order INT NOT NULL DEFAULT 0,

			UNIQUE (app_name, module_name)
		)`,
	}

	if _, err := database.Exec(ctx, `
		ALTER TABLE "tabDocType"
		ADD COLUMN IF NOT EXISTS app_name TEXT NOT NULL DEFAULT 'gogal_studio';

		ALTER TABLE "tabDocType"
		ADD COLUMN IF NOT EXISTS is_single BOOLEAN NOT NULL DEFAULT FALSE;

		ALTER TABLE "tabDocType"
		ADD COLUMN IF NOT EXISTS is_submittable BOOLEAN NOT NULL DEFAULT FALSE;

		ALTER TABLE "tabDocType"
		ADD COLUMN IF NOT EXISTS is_child_table BOOLEAN NOT NULL DEFAULT FALSE;

		ALTER TABLE "tabDocType"
		ADD COLUMN IF NOT EXISTS is_tree BOOLEAN NOT NULL DEFAULT FALSE;
	`); err != nil {
		return err
	}

	for _, statement := range statements {
		if _, err := database.Exec(ctx, statement); err != nil {
			return fmt.Errorf("failed to execute bootstrap statement: %w", err)
		}
	}

	return nil
}
func SeedCoreData(database *pgxpool.Pool, appName string, modules []string, adminPasswordHash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := database.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	now := time.Now()

	if err := coreapp.ValidateAppName(appName); err != nil {
		return fmt.Errorf("invalid app name %q: %w", appName, err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO "tabApp" (
			name,
			creation,
			modified,
			modified_by,
			owner,
			app_name,
			app_title,
			app_version,
			is_system,
			enabled
		)
		VALUES ($1, $2, $2, 'Administrator', 'Administrator', $1, 'Gogal Studio', '0.0.1', TRUE, TRUE)
		ON CONFLICT (name)
		DO UPDATE SET
			modified = EXCLUDED.modified,
			app_title = EXCLUDED.app_title,
			app_version = EXCLUDED.app_version,
			is_system = TRUE,
			enabled = TRUE
	`, appName, now)
	if err != nil {
		return fmt.Errorf("failed to seed tabApp: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO "tabInstalled App" (
			name,
			creation,
			modified,
			modified_by,
			owner,
			app_name,
			app_version,
			has_setup_wizard,
			is_setup_complete
		)
		VALUES ($1, $2, $2, 'Administrator', 'Administrator', $1, '0.0.1', FALSE, TRUE)
		ON CONFLICT (name)
		DO UPDATE SET
			modified = EXCLUDED.modified,
			app_version = EXCLUDED.app_version,
			is_setup_complete = TRUE
	`, appName, now)
	if err != nil {
		return fmt.Errorf("failed to seed tabInstalled App: %w", err)
	}

	for i, moduleName := range modules {
		if err := coremodule.ValidateModuleName(moduleName); err != nil {
			return fmt.Errorf("invalid module name %q: %w", moduleName, err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO "tabModule Def" (
				name,
				creation,
				modified,
				modified_by,
				owner,
				idx,
				module_name,
				app_name,
				custom,
				enabled
			)
			VALUES ($1, $2, $2, 'Administrator', 'Administrator', $3, $1, $4, FALSE, TRUE)
			ON CONFLICT (name)
			DO UPDATE SET
				modified = EXCLUDED.modified,
				idx = EXCLUDED.idx,
				module_name = EXCLUDED.module_name,
				app_name = EXCLUDED.app_name,
				enabled = TRUE
		`, moduleName, now, i+1, appName)
		if err != nil {
			return fmt.Errorf("failed to seed module %s: %w", moduleName, err)
		}

		installedModuleName := appName + ":" + moduleName

		_, err = tx.Exec(ctx, `
			INSERT INTO "tabInstalled Module" (
				name,
				creation,
				modified,
				modified_by,
				owner,
				idx,
				app_name,
				module_name,
				enabled,
				sort_order
			)
			VALUES ($1, $2, $2, 'Administrator', 'Administrator', $3, $4, $5, TRUE, $3)
			ON CONFLICT (name)
			DO UPDATE SET
				modified = EXCLUDED.modified,
				enabled = TRUE,
				sort_order = EXCLUDED.sort_order
		`, installedModuleName, now, i+1, appName, moduleName)
		if err != nil {
			return fmt.Errorf("failed to seed installed module %s: %w", moduleName, err)
		}
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO "tabRole" (
			name,
			creation,
			modified,
			modified_by,
			owner,
			role_name,
			desk_access,
			is_system,
			enabled
		)
		VALUES ('System Manager', $1, $1, 'Administrator', 'Administrator', 'System Manager', TRUE, TRUE, TRUE)
		ON CONFLICT (name)
		DO UPDATE SET
			modified = EXCLUDED.modified,
			desk_access = TRUE,
			is_system = TRUE,
			enabled = TRUE
	`, now)
	if err != nil {
		return fmt.Errorf("failed to seed System Manager role: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO "tabUser" (
			name,
			creation,
			modified,
			modified_by,
			owner,
			enabled,
			email,
			username,
			first_name,
			full_name,
			password_hash,
			user_type
		)
		VALUES (
			'Administrator',
			$1,
			$1,
			'Administrator',
			'Administrator',
			TRUE,
			'administrator@gogal.dev',
			'Administrator',
			'Administrator',
			'Administrator',
			$2,
			'System User'
		)
		ON CONFLICT (name)
		DO UPDATE SET
			modified = EXCLUDED.modified,
			enabled = TRUE,
			email = EXCLUDED.email,
			username = EXCLUDED.username,
			full_name = EXCLUDED.full_name,
			password_hash = EXCLUDED.password_hash,
			user_type = EXCLUDED.user_type
	`, now, adminPasswordHash)
	if err != nil {
		return fmt.Errorf("failed to seed Administrator user: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO "tabHas Role" (
			name,
			creation,
			modified,
			modified_by,
			owner,
			role,
			parent,
			parentfield,
			parenttype
		)
		VALUES (
			'Administrator-System Manager',
			$1,
			$1,
			'Administrator',
			'Administrator',
			'System Manager',
			'Administrator',
			'roles',
			'User'
		)
		ON CONFLICT (name)
		DO UPDATE SET
			modified = EXCLUDED.modified,
			role = EXCLUDED.role,
			parent = EXCLUDED.parent,
			parentfield = EXCLUDED.parentfield,
			parenttype = EXCLUDED.parenttype
	`, now)
	if err != nil {
		return fmt.Errorf("failed to assign System Manager role: %w", err)
	}

	return tx.Commit(ctx)
}
