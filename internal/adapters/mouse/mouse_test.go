package mouse

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestParseEventScrollUp(t *testing.T) {
	msg := tea.MouseMsg{
		X:      10,
		Y:      5,
		Button: tea.MouseButtonWheelUp,
	}

	evt := ParseEvent(msg)
	if evt.Type != EventScrollUp {
		t.Errorf("expected EventScrollUp, got %v", evt.Type)
	}
	if evt.X != 10 || evt.Y != 5 {
		t.Errorf("expected X=10, Y=5, got X=%d, Y=%d", evt.X, evt.Y)
	}
}

func TestParseEventScrollDown(t *testing.T) {
	msg := tea.MouseMsg{
		X:      10,
		Y:      5,
		Button: tea.MouseButtonWheelDown,
	}

	evt := ParseEvent(msg)
	if evt.Type != EventScrollDown {
		t.Errorf("expected EventScrollDown, got %v", evt.Type)
	}
}

func TestEnableMouseCmd(t *testing.T) {
	cmd := EnableMouseCmd()
	if cmd == nil {
		t.Errorf("expected non-nil tea.Cmd from EnableMouseCmd")
	}
}
