package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gitprofile",
		Short: "Manage your git profiles",
		Long: `A CLI tool to manage different git profiles (name, email, GPG key, commit signing)
and activate them per project.`,
	}

	rootCmd.AddCommand(newAddCmd())
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newUseCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
