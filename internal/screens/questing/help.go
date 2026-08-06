package questing

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/pomquest/internal/ui"
)

// HelpComponent wraps ui.HelpComponent for questing keybindings.
type HelpComponent struct {
	theme ui.Theme
	keys  []key.Binding
}

func NewHelpComponent(theme ui.Theme, keys []key.Binding) HelpComponent {
	return HelpComponent{
		theme: theme,
		keys:  keys,
	}
}

func (c HelpComponent) View() string {
	return ui.NewHelpComponent(c.theme, c.keys).View()
}
