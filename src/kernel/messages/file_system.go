package messages

import "github.com/markel1974/c64emu/src/kernel/interfaces"

// --- Messaggi di Richiesta ---

type MessageCWDSet struct {
	interfaces.Message
	Path string
}

func NewMessageCWDSet(path string) *MessageCWDSet {
	return &MessageCWDSet{
		Message: *interfaces.NewMessage(interfaces.MessageTypeCWDSet),
		Path:    path,
	}
}

type MessageCWDGetRequest struct {
	interfaces.Message
}

func NewMessageCWDGetRequest(originatorPID int) *MessageCWDGetRequest {
	return &MessageCWDGetRequest{
		Message: *interfaces.NewMessage(interfaces.MessageTypeCWDGetRequest),
	}
}

type MessageCWDNameRequest struct {
	interfaces.Message
}

func NewMessageCWDNameRequest(originatorPID int) *MessageCWDNameRequest {
	return &MessageCWDNameRequest{
		Message: *interfaces.NewMessage(interfaces.MessageTypeCWDNameRequest),
	}
}

type MessageCWDDirectoryListingRequest struct {
	interfaces.Message
}

func NewMessageCWDDirectoryListingRequest() *MessageCWDDirectoryListingRequest {
	return &MessageCWDDirectoryListingRequest{
		Message: *interfaces.NewMessage(interfaces.MessageTypeCWDDirectoryListingRequest),
	}
}

type MessageFileSystemSuggestionRequest struct {
	interfaces.Message
	In     string
	Cursor int
}

func NewMessageFileSystemSuggestionRequest(in string, cursor int) *MessageFileSystemSuggestionRequest {
	return &MessageFileSystemSuggestionRequest{
		Message: *interfaces.NewMessage(interfaces.MessageTypeFileSystemSuggestionRequest),
		In:      in,
		Cursor:  cursor,
	}
}

type MessageFileSystemHelpRequest struct {
	interfaces.Message
	Arg string
}

func NewMessageFileSystemHelpRequest(arg string) *MessageFileSystemHelpRequest {
	return &MessageFileSystemHelpRequest{
		Message: *interfaces.NewMessage(interfaces.MessageTypeFileSystemHelpRequest),
		Arg:     arg,
	}
}

type MessageCWDGetResponse struct {
	interfaces.Message
	Path string
}

func NewMessageCWDGetResponse(path string) *MessageCWDGetResponse {
	return &MessageCWDGetResponse{
		Message: *interfaces.NewMessage(interfaces.MessageTypeCWDGetResponse),
		Path:    path,
	}
}

type MessageCWDNameResponse struct {
	interfaces.Message
	Name string
}

func NewMessageCWDNameResponse(name string) *MessageCWDNameResponse {
	return &MessageCWDNameResponse{
		Message: *interfaces.NewMessage(interfaces.MessageTypeCWDNameResponse),
		Name:    name,
	}
}

type MessageCWDDirectoryListingResponse struct {
	interfaces.Message
	Listing []string
}

func NewMessageCWDDirectoryListingResponse(listing []string) *MessageCWDDirectoryListingResponse {
	return &MessageCWDDirectoryListingResponse{
		Message: *interfaces.NewMessage(interfaces.MessageTypeCWDDirectoryListingResponse),
		Listing: listing,
	}
}

type MessageFileSystemSuggestionResponse struct {
	interfaces.Message
	Prefix      string
	Suggestions []string
	Found       bool
}

func NewMessageFileSystemSuggestionResponse(prefix string, suggestions []string, found bool) *MessageFileSystemSuggestionResponse {
	return &MessageFileSystemSuggestionResponse{
		Message:     *interfaces.NewMessage(interfaces.MessageTypeFileSystemSuggestionResponse),
		Prefix:      prefix,
		Suggestions: suggestions,
		Found:       found,
	}
}
