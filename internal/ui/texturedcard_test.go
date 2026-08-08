package ui

import (
	"strings"
	"testing"

	"github.com/pomquest/internal/domain"
)

func TestTexturedCardFourCorners(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ThemeFromConfig(cfg)

	card := NewTexturedCard(theme, 40, "Hello World")
	card.TopLeft = CornerText{Text: "top-left", Margin: 1, PaddingLeft: 1, PaddingRight: 1}
	card.TopRight = CornerText{Text: "top-right", Margin: 1, PaddingLeft: 1, PaddingRight: 1}
	card.BottomLeft = CornerText{Text: "bottom-left", Margin: 1, PaddingLeft: 1, PaddingRight: 1}
	card.BottomRight = CornerText{Text: "bottom-right", Margin: 1, PaddingLeft: 1, PaddingRight: 1}

	view := card.View()
	t.Logf("TexturedCard Four Corners View:\n%s", view)

	if !strings.Contains(view, "top-left") {
		t.Errorf("expected view to contain top-left, got:\n%s", view)
	}
	if !strings.Contains(view, "top-right") {
		t.Errorf("expected view to contain top-right, got:\n%s", view)
	}
	if !strings.Contains(view, "bottom-left") {
		t.Errorf("expected view to contain bottom-left, got:\n%s", view)
	}
	if !strings.Contains(view, "bottom-right") {
		t.Errorf("expected view to contain bottom-right, got:\n%s", view)
	}
	if !strings.Contains(view, "Hello World") {
		t.Errorf("expected view to contain content Hello World, got:\n%s", view)
	}
}

func TestTexturedCardTextWrapping(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ThemeFromConfig(cfg)

	longContent := "This is a very long line of text that exceeds the card width and should automatically be wrapped."
	card := NewTexturedCard(theme, 30, longContent).Padding(0, 1)
	view := card.View()

	if !strings.Contains(view, "This is a very long") {
		t.Errorf("expected view to contain wrapped text, got:\n%s", view)
	}
}
