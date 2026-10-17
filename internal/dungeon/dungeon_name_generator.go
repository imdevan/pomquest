package dungeon

import (
	"math/rand"
	"time"

	"github.com/pomquest/internal/sacred_texts"
	"github.com/pomquest/internal/utils"
)

// Locays is the slice of valid location types for a dungeon (dungeon, cave, cavern).
var Locays = []string{
	"Dungeon",
	"Cave",
	"Cavern",
}

// GenerateDungeonName creates a random dungeon name using template selection and location name generation.
func GenerateDungeonName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return GenerateDungeonNameWithRand(r)
}

// GenerateDungeonNameWithRand creates a random dungeon name using a provided rand.Rand instance.
func GenerateDungeonNameWithRand(r *rand.Rand) string {
	blank := sacred_texts.GenerateLocationNameWithRand(r)
	return utils.GenerateNameWithRand(r, blank, Locays)
}

// GenerateName is an alias for GenerateDungeonName for backward compatibility.
func GenerateName() string {
	return GenerateDungeonName()
}

// GenerateNameWithRand is an alias for GenerateDungeonNameWithRand for backward compatibility.
func GenerateNameWithRand(r *rand.Rand) string {
	return GenerateDungeonNameWithRand(r)
}
