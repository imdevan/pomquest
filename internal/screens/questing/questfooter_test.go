package questing

import (
	"strings"
	"testing"

	"github.com/pomquest/internal/domain"
	"github.com/pomquest/internal/ui"
)

func TestQuestFooterComponent(t *testing.T) {
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
