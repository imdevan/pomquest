package dashboard

import "github.com/pomquest/internal/ui"

// MenuComponent represents the dashboard-specific menu component.
type MenuComponent struct {
	choices      []string
	cursor       int
	focusedIndex int
	width        int
	theme        ui.Theme
}

func NewMenuComponent(theme ui.Theme, choices []string, cursor int) MenuComponent {
	return MenuComponent{
		choices:      choices,
		cursor:       cursor,
		focusedIndex: -1,
		width:        0,
		theme:        theme,
	}
}

func (c MenuComponent) WithFocusedIndex(idx int) MenuComponent {
	c.focusedIndex = idx
	return c
}

func (c MenuComponent) WithWidth(w int) MenuComponent {
	c.width = w
	return c
}

func (c MenuComponent) View() string {
	comp := ui.NewMenuComponent(c.theme, c.choices, c.cursor)
	comp.IDPrefix = "dashboard_menu"
	comp.FocusedIndex = c.focusedIndex
	comp.Width = c.width
	return comp.View()
}

