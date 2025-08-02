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

// ProcessState represents the state of a task, defined as an integer-based enumerated type.
type ProcessState int

// TaskStateSetup represents the initial setup state of a task.
// TaskStateRunning represents the state where a task is currently running.
const (
	ProcessStateSetup   ProcessState = iota
	ProcessStateRunning ProcessState = iota
)

// ProcessDescription represents the details of a system process, including its name and process ID.
type ProcessDescription struct {
	Name string
	Pid  int
	Line string
}

// WindowOptions represents configurable parameters for a task, including offsets, scaling, and associated command line.
type WindowOptions struct {
	OffsetY int
	OffsetX int
	Scale   float64
	Line    string
}

func NewWindowOptions(offsetX int, offsetY int, scale float64, line string) *WindowOptions {
	return &WindowOptions{
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

	SetId(i int)

	Line() string

	Description() *ProcessDescription

	Options() *WindowOptions

	SetWindowOption(option rune, value float64)

	SetWindowOptions(options *WindowOptions)

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

	DeactivateForeground() bool

	DeactivateAll(name string) int

	CWDName() string

	CWDSet(arg string) bool

	CWDGet() string

	CWDPath() []string

	CWDDirectoryListing() []string

	GetScreenSize() (int, int)

	Paint(surface ISurface)

	PaintRequest() bool

	SetState(state ProcessState)

	SetCaption(caption string) bool

	ProcessExec(line string) (bool, error)

	WindowsSelectionBegin()

	WindowsSelectionEnd()

	WindowsSelectionOptions(option rune, value float64) bool

	WindowsSelectionNext()

	WindowsSelectionPrevious()

	ProcessList() []*ProcessDescription

	ProcessSetFg(pid int) bool

	Write(data string)

	WritePromptEOL(prompt string, eol bool)

	WritePromptLine(prompt string, line string)

	WriteLn(data string)

	WriteColor(data string, fg ColorDef, bg ColorDef, mode ColorMode)

	WriteColorLn(data string, fg ColorDef, bg ColorDef, mode ColorMode)

	WriteNormal(data string)

	WriteHighlights(data string)

	WriteCritical(data string)

	MoveCursorLeft()

	MoveCursorRight()

	SaveCursor()

	RestoreCursor()

	ClearScreen()

	SetExit()

	Suggestion(in string, cursor int) (string, []string, bool)

	Help(arg string) (string, error)
}
