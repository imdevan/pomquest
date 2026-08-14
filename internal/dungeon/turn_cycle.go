package dungeon

import (
	"fmt"
	"math/rand"

	"github.com/pomquest/internal/app"
)

// Combatant represents an entity (Hero or Monster) participating in combat turn cycle.
type Combatant struct {
	IsHero     bool
	HeroIdx    int
	MonsterIdx int
	Dexterity  int
	Name       string
}

// TurnResult captures the output of executing a single combat turn cycle step.
type TurnResult struct {
	ActiveHeroIdx int
	DungeonState  string
	ActionLogs    []string
	QueueReset    bool
}

// BuildTurn constructs a DEX-sorted turn queue of living heroes and monsters.
func BuildTurn(heroes []app.Hero, monsters []app.Monster, rng *rand.Rand) []Combatant {
	var queue []Combatant

	for i, h := range heroes {
		if !h.IsDead && h.HP > 0 {
			queue = append(queue, Combatant{
				IsHero:    true,
				HeroIdx:   i,
				Dexterity: h.Dex,
				Name:      h.Name,
			})
		}
	}

	for i, mon := range monsters {
		if mon.HP > 0 {
			queue = append(queue, Combatant{
				IsHero:     false,
				MonsterIdx: i,
				Dexterity:  mon.Dex,
				Name:       mon.Name,
			})
		}
	}

	if len(queue) == 0 {
		return queue
	}

	// Sort by Dex descending with RNG tie-breaking
	for i := 0; i < len(queue); i++ {
		for j := i + 1; j < len(queue); j++ {
			if queue[j].Dexterity > queue[i].Dexterity || (queue[j].Dexterity == queue[i].Dexterity && rng != nil && rng.Intn(2) == 0) {
				queue[i], queue[j] = queue[j], queue[i]
			}
		}
	}

	return queue
}

// HasLivingHeroes returns true if at least one hero in the slice is alive.
func HasLivingHeroes(heroes []app.Hero) bool {
	for _, h := range heroes {
		if !h.IsDead && h.HP > 0 {
			return true
		}
	}
	return false
}

