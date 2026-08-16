package ui

import (
	"math"
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"

	"github.com/pomquest/internal/adapters/icon"
	"github.com/pomquest/internal/adapters/mouse"
)

// PathCardData holds values for rendering a dungeon path option card.
type PathCardData struct {
	Name               string
	Biome              string
	Description        string
	DifficultyModifier float64
	Icon               string
	IsActive           bool
}

// PathCardComponent renders path selection cards according to design specs.
type PathCardComponent struct {
	ZoneID     string
	Theme      Theme
	DataModel  PathCardData
	Responsive *ResponsiveManager
	Width      int
	Height     int
}

// NewPathCardComponent creates a new PathCardComponent.
func NewPathCardComponent(theme Theme, data PathCardData, responsive ...*ResponsiveManager) PathCardComponent {
	var rm *ResponsiveManager
	if len(responsive) > 0 && responsive[0] != nil {
		rm = responsive[0]
	}
	return PathCardComponent{
		Theme:      theme,
		DataModel:  data,
		Responsive: rm,
	}
}

// NewPathCardComponentSized creates a new PathCardComponent with explicit width and height.
func NewPathCardComponentSized(theme Theme, data PathCardData, width, height int, responsive ...*ResponsiveManager) PathCardComponent {
	var rm *ResponsiveManager
	if len(responsive) > 0 && responsive[0] != nil {
		rm = responsive[0]
	}
	return PathCardComponent{
		Theme:      theme,
		DataModel:  data,
		Responsive: rm,
		Width:      width,
		Height:     height,
	}
}

// View renders the path card string representation with active box hover effect.
func (c PathCardComponent) View() string {
	rm := c.Responsive
	if rm == nil {
		rm = NewResponsiveManager(80)
	}

	cardWidth := c.Width
	cardHeight := c.Height

	if cardWidth <= 0 && cardHeight <= 0 {
		switch {
		case rm.IsAtLeast(BreakpointXL):
			cardHeight = 18
		case rm.IsAtLeast(BreakpointLG):
			cardHeight = 15
		case rm.IsAtLeast(BreakpointMD):
			cardHeight = 12
		default:
			cardHeight = 10
		}
		cardWidth = int(math.Round(float64(cardHeight) * 1.5))
	} else if cardWidth <= 0 {
		cardWidth = int(math.Round(float64(cardHeight) * 1.5))
	} else if cardHeight <= 0 {
		cardHeight = int(math.Round(float64(cardWidth) / 1.5))
	}

	frontBoxHeight := cardHeight - 1
	if frontBoxHeight < 8 {
		frontBoxHeight = 8
	}

	paddingV := 1
	paddingH := 1

	if cardWidth >= 28 {
		paddingH = 2
	}
	if frontBoxHeight < 10 {
		paddingV = 0
	}

	styleWidth := cardWidth - 2
	if styleWidth < 10 {
		styleWidth = 10
	}
	styleHeight := frontBoxHeight - 2
	if styleHeight < 5 {
		styleHeight = 5
	}

	innerWidth := styleWidth - (paddingH * 2)
	if innerWidth < 8 {
		innerWidth = 8
	}

	contentHeight := styleHeight - (paddingV * 2)
	if contentHeight < 3 {
		contentHeight = 3
	}
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(c.Theme.Border).
		Padding(paddingV, paddingH).
		Width(styleWidth).
		Height(styleHeight).
		Align(lipgloss.Center, lipgloss.Center)

	if c.DataModel.IsActive {
		cardStyle = cardStyle.BorderForeground(c.Theme.Primary)
	}

	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(c.Theme.Primary).Align(lipgloss.Center)
	biomeStyle := lipgloss.NewStyle().Foreground(c.Theme.Tags).Align(lipgloss.Center)
	tildeStyle := lipgloss.NewStyle().Foreground(c.Theme.Muted).Align(lipgloss.Center)

	nameStr := c.DataModel.Name
	nameLines := WrapText(nameStr, innerWidth)
	nameContent := nameStyle.Render(strings.Join(nameLines, "\n"))

	var bodyContent string
	if c.DataModel.Biome != "" {
		iconStr := c.DataModel.Icon
		if iconStr == "" {
			iconStr = icon.GetBiomeIcon(c.DataModel.Biome, icon.ModeNerdFont)
		}
		biomeTitle := toTitleCase(c.DataModel.Biome)
		biomeContent := biomeStyle.Render(biomeTitle)
		iconContent := biomeStyle.Render(iconStr)
		tildeContent := tildeStyle.Render("~")

		bodyContent = nameContent + "\n\n\n" + tildeContent + "\n" + biomeContent + "\n" + iconContent
	} else {
		bodyContent = nameContent
	}

	frontCard := cardStyle.Render(bodyContent)

	var cardResult string
	if c.DataModel.IsActive {
		cardResult = compositeActiveOffsetCard2x1(frontCard, c.Theme.Border)
	} else {
		cardResult = "  \n" + frontCard
	}

	if c.ZoneID != "" {
		cardResult = mouse.Mark(c.ZoneID, cardResult)
	}

	return cardResult
}

func toTitleCase(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			r := []rune(w)
			r[0] = unicode.ToUpper(r[0])
			for j := 1; j < len(r); j++ {
				r[j] = unicode.ToLower(r[j])
			}
			words[i] = string(r)
		}
	}
	return strings.Join(words, " ")
}

