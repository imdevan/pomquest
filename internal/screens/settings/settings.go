package settings

import (
	"fmt"
	"strings"

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
	StateAudio
	StateCredits
)

// Model represents the Settings UI screen.
type Model struct {
	Theme           ui.Theme
	Responsive      *ui.ResponsiveManager
	DBPath          string
	Choices         []string
	Cursor          int
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
	return Model{
		Theme:      theme,
		Responsive: ui.NewResponsiveManager(80),
		DBPath:     dbPath,
		Choices: []string{
			"1. ♻  Reset Game Data (SQL DB)",
			"2. 🔊 Audio Settings (Stub)",
			"3. 📜 Credits & License (Stub)",
			"4. ⬅  Back to Hearth",
		},
		Cursor:          0,
		State:           StateMain,
		ConfirmReset:    ui.NewConfirmationModel("Reset Game Data?", "Are you sure you want to reset all game data and auto-migrate the database?", theme).WithBorderless(true),
		AudioSFXEnabled: true,
		AudioMusicMuted: true,
		Width:           80,
		Height:          24,
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
	}

	switch m.State {
	case StateMain:
		return m.updateMain(msg)
	case StateConfirmReset:
		return m.updateConfirmReset(msg)
	case StateAudio:
		return m.updateAudio(msg)
	case StateCredits:
		return m.updateCredits(msg)
	}

	return m, nil
}

func (m Model) updateMain(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "b", "backspace", "ctrl+c":
			m.Quitted = true
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

		case "1":
			m.Cursor = 0
			m.selectCurrentChoice()
		case "2":
			m.Cursor = 1
			m.selectCurrentChoice()
		case "3":
			m.Cursor = 2
			m.selectCurrentChoice()
		case "4":
			m.Cursor = 3
			m.selectCurrentChoice()

		case "enter", "space":
			m.selectCurrentChoice()
		}
	}
	return m, nil
}

func (m *Model) selectCurrentChoice() {
	switch m.Cursor {
	case 0: // Reset Game Data
		m.State = StateConfirmReset
		m.ConfirmReset = ui.NewConfirmationModel(
			"Reset Game Data?",
			"Are you sure you want to reset all game data and auto-migrate the database?",
			m.Theme,
		).WithBorderless(true)
	case 1: // Audio
		m.State = StateAudio
	case 2: // Credits
		m.State = StateCredits
	case 3: // Back
		m.Quitted = true
	}
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

func (m Model) updateAudio(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "b", "backspace":
			m.State = StateMain
			m.StatusNotice = "Audio settings saved (stub)."
			return m, nil
		case "1", "s":
			m.AudioSFXEnabled = !m.AudioSFXEnabled
			m.StatusNotice = fmt.Sprintf("Sound Effects: %s (stub)", sfxStatus(m.AudioSFXEnabled))
		case "2", "m":
			m.AudioMusicMuted = !m.AudioMusicMuted
			m.StatusNotice = fmt.Sprintf("Music: %s (stub)", musicStatus(m.AudioMusicMuted))
		case "enter", "space":
			m.AudioSFXEnabled = !m.AudioSFXEnabled
			m.StatusNotice = fmt.Sprintf("Sound Effects: %s (stub)", sfxStatus(m.AudioSFXEnabled))
		}
	}
	return m, nil
}

func (m Model) updateCredits(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "b", "backspace", "enter", "space":
			m.State = StateMain
			return m, nil
		}
	}
	return m, nil
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

	header := ui.NewHeaderComponent(m.Theme, "⚙ Settings & Configuration", "Manage local database, audio preferences, and view credits").View()

	var body string
	switch m.State {
	case StateMain:
		body = m.viewMain(w)
	case StateConfirmReset:
		body = m.ConfirmReset.View()
	case StateAudio:
		body = m.viewAudio(w)
	case StateCredits:
		body = m.viewCredits(w)
	}

	paneWidth := min(w-4, 76)
	if paneWidth < 40 {
		paneWidth = 40
	}

	containerStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Width(paneWidth)

	if m.State != StateConfirmReset {
		containerStyle = containerStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(m.Theme.Border)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		containerStyle.Render(body),
	)

	if m.Height > 0 {
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}

