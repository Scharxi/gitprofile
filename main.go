package main

import (
	"fmt"
	"os"

	"github.com/Scharxi/gitprofile/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gitprofile",
		Short: "A tool for managing multiple git profiles",
		Long: `gitprofile helps you manage multiple git configurations for different projects.
It allows you to save and switch between different git profiles with different names, emails, and GPG keys.`,
	}

	rootCmd.AddCommand(
		cmd.NewAddCmd(),
		cmd.NewListCmd(),
		cmd.NewUseCmd(),
		cmd.NewStatusCmd(),
		cmd.NewDeleteCmd(),
		cmd.NewCompletionCmd(),
		cmd.NewTUICmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
