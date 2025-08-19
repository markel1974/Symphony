package messages

import "github.com/markel1974/c64emu/src/kernel/interfaces"

// MessageGetScreenSizeRequest represents a request message to retrieve the screen size. It implements the IMessage interface.
type MessageGetScreenSizeRequest struct {
	interfaces.IMessage
}

// MessageGetScreenSizeResponse represents the response message containing the screen size dimensions.
// It includes a width and height, along with the base message interface IMessage functionalities.
type MessageGetScreenSizeResponse struct {
	interfaces.IMessage
	width  int
	height int
}

// NewMessageGetScreenSizeRequest initializes a new MessageGetScreenSizeRequest with the specified acknowledgment channel.
func NewMessageGetScreenSizeRequest(source int, destination int, ack chan bool) *MessageGetScreenSizeRequest {
	return &MessageGetScreenSizeRequest{
		IMessage: interfaces.NewMessageRequest(destination, source, interfaces.MessageTypeGetScreenSizeRequest, ack),
	}
}

// CreateResponse creates and sets a response message with the given width and height for the current request instance.
func (m *MessageGetScreenSizeRequest) CreateResponse(responder int, width int, height int) {
	base := m.PrepareResponse(responder, interfaces.MessageTypeGetScreenSizeResponse)
	msg := &MessageGetScreenSizeResponse{IMessage: base, width: width, height: height}
	m.SetResponse(msg)
}

// GetResponse returns the width and height values stored in the MessageGetScreenSizeResponse instance.
func (m *MessageGetScreenSizeResponse) GetResponse() (int, int) {
	return m.width, m.height
}

// MessageSetScreenSize represents a message for setting the screen size with specified width and height.
type MessageSetScreenSize struct {
	interfaces.IMessage
	width  int
	height int
}

// NewMessageSetScreenSize creates a new MessageSetScreenSize with the specified width and height for screen resizing.
func NewMessageSetScreenSize(source int, destination int, width int, height int) *MessageSetScreenSize {
	return &MessageSetScreenSize{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeSetScreenSize),
		width:    width,
		height:   height,
	}
}

// Width returns the width value of the MessageSetScreenSize instance.
func (m *MessageSetScreenSize) Width() int {
	return m.width
}

// Height returns the height value of the MessageSetScreenSize instance.
func (m *MessageSetScreenSize) Height() int {
	return m.height
}

// MessageWrite represents a write operation message that includes string data and an end-of-line (eol) flag.
type MessageWrite struct {
	interfaces.IMessage
	data string
	eol  bool
}

// NewMessageWrite creates a new MessageWrite instance with the specified data and end-of-line flag.
func NewMessageWrite(source int, destination int, data string, eol bool) *MessageWrite {
	return &MessageWrite{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeWrite),
		data:     data,
		eol:      eol,
	}
}

// Data returns the data string stored within the MessageWrite instance.
func (m *MessageWrite) Data() string {
	return m.data
}

// Eol returns a boolean indicating whether the message ends with an end-of-line character.
func (m *MessageWrite) Eol() bool {
	return m.eol
}

// MessageWriteLn represents a message that writes a line of data, embedding IMessage for message-related behavior.
type MessageWriteLn struct {
	interfaces.IMessage
	Data string
}

// MessageWriteColor represents a message for writing a string with specified foreground and background colors, mode, and EOL flag.
type MessageWriteColor struct {
	interfaces.IMessage
	data string
	fg   interfaces.ColorDef
	bg   interfaces.ColorDef
	mode interfaces.ColorMode
	eol  bool
}

// NewMessageWriteColor creates a new MessageWriteColor with specified text data, colors, mode, and end-of-line flag.
func NewMessageWriteColor(source int, destination int, data string, fg, bg interfaces.ColorDef, mode interfaces.ColorMode, eol bool) *MessageWriteColor {
	return &MessageWriteColor{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeWriteColor),
		data:     data,
		fg:       fg,
		bg:       bg,
		mode:     mode,
		eol:      eol,
	}
}

// Data returns the data string associated with the MessageWriteColor instance.
func (m *MessageWriteColor) Data() string {
	return m.data
}

