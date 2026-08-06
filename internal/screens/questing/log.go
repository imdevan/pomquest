package questing

import "github.com/pomquest/internal/ui"

// LogComponent wraps ui.LogViewComponent for event scrollback.
type LogComponent struct {
	theme    ui.Theme
	messages []string
}

func NewLogComponent(theme ui.Theme, messages []string) LogComponent {
	return LogComponent{
		theme:    theme,
		messages: messages,
	}
}

func (c LogComponent) View() string {
	return ui.NewLogViewComponent(c.theme, "📜 DUNGEON LOG", c.messages, 7).View()
}
