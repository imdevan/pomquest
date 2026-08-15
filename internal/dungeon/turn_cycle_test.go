package dungeon

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/pomquest/internal/adapters/icon"
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
		{Creature: app.Creature{Noun: app.Noun{Name: "Fighter"}, HP: 10, MaxHP: 10, Attack: 3}, Class: "Fighter", Stats: app.Stats{Dex: 4}},
	}

	monsters := []app.Monster{
		{Creature: app.Creature{Noun: app.Noun{Name: "Goblin", Icon: "👺"}, HP: 5, MaxHP: 5, Armor: 1}, Stats: app.Stats{Dex: 2}},
	}

	queue := BuildTurn(heroes, monsters, rng)
	currentIdx := 0
	stepPhase := 0
	state := TurnCycleState{}

	// Phase 0: Declaration
	res0 := ExecuteCombatTurnCycle(&queue, &currentIdx, &stepPhase, &state, heroes, &monsters, rng)
	if len(res0.ActionLogs) == 0 || !strings.Contains(res0.ActionLogs[0], "attacks") {
		t.Errorf("expected declaration log containing 'attacks', got: %v", res0.ActionLogs)
	}

	// Phase 1: Roll
	res1 := ExecuteCombatTurnCycle(&queue, &currentIdx, &stepPhase, &state, heroes, &monsters, rng)
	if len(res1.ActionLogs) == 0 || !strings.Contains(res1.ActionLogs[0], "rolls") {
		t.Errorf("expected roll log containing 'rolls', got: %v", res1.ActionLogs)
	}
}

func TestExecuteCombatTurnCycle_SlayMonster(t *testing.T) {
	rng := rand.New(rand.NewSource(2))

	heroes := []app.Hero{
		{Creature: app.Creature{Noun: app.Noun{Name: "Hero"}, HP: 10, MaxHP: 10, Attack: 100}, Class: "Fighter", Stats: app.Stats{Dex: 10}},
	}

	monsters := []app.Monster{
		{Creature: app.Creature{Noun: app.Noun{Name: "Weak Goblin", Icon: "👺"}, HP: 1, MaxHP: 1}, Stats: app.Stats{Dex: 1}},
	}

	queue := BuildTurn(heroes, monsters, rng)
	currentIdx := 0
	stepPhase := 0
	state := TurnCycleState{}

	// Step through sub-phases until monster is slain
	var lastRes TurnResult
	for i := 0; i < 10; i++ {
		lastRes = ExecuteCombatTurnCycle(&queue, &currentIdx, &stepPhase, &state, heroes, &monsters, rng)
		if lastRes.QueueReset {
			break
		}
	}

	if len(monsters) != 0 {
		t.Errorf("expected monster to be slain and removed, remaining: %d", len(monsters))
	}
	if lastRes.DungeonState != "Victory" {
		t.Errorf("expected DungeonState 'Victory', got %q", lastRes.DungeonState)
	}
}

