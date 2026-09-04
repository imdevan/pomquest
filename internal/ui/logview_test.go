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

func TestLogViewComponentScrolling(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ThemeFromConfig(cfg)

	messages := []string{
		"Line 1: Log event alpha",
		"Line 2: Log event beta",
		"Line 3: Log event gamma",
		"Line 4: Log event delta",
		"Line 5: Log event epsilon",
	}

	// MaxRows = 2. Total messages = 5.
	// Offset = 0: Should show Line 4 and Line 5
	comp0 := NewLogViewComponent(theme, "Log", messages, 2, 50)
	view0 := comp0.View()
	if !strings.Contains(view0, "Line 4") || !strings.Contains(view0, "Line 5") {
		t.Errorf("offset=0 expected Line 4 and Line 5, got:\n%s", view0)
	}
	if strings.Contains(view0, "Line 1") {
		t.Errorf("offset=0 should not contain Line 1, got:\n%s", view0)
	}

	// Offset = 3 (scrolled back 3 lines): Should show Line 1 and Line 2
	comp3 := NewLogViewComponent(theme, "Log", messages, 2, 50)
	comp3.ScrollOffset = 3
	view3 := comp3.View()
	if !strings.Contains(view3, "Line 1") || !strings.Contains(view3, "Line 2") {
		t.Errorf("offset=3 expected Line 1 and Line 2, got:\n%s", view3)
	}
	if strings.Contains(view3, "Line 5") {
		t.Errorf("offset=3 should not contain Line 5, got:\n%s", view3)
	}
}

func TestLogViewComponentFixedHeightPadding(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ThemeFromConfig(cfg)

	// An 8-row viewport should always have exactly 10 lines (1 top border + 8 rows + 1 bottom border)
	for msgCount := 0; msgCount <= 12; msgCount++ {
		var msgs []string
		for i := 0; i < msgCount; i++ {
			msgs = append(msgs, "Log event message")
		}
		comp := NewLogViewComponent(theme, "Log", msgs, 8, 60)
		view := comp.View()
		lines := strings.Split(strings.TrimRight(view, "\n"), "\n")
		if len(lines) != 10 {
			t.Errorf("expected LogViewComponent with MaxRows=8 to render exactly 10 lines for msgCount=%d, got %d lines:\n%s", msgCount, len(lines), view)
		}
	}
}
