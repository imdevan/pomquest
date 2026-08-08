package questing

import "github.com/pomquest/internal/ui"

// TimerComponent wraps ui.ProgressBarComponent for quest destination and pomodoro progress.
type TimerComponent struct {
	theme ui.Theme
	data  ui.QuestProgressData
	width int
}

func NewTimerComponent(theme ui.Theme, data ui.QuestProgressData, width int) TimerComponent {
	if width <= 0 {
		width = DefaultContainerWidth
	}
	data.Width = width
	return TimerComponent{
		theme: theme,
		data:  data,
		width: width,
	}
}

func (c TimerComponent) View() string {
	return ui.NewProgressBarComponent(c.theme, c.data).View()
}
