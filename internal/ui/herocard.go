package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pomquest/internal/adapters/icon"
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
	Theme      Theme
	DataModel  HeroCardData
	Responsive *ResponsiveManager
}

func NewHeroCardComponent(theme Theme, data HeroCardData, responsive ...*ResponsiveManager) HeroCardComponent {
	var rm *ResponsiveManager
	if len(responsive) > 0 && responsive[0] != nil {
		rm = responsive[0]
	}
	return HeroCardComponent{
		DataModel:  data,
		Theme:      theme,
		Responsive: rm,
	}
}

func (c HeroCardComponent) View() string {
	rm := c.Responsive
	if rm == nil {
		rm = NewResponsiveManager(80)
	}
	isMD := !rm.IsAtLeast(BreakpointLG)

	innerWidth := 14
	paddingH := 1
	if isMD {
		innerWidth = 10
		paddingH = 0
	}

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(c.Theme.Border).
		Padding(0, paddingH).
		Width(innerWidth + paddingH*2)

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

	heroName := c.DataModel.Name
	if isMD {
		heroName = abbreviateName(heroName)
	}

	// Line 1: Name (Top Left), Level (Top Right)
	var nameStr, levelStr string
	if c.DataModel.IsDead {
		nameStr = deadStyle.Render(heroName)
		levelStr = deadStyle.Render(fmt.Sprintf("%d", c.DataModel.Level))
	} else {
		nameStr = nameStyle.Render(heroName)
		levelStr = classStyle.Render(fmt.Sprintf("%d", c.DataModel.Level))
	}
	topLine := spaceBetween(nameStr, levelStr, innerWidth)

	// Line 2: Class with Icon (Top Left)
	classIcon := c.DataModel.Icon
	if classIcon == "" {
		classIcon = icon.GetClassIcon(c.DataModel.Class, icon.ModeNerdFont)
	}

	var classDisplay string
	if isMD {
		// On MD: only show class icon
		classDisplay = classIcon
	} else {
		// On LG+: show class icon and class name
		if classIcon != "" {
			classDisplay = fmt.Sprintf("%s %s", classIcon, c.DataModel.Class)
		} else {
			classDisplay = c.DataModel.Class
		}
	}

	var classStr string
	if c.DataModel.IsDead {
		classStr = deadStyle.Render(classDisplay)
	} else {
		classStr = classStyle.Render(classDisplay)
	}

	// Line 3: Status effects line (line above hp / atk / def row)
	line3 := ""
	if c.DataModel.Status != "" && c.DataModel.Status != "Healthy" {
		statusText := c.DataModel.Status
		if isMD && len(statusText) > innerWidth {
			statusText = statusText[:innerWidth]
		}
		if c.DataModel.IsDead {
			line3 = deadStyle.Render(statusText)
		} else {
			line3 = statusStyle.Render(statusText)
		}
	}

	// Line 4: Atk & Def (Bottom Left), HP (Bottom Right)
	atkIcon := icon.GetNerdFont(icon.KeyAttack)
	armIcon := icon.GetNerdFont(icon.KeyArmor)

	var atkDefFmt string
	if isMD {
		// On MD breakpoint: remove space between attack icon & atk value, and defense icon & def value
		atkDefFmt = "%s%d%s%d"
	} else {
		// On LG+: keep spaces
		atkDefFmt = "%s %d %s %d"
	}

	var atkDefStr, hpStr string
	if c.DataModel.IsDead {
		atkDefStr = deadStyle.Render(fmt.Sprintf(atkDefFmt, atkIcon, c.DataModel.Attack, armIcon, c.DataModel.Armor))
		hpStr = deadStyle.Render(fmt.Sprintf("%d/%d", c.DataModel.HP, c.DataModel.MaxHP))
	} else {
		atkDefStr = statStyle.Render(fmt.Sprintf(atkDefFmt, atkIcon, c.DataModel.Attack, armIcon, c.DataModel.Armor))
		hpStr = statStyle.Render(fmt.Sprintf("%d/%d", c.DataModel.HP, c.DataModel.MaxHP))
	}
	bottomLine := spaceBetween(atkDefStr, hpStr, innerWidth)

	var sb strings.Builder
	sb.WriteString(topLine + "\n")
	sb.WriteString(classStr + "\n")
	sb.WriteString("\n")          // Extra line of space between hero class and status
	sb.WriteString(line3 + "\n") // Line reserved for status effects
	sb.WriteString(bottomLine)

	frontCard := cardStyle.Render(sb.String())

	if c.DataModel.IsActive {
		return compositeActiveOffsetCard(frontCard, c.Theme.Border)
	}

	// Non-active cards in a grid with an active card: top row offset by 1 empty line so active card raises above them
	return "\n" + frontCard
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

	cardWidth := lipgloss.Width(frontLines[0])
	borderStyle := lipgloss.NewStyle().Foreground(borderColor)

	var result []string
	// Row 0: top line of front card (raised 1 row up, 1 col left)
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

func abbreviateName(name string) string {
	if len(name) > 4 {
		return name[:3] + "."
	}
	return name
}
