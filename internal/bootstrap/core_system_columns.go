package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func EnsureCoreSystemColumns(database *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tables := []string{
		"tabInstalled App",
		"tabInstalled Module",
		"tabModule Def",
		"tabDocType",
		"tabDocField",
		"tabUser",
		"tabRole",
		"tabHas Role",
		"tabDocPerm",
		"tabNaming Series",
	}

	for _, table := range tables {
		statements := []string{
			fmt.Sprintf(`ALTER TABLE %s ADD COLUMN IF NOT EXISTS owner TEXT NOT NULL DEFAULT 'Administrator'`, quoteIdent(table)),
			fmt.Sprintf(`ALTER TABLE %s ADD COLUMN IF NOT EXISTS creation TIMESTAMPTZ NOT NULL DEFAULT NOW()`, quoteIdent(table)),
			fmt.Sprintf(`ALTER TABLE %s ADD COLUMN IF NOT EXISTS modified TIMESTAMPTZ NOT NULL DEFAULT NOW()`, quoteIdent(table)),
			fmt.Sprintf(`ALTER TABLE %s ADD COLUMN IF NOT EXISTS modified_by TEXT NOT NULL DEFAULT 'Administrator'`, quoteIdent(table)),
			fmt.Sprintf(`ALTER TABLE %s ADD COLUMN IF NOT EXISTS docstatus INT NOT NULL DEFAULT 0`, quoteIdent(table)),
			fmt.Sprintf(`ALTER TABLE %s ADD COLUMN IF NOT EXISTS idx INT NOT NULL DEFAULT 0`, quoteIdent(table)),
		}

		for _, statement := range statements {
			if _, err := database.Exec(ctx, statement); err != nil {
				return fmt.Errorf("system column migration failed for %s: %w\nstatement: %s", table, err, statement)
			}
		}
	}

	return nil
}

func quoteIdent(identifier string) string {
	return `"` + identifier + `"`
}
