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

func NewMessageCWDGetRequest(router interfaces.IRouter, ack chan bool) *MessageCWDGetRequest {
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
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeCWDDirectoryListingRequest),
	}
}

type MessageFileSystemSuggestionRequest struct {
	interfaces.Message
	In     string
	Cursor int
}

func NewMessageFileSystemSuggestionRequest(router interfaces.IRouter, in string, cursor int) *MessageFileSystemSuggestionRequest {
	return &MessageFileSystemSuggestionRequest{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeFileSystemSuggestionRequest),
		In:      in,
		Cursor:  cursor,
	}
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

type MessageCWDDirectoryListingResponse struct {
	interfaces.Message
	Listing []string
}

func NewMessageCWDDirectoryListingResponse(router interfaces.IRouter, listing []string) *MessageCWDDirectoryListingResponse {
	return &MessageCWDDirectoryListingResponse{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeCWDDirectoryListingResponse),
		Listing: listing,
	}
}

type MessageFileSystemSuggestionResponse struct {
	interfaces.Message
	Prefix      string
	Suggestions []string
	Found       bool
}

func NewMessageFileSystemSuggestionResponse(router interfaces.IRouter, prefix string, suggestions []string, found bool) *MessageFileSystemSuggestionResponse {
	return &MessageFileSystemSuggestionResponse{
		Message:     *interfaces.NewMessage(router, interfaces.MessageTypeFileSystemSuggestionResponse),
		Prefix:      prefix,
		Suggestions: suggestions,
		Found:       found,
	}
}
