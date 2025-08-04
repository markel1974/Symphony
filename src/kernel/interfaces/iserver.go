package interfaces

type IServer interface {
	NotifyProcessCreation(desc *ProcessDescription)

	NotifyProcessTermination(desc *ProcessDescription)
}
