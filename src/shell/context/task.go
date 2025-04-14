package context

import (
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"strconv"
)

// TaskState represents the state of a task within a task management or process control system.
type TaskState int

// TaskStateSetup represents the state where a task is being prepared for execution.
// TaskStateRunning represents the state where a task is currently executing.
const (
	TaskStateSetup   TaskState = iota
	TaskStateRunning TaskState = iota
)

// TaskOptions defines configuration options for a task, including positional offsets, scaling, and associated command line.
type TaskOptions struct {
	OffsetY int
	OffsetX int
	Scale   float64
	Line    string
}

// Task represents a unit of work associated with a command, managing state, context, timers, offsets, and scaling.
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

// NewTask initializes and returns a new Task instance with the given Kernel, ICommand, and command line string input.
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

// SetOptions updates the Task's scale, offsetY, and offsetX properties based on the provided TaskOptions.
func (t *Task) SetOptions(options *TaskOptions) {
	if options == nil {
		return
	}
	t.scale = options.Scale
	t.offsetY = options.OffsetY
	t.offsetX = options.OffsetX
}

// Options returns a pointer to a TaskOptions struct populated with the current task's offset, scale, and line properties.```go
// Options returns a pointer to a TaskOptions object with values set based on the current Task's properties.
func (t *Task) Options() *TaskOptions {
	return &TaskOptions{
		OffsetY: t.offsetY,
		OffsetX: t.offsetX,
		Scale:   t.scale,
		Line:    t.line,
	}
}

// OffsetX returns the current horizontal offset of the task.
func (t *Task) OffsetX() int {
	return t.offsetX
}

// SetOffsetX sets the horizontal offset of the Task to the specified value x.
func (t *Task) SetOffsetX(x int) {
	t.offsetX = x
}

// OffsetY returns the vertical offset value for the task.
func (t *Task) OffsetY() int {
	return t.offsetY
}

// SetOffsetY sets the vertical offset value for the task to the specified value y.
func (t *Task) SetOffsetY(y int) {
	t.offsetY = y
}

// Scale returns the current scale factor of the Task.
func (t *Task) Scale() float64 {
	return t.scale
}

// SetScale updates the scale attribute of the task to the specified value.
func (t *Task) SetScale(scale float64) {
	t.scale = scale
}

// Line retrieves the line property of the Task, returning it as a string.
func (t *Task) Line() string {
	return t.line
}

// PID returns the process ID (PID) associated with the Task instance.
func (t *Task) PID() int {
	return t.pid
}

// GetCommand retrieves the command associated with the task, implementing the ICommand interface.
func (t *Task) GetCommand() interfaces.ICommand {
	return t.cmd
}

// SetContext sets the context of the task to the given interface value.
func (t *Task) SetContext(ctx interface{}) {
	t.context = ctx
}

// GetContext retrieves the context object associated with the task.
func (t *Task) GetContext() interface{} {
	return t.context
}

// CreateTimer creates and starts a timer with specified initial delay, interval, and repetition count for the task.
func (t *Task) CreateTimer(first int, interval int, count int) bool {
	return t.kernel.CreateTimer(t.pid, first, interval, count)
}

// StopTimer stops a timer associated with the given timer ID (tid) for the current task and returns true if successful.
func (t *Task) StopTimer(tid int) bool {
	return t.kernel.StopTimer(t.pid, tid)
}

// IsActive checks if a process with the provided pid is active in the kernel context and returns true if active.
func (t *Task) IsActive(pid int) bool {
	return t.kernel.IsActive(pid)
}

// Deactivate terminates the task identified by the given pid, returning true if successful or false otherwise.
func (t *Task) Deactivate(pid int) bool {
	return t.kernel.Kill(pid)
}

// DeactivateAll terminates all tasks that match the specified name and returns the count of tasks deactivated.
func (t *Task) DeactivateAll(name string) int {
	return t.kernel.KillAll(name)
}

// SaveTasks saves the current tasks under the specified name and returns true if successful or false otherwise.
func (t *Task) SaveTasks(name string) bool {
	return t.kernel.SaveTasks(name)
}

// RestoreTasks attempts to restore previously saved task states associated with the specified name. Returns true on success.
func (t *Task) RestoreTasks(name string) bool {
	return t.kernel.RestoreTasks(name)
}

// ListTasks retrieves the list of all currently active task names from the kernel.
func (t *Task) ListTasks() []string {
	return t.kernel.ListTasks()
}

