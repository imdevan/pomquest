package settings

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

// ViewState tracks which sub-view is active on the Settings screen.
type ViewState int

const (
	StateMain ViewState = iota
	StateConfirmReset
)

// Choice represents an item selected on the settings menu.
type Choice int

const (
	ChoiceReset Choice = iota
	ChoiceSFX
	ChoiceMusic
	ChoiceCredits
	ChoiceBack
)

type settingsSelectMsg struct {
	choice Choice
}

// Model represents the Settings UI screen.
type Model struct {
	Theme           ui.Theme
	Responsive      *ui.ResponsiveManager
	DBPath          string
	Choices         []string
	Cursor          int
	FocusedIndex    int
	Selected        Choice
	State           ViewState
	ConfirmReset    ui.ConfirmationModel
	AudioSFXEnabled bool
	AudioMusicMuted bool
	StatusNotice    string
	Quitted         bool
	Width           int
	Height          int
}

// NewModel initializes the Settings screen with default options.
func NewModel(cfg domain.Config) Model {
	return NewModelWithDBPath(cfg, app.DefaultDBPath())
}

// NewModelWithDBPath initializes Settings screen with a custom SQLite database path.
func NewModelWithDBPath(cfg domain.Config, dbPath string) Model {
	theme := ui.ThemeFromConfig(cfg)
	m := Model{
		Theme:           theme,
		Responsive:      ui.NewResponsiveManager(80),
		DBPath:          dbPath,
		Cursor:          0,
		FocusedIndex:    -1,
		Selected:        -1,
		State:           StateMain,
		ConfirmReset:    ui.NewConfirmationModel("Reset Game Data?", "Are you sure you want to reset all game data and auto-migrate the database?", theme).WithBorderless(true),
		AudioSFXEnabled: true,
		AudioMusicMuted: true,
		Width:           80,
		Height:          24,
	}
	m.updateChoices()
	return m
}

func (m *Model) updateChoices() {
	m.Choices = []string{
		"1. ♻  Reset Game Data (SQL DB)",
		fmt.Sprintf("2. 🔊 Sound Effects: %s", sfxStatus(m.AudioSFXEnabled)),
		fmt.Sprintf("3. 🎵 Background Music: %s", musicStatus(m.AudioMusicMuted)),
		"4. 📜 Credits & License",
		"5. ⬅  Back to Hearth",
	}
}

// SetDimensions updates responsive dimensions for the settings view.
func (m *Model) SetDimensions(width, height int) {
	m.Width = width
	m.Height = height
	if m.Responsive != nil {
		m.Responsive.SetWidth(width)
	}
}

// Init clears the screen when entering the settings view.
func (m Model) Init() tea.Cmd {
	return tea.ClearScreen
}

// Update handles messages and input for the settings screen.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Responsive.SetWidth(msg.Width)
		return m, nil
	case settingsSelectMsg:
		m.FocusedIndex = -1
		return m.handleChoiceSelection(msg.choice)
	}

	switch m.State {
	case StateMain:
		return m.updateMain(msg)
	case StateConfirmReset:
		return m.updateConfirmReset(msg)
	}

	return m, nil
}

func (m Model) updateMain(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if m.FocusedIndex != -1 {
			return m, nil
		}
		menuComp := ui.NewMenuComponent(m.Theme, m.Choices, m.Cursor)
		menuComp.IDPrefix = "settings_menu"
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
		case "q", "esc", "b", "backspace", "ctrl+c":
			m.Quitted = true
			m.Selected = ChoiceBack
			return m, nil

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

		case "1", "r":
			m.Cursor = 0
			return m.triggerChoiceSelection(ChoiceReset)
		case "2", "s":
			m.Cursor = 1
			return m.triggerChoiceSelection(ChoiceSFX)
		case "3", "m":
			m.Cursor = 2
			return m.triggerChoiceSelection(ChoiceMusic)
		case "4", "c":
			m.Cursor = 3
			return m.triggerChoiceSelection(ChoiceCredits)
		case "5":
			m.Cursor = 4
			return m.triggerChoiceSelection(ChoiceBack)

		case "enter", "space":
			return m.triggerChoiceSelection(Choice(m.Cursor))
		}
	}
	return m, nil
}

