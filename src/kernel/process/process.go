package process

import (
	"log"
	"time"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/messages"
)

const (
	gatekeeperQueueLen = 1024
	gatekeeperQueueMax = gatekeeperQueueLen - 1

	executorQueueLen = 128
	executorQueueMax = executorQueueLen - 1
)

// Process represents a process or job in the system, including its context, state, associated options, and execution details.
type Process struct {
	kRouter          interfaces.IKernelRouter
	cmd              interfaces.ICommand
	user             string
	context          interface{}
	timers           []int
	pid              int
	state            interfaces.ProcessState
	gatekeeperChan   chan interfaces.IMessage
	executorChan     chan interfaces.IMessage
	executorWaitChan chan bool
	timeout          time.Duration
}

// NewProcess initializes and returns a new Process instance with the provided kRouter, command, and command line data.
func NewProcess(pid int, user string, cmd interfaces.ICommand) *Process {
	t := &Process{
		pid:              pid,
		user:             user,
		cmd:              cmd,
		context:          nil,
		state:            interfaces.ProcessStateSetup,
		gatekeeperChan:   make(chan interfaces.IMessage, gatekeeperQueueLen),
		executorChan:     make(chan interfaces.IMessage, executorQueueLen),
		executorWaitChan: make(chan bool, 1),
		timeout:          10 * time.Second,
	}
	return t
}

func (t *Process) Bind(kRouter interfaces.IKernelRouter) {
	t.kRouter = kRouter
}

// Process returns the Process instance, implementing the IUserProcess interface.
func (t *Process) Process() interfaces.IUserProcess {
	return t
}

// Start begins the process by setting its state to running and initiating its event loop asynchronously.
func (t *Process) Start() {
	a := make(chan bool)
	b := make(chan bool)
	t.gatekeeperLoop(a)
	_ = <-a
	t.executorLoop(b)
	_ = <-b
	t.state = interfaces.ProcessStateRunning
}

// PID returns the process ID (PID) associated with the process.
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

// SetContext sets the context for the process, storing the provided context object in the process internal context field.
func (t *Process) SetContext(ctx interface{}) {
	t.context = ctx
}

// GetContext retrieves the context associated with the Process, returning it as an interface{}.
func (t *Process) GetContext() interface{} {
	return t.context
}

// CreateTimer initializes a timer with a specified start delay, repeat interval, and count for the current process.
func (t *Process) CreateTimer(first int, interval int, count int) {
	t.kRouter.PostKernelMessage(messages.NewMessageTimerCreate(first, interval, count))
}

// StopTimer stops a timer identified by the given timer ID (tid) for the current process and returns true if successful.
func (t *Process) StopTimer(tid int) {
	t.kRouter.PostKernelMessage(messages.NewMessageTimerStop(tid))
}

// IsActive checks if the process with the specified PID is currently active in the kRouter.
func (t *Process) IsActive(pid int) bool {
	msg := messages.NewMessageProcessIsRunning(pid, t.executorWaitChan)
	t.kRouter.PostKernelMessage(msg)
	select {
	case <-t.executorWaitChan:
		return msg.GetResponse()
	case <-time.After(t.timeout):
		log.Printf("IsActive timeout waiting for response")
		return msg.GetResponse()
	}
}

// Kill attempts to terminate the process associated with the specified pid and returns true if successful.
func (t *Process) Kill(pid int) {
	t.kRouter.PostKernelMessage(messages.NewMessageProcessKill(pid))
}

// KillForeground removes the process from the foreground state and returns true if the operation succeeds.
func (t *Process) KillForeground() {
	t.kRouter.PostKernelMessage(messages.NewMessageProcessKillForeground())
}

// KillAll terminates all processes matching the provided name and returns the count of deactivated processes.
func (t *Process) KillAll(name string) {
	t.kRouter.PostKernelMessage(messages.NewMessageProcessKillAll(name))
}

// ProcessSetForeground sets the foreground process by specifying its PID and returns true if successfully set.
func (t *Process) ProcessSetForeground(pid int) {
	t.kRouter.PostKernelMessage(messages.NewMessageProcessSetForeground(pid))
}

// ProcessList returns a string representation of the process list from the kRouter.
func (t *Process) ProcessList() []*interfaces.ProcessDescription {
	msg := messages.NewMessageProcessList(t.executorWaitChan)
	t.kRouter.PostKernelMessage(msg)
	select {
	case <-t.executorWaitChan:
		return msg.GetResponse()
	case <-time.After(t.timeout):
		log.Printf("ProcessList timeout waiting for response")
		return msg.GetResponse()
	}
}

