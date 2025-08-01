package interfaces

// IKernel defines the core interface for managing tasks, input, system commands, and rendering operations in a system.
type IKernel interface {
	CallTaskExec(line string, options *ProcessOptions) (bool, error)

	CallTaskKill(pid int) bool

	CallTaskKillAll(name string) int

	CallTaskGetForegroundName() (int, string)

	CallTaskSetBackground() bool

	CallTaskKillForeground()

	CallTaskSaveAll(name string) bool

	CallTaskRestoreAll(name string) bool

	CallTaskList() string

	CallTaskSavedList() []string

	CallTaskSelection(pid int)

	CallTaskSelectionPrevious()

	CallTaskSelectionNext()

	CallTaskSelectionOptions(option rune, value float64) bool

	CallPaintRequest() bool

	CallWrite(data string)

	CallWriteLn(data string)

	CallWriteColor(data string, fg ColorDef, bg ColorDef, mode ColorMode)

	CallWriteColorLn(data string, fg ColorDef, bg ColorDef, mode ColorMode)

	CallClearScreen()

	CallScreenSize() (int, int)

	CallCWDSet(arg string) bool

	CallCWDGet() string

	CallCWDPath() []string

	CallCWDDirectoryListing() []string

	CallHistory(verb HistoryAction, idx int)

	CallHelp(arg string) (string, error)

	CallSetScreenSize(w int, h int)

	CallExitRequested()

	CallSetFg(pid int) bool

	CallCreateTimer(pid int, first int, interval int, count int) bool

	CallStopTimer(pid int, tid int) bool

	CallIsActive(pid int) bool
}
