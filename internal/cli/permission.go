package cli

import (
	"context"
	"fmt"

	"github.com/galaxylabs/gogal-studio/internal/core/permission"
	"github.com/galaxylabs/gogal-studio/internal/db"
	"github.com/spf13/cobra"
)

func NewPermissionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "permission",
		Short: "Permission utilities",
	}

	cmd.AddCommand(newPermissionCheckCommand())
	cmd.AddCommand(newPermissionCheckUserCommand())

	return cmd
}

func newPermissionCheckCommand() *cobra.Command {
	var role string
	var user string
	var doctype string

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check role or user permissions for a DocType",
		RunE: func(cmd *cobra.Command, args []string) error {
			if role == "" && user == "" {
				role = "System Manager"
			}

			if doctype == "" {
				return fmt.Errorf("--doctype is required")
			}

			cfg, err := loadListSiteConfig()
			if err != nil {
				return err
			}

			database, err := db.Connect(cfg.DatabaseURL())
			if err != nil {
				return err
			}
			defer database.Close()

			checker := permission.NewChecker(database)
			ctx := context.Background()

			if user != "" {
				roles, err := checker.GetUserRoles(ctx, user)
				if err != nil {
					return err
				}

				canRead, err := checker.CanUserRead(ctx, user, doctype)
				if err != nil {
					return err
				}

				canCreate, err := checker.CanUserCreate(ctx, user, doctype)
				if err != nil {
					return err
				}

				canWrite, err := checker.CanUserWrite(ctx, user, doctype)
				if err != nil {
					return err
				}

				canDelete, err := checker.CanUserDelete(ctx, user, doctype)
				if err != nil {
					return err
				}

				fmt.Println("User Permission Check")
				fmt.Println("---------------------")
				fmt.Println("User:    ", user)
				fmt.Println("Roles:   ", roles)
				fmt.Println("DocType: ", doctype)
				fmt.Println()
				fmt.Println("Read:    ", canRead)
				fmt.Println("Create:  ", canCreate)
				fmt.Println("Write:   ", canWrite)
				fmt.Println("Delete:  ", canDelete)

				return nil
			}

			canRead, err := checker.CanRead(ctx, doctype, []string{role})
			if err != nil {
				return err
			}

			canCreate, err := checker.CanCreate(ctx, doctype, []string{role})
			if err != nil {
				return err
			}

			canWrite, err := checker.CanWrite(ctx, doctype, []string{role})
			if err != nil {
				return err
			}

			canDelete, err := checker.CanDelete(ctx, doctype, []string{role})
			if err != nil {
				return err
			}

			fmt.Println("Permission Check")
			fmt.Println("----------------")
			fmt.Println("Role:    ", role)
			fmt.Println("DocType: ", doctype)
			fmt.Println()
			fmt.Println("Read:    ", canRead)
			fmt.Println("Create:  ", canCreate)
			fmt.Println("Write:   ", canWrite)
			fmt.Println("Delete:  ", canDelete)

			return nil
		},
	}

	cmd.Flags().StringVar(&role, "role", "System Manager", "Role name")
	cmd.Flags().StringVar(&user, "user", "", "User name")
	cmd.Flags().StringVar(&doctype, "doctype", "", "DocType name")

	return cmd
}

func newPermissionCheckUserCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "check-user [username] [doctype] [action]",
		Short: "Check permission for one user and action",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]
			doctypeName := args[1]
			action := permission.Action(args[2])

			cfg, err := loadListSiteConfig()
			if err != nil {
				return err
			}

			database, err := db.Connect(cfg.DatabaseURL())
			if err != nil {
				return err
			}
			defer database.Close()

			checker := permission.NewChecker(database)

			roles, err := permission.GetRolesForUser(context.Background(), database, username)
			if err != nil {
				return err
			}

			allowed, err := checker.CanUser(context.Background(), username, doctypeName, action)
			if err != nil {
				return err
			}

			fmt.Println("User Permission Check")
			fmt.Println("---------------------")
			fmt.Println("User:   ", username)
			fmt.Println("Roles:  ", roles)
			fmt.Println("DocType:", doctypeName)
			fmt.Println("Action: ", action)
			fmt.Println("Allowed:", allowed)

			return nil
		},
	}
}
