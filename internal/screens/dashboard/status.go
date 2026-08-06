package dashboard

import (
	"fmt"

	"github.com/pomquest/internal/ui"
)

// StatusComponent represents the dashboard-specific status component.
type StatusComponent struct {
	location  string
	hoardGold int
	theme     ui.Theme
}

func NewStatusComponent(theme ui.Theme, location string, hoardGold int) StatusComponent {
	return StatusComponent{
		location:  location,
		hoardGold: hoardGold,
		theme:     theme,
	}
}

func (c StatusComponent) View() string {
	text := fmt.Sprintf("📍 Location: %s | 🐉 Hoard: %d Gold", c.location, c.hoardGold)
	return ui.NewStatusComponent(c.theme, text).View()
}
