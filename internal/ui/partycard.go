package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
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

// PartyCardComponent is a generic party resources card widget rendered using TexturedCard.
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

	horizontalPadding := 2
	paddingTopBottom := 1

	coinStr := fmt.Sprintf("coins: %dg %ds %dc", c.Data.Gold, c.Data.Silver, c.Data.Copper)
	foodStr := fmt.Sprintf("food: %slbs", formatNumberWithCommas(c.Data.FoodLbs))
	oxStr := fmt.Sprintf("ox: %d", c.Data.Oxen)
	partsStr := fmt.Sprintf("spare parts: %d", c.Data.SpareParts)

	dietStr := fmt.Sprintf("diet: %s", c.Data.Diet)
	paceStr := fmt.Sprintf("pace: %s", c.Data.Pace)

	magicStr := fmt.Sprintf("magic items: %d", c.Data.MagicItems)
	gemsStr := fmt.Sprintf("gems: %d", c.Data.Gems)
	chestStr := fmt.Sprintf("treasure chests: %d", c.Data.TreasureChests)

	// Inner table width takes into account outer borders (2) and horizontal padding (2 * horizontalPadding)
	innerTableWidth := c.Width - 2 - (2 * horizontalPadding)
	if innerTableWidth < 10 {
		innerTableWidth = 10
	}

	emptyBorder := lipgloss.Border{}
	t := table.New().
		Border(emptyBorder).
		BorderRow(false).
		BorderColumn(false).
		Width(innerTableWidth).
		Rows(
			[]string{coinStr, dietStr, ""},
			[]string{foodStr, paceStr, magicStr},
			[]string{oxStr, "", gemsStr},
			[]string{partsStr, "", chestStr},
		).
		StyleFunc(func(row, col int) lipgloss.Style {
			if col == 2 {
				return lipgloss.NewStyle().Foreground(c.Theme.Text).Align(lipgloss.Right)
			}
			return lipgloss.NewStyle().Foreground(c.Theme.Text).Align(lipgloss.Left)
		})

	tableContent := t.Render()

	// Wrap party table inside TexturedCard with TopLeft: "Party" header and fluent .Padding(0, 2)
	card := NewTexturedCard(c.Theme, c.Width, tableContent).Padding(paddingTopBottom, horizontalPadding)
	card.TopLeft = CornerText{
		Text:         "Party",
		Margin:       1,
		PaddingLeft:  1,
		PaddingRight: 1,
	}

	return card.View()
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
