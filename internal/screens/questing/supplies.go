package questing

import "github.com/pomquest/internal/ui"

// SuppliesComponent wraps ui.PartyCardComponent rendering party card table.
type SuppliesComponent struct {
	theme ui.Theme
	data  ui.PartyCardData
}

func NewSuppliesComponent(theme ui.Theme, data ui.PartyCardData) SuppliesComponent {
	return SuppliesComponent{
		theme: theme,
		data:  data,
	}
}

func (c SuppliesComponent) View() string {
	return ui.NewPartyCardComponent(c.theme, c.data, 72).View()
}
