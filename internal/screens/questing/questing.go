package questing

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/pomquest/internal/adapters/icon"
	"github.com/pomquest/internal/app"
	"github.com/pomquest/internal/domain"
	"github.com/pomquest/internal/dungeon"
	"github.com/pomquest/internal/ui"
)

type tickMsg time.Time
type charTickMsg struct{}

type Combatant struct {
	IsHero     bool
	HeroIdx    int
	MonsterIdx int
	Dexterity  int
	Name       string
}

// Model represents the active questing UI screen.
type Model struct {
	Theme           ui.Theme
	Responsive      *ui.ResponsiveManager
	LogTitle        string
	LogMessages     []string
	Heroes          []ui.HeroCardData
	Supplies        ui.PartyCardData
	ProgressData    QuestProgressData
	Paused          bool
	Quitted         bool
	Width           int
	Height          int
	DBPath          string
	CharDelay       int
	activeAction     string
	activeCharIdx    int
	isTyping         bool
	pendingActions   []string
	DungeonState    string
	TurnPhase       string
	CurrentRoom     int
	ActiveMonsters  []app.Monster
	ActiveHeroIdx   int
	AutoTickCounter int
	TurnTickSecs    int
	PartyID         uint
	GameStateID     uint
	QuestID         uint
	partyHeroes     []app.Hero
	rng             *rand.Rand
	TurnQueue       []Combatant
	CurrentTurnIdx  int
}

// NewModel initializes the questing screen using default DB path.
func NewModel(cfg domain.Config) Model {
	return NewModelWithDBPath(cfg, app.DefaultDBPath())
}

