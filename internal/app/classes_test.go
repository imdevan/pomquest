package app

import (
	"testing"

	"github.com/pomquest/internal/adapters/icon"
)

func TestClassCatalog(t *testing.T) {
	catalog := ClassCatalog()
	if len(catalog) != 12 {
		t.Errorf("expected 12 class archetypes in catalog, got %d", len(catalog))
	}

	fighter := GetClassArchetype("fighter")
	if fighter.Name != "fighter" {
		t.Errorf("expected class name 'fighter', got %s", fighter.Name)
	}

	// Verify icon resolution using central icon library
	fighterIcon := fighter.Icon(icon.ModeNerdFont)
	expectedIcon := icon.GetClassIcon("fighter", icon.ModeNerdFont)
	if fighterIcon != expectedIcon {
		t.Errorf("expected fighter icon %q, got %q", expectedIcon, fighterIcon)
	}

	for _, class := range catalog {
		s := class.Stats
		totalStats := s.Str + s.Dex + s.Cha + s.Sur + s.Con + s.Wil + s.Int + s.Per + s.Hel + s.Luc
		if totalStats != 4 {
			t.Errorf("expected class %s to have total 4 stat points, got %d", class.Name, totalStats)
		}

		nonZeroCount := 0
		for _, val := range []int{s.Str, s.Dex, s.Cha, s.Sur, s.Con, s.Wil, s.Int, s.Per, s.Hel, s.Luc} {
			if val > 0 {
				nonZeroCount++
			}
		}
		if nonZeroCount > 4 {
			t.Errorf("expected class %s to have at most 4 non-zero stats, got %d", class.Name, nonZeroCount)
		}
	}
}
