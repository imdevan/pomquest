package dungeon

import (
	"math/rand"
	"time"
)

// BookOfLocationNames contains a curated collection of location base names.
var BookOfLocationNames = []string{
	"Izelliador",
	"Malakor",
	"Zalthor",
	"Gloomhaven",
	"Vaelgor",
	"Kardoth",
	"Thulgor",
	"Baelroth",
	"Dreadfang",
	"Shadowveil",
	"Obsidian",
	"Bloodspire",
	"Voragath",
	"Kaelmon",
	"Drakfell",
	"Arakmire",
	"Xalnoth",
	"Grimhold",
	"Frostcrag",
	"Ironreach",
	"Stormforge",
	"Voidvault",
	"Nightmar",
	"Skymor",
	"Emberfall",
	"Blightwood",
	"Ashridge",
	"Morvath",
}

// LocationPrefixes and LocationSuffixes allow generating donjon-style location names.
var LocationPrefixes = []string{
	"Bael", "Zal", "Kard", "Thul", "Vael", "Gloom", "Dread", "Gor", "Mael", "Mor",
	"Rath", "Drak", "Khar", "Grim", "Vor", "Shadow", "Ash", "Storm", "Iron", "Frost",
	"Void", "Night", "Mal", "Izel", "Arak", "Xal", "Kael", "Bal", "Krag", "Vex",
}

var LocationSuffixes = []string{
	"roth", "thor", "doth", "gor", "haven", "fang", "veil", "spire", "mire", "hold",
	"keep", "fell", "crag", "reach", "depth", "forge", "gate", "ridge", "vault", "mar",
	"dor", "mon", "gath", "noth", "liador", "heim", "crest", "glen",
}

// GenerateLocationName returns a random location name (either from BookOfLocationNames or composed like donjon).
func GenerateLocationName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return GenerateLocationNameWithRand(r)
}

// GenerateLocationNameWithRand returns a random location name using the provided rand.Rand instance.
func GenerateLocationNameWithRand(r *rand.Rand) string {
	if r.Intn(2) == 0 {
		return BookOfLocationNames[r.Intn(len(BookOfLocationNames))]
	}
	prefix := LocationPrefixes[r.Intn(len(LocationPrefixes))]
	suffix := LocationSuffixes[r.Intn(len(LocationSuffixes))]
	return prefix + suffix
}

// Backward compatibility aliases
var (
	BookOfFantasyNames = BookOfLocationNames
	FantasyPrefixes    = LocationPrefixes
	FantasySuffixes    = LocationSuffixes
)

func GenerateFantasyName() string {
	return GenerateLocationName()
}

func GenerateFantasyNameWithRand(r *rand.Rand) string {
	return GenerateLocationNameWithRand(r)
}
