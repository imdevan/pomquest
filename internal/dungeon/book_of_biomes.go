package dungeon

import (
	"math/rand"
	"time"

	"github.com/pomquest/internal/app"
)

// Biome represents an environment type for a dungeon location.
type Biome struct {
	ID                string
	Name              string
	HazardWeights     map[string]int
	TrapRate          float64 // 0.0 to 1.0 likelihood of traps in rooms
	ResourceDropRates map[string]float64
}

// BookOfBiomes contains predefined biomes for dungeon locations.
var BookOfBiomes = map[string]Biome{
	"forest": {
		ID:   "forest",
		Name: "Forest",
		HazardWeights: map[string]int{
			"Entangling Vines": 40,
			"Poison Ivy":       30,
			"Falling Branch":   30,
		},
		TrapRate: 0.15,
		ResourceDropRates: map[string]float64{
			"food":    0.40,
			"rope":    0.30,
			"potions": 0.20,
			"coin":    0.10,
		},
	},
	"swamp": {
		ID:   "swamp",
		Name: "Swamp",
		HazardWeights: map[string]int{
			"Mire Pit":     50,
			"Toxic Gas":    30,
			"Leech Infest": 20,
		},
		TrapRate: 0.25,
		ResourceDropRates: map[string]float64{
			"potions": 0.35,
			"shovels": 0.30,
			"rope":    0.20,
			"loot":    0.15,
		},
	},
	"plains": {
		ID:   "plains",
		Name: "Plains",
		HazardWeights: map[string]int{
			"Wildfire":       40,
			"Stampede":       35,
			"Concealed Hole": 25,
		},
		TrapRate: 0.10,
		ResourceDropRates: map[string]float64{
			"food":       0.50,
			"ox":         0.20,
			"spareparts": 0.20,
			"coin":       0.10,
		},
	},
	"mountains": {
		ID:   "mountains",
		Name: "Mountains",
		HazardWeights: map[string]int{
			"Rockslide":      50,
			"Freezing Winds": 30,
			"Chasm Collapse": 20,
		},
		TrapRate: 0.20,
		ResourceDropRates: map[string]float64{
			"rope":    0.35,
			"bombs":   0.30,
			"shovels": 0.20,
			"loot":    0.15,
		},
	},
	"island": {
		ID:   "island",
		Name: "Island",
		HazardWeights: map[string]int{
			"Riptide":      40,
			"Quick Sand":   30,
			"Volcanic Ash": 30,
		},
		TrapRate: 0.18,
		ResourceDropRates: map[string]float64{
			"food":    0.40,
			"loot":    0.30,
			"potions": 0.20,
			"coin":    0.10,
		},
	},
	"desert": {
		ID:   "desert",
		Name: "Desert",
		HazardWeights: map[string]int{
			"Sandstorm":   45,
			"Heat Stroke": 35,
			"Quicksand":   20,
		},
		TrapRate: 0.20,
		ResourceDropRates: map[string]float64{
			"shovels":    0.35,
			"spareparts": 0.25,
			"loot":       0.25,
			"coin":       0.15,
		},
	},
}

// GetBiomeHazardDescription returns a summary string of hazards in a biome.
func GetBiomeHazardDescription(b Biome) string {
	var hazards []string
	for h := range b.HazardWeights {
		hazards = append(hazards, h)
	}
	if len(hazards) == 0 {
		return "unknown hazards"
	}
	return "hazards: " + hazards[0]
}

// AllBiomes is a list of all defined biome keys.
var AllBiomes = []string{
	"forest",
	"swamp",
	"plains",
	"mountains",
	"island",
	"desert",
}

// GetBiome returns a Biome by ID/key, falling back to "forest" if not found.
func GetBiome(id string) Biome {
	if b, ok := BookOfBiomes[id]; ok {
		return b
	}
	return BookOfBiomes["forest"]
}

// MonstersForBiome returns monster templates associated with a given biome ID.
func MonstersForBiome(biomeID string) []MonsterTemplate {
	var matches []MonsterTemplate
	for _, m := range BookOfMonsters {
		for _, b := range m.Biome {
			if b == biomeID || b == "Any" {
				matches = append(matches, m)
				break
			}
		}
	}
	if len(matches) == 0 {
		return BookOfMonsters
	}
	return matches
}

// GenerateMonstersForBiomeRoom picks monsters matching the biome and CR budget for a room encounter.
func GenerateMonstersForBiomeRoom(biomeID string, crBudget int, r *rand.Rand, maxCR int) []app.Monster {
	if r == nil {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	if crBudget < 1 {
		crBudget = 1
	}

	candidates := MonstersForBiome(biomeID)
	var filtered []MonsterTemplate
	for _, m := range candidates {
		if maxCR <= 0 || m.CR <= maxCR {
			filtered = append(filtered, m)
		}
	}
	if len(filtered) == 0 {
		for _, m := range BookOfMonsters {
			if maxCR <= 0 || m.CR <= maxCR {
				filtered = append(filtered, m)
			}
		}
	}
	if len(filtered) == 0 {
		filtered = BookOfMonsters
	}

	var spawned []app.Monster
	currentCR := 0

	for currentCR < crBudget && len(spawned) < 3 {
		tmpl := filtered[r.Intn(len(filtered))]
		name := tmpl.Name
		if tmpl.Adaptable && biomeID != "" {
			b := GetBiome(biomeID)
			if b.Name != "" {
				name = b.Name + " " + tmpl.Name
			}
		}
		monster := app.Monster{
			Creature: app.Creature{
				Noun: app.Noun{
					Name: name,
					Icon: tmpl.Icon,
				},
				HP:     tmpl.HP,
				MaxHP:  tmpl.MaxHP,
				Attack: tmpl.Attack,
				Armor:  tmpl.Armor,
			},
			XPValue:  tmpl.XPValue,
			LootGold: tmpl.LootGold,
			Biome:    biomeID,
		}
		spawned = append(spawned, monster)
		currentCR += tmpl.CR
	}

	return spawned
}