func TestTurnCyclePhasesAndIconography(t *testing.T) {
	rng := rand.New(rand.NewSource(10))

	heroes := []app.Hero{
		{
			Creature: app.Creature{
				Noun:   app.Noun{Name: "Arsel", Icon: "⚔️"},
				HP:     10,
				MaxHP:  10,
				Attack: 2,
			},
			Class: "Fighter",
			Stats: app.Stats{Dex: 5, Con: 2},
		},
	}

	monsters := []app.Monster{
		{
			Creature: app.Creature{
				Noun:   app.Noun{Name: "bat", Icon: "🦇"},
				HP:     2,
				MaxHP:  2,
				Attack: 1,
			},
			Stats: app.Stats{Dex: 1, Sur: 1, Con: 0},
		},
	}

	queue := BuildTurn(heroes, monsters, rng)
	currentIdx := 0
	stepPhase := 0
	state := TurnCycleState{}

	// Phase 0: declaration
	res0 := ExecuteCombatTurnCycle(&queue, &currentIdx, &stepPhase, &state, heroes, &monsters, rng)
	if len(res0.ActionLogs) == 0 || res0.ActionLogs[0] != "⚔️ Arsel attacks 🦇 bat" {
		t.Errorf("Phase 0 declaration log mismatch: got %v", res0.ActionLogs)
	}

	// Phase 1: roll
	res1 := ExecuteCombatTurnCycle(&queue, &currentIdx, &stepPhase, &state, heroes, &monsters, rng)
	if len(res1.ActionLogs) == 0 || !strings.HasPrefix(res1.ActionLogs[0], "⚔️ Arsel rolls ") {
		t.Errorf("Phase 1 roll log mismatch: got %v", res1.ActionLogs)
	}
	d1Icon := icon.GetD6Icon(state.AtkD1, icon.ModeNerdFont)
	d2Icon := icon.GetD6Icon(state.AtkD2, icon.ModeNerdFont)
	if !strings.Contains(res1.ActionLogs[0], d1Icon) || !strings.Contains(res1.ActionLogs[0], d2Icon) {
		t.Errorf("Phase 1 roll log missing D6 icons (%s, %s): got %s", d1Icon, d2Icon, res1.ActionLogs[0])
	}

	// Phase 2: target reaction (or phase 4 if miss/crit)
	if stepPhase == 2 {
		res2 := ExecuteCombatTurnCycle(&queue, &currentIdx, &stepPhase, &state, heroes, &monsters, rng)
		if len(res2.ActionLogs) == 0 || !strings.HasPrefix(res2.ActionLogs[0], "🦇 bat attempts to ") {
			t.Errorf("Phase 2 reaction log mismatch: got %v", res2.ActionLogs)
		}

		// Phase 3: target roll
		res3 := ExecuteCombatTurnCycle(&queue, &currentIdx, &stepPhase, &state, heroes, &monsters, rng)
		if len(res3.ActionLogs) == 0 || !strings.HasPrefix(res3.ActionLogs[0], "🦇 bat rolls ") {
			t.Errorf("Phase 3 target roll log mismatch: got %v", res3.ActionLogs)
		}
		defDIcon := icon.GetD6Icon(state.DefD, icon.ModeNerdFont)
		if !strings.Contains(res3.ActionLogs[0], defDIcon) {
			t.Errorf("Phase 3 target roll log missing D6 icon (%s): got %s", defDIcon, res3.ActionLogs[0])
		}
	}
}

func TestCriticalHitMissResolution(t *testing.T) {
	heroes := []app.Hero{
		{Creature: app.Creature{Noun: app.Noun{Name: "Arsel", Icon: "⚔️"}, HP: 10, MaxHP: 10, Attack: 3, Armor: 2}, Class: "Warrior", Stats: app.Stats{Dex: 5}},
	}
	monsters := []app.Monster{
		{Creature: app.Creature{Noun: app.Noun{Name: "bat", Icon: "🦇"}, HP: 10, MaxHP: 10, Attack: 2, Armor: 1}, Stats: app.Stats{Dex: 3}},
	}

	queue := []Combatant{
		{IsHero: true, HeroIdx: 0, Name: "Arsel", Dexterity: 5},
	}
	currentIdx := 0
	stepPhase := 4
	state := TurnCycleState{
		IsCritHit: true,
		Dmg:       8,
	}

	res := ExecuteCombatTurnCycle(&queue, &currentIdx, &stepPhase, &state, heroes, &monsters, rand.New(rand.NewSource(1)))
	if len(res.ActionLogs) == 0 {
		t.Fatalf("expected action logs in phase 4 for critical hit")
	}

	critLog := res.ActionLogs[0]
	if !strings.Contains(critLog, "⚡ CRITICAL HIT!") {
		t.Errorf("expected resolution log to contain '⚡ CRITICAL HIT!', got: %s", critLog)
	}
	if strings.Contains(critLog, "hits") {
		t.Errorf("critical hit log should replace standard hits message, got: %s", critLog)
	}
}


