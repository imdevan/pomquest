package questing

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/pomquest/internal/ui"
)

// PartyGridComponent renders hero cards and turn indicator header.
type PartyGridComponent struct {
	theme  ui.Theme
	heroes []ui.HeroCardData
}

func NewPartyGridComponent(theme ui.Theme, heroes []ui.HeroCardData) PartyGridComponent {
	return PartyGridComponent{
		theme:  theme,
		heroes: heroes,
	}
}

func (c PartyGridComponent) View() string {
	turnIndicatorStyle := lipgloss.NewStyle().
		Foreground(c.theme.Primary).
		Bold(true)

	// Turn damage indicator above active hero card
	turnHeader := "      " + turnIndicatorStyle.Render("7!")

	var cards []string
	for _, h := range c.heroes {
		card := ui.NewHeroCardComponent(c.theme, h).View()
		cards = append(cards, card)
	}

	// 2 spaces in between cards to account for offset active card (Requirement 5.5.1)
	var spacerCards []string
	for i, card := range cards {
		spacerCards = append(spacerCards, card)
		if i < len(cards)-1 {
			spacerCards = append(spacerCards, "  ")
		}
	}

	grid := lipgloss.JoinHorizontal(lipgloss.Top, spacerCards...)
	return turnHeader + "\n" + grid
}
