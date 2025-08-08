package interfaces

// IKernel defines the core interface for managing processes, input, system commands, and rendering operations in a system.
type IKernel interface {
	IRouter

	SetScreenSize(w int, h int)

	Setup(server []IServer) error

	CallExitRequested(process IRouter)

	CallProcessIsActive(process IRouter, pid int) bool
}
