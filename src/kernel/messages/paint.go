package messages

import "github.com/markel1974/c64emu/src/kernel/interfaces"

// MessagePaint represents a message used to trigger a paint operation. It embeds the Message struct.
type MessagePaint struct {
	interfaces.Message
}

// NewMessagePaint creates a new instance of MessagePaint with the MessageType set to MessageTypePaint.
func NewMessagePaint() *MessagePaint {
	return &MessagePaint{
		Message: *interfaces.NewMessage(interfaces.MessageTypePaint),
	}
}
