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
}

// PartyBarComponent is a generic party resources bar widget.
type PartyBarComponent struct {
	Data  PartyBarData
	Theme Theme
}

func NewPartyBarComponent(theme Theme, data PartyBarData) PartyBarComponent {
	return PartyBarComponent{
		Data:  data,
		Theme: theme,
	}
}

func (c PartyBarComponent) View() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(c.Theme.Border).
		Padding(0, 1)

	textStyle := lipgloss.NewStyle().Foreground(c.Theme.Text)
	goldStyle := lipgloss.NewStyle().Foreground(c.Theme.Tags).Bold(true)

	coinStr := goldStyle.Render(fmt.Sprintf("🪙 %dg %ds %dc", c.Data.Gold, c.Data.Silver, c.Data.Copper))
	foodStr := textStyle.Render(fmt.Sprintf("🌾 food: %dlbs (%s)", c.Data.FoodLbs, c.Data.Diet))
	suppliesStr := textStyle.Render(fmt.Sprintf("🐂 oxen: %d |  parts: %d", c.Data.Oxen, c.Data.SpareParts))

	return boxStyle.Render(strings.Join([]string{coinStr, foodStr, suppliesStr}, "  •  "))
}
