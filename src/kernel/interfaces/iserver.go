package interfaces

type IServer interface {
	IRouter

	Start()

	Register(router IRouter) []MessageType

	NotifyProcessCreation(desc *ProcessDescription)

	NotifyProcessTermination(desc *ProcessDescription)

	NotifyProcessForeground(desc *ProcessDescription)
}