// Fg returns the foreground color definition of the MessageWriteColor instance.
func (m *MessageWriteColor) Fg() interfaces.ColorDef {
	return m.fg
}

// Bg returns the background color of the MessageWriteColor instance.
func (m *MessageWriteColor) Bg() interfaces.ColorDef {
	return m.bg
}

// Mode returns the color mode used for rendering the message.
func (m *MessageWriteColor) Mode() interfaces.ColorMode {
	return m.mode
}

// Eol returns a boolean indicating whether the message should output a line ending after its content.
func (m *MessageWriteColor) Eol() bool {
	return m.eol
}

// MessageClearScreen represents a message type used to indicate a request to clear the screen in the messaging system.
type MessageClearScreen struct {
	interfaces.IMessage
}

// NewMessageClearScreen creates a new MessageClearScreen instance representing a request to clear the screen.
func NewMessageClearScreen(source int, destination int) *MessageClearScreen {
	return &MessageClearScreen{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeClearScreen),
	}
}

// MessageClearLine represents a message for clearing a specific line within the system.
// It embeds the IMessage interface and includes a line string attribute to specify the line to be cleared.
type MessageClearLine struct {
	interfaces.IMessage
	line string
}

// NewMessageClearLine creates a new MessageClearLine instance to request clearing a specific line on the screen.
func NewMessageClearLine(source int, destination int, line string) *MessageClearLine {
	return &MessageClearLine{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeClearLine),
		line:     line,
	}
}

// Line returns the stored line message from the MessageClearLine instance.
func (m *MessageClearLine) Line() string {
	return m.line
}

// MessageSaveCursor represents a message type for saving the current cursor position in the terminal or screen state.
type MessageSaveCursor struct {
	interfaces.IMessage
}

// NewMessageSaveCursor creates a new MessageSaveCursor with the MessageTypeSaveCursor type.
func NewMessageSaveCursor(source int, destination int) *MessageSaveCursor {
	return &MessageSaveCursor{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeSaveCursor),
	}
}

// MessageRestoreCursor represents a message type for restoring the cursor position in the rendering system.
type MessageRestoreCursor struct {
	interfaces.IMessage
}

// NewMessageRestoreCursor creates a new MessageRestoreCursor instance with the MessageTypeRestoreCursor message type.
func NewMessageRestoreCursor(source int, destination int) *MessageRestoreCursor {
	return &MessageRestoreCursor{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeRestoreCursor),
	}
}

// MessageMoveCursorLeft represents a message indicating a request to move the cursor one position to the left.
type MessageMoveCursorLeft struct {
	interfaces.IMessage
}

// NewMessageMoveCursorLeft creates a new MessageMoveCursorLeft instance to signal a cursor move to the left.
func NewMessageMoveCursorLeft(source int, destination int) *MessageMoveCursorLeft {
	return &MessageMoveCursorLeft{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeMoveCursorLeft),
	}
}

// MessageMoveCursorRight represents a message instructing the system to move the cursor one step to the right.
type MessageMoveCursorRight struct {
	interfaces.IMessage
}

// NewMessageMoveCursorRight creates a new MessageMoveCursorRight instance with MessageTypeMoveCursorRight.
func NewMessageMoveCursorRight(source int, destination int) *MessageMoveCursorRight {
	return &MessageMoveCursorRight{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeMoveCursorRight),
	}
}

// MessageMoveCursor represents a message to move the cursor in a graphical or textual interface.
// It specifies the new cursor position with row and column coordinates.
// Implements the IMessage interface for structured messaging.
type MessageMoveCursor struct {
	interfaces.IMessage
	row    int
	column int
}

// NewMessageMoveCursor creates a new MessageMoveCursor instance with the specified source, destination, row, and column values.
func NewMessageMoveCursor(source int, destination int, row int, column int) *MessageMoveCursor {
	return &MessageMoveCursor{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeMoveCursor),
		row:      row,
		column:   column,
	}
}

// Row returns the row index associated with the MessageMoveCursor instance.
func (m *MessageMoveCursor) Row() int {
	return m.row
}

// Column returns the column index associated with the MessageMoveCursor.
func (m *MessageMoveCursor) Column() int {
	return m.column
}
