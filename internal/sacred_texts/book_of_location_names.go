package sacred_texts

import (
	"math/rand"
)

// Backward compatibility aliases for location names within sacred_texts package.
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
