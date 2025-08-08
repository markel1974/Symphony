package interfaces

// MessageType represents the type of a message in the system, typically defined by constant values.
type MessageType int

const (
	MessageTypeError MessageType = iota
	MessageTypeQuit
	MessageTypeRead
	MessageTypeTimer
	MessageTypeGetScreenSize
	MessageTypeSetScreenSize
	MessageTypeTimerCreate
	MessageTypeTimerCreated
	MessageTypeTimerStop
	MessageTypeExitRequested
	MessageTypeTimedMessage
	MessageTypePaintRequest
	MessageTypePaintPrepare
	MessageTypePaintApply
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
	MessageTypeProcessExec
	MessageTypeProcessStart
	MessageTypeProcessActivate
	MessageTypeProcessExit
	MessageTypeProcessKill
	MessageTypeProcessKillAll
	MessageTypeProcessKillForeground
	MessageTypeProcessSetForeground
	MessageTypeProcessList
	MessageTypeFileSystemCWDSet
	MessageTypeFileSystemCWDGet
	MessageTypeFileSystemCWDDirectoryListing
	MessageTypeFileSystemCWDPath
	MessageTypeFileSystemSuggestion
	MessageTypeFileSystemHelp
	MessageTypeFileSystemFindRequest
	MessageTypeFileSystemFindResponse
)

// IMessage defines the interface for messages used within the system, requiring a method to retrieve the message type.
type IMessage interface {
	GetType() MessageType

	PID() int

	Response() bool

	Ack() bool
}

// Message represents a basic unit containing a MessageType to define its specific behavior or category.
type Message struct {
	originatorPID int
	kind          MessageType
	response      bool
}

// NewMessage creates a new Message instance with the specified MessageType.
func NewMessage(originatorPID int, kind MessageType) *Message {
	return &Message{
		originatorPID: originatorPID,
		kind:          kind,
		response:      false,
	}
}

// GetType returns the MessageType of the current Message instance.
func (m *Message) GetType() MessageType {
	return m.kind
}

// PID returns the originator's process ID (PID) associated with the Message instance.
func (m *Message) PID() int {
	return m.originatorPID
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
