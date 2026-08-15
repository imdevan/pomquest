package questing

import "github.com/pomquest/internal/ui"

// LogComponent wraps ui.LogViewComponent for the top dungeon event log.
type LogComponent struct {
	theme       ui.Theme
	date        string
	dungeonName string
	roomName    string
	messages    []string
	width       int
	pageHeight  int
}

func NewLogComponent(
	theme ui.Theme,
	date string,
	dungeonName string,
	roomName string,
	messages []string,
	width int,
	pageHeight ...int,
) LogComponent {
	if width <= 0 {
		width = DefaultContainerWidth
	}
	if date == "" {
		date = "12, Jan, 1000 A.P."
	}
	ph := 0
	if len(pageHeight) > 0 {
		ph = pageHeight[0]
	}
	return LogComponent{
		theme:       theme,
		date:        date,
		dungeonName: dungeonName,
		roomName:    roomName,
		messages:    messages,
		width:       width,
		pageHeight:  ph,
	}
}

func (c LogComponent) View() string {
	maxRows := 12
	if c.pageHeight > 0 {
		maxRows = c.pageHeight - 31
		if maxRows < 4 {
			maxRows = 4
		} else if maxRows > 15 {
			maxRows = 15
		}
	}
	logView := ui.LogViewComponent{
		Theme:      c.theme,
		TopTopLeft: c.date,
		TopLeft:    c.dungeonName,
		TopRight:   c.roomName,
		Messages:   c.messages,
		MaxRows:    maxRows,
		Width:      c.width,
	}
	return logView.View()
}
