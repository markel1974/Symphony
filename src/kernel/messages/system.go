package messages

import "github.com/markel1974/c64emu/src/kernel/interfaces"

type MessageExitRequested struct {
	interfaces.Message
}

func NewMessageExitRequested(originatorPID int) *MessageExitRequested {
	return &MessageExitRequested{
		Message: *interfaces.NewMessage(originatorPID, interfaces.MessageTypeExitRequested),
	}
}
