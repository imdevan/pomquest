package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

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
