package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show current git profile",
		Long:  `Display the active git profile in the current repository`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			profile, profileName, err := GetCurrentProfile()
			if err != nil {
				return fmt.Errorf("failed to get current profile: %w", err)
			}

			w := cmd.OutOrStdout()

			if profile == nil {
				fmt.Fprintln(w, "No active profile found")
				return nil
			}

			fmt.Fprintf(w, "Active profile: %s\n", profileName)
			fmt.Fprintf(w, "  Name: %s\n", profile.Name)
			fmt.Fprintf(w, "  Email: %s\n", profile.Email)
			if profile.GPGKey != "" {
				fmt.Fprintf(w, "  GPG Key: %s\n", profile.GPGKey)
			}
			if profile.SSHKey != "" {
				fmt.Fprintf(w, "  SSH Key: %s\n", profile.SSHKey)
			}
			fmt.Fprintf(w, "  Sign Commits: %v\n", profile.SignCommits)

			return nil
		},
	}

	return cmd
}
