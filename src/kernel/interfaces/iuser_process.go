package interfaces

// OnError defines a function type that handles errors that occur during process execution.
type OnError func(process IUserProcess, err error)

// OnRun defines a function type that performs a process execution with given arguments and returns an error if any occurs.
type OnRun func(process IUserProcess, args []string) error

// OnTimer defines a function type for processes invoked at regular intervals, receiving the process, timer id, and interval.
type OnTimer func(process IUserProcess, tid int, interval int)

// OnRead defines a function type invoked for processing input events with a process, an event code, and a key character.
type OnRead func(process IUserProcess, code int, key rune)

// OnPaint defines a function type used to handle painting processes on a specified surface in the context of a process.
type OnPaint func(process IUserProcess, surface ISurface)

// OnActivate defines a function type executed when a process is activated, receiving the process as a parameter.
type OnActivate func(process IUserProcess)

// IUserProcess defines an interface for process management, process handling, interaction, and rendering within a system.
type IUserProcess interface {
	IRouter

	PID() int

	Bind(router IKernelRequestRouter)

	User() string

	Start()

	GetCommand() ICommand

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
