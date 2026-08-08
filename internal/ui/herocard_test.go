package ui

import (
	"strings"
	"testing"

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

	// On LG+: expect space between icon and attack value ("⚔ 2 🛡 2")
	if !strings.Contains(view, "2 🛡 2") {
		t.Errorf("expected LG hero card view to contain '2 🛡 2' with spaces, got:\n%s", view)
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

	// On MD: expect no space between icon and attack/armor value ("⚔2 🛡2")
	if !strings.Contains(view, "2 🛡2") {
		t.Errorf("expected MD hero card view to contain '2 🛡2' without space between icon & val, got:\n%s", view)
	}
	if strings.Contains(view, "2 🛡 2") {
		t.Errorf("expected MD hero card view NOT to contain '2 🛡 2', got:\n%s", view)
	}
}
