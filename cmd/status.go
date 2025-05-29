package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current git profile",
		Long:  `Show the active git profile in the current repository`,
		RunE: func(cmd *cobra.Command, args []string) error {
			profile, name, err := GetCurrentProfile()
			if err != nil {
				return fmt.Errorf("failed to get current profile: %w", err)
			}
			
			if profile == nil {
				fmt.Println("No profile active in current repository")
				return nil
			}
			
			fmt.Printf("Active profile: %s\n", name)
			fmt.Printf("  Name: %s\n", profile.Name)
			fmt.Printf("  Email: %s\n", profile.Email)
			if profile.GPGKey != "" {
				fmt.Printf("  GPG Key: %s\n", profile.GPGKey)
			}
			fmt.Printf("  Sign Commits: %v\n", profile.SignCommits)
			
			return nil
		},
	}
} 