package settings

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/pomquest/internal/ui"
)

// HelpComponent represents the settings-specific help footer component.
type HelpComponent struct {
	keys  []key.Binding
	theme ui.Theme
}

func NewHelpComponent(theme ui.Theme, keys []key.Binding) HelpComponent {
	return HelpComponent{
		keys:  keys,
		theme: theme,
	}
}

func (c HelpComponent) View() string {
	return ui.NewHelpComponent(c.theme, c.keys).View()
}
