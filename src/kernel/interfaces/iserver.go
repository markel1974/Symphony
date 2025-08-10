package interfaces

type IServer interface {
	IRouter

	Register() []MessageType

	Setup(router IKernelResponseRouter, pid int, process IUserProcess) error

	Name() string

	PID() int
}
