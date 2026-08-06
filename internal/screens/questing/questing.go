package questing

import (
	"fmt"
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
	LocationName string
	TotalSecs    int
	RemainingSec int
	LogMessages  []string
	Heroes       []ui.HeroCardData
	Supplies     ui.PartyBarData
	Paused       bool
	Quitted      bool
	Width        int
	Height       int
}

// NewModel initializes the questing screen with dummy prototype data.
func NewModel(cfg domain.Config) Model {
	theme := ui.ThemeFromConfig(cfg)
	return Model{
		Theme:        theme,
		Responsive:   ui.NewResponsiveManager(80),
		LocationName: "Dungeon of the Sacred Mirror",
		TotalSecs:    1500, // 25 minutes
		RemainingSec: 1420,
		Heroes: []ui.HeroCardData{
			{Name: "Jimmothy", Class: "Fighter", Level: 1, HP: 10, MaxHP: 10, Attack: 2, Armor: 2, Status: "Healthy"},
			{Name: "Applebee", Class: "Wizard", Level: 1, HP: 6, MaxHP: 6, Attack: 3, Armor: 0, Status: "Sick"},
			{Name: "Pam", Class: "Ranger", Level: 1, HP: 6, MaxHP: 8, Attack: 3, Armor: 0, Status: "Healthy"},
			{Name: "Sam", Class: "Rogue", Level: 1, HP: 0, MaxHP: 8, Attack: 3, Armor: 0, Status: "Dead", IsDead: true},
		},
		Supplies: ui.PartyBarData{
			Gold:       1,
			Silver:     60,
			Copper:     0,
			FoodLbs:    400,
			Diet:       "balanced",
			SpareParts: 2,
			Oxen:       4,
		},
		LogMessages: []string{
			"The party arrives at a wooden dungeon door...",
			"Applebee opens the door...",
			"3 bats attack the party!",
			"Pam shoots at bats... 󰇍 󰇌 hits 2 bats for 󰇊 and 󰇋 dmg!",
			"Sam strikes at 1 bat... 󰇏 󰇎 crit! for 󰇏 dmg!",
			"The party loots the room...",
			"Pam finds 2s!",
			"Sam finds 1 regent!",
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
		if !m.Paused && m.RemainingSec > 0 {
			m.RemainingSec--
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
			// Roll fake dice jingle event into log
			m.LogMessages = append(m.LogMessages, "Pam rolls dice... 󰇍 󰇌 hits for 󰇏 dmg!")
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.Quitted {
		return ""
	}

	// Format timer string (MM:SS)
	mins := m.RemainingSec / 60
	secs := m.RemainingSec % 60
	timerLabel := fmt.Sprintf("⏱  FOCUS TIMER (%02d:%02d)", mins, secs)
	if m.Paused {
		timerLabel += " [PAUSED]"
	}

	// Component hierarchy: questing model -> questing component -> generic ui component
	timerComp := NewTimerComponent(m.Theme, timerLabel, m.TotalSecs-m.RemainingSec, m.TotalSecs)
	logComp := NewLogComponent(m.Theme, m.LogMessages)
	partyComp := NewPartyGridComponent(m.Theme, m.Heroes)
	suppliesComp := NewSuppliesComponent(m.Theme, m.Supplies)

	helpKeys := []key.Binding{
		key.NewBinding(key.WithKeys("space"), key.WithHelp("space", "roll dice 󰇍")),
		key.NewBinding(key.WithKeys("p/esc"), key.WithHelp("p", "pause")),
		key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
	}
	helpComp := NewHelpComponent(m.Theme, helpKeys)

	var sb strings.Builder
	sb.WriteString(timerComp.View() + "\n\n")
	sb.WriteString(partyComp.View() + "\n\n")
	sb.WriteString(logComp.View() + "\n\n")
	sb.WriteString(suppliesComp.View() + "\n")
	sb.WriteString(helpComp.View())

	contentStyle := lipgloss.NewStyle().Align(lipgloss.Left)
	content := contentStyle.Render(sb.String())

	if m.Width > 0 && m.Height > 0 {
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}
