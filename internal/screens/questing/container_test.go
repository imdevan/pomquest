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

func TestQuestContainerComponent_SmallScreenUnsupported(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ui.ThemeFromConfig(cfg)

	rmSM := ui.NewResponsiveManager(50)
	c := NewQuestContainerComponent(theme, 50, rmSM)

	view := c.Render("Section 1 Content")

	if !strings.Contains(view, "terminal width is not currently") {
		t.Errorf("expected container view on SM breakpoint to display unsupported width message, got:\n%s", view)
	}
	// Verify border characters exist
	if !strings.Contains(view, "╭") && !strings.Contains(view, "─") {
		t.Errorf("expected container view to render full pane border, got:\n%s", view)
	}
}

func TestQuestContainerComponent_FourSectionsSpacing(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ui.ThemeFromConfig(cfg)
	c := NewQuestContainerComponent(theme, 96)

	s0 := "Quest Log Section"
	s1 := "Hero Card Section"
	s2 := "Party Card Section"
	s3 := "Quest Summary Section"

	view := c.Render(s0, s1, s2, s3)

	rawView := s0 + "\n\n" + s1 + "\n" + s2 + "\n\n" + s3
	lines := strings.Split(view, "\n")
	var trimmedLines []string
	for _, l := range lines {
		trimmedLines = append(trimmedLines, strings.TrimSpace(l))
	}
	compactView := strings.Join(trimmedLines, "\n")

	if !strings.Contains(compactView, rawView) {
		t.Errorf("expected 4-section view layout to match:\n%s\ngot:\n%s", rawView, compactView)
	}
}

