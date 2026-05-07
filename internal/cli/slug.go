package cli

import (
	"fmt"

	"github.com/galaxylabs/gogal-studio/internal/core/slug"
	"github.com/spf13/cobra"
)

func NewSlugCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slug",
		Short: "Generate safe Gogal folder slugs",
	}

	cmd.AddCommand(newSlugAppCommand())
	cmd.AddCommand(newSlugModuleCommand())
	cmd.AddCommand(newSlugDocTypeCommand())
	cmd.AddCommand(newSlugDocTypePathCommand())

	return cmd
}

func newSlugAppCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "app [app_name]",
		Short: "Generate app slug",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(slug.FromAppName(args[0]))
		},
	}
}

func newSlugModuleCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "module [module_name]",
		Short: "Generate module slug",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(slug.FromModuleName(args[0]))
		},
	}
}

func newSlugDocTypeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "doctype [doctype_name]",
		Short: "Generate DocType slug",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(slug.FromDocTypeName(args[0]))
		},
	}
}

func newSlugDocTypePathCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "doctype-path [module_name] [doctype_name]",
		Short: "Generate DocType folder and JSON path",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Folder:", slug.DocTypeFolderPath(args[0], args[1]))
			fmt.Println("JSON:  ", slug.DocTypeJSONPath(args[0], args[1]))
		},
	}
}
