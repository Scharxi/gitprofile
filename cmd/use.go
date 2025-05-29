package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewUseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use [profile-name]",
		Short: "Use a git profile in the current repository",
		Long: `Set the git configuration for the current repository using a saved profile.
This will set user.name, user.email, GPG signing, and SSH key configuration.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if we're in a git repository
			_, err := runGitCommand("rev-parse", "--git-dir")
			if err != nil {
				return fmt.Errorf("not a git repository (or any of the parent directories)")
			}

			profileName := args[0]
			profiles, err := LoadProfiles()
			if err != nil {
				return fmt.Errorf("failed to load profiles: %w", err)
			}

			profile, exists := profiles[profileName]
			if !exists {
				return fmt.Errorf("profile '%s' not found", profileName)
			}

			// Set user.name
			_, err = runGitCommand("config", "--local", "user.name", profile.Name)
			if err != nil {
				return fmt.Errorf("failed to set user.name: %w", err)
			}

			// Set user.email
			_, err = runGitCommand("config", "--local", "user.email", profile.Email)
			if err != nil {
				return fmt.Errorf("failed to set user.email: %w", err)
			}

			// Handle GPG settings if configured
			if profile.GPGKey != "" {
				_, err = runGitCommand("config", "--local", "user.signingkey", profile.GPGKey)
				if err != nil {
					return fmt.Errorf("failed to set signing key: %w", err)
				}
			}

			// Set commit signing
			signValue := "false"
			if profile.SignCommits {
				signValue = "true"
			}
			_, err = runGitCommand("config", "--local", "commit.gpgsign", signValue)
			if err != nil {
				return fmt.Errorf("failed to set commit signing: %w", err)
			}

			// Configure SSH key if specified
			if profile.SSHKey != "" {
				_, err = runGitCommand("config", "--local", "core.sshCommand", fmt.Sprintf("ssh -i %s", profile.SSHKey))
				if err != nil {
					return fmt.Errorf("failed to set SSH key: %w", err)
				}
			}

			fmt.Printf("Successfully activated profile '%s' in current repository\n", profileName)
			return nil
		},
		ValidArgsFunction: ValidProfileArgsForUse,
	}

	return cmd
}
