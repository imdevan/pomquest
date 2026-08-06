package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// PartyBarData holds resources and supplies info for the party.
type PartyBarData struct {
	Gold       int
	Silver     int
	Copper     int
	FoodLbs    int
	Diet       string
	SpareParts int
	Oxen       int
	Regents    int
	SickCount  int
	Pace       string
}

// PartyBarComponent is a generic party resources bar widget matching mockup design.
type PartyBarComponent struct {
	Data  PartyBarData
	Width int
	Theme Theme
}

func NewPartyBarComponent(theme Theme, data PartyBarData, width int) PartyBarComponent {
	if width <= 0 {
		width = 72
	}
	return PartyBarComponent{
		Data:  data,
		Width: width,
		Theme: theme,
	}
}

func (c PartyBarComponent) View() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(c.Theme.Border).
		Width(c.Width)

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(c.Theme.Primary)
	textStyle := lipgloss.NewStyle().Foreground(c.Theme.Text)
	goldStyle := lipgloss.NewStyle().Foreground(c.Theme.Tags).Bold(true)
	mutedStyle := lipgloss.NewStyle().Foreground(c.Theme.Muted)

	if c.Data.Pace == "" {
		c.Data.Pace = "fast"
	}
	if c.Data.Diet == "" {
		c.Data.Diet = "balanced"
	}

	header := headerStyle.Render(" Party:") + strings.Repeat(" ", c.Width-16) + mutedStyle.Render("x0 ")
	divider := mutedStyle.Render("├" + strings.Repeat("─", c.Width-2) + "┤")

	line1 := fmt.Sprintf("  %-30s %-25s %s",
		goldStyle.Render(fmt.Sprintf("%dg %ds %dc", c.Data.Gold, c.Data.Silver, c.Data.Copper)),
		textStyle.Render(fmt.Sprintf("sick: %d", c.Data.SickCount)),
		mutedStyle.Render(""),
	)

	line2 := fmt.Sprintf("  %-30s %-25s %s",
		textStyle.Render(fmt.Sprintf("food: %dlbs", c.Data.FoodLbs)),
		textStyle.Render(fmt.Sprintf("diet: %s", c.Data.Diet)),
		mutedStyle.Render("x1"),
	)

	line3 := fmt.Sprintf("  %-30s %-25s %s",
		textStyle.Render(fmt.Sprintf("spare parts: %d", c.Data.SpareParts)),
		textStyle.Render(fmt.Sprintf("pace: %s", c.Data.Pace)),
		mutedStyle.Render("x0"),
	)

	line4 := fmt.Sprintf("  %-30s %-25s %s",
		textStyle.Render(fmt.Sprintf("ox: %d", c.Data.Oxen)),
		textStyle.Render(fmt.Sprintf("regents: %d", c.Data.Regents)),
		mutedStyle.Render("x0"),
	)

	content := strings.Join([]string{header, divider, line1, line2, line3, line4}, "\n")
	return boxStyle.Render(content)
}
