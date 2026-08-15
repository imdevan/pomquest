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
	// ── Forest Biome ─────────────────────────────────────────────────────

	// CR 1 — Forest
	{Name: "Twig Blight", Biome: []string{"forest"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🌿", Stats: app.Stats{Str: -2, Dex: 1, Con: -1, Wil: -2, Int: -4, Sur: 1}},
	{Name: "Stoat", Biome: []string{"forest", "plains"}, Adaptable: true, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🦦", Stats: app.Stats{Str: -2, Dex: 3, Con: -2, Wil: -1, Int: -2, Sur: 2}},
	{Name: "Bat", Biome: []string{"forest", "mountains"}, Adaptable: true, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🦇", Stats: app.Stats{Str: -3, Dex: 2, Con: -2, Wil: -1, Int: -3, Sur: 1}},
	{Name: "Spider", Biome: []string{"forest", "swamp", "mountains"}, Adaptable: true, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🕷️", Stats: app.Stats{Str: -2, Dex: 2, Con: -2, Wil: -1, Int: -3, Sur: 1}},
	{Name: "Stinkbug", Biome: []string{"forest"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🪲", Stats: app.Stats{Str: -3, Dex: 0, Con: 0, Wil: -2, Int: -4, Sur: 0}},
	{Name: "Pixie", Biome: []string{"forest"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "✨", Stats: app.Stats{Str: -4, Dex: 3, Con: -3, Wil: 1, Int: 1, Sur: 2}},
	{Name: "Fungal Sprout", Biome: []string{"forest", "swamp"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🍄", Stats: app.Stats{Str: -2, Dex: -1, Con: 1, Wil: -1, Int: -4, Sur: 0}},
	{Name: "Wolf Pup", Biome: []string{"forest"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐺", Stats: app.Stats{Str: -1, Dex: 2, Con: -1, Wil: -1, Int: -3, Sur: 1}},
	{Name: "Adder", Biome: []string{"forest", "swamp"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐍", Stats: app.Stats{Str: -3, Dex: 3, Con: -2, Wil: -1, Int: -3, Sur: 2}},
	{Name: "Moss Mite", Biome: []string{"forest"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐛", Stats: app.Stats{Str: -3, Dex: 0, Con: -1, Wil: -3, Int: -5, Sur: 0}},

	// CR 4 — Forest
	{Name: "Goblin Sentry", Biome: []string{"forest"}, CR: 4, HP: 6, MaxHP: 6, Attack: 2, Armor: 0, XPValue: 5, LootGold: 2, Icon: "👺", Stats: app.Stats{Str: 0, Dex: 2, Con: 0, Wil: -1, Int: 0, Sur: 1}},
	{Name: "Skeleton", Biome: []string{"forest"}, CR: 4, HP: 8, MaxHP: 8, Attack: 2, Armor: 1, XPValue: 8, LootGold: 3, Icon: "💀", Stats: app.Stats{Str: 1, Dex: 1, Con: 1, Wil: 0, Int: -2, Sur: 0}},
	{Name: "Ghoul", Biome: []string{"forest"}, CR: 4, HP: 12, MaxHP: 12, Attack: 3, Armor: 1, XPValue: 12, LootGold: 5, Icon: "🧟", Stats: app.Stats{Str: 2, Dex: 2, Con: 2, Wil: 1, Int: -1, Sur: 1}},
	{Name: "Zombie", Biome: []string{"forest"}, CR: 4, HP: 14, MaxHP: 14, Attack: 3, Armor: 2, XPValue: 14, LootGold: 4, Icon: "🧟‍♂️", Stats: app.Stats{Str: 3, Dex: -1, Con: 3, Wil: 0, Int: -3, Sur: 0}},
	{Name: "Owlbear Cub", Biome: []string{"forest"}, CR: 4, HP: 16, MaxHP: 16, Attack: 4, Armor: 1, XPValue: 18, LootGold: 6, Icon: "🦉", Stats: app.Stats{Str: 3, Dex: 1, Con: 2, Wil: 0, Int: -2, Sur: 1}},
	{Name: "Displacer Stalker", Biome: []string{"forest"}, CR: 4, HP: 10, MaxHP: 10, Attack: 3, Armor: 0, XPValue: 12, LootGold: 5, Icon: "🐆", Stats: app.Stats{Str: 1, Dex: 4, Con: 1, Wil: 1, Int: 0, Sur: 3}},
	{Name: "Brigand", Biome: []string{"forest", "plains"}, CR: 4, HP: 9, MaxHP: 9, Attack: 3, Armor: 1, XPValue: 10, LootGold: 7, Icon: "🗡️", Stats: app.Stats{Str: 2, Dex: 2, Con: 1, Wil: 0, Int: 0, Sur: 1}},
	{Name: "Dire Wolf", Biome: []string{"forest"}, CR: 4, HP: 18, MaxHP: 18, Attack: 4, Armor: 1, XPValue: 20, LootGold: 4, Icon: "🐕", Stats: app.Stats{Str: 3, Dex: 3, Con: 2, Wil: 1, Int: -2, Sur: 2}},
	{Name: "Bone Soldier", Biome: []string{"forest"}, CR: 4, HP: 10, MaxHP: 10, Attack: 3, Armor: 2, XPValue: 12, LootGold: 5, Icon: "☠️", Stats: app.Stats{Str: 2, Dex: 1, Con: 2, Wil: 1, Int: -2, Sur: 0}},
	{Name: "Cultist Acolyte", Biome: []string{"forest", "swamp"}, CR: 4, HP: 8, MaxHP: 8, Attack: 2, Armor: 0, XPValue: 10, LootGold: 4, Icon: "🔮", Stats: app.Stats{Str: 0, Dex: 1, Con: 0, Wil: 3, Int: 2, Sur: 1}},

	// CR 5 — Forest
	{Name: "Owlbear", Biome: []string{"forest"}, CR: 5, HP: 34, MaxHP: 34, Attack: 6, Armor: 2, XPValue: 45, LootGold: 18, Icon: "🦅", Stats: app.Stats{Str: 5, Dex: 2, Con: 4, Wil: 1, Int: -1, Sur: 2}},
	{Name: "Treant", Biome: []string{"forest"}, CR: 5, HP: 40, MaxHP: 40, Attack: 5, Armor: 4, XPValue: 50, LootGold: 15, Icon: "🌳", Stats: app.Stats{Str: 5, Dex: -2, Con: 5, Wil: 3, Int: 1, Sur: 0}},
	{Name: "Forest Lich", Biome: []string{"forest"}, CR: 5, HP: 28, MaxHP: 28, Attack: 7, Armor: 2, XPValue: 50, LootGold: 25, Icon: "👻", Stats: app.Stats{Str: 1, Dex: 2, Con: 2, Wil: 5, Int: 5, Sur: 2}},
	{Name: "Mantis God", Biome: []string{"forest"}, CR: 5, HP: 30, MaxHP: 30, Attack: 7, Armor: 2, XPValue: 48, LootGold: 20, Icon: "🦗", Stats: app.Stats{Str: 4, Dex: 5, Con: 3, Wil: 2, Int: 1, Sur: 3}},
	{Name: "Grizzly Ancient", Biome: []string{"forest"}, CR: 5, HP: 38, MaxHP: 38, Attack: 6, Armor: 3, XPValue: 48, LootGold: 16, Icon: "🐻", Stats: app.Stats{Str: 5, Dex: 1, Con: 5, Wil: 2, Int: -1, Sur: 1}},
	{Name: "The Collector", Biome: []string{"forest"}, CR: 5, HP: 30, MaxHP: 30, Attack: 6, Armor: 2, XPValue: 45, LootGold: 22, Icon: "🎭", Stats: app.Stats{Str: 3, Dex: 3, Con: 3, Wil: 4, Int: 3, Sur: 2}},
	{Name: "Wraith Lord", Biome: []string{"forest"}, CR: 5, HP: 28, MaxHP: 28, Attack: 6, Armor: 3, XPValue: 46, LootGold: 20, Icon: "👤", Stats: app.Stats{Str: 2, Dex: 4, Con: 2, Wil: 4, Int: 3, Sur: 3}},
	{Name: "Shambler", Biome: []string{"forest"}, CR: 5, HP: 36, MaxHP: 36, Attack: 5, Armor: 3, XPValue: 50, LootGold: 24, Icon: "🌑", Stats: app.Stats{Str: 4, Dex: 1, Con: 4, Wil: 5, Int: 2, Sur: 1}},
	{Name: "Green Dragon Wyrmling", Biome: []string{"forest"}, CR: 5, HP: 32, MaxHP: 32, Attack: 6, Armor: 3, XPValue: 48, LootGold: 22, Icon: "🐲", Stats: app.Stats{Str: 4, Dex: 3, Con: 4, Wil: 3, Int: 3, Sur: 2}},
	{Name: "Urayuli", Biome: []string{"forest"}, CR: 5, HP: 35, MaxHP: 35, Attack: 6, Armor: 2, XPValue: 46, LootGold: 18, Icon: "🦍", Stats: app.Stats{Str: 5, Dex: 2, Con: 4, Wil: 2, Int: 0, Sur: 2}},

	// ── Swamp Biome ──────────────────────────────────────────────────────

	// CR 1 — Swamp
	{Name: "Mote", Biome: []string{"swamp"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "👾", Stats: app.Stats{Str: -3, Dex: 3, Con: -2, Wil: 0, Int: -2, Sur: 1}},
	{Name: "Leech", Biome: []string{"swamp"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🪱", Stats: app.Stats{Str: -3, Dex: 1, Con: 0, Wil: -2, Int: -5, Sur: 0}},
	{Name: "Maggot", Biome: []string{"swamp"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐛", Stats: app.Stats{Str: -3, Dex: 0, Con: 0, Wil: -3, Int: -5, Sur: -1}},
	{Name: "Bog Frog", Biome: []string{"swamp"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐸", Stats: app.Stats{Str: -2, Dex: 2, Con: -1, Wil: -2, Int: -3, Sur: 1}},
	{Name: "Corpse Fly", Biome: []string{"swamp"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🪰", Stats: app.Stats{Str: -4, Dex: 3, Con: -2, Wil: -2, Int: -4, Sur: 1}},
	{Name: "Mud Crab", Biome: []string{"swamp", "island"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🦀", Stats: app.Stats{Str: -1, Dex: -1, Con: 0, Wil: -2, Int: -4, Sur: 0}},
	{Name: "Newt", Biome: []string{"swamp"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🦎", Stats: app.Stats{Str: -3, Dex: 2, Con: -1, Wil: -2, Int: -3, Sur: 1}},
	{Name: "Gas Spore", Biome: []string{"swamp"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "💨", Stats: app.Stats{Str: -4, Dex: -2, Con: -1, Wil: -1, Int: -5, Sur: -1}},
	{Name: "Beetle", Biome: []string{"swamp", "mountains"}, Adaptable: true, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🪲", Stats: app.Stats{Str: -2, Dex: -1, Con: -1, Wil: -2, Int: -4, Sur: 0}},
	{Name: "Fungal Crawler", Biome: []string{"swamp"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🍄", Stats: app.Stats{Str: -2, Dex: 0, Con: 1, Wil: -1, Int: -4, Sur: 0}},

	// CR 4 — Swamp
	{Name: "Swine Wretch", Biome: []string{"swamp"}, CR: 4, HP: 12, MaxHP: 12, Attack: 3, Armor: 1, XPValue: 12, LootGold: 4, Icon: "🐷", Stats: app.Stats{Str: 3, Dex: 0, Con: 3, Wil: 0, Int: -2, Sur: 0}},
	{Name: "Pelagic Guardian", Biome: []string{"swamp", "island"}, CR: 4, HP: 14, MaxHP: 14, Attack: 3, Armor: 2, XPValue: 14, LootGold: 5, Icon: "🐟", Stats: app.Stats{Str: 2, Dex: 1, Con: 3, Wil: 1, Int: -1, Sur: 1}},
	{Name: "Bog Witch", Biome: []string{"swamp"}, CR: 4, HP: 8, MaxHP: 8, Attack: 4, Armor: 0, XPValue: 12, LootGold: 6, Icon: "🧙", Stats: app.Stats{Str: 0, Dex: 1, Con: 1, Wil: 4, Int: 3, Sur: 1}},
	{Name: "Carrion Feeder", Biome: []string{"swamp"}, CR: 4, HP: 10, MaxHP: 10, Attack: 3, Armor: 1, XPValue: 10, LootGold: 3, Icon: "🦅", Stats: app.Stats{Str: 2, Dex: 2, Con: 2, Wil: 0, Int: -2, Sur: 1}},
	{Name: "Fungal Hulk", Biome: []string{"swamp"}, CR: 4, HP: 20, MaxHP: 20, Attack: 3, Armor: 2, XPValue: 18, LootGold: 5, Icon: "🟤", Stats: app.Stats{Str: 3, Dex: -1, Con: 4, Wil: 1, Int: -3, Sur: 0}},
	{Name: "Louse Swarm", Biome: []string{"swamp"}, CR: 4, HP: 8, MaxHP: 8, Attack: 2, Armor: 0, XPValue: 8, LootGold: 2, Icon: "🐜", Stats: app.Stats{Str: -1, Dex: 3, Con: 0, Wil: -1, Int: -4, Sur: 1}},
	{Name: "Drowned Sailor", Biome: []string{"swamp"}, CR: 4, HP: 12, MaxHP: 12, Attack: 3, Armor: 1, XPValue: 12, LootGold: 6, Icon: "🧟", Stats: app.Stats{Str: 2, Dex: 1, Con: 2, Wil: 1, Int: -2, Sur: 0}},
	{Name: "Crocodilian", Biome: []string{"swamp"}, CR: 4, HP: 18, MaxHP: 18, Attack: 4, Armor: 2, XPValue: 20, LootGold: 5, Icon: "🐊", Stats: app.Stats{Str: 4, Dex: 1, Con: 3, Wil: 0, Int: -3, Sur: 1}},
	{Name: "Virulent Blob", Biome: []string{"swamp"}, CR: 4, HP: 16, MaxHP: 16, Attack: 2, Armor: 1, XPValue: 14, LootGold: 3, Icon: "🟣", Stats: app.Stats{Str: 1, Dex: -2, Con: 4, Wil: 0, Int: -4, Sur: -1}},
	{Name: "Swamp Cultist", Biome: []string{"swamp"}, CR: 4, HP: 10, MaxHP: 10, Attack: 3, Armor: 0, XPValue: 10, LootGold: 5, Icon: "🔮", Stats: app.Stats{Str: 1, Dex: 1, Con: 1, Wil: 3, Int: 2, Sur: 1}},

	// CR 5 — Swamp
	{Name: "Swine Prince", Biome: []string{"swamp"}, CR: 5, HP: 36, MaxHP: 36, Attack: 6, Armor: 3, XPValue: 48, LootGold: 20, Icon: "🐗", Stats: app.Stats{Str: 5, Dex: 0, Con: 5, Wil: 1, Int: -1, Sur: 0}},
	{Name: "Hag Matron", Biome: []string{"swamp"}, CR: 5, HP: 28, MaxHP: 28, Attack: 7, Armor: 2, XPValue: 50, LootGold: 22, Icon: "👹", Stats: app.Stats{Str: 2, Dex: 2, Con: 3, Wil: 5, Int: 4, Sur: 2}},
	{Name: "Hydra", Biome: []string{"swamp"}, CR: 5, HP: 38, MaxHP: 38, Attack: 6, Armor: 2, XPValue: 48, LootGold: 20, Icon: "🐉", Stats: app.Stats{Str: 5, Dex: 1, Con: 5, Wil: 2, Int: -1, Sur: 1}},
	{Name: "Shambling Mound", Biome: []string{"swamp"}, CR: 5, HP: 36, MaxHP: 36, Attack: 5, Armor: 3, XPValue: 45, LootGold: 16, Icon: "🌿", Stats: app.Stats{Str: 4, Dex: -1, Con: 5, Wil: 2, Int: -2, Sur: 0}},
	{Name: "Pelagic Abomination", Biome: []string{"swamp"}, CR: 5, HP: 32, MaxHP: 32, Attack: 6, Armor: 3, XPValue: 46, LootGold: 18, Icon: "🐙", Stats: app.Stats{Str: 4, Dex: 2, Con: 4, Wil: 3, Int: 1, Sur: 2}},
	{Name: "The Prophet", Biome: []string{"swamp"}, CR: 5, HP: 28, MaxHP: 28, Attack: 7, Armor: 2, XPValue: 50, LootGold: 25, Icon: "📿", Stats: app.Stats{Str: 1, Dex: 1, Con: 2, Wil: 5, Int: 5, Sur: 3}},
	{Name: "Beholder Spawn", Biome: []string{"swamp"}, CR: 5, HP: 30, MaxHP: 30, Attack: 7, Armor: 2, XPValue: 48, LootGold: 22, Icon: "👁️", Stats: app.Stats{Str: 1, Dex: 3, Con: 3, Wil: 4, Int: 5, Sur: 3}},
	{Name: "Drowned Crew", Biome: []string{"swamp"}, CR: 5, HP: 34, MaxHP: 34, Attack: 5, Armor: 3, XPValue: 45, LootGold: 20, Icon: "⚓", Stats: app.Stats{Str: 3, Dex: 2, Con: 4, Wil: 2, Int: 0, Sur: 1}},
	{Name: "Rot Wyrm", Biome: []string{"swamp"}, CR: 5, HP: 35, MaxHP: 35, Attack: 6, Armor: 3, XPValue: 47, LootGold: 19, Icon: "🐛", Stats: app.Stats{Str: 5, Dex: 1, Con: 5, Wil: 1, Int: -2, Sur: 1}},
	{Name: "Amalgam", Biome: []string{"swamp"}, CR: 5, HP: 32, MaxHP: 32, Attack: 6, Armor: 2, XPValue: 46, LootGold: 18, Icon: "🧬", Stats: app.Stats{Str: 4, Dex: 2, Con: 4, Wil: 2, Int: 1, Sur: 2}},

	// ── Plains Biome ─────────────────────────────────────────────────────

	// CR 1 — Plains
	{Name: "Prairie Dog", Biome: []string{"plains"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐿️", Stats: app.Stats{Str: -3, Dex: 2, Con: -2, Wil: -2, Int: -3, Sur: 2}},
	{Name: "Field Mouse", Biome: []string{"plains"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐁", Stats: app.Stats{Str: -3, Dex: 3, Con: -2, Wil: -2, Int: -3, Sur: 1}},
	{Name: "Crow", Biome: []string{"plains", "forest"}, Adaptable: true, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐦‍⬛", Stats: app.Stats{Str: -3, Dex: 2, Con: -2, Wil: -1, Int: -2, Sur: 2}},
	{Name: "Grass Snake", Biome: []string{"plains"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐍", Stats: app.Stats{Str: -2, Dex: 2, Con: -1, Wil: -2, Int: -3, Sur: 1}},
	{Name: "Coyote Pup", Biome: []string{"plains"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🦊", Stats: app.Stats{Str: -1, Dex: 2, Con: -1, Wil: -1, Int: -2, Sur: 2}},
	{Name: "Dust Devil", Biome: []string{"plains", "desert"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🌪️", Stats: app.Stats{Str: -3, Dex: 3, Con: -3, Wil: 0, Int: -3, Sur: 0}},
	{Name: "Locust", Biome: []string{"plains"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🦗", Stats: app.Stats{Str: -4, Dex: 2, Con: -2, Wil: -3, Int: -5, Sur: 0}},
	{Name: "Jack Rabbit", Biome: []string{"plains"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐇", Stats: app.Stats{Str: -2, Dex: 3, Con: -1, Wil: -2, Int: -3, Sur: 2}},
	{Name: "Hawk Fledgling", Biome: []string{"plains"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🦅", Stats: app.Stats{Str: -2, Dex: 3, Con: -2, Wil: -1, Int: -2, Sur: 2}},
	{Name: "Fire Ant Swarm", Biome: []string{"plains"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐜", Stats: app.Stats{Str: -3, Dex: 1, Con: 0, Wil: -2, Int: -4, Sur: 0}},

	// CR 4 — Plains
	{Name: "Coyote", Biome: []string{"plains"}, CR: 4, HP: 10, MaxHP: 10, Attack: 3, Armor: 0, XPValue: 10, LootGold: 3, Icon: "🦊", Stats: app.Stats{Str: 2, Dex: 3, Con: 1, Wil: 0, Int: -1, Sur: 2}},
	{Name: "Jaw Worm", Biome: []string{"plains"}, CR: 4, HP: 12, MaxHP: 12, Attack: 3, Armor: 1, XPValue: 12, LootGold: 4, Icon: "🪱", Stats: app.Stats{Str: 2, Dex: 1, Con: 3, Wil: 0, Int: -3, Sur: 0}},
	{Name: "Centaur Raider", Biome: []string{"plains"}, CR: 4, HP: 14, MaxHP: 14, Attack: 4, Armor: 1, XPValue: 16, LootGold: 6, Icon: "🏇", Stats: app.Stats{Str: 3, Dex: 3, Con: 2, Wil: 0, Int: 0, Sur: 1}},
	{Name: "Gnoll Hunter", Biome: []string{"plains"}, CR: 4, HP: 12, MaxHP: 12, Attack: 3, Armor: 1, XPValue: 12, LootGold: 5, Icon: "🐺", Stats: app.Stats{Str: 2, Dex: 2, Con: 2, Wil: -1, Int: -1, Sur: 2}},
	{Name: "Plains Slaver", Biome: []string{"plains"}, CR: 4, HP: 10, MaxHP: 10, Attack: 3, Armor: 1, XPValue: 11, LootGold: 7, Icon: "⛓️", Stats: app.Stats{Str: 2, Dex: 2, Con: 1, Wil: 1, Int: 0, Sur: 1}},
	{Name: "Bulette Pup", Biome: []string{"plains"}, CR: 4, HP: 18, MaxHP: 18, Attack: 4, Armor: 2, XPValue: 20, LootGold: 5, Icon: "🦈", Stats: app.Stats{Str: 4, Dex: 1, Con: 3, Wil: 0, Int: -3, Sur: 0}},
	{Name: "Warhawk", Biome: []string{"plains", "mountains"}, CR: 4, HP: 8, MaxHP: 8, Attack: 4, Armor: 0, XPValue: 10, LootGold: 3, Icon: "🦅", Stats: app.Stats{Str: 2, Dex: 4, Con: 1, Wil: 0, Int: -1, Sur: 3}},
	{Name: "Dust Wraith", Biome: []string{"plains"}, CR: 4, HP: 10, MaxHP: 10, Attack: 3, Armor: 1, XPValue: 12, LootGold: 4, Icon: "💨", Stats: app.Stats{Str: 1, Dex: 3, Con: 1, Wil: 2, Int: 0, Sur: 2}},
	{Name: "Kobold Ambusher", Biome: []string{"plains"}, CR: 4, HP: 6, MaxHP: 6, Attack: 2, Armor: 0, XPValue: 6, LootGold: 3, Icon: "🦎", Stats: app.Stats{Str: 0, Dex: 3, Con: 0, Wil: -1, Int: 1, Sur: 2}},
	{Name: "Stampeding Bull", Biome: []string{"plains"}, CR: 4, HP: 20, MaxHP: 20, Attack: 5, Armor: 1, XPValue: 22, LootGold: 4, Icon: "🐂", Stats: app.Stats{Str: 4, Dex: 2, Con: 3, Wil: -1, Int: -3, Sur: 0}},

	// CR 5 — Plains
	{Name: "Bulette", Biome: []string{"plains"}, CR: 5, HP: 36, MaxHP: 36, Attack: 6, Armor: 4, XPValue: 48, LootGold: 18, Icon: "🦈", Stats: app.Stats{Str: 5, Dex: 2, Con: 5, Wil: 0, Int: -3, Sur: 0}},
	{Name: "Warlord", Biome: []string{"plains"}, CR: 5, HP: 30, MaxHP: 30, Attack: 6, Armor: 3, XPValue: 46, LootGold: 22, Icon: "⚔️", Stats: app.Stats{Str: 4, Dex: 3, Con: 3, Wil: 3, Int: 2, Sur: 1}},
	{Name: "Gorgon", Biome: []string{"plains"}, CR: 5, HP: 38, MaxHP: 38, Attack: 6, Armor: 4, XPValue: 50, LootGold: 20, Icon: "🐃", Stats: app.Stats{Str: 5, Dex: 1, Con: 5, Wil: 2, Int: -1, Sur: 0}},
	{Name: "Chimera", Biome: []string{"plains"}, CR: 5, HP: 34, MaxHP: 34, Attack: 7, Armor: 2, XPValue: 48, LootGold: 20, Icon: "🦁", Stats: app.Stats{Str: 4, Dex: 3, Con: 4, Wil: 2, Int: 0, Sur: 2}},
	{Name: "Lagavulin", Biome: []string{"plains"}, CR: 5, HP: 40, MaxHP: 40, Attack: 5, Armor: 4, XPValue: 50, LootGold: 18, Icon: "😈", Stats: app.Stats{Str: 5, Dex: 0, Con: 5, Wil: 3, Int: 0, Sur: 0}},
	{Name: "War Mammoth", Biome: []string{"plains"}, CR: 5, HP: 40, MaxHP: 40, Attack: 6, Armor: 3, XPValue: 50, LootGold: 16, Icon: "🦣", Stats: app.Stats{Str: 5, Dex: -1, Con: 5, Wil: 2, Int: -2, Sur: 0}},
	{Name: "Giant Scorpion", Biome: []string{"plains", "desert"}, CR: 5, HP: 32, MaxHP: 32, Attack: 6, Armor: 3, XPValue: 45, LootGold: 17, Icon: "🦂", Stats: app.Stats{Str: 4, Dex: 3, Con: 4, Wil: 0, Int: -2, Sur: 1}},
	{Name: "Gnoll Overlord", Biome: []string{"plains"}, CR: 5, HP: 30, MaxHP: 30, Attack: 6, Armor: 2, XPValue: 45, LootGold: 20, Icon: "🐺", Stats: app.Stats{Str: 4, Dex: 3, Con: 3, Wil: 2, Int: 1, Sur: 2}},
	{Name: "Reptomancer", Biome: []string{"plains"}, CR: 5, HP: 28, MaxHP: 28, Attack: 7, Armor: 2, XPValue: 48, LootGold: 24, Icon: "🐊", Stats: app.Stats{Str: 2, Dex: 3, Con: 3, Wil: 4, Int: 4, Sur: 2}},
	{Name: "Thunderbird", Biome: []string{"plains"}, CR: 5, HP: 32, MaxHP: 32, Attack: 7, Armor: 2, XPValue: 48, LootGold: 20, Icon: "⚡", Stats: app.Stats{Str: 4, Dex: 4, Con: 3, Wil: 3, Int: 1, Sur: 3}},

	// ── Mountains Biome ──────────────────────────────────────────────────

	// CR 1 — Mountains
	{Name: "Snow Hare", Biome: []string{"mountains"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐇", Stats: app.Stats{Str: -3, Dex: 3, Con: -1, Wil: -2, Int: -3, Sur: 2}},
	{Name: "Rock Beetle", Biome: []string{"mountains"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🪲", Stats: app.Stats{Str: -1, Dex: -1, Con: 0, Wil: -2, Int: -4, Sur: 0}},
	{Name: "Ice Mote", Biome: []string{"mountains"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "❄️", Stats: app.Stats{Str: -4, Dex: 2, Con: -2, Wil: 0, Int: -2, Sur: 1}},
	{Name: "Cave Bat", Biome: []string{"mountains"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🦇", Stats: app.Stats{Str: -3, Dex: 2, Con: -2, Wil: -1, Int: -3, Sur: 1}},
	{Name: "Pebbleling", Biome: []string{"mountains"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🪨", Stats: app.Stats{Str: -1, Dex: -2, Con: 1, Wil: -2, Int: -4, Sur: -1}},
	{Name: "Wind Sprite", Biome: []string{"mountains"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "💨", Stats: app.Stats{Str: -4, Dex: 3, Con: -3, Wil: 0, Int: -1, Sur: 1}},
	{Name: "Cliff Swallow", Biome: []string{"mountains"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐦", Stats: app.Stats{Str: -3, Dex: 3, Con: -2, Wil: -1, Int: -2, Sur: 2}},
	{Name: "Kobold Runt", Biome: []string{"mountains"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🦎", Stats: app.Stats{Str: -2, Dex: 1, Con: -1, Wil: -1, Int: -1, Sur: 1}},
	{Name: "Frost Grub", Biome: []string{"mountains"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐛", Stats: app.Stats{Str: -3, Dex: -1, Con: 0, Wil: -2, Int: -5, Sur: -1}},
	{Name: "Mountain Viper", Biome: []string{"mountains"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐍", Stats: app.Stats{Str: -2, Dex: 3, Con: -1, Wil: -1, Int: -3, Sur: 2}},

	// CR 4 — Mountains
	{Name: "Arachnid", Biome: []string{"mountains"}, CR: 4, HP: 12, MaxHP: 12, Attack: 3, Armor: 1, XPValue: 12, LootGold: 5, Icon: "🕷️", Stats: app.Stats{Str: 1, Dex: 3, Con: 1, Wil: 0, Int: -2, Sur: 2}},
	{Name: "Lurker", Biome: []string{"mountains"}, CR: 4, HP: 18, MaxHP: 18, Attack: 4, Armor: 2, XPValue: 20, LootGold: 8, Icon: "👥", Stats: app.Stats{Str: 3, Dex: 2, Con: 3, Wil: 1, Int: 0, Sur: 2}},
	{Name: "Troll", Biome: []string{"mountains"}, CR: 4, HP: 26, MaxHP: 26, Attack: 5, Armor: 3, XPValue: 30, LootGold: 12, Icon: "🧌", Stats: app.Stats{Str: 4, Dex: 0, Con: 4, Wil: 1, Int: -2, Sur: 1}},
	{Name: "Bugbear", Biome: []string{"mountains"}, CR: 4, HP: 14, MaxHP: 14, Attack: 4, Armor: 1, XPValue: 16, LootGold: 6, Icon: "👹", Stats: app.Stats{Str: 3, Dex: 2, Con: 2, Wil: 0, Int: -1, Sur: 2}},
	{Name: "Harpy", Biome: []string{"mountains"}, CR: 4, HP: 10, MaxHP: 10, Attack: 3, Armor: 0, XPValue: 12, LootGold: 4, Icon: "🦅", Stats: app.Stats{Str: 1, Dex: 3, Con: 1, Wil: 2, Int: 0, Sur: 2}},
	{Name: "Stone Golem", Biome: []string{"mountains"}, CR: 4, HP: 22, MaxHP: 22, Attack: 4, Armor: 3, XPValue: 25, LootGold: 8, Icon: "🗿", Stats: app.Stats{Str: 4, Dex: -1, Con: 4, Wil: 1, Int: -2, Sur: 0}},
	{Name: "Frost Gremlin", Biome: []string{"mountains"}, CR: 4, HP: 8, MaxHP: 8, Attack: 3, Armor: 0, XPValue: 8, LootGold: 4, Icon: "👺", Stats: app.Stats{Str: 1, Dex: 3, Con: 0, Wil: 1, Int: 1, Sur: 2}},
	{Name: "Yeti Scout", Biome: []string{"mountains"}, CR: 4, HP: 16, MaxHP: 16, Attack: 4, Armor: 2, XPValue: 18, LootGold: 6, Icon: "🦍", Stats: app.Stats{Str: 3, Dex: 1, Con: 3, Wil: 1, Int: -2, Sur: 1}},
	{Name: "Wyvern Hatchling", Biome: []string{"mountains"}, CR: 4, HP: 14, MaxHP: 14, Attack: 4, Armor: 1, XPValue: 16, LootGold: 7, Icon: "🐉", Stats: app.Stats{Str: 3, Dex: 3, Con: 2, Wil: 0, Int: -1, Sur: 2}},
	{Name: "Cave Ogre", Biome: []string{"mountains"}, CR: 4, HP: 24, MaxHP: 24, Attack: 5, Armor: 2, XPValue: 28, LootGold: 10, Icon: "👹", Stats: app.Stats{Str: 4, Dex: 0, Con: 4, Wil: 0, Int: -2, Sur: 0}},

	// CR 5 — Mountains
	{Name: "Frost Giant", Biome: []string{"mountains"}, CR: 5, HP: 40, MaxHP: 40, Attack: 7, Armor: 3, XPValue: 50, LootGold: 22, Icon: "🧊", Stats: app.Stats{Str: 5, Dex: 0, Con: 5, Wil: 2, Int: 0, Sur: 0}},
	{Name: "Elder Wyvern", Biome: []string{"mountains"}, CR: 5, HP: 34, MaxHP: 34, Attack: 6, Armor: 3, XPValue: 48, LootGold: 20, Icon: "🐉", Stats: app.Stats{Str: 4, Dex: 4, Con: 4, Wil: 2, Int: 0, Sur: 2}},
	{Name: "Roc", Biome: []string{"mountains"}, CR: 5, HP: 36, MaxHP: 36, Attack: 7, Armor: 2, XPValue: 48, LootGold: 18, Icon: "🦅", Stats: app.Stats{Str: 5, Dex: 3, Con: 4, Wil: 1, Int: -1, Sur: 2}},
	{Name: "Ice King Revenant", Biome: []string{"mountains"}, CR: 5, HP: 30, MaxHP: 30, Attack: 7, Armor: 3, XPValue: 50, LootGold: 25, Icon: "👑", Stats: app.Stats{Str: 3, Dex: 2, Con: 3, Wil: 5, Int: 4, Sur: 2}},
	{Name: "Bronze Automaton", Biome: []string{"mountains"}, CR: 5, HP: 38, MaxHP: 38, Attack: 6, Armor: 4, XPValue: 50, LootGold: 20, Icon: "🤖", Stats: app.Stats{Str: 5, Dex: 1, Con: 5, Wil: 2, Int: 2, Sur: 0}},
	{Name: "Mountain Behemoth", Biome: []string{"mountains"}, CR: 5, HP: 40, MaxHP: 40, Attack: 6, Armor: 4, XPValue: 50, LootGold: 18, Icon: "🦏", Stats: app.Stats{Str: 5, Dex: -1, Con: 5, Wil: 2, Int: -2, Sur: 0}},
	{Name: "Yeti Alpha", Biome: []string{"mountains"}, CR: 5, HP: 36, MaxHP: 36, Attack: 6, Armor: 3, XPValue: 48, LootGold: 20, Icon: "🦍", Stats: app.Stats{Str: 5, Dex: 2, Con: 5, Wil: 2, Int: -1, Sur: 1}},
	{Name: "Awakened Stone", Biome: []string{"mountains"}, CR: 5, HP: 38, MaxHP: 38, Attack: 5, Armor: 4, XPValue: 48, LootGold: 16, Icon: "🗿", Stats: app.Stats{Str: 5, Dex: -2, Con: 5, Wil: 3, Int: 1, Sur: 0}},
	{Name: "Snecko", Biome: []string{"mountains"}, CR: 5, HP: 32, MaxHP: 32, Attack: 6, Armor: 2, XPValue: 46, LootGold: 20, Icon: "🐍", Stats: app.Stats{Str: 3, Dex: 4, Con: 3, Wil: 3, Int: 2, Sur: 3}},
	{Name: "Mind Flayer", Biome: []string{"mountains"}, CR: 5, HP: 28, MaxHP: 28, Attack: 7, Armor: 2, XPValue: 50, LootGold: 25, Icon: "🧠", Stats: app.Stats{Str: 2, Dex: 2, Con: 3, Wil: 5, Int: 5, Sur: 2}},

	// ── Island Biome ─────────────────────────────────────────────────────

	// CR 1 — Island
	{Name: "Hermit Crab", Biome: []string{"island"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🦀", Stats: app.Stats{Str: -2, Dex: 0, Con: 0, Wil: -2, Int: -4, Sur: 0}},
	{Name: "Reef Fish", Biome: []string{"island"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐟", Stats: app.Stats{Str: -3, Dex: 2, Con: -2, Wil: -2, Int: -4, Sur: 1}},
	{Name: "Sand Flea", Biome: []string{"island"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🦟", Stats: app.Stats{Str: -4, Dex: 3, Con: -3, Wil: -2, Int: -5, Sur: 1}},
	{Name: "Coconut Imp", Biome: []string{"island"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🥥", Stats: app.Stats{Str: -2, Dex: 1, Con: -1, Wil: -1, Int: -2, Sur: 1}},
	{Name: "Sea Slug", Biome: []string{"island"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐌", Stats: app.Stats{Str: -3, Dex: -2, Con: 0, Wil: -2, Int: -5, Sur: -1}},
	{Name: "Jungle Rat", Biome: []string{"island"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐀", Stats: app.Stats{Str: -2, Dex: 2, Con: -2, Wil: -2, Int: -3, Sur: 1}},
	{Name: "Pelican", Biome: []string{"island"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐦", Stats: app.Stats{Str: -2, Dex: 2, Con: -1, Wil: -1, Int: -3, Sur: 1}},
	{Name: "Starfish", Biome: []string{"island"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "⭐", Stats: app.Stats{Str: -3, Dex: -2, Con: 1, Wil: -2, Int: -5, Sur: -1}},
	{Name: "Tidepool Eel", Biome: []string{"island"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐍", Stats: app.Stats{Str: -2, Dex: 2, Con: -1, Wil: -1, Int: -3, Sur: 1}},
	{Name: "Barnacle", Biome: []string{"island"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐚", Stats: app.Stats{Str: -3, Dex: -3, Con: 1, Wil: -2, Int: -5, Sur: -1}},

	// CR 4 — Island
	{Name: "Pirate Corsair", Biome: []string{"island"}, CR: 4, HP: 10, MaxHP: 10, Attack: 3, Armor: 1, XPValue: 10, LootGold: 8, Icon: "🏴‍☠️", Stats: app.Stats{Str: 2, Dex: 3, Con: 1, Wil: 0, Int: 0, Sur: 2}},
	{Name: "Lake Knight", Biome: []string{"island"}, CR: 4, HP: 14, MaxHP: 14, Attack: 3, Armor: 2, XPValue: 14, LootGold: 6, Icon: "🛡️", Stats: app.Stats{Str: 2, Dex: 1, Con: 2, Wil: 2, Int: 1, Sur: 1}},
	{Name: "Reef Lurker", Biome: []string{"island"}, CR: 4, HP: 12, MaxHP: 12, Attack: 3, Armor: 1, XPValue: 12, LootGold: 5, Icon: "🦑", Stats: app.Stats{Str: 2, Dex: 2, Con: 2, Wil: 1, Int: -1, Sur: 2}},
	{Name: "Siren", Biome: []string{"island"}, CR: 4, HP: 8, MaxHP: 8, Attack: 4, Armor: 0, XPValue: 12, LootGold: 6, Icon: "🧜", Stats: app.Stats{Str: 0, Dex: 2, Con: 1, Wil: 4, Int: 2, Sur: 1}},
	{Name: "Cannon Crab", Biome: []string{"island"}, CR: 4, HP: 16, MaxHP: 16, Attack: 4, Armor: 2, XPValue: 18, LootGold: 5, Icon: "🦞", Stats: app.Stats{Str: 3, Dex: 0, Con: 3, Wil: 0, Int: -2, Sur: 0}},
	{Name: "Sand Golem", Biome: []string{"island", "desert"}, CR: 4, HP: 20, MaxHP: 20, Attack: 3, Armor: 3, XPValue: 18, LootGold: 5, Icon: "🏖️", Stats: app.Stats{Str: 3, Dex: -1, Con: 4, Wil: 1, Int: -3, Sur: 0}},
	{Name: "Jungle Stalker", Biome: []string{"island"}, CR: 4, HP: 12, MaxHP: 12, Attack: 4, Armor: 0, XPValue: 14, LootGold: 4, Icon: "🐆", Stats: app.Stats{Str: 2, Dex: 4, Con: 1, Wil: 0, Int: -1, Sur: 3}},
	{Name: "Shipwreck Phantom", Biome: []string{"island"}, CR: 4, HP: 10, MaxHP: 10, Attack: 3, Armor: 1, XPValue: 12, LootGold: 7, Icon: "👻", Stats: app.Stats{Str: 1, Dex: 2, Con: 1, Wil: 3, Int: 1, Sur: 2}},
	{Name: "Tiki Brute", Biome: []string{"island"}, CR: 4, HP: 18, MaxHP: 18, Attack: 4, Armor: 2, XPValue: 20, LootGold: 6, Icon: "🗿", Stats: app.Stats{Str: 3, Dex: 1, Con: 3, Wil: 1, Int: -1, Sur: 1}},
	{Name: "Sahuagin", Biome: []string{"island"}, CR: 4, HP: 14, MaxHP: 14, Attack: 3, Armor: 1, XPValue: 14, LootGold: 5, Icon: "🐟", Stats: app.Stats{Str: 2, Dex: 2, Con: 2, Wil: 1, Int: 0, Sur: 2}},

	// CR 5 — Island
	{Name: "Kraken Tendril", Biome: []string{"island"}, CR: 5, HP: 36, MaxHP: 36, Attack: 6, Armor: 3, XPValue: 48, LootGold: 22, Icon: "🐙", Stats: app.Stats{Str: 5, Dex: 3, Con: 4, Wil: 2, Int: 1, Sur: 2}},
	{Name: "Pirate King", Biome: []string{"island"}, CR: 5, HP: 30, MaxHP: 30, Attack: 6, Armor: 3, XPValue: 46, LootGold: 25, Icon: "☠️", Stats: app.Stats{Str: 4, Dex: 3, Con: 3, Wil: 2, Int: 2, Sur: 2}},
	{Name: "Sea Serpent", Biome: []string{"island"}, CR: 5, HP: 38, MaxHP: 38, Attack: 6, Armor: 3, XPValue: 48, LootGold: 20, Icon: "🐍", Stats: app.Stats{Str: 5, Dex: 3, Con: 4, Wil: 1, Int: -1, Sur: 2}},
	{Name: "Leviathan Spawn", Biome: []string{"island"}, CR: 5, HP: 40, MaxHP: 40, Attack: 7, Armor: 3, XPValue: 50, LootGold: 22, Icon: "🐋", Stats: app.Stats{Str: 5, Dex: 1, Con: 5, Wil: 2, Int: 0, Sur: 1}},
	{Name: "Volcanic Elemental", Biome: []string{"island"}, CR: 5, HP: 34, MaxHP: 34, Attack: 7, Armor: 2, XPValue: 48, LootGold: 18, Icon: "🌋", Stats: app.Stats{Str: 4, Dex: 2, Con: 4, Wil: 3, Int: 2, Sur: 1}},
	{Name: "Aboleth", Biome: []string{"island"}, CR: 5, HP: 32, MaxHP: 32, Attack: 6, Armor: 2, XPValue: 50, LootGold: 24, Icon: "🧠", Stats: app.Stats{Str: 3, Dex: 2, Con: 3, Wil: 5, Int: 5, Sur: 2}},
	{Name: "Coral Titan", Biome: []string{"island"}, CR: 5, HP: 38, MaxHP: 38, Attack: 5, Armor: 4, XPValue: 48, LootGold: 18, Icon: "🪸", Stats: app.Stats{Str: 5, Dex: -1, Con: 5, Wil: 2, Int: 0, Sur: 0}},
	{Name: "Tempest Djinn", Biome: []string{"island"}, CR: 5, HP: 28, MaxHP: 28, Attack: 7, Armor: 2, XPValue: 48, LootGold: 22, Icon: "🌊", Stats: app.Stats{Str: 3, Dex: 4, Con: 3, Wil: 4, Int: 3, Sur: 3}},
	{Name: "Ancient Megalodon", Biome: []string{"island"}, CR: 5, HP: 40, MaxHP: 40, Attack: 7, Armor: 3, XPValue: 50, LootGold: 20, Icon: "🦈", Stats: app.Stats{Str: 5, Dex: 2, Con: 5, Wil: 1, Int: -2, Sur: 1}},
	{Name: "Ghost Ship Captain", Biome: []string{"island"}, CR: 5, HP: 30, MaxHP: 30, Attack: 6, Armor: 3, XPValue: 46, LootGold: 24, Icon: "⚓", Stats: app.Stats{Str: 3, Dex: 3, Con: 3, Wil: 4, Int: 3, Sur: 2}},

	// ── Desert Biome ─────────────────────────────────────────────────────

	// CR 1 — Desert
	{Name: "Sprite", Biome: []string{"desert"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "💨", Stats: app.Stats{Str: -3, Dex: 3, Con: -2, Wil: 0, Int: 0, Sur: 1}},
	{Name: "Sand Beetle", Biome: []string{"desert"}, Adaptable: true, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🪲", Stats: app.Stats{Str: -2, Dex: -1, Con: -1, Wil: -2, Int: -4, Sur: 0}},
	{Name: "Scarab", Biome: []string{"desert"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🪲", Stats: app.Stats{Str: -3, Dex: 1, Con: 0, Wil: -2, Int: -4, Sur: 0}},
	{Name: "Sidewinder", Biome: []string{"desert"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐍", Stats: app.Stats{Str: -2, Dex: 3, Con: -1, Wil: -1, Int: -3, Sur: 2}},
	{Name: "Vulture Chick", Biome: []string{"desert"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🦅", Stats: app.Stats{Str: -2, Dex: 2, Con: -2, Wil: -1, Int: -3, Sur: 1}},
	{Name: "Sand Flea", Biome: []string{"desert"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🦟", Stats: app.Stats{Str: -4, Dex: 2, Con: -2, Wil: -2, Int: -5, Sur: 0}},
	{Name: "Dust Mephit", Biome: []string{"desert"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🌪️", Stats: app.Stats{Str: -3, Dex: 2, Con: -2, Wil: 0, Int: -1, Sur: 1}},
	{Name: "Scorpion", Biome: []string{"desert"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🦂", Stats: app.Stats{Str: -2, Dex: 1, Con: 0, Wil: -2, Int: -4, Sur: 0}},
	{Name: "Cactus Imp", Biome: []string{"desert"}, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🌵", Stats: app.Stats{Str: -2, Dex: 1, Con: 0, Wil: -1, Int: -2, Sur: 1}},
	{Name: "Tomb Scarab", Biome: []string{"desert"}, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "💀", Stats: app.Stats{Str: -3, Dex: 1, Con: -1, Wil: 0, Int: -3, Sur: 0}},

	// CR 4 — Desert
	{Name: "Goblin", Biome: []string{"desert"}, CR: 4, HP: 8, MaxHP: 8, Attack: 3, Armor: 0, XPValue: 8, LootGold: 4, Icon: "👺", Stats: app.Stats{Str: 0, Dex: 2, Con: 0, Wil: -1, Int: 0, Sur: 1}},
	{Name: "Elemental", Biome: []string{"desert"}, CR: 4, HP: 20, MaxHP: 20, Attack: 4, Armor: 2, XPValue: 22, LootGold: 10, Icon: "🔥", Stats: app.Stats{Str: 3, Dex: 2, Con: 3, Wil: 2, Int: 1, Sur: 1}},
	{Name: "Hellhound", Biome: []string{"desert"}, CR: 4, HP: 18, MaxHP: 18, Attack: 5, Armor: 1, XPValue: 25, LootGold: 8, Icon: "🐺", Stats: app.Stats{Str: 3, Dex: 3, Con: 2, Wil: 1, Int: -1, Sur: 2}},
	{Name: "Mummy", Biome: []string{"desert"}, CR: 4, HP: 16, MaxHP: 16, Attack: 3, Armor: 2, XPValue: 16, LootGold: 8, Icon: "🧟", Stats: app.Stats{Str: 3, Dex: 0, Con: 3, Wil: 2, Int: 0, Sur: 0}},
	{Name: "Sand Wraith", Biome: []string{"desert"}, CR: 4, HP: 10, MaxHP: 10, Attack: 4, Armor: 0, XPValue: 14, LootGold: 6, Icon: "👻", Stats: app.Stats{Str: 1, Dex: 3, Con: 1, Wil: 3, Int: 1, Sur: 2}},
	{Name: "Tomb Guardian", Biome: []string{"desert"}, CR: 4, HP: 22, MaxHP: 22, Attack: 3, Armor: 3, XPValue: 22, LootGold: 8, Icon: "⚱️", Stats: app.Stats{Str: 3, Dex: 0, Con: 4, Wil: 2, Int: -1, Sur: 0}},
	{Name: "Fire Wolf", Biome: []string{"desert"}, CR: 4, HP: 12, MaxHP: 12, Attack: 4, Armor: 0, XPValue: 14, LootGold: 5, Icon: "🔥", Stats: app.Stats{Str: 2, Dex: 3, Con: 2, Wil: 1, Int: -1, Sur: 2}},
	{Name: "Gnoll Raider", Biome: []string{"desert", "plains"}, CR: 4, HP: 12, MaxHP: 12, Attack: 3, Armor: 1, XPValue: 12, LootGold: 6, Icon: "🗡️", Stats: app.Stats{Str: 2, Dex: 2, Con: 2, Wil: 0, Int: -1, Sur: 1}},
	{Name: "Dust Stalker", Biome: []string{"desert"}, CR: 4, HP: 14, MaxHP: 14, Attack: 3, Armor: 1, XPValue: 14, LootGold: 5, Icon: "🌫️", Stats: app.Stats{Str: 2, Dex: 3, Con: 2, Wil: 1, Int: 0, Sur: 3}},
	{Name: "Anubite", Biome: []string{"desert"}, CR: 4, HP: 14, MaxHP: 14, Attack: 4, Armor: 1, XPValue: 16, LootGold: 7, Icon: "🐕", Stats: app.Stats{Str: 3, Dex: 2, Con: 2, Wil: 2, Int: 1, Sur: 1}},

	// CR 5 — Desert
	{Name: "Drake", Biome: []string{"desert"}, CR: 5, HP: 32, MaxHP: 32, Attack: 6, Armor: 3, XPValue: 45, LootGold: 20, Icon: "🐉", Stats: app.Stats{Str: 5, Dex: 2, Con: 4, Wil: 3, Int: 2, Sur: 2}},
	{Name: "Sphinx", Biome: []string{"desert"}, CR: 5, HP: 30, MaxHP: 30, Attack: 6, Armor: 3, XPValue: 48, LootGold: 24, Icon: "🦁", Stats: app.Stats{Str: 4, Dex: 3, Con: 3, Wil: 4, Int: 5, Sur: 2}},
	{Name: "Mummy Lord", Biome: []string{"desert"}, CR: 5, HP: 34, MaxHP: 34, Attack: 6, Armor: 3, XPValue: 48, LootGold: 22, Icon: "👑", Stats: app.Stats{Str: 4, Dex: 1, Con: 4, Wil: 4, Int: 3, Sur: 1}},
	{Name: "Sandstorm Djinn", Biome: []string{"desert"}, CR: 5, HP: 28, MaxHP: 28, Attack: 7, Armor: 2, XPValue: 50, LootGold: 22, Icon: "🌪️", Stats: app.Stats{Str: 3, Dex: 4, Con: 3, Wil: 4, Int: 3, Sur: 3}},
	{Name: "Bone Colossus", Biome: []string{"desert"}, CR: 5, HP: 40, MaxHP: 40, Attack: 6, Armor: 4, XPValue: 50, LootGold: 20, Icon: "💀", Stats: app.Stats{Str: 5, Dex: 0, Con: 5, Wil: 2, Int: -1, Sur: 0}},
	{Name: "Efreeti", Biome: []string{"desert"}, CR: 5, HP: 32, MaxHP: 32, Attack: 7, Armor: 2, XPValue: 48, LootGold: 22, Icon: "🔥", Stats: app.Stats{Str: 4, Dex: 3, Con: 4, Wil: 3, Int: 3, Sur: 2}},
	{Name: "Purple Worm", Biome: []string{"desert"}, CR: 5, HP: 40, MaxHP: 40, Attack: 7, Armor: 3, XPValue: 50, LootGold: 18, Icon: "🪱", Stats: app.Stats{Str: 5, Dex: 0, Con: 5, Wil: 1, Int: -3, Sur: 0}},
	{Name: "Pharaoh's Shade", Biome: []string{"desert"}, CR: 5, HP: 30, MaxHP: 30, Attack: 6, Armor: 2, XPValue: 46, LootGold: 24, Icon: "👤", Stats: app.Stats{Str: 3, Dex: 3, Con: 3, Wil: 5, Int: 4, Sur: 2}},
	{Name: "Awakened One", Biome: []string{"desert"}, CR: 5, HP: 36, MaxHP: 36, Attack: 6, Armor: 3, XPValue: 48, LootGold: 20, Icon: "🗿", Stats: app.Stats{Str: 4, Dex: 2, Con: 4, Wil: 4, Int: 3, Sur: 1}},
	{Name: "Flame Lich", Biome: []string{"desert"}, CR: 5, HP: 28, MaxHP: 28, Attack: 7, Armor: 2, XPValue: 50, LootGold: 25, Icon: "🔥", Stats: app.Stats{Str: 2, Dex: 2, Con: 3, Wil: 5, Int: 5, Sur: 2}},

	// ── Any Biome (Adaptable) ────────────────────────────────────────────

	// CR 1 — Any
	{Name: "Rat", Biome: []string{"Any"}, Adaptable: true, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐀", Stats: app.Stats{Str: -2, Dex: 1, Con: -2, Wil: -2, Int: -3, Sur: 0}},
	{Name: "Slime", Biome: []string{"Any"}, Adaptable: true, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🟢", Stats: app.Stats{Str: -2, Dex: -2, Con: 0, Wil: -3, Int: -5, Sur: -1}},
	{Name: "Packrat", Biome: []string{"Any"}, Adaptable: true, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐀", Stats: app.Stats{Str: -2, Dex: 2, Con: -1, Wil: -2, Int: -2, Sur: 1}},
	{Name: "Wisp", Biome: []string{"Any"}, Adaptable: true, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "✨", Stats: app.Stats{Str: -4, Dex: 3, Con: -3, Wil: 0, Int: 0, Sur: 1}},
	{Name: "Shadow Mite", Biome: []string{"Any"}, Adaptable: true, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🕳️", Stats: app.Stats{Str: -3, Dex: 2, Con: -2, Wil: -1, Int: -3, Sur: 1}},
	{Name: "Whelp", Biome: []string{"Any"}, Adaptable: true, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🐶", Stats: app.Stats{Str: -1, Dex: 1, Con: -1, Wil: -1, Int: -3, Sur: 1}},
	{Name: "Imp", Biome: []string{"Any"}, Adaptable: true, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "😈", Stats: app.Stats{Str: -3, Dex: 2, Con: -2, Wil: 0, Int: -1, Sur: 1}},
	{Name: "Crawling Hand", Biome: []string{"Any"}, Adaptable: true, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🖐️", Stats: app.Stats{Str: -2, Dex: 1, Con: -2, Wil: -1, Int: -4, Sur: 0}},
	{Name: "Gnat Swarm", Biome: []string{"Any"}, Adaptable: true, CR: 1, HP: 1, MaxHP: 1, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🦟", Stats: app.Stats{Str: -4, Dex: 3, Con: -3, Wil: -2, Int: -5, Sur: 0}},
	{Name: "Shroom", Biome: []string{"Any"}, Adaptable: true, CR: 1, HP: 2, MaxHP: 2, Attack: 1, Armor: 0, XPValue: 1, LootGold: 0, Icon: "🍄", Stats: app.Stats{Str: -3, Dex: -1, Con: 1, Wil: -1, Int: -4, Sur: 0}},

	// CR 4 — Any
	{Name: "Marauder", Biome: []string{"Any"}, Adaptable: true, CR: 4, HP: 9, MaxHP: 9, Attack: 2, Armor: 1, XPValue: 8, LootGold: 5, Icon: "🥷", Stats: app.Stats{Str: 1, Dex: 2, Con: 1, Wil: 0, Int: 0, Sur: 1}},
	{Name: "Gargoyle", Biome: []string{"Any"}, Adaptable: true, CR: 4, HP: 15, MaxHP: 15, Attack: 3, Armor: 3, XPValue: 15, LootGold: 6, Icon: "🗿", Stats: app.Stats{Str: 3, Dex: 1, Con: 3, Wil: 1, Int: -1, Sur: 1}},
	{Name: "Mimic", Biome: []string{"Any"}, Adaptable: true, CR: 4, HP: 14, MaxHP: 14, Attack: 4, Armor: 1, XPValue: 16, LootGold: 10, Icon: "📦", Stats: app.Stats{Str: 2, Dex: 1, Con: 2, Wil: 0, Int: 1, Sur: 0}},
	{Name: "Gelatinous Cube", Biome: []string{"Any"}, Adaptable: true, CR: 4, HP: 18, MaxHP: 18, Attack: 2, Armor: 1, XPValue: 14, LootGold: 4, Icon: "🟩", Stats: app.Stats{Str: 2, Dex: -3, Con: 4, Wil: 0, Int: -5, Sur: -1}},
	{Name: "Bandit Captain", Biome: []string{"Any"}, Adaptable: true, CR: 4, HP: 12, MaxHP: 12, Attack: 3, Armor: 1, XPValue: 12, LootGold: 8, Icon: "🗡️", Stats: app.Stats{Str: 2, Dex: 2, Con: 1, Wil: 1, Int: 1, Sur: 1}},
	{Name: "Shade", Biome: []string{"Any"}, Adaptable: true, CR: 4, HP: 10, MaxHP: 10, Attack: 3, Armor: 0, XPValue: 10, LootGold: 4, Icon: "👤", Stats: app.Stats{Str: 1, Dex: 3, Con: 1, Wil: 2, Int: 0, Sur: 2}},
	{Name: "War Hound", Biome: []string{"Any"}, Adaptable: true, CR: 4, HP: 12, MaxHP: 12, Attack: 3, Armor: 1, XPValue: 12, LootGold: 3, Icon: "🐕", Stats: app.Stats{Str: 2, Dex: 3, Con: 2, Wil: 0, Int: -2, Sur: 2}},
	{Name: "Cursed Knight", Biome: []string{"Any"}, Adaptable: true, CR: 4, HP: 16, MaxHP: 16, Attack: 4, Armor: 2, XPValue: 18, LootGold: 7, Icon: "⚔️", Stats: app.Stats{Str: 3, Dex: 1, Con: 2, Wil: 2, Int: 0, Sur: 1}},
	{Name: "Basilisk Runt", Biome: []string{"Any"}, Adaptable: true, CR: 4, HP: 14, MaxHP: 14, Attack: 3, Armor: 2, XPValue: 14, LootGold: 5, Icon: "🦎", Stats: app.Stats{Str: 2, Dex: 1, Con: 3, Wil: 1, Int: -2, Sur: 1}},
	{Name: "Candy Golem", Biome: []string{"Any"}, Adaptable: true, CR: 4, HP: 20, MaxHP: 20, Attack: 3, Armor: 2, XPValue: 16, LootGold: 6, Icon: "🍬", Stats: app.Stats{Str: 3, Dex: 0, Con: 3, Wil: 0, Int: -3, Sur: 0}},

	// CR 5 — Any
	{Name: "Beholder", Biome: []string{"Any"}, Adaptable: true, CR: 5, HP: 30, MaxHP: 30, Attack: 7, Armor: 2, XPValue: 50, LootGold: 25, Icon: "👁️", Stats: app.Stats{Str: 2, Dex: 3, Con: 3, Wil: 5, Int: 5, Sur: 3}},
	{Name: "Lich", Biome: []string{"Any"}, Adaptable: true, CR: 5, HP: 28, MaxHP: 28, Attack: 7, Armor: 2, XPValue: 50, LootGold: 25, Icon: "💀", Stats: app.Stats{Str: 1, Dex: 2, Con: 2, Wil: 5, Int: 5, Sur: 2}},
	{Name: "Vampire Lord", Biome: []string{"Any"}, Adaptable: true, CR: 5, HP: 32, MaxHP: 32, Attack: 6, Armor: 3, XPValue: 48, LootGold: 22, Icon: "🧛", Stats: app.Stats{Str: 4, Dex: 4, Con: 3, Wil: 3, Int: 3, Sur: 3}},
	{Name: "Death Knight", Biome: []string{"Any"}, Adaptable: true, CR: 5, HP: 36, MaxHP: 36, Attack: 6, Armor: 4, XPValue: 48, LootGold: 20, Icon: "⚔️", Stats: app.Stats{Str: 5, Dex: 2, Con: 4, Wil: 3, Int: 1, Sur: 1}},
	{Name: "Book of Stabbing", Biome: []string{"Any"}, Adaptable: true, CR: 5, HP: 30, MaxHP: 30, Attack: 7, Armor: 2, XPValue: 46, LootGold: 20, Icon: "📕", Stats: app.Stats{Str: 3, Dex: 3, Con: 3, Wil: 3, Int: 3, Sur: 2}},
	{Name: "Basilisk", Biome: []string{"Any"}, Adaptable: true, CR: 5, HP: 34, MaxHP: 34, Attack: 6, Armor: 3, XPValue: 46, LootGold: 18, Icon: "🦎", Stats: app.Stats{Str: 4, Dex: 2, Con: 5, Wil: 2, Int: -1, Sur: 1}},
	{Name: "Iron Golem", Biome: []string{"Any"}, Adaptable: true, CR: 5, HP: 40, MaxHP: 40, Attack: 6, Armor: 4, XPValue: 50, LootGold: 18, Icon: "🤖", Stats: app.Stats{Str: 5, Dex: -1, Con: 5, Wil: 2, Int: 0, Sur: 0}},
	{Name: "Nightwalker", Biome: []string{"Any"}, Adaptable: true, CR: 5, HP: 32, MaxHP: 32, Attack: 7, Armor: 2, XPValue: 48, LootGold: 22, Icon: "🌑", Stats: app.Stats{Str: 4, Dex: 3, Con: 3, Wil: 4, Int: 2, Sur: 2}},
	{Name: "Elder Dragon", Biome: []string{"Any"}, Adaptable: true, CR: 5, HP: 40, MaxHP: 40, Attack: 7, Armor: 4, XPValue: 50, LootGold: 25, Icon: "🐉", Stats: app.Stats{Str: 5, Dex: 2, Con: 5, Wil: 4, Int: 4, Sur: 2}},
	{Name: "The Shambler", Biome: []string{"Any"}, Adaptable: true, CR: 5, HP: 38, MaxHP: 38, Attack: 6, Armor: 3, XPValue: 50, LootGold: 22, Icon: "🌑", Stats: app.Stats{Str: 5, Dex: 1, Con: 5, Wil: 4, Int: 2, Sur: 1}},
}

// GenerateMonstersForRoom picks monsters matching biome and CR budget for a room encounter.
// If maxCR > 0, monsters with CR > maxCR will be excluded from spawning.
func GenerateMonstersForRoom(biome string, crBudget int, r *rand.Rand, maxCR int) []app.Monster {
	return GenerateMonstersForBiomeRoom(biome, crBudget, r, maxCR)
}