// PaintRequest sends a request to repaint the task and returns true if the request was successfully processed.
func (t *Process) PaintRequest() {
	t.kRouter.PostKernelMessage(messages.NewMessagePaintRequest())
}

// GetScreenSize returns the width and height of the screen as integers.
func (t *Process) GetScreenSize() (int, int) {
	msg := messages.NewMessageGetScreenSize(t.executorWaitChan)
	t.kRouter.PostKernelMessage(msg)
	select {
	case <-t.executorWaitChan:
		return msg.GetResponse()
	case <-time.After(t.timeout):
		log.Printf("GetScreenSize timeout waiting for response")
		return msg.GetResponse()
	}
}

// CWDSet sets the current working directory to the specified path and returns true if the operation is successful.
func (t *Process) CWDSet(arg string) bool {
	msg := messages.NewMessageFileSystemCWDSet(arg, t.executorWaitChan)
	t.kRouter.PostKernelMessage(msg)
	select {
	case <-t.executorWaitChan:
		return msg.GetResponse()
	case <-time.After(t.timeout):
		log.Printf("CWDSet timeout waiting for response")
		return msg.GetResponse()
	}
}

// CWDName returns the current working directory name by invoking a kRouter-level method.
func (t *Process) CWDName() string {
	msg := messages.NewMessageFileSystemCWDGet(t.executorWaitChan)
	t.kRouter.PostKernelMessage(msg)
	select {
	case <-t.executorWaitChan:
		return msg.GetResponse()
	case <-time.After(t.timeout):
		log.Printf("CWDName timeout waiting for response")
		return msg.GetResponse()
	}
}

// CWDPath retrieves the current working directory as a string from the associated kRouter instance.
func (t *Process) CWDPath() string {
	msg := messages.NewMessageFileSystemCWDPath(t.executorWaitChan)
	t.kRouter.PostKernelMessage(msg)
	select {
	case <-t.executorWaitChan:
		return msg.GetResponse()
	case <-time.After(t.timeout):
		log.Printf("CWDPath timeout waiting for response")
		return msg.GetResponse()
	}
}

// CWDDirectoryListing retrieves a slice of strings representing the child nodes of the current working directory (CWD).
func (t *Process) CWDDirectoryListing() []string {
	msg := messages.NewMessageFileSystemCWDDirectoryListing(t.executorWaitChan)
	t.kRouter.PostKernelMessage(msg)
	select {
	case <-t.executorWaitChan:
		return msg.GetResponse()
	case <-time.After(t.timeout):
		log.Printf("CWDDirectoryListing timeout waiting for response")
		return msg.GetResponse()
	}
}

// Suggestion provides auto-completion suggestions based on the input string and cursor position. Returns prefix, suggestions, and a success flag.
func (t *Process) Suggestion(in string, cursor int) (string, []string, bool) {
	msg := messages.NewMessageFileSystemSuggestion(in, cursor, t.executorWaitChan)
	t.kRouter.PostKernelMessage(msg)
	select {
	case <-t.executorWaitChan:
		return msg.GetResponse()
	case <-time.After(t.timeout):
		log.Printf("Suggestion timeout waiting for response")
		return msg.GetResponse()
	}
}

// Help calls the kRouter's Help method with the provided argument and returns the result or an error.
func (t *Process) Help(arg string) (string, error) {
	msg := messages.NewMessageFileSystemHelp(arg, t.executorWaitChan)
	t.kRouter.PostKernelMessage(msg)
	select {
	case <-t.executorWaitChan:
		return msg.GetResponse()
	case <-time.After(t.timeout):
		log.Printf("Help timeout waiting for response")
		return msg.GetResponse()
	}
}

// ProcessExec executes a task based on the provided command line input and returns a success status and any execution error.
func (t *Process) ProcessExec(line string) {
	t.kRouter.PostKernelMessage(messages.NewMessageProcessExec(line))
}

// WindowsSelectionBegin updates the task selection for the given process ID by invoking the kRouter's task selection method.
func (t *Process) WindowsSelectionBegin() {
	t.kRouter.PostKernelMessage(messages.NewMessageWindowsSelectionBegin())
}

