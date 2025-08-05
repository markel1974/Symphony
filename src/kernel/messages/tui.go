package messages

import "github.com/markel1974/c64emu/src/kernel/interfaces"

// MessagePaint represents a message used to trigger a paint operation. It embeds the Message struct.
//type MessagePaint struct {
//	interfaces.Message
//}

// NewMessagePaint creates a new instance of MessagePaint with the MessageType set to MessageTypePaint.
//func NewMessagePaint() *MessagePaint {
//	return &MessagePaint{
//		Message: *interfaces.NewMessage(interfaces.MessageTypePaint),
//	}
//}

// MessagePaintRequest represents a message used to trigger a paint operation. It embeds the Message struct.
type MessagePaintRequest struct {
	interfaces.Message
}

// NewMessagePaintRequest creates a new instance of MessagePaint with the MessageType set to MessageTypePaint.
func NewMessagePaintRequest() *MessagePaintRequest {
	return &MessagePaintRequest{
		Message: *interfaces.NewMessage(interfaces.MessageTypePaintRequest),
	}
}

// MessageWindowsSelectionBegin represents a message signaling the start of a Windows selection operation.
type MessageWindowsSelectionBegin struct {
	interfaces.Message
}

// NewMessageWindowsSelectionBegin creates and returns a new instance of MessageWindowsSelectionBegin with the proper message type.
func NewMessageWindowsSelectionBegin() *MessageWindowsSelectionBegin {
	return &MessageWindowsSelectionBegin{
		Message: *interfaces.NewMessage(interfaces.MessageTypeWindowsSelectionBegin),
	}
}

// MessageWindowsSelectionEnd represents a message signaling the end of a windows selection operation.
type MessageWindowsSelectionEnd struct {
	interfaces.Message
}

// NewMessageWindowsSelectionEnd creates and returns a pointer to a new MessageWindowsSelectionEnd instance with the appropriate type.
func NewMessageWindowsSelectionEnd() *MessageWindowsSelectionEnd {
	return &MessageWindowsSelectionEnd{
		Message: *interfaces.NewMessage(interfaces.MessageTypeWindowsSelectionEnd),
	}
}

// MessageWindowsSelectionNext defines a message type used to indicate moving to the next selection in a Windows environment.
type MessageWindowsSelectionNext struct {
	interfaces.Message
}

// NewMessageWindowsSelectionNext creates and initializes a new MessageWindowsSelectionNext instance with the appropriate message type.
func NewMessageWindowsSelectionNext() *MessageWindowsSelectionNext {
	return &MessageWindowsSelectionNext{
		Message: *interfaces.NewMessage(interfaces.MessageTypeWindowsSelectionNext),
	}
}

// MessageWindowsSelectionPrevious is a specialized message type indicating a request to move to the previous selection window.
type MessageWindowsSelectionPrevious struct {
	interfaces.Message
}

// NewMessageWindowsSelectionPrevious creates a new instance of MessageWindowsSelectionPrevious with predefined message type.
func NewMessageWindowsSelectionPrevious() *MessageWindowsSelectionPrevious {
	return &MessageWindowsSelectionPrevious{
		Message: *interfaces.NewMessage(interfaces.MessageTypeWindowsSelectionPrevious),
	}
}

// MessageWindowsSelectionOptions represents a message detailing options for Windows selection with additional parameters.
type MessageWindowsSelectionOptions struct {
	interfaces.Message
	Option rune
	Value  float64
}

// NewMessageWindowsSelectionOptions initializes a MessageWindowsSelectionOptions with the given option and value.
func NewMessageWindowsSelectionOptions(option rune, value float64) *MessageWindowsSelectionOptions {
	return &MessageWindowsSelectionOptions{
		Message: *interfaces.NewMessage(interfaces.MessageTypeWindowsSelectionOptions),
		Option:  option,
		Value:   value,
	}
}
