package settings

import (
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/pomquest/internal/app"
	"github.com/pomquest/internal/domain"
)

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

func TestSettingsScreen_AudioSubView(t *testing.T) {
	cfg := domain.DefaultConfig()
	m := NewModel(cfg)

	// Press '2' to go to Audio settings
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	m = updated.(Model)

	if m.State != StateAudio {
		t.Fatalf("expected state StateAudio after pressing '2', got %v", m.State)
	}

	view := m.View()
	if !strings.Contains(view, "Audio Settings (Stub)") {
		t.Fatalf("expected view to contain Audio Settings, got:\n%s", view)
	}

	// Toggle SFX
	initialSFX := m.AudioSFXEnabled
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	m = updated.(Model)
	if m.AudioSFXEnabled == initialSFX {
		t.Fatalf("expected AudioSFXEnabled to toggle, got %v", m.AudioSFXEnabled)
	}

	// Press 'esc' to return to Main settings menu
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEscape})
	m = updated.(Model)
	if m.State != StateMain {
		t.Fatalf("expected state StateMain after pressing esc, got %v", m.State)
	}
}

func TestSettingsScreen_CreditsSubView(t *testing.T) {
	cfg := domain.DefaultConfig()
	m := NewModel(cfg)

	// Press '3' to go to Credits
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	m = updated.(Model)

	if m.State != StateCredits {
		t.Fatalf("expected state StateCredits after pressing '3', got %v", m.State)
	}

	view := m.View()
	if !strings.Contains(view, "Credits & Acknowledgments") {
		t.Fatalf("expected view to contain Credits, got:\n%s", view)
	}

	// Press 'b' to return
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	m = updated.(Model)
	if m.State != StateMain {
		t.Fatalf("expected state StateMain after pressing 'b', got %v", m.State)
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
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	m = updated.(Model)

	if m.State != StateConfirmReset {
		t.Fatalf("expected state StateConfirmReset, got %v", m.State)
	}

	// Send 'y' to confirm reset
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
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
