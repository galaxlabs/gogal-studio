package permission

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetRolesForUser(ctx context.Context, db *pgxpool.Pool, username string) ([]string, error) {
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := db.Query(ctx, `
		SELECT role
		FROM "tabHas Role"
		WHERE parent = $1
		ORDER BY idx, role
	`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := []string{}

	for rows.Next() {
		var role string

		if err := rows.Scan(&role); err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func (c *Checker) GetUserRoles(ctx context.Context, user string) ([]string, error) {
	return GetRolesForUser(ctx, c.DB, user)
}
