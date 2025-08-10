package messages

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
)

// MessageError represents a specialized message that encapsulates an error within the messaging system.
type MessageError struct {
	interfaces.IMessage
	err error
}

// NewMessageError creates and returns a new MessageError instance with the provided error and sets its type to MessageTypeError.
func NewMessageError(source int, destination int, err error) *MessageError {
	return &MessageError{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeError),
		err:      err,
	}
}

// Error returns the underlying error associated with the MessageError instance.
func (m *MessageError) Error() error {
	return m.err
}

// MessageQuit represents a specific IMessage implementation used to signal a quit or termination event in the system.
type MessageQuit struct {
	interfaces.IMessage
}

// NewMessageQuit creates and returns a new instance of MessageQuit with type MessageTypeQuit and no acknowledgment.
func NewMessageQuit(source int, destination int) *MessageQuit {
	return &MessageQuit{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeQuit),
	}
}

// MessageRead represents a read operation message encapsulating the key type, character data, and broadcast behavior.
type MessageRead struct {
	interfaces.IMessage
	kind      interfaces.KeyType
	data      rune
	broadcast bool
}

// NewMessageRead creates and initializes a new MessageRead instance with the specified key type, data, and broadcast flag.
func NewMessageRead(source int, destination int, kind interfaces.KeyType, data rune, broadcast bool) *MessageRead {
	return &MessageRead{
		IMessage:  interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeRead),
		kind:      kind,
		data:      data,
		broadcast: broadcast,
	}
}

// Kind returns the key type associated with the MessageRead instance.
func (m *MessageRead) Kind() interfaces.KeyType {
	return m.kind
}

// Data returns the rune data associated with the MessageRead instance.
func (m *MessageRead) Data() rune {
	return m.data
}

// Broadcast returns whether the message is marked as a broadcast.
func (m *MessageRead) Broadcast() bool {
	return m.broadcast
}

// MessageProcessExec represents a message encapsulating process execution details.
// It includes the IMessage interface and a command line to execute.
type MessageProcessExec struct {
	interfaces.IMessage
	line string
}

// NewMessageProcessExec creates a new MessageProcessExec instance with the provided command line and a process execution type.
func NewMessageProcessExec(source int, destination int, line string) *MessageProcessExec {
	return &MessageProcessExec{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeProcessExec),
		line:     line,
	}
}

// Line retrieves the stored line string from the MessageProcessExec instance.
func (m *MessageProcessExec) Line() string {
	return m.line
}

// MessageProcessStart represents a message to initiate the start of a process with specified arguments.
type MessageProcessStart struct {
	interfaces.IMessage
	args []string
}

// NewMessageProcessStart creates a new MessageProcessStart instance with the provided arguments.
func NewMessageProcessStart(source int, destination int, args []string) *MessageProcessStart {
	return &MessageProcessStart{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeProcessStart),
		args:     args,
	}
}

// MessageProcessActivate represents a message type for activating an existing process in the messaging system.
type MessageProcessActivate struct {
	interfaces.IMessage
}

// NewMessageProcessActivate creates a new message of type MessageTypeProcessActivate without acknowledgment requirements.
func NewMessageProcessActivate(source int, destination int) *MessageProcessActivate {
	return &MessageProcessActivate{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeProcessActivate),
	}
}

// Args returns the list of arguments associated with the MessageProcessStart instance.
func (m *MessageProcessStart) Args() []string {
	return m.args
}

// MessageProcessKill represents a message used to request termination of a process, identified by its process ID (pid).
type MessageProcessKill struct {
	interfaces.IMessage
	pid int
}

// NewMessageProcessKill creates a new MessageProcessKill instance to signal termination of a specific process by its pid.
func NewMessageProcessKill(source int, destination int, pid int) *MessageProcessKill {
	return &MessageProcessKill{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeProcessKill),
		pid:      pid,
	}
}

// MessageProcessKillAll represents a message designed to instruct the termination of all processes or a subset by name.
type MessageProcessKillAll struct {
	interfaces.IMessage
	name string
}

// NewMessageProcessKillAll creates a new MessageProcessKillAll instance for terminating all processes by name.
func NewMessageProcessKillAll(source int, destination int, name string) *MessageProcessKillAll {
	return &MessageProcessKillAll{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeProcessKillAll),
		name:     name,
	}
}

