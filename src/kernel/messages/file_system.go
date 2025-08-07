package messages

import "github.com/markel1974/c64emu/src/kernel/interfaces"

type MessageFileSystemFindResponse struct {
	interfaces.Message
	parent    interfaces.IRouter
	line      string
	protected bool
	cmd       interfaces.ICommand
	args      []string
	err       error
}

type MessageFileSystemFindRequest struct {
	interfaces.Message
	parent    interfaces.IRouter
	line      string
	protected bool
}

func NewMessageFileSystemFindRequest(router interfaces.IRouter, parent interfaces.IRouter, line string, protected bool) *MessageFileSystemFindRequest {
	return &MessageFileSystemFindRequest{
		Message:   *interfaces.NewMessage(router, interfaces.MessageTypeFileSystemFindRequest),
		parent:    parent,
		line:      line,
		protected: protected,
	}
}

func (m *MessageFileSystemFindRequest) Line() string {
	return m.line
}

func (m *MessageFileSystemFindRequest) CreateResponse(cmd interfaces.ICommand, args []string, err error) *MessageFileSystemFindResponse {
	return &MessageFileSystemFindResponse{
		Message:   *interfaces.NewMessage(m.Router(), interfaces.MessageTypeFileSystemFindResponse),
		parent:    m.parent,
		line:      m.line,
		protected: m.protected,
		cmd:       cmd,
		args:      args,
		err:       err,
	}
}

func (m *MessageFileSystemFindResponse) Parent() interfaces.IRouter {
	return m.parent
}

func (m *MessageFileSystemFindResponse) Line() string {
	return m.line
}

func (m *MessageFileSystemFindResponse) Protected() bool {
	return m.protected
}

func (m *MessageFileSystemFindResponse) GetResult() (interfaces.ICommand, []string, error) {
	return m.cmd, m.args, m.err
}

type MessageFileSystemCWDSet struct {
	interfaces.Message
	path   string
	ack    chan bool
	result bool
}

func NewMessageFileSystemCWDSet(router interfaces.IRouter, path string, ack chan bool) *MessageFileSystemCWDSet {
	return &MessageFileSystemCWDSet{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeFileSystemCWDSet),
		path:    path,
		ack:     ack,
	}
}

func (m *MessageFileSystemCWDSet) Ack() bool {
	m.ack <- true
	return true
}

func (m *MessageFileSystemCWDSet) Path() string {
	return m.path
}

func (m *MessageFileSystemCWDSet) SetResult(result bool) {
	m.MakeResponse()
	m.result = result
}

func (m *MessageFileSystemCWDSet) GetResult() bool {
	return m.result
}

type MessageFileSystemCWDGet struct {
	interfaces.Message
	result string
	ack    chan bool
}

func NewMessageFileSystemCWDGet(router interfaces.IRouter, ack chan bool) *MessageFileSystemCWDGet {
	return &MessageFileSystemCWDGet{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeFileSystemCWDGet),
		ack:     ack,
	}
}

func (m *MessageFileSystemCWDGet) SetResult(result string) {
	m.MakeResponse()
	m.result = result
}

func (m *MessageFileSystemCWDGet) GetResult() string {
	return m.result
}

func (m *MessageFileSystemCWDGet) Ack() bool {
	m.ack <- true
	return true
}

type MessageFileSystemCWDPath struct {
	interfaces.Message
	ack    chan bool
	result string
}

func NewMessageFileSystemCWDPath(router interfaces.IRouter, ack chan bool) *MessageFileSystemCWDPath {
	return &MessageFileSystemCWDPath{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeFileSystemCWDPath),
		ack:     ack,
	}
}

func (m *MessageFileSystemCWDPath) Ack() bool {
	m.ack <- true
	return true
}

func (m *MessageFileSystemCWDPath) SetResult(result string) {
	m.MakeResponse()
	m.result = result
}

func (m *MessageFileSystemCWDPath) GetResult() string {
	return m.result
}

type MessageFileSystemCWDDirectoryListing struct {
	interfaces.Message
	ack    chan bool
	result []string
}

func NewMessageFileSystemCWDDirectoryListing(router interfaces.IRouter, ack chan bool) *MessageFileSystemCWDDirectoryListing {
	return &MessageFileSystemCWDDirectoryListing{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeFileSystemCWDDirectoryListing),
		ack:     ack,
	}
}

func (m *MessageFileSystemCWDDirectoryListing) Ack() bool {
	m.ack <- true
	return true
}

func (m *MessageFileSystemCWDDirectoryListing) SetResult(result []string) {
	m.MakeResponse()
	m.result = result
}

func (m *MessageFileSystemCWDDirectoryListing) GetResult() []string {
	return m.result
}

// MessageFileSystemSuggestion represents a message for requesting filesystem suggestions based on user input.
// It captures the input string, cursor position, and provides an acknowledgment channel for processing confirmation.
// The structure enables setting and retrieving responses containing prefix, suggestions, and their validity.
type MessageFileSystemSuggestion struct {
	interfaces.Message
	in         string
	cursor     int
	ack        chan bool
	prefix     string
	suggestion []string
	valid      bool
}

// NewMessageFileSystemSuggestion creates a new MessageFileSystemSuggestion with provided router, input text, cursor position, and acknowledgment channel.
func NewMessageFileSystemSuggestion(router interfaces.IRouter, in string, cursor int, ack chan bool) *MessageFileSystemSuggestion {
	return &MessageFileSystemSuggestion{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeFileSystemSuggestion),
		in:      in,
		cursor:  cursor,
		ack:     ack,
	}
}

// Ack acknowledges the message by signaling through the ack channel and returns true to confirm the acknowledgment.
func (m *MessageFileSystemSuggestion) Ack() bool {
	m.ack <- true
	return true
}

// In returns the `in` field of the MessageFileSystemSuggestion instance as a string.
func (m *MessageFileSystemSuggestion) In() string {
	return m.in
}

// Cursor returns the current cursor position as an integer.
func (m *MessageFileSystemSuggestion) Cursor() int {
	return m.cursor
}

// SetResponse sets the response data with the given prefix, suggestions, and validation state.
func (m *MessageFileSystemSuggestion) SetResponse(prefix string, suggestion []string, valid bool) {
	m.MakeResponse()
	m.prefix = prefix
	m.suggestion = suggestion
	m.valid = valid
}

// GetResponse retrieves the suggestion result, including the prefix, suggestions list, and validity flag.
func (m *MessageFileSystemSuggestion) GetResponse() (prefix string, suggestion []string, valid bool) {
	return m.prefix, m.suggestion, m.valid
}

type MessageFileSystemHelp struct {
	interfaces.Message
	path   string
	ack    chan bool
	result string
	err    error
}

func NewMessageFileSystemHelp(router interfaces.IRouter, path string, ack chan bool) *MessageFileSystemHelp {
	return &MessageFileSystemHelp{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeFileSystemHelp),
		path:    path,
		ack:     ack,
	}
}

func (m *MessageFileSystemHelp) Ack() bool {
	m.ack <- true
	return true
}

func (m *MessageFileSystemHelp) Path() string {
	return m.path
}

func (m *MessageFileSystemHelp) SetResponse(result string, err error) {
	m.MakeResponse()
	m.result = result
	m.err = err
}

func (m *MessageFileSystemHelp) GetResponse() (string, error) {
	return m.result, m.err
}
