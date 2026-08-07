package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/pomquest/internal/adapters/icon"
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
	Loot           int
	Rope           int
	Bombs          int
	Potions        int
	Shovels        int
	ItemUsage      string
	Greed          string
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
	if c.Data.ItemUsage == "" {
		c.Data.ItemUsage = "none"
	}
	if c.Data.Greed == "" {
		c.Data.Greed = "none"
	}

	horizontalPadding := 2
	paddingTopBottom := 1

	foodIcon := icon.GetNerdFont(icon.KeyFood)
	lootIcon := icon.GetNerdFont(icon.KeyLoot)
	partsIcon := icon.GetNerdFont(icon.KeySpareParts)
	oxIcon := icon.GetNerdFont(icon.KeyOx)

	itemUsageStatIcon := icon.GetNerdFont(icon.KeyItemUsage)
	dietStatIcon := icon.GetNerdFont(icon.KeyDiet)
	paceStatIcon := icon.GetNerdFont(icon.KeyPace)
	greedStatIcon := icon.GetNerdFont(icon.KeyGreed)

	usageValIcon := icon.GetItemUsageIcon(c.Data.ItemUsage, icon.ModeNerdFont)
	dietValIcon := icon.GetDietIcon(c.Data.Diet, icon.ModeNerdFont)
	paceValIcon := icon.GetPaceIcon(c.Data.Pace, icon.ModeNerdFont)
	greedValIcon := icon.GetGreedIcon(c.Data.Greed, icon.ModeNerdFont)

	potionIcon := icon.GetNerdFont(icon.KeyPotion)
	bombIcon := icon.GetNerdFont(icon.KeyBomb)
	shovelIcon := icon.GetNerdFont(icon.KeyShovel)
	ropeIcon := icon.GetNerdFont(icon.KeyRope)

	// Column 1: food, loot, spare parts, ox
	foodStr := fmt.Sprintf("%s food: %slbs", foodIcon, formatNumberWithCommas(c.Data.FoodLbs))
	lootStr := fmt.Sprintf("%s loot: %d", lootIcon, c.Data.Loot)
	partsStr := fmt.Sprintf("%s spare parts: %d", partsIcon, c.Data.SpareParts)
	oxStr := fmt.Sprintf("%s ox: %d", oxIcon, c.Data.Oxen)

	// Column 2: item usage, diet, pace, greed
	usageStr := fmt.Sprintf("%s item usage: %s %s", itemUsageStatIcon, c.Data.ItemUsage, usageValIcon)
	dietStr := fmt.Sprintf("%s diet: %s %s", dietStatIcon, c.Data.Diet, dietValIcon)
	paceStr := fmt.Sprintf("%s pace: %s %s", paceStatIcon, c.Data.Pace, paceValIcon)
	greedStr := fmt.Sprintf("%s greed: %s %s", greedStatIcon, c.Data.Greed, greedValIcon)

	// Column 3: potions, bombs, shovels, rope
	potionsStr := fmt.Sprintf("%s potions: %d", potionIcon, c.Data.Potions)
	bombsStr := fmt.Sprintf("%s bombs: %d", bombIcon, c.Data.Bombs)
	shovelsStr := fmt.Sprintf("%s shovels: %d", shovelIcon, c.Data.Shovels)
	ropeStr := fmt.Sprintf("%s rope: %d", ropeIcon, c.Data.Rope)

	// Gold placed in bottom-right corner of TexturedCard container
	goldStr := fmt.Sprintf("%dg %ds %dc", c.Data.Gold, c.Data.Silver, c.Data.Copper)

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
			[]string{foodStr, usageStr, potionsStr},
			[]string{lootStr, dietStr, bombsStr},
			[]string{partsStr, paceStr, shovelsStr},
			[]string{oxStr, greedStr, ropeStr},
		).
		StyleFunc(func(row, col int) lipgloss.Style {
			if col == 2 {
				return lipgloss.NewStyle().Foreground(c.Theme.Text).Align(lipgloss.Right)
			}
			return lipgloss.NewStyle().Foreground(c.Theme.Text).Align(lipgloss.Left)
		})

	tableContent := t.Render()

	// Wrap party table inside TexturedCard with TopLeft: "Party" header and BottomRight: Gold
	card := NewTexturedCard(c.Theme, c.Width, tableContent).Padding(paddingTopBottom, horizontalPadding)
	card.TopLeft = CornerText{
		Text:         "Party",
		Margin:       1,
		PaddingLeft:  1,
		PaddingRight: 1,
	}
	card.BottomRight = CornerText{
		Text:         goldStr,
		Margin:       1,
		PaddingLeft:  1,
		PaddingRight: 1,
		Style:        lipgloss.NewStyle().Bold(true).Foreground(c.Theme.Tags),
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