// Name retrieves the name of the MessageProcessKillAll instance.
func (m *MessageProcessKillAll) Name() string {
	return m.name
}

// MessageProcessKillForeground is a message type used to signal the termination of the currently active foreground process.
type MessageProcessKillForeground struct {
	interfaces.IMessage
}

// NewMessageProcessKillForeground creates a new MessageProcessKillForeground instance with a no-acknowledgment message type.
func NewMessageProcessKillForeground(source int, destination int) *MessageProcessKillForeground {
	return &MessageProcessKillForeground{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeProcessKillForeground),
	}
}

// MessageProcessExit represents a message indicating a process exit event, containing its process ID and message metadata.
type MessageProcessExit struct {
	interfaces.IMessage
	pid int
}

// NewMessageMessageProcessExit creates and returns a new MessageProcessExit instance with the MessageTypeProcessExit type.
func NewMessageMessageProcessExit(source int, destination int) *MessageProcessExit {
	return &MessageProcessExit{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeProcessExit),
	}
}

// MessageProcessSetSelfForeground represents a message to set a process as the foreground process for itself.
// It embeds the IMessage interface for standard message operations, including type retrieval and response handling.
type MessageProcessSetSelfForeground struct {
	interfaces.IMessage
}

// NewMessageProcessSetSelfForeground creates and returns a new MessageProcessSetSelfForeground instance.
// It uses MessageTypeProcessSetSelfForeground as its type and does not initialize an acknowledgment channel.
func NewMessageProcessSetSelfForeground(source int, destination int) *MessageProcessSetSelfForeground {
	return &MessageProcessSetSelfForeground{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeProcessSetSelfForeground),
	}
}

// MessageProcessSetForeground represents a message that sets a process as the foreground process.
// It embeds the IMessage interface for messaging capabilities and contains the process ID.
type MessageProcessSetForeground struct {
	interfaces.IMessage
	pid int
}

// NewMessageProcessSetForeground creates a new MessageProcessSetForeground with the specified process ID (PID).
func NewMessageProcessSetForeground(source int, destination int, pid int) *MessageProcessSetForeground {
	return &MessageProcessSetForeground{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeProcessSetForeground),
		pid:      pid,
	}
}

// PID returns the process identifier associated with the MessageProcessSetForeground instance.
func (m *MessageProcessSetForeground) PID() int {
	return m.pid
}

// MessageProcessListRequest represents a message type used to request the current list of active processes.
// It implements the IMessage interface, allowing message operations such as setting type, originator, and acknowledgment.
type MessageProcessListRequest struct {
	interfaces.IMessage
}

// MessageProcessListResponse represents a response containing a list of process descriptions associated with a message.
type MessageProcessListResponse struct {
	interfaces.IMessage
	processes []*interfaces.ProcessDescription
}

// NewMessageProcessListRequest creates a new MessageProcessListRequest with the provided acknowledgment channel.
func NewMessageProcessListRequest(source int, destination int, ackChan chan bool) *MessageProcessListRequest {
	return &MessageProcessListRequest{
		IMessage: interfaces.NewMessageRequest(source, destination, interfaces.MessageTypeProcessListRequest, ackChan),
	}
}

// CreateResponse constructs a MessageProcessListResponse, clones the original message, and sets it as a response.
func (m *MessageProcessListRequest) CreateResponse(responder int, processes []*interfaces.ProcessDescription) {
	base := m.PrepareResponse(responder, interfaces.MessageTypeProcessListResponse)
	msg := &MessageProcessListResponse{IMessage: base, processes: processes}
	m.SetResponse(msg)
}

// GetResponse returns the list of ProcessDescription objects associated with the MessageProcessListResponse instance.
func (m *MessageProcessListResponse) GetResponse() []*interfaces.ProcessDescription {
	return m.processes
}

// MessageProcessIsRunningRequest represents a request to verify if a specific process is currently running.
// It embeds interfaces.IMessage and includes a process ID (verifyPid) to check.
type MessageProcessIsRunningRequest struct {
	interfaces.IMessage
	verifyPid int
}

