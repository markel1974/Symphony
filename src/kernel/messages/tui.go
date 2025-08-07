package messages

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
)

// MessagePaintRequest represents a message used to trigger a paint operation. It embeds the Message struct.
type MessagePaintRequest struct {
	interfaces.Message
}

// NewMessagePaintRequest creates a new instance of MessagePaint with the MessageType set to MessageTypePaint.
func NewMessagePaintRequest(router interfaces.IRouter) *MessagePaintRequest {
	return &MessagePaintRequest{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypePaintRequest),
	}
}

// MessagePaintPrepare represents a message used to trigger a paint operation. It embeds the Message struct.
type MessagePaintPrepare struct {
	interfaces.Message
	surface interfaces.ISurface
}

// NewMessagePaintPrepare creates a new instance of MessagePaint with the MessageType set to MessageTypePaint.
func NewMessagePaintPrepare(router interfaces.IRouter, surface interfaces.ISurface) *MessagePaintPrepare {
	mp := &MessagePaintPrepare{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypePaintPrepare),
		surface: surface,
	}
	mp.SetReply()
	return mp
}

// Surface returns the ISurface associated with the MessagePaint instance.
func (m *MessagePaintPrepare) Surface() interfaces.ISurface {
	return m.surface
}

// MessagePaintApply represents a message used to apply painting or rendering operations on an ISurface.
// It embeds the generic Message type and associates an ISurface for graphical or textual manipulations.
// Surface provides access to the ISurface instance associated with this message for drawing and rendering tasks.
type MessagePaintApply struct {
	interfaces.Message
	surface interfaces.ISurface
}

// NewMessagePaintApply creates a new MessagePaintApply instance with the specified router and surface for rendering tasks.
func NewMessagePaintApply(router interfaces.IRouter, surface interfaces.ISurface) *MessagePaintApply {
	mp := &MessagePaintApply{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypePaintApply),
		surface: surface,
	}
	return mp
}

// Surface returns the ISurface associated with the MessagePaint instance.
func (m *MessagePaintApply) Surface() interfaces.ISurface {
	return m.surface
}

// MessageWindowsSelectionBegin represents a message signaling the start of a Windows selection operation.
type MessageWindowsSelectionBegin struct {
	interfaces.Message
}

// NewMessageWindowsSelectionBegin creates and returns a new instance of MessageWindowsSelectionBegin with the proper message type.
func NewMessageWindowsSelectionBegin(router interfaces.IRouter) *MessageWindowsSelectionBegin {
	return &MessageWindowsSelectionBegin{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeWindowsSelectionBegin),
	}
}

// MessageWindowsSelectionEnd represents a message signaling the end of a windows selection operation.
type MessageWindowsSelectionEnd struct {
	interfaces.Message
}

// NewMessageWindowsSelectionEnd creates and returns a pointer to a new MessageWindowsSelectionEnd instance with the appropriate type.
func NewMessageWindowsSelectionEnd(router interfaces.IRouter) *MessageWindowsSelectionEnd {
	return &MessageWindowsSelectionEnd{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeWindowsSelectionEnd),
	}
}

// MessageWindowsSelectionNext defines a message type used to indicate moving to the next selection in a Windows environment.
type MessageWindowsSelectionNext struct {
	interfaces.Message
}

// NewMessageWindowsSelectionNext creates and initializes a new MessageWindowsSelectionNext instance with the appropriate message type.
func NewMessageWindowsSelectionNext(router interfaces.IRouter) *MessageWindowsSelectionNext {
	return &MessageWindowsSelectionNext{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeWindowsSelectionNext),
	}
}

// MessageWindowsSelectionPrevious is a specialized message type indicating a request to move to the previous selection window.
type MessageWindowsSelectionPrevious struct {
	interfaces.Message
}

// NewMessageWindowsSelectionPrevious creates a new instance of MessageWindowsSelectionPrevious with predefined message type.
func NewMessageWindowsSelectionPrevious(router interfaces.IRouter) *MessageWindowsSelectionPrevious {
	return &MessageWindowsSelectionPrevious{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeWindowsSelectionPrevious),
	}
}

// MessageWindowsSelectionOptions represents a message detailing options for Windows selection with additional parameters.
type MessageWindowsSelectionOptions struct {
	interfaces.Message
	option rune
	value  float64
}

// NewMessageWindowsSelectionOptions initializes a MessageWindowsSelectionOptions with the given option and value.
func NewMessageWindowsSelectionOptions(router interfaces.IRouter, option rune, value float64) *MessageWindowsSelectionOptions {
	return &MessageWindowsSelectionOptions{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeWindowsSelectionOptions),
		option:  option,
		value:   value,
	}
}

// Option returns the rune representing the selected option in MessageWindowsSelectionOptions.
func (m *MessageWindowsSelectionOptions) Option() rune {
	return m.option
}

// Value returns the float64 value associated with the MessageWindowsSelectionOptions instance.
func (m *MessageWindowsSelectionOptions) Value() float64 {
	return m.value
}
