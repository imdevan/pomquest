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
	if c.Responsive != nil && c.Responsive.IsAtMost(ui.BreakpointSM) {
		msgStyle := lipgloss.NewStyle().
			Foreground(c.Theme.Text).
			Align(lipgloss.Center)

		w := c.Width
		if c.Responsive.Width() > 0 {
			w = c.Responsive.Width()
		}

		paneStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(c.Theme.Border).
			Padding(1, 2).
			Align(lipgloss.Center)

		if w > 4 {
			paneStyle = paneStyle.Width(w - 2)
		}

		return paneStyle.Render(msgStyle.Render("This terminal width is not currently supported."))
	}

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
