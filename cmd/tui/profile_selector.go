package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProfileSelector struct {
	profiles []string
	cursor   int
	selected string
	err      error
	quitting bool
	editing  bool
}

func NewProfileSelector(profiles []string) *ProfileSelector {
	return &ProfileSelector{
		profiles: profiles,
	}
}

func (m *ProfileSelector) Init() tea.Cmd {
	return nil
}

func (m *ProfileSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.profiles)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.profiles[m.cursor]
			return m, tea.Quit
		case "e":
			m.selected = m.profiles[m.cursor]
			m.editing = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *ProfileSelector) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	s.WriteString(titleStyle.Render("Git Profiles"))
	s.WriteString("\n\n")

	for i, profile := range m.profiles {
		if i == m.cursor {
			s.WriteString(selectedItemStyle.Render(profile))
		} else {
			s.WriteString(itemStyle.Render(profile))
		}
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Foreground(subtle).Render("↑/↓: Navigate • enter: Select • e: Edit • q: Quit"))

	if m.err != nil {
		s.WriteString("\n\n")
		s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	}

	return s.String()
}

func (m *ProfileSelector) Selected() string {
	return m.selected
}

func (m *ProfileSelector) IsEditing() bool {
	return m.editing
}
