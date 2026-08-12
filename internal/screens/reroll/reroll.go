package reroll

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
	"github.com/pomquest/internal/namegen"
	"github.com/pomquest/internal/ui"
)

// Model represents the reroll screen state and interactions.
type Model struct {
	Theme         ui.Theme
	Responsive    *ui.ResponsiveManager
	TargetHeroIdx int
	Cards         [4]ui.HeroCardData
	ActiveIndex   int
	RerollsLeft   int
	EnteringName  bool
	NameInput     string
	SelectedArch  app.ClassArchetype
	Completed     bool
	NewHero       ui.HeroCardData
	Quitted       bool
	Width         int
	Height        int
	rng           *rand.Rand
}

// NewModel creates a new reroll screen instance.
func NewModel(cfg domain.Config, targetHeroIdx int) Model {
	theme := ui.ThemeFromConfig(cfg)
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	m := Model{
		Theme:         theme,
		Responsive:    ui.NewResponsiveManager(80),
		TargetHeroIdx: targetHeroIdx,
		ActiveIndex:   0,
		RerollsLeft:   2,
		Width:         80,
		Height:        24,
		rng:           rng,
	}

	m.generateCards()
	return m
}

func (m *Model) generateCards() {
	catalog := app.ClassCatalog()

	// Card 0: Random Card with unknown stats
	m.Cards[0] = ui.HeroCardData{
		Name:           "?????",
		Class:          "random",
		Level:          1,
		Icon:           icon.GetNerdFont(icon.KeyDice),
		ShowFullStats:  true,
		IsUnknownStats: true,
	}

	// Pick 3 random distinct classes from catalog for cards 1, 2, 3
	perm := m.rng.Perm(len(catalog))
	for i := 0; i < 3; i++ {
		arch := catalog[perm[i]]
		m.Cards[i+1] = ui.HeroCardData{
			Name:          namegen.GenerateCharacterNameWithRand(m.rng),
			Class:         arch.Name,
			Level:         1,
			HP:            arch.HP,
			MaxHP:         arch.HP,
			Attack:        arch.Attack,
			Armor:         arch.Armor,
			Status:        "Healthy",
			Icon:          arch.Icon(icon.ModeNerdFont),
			Stats:         arch.Stats,
			ShowFullStats: true,
		}
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
	catalog := app.ClassCatalog()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Responsive.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		if m.EnteringName {
			switch msg.String() {
			case "esc":
				m.EnteringName = false
				return m, nil
			case "enter":
				name := strings.TrimSpace(m.NameInput)
				if name == "" {
					name = "Hero"
				}
				m.NewHero = ui.HeroCardData{
					Name:          name,
					Class:         m.SelectedArch.Name,
					Level:         1,
					HP:            m.SelectedArch.HP,
					MaxHP:         m.SelectedArch.HP,
					Attack:        m.SelectedArch.Attack,
					Armor:         m.SelectedArch.Armor,
					Status:        "Healthy",
					Icon:          m.SelectedArch.Icon(icon.ModeNerdFont),
					Stats:         m.SelectedArch.Stats,
					ShowFullStats: true,
					IsDead:        false,
				}
				m.Completed = true
				return m, tea.Quit
			case "backspace":
				if len(m.NameInput) > 0 {
					m.NameInput = m.NameInput[:len(m.NameInput)-1]
				}
			default:
				if len(msg.String()) == 1 && len(m.NameInput) < 12 {
					m.NameInput += msg.String()
				}
			}
			return m, nil
		}

		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.Quitted = true
			return m, tea.Quit

		case "r":
			if m.RerollsLeft > 0 {
				m.RerollsLeft--
				// Reroll cards 1, 2, 3
				perm := m.rng.Perm(len(catalog))
				for i := 0; i < 3; i++ {
					arch := catalog[perm[i]]
					m.Cards[i+1] = ui.HeroCardData{
						Name:          namegen.GenerateCharacterNameWithRand(m.rng),
						Class:         arch.Name,
						Level:         1,
						HP:            arch.HP,
						MaxHP:         arch.HP,
						Attack:        arch.Attack,
						Armor:         arch.Armor,
						Status:        "Healthy",
						Icon:          arch.Icon(icon.ModeNerdFont),
						Stats:         arch.Stats,
						ShowFullStats: true,
					}
				}
			}

		case "left", "h", "up", "k":
			if m.ActiveIndex > 0 {
				m.ActiveIndex--
			} else {
				m.ActiveIndex = 3
			}

		case "right", "l", "down", "j":
			if m.ActiveIndex < 3 {
				m.ActiveIndex++
			} else {
				m.ActiveIndex = 0
			}

		case "1":
			m.ActiveIndex = 0
		case "2":
			m.ActiveIndex = 1
		case "3":
			m.ActiveIndex = 2
		case "4":
			m.ActiveIndex = 3

		case "enter", "space", " ":
			if m.ActiveIndex == 0 {
				// Random card chosen
				rndIdx := m.rng.Intn(len(catalog))
				m.SelectedArch = catalog[rndIdx]
				m.NameInput = namegen.GenerateCharacterNameWithRand(m.rng)
			} else {
				// Find matching archetype from selected card
				cardClass := m.Cards[m.ActiveIndex].Class
				m.SelectedArch = app.GetClassArchetype(cardClass)
				m.NameInput = m.Cards[m.ActiveIndex].Name
			}
			m.EnteringName = true
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

	subTitleStyle := lipgloss.NewStyle().
		Foreground(m.Theme.Tags)

	headerText := fmt.Sprintf("🎲 REROLL HERO CLASS  •  Rerolls left: %d/2", m.RerollsLeft)
	header := titleStyle.Render(headerText)

	var cards []string
	for i := 0; i < 4; i++ {
		h := m.Cards[i]
		h.IsActive = (i == m.ActiveIndex)
		cardView := ui.NewHeroCardComponent(m.Theme, h, m.Responsive).View()
		cards = append(cards, cardView)
	}

	var spacerCards []string
	for i, card := range cards {
		spacerCards = append(spacerCards, card)
		if i < 3 {
			spacerCards = append(spacerCards, "  ")
		}
	}

	grid := lipgloss.JoinHorizontal(lipgloss.Top, spacerCards...)

	var sb strings.Builder
	sb.WriteString(header + "\n\n")
	sb.WriteString(grid + "\n\n")

	if m.EnteringName {
		classIcon := m.SelectedArch.Icon(icon.ModeNerdFont)
		promptTitle := fmt.Sprintf("SELECTED CLASS: %s %s", classIcon, strings.ToUpper(m.SelectedArch.Name))
		boxStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(m.Theme.Primary).
			Padding(1, 2).
			Width(50)

		inputDisplay := m.NameInput + "█"
		promptContent := subTitleStyle.Render(promptTitle) + "\n\n" +
			"Enter Hero Name: " + lipgloss.NewStyle().Bold(true).Foreground(m.Theme.Primary).Render(inputDisplay) + "\n\n" +
			lipgloss.NewStyle().Foreground(m.Theme.Muted).Render("[Enter] Confirm Hero  •  [Esc] Cancel")

		sb.WriteString(boxStyle.Render(promptContent) + "\n\n")
	} else {
		helpKeys := []key.Binding{
			key.NewBinding(key.WithKeys("←/→", "h/l"), key.WithHelp("←/→", "select class")),
			key.NewBinding(key.WithKeys("r"), key.WithHelp("r", fmt.Sprintf("reroll choices (%d/2)", m.RerollsLeft))),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "choose class")),
			key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "cancel")),
		}
		helpComp := ui.NewHelpComponent(m.Theme, helpKeys)
		sb.WriteString(helpComp.View())
	}

	contentStyle := lipgloss.NewStyle().Align(lipgloss.Left)
	content := contentStyle.Render(sb.String())

	if m.Width > 0 && m.Height > 0 {
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}
