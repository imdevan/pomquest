package party

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/pomquest/internal/domain"
)

func TestPartyScreen(t *testing.T) {
	cfg := domain.DefaultConfig()
	model := NewModel(cfg)

	// Verify title renders
	view := model.View()
	if !strings.Contains(view, "PARTY MEMBERS") {
		t.Errorf("expected view to contain 'PARTY MEMBERS', got:\n%s", view)
	}

	// Initial active index is 2 (pam)
	if model.ActiveIndex != 2 {
		t.Errorf("expected default active index 2, got %d", model.ActiveIndex)
	}

	// Test navigation left
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m := updated.(Model)
	if m.ActiveIndex != 1 {
		t.Errorf("expected active index 1 after pressing 'h', got %d", m.ActiveIndex)
	}

	// Test quick key jump
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	m = updated.(Model)
	if m.ActiveIndex != 0 {
		t.Errorf("expected active index 0 after pressing '1', got %d", m.ActiveIndex)
	}

	// Test quit key
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m = updated.(Model)
	if !m.Quitted {
		t.Errorf("expected Quitted to be true after pressing 'q'")
	}
	if cmd == nil {
		t.Errorf("expected quit tea.Cmd, got nil")
	}
}
