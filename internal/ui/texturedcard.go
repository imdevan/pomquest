package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// CornerText specifies text content, padding, margins, and style for a card border corner.
type CornerText struct {
	Text         string
	PaddingLeft  int
	PaddingRight int
	Margin       int
	Style        lipgloss.Style
}

// TexturedCard is a custom border card container supporting embedded text on any of the 4 corners and inner padding.
type TexturedCard struct {
	Width         int
	Height        int
	PaddingTop    int
	PaddingRight  int
	PaddingBottom int
	PaddingLeft   int
	Border        lipgloss.Border
	BorderStyle   lipgloss.Style
	Content       string
	TopLeft       CornerText
	TopRight      CornerText
	BottomLeft    CornerText
	BottomRight   CornerText
	Theme         Theme
}

// NewTexturedCard creates a TexturedCard with default rounded border and theme colors.
func NewTexturedCard(theme Theme, width int, content string) TexturedCard {
	if width <= 0 {
		width = 72
	}
	return TexturedCard{
		Width:       width,
		Border:      lipgloss.RoundedBorder(),
		BorderStyle: lipgloss.NewStyle().Foreground(theme.Border),
		Content:     content,
		Theme:       theme,
	}
}

// Padding sets inner padding values (top, right, bottom, left) on the card container matching Lip Gloss syntax.
func (c TexturedCard) Padding(args ...int) TexturedCard {
	switch len(args) {
	case 1:
		c.PaddingTop = args[0]
		c.PaddingRight = args[0]
		c.PaddingBottom = args[0]
		c.PaddingLeft = args[0]
	case 2:
		c.PaddingTop = args[0]
		c.PaddingBottom = args[0]
		c.PaddingRight = args[1]
		c.PaddingLeft = args[1]
	case 4:
		c.PaddingTop = args[0]
		c.PaddingRight = args[1]
		c.PaddingBottom = args[2]
		c.PaddingLeft = args[3]
	}
	return c
}

func (c TexturedCard) View() string {
	width := c.Width
	if width <= 4 {
		width = 72
	}

	borderStyle := c.BorderStyle
	if borderStyle.GetForeground() == nil {
		borderStyle = lipgloss.NewStyle().Foreground(c.Theme.Border)
	}

	topLine := c.renderBorderLine(
		c.Border.TopLeft,
		c.Border.TopRight,
		c.Border.Top,
		c.TopLeft,
		c.TopRight,
		width,
		borderStyle,
	)

	bottomLine := c.renderBorderLine(
		c.Border.BottomLeft,
		c.Border.BottomRight,
		c.Border.Bottom,
		c.BottomLeft,
		c.BottomRight,
		width,
		borderStyle,
	)

	leftBorder := borderStyle.Render(c.Border.Left)
	rightBorder := borderStyle.Render(c.Border.Right)

	leftPadStr := strings.Repeat(" ", c.PaddingLeft)
	rightPadStr := strings.Repeat(" ", c.PaddingRight)
	innerAvailableWidth := width - 2 - c.PaddingLeft - c.PaddingRight
	if innerAvailableWidth < 0 {
		innerAvailableWidth = 0
	}

	emptyPaddingLine := leftBorder + strings.Repeat(" ", width-2) + rightBorder

	var middleLines []string

	// Padding top lines
	for i := 0; i < c.PaddingTop; i++ {
		middleLines = append(middleLines, emptyPaddingLine)
	}

	trimmedContent := strings.TrimSpace(c.Content)
	contentLines := strings.Split(trimmedContent, "\n")
	for _, rawLine := range contentLines {
		rawLine = strings.TrimRight(rawLine, "\r")
		if strings.TrimSpace(rawLine) == "" {
			continue
		}
		linesToRender := []string{rawLine}
		if innerAvailableWidth > 0 && lipgloss.Width(rawLine) > innerAvailableWidth {
			linesToRender = WrapLine(rawLine, innerAvailableWidth)
		}

		for _, line := range linesToRender {
			lineWidth := lipgloss.Width(line)
			neededSpaces := innerAvailableWidth - lineWidth
			if neededSpaces < 0 {
				neededSpaces = 0
			}
			formattedLine := leftBorder + leftPadStr + line + strings.Repeat(" ", neededSpaces) + rightPadStr + rightBorder
			middleLines = append(middleLines, formattedLine)
		}
	}

	// Padding bottom lines
	for i := 0; i < c.PaddingBottom; i++ {
		middleLines = append(middleLines, emptyPaddingLine)
	}

	var result []string
	result = append(result, topLine)
	result = append(result, middleLines...)
	result = append(result, bottomLine)

	return strings.Join(result, "\n")
}

func (c TexturedCard) renderBorderLine(
	leftCornerChar, rightCornerChar, dashChar string,
	leftCorner, rightCorner CornerText,
	totalWidth int,
	borderStyle lipgloss.Style,
) string {
	innerWidth := totalWidth - 2
	if innerWidth < 0 {
		innerWidth = 0
	}

	leftSeg, leftLen := c.renderCornerSegment(leftCorner, true, dashChar, borderStyle)
	rightSeg, rightLen := c.renderCornerSegment(rightCorner, false, dashChar, borderStyle)

	middleDashesLen := innerWidth - leftLen - rightLen
	if middleDashesLen < 0 {
		middleDashesLen = 0
	}

	middleDashes := borderStyle.Render(strings.Repeat(dashChar, middleDashesLen))

	return borderStyle.Render(leftCornerChar) + leftSeg + middleDashes + rightSeg + borderStyle.Render(rightCornerChar)
}

func (c TexturedCard) renderCornerSegment(
	ct CornerText,
	isLeft bool,
	dashChar string,
	borderStyle lipgloss.Style,
) (string, int) {
	if ct.Text == "" {
		return "", 0
	}

	textStyle := ct.Style
	if textStyle.GetForeground() == nil {
		textStyle = lipgloss.NewStyle().Bold(true).Foreground(c.Theme.Primary)
	}

	textRendered := textStyle.Render(ct.Text)
	textLen := lipgloss.Width(ct.Text)

	marginDashes := borderStyle.Render(strings.Repeat(dashChar, ct.Margin))
	leftPadding := strings.Repeat(" ", ct.PaddingLeft)
	rightPadding := strings.Repeat(" ", ct.PaddingRight)

	totalLen := ct.Margin + ct.PaddingLeft + textLen + ct.PaddingRight

	if isLeft {
		// Left corner: [Margin dashes][PaddingLeft][Text][PaddingRight]
		return marginDashes + leftPadding + textRendered + rightPadding, totalLen
	}

	// Right corner: [PaddingLeft][Text][PaddingRight][Margin dashes]
	return leftPadding + textRendered + rightPadding + marginDashes, totalLen
}
