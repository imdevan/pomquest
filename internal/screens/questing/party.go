package questing

import (
	"github.com/pomquest/internal/ui"
)

// PartyGridComponent wraps HeroCardsComponent for party display.
type PartyGridComponent struct {
	heroCards HeroCardsComponent
}

func NewPartyGridComponent(theme ui.Theme, heroes []ui.HeroCardData, width ...int) PartyGridComponent {
	return PartyGridComponent{
		heroCards: NewHeroCardsComponent(theme, heroes, width...),
	}
}

func (c PartyGridComponent) View() string {
	return c.heroCards.View()
}
