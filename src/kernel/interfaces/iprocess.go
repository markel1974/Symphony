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

	Setup()

	GetCommand() ICommand

	Line() string

	Description() *ProcessDescription

	SetContext(ctx interface{})

	GetContext() interface{}

	CreateTimer(first int, interval int, count int)

	IsActive(pid int) bool

	Kill(pid int)

	KillForeground()

	KillAll(name string)

	CWDSet(arg string) bool

	CWDName() string

	CWDPath() string

	//CWDPathEntries() []string

	CWDDirectoryListing() []string

	GetScreenSize() (int, int)

	PaintRequest()

	ProcessExec(line string)

	WindowsSelectionBegin()

	WindowsSelectionEnd()

	WindowsSelectionOptions(option rune, value float64)

	WindowsSelectionNext()

	WindowsSelectionPrevious()

	ProcessList() []*ProcessDescription

	ProcessSetForeground(pid int)

	Write(data string, eol bool)

	WritePromptEOL(prompt string, eol bool)

	WritePromptLine(prompt string, line string)

	WriteColor(data string, fg ColorDef, bg ColorDef, mode ColorMode, eol bool)

	WriteForeground(data string, color ColorDef, eol bool)

	MoveCursorLeft()

	MoveCursorRight()

	SaveCursor()

	RestoreCursor()

	ClearScreen()

	SetExit()

	Suggestion(in string, cursor int) (string, []string, bool)

	Help(arg string) (string, error)
}
