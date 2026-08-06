package questing

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/pomquest/internal/ui"
)

// PartyGridComponent renders hero cards side-by-side using ui.HeroCardComponent.
type PartyGridComponent struct {
	theme ui.Theme
	heroes []ui.HeroCardData
}

func NewPartyGridComponent(theme ui.Theme, heroes []ui.HeroCardData) PartyGridComponent {
	return PartyGridComponent{
		theme: theme,
		heroes: heroes,
	}
}

func (c PartyGridComponent) View() string {
	var cards []string
	for _, h := range c.heroes {
		card := ui.NewHeroCardComponent(c.theme, h).View()
		cards = append(cards, card)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, cards...)
}
