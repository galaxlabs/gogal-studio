package cli

import (
	"fmt"
	"sort"

	"github.com/galaxylabs/gogal-studio/internal/core/fieldtype"
	"github.com/spf13/cobra"
)

func NewFieldTypeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fieldtype",
		Short: "Field type registry utilities",
	}

	cmd.AddCommand(newListFieldTypesCommand())
	cmd.AddCommand(newShowFieldTypeCommand())

	return cmd
}

func newListFieldTypesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List registered field types",
		RunE: func(cmd *cobra.Command, args []string) error {
			names := fieldtype.Names()
			sort.Strings(names)

			fmt.Println("Field Types")
			fmt.Println("-----------")

			for _, name := range names {
				def, _ := fieldtype.Get(name)

				fmt.Printf(
					"%s | control=%s | sql=%s | options=%s | columns=%d\n",
					def.Name,
					def.Control,
					def.SQLType,
					def.OptionMode,
					def.DefaultColumns,
				)
			}

			return nil
		},
	}
}

func newShowFieldTypeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show [fieldtype]",
		Short: "Show one field type definition",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			def, err := fieldtype.MustGet(args[0])
			if err != nil {
				return err
			}

			fmt.Println("Field Type")
			fmt.Println("----------")
			fmt.Println("Name:             ", def.Name)
			fmt.Println("Control:          ", def.Control)
			fmt.Println("SQL Type:         ", def.SQLType)
			fmt.Println("Option Mode:      ", def.OptionMode)
			fmt.Println("Requires Options: ", def.RequiresOptions)
			fmt.Println("Is Layout:        ", def.IsLayout)
			fmt.Println("Is Numeric:       ", def.IsNumeric)
			fmt.Println("Is DateTime:      ", def.IsDateTime)
			fmt.Println("Is Text:          ", def.IsText)
			fmt.Println("Is Attach:        ", def.IsAttach)
			fmt.Println("Default Columns:  ", def.DefaultColumns)
			fmt.Println("Default Length:   ", def.DefaultLength)
			fmt.Println("Default Precision:", def.DefaultPrecision)
			fmt.Println("Description:      ", def.Description)

			return nil
		},
	}
}
