package path

import (
	"math"
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

// Model represents the path screen state and interaction.
type Model struct {
	Theme            ui.Theme
	Responsive       *ui.ResponsiveManager
	Cards            [3]ui.PathCardData
	Options          [3]app.Location
	ActiveIndex      int
	SelectedLocation app.Location
	Completed        bool
	Quitted          bool
	Width            int
	Height           int
	DBPath           string
	PartySize        int
	StatusNotice     string
	rng              *rand.Rand
}

// NewModel creates a new path screen instance using default DB path.
func NewModel(cfg domain.Config) Model {
	return NewModelWithDBPath(cfg, app.DefaultDBPath())
}

// NewModelWithDBPath creates a new path screen instance with a specified DB path.
func NewModelWithDBPath(cfg domain.Config, dbPath string) Model {
	theme := ui.ThemeFromConfig(cfg)
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	m := Model{
		Theme:       theme,
		Responsive:  ui.NewResponsiveManager(80),
		ActiveIndex: 0,
		Width:       80,
		Height:      24,
		DBPath:      dbPath,
		rng:         rng,
	}

	party, err := app.LoadActivePartyFromPath(dbPath)
	if err == nil {
		m.PartySize = party.CurrentPartySize
	} else {
		m.PartySize = 0
	}

	if m.PartySize <= 0 {
		m.StatusNotice = "⚠️ Cannot start a quest with 0 party members! Recruit heroes in Party screen first."
	}

	m.generateOptions()
	return m
}

func (m *Model) generateOptions() {
	biomeKeys := dungeon.AllBiomes

	generatedNames := make(map[string]bool)
	for i := 0; i < 3; i++ {
		name := dungeon.GenerateDungeonNameWithRand(m.rng)
		for generatedNames[name] {
			name = dungeon.GenerateDungeonNameWithRand(m.rng)
		}
		generatedNames[name] = true

		biomeKey := biomeKeys[(i+m.rng.Intn(len(biomeKeys)))%len(biomeKeys)]
		biomeObj := dungeon.GetBiome(biomeKey)

		diff := 1.0 + float64(i)*0.25
		loc := app.Location{
			Name:               name,
			Biome:              biomeObj.ID,
			Description:        biomeObj.Name + " environment with " + dungeon.GetBiomeHazardDescription(biomeObj),
			DifficultyModifier: diff,
		}
		m.Options[i] = loc

		m.Cards[i] = ui.PathCardData{
			Name:               loc.Name,
			Biome:              biomeObj.Name,
			Description:        loc.Description,
			DifficultyModifier: loc.DifficultyModifier,
			Icon:               icon.GetBiomeIcon(biomeObj.ID, icon.ModeNerdFont),
		}
	}
}

// SetDimensions sets the responsive layout dimensions.
func (m *Model) SetDimensions(width, height int) {
	m.Width = width
	m.Height = height
	if m.Responsive != nil {
		m.Responsive.SetWidth(width)
	}
}

// Init initializes bubbletea commands.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles key messages and state changes for the path screen.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.Quitted = true
			return m, nil

		case "left", "h", "up", "k":
			if m.ActiveIndex > 0 {
				m.ActiveIndex--
			} else {
				m.ActiveIndex = 2
			}

		case "right", "l", "down", "j":
			if m.ActiveIndex < 2 {
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

		case "r":
			m.generateOptions()

		case "enter", "space", " ":
			if m.PartySize <= 0 {
				m.StatusNotice = "⚠️ Cannot start a quest with 0 party members! Recruit heroes in Party screen first."
				return m, nil
			}

			selected := m.Options[m.ActiveIndex]
			savedLoc, err := app.SaveSelectedLocationToPath(m.DBPath, selected)
			if err == nil {
				m.SelectedLocation = savedLoc
			} else {
				m.SelectedLocation = selected
			}

			newQuest := app.Quest{
				Title:        m.SelectedLocation.Name,
				Biome:        m.SelectedLocation.Biome,
				PathLength:   1500,
				TotalSecs:    1500,
				ElapsedSecs:  0,
				CurrentRoom:  1,
				DungeonState: "Entering",
				Status:       "Active",
			}
			_, _ = app.SaveActiveQuestToPath(m.DBPath, newQuest)

			m.Completed = true
		}
	}

	return m, nil
}

