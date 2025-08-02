package process

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"strconv"
)

// Process represents a task or job in the system, including its context, state, associated options, and execution details.
type Process struct {
	*interfaces.WindowOptions
	kernel  interfaces.IKernel
	cmd     interfaces.ICommand
	context interface{}
	timers  []int
	pid     int
	state   interfaces.ProcessState
	label   string
	caption string
	offsetY int
	offsetX int
	scale   float64
	line    string
}

// NewProcess initializes and returns a new Process instance with the provided kernel, command, and command line data.
func NewProcess(kernel interfaces.IKernel, cmd interfaces.ICommand, line string) *Process {
	t := &Process{
		kernel:  kernel,
		cmd:     cmd,
		context: nil,
		state:   interfaces.ProcessStateSetup,
		label:   cmd.Name(),
		line:    line,
		offsetX: 0,
		offsetY: 0,
		scale:   1.0,
	}
	return t
}

// SetState updates the current state of the Process to the provided State.
func (t *Process) SetState(state interfaces.ProcessState) {
	t.state = state
}

// Timers returns a slice of integers representing the current timers associated with the Process instance.
func (t *Process) Timers() []int {
	return t.timers
}

// TimersIterator iterates over all active timer IDs in the Process and executes the callback function for each one.
func (t *Process) TimersIterator(callback func(tid int) bool) {
	for _, tid := range t.timers {
		if callback(tid) {
			break
		}
	}
}

// AddTimer adds a timer identifier to the list of timers for the process.
func (t *Process) AddTimer(tid int) {
	t.timers = append(t.timers, tid)
}

// SetWindowOptions updates the Process's properties using the provided TaskOptions. It returns immediately if options is nil.
func (t *Process) SetWindowOptions(options *interfaces.WindowOptions) {
	if options == nil {
		return
	}
	t.scale = options.Scale
	t.offsetY = options.OffsetY
	t.offsetX = options.OffsetX
}

// SetWindowOption updates the task's X, Y offsets or Scale based on the given option ('x', 'y', or 'z') and value.
func (t *Process) SetWindowOption(option rune, value float64) {
	switch option {
	case 'y':
		t.SetOffsetY(t.OffsetY() + int(value))
	case 'x':
		t.SetOffsetX(t.OffsetX() + int(value))
	case 'z':
		if scale := t.Scale() + value; scale >= 0.2 && scale <= 1 {
			t.SetScale(scale)
		}
	}
}

// Description provides a brief summary of the process including its name, PID, and line information.
func (t *Process) Description() *interfaces.ProcessDescription {
	return &interfaces.ProcessDescription{
		Name: t.cmd.Name(),
		Pid:  t.pid,
		Line: t.line,
	}
}

// Options return the current task options, including offset, scale, and line settings.
func (t *Process) Options() *interfaces.WindowOptions {
	return interfaces.NewWindowOptions(t.offsetY, t.offsetX, t.scale, t.line)
}

// OffsetX returns the X-axis offset value for the Process.
func (t *Process) OffsetX() int {
	return t.offsetX
}

// SetOffsetX sets the horizontal offset (offsetX) of the task to the specified value x.
func (t *Process) SetOffsetX(x int) {
	t.offsetX = x
}

// OffsetY returns the current vertical offset value for the task.
func (t *Process) OffsetY() int {
	return t.offsetY
}

// SetOffsetY sets the vertical offset value for the task.
func (t *Process) SetOffsetY(y int) {
	t.offsetY = y
}

// Scale returns the current scale factor of the task. It determines the zoom level or relative size adjustment.
func (t *Process) Scale() float64 {
	return t.scale
}

// SetScale sets the scale factor for the Process object to the specified value.
func (t *Process) SetScale(scale float64) {
	t.scale = scale
}

// Line returns the line configuration of the Process as a string.
func (t *Process) Line() string {
	return t.line
}

