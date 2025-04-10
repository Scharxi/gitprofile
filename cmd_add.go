package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	var name, email, gpgKey string
	var signCommits bool

	cmd := &cobra.Command{
		Use:   "add [profile-name]",
		Short: "Add a new git profile",
		Long: `Add a new git profile with name, email, and optional GPG key settings.
The profile will be saved in ~/.gitprofiles.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]
			profiles, err := loadProfiles()
			if err != nil {
				return fmt.Errorf("failed to load profiles: %w", err)
			}

			profiles[profileName] = Profile{
				Name:        name,
				Email:       email,
				GPGKey:      gpgKey,
				SignCommits: signCommits,
			}

			if err := saveProfiles(profiles); err != nil {
				return fmt.Errorf("failed to save profiles: %w", err)
			}

			fmt.Printf("Profile '%s' added successfully\n", profileName)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Git user name")
	cmd.Flags().StringVar(&email, "email", "", "Git email")
	cmd.Flags().StringVar(&gpgKey, "gpg-key", "", "GPG key ID")
	cmd.Flags().BoolVar(&signCommits, "sign", false, "Enable commit signing")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("email")

	return cmd
} 