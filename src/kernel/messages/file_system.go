package messages

import "github.com/markel1974/c64emu/src/kernel/interfaces"

// MessageFileSystemFindRequest represents a request to find a file system entry based on the provided query parameters.
// It embeds IMessage and includes specific properties like requestor PID, requested line, and a protected flag.
type MessageFileSystemFindRequest struct {
	interfaces.IMessage
	requestorPID int
	line         string
	protected    bool
}

// MessageFileSystemFindResponse represents the response to a file system find request, including command, arguments, and error state.
type MessageFileSystemFindResponse struct {
	interfaces.IMessage
	cmd  interfaces.ICommand
	args []string
	err  error
}

// NewMessageFileSystemFindRequest creates a new MessageFileSystemFindRequest for finding files in the file system.
// It initializes the message with the provided acknowledgment channel, requestor PID, query line, and protection flag.
// Returns a pointer to the created MessageFileSystemFindRequest instance.
func NewMessageFileSystemFindRequest(ack chan bool, requestorPID int, line string, protected bool) *MessageFileSystemFindRequest {
	return &MessageFileSystemFindRequest{
		IMessage:     interfaces.NewMessageRequest(interfaces.MessageTypeFileSystemFindRequest, ack),
		requestorPID: requestorPID,
		line:         line,
		protected:    protected,
	}
}

// Protected returns the value of the `protected` field, indicating whether the file system find request is protected.
func (m *MessageFileSystemFindRequest) Protected() bool {
	return m.protected
}

// RequestorPID retrieves the process ID of the entity that initiated the file system find request.
func (m *MessageFileSystemFindRequest) RequestorPID() int {
	return m.requestorPID
}

// Line returns the line string associated with the MessageFileSystemFindRequest instance.
func (m *MessageFileSystemFindRequest) Line() string {
	return m.line
}

// CreateResponse initializes a response to the current message with the given command, arguments, and error, and sets it.
func (m *MessageFileSystemFindRequest) CreateResponse(cmd interfaces.ICommand, args []string, err error) {
	r := &MessageFileSystemFindResponse{IMessage: m.Clone(), cmd: cmd, args: args, err: err}
	m.SetResponse(interfaces.MessageTypeFileSystemFindResponse, r)
}

// GetResult retrieves the command, arguments, and error from the MessageFileSystemFindResponse instance.
func (m *MessageFileSystemFindResponse) GetResult() (interfaces.ICommand, []string, error) {
	return m.cmd, m.args, m.err
}

// MessageFileSystemCWDSetRequest represents a request to set the current working directory in the file system.
type MessageFileSystemCWDSetRequest struct {
	interfaces.IMessage
	path string
}

// MessageFileSystemCWDSetResponse represents the response to a request for setting the current working directory.
type MessageFileSystemCWDSetResponse struct {
	interfaces.IMessage
	result bool
}

// NewMessageFileSystemCWDSetRequest creates a new MessageFileSystemCWDSetRequest with the given acknowledgment channel and path.
func NewMessageFileSystemCWDSetRequest(ack chan bool, path string) *MessageFileSystemCWDSetRequest {
	return &MessageFileSystemCWDSetRequest{
		IMessage: interfaces.NewMessageRequest(interfaces.MessageTypeFileSystemCWDSetRequest, ack),
		path:     path,
	}
}

// Path returns the stored path associated with the MessageFileSystemCWDSetRequest instance.
func (m *MessageFileSystemCWDSetRequest) Path() string {
	return m.path
}

// CreateResponse creates and sets a response message of type MessageFileSystemCWDSetResponse with the provided result.
func (m *MessageFileSystemCWDSetRequest) CreateResponse(result bool) {
	msg := &MessageFileSystemCWDSetResponse{IMessage: m.Clone(), result: result}
	m.SetResponse(interfaces.MessageTypeFileSystemCWDSetResponse, msg)
}

// GetResponse returns the result of the MessageFileSystemCWDSetResponse indicating the operation's success state.
func (m *MessageFileSystemCWDSetResponse) GetResponse() bool {
	return m.result
}

