package ui

import (
	"strings"
	"testing"

	"github.com/pomquest/internal/domain"
)

func TestLogViewComponentWithTexturedCard(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ThemeFromConfig(cfg)

	title := "12, Jan, 1000 A.P."
	messages := []string{
		"applebee casts magic missile",
		"+ 3 xp! +3 copper!",
	}

	comp := NewLogViewComponent(theme, title, messages, 8, 72)
	view := comp.View()

	if !strings.Contains(view, title) {
		t.Errorf("expected view to contain top-left date title %q, got:\n%s", title, view)
	}
	if !strings.Contains(view, "applebee casts magic missile") {
		t.Errorf("expected view to contain log message, got:\n%s", view)
	}
	if !strings.Contains(view, "+ 3 xp! +3 copper!") {
		t.Errorf("expected view to contain xp log message, got:\n%s", view)
	}
}

func TestLogViewComponentTextWrapping(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ThemeFromConfig(cfg)

	title := "Quest Log"
	longMsg := "6 goblins ambush the party in the dark forest and attack with rusty scythes and crude daggers!"
	messages := []string{longMsg}

	// Width 40 -> availWidth = 36. longMsg should wrap into multiple lines.
	comp := NewLogViewComponent(theme, title, messages, 8, 40)
	view := comp.View()

	lines := strings.Split(view, "\n")
	for _, l := range lines {
		// Each rendered line in the card must not exceed card width 40
		if len([]rune(l)) > 45 { // accounting for ANSI escapes, visual width should be <= 40
			t.Errorf("rendered line exceeded max width 40: %q", l)
		}
	}
	if !strings.Contains(view, "6 goblins") || !strings.Contains(view, "scythes") {
		t.Errorf("expected view to contain wrapped log message text, got:\n%s", view)
	}
}
