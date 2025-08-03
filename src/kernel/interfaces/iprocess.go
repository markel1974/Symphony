package interfaces

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

	SetWindowOption(option rune, value float64)

	SetWindowOptions(options *WindowOptions)

	GetCommand() ICommand

	SetContext(ctx interface{})

	GetContext() interface{}

	CreateTimer(first int, interval int, count int)

	StopTimer(tid int)

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

	PaintRequest()

	Start()

	SetCaption(caption string) bool

	ProcessExec(line string) (bool, error)

	WindowsSelectionBegin()

	WindowsSelectionEnd()

	WindowsSelectionOptions(option rune, value float64)

	WindowsSelectionNext()

	WindowsSelectionPrevious()

	ProcessList() []*ProcessDescription

	ProcessSetForeground(pid int) bool

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

	PostMessage(msg IMessage)
}
