package questing

import "github.com/pomquest/internal/ui"

// SuppliesComponent wraps ui.PartyBarComponent matching mockup specs.
type SuppliesComponent struct {
	theme ui.Theme
	data  ui.PartyBarData
}

func NewSuppliesComponent(theme ui.Theme, data ui.PartyBarData) SuppliesComponent {
	return SuppliesComponent{
		theme: theme,
		data:  data,
	}
}

func (c SuppliesComponent) View() string {
	return ui.NewPartyBarComponent(c.theme, c.data, 72).View()
}
