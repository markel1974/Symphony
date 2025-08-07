package process

import (
	"log"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/messages"
)

// Process represents a task or job in the system, including its context, state, associated options, and execution details.
type Process struct {
	kernel      interfaces.IKernel
	cmd         interfaces.ICommand
	user        string
	context     interface{}
	timers      []int
	pid         int
	state       interfaces.ProcessState
	line        string
	messageChan chan interfaces.IMessage
}

// NewProcess initializes and returns a new Process instance with the provided kernel, command, and command line data.
func NewProcess(kernel interfaces.IKernel, pid int, user string, cmd interfaces.ICommand, line string) *Process {
	t := &Process{
		kernel:      kernel,
		pid:         pid,
		user:        user,
		cmd:         cmd,
		context:     nil,
		state:       interfaces.ProcessStateSetup,
		line:        line,
		messageChan: make(chan interfaces.IMessage, 128),
	}
	return t
}

// Process returns the Process instance, implementing the IProcess interface.
func (t *Process) Process() interfaces.IProcess {
	return t
}

// Setup begins the process by setting its state to running and initiating its event loop asynchronously.
func (t *Process) Setup() {
	c := make(chan bool)
	t.eventLoop(c)
	_ = <-c
	t.state = interfaces.ProcessStateRunning
}

// Description provides a brief summary of the process including its name, PID, and line information.
func (t *Process) Description() *interfaces.ProcessDescription {
	return interfaces.NewProcessDescription(t.cmd.Name(), t.pid, t.line)
}

// Line returns the line configuration of the Process as a string.
func (t *Process) Line() string {
	return t.line
}

// PID returns the process ID (PID) associated with the task.
func (t *Process) PID() int {
	return t.pid
}

func (t *Process) User() string {
	return t.user
}

// GetCommand returns the ICommand instance associated with the Process.
func (t *Process) GetCommand() interfaces.ICommand {
	return t.cmd
}

// SetContext sets the context for the task, storing the provided context object in the task's internal context field.
func (t *Process) SetContext(ctx interface{}) {
	t.context = ctx
}

// GetContext retrieves the context associated with the Process, returning it as an interface{}.
func (t *Process) GetContext() interface{} {
	return t.context
}

// CreateTimer initializes a timer with a specified start delay, repeat interval, and count for the current task.
func (t *Process) CreateTimer(first int, interval int, count int) {
	t.kernel.PostMessage(messages.NewMessageTimerCreate(t, first, interval, count))
}

// StopTimer stops a timer identified by the given timer ID (tid) for the current task and returns true if successful.
func (t *Process) StopTimer(tid int) {
	t.kernel.PostMessage(messages.NewMessageTimerStop(t, tid))
}

// IsActive checks if the process with the specified PID is currently active in the kernel.
func (t *Process) IsActive(pid int) bool {
	return t.kernel.CallProcessIsActive(t, pid)
}

// Kill attempts to terminate the task associated with the specified pid and returns true if successful.
func (t *Process) Kill(pid int) {
	t.kernel.PostMessage(messages.NewMessageProcessKill(t, pid))
}

// KillForeground removes the process from the foreground state and returns true if the operation succeeds.
func (t *Process) KillForeground() {
	t.kernel.PostMessage(messages.NewMessageProcessKillForeground(t))
}

// KillAll terminates all tasks matching the provided name and returns the count of deactivated tasks.
func (t *Process) KillAll(name string) {
	t.kernel.PostMessage(messages.NewMessageProcessKillAll(t, name))
}

// ProcessSetForeground sets the foreground task by specifying its PID and returns true if successfully set.
func (t *Process) ProcessSetForeground(pid int) {
	t.kernel.PostMessage(messages.NewMessageProcessSetForeground(t, pid))
}

// ProcessList returns a string representation of the task list from the kernel.
func (t *Process) ProcessList() []*interfaces.ProcessDescription {
	return t.kernel.CallProcessList(t)
}

// PaintRequest sends a request to repaint the task and returns true if the request was successfully processed.
func (t *Process) PaintRequest() {
	t.kernel.PostMessage(messages.NewMessagePaintRequest(t))
}

// GetScreenSize returns the width and height of the screen as integers.
func (t *Process) GetScreenSize() (int, int) {
	return t.kernel.CallScreenSize(t)
}

// CWDSet sets the current working directory to the specified path and returns true if the operation is successful.
func (t *Process) CWDSet(arg string) bool {
	return t.kernel.CallCWDSet(t, arg)
}

