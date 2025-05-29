package cmd

import (
	"github.com/spf13/cobra"
)

// ValidProfileArgs returns a list of valid profile names for completion
func ValidProfileArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	profiles, err := LoadProfiles()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	profileNames := make([]string, 0, len(profiles))
	for name := range profiles {
		profileNames = append(profileNames, name)
	}

	return profileNames, cobra.ShellCompDirectiveNoFileComp
}

// ValidProfileArgsForUse returns a list of valid profile names for the use command
func ValidProfileArgsForUse(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return ValidProfileArgs(cmd, args, toComplete)
} 