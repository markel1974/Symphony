package interfaces

type IRouter interface {
	IReceiver

	PostTimedMessage(msg IMessage, first int64, interval int64, count int64)
}
