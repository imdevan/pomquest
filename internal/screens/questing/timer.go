package questing

import "github.com/pomquest/internal/ui"

// TimerComponent wraps ui.ProgressBarComponent for quest destination and pomodoro progress.
type TimerComponent struct {
	theme ui.Theme
	data  ui.QuestProgressData
}

func NewTimerComponent(theme ui.Theme, data ui.QuestProgressData) TimerComponent {
	return TimerComponent{
		theme: theme,
		data:  data,
	}
}

func (c TimerComponent) View() string {
	return ui.NewProgressBarComponent(c.theme, c.data).View()
}
