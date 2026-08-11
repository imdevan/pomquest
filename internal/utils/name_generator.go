package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateName creates a random formatted name from a base name and a list of type descriptors using name templates.
func GenerateName(baseName string, types []string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return GenerateNameWithRand(r, baseName, types)
}

// GenerateNameWithRand creates a random formatted name using a provided rand.Rand instance, base name, and type descriptors.
// Templates:
// 1. <type> of <baseName> ("Dungeon of Izelliador")
// 2. <baseName>'s <type> ("Izelliador's Cave")
// 3. <baseName> <type> ("Izelliador Cavern")
func GenerateNameWithRand(r *rand.Rand, baseName string, types []string) string {
	if len(types) == 0 {
		return baseName
	}
	t := types[r.Intn(len(types))]

	templateChoice := r.Intn(3)
	switch templateChoice {
	case 0:
		// type of baseName
		return fmt.Sprintf("%s of %s", t, baseName)
	case 1:
		// baseName's type
		return fmt.Sprintf("%s's %s", baseName, t)
	default:
		// baseName type
		return fmt.Sprintf("%s %s", baseName, t)
	}
}
