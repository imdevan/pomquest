package dungeon

import (
	"math/rand"
	"testing"
)

func TestBookOfBiomes(t *testing.T) {
	expectedBiomes := []string{"forest", "swamp", "plains", "mountains", "island", "desert"}
	for _, key := range expectedBiomes {
		b := GetBiome(key)
		if b.ID != key {
			t.Errorf("expected biome ID %q, got %q", key, b.ID)
		}
		if b.Name == "" {
			t.Errorf("expected non-empty name for biome %q", key)
		}
		if b.TrapRate <= 0.0 || b.TrapRate > 1.0 {
			t.Errorf("expected TrapRate between 0 and 1 for biome %q, got %f", key, b.TrapRate)
		}
		if len(b.HazardWeights) == 0 {
			t.Errorf("expected hazards for biome %q", key)
		}
		validSupplies := map[string]bool{
			"coin": true, "food": true, "spareparts": true, "ox": true,
			"loot": true, "potions": true, "rope": true, "bombs": true, "shovels": true,
		}
		if len(b.ResourceDropRates) == 0 {
			t.Errorf("expected resource drop rates for biome %q", key)
		}
		for res := range b.ResourceDropRates {
			if !validSupplies[res] {
				t.Errorf("biome %q has invalid resource drop type %q", key, res)
			}
		}
	}
}

func TestGetBiomeFallback(t *testing.T) {
	b := GetBiome("nonexistent")
	if b.ID != "forest" {
		t.Errorf("expected fallback biome 'forest', got %q", b.ID)
	}
}

func TestGenerateMonstersForBiomeRoom(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	monsters := GenerateMonstersForBiomeRoom("forest", 4, r, 0)
	if len(monsters) == 0 {
		t.Fatalf("expected at least 1 monster spawned for forest biome")
	}
	for _, m := range monsters {
		if m.Name == "" {
			t.Errorf("expected non-empty monster name")
		}
	}
}
