package interfaces

type IRouter interface {
	PostMessage(m IMessage)

	PostTimedMessage(msg IMessage, first int64, interval int64, count int64)
}
