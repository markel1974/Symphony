package interfaces

type IServer interface {
	IRouter

	Register(router IRouter) []MessageType

	Setup(process IProcess) error

	Name() string
}
