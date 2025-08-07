package messages

import "github.com/markel1974/c64emu/src/kernel/interfaces"

// MessageSetScreenSize represents a message to set the screen's width and height in the rendering system.
// It embeds the base Message type and includes width and height properties to define the screen dimensions.
type MessageSetScreenSize struct {
	interfaces.Message
	width  int
	height int
}

// NewMessageSetScreenSize creates a new MessageSetScreenSize instance for setting the screen width and height.
func NewMessageSetScreenSize(router interfaces.IRouter, width int, height int) *MessageSetScreenSize {
	return &MessageSetScreenSize{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeSetScreenSize),
		width:   width,
		height:  height,
	}
}

// Width returns the width of the screen as an integer.
func (m *MessageSetScreenSize) Width() int {
	return m.width
}

// Height returns the height value of the screen size stored in the MessageSetScreenSize instance.
func (m *MessageSetScreenSize) Height() int {
	return m.height
}

// MessageWrite represents a message designed to carry string data for writing operations with an optional end-of-line flag.
type MessageWrite struct {
	interfaces.Message
	data string
	eol  bool
}

// NewMessageWrite creates a new instance of MessageWrite with the given router, data, and end-of-line flag.
func NewMessageWrite(router interfaces.IRouter, data string, eol bool) *MessageWrite {
	return &MessageWrite{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeWrite),
		data:    data,
		eol:     eol,
	}
}

// Data returns the internal data string of the MessageWrite instance.
func (m *MessageWrite) Data() string {
	return m.data
}

// Eol returns a boolean indicating whether the end of the line (EOL) flag is set for the message.
func (m *MessageWrite) Eol() bool {
	return m.eol
}

// MessageWriteLn represents a message type for writing a line of data to an output, embedding a base Message structure.
// Data contains the string payload to be written as part of the message.
type MessageWriteLn struct {
	interfaces.Message
	Data string
}

// NewMessageWriteLn initializes and returns a pointer to a MessageWriteLn with the provided data.
// It sets the Message type to MessageTypeWriteLn.
func NewMessageWriteLn(router interfaces.IRouter, data string) *MessageWriteLn {
	return &MessageWriteLn{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeWriteLn),
		Data:    data,
	}
}

// MessageWriteColor represents a message containing text data with specified foreground, background colors, and color mode.
type MessageWriteColor struct {
	interfaces.Message
	data string
	fg   interfaces.ColorDef
	bg   interfaces.ColorDef
	mode interfaces.ColorMode
	eol  bool
}

// NewMessageWriteColor creates a new MessageWriteColor instance with specified data, foreground and background colors, and mode.
func NewMessageWriteColor(router interfaces.IRouter, data string, fg, bg interfaces.ColorDef, mode interfaces.ColorMode, eol bool) *MessageWriteColor {
	return &MessageWriteColor{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeWriteColor),
		data:    data,
		fg:      fg,
		bg:      bg,
		mode:    mode,
		eol:     eol,
	}
}

// Data returns the textual data contained in the MessageWriteColor instance.
func (m *MessageWriteColor) Data() string {
	return m.data
}

// Fg returns the foreground color (interfaces.ColorDef) associated with the MessageWriteColor instance.
func (m *MessageWriteColor) Fg() interfaces.ColorDef {
	return m.fg
}

// Bg returns the background color of the MessageWriteColor instance as an interfaces.ColorDef.
func (m *MessageWriteColor) Bg() interfaces.ColorDef {
	return m.bg
}

// Mode retrieves the color mode (interfaces.ColorMode) associated with the MessageWriteColor instance.
func (m *MessageWriteColor) Mode() interfaces.ColorMode {
	return m.mode
}

// Eol returns a boolean indicating whether the message should end with a line terminator.
func (m *MessageWriteColor) Eol() bool {
	return m.eol
}

// MessageClearScreen represents a message used to request clearing of the screen in the system.
// It embeds the Message type, inheriting its behavior and properties.
type MessageClearScreen struct {
	interfaces.Message
}

// NewMessageClearScreen creates a new instance of MessageClearScreen with the MessageType set to ClearScreen.
func NewMessageClearScreen(router interfaces.IRouter) *MessageClearScreen {
	return &MessageClearScreen{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeClearScreen),
	}
}

// MessageClearLine represents a message for clearing the specified line in the system.
type MessageClearLine struct {
	interfaces.Message
	line string
}

// NewMessageClearLine creates a new MessageClearLine instance to represent a clear line action in the system.
func NewMessageClearLine(router interfaces.IRouter, line string) *MessageClearLine {
	return &MessageClearLine{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeClearLine),
	}
}

// Line returns the line associated with the MessageClearLine instance.
func (m *MessageClearLine) Line() string {
	return m.line
}

// MessageSaveCursor extends the Message type and represents a specific message type used to save the cursor position.
type MessageSaveCursor struct {
	interfaces.Message
}

// NewMessageSaveCursor creates a new MessageSaveCursor instance with the MessageTypeSaveCursor message type.
func NewMessageSaveCursor(router interfaces.IRouter) *MessageSaveCursor {
	return &MessageSaveCursor{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeSaveCursor),
	}
}

// MessageRestoreCursor represents a specialized message type for restoring the cursor position on the screen.
type MessageRestoreCursor struct {
	interfaces.Message
}

// NewMessageRestoreCursor creates a new MessageRestoreCursor instance with a MessageTypeRestoreCursor type.
func NewMessageRestoreCursor(router interfaces.IRouter) *MessageRestoreCursor {
	return &MessageRestoreCursor{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeRestoreCursor),
	}
}

// MessageMoveCursorLeft is a message type used to represent the action of moving the cursor to the left in the system.
type MessageMoveCursorLeft struct {
	interfaces.Message
}

// NewMessageMoveCursorLeft creates a new MessageMoveCursorLeft instance to represent a cursor-left movement message.
func NewMessageMoveCursorLeft(router interfaces.IRouter) *MessageMoveCursorLeft {
	return &MessageMoveCursorLeft{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeMoveCursorLeft),
	}
}

// MessageMoveCursorRight represents a message that instructs the system to move the cursor one position to the right.
// It embeds the Message type, inheriting its properties and behavior within the messaging system.
type MessageMoveCursorRight struct {
	interfaces.Message
}

// NewMessageMoveCursorRight creates a new instance of MessageMoveCursorRight with the specified router.
func NewMessageMoveCursorRight(router interfaces.IRouter) *MessageMoveCursorRight {
	return &MessageMoveCursorRight{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeMoveCursorRight),
	}
}
