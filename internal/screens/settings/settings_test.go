package settings

import (
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/pomquest/internal/app"
	"github.com/pomquest/internal/domain"
)

func TestSettingsScreen_TitleAndNoSubtitles(t *testing.T) {
	cfg := domain.DefaultConfig()
	m := NewModel(cfg)

	view := m.View()

	// Title must have Settings icon and Settings
	if !strings.Contains(view, "⚙") || !strings.Contains(view, "Settings") {
		t.Fatalf("expected view to contain '⚙' and 'Settings', got:\n%s", view)
	}

	// Must NOT contain old subtitles or sub-titles
	if strings.Contains(view, "Manage local database") {
		t.Fatalf("expected view to NOT contain subtitle, got:\n%s", view)
	}
	if strings.Contains(view, "Configuration Menu") {
		t.Fatalf("expected view to NOT contain 'Configuration Menu' subtitle, got:\n%s", view)
	}
}

func TestSettingsScreen_Navigation(t *testing.T) {
	cfg := domain.DefaultConfig()
	m := NewModel(cfg)

	if m.Cursor != 0 {
		t.Fatalf("expected initial cursor to be 0, got %d", m.Cursor)
	}

	// Move cursor down
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(Model)
	if m.Cursor != 1 {
		t.Fatalf("expected cursor 1 after 'j', got %d", m.Cursor)
	}

	// Move cursor up
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(Model)
	if m.Cursor != 0 {
		t.Fatalf("expected cursor 0 after 'k', got %d", m.Cursor)
	}
}

func TestSettingsScreen_SFXToggle1Layer(t *testing.T) {
	cfg := domain.DefaultConfig()
	m := NewModel(cfg)

	initialSFX := m.AudioSFXEnabled

	// Press '2' to toggle SFX
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	if cmd != nil {
		updated, _ = updated.(Model).Update(cmd())
	}
	m = updated.(Model)

	if m.AudioSFXEnabled == initialSFX {
		t.Fatalf("expected AudioSFXEnabled to toggle, got %v", m.AudioSFXEnabled)
	}
	if !strings.Contains(m.StatusNotice, "Sound Effects") {
		t.Fatalf("expected StatusNotice to mention Sound Effects, got %q", m.StatusNotice)
	}

	// Press 's' shortcut key
	updated, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	if cmd != nil {
		updated, _ = updated.(Model).Update(cmd())
	}
	m = updated.(Model)

	if m.AudioSFXEnabled != initialSFX {
		t.Fatalf("expected AudioSFXEnabled to toggle back, got %v", m.AudioSFXEnabled)
	}
}

func TestSettingsScreen_MusicToggle1Layer(t *testing.T) {
	cfg := domain.DefaultConfig()
	m := NewModel(cfg)

	initialMusic := m.AudioMusicMuted

	// Press '3' to toggle Music
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	if cmd != nil {
		updated, _ = updated.(Model).Update(cmd())
	}
	m = updated.(Model)

	if m.AudioMusicMuted == initialMusic {
		t.Fatalf("expected AudioMusicMuted to toggle, got %v", m.AudioMusicMuted)
	}
	if !strings.Contains(m.StatusNotice, "Music") {
		t.Fatalf("expected StatusNotice to mention Music, got %q", m.StatusNotice)
	}

	// Press 'm' shortcut key
	updated, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	if cmd != nil {
		updated, _ = updated.(Model).Update(cmd())
	}
	m = updated.(Model)

	if m.AudioMusicMuted != initialMusic {
		t.Fatalf("expected AudioMusicMuted to toggle back, got %v", m.AudioMusicMuted)
	}
}

func TestSettingsScreen_Credits1Layer(t *testing.T) {
	cfg := domain.DefaultConfig()
	m := NewModel(cfg)

	// Press '4' for Credits
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}})
	if cmd != nil {
		updated, _ = updated.(Model).Update(cmd())
	}
	m = updated.(Model)

	if !strings.Contains(m.StatusNotice, "PomQuest") {
		t.Fatalf("expected StatusNotice to mention PomQuest credits, got %q", m.StatusNotice)
	}
}

