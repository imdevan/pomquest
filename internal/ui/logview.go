package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// LogViewComponent is a generic auto-scrolling log stream widget with rounded border.
type LogViewComponent struct {
	Title    string
	Messages []string
	MaxRows  int
	Width    int
	Theme    Theme
}

func NewLogViewComponent(theme Theme, title string, messages []string, maxRows, width int) LogViewComponent {
	if maxRows <= 0 {
		maxRows = 8
	}
	if width <= 0 {
		width = 72
	}
	return LogViewComponent{
		Title:    title,
		Messages: messages,
		MaxRows:  maxRows,
		Width:    width,
		Theme:    theme,
	}
}

func (c LogViewComponent) View() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(c.Theme.Border).
		Padding(0, 1).
		Width(c.Width)

	start := 0
	if len(c.Messages) > c.MaxRows {
		start = len(c.Messages) - c.MaxRows
	}

	var sb strings.Builder
	visible := c.Messages[start:]
	for _, msg := range visible {
		sb.WriteString(msg + "\n")
	}

	content := strings.TrimRight(sb.String(), "\n")
	
	// Apply custom title top border if title present
	if c.Title != "" {
		boxStyle = boxStyle.BorderTop(true)
	}

	return boxStyle.Render(content)
}
