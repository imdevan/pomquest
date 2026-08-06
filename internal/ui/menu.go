package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// MenuComponent is a generic choice list selector widget.
type MenuComponent struct {
	Choices []string
	Cursor  int
	Theme   Theme
}

func NewMenuComponent(theme Theme, choices []string, cursor int) MenuComponent {
	return MenuComponent{
		Choices: choices,
		Cursor:  cursor,
		Theme:   theme,
	}
}

func (c MenuComponent) View() string {
	cursorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(c.Theme.Primary)

	normalStyle := lipgloss.NewStyle().
		Foreground(c.Theme.Text)

	var sb strings.Builder
	for i, choice := range c.Choices {
		if c.Cursor == i {
			sb.WriteString(cursorStyle.Render(fmt.Sprintf("▶ %s", choice)) + "\n")
		} else {
			sb.WriteString(normalStyle.Render(fmt.Sprintf("  %s", choice)) + "\n")
		}
	}
	return sb.String()
}
