package dungeon

import (
	"math/rand"
	"time"

	"github.com/pomquest/internal/app"
)

// MonsterTemplate defines raw stats for generating room combatants.
type MonsterTemplate struct {
	Name     string
	Biome    string
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
	// Ancient Catacombs - CR 1 Weaklings
	{Name: "Dungeon Rat", Biome: "Ancient Catacombs", CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐀"},
	{Name: "Catacomb Bat", Biome: "Ancient Catacombs", CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🦇"},

	// Ancient Catacombs - CR 4+ Standard & Elite
	{Name: "Goblin Sentry", Biome: "Ancient Catacombs", CR: 4, HP: 6, MaxHP: 6, Attack: 2, Armor: 0, XPValue: 5, LootGold: 2, Icon: "👺"},
	{Name: "Skeletal Soldier", Biome: "Ancient Catacombs", CR: 4, HP: 8, MaxHP: 8, Attack: 2, Armor: 1, XPValue: 8, LootGold: 3, Icon: "💀"},
	{Name: "Crypt Ghoul", Biome: "Ancient Catacombs", CR: 4, HP: 12, MaxHP: 12, Attack: 3, Armor: 1, XPValue: 12, LootGold: 5, Icon: "🧟"},
	{Name: "Zombie Warrior", Biome: "Ancient Catacombs", CR: 4, HP: 14, MaxHP: 14, Attack: 3, Armor: 2, XPValue: 14, LootGold: 4, Icon: "🧟‍♂️"},

	// Shadow Caverns - CR 1 Weaklings
	{Name: "Cave Beetle", Biome: "Shadow Caverns", CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🪲"},
	{Name: "Shadow Mote", Biome: "Shadow Caverns", CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "👾"},

	// Shadow Caverns - CR 4+ Standard & Elite
	{Name: "Cave Spider", Biome: "Shadow Caverns", CR: 4, HP: 7, MaxHP: 7, Attack: 2, Armor: 0, XPValue: 6, LootGold: 2, Icon: "🕷️"},
	{Name: "Giant Arachnid", Biome: "Shadow Caverns", CR: 4, HP: 12, MaxHP: 12, Attack: 3, Armor: 1, XPValue: 12, LootGold: 5, Icon: "🕷️"},
	{Name: "Shadow Lurker", Biome: "Shadow Caverns", CR: 4, HP: 18, MaxHP: 18, Attack: 4, Armor: 2, XPValue: 20, LootGold: 8, Icon: "👥"},
	{Name: "Subterranean Troll", Biome: "Shadow Caverns", CR: 4, HP: 26, MaxHP: 26, Attack: 5, Armor: 3, XPValue: 30, LootGold: 12, Icon: "🧌"},

	// Brimstone Keep - CR 1 Weaklings
	{Name: "Ash Sprite", Biome: "Brimstone Keep", CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "💨"},
	{Name: "Fire Beetle", Biome: "Brimstone Keep", CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🪲"},

	// Brimstone Keep - CR 4+ Standard & Elite
	{Name: "Fire Goblin", Biome: "Brimstone Keep", CR: 4, HP: 8, MaxHP: 8, Attack: 3, Armor: 0, XPValue: 8, LootGold: 4, Icon: "👺"},
	{Name: "Lava Elemental", Biome: "Brimstone Keep", CR: 4, HP: 20, MaxHP: 20, Attack: 4, Armor: 2, XPValue: 22, LootGold: 10, Icon: "🔥"},
	{Name: "Hellhound", Biome: "Brimstone Keep", CR: 4, HP: 18, MaxHP: 18, Attack: 5, Armor: 1, XPValue: 25, LootGold: 8, Icon: "🐺"},
	{Name: "Brimstone Drake", Biome: "Brimstone Keep", CR: 5, HP: 32, MaxHP: 32, Attack: 6, Armor: 3, XPValue: 45, LootGold: 20, Icon: "🐉"},

	// Generic / Any Biome - CR 1 Weaklings
	{Name: "Sewer Rat", Biome: "Any", CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐀"},
	{Name: "Tiny Slime", Biome: "Any", CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🟢"},

	// Generic / Any Biome - CR 4+ Standard & Elite
	{Name: "Bandit Marauder", Biome: "Any", CR: 4, HP: 9, MaxHP: 9, Attack: 2, Armor: 1, XPValue: 8, LootGold: 5, Icon: "🥷"},
	{Name: "Dungeon Gargoyle", Biome: "Any", CR: 4, HP: 15, MaxHP: 15, Attack: 3, Armor: 3, XPValue: 15, LootGold: 6, Icon: "🗿"},
}

// GenerateMonstersForRoom picks monsters matching biome and CR budget for a room encounter.
func GenerateMonstersForRoom(biome string, crBudget int, r *rand.Rand) []app.Monster {
	if r == nil {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	if crBudget < 1 {
		crBudget = 1
	}

	var candidatePool []MonsterTemplate
	for _, m := range BookOfMonsters {
		if m.Biome == biome || m.Biome == "Any" {
			candidatePool = append(candidatePool, m)
		}
	}
	if len(candidatePool) == 0 {
		candidatePool = BookOfMonsters
	}

	var spawned []app.Monster
	currentCR := 0

	for currentCR < crBudget && len(spawned) < 3 {
		tmpl := candidatePool[r.Intn(len(candidatePool))]
		monster := app.Monster{
			Creature: app.Creature{
				Noun: app.Noun{
					Name: tmpl.Name,
					Icon: tmpl.Icon,
				},
				HP:     tmpl.HP,
				MaxHP:  tmpl.MaxHP,
				Attack: tmpl.Attack,
				Armor:  tmpl.Armor,
			},
			XPValue:  tmpl.XPValue,
			LootGold: tmpl.LootGold,
			Biome:    tmpl.Biome,
		}
		spawned = append(spawned, monster)
		currentCR += tmpl.CR
	}

	return spawned
}
