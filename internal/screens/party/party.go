package party

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/pomquest/internal/adapters/icon"
	"github.com/pomquest/internal/app"
	"github.com/pomquest/internal/domain"
	"github.com/pomquest/internal/ui"
)

// Model represents the dedicated Party Screen UI.
type Model struct {
	Theme             ui.Theme
	Responsive        *ui.ResponsiveManager
	Heroes            []ui.HeroCardData
	ActiveIndex       int
	NavToReroll       bool
	RerollTargetIndex int
	Quitted           bool
	Width             int
	Height            int
}

// NewModel initializes the party screen model with default party roster data.
func NewModel(cfg domain.Config) Model {
	theme := ui.ThemeFromConfig(cfg)
	heroes := []ui.HeroCardData{
		{
			Name:          "jimmothy",
			Class:         "fighter",
			Level:         2,
			HP:            10,
			MaxHP:         10,
			Attack:        2,
			Armor:         2,
			Status:        "Healthy",
			Icon:          app.GetClassArchetype("fighter").Icon(icon.ModeNerdFont),
			Stats:         app.Stats{Str: 2, Dex: 1, Cha: 0, Sur: 0, Con: 2, Wil: 0, Int: 0, Per: 0, Hel: 2, Luc: 0},
			ShowFullStats: true,
		},
		{
			Name:          "applebee",
			Class:         "wizard",
			Level:         1,
			HP:            6,
			MaxHP:         6,
			Attack:        3,
			Armor:         0,
			Status:        "Healthy",
			Icon:          app.GetClassArchetype("wizard").Icon(icon.ModeNerdFont),
			Stats:         app.Stats{Str: 0, Dex: 1, Cha: 0, Sur: 0, Con: 0, Wil: 1, Int: 3, Per: 1, Hel: 0, Luc: 1},
			ShowFullStats: true,
		},
		{
			Name:          "pam",
			Class:         "ranger",
			Level:         1,
			HP:            6,
			MaxHP:         8,
			Attack:        3,
			Armor:         0,
			Status:        "Healthy",
			Icon:          app.GetClassArchetype("ranger").Icon(icon.ModeNerdFont),
			Stats:         app.Stats{Str: 0, Dex: 1, Cha: 0, Sur: 1, Con: 0, Wil: 0, Int: 0, Per: 1, Hel: 1, Luc: 0},
			ShowFullStats: true,
		},
		{
			Name:          "sam",
			Class:         "rogue",
			Level:         2,
			HP:            0,
			MaxHP:         8,
			Attack:        3,
			Armor:         0,
			Status:        "Dead",
			IsDead:        true,
			Icon:          app.GetClassArchetype("rogue").Icon(icon.ModeNerdFont),
			Stats:         app.Stats{Str: 0, Dex: 3, Cha: 1, Sur: 0, Con: 0, Wil: 0, Int: 0, Per: 1, Hel: 0, Luc: 1},
			ShowFullStats: true,
		},
	}

	return Model{
		Theme:       theme,
		Responsive:  ui.NewResponsiveManager(80),
		Heroes:      heroes,
		ActiveIndex: 2, // pam is default active
		Width:       80,
		Height:      24,
	}
}

func (m *Model) ReplaceHero(index int, newHero ui.HeroCardData) {
	if index >= 0 && index < len(m.Heroes) {
		m.Heroes[index] = newHero
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
			return m, tea.Quit

		case "r", "enter", "space":
			if m.ActiveIndex >= 0 && m.ActiveIndex < len(m.Heroes) && m.Heroes[m.ActiveIndex].IsDead {
				m.NavToReroll = true
				m.RerollTargetIndex = m.ActiveIndex
				return m, tea.Quit
			}

		case "left", "h", "up", "k":
			if m.ActiveIndex > 0 {
				m.ActiveIndex--
			} else {
				m.ActiveIndex = len(m.Heroes) - 1
			}

		case "right", "l", "down", "j":
			if m.ActiveIndex < len(m.Heroes)-1 {
				m.ActiveIndex++
			} else {
				m.ActiveIndex = 0
			}

		case "1":
			if len(m.Heroes) > 0 {
				m.ActiveIndex = 0
			}
		case "2":
			if len(m.Heroes) > 1 {
				m.ActiveIndex = 1
			}
		case "3":
			if len(m.Heroes) > 2 {
				m.ActiveIndex = 2
			}
		case "4":
			if len(m.Heroes) > 3 {
				m.ActiveIndex = 3
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.Quitted {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.Theme.Primary).
		MarginBottom(1)

	header := titleStyle.Render("👥 PARTY MEMBERS & HERO STATS")

	var cards []string
	for i, h := range m.Heroes {
		h.IsActive = (i == m.ActiveIndex)
		cardView := ui.NewHeroCardComponent(m.Theme, h, m.Responsive).View()
		cards = append(cards, cardView)
	}

	var spacerCards []string
	for i, card := range cards {
		spacerCards = append(spacerCards, card)
		if i < len(cards)-1 {
			spacerCards = append(spacerCards, "  ")
		}
	}

	grid := lipgloss.JoinHorizontal(lipgloss.Top, spacerCards...)

	helpKeys := []key.Binding{
		key.NewBinding(key.WithKeys("←/→", "h/l"), key.WithHelp("←/→", "select hero")),
	}
	if m.ActiveIndex >= 0 && m.ActiveIndex < len(m.Heroes) && m.Heroes[m.ActiveIndex].IsDead {
		helpKeys = append(helpKeys, key.NewBinding(key.WithKeys("r/enter"), key.WithHelp("r", "reroll dead hero")))
	}
	helpKeys = append(helpKeys,
		key.NewBinding(key.WithKeys("1-4"), key.WithHelp("1-4", "jump hero")),
		key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "back to hearth")),
	)

	helpComp := ui.NewHelpComponent(m.Theme, helpKeys)

	var sb strings.Builder
	sb.WriteString(header + "\n\n")
	sb.WriteString(grid + "\n\n")
	sb.WriteString(helpComp.View())

	contentStyle := lipgloss.NewStyle().Align(lipgloss.Left)
	content := contentStyle.Render(sb.String())

	if m.Width > 0 && m.Height > 0 {
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}
