package interfaces

type IServer interface {
	Start()

	SetRouter(router IRouter)

	PostMessage(msg IMessage)

	Register() []MessageType

	NotifyProcessCreation(desc *ProcessDescription)

	NotifyProcessTermination(desc *ProcessDescription)

	NotifyProcessForeground(desc *ProcessDescription)
}
