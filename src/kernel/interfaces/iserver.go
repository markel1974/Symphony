package interfaces

type IServer interface {
	IReceiver

	Start()

	SetRouter(router IRouter)

	Register() []MessageType

	NotifyProcessCreation(desc *ProcessDescription)

	NotifyProcessTermination(desc *ProcessDescription)

	NotifyProcessForeground(desc *ProcessDescription)
}
