package main

import (
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/pomquest/internal/app"
	"github.com/pomquest/internal/domain"
)

func TestRootCommandHasVersion(t *testing.T) {
	cmd := newRootCmd()
	versionFlag := cmd.Flags().Lookup("version")
	if versionFlag == nil {
		t.Fatal("expected --version flag to be registered")
	}
	if versionFlag.Shorthand != "v" {
		t.Fatalf("expected shorthand -v, got %q", versionFlag.Shorthand)
	}
}

func TestRootCommandHasConfig(t *testing.T) {
	cmd := newRootCmd()
	configFlag := cmd.Flags().Lookup("config")
	if configFlag == nil {
		t.Fatal("expected --config flag to be registered")
	}
}

func TestResolvedVersion(t *testing.T) {
	ver := resolvedVersion()
	if ver == "" {
		t.Fatal("expected version to be non-empty")
	}
}

func TestAppModel_SettingsNavigation(t *testing.T) {
	cfg := domain.DefaultConfig()
	app := newAppModel(cfg)

	if app.state != screenDashboard {
		t.Fatalf("expected initial state screenDashboard, got %v", app.state)
	}

	// Move to ChoiceSettings (index 4) on dashboard
	for i := 0; i < 4; i++ {
		updated, _ := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		app = updated.(appModel)
	}

	// Press Enter to navigate to Settings
	updated, _ := app.Update(tea.KeyMsg{Type: tea.KeyEnter})
	app = updated.(appModel)

	if app.state != screenSettings {
		t.Fatalf("expected state screenSettings after selecting settings choice, got %v", app.state)
	}

	// Press Esc to exit settings and return to Dashboard
	updated, _ = app.Update(tea.KeyMsg{Type: tea.KeyEscape})
	app = updated.(appModel)

	if app.state != screenDashboard {
		t.Fatalf("expected state screenDashboard after quitting settings, got %v", app.state)
	}

	// Press 'j' on Dashboard — verify user remains on Dashboard and isn't bounced back into Settings
	updated, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	app = updated.(appModel)

	if app.state != screenDashboard {
		t.Fatalf("expected state to remain screenDashboard on keypress after returning, got %v", app.state)
	}
}

func TestAppModel_PathNavigation(t *testing.T) {
	cfg := domain.DefaultConfig()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "path_nav_test.db")

	appModelObj := newAppModelWithDB(cfg, dbPath)

	if appModelObj.state != screenDashboard {
		t.Fatalf("expected initial state screenDashboard, got %v", appModelObj.state)
	}

	// 1. Choice 0 is ChoiceQuest. Press '1' to trigger choice quest when party size is 0
	updated, _ := appModelObj.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	appModelObj = updated.(appModel)

	if appModelObj.state != screenDashboard {
		t.Fatalf("expected state to remain screenDashboard when party size is 0, got %v", appModelObj.state)
	}

	// 2. Add hero to party
	db, err := app.InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	hero := app.Hero{
		Creature: app.Creature{
			Noun:  app.Noun{Name: "Conan"},
			HP:    10,
			MaxHP: 10,
		},
		Class: "barbarian",
	}
	if err := app.ReplaceOrAddHeroInParty(db, 0, hero); err != nil {
		t.Fatalf("ReplaceOrAddHeroInParty failed: %v", err)
	}

	// Reload appModel with hero in party
	appModelObj = newAppModelWithDB(cfg, dbPath)
	updated, _ = appModelObj.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	appModelObj = updated.(appModel)

	if appModelObj.state != screenPath {
		t.Fatalf("expected state screenPath after selecting ChoiceQuest with active party, got %v", appModelObj.state)
	}
}

