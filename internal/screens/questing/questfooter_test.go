package questing

import (
	"strings"
	"testing"

	"github.com/pomquest/internal/domain"
	"github.com/pomquest/internal/ui"
)

func TestQuestFooterComponent_LG(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ui.ThemeFromConfig(cfg)

	data := QuestProgressData{
		Destination: "Forest of Izelliador",
		CurrentSecs:  1260,
		TotalSecs:    1500,
		Weather:      "☀️",
		Terrain:      "🌲",
		Hazards:      "mud, rain",
		DaysPassed:   12,
		Width:        96,
	}

	comp := NewQuestFooterComponent(theme, data, 96)
	view := comp.View()

	if !strings.Contains(view, "Forest of Izelliador") {
		t.Errorf("expected view to contain destination, got:\n%s", view)
	}
	if !strings.Contains(view, "21:00/25:00") {
		t.Errorf("expected view to contain time 21:00/25:00, got:\n%s", view)
	}
	if !strings.Contains(view, "days passed: 12") {
		t.Errorf("expected view to contain days passed, got:\n%s", view)
	}
	if !strings.Contains(view, "Weather: ☀️ | Terrain: 🌲 | Hazards: mud, rain") {
		t.Errorf("expected view to contain full environmental labels on LG+, got:\n%s", view)
	}

	lines := strings.Split(view, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines in quest footer view, got %d", len(lines))
	}

	if !strings.HasSuffix(lines[0], "21:00/25:00") {
		t.Errorf("expected line 1 to end with right-aligned time, got:\n%s", lines[0])
	}
}

func TestQuestFooterComponent_MD(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ui.ThemeFromConfig(cfg)

	data := QuestProgressData{
		Destination: "Forest of Izelliador",
		CurrentSecs:  1260,
		TotalSecs:    1500,
		Weather:      "☀️",
		Terrain:      "🌲",
		Hazards:      "mud, rain",
		DaysPassed:   12,
		Width:        72,
	}

	comp := NewQuestFooterComponent(theme, data, 72)
	view := comp.View()

	// On MD: expect "☀️ | 🌲 | mud, rain" without Weather:, Terrain:, Hazards: labels
	if !strings.Contains(view, "☀️ | 🌲 | mud, rain") {
		t.Errorf("expected view to contain '☀️ | 🌲 | mud, rain' on MD breakpoint, got:\n%s", view)
	}
	if strings.Contains(view, "Weather:") {
		t.Errorf("expected view NOT to contain 'Weather:' label on MD breakpoint, got:\n%s", view)
	}
}
