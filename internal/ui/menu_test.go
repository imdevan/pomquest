package ui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pomquest/internal/adapters/mouse"
	"github.com/pomquest/internal/domain"
)

func TestMenuComponent_HoverAndClick(t *testing.T) {
	mouse.InitManager()

	theme := ThemeFromConfig(domain.DefaultConfig())
	choices := []string{"Option 1", "Option 2", "Option 3"}
	menu := NewMenuComponent(theme, choices, 0)
	menu.IDPrefix = "test_menu"

	// Render view so bubblezone records zone coordinates
	viewStr := menu.View()
	_ = mouse.Scan(viewStr)

	time.Sleep(20 * time.Millisecond)

	zInfo1 := mouse.Get("test_menu_item_1")
	if zInfo1 == nil {
		t.Fatalf("expected zone info for test_menu_item_1")
	}

	// 1. Test mouse hover over Option 2 (index 1) highlights cursor
	hoverMsg := tea.MouseMsg{
		X:      zInfo1.StartX,
		Y:      zInfo1.StartY,
		Type:   tea.MouseMotion,
		Action: tea.MouseActionMotion,
	}

	newCursor, clicked := menu.HandleMouse(hoverMsg)
	if newCursor != 1 {
		t.Errorf("expected cursor to highlight index 1 on hover, got %d", newCursor)
	}
	if clicked {
		t.Errorf("expected clicked to be false on hover")
	}

	// 2. Test mouse click over Option 3 (index 2) selects item
	zInfo2 := mouse.Get("test_menu_item_2")
	if zInfo2 == nil {
		t.Fatalf("expected zone info for test_menu_item_2")
	}

	clickMsg := tea.MouseMsg{
		X:      zInfo2.StartX,
		Y:      zInfo2.StartY,
		Button: tea.MouseButtonLeft,
	}

	newCursor2, clicked2 := menu.HandleMouse(clickMsg)
	if newCursor2 != 2 {
		t.Errorf("expected cursor to be 2 on click, got %d", newCursor2)
	}
	if !clicked2 {
		t.Errorf("expected clicked to be true on mouse click")
	}
}
