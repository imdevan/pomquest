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
// HP is evaluated in 10-HP and 5-HP increments using icon.KeyHeart...
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

	// Round maxHP up to nearest multiple of 5 (total HP capacity)
	roundedMax := (maxHP + 4) / 5 * 5
	if roundedMax < 5 {
		roundedMax = 5
	}

	num10s := roundedMax / 10
	hasRem5 := (roundedMax % 10) >= 5

	var hearts []string
	hpLeft := currentHP

	for i := 0; i < num10s; i++ {
		if hpLeft >= 8 {
			hearts = append(hearts, icon.Get(icon.KeyHeart10Full, mode))
			hpLeft -= 10
			if hpLeft < 0 {
				hpLeft = 0
			}
		} else if hpLeft >= 3 {
			hearts = append(hearts, icon.Get(icon.KeyHeart10Half, mode))
			hpLeft = 0
		} else {
			hearts = append(hearts, icon.Get(icon.KeyHeart10Empty, mode))
		}
	}

	if hasRem5 {
		if hpLeft >= 3 {
			hearts = append(hearts, icon.Get(icon.KeyHeart5Full, mode))
		} else {
			hearts = append(hearts, icon.Get(icon.KeyHeart5Empty, mode))
		}
	}

	return strings.Join(hearts, " ")
}
