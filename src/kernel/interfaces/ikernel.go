package interfaces

// IKernel defines the core interface for managing tasks, input, system commands, and rendering operations in a system.
type IKernel interface {
	IRouter

	SetScreenSize(w int, h int)

	AddServer(server IServer)

	CallProcessList(process IRouter) []*ProcessDescription

	CallWindowsSelectionBegin(process IRouter)

	CallWindowsSelectionOptions(process IRouter, option rune, value float64)

	CallWindowsSelectionPrevious(process IRouter)

	CallWindowsSelectionNext(process IRouter)

	CallWindowsSelectionEnd(process IRouter)

	CallWritePromptEOL(process IRouter, prompt string, eol bool)

	CallWritePromptLine(process IRouter, prompt string, line string)

	CallWrite(process IRouter, data string)

	CallWriteLn(process IRouter, data string)

	CallWriteColor(process IRouter, data string, fg ColorDef, bg ColorDef, mode ColorMode)

	CallWriteColorLn(process IRouter, data string, fg ColorDef, bg ColorDef, mode ColorMode)

	CallWriteNormal(process IRouter, data string)

	CallWriteHighlights(process IRouter, data string)

	CallWriteCritical(process IRouter, data string)

	CallMoveCursorLeft(process IRouter)

	CallMoveCursorRight(process IRouter)

	CallSaveCursor(process IRouter)

	CallRestoreCursor(process IRouter)

	CallClearScreen(process IRouter)

	CallScreenSize(process IRouter) (int, int)

	CallFileSystemSuggestion(process IRouter, in string, cursor int) (string, []string, bool)

	CallCWDSet(process IRouter, arg string) bool

	CallCWDGet(process IRouter) string

	CallCWDPath(process IRouter) []string

	CallCWDName(process IRouter) string

	CallCWDDirectoryListing(process IRouter) []string

	CallFileSystemHelp(process IRouter, arg string) (string, error)

	CallExitRequested(process IRouter)

	//CallTimerCreate(process IRouter, pid int, first int, interval int, count int)

	//CallTimerStop(process IRouter, pid int, tid int)

	CallProcessIsActive(process IRouter, pid int) bool
}
