package app

import "fmt"

// TotalCoins returns the total coin value in copper for a given Party.
// Conversion rate: 1g = 10s = 1,000c; 1s = 100c.
func TotalCoins(party Party) int {
	return party.Gold*1000 + party.Silver*100 + party.Copper
}

// BreakdownCoins converts a total coin value into gold (1000c), silver (100c), and copper (1c).
func BreakdownCoins(totalCoins int) (gold, silver, copper int) {
	if totalCoins < 0 {
		totalCoins = 0
	}
	gold = totalCoins / 1000
	silver = (totalCoins % 1000) / 100
	copper = totalCoins % 100
	return gold, silver, copper
}

// SetPartyCoins updates party Gold, Silver, and Copper fields from total coins.
func SetPartyCoins(party *Party, totalCoins int) {
	g, s, c := BreakdownCoins(totalCoins)
	party.Gold = g
	party.Silver = s
	party.Copper = c
}

// FormatCoins returns a formatted string representation of total coins (e.g. "1g 2s 50c").
func FormatCoins(totalCoins int) string {
	g, s, c := BreakdownCoins(totalCoins)
	if g > 0 {
		return fmt.Sprintf("%dg %ds %dc", g, s, c)
	}
	if s > 0 {
		return fmt.Sprintf("%ds %dc", s, c)
	}
	return fmt.Sprintf("%dc", c)
}
