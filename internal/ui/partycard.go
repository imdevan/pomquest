package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/table"
)

// PartyCardData holds resources and supplies info for the party card.
type PartyCardData struct {
	Gold           int
	Silver         int
	Copper         int
	FoodLbs        int
	Diet           string
	SpareParts     int
	Oxen           int
	Regents        int
	SickCount      int
	Pace           string
	MagicItems     int
	Gems           int
	TreasureChests int
}

// PartyCardComponent is a generic party resources card widget rendered with default lipgloss/table.
type PartyCardComponent struct {
	Data  PartyCardData
	Width int
	Theme Theme
}

func NewPartyCardComponent(theme Theme, data PartyCardData, width int) PartyCardComponent {
	if width <= 0 {
		width = 72
	}
	return PartyCardComponent{
		Data:  data,
		Width: width,
		Theme: theme,
	}
}

func (c PartyCardComponent) View() string {
	if c.Data.Pace == "" {
		c.Data.Pace = "fast"
	}
	if c.Data.Diet == "" {
		c.Data.Diet = "balanced"
	}

	coinStr := fmt.Sprintf("coins: %dg %ds %dc", c.Data.Gold, c.Data.Silver, c.Data.Copper)
	foodStr := fmt.Sprintf("food: %slbs", formatNumberWithCommas(c.Data.FoodLbs))
	oxStr := fmt.Sprintf("ox: %d", c.Data.Oxen)
	partsStr := fmt.Sprintf("spare parts: %d", c.Data.SpareParts)

	dietStr := fmt.Sprintf("diet: %s", c.Data.Diet)
	paceStr := fmt.Sprintf("pace: %s", c.Data.Pace)

	magicStr := fmt.Sprintf("magic items: %d", c.Data.MagicItems)
	gemsStr := fmt.Sprintf("gems: %d", c.Data.Gems)
	chestStr := fmt.Sprintf("treasure chests: %d", c.Data.TreasureChests)

	t := table.New().
		Headers("Party:", "", "").
		Rows(
			[]string{coinStr, dietStr, ""},
			[]string{foodStr, paceStr, magicStr},
			[]string{oxStr, "", gemsStr},
			[]string{partsStr, "", chestStr},
		)

	return t.Render()
}

func formatNumberWithCommas(n int) string {
	in := fmt.Sprintf("%d", n)
	out := make([]byte, 0, len(in)+(len(in)-1)/3)
	for i, c := range in {
		if i > 0 && (len(in)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, byte(c))
	}
	return string(out)
}
