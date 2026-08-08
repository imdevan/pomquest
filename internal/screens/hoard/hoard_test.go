package hoard

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pomquest/internal/domain"
)

func TestHoardScreenView(t *testing.T) {
	cfg := domain.DefaultConfig()
	model := NewModel(cfg)
	model.SetDimensions(100, 40)

	view := model.View()

	// Check for "Hoard" textured card top-left corner
	if !strings.Contains(view, "Hoard") {
		t.Errorf("expected hoard view to contain textured card header 'Hoard', got:\n%s", view)
	}

	// Check 3-column stats
	expectedStats := []string{
		"total gold",
		"heroes hired",
		"quests started",
		"total loot",
		"heroes lost",
		"quests completed",
		"total treasure",
		"places traveled",
		"time focused",
	}

	for _, stat := range expectedStats {
		if !strings.Contains(view, stat) {
			t.Errorf("expected hoard view to contain stat %q, got:\n%s", stat, view)
		}
	}
}

func TestHoardScreen_SmallScreenUnsupported(t *testing.T) {
	cfg := domain.DefaultConfig()
	model := NewModel(cfg)
	model.SetDimensions(45, 20) // SM breakpoint

	view := model.View()

	if !strings.Contains(view, "terminal width is not currently") {
		t.Errorf("expected small screen hoard view to show unsupported message, got:\n%s", view)
	}
}

func TestHoardScreen_Quit(t *testing.T) {
	cfg := domain.DefaultConfig()
	model := NewModel(cfg)

	// Simulate pressing 'q'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	updated, _ := model.Update(msg)
	updatedModel := updated.(Model)

	if !updatedModel.Quitted {
		t.Errorf("expected model.Quitted to be true after pressing 'q'")
	}
	if updatedModel.View() != "" {
		t.Errorf("expected View() to return empty string when Quitted is true")
	}
}
