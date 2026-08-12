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

func TestMonstersForBiomeMultiBiome(t *testing.T) {
	// Temporarily add a test monster with multiple biomes
	testMonster := MonsterTemplate{
		Name:  "Test Multi Monster",
		Biome: []string{"forest", "desert"},
		CR:    1,
		HP:    5, MaxHP: 5, Attack: 1, Armor: 0, XPValue: 5, LootGold: 1, Icon: "👾",
	}
	BookOfMonsters = append(BookOfMonsters, testMonster)
	defer func() {
		BookOfMonsters = BookOfMonsters[:len(BookOfMonsters)-1]
	}()

	forestMonsters := MonstersForBiome("forest")
	desertMonsters := MonstersForBiome("desert")

	foundInForest := false
	for _, m := range forestMonsters {
		if m.Name == "Test Multi Monster" {
			foundInForest = true
			break
		}
	}
	if !foundInForest {
		t.Errorf("expected multi-biome monster to be found in forest biome")
	}

	foundInDesert := false
	for _, m := range desertMonsters {
		if m.Name == "Test Multi Monster" {
			foundInDesert = true
			break
		}
	}
	if !foundInDesert {
		t.Errorf("expected multi-biome monster to be found in desert biome")
	}
}

func TestAdaptableMonsterName(t *testing.T) {
	testMonster := MonsterTemplate{
		Name:      "TestAdaptable",
		Biome:     []string{"Any"},
		Adaptable: true,
		CR:        1,
		HP:        1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "👾",
	}
	BookOfMonsters = append(BookOfMonsters, testMonster)
	defer func() {
		BookOfMonsters = BookOfMonsters[:len(BookOfMonsters)-1]
	}()

	r := rand.New(rand.NewSource(1))
	// Force candidate selection to only testMonster by temporarily filtering BookOfMonsters
	oldBook := BookOfMonsters
	BookOfMonsters = []MonsterTemplate{testMonster}
	defer func() {
		BookOfMonsters = oldBook
	}()

	monsters := GenerateMonstersForBiomeRoom("forest", 1, r, 0)
	if len(monsters) == 0 {
		t.Fatalf("expected 1 monster spawned")
	}
	if monsters[0].Name != "Forest TestAdaptable" {
		t.Errorf("expected monster name 'Forest TestAdaptable', got %q", monsters[0].Name)
	}

	monstersDesert := GenerateMonstersForBiomeRoom("desert", 1, r, 0)
	if len(monstersDesert) == 0 {
		t.Fatalf("expected 1 monster spawned")
	}
	if monstersDesert[0].Name != "Desert TestAdaptable" {
		t.Errorf("expected monster name 'Desert TestAdaptable', got %q", monstersDesert[0].Name)
	}
}
