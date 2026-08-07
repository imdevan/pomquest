package icon

import "testing"

func TestIconRegistry(t *testing.T) {
	atkNerd := GetNerdFont(KeyAttack)
	if atkNerd != "󰓥" {
		t.Errorf("expected KeyAttack NerdFont to be 󰓥, got %s", atkNerd)
	}

	atkEmoji := GetEmoji(KeyAttack)
	if atkEmoji != "⚔️" {
		t.Errorf("expected KeyAttack Emoji to be ⚔️, got %s", atkEmoji)
	}

	weather := Get(KeyWeather, ModeNerdFont)
	if weather != "☀️" {
		t.Errorf("expected KeyWeather to be ☀️, got %s", weather)
	}

	// Test party stat static icons
	if got := GetNerdFont(KeyFood); got != "󾠅" {
		t.Errorf("expected KeyFood NerdFont to be 󾠅, got %s", got)
	}
	if got := GetNerdFont(KeyOx); got != "󰄛" {
		t.Errorf("expected KeyOx NerdFont to be 󰄛, got %s", got)
	}
	if got := GetNerdFont(KeySpareParts); got != "󰏖" {
		t.Errorf("expected KeySpareParts NerdFont to be 󰏖, got %s", got)
	}
	if got := GetNerdFont(KeyLoot); got != "󱒔" {
		t.Errorf("expected KeyLoot NerdFont to be 󱒔, got %s", got)
	}
	if got := GetNerdFont(KeyRope); got != "󰈤" {
		t.Errorf("expected KeyRope NerdFont to be 󰈤, got %s", got)
	}
	if got := GetNerdFont(KeyBomb); got != "󰄬" {
		t.Errorf("expected KeyBomb NerdFont to be 󰄬, got %s", got)
	}
	if got := GetNerdFont(KeyPotion); got != "󰅠" {
		t.Errorf("expected KeyPotion NerdFont to be 󰅠, got %s", got)
	}
	if got := GetNerdFont(KeyShovel); got != "󱇨" {
		t.Errorf("expected KeyShovel NerdFont to be 󱇨, got %s", got)
	}

	// Test emoji mode variants
	if got := GetEmoji(KeyLoot); got != "💰" {
		t.Errorf("expected KeyLoot Emoji to be 💰, got %s", got)
	}
	if got := GetEmoji(KeyRope); got != "🪢" {
		t.Errorf("expected KeyRope Emoji to be 🪢, got %s", got)
	}
}

func TestPartyStatValueBasedIcons(t *testing.T) {
	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{"item usage none", GetItemUsageIcon("none", ModeNerdFont), "󰅖"},
		{"item usage sparse", GetItemUsageIcon("sparse", ModeNerdFont), "󰓎"},
		{"item usage balanced", GetItemUsageIcon("balanced", ModeNerdFont), "󰓍"},
		{"item usage heavy", GetItemUsageIcon("heavy", ModeNerdFont), "󰓌"},

		{"diet none", GetDietIcon("none", ModeNerdFont), "󰅖"},
		{"diet sparse", GetDietIcon("sparse", ModeNerdFont), "󰓎"},
		{"diet balanced", GetDietIcon("balanced", ModeNerdFont), "󰓍"},
		{"diet filling", GetDietIcon("filling", ModeNerdFont), "󰓌"},

		{"greed none", GetGreedIcon("none", ModeNerdFont), "󰅖"},
		{"greed sparse", GetGreedIcon("sparse", ModeNerdFont), "󰓎"},
		{"greed balanced", GetGreedIcon("balanced", ModeNerdFont), "󰓍"},
		{"greed heavy", GetGreedIcon("heavy", ModeNerdFont), "󰓌"},

		{"pace slow", GetPaceIcon("slow", ModeNerdFont), "󰓎"},
		{"pace balanced", GetPaceIcon("balanced", ModeNerdFont), "󰓍"},
		{"pace fast", GetPaceIcon("fast", ModeNerdFont), "󰓌"},

		// Emoji mode tests for value based icons
		{"item usage heavy emoji", GetItemUsageIcon("heavy", ModeEmoji), "🔥"},
		{"diet balanced emoji", GetDietIcon("balanced", ModeEmoji), "🍲"},
		{"greed sparse emoji", GetGreedIcon("sparse", ModeEmoji), "🪙"},
		{"pace fast emoji", GetPaceIcon("fast", ModeEmoji), "🏃"},
	}

	for _, tt := range tests {
		if tt.got != tt.expected {
			t.Errorf("test %s failed: expected %s, got %s", tt.name, tt.expected, tt.got)
		}
	}
}
