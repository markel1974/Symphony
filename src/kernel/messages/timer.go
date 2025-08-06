package messages

import "github.com/markel1974/c64emu/src/kernel/interfaces"

// MessageTimer represents a timed message with an associated process ID, timer ID, and interval.
type MessageTimer struct {
	interfaces.Message
	pid      int
	tid      int
	interval int
}

// NewMessageTimer creates a new MessageTimer instance with the specified process ID and interval.
func NewMessageTimer(router interfaces.IRouter, pid int, interval int) *MessageTimer {
	return &MessageTimer{
		Message:  *interfaces.NewMessage(router, interfaces.MessageTypeTimer),
		pid:      pid,
		interval: interval,
	}
}

// PID returns the process ID associated with the MessageTimer instance.
func (m *MessageTimer) PID() int {
	return m.pid
}

// SetTID assigns a new timer ID (tid) to the MessageTimer instance.
func (m *MessageTimer) SetTID(tid int) {
	m.tid = tid
}

// TID returns the process ID associated with the MessageTimer instance.
func (m *MessageTimer) TID() int {
	return m.tid
}

// Interval returns the interval value associated with the MessageTimer instance.
func (m *MessageTimer) Interval() int {
	return m.interval
}

// MessageTimedMessage represents a timed message with properties for initial delay, interval, and remaining count.
type MessageTimedMessage struct {
	interfaces.IMessage
	msg      interfaces.IMessage
	first    int64
	interval int64
	count    int64
}

// NewMessageTimedMessage creates and returns a new instance of MessageTimedMessage with the specified timing settings.
func NewMessageTimedMessage(router interfaces.IRouter, msg interfaces.IMessage, first int64, interval int64, count int64) *MessageTimedMessage {
	return &MessageTimedMessage{
		IMessage: interfaces.NewMessage(router, interfaces.MessageTypeTimedMessage),
		msg:      msg,
		first:    first,
		interval: interval,
		count:    count,
	}
}

// Message returns the encapsulated IMessage instance within the MessageTimedMessage structure.
func (m *MessageTimedMessage) Message() interfaces.IMessage {
	return m.msg
}

// First returns the first timestamp associated with the MessageTimedMessage instance.
func (m *MessageTimedMessage) First() int64 {
	return m.first
}

// Interval returns the interval value (in int64) of the MessageTimedMessage.
func (m *MessageTimedMessage) Interval() int64 {
	return m.interval
}

// Count retrieves the current value of the count field from the MessageTimedMessage instance.
func (m *MessageTimedMessage) Count() int64 {
	return m.count
}

type MessageTimerCreate struct {
	interfaces.IMessage
	first    int
	interval int
	count    int
}

// NewMessageTimerCreate initializes a new MessageTimerCreate instance with the given router, first, interval, and count values.
// It sets the MessageType to MessageTypeTimerCreate.
func NewMessageTimerCreate(router interfaces.IRouter, first int, interval int, count int) *MessageTimerCreate {
	return &MessageTimerCreate{
		IMessage: interfaces.NewMessage(router, interfaces.MessageTypeTimerCreate),
		first:    first,
		interval: interval,
		count:    count,
	}
}

// First returns the initial delay value in milliseconds for the MessageTimerCreate instance.
func (m *MessageTimerCreate) First() int {
	return m.first
}

// Interval retrieves the interval value associated with the MessageTimerCreate instance.
func (m *MessageTimerCreate) Interval() int {
	return m.interval
}

// Count returns the count value associated with the MessageTimerCreate instance.
func (m *MessageTimerCreate) Count() int {
	return m.count
}

// MessageTimerStop represents a message to stop a specific timer identified by its TID in the system.
type MessageTimerStop struct {
	interfaces.IMessage
	tid int
}

// NewMessageTimerStop creates a new MessageTimerStop instance with the specified router and timer ID (tid).
func NewMessageTimerStop(router interfaces.IRouter, tid int) *MessageTimerStop {
	return &MessageTimerStop{
		IMessage: interfaces.NewMessage(router, interfaces.MessageTypeTimerStop),
		tid:      tid,
	}
}

// TID returns the timer ID associated with the MessageTimerStop instance.
func (m *MessageTimerStop) TID() int {
	return m.tid
}

// MessageTimerCreated represents a message indicating that a timer has been successfully created within the system.
type MessageTimerCreated struct {
	interfaces.IMessage
	tid int
}

// NewMessageTimerCreated creates a new instance of MessageTimerCreated with the specified router and timer ID.
func NewMessageTimerCreated(router interfaces.IRouter, tid int) *MessageTimerCreated {
	return &MessageTimerCreated{
		IMessage: interfaces.NewMessage(router, interfaces.MessageTypeTimerCreated),
		tid:      tid,
	}
}

// TID retrieves the timer identifier associated with the MessageTimerCreated instance.
func (m *MessageTimerCreated) TID() int {
	return m.tid
}
