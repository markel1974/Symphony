package interfaces

// IKernel defines the core interface for managing tasks, input, system commands, and rendering operations in a system.
type IKernel interface {
	IRouter

	SetScreenSize(w int, h int)

	AddServer(server IServer)

	CallExitRequested(process IRouter)

	CallProcessIsActive(process IRouter, pid int) bool

	CallProcessList(process IRouter) []*ProcessDescription

	CallWrite(process IRouter, data string, eol bool)

	CallWriteColor(process IRouter, data string, fg ColorDef, bg ColorDef, mode ColorMode, eol bool)

	CallMoveCursorLeft(process IRouter)

	CallMoveCursorRight(process IRouter)

	CallSaveCursor(process IRouter)

	CallRestoreCursor(process IRouter)

	CallClearLine(router IRouter, line string)

	CallClearScreen(process IRouter)

	CallScreenSize(process IRouter) (int, int)

	CallCWDSet(process IRouter, arg string) bool

	CallCWDPath(process IRouter) string

	CallCWDName(process IRouter) string

	//CallCWDPathEntries(process IRouter) []string

	CallCWDDirectoryListing(process IRouter) []string

	CallFileSystemHelp(process IRouter, arg string) (string, error)

	CallFileSystemSuggestion(process IRouter, in string, cursor int) (string, []string, bool)
}
