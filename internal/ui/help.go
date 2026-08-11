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
	Width int
}

func NewHelpComponent(theme Theme, keys []key.Binding) HelpComponent {
	return HelpComponent{
		Keys:  keys,
		Theme: theme,
	}
}

func NewHelpComponentSized(theme Theme, keys []key.Binding, width int) HelpComponent {
	return HelpComponent{
		Keys:  keys,
		Theme: theme,
		Width: width,
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

	maxWidth := c.Width
	if maxWidth <= 0 {
		return helpStyle.Render(strings.Join(helpStrs, "  •  "))
	}

	var lines []string
	var currentLine []string
	currentLen := 0

	for _, item := range helpStrs {
		itemLen := lipgloss.Width(item)
		if len(currentLine) == 0 {
			currentLine = append(currentLine, item)
			currentLen = itemLen
		} else if currentLen+5+itemLen <= maxWidth {
			currentLine = append(currentLine, item)
			currentLen += 5 + itemLen
		} else {
			lines = append(lines, strings.Join(currentLine, "  •  "))
			currentLine = []string{item}
			currentLen = itemLen
		}
	}
	if len(currentLine) > 0 {
		lines = append(lines, strings.Join(currentLine, "  •  "))
	}

	return helpStyle.Render(strings.Join(lines, "\n"))
}
