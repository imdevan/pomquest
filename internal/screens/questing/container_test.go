package questing

import (
	"strings"
	"testing"

	"github.com/pomquest/internal/domain"
	"github.com/pomquest/internal/ui"
)

func TestQuestContainerComponent(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ui.ThemeFromConfig(cfg)

	c := NewQuestContainerComponent(theme, 96)

	if c.Width != 96 {
		t.Errorf("expected container width to be 96, got %d", c.Width)
	}

	section1 := "Section 1 Content"
	section2 := "Section 2 Content"

	view := c.Render(section1, section2)

	if !strings.Contains(view, section1) {
		t.Errorf("expected container view to contain section1, got:\n%s", view)
	}
	if !strings.Contains(view, section2) {
		t.Errorf("expected container view to contain section2, got:\n%s", view)
	}
}
