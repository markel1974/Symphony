package interfaces

// MessageType represents the type of a message exchanged within the system, identifying the message's purpose or category.
type MessageType int

// MessageTypeError represents an error message type.
// MessageTypeQuit represents a quit message type.
// MessageTypeRead represents a read operation message type.
// MessageTypeTimer represents a timer-related message type.
// MessageTypeGetScreenSizeRequest requests the current screen size.
// MessageTypeGetScreenSizeResponse provides the screen size as a response.
// MessageTypeSetScreenSize sets the screen size.
// MessageTypeTimerCreate requests creation of a timer.
// MessageTypeTimerCreated indicates a timer has been created.
// MessageTypeTimerStop stops an existing timer.
// MessageTypeExitRequested signals a request to exit.
// MessageTypePaintRequest requests a paint operation.
// MessageTypePaintPrepare prepares for a paint operation.
// MessageTypePaintApply applies a paint operation.
// MessageTypeWrite writes a message.
// MessageTypeWriteLn writes a message with a newline.
// MessageTypeWriteColor writes a colored message.
// MessageTypeClearScreen clears the screen.
// MessageTypeClearLine clears the current line.
// MessageTypeMoveCursorLeft moves the cursor to the left.
// MessageTypeMoveCursorRight moves the cursor to the right.
// MessageTypeSaveCursor saves the cursor position.
// MessageTypeRestoreCursor restores the cursor position.
// MessageTypeWindowsSelectionBegin begins a window selection.
// MessageTypeWindowsSelectionEnd ends a window selection.
// MessageTypeWindowsSelectionNext moves to the next window in selection.
// MessageTypeWindowsSelectionPrevious moves to the previous window in selection.
// MessageTypeWindowsSelectionOptions retrieves selection options for windows.
// MessageTypeProcessExec executes a process.
// MessageTypeProcessStart starts a process.
// MessageTypeProcessActivate activates an existing process.
// MessageTypeProcessIsRunningRequest checks if a process is running.
// MessageTypeProcessIsRunningResponse responds whether a process is running.
// MessageTypeProcessExit indicates a process has exited.
// MessageTypeProcessKill kills a process.
// MessageTypeProcessKillAll kills all processes.
// MessageTypeProcessKillForeground kills the foreground process.
// MessageTypeProcessSetForeground sets a process as the foreground process.
// MessageTypeProcessListRequest requests a list of processes.
// MessageTypeProcessListResponse provides a list of processes.
// MessageTypeFileSystemCWDSetRequest requests to set the current working directory.
// MessageTypeFileSystemCWDSetResponse responds to the CWD set request.
// MessageTypeFileSystemCWDGetRequest requests the current working directory.
// MessageTypeFileSystemCWDGetResponse provides the current working directory.
// MessageTypeFileSystemCWDDirectoryListingRequest requests a directory listing in CWD.
// MessageTypeFileSystemCWDDirectoryListingResponse provides a directory listing in CWD.
// MessageTypeFileSystemCWDPathRequest requests the path in the current working directory.
// MessageTypeFileSystemCWDPathResponse provides the path in the current working directory.
// MessageTypeFileSystemSuggestionRequest requests file system suggestions.
// MessageTypeFileSystemSuggestionResponse provides file system suggestions.
// MessageTypeFileSystemHelpRequest requests help information related to the file system.
// MessageTypeFileSystemHelpResponse provides help information related to the file system.
// MessageTypeFileSystemFindRequest requests to find files in the file system.
// MessageTypeFileSystemFindResponse provides results for the file system find request.
// MessageTypeNotifyProcessCreate notifies about the creation of a process.
// MessageTypeNotifyProcessForeground notifies about a process moving to the foreground.
// MessageTypeNotifyProcessTerminate notifies about the termination of a process.
const (
	MessageTypeError MessageType = iota
	MessageTypeQuit
	MessageTypeRead
	MessageTypeTimer
	MessageTypeGetScreenSizeRequest
	MessageTypeGetScreenSizeResponse
	MessageTypeSetScreenSize
	MessageTypeTimerCreate
	MessageTypeTimerCreated
	MessageTypeTimerStop
	MessageTypeExitRequested
	MessageTypePaintRequest
	MessageTypePaintPrepare
	MessageTypePaintApply
	MessageTypeWrite
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
	MessageTypeProcessIsRunningRequest
	MessageTypeProcessIsRunningResponse
	MessageTypeProcessExit
	MessageTypeProcessKill
	MessageTypeProcessKillAll
	MessageTypeProcessKillForeground
	MessageTypeProcessSetForeground
	MessageTypeProcessSetSelfForeground
	MessageTypeProcessListRequest
	MessageTypeProcessListResponse
	MessageTypeFileSystemCWDSetRequest
	MessageTypeFileSystemCWDSetResponse
	MessageTypeFileSystemCWDGetRequest
	MessageTypeFileSystemCWDGetResponse
	MessageTypeFileSystemCWDDirectoryListingRequest
	MessageTypeFileSystemCWDDirectoryListingResponse
	MessageTypeFileSystemCWDPathRequest
	MessageTypeFileSystemCWDPathResponse
	MessageTypeFileSystemSuggestionRequest
	MessageTypeFileSystemSuggestionResponse
	MessageTypeFileSystemHelpRequest
	MessageTypeFileSystemHelpResponse
	MessageTypeFileSystemFindRequest
	MessageTypeFileSystemFindResponse
	MessageTypeNotifyProcessCreate
	MessageTypeNotifyProcessForeground
	MessageTypeNotifyProcessTerminate
	//MessageTypeTimedMessage
)

