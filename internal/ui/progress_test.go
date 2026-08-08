package ui

import (
	"strings"
	"testing"

	"github.com/pomquest/internal/domain"
)

func TestProgressBarComponent(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ThemeFromConfig(cfg)

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

	comp := NewProgressBarComponent(theme, data)
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
}