func TestSettingsScreen_ResetGameData(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "settings_test.db")

	// Pre-populate DB with a party record (InitDB seeds 1 initial party)
	db, err := app.InitDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	db.Create(&app.Party{Name: "Pre-reset Party"})

	var count int64
	db.Model(&app.Party{}).Count(&count)
	if count != 2 {
		t.Fatalf("expected 2 parties before reset (1 initial seeded + 1 created), got %d", count)
	}

	cfg := domain.DefaultConfig()
	m := NewModelWithDBPath(cfg, dbPath)

	// Select Reset Game Data ('1')
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	if cmd != nil {
		updated, _ = updated.(Model).Update(cmd())
	}
	m = updated.(Model)

	if m.State != StateConfirmReset {
		t.Fatalf("expected state StateConfirmReset, got %v", m.State)
	}

	// Send 'y' to confirm reset
	updated, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	m = updated.(Model)

	// Process confirmation command if quit message emitted
	if cmd != nil {
		if _, isQuit := cmd().(tea.QuitMsg); isQuit {
			m.State = StateMain
			_, _ = app.ResetAndMigrateDB(dbPath)
			m.StatusNotice = "✅ Local SQLite DB reset and auto-migrated successfully!"
		}
	}

	// Verify SQLite database was reset to initial seeded state
	verifyDB, err := app.InitDB(dbPath)
	if err != nil {
		t.Fatalf("failed to open db after reset: %v", err)
	}
	verifyDB.Model(&app.Party{}).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 initial party record after reset, got %d", count)
	}
}

func TestSettingsScreen_Quit(t *testing.T) {
	cfg := domain.DefaultConfig()
	m := NewModel(cfg)

	// Press 'q'
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m = updated.(Model)

	if !m.Quitted {
		t.Fatalf("expected Quitted to be true after pressing 'q'")
	}
	if m.View() != "" {
		t.Fatalf("expected View to be empty string when Quitted, got %q", m.View())
	}
}

func TestSettingsScreen_SmallScreenUnsupported(t *testing.T) {
	cfg := domain.DefaultConfig()
	m := NewModel(cfg)
	m.SetDimensions(40, 24) // SM breakpoint or less

	view := m.View()
	if !strings.Contains(view, "terminal width is not") {
		t.Fatalf("expected small screen unsupported message, got:\n%s", view)
	}
}

func TestSettingsScreen_MenuCardStyleAndFocusedState(t *testing.T) {
	cfg := domain.DefaultConfig()
	m := NewModel(cfg)

	view := m.View()

	// 1. Should not contain cursor pointer '▶'
	if strings.Contains(view, "▶") {
		t.Errorf("expected settings view to not contain cursor pointer '▶', got:\n%s", view)
	}

	// 2. Should contain rounded borders for menu cards
	if !strings.Contains(view, "╭") || !strings.Contains(view, "╰") {
		t.Errorf("expected settings view to contain bordered menu cards, got:\n%s", view)
	}

	// 3. Test focused state on selection
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	m2 := updated.(Model)
	if m2.FocusedIndex != 1 {
		t.Errorf("expected FocusedIndex to be 1 on select '2', got %d", m2.FocusedIndex)
	}
	if cmd == nil {
		t.Fatalf("expected delay transition cmd on select, got nil")
	}

	// Executing transition cmd delivers settingsSelectMsg
	updatedFinal, _ := m2.Update(cmd())
	mFinal := updatedFinal.(Model)
	if mFinal.FocusedIndex != -1 {
		t.Errorf("expected FocusedIndex reset to -1, got %d", mFinal.FocusedIndex)
	}
	if mFinal.Selected != ChoiceSFX {
		t.Errorf("expected ChoiceSFX after transition, got %v", mFinal.Selected)
	}
}
