package messages

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
