package reroll

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/pomquest/internal/adapters/mouse"
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

	// Verify all cards have placeholder names
	for i := 0; i < 4; i++ {
		if model.Cards[i].Name == "" {
			t.Errorf("expected Card %d to have a placeholder name", i)
		}
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
	expectedName := m.Cards[1].Name
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	m = updated.(Model)

	if !m.EnteringName {
		t.Errorf("expected EnteringName to be true after selecting class")
	}

	// Verify the placeholder name matches the selected hero card's name
	if m.NameInput != expectedName {
		t.Errorf("expected NameInput to equal card name %q, got %q", expectedName, m.NameInput)
	}

	// Clear the pre-filled name with backspace
	for len(m.NameInput) > 0 {
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		m = updated.(Model)
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

func TestRerollScreen_MouseHoverAndClick(t *testing.T) {
	mouse.InitManager()
	cfg := domain.DefaultConfig()
	model := NewModel(cfg, 1)

	// Render view so bubblezone registers zones
	_ = mouse.Scan(model.View())
	time.Sleep(20 * time.Millisecond)

	zInfo := mouse.Get("reroll_card_1")
	if zInfo == nil {
		t.Fatalf("expected zone info for reroll_card_1")
	}

	// 1. Mouse hover over card 1 -> ActiveIndex becomes 1
	hoverMsg := tea.MouseMsg{
		X:      zInfo.StartX + 5,
		Y:      zInfo.StartY + 2,
		Type:   tea.MouseMotion,
		Action: tea.MouseActionMotion,
	}
	updated, _ := model.Update(hoverMsg)
	m := updated.(Model)
	if m.ActiveIndex != 1 {
		t.Errorf("expected ActiveIndex 1 on mouse hover, got %d", m.ActiveIndex)
	}

	// 2. Mouse click (EventClick) -> PressedCardIndex becomes 1 (active_focused) & returns next cycle cmd
	pressMsg := tea.MouseMsg{
		X:      zInfo.StartX + 5,
		Y:      zInfo.StartY + 2,
		Button: tea.MouseButtonLeft,
	}
	updated, cmd := m.Update(pressMsg)
	m = updated.(Model)
	if m.PressedCardIndex != 1 {
		t.Errorf("expected PressedCardIndex 1 on mouse press, got %d", m.PressedCardIndex)
	}
	if cmd == nil {
		t.Fatalf("expected next cycle cmd on mouse click")
	}

	// 3. Next cycle message -> EnteringName becomes true, PressedCardIndex resets to -1
	nextCycleMsg := cmd()
	updated, _ = m.Update(nextCycleMsg)
	m = updated.(Model)
	if m.PressedCardIndex != -1 {
		t.Errorf("expected PressedCardIndex -1 on next cycle, got %d", m.PressedCardIndex)
	}
	if !m.EnteringName {
		t.Errorf("expected EnteringName to be true on next cycle after mouse click")
	}
}

func TestRerollScreen_SmallScreenUnsupported(t *testing.T) {
	cfg := domain.DefaultConfig()
	m := NewModel(cfg, 1)
	m.SetDimensions(40, 24) // SM breakpoint or less

	view := m.View()
	if !strings.Contains(view, "terminal width is not") {
		t.Fatalf("expected small screen unsupported message, got:\n%s", view)
	}
}

