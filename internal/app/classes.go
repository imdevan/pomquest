package app

import (
	"github.com/pomquest/internal/adapters/icon"
)

// ClassArchetype defines standard RPG class stats, HP, armor, and combat attributes.
type ClassArchetype struct {
	Name   string
	HP     int
	Attack int
	Armor  int
	Stats  Stats
}

// Icon returns the class icon from the central icon library.
func (c ClassArchetype) Icon(mode ...icon.Mode) string {
	m := icon.ModeNerdFont
	if len(mode) > 0 {
		m = mode[0]
	}
	return icon.GetClassIcon(c.Name, m)
}

// ClassCatalog returns the standard predefined set of available RPG class archetypes.
// each class has 4 stat points spread over up to 4 stats (preferrably 2, depending on class flavor)
func ClassCatalog() []ClassArchetype {
	return []ClassArchetype{
		{Name: "fighter", HP: 10, Attack: 3, Armor: 2, Stats: Stats{Str: 2, Con: 2}},
		{Name: "wizard", HP: 6, Attack: 4, Armor: 0, Stats: Stats{Int: 3, Wil: 1}},
		{Name: "warlock", HP: 7, Attack: 3, Armor: 1, Stats: Stats{Cha: 2, Int: 2}},
		{Name: "ranger", HP: 8, Attack: 3, Armor: 1, Stats: Stats{Dex: 2, Sur: 2}},
		{Name: "rogue", HP: 8, Attack: 3, Armor: 0, Stats: Stats{Dex: 3, Luc: 1}},
		{Name: "cleric", HP: 9, Attack: 2, Armor: 2, Stats: Stats{Wil: 2, Hel: 2}},
		{Name: "paladin", HP: 10, Attack: 3, Armor: 3, Stats: Stats{Str: 2, Con: 1, Wil: 1}},
		{Name: "bard", HP: 8, Attack: 2, Armor: 1, Stats: Stats{Cha: 3, Luc: 1}},
		{Name: "druid", HP: 8, Attack: 2, Armor: 1, Stats: Stats{Sur: 2, Con: 2}},
		{Name: "barbarian", HP: 12, Attack: 4, Armor: 1, Stats: Stats{Str: 3, Con: 1}},
		{Name: "monk", HP: 9, Attack: 3, Armor: 1, Stats: Stats{Dex: 2, Wil: 2}},
		{Name: "guard", HP: 11, Attack: 2, Armor: 4, Stats: Stats{Con: 3, Str: 1}},
	}
}

// GetClassArchetype returns a ClassArchetype by name or default fighter if not found.
func GetClassArchetype(name string) ClassArchetype {
	catalog := ClassCatalog()
	for _, arch := range catalog {
		if arch.Name == name {
			return arch
		}
	}
	return catalog[0]
}
