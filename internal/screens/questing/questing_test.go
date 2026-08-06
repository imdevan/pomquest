package questing

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/pomquest/internal/domain"
)

func TestQuestingScreen(t *testing.T) {
	cfg := domain.DefaultConfig()
	model := NewModel(cfg)

	// Test View rendering
	view := model.View()
	if !strings.Contains(view, "Forest of Izelliador") {
		t.Errorf("expected view to contain Forest of Izelliador, got:\n%s", view)
	}
	if !strings.Contains(view, "jimmothy") {
		t.Errorf("expected view to contain hero jimmothy, got:\n%s", view)
	}

	// Test pressing space to add attack action log entry
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	m := updated.(Model)
	lastLog := m.LogMessages[len(m.LogMessages)-1]
	if !strings.Contains(lastLog, "pam attacks") {
		t.Errorf("expected new log entry to contain pam attacks, got '%s'", lastLog)
	}

	// Test pause key
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	m = updated.(Model)
	if !m.Paused {
		t.Errorf("expected model to be paused after pressing 'p'")
	}
}
