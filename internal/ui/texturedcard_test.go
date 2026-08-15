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

func TestTexturedCardTopTopCorners(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ThemeFromConfig(cfg)

	card := NewTexturedCard(theme, 50, "Secondary Header Test Content")
	card.TopTopLeft = CornerText{Text: "2026-08-15", Margin: 1, PaddingLeft: 1, PaddingRight: 1}
	card.TopTopRight = CornerText{Text: "Dungeon Level 1", Margin: 1, PaddingLeft: 1, PaddingRight: 1}
	card.TopLeft = CornerText{Text: "Dark Caverns", Margin: 1, PaddingLeft: 1, PaddingRight: 1}
	card.TopRight = CornerText{Text: "Room 3", Margin: 1, PaddingLeft: 1, PaddingRight: 1}

	view := card.View()
	t.Logf("TexturedCard TopTop Corners View:\n%s", view)

	if !strings.Contains(view, "2026-08-15") {
		t.Errorf("expected view to contain TopTopLeft 2026-08-15, got:\n%s", view)
	}
	if !strings.Contains(view, "Dungeon Level 1") {
		t.Errorf("expected view to contain TopTopRight Dungeon Level 1, got:\n%s", view)
	}
	if !strings.Contains(view, "Dark Caverns") {
		t.Errorf("expected view to contain TopLeft Dark Caverns, got:\n%s", view)
	}
	if !strings.Contains(view, "Room 3") {
		t.Errorf("expected view to contain TopRight Room 3, got:\n%s", view)
	}

	lines := strings.Split(view, "\n")
	// Expected lines: top-top header line, secondary header line, content line, bottom border line
	if len(lines) < 4 {
		t.Errorf("expected at least 4 lines for card with secondary header row, got %d:\n%s", len(lines), view)
	}
}
