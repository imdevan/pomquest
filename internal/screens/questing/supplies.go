package questing

import "github.com/pomquest/internal/ui"

// SuppliesComponent wraps ui.PartyCardComponent rendering party card table.
type SuppliesComponent struct {
	theme      ui.Theme
	data       ui.PartyCardData
	width      int
	responsive *ui.ResponsiveManager
}

func NewSuppliesComponent(theme ui.Theme, data ui.PartyCardData, width int, responsive ...*ui.ResponsiveManager) SuppliesComponent {
	if width <= 0 {
		width = DefaultContainerWidth
	}
	var rm *ui.ResponsiveManager
	if len(responsive) > 0 && responsive[0] != nil {
		rm = responsive[0]
	} else {
		rm = ui.NewResponsiveManager(width)
	}
	return SuppliesComponent{
		theme:      theme,
		data:       data,
		width:      width,
		responsive: rm,
	}
}

func (c SuppliesComponent) View() string {
	return ui.NewPartyCardComponent(c.theme, c.data, c.width, c.responsive).View()
}
