package provisions

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/pomquest/internal/adapters/icon"
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
	Cart         map[int]int
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
		Cart:       make(map[int]int),
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
		totalRows := len(m.Catalog) + 1 // items (0..6) + Purchase row (7)

		switch msg.String() {
		case "q", "esc", "b", "backspace", "ctrl+c":
			m.Quitted = true
			return m, nil

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			} else {
				m.Cursor = totalRows - 1
			}

		case "down", "j":
			if m.Cursor < totalRows-1 {
				m.Cursor++
			} else {
				m.Cursor = 0
			}

		case "+", "right", "l", "d":
			if m.Cursor < len(m.Catalog) {
				m.Cart[m.Cursor]++
			}

		case "-", "left", "h", "a":
			if m.Cursor < len(m.Catalog) {
				if m.Cart[m.Cursor] > 0 {
					m.Cart[m.Cursor]--
				}
			}

		case "enter", "space":
			if m.Cursor == len(m.Catalog) {
				m.executePurchase()
			} else {
				m.Cart[m.Cursor]++
			}

		default:
			// If key is a digit 0-9 and cursor is on an item row
			if len(msg.String()) == 1 && msg.String()[0] >= '0' && msg.String()[0] <= '9' {
				val, err := strconv.Atoi(msg.String())
				if err == nil {
					if m.Cursor < len(m.Catalog) {
						m.Cart[m.Cursor] = val
					} else if val >= 1 && val <= len(m.Catalog) {
						m.Cursor = val - 1
					}
				}
			}
		}
	}

	return m, nil
}

func (m Model) calculateGrandTotal() int {
	total := 0
	for idx, qty := range m.Cart {
		if idx >= 0 && idx < len(m.Catalog) && qty > 0 {
			total += qty * m.Catalog[idx].BaseCostCopper
		}
	}
	return total
}

