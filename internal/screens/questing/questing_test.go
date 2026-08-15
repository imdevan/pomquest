package questing

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/pomquest/internal/app"
	"github.com/pomquest/internal/domain"
	"github.com/pomquest/internal/ui"
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

func TestHeroAttackAnimationsFromLog(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_anim.db")

	cfg := domain.DefaultConfig()
	model := NewModelWithDBPath(cfg, dbPath)

	// 1. Attack animation log: " jimmothy hits 🦇 Bat for 3 dmg!"
	log1 := " jimmothy hits 🦇 Bat for 3 dmg!"
	_ = model.appendActionLog(log1)

	if model.Heroes[0].AnimState.Type != ui.AnimAttack {
		t.Errorf("expected jimmothy card to have AnimAttack state, got %v", model.Heroes[0].AnimState.Type)
	}
	if model.Heroes[0].AnimState.DamageText != "3!" {
		t.Errorf("expected damage text '3!', got %q", model.Heroes[0].AnimState.DamageText)
	}

	// View rendering contains damage row indicator "⚔️ 3!"
	view := model.View()
	if !strings.Contains(view, "3!") {
		t.Errorf("expected view to render '3!' in damage row, got:\n%s", view)
	}

	// 2. Separator line (. . . ) keeps animation state active!
	model.updatePacingDots()
	if model.Heroes[0].AnimState.Type != ui.AnimAttack {
		t.Errorf("expected animation state to remain AnimAttack during pacing dot separator line, got %v", model.Heroes[0].AnimState.Type)
	}

	// 3. Next phase line clears previous animation state!
	nextPhaseLog := "💀 🦇 Bat dies 😵"
	_ = model.appendActionLog(nextPhaseLog)
	if model.Heroes[0].AnimState.Type != ui.AnimNone {
		t.Errorf("expected hero AnimState to reset to AnimNone when next phase line begins, got %v", model.Heroes[0].AnimState.Type)
	}

	// 4. Block animation log: "jimmothy blocks the attack!"
	log2 := "jimmothy blocks the attack!"
	_ = model.appendActionLog(log2)
	if model.Heroes[0].AnimState.Type != ui.AnimBlock {
		t.Errorf("expected hero card to have AnimBlock state, got %v", model.Heroes[0].AnimState.Type)
	}

	// 5. Dodge animation log: "jimmothy dodges the attack!"
	log3 := "jimmothy dodges the attack!"
	_ = model.appendActionLog(log3)
	if model.Heroes[0].AnimState.Type != ui.AnimDodge {
		t.Errorf("expected hero card to have AnimDodge state, got %v", model.Heroes[0].AnimState.Type)
	}

	// 6. Take damage animation log: "🦇 Bat hits jimmothy for 2 dmg!"
	log4 := "🦇 Bat hits jimmothy for 2 dmg!"
	_ = model.appendActionLog(log4)
	if model.Heroes[0].AnimState.Type != ui.AnimTakeDamage {
		t.Errorf("expected hero card to have AnimTakeDamage state, got %v", model.Heroes[0].AnimState.Type)
	}
}

func TestHeroCardsComponent_HeightStabilityOnTakeDamage(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ui.ThemeFromConfig(cfg)
	rm := ui.NewResponsiveManager(80)

	normalHeroes := []ui.HeroCardData{
		{Name: "jimmothy", Class: "fighter", HP: 10, MaxHP: 10, Level: 1},
		{Name: "applebee", Class: "wizard", HP: 6, MaxHP: 6, Level: 1},
	}

	normalComp := NewHeroCardsComponent(theme, normalHeroes, 80, rm)
	normalView := normalComp.View()
	normalLineCount := strings.Count(normalView, "\n") + 1

	hitHeroes := []ui.HeroCardData{
		{
			Name:   "jimmothy",
			Class:  "fighter",
			HP:     8,
			MaxHP:  10,
			Level:  1,
			AnimState: ui.HeroAnimationState{
				Type:       ui.AnimTakeDamage,
				DamageText: "2",
				Icon:       "👾",
			},
		},
		{Name: "applebee", Class: "wizard", HP: 6, MaxHP: 6, Level: 1},
	}

	hitComp := NewHeroCardsComponent(theme, hitHeroes, 80, rm)
	hitView := hitComp.View()
	hitLineCount := strings.Count(hitView, "\n") + 1

	if normalLineCount != hitLineCount {
		t.Errorf("expected hero cards view line count to remain constant when hero takes damage, got normal=%d hit=%d", normalLineCount, hitLineCount)
	}
}

