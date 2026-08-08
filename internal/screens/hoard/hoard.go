package hoard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/pomquest/internal/adapters/icon"
	"github.com/pomquest/internal/domain"
	"github.com/pomquest/internal/ui"
)

// TreasureHoardASCII contains the ASCII art representation for the dragon hoard screen (from plan/ascii-art.md:575).
const TreasureHoardASCII = `⠀⢀⠀⠀⣀⡤⠤⠚⠳⢤⣀⣀⣀⡀⠀⠀⠀⣠⣎⣁⠀⢀⡀⠀⠄⠂⠀⠀⢀⠤⠀⠀⠑⢝⢟⣮⣣⡀⠀⠀⣠⠄⠢⡀⠀⠈⠺⣒⢮⣑⢄⢺⣦⡀⠢⢄⣰⠑⢄⡀⢀⠀⢀⠀⠀
⠀⡇⢀⡞⠁⠀⠀⢀⣀⣀⣀⣀⠀⠑⣄⣆⡼⠉⠁⢀⣉⣙⢆⣀⡀⠀⠀⢀⣯⠀⣠⣰⠗⠉⠑⣝⡙⣿⣦⠊⠀⠀⠀⠸⡑⢤⡀⠈⠳⣌⣓⠓⣷⣽⣦⣀⣉⣢⣀⣉⣙⡄⠀⡗⠀
⡰⢣⠎⠀⠀⣠⡴⠚⢿⠀⠀⠀⡝⢤⡈⠁⢠⠋⠑⡞⠀⠜⠈⢣⣱⣄⣰⣋⣜⡤⠛⠛⡉⠭⠉⠛⠻⣾⣧⡏⠉⠭⠭⣭⣙⢿⣛⣫⣤⣌⣻⣷⣌⠪⣻⣆⣋⣉⣩⠄⣥⣹⣤⡇⠀
⢀⡜⠀⢀⣼⠛⢦⣀⠬⠷⣶⠴⢅⢀⠟⢆⢰⢳⡖⢃⣼⠀⠐⢢⢧⣴⣶⣢⡟⠛⢻⣛⠟⠉⠉⠉⠉⠉⠉⠉⡍⠉⠉⠉⠉⠉⢫⠙⠉⣉⣉⠉⢹⠁⣇⣿⠀⡀⣀⡯⢀⠀⢹⠸⡄
⢈⡇⢀⣾⢇⣀⡼⠁⠀⠀⠀⠉⣎⠙⣄⢈⣏⠉⠀⢨⣿⡆⢸⠉⣾⣽⣿⣿⠀⣠⠏⢠⣤⠀⡀⢀⣤⠀⠀⠀⠁⠀⠀⠀⠀⠀⠀⠉⡶⢥⣤⣠⢸⢰⡿⢹⢸⣀⣼⠀⣻⣷⠘⣇⢣
⠎⠀⢸⡇⠀⣸⠀⠀⠀⠀⠀⠀⠘⡆⠘⡁⢸⡄⠀⢸⣿⡟⠈⠉⢹⣿⣿⣿⠑⠁⠀⡏⢀⣿⠛⣶⠙⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⠤⠴⡧⢤⠘⣼⡇⣿⠉⢀⢿⠀⢿⣻⡇⠘⡆
⠀⡤⢼⡏⠉⡇⠀⠀⠀⠀⠀⠀⠀⡗⠚⡟⠛⡇⢸⣿⣿⠁⠀⠀⢸⢹⡿⣿⠄⢄⠛⣿⡟⠓⠛⠛⣿⡟⠀⠀⠀⠀⠀⠀⠀⠀⢀⠤⠐⠂⢌⣉⢰⠈⡵⣼⠉⡹⠀⢣⠈⢻⣷⣶⠘
⡔⠛⣿⡇⠀⡇⠀⠀⠀⠀⠀⠀⠀⡧⠤⡧⠤⡇⠘⢻⡇⣀⣸⣉⣿⡟⣧⣿⢩⠋⠒⣬⠥⡄⠀⠀⠈⠁⠀⠀⠀⢠⣤⡤⠒⠊⠁⠀⠀⠀⠀⠀⠙⠢⣇⣿⠛⠁⣀⡠⢷⡀⢻⣿⠀
⢇⣀⣸⡗⠒⡇⠀⠀⠀⠀⠀⠀⠀⣇⣀⣇⣀⡧⢤⡰⠇⠀⠧⢼⣇⡏⣟⠘⠓⣾⠉⣇⠁⠀⠀⣀⣼⡿⣟⣛⣿⣛⣻⠿⠶⣶⣦⡀⠀⠀⠀⣤⠶⠬⣽⣶⣎⠁⠀⠀⢳⣒⠭⠆
⠀⠀⢘⠃⠀⡇⠀⠀⠀⠀⠀⢀⡠⡇⠈⡇⠀⡧⠼⢒⣳⡖⣦⣠⡞⣏⣏⡏⣴⠶⣿⣉⣉⠗⠒⠈⠀⣿⣷⣿⣴⣿⣒⣒⣒⣸⣿⢿⡇⠀⠀⠴⢟⣿⢾⡿⣿⡈⠁⠀⠀⠀⠈⠑⠤
⠀⠀⠸⡧⠐⡇⠀⠀⠀⣀⡴⠊⠀⡏⠑⡖⠊⣇⣶⢿⣿⠃⠙⡏⢸⠟⢫⣻⢿⢛⡿⠟⠁⢀⣠⣾⡿⣯⠳⢾⣷⣿⣛⠟⢦⣾⣿⡟⠁⠀⠀⠀⠀⠀⠉⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⢰⡇⠀⡿⣩⠖⠋⠀⠀⠀⠀⡗⠂⣇⠒⠇⢻⠀⠘⣦⠤⠓⠛⠒⠚⠉⠺⢅⠀⠀⠀⡿⣟⡛⠿⠿⠷⣶⣻⡮⣿⣾⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀
⠤⢄⡘⡏⠀⢱⠀⠈⠉⠑⠢⢤⣀⣇⣀⡏⡠⠔⠉⠓⠚⠁⠀⠀⠀⠀⠀⠀⠀⣀⠤⠒⠊⣿⣧⣤⠉⣒⣤⡀⠒⠛⣿⣻⢿⡏⠑⢄⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⠴⠊⠉⠀
⠀⢀⣽⡷⠶⠼⠶⢒⣀⣈⡭⢽⣦⣒⣊⣉⢩⣤⠴⠒⢒⣒⠦⠤⢄⣠⠤⠖⠉⠀⠀⠀⠀⠈⠉⠉⠛⠻⠶⠿⢿⣿⣿⠿⠋⠁⠀⠀⠈⠳⢤⡀⠀⣀⣀⣀⠀⠤⠤⠾⠷⣦⡀⠀⠀
⠀⣾⡏⠀⠀⠀⢠⣺⠟⠛⢻⣶⡈⠀⠀⢀⣠⠤⠝⠿⢥⣤⠠⠦⣍⣀⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣼⡟⠘⣿⢧⠀⠀⠀⠀⠘⡿⡄⠀
⠀⣿⣿⣿⣶⣦⣼⡿⣴⣾⣿⣿⠈⠉⠛⠦⢄⡀⠀⠀⠀⠀⠀⢀⡠⠄⠒⠉⠉⠉⣳⠤⠤⣄⣀⣀⣀⣀⣤⣀⣀⠀⠀⠀⠀⠀⠀⣀⠤⠒⠊⠉⣿⣿⣶⡼⡇⠦⣔⣲⠾⣿⣿⡇⠀
⠀⠸⣇⡀⠈⠙⢻⡿⠋⠁⠀⣿⠀⠀⠀⠀⠀⠈⠉⠙⠓⠒⠋⢄⡀⢀⣀⣠⣞⣉⡀⠀⣀⣠⠴⣮⣅⣀⣀⠀⠀⠀⣀⡭⠓⠋⠉⠀⠀⠀⠀⠀⣿⠀⠉⢻⡟⠛⠋⠉⠀⢰⣿⣁⠀
⠀⠀⠀⠈⠁⠲⢬⣹⣆⡴⠋⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⠯⣁⡀⠀⢉⣩⡿⠛⠉⠀⠀⣀⡠⠬⠵⠛⠛⠉⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠢⣼⣷⡠⠔⠒⠊⠉⠉⠀⠀`

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

	// 1. Top Section: ASCII Art
	asciiStyle := lipgloss.NewStyle().
		Foreground(m.Theme.Secondary).
		Align(lipgloss.Center)
	asciiView := asciiStyle.Render(TreasureHoardASCII)

	// 2. Bottom Section: Textured Card with 3-column table
	goldIcon := icon.GetNerdFont(icon.KeyCoins)
	lootIcon := icon.GetNerdFont(icon.KeyLoot)
	treasureIcon := icon.GetNerdFont(icon.KeyTreasure)
	hiredIcon := icon.GetNerdFont(icon.KeyHeroesHired)
	lostIcon := icon.GetNerdFont(icon.KeyHeroesLost)
	travelIcon := icon.GetNerdFont(icon.KeyPlacesTraveled)
	startedIcon := icon.GetNerdFont(icon.KeyQuestsStarted)
	completedIcon := icon.GetNerdFont(icon.KeyQuestsCompleted)
	timeIcon := icon.GetNerdFont(icon.KeyTimeFocused)

	// Formatted values
	goldStr := fmt.Sprintf("%s total gold: %sg", goldIcon, formatNumber(m.Data.TotalGold))
	lootStr := fmt.Sprintf("%s total loot: %d", lootIcon, m.Data.TotalLoot)
	treasureStr := fmt.Sprintf("%s total treasure: %d", treasureIcon, m.Data.TotalTreasure)

	hiredStr := fmt.Sprintf("%s heroes hired: %d", hiredIcon, m.Data.HeroesHired)
	lostStr := fmt.Sprintf("%s heroes lost: %d", lostIcon, m.Data.HeroesLost)
	travelStr := fmt.Sprintf("%s places traveled: %d", travelIcon, m.Data.PlacesTraveled)

	startedStr := fmt.Sprintf("%s quests started: %d", startedIcon, m.Data.QuestsStarted)
	completedStr := fmt.Sprintf("%s quests completed: %d", completedIcon, m.Data.QuestsCompleted)
	timeStr := fmt.Sprintf("%s time focused: %dh %dm", timeIcon, m.Data.TimeFocusedMins/60, m.Data.TimeFocusedMins%60)

	cardWidth := min(w, 80)
	if cardWidth < 50 {
		cardWidth = 50
	}

	innerTableWidth := cardWidth - 6
	if innerTableWidth < 30 {
		innerTableWidth = 30
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

	card := ui.NewTexturedCard(m.Theme, cardWidth, t.Render()).Padding(1, 2)
	card.TopLeft = ui.CornerText{
		Text:         "Hoard",
		Margin:       1,
		PaddingLeft:  1,
		PaddingRight: 1,
	}
	cardView := card.View()

	// 3. Help component
	helpKeys := []key.Binding{
		key.NewBinding(key.WithKeys("q", "esc", "b"), key.WithHelp("q/esc", "back to menu")),
	}
	helpComp := ui.NewHelpComponent(m.Theme, helpKeys)
	helpView := helpComp.View()

	var sb strings.Builder
	sb.WriteString(asciiView)
	sb.WriteString("\n\n")
	sb.WriteString(cardView)
	sb.WriteString("\n")
	sb.WriteString(helpView)

	contentStyle := lipgloss.NewStyle().Align(lipgloss.Center)
	content := contentStyle.Render(sb.String())

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
