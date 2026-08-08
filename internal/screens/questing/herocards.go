package questing

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pomquest/internal/ui"
)

// HeroCardsComponent renders party hero cards centered within the container width.
type HeroCardsComponent struct {
	theme      ui.Theme
	heroes     []ui.HeroCardData
	width      int
	responsive *ui.ResponsiveManager
}

func NewHeroCardsComponent(theme ui.Theme, heroes []ui.HeroCardData, width int, responsive ...*ui.ResponsiveManager) HeroCardsComponent {
	w := DefaultContainerWidth
	if width > 0 {
		w = width
	}
	var rm *ui.ResponsiveManager
	if len(responsive) > 0 && responsive[0] != nil {
		rm = responsive[0]
	} else {
		rm = ui.NewResponsiveManager(w)
	}
	return HeroCardsComponent{
		theme:      theme,
		heroes:     heroes,
		width:      w,
		responsive: rm,
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
		cardView := ui.NewHeroCardComponent(c.theme, h, c.responsive).View()
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