// CWDPathEntries retrieves the current working directory path as a slice of strings from the kernel.
//func (t *Process) CWDPathEntries() []string {
//	return t.kernel.CallCWDPathEntries(t)
//}

// CWDName returns the current working directory name by invoking a kernel-level method.
func (t *Process) CWDName() string {
	return t.kernel.CallCWDName(t)
}

// CWDPath retrieves the current working directory as a string from the associated kernel instance.
func (t *Process) CWDPath() string {
	return t.kernel.CallCWDPath(t)
}

// CWDDirectoryListing retrieves a slice of strings representing the child nodes of the current working directory (CWD).
func (t *Process) CWDDirectoryListing() []string {
	return t.kernel.CallCWDDirectoryListing(t)
}

// Suggestion provides auto-completion suggestions based on the input string and cursor position. Returns prefix, suggestions, and a success flag.
func (t *Process) Suggestion(in string, cursor int) (string, []string, bool) {
	return t.kernel.CallFileSystemSuggestion(t, in, cursor)
}

// Help calls the kernel's Help method with the provided argument and returns the result or an error.
func (t *Process) Help(arg string) (string, error) {
	return t.kernel.CallFileSystemHelp(t, arg)
}

// ProcessExec executes a task based on the provided command line input and returns a success status and any execution error.
func (t *Process) ProcessExec(line string) {
	t.kernel.PostMessage(messages.NewMessageProcessExec(t, line))
}

// WindowsSelectionBegin updates the task selection for the given process ID by invoking the kernel's task selection method.
func (t *Process) WindowsSelectionBegin() {
	t.kernel.PostMessage(messages.NewMessageWindowsSelectionBegin(t))
}

// WindowsSelectionEnd finalizes the current text selection process within the Windows environment for the associated process.
func (t *Process) WindowsSelectionEnd() {
	t.kernel.PostMessage(messages.NewMessageWindowsSelectionEnd(t))
}

// WindowsSelectionPrevious moves the task selection pointer to the previous task in the list within the Process.
func (t *Process) WindowsSelectionPrevious() {
	t.kernel.PostMessage(messages.NewMessageWindowsSelectionPrevious(t))
}

// WindowsSelectionNext moves the task selection to the next task in the sequence by invoking the kernel method.
func (t *Process) WindowsSelectionNext() {
	t.kernel.PostMessage(messages.NewMessageWindowsSelectionNext(t))
}

// WindowsSelectionOptions configures selection behavior for the task based on the provided option and value.
func (t *Process) WindowsSelectionOptions(option rune, value float64) {
	t.kernel.PostMessage(messages.NewMessageWindowsSelectionOptions(t, option, value))
}

// WritePromptEOL writes the provided prompt followed by an end-of-line character if the eol parameter is true.
func (t *Process) WritePromptEOL(prompt string, eol bool) {
	t.kernel.CallWriteColor(t, "", interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal, eol)
	t.kernel.CallWriteColor(t, prompt, interfaces.ColorGreenDef, interfaces.ColorNoneDef, interfaces.ModeNormal, false)
}

// WritePromptLine sends a specified prompt and line string to the kernel for handling the output display.
func (t *Process) WritePromptLine(prompt string, line string) {
	t.kernel.CallClearLine(t, line)
	t.kernel.CallWriteColor(t, prompt, interfaces.ColorGreenDef, interfaces.ColorNoneDef, interfaces.ModeNormal, false)
	t.kernel.CallWriteColor(t, line, interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal, false)
}

// Write sends the provided string data to the kernel's write mechanism associated with the task.
func (t *Process) Write(data string, eol bool) {
	t.kernel.CallWrite(t, data, eol)
}

// WriteColor writes a string to the output with specified foreground, background colors, and a color mode.
func (t *Process) WriteColor(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode, eol bool) {
	t.kernel.CallWriteColor(t, data, fg, bg, mode, eol)
}

// WriteForeground writes the given data to the foreground with the specified color using the kernel's functionality.
func (t *Process) WriteForeground(data string, color interfaces.ColorDef, eol bool) {
	t.kernel.CallWriteColor(t, data, color, interfaces.ColorNoneDef, interfaces.ModeNormal, eol)
}

// MoveCursorLeft moves the cursor one position to the left within the render context.
func (t *Process) MoveCursorLeft() {
	t.kernel.CallMoveCursorLeft(t)
}

