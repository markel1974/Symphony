package core

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"strconv"
)

// TaskState represents the state of a task, defined as an integer-based enumerated type.
type TaskState int

// TaskStateSetup represents the initial setup state of a task.
// TaskStateRunning represents the state where a task is currently running.
const (
	TaskStateSetup   TaskState = iota
	TaskStateRunning TaskState = iota
)

// TaskOptions represents configurable parameters for a task, including offsets, scaling, and associated command line.
type TaskOptions struct {
	OffsetY int
	OffsetX int
	Scale   float64
	Line    string
}

// Task represents a task instance embedding TaskOptions and providing states, labels, and related configurations.
type Task struct {
	*TaskOptions
	kernel  *Kernel
	cmd     interfaces.ICommand
	context interface{}
	timers  []int
	pid     int
	state   TaskState
	label   string
	caption string
	offsetY int
	offsetX int
	scale   float64
	line    string
}

// NewTask creates a new Task instance, initializes it with the provided Kernel, ICommand, and input line, and returns it.
func NewTask(kernel *Kernel, cmd interfaces.ICommand, line string) *Task {
	t := &Task{
		kernel:  kernel,
		cmd:     cmd,
		context: nil,
		state:   TaskStateSetup,
		label:   cmd.Name(),
		line:    line,
		offsetX: 0,
		offsetY: 0,
		scale:   1.0,
	}
	return t
}

// SetOptions updates the Task's properties using the provided TaskOptions. It returns immediately if options is nil.
func (t *Task) SetOptions(options *TaskOptions) {
	if options == nil {
		return
	}
	t.scale = options.Scale
	t.offsetY = options.OffsetY
	t.offsetX = options.OffsetX
}

// Options returns the current task options, including offset, scale, and line settings.
func (t *Task) Options() *TaskOptions {
	return &TaskOptions{
		OffsetY: t.offsetY,
		OffsetX: t.offsetX,
		Scale:   t.scale,
		Line:    t.line,
	}
}

// OffsetX returns the X-axis offset value for the Task.
func (t *Task) OffsetX() int {
	return t.offsetX
}

// SetOffsetX sets the horizontal offset (offsetX) of the task to the specified value x.
func (t *Task) SetOffsetX(x int) {
	t.offsetX = x
}

// OffsetY returns the current vertical offset value for the task.
func (t *Task) OffsetY() int {
	return t.offsetY
}

// SetOffsetY sets the vertical offset value for the task.
func (t *Task) SetOffsetY(y int) {
	t.offsetY = y
}

// Scale returns the current scale factor of the task. It determines the zoom level or relative size adjustment.
func (t *Task) Scale() float64 {
	return t.scale
}

// SetScale sets the scale factor for the Task object to the specified value.
func (t *Task) SetScale(scale float64) {
	t.scale = scale
}

// Line returns the line configuration of the Task as a string.
func (t *Task) Line() string {
	return t.line
}

// PID returns the process ID (PID) associated with the task.
func (t *Task) PID() int {
	return t.pid
}

// GetCommand returns the ICommand instance associated with the Task.
func (t *Task) GetCommand() interfaces.ICommand {
	return t.cmd
}

// SetContext sets the context for the task, storing the provided context object in the task's internal context field.
func (t *Task) SetContext(ctx interface{}) {
	t.context = ctx
}

// GetContext retrieves the context associated with the Task, returning it as an interface{}.
func (t *Task) GetContext() interface{} {
	return t.context
}

// CreateTimer initializes a timer with a specified start delay, repeat interval, and count for the current task.
func (t *Task) CreateTimer(first int, interval int, count int) bool {
	return t.kernel.CreateTimer(t.pid, first, interval, count)
}

// StopTimer stops a timer identified by the given timer ID (tid) for the current task and returns true if successful.
func (t *Task) StopTimer(tid int) bool {
	return t.kernel.StopTimer(t.pid, tid)
}

// IsActive checks if the process with the specified PID is currently active in the kernel.
func (t *Task) IsActive(pid int) bool {
	return t.kernel.IsActive(pid)
}

// Deactivate attempts to terminate the task associated with the specified pid and returns true if successful.
func (t *Task) Deactivate(pid int) bool {
	return t.kernel.Kill(pid)
}

// DeactivateAll terminates all tasks matching the provided name and returns the count of deactivated tasks.
func (t *Task) DeactivateAll(name string) int {
	return t.kernel.KillAll(name)
}

// SaveTasks saves the current state of tasks with the provided name and returns true if successful, false otherwise.
func (t *Task) SaveTasks(name string) bool {
	return t.kernel.SaveTasks(name)
}

