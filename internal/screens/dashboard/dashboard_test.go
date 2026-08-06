package dashboard

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/pomquest/internal/domain"
)

func TestDashboardScreen(t *testing.T) {
	cfg := domain.DefaultConfig()
	model := NewModel(cfg)

	// Verify initial rendering contains title and options
	view := model.View()
	if !strings.Contains(view, "POMQUEST - HEARTH DASHBOARD") {
		t.Errorf("expected view to contain header title, got:\n%s", view)
	}
	if !strings.Contains(view, "Quest!") {
		t.Errorf("expected view to contain Quest option, got:\n%s", view)
	}

	// Test navigation down
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m := updated.(Model)
	if m.Cursor != 1 {
		t.Errorf("expected cursor to be 1 after pressing 'j', got %d", m.Cursor)
	}

	// Test direct numeric key selection
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	m = updated.(Model)
	if m.Selected != ChoiceQuest {
		t.Errorf("expected Selected to be ChoiceQuest, got %v", m.Selected)
	}
	if cmd == nil {
		t.Errorf("expected quit command on selection, got nil")
	}
}