// NewModelWithDBPath initializes the questing screen using data loaded from SQLite DB.
func NewModelWithDBPath(cfg domain.Config, dbPath string) Model {
	theme := ui.ThemeFromConfig(cfg)
	if dbPath == "" {
		dbPath = app.DefaultDBPath()
	}

	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	party, _ := app.LoadActivePartyFromPath(dbPath)
	gameState, _ := app.LoadActiveGameStateFromPath(dbPath)

	var heroes []ui.HeroCardData
	var totalCurrentHP, totalMaxHP int

	partyHeroes := party.Heroes
	for _, h := range partyHeroes {
		iconStr := h.Icon
		if iconStr == "" {
			iconStr = app.GetClassArchetype(h.Class).Icon(icon.ModeNerdFont)
		}
		status := h.Status
		if status == "" {
			status = "Healthy"
		}
		if h.IsDead || h.HP <= 0 {
			status = "Dead"
		}
		totalCurrentHP += h.HP
		totalMaxHP += h.MaxHP

		heroes = append(heroes, ui.HeroCardData{
			Name:     h.Name,
			Class:    h.Class,
			Level:    h.Level,
			HP:       h.HP,
			MaxHP:    h.MaxHP,
			Attack:   h.Attack,
			Armor:    h.Armor,
			Status:   status,
			IsDead:   h.IsDead || h.HP <= 0,
			Icon:     iconStr,
			IsActive: false,
		})
	}

	if len(heroes) == 0 {
		h1 := app.Hero{Creature: app.Creature{Noun: app.Noun{Name: "jimmothy", Icon: "🛡️"}, HP: 10, MaxHP: 10, Attack: 2, Armor: 2}, Class: "fighter", Level: 2, Status: "Healthy"}
		h2 := app.Hero{Creature: app.Creature{Noun: app.Noun{Name: "applebee", Icon: "🧙"}, HP: 6, MaxHP: 6, Attack: 3, Armor: 0}, Class: "wizard", Level: 1, Status: "Healthy"}
		h3 := app.Hero{Creature: app.Creature{Noun: app.Noun{Name: "pam", Icon: "🏹"}, HP: 8, MaxHP: 8, Attack: 3, Armor: 0}, Class: "ranger", Level: 1, Status: "Healthy"}
		h4 := app.Hero{Creature: app.Creature{Noun: app.Noun{Name: "sam", Icon: "🗡️"}, HP: 0, MaxHP: 8, Attack: 3, Armor: 0}, Class: "rouge", Level: 2, Status: "Dead", IsDead: true}
		partyHeroes = []app.Hero{h1, h2, h3, h4}
		for _, h := range partyHeroes {
			heroes = append(heroes, ui.HeroCardData{
				Name:     h.Name,
				Class:    h.Class,
				Level:    h.Level,
				HP:       h.HP,
				MaxHP:    h.MaxHP,
				Attack:   h.Attack,
				Armor:    h.Armor,
				Status:   h.Status,
				IsDead:   h.IsDead,
				Icon:     h.Icon,
				IsActive: false,
			})
		}
		totalCurrentHP = 24
		totalMaxHP = 32
	}

	var magicItems, treasureChests int
	if gameState.Hoard != nil {
		magicItems = gameState.Hoard.MagicItems
		treasureChests = gameState.Hoard.TreasureChests
	}

	supplies := ui.PartyCardData{
		Coins:            party.Coins,
		FoodLbs:          party.FoodLbs,
		Diet:             party.Diet,
		Pace:             party.Pace,
		SpareParts:       party.SpareParts,
		Oxen:             party.Oxen,
		Regents:          party.Regents,
		Potions:          party.Potions,
		Loot:             party.Loot,
		Rope:             party.Rope,
		Bombs:            party.Bombs,
		Shovels:          party.Shovels,
		ItemUsage:        party.ItemUsage,
		Greed:            party.Greed,
		MagicItems:       magicItems,
		TreasureChests:   treasureChests,
		CurrentHP:        totalCurrentHP,
		MaxHP:            totalMaxHP,
		CurrentPartySize: party.CurrentPartySize,
		MaxPartySize:     party.MaxPartySize,
	}

	destination := ""
	terrain := "🌲"
	weather := "☀️"
	hazards := "mud, rain"

	currentSecs := 0
	totalSecs := 1500
	currentRoom := 1
	dungeonState := "Entering"
	var questID uint

	if gameState.CurrentQuest != nil {
		questID = gameState.CurrentQuest.ID
		if gameState.CurrentQuest.Title != "" {
			destination = gameState.CurrentQuest.Title
		}
		if gameState.CurrentQuest.Biome != "" {
			terrain = gameState.CurrentQuest.Biome
		}
		if gameState.CurrentQuest.ElapsedSecs > 0 {
			currentSecs = gameState.CurrentQuest.ElapsedSecs
		}
		if gameState.CurrentQuest.TotalSecs > 0 {
			totalSecs = gameState.CurrentQuest.TotalSecs
		}
		if gameState.CurrentQuest.CurrentRoom > 0 {
			currentRoom = gameState.CurrentQuest.CurrentRoom
		}
		if gameState.CurrentQuest.DungeonState != "" {
			dungeonState = gameState.CurrentQuest.DungeonState
		}
	}

	if destination == "" {
		if gameState.CurrentLocation != nil && gameState.CurrentLocation.Name != "" {
			destination = gameState.CurrentLocation.Name
			if gameState.CurrentLocation.Biome != "" {
				terrain = gameState.CurrentLocation.Biome
			}
		}
	}
	if destination == "" {
		destination = dungeon.GenerateDungeonNameWithRand(rng)
	}

	progressData := QuestProgressData{
		Destination: destination,
		CurrentSecs: currentSecs,
		TotalSecs:   totalSecs,
		Weather:     weather,
		Terrain:     terrain,
		Hazards:     hazards,
		DaysPassed:  12,
		Width:       72,
	}

	logs, _ := app.LoadGameLogsFromPath(dbPath, gameState.ID)
	if len(logs) == 0 {
		logs = []string{
			fmt.Sprintf("🏰 The Party enters %s...", destination),
		}
	}

	m := Model{
		Theme:         theme,
		Responsive:    ui.NewResponsiveManager(80),
		DBPath:        dbPath,
		CharDelay:     cfg.CharDelay,
		LogTitle:      "12, Jan, 1000 A.P.",
		LogMessages:   logs,
		Heroes:        heroes,
		Supplies:      supplies,
		ProgressData:  progressData,
		Width:         80,
		Height:        24,
		DungeonState:  dungeonState,
		TurnPhase:     "hero",
		CurrentRoom:   currentRoom,
		ActiveHeroIdx: 0,
		TurnTickSecs:  0,
		PartyID:       party.ID,
		GameStateID:   gameState.ID,
		QuestID:       questID,
		partyHeroes:   partyHeroes,
		rng:           rng,
	}

	if cfg.CharDelay > 0 && len(logs) > 0 {
		m.activeAction = logs[len(logs)-1]
		m.activeCharIdx = 1
		m.isTyping = true
	}

	return m
}

