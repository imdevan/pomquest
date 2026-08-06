package questing

import (
	"strings"

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
	hasActive := false
	for _, h := range c.heroes {
		if h.IsActive {
			hasActive = true
			break
		}
	}

	turnIndicatorStyle := lipgloss.NewStyle().
		Foreground(c.theme.Primary).
		Bold(true)

	var turnHeader string
	if hasActive {
		turnHeader = "      " + turnIndicatorStyle.Render("7!") + "\n"
	}

	var cards []string
	for _, h := range c.heroes {
		cardView := ui.NewHeroCardComponent(c.theme, h).View()
		if !hasActive {
			cardView = strings.TrimPrefix(cardView, "\n")
		}
		cards = append(cards, cardView)
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
	return turnHeader + grid
}
