package dashboard

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/pomquest/internal/domain"
	"github.com/pomquest/internal/ui"
)

// Choice represents an item selected on the dashboard menu.
type Choice int

const (
	ChoiceQuest Choice = iota
	ChoiceHoard
	ChoiceProvision
	ChoiceSettings
	ChoiceQuit
)

// Model represents the main menu / hearth dashboard UI screen.
type Model struct {
	Theme        ui.Theme
	Responsive   *ui.ResponsiveManager
	Choices      []string
	Cursor       int
	Selected     Choice
	Quitted      bool
	StatusNotice string
	Width        int
	Height       int
}

// NewModel initializes a dashboard UI with default menu items.
func NewModel(cfg domain.Config) Model {
	theme := ui.ThemeFromConfig(cfg)
	return Model{
		Theme:      theme,
		Responsive: ui.NewResponsiveManager(80),
		Choices: []string{
			"1. ⚔  Quest! (Start Journey)",
			"2. 💎 View Dragon Hoard",
			"3. 🎒 Provisions & Party",
			"4. ⚙  Settings & Config",
			"5. 🚪 Quit",
		},
		Cursor:   0,
		Selected: -1,
		Width:    80,
		Height:   24,
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
			m.Selected = ChoiceQuest
			return m, tea.Quit
		case "2":
			m.Cursor = 1
			m.Selected = ChoiceHoard
			return m, tea.Quit
		case "3":
			m.Cursor = 2
			m.Selected = ChoiceProvision
			return m, tea.Quit
		case "4":
			m.Cursor = 3
			m.Selected = ChoiceSettings
			return m, tea.Quit
		case "5":
			m.Cursor = 4
			m.Selected = ChoiceQuit
			return m, tea.Quit

		case "enter", "space":
			m.Selected = Choice(m.Cursor)
			return m, tea.Quit
		}
	}

	return m, nil
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
		key.NewBinding(key.WithKeys("1-5"), key.WithHelp("1-5", "quick select")),
		key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
	}

	helpComp := NewHelpComponent(m.Theme, helpKeys)

	var sb strings.Builder
	sb.WriteString(headerComp.View())
	sb.WriteString(statusComp.View())
	sb.WriteString(menuComp.View())

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
