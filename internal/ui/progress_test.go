package ui

import (
	"testing"

	"github.com/pomquest/internal/domain"
)

func TestProgressBarComponent(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ThemeFromConfig(cfg)

	comp := NewProgressBarComponent(theme, 0.5, 40)
	view := comp.View()

	if len(view) == 0 {
		t.Errorf("expected view to render progress bar, got empty string")
	}
}
