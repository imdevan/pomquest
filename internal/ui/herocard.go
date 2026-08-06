package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// HeroCardData holds values for a single hero card display.
type HeroCardData struct {
	Name   string
	Class  string
	Level  int
	HP     int
	MaxHP  int
	Attack int
	Armor  int
	Status string
	Icon   string
	IsDead bool
}

// HeroCardComponent is a generic hero card widget.
type HeroCardComponent struct {
	Data  HeroCardData
	Theme Theme
}

func NewHeroCardComponent(theme Theme, data HeroCardData) HeroCardComponent {
	return HeroCardComponent{
		Data:  data,
		Theme: theme,
	}
}

func (c HeroCardComponent) View() string {
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(c.Theme.Border).
		Padding(0, 1).
		Width(18)

	if c.Data.IsDead {
		cardStyle = cardStyle.BorderForeground(c.Theme.Muted)
	}

	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(c.Theme.Primary)
	classStyle := lipgloss.NewStyle().Italic(true).Foreground(c.Theme.Secondary)
	statStyle := lipgloss.NewStyle().Foreground(c.Theme.Text)
	statusStyle := lipgloss.NewStyle().Foreground(c.Theme.Tags)

	nameStr := c.Data.Name
	if c.Data.IsDead {
		nameStr = fmt.Sprintf("~%s~ (ded)", c.Data.Name)
	}

	var sb strings.Builder
	sb.WriteString(nameStyle.Render(nameStr) + "\n")
	sb.WriteString(classStyle.Render(fmt.Sprintf("%s Lvl %d", c.Data.Class, c.Data.Level)) + "\n")

	if c.Data.Status != "" && c.Data.Status != "Healthy" {
		sb.WriteString(statusStyle.Render("󰗑 "+c.Data.Status) + "\n")
	}

	sb.WriteString(statStyle.Render(fmt.Sprintf("⚔%d 🛡%d", c.Data.Attack, c.Data.Armor)) + "\n")
	sb.WriteString(statStyle.Render(fmt.Sprintf(" %d/%d HP", c.Data.HP, c.Data.MaxHP)))

	return cardStyle.Render(sb.String())
}
