package messages

import "github.com/markel1974/symphony/src/kernel/interfaces"

// MessageTimer represents a timer message in the system, encapsulating process Id, timer Id, and interval details.
type MessageTimer struct {
	interfaces.IMessage
	tid      int
	interval int
}

// NewMessageTimer creates and returns a new MessageTimer with the specified process Id and interval duration.
func NewMessageTimer(source int, destination int, interval int) *MessageTimer {
	return &MessageTimer{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeTimer),
		interval: interval,
	}
}

// SetTID sets the timer Id (tid) for the MessageTimer instance.
func (m *MessageTimer) SetTID(tid int) {
	m.tid = tid
}

// TID returns the timer Id associated with the MessageTimer instance.
func (m *MessageTimer) TID() int {
	return m.tid
}

// Interval returns the interval value for the MessageTimer instance.
func (m *MessageTimer) Interval() int {
	return m.interval
}

// MessageTimedMessage wraps an IMessage with timing properties for scheduled or repeated messaging functionality.
type MessageTimedMessage struct {
	interfaces.IMessage
	msg      interfaces.IMessage
	first    int64
	interval int64
	count    int64
}

// MessageTimerCreate represents a request to create a timer with specific parameters for delay, interval, and repetition.
type MessageTimerCreate struct {
	interfaces.IMessage
	first    int
	interval int
	count    int
}

// NewMessageTimerCreate creates a MessageTimerCreate object with specified first delay, interval, and execution count.
func NewMessageTimerCreate(source int, destination int, first int, interval int, count int) *MessageTimerCreate {
	return &MessageTimerCreate{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeTimerCreate),
		first:    first,
		interval: interval,
		count:    count,
	}
}

// First returns the initial delay in milliseconds before the timer starts.
func (m *MessageTimerCreate) First() int {
	return m.first
}

// Interval returns the interval value of the MessageTimerCreate instance.
func (m *MessageTimerCreate) Interval() int {
	return m.interval
}

// Count returns the total count value associated with the MessageTimerCreate instance.
func (m *MessageTimerCreate) Count() int {
	return m.count
}

// MessageTimerStop represents a message to stop a previously created timer, identified by its timer Id (tid).
type MessageTimerStop struct {
	interfaces.IMessage
	tid int
}

// NewMessageTimerStop creates a new MessageTimerStop instance with the specified timer Id (tid).
// It initializes the IMessage field with a no-acknowledgment message of type MessageTypeTimerStop.
// tid is the identifier of the timer to be stopped.
func NewMessageTimerStop(source int, destination int, tid int) *MessageTimerStop {
	return &MessageTimerStop{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeTimerStop),
		tid:      tid,
	}
}

// TID returns the timer Id associated with the MessageTimerStop instance.
func (m *MessageTimerStop) TID() int {
	return m.tid
}

// MessageTimerCreated is a message type signaling the successful creation of a timer, encapsulating the timer's Id (tid).
type MessageTimerCreated struct {
	interfaces.IMessage
	tid int
}

// NewMessageTimerCreated creates a new instance of MessageTimerCreated with the specified timer Id (tid).
func NewMessageTimerCreated(source int, destination int, tid int) *MessageTimerCreated {
	return &MessageTimerCreated{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeTimerCreated),
		tid:      tid,
	}
}

// TID retrieves the timer Id associated with the MessageTimerCreated instance.
func (m *MessageTimerCreated) TID() int {
	return m.tid
}
