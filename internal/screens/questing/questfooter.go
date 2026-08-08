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
	theme      ui.Theme
	data       QuestProgressData
	width      int
	responsive *ui.ResponsiveManager
}

func NewQuestFooterComponent(theme ui.Theme, data QuestProgressData, width int, responsive ...*ui.ResponsiveManager) QuestFooterComponent {
	if width <= 0 {
		width = DefaultContainerWidth
	}
	data.Width = width
	var rm *ui.ResponsiveManager
	if len(responsive) > 0 && responsive[0] != nil {
		rm = responsive[0]
	} else {
		rm = ui.NewResponsiveManager(width)
	}
	return QuestFooterComponent{
		theme:      theme,
		data:       data,
		width:      width,
		responsive: rm,
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

	// Line 1: Dest on left, Timer aligned fully to right of container width
	leftDest := fmt.Sprintf("Dest: %s", destStyle.Render(c.data.Destination))
	rightTime := timeStyle.Render(timeStr)
	destSpacing := c.width - lipgloss.Width(leftDest) - lipgloss.Width(rightTime)
	if destSpacing < 1 {
		destSpacing = 1
	}
	destLine := leftDest + strings.Repeat(" ", destSpacing) + rightTime

	// Line 2: Progress bar matching full width of container
	ratio := float64(c.data.CurrentSecs) / float64(c.data.TotalSecs)
	progBar := ui.NewProgressBarComponent(c.theme, ratio, c.width)
	barLine := progBar.View()

	// Line 3: Environment on left, Days passed aligned fully to right of container width
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

	rm := c.responsive
	if rm == nil {
		rm = ui.NewResponsiveManager(c.width)
	}

	var leftEnv string
	if !rm.IsAtLeast(ui.BreakpointLG) {
		// On MD or lower: hide "Weather:", "Terrain:", and "Hazards:" labels, keeping "|" separator
		leftEnv = fmt.Sprintf("%s | %s | %s", c.data.Weather, c.data.Terrain, c.data.Hazards)
	} else {
		// On LG+: full labels
		leftEnv = fmt.Sprintf("Weather: %s | Terrain: %s | Hazards: %s", c.data.Weather, c.data.Terrain, c.data.Hazards)
	}

	rightEnv := fmt.Sprintf("days passed: %d", c.data.DaysPassed)
	envSpacing := c.width - lipgloss.Width(leftEnv) - lipgloss.Width(rightEnv)
	if envSpacing < 1 {
		envSpacing = 1
	}
	envLine := leftEnv + strings.Repeat(" ", envSpacing) + rightEnv

	return strings.Join([]string{destLine, barLine, envLine}, "\n")
}