// MessageProcessIsRunningResponse represents the response message confirming whether a process is currently running.
type MessageProcessIsRunningResponse struct {
	interfaces.IMessage
	result bool
}

// NewMessageProcessIsRunningRequest creates a new MessageProcessIsRunningRequest with a verify PID and acknowledgment channel.
func NewMessageProcessIsRunningRequest(source int, destination, verifyPid int, ackChan chan bool) *MessageProcessIsRunningRequest {
	return &MessageProcessIsRunningRequest{
		IMessage:  interfaces.NewMessageRequest(source, destination, interfaces.MessageTypeProcessIsRunningRequest, ackChan),
		verifyPid: verifyPid,
	}
}

// VerifyPID retrieves the process ID associated with the MessageProcessIsRunningRequest instance.
func (m *MessageProcessIsRunningRequest) VerifyPID() int {
	return m.verifyPid
}

// CreateResponse generates a MessageProcessIsRunningResponse based on the result and sets it as the response for the request.
func (m *MessageProcessIsRunningRequest) CreateResponse(responder int, result bool) {
	base := m.PrepareResponse(responder, interfaces.MessageTypeProcessIsRunningResponse)
	msg := &MessageProcessIsRunningResponse{IMessage: base, result: result}
	m.SetResponse(msg)
}

// GetResponse returns the boolean result indicating if the message process is running.
func (m *MessageProcessIsRunningResponse) GetResponse() bool {
	return m.result
}

// MessageNotifyProcessCreate is a message type used to notify the creation of a new process.
// It includes the process ID and name of the created process as fields.
// Implements the IMessage interface for messaging system integration.
type MessageNotifyProcessCreate struct {
	interfaces.IMessage
	createdPID int
	name       string
}

// NewMessageNotifyProcessCreate creates a new MessageNotifyProcessCreate instance with the given process ID and name.
func NewMessageNotifyProcessCreate(source int, destination int, createdPID int, name string) *MessageNotifyProcessCreate {
	return &MessageNotifyProcessCreate{
		IMessage:   interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeNotifyProcessCreate),
		createdPID: createdPID,
		name:       name,
	}
}

// CreatedPID returns the process ID of the created process stored in the MessageNotifyProcessCreate instance.
func (m *MessageNotifyProcessCreate) CreatedPID() int {
	return m.createdPID
}

// Name returns the name string associated with the MessageNotifyProcessCreate instance.
func (m *MessageNotifyProcessCreate) Name() string {
	return m.name
}

// MessageNotifyProcessTerminate notifies when a specific process has been terminated in the system.
// It implements the IMessage interface for inter-process communication.
// The terminatedPID field identifies the process ID of the terminated process.
type MessageNotifyProcessTerminate struct {
	interfaces.IMessage
	terminatedPID int
}

// NewMessageMessageNotifyProcessTerminate creates a new MessageNotifyProcessTerminate for a terminated process with the given PID.
func NewMessageMessageNotifyProcessTerminate(source int, destination int, terminatedPID int) *MessageNotifyProcessTerminate {
	return &MessageNotifyProcessTerminate{
		IMessage:      interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeNotifyProcessTerminate),
		terminatedPID: terminatedPID,
	}
}

// TerminatedPID returns the PID of the process that has terminated.
func (m *MessageNotifyProcessTerminate) TerminatedPID() int {
	return m.terminatedPID
}

// MessageNotifyProcessForeground represents a message indicating that a process has entered the foreground.
// It embeds the IMessage interface and includes the process ID of the foreground process.
type MessageNotifyProcessForeground struct {
	interfaces.IMessage
	foregroundPID int
}

// NewMessageNotifyProcessForeground creates a new instance of MessageNotifyProcessForeground with the specified foreground PID.
func NewMessageNotifyProcessForeground(source int, destination int, foregroundPID int) *MessageNotifyProcessForeground {
	return &MessageNotifyProcessForeground{
		IMessage:      interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeNotifyProcessForeground),
		foregroundPID: foregroundPID,
	}
}

// ForegroundPID returns the process ID of the foreground process associated with the message.
func (m *MessageNotifyProcessForeground) ForegroundPID() int {
	return m.foregroundPID
}
