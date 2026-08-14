package dungeon

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/pomquest/internal/adapters/icon"
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

// TurnResult captures the output of executing a single combat turn cycle sub-phase step.
type TurnResult struct {
	ActiveHeroIdx int
	DungeonState  string
	ActionLogs    []string
	QueueReset    bool
}

// TurnCycleState tracks sub-phase progress for the active combat turn step.
type TurnCycleState struct {
	StepPhase   int // 0: declaration, 1: roll, 2: target_reaction, 3: target_roll, 4: resolution, 5: death_declaration
	AtkD1       int
	AtkD2       int
	AtkMod      int
	AtkTotal    int
	IsCritMiss  bool
	IsCritHit   bool
	IsMiss      bool
	DefD        int
	DefStatVal  int
	DefTotal    int
	IsDodge     bool
	Blocked     bool
	Dmg         int
	TargetSlain bool
}

// FormatEntityName prepends the entity icon to its name if not already present.
func FormatEntityName(name string, iconStr string) string {
	name = strings.TrimSpace(name)
	if iconStr == "" {
		return name
	}
	if strings.HasPrefix(name, iconStr) {
		return name
	}
	return iconStr + " " + name
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

// ExecuteCombatTurnCycle executes a single sub-phase step for the current combatant in the cycle queue.
func ExecuteCombatTurnCycle(
	queue *[]Combatant,
	currentIdx *int,
	stepPhase *int,
	state *TurnCycleState,
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
	if stepPhase == nil {
		dummy := 0
		stepPhase = &dummy
	}
	if state == nil {
		state = &TurnCycleState{}
	}

	combatant := (*queue)[*currentIdx]

	var subjectName, subjectIcon string
	var targetName, targetIcon string
	var isSubjectHero bool

	if combatant.IsHero {
		if combatant.HeroIdx >= len(heroes) {
			return res
		}
		hero := &heroes[combatant.HeroIdx]
		if hero.IsDead || hero.HP <= 0 {
			*currentIdx++
			*stepPhase = 0
			*state = TurnCycleState{}
			return res
		}

		res.ActiveHeroIdx = combatant.HeroIdx
		isSubjectHero = true

		if len(*monsters) == 0 {
			res.DungeonState = "Victory"
			return res
		}

		subjectName = hero.Name
		subjectIcon = hero.Icon
		if subjectIcon == "" {
			subjectIcon = icon.GetClassIcon(hero.Class, icon.ModeNerdFont)
		}

		targetMon := &(*monsters)[0]
		targetName = targetMon.Name
		targetIcon = targetMon.Icon
		if targetIcon == "" {
			targetIcon = "👾"
		}
	} else {
		if combatant.MonsterIdx >= len(*monsters) {
			*currentIdx++
			*stepPhase = 0
			*state = TurnCycleState{}
			return res
		}
		activeMonster := &(*monsters)[combatant.MonsterIdx]
		if activeMonster.HP <= 0 {
			*currentIdx++
			*stepPhase = 0
			*state = TurnCycleState{}
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

		subjectName = activeMonster.Name
		subjectIcon = activeMonster.Icon
		if subjectIcon == "" {
			subjectIcon = "👾"
		}

		targetHero := &heroes[livingHeroIdx]
		targetName = targetHero.Name
		targetIcon = targetHero.Icon
		if targetIcon == "" {
			targetIcon = icon.GetClassIcon(targetHero.Class, icon.ModeNerdFont)
		}
	}

	subjectDisp := FormatEntityName(subjectName, subjectIcon)
	targetDisp := FormatEntityName(targetName, targetIcon)

	switch *stepPhase {
	case 0: // Declaration phase: subject attacks target
		res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("%s attacks %s", subjectDisp, targetDisp))
		*stepPhase = 1

	case 1: // Roll phase: subject rolls 1d6 + 1d6 + atk_mod = total_atk 󰓥
		d1 := rng.Intn(6) + 1
		d2 := rng.Intn(6) + 1
		atkRollSum := d1 + d2

		var atkMod int
		if isSubjectHero {
			atkMod = heroes[combatant.HeroIdx].Attack
		} else {
			atkMod = (*monsters)[combatant.MonsterIdx].Attack
		}

		atkTotal := atkRollSum + atkMod
		state.AtkD1 = d1
		state.AtkD2 = d2
		state.AtkMod = atkMod
		state.AtkTotal = atkTotal

		atkIcon := icon.GetNerdFont(icon.KeyAttack)
		d1Icon := icon.GetD6Icon(d1, icon.ModeNerdFont)
		d2Icon := icon.GetD6Icon(d2, icon.ModeNerdFont)
		var atkExpr string
		if atkMod >= 0 {
			atkExpr = fmt.Sprintf("%s + %s + %d", d1Icon, d2Icon, atkMod)
		} else {
			atkExpr = fmt.Sprintf("%s + %s - %d", d1Icon, d2Icon, -atkMod)
		}
		res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("%s rolls %s = %d %s", subjectDisp, atkExpr, atkTotal, atkIcon))

		if d1 == 1 && d2 == 1 {
			// Critical miss
			state.IsCritMiss = true
			if isSubjectHero {
				hero := &heroes[combatant.HeroIdx]
				hero.HP--
				if hero.HP <= 0 {
					hero.HP = 0
					hero.IsDead = true
					hero.Status = "Dead"
				}
			} else {
				mon := &(*monsters)[combatant.MonsterIdx]
				mon.HP--
			}
			res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("💥 CRITICAL MISS! %s stumbles and takes 1 recoil dmg!", subjectDisp))
			state.IsMiss = true
			*stepPhase = 4
		} else if d1 == 6 && d2 == 6 {
			// Critical hit
			state.IsCritHit = true
			if isSubjectHero {
				state.Dmg = heroes[combatant.HeroIdx].Attack*2 + 2
			} else {
				state.Dmg = (*monsters)[combatant.MonsterIdx].Attack*2 + 2
			}
			res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("⚡ CRITICAL HIT! %s strikes for %d max dmg!", subjectDisp, state.Dmg))
			*stepPhase = 4
		} else if atkRollSum <= 4 {
			// Normal miss
			state.IsMiss = true
			*stepPhase = 4
		} else {
			*stepPhase = 2
		}

	case 2: // Target reaction phase: target attempts to block or dodge
		var defStat, dodgeStat int
		if isSubjectHero {
			targetMon := &(*monsters)[0]
			defStat = targetMon.Con
			if targetMon.Wil > defStat {
				defStat = targetMon.Wil
			}
			dodgeStat = targetMon.Dex
			if targetMon.Sur > dodgeStat {
				dodgeStat = targetMon.Sur
			}
		} else {
			targetHero := &heroes[0]
			defStat = targetHero.Con
			if targetHero.Wil > defStat {
				defStat = targetHero.Wil
			}
			dodgeStat = targetHero.Dex
			if targetHero.Sur > dodgeStat {
				dodgeStat = targetHero.Sur
			}
		}

		isDodge := dodgeStat > defStat
		if !isDodge && dodgeStat == defStat {
			isDodge = rng.Intn(2) == 0
		}
		state.IsDodge = isDodge

		reactionType := "block"
		if isDodge {
			reactionType = "dodge"
		}
		res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("%s attempts to %s", targetDisp, reactionType))
		*stepPhase = 3

	case 3: // Target roll phase: target rolls DC / defense total
		defD := rng.Intn(6) + 1
		state.DefD = defD

		var bestDefStat int
		if isSubjectHero {
			targetMon := &(*monsters)[0]
			if state.IsDodge {
				bestDefStat = targetMon.Dex
				if targetMon.Sur > bestDefStat {
					bestDefStat = targetMon.Sur
				}
			} else {
				bestDefStat = targetMon.Con
				if targetMon.Wil > bestDefStat {
					bestDefStat = targetMon.Wil
				}
			}
		} else {
			targetHero := &heroes[0]
			if state.IsDodge {
				bestDefStat = targetHero.Dex
				if targetHero.Sur > bestDefStat {
					bestDefStat = targetHero.Sur
				}
			} else {
				bestDefStat = targetHero.Con
				if targetHero.Wil > bestDefStat {
					bestDefStat = targetHero.Wil
				}
			}
		}
		state.DefStatVal = bestDefStat

		defTotal := defD + bestDefStat
		state.DefTotal = defTotal

		defIcon := icon.GetNerdFont(icon.KeyBlock)
		if state.IsDodge {
			defIcon = icon.GetNerdFont(icon.KeyDodge)
		}

		defDIcon := icon.GetD6Icon(defD, icon.ModeNerdFont)
		var defExpr string
		if bestDefStat >= 0 {
			defExpr = fmt.Sprintf("%s + %d", defDIcon, bestDefStat)
		} else {
			defExpr = fmt.Sprintf("%s - %d", defDIcon, -bestDefStat)
		}

		res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("%s rolls %s = %d %s", targetDisp, defExpr, defTotal, defIcon))

		var rawDmg int
		if isSubjectHero {
			rawDmg = heroes[combatant.HeroIdx].Attack
		} else {
			rawDmg = (*monsters)[combatant.MonsterIdx].Attack
		}
		if rawDmg < 1 {
			rawDmg = 1
		}

		if defTotal >= state.AtkTotal {
			state.Blocked = true
		} else if defTotal <= state.AtkTotal/2 {
			state.Dmg = rawDmg
		} else {
			dmg := rawDmg / 2
			if dmg < 1 {
				dmg = 1
			}
			state.Dmg = dmg
		}
		*stepPhase = 4

	case 4: // Resolution phase: subject hits target for X dmg or misses/blocked
		if state.IsMiss {
			res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("%s attacks %s, but misses!", subjectDisp, targetDisp))
			*stepPhase = 0
			*state = TurnCycleState{}
			*currentIdx++
		} else if state.Blocked {
			actionType := "blocks"
			if state.IsDodge {
				actionType = "dodges"
			}
			res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("%s %s the attack!", targetDisp, actionType))
			*stepPhase = 0
			*state = TurnCycleState{}
			*currentIdx++
		} else {
			if isSubjectHero {
				targetMon := &(*monsters)[0]
				targetMon.HP -= state.Dmg
				res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("%s hits %s for %d dmg!", subjectDisp, targetDisp, state.Dmg))
				if targetMon.HP <= 0 {
					state.TargetSlain = true
					*stepPhase = 5
				} else {
					*stepPhase = 0
					*state = TurnCycleState{}
					*currentIdx++
				}
			} else {
				livingHeroIdx := -1
				for i, h := range heroes {
					if !h.IsDead && h.HP > 0 {
						livingHeroIdx = i
						break
					}
				}
				if livingHeroIdx != -1 {
					targetHero := &heroes[livingHeroIdx]
					targetHero.HP -= state.Dmg
					if targetHero.HP <= 0 {
						targetHero.HP = 0
						targetHero.IsDead = true
						targetHero.Status = "Dead"
						state.TargetSlain = true
					}
				}
				res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("%s hits %s for %d dmg!", subjectDisp, targetDisp, state.Dmg))
				if state.TargetSlain {
					*stepPhase = 5
				} else {
					*stepPhase = 0
					*state = TurnCycleState{}
					*currentIdx++
				}
			}
		}

	case 5: // Death announcement phase: optional target is slain
		res.ActionLogs = append(res.ActionLogs, fmt.Sprintf("💀 %s dies 😵", targetDisp))
		if isSubjectHero {
			if len(*monsters) > 0 {
				*monsters = (*monsters)[1:]
			}
		}
		*queue = nil
		*currentIdx = 0
		*stepPhase = 0
		*state = TurnCycleState{}
		res.QueueReset = true

		if len(*monsters) == 0 {
			res.DungeonState = "Victory"
		} else if !HasLivingHeroes(heroes) {
			res.DungeonState = "Defeat"
		}
	}

	return res
}
