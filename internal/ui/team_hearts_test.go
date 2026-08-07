package ui

import (
	"testing"

	"github.com/pomquest/internal/adapters/icon"
)

func TestTeamHeartsComponent(t *testing.T) {
	tests := []struct {
		name      string
		currentHP int
		maxHP     int
		mode      icon.Mode
		expected  string
	}{
		{
			name:      "32 max hp 24 current hp",
			currentHP: 24,
			maxHP:     32,
			mode:      icon.ModeNerdFont,
			expected:  "󰋑 󰋑 󰛞 󰛠",
		},
		{
			name:      "35 max hp full 35 current hp",
			currentHP: 35,
			maxHP:     35,
			mode:      icon.ModeNerdFont,
			expected:  "󰋑 󰋑 󰋑 󰛟",
		},
		{
			name:      "35 max hp 30 current hp",
			currentHP: 30,
			maxHP:     35,
			mode:      icon.ModeNerdFont,
			expected:  "󰋑 󰋑 󰋑 󰛠",
		},
		{
			name:      "35 max hp 25 current hp",
			currentHP: 25,
			maxHP:     35,
			mode:      icon.ModeNerdFont,
			expected:  "󰋑 󰋑 󰛞 󰛠",
		},
		{
			name:      "40 max hp 40 current hp",
			currentHP: 40,
			maxHP:     40,
			mode:      icon.ModeNerdFont,
			expected:  "󰋑 󰋑 󰋑 󰋑",
		},
		{
			name:      "40 max hp 35 current hp",
			currentHP: 35,
			maxHP:     40,
			mode:      icon.ModeNerdFont,
			expected:  "󰋑 󰋑 󰋑 󰛞",
		},
		{
			name:      "0 max hp",
			currentHP: 0,
			maxHP:     0,
			mode:      icon.ModeNerdFont,
			expected:  "",
		},
		{
			name:      "emoji mode 35/35",
			currentHP: 35,
			maxHP:     35,
			mode:      icon.ModeEmoji,
			expected:  "❤️ ❤️ ❤️ 💗",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderTeamHearts(tt.currentHP, tt.maxHP, tt.mode)
			if got != tt.expected {
				t.Errorf("RenderTeamHearts(%d, %d) = %q, expected %q", tt.currentHP, tt.maxHP, got, tt.expected)
			}
		})
	}
}
