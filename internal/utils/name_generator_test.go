package utils

import (
	"math/rand"
	"testing"
)

func TestGenerateName(t *testing.T) {
	types := []string{"Dungeon", "Cave", "Cavern"}
	for i := 0; i < 50; i++ {
		name := GenerateName("TestBase", types)
		if name == "" {
			t.Errorf("expected non-empty generated name")
		}
	}
}

func TestGenerateNameWithRandEmptyTypes(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	name := GenerateNameWithRand(r, "Standalone", nil)
	if name != "Standalone" {
		t.Errorf("expected 'Standalone', got %q", name)
	}
}

func TestGenerateNameTemplates(t *testing.T) {
	types := []string{"Dungeon", "Cave", "Cavern"}
	baseName := "Izelliador"

	for seed := int64(1); seed <= 100; seed++ {
		r := rand.New(rand.NewSource(seed))
		name := GenerateNameWithRand(r, baseName, types)
		if name == "" {
			t.Fatalf("expected non-empty name for seed %d", seed)
		}

		matchedTemplate := false

		// Template 1: "<type> of <baseName>"
		for _, loc := range types {
			if name == loc+" of "+baseName {
				matchedTemplate = true
				break
			}
		}

		// Template 2: "<baseName>'s <type>"
		if !matchedTemplate {
			for _, loc := range types {
				if name == baseName+"'s "+loc {
					matchedTemplate = true
					break
				}
			}
		}

		// Template 3: "<baseName> <type>"
		if !matchedTemplate {
			for _, loc := range types {
				if name == baseName+" "+loc {
					matchedTemplate = true
					break
				}
			}
		}

		if !matchedTemplate {
			t.Errorf("name %q did not match expected templates", name)
		}
	}
}
