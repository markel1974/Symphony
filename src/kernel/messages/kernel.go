package messages

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
)

type MessageError struct {
	interfaces.Message
	err error
}

// NewMessageError creates and returns a new MessageQuit instance with the MessageType set to MessageTypeQuit.
func NewMessageError(originatorPID int, err error) *MessageError {
	return &MessageError{
		Message: *interfaces.NewMessage(originatorPID, interfaces.MessageTypeError),
		err:     err,
	}
}

// Error returns the error associated with the MessageError instance.
func (m *MessageError) Error() error {
	return m.err
}

// MessageQuit represents a message signaling a quit operation or system termination.
// It embeds the base Message type and sets its kind to MessageTypeQuit.
type MessageQuit struct {
	interfaces.Message
}

// NewMessageQuit creates and returns a new MessageQuit instance with the MessageType set to MessageTypeQuit.
func NewMessageQuit(originatorPID int) *MessageQuit {
	return &MessageQuit{
		Message: *interfaces.NewMessage(originatorPID, interfaces.MessageTypeQuit),
	}
}

// MessageRead represents a specific type of Message containing read operation data. It embeds Message and includes a data field.
type MessageRead struct {
	interfaces.Message
	kind      interfaces.KeyType
	data      rune
	broadcast bool
}

// NewMessageRead creates a new MessageRead instance with provided data and limits its length to n if necessary.
func NewMessageRead(originatorPID int, kind interfaces.KeyType, data rune, broadcast bool) *MessageRead {
	return &MessageRead{
		Message:   *interfaces.NewMessage(originatorPID, interfaces.MessageTypeRead),
		kind:      kind,
		data:      data,
		broadcast: broadcast,
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

// Broadcast returns a boolean indicating whether the message should be broadcasted.
func (m *MessageRead) Broadcast() bool {
	return m.broadcast
}

// MessageProcessExec represents a message type designed to execute a process with a given command line.
type MessageProcessExec struct {
	interfaces.Message
	line string
}

// NewMessageProcessExec creates a new MessageProcessExec instance with a provided line and MessageTypeProcessExec type.
func NewMessageProcessExec(originatorPID int, line string) *MessageProcessExec {
	return &MessageProcessExec{
		Message: *interfaces.NewMessage(originatorPID, interfaces.MessageTypeProcessExec),
		line:    line,
	}
}

func (m *MessageProcessExec) Line() string {
	return m.line
}

// MessageProcessStart represents a message type designed to execute a process with a given command line.
type MessageProcessStart struct {
	interfaces.Message
	args []string
}

// NewMessageProcessStart creates a new MessageProcessExec instance with a provided line and MessageTypeProcessExec type.
func NewMessageProcessStart(originatorPID int, args []string) *MessageProcessStart {
	return &MessageProcessStart{
		Message: *interfaces.NewMessage(originatorPID, interfaces.MessageTypeProcessStart),
		args:    args,
	}
}

// MessageProcessActivate represents a message type designed to execute a process with a given command line.
type MessageProcessActivate struct {
	interfaces.Message
}

// NewMessageProcessActivate creates and returns a new instance of MessageProcessActivate initialized with MessageTypeProcessActivate.
func NewMessageProcessActivate(originatorPID int) *MessageProcessActivate {
	return &MessageProcessActivate{
		Message: *interfaces.NewMessage(originatorPID, interfaces.MessageTypeProcessActivate),
	}
}

// Args returns the slice of arguments associated with the MessageProcessStart instance.
func (m *MessageProcessStart) Args() []string {
	return m.args
}

// MessageProcessKill represents a message used to signal the termination of a process identified by its PID.
// It embeds the base Message type and includes a PID field for specifying the process to be killed.
type MessageProcessKill struct {
	interfaces.Message
	pid int
}

// NewMessageProcessKill creates a new MessageProcessKill object with the specified PID and a ProcessKill message type.
func NewMessageProcessKill(originatorPID int, pid int) *MessageProcessKill {
	return &MessageProcessKill{
		Message: *interfaces.NewMessage(originatorPID, interfaces.MessageTypeProcessKill),
		pid:     pid,
	}
}

// MessageProcessKillAll represents a message type used to signal the termination of all processes in the system.
// It embeds the Message struct and includes a Name field to specify the associated process group or identifier.
type MessageProcessKillAll struct {
	interfaces.Message
	name string
}

// NewMessageProcessKillAll creates a new MessageProcessKillAll instance with the specified name and appropriate message type.
func NewMessageProcessKillAll(originatorPID int, name string) *MessageProcessKillAll {
	return &MessageProcessKillAll{
		Message: *interfaces.NewMessage(originatorPID, interfaces.MessageTypeProcessKillAll),
		name:    name,
	}
}

// Name returns the value of the name field associated with the MessageProcessKillAll instance.
func (m *MessageProcessKillAll) Name() string {
	return m.name
}

// MessageProcessKillForeground represents a message that signals the system to terminate all foreground processes.
// It extends the base Message type to specialize for foreground process kill behavior.
type MessageProcessKillForeground struct {
	interfaces.Message
}

// NewMessageProcessKillForeground creates a new MessageProcessKillForeground instance with a router and MessageTypeProcessKillForeground.
func NewMessageProcessKillForeground(originatorPID int) *MessageProcessKillForeground {
	return &MessageProcessKillForeground{
		Message: *interfaces.NewMessage(originatorPID, interfaces.MessageTypeProcessKillForeground),
	}
}

// MessageProcessExit represents a message to set a specific process to the foreground.
// It embeds interfaces.Message and includes a PID field for the target process identifier.
type MessageProcessExit struct {
	interfaces.Message
	pid int
}

// NewMessageMessageProcessExit creates and returns a new MessageProcessExit instance with the ProcessExit message type.
func NewMessageMessageProcessExit(originatorPID int) *MessageProcessExit {
	return &MessageProcessExit{
		Message: *interfaces.NewMessage(originatorPID, interfaces.MessageTypeProcessExit),
	}
}

// MessageProcessSetForeground represents a message to set a specific process to the foreground.
// It embeds interfaces.Message and includes a PID field for the target process identifier.
type MessageProcessSetForeground struct {
	interfaces.Message
	pid int
}

// NewMessageProcessSetForeground creates and returns a new MessageProcessSetForeground with the specified process ID.
// It initializes the embedded Message with the MessageTypeProcessSetForeground constant.
func NewMessageProcessSetForeground(originatorPID int, pid int) *MessageProcessSetForeground {
	return &MessageProcessSetForeground{
		Message: *interfaces.NewMessage(originatorPID, interfaces.MessageTypeProcessSetForeground),
		pid:     pid,
	}
}

func (m *MessageProcessSetForeground) PID() int {
	return m.pid
}

// MessageProcessList represents a request message to retrieve a list of processes from the system.
// It embeds Message for base message functionality and uses ReplyTo to specify the recipient for the response.
type MessageProcessList struct {
	interfaces.Message
	processes []*interfaces.ProcessDescription
	ackChan   chan bool
}

// NewMessageProcessList creates a new MessageProcessList with specified replyTo as an IRouter.
func NewMessageProcessList(originatorPID int, ackChan chan bool) *MessageProcessList {
	return &MessageProcessList{
		Message: *interfaces.NewMessage(originatorPID, interfaces.MessageTypeProcessList),
		ackChan: ackChan,
	}
}

// SetResult sets the list of processes and sends the message to the specified recipient via ReplyTo.
func (m *MessageProcessList) SetResult(processes []*interfaces.ProcessDescription) {
	m.MakeResponse()
	m.processes = processes
}

// GetResult returns the list of ProcessDescription objects associated with the MessageProcessListResponse instance.
func (m *MessageProcessList) GetResult() []*interfaces.ProcessDescription {
	return m.processes
}

// Ack signals the acknowledgment channel and closes it to indicate process completion or unblock synchronous calls.
func (m *MessageProcessList) Ack() bool {
	m.ackChan <- true
	return true
}

type MessageProcessIsRunning struct {
	interfaces.Message
	verifyPid int
	ackChan   chan bool
	result    bool
}

// NewMessageProcessIsRunning creates a new MessageProcessIsRunning with the given originatorPID, pid, and acknowledgment channel.
func NewMessageProcessIsRunning(originatorPID int, verifyPid int, ackChan chan bool) *MessageProcessIsRunning {
	return &MessageProcessIsRunning{
		Message:   *interfaces.NewMessage(originatorPID, interfaces.MessageTypeProcessIsRunning),
		verifyPid: verifyPid,
		ackChan:   ackChan,
	}
}

// Ack acknowledges that the processing of the message is complete by sending a signal to the ackChan and returns true.
func (m *MessageProcessIsRunning) Ack() bool {
	m.ackChan <- true
	return true
}

// VerifyPID returns the PID to be verified.
func (m *MessageProcessIsRunning) VerifyPID() int {
	return m.verifyPid
}

// SetResult assigns the processing result to the result field and marks the message as a response.
func (m *MessageProcessIsRunning) SetResult(result bool) {
	m.MakeResponse()
	m.result = result
}

// GetResult returns the result of the message process as a boolean value.
func (m *MessageProcessIsRunning) GetResult() bool {
	return m.result
}

// MessageNotifyProcessCreate represents a message notifying the creation of a process in the system.
// It embeds the base Message type and includes details about the created process's PID and name.
type MessageNotifyProcessCreate struct {
	interfaces.Message
	createdPID int
	name       string
}

// NewMessageNotifyProcessCreate creates a new instance of MessageNotifyProcessCreate with the given parameters.
// The `originatorPID` specifies the process ID of the message sender.
// The `createdPID` specifies the process ID of the newly created process.
// The `name` specifies the name associated with the created process.
func NewMessageNotifyProcessCreate(originatorPID int, createdPID int, name string) *MessageNotifyProcessCreate {
	return &MessageNotifyProcessCreate{
		Message:    *interfaces.NewMessage(originatorPID, interfaces.MessageTypeNotifyProcessCreate),
		createdPID: createdPID,
		name:       name,
	}
}

// CreatedPID returns the process ID (PID) created and associated with this MessageNotifyProcessCreate instance.
func (m *MessageNotifyProcessCreate) CreatedPID() int {
	return m.createdPID
}

// Name returns the name associated with the MessageNotifyProcessCreate instance.
func (m *MessageNotifyProcessCreate) Name() string {
	return m.name
}

// MessageNotifyProcessTerminate is a message type used to notify about the termination of a specific process.
// It embeds the Message struct to inherit general message features and includes the terminated process ID.
// TerminatedPID returns the ID of the process that was terminated.
type MessageNotifyProcessTerminate struct {
	interfaces.Message
	terminatedPID int
}

// NewMessageMessageNotifyProcessTerminate creates a new MessageNotifyProcessTerminate with the given PIDs.
func NewMessageMessageNotifyProcessTerminate(originatorPID int, terminatedPID int) *MessageNotifyProcessTerminate {
	return &MessageNotifyProcessTerminate{
		Message:       *interfaces.NewMessage(originatorPID, interfaces.MessageTypeNotifyProcessTerminate),
		terminatedPID: terminatedPID,
	}
}

// TerminatedPID retrieves the process ID (PID) of the terminated process from the MessageNotifyProcessTerminate instance.
func (m *MessageNotifyProcessTerminate) TerminatedPID() int {
	return m.terminatedPID
}

// MessageNotifyProcessForeground represents a notification message to indicate a process is moved to the foreground.
// It embeds the Message type and includes an additional foreground process ID field, foregroundPID.
type MessageNotifyProcessForeground struct {
	interfaces.Message
	foregroundPID int
}

// NewMessageNotifyProcessForeground creates a new MessageNotifyProcessForeground with the given PIDs.
func NewMessageNotifyProcessForeground(originatorPID int, foregroundPID int) *MessageNotifyProcessForeground {
	return &MessageNotifyProcessForeground{
		Message:       *interfaces.NewMessage(originatorPID, interfaces.MessageTypeNotifyProcessForeground),
		foregroundPID: foregroundPID,
	}
}

// ForegroundPID returns the process ID (PID) of the process currently marked as running in the foreground.
func (m *MessageNotifyProcessForeground) ForegroundPID() int {
	return m.foregroundPID
}
