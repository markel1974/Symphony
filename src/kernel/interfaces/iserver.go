package interfaces

type IServer interface {
	IRouter

	SetProcess(process IProcess)

	Name() string

	Start()

	Register(router IRouter) []MessageType
}
