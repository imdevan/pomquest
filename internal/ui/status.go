package ui

import "github.com/charmbracelet/lipgloss"

// StatusComponent is a generic status line message banner widget.
type StatusComponent struct {
	Text  string
	Theme Theme
}

func NewStatusComponent(theme Theme, text string) StatusComponent {
	return StatusComponent{
		Text:  text,
		Theme: theme,
	}
}

func (c StatusComponent) View() string {
	statusStyle := lipgloss.NewStyle().
		Foreground(c.Theme.Muted).
		MarginBottom(1)

	return statusStyle.Render(c.Text) + "\n"
}