// PID returns the process ID (PID) associated with the task.
func (t *Process) PID() int {
	return t.pid
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
func (t *Process) CreateTimer(first int, interval int, count int) bool {
	return t.kernel.CallCreateTimer(t.pid, first, interval, count)
}

// StopTimer stops a timer identified by the given timer ID (tid) for the current task and returns true if successful.
func (t *Process) StopTimer(tid int) bool {
	return t.kernel.CallStopTimer(t.pid, tid)
}

// IsActive checks if the process with the specified PID is currently active in the kernel.
func (t *Process) IsActive(pid int) bool {
	return t.kernel.CallIsActive(pid)
}

// DeactivateForeground removes the process from the foreground state and returns true if the operation succeeds.
func (t *Process) DeactivateForeground() bool {
	return t.kernel.CallProcessKillForeground()
}

// Deactivate attempts to terminate the task associated with the specified pid and returns true if successful.
func (t *Process) Deactivate(pid int) bool {
	return t.kernel.CallProcessKill(pid)
}

// DeactivateAll terminates all tasks matching the provided name and returns the count of deactivated tasks.
func (t *Process) DeactivateAll(name string) int {
	return t.kernel.CallProcessKillAll(name)
}

// SaveTasks saves the current state of tasks with the provided name and returns true if successful, false otherwise.
//func (t *Process) SaveTasks(name string) bool {
//	return t.kernel.CallProcessSaveAll(name)
//}

// RestoreTasks restores the state of tasks from the given name and returns true if successful, false otherwise.
//func (t *Process) RestoreTasks(name string) bool {
//	return t.kernel.CallTaskRestoreAll(name)
//}

// ListTasks retrieves a list of all currently active task names managed by the kernel.
//func (t *Process) ListTasks() []string {
//	return t.kernel.CallTaskSavedList()
//}

// SetCaption updates the task's caption using a provided string and task ID, returning true to indicate successful update.
func (t *Process) SetCaption(caption string) bool {
	t.label = caption
	t.caption = strconv.Itoa(t.pid)
	if len(t.label) > 0 {
		t.caption += " - " + t.label
	}
	return true
}

// PaintRequest sends a request to repaint the task and returns true if the request was successfully processed.
func (t *Process) PaintRequest() bool {
	return t.kernel.CallPaintRequest()
}

// GetScreenSize returns the width and height of the screen as integers.
func (t *Process) GetScreenSize() (int, int) {
	return t.kernel.CallScreenSize()
}

// CWDSet sets the current working directory to the specified path and returns true if the operation is successful.
func (t *Process) CWDSet(arg string) bool {
	return t.kernel.CallCWDSet(arg)
}

// CWDName returns the current working directory name by invoking a kernel-level method.
func (t *Process) CWDName() string {
	return t.kernel.CallCWDName()
}

// CWDGet retrieves the current working directory as a string from the associated kernel instance.
func (t *Process) CWDGet() string {
	return t.kernel.CallCWDGet()
}

// CWDPath retrieves the current working directory path as a slice of strings from the kernel.
func (t *Process) CWDPath() []string {
	return t.kernel.CallCWDPath()
}

// CWDDirectoryListing retrieves a slice of strings representing the child nodes of the current working directory (CWD).
func (t *Process) CWDDirectoryListing() []string {
	return t.kernel.CallCWDDirectoryListing()
}

// Suggestion provides auto-completion suggestions based on the input string and cursor position. Returns prefix, suggestions, and a success flag.
func (t *Process) Suggestion(in string, cursor int) (string, []string, bool) {
	return t.kernel.CallSuggestion(in, cursor)
}

// Help calls the kernel's Help method with the provided argument and returns the result or an error.
func (t *Process) Help(arg string) (string, error) {
	return t.kernel.CallHelp(arg)
}

// ProcessExec executes a task based on the provided command line input and returns a success status and any execution error.
func (t *Process) ProcessExec(line string) (bool, error) {
	return t.kernel.CallProcessExec(line, nil)
}

// WindowsSelectionBegin updates the task selection for the given process ID by invoking the kernel's task selection method.
func (t *Process) WindowsSelectionBegin() {
	t.kernel.CallWindowsSelectionBegin()
}

// WindowsSelectionEnd finalizes the current text selection process within the Windows environment for the associated process.
func (t *Process) WindowsSelectionEnd() {
	t.kernel.CallWindowsSelectionEnd()
}

// WindowsSelectionPrevious moves the task selection pointer to the previous task in the list within the Process.
func (t *Process) WindowsSelectionPrevious() {
	t.kernel.CallWindowsSelectionPrevious()
}

// WindowsSelectionNext moves the task selection to the next task in the sequence by invoking the kernel method.
func (t *Process) WindowsSelectionNext() {
	t.kernel.CallWindowsSelectionNext()
}

// WindowsSelectionOptions configures selection behavior for the task based on the provided option and value.
func (t *Process) WindowsSelectionOptions(option rune, value float64) bool {
	return t.kernel.CallWindowsSelectionOptions(option, value)
}

// SetId sets the task's process ID, updates the caption, and appends the label if it exists.
func (t *Process) SetId(id int) {
	t.pid = id
	t.caption = strconv.Itoa(t.pid)
	if len(t.label) > 0 {
		t.caption += " - " + t.label
	}
}

// Paint executes the rendering logic for the task on the provided surface by invoking a paint function if defined.
func (t *Process) Paint(surface interfaces.ISurface) {
	fn := t.cmd.PaintEvent()
	if fn == nil {
		return
	}
	surface.SetOffsetX(t.offsetX)
	surface.SetOffsetY(t.offsetY)
	surface.SetScale(t.scale)
	surface.SetCaption(t.caption)
	surface.Begin()
	fn(t, surface)
	surface.End()
}

// ProcessSetFg sets the foreground task by specifying its PID and returns true if successfully set.
func (t *Process) ProcessSetFg(pid int) bool {
	return t.kernel.CallSetFg(pid)
}

// ProcessList returns a string representation of the task list from the kernel.
func (t *Process) ProcessList() []*interfaces.ProcessDescription {
	return t.kernel.CallProcessList()
}

// WritePromptEOL writes the provided prompt followed by an end-of-line character if the eol parameter is true.
func (t *Process) WritePromptEOL(prompt string, eol bool) {
	t.kernel.CallWritePromptEOL(prompt, eol)
}

// WritePromptLine sends a specified prompt and line string to the kernel for handling the output display.
func (t *Process) WritePromptLine(prompt string, line string) {
	t.kernel.CallWritePromptLine(prompt, line)
}

// Write sends the provided string data to the kernel's write mechanism associated with the task.
func (t *Process) Write(data string) {
	t.kernel.CallWrite(data)
}

// WriteLn writes the specified data followed by a new line to the task's output stream via the kernel.
func (t *Process) WriteLn(data string) {
	t.kernel.CallWriteLn(data)
}

// WriteColor writes a string to the output with specified foreground, background colors, and a color mode.
func (t *Process) WriteColor(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	t.kernel.CallWriteColor(data, fg, bg, mode)
}

// WriteColorLn writes the provided data as a line with specified foreground and background colors and color mode.
func (t *Process) WriteColorLn(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	t.kernel.CallWriteColorLn(data, fg, bg, mode)
}

// WriteNormal sends the provided data string to the kernel's CallWriteNormal method for processing.
func (t *Process) WriteNormal(data string) {
	t.kernel.CallWriteNormal(data)
}

// WriteHighlights writes the provided highlight information to the underlying kernel for further processing.
func (t *Process) WriteHighlights(data string) {
	t.kernel.CallWriteHighlights(data)
}

// WriteCritical sends critical data to the kernel for logging or processing in a thread-safe manner.
func (t *Process) WriteCritical(data string) {
	t.kernel.CallWriteCritical(data)
}

// MoveCursorLeft moves the cursor one position to the left within the render context.
func (t *Process) MoveCursorLeft() {
	t.kernel.CallMoveCursorLeft()
}

// MoveCursorRight moves the cursor one position to the right by invoking the render's MoveCursorRight method.
func (t *Process) MoveCursorRight() {
	t.kernel.CallMoveCursorRight()
}

// SaveCursor saves the current cursor state by invoking the SaveCursor method on the associated renderer.
func (t *Process) SaveCursor() {
	t.kernel.CallSaveCursor()
}

// RestoreCursor restores the cursor to its previous position using the render instance of the Kernel.
func (t *Process) RestoreCursor() {
	t.kernel.CallRestoreCursor()
}

// ClearScreen clears the task's screen by delegating the request to the associated kernel.
func (t *Process) ClearScreen() {
	t.kernel.CallClearScreen()
}

// SetExit signals the kernel that an exit is requested for the task.
func (t *Process) SetExit() {
	t.kernel.CallExitRequested()
}

// History triggers a historical operation on the task using the specified action and index.
func (t *Process) History(verb interfaces.HistoryAction, idx int) {
	t.kernel.CallHistory(verb, idx)
}
