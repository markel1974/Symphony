package interfaces

// IKernel defines the core interface for managing tasks, input, system commands, and rendering operations in a system.
type IKernel interface {
	IRouter

	AddServer(server IServer)

	CallProcessExec(line string) (bool, error)

	CallProcessKill(pid int) bool

	CallProcessKillForeground() bool

	CallProcessKillAll(name string) int

	CallProcessList() []*ProcessDescription

	CallWindowsSelectionBegin()

	CallWindowsSelectionEnd()

	CallWindowsSelectionPrevious()

	CallWindowsSelectionNext()

	CallWindowsSelectionOptions(option rune, value float64)

	CallWritePromptEOL(prompt string, eol bool)

	CallWritePromptLine(prompt string, line string)

	CallWrite(data string)

	CallWriteLn(data string)

	CallWriteColor(data string, fg ColorDef, bg ColorDef, mode ColorMode)

	CallWriteColorLn(data string, fg ColorDef, bg ColorDef, mode ColorMode)

	CallWriteNormal(data string)

	CallWriteHighlights(data string)

	CallWriteCritical(data string)

	CallMoveCursorLeft()

	CallMoveCursorRight()

	CallSaveCursor()

	CallRestoreCursor()

	CallClearScreen()

	CallScreenSize() (int, int)

	CallFileSystemSuggestion(in string, cursor int) (string, []string, bool)

	CallCWDSet(arg string) bool

	CallCWDGet() string

	CallCWDPath() []string

	CallCWDName() string

	CallCWDDirectoryListing() []string

	CallFileSystemHelp(arg string) (string, error)

	CallSetScreenSize(w int, h int)

	CallExitRequested()

	CallProcessSetForeground(pid int) bool

	CallTimerCreate(pid int, first int, interval int, count int)

	CallTimerStop(pid int, tid int)

	CallProcessIsActive(pid int) bool
}
