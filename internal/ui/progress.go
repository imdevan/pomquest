package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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

// ProgressBarComponent is a generic progress/gauge bar widget matching mockup design.
type ProgressBarComponent struct {
	Data  QuestProgressData
	Theme Theme
}

func NewProgressBarComponent(theme Theme, data QuestProgressData) ProgressBarComponent {
	if data.Width <= 0 {
		data.Width = 72
	}
	return ProgressBarComponent{
		Data:  data,
		Theme: theme,
	}
}

func (c ProgressBarComponent) View() string {
	if c.Data.TotalSecs <= 0 {
		c.Data.TotalSecs = 1500
	}

	curMins := c.Data.CurrentSecs / 60
	curSecs := c.Data.CurrentSecs % 60
	totMins := c.Data.TotalSecs / 60
	totSecs := c.Data.TotalSecs % 60

	timeStr := fmt.Sprintf("%02d:%02d/%02d:%02d", curMins, curSecs, totMins, totSecs)

	destStyle := lipgloss.NewStyle().Bold(true).Foreground(c.Theme.Primary)
	timeStyle := lipgloss.NewStyle().Foreground(c.Theme.Secondary)
	mutedStyle := lipgloss.NewStyle().Foreground(c.Theme.Muted)

	// Line 1: Dest & Timer
	destLine := fmt.Sprintf("  Dest: %-42s %s", destStyle.Render(c.Data.Destination), timeStyle.Render(timeStr))

	// Line 2: Progress bar
	barWidth := c.Data.Width - 6
	ratio := float64(c.Data.CurrentSecs) / float64(c.Data.TotalSecs)
	if ratio > 1.0 {
		ratio = 1.0
	} else if ratio < 0 {
		ratio = 0
	}

	filledLen := int(ratio * float64(barWidth))
	emptyLen := barWidth - filledLen

	filledStyle := lipgloss.NewStyle().Foreground(c.Theme.Primary)
	emptyStyle := lipgloss.NewStyle().Foreground(c.Theme.Muted)

	barLine := "  " + filledStyle.Render(strings.Repeat("█", filledLen)) + emptyStyle.Render(strings.Repeat(" ", emptyLen)) + mutedStyle.Render("]")

	// Line 3: Environment & Days
	if c.Data.Weather == "" {
		c.Data.Weather = "☀️"
	}
	if c.Data.Terrain == "" {
		c.Data.Terrain = "🌲"
	}
	if c.Data.Hazards == "" {
		c.Data.Hazards = "mud, rain"
	}
	if c.Data.DaysPassed <= 0 {
		c.Data.DaysPassed = 12
	}

	envLine := fmt.Sprintf("  Weather: %s | Terrain: %s | Hazards: %-18s days passed: %d",
		c.Data.Weather, c.Data.Terrain, c.Data.Hazards, c.Data.DaysPassed)

	return strings.Join([]string{destLine, barLine, envLine}, "\n")
}
