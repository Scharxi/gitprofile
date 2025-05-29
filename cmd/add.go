package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewAddCmd() *cobra.Command {
	var name, email, gpgKey, sshKey string
	var signCommits bool

	cmd := &cobra.Command{
		Use:   "add [profile-name]",
		Short: "Add a new git profile",
		Long: `Add a new git profile with name, email, and optional GPG key and SSH key settings.
The profile will be saved in ~/.gitprofiles.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]
			profiles, err := LoadProfiles()
			if err != nil {
				return fmt.Errorf("failed to load profiles: %w", err)
			}

			profiles[profileName] = Profile{
				Name:        name,
				Email:       email,
				GPGKey:      gpgKey,
				SignCommits: signCommits,
				SSHKey:      sshKey,
			}

			if err := SaveProfiles(profiles); err != nil {
				return fmt.Errorf("failed to save profiles: %w", err)
			}

			fmt.Printf("Profile '%s' added successfully\n", profileName)
			return nil
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) > 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Git user name")
	cmd.Flags().StringVar(&email, "email", "", "Git email")
	cmd.Flags().StringVar(&gpgKey, "gpg-key", "", "GPG key ID")
	cmd.Flags().StringVar(&sshKey, "ssh-key", "", "SSH key file path (e.g., ~/.ssh/id_rsa)")
	cmd.Flags().BoolVar(&signCommits, "sign", false, "Enable commit signing")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("email")

	return cmd
}
