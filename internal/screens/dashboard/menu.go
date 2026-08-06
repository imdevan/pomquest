package dashboard

import "github.com/pomquest/internal/ui"

// MenuComponent represents the dashboard-specific menu component.
type MenuComponent struct {
	choices []string
	cursor  int
	theme   ui.Theme
}

func NewMenuComponent(theme ui.Theme, choices []string, cursor int) MenuComponent {
	return MenuComponent{
		choices: choices,
		cursor:  cursor,
		theme:   theme,
	}
}

func (c MenuComponent) View() string {
	return ui.NewMenuComponent(c.theme, c.choices, c.cursor).View()
}
