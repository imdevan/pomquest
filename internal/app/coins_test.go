package app

import "testing"

func TestCoinConversions(t *testing.T) {
	// 1g = 10s = 1000c
	p := Party{Coins: 1250}
	total := TotalCoins(p)
	if total != 1250 {
		t.Fatalf("expected 1250 total coins, got %d", total)
	}

	g, s, c := BreakdownCoins(total)
	if g != 1 || s != 2 || c != 50 {
		t.Fatalf("expected 1g 2s 50c, got %dg %ds %dc", g, s, c)
	}

	formatted := FormatCoins(total)
	if formatted != "1g 2s 50c" {
		t.Fatalf("expected '1g 2s 50c', got %q", formatted)
	}

	// Test silver-only formatting
	formattedSilver := FormatCoins(350)
	if formattedSilver != "3s 50c" {
		t.Fatalf("expected '3s 50c', got %q", formattedSilver)
	}

	// Test copper-only formatting
	formattedCopper := FormatCoins(75)
	if formattedCopper != "75c" {
		t.Fatalf("expected '75c', got %q", formattedCopper)
	}
}