// MoveCursorRight moves the cursor one position to the right by invoking the render's MoveCursorRight method.
func (t *Process) MoveCursorRight() {
	t.kernel.CallMoveCursorRight(t)
}

// SaveCursor saves the current cursor state by invoking the SaveCursor method on the associated renderer.
func (t *Process) SaveCursor() {
	t.kernel.CallSaveCursor(t)
}

// RestoreCursor restores the cursor to its previous position using the render instance of the Kernel.
func (t *Process) RestoreCursor() {
	t.kernel.CallRestoreCursor(t)
}

// ClearScreen clears the task's screen by delegating the request to the associated kernel.
func (t *Process) ClearScreen() {
	t.kernel.CallClearScreen(t)
}

func (t *Process) RequestProcessList() []*interfaces.ProcessDescription {
	ackChan := make(chan bool, 1)
	request := messages.NewMessageProcessList(t, ackChan)
	t.kernel.PostMessage(request)
	<-ackChan
	return nil
}

// SetExit signals the kernel that an exit is requested for the task.
func (t *Process) SetExit() {
	t.kernel.CallExitRequested(t)
}

// PostMessage sends the provided message to the message channel for processing.
func (t *Process) PostMessage(msg interfaces.IMessage) {
	t.messageChan <- msg
}

// evenLoop continuously listens on the message channel and processes incoming messages until a quit message is received.
func (t *Process) eventLoop(r chan bool) {
	go func() {
		r <- true
		for {
			select {
			case m, ok := <-t.messageChan:
				if !ok {
					return
				}
				m.Ack()
				if m.GetType() == interfaces.MessageTypeQuit {
					close(t.messageChan)
					return
				}
				t.handleMessage(m)
			}
		}
	}()
}

// handleMessage processes incoming messages based on their type and executes corresponding logic or logs an error.
func (t *Process) handleMessage(msg interfaces.IMessage) {
	switch msg.GetType() {
	case interfaces.MessageTypeError:
		mt, ok := msg.(*messages.MessageError)
		if !ok {
			return
		}
		if onError := t.cmd.OnError(); onError != nil {
			onError(t, mt.Error())
		}
	case interfaces.MessageTypeTimer:
		mt, ok := msg.(*messages.MessageTimer)
		if !ok {
			return
		}
		if timerEvent := t.cmd.OnTimer(); timerEvent != nil {
			timerEvent(t, mt.TID(), mt.Interval())
		}
	case interfaces.MessageTypeRead:
		mt, ok := msg.(*messages.MessageRead)
		if !ok {
			return
		}
		if mt.Broadcast() {
			if readBroadcastEvent := t.cmd.OnReadBroadcast(); readBroadcastEvent != nil {
				readBroadcastEvent(t, int(mt.Kind()), mt.Data())
			}
		} else {
			if readEvent := t.cmd.OnRead(); readEvent != nil {
				readEvent(t, int(mt.Kind()), mt.Data())
			}
		}
	case interfaces.MessageTypeProcessStart:
		mt, ok := msg.(*messages.MessageProcessStart)
		if !ok {
			return
		}
		if !t.cmd.Background() {
			t.kernel.PostMessage(messages.NewMessageProcessSetForeground(t, t.PID()))
		}
		_ = t.cmd.Execute(t, mt.Args())
		if !t.cmd.Daemon() {
			t.kernel.PostMessage(messages.NewMessageMessageProcessExit(t))
			return
		}
	case interfaces.MessageTypeProcessActivate:
		_, ok := msg.(*messages.MessageProcessActivate)
		if !ok {
			return
		}
		if activate := t.cmd.OnActivate(); activate != nil {
			activate(t)
		}
	case interfaces.MessageTypeTimerCreated:
		mt, ok := msg.(*messages.MessageTimerCreated)
		if !ok {
			return
		}
		t.timers = append(t.timers, mt.TID())
	case interfaces.MessageTypePaintPrepare:
		mt, ok := msg.(*messages.MessagePaintPrepare)
		if !ok {
			return
		}
		if paintEvent := t.cmd.OnPaint(); paintEvent != nil {
			mt.Surface().Begin()
			paintEvent(t, mt.Surface())
			//mt.Surface().Paint(surface)
			mt.Surface().End()

			ma := messages.NewMessagePaintApply(t, mt.Surface())
			t.kernel.PostMessage(ma)
		}
	default:
		log.Printf("unknown message type: %d", msg.GetType())
	}
}
