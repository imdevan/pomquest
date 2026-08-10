package provisions

import (
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/pomquest/internal/app"
	"github.com/pomquest/internal/domain"
)

func TestBookOfItems(t *testing.T) {
	catalog := app.BookOfItems()

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

func TestProvisionsScreen_ViewAndProps(t *testing.T) {
	cfg := domain.DefaultConfig()
	m := NewModel(cfg)

	view := m.View()

	// 24.2: props name: Soand Sow's, type: general, sells: all items at base cost
	if !strings.Contains(view, "Soand Sow's General Store") {
		t.Errorf("expected store view to contain name 'Soand Sow's General Store', got:\n%s", view)
	}
	if !strings.Contains(view, "General Store • All items sold at base cost") {
		t.Errorf("expected store view to contain subtitle 'General Store • All items sold at base cost', got:\n%s", view)
	}
	if !strings.Contains(view, "Food") || !strings.Contains(view, "Spare Parts") {
		t.Errorf("expected store view to list catalog items, got:\n%s", view)
	}
}

func TestProvisionsScreen_PurchaseItem(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "provisions_test.db")

	db, err := app.InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	// Give party 1 Gold (100c)
	party, _ := app.LoadOrCreateActiveParty(db)
	party.Gold = 1
	party.Silver = 0
	party.Copper = 0
	_ = app.SaveParty(db, party)

	cfg := domain.DefaultConfig()
	m := NewModelWithDBPath(cfg, dbPath)

	// Buy item 1 (Food - 10c / lb)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	m = updated.(Model)

	if m.Party.FoodLbs != 101 { // starting 100 + 1 bought
		t.Errorf("expected FoodLbs to be 101, got %d", m.Party.FoodLbs)
	}
	// 100c - 10c = 90c (0g 9s 0c)
	if m.Party.Gold != 0 || m.Party.Silver != 9 || m.Party.Copper != 0 {
		t.Errorf("expected remaining purse to be 0g 9s 0c, got %dg %ds %dc", m.Party.Gold, m.Party.Silver, m.Party.Copper)
	}
	if !strings.Contains(m.StatusNotice, "Purchased Food") {
		t.Errorf("expected status notice for purchase, got %q", m.StatusNotice)
	}

	// Verify database was updated
	reloadedParty, _ := app.LoadActivePartyFromPath(dbPath)
	if reloadedParty.FoodLbs != 101 {
		t.Errorf("expected reloaded party FoodLbs to be 101, got %d", reloadedParty.FoodLbs)
	}
}

func TestProvisionsScreen_InsufficientCoin(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "insufficient_test.db")

	db, err := app.InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	// Give party 0 coin
	party, _ := app.LoadOrCreateActiveParty(db)
	party.Gold = 0
	party.Silver = 0
	party.Copper = 0
	_ = app.SaveParty(db, party)

	cfg := domain.DefaultConfig()
	m := NewModelWithDBPath(cfg, dbPath)

	// Try buying Food (10c)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	m = updated.(Model)

	if !strings.Contains(m.StatusNotice, "Not enough coin") {
		t.Errorf("expected warning for insufficient coin, got %q", m.StatusNotice)
	}
}
