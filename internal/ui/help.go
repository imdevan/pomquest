package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

// HelpComponent is a generic keybindings footer widget.
type HelpComponent struct {
	Keys  []key.Binding
	Theme Theme
}

func NewHelpComponent(theme Theme, keys []key.Binding) HelpComponent {
	return HelpComponent{
		Keys:  keys,
		Theme: theme,
	}
}

func (c HelpComponent) View() string {
	helpStyle := lipgloss.NewStyle().
		Foreground(c.Theme.Muted).
		MarginTop(1)

	var helpStrs []string
	for _, k := range c.Keys {
		helpStrs = append(helpStrs, fmt.Sprintf("%s: %s", k.Help().Key, k.Help().Desc))
	}
	return helpStyle.Render(strings.Join(helpStrs, "  •  "))
}
