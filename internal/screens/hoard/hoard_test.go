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

func TestHoardScreenView_MD(t *testing.T) {
	cfg := domain.DefaultConfig()
	model := NewModel(cfg)
	model.SetDimensions(70, 30) // MD breakpoint

	view := model.View()

	if !strings.Contains(view, "1,250g") {
		t.Errorf("expected MD hoard view to contain '1,250g', got:\n%s", view)
	}

	// In MD mode, full text labels like "total gold:" should not be present
	if strings.Contains(view, "total gold:") {
		t.Errorf("expected MD hoard view NOT to contain 'total gold:', got:\n%s", view)
	}
}

func TestHoardScreenView_58Width(t *testing.T) {
	cfg := domain.DefaultConfig()
	model := NewModel(cfg)
	model.SetDimensions(58, 25)

	view := model.View()

	if !strings.Contains(view, "1,250g") {
		t.Errorf("expected 58 width hoard view to contain '1,250g', got:\n%s", view)
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

func TestCropAsciiArt(t *testing.T) {
	sample := "12345678901234567890" // 20 runes
	cropped := cropAsciiArt(sample, 10)
	if len([]rune(cropped)) != 10 {
		t.Errorf("expected cropped rune length to be 10, got %d", len([]rune(cropped)))
	}
}
