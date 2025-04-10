package interfaces

import "io"

// PathSeparator defines the standard character used to separate elements in a file system path.
const PathSeparator = "/"

// RunFn defines a function type executed with a given context, command, process ID, and arguments, returning an error if any.
type RunFn func(task ITask, args []string) error

// TimerFn defines a function type for handling timer-based events within a context, providing arguments for event details.
type TimerFn func(task ITask, tid int, interval int)

// ReadFn defines a function type for handling read events, providing context, command, process ID, and input details.
type ReadFn func(task ITask, code int, key rune)

// PaintFn defines a function type responsible for handling paint events given the context, command, process ID, user context, and surface.
type PaintFn func(task ITask, surface ISurface)

type ITask interface {
	PID() int

	GetCommand() ICommand

	SetContext(ctx interface{})

	GetContext() interface{}

	CreateTimer(first int, interval int, count int) bool

	StopTimer(tid int) bool

	IsActive(pid int) bool

	Deactivate(pid int) bool

	DeactivateAll(name string) int

	CWDSet(arg string) bool

	CWDGet() string

	CWDPath() []string

	CWDChilds() []string

	GetScreenSize() (int, int)

	PaintRequest() bool

	SetCaption(caption string) bool

	SetSelectionMode(int)

	SetSelectionOptions(option rune, value float64) bool

	SetSelectionModeNext()

	SetSelectionModePrevious()

	TaskList() string

	SaveTasks(name string) bool

	RestoreTasks(name string) bool

	ListTasks() []string

	SetFg(pid int) bool

	GetWriter() io.Writer

	Write(data string)

	WriteLn(data string)

	WriteColor(data string, fg ColorDef, bg ColorDef, mode ColorMode)

	WriteColorLn(data string, fg ColorDef, bg ColorDef, mode ColorMode)

	History(verb HistoryAction, idx int)

	ClearScreen()

	SetExit()
}
