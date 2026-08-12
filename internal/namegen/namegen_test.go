package namegen

import (
	"math/rand"
	"testing"
)

func TestGenerateCharacterName(t *testing.T) {
	for i := 0; i < 50; i++ {
		name := GenerateCharacterName()
		if name == "" {
			t.Errorf("expected non-empty character name")
		}
	}
}

func TestGenerateCharacterNameForEachRace(t *testing.T) {
	for _, race := range AllRaces {
		for i := 0; i < 10; i++ {
			name := GenerateCharacterNameForRace(race)
			if name == "" {
				t.Errorf("expected non-empty name for race %s", race)
			}
		}
	}
}

func TestGenerateCharacterNameForEachRaceAndGender(t *testing.T) {
	genders := []Gender{GenderMale, GenderFemale}
	for _, race := range AllRaces {
		for _, gender := range genders {
			for i := 0; i < 5; i++ {
				name := GenerateCharacterNameFull(race, gender)
				if name == "" {
					t.Errorf("expected non-empty name for race %s, gender %s", race, gender)
				}
			}
		}
	}
}

func TestRandomRace(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	validRaces := make(map[Race]bool)
	for _, race := range AllRaces {
		validRaces[race] = true
	}

	for i := 0; i < 100; i++ {
		race := RandomRace(r)
		if !validRaces[race] {
			t.Errorf("RandomRace returned invalid race: %s", race)
		}
	}
}

func TestRandomGender(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	seenMale := false
	seenFemale := false

	for i := 0; i < 100; i++ {
		gender := RandomGender(r)
		switch gender {
		case GenderMale:
			seenMale = true
		case GenderFemale:
			seenFemale = true
		default:
			t.Errorf("RandomGender returned unexpected gender: %s", gender)
		}
	}

	if !seenMale {
		t.Errorf("expected to see GenderMale in 100 iterations")
	}
	if !seenFemale {
		t.Errorf("expected to see GenderFemale in 100 iterations")
	}
}

func TestGenerateLocationName(t *testing.T) {
	for i := 0; i < 50; i++ {
		name := GenerateLocationName()
		if name == "" {
			t.Errorf("expected non-empty location name")
		}
	}
}

func TestDeterministicWithSeed(t *testing.T) {
	seed := int64(12345)
	r1 := rand.New(rand.NewSource(seed))
	r2 := rand.New(rand.NewSource(seed))

	name1 := GenerateCharacterNameWithRand(r1)
	name2 := GenerateCharacterNameWithRand(r2)

	if name1 != name2 {
		t.Errorf("expected same name for same seed, got %q and %q", name1, name2)
	}
}

func TestFallbackToHumanForUnknownRace(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	name := GenerateCharacterNameFullWithRand(r, Race("alien"), GenderMale)
	if name == "" {
		t.Errorf("expected non-empty name for fallback race")
	}
}
