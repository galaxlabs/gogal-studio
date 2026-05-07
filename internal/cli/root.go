package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "gogal",
	Short: "Gogal Studio CLI",
	Long:  "Gogal Studio bench-like CLI for installing, managing, and running Gogal sites.",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(NewDoctorCommand())
	rootCmd.AddCommand(NewStartCommand())
	rootCmd.AddCommand(NewListCommand())
	rootCmd.AddCommand(NewNamingCommand())
	rootCmd.AddCommand(NewShowCommand())
	rootCmd.AddCommand(NewLifecycleCommand())
	rootCmd.AddCommand(NewFieldTypeCommand())
	rootCmd.AddCommand(NewSlugCommand())
}
