package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// HeroCardData holds values for a single hero card display.
type HeroCardData struct {
	Name     string
	Class    string
	Level    int
	HP       int
	MaxHP    int
	Attack   int
	Armor    int
	Status   string
	Icon     string
	IsActive bool
	IsDead   bool
}

// HeroCardComponent renders hero cards according to design specs.
type HeroCardComponent struct {
	Theme     Theme
	DataModel HeroCardData
}

func NewHeroCardComponent(theme Theme, data HeroCardData) HeroCardComponent {
	return HeroCardComponent{
		DataModel: data,
		Theme:     theme,
	}
}

func (c HeroCardComponent) View() string {
	innerWidth := 16

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(c.Theme.Border).
		Padding(0, 1).
		Width(innerWidth + 2)

	if c.DataModel.IsActive {
		cardStyle = cardStyle.BorderForeground(c.Theme.Primary)
	} else if c.DataModel.IsDead {
		cardStyle = cardStyle.BorderForeground(c.Theme.Muted)
	}

	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(c.Theme.Primary)
	classStyle := lipgloss.NewStyle().Foreground(c.Theme.Secondary)
	statStyle := lipgloss.NewStyle().Foreground(c.Theme.Text)
	statusStyle := lipgloss.NewStyle().Foreground(c.Theme.Tags)
	deadStyle := lipgloss.NewStyle().Foreground(c.Theme.Muted).Strikethrough(true)

	// Line 1: Name (Top Left), Level (Top Right)
	var nameStr, levelStr string
	if c.DataModel.IsDead {
		nameStr = deadStyle.Render(c.DataModel.Name)
		levelStr = deadStyle.Render(fmt.Sprintf("%d", c.DataModel.Level))
	} else {
		nameStr = nameStyle.Render(c.DataModel.Name)
		levelStr = classStyle.Render(fmt.Sprintf("%d", c.DataModel.Level))
	}
	topLine := spaceBetween(nameStr, levelStr, innerWidth)

	// Line 2: Class (Top Left)
	var classStr string
	if c.DataModel.IsDead {
		classStr = deadStyle.Render(c.DataModel.Class)
	} else {
		classStr = classStyle.Render(c.DataModel.Class)
	}

	// Line 3: Status effects line (line above hp / atk / def row)
	line3 := ""
	if c.DataModel.Status != "" && c.DataModel.Status != "Healthy" {
		if c.DataModel.IsDead {
			line3 = deadStyle.Render(c.DataModel.Status)
		} else {
			line3 = statusStyle.Render(c.DataModel.Status)
		}
	}

	// Line 4: Atk & Def (Bottom Left), HP (Bottom Right)
	// Format: icon, space, number, space, icon, space, number -> 󰓥 2 🛡 2
	var atkDefStr, hpStr string
	if c.DataModel.IsDead {
		atkDefStr = deadStyle.Render(fmt.Sprintf("󰓥 %d 🛡 %d", c.DataModel.Attack, c.DataModel.Armor))
		hpStr = deadStyle.Render(fmt.Sprintf("%d/%d", c.DataModel.HP, c.DataModel.MaxHP))
	} else {
		atkDefStr = statStyle.Render(fmt.Sprintf("󰓥 %d 🛡 %d", c.DataModel.Attack, c.DataModel.Armor))
		hpStr = statStyle.Render(fmt.Sprintf("%d/%d", c.DataModel.HP, c.DataModel.MaxHP))
	}
	bottomLine := spaceBetween(atkDefStr, hpStr, innerWidth)

	var sb strings.Builder
	sb.WriteString(topLine + "\n")
	sb.WriteString(classStr + "\n")
	sb.WriteString(line3 + "\n") // Line reserved for status effects
	sb.WriteString(bottomLine)

	frontCard := cardStyle.Render(sb.String())

	if c.DataModel.IsActive {
		return compositeActiveOffsetCard(frontCard, c.Theme.Border)
	}

	return frontCard
}

func spaceBetween(left, right string, width int) string {
	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	spaces := width - leftWidth - rightWidth
	if spaces < 1 {
		spaces = 1
	}
	return left + strings.Repeat(" ", spaces) + right
}

// compositeActiveOffsetCard overlaps frontCard (shifted 1 up, 1 left) over a base border box (standard position)
func compositeActiveOffsetCard(frontCard string, borderColor lipgloss.Color) string {
	frontLines := strings.Split(frontCard, "\n")
	height := len(frontLines)
	if height == 0 {
		return frontCard
	}

	// Determine physical visual width of frontCard line
	cardWidth := lipgloss.Width(frontLines[0])
	borderStyle := lipgloss.NewStyle().Foreground(borderColor)

	var result []string
	// Row 0: top line of front card (shifted 1 up, 1 left)
	result = append(result, frontLines[0]+" ")

	// Row 1: front line 1 + base box top-right corner ╮
	if height > 1 {
		result = append(result, frontLines[1]+borderStyle.Render("╮"))
	}

	// Rows 2..height-2: front line + base box vertical right border │
	for i := 2; i < height-1; i++ {
		result = append(result, frontLines[i]+borderStyle.Render("│"))
	}

	// Row height-1: front bottom line + base box vertical right border │
	if height > 2 {
		result = append(result, frontLines[height-1]+borderStyle.Render("│"))
	}

	// Bottom Row (row height): 1 space + base box bottom border ╰──────────────╯
	if cardWidth > 2 {
		bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", cardWidth-2) + "╯")
		result = append(result, " "+bottomBorder)
	}

	return strings.Join(result, "\n")
}