func (m *Model) SetDimensions(width, height int) {
	m.Width = width
	m.Height = height
	if m.Responsive != nil {
		m.Responsive.SetWidth(width)
	}
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		tea.ClearScreen,
		tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}),
	}
	if m.isTyping && m.CharDelay > 0 {
		cmds = append(cmds, tea.Tick(time.Duration(m.CharDelay)*time.Millisecond, func(t time.Time) tea.Msg {
			return charTickMsg{}
		}))
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Responsive.SetWidth(msg.Width)
		return m, nil

	case charTickMsg:
		if m.isTyping && m.CharDelay > 0 {
			runes := []rune(m.activeAction)
			if m.activeCharIdx < len(runes) {
				m.activeCharIdx++
				return m, tea.Tick(time.Duration(m.CharDelay)*time.Millisecond, func(t time.Time) tea.Msg {
					return charTickMsg{}
				})
			}

			if len(m.pendingActions) > 0 {
				nextMsg := m.pendingActions[0]
				m.pendingActions = m.pendingActions[1:]
				m.LogMessages = append(m.LogMessages, nextMsg)
				m.activeAction = nextMsg
				m.activeCharIdx = 1
				return m, tea.Tick(time.Duration(m.CharDelay)*time.Millisecond, func(t time.Time) tea.Msg {
					return charTickMsg{}
				})
			}

			m.isTyping = false
		}
		return m, nil

	case tickMsg:
		if !m.Paused && m.ProgressData.CurrentSecs < m.ProgressData.TotalSecs {
			if m.isTyping {
				return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
					return tickMsg(t)
				})
			}

			m.ProgressData.CurrentSecs++
			m.AutoTickCounter++
			m.TurnTickSecs++

			var cmd tea.Cmd
			if m.TurnTickSecs < 3 {
				m.updatePacingDots()
			} else {
				m.updatePacingDots()
				cmd = m.advanceDungeonCycle()
				m.TurnTickSecs = 0
			}
			m.saveCurrentState()
			return m, tea.Batch(cmd, tea.Tick(time.Second, func(t time.Time) tea.Msg {
				return tickMsg(t)
			}))
		}
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.saveCurrentState()
			_ = app.ClearGameLogsToPath(m.DBPath, m.GameStateID)
			m.LogMessages = nil
			m.Quitted = true
			return m, tea.Batch(tea.ClearScreen, tea.Quit)

		case "p", "esc":
			m.Paused = !m.Paused

		case " ", "space":
			if !m.Paused {
				if m.isTyping {
					for len(m.pendingActions) > 0 {
						nextMsg := m.pendingActions[0]
						m.pendingActions = m.pendingActions[1:]
						m.LogMessages = append(m.LogMessages, nextMsg)
					}
					m.activeCharIdx = len([]rune(m.activeAction))
					m.isTyping = false
				} else {
					m.TurnTickSecs = 3
					m.updatePacingDots()
					cmd := m.advanceDungeonCycle()
					m.TurnTickSecs = 0
					m.saveCurrentState()
					return m, cmd
				}
			}
		}
	}

	return m, nil
}

func (m *Model) updatePacingDots() {
	if len(m.LogMessages) == 0 {
		m.LogMessages = append(m.LogMessages, ". ")
		return
	}

	lastIdx := len(m.LogMessages) - 1
	lastMsg := m.LogMessages[lastIdx]

	dotContent := ". "
	if m.TurnTickSecs == 2 {
		dotContent = ". . "
	} else if m.TurnTickSecs >= 3 {
		dotContent = ". . . "
	}

	if lastMsg == ". " || lastMsg == ". . " || lastMsg == ". . . " {
		m.LogMessages[lastIdx] = dotContent
	} else {
		m.LogMessages = append(m.LogMessages, dotContent)
	}
}

