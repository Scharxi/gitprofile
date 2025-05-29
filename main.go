package main

import (
	"os"

	"github.com/Scharxi/gitprofile/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gitprofile",
		Short: "Manage your git profiles",
		Long: `A CLI tool to manage different git profiles (name, email, GPG key, commit signing)
and activate them per project.`,
	}

	rootCmd.AddCommand(cmd.NewAddCmd())
	rootCmd.AddCommand(cmd.NewListCmd())
	rootCmd.AddCommand(cmd.NewUseCmd())
	rootCmd.AddCommand(cmd.NewStatusCmd())
	rootCmd.AddCommand(cmd.NewCompletionCmd())
	rootCmd.AddCommand(cmd.NewTUICmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
