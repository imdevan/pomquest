package ui

import "github.com/charmbracelet/lipgloss"

// HeaderComponent is a generic title & subtitle banner widget.
type HeaderComponent struct {
	Title    string
	Subtitle string
	Theme    Theme
}

func NewHeaderComponent(theme Theme, title, subtitle string) HeaderComponent {
	return HeaderComponent{
		Title:    title,
		Subtitle: subtitle,
		Theme:    theme,
	}
}

func (c HeaderComponent) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(c.Theme.Primary).
		MarginBottom(1)

	subTitleStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(c.Theme.Secondary)

	var res string
	if c.Title != "" {
		res += titleStyle.Render(c.Title) + "\n"
	}
	if c.Subtitle != "" {
		res += subTitleStyle.Render(c.Subtitle) + "\n"
	}
	return res
}
