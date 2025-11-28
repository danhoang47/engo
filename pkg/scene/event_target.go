package scene

type EventType int

const (
	MouseEvent EventType = iota
	KeyboardEvent
)

type Event struct {
	Type EventType

	Target  *Node
	Current *Node

	// Propagation bool
}

type EventListener func(event Event)

type EventTarget interface {
	DispatchEvent(event Event) bool
	AddEventListener(eventType EventType, listener EventListener)
}