func (m Model) triggerChoiceSelection(choice Choice) (tea.Model, tea.Cmd) {
	m.Cursor = int(choice)
	m.FocusedIndex = int(choice)
	m.Selected = choice
	return m, tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg {
		return settingsSelectMsg{choice: choice}
	})
}

func (m Model) handleChoiceSelection(choice Choice) (tea.Model, tea.Cmd) {
	m.Selected = choice
	switch choice {
	case ChoiceReset:
		m.State = StateConfirmReset
		m.ConfirmReset = ui.NewConfirmationModel(
			"Reset Game Data?",
			"Are you sure you want to reset all game data and auto-migrate the database?",
			m.Theme,
		).WithBorderless(true)
	case ChoiceSFX:
		m.AudioSFXEnabled = !m.AudioSFXEnabled
		m.StatusNotice = fmt.Sprintf("🔊 Sound Effects: %s", sfxStatus(m.AudioSFXEnabled))
		m.updateChoices()
	case ChoiceMusic:
		m.AudioMusicMuted = !m.AudioMusicMuted
		m.StatusNotice = fmt.Sprintf("🎵 Background Music: %s", musicStatus(m.AudioMusicMuted))
		m.updateChoices()
	case ChoiceCredits:
		m.StatusNotice = "📜 PomQuest • Built with Charm Bubble Tea & LipGloss • MIT License"
	case ChoiceBack:
		m.Quitted = true
	}
	return m, nil
}

func (m Model) updateConfirmReset(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := m.ConfirmReset.Update(msg)
	if updatedConfirm, ok := updated.(ui.ConfirmationModel); ok {
		m.ConfirmReset = updatedConfirm
		if cmd != nil {
			if _, isQuit := cmd().(tea.QuitMsg); isQuit {
				confirmed := m.ConfirmReset.ChoiceValue()
				if confirmed {
					_, err := app.ResetAndMigrateDB(m.DBPath)
					if err != nil {
						m.StatusNotice = fmt.Sprintf("❌ Failed to reset database: %v", err)
					} else {
						m.StatusNotice = "✅ Local SQLite DB reset and auto-migrated successfully!"
					}
				} else {
					m.StatusNotice = "Reset cancelled."
				}
				m.State = StateMain
				return m, nil
			}
		}
	}
	return m, cmd
}

func sfxStatus(enabled bool) string {
	if enabled {
		return "Enabled"
	}
	return "Disabled"
}

func musicStatus(muted bool) string {
	if muted {
		return "Muted"
	}
	return "Playing"
}

// View renders the Settings screen view.
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

	headerComp := NewHeaderComponent(m.Theme)

	helpKeys := []key.Binding{
		key.NewBinding(key.WithKeys("up/down", "j/k"), key.WithHelp("↑/↓", "navigate")),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
		key.NewBinding(key.WithKeys("1-5"), key.WithHelp("1-5", "quick select")),
		key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "back")),
	}

	helpComp := NewHelpComponent(m.Theme, helpKeys)

	if m.State == StateConfirmReset {
		var sb strings.Builder
		sb.WriteString(headerComp.View())
		sb.WriteString(m.ConfirmReset.View() + "\n")
		sb.WriteString(helpComp.View())

		contentStyle := lipgloss.NewStyle().Align(lipgloss.Left)
		content := contentStyle.Render(sb.String())

		if m.Width > 0 && m.Height > 0 {
			return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
		}
		return content
	}

	// Find the maximum line width of all other text on the page (sans margin)
	otherTexts := []string{
		headerComp.View(),
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
	sb.WriteString(menuComp.View() + "\n")

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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
