package messages

import "github.com/markel1974/c64emu/src/kernel/interfaces"

type MessageExitRequested struct {
	interfaces.Message
}

func NewMessageExitRequested(router interfaces.IRouter) *MessageExitRequested {
	return &MessageExitRequested{
		Message: *interfaces.NewMessage(router, interfaces.MessageTypeExitRequested),
	}
}
