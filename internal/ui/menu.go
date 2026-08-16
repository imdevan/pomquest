package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pomquest/internal/adapters/mouse"
)

// MenuComponent is a generic choice list selector widget.
type MenuComponent struct {
	IDPrefix string
	Choices  []string
	Cursor   int
	Theme    Theme
}

func NewMenuComponent(theme Theme, choices []string, cursor int) MenuComponent {
	return MenuComponent{
		IDPrefix: "menu",
		Choices:  choices,
		Cursor:   cursor,
		Theme:    theme,
	}
}

func (c MenuComponent) View() string {
	cursorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(c.Theme.Primary)

	normalStyle := lipgloss.NewStyle().
		Foreground(c.Theme.Text)

	prefix := c.IDPrefix
	if prefix == "" {
		prefix = "menu"
	}

	var sb strings.Builder
	for i, choice := range c.Choices {
		var line string
		if c.Cursor == i {
			line = cursorStyle.Render(fmt.Sprintf("▶ %s", choice))
		} else {
			line = normalStyle.Render(fmt.Sprintf("  %s", choice))
		}
		itemID := fmt.Sprintf("%s_item_%d", prefix, i)
		sb.WriteString(mouse.Mark(itemID, line) + "\n")
	}
	return sb.String()
}

// HandleMouse processes mouse events over menu choices.
// On hover: returns (hoveredIndex, false), highlighting item as if navigated by keyboard.
// On click: returns (clickedIndex, true), selecting item as if pressing enter.
func (c MenuComponent) HandleMouse(msg tea.Msg) (int, bool) {
	prefix := c.IDPrefix
	if prefix == "" {
		prefix = "menu"
	}

	mouseEvt, ok := msg.(tea.MouseMsg)
	if !ok {
		return c.Cursor, false
	}

	parsed := mouse.ParseEvent(mouseEvt)

	for i := range c.Choices {
		itemID := fmt.Sprintf("%s_item_%d", prefix, i)
		if mouse.InBounds(itemID, msg) {
			if parsed.Type == mouse.EventClick {
				return i, true
			}
			return i, false
		}
	}

	return c.Cursor, false
}
