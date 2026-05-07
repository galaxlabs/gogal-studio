package cli

import (
	"fmt"
	"strconv"

	"github.com/galaxylabs/gogal-studio/internal/core/lifecycle"
	"github.com/spf13/cobra"
)

func NewLifecycleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lifecycle",
		Short: "Document lifecycle utilities",
	}

	cmd.AddCommand(newDocStatusCommand())

	return cmd
}

func newDocStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "docstatus [value]",
		Short: "Inspect docstatus value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			value, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			status, err := lifecycle.Parse(value)
			if err != nil {
				return err
			}

			fmt.Println("DocStatus")
			fmt.Println("---------")
			fmt.Println("Value:      ", status.Int())
			fmt.Println("Label:      ", status.String())
			fmt.Println("Can Edit:   ", lifecycle.CanEdit(status.Int(), true))
			fmt.Println("Can Submit: ", lifecycle.CanSubmit(status.Int()))
			fmt.Println("Can Cancel: ", lifecycle.CanCancel(status.Int()))
			fmt.Println("Can Amend:  ", lifecycle.CanAmend(status.Int()))
			fmt.Println("Actions:    ", lifecycle.AvailableActions(status.Int(), true))

			return nil
		},
	}
}
