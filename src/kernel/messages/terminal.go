package messages

import "github.com/markel1974/c64emu/src/kernel/interfaces"

// MessageWrite represents a message intended for write operations, containing data as a string payload.
type MessageWrite struct {
	interfaces.Message
	Data string
}

// NewMessageWrite creates a new MessageWrite instance with the provided data and MessageTypeWrite.
// It initializes the embedded Message field using interfaces.NewMessage.
// Returns a pointer to the new MessageWrite instance.
func NewMessageWrite(router interfaces.IRouter, data string) *MessageWrite {
	return &MessageWrite{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeWrite),
		Data:    data,
	}
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
	Data string
	Fg   interfaces.ColorDef
	Bg   interfaces.ColorDef
	Mode interfaces.ColorMode
}

// NewMessageWriteColor creates a new MessageWriteColor instance with specified data, foreground and background colors, and mode.
func NewMessageWriteColor(router interfaces.IRouter, data string, fg, bg interfaces.ColorDef, mode interfaces.ColorMode) *MessageWriteColor {
	return &MessageWriteColor{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeWriteColor),
		Data:    data,
		Fg:      fg,
		Bg:      bg,
		Mode:    mode,
	}
}

// MessageWritePromptLine represents a message type containing a user prompt and a single line entry.
type MessageWritePromptLine struct {
	interfaces.Message
	Prompt string
	Line   string
}

// NewMessageWritePromptLine creates a new MessageWritePromptLine with the specified prompt and line content.
func NewMessageWritePromptLine(router interfaces.IRouter, prompt, line string) *MessageWritePromptLine {
	return &MessageWritePromptLine{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeWritePromptLine),
		Prompt:  prompt,
		Line:    line,
	}
}

// MessageWritePromptEOL represents a message containing a prompt string and an end-of-line flag.
// It embeds the Message struct from the interfaces package.
type MessageWritePromptEOL struct {
	interfaces.Message
	Prompt string
	Eol    bool
}

// NewMessageWritePromptEOL creates a new MessageWritePromptEOL with the specified prompt string and EOL flag.
func NewMessageWritePromptEOL(router interfaces.IRouter, prompt string, eol bool) *MessageWritePromptEOL {
	return &MessageWritePromptEOL{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeWritePromptEOL),
		Prompt:  prompt,
		Eol:     eol,
	}
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
