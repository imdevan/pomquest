package dashboard

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/pomquest/internal/app"
	"github.com/pomquest/internal/domain"
)

func TestDashboardScreen_Navigation(t *testing.T) {
	cfg := domain.DefaultConfig()
	model := NewModel(cfg)

	// Verify initial rendering contains title and options
	view := model.View()
	if !strings.Contains(view, "POMQUEST - HEARTH DASHBOARD") {
		t.Errorf("expected view to contain header title, got:\n%s", view)
	}
	if !strings.Contains(view, "Quest!") {
		t.Errorf("expected view to contain Quest option, got:\n%s", view)
	}

	// Test navigation down
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m := updated.(Model)
	if m.Cursor != 1 {
		t.Errorf("expected cursor to be 1 after pressing 'j', got %d", m.Cursor)
	}
}

func TestDashboardScreen_PartyDisplayAndCoins(t *testing.T) {
	cfg := domain.DefaultConfig()
	party := app.Party{
		Name:             "Heroes Guild",
		MaxPartySize:     4,
		Coins:            150,
		CurrentPartySize: 2,
		Heroes: []app.Hero{
			{Creature: app.Creature{Noun: app.Noun{Name: "Gimli"}}, Class: "barbarian"},
			{Creature: app.Creature{Noun: app.Noun{Name: "Legolas"}}, Class: "ranger"},
		},
	}

	model := NewModelFromParty(cfg, party)
	view := model.View()

	// 28.2: show party current and max size on dashboard under options
	if !strings.Contains(view, "👥 Party: 2/4 Members") {
		t.Errorf("expected view to display party size '👥 Party: 2/4 Members', got:\n%s", view)
	}

	// 28.3: show party coins to right of party hero count
	if !strings.Contains(view, "💰 Coins: 1s 50c") {
		t.Errorf("expected view to display party coins '💰 Coins: 1s 50c', got:\n%s", view)
	}
}

func TestDashboardScreen_CannotQuestWithZeroPartyMembers(t *testing.T) {
	cfg := domain.DefaultConfig()
	model := NewModel(cfg) // CurrentPartySize is 0

	// Attempt to select "1. Quest!" with 0 party members
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	m := updated.(Model)

	// 28.4: do not allow players with 0 current party members to "Quest!"
	if m.Selected == ChoiceQuest {
		t.Errorf("expected Selected to NOT be ChoiceQuest when party size is 0, got ChoiceQuest")
	}
	if cmd != nil {
		t.Errorf("expected nil cmd (no quit) when Quest selection is blocked, got %v", cmd)
	}
	if !strings.Contains(m.StatusNotice, "Cannot start a quest with 0 party members") {
		t.Errorf("expected status notice warning about 0 party members, got %q", m.StatusNotice)
	}

	// Now update party data to have 1 member
	party := app.Party{
		MaxPartySize: 4,
		Heroes: []app.Hero{
			{Creature: app.Creature{Noun: app.Noun{Name: "Aragorn"}}, Class: "paladin"},
		},
	}
	m.UpdatePartyData(party)

	// Select "1. Quest!" now with 1 party member
	updated, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	m = updated.(Model)
	if m.Selected != ChoiceQuest {
		t.Errorf("expected Selected to be ChoiceQuest when party size > 0, got %v", m.Selected)
	}
	if cmd == nil {
		t.Errorf("expected quit cmd to transition screen when Quest is allowed")
	}
}

func TestDashboardScreen_MenuCardStyleAndFocusedState(t *testing.T) {
	cfg := domain.DefaultConfig()
	model := NewModel(cfg)

	view := model.View()

	// 1. Should not contain cursor pointer '▶'
	if strings.Contains(view, "▶") {
		t.Errorf("expected dashboard view to not contain cursor pointer '▶', got:\n%s", view)
	}

	// 2. Should contain rounded borders for menu cards
	if !strings.Contains(view, "╭") || !strings.Contains(view, "╰") {
		t.Errorf("expected dashboard view to contain bordered menu cards, got:\n%s", view)
	}

	// 3. Test focused state on selection
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	m := updated.(Model)
	if m.FocusedIndex != 1 {
		t.Errorf("expected FocusedIndex to be 1 on quick select '2', got %d", m.FocusedIndex)
	}
	if cmd == nil {
		t.Fatalf("expected delay transition cmd on quick select, got nil")
	}

	// Executing transition cmd delivers menuSelectMsg which quits
	quitModel, quitCmd := m.Update(cmd())
	qm := quitModel.(Model)
	if qm.FocusedIndex != -1 {
		t.Errorf("expected FocusedIndex reset to -1 after transition, got %d", qm.FocusedIndex)
	}
	if quitCmd == nil {
		t.Errorf("expected tea.Quit command after transition")
	}
}

func TestDashboardScreen_MainMenuWidthMatchesOtherText(t *testing.T) {
	cfg := domain.DefaultConfig()
	model := NewModel(cfg)
	model.SetDimensions(80, 24)

	// In default view, find max width of other text lines (sans margin)
	headerComp := NewHeaderComponent(model.Theme)
	statusComp := NewStatusComponent(model.Theme, "Sacred Mirror / Dark Hearth", 0)
	partyInfoText := "👥 Party: 0/1 Members   💰 Coins: 10s 00c"

	otherTexts := []string{
		headerComp.View(),
		statusComp.View(),
		partyInfoText,
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

	view := model.View()
	lines := strings.Split(view, "\n")

	// Find rendered menu card border top line `╭───────...╮`
	var cardBorderWidth int
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "╭") && strings.Contains(trimmed, "╮") {
			// Measure width of front card box (sans margin)
			endIdx := strings.Index(trimmed, "╮")
			cardBox := trimmed[:endIdx+len("╮")]
			cardBorderWidth = lipgloss.Width(cardBox)
			break
		}
	}

	if cardBorderWidth < maxOtherWidth {
		t.Errorf("expected main menu card width (%d) to be at least as wide as other text on page (%d)", cardBorderWidth, maxOtherWidth)
	}

	// Also verify with status notice when width is wide enough
	model.SetDimensions(100, 24)
	model.StatusNotice = "⚠️ Cannot start a quest with 0 party members! Recruit heroes in Party screen first."
	noticeWidth := lipgloss.Width(model.StatusNotice)
	viewWithNotice := model.View()
	linesWithNotice := strings.Split(viewWithNotice, "\n")

	var cardBorderWidthWithNotice int
	for _, l := range linesWithNotice {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "╭") && strings.Contains(trimmed, "╮") {
			endIdx := strings.Index(trimmed, "╮")
			cardBox := trimmed[:endIdx+len("╮")]
			cardBorderWidthWithNotice = lipgloss.Width(cardBox)
			break
		}
	}

	if cardBorderWidthWithNotice < noticeWidth {
		t.Errorf("expected main menu card width with notice (%d) to expand to at least notice width (%d)", cardBorderWidthWithNotice, noticeWidth)
	}
}

