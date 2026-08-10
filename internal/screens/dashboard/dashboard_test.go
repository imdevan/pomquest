package dashboard

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/pomquest/internal/app"
	"github.com/pomquest/internal/domain"
)

func TestDashboardScreen_Navigation(t *testing.T) {
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
}

func TestDashboardScreen_PartyDisplayAndCoins(t *testing.T) {
	cfg := domain.DefaultConfig()
	party := app.Party{
		Name:             "Heroes Guild",
		MaxPartySize:     4,
		Gold:             150,
		CurrentPartySize: 2,
		Heroes: []app.Hero{
			{Creature: app.Creature{Noun: app.Noun{Name: "Gimli"}}, Class: "barbarian"},
			{Creature: app.Creature{Noun: app.Noun{Name: "Legolas"}}, Class: "ranger"},
		},
	}

	model := NewModelFromParty(cfg, party)
	view := model.View()

	// 28.2: show party current and max size on dashboard under options
	if !strings.Contains(view, "👥 Party: 2/4 Members") {
		t.Errorf("expected view to display party size '👥 Party: 2/4 Members', got:\n%s", view)
	}

	// 28.3: show party coins to right of party hero count
	if !strings.Contains(view, "💰 Coins: 150 Gold") {
		t.Errorf("expected view to display party coins '💰 Coins: 150 Gold', got:\n%s", view)
	}
}

func TestDashboardScreen_CannotQuestWithZeroPartyMembers(t *testing.T) {
	cfg := domain.DefaultConfig()
	model := NewModel(cfg) // CurrentPartySize is 0

	// Attempt to select "1. Quest!" with 0 party members
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	m := updated.(Model)

	// 28.4: do not allow players with 0 current party members to "Quest!"
	if m.Selected == ChoiceQuest {
		t.Errorf("expected Selected to NOT be ChoiceQuest when party size is 0, got ChoiceQuest")
	}
	if cmd != nil {
		t.Errorf("expected nil cmd (no quit) when Quest selection is blocked, got %v", cmd)
	}
	if !strings.Contains(m.StatusNotice, "Cannot start a quest with 0 party members") {
		t.Errorf("expected status notice warning about 0 party members, got %q", m.StatusNotice)
	}

	// Now update party data to have 1 member
	party := app.Party{
		MaxPartySize: 4,
		Heroes: []app.Hero{
			{Creature: app.Creature{Noun: app.Noun{Name: "Aragorn"}}, Class: "paladin"},
		},
	}
	m.UpdatePartyData(party)

	// Select "1. Quest!" now with 1 party member
	updated, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	m = updated.(Model)
	if m.Selected != ChoiceQuest {
		t.Errorf("expected Selected to be ChoiceQuest when party size > 0, got %v", m.Selected)
	}
	if cmd == nil {
		t.Errorf("expected quit cmd to transition screen when Quest is allowed")
	}
}
