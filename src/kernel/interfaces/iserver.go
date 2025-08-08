package interfaces

type IServer interface {
	IUserRouter

	Register() []MessageType

	Setup(router IKernelServerRouter, process IUserProcess) error

	Name() string
}
