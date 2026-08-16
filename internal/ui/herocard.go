package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pomquest/internal/adapters/icon"
	"github.com/pomquest/internal/adapters/mouse"
	"github.com/pomquest/internal/app"
)

// AnimationType represents combat animation states for hero cards.
type AnimationType int

const (
	AnimNone AnimationType = iota
	AnimAttack
	AnimBlock
	AnimDodge
	AnimTakeDamage
)

// HeroAnimationState defines active combat animation parameters for a hero card.
type HeroAnimationState struct {
	Type             AnimationType
	DamageText       string
	Icon             string
	IsCrit           bool
	HorizontalOffset int
}

// HeroCardData holds values for a single hero card display.
type HeroCardData struct {
	Name           string
	Class          string
	Level          int
	HP             int
	MaxHP          int
	Attack         int
	Armor          int
	Status         string
	Icon           string
	IsActive       bool
	IsDead         bool
	Stats          app.Stats
	ShowFullStats  bool
	IsUnknownStats bool
	IsEmptySlot    bool
	AnimState      HeroAnimationState
}

// HeroCardComponent renders hero cards according to design specs.
type HeroCardComponent struct {
	ZoneID     string
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
	if c.DataModel.IsEmptySlot {
		return c.renderEmptySlotCard()
	}
	if c.DataModel.ShowFullStats {
		return c.renderFullStatCard()
	}
	rm := c.Responsive
	if rm == nil {
		rm = NewResponsiveManager(80)
	}
	isLG := rm.IsAtLeast(BreakpointLG)
	isMD := !isLG

	innerWidth := 14
	paddingV := 0
	paddingH := 1
	if !isLG {
		innerWidth = 10
		paddingH = 0
	} else {
		paddingV = 1
		paddingH = 2
	}

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(c.Theme.Border).
		Padding(paddingV, paddingH).
		Width(innerWidth + paddingH*2)

	if c.DataModel.AnimState.Type == AnimDodge {
		borderColor := c.Theme.Flags
		if string(borderColor) == "" {
			borderColor = lipgloss.Color("12")
		}
		cardStyle = cardStyle.BorderForeground(borderColor)
	} else if c.DataModel.AnimState.Type == AnimTakeDamage {
		borderColor := lipgloss.Color("01")
		cardStyle = cardStyle.BorderForeground(borderColor)
	} else if c.DataModel.IsActive || c.DataModel.AnimState.Type == AnimAttack {
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
	if classIcon != "" {
		classDisplay = fmt.Sprintf("%s %s", classIcon, c.DataModel.Class)
	} else {
		classDisplay = c.DataModel.Class
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
	sb.WriteString("\n")         // Extra line of space between hero class and status
	sb.WriteString(line3 + "\n") // Line reserved for status effects
	sb.WriteString(bottomLine)

	frontCard := cardStyle.Render(sb.String())

	if c.DataModel.AnimState.Type == AnimAttack {
		// Attack animation: front card moves 1 row UP while base box remains in standard active position (H + 2 + 1 = H + 3 lines)
		return compositeAttackOffsetCard(frontCard, c.Theme.Border) + "\n"
	}

	if c.DataModel.AnimState.Type == AnimBlock || c.DataModel.AnimState.Type == AnimTakeDamage {
		// Block or Take Damage animation: card moves 1 row DOWN (3 leading newlines) with under border showing (3 + H = H + 3 lines)
		return "\n\n\n" + compositeSameRowOffsetCard(frontCard, c.Theme.Border)
	}

	if c.DataModel.AnimState.Type == AnimDodge {
		// Dodge animation: card moves 1 left or right, flashing blue (1 + H + 2 = H + 3 lines)
		res := "\n" + frontCard + "\n\n"
		if c.DataModel.AnimState.HorizontalOffset != 0 {
			lines := strings.Split(res, "\n")
			for i := range lines {
				if c.DataModel.AnimState.HorizontalOffset > 0 {
					lines[i] = " " + lines[i]
				} else {
					lines[i] = strings.TrimPrefix(lines[i], " ")
				}
			}
			res = strings.Join(lines, "\n")
		}
		return res
	}

	if c.DataModel.IsActive {
		// Standard active card: 1 row UP relative to default (1 leading newline, 1 trailing newline: 1 + H + 1 + 1 = H + 3 lines)
		return "\n" + compositeActiveOffsetCard(frontCard, c.Theme.Border) + "\n"
	}

	// Standard non-active card: default position (2 leading newlines, 1 trailing newline: 2 + H + 1 = H + 3 lines)
	return "\n\n" + frontCard + "\n"
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

// compositeSameRowOffsetCard overlaps frontCard (shifted 1 col left, same row) over a base border box so under border shows (╮ │ ╯)
func compositeSameRowOffsetCard(frontCard string, borderColor lipgloss.Color) string {
	frontLines := strings.Split(frontCard, "\n")
	height := len(frontLines)
	if height == 0 {
		return frontCard
	}

	borderStyle := lipgloss.NewStyle().Foreground(borderColor)

	var result []string
	// Row 0: top line of front card + base box top-right corner ╮
	result = append(result, frontLines[0]+borderStyle.Render("╮"))

	// Rows 1..height-2: front line + base box vertical right border │
	for i := 1; i < height-1; i++ {
		result = append(result, frontLines[i]+borderStyle.Render("│"))
	}

	// Row height-1: front bottom line + base box bottom-right corner ╯
	if height > 1 {
		result = append(result, frontLines[height-1]+borderStyle.Render("╯"))
	}

	return strings.Join(result, "\n")
}

// compositeAttackOffsetCard overlaps frontCard (shifted 1 row UP relative to base box) while base box remains in standard active position
func compositeAttackOffsetCard(frontCard string, borderColor lipgloss.Color) string {
	frontLines := strings.Split(frontCard, "\n")
	height := len(frontLines)
	if height == 0 {
		return frontCard
	}

	cardWidth := lipgloss.Width(frontLines[0])
	borderStyle := lipgloss.NewStyle().Foreground(borderColor)

	var result []string
	// Line 0: top border of front card (shifted 1 row up)
	result = append(result, frontLines[0]+" ")

	// Line 1: front line 1 + base box top-right corner ╮
	if height > 1 {
		result = append(result, frontLines[1]+borderStyle.Render("╮"))
	}

	// Lines 2..height-2: front line + base box vertical right border │
	for i := 2; i < height-1; i++ {
		result = append(result, frontLines[i]+borderStyle.Render("│"))
	}

	// Line height-1: front bottom border + base box vertical right border │
	if height > 1 {
		result = append(result, frontLines[height-1]+borderStyle.Render("│"))
	}

	// Line height: empty space inside base box walls  │          │
	if cardWidth > 2 {
		emptyRow := borderStyle.Render("│") + strings.Repeat(" ", cardWidth-2) + borderStyle.Render("│")
		result = append(result, " "+emptyRow)
	}

	// Line height+1: base box bottom border  ╰──────────╯
	if cardWidth > 2 {
		bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", cardWidth-2) + "╯")
		result = append(result, " "+bottomBorder)
	}

	return strings.Join(result, "\n")
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

func (c HeroCardComponent) renderFullStatCard() string {
	rm := c.Responsive
	if rm == nil {
		rm = NewResponsiveManager(80)
	}

	innerWidth := 13
	paddingV := 0
	paddingH := 0
	if rm.IsAtLeast(BreakpointLG) {
		paddingV = 1
		paddingH = 2
	}

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(c.Theme.Border).
		Padding(paddingV, paddingH).
		Width(innerWidth + paddingH*2)

	if c.DataModel.IsActive {
		cardStyle = cardStyle.BorderForeground(c.Theme.Primary)
	} else if c.DataModel.IsDead {
		cardStyle = cardStyle.BorderForeground(c.Theme.Muted)
	}

	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(c.Theme.Primary)
	classStyle := lipgloss.NewStyle().Foreground(c.Theme.Secondary)
	statStyle := lipgloss.NewStyle().Foreground(c.Theme.Text)
	deadStyle := lipgloss.NewStyle().Foreground(c.Theme.Muted).Strikethrough(true)

	heroName := c.DataModel.Name
	if len(heroName) > 8 {
		heroName = heroName[:8]
	}

	var nameStr, levelStr string
	if c.DataModel.IsDead {
		nameStr = deadStyle.Render(heroName)
		levelStr = deadStyle.Render(fmt.Sprintf("%d", c.DataModel.Level))
	} else {
		nameStr = nameStyle.Render(heroName)
		levelStr = classStyle.Render(fmt.Sprintf("%d", c.DataModel.Level))
	}
	topLine := spaceBetween(nameStr, levelStr, innerWidth)

	classIcon := c.DataModel.Icon
	if classIcon == "" {
		classIcon = icon.GetClassIcon(c.DataModel.Class, icon.ModeNerdFont)
	}

	var classDisplay string
	if classIcon != "" {
		classDisplay = fmt.Sprintf("%s %s", classIcon, c.DataModel.Class)
	} else {
		classDisplay = c.DataModel.Class
	}

	var classStr string
	if c.DataModel.IsDead {
		classStr = deadStyle.Render(classDisplay)
	} else {
		classStr = classStyle.Render(classDisplay)
	}
	if len(classDisplay) < innerWidth {
		classStr = classStr + strings.Repeat(" ", innerWidth-len(classDisplay))
	}

	divider := lipgloss.NewStyle().Foreground(c.Theme.Border).Render(strings.Repeat("─", innerWidth))

	var statRow1, statRow2, statRow3, statRow4, statRow5 string
	var leftFooter, rightFooter string
	armIcon := icon.GetNerdFont(icon.KeyArmor)

	if c.DataModel.IsUnknownStats {
		statRow1 = statStyle.Render("str ?   dex ?")
		statRow2 = statStyle.Render("cha ?   sur ?")
		statRow3 = statStyle.Render("con ?   wil ?")
		statRow4 = statStyle.Render("int ?   per ?")
		statRow5 = statStyle.Render("hel ?   luc ?")
		if rm.IsAtLeast(BreakpointLG) {
			leftFooter = fmt.Sprintf("%s ? %s ?", classIcon, armIcon)
		} else {
			leftFooter = fmt.Sprintf("%s ?%s?", classIcon, armIcon)
		}
		rightFooter = "?/?"
	} else {
		st := c.DataModel.Stats
		statRow1 = statStyle.Render(fmt.Sprintf("str %-1d   dex %-1d", st.Str, st.Dex))
		statRow2 = statStyle.Render(fmt.Sprintf("cha %-1d   sur %-1d", st.Cha, st.Sur))
		statRow3 = statStyle.Render(fmt.Sprintf("con %-1d   wil %-1d", st.Con, st.Wil))
		statRow4 = statStyle.Render(fmt.Sprintf("int %-1d   per %-1d", st.Int, st.Per))
		statRow5 = statStyle.Render(fmt.Sprintf("hel %-1d   luc %-1d", st.Hel, st.Luc))
		if rm.IsAtLeast(BreakpointLG) {
			leftFooter = fmt.Sprintf("%s %d %s %d", classIcon, c.DataModel.Attack, armIcon, c.DataModel.Armor)
		} else {
			leftFooter = fmt.Sprintf("%s %d%s%d", classIcon, c.DataModel.Attack, armIcon, c.DataModel.Armor)
		}
		rightFooter = fmt.Sprintf("%d/%d", c.DataModel.HP, c.DataModel.MaxHP)
	}

	blankRow := strings.Repeat(" ", innerWidth)
	footerRow := spaceBetween(statStyle.Render(leftFooter), statStyle.Render(rightFooter), innerWidth)

	var sb strings.Builder
	sb.WriteString(topLine + "\n")
	sb.WriteString(classStr + "\n")
	sb.WriteString(divider + "\n")
	sb.WriteString(statRow1 + "\n")
	sb.WriteString(statRow2 + "\n")
	sb.WriteString(statRow3 + "\n")
	sb.WriteString(statRow4 + "\n")
	sb.WriteString(statRow5 + "\n")
	sb.WriteString(blankRow + "\n")
	sb.WriteString(footerRow)

	frontCard := cardStyle.Render(sb.String())

	var cardResult string
	if c.DataModel.IsActive {
		cardResult = compositeActiveOffsetCard2x1(frontCard, c.Theme.Border)
	} else {
		cardResult = "  \n" + frontCard
	}

	var statusIcon string
	if c.DataModel.Status != "" && c.DataModel.Status != "Healthy" {
		if c.DataModel.Status == "Sick" {
			statusIcon = icon.GetNerdFont(icon.KeySick)
		} else {
			statusIcon = c.DataModel.Status
		}
	}

	indicatorStyle := lipgloss.NewStyle().Foreground(c.Theme.Primary)
	if c.DataModel.IsActive && statusIcon != "" {
		activeIndicatorRow := "\n " + indicatorStyle.Render(statusIcon)
		cardResult += activeIndicatorRow
	} else {
		cardResult += "\n"
	}

	if c.ZoneID != "" {
		cardResult = mouse.Mark(c.ZoneID, cardResult)
	}

	return cardResult
}

// compositeActiveOffsetCard2x1 overlaps frontCard (shifted 1 up, 2 left) over a base border box (standard position)
func compositeActiveOffsetCard2x1(frontCard string, borderColor lipgloss.Color) string {
	frontLines := strings.Split(frontCard, "\n")
	height := len(frontLines)
	if height == 0 {
		return frontCard
	}

	cardWidth := lipgloss.Width(frontLines[0])
	borderStyle := lipgloss.NewStyle().Foreground(borderColor)

	var result []string
	// Line 0: top line of front card (raised 1 row up, 2 cols left)
	result = append(result, frontLines[0]+"  ")

	// Line 1: front line 1 + space + base box top-right corner ╮
	if height > 1 {
		result = append(result, frontLines[1]+" "+borderStyle.Render("╮"))
	}

	// Lines 2..height-2: front line + space + base box vertical right border │
	for i := 2; i < height-1; i++ {
		result = append(result, frontLines[i]+" "+borderStyle.Render("│"))
	}

	// Line height-1: front bottom line + space + base box vertical right border │
	if height > 2 {
		result = append(result, frontLines[height-1]+" "+borderStyle.Render("│"))
	}

	// Bottom Line (row height): 2 spaces + base box bottom border ╰──────────────╯
	if cardWidth > 2 {
		bottomBorder := borderStyle.Render("╰" + strings.Repeat("─", cardWidth-2) + "╯")
		result = append(result, "  "+bottomBorder)
	}

	return strings.Join(result, "\n")
}

func (c HeroCardComponent) renderEmptySlotCard() string {
	rm := c.Responsive
	if rm == nil {
		rm = NewResponsiveManager(80)
	}

	innerWidth := 13
	paddingV := 0
	paddingH := 0
	if rm.IsAtLeast(BreakpointLG) {
		paddingV = 1
		paddingH = 2
	}

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.Border{
			Top:         "┈",
			Bottom:      "┈",
			Left:        "┊",
			Right:       "┊",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomLeft:  "╰",
			BottomRight: "╯",
		}).
		BorderForeground(c.Theme.Border).
		Padding(paddingV, paddingH).
		Width(innerWidth + paddingH*2)

	if c.DataModel.IsActive {
		cardStyle = cardStyle.BorderForeground(c.Theme.Primary)
	}

	var btnBorder lipgloss.Style
	if c.DataModel.IsActive {
		btnBorder = lipgloss.NewStyle().Bold(true).Foreground(c.Theme.Primary)
	} else {
		btnBorder = lipgloss.NewStyle().Foreground(c.Theme.Muted)
	}

	bTop := btnBorder.Render("┌───────────┐")
	bMid1 := btnBorder.Render("│  Recruit  │")
	bMid2 := btnBorder.Render("│   Hero    │")
	bBot := btnBorder.Render("└───────────┘")

	blankRow := strings.Repeat(" ", innerWidth)

	var sb strings.Builder
	sb.WriteString(blankRow + "\n")
	sb.WriteString(blankRow + "\n")
	sb.WriteString(blankRow + "\n")
	sb.WriteString(bTop + "\n")
	sb.WriteString(bMid1 + "\n")
	sb.WriteString(bMid2 + "\n")
	sb.WriteString(bBot + "\n")
	sb.WriteString(blankRow + "\n")
	sb.WriteString(blankRow + "\n")
	sb.WriteString(blankRow)

	frontCard := cardStyle.Render(sb.String())

	var cardResult string
	if c.DataModel.IsActive {
		cardResult = compositeActiveOffsetCard2x1(frontCard, c.Theme.Border)
	} else {
		cardResult = "  \n" + frontCard
	}

	cardResult += "\n"
	return cardResult
}
