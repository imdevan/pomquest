package ui

import (
	"strings"
	"testing"

	"github.com/pomquest/internal/domain"
)

func TestPartyCardComponentView(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ThemeFromConfig(cfg)

	data := PartyCardData{
		Gold:           0,
		Silver:         10,
		Copper:         80,
		FoodLbs:        2000,
		Diet:           "balanced",
		Pace:           "fast",
		SpareParts:     3,
		Oxen:           4,
		MagicItems:     0,
		Gems:           0,
		TreasureChests: 1,
	}

	comp := NewPartyCardComponent(theme, data, 72)
	view := comp.View()

	if !strings.Contains(view, "Party:") {
		t.Errorf("expected view to contain Party:, got:\n%s", view)
	}
	if !strings.Contains(view, "coins: 0g 10s 80c") {
		t.Errorf("expected view to contain coins info, got:\n%s", view)
	}
	if !strings.Contains(view, "food: 2,000lbs") {
		t.Errorf("expected view to contain food info, got:\n%s", view)
	}
}
