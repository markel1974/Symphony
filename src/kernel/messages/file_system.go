package messages

import "github.com/markel1974/c64emu/src/kernel/interfaces"

type MessageCWDSet struct {
	interfaces.Message
	Path string
}

func NewMessageCWDSet(router interfaces.IRouter, path string) *MessageCWDSet {
	return &MessageCWDSet{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeCWDSet),
		Path:    path,
	}
}

type MessageCWDGetRequest struct {
	interfaces.Message
	result string
	ack    chan bool
}

func NewMessageCWDGet(router interfaces.IRouter, ack chan bool) *MessageCWDGetRequest {
	return &MessageCWDGetRequest{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeCWDGetRequest),
		ack:     ack,
	}
}

func (m *MessageCWDGetRequest) SetResult(result string) {
	m.MakeResponse()
	m.result = result
}

func (m *MessageCWDGetRequest) Result() string {
	return m.result
}

func (m *MessageCWDGetRequest) Ack() bool {
	m.ack <- true
	return true
}

type MessageCWDNameRequest struct {
	interfaces.Message
}

func NewMessageCWDNameRequest(router interfaces.IRouter) *MessageCWDNameRequest {
	return &MessageCWDNameRequest{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeCWDNameRequest),
	}
}

type MessageCWDDirectoryListingRequest struct {
	interfaces.Message
}

func NewMessageCWDDirectoryListingRequest(router interfaces.IRouter) *MessageCWDDirectoryListingRequest {
	return &MessageCWDDirectoryListingRequest{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeCWDDirectoryListing),
	}
}

// MessageFileSystemSuggestionRequest represents a message for requesting filesystem suggestions based on user input.
// It captures the input string, cursor position, and provides an acknowledgment channel for processing confirmation.
// The structure enables setting and retrieving responses containing prefix, suggestions, and their validity.
type MessageFileSystemSuggestionRequest struct {
	interfaces.Message
	in         string
	cursor     int
	ack        chan bool
	prefix     string
	suggestion []string
	valid      bool
}

// NewMessageFileSystemSuggestion creates a new MessageFileSystemSuggestionRequest with provided router, input text, cursor position, and acknowledgment channel.
func NewMessageFileSystemSuggestion(router interfaces.IRouter, in string, cursor int, ack chan bool) *MessageFileSystemSuggestionRequest {
	return &MessageFileSystemSuggestionRequest{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeFileSystemSuggestion),
		in:      in,
		cursor:  cursor,
		ack:     ack,
	}
}

// Ack acknowledges the message by signaling through the ack channel and returns true to confirm the acknowledgment.
func (m *MessageFileSystemSuggestionRequest) Ack() bool {
	m.ack <- true
	return true
}

// In returns the `in` field of the MessageFileSystemSuggestionRequest instance as a string.
func (m *MessageFileSystemSuggestionRequest) In() string {
	return m.in
}

// Cursor returns the current cursor position as an integer.
func (m *MessageFileSystemSuggestionRequest) Cursor() int {
	return m.cursor
}

// SetResponse sets the response data with the given prefix, suggestions, and validation state.
func (m *MessageFileSystemSuggestionRequest) SetResponse(prefix string, suggestion []string, valid bool) {
	m.MakeResponse()
	m.prefix = prefix
	m.suggestion = suggestion
	m.valid = valid
}

// GetResponse retrieves the suggestion result, including the prefix, suggestions list, and validity flag.
func (m *MessageFileSystemSuggestionRequest) GetResponse() (prefix string, suggestion []string, valid bool) {
	return m.prefix, m.suggestion, m.valid
}

type MessageFileSystemHelpRequest struct {
	interfaces.Message
	Arg string
}

func NewMessageFileSystemHelpRequest(router interfaces.IRouter, arg string) *MessageFileSystemHelpRequest {
	return &MessageFileSystemHelpRequest{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeFileSystemHelpRequest),
		Arg:     arg,
	}
}
