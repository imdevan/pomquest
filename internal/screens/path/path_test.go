package path

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pomquest/internal/adapters/mouse"
	"github.com/pomquest/internal/app"
	"github.com/pomquest/internal/domain"
)

func TestPathScreen_GenerationAndNavigation(t *testing.T) {
	cfg := domain.DefaultConfig()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "path_screen_test.db")

	m := NewModelWithDBPath(cfg, dbPath)

	// Check 3 generated options
	if len(m.Options) != 3 {
		t.Fatalf("expected 3 options, got %d", len(m.Options))
	}
	for i, opt := range m.Options {
		if opt.Name == "" {
			t.Errorf("option %d name is empty", i)
		}
	}

	// Initial active index is 0
	if m.ActiveIndex != 0 {
		t.Errorf("expected initial ActiveIndex 0, got %d", m.ActiveIndex)
	}

	// Move right
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	pm := newM.(Model)
	if pm.ActiveIndex != 1 {
		t.Errorf("expected ActiveIndex 1 after key right, got %d", pm.ActiveIndex)
	}

	// Press key 3
	newM, _ = pm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("3")})
	pm = newM.(Model)
	if pm.ActiveIndex != 2 {
		t.Errorf("expected ActiveIndex 2 after pressing '3', got %d", pm.ActiveIndex)
	}
}

func TestPathScreen_PartyCheckAndSelection(t *testing.T) {
	cfg := domain.DefaultConfig()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "path_selection_test.db")

	// 1. Initialized DB has party size 0
	m := NewModelWithDBPath(cfg, dbPath)
	if m.PartySize != 0 {
		t.Fatalf("expected party size 0, got %d", m.PartySize)
	}

	// Attempting selection with 0 party size should fail and set StatusNotice
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	pm := newM.(Model)
	if pm.Completed {
		t.Errorf("expected Completed to be false when party size is 0")
	}
	if pm.StatusNotice == "" {
		t.Errorf("expected StatusNotice to be non-empty when party size is 0")
	}

	// 2. Add hero to party so party size > 0
	db, err := app.InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	hero := app.Hero{
		Creature: app.Creature{
			Noun:  app.Noun{Name: "Hero 1"},
			HP:    10,
			MaxHP: 10,
		},
		Class: "warrior",
	}
	if err := app.ReplaceOrAddHeroInParty(db, 0, hero); err != nil {
		t.Fatalf("ReplaceOrAddHeroInParty failed: %v", err)
	}

	// Reload path model with updated party
	m = NewModelWithDBPath(cfg, dbPath)
	if m.PartySize != 1 {
		t.Fatalf("expected party size 1, got %d", m.PartySize)
	}

	// Press enter to select path option 1 (ActiveIndex 0)
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	pm = newM.(Model)
	if !pm.Completed {
		t.Errorf("expected Completed to be true after selecting path")
	}
	if pm.SelectedLocation.Name == "" {
		t.Errorf("expected SelectedLocation to be set")
	}

	// Verify Location was stored in GameState in DB
	state, err := app.LoadActiveGameStateFromPath(dbPath)
	if err != nil {
		t.Fatalf("LoadActiveGameStateFromPath failed: %v", err)
	}
	if state.CurrentLocationID == nil {
		t.Fatalf("expected CurrentLocationID in GameState to be non-nil")
	}
	if state.CurrentLocation == nil || state.CurrentLocation.Name != pm.SelectedLocation.Name {
		t.Fatalf("expected GameState CurrentLocation name '%s', got '%v'", pm.SelectedLocation.Name, state.CurrentLocation)
	}
}

func TestPathScreen_ResponsiveCardSizingAndCentering(t *testing.T) {
	cfg := domain.DefaultConfig()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "path_responsive_test.db")

	m := NewModelWithDBPath(cfg, dbPath)

	// 1. MD Breakpoint (60 cols)
	m.SetDimensions(60, 24)
	viewMD := m.View()
	if viewMD == "" {
		t.Errorf("expected view on MD breakpoint to be non-empty")
	}

	// Verify each line of viewMD is formatted to screen bounds
	linesMD := strings.Split(viewMD, "\n")
	if len(linesMD) != 24 {
		t.Errorf("expected 24 lines on MD breakpoint view, got %d", len(linesMD))
	}

	// 2. LG Breakpoint (80 cols)
	m.SetDimensions(80, 24)
	viewLG := m.View()
	if viewLG == "" {
		t.Errorf("expected view on LG breakpoint to be non-empty")
	}

	// 3. XL Breakpoint (120 cols)
	m.SetDimensions(120, 40)
	viewXL := m.View()
	if viewXL == "" {
		t.Errorf("expected view on XL breakpoint to be non-empty")
	}
}

func TestPathScreen_MouseHoverAndClick(t *testing.T) {
	mouse.InitManager()
	cfg := domain.DefaultConfig()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "path_mouse_test.db")

	model := NewModelWithDBPath(cfg, dbPath)

	// Render view so bubblezone registers zones
	_ = mouse.Scan(model.View())
	time.Sleep(20 * time.Millisecond)

	zInfo := mouse.Get("path_card_1")
	if zInfo == nil {
		t.Fatalf("expected zone info for path_card_1")
	}

	// 1. Mouse hover over path card 1 -> ActiveIndex becomes 1
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

	// 2. Mouse click (EventClick) -> PressedPathIndex becomes 1 (active_focused) & returns next cycle cmd
	pressMsg := tea.MouseMsg{
		X:      zInfo.StartX + 5,
		Y:      zInfo.StartY + 2,
		Button: tea.MouseButtonLeft,
	}
	updated, cmd := m.Update(pressMsg)
	m = updated.(Model)
	if m.PressedPathIndex != 1 {
		t.Errorf("expected PressedPathIndex 1 on mouse press, got %d", m.PressedPathIndex)
	}
	if cmd == nil {
		t.Fatalf("expected next cycle cmd on mouse click")
	}

	// 3. Next cycle message -> PressedPathIndex resets to -1 and confirms path selection
	nextCycleMsg := cmd()
	updated, _ = m.Update(nextCycleMsg)
	m = updated.(Model)
	if m.PressedPathIndex != -1 {
		t.Errorf("expected PressedPathIndex -1 on next cycle, got %d", m.PressedPathIndex)
	}
}
