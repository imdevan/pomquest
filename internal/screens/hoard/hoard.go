package hoard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/pomquest/internal/adapters/ascii"
	"github.com/pomquest/internal/adapters/icon"
	"github.com/pomquest/internal/domain"
	"github.com/pomquest/internal/ui"
)

// HoardData holds overall statistics for the Dragon Hoard.
type HoardData struct {
	TotalGold       int
	TotalLoot       int
	TotalTreasure   int
	HeroesHired     int
	HeroesLost      int
	PlacesTraveled  int
	QuestsStarted   int
	QuestsCompleted int
	TimeFocusedMins int
}

// Model represents the Dragon Hoard screen.
type Model struct {
	Theme      ui.Theme
	Responsive *ui.ResponsiveManager
	Data       HoardData
	Width      int
	Height     int
	Quitted    bool
}

// NewModel creates a new Hoard screen model initialized with default demo values.
func NewModel(cfg domain.Config) Model {
	theme := ui.ThemeFromConfig(cfg)
	return Model{
		Theme:      theme,
		Responsive: ui.NewResponsiveManager(80),
		Data: HoardData{
			TotalGold:       1250,
			TotalLoot:       67,
			TotalTreasure:   14,
			HeroesHired:     12,
			HeroesLost:      2,
			PlacesTraveled:  8,
			QuestsStarted:   15,
			QuestsCompleted: 11,
			TimeFocusedMins: 340,
		},
		Width:  80,
		Height: 24,
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
		case "q", "esc", "ctrl+c", "b", "backspace", "enter", "space":
			m.Quitted = true
			return m, nil
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

	// Small screen fallback check
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

	// 1. Top Section: ASCII Art Container
	asciiArt := ascii.Get(ascii.KeyTreasureHoard)
	croppedArt := cropAsciiArt(asciiArt, w)
	asciiContainer := lipgloss.NewStyle().
		Width(w).
		Align(lipgloss.Center).
		Render(lipgloss.NewStyle().Foreground(m.Theme.Secondary).Render(croppedArt))

	// 2. Bottom Section: Textured Card Container with 3-column table
	goldIcon := icon.GetNerdFont(icon.KeyCoins)
	lootIcon := icon.GetNerdFont(icon.KeyLoot)
	treasureIcon := icon.GetNerdFont(icon.KeyTreasure)
	hiredIcon := icon.GetNerdFont(icon.KeyHeroesHired)
	lostIcon := icon.GetNerdFont(icon.KeyHeroesLost)
	travelIcon := icon.GetNerdFont(icon.KeyPlacesTraveled)
	startedIcon := icon.GetNerdFont(icon.KeyQuestsStarted)
	completedIcon := icon.GetNerdFont(icon.KeyQuestsCompleted)
	timeIcon := icon.GetNerdFont(icon.KeyTimeFocused)

	var goldStr, lootStr, treasureStr string
	var hiredStr, lostStr, travelStr string
	var startedStr, completedStr, timeStr string

	if !rm.IsAtLeast(ui.BreakpointLG) {
		// On MD and smaller: show icon: value
		goldStr = fmt.Sprintf("%s: %sg", goldIcon, formatNumber(m.Data.TotalGold))
		lootStr = fmt.Sprintf("%s: %d", lootIcon, m.Data.TotalLoot)
		treasureStr = fmt.Sprintf("%s: %d", treasureIcon, m.Data.TotalTreasure)

		hiredStr = fmt.Sprintf("%s: %d", hiredIcon, m.Data.HeroesHired)
		lostStr = fmt.Sprintf("%s: %d", lostIcon, m.Data.HeroesLost)
		travelStr = fmt.Sprintf("%s: %d", travelIcon, m.Data.PlacesTraveled)

		startedStr = fmt.Sprintf("%s: %d", startedIcon, m.Data.QuestsStarted)
		completedStr = fmt.Sprintf("%s: %d", completedIcon, m.Data.QuestsCompleted)
		timeStr = fmt.Sprintf("%s: %dh %dm", timeIcon, m.Data.TimeFocusedMins/60, m.Data.TimeFocusedMins%60)
	} else {
		// On LG+: show full label names
		goldStr = fmt.Sprintf("%s total gold: %sg", goldIcon, formatNumber(m.Data.TotalGold))
		lootStr = fmt.Sprintf("%s total loot: %d", lootIcon, m.Data.TotalLoot)
		treasureStr = fmt.Sprintf("%s total treasure: %d", treasureIcon, m.Data.TotalTreasure)

		hiredStr = fmt.Sprintf("%s heroes hired: %d", hiredIcon, m.Data.HeroesHired)
		lostStr = fmt.Sprintf("%s heroes lost: %d", lostIcon, m.Data.HeroesLost)
		travelStr = fmt.Sprintf("%s places traveled: %d", travelIcon, m.Data.PlacesTraveled)

		startedStr = fmt.Sprintf("%s quests started: %d", startedIcon, m.Data.QuestsStarted)
		completedStr = fmt.Sprintf("%s quests completed: %d", completedIcon, m.Data.QuestsCompleted)
		timeStr = fmt.Sprintf("%s time focused: %dh %dm", timeIcon, m.Data.TimeFocusedMins/60, m.Data.TimeFocusedMins%60)
	}

	padH := 2
	cardMargin := 1
	if !rm.IsAtLeast(ui.BreakpointLG) {
		padH = 1
		cardMargin = 0
	}

	cardWidth := min(w, 80)
	if !rm.IsAtLeast(ui.BreakpointLG) {
		cardWidth = w
	}

	innerTableWidth := cardWidth - 2 - (2 * padH)
	if innerTableWidth < 20 {
		innerTableWidth = 20
	}

	emptyBorder := lipgloss.Border{}
	t := table.New().
		Border(emptyBorder).
		BorderRow(false).
		BorderColumn(false).
		Width(innerTableWidth).
		Rows(
			[]string{goldStr, hiredStr, startedStr},
			[]string{lootStr, lostStr, completedStr},
			[]string{treasureStr, travelStr, timeStr},
		).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Foreground(m.Theme.Text).Align(lipgloss.Left)
		})

	card := ui.NewTexturedCard(m.Theme, cardWidth, t.Render()).Padding(1, padH)
	card.TopLeft = ui.CornerText{
		Text:         "Hoard",
		Margin:       cardMargin,
		PaddingLeft:  1,
		PaddingRight: 1,
	}
	cardView := card.View()
	cardContainer := lipgloss.NewStyle().
		Width(w).
		Align(lipgloss.Center).
		Render(cardView)

	// 3. Help component Container
	helpKeys := []key.Binding{
		key.NewBinding(key.WithKeys("q", "esc", "b"), key.WithHelp("q/esc", "back to menu")),
	}
	helpComp := ui.NewHelpComponent(m.Theme, helpKeys)
	helpContainer := lipgloss.NewStyle().
		Width(w).
		Align(lipgloss.Center).
		Render(helpComp.View())

	var sb strings.Builder
	sb.WriteString(asciiContainer)
	sb.WriteString("\n\n")
	sb.WriteString(cardContainer)
	sb.WriteString("\n")
	sb.WriteString(helpContainer)

	content := sb.String()

	if m.Width > 0 && m.Height > 0 {
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}

func formatNumber(n int) string {
	in := fmt.Sprintf("%d", n)
	out := make([]byte, 0, len(in)+(len(in)-1)/3)
	for i, c := range in {
		if i > 0 && (len(in)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, byte(c))
	}
	return string(out)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func cropAsciiArt(art string, maxWidth int) string {
	if maxWidth <= 0 {
		return art
	}
	lines := strings.Split(art, "\n")
	var cropped []string
	for _, line := range lines {
		runes := []rune(line)
		if len(runes) > maxWidth {
			excess := len(runes) - maxWidth
			start := excess / 2
			end := start + maxWidth
			if end > len(runes) {
				end = len(runes)
			}
			line = string(runes[start:end])
		}
		cropped = append(cropped, line)
	}
	return strings.Join(cropped, "\n")
}
