package ui

import (
	"strings"
	"testing"

	"github.com/pomquest/internal/domain"
)

func TestPartyCardComponentView_LG(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ThemeFromConfig(cfg)

	data := PartyCardData{
		Gold:           0,
		Silver:         10,
		Copper:         80,
		FoodLbs:        2000,
		Loot:           67,
		SpareParts:     3,
		Oxen:           4,
		ItemUsage:      "heavy",
		Diet:           "balanced",
		Pace:           "fast",
		Greed:          "sparse",
		Potions:        2,
		Bombs:          5,
		Shovels:        1,
		Rope:           3,
		MagicItems:     0,
		Gems:           0,
		TreasureChests: 1,
		CurrentHP:      30,
		MaxHP:          35,
	}

	comp := NewPartyCardComponent(theme, data, 96)
	view := comp.View()

	if !strings.Contains(view, "Party") {
		t.Errorf("expected view to contain Party header, got:\n%s", view)
	}
	if !strings.Contains(view, "󰋑 󰋑 󰋑 󰛠") {
		t.Errorf("expected view to contain team hearts in top-right, got:\n%s", view)
	}
	if !strings.Contains(view, "1g 0s 80c") {
		t.Errorf("expected view to contain gold info 1g 0s 80c, got:\n%s", view)
	}
	if !strings.Contains(view, "󱐟 food: 2,000lbs") {
		t.Errorf("expected view to contain food info with icon, got:\n%s", view)
	}
	if !strings.Contains(view, "󱢏 loot: 67") {
		t.Errorf("expected view to contain loot info with icon, got:\n%s", view)
	}
	if !strings.Contains(view, " spare parts: 3") {
		t.Errorf("expected view to contain spare parts with icon, got:\n%s", view)
	}
	if !strings.Contains(view, " ox: 4") {
		t.Errorf("expected view to contain ox with icon, got:\n%s", view)
	}
	if !strings.Contains(view, " item usage: heavy ") {
		t.Errorf("expected view to contain item usage heavy with icons, got:\n%s", view)
	}
	if !strings.Contains(view, " diet: balanced ") {
		t.Errorf("expected view to contain diet balanced with icons, got:\n%s", view)
	}
	if !strings.Contains(view, "󰖃 pace: fast 󰑮") {
		t.Errorf("expected view to contain pace fast with icons, got:\n%s", view)
	}
	if !strings.Contains(view, " greed: sparse ") {
		t.Errorf("expected view to contain greed sparse with icons, got:\n%s", view)
	}
	if !strings.Contains(view, " 2") {
		t.Errorf("expected view to contain potions with icon only, got:\n%s", view)
	}
	if !strings.Contains(view, "󰚑 5") {
		t.Errorf("expected view to contain bombs with icon only, got:\n%s", view)
	}
	if !strings.Contains(view, "󰜐 1") {
		t.Errorf("expected view to contain shovels with icon only, got:\n%s", view)
	}
	if !strings.Contains(view, "󰼃 3") {
		t.Errorf("expected view to contain rope with icon only, got:\n%s", view)
	}
	if strings.Contains(view, "potions:") {
		t.Errorf("expected view NOT to contain 'potions:' label in LG mode, got:\n%s", view)
	}
}

func TestPartyCardComponentView_MD(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ThemeFromConfig(cfg)

	data := PartyCardData{
		Gold:           0,
		Silver:         10,
		Copper:         80,
		FoodLbs:        2000,
		Loot:           67,
		SpareParts:     3,
		Oxen:           4,
		ItemUsage:      "heavy",
		Diet:           "balanced",
		Pace:           "fast",
		Greed:          "sparse",
		Potions:        2,
		Bombs:          5,
		Shovels:        1,
		Rope:           3,
		CurrentHP:      30,
		MaxHP:          35,
	}

	comp := NewPartyCardComponent(theme, data, 72)
	view := comp.View()

	if !strings.Contains(view, "󱐟 2,000lbs") {
		t.Errorf("expected view to contain food info without label name or colon in MD mode, got:\n%s", view)
	}
	if !strings.Contains(view, "󱢏 67") {
		t.Errorf("expected view to contain loot info without label name or colon in MD mode, got:\n%s", view)
	}
	if !strings.Contains(view, " heavy ") {
		t.Errorf("expected view to contain item usage without label name or colon in MD mode, got:\n%s", view)
	}
	if strings.Contains(view, "food:") {
		t.Errorf("expected view NOT to contain 'food:' label in MD mode, got:\n%s", view)
	}
	if strings.Contains(view, "󱐟:") {
		t.Errorf("expected view NOT to contain colon after icon in MD mode, got:\n%s", view)
	}
}
