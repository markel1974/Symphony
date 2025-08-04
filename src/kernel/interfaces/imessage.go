package interfaces

// MessageType represents the type of a message in the system, typically defined by constant values.
type MessageType int

const (
	MessageTypeRead MessageType = iota
	MessageTypeTimer
	MessageTypePaint
	MessageTypePaintRequest
	MessageTypeQuit
	MessageTypeProcessExec
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

	Ack()
}

// Message represents a basic unit containing a MessageType to define its specific behavior or category.
type Message struct {
	kind MessageType
}

func NewMessage(kind MessageType) *Message {
	return &Message{
		kind: kind,
	}
}

// GetType returns the MessageType of the current Message instance.
func (m *Message) GetType() MessageType {
	return m.kind
}

// Ack acknowledges the processing of the message, preventing further handling or retries by the system.
func (m *Message) Ack() {

}