// ExecuteCombatTurnCycle executes a single turn step for the current combatant in the cycle queue.
func ExecuteCombatTurnCycle(
	queue *[]Combatant,
	currentIdx *int,
	heroes []app.Hero,
	monsters *[]app.Monster,
	rng *rand.Rand,
) TurnResult {
	res := TurnResult{
		ActiveHeroIdx: -1,
		DungeonState:  "Fighting",
	}

	if queue == nil || len(*queue) == 0 || currentIdx == nil || *currentIdx >= len(*queue) {
		return res
	}

	if rng == nil {
		rng = rand.New(rand.NewSource(1))
	}

	combatant := (*queue)[*currentIdx]
	*currentIdx++

	if combatant.IsHero {
		if combatant.HeroIdx >= len(heroes) {
			return res
		}
		hero := &heroes[combatant.HeroIdx]
		if hero.IsDead || hero.HP <= 0 {
			return res
		}

		res.ActiveHeroIdx = combatant.HeroIdx

		if len(*monsters) == 0 {
			res.DungeonState = "Victory"
			return res
		}

		targetMon := &(*monsters)[0]

		d1 := rng.Intn(6) + 1
		d2 := rng.Intn(6) + 1
		atkRollSum := d1 + d2

		if d1 == 1 && d2 == 1 {
			hero.HP--
			if hero.HP <= 0 {
				hero.HP = 0
				hero.IsDead = true
				hero.Status = "Dead"
			}
			res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("💥 CRITICAL MISS! 🎲(1,1) %s stumbles and takes 1 recoil dmg!", hero.Name))
		} else if d1 == 6 && d2 == 6 {
			maxDmg := hero.Attack*2 + 2
			targetMon.HP -= maxDmg
			res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("⚡ CRITICAL HIT! 🎲(6,6) %s strikes %s for %d max dmg!", hero.Name, targetMon.Name, maxDmg))
		} else if atkRollSum <= 4 {
			res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("🛡️ 🎲(%d+%d=%d) %s attacks %s, but misses!", d1, d2, atkRollSum, hero.Name, targetMon.Name))
		} else {
			defStat := targetMon.Con
			if targetMon.Wil > defStat {
				defStat = targetMon.Wil
			}
			dodgeStat := targetMon.Dex
			if targetMon.Sur > dodgeStat {
				dodgeStat = targetMon.Sur
			}
			isDodge := dodgeStat > defStat
			if !isDodge && dodgeStat == defStat {
				isDodge = rng.Intn(2) == 0
			}
			bestDefStat := defStat
			if isDodge {
				bestDefStat = dodgeStat
			}

			defD := rng.Intn(6) + 1
			defTotal := defD + bestDefStat
			atkTotal := atkRollSum + hero.Attack

			rawDmg := hero.Attack
			if rawDmg < 1 {
				rawDmg = 1
			}

			if defTotal >= atkTotal {
				actionType := "blocks"
				if isDodge {
					actionType = "dodges"
				}
				res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("🛡️ 🎲(%d+%d=%d vs %d) %s attacks, but %s %s!", d1, d2, atkTotal, defTotal, hero.Name, targetMon.Name, actionType))
			} else if defTotal <= atkTotal/2 {
				targetMon.HP -= rawDmg
				res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("🗡️ 🎲(%d+%d=%d vs %d) %s hits %s for %d dmg!", d1, d2, atkTotal, defTotal, hero.Name, targetMon.Name, rawDmg))
			} else {
				dmg := rawDmg / 2
				if dmg < 1 {
					dmg = 1
				}
				targetMon.HP -= dmg
				res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("⚔️ 🎲(%d+%d=%d vs %d) %s partially hits %s for %d dmg!", d1, d2, atkTotal, defTotal, hero.Name, targetMon.Name, dmg))
			}
		}

		if targetMon.HP <= 0 {
			res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("☠️ %s is slain!", targetMon.Name))
			*monsters = (*monsters)[1:]
			*queue = nil
			*currentIdx = 0
			res.QueueReset = true
		}

		if len(*monsters) == 0 {
			res.DungeonState = "Victory"
		}

	} else {
		// Monster Action
		if combatant.MonsterIdx >= len(*monsters) {
			return res
		}
		activeMonster := &(*monsters)[combatant.MonsterIdx]
		if activeMonster.HP <= 0 {
			return res
		}

		livingHeroIdx := -1
		for i, h := range heroes {
			if !h.IsDead && h.HP > 0 {
				livingHeroIdx = i
				break
			}
		}

		if livingHeroIdx == -1 {
			res.DungeonState = "Defeat"
			res.ActionLogs = append(res.ActionLogs, "💀 THE PARTY HAS BEEN DEFEATED! All heroes have fallen.")
			return res
		}

		targetHero := &heroes[livingHeroIdx]

		mD1 := rng.Intn(6) + 1
		mD2 := rng.Intn(6) + 1
		mAtkSum := mD1 + mD2

		if mD1 == 1 && mD2 == 1 {
			activeMonster.HP--
			res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("💥 CRITICAL MISS! 🎲(1,1) %s stumbles and takes 1 recoil dmg!", activeMonster.Name))
			if activeMonster.HP <= 0 {
				res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("☠️ %s is slain!", activeMonster.Name))
				*monsters = (*monsters)[1:]
				*queue = nil
				*currentIdx = 0
				res.QueueReset = true
			}
		} else if mD1 == 6 && mD2 == 6 {
			mDmg := activeMonster.Attack*2 + 2
			targetHero.HP -= mDmg
			if targetHero.HP <= 0 {
				targetHero.HP = 0
				targetHero.IsDead = true
				targetHero.Status = "Dead"
				res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("⚡ CRITICAL HIT! 🎲(6,6) %s strikes %s for %d critical dmg! %s HAS FALLEN!", activeMonster.Name, targetHero.Name, mDmg, targetHero.Name))
			} else {
				res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("⚡ CRITICAL HIT! 🎲(6,6) %s strikes %s for %d critical dmg!", activeMonster.Name, targetHero.Name, mDmg))
			}
		} else if mAtkSum <= 4 {
			res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("🛡️ 🎲(%d+%d=%d) %s attacks %s, but misses!", mD1, mD2, mAtkSum, activeMonster.Name, targetHero.Name))
		} else {
			defStat := targetHero.Con
			if targetHero.Wil > defStat {
				defStat = targetHero.Wil
			}
			dodgeStat := targetHero.Dex
			if targetHero.Sur > dodgeStat {
				dodgeStat = targetHero.Sur
			}
			isDodge := dodgeStat > defStat
			if !isDodge && dodgeStat == defStat {
				isDodge = rng.Intn(2) == 0
			}
			bestDefStat := defStat
			if isDodge {
				bestDefStat = dodgeStat
			}

			hDefD := rng.Intn(6) + 1
			hDefTotal := hDefD + bestDefStat
			mAtkTotal := mAtkSum + activeMonster.Attack

			rawDmg := activeMonster.Attack
			if rawDmg < 1 {
				rawDmg = 1
			}

			if hDefTotal >= mAtkTotal {
				actionType := "blocks"
				if isDodge {
					actionType = "dodges"
				}
				res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("🛡️ 🎲(%d+%d=%d vs %d) %s attacks, but %s %s!", mD1, mD2, mAtkTotal, hDefTotal, activeMonster.Name, targetHero.Name, actionType))
			} else if hDefTotal <= mAtkTotal/2 {
				targetHero.HP -= rawDmg
				if targetHero.HP <= 0 {
					targetHero.HP = 0
					targetHero.IsDead = true
					targetHero.Status = "Dead"
					res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("🩸 🎲(%d+%d=%d vs %d) %s strikes %s for %d dmg! %s HAS FALLEN!", mD1, mD2, mAtkTotal, hDefTotal, activeMonster.Name, targetHero.Name, rawDmg, targetHero.Name))
				} else {
					res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("🐾 🎲(%d+%d=%d vs %d) %s attacks %s for %d dmg!", mD1, mD2, mAtkTotal, hDefTotal, activeMonster.Name, targetHero.Name, rawDmg))
				}
			} else {
				dmg := rawDmg / 2
				if dmg < 1 {
					dmg = 1
				}
				targetHero.HP -= dmg
				if targetHero.HP <= 0 {
					targetHero.HP = 0
					targetHero.IsDead = true
					targetHero.Status = "Dead"
					res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("🩸 🎲(%d+%d=%d vs %d) %s strikes %s for %d dmg! %s HAS FALLEN!", mD1, mD2, mAtkTotal, hDefTotal, activeMonster.Name, targetHero.Name, dmg, targetHero.Name))
				} else {
					res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("🐾 🎲(%d+%d=%d vs %d) %s attacks %s for %d dmg!", mD1, mD2, mAtkTotal, hDefTotal, activeMonster.Name, targetHero.Name, dmg))
				}
			}
		}
	}

	return res
}
