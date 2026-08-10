package reroll

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/pomquest/internal/domain"
)

func TestRerollScreen(t *testing.T) {
	cfg := domain.DefaultConfig()
	model := NewModel(cfg, 3)

	// Verify header renders
	view := model.View()
	if !strings.Contains(view, "REROLL HERO CLASS") {
		t.Errorf("expected view to contain 'REROLL HERO CLASS', got:\n%s", view)
	}

	// Verify initial cards: Card 0 is random
	if model.Cards[0].Class != "random" {
		t.Errorf("expected Card 0 class to be 'random', got %s", model.Cards[0].Class)
	}
	if !model.Cards[0].IsUnknownStats {
		t.Errorf("expected Card 0 IsUnknownStats to be true")
	}
	if model.Cards[0].Level != 1 {
		t.Errorf("expected Card 0 Level to be 1, got %d", model.Cards[0].Level)
	}

	// Initial rerolls left is 2
	if model.RerollsLeft != 2 {
		t.Errorf("expected RerollsLeft = 2, got %d", model.RerollsLeft)
	}

	// Test rerolling options with 'r'
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	m := updated.(Model)
	if m.RerollsLeft != 1 {
		t.Errorf("expected RerollsLeft = 1 after pressing 'r', got %d", m.RerollsLeft)
	}

	// Test selecting Card 1 (specific class)
	m.ActiveIndex = 1
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	m = updated.(Model)

	if !m.EnteringName {
		t.Errorf("expected EnteringName to be true after selecting class")
	}

	// Type name 'arthur' and submit with Enter
	for _, ch := range "arthur" {
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		m = updated.(Model)
	}

	if m.NameInput != "arthur" {
		t.Errorf("expected NameInput to be 'arthur', got %s", m.NameInput)
	}

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)

	if !m.Completed {
		t.Errorf("expected Completed to be true after submitting name")
	}
	if m.NewHero.Name != "arthur" {
		t.Errorf("expected NewHero.Name to be 'arthur', got %s", m.NewHero.Name)
	}
	if m.NewHero.IsDead {
		t.Errorf("expected NewHero.IsDead to be false")
	}
	if cmd == nil {
		t.Errorf("expected quit tea.Cmd on completion, got nil")
	}
}
