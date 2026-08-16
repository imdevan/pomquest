package ui

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pomquest/internal/adapters/mouse"
	"github.com/pomquest/internal/domain"
)

func TestMenuComponent_CardRenderingAndSpacing(t *testing.T) {
	theme := ThemeFromConfig(domain.DefaultConfig())
	choices := []string{"1. Option One", "2. Option Two", "3. Option Three"}
	menu := NewMenuComponent(theme, choices, 0)
	menu.Width = 40

	viewStr := menu.View()

	// 1. Should not contain old cursor pointer '▶'
	if strings.Contains(viewStr, "▶") {
		t.Errorf("expected view to not contain cursor pointer '▶', got:\n%s", viewStr)
	}

	// 2. Should contain rounded and shadow border characters (╭, ╮, ╰, ╯, │)
	if !strings.Contains(viewStr, "╭") || !strings.Contains(viewStr, "╰") || !strings.Contains(viewStr, "╮") {
		t.Errorf("expected view to render bordered cards with shadow border, got:\n%s", viewStr)
	}

	// 3. Spacing: lines between cards include the shadow border row and bottom margin line
	lines := strings.Split(viewStr, "\n")
	// Each card has 4 lines (elevated card with shadow border, or standard position with top offset).
	// With Margin(0, 1, 1, 1) on non-last cards (5 lines each) and Margin(0, 1) on the last card (4 lines):
	// Total lines = 5 + 5 + 4 = 14 lines
	if len(lines) != 14 {
		t.Errorf("expected 14 total lines (3 cards with bottom margins on non-last and elevated active card), got %d lines:\n%s", len(lines), viewStr)
	}

	// 4. Test focused state rendering
	menu.FocusedIndex = 1
	focusedView := menu.View()
	if !strings.Contains(focusedView, "2. Option Two") {
		t.Errorf("expected focused view to contain choice text, got:\n%s", focusedView)
	}
}

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
