package dashboard

import "github.com/pomquest/internal/ui"

// HeaderComponent represents the dashboard-specific header component.
type HeaderComponent struct {
	theme ui.Theme
}

func NewHeaderComponent(theme ui.Theme) HeaderComponent {
	return HeaderComponent{theme: theme}
}

func (c HeaderComponent) View() string {
	return ui.NewHeaderComponent(
		c.theme,
		"🏰 POMQUEST - HEARTH DASHBOARD",
		"Welcome back, Maia",
	).View()
}