// MessageFileSystemCWDGetRequest represents a request message to retrieve the current working directory.
type MessageFileSystemCWDGetRequest struct {
	interfaces.IMessage
}

// MessageFileSystemCWDGetResponse represents a response message for retrieving the current working directory in the file system.
// It contains the result as a string and embeds the IMessage interface for common message functionalities.
type MessageFileSystemCWDGetResponse struct {
	interfaces.IMessage
	result string
}

// NewMessageFileSystemCWDGetRequest creates a new MessageFileSystemCWDGetRequest with the specified acknowledgment channel.
func NewMessageFileSystemCWDGetRequest(ack chan bool) *MessageFileSystemCWDGetRequest {
	return &MessageFileSystemCWDGetRequest{
		IMessage: interfaces.NewMessageRequest(interfaces.MessageTypeFileSystemCWDGetRequest, ack),
	}
}

// CreateResponse generates a response message with the provided result and assigns it to the request's response field.
func (m *MessageFileSystemCWDGetRequest) CreateResponse(result string) {
	msg := &MessageFileSystemCWDGetResponse{IMessage: m.Clone(), result: result}
	m.SetResponse(interfaces.MessageTypeFileSystemCWDGetResponse, msg)
}

// GetResponse retrieves the result string associated with the MessageFileSystemCWDGetResponse instance.
func (m *MessageFileSystemCWDGetResponse) GetResponse() string {
	return m.result
}

// MessageFileSystemCWDPathRequest represents a message requesting the current working directory path in the file system.
type MessageFileSystemCWDPathRequest struct {
	interfaces.IMessage
}

// MessageFileSystemCWDPathResponse represents the response containing the path of the current working directory.
// It implements the IMessage interface for message handling and encapsulates the result as a string.
type MessageFileSystemCWDPathResponse struct {
	interfaces.IMessage
	result string
}

// NewMessageFileSystemCWDPathRequest creates a new MessageFileSystemCWDPathRequest with the specified acknowledgment channel.
func NewMessageFileSystemCWDPathRequest(ack chan bool) *MessageFileSystemCWDPathRequest {
	return &MessageFileSystemCWDPathRequest{
		IMessage: interfaces.NewMessageRequest(interfaces.MessageTypeFileSystemCWDPathRequest, ack),
	}
}

// CreateResponse constructs a response message with the given result string and assigns it to the current message.
func (m *MessageFileSystemCWDPathRequest) CreateResponse(result string) {
	msg := &MessageFileSystemCWDPathResponse{IMessage: m.Clone(), result: result}
	m.SetResponse(interfaces.MessageTypeFileSystemCWDPathResponse, msg)
}

// GetResponse returns the result string stored in the MessageFileSystemCWDPathResponse instance.
func (m *MessageFileSystemCWDPathResponse) GetResponse() string {
	return m.result
}

// MessageFileSystemCWDDirectoryListingRequest represents a request for listing contents of the current working directory.
type MessageFileSystemCWDDirectoryListingRequest struct {
	interfaces.IMessage
}

// MessageFileSystemCWDDirectoryListingResponse represents a response containing a directory listing of the current working directory.
// This type implements the IMessage interface for handling message operations.
// The result field provides the list of directory contents as a slice of strings.
type MessageFileSystemCWDDirectoryListingResponse struct {
	interfaces.IMessage
	result []string
}

// NewMessageFileSystemCWDDirectoryListingRequest creates a request message for listing the current working directory.
func NewMessageFileSystemCWDDirectoryListingRequest(ack chan bool) *MessageFileSystemCWDDirectoryListingRequest {
	return &MessageFileSystemCWDDirectoryListingRequest{
		IMessage: interfaces.NewMessageRequest(interfaces.MessageTypeFileSystemCWDDirectoryListingRequest, ack),
	}
}

// CreateResponse generates a response message with the provided directory listing result and sets it as the response.
func (m *MessageFileSystemCWDDirectoryListingRequest) CreateResponse(result []string) {
	msg := &MessageFileSystemCWDDirectoryListingResponse{IMessage: m.Clone(), result: result}
	m.SetResponse(interfaces.MessageTypeFileSystemCWDDirectoryListingResponse, msg)
}

