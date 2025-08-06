package interfaces

// OnError defines a function type that handles errors that occur during task execution.
type OnError func(task IProcess, err error)

// OnRun defines a function type that performs a task execution with given arguments and returns an error if any occurs.
type OnRun func(task IProcess, args []string) error

// OnTimer defines a function type for tasks invoked at regular intervals, receiving the task, timer id, and interval.
type OnTimer func(task IProcess, tid int, interval int)

// OnRead defines a function type invoked for processing input events with a task, an event code, and a key character.
type OnRead func(task IProcess, code int, key rune)

// OnPaint defines a function type used to handle painting tasks on a specified surface in the context of a task.
type OnPaint func(task IProcess, surface ISurface)

// OnActivate defines a function type executed when a process is activated, receiving the task as a parameter.
type OnActivate func(task IProcess)

// IProcess defines an interface for process management, task handling, interaction, and rendering within a system.
type IProcess interface {
	IRouter

	PID() int

	SetId(i int)

	Parent() IProcess

	Protected() bool

	Line() string

	Description() *ProcessDescription

	GetCommand() ICommand

	SetContext(ctx interface{})

	GetContext() interface{}

	CreateTimer(first int, interval int, count int)

	StopTimer(tid int)

	Timers() []int

	TimersIterator(callback func(tid int) bool)

	AddTimer(tid int)

	IsActive(pid int) bool

	Kill(pid int)

	KillForeground()

	KillAll(name string)

	CWDName() string

	CWDSet(arg string) bool

	CWDGet() string

	CWDPath() []string

	CWDDirectoryListing() []string

	GetScreenSize() (int, int)

	Paint(surface ISurface)

	PaintRequest()

	Setup()

	ProcessExec(line string)

	WindowsSelectionBegin()

	WindowsSelectionEnd()

	WindowsSelectionOptions(option rune, value float64)

	WindowsSelectionNext()

	WindowsSelectionPrevious()

	ProcessList() []*ProcessDescription

	ProcessSetForeground(pid int)

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
