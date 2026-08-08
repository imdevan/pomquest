package ui

import "strings"

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

	card := NewTexturedCard(c.Theme, c.Width, content).Padding(0, 1)

	if c.Title != "" {
		card.TopLeft = CornerText{
			Text:         c.Title,
			Margin:       1,
			PaddingLeft:  1,
			PaddingRight: 1,
		}
	}

	return card.View()
}
