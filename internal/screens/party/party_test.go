package party

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/pomquest/internal/adapters/mouse"
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

func TestPartyScreen_HeroCardSectionCentered(t *testing.T) {
	cfg := domain.DefaultConfig()

	// Test 1 hero card
	p1 := app.Party{
		MaxPartySize: 1,
		Heroes: []app.Hero{
			{Creature: app.Creature{Noun: app.Noun{Name: "SoloHero"}, HP: 10, MaxHP: 10}, Class: "fighter"},
		},
	}
	m1 := NewModelFromParty(cfg, p1)
	m1.SetDimensions(80, 24)
	v1 := m1.View()

	if !strings.Contains(v1, "SoloHero") {
		t.Errorf("expected view v1 to contain 'SoloHero'")
	}
	if !strings.Contains(v1, "PARTY MEMBERS") {
		t.Errorf("expected view v1 to contain header 'PARTY MEMBERS'")
	}

	// Verify centering: find lines containing card box border e.g. "SoloHero"
	lines1 := strings.Split(v1, "\n")
	var foundCardLine1 string
	for _, l := range lines1 {
		if strings.Contains(l, "SoloHero") {
			foundCardLine1 = l
			break
		}
	}
	if foundCardLine1 == "" {
		t.Fatalf("card line with SoloHero not found")
	}
	// With 80 cols width, "SoloHero" card should have significant leading spaces centering it
	leadingSpaces1 := len(foundCardLine1) - len(strings.TrimLeft(foundCardLine1, " "))
	if leadingSpaces1 < 10 {
		t.Errorf("expected 1 hero card to be centered with leading padding, got line:\n%q", foundCardLine1)
	}

	// Test 4 hero cards
	m4 := NewModel(cfg) // default has 4 heroes
	m4.SetDimensions(120, 24)
	v4 := m4.View()

	if !strings.Contains(v4, "jimmothy") || !strings.Contains(v4, "applebee") {
		t.Errorf("expected 4-hero view to contain hero names")
	}

	lines4 := strings.Split(v4, "\n")
	var foundCardLine4 string
	for _, l := range lines4 {
		if strings.Contains(l, "jimmothy") {
			foundCardLine4 = l
			break
		}
	}
	if foundCardLine4 == "" {
		t.Fatalf("card line with jimmothy not found")
	}
	leadingSpaces4 := len(foundCardLine4) - len(strings.TrimLeft(foundCardLine4, " "))
	if leadingSpaces4 < 5 {
		t.Errorf("expected 4 hero cards on 120-width screen to be centered with leading padding, got line:\n%q", foundCardLine4)
	}
}

func TestPartyScreen_MouseHoverAndClick(t *testing.T) {
	mouse.InitManager()
	cfg := domain.DefaultConfig()
	model := NewModel(cfg)

	// Render view so bubblezone registers zones
	_ = mouse.Scan(model.View())
	time.Sleep(20 * time.Millisecond)

	zInfo := mouse.Get("party_hero_0")
	if zInfo == nil {
		t.Fatalf("expected zone info for party_hero_0")
	}

	// 1. Mouse hover over hero 0 -> ActiveIndex becomes 0
	hoverMsg := tea.MouseMsg{
		X:      zInfo.StartX + 5,
		Y:      zInfo.StartY + 2,
		Type:   tea.MouseMotion,
		Action: tea.MouseActionMotion,
	}
	updated, _ := model.Update(hoverMsg)
	m := updated.(Model)
	if m.ActiveIndex != 0 {
		t.Errorf("expected ActiveIndex 0 on mouse hover, got %d", m.ActiveIndex)
	}

	// 2. Mouse click down (press) -> PressedHeroIndex becomes 0 (active_focused)
	pressMsg := tea.MouseMsg{
		X:      zInfo.StartX + 5,
		Y:      zInfo.StartY + 2,
		Button: tea.MouseButtonLeft,
	}
	updated, _ = m.Update(pressMsg)
	m = updated.(Model)
	if m.PressedHeroIndex != 0 {
		t.Errorf("expected PressedHeroIndex 0 on mouse press, got %d", m.PressedHeroIndex)
	}

	// 3. Mouse release -> PressedHeroIndex resets to -1
	releaseMsg := tea.MouseMsg{
		X:      zInfo.StartX + 5,
		Y:      zInfo.StartY + 2,
		Type:   tea.MouseRelease,
		Action: tea.MouseActionRelease,
	}
	updated, _ = m.Update(releaseMsg)
	m = updated.(Model)
	if m.PressedHeroIndex != -1 {
		t.Errorf("expected PressedHeroIndex -1 on mouse release, got %d", m.PressedHeroIndex)
	}
}

func TestPartyScreen_MouseClickRecruitHeroNavigatesToReroll(t *testing.T) {
	mouse.InitManager()
	cfg := domain.DefaultConfig()
	model := NewModel(cfg)
	model.Heroes[3].IsEmptySlot = true

	_ = mouse.Scan(model.View())
	time.Sleep(20 * time.Millisecond)

	zInfo := mouse.Get("party_hero_3")
	if zInfo == nil {
		t.Fatalf("expected zone info for party_hero_3")
	}

	pressMsg := tea.MouseMsg{
		X:      zInfo.StartX + 5,
		Y:      zInfo.StartY + 2,
		Button: tea.MouseButtonLeft,
	}
	updated, cmd := model.Update(pressMsg)
	m := updated.(Model)

	if m.PressedHeroIndex != 3 {
		t.Errorf("expected PressedHeroIndex 3 on mouse press, got %d", m.PressedHeroIndex)
	}
	if cmd == nil {
		t.Fatalf("expected next cycle cmd on mouse click")
	}

	nextCycleMsg := cmd()
	updated, cmd2 := m.Update(nextCycleMsg)
	m = updated.(Model)

	if m.PressedHeroIndex != -1 {
		t.Errorf("expected PressedHeroIndex -1 on next cycle, got %d", m.PressedHeroIndex)
	}
	if !m.NavToReroll {
		t.Errorf("expected NavToReroll to be true on next cycle")
	}
	if m.RerollTargetIndex != 3 {
		t.Errorf("expected RerollTargetIndex to be 3, got %d", m.RerollTargetIndex)
	}
	if cmd2 == nil {
		t.Errorf("expected quit tea.Cmd when navigating to reroll screen")
	}
}

func TestPartyScreen_SmallScreenUnsupported(t *testing.T) {
	cfg := domain.DefaultConfig()
	m := NewModel(cfg)
	m.SetDimensions(40, 24) // SM breakpoint or less

	view := m.View()
	if !strings.Contains(view, "terminal width is not") {
		t.Fatalf("expected small screen unsupported message, got:\n%s", view)
	}
}

