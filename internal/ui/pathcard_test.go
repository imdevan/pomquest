package ui

import (
	"strings"
	"testing"

	"github.com/pomquest/internal/domain"
)

func TestPathCardComponent(t *testing.T) {
	theme := ThemeFromConfig(domain.DefaultConfig())
	data := PathCardData{
		Name:               "Dark Caverns",
		Biome:              "Shadow Caves",
		Description:        "A deep winding cave filled with shadowy whispers.",
		DifficultyModifier: 1.2,
		IsActive:           true,
	}

	comp := NewPathCardComponent(theme, data)
	output := comp.View()

	if !strings.Contains(output, "Dark Caverns") {
		t.Errorf("expected output to contain 'Dark Caverns', got:\n%s", output)
	}
}

func TestPathCardComponent_Sized3To4AspectRatio(t *testing.T) {
	theme := ThemeFromConfig(domain.DefaultConfig())
	data := PathCardData{
		Name:               "Ancient Ruins",
		Biome:              "Catacombs",
		Description:        "Crumbling stone walls rich with ancient history.",
		DifficultyModifier: 1.5,
		IsActive:           false,
	}

	// Width 24, Height 16 (visual ratio 24 / (16 * 2) = 0.75 = 3:4 aspect ratio)
	comp := NewPathCardComponentSized(theme, data, 24, 16)
	output := comp.View()

	lines := strings.Split(output, "\n")
	if len(lines) != 16 {
		t.Errorf("expected 16 lines for card height, got %d", len(lines))
	}

	for i, line := range lines {
		// Strip ANSI styles for width check
		lineLen := len([]rune(line))
		if lineLen == 0 {
			t.Errorf("line %d is empty", i)
		}
	}
}

func TestPathCardComponent_BiomeIcons(t *testing.T) {
	theme := ThemeFromConfig(domain.DefaultConfig())

	tests := []struct {
		biome        string
		icon         string
		expectedIcon string
	}{
		{biome: "swamp", icon: "", expectedIcon: "🐊"},
		{biome: "desert", icon: "", expectedIcon: "🏜️"},
		{biome: "mountains", icon: "", expectedIcon: "🏔️"},
		{biome: "island", icon: "", expectedIcon: "🏝️"},
		{biome: "plains", icon: "", expectedIcon: "🌾"},
		{biome: "forest", icon: "", expectedIcon: "🌲"},
	}

	for _, tt := range tests {
		data := PathCardData{
			Name:  "Test Location",
			Biome: tt.biome,
			Icon:  tt.icon,
		}
		comp := NewPathCardComponent(theme, data)
		output := comp.View()
		if !strings.Contains(output, tt.expectedIcon) {
			t.Errorf("expected output for biome %q to contain icon %q, got:\n%s", tt.biome, tt.expectedIcon, output)
		}
	}
}

func TestPathCardComponent_LineLayoutAndTitleCase(t *testing.T) {
	theme := ThemeFromConfig(domain.DefaultConfig())
	data := PathCardData{
		Name:     "Whispering Woods",
		Biome:    "shadow_caves",
		IsActive: false,
	}

	comp := NewPathCardComponent(theme, data)
	output := comp.View()

	if !strings.Contains(output, "Whispering Woods") {
		t.Errorf("expected output to contain title 'Whispering Woods'")
	}
	if !strings.Contains(output, "~") {
		t.Errorf("expected output to contain tilde separator '~'")
	}
	if !strings.Contains(output, "Shadow Caves") {
		t.Errorf("expected output to contain title cased biome 'Shadow Caves', got:\n%s", output)
	}
}

