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
}
