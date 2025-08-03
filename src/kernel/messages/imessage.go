package messages

// MessageType represents the type of a message in the system, typically defined by constant values.
type MessageType int

// MessageTypeRead represents a message type for read operations.
// MessageTypeTimer represents a message type for timer operations.
// MessageTypePaint represents a message type for paint operations.
// MessageTypeQuit represents a message type for quit operations.
const (
	//MessageTypeIORead MessageType = iota
	MessageTypeRead MessageType = iota
	MessageTypeTimer
	MessageTypePaint
	MessageTypeQuit
)

// IMessage defines the interface for messages used within the system, requiring a method to retrieve the message type.
type IMessage interface {
	GetType() MessageType
}

// Message represents a basic unit containing a MessageType to define its specific behavior or category.
type Message struct {
	kind MessageType
}

// GetType returns the MessageType of the current Message instance.
func (m *Message) GetType() MessageType {
	return m.kind
}
