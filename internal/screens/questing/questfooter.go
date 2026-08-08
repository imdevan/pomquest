package questing

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pomquest/internal/adapters/icon"
	"github.com/pomquest/internal/ui"
)

// QuestProgressData holds destination, progress timer, and environmental info.
type QuestProgressData struct {
	Destination string
	CurrentSecs int
	TotalSecs   int
	Weather     string
	Terrain     string
	Hazards     string
	DaysPassed  int
	Width       int
}

// QuestFooterComponent renders destination, timer, progress bar, and environmental info.
type QuestFooterComponent struct {
	theme ui.Theme
	data  QuestProgressData
	width int
}

func NewQuestFooterComponent(theme ui.Theme, data QuestProgressData, width int) QuestFooterComponent {
	if width <= 0 {
		width = DefaultContainerWidth
	}
	data.Width = width
	return QuestFooterComponent{
		theme: theme,
		data:  data,
		width: width,
	}
}

func (c QuestFooterComponent) View() string {
	if c.data.TotalSecs <= 0 {
		c.data.TotalSecs = 1500
	}

	curMins := c.data.CurrentSecs / 60
	curSecs := c.data.CurrentSecs % 60
	totMins := c.data.TotalSecs / 60
	totSecs := c.data.TotalSecs % 60

	timeStr := fmt.Sprintf("%02d:%02d/%02d:%02d", curMins, curSecs, totMins, totSecs)

	destStyle := lipgloss.NewStyle().Bold(true).Foreground(c.theme.Primary)
	timeStyle := lipgloss.NewStyle().Foreground(c.theme.Secondary)

	// Line 1: Dest & Timer flexed to width
	destPadding := c.width - len("  Dest: ") - len(timeStr) - 3
	if destPadding < 10 {
		destPadding = 10
	}
	destLine := fmt.Sprintf("  Dest: %-*s %s", destPadding, destStyle.Render(c.data.Destination), timeStyle.Render(timeStr))

	// Line 2: Progress bar using generic ui.ProgressBarComponent
	barWidth := c.width - 6
	if barWidth < 10 {
		barWidth = 10
	}
	ratio := float64(c.data.CurrentSecs) / float64(c.data.TotalSecs)

	progBar := ui.NewProgressBarComponent(c.theme, ratio, barWidth)
	barLine := "  " + progBar.View()

	// Line 3: Environment & Days flexed to width
	if c.data.Weather == "" {
		c.data.Weather = icon.GetNerdFont(icon.KeyWeather)
	}
	if c.data.Terrain == "" {
		c.data.Terrain = icon.GetNerdFont(icon.KeyTerrain)
	}
	if c.data.Hazards == "" {
		c.data.Hazards = "mud, rain"
	}
	if c.data.DaysPassed <= 0 {
		c.data.DaysPassed = 12
	}

	hazardsPadding := c.width - 45
	if hazardsPadding < 10 {
		hazardsPadding = 10
	}

	envLine := fmt.Sprintf("  Weather: %s | Terrain: %s | Hazards: %-*s days passed: %d",
		c.data.Weather, c.data.Terrain, hazardsPadding, c.data.Hazards, c.data.DaysPassed)

	return strings.Join([]string{destLine, barLine, envLine}, "\n")
}
