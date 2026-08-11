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
			fmt.Sprintf("🏰 Party enters location: %s (%s)", destination, terrain),
			fmt.Sprintf("🗺️ Entering Room #%d...", currentRoom),
		}
	}

	m := Model{
		Theme:         theme,
		Responsive:    ui.NewResponsiveManager(80),
		DBPath:        dbPath,
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
	return tea.Batch(
		tea.ClearScreen,
		tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Responsive.SetWidth(msg.Width)
		return m, nil

	case tickMsg:
		if !m.Paused && m.ProgressData.CurrentSecs < m.ProgressData.TotalSecs {
			m.ProgressData.CurrentSecs++
			m.AutoTickCounter++
			m.TurnTickSecs++

			if m.TurnTickSecs < 3 {
				m.updatePacingDots()
			} else {
				m.updatePacingDots()
				m.advanceDungeonCycle()
				m.TurnTickSecs = 0
			}
			m.saveCurrentState()
		}
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.saveCurrentState()
			m.Quitted = true
			return m, tea.Quit

		case "p", "esc":
			m.Paused = !m.Paused

		case " ", "space":
			if !m.Paused {
				m.TurnTickSecs = 3
				m.updatePacingDots()
				m.advanceDungeonCycle()
				m.TurnTickSecs = 0
				m.saveCurrentState()
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

func (m *Model) advanceDungeonCycle() {
	if m.ProgressData.CurrentSecs >= m.ProgressData.TotalSecs {
		if m.DungeonState != "Completed" {
			m.DungeonState = "Completed"
			m.LogMessages = append(m.LogMessages, "🏆 25:00 FOCUS TIMER EXPIRED! Quest completed successfully!")
		}
		return
	}

	switch m.DungeonState {
	case "Entering":
		m.DungeonState = "InRoom"
		m.LogMessages = append(m.LogMessages, fmt.Sprintf("🗺️ Party enters Room #%d...", m.CurrentRoom))

	case "InRoom":
		crBudget := (m.CurrentRoom / 2) + 1
		m.ActiveMonsters = dungeon.GenerateMonstersForRoom(m.ProgressData.Terrain, crBudget, m.rng)
		if len(m.ActiveMonsters) == 0 {
			m.DungeonState = "Victory"
			m.TurnPhase = "hero"
			return
		}
		var names []string
		for _, mon := range m.ActiveMonsters {
			names = append(names, fmt.Sprintf("%s %s", mon.Icon, mon.Name))
		}
		m.LogMessages = append(m.LogMessages, fmt.Sprintf("⚔️ Room #%d: Encountered %s!", m.CurrentRoom, strings.Join(names, ", ")))
		m.DungeonState = "Fighting"
		m.TurnPhase = "hero"
		m.ActiveHeroIdx = 0

	case "Fighting":
		if len(m.ActiveMonsters) == 0 {
			m.DungeonState = "Victory"
			m.TurnPhase = "hero"
			return
		}

		// Find next living hero
		livingHeroIdx := -1
		for i := 0; i < len(m.partyHeroes); i++ {
			idx := (m.ActiveHeroIdx + i) % len(m.partyHeroes)
			if !m.partyHeroes[idx].IsDead && m.partyHeroes[idx].HP > 0 {
				livingHeroIdx = idx
				break
			}
		}

		if livingHeroIdx == -1 {
			m.DungeonState = "Defeat"
			m.LogMessages = append(m.LogMessages, "💀 THE PARTY HAS BEEN DEFEATED! All heroes have fallen.")
			return
		}

		if m.TurnPhase == "monster" {
			// --- MONSTER ACTION STEP ---
			activeMonster := m.ActiveMonsters[0]
			targetHero := &m.partyHeroes[livingHeroIdx]

			mD1 := m.rng.Intn(6) + 1
			mD2 := m.rng.Intn(6) + 1
			mRoll := mD1 + mD2 + activeMonster.Attack
			heroDC := targetHero.Armor + 4

			if mRoll >= heroDC {
				mDmg := activeMonster.Attack
				if mDmg < 1 {
					mDmg = 1
				}
				targetHero.HP -= mDmg
				if targetHero.HP <= 0 {
					targetHero.HP = 0
					targetHero.IsDead = true
					targetHero.Status = "Dead"
					m.LogMessages = append(m.LogMessages, fmt.Sprintf("🩸 %s strikes %s for %d dmg! %s HAS FALLEN!", activeMonster.Name, targetHero.Name, mDmg, targetHero.Name))
				} else {
					m.LogMessages = append(m.LogMessages, fmt.Sprintf("🐾 %s attacks %s for %d dmg!", activeMonster.Name, targetHero.Name, mDmg))
				}
			} else {
				m.LogMessages = append(m.LogMessages, fmt.Sprintf("🛡️ %s attacks %s, but misses!", activeMonster.Name, targetHero.Name))
			}

			m.TurnPhase = "hero"
			m.ActiveHeroIdx = (livingHeroIdx + 1) % len(m.partyHeroes)

		} else {
			// --- HERO ACTION STEP ---
			m.ActiveHeroIdx = livingHeroIdx
			hero := &m.partyHeroes[m.ActiveHeroIdx]

			// Update UI active state
			for i := range m.Heroes {
				m.Heroes[i].IsActive = (i == m.ActiveHeroIdx)
			}

			targetMon := &m.ActiveMonsters[0]
			d1 := m.rng.Intn(6) + 1
			d2 := m.rng.Intn(6) + 1
			atkMod := hero.Attack

			if d1 == 1 && d2 == 1 {
				// Crit Miss
				hero.HP--
				if hero.HP <= 0 {
					hero.HP = 0
					hero.IsDead = true
					hero.Status = "Dead"
				}
				m.LogMessages = append(m.LogMessages, fmt.Sprintf("💥 CRITICAL MISS! 🎲(1,1) %s stumbles and suffers 1 recoil dmg!", hero.Name))
			} else if d1 == 6 && d2 == 6 {
				// Crit Hit
				dmg := hero.Attack*2 + 2
				targetMon.HP -= dmg
				m.LogMessages = append(m.LogMessages, fmt.Sprintf("⚡ CRITICAL HIT! 🎲(6,6) %s strikes %s for %d critical dmg!", hero.Name, targetMon.Name, dmg))
			} else {
				totalRoll := d1 + d2 + atkMod
				dc := targetMon.Armor + 4
				if totalRoll >= dc {
					dmg := hero.Attack
					if dmg < 1 {
						dmg = 1
					}
					targetMon.HP -= dmg
					m.LogMessages = append(m.LogMessages, fmt.Sprintf("🗡️ 🎲(%d+%d+%d=%d vs DC %d) %s hits %s for %d dmg!", d1, d2, atkMod, totalRoll, dc, hero.Name, targetMon.Name, dmg))
				} else {
					m.LogMessages = append(m.LogMessages, fmt.Sprintf("🛡️ 🎲(%d+%d+%d=%d vs DC %d) %s attacks %s, but armor deflects!", d1, d2, atkMod, totalRoll, dc, hero.Name, targetMon.Name))
				}
			}

			if targetMon.HP <= 0 {
				m.LogMessages = append(m.LogMessages, fmt.Sprintf("☠️ %s is slain!", targetMon.Name))
				m.ActiveMonsters = m.ActiveMonsters[1:]
			}

			if len(m.ActiveMonsters) == 0 {
				m.DungeonState = "Victory"
				m.TurnPhase = "hero"
			} else {
				m.TurnPhase = "monster"
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
					m.LogMessages = append(m.LogMessages, fmt.Sprintf("🌟 %s LEVELED UP to Level %d!", m.partyHeroes[i].Name, m.partyHeroes[i].Level))
				}
			}
		}

		m.LogMessages = append(m.LogMessages, fmt.Sprintf("✨ Room #%d Cleared! Party gains +%d Gold and +%d XP!", m.CurrentRoom, gold, xp))
		m.CurrentRoom++
		m.DungeonState = "InRoom"
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

	logComp := NewLogComponent(m.Theme, m.LogTitle, m.LogMessages, containerWidth)
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
