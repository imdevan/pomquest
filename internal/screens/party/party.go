package party

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/pomquest/internal/adapters/icon"
	"github.com/pomquest/internal/adapters/mouse"
	"github.com/pomquest/internal/app"
	"github.com/pomquest/internal/domain"
	"github.com/pomquest/internal/ui"
)

// Model represents the dedicated Party Screen UI.
type Model struct {
	Theme             ui.Theme
	Responsive        *ui.ResponsiveManager
	Heroes            []ui.HeroCardData
	MaxPartySize      int
	CurrentPartySize  int
	ActiveIndex       int
	PressedHeroIndex  int
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
		Theme:            theme,
		Responsive:       ui.NewResponsiveManager(80),
		Heroes:           heroes,
		MaxPartySize:     len(heroes),
		CurrentPartySize: len(heroes),
		ActiveIndex:      2, // pam is default active
		PressedHeroIndex: -1,
		Width:            80,
		Height:           24,
	}
}

// NewModelFromParty initializes the party screen model based on SQL database Party model.
func NewModelFromParty(cfg domain.Config, party app.Party) Model {
	theme := ui.ThemeFromConfig(cfg)
	maxSize := party.MaxPartySize
	if maxSize <= 0 {
		maxSize = 1
	}

	var heroes []ui.HeroCardData
	for i, h := range party.Heroes {
		if i >= maxSize {
			break
		}
		iconStr := h.Icon
		if iconStr == "" {
			iconStr = app.GetClassArchetype(h.Class).Icon(icon.ModeNerdFont)
		}
		heroes = append(heroes, ui.HeroCardData{
			Name:          h.Name,
			Class:         h.Class,
			Level:         h.Level,
			HP:            h.HP,
			MaxHP:         h.MaxHP,
			Attack:        h.Attack,
			Armor:         h.Armor,
			Status:        h.Status,
			IsDead:        h.IsDead,
			Icon:          iconStr,
			Stats:         h.Stats,
			ShowFullStats: true,
		})
	}

	// Fill remaining slots up to party.MaxPartySize with empty hero slots
	for len(heroes) < maxSize {
		heroes = append(heroes, ui.HeroCardData{
			IsEmptySlot:   true,
			ShowFullStats: true,
		})
	}

	return Model{
		Theme:            theme,
		Responsive:       ui.NewResponsiveManager(80),
		Heroes:           heroes,
		MaxPartySize:     maxSize,
		CurrentPartySize: len(party.Heroes),
		ActiveIndex:      0,
		PressedHeroIndex: -1,
		Width:            80,
		Height:           24,
	}
}

func (m *Model) ReplaceHero(index int, newHero ui.HeroCardData) {
	if index >= 0 && index < len(m.Heroes) {
		newHero.ShowFullStats = true
		newHero.IsEmptySlot = false
		m.Heroes[index] = newHero

		var count int
		for _, h := range m.Heroes {
			if !h.IsEmptySlot && !h.IsDead {
				count++
			}
		}
		m.CurrentPartySize = count
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

type navToRerollMsg struct {
	TargetIndex int
}

func navToRerollCmd(targetIndex int) tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return navToRerollMsg{TargetIndex: targetIndex}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case navToRerollMsg:
		m.PressedHeroIndex = -1
		m.NavToReroll = true
		m.RerollTargetIndex = msg.TargetIndex
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Responsive.SetWidth(msg.Width)
		return m, nil

	case tea.MouseMsg:
		evt := mouse.ParseEvent(msg)
		for i := range m.Heroes {
			zoneID := fmt.Sprintf("party_hero_%d", i)
			if mouse.InBounds(zoneID, msg) {
				m.ActiveIndex = i

				if evt.Type == mouse.EventClick {
					m.PressedHeroIndex = i
					if m.Heroes[i].IsDead || m.Heroes[i].IsEmptySlot {
						return m, navToRerollCmd(i)
					}
				}
			}
		}
		if evt.Type == mouse.EventRelease {
			m.PressedHeroIndex = -1
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.Quitted = true
			return m, tea.Quit

		case "r", "enter", "space":
			if m.ActiveIndex >= 0 && m.ActiveIndex < len(m.Heroes) && (m.Heroes[m.ActiveIndex].IsDead || m.Heroes[m.ActiveIndex].IsEmptySlot) {
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

	w := m.Width
	if w <= 0 {
		w = 80
	}

	rm := m.Responsive
	if rm == nil {
		rm = ui.NewResponsiveManager(w)
	}

	if rm.IsAtMost(ui.BreakpointSM) {
		msgStyle := lipgloss.NewStyle().
			Foreground(m.Theme.Text).
			Align(lipgloss.Center)

		paneStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(m.Theme.Border).
			Padding(1, 2).
			Width(max(w-2, 10)).
			Align(lipgloss.Center)

		content := paneStyle.Render(msgStyle.Render("This terminal width is not currently supported."))
		if m.Height > 0 {
			return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
		}
		return content
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.Theme.Primary).
		MarginBottom(1)

	header := titleStyle.Render("👥 PARTY MEMBERS & HERO STATS")

	var cards []string
	for i, h := range m.Heroes {
		h.IsActive = (i == m.ActiveIndex)
		h.IsActiveFocused = (i == m.PressedHeroIndex)
		cardComp := ui.NewHeroCardComponent(m.Theme, h, m.Responsive)
		cardComp.ZoneID = fmt.Sprintf("party_hero_%d", i)
		cards = append(cards, cardComp.View())
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
	if m.ActiveIndex >= 0 && m.ActiveIndex < len(m.Heroes) {
		if m.Heroes[m.ActiveIndex].IsEmptySlot {
			helpKeys = append(helpKeys, key.NewBinding(key.WithKeys("r/enter"), key.WithHelp("r", "recruit hero")))
		} else if m.Heroes[m.ActiveIndex].IsDead {
			helpKeys = append(helpKeys, key.NewBinding(key.WithKeys("r/enter"), key.WithHelp("r", "reroll dead hero")))
		}
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
