package dungeon

import (
	"math/rand"

	"github.com/pomquest/internal/app"
)

// MonsterTemplate defines raw stats for generating room combatants.
type MonsterTemplate struct {
	Name     string
	Biome    []string
	CR       int
	HP       int
	MaxHP    int
	Attack   int
	Armor    int
	XPValue  int
	LootGold int
	Icon     string
}

// BookOfMonsters contains predefined monster archetypes by biome.
var BookOfMonsters = []MonsterTemplate{
	// Forest Biome
	{Name: "Rat", Biome: []string{"Any"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐀"},
	{Name: "Bat", Biome: []string{"forest", "mountains"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🦇"},
	{Name: "Goblin Sentry", Biome: []string{"forest"}, CR: 4, HP: 6, MaxHP: 6, Attack: 2, Armor: 0, XPValue: 5, LootGold: 2, Icon: "👺"},
	{Name: "Skeleton", Biome: []string{"forest"}, CR: 4, HP: 8, MaxHP: 8, Attack: 2, Armor: 1, XPValue: 8, LootGold: 3, Icon: "💀"},
	{Name: "Ghoul", Biome: []string{"forest"}, CR: 4, HP: 12, MaxHP: 12, Attack: 3, Armor: 1, XPValue: 12, LootGold: 5, Icon: "🧟"},
	{Name: "Zombie", Biome: []string{"forest"}, CR: 4, HP: 14, MaxHP: 14, Attack: 3, Armor: 2, XPValue: 14, LootGold: 4, Icon: "🧟‍♂️"},

	// Swamp / Mountains Biome
	{Name: "Beetle", Biome: []string{"swamp"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🪲"},
	{Name: "Mote", Biome: []string{"swamp"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "👾"},
	{Name: "Spider", Biome: []string{"mountains"}, CR: 4, HP: 7, MaxHP: 7, Attack: 2, Armor: 0, XPValue: 6, LootGold: 2, Icon: "🕷️"},
	{Name: "Arachnid", Biome: []string{"mountains"}, CR: 4, HP: 12, MaxHP: 12, Attack: 3, Armor: 1, XPValue: 12, LootGold: 5, Icon: "🕷️"},
	{Name: "Lurker", Biome: []string{"mountains"}, CR: 4, HP: 18, MaxHP: 18, Attack: 4, Armor: 2, XPValue: 20, LootGold: 8, Icon: "👥"},
	{Name: "Troll", Biome: []string{"mountains"}, CR: 4, HP: 26, MaxHP: 26, Attack: 5, Armor: 3, XPValue: 30, LootGold: 12, Icon: "🧌"},

	// Desert / Island Biome
	{Name: "Sprite", Biome: []string{"desert"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "💨"},
	{Name: "Beetle", Biome: []string{"desert"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🪲"},
	{Name: "Goblin", Biome: []string{"desert"}, CR: 4, HP: 8, MaxHP: 8, Attack: 3, Armor: 0, XPValue: 8, LootGold: 4, Icon: "👺"},
	{Name: "Elemental", Biome: []string{"desert"}, CR: 4, HP: 20, MaxHP: 20, Attack: 4, Armor: 2, XPValue: 22, LootGold: 10, Icon: "🔥"},
	{Name: "Hellhound", Biome: []string{"desert"}, CR: 4, HP: 18, MaxHP: 18, Attack: 5, Armor: 1, XPValue: 25, LootGold: 8, Icon: "🐺"},
	{Name: "Drake", Biome: []string{"desert"}, CR: 5, HP: 32, MaxHP: 32, Attack: 6, Armor: 3, XPValue: 45, LootGold: 20, Icon: "🐉"},

	// Generic / Any Biome
	{Name: "Rat", Biome: []string{"Any"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐀"},
	{Name: "Slime", Biome: []string{"Any"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🟢"},
	{Name: "Marauder", Biome: []string{"Any"}, CR: 4, HP: 9, MaxHP: 9, Attack: 2, Armor: 1, XPValue: 8, LootGold: 5, Icon: "🥷"},
	{Name: "Gargoyle", Biome: []string{"Any"}, CR: 4, HP: 15, MaxHP: 15, Attack: 3, Armor: 3, XPValue: 15, LootGold: 6, Icon: "🗿"},
}

// GenerateMonstersForRoom picks monsters matching biome and CR budget for a room encounter.
// If maxCR > 0, monsters with CR > maxCR will be excluded from spawning.
func GenerateMonstersForRoom(biome string, crBudget int, r *rand.Rand, maxCR int) []app.Monster {
	return GenerateMonstersForBiomeRoom(biome, crBudget, r, maxCR)
}
