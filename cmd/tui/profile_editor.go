package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Field struct {
	name     string
	value    string
	editable bool
}

type ProfileEditor struct {
	profileName string
	fields      []Field
	cursor      int
	editing     bool
	quitting    bool
	err         error
	saved       bool
}

func NewProfileEditor(name, email, gpgKey, sshKey string, signCommits bool) *ProfileEditor {
	fields := []Field{
		{name: "Name", value: name, editable: true},
		{name: "Email", value: email, editable: true},
		{name: "GPG Key", value: gpgKey, editable: true},
		{name: "SSH Key", value: sshKey, editable: true},
		{name: "Sign Commits", value: fmt.Sprintf("%v", signCommits), editable: true},
	}

	return &ProfileEditor{
		fields: fields,
	}
}

func (m *ProfileEditor) Init() tea.Cmd {
	return nil
}

func (m *ProfileEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.editing {
				m.editing = false
				return m, nil
			}
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if !m.editing && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if !m.editing && m.cursor < len(m.fields)-1 {
				m.cursor++
			}
		case "enter":
			if !m.editing {
				if m.fields[m.cursor].editable {
					m.editing = true
				}
				return m, nil
			}
			m.editing = false
		case "backspace":
			if m.editing && len(m.fields[m.cursor].value) > 0 {
				m.fields[m.cursor].value = m.fields[m.cursor].value[:len(m.fields[m.cursor].value)-1]
			}
		case "s":
			if !m.editing && msg.Alt {
				m.saved = true
				return m, tea.Quit
			}
		default:
			if m.editing && msg.Type == tea.KeyRunes {
				m.fields[m.cursor].value += string(msg.Runes)
			}
		}
	}

	return m, nil
}

func (m *ProfileEditor) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	s.WriteString(titleStyle.Render("Edit Profile"))
	s.WriteString("\n\n")

	for i, field := range m.fields {
		var fieldStyle lipgloss.Style
		if i == m.cursor {
			if m.editing {
				fieldStyle = selectedItemStyle.Copy().
					Foreground(lipgloss.Color("#FF00FF")).
					SetString("* ")
			} else {
				fieldStyle = selectedItemStyle
			}
		} else {
			fieldStyle = itemStyle
		}

		s.WriteString(fieldStyle.Render(fmt.Sprintf("%s: %s", field.name, field.value)))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	if m.editing {
		s.WriteString(lipgloss.NewStyle().Foreground(subtle).Render("enter: Save field • esc: Cancel editing"))
	} else {
		s.WriteString(lipgloss.NewStyle().Foreground(subtle).Render("↑/↓: Navigate • enter: Edit field • alt+s: Save profile • esc: Cancel"))
	}

	if m.err != nil {
		s.WriteString("\n\n")
		s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	}

	return s.String()
}

func (m *ProfileEditor) GetFields() map[string]string {
	result := make(map[string]string)
	for _, field := range m.fields {
		result[field.name] = field.value
	}
	return result
}

func (m *ProfileEditor) IsSaved() bool {
	return m.saved
}
