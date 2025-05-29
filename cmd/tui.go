package cmd

import (
	"fmt"
	"strconv"

	"github.com/Scharxi/gitprofile/cmd/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func NewTUICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tui",
		Short: "Start the terminal user interface",
		Long:  `Launch an interactive terminal user interface to manage git profiles.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			profiles, err := LoadProfiles()
			if err != nil {
				return fmt.Errorf("failed to load profiles: %w", err)
			}

			// Get profile names
			var profileNames []string
			for name := range profiles {
				profileNames = append(profileNames, name)
			}

			if len(profileNames) == 0 {
				return fmt.Errorf("no profiles found. Add a profile first using 'gitprofile add'")
			}

			// Create and run the TUI
			model := tui.NewProfileSelector(profileNames)
			p := tea.NewProgram(model)

			m, err := p.Run()
			if err != nil {
				return fmt.Errorf("error running TUI: %w", err)
			}

			selector := m.(*tui.ProfileSelector)
			selected := selector.Selected()

			if selected == "" {
				return nil
			}

			if selector.IsEditing() {
				// Get the profile to edit
				profile := profiles[selected]

				// Create and run the editor
				editor := tui.NewProfileEditor(
					profile.Name,
					profile.Email,
					profile.GPGKey,
					profile.SSHKey,
					profile.SignCommits,
				)

				p = tea.NewProgram(editor)
				m, err = p.Run()
				if err != nil {
					return fmt.Errorf("error running editor: %w", err)
				}

				editorModel := m.(*tui.ProfileEditor)
				if editorModel.IsSaved() {
					// Update the profile with new values
					fields := editorModel.GetFields()
					signCommits, _ := strconv.ParseBool(fields["Sign Commits"])

					profiles[selected] = Profile{
						Name:        fields["Name"],
						Email:       fields["Email"],
						GPGKey:      fields["GPG Key"],
						SSHKey:      fields["SSH Key"],
						SignCommits: signCommits,
					}

					if err := SaveProfiles(profiles); err != nil {
						return fmt.Errorf("failed to save profiles: %w", err)
					}

					fmt.Printf("Profile '%s' updated successfully\n", selected)
				}

				return nil
			}

			// Use the selected profile
			useCmd := NewUseCmd()
			return useCmd.RunE(cmd, []string{selected})
		},
	}

	return cmd
}
