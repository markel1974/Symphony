package messages

import "github.com/markel1974/c64emu/src/kernel/interfaces"

// MessageRead represents a specific type of Message containing read operation data. It embeds Message and includes a data field.
type MessageRead struct {
	Message
	kind interfaces.KeyType
	data rune
}

// NewMessageRead creates a new MessageRead instance with provided data and limits its length to n if necessary.
func NewMessageRead(kind interfaces.KeyType, data rune) *MessageRead {
	return &MessageRead{
		Message: Message{MessageTypeRead},
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
