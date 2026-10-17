package provisions

import (
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/pomquest/internal/app"
	"github.com/pomquest/internal/domain"
	"github.com/pomquest/internal/sacred_texts"
)

func TestBookOfItems(t *testing.T) {
	catalog := sacred_texts.BookOfItems()

	if len(catalog) != 7 {
		t.Fatalf("expected 7 items in catalog, got %d", len(catalog))
	}

	expectedPrices := map[string]string{
		"food":        "10c / lb",
		"spare_parts": "60c / part",
		"ox":          "20s / yoke (3 ox)",
		"potions":     "20s / pot",
		"bombs":       "50c / bomb",
		"shovels":     "20c / shovel",
		"rope":        "10c / 50ft",
	}

	for _, item := range catalog {
		price, ok := expectedPrices[item.ID]
		if !ok {
			t.Errorf("unexpected item ID in catalog: %s", item.ID)
		} else if item.DisplayCost != price {
			t.Errorf("expected price %s for %s, got %s", price, item.ID, item.DisplayCost)
		}
	}
}

func TestProvisionsScreen_ViewAndColumns(t *testing.T) {
	cfg := domain.DefaultConfig()
	m := NewModel(cfg)

	view := m.View()

	// 24.2 & 24.9: borderless store view with title and table columns
	if !strings.Contains(view, "Soand Sow's General Store") {
		t.Errorf("expected store view to contain header name, got:\n%s", view)
	}

	// 24.6: 3 cols: Item, Cost, Subtotal
	if !strings.Contains(view, "Item") || !strings.Contains(view, "Cost") || !strings.Contains(view, "Subtotal") {
		t.Errorf("expected 3 table columns (Item, Cost, Subtotal), got:\n%s", view)
	}

	// 24.8: Last option is Complete Purchase
	if !strings.Contains(view, "Complete Purchase") {
		t.Errorf("expected last option Complete Purchase, got:\n%s", view)
	}

	// 24.4: Party card rendered under store
	if !strings.Contains(view, "Party") {
		t.Errorf("expected party card under store view, got:\n%s", view)
	}
}

func TestProvisionsScreen_DeferredPurchaseExecution(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "provisions_test.db")

	db, err := app.InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	// Give party 1 Gold (1000c)
	party, _ := app.LoadOrCreateActiveParty(db)
	party.Coins = 1000
	_ = app.SaveParty(db, party)

	cfg := domain.DefaultConfig()
	m := NewModelWithDBPath(cfg, dbPath)

	// Increment Food Qty by pressing '+' twice (2 lbs = 20c)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'+'}})
	m = updated.(Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'+'}})
	m = updated.(Model)

	// 24.8: Party coins should NOT change before purchase
	if m.Party.Coins != 1000 {
		t.Fatalf("expected coins unchanged before purchase, got %d coins", m.Party.Coins)
	}

	// Subtotal should reflect 20c (2 * 10c)
	if m.calculateGrandTotal() != 20 {
		t.Fatalf("expected cart grand total 20c, got %d", m.calculateGrandTotal())
	}

	// Navigate to Purchase row (row 7) and press Enter
	for i := 0; i < 7; i++ {
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		m = updated.(Model)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)

	// 1000c - 20c = 980c -> 9s 80c
	if m.Party.Coins != 980 {
		t.Errorf("expected remaining purse 980 coins, got %d", m.Party.Coins)
	}
	if m.Party.FoodLbs != 102 { // 100 + 2
		t.Errorf("expected FoodLbs to be 102, got %d", m.Party.FoodLbs)
	}
	if !strings.Contains(m.StatusNotice, "Purchase complete") {
		t.Errorf("expected status notice for purchase, got %q", m.StatusNotice)
	}
}

func TestProvisionsScreen_PartyCardHeartsFromSavedHeroes(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "provisions_hearts_test.db")

	db, err := app.InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	party, _ := app.LoadOrCreateActiveParty(db)
	party.Coins = 2000
	_ = app.SaveParty(db, party)

	// Add 2 heroes with specific HP values
	hero1 := app.Hero{
		Creature: app.Creature{
			Noun:  app.Noun{Name: "Gimli"},
			HP:    20,
			MaxHP: 20,
		},
		Class: "fighter",
	}
	hero2 := app.Hero{
		Creature: app.Creature{
			Noun:  app.Noun{Name: "Legolas"},
			HP:    10,
			MaxHP: 15,
		},
		Class: "ranger",
	}
	_ = app.ReplaceOrAddHeroInParty(db, 0, hero1)
	_ = app.ReplaceOrAddHeroInParty(db, 1, hero2)

	cfg := domain.DefaultConfig()
	m := NewModelWithDBPath(cfg, dbPath)

	view := m.View()
	// Party Card should be rendered
	if !strings.Contains(view, "Party") {
		t.Errorf("expected Party card in view, got:\n%s", view)
	}
}

