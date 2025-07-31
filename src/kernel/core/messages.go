package core

// MessageType represents the type of a message in the system, typically defined by constant values.
type MessageType int

// MessageTypeRead represents a message type for read operations.
// MessageTypeTimer represents a message type for timer operations.
// MessageTypePaint represents a message type for paint operations.
// MessageTypeQuit represents a message type for quit operations.
const (
	MessageTypeRead  MessageType = iota
	MessageTypeTimer MessageType = iota
	MessageTypePaint MessageType = iota
	MessageTypeQuit  MessageType = iota
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

// PostEvent sends the Message instance to the specified channel asynchronously using a goroutine.
//func (m *Message) PostEvent(ch chan IMessage) {
//	go func() { ch <- m }()
//}

// MessageRead represents a specific type of Message containing read operation data. It embeds Message and includes a data field.
type MessageRead struct {
	Message
	data []byte
}

// NewMessageRead creates a new MessageRead instance with provided data and limits its length to n if necessary.
func NewMessageRead(data []byte, n int) *MessageRead {
	if n > len(data) {
		n = len(data) - 1
	}
	x := data[:n]
	return &MessageRead{
		Message: Message{MessageTypeRead},
		data:    x,
	}
}

// MessageTimer represents a timed message with an associated process ID, timer ID, and interval.
type MessageTimer struct {
	Message
	pid      int
	tid      int
	interval int
}

// NewMessageTimer creates a new MessageTimer instance with the specified process ID and interval.
func NewMessageTimer(pid int, interval int) *MessageTimer {
	return &MessageTimer{
		Message:  Message{MessageTypeTimer},
		pid:      pid,
		interval: interval,
	}
}

// MessageQuit represents a message signaling a quit operation or system termination.
// It embeds the base Message type and sets its kind to MessageTypeQuit.
type MessageQuit struct {
	Message
}

// NewMessageQuit creates and returns a new MessageQuit instance with the MessageType set to MessageTypeQuit.
func NewMessageQuit() *MessageQuit {
	return &MessageQuit{
		Message: Message{MessageTypeQuit},
	}
}

// MessagePaint represents a message used to trigger a paint operation. It embeds the Message struct.
type MessagePaint struct {
	Message
}

// NewMessagePaint creates a new instance of MessagePaint with the MessageType set to MessageTypePaint.
func NewMessagePaint() *MessagePaint {
	return &MessagePaint{
		Message: Message{MessageTypePaint},
	}
}
