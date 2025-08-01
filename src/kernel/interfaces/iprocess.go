package interfaces

import (
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

// State represents the state of a task, defined as an integer-based enumerated type.
type ProcessState int

// TaskStateSetup represents the initial setup state of a task.
// TaskStateRunning represents the state where a task is currently running.
const (
	ProcessStateSetup   ProcessState = iota
	ProcessStateRunning ProcessState = iota
)

// ProcessOptions represents configurable parameters for a task, including offsets, scaling, and associated command line.
type ProcessOptions struct {
	OffsetY int
	OffsetX int
	Scale   float64
	Line    string
}

func NewProcessOptions(offsetX int, offsetY int, scale float64, line string) *ProcessOptions {
	return &ProcessOptions{
		OffsetY: offsetY,
		OffsetX: offsetX,
		Scale:   scale,
		Line:    line,
	}
}

// RunFn defines a function type that performs a task execution with given arguments and returns an error if any occurs.
type RunFn func(task IProcess, args []string) error

// TimerFn defines a function type for tasks invoked at regular intervals, receiving the task, timer id, and interval.
type TimerFn func(task IProcess, tid int, interval int)

// ReadFn defines a function type invoked for processing input events with a task, an event code, and a key character.
type ReadFn func(task IProcess, code int, key rune)

// PaintFn defines a function type used to handle painting tasks on a specified surface in the context of a task.
type PaintFn func(task IProcess, surface ISurface)

// IProcess defines an interface for process management, task handling, interaction, and rendering within a system.
type IProcess interface {
	PID() int

	Line() string

	Options() *ProcessOptions

	SetOption(option rune, value float64)

	SetOptions(options *ProcessOptions)

	SetId(i int)

	GetCommand() ICommand

	SetContext(ctx interface{})

	GetContext() interface{}

	CreateTimer(first int, interval int, count int) bool

	StopTimer(tid int) bool

	Timers() []int

	TimersIterator(callback func(tid int) bool)

	AddTimer(tid int)

	IsActive(pid int) bool

	Deactivate(pid int) bool

	DeactivateAll(name string) int

	CWDSet(arg string) bool

	CWDGet() string

	CWDPath() []string

	CWDDirectoryListing() []string

	GetScreenSize() (int, int)

	Paint(surface ISurface)

	PaintRequest() bool

	SetState(state ProcessState)

	SetCaption(caption string) bool

	SetTaskSelection(int)

	SetTaskSelectionOptions(option rune, value float64) bool

	SetTaskSelectionNext()

	SetTaskSelectionPrevious()

	TaskList() string

	SaveTasks(name string) bool

	RestoreTasks(name string) bool

	ListTasks() []string

	SetFg(pid int) bool

	Write(data string)

	WriteLn(data string)

	WriteColor(data string, fg ColorDef, bg ColorDef, mode ColorMode)

	WriteColorLn(data string, fg ColorDef, bg ColorDef, mode ColorMode)

	History(verb HistoryAction, idx int)

	ClearScreen()

	SetExit()

	Help(arg string) (string, error)
}
