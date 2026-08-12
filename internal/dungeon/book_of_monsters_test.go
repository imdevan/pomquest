package dungeon

import (
	"math/rand"
	"testing"
)

func TestGenerateMonstersForRoom(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	monsters := GenerateMonstersForRoom("Ancient Catacombs", 3, r, 0)

	if len(monsters) == 0 {
		t.Fatalf("expected at least 1 monster generated, got 0")
	}

	for _, m := range monsters {
		if m.Name == "" {
			t.Errorf("expected monster to have a non-empty name")
		}
		if m.MaxHP <= 0 {
			t.Errorf("expected monster HP > 0, got %d", m.MaxHP)
		}
	}
}

func TestGenerateMonstersForRoomMaxCR(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	for i := 0; i < 50; i++ {
		monsters := GenerateMonstersForRoom("Ancient Catacombs", 10, r, 1)
		for _, m := range monsters {
			for _, tmpl := range BookOfMonsters {
				if tmpl.Name == m.Name && tmpl.CR > 1 {
					t.Errorf("expected only CR 1 monsters when maxCR=1, got %s (CR %d)", m.Name, tmpl.CR)
				}
			}
		}
	}
}

func TestBookOfMonstersCRRules(t *testing.T) {
	for _, m := range BookOfMonsters {
		if m.CR == 1 {
			if m.HP < 1 || m.HP > 2 || m.MaxHP < 1 || m.MaxHP > 2 {
				t.Errorf("CR 1 monster %s expected 1-2 HP, got HP=%d MaxHP=%d", m.Name, m.HP, m.MaxHP)
			}
			if m.XPValue != 1 {
				t.Errorf("CR 1 monster %s expected 1 XP, got %d", m.Name, m.XPValue)
			}
			if m.LootGold != 0 {
				t.Errorf("CR 1 monster %s expected 0 gold, got %d", m.Name, m.LootGold)
			}
		} else if m.CR < 4 {
			t.Errorf("monster %s has CR %d, expected CR >= 4 or CR == 1", m.Name, m.CR)
		}
	}
}
