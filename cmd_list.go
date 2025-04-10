package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all saved git profiles",
		Long:  `List all git profiles stored in ~/.gitprofiles.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			profiles, err := loadProfiles()
			if err != nil {
				return fmt.Errorf("failed to load profiles: %w", err)
			}

			if len(profiles) == 0 {
				fmt.Println("No profiles found")
				return nil
			}

			for name, profile := range profiles {
				fmt.Printf("\nProfile: %s\n", name)
				fmt.Printf("  Name: %s\n", profile.Name)
				fmt.Printf("  Email: %s\n", profile.Email)
				if profile.GPGKey != "" {
					fmt.Printf("  GPG Key: %s\n", profile.GPGKey)
				}
				fmt.Printf("  Sign Commits: %v\n", profile.SignCommits)
			}

			return nil
		},
	}
}
