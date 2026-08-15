package questing

import (
	"fmt"
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
	rm := c.responsive
	if rm == nil {
		rm = ui.NewResponsiveManager(c.width)
	}

	cardWidth := 21
	if !rm.IsAtLeast(ui.BreakpointLG) {
		cardWidth = 13
	}
	spacerWidth := 2

	critStyle := lipgloss.NewStyle().Foreground(c.theme.Primary).Bold(true)
	atkStyle := lipgloss.NewStyle().Foreground(c.theme.Primary).Bold(true)
	blockStyle := lipgloss.NewStyle().Foreground(c.theme.Secondary).Bold(true)
	dodgeStyle := lipgloss.NewStyle().Foreground(c.theme.Flags).Bold(true)
	takeDmgStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("01")).Bold(true)

	var damageElements []string
	for i, h := range c.heroes {
		var indicator string
		if h.AnimState.Type != ui.AnimNone {
			switch h.AnimState.Type {
			case ui.AnimAttack:
				atkIcon := h.AnimState.Icon
				if atkIcon == "" {
					atkIcon = "⚔️"
				}
				val := h.AnimState.DamageText
				if val != "" {
					suffix := "!"
					if h.AnimState.IsCrit {
						suffix = "!!!"
					}
					if !strings.HasSuffix(val, "!") {
						val += suffix
					}
				}
				text := strings.TrimSpace(atkIcon + " " + val)
				if h.AnimState.IsCrit {
					indicator = critStyle.Render(text)
				} else {
					indicator = atkStyle.Render(text)
				}

			case ui.AnimBlock:
				blockIcon := h.AnimState.Icon
				if blockIcon == "" {
					blockIcon = "🛡️"
				}
				val := h.AnimState.DamageText
				if val == "" {
					val = "0"
				}
				indicator = blockStyle.Render(fmt.Sprintf("%s %s", blockIcon, val))

			case ui.AnimDodge:
				dodgeIcon := h.AnimState.Icon
				if dodgeIcon == "" {
					dodgeIcon = "󰆷"
				}
				val := h.AnimState.DamageText
				if val == "" {
					val = "0"
				}
				indicator = dodgeStyle.Render(fmt.Sprintf("%s %s", dodgeIcon, val))

			case ui.AnimTakeDamage:
				enemyIcon := h.AnimState.Icon
				if enemyIcon == "" {
					enemyIcon = "👾"
				}
				val := h.AnimState.DamageText
				if val == "" {
					val = "0"
				}
				indicator = takeDmgStyle.Render(fmt.Sprintf("%s %s", enemyIcon, val))
			}
		}

		indLen := lipgloss.Width(indicator)
		leftPad := (cardWidth - indLen) / 2
		if leftPad < 0 {
			leftPad = 0
		}
		rightPad := cardWidth - indLen - leftPad
		if rightPad < 0 {
			rightPad = 0
		}

		paddedInd := strings.Repeat(" ", leftPad) + indicator + strings.Repeat(" ", rightPad)
		damageElements = append(damageElements, paddedInd)
		if i < len(c.heroes)-1 {
			damageElements = append(damageElements, strings.Repeat(" ", spacerWidth))
		}
	}

	damageLine := strings.Join(damageElements, "")

	var cards []string
	for _, h := range c.heroes {
		cardView := ui.NewHeroCardComponent(c.theme, h, c.responsive).View()
		cards = append(cards, cardView)
	}

	var spacerCards []string
	for i, card := range cards {
		spacerCards = append(spacerCards, card)
		if i < len(cards)-1 {
			spacerCards = append(spacerCards, "  ")
		}
	}

	grid := lipgloss.JoinHorizontal(lipgloss.Top, spacerCards...)
	content := damageLine + "\n" + grid

	if c.width > 0 {
		return lipgloss.NewStyle().Width(c.width).Align(lipgloss.Center).Render(content)
	}
	return content
}
