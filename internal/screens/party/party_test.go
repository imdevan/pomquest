package party

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/pomquest/internal/app"
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

func TestPartyScreen_SQLData(t *testing.T) {
	cfg := domain.DefaultConfig()

	sqlParty := app.Party{
		Name:             "Dungeon Crawlers",
		MaxPartySize:     3,
		CurrentPartySize: 1,
		Heroes: []app.Hero{
			{
				Creature: app.Creature{
					Noun:   app.Noun{Name: "Gimli"},
					HP:     12,
					MaxHP:  12,
					Attack: 3,
					Armor:  2,
				},
				Class:  "barbarian",
				Level:  1,
				Status: "Healthy",
			},
		},
	}

	model := NewModelFromParty(cfg, sqlParty)

	// 22.4.1: show only number of hero cards based on party max size (3)
	if len(model.Heroes) != 3 {
		t.Fatalf("expected 3 hero slots based on MaxPartySize=3, got %d", len(model.Heroes))
	}
	if model.Heroes[0].IsEmptySlot {
		t.Errorf("expected slot 0 to be filled hero")
	}
	if !model.Heroes[1].IsEmptySlot || !model.Heroes[2].IsEmptySlot {
		t.Errorf("expected slots 1 and 2 to be empty hero slots")
	}

	// 22.4.2 & 22.4.3: render empty hero slot card with "Recruit Hero" button inside
	view := model.View()
	if !strings.Contains(view, "Recruit") || !strings.Contains(view, "Hero") {
		t.Errorf("expected party screen view to contain 'Recruit Hero' inside empty card, got:\n%s", view)
	}

	// 22.4.4: selecting an empty slot and pressing enter or r directs user to reroll screen
	model.ActiveIndex = 1
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	m := updated.(Model)
	if !m.NavToReroll {
		t.Errorf("expected NavToReroll to be true when pressing 'r' on empty slot")
	}
	if m.RerollTargetIndex != 1 {
		t.Errorf("expected RerollTargetIndex to be 1, got %d", m.RerollTargetIndex)
	}
	if cmd == nil {
		t.Errorf("expected quit tea.Cmd when navigating to reroll screen")
	}
}