func (m *Model) appendActionLog(msg string) tea.Cmd {
	if m.CharDelay <= 0 {
		m.LogMessages = append(m.LogMessages, msg)
		m.activeAction = msg
		m.activeCharIdx = len([]rune(msg))
		m.isTyping = false
		return nil
	}

	if !m.isTyping {
		m.LogMessages = append(m.LogMessages, msg)
		m.activeAction = msg
		m.activeCharIdx = 1
		m.isTyping = true
		return tea.Tick(time.Duration(m.CharDelay)*time.Millisecond, func(t time.Time) tea.Msg {
			return charTickMsg{}
		})
	}

	m.pendingActions = append(m.pendingActions, msg)
	return nil
}

func (m Model) displayLogMessages() []string {
	if !m.isTyping || len(m.LogMessages) == 0 {
		return m.LogMessages
	}

	display := make([]string, len(m.LogMessages))
	copy(display, m.LogMessages)

	runes := []rune(m.activeAction)
	if m.activeCharIdx < len(runes) {
		display[len(display)-1] = string(runes[:m.activeCharIdx])
	}
	return display
}

func (m *Model) advanceDungeonCycle() tea.Cmd {
	var cmds []tea.Cmd
	if m.ProgressData.CurrentSecs >= m.ProgressData.TotalSecs {
		if m.DungeonState != "Completed" {
			m.DungeonState = "Completed"
			cmds = append(cmds, m.appendActionLog("🏆 25:00 FOCUS TIMER EXPIRED! Quest completed successfully!"))
		}
		return tea.Batch(cmds...)
	}

	switch m.DungeonState {
	case "Entering":
		m.DungeonState = "InRoom"
		cmds = append(cmds, m.appendActionLog(fmt.Sprintf("🗺️ The Party enters Room #%d...", m.CurrentRoom)))

	case "InRoom":
		crBudget := (m.CurrentRoom / 2) + 1
		maxHeroLevel := 1
		for _, h := range m.partyHeroes {
			if !h.IsDead && h.Level > maxHeroLevel {
				maxHeroLevel = h.Level
			}
		}
		maxCR := maxHeroLevel
		m.ActiveMonsters = dungeon.GenerateMonstersForRoom(m.ProgressData.Terrain, crBudget, m.rng, maxCR)

		if len(m.ActiveMonsters) == 0 {
			m.DungeonState = "Victory"
			m.TurnPhase = "hero"
			return tea.Batch(cmds...)
		}
		var names []string
		for _, mon := range m.ActiveMonsters {
			names = append(names, fmt.Sprintf("%s %s", mon.Icon, mon.Name))
		}
		cmds = append(cmds, m.appendActionLog(fmt.Sprintf("⚔️ Room #%d: Encountered %s!", m.CurrentRoom, strings.Join(names, ", "))))
		m.DungeonState = "Fighting"
		m.TurnPhase = "hero"
		m.ActiveHeroIdx = 0

	case "Fighting":
		if len(m.ActiveMonsters) == 0 {
			m.DungeonState = "Victory"
			m.TurnPhase = "hero"
			return tea.Batch(cmds...)
		}

		// Re-calculate turn order if queue is empty or index exceeds queue length
		if len(m.TurnQueue) == 0 || m.CurrentTurnIdx >= len(m.TurnQueue) {
			m.TurnQueue = nil
			m.CurrentTurnIdx = 0

			// Add active living heroes
			for i, h := range m.partyHeroes {
				if !h.IsDead && h.HP > 0 {
					m.TurnQueue = append(m.TurnQueue, Combatant{
						IsHero:    true,
						HeroIdx:   i,
						Dexterity: h.Dex,
						Name:      h.Name,
					})
				}
			}

			// Add active monsters
			for i, mon := range m.ActiveMonsters {
				if mon.HP > 0 {
					m.TurnQueue = append(m.TurnQueue, Combatant{
						IsHero:     false,
						MonsterIdx: i,
						Dexterity:  mon.Dex,
						Name:       mon.Name,
					})
				}
			}

			if len(m.TurnQueue) == 0 {
				// No combatants left
				livingHeroes := false
				for _, h := range m.partyHeroes {
					if !h.IsDead && h.HP > 0 {
						livingHeroes = true
						break
					}
				}
				if !livingHeroes {
					m.DungeonState = "Defeat"
					cmds = append(cmds, m.appendActionLog("💀 THE PARTY HAS BEEN DEFEATED! All heroes have fallen."))
				} else {
					m.DungeonState = "Victory"
				}
				return tea.Batch(cmds...)
			}

			// Sort by Dex descending using stable sort / rng tie-breaking
			for i := 0; i < len(m.TurnQueue); i++ {
				for j := i + 1; j < len(m.TurnQueue); j++ {
					if m.TurnQueue[j].Dexterity > m.TurnQueue[i].Dexterity || (m.TurnQueue[j].Dexterity == m.TurnQueue[i].Dexterity && m.rng.Intn(2) == 0) {
						m.TurnQueue[i], m.TurnQueue[j] = m.TurnQueue[j], m.TurnQueue[i]
					}
				}
			}
		}

		combatant := m.TurnQueue[m.CurrentTurnIdx]
		m.CurrentTurnIdx++

		if combatant.IsHero {
			hero := &m.partyHeroes[combatant.HeroIdx]
			if hero.IsDead || hero.HP <= 0 {
				return tea.Batch(cmds...)
			}

			// Update UI active card state
			for i := range m.Heroes {
				m.Heroes[i].IsActive = (i == combatant.HeroIdx)
			}
			m.ActiveHeroIdx = combatant.HeroIdx

			if len(m.ActiveMonsters) == 0 {
				m.DungeonState = "Victory"
				return tea.Batch(cmds...)
			}

			targetMon := &m.ActiveMonsters[0]

			// Attacker (Hero) rolls 2d6
			d1 := m.rng.Intn(6) + 1
			d2 := m.rng.Intn(6) + 1
			atkRollSum := d1 + d2

			if d1 == 1 && d2 == 1 {
				// Critical Miss: miss and take 1 recoil damage
				hero.HP--
				if hero.HP <= 0 {
					hero.HP = 0
					hero.IsDead = true
					hero.Status = "Dead"
				}
				cmds = append(cmds, m.appendActionLog(fmt.Sprintf("💥 CRITICAL MISS! 🎲(1,1) %s stumbles and takes 1 recoil dmg!", hero.Name)))
			} else if d1 == 6 && d2 == 6 {
				// Critical Hit: deal max damage
				maxDmg := hero.Attack*2 + 2
				targetMon.HP -= maxDmg
				cmds = append(cmds, m.appendActionLog(fmt.Sprintf("⚡ CRITICAL HIT! 🎲(6,6) %s strikes %s for %d max dmg!", hero.Name, targetMon.Name, maxDmg)))
			} else if atkRollSum <= 4 {
				// Miss on 4 or less
				cmds = append(cmds, m.appendActionLog(fmt.Sprintf("🛡️ 🎲(%d+%d=%d) %s attacks %s, but misses!", d1, d2, atkRollSum, hero.Name, targetMon.Name)))
			} else {
				// Normal Hit attempt -> Defender (Monster) defense roll
				// Defender uses block (con/wil) or dodge (dex/sur) based on highest stat
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
					isDodge = m.rng.Intn(2) == 0
				}
				bestDefStat := defStat
				if isDodge {
					bestDefStat = dodgeStat
				}

				defD := m.rng.Intn(6) + 1
				defTotal := defD + bestDefStat
				atkTotal := atkRollSum + hero.Attack

				rawDmg := hero.Attack
				if rawDmg < 1 {
					rawDmg = 1
				}

				if defTotal >= atkTotal {
					// Full block / dodge
					actionType := "blocks"
					if isDodge {
						actionType = "dodges"
					}
					cmds = append(cmds, m.appendActionLog(fmt.Sprintf("🛡️ 🎲(%d+%d=%d vs %d) %s attacks, but %s %s!", d1, d2, atkTotal, defTotal, hero.Name, targetMon.Name, actionType)))
				} else if defTotal <= atkTotal/2 {
					// Half block/dodge -> full damage taken
					targetMon.HP -= rawDmg
					cmds = append(cmds, m.appendActionLog(fmt.Sprintf("🗡️ 🎲(%d+%d=%d vs %d) %s hits %s for %d dmg!", d1, d2, atkTotal, defTotal, hero.Name, targetMon.Name, rawDmg)))
				} else {
					// Partial defense -> half damage rounded down (min 1)
					dmg := rawDmg / 2
					if dmg < 1 {
						dmg = 1
					}
					targetMon.HP -= dmg
					cmds = append(cmds, m.appendActionLog(fmt.Sprintf("⚔️ 🎲(%d+%d=%d vs %d) %s partially hits %s for %d dmg!", d1, d2, atkTotal, defTotal, hero.Name, targetMon.Name, dmg)))
				}
			}

			// Monster Slay Resolution
			if targetMon.HP <= 0 {
				cmds = append(cmds, m.appendActionLog(fmt.Sprintf("☠️ %s is slain!", targetMon.Name)))
				m.ActiveMonsters = m.ActiveMonsters[1:]
				// Reset turn queue to adjust monster indices
				m.TurnQueue = nil
				m.CurrentTurnIdx = 0
			}

			if len(m.ActiveMonsters) == 0 {
				m.DungeonState = "Victory"
			}

		} else {
			// Monster Action Step
			if combatant.MonsterIdx >= len(m.ActiveMonsters) {
				return tea.Batch(cmds...)
			}
			activeMonster := &m.ActiveMonsters[combatant.MonsterIdx]
			if activeMonster.HP <= 0 {
				return tea.Batch(cmds...)
			}

			// Find target hero (first living hero or tank)
			livingHeroIdx := -1
			for i, h := range m.partyHeroes {
				if !h.IsDead && h.HP > 0 {
					livingHeroIdx = i
					break
				}
			}

			if livingHeroIdx == -1 {
				m.DungeonState = "Defeat"
				cmds = append(cmds, m.appendActionLog("💀 THE PARTY HAS BEEN DEFEATED! All heroes have fallen."))
				return tea.Batch(cmds...)
			}

			targetHero := &m.partyHeroes[livingHeroIdx]

			mD1 := m.rng.Intn(6) + 1
			mD2 := m.rng.Intn(6) + 1
			mAtkSum := mD1 + mD2

			if mD1 == 1 && mD2 == 1 {
				// Monster Critical Miss: takes 1 recoil dmg
				activeMonster.HP--
				cmds = append(cmds, m.appendActionLog(fmt.Sprintf("💥 CRITICAL MISS! 🎲(1,1) %s stumbles and takes 1 recoil dmg!", activeMonster.Name)))
				if activeMonster.HP <= 0 {
					cmds = append(cmds, m.appendActionLog(fmt.Sprintf("☠️ %s is slain!", activeMonster.Name)))
					m.ActiveMonsters = m.ActiveMonsters[1:]
					m.TurnQueue = nil
					m.CurrentTurnIdx = 0
				}
			} else if mD1 == 6 && mD2 == 6 {
				// Monster Critical Hit: deal max damage
				mDmg := activeMonster.Attack*2 + 2
				targetHero.HP -= mDmg
				if targetHero.HP <= 0 {
					targetHero.HP = 0
					targetHero.IsDead = true
					targetHero.Status = "Dead"
					cmds = append(cmds, m.appendActionLog(fmt.Sprintf("⚡ CRITICAL HIT! 🎲(6,6) %s strikes %s for %d critical dmg! %s HAS FALLEN!", activeMonster.Name, targetHero.Name, mDmg, targetHero.Name)))
				} else {
					cmds = append(cmds, m.appendActionLog(fmt.Sprintf("⚡ CRITICAL HIT! 🎲(6,6) %s strikes %s for %d critical dmg!", activeMonster.Name, targetHero.Name, mDmg)))
				}
			} else if mAtkSum <= 4 {
				// Miss on 4 or less
				cmds = append(cmds, m.appendActionLog(fmt.Sprintf("🛡️ 🎲(%d+%d=%d) %s attacks %s, but misses!", mD1, mD2, mAtkSum, activeMonster.Name, targetHero.Name)))
			} else {
				// Normal monster attack -> Hero defense roll (1d6 + best stat)
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
					isDodge = m.rng.Intn(2) == 0
				}
				bestDefStat := defStat
				if isDodge {
					bestDefStat = dodgeStat
				}

				hDefD := m.rng.Intn(6) + 1
				hDefTotal := hDefD + bestDefStat
				mAtkTotal := mAtkSum + activeMonster.Attack

				rawDmg := activeMonster.Attack
				if rawDmg < 1 {
					rawDmg = 1
				}

				if hDefTotal >= mAtkTotal {
					// Full block / dodge
					actionType := "blocks"
					if isDodge {
						actionType = "dodges"
					}
					cmds = append(cmds, m.appendActionLog(fmt.Sprintf("🛡️ 🎲(%d+%d=%d vs %d) %s attacks, but %s %s!", mD1, mD2, mAtkTotal, hDefTotal, activeMonster.Name, targetHero.Name, actionType)))
				} else if hDefTotal <= mAtkTotal/2 {
					// Full hit taken
					targetHero.HP -= rawDmg
					if targetHero.HP <= 0 {
						targetHero.HP = 0
						targetHero.IsDead = true
						targetHero.Status = "Dead"
						cmds = append(cmds, m.appendActionLog(fmt.Sprintf("🩸 🎲(%d+%d=%d vs %d) %s strikes %s for %d dmg! %s HAS FALLEN!", mD1, mD2, mAtkTotal, hDefTotal, activeMonster.Name, targetHero.Name, rawDmg, targetHero.Name)))
					} else {
						cmds = append(cmds, m.appendActionLog(fmt.Sprintf("🐾 🎲(%d+%d=%d vs %d) %s attacks %s for %d dmg!", mD1, mD2, mAtkTotal, hDefTotal, activeMonster.Name, targetHero.Name, rawDmg)))
					}
				} else {
					// Partial defense -> half damage
					dmg := rawDmg / 2
					if dmg < 1 {
						dmg = 1
					}
					targetHero.HP -= dmg
					if targetHero.HP <= 0 {
						targetHero.HP = 0
						targetHero.IsDead = true
						targetHero.Status = "Dead"
						cmds = append(cmds, m.appendActionLog(fmt.Sprintf("🩸 🎲(%d+%d=%d vs %d) %s strikes %s for %d dmg! %s HAS FALLEN!", mD1, mD2, mAtkTotal, hDefTotal, activeMonster.Name, targetHero.Name, dmg, targetHero.Name)))
					} else {
						cmds = append(cmds, m.appendActionLog(fmt.Sprintf("🐾 🎲(%d+%d=%d vs %d) %s attacks %s for %d dmg!", mD1, mD2, mAtkTotal, hDefTotal, activeMonster.Name, targetHero.Name, dmg)))
					}
				}
			}
		}

	case "Victory":
		gold := 10 + m.rng.Intn(15)
		xp := 15 + m.rng.Intn(10)

		m.Supplies.Coins += gold
		m.Supplies.Loot++

		for i := range m.partyHeroes {
			if !m.partyHeroes[i].IsDead && m.partyHeroes[i].HP > 0 {
				m.partyHeroes[i].XP += xp
				if m.partyHeroes[i].XP >= m.partyHeroes[i].Level*25 {
					m.partyHeroes[i].Level++
					m.partyHeroes[i].MaxHP += 2
					m.partyHeroes[i].HP = m.partyHeroes[i].MaxHP
					cmds = append(cmds, m.appendActionLog(fmt.Sprintf("🌟 %s LEVELED UP to Level %d!", m.partyHeroes[i].Name, m.partyHeroes[i].Level)))
				}
			}
		}

		cmds = append(cmds, m.appendActionLog(fmt.Sprintf("✨ Room #%d Cleared! Party gains +%d Gold and +%d XP!", m.CurrentRoom, gold, xp)))
		m.CurrentRoom++
		m.DungeonState = "Entering"
		m.TurnPhase = "hero"

	case "Defeat":
		// Party defeated
	}

	// Keep log size bounded
	if len(m.LogMessages) > 50 {
		m.LogMessages = m.LogMessages[len(m.LogMessages)-50:]
	}

	// Refresh UI cards
	var totalCurrentHP, totalMaxHP int
	for i, h := range m.partyHeroes {
		totalCurrentHP += h.HP
		totalMaxHP += h.MaxHP
		if i < len(m.Heroes) {
			m.Heroes[i].HP = h.HP
			m.Heroes[i].MaxHP = h.MaxHP
			m.Heroes[i].Level = h.Level
			m.Heroes[i].Status = h.Status
			m.Heroes[i].IsDead = h.IsDead || h.HP <= 0
		}
	}
	m.Supplies.CurrentHP = totalCurrentHP
	m.Supplies.MaxHP = totalMaxHP
	return tea.Batch(cmds...)
}

