package dashboard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/pomquest/internal/app"
	"github.com/pomquest/internal/domain"
	"github.com/pomquest/internal/ui"
)

// Choice represents an item selected on the dashboard menu.
type Choice int

const (
	ChoiceQuest Choice = iota
	ChoiceParty
	ChoiceProvision
	ChoiceHoard
	ChoiceSettings
	ChoiceQuit
)

// Model represents the main menu / hearth dashboard UI screen.
type Model struct {
	Theme            ui.Theme
	Responsive       *ui.ResponsiveManager
	Choices          []string
	Cursor           int
	Selected         Choice
	CurrentPartySize int
	MaxPartySize     int
	PartyGold        int
	Quitted          bool
	StatusNotice     string
	Width            int
	Height           int
}

// NewModel initializes a dashboard UI with default menu items.
func NewModel(cfg domain.Config) Model {
	theme := ui.ThemeFromConfig(cfg)
	return Model{
		Theme:      theme,
		Responsive: ui.NewResponsiveManager(80),
		Choices: []string{
			"1. ⚔  Quest! (Start Journey)",
			"2. 👥 Party",
			"3. 🎒 Provisions",
			"4. 💎 View Dragon Hoard",
			"5. ⚙  Settings & Config",
			"6. 🚪 Quit",
		},
		Cursor:           0,
		Selected:         -1,
		CurrentPartySize: 0,
		MaxPartySize:     1,
		PartyGold:        1,
		Width:            80,
		Height:           24,
	}
}

// NewModelFromParty initializes a dashboard UI with data populated from DB Party model.
func NewModelFromParty(cfg domain.Config, party app.Party) Model {
	m := NewModel(cfg)
	m.UpdatePartyData(party)
	return m
}

// UpdatePartyData updates party stats displayed on the dashboard.
func (m *Model) UpdatePartyData(party app.Party) {
	var activeCount int
	for _, h := range party.Heroes {
		if !h.IsDead {
			activeCount++
		}
	}
	m.CurrentPartySize = activeCount
	maxSize := party.MaxPartySize
	if maxSize <= 0 {
		maxSize = 1
	}
	m.MaxPartySize = maxSize
	m.PartyGold = party.Gold
}

func (m *Model) SetDimensions(width, height int) {
	m.Width = width
	m.Height = height
	if m.Responsive != nil {
		m.Responsive.SetWidth(width)
	}
}

func (m Model) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Responsive.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.Quitted = true
			m.Selected = ChoiceQuit
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			} else {
				m.Cursor = len(m.Choices) - 1
			}

		case "down", "j":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			} else {
				m.Cursor = 0
			}

		case "1":
			m.Cursor = 0
			return m.handleChoiceSelection(ChoiceQuest)
		case "2":
			m.Cursor = 1
			return m.handleChoiceSelection(ChoiceParty)
		case "3":
			m.Cursor = 2
			return m.handleChoiceSelection(ChoiceProvision)
		case "4":
			m.Cursor = 3
			return m.handleChoiceSelection(ChoiceHoard)
		case "5":
			m.Cursor = 4
			return m.handleChoiceSelection(ChoiceSettings)
		case "6":
			m.Cursor = 5
			return m.handleChoiceSelection(ChoiceQuit)

		case "enter", "space":
			return m.handleChoiceSelection(Choice(m.Cursor))
		}
	}

	return m, nil
}

func (m Model) handleChoiceSelection(choice Choice) (tea.Model, tea.Cmd) {
	if choice == ChoiceQuest && m.CurrentPartySize <= 0 {
		m.Selected = -1
		m.StatusNotice = "⚠️ Cannot start a quest with 0 party members! Recruit heroes in Party screen first."
		return m, nil
	}

	m.Selected = choice
	return m, tea.Quit
}

func (m Model) View() string {
	if m.Quitted {
		return ""
	}

	// 3-tier hierarchy: dashboard screen -> dashboard component -> generic ui component
	headerComp := NewHeaderComponent(m.Theme)
	statusComp := NewStatusComponent(m.Theme, "Sacred Mirror / Dark Hearth", 0)
	menuComp := NewMenuComponent(m.Theme, m.Choices, m.Cursor)

	helpKeys := []key.Binding{
		key.NewBinding(key.WithKeys("up/down", "j/k"), key.WithHelp("↑/↓", "navigate")),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
		key.NewBinding(key.WithKeys("1-6"), key.WithHelp("1-6", "quick select")),
		key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
	}

	helpComp := NewHelpComponent(m.Theme, helpKeys)

	var sb strings.Builder
	sb.WriteString(headerComp.View())
	sb.WriteString(statusComp.View())
	sb.WriteString(menuComp.View())

	// Show party current/max size and coins right under options
	partyInfoStyle := lipgloss.NewStyle().
		MarginTop(1).
		Foreground(m.Theme.Secondary)
	partyInfoText := fmt.Sprintf("👥 Party: %d/%d Members   💰 Coins: %d Gold", m.CurrentPartySize, m.MaxPartySize, m.PartyGold)
	sb.WriteString("\n" + partyInfoStyle.Render(partyInfoText) + "\n")

	if m.StatusNotice != "" {
		sb.WriteString("\n" + lipgloss.NewStyle().Foreground(m.Theme.TextHighlight).Render(m.StatusNotice) + "\n")
	}

	sb.WriteString(helpComp.View())

	contentStyle := lipgloss.NewStyle().Align(lipgloss.Left)
	content := contentStyle.Render(sb.String())

	if m.Width > 0 && m.Height > 0 {
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}
