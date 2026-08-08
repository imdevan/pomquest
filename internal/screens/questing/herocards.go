package questing

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pomquest/internal/ui"
)

// HeroCardsComponent renders party hero cards centered within the container width.
type HeroCardsComponent struct {
	theme  ui.Theme
	heroes []ui.HeroCardData
	width  int
}

func NewHeroCardsComponent(theme ui.Theme, heroes []ui.HeroCardData, width ...int) HeroCardsComponent {
	w := DefaultContainerWidth
	if len(width) > 0 && width[0] > 0 {
		w = width[0]
	}
	return HeroCardsComponent{
		theme:  theme,
		heroes: heroes,
		width:  w,
	}
}

func (c HeroCardsComponent) View() string {
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

	// Line 1: Damage section
	var damageLine string
	if hasActive {
		damageLine = "      " + turnIndicatorStyle.Render("7!")
	} else {
		damageLine = ""
	}

	// Line 2: Active card buffer space
	bufferLine := ""

	topHeader := damageLine + "\n" + bufferLine + "\n"

	var cards []string
	for _, h := range c.heroes {
		cardView := ui.NewHeroCardComponent(c.theme, h).View()
		if !hasActive {
			cardView = strings.TrimPrefix(cardView, "\n")
		}
		cards = append(cards, cardView)
	}

	// 2 spaces in between cards to account for offset active card
	var spacerCards []string
	for i, card := range cards {
		spacerCards = append(spacerCards, card)
		if i < len(cards)-1 {
			spacerCards = append(spacerCards, "  ")
		}
	}

	grid := lipgloss.JoinHorizontal(lipgloss.Top, spacerCards...)
	content := topHeader + grid

	if c.width > 0 {
		return lipgloss.NewStyle().Width(c.width).Align(lipgloss.Center).Render(content)
	}
	return content
}
