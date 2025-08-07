package interfaces

// MessageType represents the type of a message in the system, typically defined by constant values.
type MessageType int

const (
	MessageTypeError MessageType = iota
	MessageTypeRead
	MessageTypeTimer
	MessageTypeTimedMessage
	MessageTypePaintRequest
	MessageTypePaintPrepare
	MessageTypePaintApply
	MessageTypeQuit
	MessageTypeProcessExec
	MessageTypeProcessStart
	MessageTypeProcessActivate
	MessageTypeProcessExit
	MessageTypeProcessKill
	MessageTypeProcessKillAll
	MessageTypeProcessKillForeground
	MessageTypeProcessSetForeground
	MessageTypeProcessIsActiveRequest
	MessageTypeProcessIsActiveResponse
	MessageTypeProcessListRequest
	MessageTypeProcessListResponse
	MessageTypeCWDSet
	MessageTypeCWDGetRequest
	MessageTypeCWDPathRequest
	MessageTypeCWDPathResponse
	MessageTypeCWDNameRequest
	MessageTypeCWDNameResponse
	MessageTypeCWDDirectoryListing
	MessageTypeFileSystemSuggestion
	MessageTypeFileSystemHelpRequest
	MessageTypeFileSystemHelpResponse
	MessageTypeWrite
	MessageTypeWriteLn
	MessageTypeWriteColor
	MessageTypeClearScreen
	MessageTypeClearLine
	MessageTypeMoveCursorLeft
	MessageTypeMoveCursorRight
	MessageTypeSaveCursor
	MessageTypeRestoreCursor
	MessageTypeWindowsSelectionBegin
	MessageTypeWindowsSelectionEnd
	MessageTypeWindowsSelectionNext
	MessageTypeWindowsSelectionPrevious
	MessageTypeWindowsSelectionOptions
	MessageTypeScreenSizeRequest
	MessageTypeScreenSizeResponse
	MessageTypeSetScreenSize
	MessageTypeTimerCreate
	MessageTypeTimerCreated
	MessageTypeTimerStop
	MessageTypeExitRequested
)

// IMessage defines the interface for messages used within the system, requiring a method to retrieve the message type.
type IMessage interface {
	GetType() MessageType

	Router() IRouter

	Response() bool

	Ack() bool
}

// Message represents a basic unit containing a MessageType to define its specific behavior or category.
type Message struct {
	router   IRouter
	kind     MessageType
	response bool
}

// NewMessage creates a new Message instance with the specified MessageType.
func NewMessage(router IRouter, kind MessageType) *Message {
	return &Message{
		router:   router,
		kind:     kind,
		response: false,
	}
}

// GetType returns the MessageType of the current Message instance.
func (m *Message) GetType() MessageType {
	return m.kind
}

// Router returns the IRouter instance associated with the current Message instance.
func (m *Message) Router() IRouter {
	return m.router
}

// Response returns a boolean indicating whether the message should be responded to.
func (m *Message) Response() bool {
	return m.response
}

// MakeResponse marks the message as a response to be sent back to the sender.
func (m *Message) MakeResponse() {
	m.response = true
}

// Ack acknowledges the processing of the message, preventing further handling or retries by the system.
func (m *Message) Ack() bool {
	return false
}
