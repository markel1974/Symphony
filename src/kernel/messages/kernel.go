package messages

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
)

// MessageQuit represents a message signaling a quit operation or system termination.
// It embeds the base Message type and sets its kind to MessageTypeQuit.
type MessageQuit struct {
	interfaces.Message
}

// NewMessageQuit creates and returns a new MessageQuit instance with the MessageType set to MessageTypeQuit.
func NewMessageQuit() *MessageQuit {
	return &MessageQuit{
		Message: *interfaces.NewMessage(interfaces.MessageTypeQuit),
	}
}

// MessageRead represents a specific type of Message containing read operation data. It embeds Message and includes a data field.
type MessageRead struct {
	interfaces.Message
	kind interfaces.KeyType
	data rune
}

// NewMessageRead creates a new MessageRead instance with provided data and limits its length to n if necessary.
func NewMessageRead(kind interfaces.KeyType, data rune) *MessageRead {
	return &MessageRead{
		Message: *interfaces.NewMessage(interfaces.MessageTypeRead),
		kind:    kind,
		data:    data,
	}
}

// Kind returns the key type associated with the MessageRead instance.
func (m *MessageRead) Kind() interfaces.KeyType {
	return m.kind
}

// Data returns the data payload contained within the MessageRead instance.
func (m *MessageRead) Data() rune {
	return m.data
}

// MessageProcessExec represents a message type designed to execute a process with a given command line.
type MessageProcessExec struct {
	interfaces.Message
	Line string
}

// NewMessageProcessExec creates a new MessageProcessExec instance with a provided line and MessageTypeProcessExec type.
func NewMessageProcessExec(line string) *MessageProcessExec {
	return &MessageProcessExec{
		Message: *interfaces.NewMessage(interfaces.MessageTypeProcessExec),
		Line:    line,
	}
}

// MessageProcessKill represents a message used to signal the termination of a process identified by its PID.
// It embeds the base Message type and includes a PID field for specifying the process to be killed.
type MessageProcessKill struct {
	interfaces.Message
	PID int
}

// NewMessageProcessKill creates a new MessageProcessKill object with the specified PID and a ProcessKill message type.
func NewMessageProcessKill(pid int) *MessageProcessKill {
	return &MessageProcessKill{
		Message: *interfaces.NewMessage(interfaces.MessageTypeProcessKill),
		PID:     pid,
	}
}

// MessageProcessKillAll represents a message type used to signal the termination of all processes in the system.
// It embeds the Message struct and includes a Name field to specify the associated process group or identifier.
type MessageProcessKillAll struct {
	interfaces.Message
	Name string
}

// NewMessageProcessKillAll creates a new MessageProcessKillAll instance with the specified name and appropriate message type.
func NewMessageProcessKillAll(name string) *MessageProcessKillAll {
	return &MessageProcessKillAll{
		Message: *interfaces.NewMessage(interfaces.MessageTypeProcessKillAll),
		Name:    name,
	}
}

// MessageProcessSetForeground represents a message to set a specific process to the foreground.
// It embeds interfaces.Message and includes a PID field for the target process identifier.
type MessageProcessSetForeground struct {
	interfaces.Message
	PID int
}

// NewMessageProcessSetForeground creates and returns a new MessageProcessSetForeground with the specified process ID.
// It initializes the embedded Message with the MessageTypeProcessSetForeground constant.
func NewMessageProcessSetForeground(pid int) *MessageProcessSetForeground {
	return &MessageProcessSetForeground{
		Message: *interfaces.NewMessage(interfaces.MessageTypeProcessSetForeground),
		PID:     pid,
	}
}

// MessageProcessListRequest represents a request message to retrieve a list of processes from the system.
// It embeds Message for base message functionality and uses ReplyTo to specify the recipient for the response.
type MessageProcessListRequest struct {
	interfaces.Message
	replyTo   interfaces.IReceiver
	processes []*interfaces.ProcessDescription
	ackChan   chan bool
}

// NewMessageProcessList creates a new MessageProcessListRequest with specified replyTo as an IReceiver.
func NewMessageProcessList(replyTo interfaces.IReceiver, ackChan chan bool) *MessageProcessListRequest {
	return &MessageProcessListRequest{
		Message: *interfaces.NewMessage(interfaces.MessageTypeProcessListRequest),
		replyTo: replyTo,
		ackChan: ackChan,
	}
}

// Respond sets the list of processes and sends the message to the specified recipient via ReplyTo.
func (m *MessageProcessListRequest) Respond(processes []*interfaces.ProcessDescription) {
	m.processes = processes
	m.replyTo.PostMessage(m)
}

// Processes returns the list of ProcessDescription objects associated with the MessageProcessListResponse instance.
func (m *MessageProcessListRequest) Processes() []*interfaces.ProcessDescription {
	return m.processes
}

// Ack signals the acknowledgment channel and closes it to indicate task completion or unblock synchronous calls.
func (m *MessageProcessListRequest) Ack() {
	m.ackChan <- true
	close(m.ackChan)
}
