package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

func NewUseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use [profile-name]",
		Short: "Use a git profile in the current repository",
		Long: `Set the git configuration for the current repository using a saved profile.
This will set user.name, user.email, and optionally configure GPG signing.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if we're in a git repository
			if err := runGitCommand("rev-parse", "--git-dir"); err != nil {
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
			if err := runGitCommand("config", "--local", "user.name", profile.Name); err != nil {
				return fmt.Errorf("failed to set user.name: %w", err)
			}

			// Set user.email
			if err := runGitCommand("config", "--local", "user.email", profile.Email); err != nil {
				return fmt.Errorf("failed to set user.email: %w", err)
			}

			// Handle GPG settings if configured
			if profile.GPGKey != "" {
				if err := runGitCommand("config", "--local", "user.signingkey", profile.GPGKey); err != nil {
					return fmt.Errorf("failed to set signing key: %w", err)
				}
			}

			// Set commit signing
			signValue := "false"
			if profile.SignCommits {
				signValue = "true"
			}
			if err := runGitCommand("config", "--local", "commit.gpgsign", signValue); err != nil {
				return fmt.Errorf("failed to set commit signing: %w", err)
			}

			fmt.Printf("Successfully activated profile '%s' in current repository\n", profileName)
			return nil
		},
		ValidArgsFunction: ValidProfileArgsForUse,
	}

	return cmd
}

func runGitCommand(args ...string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("git.exe", args...)
	} else {
		cmd = exec.Command("git", args...)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), strings.TrimSpace(string(output)))
	}
	return nil
} 