package interfaces

type IServer interface {
	IRouter

	Start()

	SetRouter(router IRouter)

	Register() []MessageType

	NotifyProcessCreation(desc *ProcessDescription)

	NotifyProcessTermination(desc *ProcessDescription)

	NotifyProcessForeground(desc *ProcessDescription)
}