// WindowsSelectionEnd finalizes the current text selection process within the Windows environment for the associated process.
func (t *Process) WindowsSelectionEnd() {
	t.kRouter.PostKernelMessage(messages.NewMessageWindowsSelectionEnd())
}

// WindowsSelectionPrevious moves the task selection pointer to the previous task in the list within the Process.
func (t *Process) WindowsSelectionPrevious() {
	t.kRouter.PostKernelMessage(messages.NewMessageWindowsSelectionPrevious())
}

// WindowsSelectionNext moves the task selection to the next task in the sequence by invoking the kRouter method.
func (t *Process) WindowsSelectionNext() {
	t.kRouter.PostKernelMessage(messages.NewMessageWindowsSelectionNext())
}

// WindowsSelectionOptions configures selection behavior for the task based on the provided option and value.
func (t *Process) WindowsSelectionOptions(option rune, value float64) {
	t.kRouter.PostKernelMessage(messages.NewMessageWindowsSelectionOptions(option, value))
}

// WritePromptEOL writes the provided prompt followed by an end-of-line character if the eol parameter is true.
func (t *Process) WritePromptEOL(prompt string, eol bool) {
	t.kRouter.PostKernelMessage(messages.NewMessageWriteColor("", interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal, eol))
	t.kRouter.PostKernelMessage(messages.NewMessageWriteColor(prompt, interfaces.ColorGreenDef, interfaces.ColorNoneDef, interfaces.ModeNormal, false))
}

// WritePromptLine sends a specified prompt and line string to the kRouter for handling the output display.
func (t *Process) WritePromptLine(prompt string, line string) {
	t.kRouter.PostKernelMessage(messages.NewMessageClearLine(line))
	t.kRouter.PostKernelMessage(messages.NewMessageWriteColor(prompt, interfaces.ColorGreenDef, interfaces.ColorNoneDef, interfaces.ModeNormal, false))
	t.kRouter.PostKernelMessage(messages.NewMessageWriteColor(line, interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal, false))
}

// Write sends the provided string data to the kRouter's write mechanism associated with the task.
func (t *Process) Write(data string, eol bool) {
	t.kRouter.PostKernelMessage(messages.NewMessageWrite(data, eol))
}

// WriteColor writes a string to the output with specified foreground, background colors, and a color mode.
func (t *Process) WriteColor(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode, eol bool) {
	t.kRouter.PostKernelMessage(messages.NewMessageWriteColor(data, fg, bg, mode, eol))
}

// WriteForeground writes the given data to the foreground with the specified color using the kRouter's functionality.
func (t *Process) WriteForeground(data string, color interfaces.ColorDef, eol bool) {
	t.kRouter.PostKernelMessage(messages.NewMessageWriteColor(data, color, interfaces.ColorNoneDef, interfaces.ModeNormal, eol))
}

// MoveCursorLeft moves the cursor one position to the left within the render context.
func (t *Process) MoveCursorLeft() {
	t.kRouter.PostKernelMessage(messages.NewMessageMoveCursorLeft())
}

// MoveCursorRight moves the cursor one position to the right by invoking the render's MoveCursorRight method.
func (t *Process) MoveCursorRight() {
	t.kRouter.PostKernelMessage(messages.NewMessageMoveCursorRight())
}

// SaveCursor saves the current cursor state by invoking the SaveCursor method on the associated renderer.
func (t *Process) SaveCursor() {
	t.kRouter.PostKernelMessage(messages.NewMessageSaveCursor())
}

// RestoreCursor restores the cursor to its previous position using the render instance of the Kernel.
func (t *Process) RestoreCursor() {
	t.kRouter.PostKernelMessage(messages.NewMessageRestoreCursor())
}

// ClearScreen clears the task's screen by delegating the request to the associated kRouter.
func (t *Process) ClearScreen() {
	t.kRouter.PostKernelMessage(messages.NewMessageClearScreen())
}

// SetExit signals the kRouter that an exit is requested for the task.
func (t *Process) SetExit() {
	t.kRouter.PostKernelMessage(messages.NewMessageExitRequested())
}

// PostUserMessage sends the provided message to the message channel for processing.
func (t *Process) PostUserMessage(msg interfaces.IMessage) {
	if len(t.gatekeeperChan) >= gatekeeperQueueMax {
		log.Printf("Process Gatekeeper: message queue full, dropping message: %d", msg.GetType())
		return
	}
	t.gatekeeperChan <- msg
}

