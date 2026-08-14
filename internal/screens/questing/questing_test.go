package questing

import (
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/pomquest/internal/app"
	"github.com/pomquest/internal/domain"
)

func TestQuestingScreen(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_questing_default.db")

	cfg := domain.DefaultConfig()
	model := NewModelWithDBPath(cfg, dbPath)

	// Test View rendering
	view := model.View()
	if !strings.Contains(view, model.ProgressData.Destination) {
		t.Errorf("expected view to contain generated destination %s, got:\n%s", model.ProgressData.Destination, view)
	}
	if !strings.Contains(view, "jimmothy") {
		t.Errorf("expected view to contain hero jimmothy, got:\n%s", view)
	}

	// Test pressing space to step automated dungeon cycle
	prevLogsCount := len(model.LogMessages)
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	m := updated.(Model)
	if len(m.LogMessages) < prevLogsCount {
		t.Errorf("expected log messages count to increase after space keypress, got %d", len(m.LogMessages))
	}

	// Test pause key
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	m = updated.(Model)
	if !m.Paused {
		t.Errorf("expected model to be paused after pressing 'p'")
	}
}

func TestQuestingScreenWithDB(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_questing.db")

	db, err := app.InitDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}

	// 1. Save party with custom supplies & coins
	party, err := app.LoadOrCreateActiveParty(db)
	if err != nil {
		t.Fatalf("failed to load active party: %v", err)
	}
	party.Coins = 4500
	party.FoodLbs = 1200
	party.SpareParts = 7
	party.Oxen = 3
	party.Diet = "heavy"
	party.Pace = "balanced"
	if err := app.SaveParty(db, party); err != nil {
		t.Fatalf("failed to save party: %v", err)
	}

	// 2. Save party heroes
	hero1 := app.Hero{
		Creature: app.Creature{
			Noun:   app.Noun{Name: "Thorin", Icon: "🛡️"},
			HP:     15,
			MaxHP:  20,
			Attack: 4,
			Armor:  3,
		},
		Class:  "fighter",
		Level:  3,
		Status: "Healthy",
	}
	hero2 := app.Hero{
		Creature: app.Creature{
			Noun:   app.Noun{Name: "Gandalf", Icon: "🧙"},
			HP:     12,
			MaxHP:  12,
			Attack: 5,
			Armor:  1,
		},
		Class:  "wizard",
		Level:  4,
		Status: "Healthy",
	}
	if err := app.ReplaceOrAddHeroInParty(db, 0, hero1); err != nil {
		t.Fatalf("failed to add hero1: %v", err)
	}
	if err := app.ReplaceOrAddHeroInParty(db, 1, hero2); err != nil {
		t.Fatalf("failed to add hero2: %v", err)
	}

	// 3. Save selected location
	loc := app.Location{
		Name:  "The Obsidian Citadel",
		Biome: "Volcanic Crag",
	}
	if _, err := app.SaveSelectedLocation(db, loc); err != nil {
		t.Fatalf("failed to save selected location: %v", err)
	}

	// 4. Initialize Questing model from DB
	cfg := domain.DefaultConfig()
	model := NewModelWithDBPath(cfg, dbPath)

	// Verify DB Data loaded
	if model.ProgressData.Destination != "The Obsidian Citadel" {
		t.Errorf("expected destination 'The Obsidian Citadel', got '%s'", model.ProgressData.Destination)
	}
	if model.Supplies.Coins != 4500 {
		t.Errorf("expected 4500 coins, got %d", model.Supplies.Coins)
	}
	if model.Supplies.FoodLbs != 1200 {
		t.Errorf("expected 1200 food lbs, got %d", model.Supplies.FoodLbs)
	}
	if model.Supplies.Diet != "heavy" {
		t.Errorf("expected diet 'heavy', got '%s'", model.Supplies.Diet)
	}

	// Verify Hero cards render based on current party heroes
	if len(model.Heroes) != 2 {
		t.Fatalf("expected 2 party heroes, got %d", len(model.Heroes))
	}
	if model.Heroes[0].Name != "Thorin" || model.Heroes[1].Name != "Gandalf" {
		t.Errorf("expected heroes Thorin and Gandalf, got %s and %s", model.Heroes[0].Name, model.Heroes[1].Name)
	}

	// Verify Party card hearts based on current party hero health
	if model.Supplies.CurrentHP != 27 {
		t.Errorf("expected CurrentHP 27, got %d", model.Supplies.CurrentHP)
	}
	if model.Supplies.MaxHP != 32 {
		t.Errorf("expected MaxHP 32, got %d", model.Supplies.MaxHP)
	}

	// Verify view rendering contains hero names and location
	view := model.View()
	if !strings.Contains(view, "The Obsidian Citadel") {
		t.Errorf("expected view to contain 'The Obsidian Citadel', got:\n%s", view)
	}
	if !strings.Contains(view, "Thorin") || !strings.Contains(view, "Gandalf") {
		t.Errorf("expected view to contain hero names Thorin and Gandalf, got:\n%s", view)
	}
}

func TestAutomatedDungeonCycleAndPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "dungeon_cycle.db")

	cfg := domain.DefaultConfig()
	model := NewModelWithDBPath(cfg, dbPath)

	// Step 1: Advance cycle several times using space bar
	for i := 0; i < 5; i++ {
		updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
		model = updated.(Model)
	}

	// Verify room and state advanced
	if model.CurrentRoom < 1 {
		t.Errorf("expected CurrentRoom >= 1, got %d", model.CurrentRoom)
	}

	// Step 2: Simulate exit and re-open (resume quest from SQLite DB)
	resumedModel := NewModelWithDBPath(cfg, dbPath)

	if resumedModel.ProgressData.Destination != model.ProgressData.Destination {
		t.Errorf("expected resumed destination '%s', got '%s'", model.ProgressData.Destination, resumedModel.ProgressData.Destination)
	}
	if resumedModel.CurrentRoom != model.CurrentRoom {
		t.Errorf("expected resumed room %d, got %d", model.CurrentRoom, resumedModel.CurrentRoom)
	}
	if len(resumedModel.LogMessages) == 0 {
		t.Errorf("expected restored game logs from DB, got 0")
	}
}

func TestCharacterDelay_TypingAnimation(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_char_delay.db")

	cfg := domain.DefaultConfig()
	cfg.CharDelay = 10
	model := NewModelWithDBPath(cfg, dbPath)

	if !model.isTyping {
		t.Errorf("expected isTyping to be true on initial load with CharDelay > 0")
	}

	// Step character ticks
	initialActiveLen := model.activeCharIdx
	updated, cmd := model.Update(charTickMsg{})
	model = updated.(Model)
	if cmd == nil {
		t.Errorf("expected non-nil cmd for next char tick")
	}
	if model.activeCharIdx != initialActiveLen+1 {
		t.Errorf("expected activeCharIdx to increment to %d, got %d", initialActiveLen+1, model.activeCharIdx)
	}

	// Press space to fast forward typing
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	model = updated.(Model)
	if model.isTyping {
		t.Errorf("expected isTyping to be false after space bar press fast forwards text")
	}

	// Test pacing dots do not block typing state
	model.TurnTickSecs = 1
	model.updatePacingDots()
	if model.isTyping {
		t.Errorf("expected updatePacingDots not to set isTyping")
	}

	lastLog := model.LogMessages[len(model.LogMessages)-1]
	if lastLog != ". " {
		t.Errorf("expected last log to be pacing dot '. ', got %q", lastLog)
	}
}
