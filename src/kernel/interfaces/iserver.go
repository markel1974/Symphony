package interfaces

type IServer interface {
	IRouter

	SetProcess(process IProcess)

	Name() string

	Start()

	Register(router IRouter) []MessageType

	NotifyProcessCreation(pid int, name string)

	NotifyProcessTermination(pid int)

	NotifyProcessForeground(pid int)
}
