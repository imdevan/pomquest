package questing

import "github.com/pomquest/internal/ui"

// LogComponent wraps ui.LogViewComponent for the top dungeon event log.
type LogComponent struct {
	theme      ui.Theme
	title      string
	messages   []string
	width      int
	pageHeight int
}

func NewLogComponent(theme ui.Theme, title string, messages []string, width int, pageHeight ...int) LogComponent {
	if width <= 0 {
		width = DefaultContainerWidth
	}
	if title == "" {
		title = "12, Jan, 1000 A.P."
	}
	ph := 0
	if len(pageHeight) > 0 {
		ph = pageHeight[0]
	}
	return LogComponent{
		theme:      theme,
		title:      title,
		messages:   messages,
		width:      width,
		pageHeight: ph,
	}
}

func (c LogComponent) View() string {
	maxRows := 10
	if c.pageHeight > 0 {
		maxRows = c.pageHeight - 31
		if maxRows < 6 {
			maxRows = 6
		} else if maxRows > 14 {
			maxRows = 14
		}
	}
	return ui.NewLogViewComponent(c.theme, c.title, c.messages, maxRows, c.width).View()
}
