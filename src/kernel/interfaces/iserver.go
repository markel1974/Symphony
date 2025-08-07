package interfaces

type IServer interface {
	IRouter

	Start()

	Register(router IRouter) []MessageType

	NotifyProcessCreation(pid int, name string)

	NotifyProcessTermination(pid int)

	NotifyProcessForeground(pid int)
}
