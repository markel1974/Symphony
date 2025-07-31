package messages

// MessageTimer represents a timed message with an associated process ID, timer ID, and interval.
type MessageTimer struct {
	Message
	pid      int
	tid      int
	interval int
}

// NewMessageTimer creates a new MessageTimer instance with the specified process ID and interval.
func NewMessageTimer(pid int, interval int) *MessageTimer {
	return &MessageTimer{
		Message:  Message{MessageTypeTimer},
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
