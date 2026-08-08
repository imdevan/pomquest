package questing

import "github.com/pomquest/internal/ui"

// SuppliesComponent wraps ui.PartyCardComponent rendering party card table.
type SuppliesComponent struct {
	theme ui.Theme
	data  ui.PartyCardData
	width int
}

func NewSuppliesComponent(theme ui.Theme, data ui.PartyCardData, width int) SuppliesComponent {
	if width <= 0 {
		width = DefaultContainerWidth
	}
	return SuppliesComponent{
		theme: theme,
		data:  data,
		width: width,
	}
}

func (c SuppliesComponent) View() string {
	return ui.NewPartyCardComponent(c.theme, c.data, c.width).View()
}
