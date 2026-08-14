package dungeon

import (
	"math/rand"
	"testing"

	"github.com/pomquest/internal/app"
)

func TestBuildTurn(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	heroes := []app.Hero{
		{Creature: app.Creature{Noun: app.Noun{Name: "Hero Fast"}, HP: 10, MaxHP: 10}, Stats: app.Stats{Dex: 5}},
		{Creature: app.Creature{Noun: app.Noun{Name: "Hero Slow"}, HP: 10, MaxHP: 10}, Stats: app.Stats{Dex: 1}},
		{Creature: app.Creature{Noun: app.Noun{Name: "Hero Dead"}, HP: 0, MaxHP: 10}, IsDead: true, Stats: app.Stats{Dex: 10}},
	}

	monsters := []app.Monster{
		{Creature: app.Creature{Noun: app.Noun{Name: "Monster Mid"}, HP: 8, MaxHP: 8}, Stats: app.Stats{Dex: 3}},
		{Creature: app.Creature{Noun: app.Noun{Name: "Monster Dead"}, HP: 0, MaxHP: 5}, Stats: app.Stats{Dex: 6}},
	}

	queue := BuildTurn(heroes, monsters, rng)

	if len(queue) != 3 {
		t.Fatalf("expected 3 living combatants in queue, got %d", len(queue))
	}

	// Order should be Hero Fast (Dex 5) -> Monster Mid (Dex 3) -> Hero Slow (Dex 1)
	if queue[0].Name != "Hero Fast" {
		t.Errorf("expected first combatant to be 'Hero Fast', got %q", queue[0].Name)
	}
	if queue[1].Name != "Monster Mid" {
		t.Errorf("expected second combatant to be 'Monster Mid', got %q", queue[1].Name)
	}
	if queue[2].Name != "Hero Slow" {
		t.Errorf("expected third combatant to be 'Hero Slow', got %q", queue[2].Name)
	}
}

func TestExecuteCombatTurnCycle_HeroAttack(t *testing.T) {
	rng := rand.New(rand.NewSource(1))

	heroes := []app.Hero{
		{Creature: app.Creature{Noun: app.Noun{Name: "Fighter"}, HP: 10, MaxHP: 10, Attack: 3}, Stats: app.Stats{Dex: 4}},
	}

	monsters := []app.Monster{
		{Creature: app.Creature{Noun: app.Noun{Name: "Goblin"}, HP: 5, MaxHP: 5, Armor: 1}, Stats: app.Stats{Dex: 2}},
	}

	queue := BuildTurn(heroes, monsters, rng)
	currentIdx := 0

	res := ExecuteCombatTurnCycle(&queue, &currentIdx, heroes, &monsters, rng)

	if res.ActiveHeroIdx != 0 {
		t.Errorf("expected ActiveHeroIdx to be 0, got %d", res.ActiveHeroIdx)
	}
	if len(res.ActionLogs) == 0 {
		t.Errorf("expected at least 1 action log message")
	}
}

func TestExecuteCombatTurnCycle_SlayMonster(t *testing.T) {
	rng := rand.New(rand.NewSource(2))

	heroes := []app.Hero{
		{Creature: app.Creature{Noun: app.Noun{Name: "Hero"}, HP: 10, MaxHP: 10, Attack: 100}, Stats: app.Stats{Dex: 10}},
	}

	monsters := []app.Monster{
		{Creature: app.Creature{Noun: app.Noun{Name: "Weak Goblin"}, HP: 1, MaxHP: 1}, Stats: app.Stats{Dex: 1}},
	}

	queue := BuildTurn(heroes, monsters, rng)
	currentIdx := 0

	res := ExecuteCombatTurnCycle(&queue, &currentIdx, heroes, &monsters, rng)

	if len(monsters) != 0 {
		t.Errorf("expected monster to be slain and removed, remaining: %d", len(monsters))
	}
	if res.DungeonState != "Victory" {
		t.Errorf("expected DungeonState 'Victory', got %q", res.DungeonState)
	}
	if !res.QueueReset {
		t.Errorf("expected QueueReset to be true after monster slain")
	}
}
