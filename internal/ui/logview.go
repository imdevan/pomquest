package ui

import "strings"

// LogViewComponent is a generic auto-scrolling log stream widget with rounded border.
type LogViewComponent struct {
	Title        string
	TopTopLeft   string
	TopTopRight  string
	TopLeft      string
	TopRight     string
	Messages     []string
	MaxRows      int
	Width        int
	Theme        Theme
	ScrollOffset int
}

func NewLogViewComponent(theme Theme, title string, messages []string, maxRows, width int) LogViewComponent {
	if maxRows <= 0 {
		maxRows = 12
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
	availWidth := c.Width - 4
	if availWidth <= 0 {
		availWidth = 68
	}

	var allWrapped []string
	for _, msg := range c.Messages {
		lines := strings.Split(msg, "\n")
		for _, l := range lines {
			wrapped := WrapLine(l, availWidth)
			allWrapped = append(allWrapped, wrapped...)
		}
	}

	maxScroll := len(allWrapped) - c.MaxRows
	if maxScroll < 0 {
		maxScroll = 0
	}

	offset := c.ScrollOffset
	if offset < 0 {
		offset = 0
	}
	if offset > maxScroll {
		offset = maxScroll
	}

	start := maxScroll - offset
	end := start + c.MaxRows
	if end > len(allWrapped) {
		end = len(allWrapped)
	}

	var visible []string
	if len(allWrapped) > 0 && start < len(allWrapped) {
		visible = append(visible, allWrapped[start:end]...)
	}

	for len(visible) < c.MaxRows {
		visible = append(visible, "")
	}

	var sb strings.Builder
	for i, msg := range visible {
		sb.WriteString(msg)
		if i < len(visible)-1 {
			sb.WriteString("\n")
		}
	}

	content := sb.String()

	card := NewTexturedCard(c.Theme, c.Width, content).Padding(0, 1)

	if c.TopTopLeft != "" {
		card.TopTopLeft = CornerText{
			Text:         c.TopTopLeft,
			Margin:       1,
			PaddingLeft:  1,
			PaddingRight: 1,
		}
	}
	if c.TopTopRight != "" {
		card.TopTopRight = CornerText{
			Text:         c.TopTopRight,
			Margin:       1,
			PaddingLeft:  1,
			PaddingRight: 1,
		}
	}

	if c.TopLeft != "" {
		card.TopLeft = CornerText{
			Text:         c.TopLeft,
			Margin:       1,
			PaddingLeft:  1,
			PaddingRight: 1,
		}
	} else if c.Title != "" {
		card.TopLeft = CornerText{
			Text:         c.Title,
			Margin:       1,
			PaddingLeft:  1,
			PaddingRight: 1,
		}
	}

	if c.TopRight != "" {
		card.TopRight = CornerText{
			Text:         c.TopRight,
			Margin:       1,
			PaddingLeft:  1,
			PaddingRight: 1,
		}
	}

	return card.View()
}
