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
	IDPrefix     string
	Choices      []string
	Cursor       int
	FocusedIndex int
	Width        int
	Theme        Theme
}

func NewMenuComponent(theme Theme, choices []string, cursor int) MenuComponent {
	return MenuComponent{
		IDPrefix:     "menu",
		Choices:      choices,
		Cursor:       cursor,
		FocusedIndex: -1,
		Width:        0,
		Theme:        theme,
	}
}

func (c MenuComponent) View() string {
	prefix := c.IDPrefix
	if prefix == "" {
		prefix = "menu"
	}

	// Calculate uniform card inner width
	maxLen := 0
	for _, choice := range c.Choices {
		l := lipgloss.Width(choice)
		if l > maxLen {
			maxLen = l
		}
	}

	cardWidth := c.Width
	if cardWidth > 2 {
		cardWidth -= 2 // Account for 2 cols border (1 col each side) so box width sans margin equals c.Width
	}
	if cardWidth < maxLen {
		cardWidth = maxLen
	}
	if cardWidth < 16 {
		cardWidth = 16
	}

	normalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(c.Theme.Border).
		Foreground(c.Theme.Text).
		Padding(0, 2).
		Width(cardWidth)

	activeStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(c.Theme.Primary).
		Foreground(c.Theme.Primary).
		Bold(true).
		Padding(0, 2).
		Width(cardWidth)

	focusedStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(c.Theme.Primary).
		Foreground(c.Theme.TextHighlight).
		Bold(true).
		Padding(0, 2).
		Width(cardWidth)

	var cards []string
	for i, choice := range c.Choices {
		var frontCard string
		var renderedCard string
		// isLast := i == len(c.Choices)-1

		if i == c.FocusedIndex {
			frontCard = focusedStyle.Render(choice)
			renderedCard = PadCardLeft2x1(frontCard)
		} else if i == c.Cursor {
			frontCard = activeStyle.Render(choice)
			renderedCard = CompositeActiveOffsetCard2x1(frontCard, c.Theme.Border)
		} else {
			frontCard = normalStyle.Render(choice)
			renderedCard = PadCardLeft2x1(frontCard)
		}

		marginStyle := lipgloss.NewStyle()
		// if isLast {
		// marginStyle = marginStyle.Margin(0, 1)
		// } else {
		// marginStyle = marginStyle.Margin(0, 1, 1, 1)
		// }
		// marginStyle = marginStyle.Margin(0, 1)

		finalCard := marginStyle.Render(renderedCard)
		itemID := fmt.Sprintf("%s_item_%d", prefix, i)
		cards = append(cards, mouse.Mark(itemID, finalCard))
	}

	return strings.Join(cards, "\n")
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
