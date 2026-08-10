package provisions

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/pomquest/internal/app"
	"github.com/pomquest/internal/domain"
	"github.com/pomquest/internal/ui"
)

// Model represents the provisions screen (Soand Sow's General Store).
type Model struct {
	Theme        ui.Theme
	Responsive   *ui.ResponsiveManager
	DBPath       string
	Party        app.Party
	Catalog      []app.ItemDef
	Cursor       int
	StatusNotice string
	Quitted      bool
	Width        int
	Height       int
}

// NewModel initializes the provisions store screen with active party data from default DB path.
func NewModel(cfg domain.Config) Model {
	return NewModelWithDBPath(cfg, app.DefaultDBPath())
}

// NewModelWithDBPath initializes the provisions store screen with active party data from custom DB path.
func NewModelWithDBPath(cfg domain.Config, dbPath string) Model {
	theme := ui.ThemeFromConfig(cfg)
	partyData, _ := app.LoadActivePartyFromPath(dbPath)

	return Model{
		Theme:      theme,
		Responsive: ui.NewResponsiveManager(80),
		DBPath:     dbPath,
		Party:      partyData,
		Catalog:    app.BookOfItems(),
		Cursor:     0,
		Width:      80,
		Height:     24,
	}
}

// SetDimensions sets responsive width and height for the provisions store.
func (m *Model) SetDimensions(width, height int) {
	m.Width = width
	m.Height = height
	if m.Responsive != nil {
		m.Responsive.SetWidth(width)
	}
}

// Init clears screen when entering provisions store.
func (m Model) Init() tea.Cmd {
	return tea.ClearScreen
}

// Update handles user key input and window resize events.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Responsive.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "b", "backspace", "ctrl+c":
			m.Quitted = true
			return m, nil

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			} else {
				m.Cursor = len(m.Catalog) - 1
			}

		case "down", "j":
			if m.Cursor < len(m.Catalog)-1 {
				m.Cursor++
			} else {
				m.Cursor = 0
			}

		case "1":
			m.Cursor = 0
			m.buyCurrentItem()
		case "2":
			m.Cursor = 1
			m.buyCurrentItem()
		case "3":
			m.Cursor = 2
			m.buyCurrentItem()
		case "4":
			m.Cursor = 3
			m.buyCurrentItem()
		case "5":
			m.Cursor = 4
			m.buyCurrentItem()
		case "6":
			m.Cursor = 5
			m.buyCurrentItem()
		case "7":
			m.Cursor = 6
			m.buyCurrentItem()

		case "enter", "space":
			m.buyCurrentItem()
		}
	}

	return m, nil
}

func (m *Model) buyCurrentItem() {
	if m.Cursor < 0 || m.Cursor >= len(m.Catalog) {
		return
	}
	item := m.Catalog[m.Cursor]

	totalCopper := m.Party.Gold*100 + m.Party.Silver*10 + m.Party.Copper
	if totalCopper < item.BaseCostCopper {
		m.StatusNotice = fmt.Sprintf("⚠️ Not enough coin for %s! Costs %s (Purse: %dg %ds %dc).", item.Name, item.DisplayCost, m.Party.Gold, m.Party.Silver, m.Party.Copper)
		return
	}

	remainingCopper := totalCopper - item.BaseCostCopper
	m.Party.Gold = remainingCopper / 100
	m.Party.Silver = (remainingCopper % 100) / 10
	m.Party.Copper = remainingCopper % 10

	switch item.ID {
	case "food":
		m.Party.FoodLbs += 1
	case "spare_parts":
		m.Party.SpareParts += 1
	case "ox":
		m.Party.Oxen += 3
	case "potions":
		m.Party.Potions += 1
	case "bombs":
		m.Party.Bombs += 1
	case "shovels":
		m.Party.Shovels += 1
	case "rope":
		m.Party.Rope += 1
	}

	if m.DBPath != "" {
		_ = app.SavePartyToPath(m.DBPath, m.Party)
	}

	m.StatusNotice = fmt.Sprintf("✅ Purchased %s (%s) for %s!", item.Name, item.Unit, item.DisplayCost)
}

// View renders the Provisions General Store UI screen.
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

	header := ui.NewHeaderComponent(m.Theme, "🛒 Soand Sow's General Store", "General Store • All items sold at base cost").View()

	var sb strings.Builder

	// Purse & Current Inventory bar
	purseStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.Theme.Headings).
		MarginBottom(1)
	sb.WriteString(purseStyle.Render(fmt.Sprintf("💰 Party Purse: %dg %ds %dc", m.Party.Gold, m.Party.Silver, m.Party.Copper)) + "\n")

	invStyle := lipgloss.NewStyle().
		Foreground(m.Theme.Secondary).
		MarginBottom(1)
	invText := fmt.Sprintf("📦 Owned: 🍞 %d lbs Food  ⚙ %d Parts  🐂 %d Oxen  🧪 %d Potions  💣 %d Bombs  ⛏ %d Shovels  🪢 %d Rope",
		m.Party.FoodLbs, m.Party.SpareParts, m.Party.Oxen, m.Party.Potions, m.Party.Bombs, m.Party.Shovels, m.Party.Rope)
	sb.WriteString(invStyle.Render(invText) + "\n\n")

	// Store catalog item list
	cursorStyle := lipgloss.NewStyle().Bold(true).Foreground(m.Theme.Primary)
	selectedStyle := lipgloss.NewStyle().Bold(true).Foreground(m.Theme.Primary)
	normalStyle := lipgloss.NewStyle().Foreground(m.Theme.Text)
	priceStyle := lipgloss.NewStyle().Foreground(m.Theme.Tags)
	descStyle := lipgloss.NewStyle().Foreground(m.Theme.Muted).Italic(true)

	for i, item := range m.Catalog {
		ownedStr := m.getOwnedCountString(item.ID)
		line := fmt.Sprintf("%d. %s %-18s  %-18s  %-12s", i+1, item.Icon, item.Name, priceStyle.Render(item.DisplayCost), descStyle.Render(ownedStr))

		if i == m.Cursor {
			sb.WriteString(fmt.Sprintf("%s %s\n", cursorStyle.Render("▶"), selectedStyle.Render(line)))
		} else {
			sb.WriteString(fmt.Sprintf("  %s\n", normalStyle.Render(line)))
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
	sb.WriteString("\n" + helpStyle.Render("[↑/↓/j/k] Navigate  [Enter/Space/1-7] Buy Item  [q/esc/b] Back to Hearth"))

	paneWidth := min(w-4, 78)
	if paneWidth < 40 {
		paneWidth = 40
	}

	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.Theme.Border).
		Padding(1, 2).
		Width(paneWidth)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		containerStyle.Render(sb.String()),
	)

	if m.Height > 0 {
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}

func (m Model) getOwnedCountString(itemID string) string {
	switch itemID {
	case "food":
		return fmt.Sprintf("Owned: %d lbs", m.Party.FoodLbs)
	case "spare_parts":
		return fmt.Sprintf("Owned: %d parts", m.Party.SpareParts)
	case "ox":
		return fmt.Sprintf("Owned: %d oxen", m.Party.Oxen)
	case "potions":
		return fmt.Sprintf("Owned: %d pots", m.Party.Potions)
	case "bombs":
		return fmt.Sprintf("Owned: %d bombs", m.Party.Bombs)
	case "shovels":
		return fmt.Sprintf("Owned: %d shovels", m.Party.Shovels)
	case "rope":
		return fmt.Sprintf("Owned: %d rope", m.Party.Rope)
	default:
		return ""
	}
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
