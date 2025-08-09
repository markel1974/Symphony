package messages

import "github.com/markel1974/c64emu/src/kernel/interfaces"

type MessageExitRequested struct {
	interfaces.IMessage
}

func NewMessageExitRequested() *MessageExitRequested {
	return &MessageExitRequested{
		IMessage: interfaces.NewMessageNoAck(interfaces.MessageTypeExitRequested),
	}
}
