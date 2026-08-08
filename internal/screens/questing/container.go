package questing

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pomquest/internal/ui"
)

const DefaultContainerWidth = 96

// QuestContainerComponent wraps quest screen components into a flex width container layout.
type QuestContainerComponent struct {
	Width      int
	Theme      ui.Theme
	Responsive *ui.ResponsiveManager
}

func NewQuestContainerComponent(theme ui.Theme, width int, responsive ...*ui.ResponsiveManager) QuestContainerComponent {
	if width <= 0 {
		width = DefaultContainerWidth
	}
	var rm *ui.ResponsiveManager
	if len(responsive) > 0 && responsive[0] != nil {
		rm = responsive[0]
	} else {
		rm = ui.NewResponsiveManager(width)
	}
	return QuestContainerComponent{
		Width:      width,
		Theme:      theme,
		Responsive: rm,
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
