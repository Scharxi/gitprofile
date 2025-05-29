package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [profile-name]",
		Short: "Delete a git profile",
		Long:  `Delete a saved git profile from ~/.gitprofiles.json`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]
			profiles, err := LoadProfiles()
			if err != nil {
				return fmt.Errorf("failed to load profiles: %w", err)
			}

			if _, exists := profiles[profileName]; !exists {
				return fmt.Errorf("profile '%s' not found", profileName)
			}

			delete(profiles, profileName)

			if err := SaveProfiles(profiles); err != nil {
				return fmt.Errorf("failed to save profiles: %w", err)
			}

			fmt.Printf("Profile '%s' deleted successfully\n", profileName)
			return nil
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) > 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			profiles, err := LoadProfiles()
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			var names []string
			for name := range profiles {
				names = append(names, name)
			}
			return names, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return cmd
}
