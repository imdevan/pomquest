package questing

import "github.com/pomquest/internal/ui"

// LogComponent wraps ui.LogViewComponent for the top dungeon event log.
type LogComponent struct {
	theme        ui.Theme
	date         string
	dungeonName  string
	roomName     string
	messages     []string
	width        int
	pageHeight   int
	scrollOffset int
}

func NewLogComponent(
	theme ui.Theme,
	date string,
	dungeonName string,
	roomName string,
	messages []string,
	width int,
	pageHeight int,
	scrollOffset ...int,
) LogComponent {
	if width <= 0 {
		width = DefaultContainerWidth
	}
	if date == "" {
		date = "12, Jan, 1000 A.P."
	}
	so := 0
	if len(scrollOffset) > 0 {
		so = scrollOffset[0]
	}
	return LogComponent{
		theme:        theme,
		date:         date,
		dungeonName:  dungeonName,
		roomName:     roomName,
		messages:     messages,
		width:        width,
		pageHeight:   pageHeight,
		scrollOffset: so,
	}
}

func (c LogComponent) View() string {
	maxRows := 6
	if c.pageHeight >= 43 {
		maxRows = c.pageHeight - 33
		if maxRows > 15 {
			maxRows = 15
		}
	} else if c.pageHeight > 0 {
		maxRows = c.pageHeight - 32
		if maxRows < 2 {
			maxRows = 2
		} else if maxRows > 15 {
			maxRows = 15
		}
	}
	logView := ui.LogViewComponent{
		Theme:        c.theme,
		TopTopLeft:   c.date,
		TopLeft:      c.dungeonName,
		TopRight:     c.roomName,
		Messages:     c.messages,
		MaxRows:      maxRows,
		Width:        c.width,
		ScrollOffset: c.scrollOffset,
	}
	return logView.View()
}
