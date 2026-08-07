package questing

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/pomquest/internal/domain"
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
	ProgressData ui.QuestProgressData
	Paused       bool
	Quitted      bool
	Width        int
	Height       int
}

// NewModel initializes the questing screen with exact mockup data.
func NewModel(cfg domain.Config) Model {
	theme := ui.ThemeFromConfig(cfg)
	return Model{
		Theme:      theme,
		Responsive: ui.NewResponsiveManager(80),
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
		Heroes: []ui.HeroCardData{
			{Name: "jimmothy", Class: "fighter", Level: 2, HP: 10, MaxHP: 10, Attack: 2, Armor: 2, Status: "Healthy", IsActive: false},
			{Name: "applebee", Class: "wizard", Level: 1, HP: 6, MaxHP: 6, Attack: 3, Armor: 0, Status: "Healthy"},
			{Name: "pam", Class: "ranger", Level: 1, HP: 8, MaxHP: 8, Attack: 3, Armor: 0, Status: "Healthy"},
			{Name: "sam", Class: "rouge", Level: 2, HP: 0, MaxHP: 8, Attack: 3, Armor: 0, Status: "Dead", IsDead: true},
		},
		Supplies: ui.PartyCardData{
			Gold:           0,
			Silver:         10,
			Copper:         80,
			FoodLbs:        2000,
			Diet:           "balanced",
			Pace:           "fast",
			SpareParts:     3,
			Oxen:           4,
			MagicItems:     0,
			Gems:           0,
			TreasureChests: 1,
			CurrentHP:      24,
			MaxHP:          32,
		},
		ProgressData: ui.QuestProgressData{
			Destination: "Forest of Izelliador",
			CurrentSecs: 1260, // 21:00
			TotalSecs:   1500, // 25:00
			Weather:     "☀️",
			Terrain:     "🌲",
			Hazards:     "mud, rain",
			DaysPassed:  12,
			Width:       72,
		},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
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

	logComp := NewLogComponent(m.Theme, m.LogTitle, m.LogMessages)
	partyComp := NewPartyGridComponent(m.Theme, m.Heroes)
	suppliesComp := NewSuppliesComponent(m.Theme, m.Supplies)
	progressComp := NewTimerComponent(m.Theme, m.ProgressData)

	helpKeys := []key.Binding{
		key.NewBinding(key.WithKeys("space"), key.WithHelp("space", "attack/action")),
		key.NewBinding(key.WithKeys("p/esc"), key.WithHelp("p", "pause")),
		key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
	}
	helpComp := NewHelpComponent(m.Theme, helpKeys)

	var sb strings.Builder
	sb.WriteString(logComp.View() + "\n\n")
	sb.WriteString(partyComp.View() + "\n\n")
	sb.WriteString(suppliesComp.View() + "\n")
	sb.WriteString(progressComp.View() + "\n")
	sb.WriteString(helpComp.View())

	contentStyle := lipgloss.NewStyle().Align(lipgloss.Left)
	content := contentStyle.Render(sb.String())

	if m.Width > 0 && m.Height > 0 {
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}