// executorLoop initializes a loop to process messages from the executorChan and forwards a signal when ready.
func (t *Process) executorLoop(r chan bool) {
	go func() {
		r <- true
		for {
			select {
			case m, ok := <-t.executorChan:
				if !ok {
					return
				}
				t.handleMessage(m)
			}
		}
	}()
}

// evenLoop continuously listens on the message channel and processes incoming messages until a quit message is received.
func (t *Process) gatekeeperLoop(r chan bool) {
	go func() {
		r <- true
		for {
			select {
			case m, ok := <-t.gatekeeperChan:
				if !ok {
					return
				}
				if m.GetType() == interfaces.MessageTypeQuit {
					close(t.executorWaitChan)
					close(t.executorChan)
					close(t.gatekeeperChan)
					return
				}
				if !m.Ack() {
					if len(t.executorChan) >= executorQueueMax {
						log.Printf("Process Executor: message queue full, dropping message: %d", m.GetType())
						return
					}
					t.executorChan <- m
				}
			}
		}
	}()
}

// handleMessage processes incoming messages based on their type and executes corresponding logic or logs an error.
func (t *Process) handleMessage(msg interfaces.IMessage) {
	switch msg.GetType() {
	case interfaces.MessageTypeError:
		t.handleMessageError(msg)
	case interfaces.MessageTypeTimer:
		t.handleMessageTimer(msg)
	case interfaces.MessageTypeRead:
		t.handleMessageRead(msg)
	case interfaces.MessageTypeProcessStart:
		t.handleMessageProcessStart(msg)
	case interfaces.MessageTypeProcessActivate:
		t.handleMessageProcessActivate(msg)
	case interfaces.MessageTypeTimerCreated:
		t.handleMessageTimerCreated(msg)
	case interfaces.MessageTypePaintPrepare:
		t.handleMessagePaintPrepare(msg)
	default:
		log.Printf("unknown message type: %d", msg.GetType())
	}
}

// handleMessageError processes a message of type MessageError and invokes the OnError callback if defined.
func (t *Process) handleMessageError(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageError)
	if !ok {
		return
	}
	if onError := t.cmd.OnError(); onError != nil {
		onError(t, mt.Error())
	}
}

// handleMessageTimer processes a timer message and triggers the corresponding timer event if available.
func (t *Process) handleMessageTimer(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageTimer)
	if !ok {
		return
	}
	if timerEvent := t.cmd.OnTimer(); timerEvent != nil {
		timerEvent(t, mt.TID(), mt.Interval())
	}
}

// handleMessageRead processes a message read event and triggers the appropriate read or broadcast event handler.
func (t *Process) handleMessageRead(msg interfaces.IMessage) {
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
}

// handleMessageProcessStart handles the initialization of a process start message and triggers the necessary commands.
// It ensures the process is set to foreground, executes the process, and checks if it needs to run as a daemon or exit.
func (t *Process) handleMessageProcessStart(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageProcessStart)
	if !ok {
		return
	}
	_ = t.cmd.Execute(t, mt.Args())
	if !t.cmd.Background() {
		t.kRouter.PostKernelMessage(messages.NewMessageProcessSetForeground(t.PID()))
	}
	if !t.cmd.Daemon() {
		t.kRouter.PostKernelMessage(messages.NewMessageMessageProcessExit())
		return
	}
}

// handleMessageProcessActivate handles the activation message for the process and executes the defined activation callback.
func (t *Process) handleMessageProcessActivate(msg interfaces.IMessage) {
	_, ok := msg.(*messages.MessageProcessActivate)
	if !ok {
		return
	}
	if activate := t.cmd.OnActivate(); activate != nil {
		activate(t)
	}
}

// handleMessageTimerCreated processes a MessageTimerCreated message and appends the timer ID to the Process's timer list.
func (t *Process) handleMessageTimerCreated(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageTimerCreated)
	if !ok {
		return
	}
	t.timers = append(t.timers, mt.TID())
}

// handleMessagePaintPrepare processes a paint preparation message, triggering a paint event and posting a paint apply message.
func (t *Process) handleMessagePaintPrepare(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessagePaintPrepare)
	if !ok {
		return
	}
	if paintEvent := t.cmd.OnPaint(); paintEvent != nil {
		mt.Surface().Begin()
		paintEvent(t, mt.Surface())
		mt.Surface().End()

		ma := messages.NewMessagePaintApply(mt.Surface())
		t.kRouter.PostKernelMessage(ma)
	}
}
