package dashboard

import (
	"fmt"
	"strings"
	"time"

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

type menuSelectMsg struct {
	choice Choice
}

// Model represents the main menu / hearth dashboard UI screen.
type Model struct {
	Theme            ui.Theme
	Responsive       *ui.ResponsiveManager
	Choices          []string
	Cursor           int
	FocusedIndex     int
	Selected         Choice
	CurrentPartySize int
	MaxPartySize     int
	PartyCoins       int
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
		FocusedIndex:     -1,
		Selected:         -1,
		CurrentPartySize: 0,
		MaxPartySize:     1,
		PartyCoins:       1000,
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
	m.PartyCoins = app.TotalCoins(party)
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

	case menuSelectMsg:
		m.FocusedIndex = -1
		return m.handleChoiceSelection(msg.choice)

	case tea.MouseMsg:
		if m.FocusedIndex != -1 {
			return m, nil
		}
		menuComp := ui.NewMenuComponent(m.Theme, m.Choices, m.Cursor)
		menuComp.IDPrefix = "dashboard_menu"
		if newCursor, clicked := menuComp.HandleMouse(msg); clicked {
			m.Cursor = newCursor
			return m.triggerChoiceSelection(Choice(newCursor))
		} else if newCursor != m.Cursor {
			m.Cursor = newCursor
		}
		return m, nil

	case tea.KeyMsg:
		if m.FocusedIndex != -1 {
			return m, nil
		}
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
			return m.triggerChoiceSelection(ChoiceQuest)
		case "2":
			m.Cursor = 1
			return m.triggerChoiceSelection(ChoiceParty)
		case "3":
			m.Cursor = 2
			return m.triggerChoiceSelection(ChoiceProvision)
		case "4":
			m.Cursor = 3
			return m.triggerChoiceSelection(ChoiceHoard)
		case "5":
			m.Cursor = 4
			return m.triggerChoiceSelection(ChoiceSettings)
		case "6":
			m.Cursor = 5
			return m.triggerChoiceSelection(ChoiceQuit)

		case "enter", "space":
			return m.triggerChoiceSelection(Choice(m.Cursor))
		}
	}

	return m, nil
}

func (m Model) triggerChoiceSelection(choice Choice) (tea.Model, tea.Cmd) {
	if choice == ChoiceQuest && m.CurrentPartySize <= 0 {
		m.Selected = -1
		m.StatusNotice = "⚠️ Cannot start a quest with 0 party members! Recruit heroes in Party screen first."
		return m, nil
	}

	m.FocusedIndex = int(choice)
	m.Selected = choice
	return m, tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg {
		return menuSelectMsg{choice: choice}
	})
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

	helpKeys := []key.Binding{
		key.NewBinding(key.WithKeys("up/down", "j/k"), key.WithHelp("↑/↓", "navigate")),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
		key.NewBinding(key.WithKeys("1-6"), key.WithHelp("1-6", "quick select")),
		key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
	}

	helpComp := NewHelpComponent(m.Theme, helpKeys)

	partyInfoText := fmt.Sprintf("👥 Party: %d/%d Members   💰 Coins: %s", m.CurrentPartySize, m.MaxPartySize, app.FormatCoins(m.PartyCoins))

	// Find the maximum line width of all other text on the page (sans margin)
	otherTexts := []string{
		headerComp.View(),
		statusComp.View(),
		partyInfoText,
		helpComp.View(),
	}
	if m.StatusNotice != "" {
		otherTexts = append(otherTexts, m.StatusNotice)
	}

	maxOtherWidth := 0
	for _, txt := range otherTexts {
		for _, line := range strings.Split(txt, "\n") {
			w := lipgloss.Width(strings.TrimSpace(line))
			if w > maxOtherWidth {
				maxOtherWidth = w
			}
		}
	}

	menuWidth := maxOtherWidth
	if m.Width > 0 && m.Width-4 < menuWidth {
		menuWidth = m.Width - 4
	}
	if menuWidth < 30 {
		menuWidth = 30
	}
	menuComp := NewMenuComponent(m.Theme, m.Choices, m.Cursor).WithFocusedIndex(m.FocusedIndex).WithWidth(menuWidth)

	var sb strings.Builder
	sb.WriteString(headerComp.View())
	sb.WriteString(statusComp.View())
	sb.WriteString(menuComp.View() + "\n")

	// Show party current/max size and coins right under options
	partyInfoStyle := lipgloss.NewStyle().
		MarginTop(1).
		Foreground(m.Theme.Secondary)
	sb.WriteString(partyInfoStyle.Render(partyInfoText) + "\n")

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
