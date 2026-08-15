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
	"github.com/pomquest/internal/ui/dice_animation"
)

type tickMsg time.Time
type charTickMsg struct{}
type diceTickMsg struct{}

// Model represents the active questing UI screen.
type Model struct {
	Theme              ui.Theme
	Responsive         *ui.ResponsiveManager
	LogTitle           string
	LogMessages        []string
	Heroes             []ui.HeroCardData
	Supplies           ui.PartyCardData
	ProgressData       QuestProgressData
	Paused             bool
	Quitted            bool
	Width              int
	Height             int
	CharDelay          int
	DungeonState       string
	TurnPhase          string
	CurrentRoom        int
	ActiveHeroIdx      int
	AutoTickCounter    int
	TurnTickSecs       int
	ActiveMonsters     []app.Monster
	isTyping           bool
	activeAction       string
	activeCharIdx      int
	pendingActions     []string
	DBPath             string
	PartyID            uint
	GameStateID        uint
	QuestID            uint
	partyHeroes        []app.Hero
	rng                *rand.Rand
	TurnQueue          []dungeon.Combatant
	CurrentTurnIdx     int
	StepPhase          int
	TurnCycleState     dungeon.TurnCycleState
	activeDiceAnimator *dice_animation.LineAnimator
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

	case diceTickMsg:
		if m.activeDiceAnimator != nil {
			m.activeDiceAnimator.Tick(m.activeDiceAnimator.FrameInterval())
			if m.activeDiceAnimator.IsFinished() {
				targetIdx := m.activeDiceAnimator.TargetLineIdx
				if targetIdx >= 0 && targetIdx < len(m.LogMessages) {
					m.LogMessages[targetIdx] = m.activeDiceAnimator.Render()
				}
				m.activeDiceAnimator = nil
				if m.isTyping {
					return m, tea.Tick(time.Duration(m.CharDelay)*time.Millisecond, func(t time.Time) tea.Msg {
						return charTickMsg{}
					})
				}
				return m, nil
			}
			return m, tea.Tick(m.activeDiceAnimator.FrameInterval(), func(t time.Time) tea.Msg {
				return diceTickMsg{}
			})
		}
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

			// Pause printing next line while dice are actively rolling
			if m.activeDiceAnimator != nil && !m.activeDiceAnimator.IsFinished() {
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
			if m.isTyping || (m.activeDiceAnimator != nil && !m.activeDiceAnimator.IsFinished()) {
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
				if m.activeDiceAnimator != nil && !m.activeDiceAnimator.IsFinished() {
					targetIdx := m.activeDiceAnimator.TargetLineIdx
					m.activeDiceAnimator.Finish()
					if targetIdx >= 0 && targetIdx < len(m.LogMessages) {
						m.LogMessages[targetIdx] = m.activeDiceAnimator.Render()
					}
					m.activeDiceAnimator = nil
				}

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
	if len(m.LogMessages) == 0 {
		return m.LogMessages
	}

	display := make([]string, len(m.LogMessages))
	copy(display, m.LogMessages)

	if m.isTyping {
		runes := []rune(m.activeAction)
		if m.activeCharIdx < len(runes) {
			display[len(display)-1] = string(runes[:m.activeCharIdx])
		}
	}

	if m.activeDiceAnimator != nil {
		targetIdx := m.activeDiceAnimator.TargetLineIdx
		if targetIdx >= 0 && targetIdx < len(display) {
			display[targetIdx] = m.activeDiceAnimator.Render()
		}
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

		if len(m.TurnQueue) == 0 || m.CurrentTurnIdx >= len(m.TurnQueue) {
			m.TurnQueue = dungeon.BuildTurn(m.partyHeroes, m.ActiveMonsters, m.rng)
			m.CurrentTurnIdx = 0

			if len(m.TurnQueue) == 0 {
				if !dungeon.HasLivingHeroes(m.partyHeroes) {
					m.DungeonState = "Defeat"
					cmds = append(cmds, m.appendActionLog("💀 THE PARTY HAS BEEN DEFEATED! All heroes have fallen."))
				} else {
					m.DungeonState = "Victory"
				}
				return tea.Batch(cmds...)
			}
		}

		res := dungeon.ExecuteCombatTurnCycle(&m.TurnQueue, &m.CurrentTurnIdx, &m.StepPhase, &m.TurnCycleState, m.partyHeroes, &m.ActiveMonsters, m.rng)
		if res.ActiveHeroIdx >= 0 {
			for i := range m.Heroes {
				m.Heroes[i].IsActive = (i == res.ActiveHeroIdx)
			}
			m.ActiveHeroIdx = res.ActiveHeroIdx
		}
		if res.DungeonState != "" {
			m.DungeonState = res.DungeonState
		}
		for _, log := range res.ActionLogs {
			cmd := m.appendActionLog(log)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			if strings.Contains(log, "rolls") {
				var diceTargets []int
				if m.TurnCycleState.AtkD1 > 0 && m.TurnCycleState.AtkD2 > 0 {
					diceTargets = []int{m.TurnCycleState.AtkD1, m.TurnCycleState.AtkD2}
				} else if m.TurnCycleState.DefD > 0 {
					diceTargets = []int{m.TurnCycleState.DefD}
				}
				if len(diceTargets) > 0 {
					targetIdx := len(m.LogMessages) - 1
					if targetIdx < 0 {
						targetIdx = 0
					}
					m.activeDiceAnimator = dice_animation.NewLineAnimator(log, diceTargets, targetIdx, m.rng)
					interval := m.activeDiceAnimator.FrameInterval()
					cmds = append(cmds, tea.Tick(interval, func(t time.Time) tea.Msg {
						return diceTickMsg{}
					}))
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
