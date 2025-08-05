package interfaces

type IRouter interface {
	PostMessage(msg IMessage)

	PID() int

	User() string
}
