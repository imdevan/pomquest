package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ProgressBarComponent is a generic progress/gauge bar widget.
type ProgressBarComponent struct {
	Label   string
	Current int
	Max     int
	Width   int
	Theme   Theme
}

func NewProgressBarComponent(theme Theme, label string, current, max, width int) ProgressBarComponent {
	if width <= 0 {
		width = 30
	}
	return ProgressBarComponent{
		Label:   label,
		Current: current,
		Max:     max,
		Width:   width,
		Theme:   theme,
	}
}

func (c ProgressBarComponent) View() string {
	if c.Max <= 0 {
		c.Max = 1
	}
	ratio := float64(c.Current) / float64(c.Max)
	if ratio > 1.0 {
		ratio = 1.0
	} else if ratio < 0 {
		ratio = 0
	}

	filledLen := int(ratio * float64(c.Width))
	emptyLen := c.Width - filledLen

	filledStyle := lipgloss.NewStyle().Foreground(c.Theme.Primary).Bold(true)
	emptyStyle := lipgloss.NewStyle().Foreground(c.Theme.Muted)
	labelStyle := lipgloss.NewStyle().Foreground(c.Theme.Text).Bold(true)

	bar := filledStyle.Render(strings.Repeat("█", filledLen)) + emptyStyle.Render(strings.Repeat("░", emptyLen))
	return fmt.Sprintf("%s [%s] %d/%d", labelStyle.Render(c.Label), bar, c.Current, c.Max)
}
