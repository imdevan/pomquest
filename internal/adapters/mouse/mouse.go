package mouse

import (
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

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

// EnableMouseCmd returns a tea.Cmd to enable all mouse motion events (including hover) in Bubble Tea.
func EnableMouseCmd() tea.Cmd {
	return tea.EnableMouseAllMotion
}

// InitManager initializes the global zone manager.
func InitManager() {
	zone.NewGlobal()
}

func ensureManager() {
	zone.NewGlobal()
}

// Scan processes content, recording zone coordinates and stripping zone marker ANSI codes.
func Scan(content string) string {
	ensureManager()
	return zone.Scan(content)
}

// Mark wraps content with a zone identifier marker.
func Mark(id string, content string) string {
	ensureManager()
	return zone.Mark(id, content)
}

// InBounds returns true if the mouse message occurred within the bounds of zone id.
func InBounds(id string, msg tea.Msg) bool {
	ensureManager()
	z := zone.Get(id)
	if z == nil {
		return false
	}
	if mouseMsg, ok := msg.(tea.MouseMsg); ok {
		return z.InBounds(mouseMsg)
	}
	return false
}

// Get returns zone info for the given zone ID.
func Get(id string) *zone.ZoneInfo {
	ensureManager()
	return zone.Get(id)
}

// ParseEvent converts a tea.MouseMsg into a mouse.Event.
func ParseEvent(msg tea.MouseMsg) Event {
	evt := Event{
		Type:   EventNone,
		X:      msg.X,
		Y:      msg.Y,
		Button: msg.Button,
	}

	if msg.Type == tea.MouseRelease || msg.Action == tea.MouseActionRelease {
		evt.Type = EventRelease
		return evt
	}

	if msg.Type == tea.MouseMotion || msg.Action == tea.MouseActionMotion {
		evt.Type = EventMotion
		return evt
	}

	switch msg.Button {
	case tea.MouseButtonWheelUp:
		evt.Type = EventScrollUp
	case tea.MouseButtonWheelDown:
		evt.Type = EventScrollDown
	case tea.MouseButtonLeft:
		evt.Type = EventClick
	case tea.MouseButtonNone:
		evt.Type = EventMotion
	default:
		msgStr := msg.String()
		if msgStr == "wheelup" || msgStr == "mouse wheel up" {
			evt.Type = EventScrollUp
		} else if msgStr == "wheeldown" || msgStr == "mouse wheel down" {
			evt.Type = EventScrollDown
		} else if msgStr == "release" {
			evt.Type = EventRelease
		} else if msgStr == "motion" {
			evt.Type = EventMotion
		}
	}

	return evt
}
