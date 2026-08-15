package ui

import (
	"strings"

	"github.com/pomquest/internal/adapters/icon"
)

// TeamHeartsComponent renders party health using 10 HP and 5 HP heart icons.
type TeamHeartsComponent struct {
	CurrentHP int
	MaxHP     int
	Mode      icon.Mode
}

func NewTeamHeartsComponent(currentHP, maxHP int, mode icon.Mode) TeamHeartsComponent {
	return TeamHeartsComponent{
		CurrentHP: currentHP,
		MaxHP:     maxHP,
		Mode:      mode,
	}
}

// View renders the team hearts string.
func (c TeamHeartsComponent) View() string {
	return RenderTeamHearts(c.CurrentHP, c.MaxHP, c.Mode)
}

// RenderTeamHearts generates a heart icon string representing current health out of max health.
// HP is evaluated using method: total half hearts N = (HP + 2) / 5.
func RenderTeamHearts(currentHP, maxHP int, mode icon.Mode) string {
	if maxHP <= 0 {
		return ""
	}
	if currentHP < 0 {
		currentHP = 0
	}
	if currentHP > maxHP {
		currentHP = maxHP
	}

	// Calculate total half-hearts N = (HP + 2) / 5
	maxHalfHearts := (maxHP + 2) / 5
	if maxHalfHearts < 1 {
		maxHalfHearts = 1
	}

	curHalfHearts := (currentHP + 2) / 5

	num10s := maxHalfHearts / 2
	hasRem5 := (maxHalfHearts % 2) != 0

	var hearts []string
	hpLeft := curHalfHearts

	for i := 0; i < num10s; i++ {
		if hpLeft >= 2 {
			hearts = append(hearts, icon.Get(icon.KeyHeart10Full, mode))
			hpLeft -= 2
		} else if hpLeft == 1 {
			hearts = append(hearts, icon.Get(icon.KeyHeart10Half, mode))
			hpLeft = 0
		} else {
			hearts = append(hearts, icon.Get(icon.KeyHeart10Empty, mode))
		}
	}

	if hasRem5 {
		if hpLeft >= 1 {
			hearts = append(hearts, icon.Get(icon.KeyHeart5Full, mode))
		} else {
			hearts = append(hearts, icon.Get(icon.KeyHeart5Empty, mode))
		}
	}

	return strings.Join(hearts, " ")
}