// IMessage defines the interface for structured messages, supporting operations such as cloning, acknowledgment, and responses.
type IMessage interface {
	GetType() MessageType

	Source() int

	SetSource(source int)

	Destination() int

	SetDestination(destination int)

	PrepareResponse(responder int, kind MessageType) IMessage

	Response() IMessage

	SetResponse(response IMessage)

	Ack() chan bool
}

// Message represents a message structure used for communication between entities in the system.
// It contains metadata such as source, destination, type, and response handling mechanisms.
type Message struct {
	source      int
	destination int
	kind        MessageType
	ack         chan bool
	response    IMessage
}

// NewMessageNoAck initializes a new Message instance of the specified MessageType without an acknowledgment channel.
func NewMessageNoAck(source int, destination int, kind MessageType) *Message {
	return &Message{
		kind:        kind,
		ack:         nil,
		source:      source,
		destination: destination,
	}
}

// NewMessageRequest creates a new Message instance with the specified MessageType and acknowledgment channel.
func NewMessageRequest(source int, destination int, kind MessageType, ack chan bool) *Message {
	return &Message{
		kind:        kind,
		ack:         ack,
		source:      source,
		destination: destination,
	}
}

// PrepareResponse creates a new Message with the provided responder and type, reversing the source and destination fields.
func (m *Message) PrepareResponse(responder int, kind MessageType) IMessage {
	return &Message{
		kind:        kind,
		source:      responder,
		destination: m.source,
		ack:         m.ack,
	}
}

// GetType returns the type of the message as a MessageType.
func (m *Message) GetType() MessageType {
	return m.kind
}

// SetSource sets the source identifier for the message.
func (m *Message) SetSource(source int) {
	m.source = source
}

// Source returns the source identifier of the Message as an integer.
func (m *Message) Source() int {
	return m.source
}

// Destination returns the destination identifier of the message.
func (m *Message) Destination() int {
	return m.destination
}

// SetDestination updates the destination field of the Message with the specified value.
func (m *Message) SetDestination(destination int) {
	m.destination = destination
}

// Ack sends a signal to the acknowledgment channel if it is not nil.
func (m *Message) Ack() chan bool {
	return m.ack
}

// SetResponse sets the response message associated with the current message.
func (m *Message) SetResponse(response IMessage) {
	m.response = response
}

// Response retrieves the response message associated with the current message. Returns nil if no response exists.
func (m *Message) Response() IMessage {
	return m.response
}
