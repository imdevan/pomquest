package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/pomquest/internal/domain"
)

func TestHeroCardComponentView_LG(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ThemeFromConfig(cfg)

	hero := HeroCardData{
		Name:   "jimmothy",
		Class:  "fighter",
		Level:  2,
		HP:     10,
		MaxHP:  10,
		Attack: 2,
		Armor:  2,
		Status: "Healthy",
	}

	rmLG := NewResponsiveManager(96)
	comp := NewHeroCardComponent(theme, hero, rmLG)
	view := comp.View()

	// On LG+: expect space between icon and attack value ("⚔ 2 󰒘 2")
	if !strings.Contains(view, "2 󰒘 2") {
		t.Errorf("expected LG hero card view to contain '2 󰒘 2' with spaces, got:\n%s", view)
	}
	// Expect class icon next to class name ("⚔ fighter")
	if !strings.Contains(view, "⚔ fighter") {
		t.Errorf("expected hero card view to contain class icon next to class '⚔ fighter', got:\n%s", view)
	}
}

func TestHeroCardComponentView_MD(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ThemeFromConfig(cfg)

	hero := HeroCardData{
		Name:   "jimmothy",
		Class:  "fighter",
		Level:  2,
		HP:     10,
		MaxHP:  10,
		Attack: 2,
		Armor:  2,
		Status: "Healthy",
	}

	rmMD := NewResponsiveManager(72)
	comp := NewHeroCardComponent(theme, hero, rmMD)
	view := comp.View()

	// On MD: expect no space between icon and attack/armor value ("󰓥2󰒘2")
	if !strings.Contains(view, "󰓥2󰒘2") {
		t.Errorf("expected MD hero card view to contain '󰓥2󰒘2' without spaces, got:\n%s", view)
	}
	if strings.Contains(view, "2 󰒘 2") {
		t.Errorf("expected MD hero card view NOT to contain '2 󰒘 2', got:\n%s", view)
	}

	// On MD: expect abbreviated name ("jim.")
	if !strings.Contains(view, "jim.") {
		t.Errorf("expected MD hero card view to contain abbreviated name 'jim.', got:\n%s", view)
	}

	// On MD: expect class icon next to class name ("⚔ fighter")
	if !strings.Contains(view, "⚔ fighter") {
		t.Errorf("expected MD hero card view to contain '⚔ fighter', got:\n%s", view)
	}
}

func TestHeroCardsComponent_Fit58Columns(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ThemeFromConfig(cfg)

	heroes := []HeroCardData{
		{Name: "jimmothy", Class: "fighter", Level: 2, HP: 10, MaxHP: 10, Attack: 2, Armor: 2, Status: "Healthy", IsActive: false},
		{Name: "applebee", Class: "wizard", Level: 1, HP: 6, MaxHP: 6, Attack: 3, Armor: 0, Status: "Healthy"},
		{Name: "pam", Class: "ranger", Level: 1, HP: 8, MaxHP: 8, Attack: 3, Armor: 0, Status: "Healthy"},
		{Name: "sam", Class: "rouge", Level: 2, HP: 0, MaxHP: 8, Attack: 3, Armor: 0, Status: "Dead", IsDead: true},
	}

	rmMD := NewResponsiveManager(58)
	var cards []string
	for _, h := range heroes {
		cardView := NewHeroCardComponent(theme, h, rmMD).View()
		cardView = strings.TrimPrefix(cardView, "\n")
		cards = append(cards, cardView)
	}

	var spacerCards []string
	for i, card := range cards {
		spacerCards = append(spacerCards, card)
		if i < len(cards)-1 {
			spacerCards = append(spacerCards, "  ")
		}
	}

	grid := lipgloss.JoinHorizontal(lipgloss.Top, spacerCards...)
	view := lipgloss.NewStyle().Width(58).Align(lipgloss.Center).Render(grid)

	lines := strings.Split(view, "\n")
	for i, line := range lines {
		w := lipgloss.Width(line)
		if w > 58 {
			t.Errorf("line %d width %d exceeds 58 col max limit on MD:\n%q", i, w, line)
		}
	}
}

func TestHeroCardComponent_FullStats(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ThemeFromConfig(cfg)

	hero := HeroCardData{
		Name:          "pam",
		Class:         "ranger",
		Level:         1,
		HP:            6,
		MaxHP:         8,
		Attack:        3,
		Armor:         0,
		Status:        "Healthy",
		IsActive:      true,
		ShowFullStats: true,
	}

	comp := NewHeroCardComponent(theme, hero)
	view := comp.View()

	// Verify all 10 stats are present in the full stat card view
	stats := []string{"str", "dex", "cha", "sur", "con", "wil", "int", "per", "hel", "luc"}
	for _, s := range stats {
		if !strings.Contains(view, s) {
			t.Errorf("expected full stat hero card to contain %q, got:\n%s", s, view)
		}
	}

	// Verify active offset (top right corner '╮' from old box underneath)
	if !strings.Contains(view, "╮") {
		t.Errorf("expected active full stat hero card to contain raised border '╮', got:\n%s", view)
	}

	// Verify class icon is shown next to class name ("󱡁 ranger")
	if !strings.Contains(view, "󱡁 ranger") {
		t.Errorf("expected full stat card to show class icon next to class '󱡁 ranger', got:\n%s", view)
	}
}

func TestHeroCardComponent_LG_Padding(t *testing.T) {
	cfg := domain.DefaultConfig()
	theme := ThemeFromConfig(cfg)

	hero := HeroCardData{
		Name:          "pam",
		Class:         "ranger",
		Level:         1,
		HP:            6,
		MaxHP:         8,
		Attack:        3,
		Armor:         0,
		Status:        "Healthy",
		ShowFullStats: true,
	}

	rmLG := NewResponsiveManager(80) // LG screen (80 cols)
	comp := NewHeroCardComponent(theme, hero, rmLG)
	view := comp.View()

	lines := strings.Split(view, "\n")
	// Full stat card with vertical padding 1 has top blank line inside border box
	if len(lines) < 12 {
		t.Errorf("expected LG hero card view to have vertical padding lines, got %d lines:\n%s", len(lines), view)
	}

	// Verify space between attack icon/number and armor icon/number on LG+ screens
	if !strings.Contains(view, "3 󰒘 0") {
		t.Errorf("expected LG full stat hero card to contain '3 󰒘 0' with spaces between icons and numbers, got:\n%s", view)
	}
}



