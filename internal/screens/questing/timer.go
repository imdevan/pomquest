package questing

import "github.com/pomquest/internal/ui"

// TimerComponent wraps QuestFooterComponent for quest destination and pomodoro progress.
type TimerComponent struct {
	footer QuestFooterComponent
}

func NewTimerComponent(theme ui.Theme, data QuestProgressData, width int) TimerComponent {
	return TimerComponent{
		footer: NewQuestFooterComponent(theme, data, width),
	}
}

func (c TimerComponent) View() string {
	return c.footer.View()
}
