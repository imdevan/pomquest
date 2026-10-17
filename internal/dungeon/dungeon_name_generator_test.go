package dungeon

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/pomquest/internal/sacred_texts"
)

func TestGenerateLocationName(t *testing.T) {
	for i := 0; i < 50; i++ {
		name := sacred_texts.GenerateLocationName()
		if name == "" {
			t.Errorf("expected non-empty location name")
		}
	}
}

func TestGenerateDungeonNameTemplatesAndLocays(t *testing.T) {
	locayMap := map[string]bool{
		"Dungeon": true,
		"Cave":    true,
		"Cavern":  true,
	}

	for seed := int64(1); seed <= 100; seed++ {
		r := rand.New(rand.NewSource(seed))
		name := GenerateDungeonNameWithRand(r)
		if name == "" {
			t.Fatalf("expected non-empty dungeon name for seed %d", seed)
		}

		// Check if it matches one of the 3 templates:
		// 1. "<locay> of <blank>"
		// 2. "<blank>'s <locay>"
		// 3. "<blank> <locay>"
		matchedTemplate := false

		// Template 1: "<locay> of <blank>"
		for loc := range locayMap {
			prefix := loc + " of "
			if strings.HasPrefix(name, prefix) {
				blank := strings.TrimPrefix(name, prefix)
				if len(blank) > 0 {
					matchedTemplate = true
					break
				}
			}
		}

		// Template 2: "<blank>'s <locay>"
		if !matchedTemplate {
			for loc := range locayMap {
				suffix := "'s " + loc
				if strings.HasSuffix(name, suffix) {
					blank := strings.TrimSuffix(name, suffix)
					if len(blank) > 0 {
						matchedTemplate = true
						break
					}
				}
			}
		}

		// Template 3: "<blank> <locay>"
		if !matchedTemplate {
			for loc := range locayMap {
				suffix := " " + loc
				if strings.HasSuffix(name, suffix) {
					blank := strings.TrimSuffix(name, suffix)
					if len(blank) > 0 {
						matchedTemplate = true
						break
					}
				}
			}
		}

		if !matchedTemplate {
			t.Errorf("name %q did not match any of the expected templates or locays", name)
		}
	}
}
