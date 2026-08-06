package questing

import "github.com/pomquest/internal/ui"

// LogComponent wraps ui.LogViewComponent for the top dungeon event log.
type LogComponent struct {
	theme    ui.Theme
	title    string
	messages []string
}

func NewLogComponent(theme ui.Theme, title string, messages []string) LogComponent {
	if title == "" {
		title = "12, Jan, 1000 A.P."
	}
	return LogComponent{
		theme:    theme,
		title:    title,
		messages: messages,
	}
}

func (c LogComponent) View() string {
	return ui.NewLogViewComponent(c.theme, c.title, c.messages, 8, 72).View()
}
