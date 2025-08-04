package interfaces

type IReceiver interface {
	PostMessage(msg IMessage)

	PID() int
}