func (m *Model) executePurchase() {
	grandTotal := m.calculateGrandTotal()
	if grandTotal <= 0 {
		m.StatusNotice = "⚠️ Cart is empty. Adjust item quantities to purchase supplies."
		return
	}

	partyCoins := app.TotalCoins(m.Party)
	if grandTotal > partyCoins {
		m.StatusNotice = fmt.Sprintf("⚠️ Insufficient coins! Required: %s, Purse: %s.", app.FormatCoins(grandTotal), app.FormatCoins(partyCoins))
		return
	}

	// Deduct coin values after purchase (24.8)
	newCoins := partyCoins - grandTotal
	app.SetPartyCoins(&m.Party, newCoins)

	// Add purchased items to party inventory
	for idx, qty := range m.Cart {
		if idx >= 0 && idx < len(m.Catalog) && qty > 0 {
			item := m.Catalog[idx]
			switch item.ID {
			case "food":
				m.Party.FoodLbs += qty
			case "spare_parts":
				m.Party.SpareParts += qty
			case "ox":
				m.Party.Oxen += qty * 3 // 1 yoke = 3 oxen
			case "potions":
				m.Party.Potions += qty
			case "bombs":
				m.Party.Bombs += qty
			case "shovels":
				m.Party.Shovels += qty
			case "rope":
				m.Party.Rope += qty
			}
		}
	}

	if m.DBPath != "" {
		_ = app.SavePartyToPath(m.DBPath, m.Party)
	}

	m.Cart = make(map[int]int)
	m.StatusNotice = fmt.Sprintf("✅ Purchase complete! Total spent: %s. Remaining Purse: %s.", app.FormatCoins(grandTotal), app.FormatCoins(newCoins))
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

	paneWidth := min(w-4, 78)
	if paneWidth < 40 {
		paneWidth = 40
	}

	// 24.6: Store items table with 3 left-aligned columns: item, cost, subtotal
	tableWidth := paneWidth - 4
	if tableWidth < 36 {
		tableWidth = 36
	}

	var rows [][]string

	for i, item := range m.Catalog {
		qty := m.Cart[i]
		itemIcon := icon.Get(item.IconKey, icon.ModeEmoji)
		itemCol := fmt.Sprintf("%s %s (%s)", itemIcon, item.Name, item.Unit)
		if qty > 0 {
			itemCol = fmt.Sprintf("%s %s (%s) x%d", itemIcon, item.Name, item.Unit, qty)
		}

		costCol := item.DisplayCost
		subtotalVal := item.BaseCostCopper * qty
		subtotalCol := app.FormatCoins(subtotalVal)

		if i == m.Cursor {
			itemCol = fmt.Sprintf("▶ %s", itemCol)
		} else {
			itemCol = fmt.Sprintf("  %s", itemCol)
		}

		rows = append(rows, []string{itemCol, costCol, subtotalCol})
	}

	// 24.8: Last option is Purchase row (item: purchase, cost: blank, subtotal: Total: <Grand Total>)
	grandTotal := m.calculateGrandTotal()
	purchaseCol := "🛒 Complete Purchase"
	if m.Cursor == len(m.Catalog) {
		purchaseCol = "▶ 🛒 Complete Purchase"
	} else {
		purchaseCol = "  🛒 Complete Purchase"
	}
	totalStr := fmt.Sprintf("Total: %s", app.FormatCoins(grandTotal))
	rows = append(rows, []string{purchaseCol, "", totalStr})

	emptyBorder := lipgloss.Border{}
	t := table.New().
		Border(emptyBorder).
		BorderRow(false).
		BorderColumn(false).
		Width(tableWidth).
		Headers("Item", "Cost", "Subtotal").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row < 0 { // Table Header (table.HeaderRow == -1)
				return lipgloss.NewStyle().Bold(true).Foreground(m.Theme.Headings).Align(lipgloss.Left).PaddingRight(2)
			}

			if row == m.Cursor {
				return lipgloss.NewStyle().Bold(true).Foreground(m.Theme.Primary).Align(lipgloss.Left).PaddingRight(2)
			}

			return lipgloss.NewStyle().Foreground(m.Theme.Text).Align(lipgloss.Left).PaddingRight(2)
		})

	var sb strings.Builder
	sb.WriteString(t.Render() + "\n")

	if m.StatusNotice != "" {
		noticeStyle := lipgloss.NewStyle().
			MarginTop(1).
			Italic(true).
			Foreground(m.Theme.DescriptionHighlight)
		sb.WriteString(noticeStyle.Render(m.StatusNotice) + "\n")
	}

	helpStyle := lipgloss.NewStyle().
		MarginTop(1).
		MarginBottom(1).
		Foreground(m.Theme.Muted)
	sb.WriteString(helpStyle.Render("[↑/↓/j/k] Select Row  [+/->/Enter] Qty+  [-/<-] Qty-  [q/esc/b] Back to Hearth") + "\n\n")

	// 24.4: Render shared party card component under store
	curSize := m.Party.CurrentPartySize
	if curSize == 0 && len(m.Party.Heroes) > 0 {
		for _, h := range m.Party.Heroes {
			if !h.IsDead {
				curSize++
			}
		}
	}
	maxSize := m.Party.MaxPartySize
	if maxSize <= 0 {
		maxSize = 1
	}

	partyCardData := ui.PartyCardData{
		Coins:            m.Party.Coins,
		FoodLbs:          m.Party.FoodLbs,
		Diet:             m.Party.Diet,
		SpareParts:       m.Party.SpareParts,
		Oxen:             m.Party.Oxen,
		Pace:             m.Party.Pace,
		Loot:             m.Party.Loot,
		Rope:             m.Party.Rope,
		Bombs:            m.Party.Bombs,
		Potions:          m.Party.Potions,
		Shovels:          m.Party.Shovels,
		CurrentPartySize: curSize,
		MaxPartySize:     maxSize,
	}
	cardWidth := paneWidth - 4
	if cardWidth < 36 {
		cardWidth = 36
	}
	partyCardView := ui.NewPartyCardComponent(m.Theme, partyCardData, cardWidth, rm).View()
	sb.WriteString(partyCardView)

	// 24.9: Remove border from store container
	containerStyle := lipgloss.NewStyle().
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
