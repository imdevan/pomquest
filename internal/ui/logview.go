package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// LogViewComponent is a generic auto-scrolling log stream widget.
type LogViewComponent struct {
	Title    string
	Messages []string
	MaxRows  int
	Theme    Theme
}

func NewLogViewComponent(theme Theme, title string, messages []string, maxRows int) LogViewComponent {
	if maxRows <= 0 {
		maxRows = 8
	}
	return LogViewComponent{
		Title:    title,
		Messages: messages,
		MaxRows:  maxRows,
		Theme:    theme,
	}
}

func (c LogViewComponent) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(c.Theme.Secondary).
		MarginBottom(1)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(c.Theme.Border).
		Padding(0, 1)

	start := 0
	if len(c.Messages) > c.MaxRows {
		start = len(c.Messages) - c.MaxRows
	}

	var sb strings.Builder
	if c.Title != "" {
		sb.WriteString(titleStyle.Render(c.Title) + "\n")
	}

	visible := c.Messages[start:]
	for _, msg := range visible {
		sb.WriteString("• " + msg + "\n")
	}

	return boxStyle.Render(strings.TrimRight(sb.String(), "\n"))
}
