package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all git profiles",
		Long:  `Display all saved git profiles with their settings`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			profiles, err := LoadProfiles()
			if err != nil {
				return fmt.Errorf("failed to load profiles: %w", err)
			}

			w := cmd.OutOrStdout()

			for name, profile := range profiles {
				fmt.Fprintf(w, "\nProfile: %s\n", name)
				fmt.Fprintf(w, "  Name: %s\n", profile.Name)
				fmt.Fprintf(w, "  Email: %s\n", profile.Email)
				if profile.GPGKey != "" {
					fmt.Fprintf(w, "  GPG Key: %s\n", profile.GPGKey)
				}
				if profile.SSHKey != "" {
					fmt.Fprintf(w, "  SSH Key: %s\n", profile.SSHKey)
				}
				fmt.Fprintf(w, "  Sign Commits: %v\n", profile.SignCommits)
			}

			return nil
		},
	}

	return cmd
}
