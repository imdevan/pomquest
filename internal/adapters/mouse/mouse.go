package mouse

import tea "github.com/charmbracelet/bubbletea"

// EventType represents classified mouse interaction types.
type EventType int

const (
	EventNone EventType = iota
	EventScrollUp
	EventScrollDown
	EventClick
	EventRelease
	EventMotion
)

// Event represents a parsed mouse interaction with coordinates and type.
type Event struct {
	Type   EventType
	X      int
	Y      int
	Button tea.MouseButton
}

// EnableMouseCmd returns a tea.Cmd to enable mouse cell motion events in Bubble Tea.
func EnableMouseCmd() tea.Cmd {
	return tea.EnableMouseCellMotion
}

// ParseEvent converts a tea.MouseMsg into a mouse.Event.
func ParseEvent(msg tea.MouseMsg) Event {
	evt := Event{
		Type:   EventNone,
		X:      msg.X,
		Y:      msg.Y,
		Button: msg.Button,
	}

	switch msg.Button {
	case tea.MouseButtonWheelUp:
		evt.Type = EventScrollUp
	case tea.MouseButtonWheelDown:
		evt.Type = EventScrollDown
	case tea.MouseButtonLeft:
		evt.Type = EventClick
	default:
		// Fallback check based on string representation
		msgStr := msg.String()
		if msgStr == "wheelup" || msgStr == "mouse wheel up" {
			evt.Type = EventScrollUp
		} else if msgStr == "wheeldown" || msgStr == "mouse wheel down" {
			evt.Type = EventScrollDown
		}
	}

	return evt
}