// RestoreTasks restores the state of tasks from the given name and returns true if successful, false otherwise.
func (t *Task) RestoreTasks(name string) bool {
	return t.kernel.RestoreTasks(name)
}

// ListTasks retrieves a list of all currently active task names managed by the kernel.
func (t *Task) ListTasks() []string {
	return t.kernel.ListTasks()
}

// SetCaption updates the task's caption using a provided string and task ID, returning true to indicate successful update.
func (t *Task) SetCaption(caption string) bool {
	t.label = caption
	t.caption = strconv.Itoa(t.pid)
	if len(t.label) > 0 {
		t.caption += " - " + t.label
	}
	return true
}

// PaintRequest sends a request to repaint the task and returns true if the request was successfully processed.
func (t *Task) PaintRequest() bool {
	return t.kernel.PaintRequest()
}

// GetScreenSize returns the width and height of the screen as integers.
func (t *Task) GetScreenSize() (int, int) {
	return t.kernel.GetScreenSize()
}

// CWD returns the current working directory command associated with the task.
func (t *Task) CWD() interfaces.ICommand {
	return t.kernel.CWD()
}

// CWDSet sets the current working directory to the specified path and returns true if the operation is successful.
func (t *Task) CWDSet(arg string) bool {
	b := t.kernel.CWDSet(arg)
	return b
}

// CWDGet retrieves the current working directory as a string from the associated kernel instance.
func (t *Task) CWDGet() string {
	return t.kernel.CWDGet()
}

// CWDPath retrieves the current working directory path as a slice of strings from the kernel.
func (t *Task) CWDPath() []string {
	return t.kernel.CWDPath()
}

// CWDDirectoryListing retrieves a slice of strings representing the child nodes of the current working directory (CWD).
func (t *Task) CWDDirectoryListing() []string {
	return t.kernel.CWDDirectoryListing()
}

// Help calls the kernel's Help method with the provided argument and returns the result or an error.
func (t *Task) Help(arg string) (string, error) {
	return t.kernel.Help(arg)
}

// SetSelectionMode changes the task's selection mode to the specified process ID by interacting with the kernel.
func (t *Task) SetSelectionMode(pid int) {
	t.kernel.SetSelectionMode(pid)
}

// SetSelectionModeNext switches the selection mode to the next option using the kernel's selection handling functionality.
func (t *Task) SetSelectionModeNext() {
	t.kernel.SetSelectionModePrevious()
}

// SetSelectionModePrevious sets the selection mode to the previous item using the underlying kernel functionality.
func (t *Task) SetSelectionModePrevious() {
	t.kernel.SetSelectionModePrevious()
}

// SetSelectionOptions configures selection behavior for the task based on the provided option and value.
func (t *Task) SetSelectionOptions(option rune, value float64) bool {
	return t.kernel.SetSelectionOptions(option, value)
}

// SetId sets the task's process ID, updates the caption, and appends the label if it exists.
func (t *Task) SetId(id int) {
	t.pid = id
	t.caption = strconv.Itoa(t.pid)
	if len(t.label) > 0 {
		t.caption += " - " + t.label
	}
}

// Paint executes the rendering logic for the task on the provided surface by invoking a paint function if defined.
func (t *Task) Paint(surface interfaces.ISurface) {
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

// SetFg sets the foreground task by specifying its PID and returns true if successfully set.
func (t *Task) SetFg(pid int) bool {
	return t.kernel.SetFg(pid)
}

// TaskList returns a string representation of the task list from the kernel.
func (t *Task) TaskList() string {
	return t.kernel.List()
}

// Write sends the provided string data to the kernel's write mechanism associated with the task.
func (t *Task) Write(data string) {
	t.kernel.Write(data)
}

// WriteLn writes the specified data followed by a new line to the task's output stream via the kernel.
func (t *Task) WriteLn(data string) {
	t.kernel.WriteLn(data)
}

// WriteColor writes a string to the output with specified foreground, background colors, and a color mode.
func (t *Task) WriteColor(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	t.kernel.WriteColor(data, fg, bg, mode)
}

// WriteColorLn writes the provided data as a line with specified foreground and background colors and color mode.
func (t *Task) WriteColorLn(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	t.kernel.WriteColorLn(data, fg, bg, mode)
}

// ClearScreen clears the task's screen by delegating the request to the associated kernel.
func (t *Task) ClearScreen() {
	t.kernel.ClearScreen()
}

// SetExit signals the kernel that an exit is requested for the task.
func (t *Task) SetExit() {
	t.kernel.ExitRequested()
}

// History triggers a historical operation on the task using the specified action and index.
func (t *Task) History(verb interfaces.HistoryAction, idx int) {
	t.kernel.HistoryApply(verb, idx)
}
