package dungeon

import (
	"math/rand"

	"github.com/pomquest/internal/app"
)

// MonsterTemplate defines raw stats for generating room combatants.
type MonsterTemplate struct {
	Name      string
	Biome     []string
	Adaptable bool
	CR        int
	HP        int
	MaxHP     int
	Attack    int
	Armor     int
	XPValue   int
	LootGold  int
	Icon      string
	Stats     app.Stats
}

// BookOfMonsters contains predefined monster archetypes by biome.
var BookOfMonsters = []MonsterTemplate{
	// Forest Biome
	{Name: "Rat", Biome: []string{"Any"}, Adaptable: true, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐀", Stats: app.Stats{Str: -2, Dex: 1, Con: -2, Wil: -2, Int: -3, Sur: 0}},
	{Name: "Bat", Biome: []string{"forest", "mountains"}, Adaptable: true, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🦇", Stats: app.Stats{Str: -3, Dex: 2, Con: -2, Wil: -1, Int: -3, Sur: 1}},
	{Name: "Goblin Sentry", Biome: []string{"forest"}, CR: 4, HP: 6, MaxHP: 6, Attack: 2, Armor: 0, XPValue: 5, LootGold: 2, Icon: "👺", Stats: app.Stats{Str: 0, Dex: 2, Con: 0, Wil: -1, Int: 0, Sur: 1}},
	{Name: "Skeleton", Biome: []string{"forest"}, CR: 4, HP: 8, MaxHP: 8, Attack: 2, Armor: 1, XPValue: 8, LootGold: 3, Icon: "💀", Stats: app.Stats{Str: 1, Dex: 1, Con: 1, Wil: 0, Int: -2, Sur: 0}},
	{Name: "Ghoul", Biome: []string{"forest"}, CR: 4, HP: 12, MaxHP: 12, Attack: 3, Armor: 1, XPValue: 12, LootGold: 5, Icon: "🧟", Stats: app.Stats{Str: 2, Dex: 2, Con: 2, Wil: 1, Int: -1, Sur: 1}},
	{Name: "Zombie", Biome: []string{"forest"}, CR: 4, HP: 14, MaxHP: 14, Attack: 3, Armor: 2, XPValue: 14, LootGold: 4, Icon: "🧟‍♂️", Stats: app.Stats{Str: 3, Dex: -1, Con: 3, Wil: 0, Int: -3, Sur: 0}},

	// Swamp / Mountains Biome
	{Name: "Beetle", Biome: []string{"swamp", "mountains"}, Adaptable: true, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🪲", Stats: app.Stats{Str: -2, Dex: -1, Con: -1, Wil: -2, Int: -4, Sur: 0}},
	{Name: "Mote", Biome: []string{"swamp"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "👾", Stats: app.Stats{Str: -3, Dex: 3, Con: -2, Wil: 0, Int: -2, Sur: 1}},
	{Name: "Spider", Biome: []string{"mountains", "swamp", "forest"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🕷️", Stats: app.Stats{Str: -2, Dex: 2, Con: -2, Wil: -1, Int: -3, Sur: 1}},
	{Name: "Arachnid", Biome: []string{"mountains"}, CR: 4, HP: 12, MaxHP: 12, Attack: 3, Armor: 1, XPValue: 12, LootGold: 5, Icon: "🕷️", Stats: app.Stats{Str: 1, Dex: 3, Con: 1, Wil: 0, Int: -2, Sur: 2}},
	{Name: "Lurker", Biome: []string{"mountains"}, CR: 4, HP: 18, MaxHP: 18, Attack: 4, Armor: 2, XPValue: 20, LootGold: 8, Icon: "👥", Stats: app.Stats{Str: 3, Dex: 2, Con: 3, Wil: 1, Int: 0, Sur: 2}},
	{Name: "Troll", Biome: []string{"mountains"}, CR: 4, HP: 26, MaxHP: 26, Attack: 5, Armor: 3, XPValue: 30, LootGold: 12, Icon: "🧌", Stats: app.Stats{Str: 4, Dex: 0, Con: 4, Wil: 1, Int: -2, Sur: 1}},

	// Desert / Island Biome
	{Name: "Sprite", Biome: []string{"desert"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "💨", Stats: app.Stats{Str: -3, Dex: 3, Con: -2, Wil: 0, Int: 0, Sur: 1}},
	{Name: "Beetle", Biome: []string{"desert"}, Adaptable: true, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🪲", Stats: app.Stats{Str: -2, Dex: -1, Con: -1, Wil: -2, Int: -4, Sur: 0}},
	{Name: "Goblin", Biome: []string{"desert"}, CR: 4, HP: 8, MaxHP: 8, Attack: 3, Armor: 0, XPValue: 8, LootGold: 4, Icon: "👺", Stats: app.Stats{Str: 0, Dex: 2, Con: 0, Wil: -1, Int: 0, Sur: 1}},
	{Name: "Elemental", Biome: []string{"desert"}, CR: 4, HP: 20, MaxHP: 20, Attack: 4, Armor: 2, XPValue: 22, LootGold: 10, Icon: "🔥", Stats: app.Stats{Str: 3, Dex: 2, Con: 3, Wil: 2, Int: 1, Sur: 1}},
	{Name: "Hellhound", Biome: []string{"desert"}, CR: 4, HP: 18, MaxHP: 18, Attack: 5, Armor: 1, XPValue: 25, LootGold: 8, Icon: "🐺", Stats: app.Stats{Str: 3, Dex: 3, Con: 2, Wil: 1, Int: -1, Sur: 2}},
	{Name: "Drake", Biome: []string{"desert"}, CR: 5, HP: 32, MaxHP: 32, Attack: 6, Armor: 3, XPValue: 45, LootGold: 20, Icon: "🐉", Stats: app.Stats{Str: 5, Dex: 2, Con: 4, Wil: 3, Int: 2, Sur: 2}},

	// Generic / Any Biome
	{Name: "Rat", Biome: []string{"Any"}, Adaptable: true, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐀", Stats: app.Stats{Str: -2, Dex: 1, Con: -2, Wil: -2, Int: -3, Sur: 0}},
	{Name: "Slime", Biome: []string{"Any"}, Adaptable: true, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🟢", Stats: app.Stats{Str: -2, Dex: -2, Con: 0, Wil: -3, Int: -5, Sur: -1}},
	{Name: "Marauder", Biome: []string{"Any"}, CR: 4, HP: 9, MaxHP: 9, Attack: 2, Armor: 1, XPValue: 8, LootGold: 5, Icon: "🥷", Stats: app.Stats{Str: 1, Dex: 2, Con: 1, Wil: 0, Int: 0, Sur: 1}},
	{Name: "Gargoyle", Biome: []string{"Any"}, CR: 4, HP: 15, MaxHP: 15, Attack: 3, Armor: 3, XPValue: 15, LootGold: 6, Icon: "🗿", Stats: app.Stats{Str: 3, Dex: 1, Con: 3, Wil: 1, Int: -1, Sur: 1}},
}

// GenerateMonstersForRoom picks monsters matching biome and CR budget for a room encounter.
// If maxCR > 0, monsters with CR > maxCR will be excluded from spawning.
func GenerateMonstersForRoom(biome string, crBudget int, r *rand.Rand, maxCR int) []app.Monster {
	return GenerateMonstersForBiomeRoom(biome, crBudget, r, maxCR)
}
