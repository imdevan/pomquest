package sacred_texts

import (
	"math/rand"
	"time"
)

// Race represents a fantasy race for name generation.
type Race string

const (
	RaceHuman  Race = "human"
	RaceElf    Race = "elf"
	RaceDwarf  Race = "dwarf"
	RaceHobbit Race = "hobbit"
	RaceGoblin Race = "goblin"
	RaceOrc    Race = "orc"
)

// AllRaces contains all available fantasy races.
var AllRaces = []Race{
	RaceHuman,
	RaceElf,
	RaceDwarf,
	RaceHobbit,
	RaceGoblin,
	RaceOrc,
}

// Gender represents a character gender for name generation (not displayed to user).
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

// nameData holds prefix and suffix syllables for a specific race+gender combination.
type nameData struct {
	Prefixes []string
	Suffixes []string
}



// GenerateCharacterName generates a character name with a random race and gender.
func GenerateCharacterName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return GenerateCharacterNameWithRand(r)
}

// GenerateCharacterNameWithRand generates a character name with a random race and gender
// using the provided rand.Rand instance.
func GenerateCharacterNameWithRand(r *rand.Rand) string {
	race := RandomRace(r)
	gender := RandomGender(r)
	return GenerateCharacterNameFullWithRand(r, race, gender)
}

// GenerateCharacterNameForRace generates a character name for a specific race with random gender.
func GenerateCharacterNameForRace(race Race) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return GenerateCharacterNameForRaceWithRand(r, race)
}

// GenerateCharacterNameForRaceWithRand generates a character name for a specific race
// with random gender using the provided rand.Rand instance.
func GenerateCharacterNameForRaceWithRand(r *rand.Rand, race Race) string {
	gender := RandomGender(r)
	return GenerateCharacterNameFullWithRand(r, race, gender)
}

// GenerateCharacterNameFull generates a character name for a specific race and gender.
func GenerateCharacterNameFull(race Race, gender Gender) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return GenerateCharacterNameFullWithRand(r, race, gender)
}

// GenerateCharacterNameFullWithRand generates a character name for a specific race and gender
// using the provided rand.Rand instance. This is the core name generation function.
func GenerateCharacterNameFullWithRand(r *rand.Rand, race Race, gender Gender) string {
	genderMap, ok := characterNames[race]
	if !ok {
		// Fallback to human if unknown race
		genderMap = characterNames[RaceHuman]
	}

	data, ok := genderMap[gender]
	if !ok {
		// Fallback to male if unknown gender
		data = genderMap[GenderMale]
	}

	prefix := data.Prefixes[r.Intn(len(data.Prefixes))]
	suffix := data.Suffixes[r.Intn(len(data.Suffixes))]
	return prefix + suffix
}

// RandomRace returns a random race from AllRaces.
func RandomRace(r *rand.Rand) Race {
	return AllRaces[r.Intn(len(AllRaces))]
}

// RandomGender returns a random gender.
func RandomGender(r *rand.Rand) Gender {
	if r.Intn(2) == 0 {
		return GenderMale
	}
	return GenderFemale
}

// GenerateLocationName returns a random location name.
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
