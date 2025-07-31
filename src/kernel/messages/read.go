package messages

// MessageRead represents a specific type of Message containing read operation data. It embeds Message and includes a data field.
type MessageRead struct {
	Message
	data []byte
}

// NewMessageRead creates a new MessageRead instance with provided data and limits its length to n if necessary.
func NewMessageRead(data []byte, n int) *MessageRead {
	if n > len(data) {
		n = len(data) - 1
	}
	x := data[:n]
	return &MessageRead{
		Message: Message{MessageTypeRead},
		data:    x,
	}
}

// Data returns the data payload contained within the MessageRead instance.
func (m *MessageRead) Data() []byte {
	return m.data
}
