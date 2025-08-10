package messages

import "github.com/markel1974/c64emu/src/kernel/interfaces"

type MessageExitRequested struct {
	interfaces.IMessage
}

func NewMessageExitRequested(source int, destination int) *MessageExitRequested {
	return &MessageExitRequested{
		IMessage: interfaces.NewMessageNoAck(source, destination, interfaces.MessageTypeExitRequested),
	}
}
