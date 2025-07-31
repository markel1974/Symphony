package messages

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
