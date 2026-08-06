package questing

import "github.com/pomquest/internal/ui"

// TimerComponent wraps ui.ProgressBarComponent for quest pomodoro progress.
type TimerComponent struct {
	theme       ui.Theme
	label       string
	currentSecs int
	totalSecs   int
}

func NewTimerComponent(theme ui.Theme, label string, currentSecs, totalSecs int) TimerComponent {
	return TimerComponent{
		theme:       theme,
		label:       label,
		currentSecs: currentSecs,
		totalSecs:   totalSecs,
	}
}

func (c TimerComponent) View() string {
	return ui.NewProgressBarComponent(c.theme, c.label, c.currentSecs, c.totalSecs, 30).View()
}
