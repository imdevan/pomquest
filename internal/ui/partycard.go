package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/pomquest/internal/adapters/icon"
	"github.com/pomquest/internal/app"
)

// PartyCardData holds resources and supplies info for the party card.
type PartyCardData struct {
	Coins          int
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
	CurrentHP        int
	MaxHP            int
	CurrentPartySize int
	MaxPartySize     int
}

// PartyCardComponent is a generic party resources card widget rendered using TexturedCard.
type PartyCardComponent struct {
	Data       PartyCardData
	Width      int
	Theme      Theme
	Responsive *ResponsiveManager
}

func NewPartyCardComponent(theme Theme, data PartyCardData, width int, responsive ...*ResponsiveManager) PartyCardComponent {
	if width <= 0 {
		width = 72
	}
	var rm *ResponsiveManager
	if len(responsive) > 0 && responsive[0] != nil {
		rm = responsive[0]
	} else {
		rm = NewResponsiveManager(width)
	}
	return PartyCardComponent{
		Data:       data,
		Width:      width,
		Theme:      theme,
		Responsive: rm,
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

	rm := c.Responsive
	if rm == nil {
		rm = NewResponsiveManager(c.Width)
	}

	var foodStr, lootStr, partsStr, oxStr string
	var usageStr, dietStr, paceStr, greedStr string
	var potionsStr, bombsStr, shovelsStr, ropeStr string

	if !rm.IsAtLeast(BreakpointLG) {
		// On MD or lower: remove name from party status labels. Such that party card has icon: value
		foodStr = fmt.Sprintf("%s: %slbs", foodIcon, formatNumberWithCommas(c.Data.FoodLbs))
		lootStr = fmt.Sprintf("%s: %d", lootIcon, c.Data.Loot)
		partsStr = fmt.Sprintf("%s: %d", partsIcon, c.Data.SpareParts)
		oxStr = fmt.Sprintf("%s: %d", oxIcon, c.Data.Oxen)

		usageStr = fmt.Sprintf("%s: %s %s", itemUsageStatIcon, c.Data.ItemUsage, usageValIcon)
		dietStr = fmt.Sprintf("%s: %s %s", dietStatIcon, c.Data.Diet, dietValIcon)
		paceStr = fmt.Sprintf("%s: %s %s", paceStatIcon, c.Data.Pace, paceValIcon)
		greedStr = fmt.Sprintf("%s: %s %s", greedStatIcon, c.Data.Greed, greedValIcon)

		potionsStr = fmt.Sprintf("%s: %d", potionIcon, c.Data.Potions)
		bombsStr = fmt.Sprintf("%s: %d", bombIcon, c.Data.Bombs)
		shovelsStr = fmt.Sprintf("%s: %d", shovelIcon, c.Data.Shovels)
		ropeStr = fmt.Sprintf("%s: %d", ropeIcon, c.Data.Rope)
	} else {
		// On LG+: keep full label names
		foodStr = fmt.Sprintf("%s food: %slbs", foodIcon, formatNumberWithCommas(c.Data.FoodLbs))
		lootStr = fmt.Sprintf("%s loot: %d", lootIcon, c.Data.Loot)
		partsStr = fmt.Sprintf("%s spare parts: %d", partsIcon, c.Data.SpareParts)
		oxStr = fmt.Sprintf("%s ox: %d", oxIcon, c.Data.Oxen)

		usageStr = fmt.Sprintf("%s item usage: %s %s", itemUsageStatIcon, c.Data.ItemUsage, usageValIcon)
		dietStr = fmt.Sprintf("%s diet: %s %s", dietStatIcon, c.Data.Diet, dietValIcon)
		paceStr = fmt.Sprintf("%s pace: %s %s", paceStatIcon, c.Data.Pace, paceValIcon)
		greedStr = fmt.Sprintf("%s greed: %s %s", greedStatIcon, c.Data.Greed, greedValIcon)

		potionsStr = fmt.Sprintf("%s potions: %d", potionIcon, c.Data.Potions)
		bombsStr = fmt.Sprintf("%s bombs: %d", bombIcon, c.Data.Bombs)
		shovelsStr = fmt.Sprintf("%s shovels: %d", shovelIcon, c.Data.Shovels)
		ropeStr = fmt.Sprintf("%s rope: %d", ropeIcon, c.Data.Rope)
	}

	// Gold placed in bottom-right corner of TexturedCard container
	totalCoins := c.Data.Coins
	if totalCoins == 0 && (c.Data.Gold > 0 || c.Data.Silver > 0 || c.Data.Copper > 0) {
		totalCoins = c.Data.Gold*1000 + c.Data.Silver*100 + c.Data.Copper
	}
	goldStr := app.FormatCoins(totalCoins)

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

	// Wrap party table inside TexturedCard with TopLeft: "Party" header, TopRight: team hearts, and BottomRight: Gold
	card := NewTexturedCard(c.Theme, c.Width, tableContent).Padding(paddingTopBottom, horizontalPadding)
	card.TopLeft = CornerText{
		Text:         "Party",
		Margin:       1,
		PaddingLeft:  1,
		PaddingRight: 1,
	}

	maxHP := c.Data.MaxHP
	curHP := c.Data.CurrentHP
	if maxHP <= 0 {
		maxHP = 35
		curHP = 35
	}
	teamHeartsStr := RenderTeamHearts(curHP, maxHP, icon.ModeNerdFont)
	if teamHeartsStr != "" {
		card.TopRight = CornerText{
			Text:         teamHeartsStr,
			Margin:       1,
			PaddingLeft:  1,
			PaddingRight: 1,
		}
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