// View renders the path screen UI.
func (m Model) View() string {
	if m.Quitted {
		return ""
	}

	w := m.Width
	h := m.Height
	if w <= 0 {
		w = 80
	}
	if h <= 0 {
		h = 24
	}

	rm := m.Responsive
	if rm == nil {
		rm = ui.NewResponsiveManager(w)
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.Theme.Primary).
		MarginBottom(1)

	header := titleStyle.Render("🗺️ CHOOSE YOUR PATH  •  Select a dungeon path to begin questing")

	isMDOrLarger := rm.IsAtLeast(ui.BreakpointMD)

	var cardWidth, cardHeight int
	if isMDOrLarger {
		// Overhead: Header (2), Notice (2 if present), Help (3), Margins/Padding (4)
		verticalOverhead := 9
		if m.StatusNotice != "" {
			verticalOverhead += 2
		}
		availH := h - verticalOverhead
		if availH < 10 {
			availH = 10
		}

		// Horizontal overhead for 3 cards: 2 gaps of 2 spaces (4) + outer margins (4)
		availW := w - 8
		if availW < 45 {
			availW = 45
		}
		maxCardWidth := (availW - 4) / 3

		// Aspect ratio 3:4 visual => cardWidth = round(1.5 * cardHeight)
		cardWidth = maxCardWidth
		cardHeight = int(math.Round(float64(cardWidth) / 1.5))

		if cardHeight > availH {
			cardHeight = availH
			cardWidth = int(math.Round(float64(cardHeight) * 1.5))
		}

		if cardWidth < 15 {
			cardWidth = 15
		}
		if cardHeight < 10 {
			cardHeight = 10
		}
	} else {
		// Stacked layout for smaller breakpoints (< MD)
		cardWidth = w - 4
		if cardWidth > 36 {
			cardWidth = 36
		}
		if cardWidth < 15 {
			cardWidth = 15
		}
		cardHeight = int(math.Round(float64(cardWidth) / 1.5))
		if cardHeight > 11 {
			cardHeight = 11
		}
		if cardHeight < 8 {
			cardHeight = 8
		}
	}

	var cards []string
	for i := 0; i < 3; i++ {
		c := m.Cards[i]
		c.IsActive = (i == m.ActiveIndex)
		cardComp := ui.NewPathCardComponentSized(m.Theme, c, cardWidth, cardHeight, rm)
		cards = append(cards, cardComp.View())
	}

	var grid string
	if isMDOrLarger {
		var spacerCards []string
		for i, card := range cards {
			spacerCards = append(spacerCards, card)
			if i < 2 {
				spacerCards = append(spacerCards, "  ")
			}
		}
		grid = lipgloss.JoinHorizontal(lipgloss.Top, spacerCards...)
	} else {
		var spacerCards []string
		for i, card := range cards {
			spacerCards = append(spacerCards, card)
			if i < 2 {
				spacerCards = append(spacerCards, "\n")
			}
		}
		grid = lipgloss.JoinVertical(lipgloss.Center, spacerCards...)
	}

	gridWidth := lipgloss.Width(grid)
	headerWidth := lipgloss.Width(header)
	contentWidth := w
	if contentWidth < gridWidth {
		contentWidth = gridWidth
	}
	if contentWidth < headerWidth {
		contentWidth = headerWidth
	}

	maxHelpWidth := contentWidth - 4
	if maxHelpWidth < gridWidth {
		maxHelpWidth = gridWidth
	}

	helpKeys := []key.Binding{
		key.NewBinding(key.WithKeys("←/→", "h/l"), key.WithHelp("←/→", "select path")),
		key.NewBinding(key.WithKeys("1-3"), key.WithHelp("1-3", "jump to path")),
		key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reroll paths")),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm path")),
		key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "cancel")),
	}
	helpComp := ui.NewHelpComponentSized(m.Theme, helpKeys, maxHelpWidth)

	centerStyle := lipgloss.NewStyle().Width(contentWidth).Align(lipgloss.Center)

	var sb strings.Builder
	sb.WriteString(centerStyle.Render(header) + "\n\n")

	if m.StatusNotice != "" {
		noticeStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(m.Theme.Primary).
			MarginBottom(1)
		sb.WriteString(centerStyle.Render(noticeStyle.Render(m.StatusNotice)) + "\n\n")
	}

	sb.WriteString(centerStyle.Render(grid) + "\n\n")
	sb.WriteString(centerStyle.Render(helpComp.View()))

	content := sb.String()

	if m.Width > 0 && m.Height > 0 {
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}
