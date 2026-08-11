package questing

import (
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
	Theme        ui.Theme
	Responsive   *ui.ResponsiveManager
	LogTitle     string
	LogMessages  []string
	Heroes       []ui.HeroCardData
	Supplies     ui.PartyCardData
	ProgressData QuestProgressData
	Paused       bool
	Quitted      bool
	Width        int
	Height       int
	DBPath       string
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

	party, _ := app.LoadActivePartyFromPath(dbPath)
	gameState, _ := app.LoadActiveGameStateFromPath(dbPath)

	var heroes []ui.HeroCardData
	var totalCurrentHP, totalMaxHP int

	for _, h := range party.Heroes {
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
		heroes = []ui.HeroCardData{
			{Name: "jimmothy", Class: "fighter", Level: 2, HP: 10, MaxHP: 10, Attack: 2, Armor: 2, Status: "Healthy", IsActive: false},
			{Name: "applebee", Class: "wizard", Level: 1, HP: 6, MaxHP: 6, Attack: 3, Armor: 0, Status: "Healthy"},
			{Name: "pam", Class: "ranger", Level: 1, HP: 8, MaxHP: 8, Attack: 3, Armor: 0, Status: "Healthy"},
			{Name: "sam", Class: "rouge", Level: 2, HP: 0, MaxHP: 8, Attack: 3, Armor: 0, Status: "Dead", IsDead: true},
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

	if gameState.CurrentLocation != nil && gameState.CurrentLocation.Name != "" {
		destination = gameState.CurrentLocation.Name
		if gameState.CurrentLocation.Biome != "" {
			terrain = gameState.CurrentLocation.Biome
		}
	} else if gameState.CurrentQuest != nil && gameState.CurrentQuest.Title != "" {
		destination = gameState.CurrentQuest.Title
		if gameState.CurrentQuest.Weather != "" {
			weather = gameState.CurrentQuest.Weather
		}
	}
	if destination == "" {
		destination = dungeon.GenerateDungeonName()
	}

	progressData := QuestProgressData{
		Destination: destination,
		CurrentSecs: 1260, // 21:00
		TotalSecs:   1500, // 25:00
		Weather:     weather,
		Terrain:     terrain,
		Hazards:     hazards,
		DaysPassed:  12,
		Width:       72,
	}

	return Model{
		Theme:      theme,
		Responsive: ui.NewResponsiveManager(80),
		DBPath:     dbPath,
		LogTitle:   "12, Jan, 1000 A.P.",
		LogMessages: []string{
			"...",
			"6 goblins ambush the party!",
			"...",
			"jimmothy attacks   2/2 hit for ...   dmg!",
			"...",
			"4 goblins remaining...",
			"...",
			"applebee casts magic missle     1/3 hit for ...  /2 = 3 dmg!",
			"...",
			"last remaining goblin fled!",
			"...",
			"+ 3 xp! +3 copper!",
		},
		Heroes:       heroes,
		Supplies:     supplies,
		ProgressData: progressData,
		Width:        80,
		Height:       24,
	}
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
		}
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.Quitted = true
			return m, tea.Quit

		case "p", "esc":
			m.Paused = !m.Paused

		case " ", "space":
			// Add turn action log entry
			m.LogMessages = append(m.LogMessages, "...", "pam attacks 1/1 hit for 3 dmg!")
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.Quitted {
		return ""
	}

	containerWidth := DefaultContainerWidth // 96 container width requirement (18.1.1)
	if m.Responsive != nil && m.Responsive.Width() > 0 {
		containerWidth = min(DefaultContainerWidth, m.Responsive.Width())
	}

	logComp := NewLogComponent(m.Theme, m.LogTitle, m.LogMessages, containerWidth)
	heroCardsComp := NewHeroCardsComponent(m.Theme, m.Heroes, containerWidth, m.Responsive)
	suppliesComp := NewSuppliesComponent(m.Theme, m.Supplies, containerWidth, m.Responsive)
	progressComp := NewTimerComponent(m.Theme, m.ProgressData, containerWidth, m.Responsive)

	helpKeys := []key.Binding{
		key.NewBinding(key.WithKeys("space"), key.WithHelp("space", "attack/action")),
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
