package interfaces

type IServer interface {
	IRouter

	Register() []MessageType

	Setup(router IKernelResponseRouter, process IUserProcess) error

	Name() string
}