func (m *Model) saveCurrentState() {
	if m.DBPath == "" {
		return
	}

	status := "Active"
	if m.DungeonState == "Completed" {
		status = "Completed"
	} else if m.DungeonState == "Defeat" {
		status = "Failed"
	}

	quest := app.Quest{
		Title:        m.ProgressData.Destination,
		Biome:        m.ProgressData.Terrain,
		ElapsedSecs:  m.ProgressData.CurrentSecs,
		TotalSecs:    m.ProgressData.TotalSecs,
		CurrentRoom:  m.CurrentRoom,
		DungeonState: m.DungeonState,
		Status:       status,
	}
	quest.ID = m.QuestID

	savedQuest, err := app.SaveActiveQuestToPath(m.DBPath, quest)
	if err == nil && savedQuest.ID != 0 {
		m.QuestID = savedQuest.ID
	}

	party := app.Party{
		Coins: m.Supplies.Coins,
		Loot:  m.Supplies.Loot,
	}
	party.ID = m.PartyID
	_ = app.SavePartyToPath(m.DBPath, party)

	if m.PartyID != 0 {
		_ = app.SavePartyHeroesToPath(m.DBPath, m.PartyID, m.partyHeroes)
	}

	if m.GameStateID != 0 {
		_ = app.SaveGameLogsToPath(m.DBPath, m.GameStateID, m.LogMessages)
	}
}

