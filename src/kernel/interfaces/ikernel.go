package interfaces

// IKernel defines the core interface for managing processes, input, system commands, and rendering operations in a system.
type IKernel interface {
	IRouter

	SetScreenSize(w int, h int)

	AddServer(server IServer)

	CallExitRequested(process IRouter)

	CallProcessIsActive(process IRouter, pid int) bool

	CallScreenSize(process IRouter) (int, int)

	CallCWDSet(process IRouter, arg string) bool

	CallCWDPath(process IRouter) string

	CallCWDDirectoryListing(process IRouter) []string

	CallFileSystemHelp(process IRouter, arg string) (string, error)

	//CallFileSystemSuggestion(process IRouter, in string, cursor int) (string, []string, bool)
}
