package mouse

import (
	"testing"
	"time"

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

func TestZoneHelpers(t *testing.T) {
	InitManager()

	marked := Mark("btn1", "Click Me")
	if marked == "Click Me" {
		t.Errorf("expected marked string to contain bubblezone ANSI identifiers")
	}

	scanned := Scan(marked)
	if scanned != "Click Me" {
		t.Errorf("expected scanned string to equal 'Click Me', got %q", scanned)
	}

	time.Sleep(20 * time.Millisecond)

	zInfo := Get("btn1")
	if zInfo == nil {
		t.Errorf("expected zone info for 'btn1'")
	} else if zInfo.StartX != 0 || zInfo.StartY != 0 {
		t.Errorf("expected start coords (0,0), got (%d,%d)", zInfo.StartX, zInfo.StartY)
	}

	msgIn := tea.MouseMsg{X: 0, Y: 0, Button: tea.MouseButtonLeft}
	if !InBounds("btn1", msgIn) {
		t.Errorf("expected click at (0,0) to be inside 'btn1'")
	}

	msgOut := tea.MouseMsg{X: 20, Y: 20, Button: tea.MouseButtonLeft}
	if InBounds("btn1", msgOut) {
		t.Errorf("expected click at (20,20) to be outside 'btn1'")
	}
}