// SetCaption updates the task's caption using the provided label and the task's PID. Returns true upon completion.
func (t *Task) SetCaption(caption string) bool {
	t.label = caption
	t.caption = strconv.Itoa(t.pid)
	if len(t.label) > 0 {
		t.caption += " - " + t.label
	}
	return true
}

// PaintRequest checks if a paint operation is required for the task by delegating to the underlying kernel.
func (t *Task) PaintRequest() bool {
	return t.kernel.PaintRequest()
}

// GetScreenSize returns the width and height of the screen as integers.
func (t *Task) GetScreenSize() (int, int) {
	return t.kernel.GetScreenSize()
}

// CWD retrieves the current working directory command associated with the task.
func (t *Task) CWD() interfaces.ICommand {
	return t.kernel.CWD()
}

// CWDSet sets the current working directory to the specified path and returns true if the operation is successful.
func (t *Task) CWDSet(arg string) bool {
	return t.kernel.CWDSet(arg)
}

// CWDGet retrieves the current working directory string for the task from the associated kernel.
func (t *Task) CWDGet() string {
	return t.kernel.CWDGet()
}

// CWDPath returns the current working directory path as a slice of strings.
func (t *Task) CWDPath() []string {
	return t.kernel.CWDPath()
}

// CWDChilds retrieves a list of child command names in the current working directory context of the task.
func (t *Task) CWDChilds() []string {
	return t.kernel.CWDChilds()
}

// Help calls the kernel's Help method with the provided argument and returns the resulting help text or an error.
func (t *Task) Help(arg string) (string, error) {
	return t.kernel.Help(arg)
}

// SetSelectionMode configures the selection mode for the Task using the specified process ID.
func (t *Task) SetSelectionMode(pid int) {
	t.kernel.SetSelectionMode(pid)
}

// SetSelectionModeNext switches the selection mode to the next state by delegating to the kernel's previous mode setter.
func (t *Task) SetSelectionModeNext() {
	t.kernel.SetSelectionModePrevious()
}

// SetSelectionModePrevious sets the selection mode to the previous item using the kernel configuration.
func (t *Task) SetSelectionModePrevious() {
	t.kernel.SetSelectionModePrevious()
}

// SetSelectionOptions configures selection options using the provided option rune and value. Returns true on success.
func (t *Task) SetSelectionOptions(option rune, value float64) bool {
	return t.kernel.SetSelectionOptions(option, value)
}

// SetId assigns a unique process ID to the task and updates the caption based on the ID and any associated label.
func (t *Task) SetId(id int) {
	t.pid = id
	t.caption = strconv.Itoa(t.pid)
	if len(t.label) > 0 {
		t.caption += " - " + t.label
	}
}

// Paint renders the task's output on the given ISurface, initializing offsets, scale, and caption before invoking PaintEvent.
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

// SetFg sets the specified process ID (pid) as the foreground process. Returns true if successful, or false otherwise.
func (t *Task) SetFg(pid int) bool {
	return t.kernel.SetFg(pid)
}

// TaskList returns a string representation of the task list as managed by the kernel.
func (t *Task) TaskList() string {
	return t.kernel.List()
}

// Write sends the provided data string to the underlying kernel for processing or output.
func (t *Task) Write(data string) {
	t.kernel.Write(data)
}

// WriteLn writes a line of text followed by a newline to the kernel's output stream.
func (t *Task) WriteLn(data string) {
	t.kernel.WriteLn(data)
}

// WriteColor outputs the given text in specified foreground and background colors with a specific color rendering mode.
func (t *Task) WriteColor(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	t.kernel.WriteColor(data, fg, bg, mode)
}

// WriteColorLn writes a line of text with specified foreground, background colors, and color mode to the task's output.
func (t *Task) WriteColorLn(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	t.kernel.WriteColorLn(data, fg, bg, mode)
}

// ClearScreen clears the current output screen of the task by invoking the kernel's ClearScreen method.
func (t *Task) ClearScreen() {
	t.kernel.ClearScreen()
}

// SetExit signals a request to exit the kernel associated with the task.
func (t *Task) SetExit() {
	t.kernel.ExitRequested()
}

// History performs a kernel-level history operation based on the specified action and index.
func (t *Task) History(verb interfaces.HistoryAction, idx int) {
	t.kernel.History(verb, idx)
}