func (m Model) View() string {
	if m.Quitted {
		return ""
	}

	containerWidth := DefaultContainerWidth
	if m.Responsive != nil && m.Responsive.Width() > 0 {
		containerWidth = min(DefaultContainerWidth, m.Responsive.Width())
	}

	logComp := NewLogComponent(m.Theme, m.LogTitle, m.displayLogMessages(), containerWidth)
	heroCardsComp := NewHeroCardsComponent(m.Theme, m.Heroes, containerWidth, m.Responsive)
	suppliesComp := NewSuppliesComponent(m.Theme, m.Supplies, containerWidth, m.Responsive)
	progressComp := NewTimerComponent(m.Theme, m.ProgressData, containerWidth, m.Responsive)

	helpKeys := []key.Binding{
		key.NewBinding(key.WithKeys("space"), key.WithHelp("space", "attack/step")),
		key.NewBinding(key.WithKeys("p/esc"), key.WithHelp("p", "pause")),
		key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
	}
	helpComp := NewHelpComponent(m.Theme, helpKeys)

	qContainer := NewQuestContainerComponent(m.Theme, containerWidth, m.Responsive)
	content := qContainer.Render(
		logComp.View(),
		heroCardsComp.View(),
		suppliesComp.View(),
		progressComp.View(),
		helpComp.View(),
	)

	if m.Width > 0 && m.Height > 0 {
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}