// GetResponse retrieves the directory listing results as a slice of strings.
func (m *MessageFileSystemCWDDirectoryListingResponse) GetResponse() []string {
	return m.result
}

// MessageFileSystemSuggestionRequest represents a request to suggest file system operations or paths.
// It includes the input string `in`, cursor position `cursor`, and an acknowledgment channel `ack`.
type MessageFileSystemSuggestionRequest struct {
	interfaces.IMessage
	in     string
	cursor int
	ack    chan bool
}

// MessageFileSystemSuggestionResponse represents a response message containing file system suggestions and their validity.
type MessageFileSystemSuggestionResponse struct {
	interfaces.IMessage
	prefix     string
	suggestion []string
	valid      bool
}

// NewMessageFileSystemSuggestionRequest creates a new MessageFileSystemSuggestionRequest with input, cursor, and ack channel.
func NewMessageFileSystemSuggestionRequest(in string, cursor int, ack chan bool) *MessageFileSystemSuggestionRequest {
	return &MessageFileSystemSuggestionRequest{
		IMessage: interfaces.NewMessageRequest(interfaces.MessageTypeFileSystemSuggestionRequest, ack),
		in:       in,
		cursor:   cursor,
		ack:      ack,
	}
}

// In returns the `in` property of the MessageFileSystemSuggestionRequest instance, representing input data or state.
func (m *MessageFileSystemSuggestionRequest) In() string {
	return m.in
}

// Cursor returns the current position of the cursor within the MessageFileSystemSuggestionRequest instance.
func (m *MessageFileSystemSuggestionRequest) Cursor() int {
	return m.cursor
}

// CreateResponse generates a new MessageFileSystemSuggestionResponse with the provided prefix, suggestions, and validity.
func (m *MessageFileSystemSuggestionRequest) CreateResponse(prefix string, suggestion []string, valid bool) {
	msg := &MessageFileSystemSuggestionResponse{IMessage: m.Clone(), prefix: prefix, suggestion: suggestion, valid: valid}
	m.SetResponse(interfaces.MessageTypeFileSystemSuggestionResponse, msg)
}

// GetResponse retrieves the prefix, suggestion list, and validity status from the MessageFileSystemSuggestionResponse instance.
func (m *MessageFileSystemSuggestionResponse) GetResponse() (prefix string, suggestion []string, valid bool) {
	return m.prefix, m.suggestion, m.valid
}

// MessageFileSystemHelpRequest represents a request for file system-related help, encapsulating the path for context.
type MessageFileSystemHelpRequest struct {
	interfaces.IMessage
	path string
}

// MessageFileSystemHelpResponse represents a response containing the result of a file system help request.
// It embeds the IMessage interface for message behavior and includes fields for result data and error information.
type MessageFileSystemHelpResponse struct {
	interfaces.IMessage
	ack    chan bool
	result string
	err    error
}

// NewMessageFileSystemHelpRequest creates a new MessageFileSystemHelpRequest with the specified path and acknowledgment channel.
func NewMessageFileSystemHelpRequest(path string, ack chan bool) *MessageFileSystemHelpRequest {
	return &MessageFileSystemHelpRequest{
		IMessage: interfaces.NewMessageRequest(interfaces.MessageTypeFileSystemHelpRequest, ack),
		path:     path,
	}
}

// Path returns the file system path associated with the MessageFileSystemHelpRequest instance.
func (m *MessageFileSystemHelpRequest) Path() string {
	return m.path
}

// CreateResponse generates a response message with the provided result and error, linking it to the current request.
func (m *MessageFileSystemHelpRequest) CreateResponse(result string, err error) {
	msg := &MessageFileSystemHelpResponse{IMessage: m.Clone(), result: result, err: err}
	m.SetResponse(interfaces.MessageTypeFileSystemHelpResponse, msg)
}

// GetResponse retrieves the stored result and error from the MessageFileSystemHelpResponse instance.
func (m *MessageFileSystemHelpResponse) GetResponse() (string, error) {
	return m.result, m.err
}
