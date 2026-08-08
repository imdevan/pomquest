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
