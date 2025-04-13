package interfaces

import (
	"io"
	"strings"
)

// PathSeparator is the string used to separate elements in a hierarchical path format, typically a forward slash "/".
const PathSeparator = "/"

func PathToSegments(path string) []string {
	var segments []string
	for _, part := range strings.Split(path, PathSeparator) {
		if len(part) > 0 {
			segments = append(segments, part)
		}
	}
	return segments
}

func IsPathAbsolute(path string) bool {
	isAbsolute := strings.HasPrefix(path, PathSeparator)
	return isAbsolute
}

// RunFn defines a function type that performs a task execution with given arguments and returns an error if any occurs.
type RunFn func(task ITask, args []string) error

// TimerFn defines a function type for tasks invoked at regular intervals, receiving the task, timer id, and interval.
type TimerFn func(task ITask, tid int, interval int)

// ReadFn defines a function type invoked for processing input events with a task, an event code, and a key character.
type ReadFn func(task ITask, code int, key rune)

// PaintFn defines a function type used to handle painting tasks on a specified surface in the context of a task.
type PaintFn func(task ITask, surface ISurface)

// ITask defines a task interface providing functions for process handling, context management, and terminal interactions.
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

	Help(arg string) (string, error)
}
