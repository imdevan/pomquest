package dungeon

import (
	"math/rand"
	"time"

	"github.com/pomquest/internal/namegen"
)

// BookOfLocationNames contains a curated collection of location base names.
// Data is sourced from the shared namegen package.
var BookOfLocationNames = namegen.BookOfLocationNames

// LocationPrefixes and LocationSuffixes allow generating donjon-style location names.
// Data is sourced from the shared namegen package.
var LocationPrefixes = namegen.LocationPrefixes
var LocationSuffixes = namegen.LocationSuffixes

// GenerateLocationName returns a random location name (either from BookOfLocationNames or composed like donjon).
func GenerateLocationName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return GenerateLocationNameWithRand(r)
}

// GenerateLocationNameWithRand returns a random location name using the provided rand.Rand instance.
func GenerateLocationNameWithRand(r *rand.Rand) string {
	return namegen.GenerateLocationNameWithRand(r)
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
