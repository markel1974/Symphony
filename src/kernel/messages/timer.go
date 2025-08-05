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
func NewMessageTimer(pid int, interval int) *MessageTimer {
	return &MessageTimer{
		Message:  *interfaces.NewMessage(interfaces.MessageTypeTimer),
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
func NewMessageTimedMessage(msg interfaces.IMessage, first int64, interval int64, count int64) *MessageTimedMessage {
	return &MessageTimedMessage{
		IMessage: interfaces.NewMessage(interfaces.MessageTypeTimedMessage),
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