func (m Model) viewMain(width int) string {
	var sb strings.Builder

	menuTitleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.Theme.Headings).
		MarginBottom(1)
	sb.WriteString(menuTitleStyle.Render("Configuration Menu") + "\n\n")

	cursorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.Theme.Primary)
	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.Theme.Primary)
	normalStyle := lipgloss.NewStyle().
		Foreground(m.Theme.Text)

	for i, choice := range m.Choices {
		if i == m.Cursor {
			sb.WriteString(fmt.Sprintf("%s %s\n", cursorStyle.Render("▶"), selectedStyle.Render(choice)))
		} else {
			sb.WriteString(fmt.Sprintf("  %s\n", normalStyle.Render(choice)))
		}
	}

	if m.StatusNotice != "" {
		noticeStyle := lipgloss.NewStyle().
			MarginTop(1).
			Italic(true).
			Foreground(m.Theme.DescriptionHighlight)
		sb.WriteString("\n" + noticeStyle.Render(m.StatusNotice))
	}

	helpStyle := lipgloss.NewStyle().
		MarginTop(1).
		Foreground(m.Theme.Muted)
	sb.WriteString("\n" + helpStyle.Render("[↑/↓/j/k] Navigate  [Enter/Space/1-4] Select  [q/esc/b] Back"))

	return sb.String()
}

func (m Model) viewAudio(width int) string {
	var sb strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.Theme.Headings).
		MarginBottom(1)
	sb.WriteString(titleStyle.Render("🔊 Audio Settings (Stub)"))

	infoStyle := lipgloss.NewStyle().Foreground(m.Theme.Text)
	highlightStyle := lipgloss.NewStyle().Bold(true).Foreground(m.Theme.Primary)

	sb.WriteString("\n\n" + infoStyle.Render("Configure in-game sound effects and background music:\n\n"))

	sfxLabel := fmt.Sprintf("1. Sound Effects: [%s]", sfxStatus(m.AudioSFXEnabled))
	musicLabel := fmt.Sprintf("2. Music Track:   [%s]", musicStatus(m.AudioMusicMuted))
	synthLabel := "3. Sound Engine:  [Retro Synthesizer Stub]"

	sb.WriteString(highlightStyle.Render(sfxLabel) + "\n")
	sb.WriteString(highlightStyle.Render(musicLabel) + "\n")
	sb.WriteString(infoStyle.Render(synthLabel) + "\n")

	if m.StatusNotice != "" {
		noticeStyle := lipgloss.NewStyle().
			MarginTop(1).
			Italic(true).
			Foreground(m.Theme.DescriptionHighlight)
		sb.WriteString("\n" + noticeStyle.Render(m.StatusNotice))
	}

	helpStyle := lipgloss.NewStyle().
		MarginTop(1).
		Foreground(m.Theme.Muted)
	sb.WriteString("\n" + helpStyle.Render("[1/s/Space] Toggle SFX  [2/m] Toggle Music  [q/esc/b] Back"))

	return sb.String()
}

func (m Model) viewCredits(width int) string {
	var sb strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.Theme.Headings)
	sb.WriteString(titleStyle.Render("📜 Credits & Acknowledgments"))

	bodyStyle := lipgloss.NewStyle().Foreground(m.Theme.Text)
	accentStyle := lipgloss.NewStyle().Foreground(m.Theme.Secondary).Bold(true)

	creditsText := `
⚔ PomQuest - Terminal Pomodoro RPG ⚔
------------------------------------
A command-line RPG driven by pomodoro focus loops.

Development & Architecture:
  • Lead Developer: PomQuest Adventurers
  • Terminal Engine: Charm Bubble Tea & LipGloss
  • Database: SQLite & GORM Persistence

Special Thanks:
  • ASCII Tavern Masters & Nerd Fonts
  • Open Source CLI Community

License: MIT
`
	sb.WriteString(bodyStyle.Render(creditsText))

	helpStyle := lipgloss.NewStyle().
		MarginTop(1).
		Foreground(m.Theme.Muted)
	sb.WriteString(accentStyle.Render("Press Esc, Enter, or 'b' to return to Settings."))
	sb.WriteString("\n" + helpStyle.Render("[q/esc/b/Enter] Back"))

	return sb.String()
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
