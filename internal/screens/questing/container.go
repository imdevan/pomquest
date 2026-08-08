package questing

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pomquest/internal/ui"
)

const DefaultContainerWidth = 96

// QuestContainerComponent wraps quest screen components into a flex width container layout.
type QuestContainerComponent struct {
	Width int
	Theme ui.Theme
}

func NewQuestContainerComponent(theme ui.Theme, width int) QuestContainerComponent {
	if width <= 0 {
		width = DefaultContainerWidth
	}
	return QuestContainerComponent{
		Width: width,
		Theme: theme,
	}
}

func (c QuestContainerComponent) Render(sections ...string) string {
	var sb strings.Builder
	for i, s := range sections {
		if s == "" {
			continue
		}
		sb.WriteString(s)
		if i < len(sections)-1 {
			sb.WriteString("\n\n")
		}
	}
	contentStyle := lipgloss.NewStyle().Width(c.Width).Align(lipgloss.Left)
	return contentStyle.Render(sb.String())
}
