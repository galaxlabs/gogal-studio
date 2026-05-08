package permission

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Action string

const (
	ActionRead   Action = "read"
	ActionWrite  Action = "write"
	ActionCreate Action = "create"
	ActionDelete Action = "delete"
	ActionSubmit Action = "submit"
	ActionCancel Action = "cancel"
	ActionAmend  Action = "amend"
	ActionPrint  Action = "print"
	ActionEmail  Action = "email"
	ActionExport Action = "export"
	ActionImport Action = "import"
	ActionShare  Action = "share"
	ActionReport Action = "report"
)

type Checker struct {
	DB *pgxpool.Pool
}

func NewChecker(db *pgxpool.Pool) *Checker {
	return &Checker{DB: db}
}

func (c *Checker) Can(ctx context.Context, doctype string, roles []string, action Action) (bool, error) {
	if doctype == "" {
		return false, fmt.Errorf("doctype is required")
	}

	if len(roles) == 0 {
		return false, nil
	}

	column, err := actionColumn(action)
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM "tabDocPerm"
			WHERE parent = $1
			AND role = ANY($2)
			AND %s = TRUE
		)
	`, column)

	var allowed bool

	if err := c.DB.QueryRow(ctx, query, doctype, roles).Scan(&allowed); err != nil {
		return false, err
	}

	return allowed, nil
}

func actionColumn(action Action) (string, error) {
	switch action {
	case ActionRead:
		return `"read"`, nil
	case ActionWrite:
		return `"write"`, nil
	case ActionCreate:
		return "create_perm", nil
	case ActionDelete:
		return "delete_perm", nil
	case ActionSubmit:
		return "submit_perm", nil
	case ActionCancel:
		return "cancel_perm", nil
	case ActionAmend:
		return "amend_perm", nil
	case ActionPrint:
		return "print_perm", nil
	case ActionEmail:
		return "email_perm", nil
	case ActionExport:
		return "export_perm", nil
	case ActionImport:
		return "import_perm", nil
	case ActionShare:
		return "share_perm", nil
	case ActionReport:
		return "report_perm", nil
	default:
		return "", fmt.Errorf("unknown permission action: %s", action)
	}
}

func (c *Checker) CanUser(ctx context.Context, username string, doctype string, action Action) (bool, error) {
	roles, err := GetRolesForUser(ctx, c.DB, username)
	if err != nil {
		return false, err
	}

	return c.Can(ctx, doctype, roles, action)
}