func TestQuestingScreen_BottomAlignmentAndMargins(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_questing_alignment.db")

	cfg := domain.DefaultConfig()
	model := NewModelWithDBPath(cfg, dbPath)
	model.Width = 100
	model.Height = 43

	view := model.View()

	lines := strings.Split(view, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	t.Logf("TOTAL LINES: %d\nVIEW:\n%s", len(lines), view)
	if len(lines) != 43 {
		t.Fatalf("expected view to have 43 lines matching window height, got %d lines", len(lines))
	}

	// 5.8.1: Top line should be empty (at least one line between top of screen and top of quest log)
	if strings.TrimSpace(lines[0]) != "" {
		t.Errorf("expected top line (line 0) to be empty for top margin, got: %q", lines[0])
	}

	// 5.8: Bottom line should be empty (1 row margin at bottom of screen)
	if strings.TrimSpace(lines[len(lines)-1]) != "" {
		t.Errorf("expected bottom line (line 42) to be empty for bottom margin, got: %q", lines[len(lines)-1])
	}

	// Verify quest progress summary is placed towards bottom before the last line margin
	fullText := strings.Join(lines, "\n")
	if !strings.Contains(fullText, "Dest:") {
		t.Errorf("expected view to contain quest progress destination, got:\n%s", fullText)
	}
}

func TestQuestingScreen_QuestLogHeaders(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_questing_headers.db")

	cfg := domain.DefaultConfig()
	model := NewModelWithDBPath(cfg, dbPath)
	model.LogTitle = "15, Aug, 1000 A.P."
	model.ProgressData.Destination = "Catacombs of Doom"
	model.CurrentRoom = 3

	view := model.View()
	t.Logf("Questing Screen Header View:\n%s", view)

	if !strings.Contains(view, "15, Aug, 1000 A.P.") {
		t.Errorf("expected quest log to render date (top-top-left), got:\n%s", view)
	}
	if !strings.Contains(view, "Catacombs of Doom") {
		t.Errorf("expected quest log to render dungeon name (top-left), got:\n%s", view)
	}
	if !strings.Contains(view, "Room #3") {
		t.Errorf("expected quest log to render room name (top-right), got:\n%s", view)
	}
}

func TestQuestingScreen_LogScrollingJK(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_questing_scrolling.db")

	cfg := domain.DefaultConfig()
	m := NewModelWithDBPath(cfg, dbPath)
	m.CharDelay = 0

	var logMsgs []string
	for i := 1; i <= 20; i++ {
		logMsgs = append(logMsgs, fmt.Sprintf("Event message number %d", i))
	}
	m.LogMessages = logMsgs

	// Pressing 'k' scrolls back up in log history
	updatedM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	modelK := updatedM.(Model)
	if modelK.LogScrollOffset != 1 {
		t.Errorf("expected LogScrollOffset to be 1 after pressing 'k', got %d", modelK.LogScrollOffset)
	}

	// Pressing 'j' scrolls back down towards live logs
	updatedM2, _ := modelK.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	modelJ := updatedM2.(Model)
	if modelJ.LogScrollOffset != 0 {
		t.Errorf("expected LogScrollOffset to be 0 after pressing 'j', got %d", modelJ.LogScrollOffset)
	}

	// While in scrolling mode, new log events buffer to pendingActions without printing
	(&modelK).appendActionLog("Queued Event While Scrolling")
	if len(modelK.pendingActions) != 1 || modelK.pendingActions[0] != "Queued Event While Scrolling" {
		t.Errorf("expected new log event to be queued in pendingActions while scrolling, got pending: %v", modelK.pendingActions)
	}

	// Pressing 'esc' while in scrolling mode resets LogScrollOffset to 0 and flushes queue
	if modelK.Paused {
		t.Errorf("modelK should not be paused")
	}
	escM, _ := modelK.Update(tea.KeyMsg{Type: tea.KeyEsc})
	modelEsc := escM.(Model)
	if modelEsc.LogScrollOffset != 0 {
		t.Errorf("expected LogScrollOffset to reset to 0 after pressing esc in scroll mode, got %d", modelEsc.LogScrollOffset)
	}
	if modelEsc.Paused {
		t.Errorf("expected pressing esc in scroll mode NOT to toggle Paused, but Paused was true")
	}
	if len(modelEsc.pendingActions) != 0 {
		t.Errorf("expected pendingActions to be empty after exiting scroll mode, got: %v", modelEsc.pendingActions)
	}
	lastMsg := modelEsc.LogMessages[len(modelEsc.LogMessages)-1]
	if lastMsg != "Queued Event While Scrolling" {
		t.Errorf("expected last log message to be 'Queued Event While Scrolling', got %q", lastMsg)
	}

	// Pressing 'esc' when NOT in scrolling mode toggles Paused
	escM2, _ := modelEsc.Update(tea.KeyMsg{Type: tea.KeyEsc})
	modelEsc2 := escM2.(Model)
	if !modelEsc2.Paused {
		t.Errorf("expected pressing esc when not in scroll mode to toggle Paused to true")
	}
}

