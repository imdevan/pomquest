package ui

import (
	"github.com/charmbracelet/bubbles/progress"
)

// ProgressBarComponent is a generic progress/gauge bar widget.
type ProgressBarComponent struct {
	Percent float64
	Width   int
	Theme   Theme
}

func NewProgressBarComponent(theme Theme, percent float64, width int) ProgressBarComponent {
	if width <= 0 {
		width = 72
	}
	if percent < 0 {
		percent = 0
	} else if percent > 1.0 {
		percent = 1.0
	}
	return ProgressBarComponent{
		Percent: percent,
		Width:   width,
		Theme:   theme,
	}
}

func (c ProgressBarComponent) View() string {
	barWidth := c.Width
	if barWidth < 1 {
		barWidth = 10
	}

	prog := progress.New(
		progress.WithoutPercentage(),
		progress.WithSolidFill(string(c.Theme.Primary)),
	)
	prog.Width = barWidth

	return prog.ViewAs(c.Percent)
}
