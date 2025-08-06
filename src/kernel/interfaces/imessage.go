package interfaces

// MessageType represents the type of a message in the system, typically defined by constant values.
type MessageType int

const (
	MessageTypeError MessageType = iota
	MessageTypeRead
	MessageTypeTimer
	MessageTypeTimedMessage
	MessageTypePaintRequest
	MessageTypeQuit
	MessageTypeProcessExec
	MessageTypeProcessStart
	MessageTypeProcessActivate
	MessageTypeProcessExit
	MessageTypeProcessKill
	MessageTypeProcessKillAll
	MessageTypeProcessSetForeground
	MessageTypeProcessIsActiveRequest
	MessageTypeProcessIsActiveResponse
	MessageTypeProcessListRequest
	MessageTypeProcessListResponse
	MessageTypeCWDSet
	MessageTypeCWDGetRequest
	MessageTypeCWDGetResponse
	MessageTypeCWDPathRequest
	MessageTypeCWDPathResponse
	MessageTypeCWDNameRequest
	MessageTypeCWDNameResponse
	MessageTypeCWDDirectoryListingRequest
	MessageTypeCWDDirectoryListingResponse
	MessageTypeFileSystemSuggestionRequest
	MessageTypeFileSystemSuggestionResponse
	MessageTypeFileSystemHelpRequest
	MessageTypeFileSystemHelpResponse
	MessageTypeWrite
	MessageTypeWriteLn
	MessageTypeWriteColor
	MessageTypeWriteColorLn
	MessageTypeWriteNormal
	MessageTypeWriteHighlights
	MessageTypeWriteCritical
	MessageTypeWritePromptEOL
	MessageTypeWritePromptLine
	MessageTypeClearScreen
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
	MessageTypeTimerStop
	MessageTypeExitRequested
)

// IMessage defines the interface for messages used within the system, requiring a method to retrieve the message type.
type IMessage interface {
	GetType() MessageType

	Router() IRouter

	Ack()
}

// Message represents a basic unit containing a MessageType to define its specific behavior or category.
type Message struct {
	router IRouter
	kind   MessageType
}

// NewMessage creates a new Message instance with the specified MessageType.
func NewMessage(router IRouter, kind MessageType) *Message {
	return &Message{
		router: router,
		kind:   kind,
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

// Ack acknowledges the processing of the message, preventing further handling or retries by the system.
func (m *Message) Ack() {

}
