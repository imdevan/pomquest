package settings

import "github.com/pomquest/internal/ui"

// HeaderComponent represents the settings-specific header component.
type HeaderComponent struct {
	theme ui.Theme
}

func NewHeaderComponent(theme ui.Theme) HeaderComponent {
	return HeaderComponent{theme: theme}
}

func (c HeaderComponent) View() string {
	return ui.NewHeaderComponent(
		c.theme,
		"⚙  Settings",
		"",
	).View()
}
