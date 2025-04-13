package context

import (
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"io"
	"strconv"
)

// TaskState represents the state of a task, typically managed by a TaskManager within an application.
type TaskState int

// TaskStateSetup represents the initial setup state of a task.
// TaskStateRunning represents the running state of a task.
const (
	TaskStateSetup   TaskState = iota
	TaskStateRunning TaskState = iota
)

type TaskOptions struct {
	OffsetY int
	OffsetX int
	Scale   float64
	Line    string
}

// Task represents a unit of work managed by a TaskManager and associated with a specific context, command, and state.
type Task struct {
	*TaskOptions
	tasks   *TaskManager
	ctx     *Context
	cmd     interfaces.ICommand
	context interface{}
	timers  []int
	pid     int
	state   TaskState
	caption string
	offsetY int
	offsetX int
	scale   float64
	line    string
}

// NewTask creates a new Task instance associated with TaskManager, Context, a command, and initializes its state.
func NewTask(tasks *TaskManager, ctx *Context, cmd interfaces.ICommand, line string) *Task {
	t := &Task{
		tasks:   tasks,
		ctx:     ctx,
		cmd:     cmd,
		context: nil,
		state:   TaskStateSetup,
		caption: "",
		line:    line,
		offsetX: 0,
		offsetY: 0,
		scale:   1.0,
	}
	return t
}

func (t *Task) SetOptions(options *TaskOptions) {
	if options == nil {
		return
	}
	t.scale = options.Scale
	t.offsetY = options.OffsetY
	t.offsetX = options.OffsetX
}

func (t *Task) Options() *TaskOptions {
	return &TaskOptions{
		OffsetY: t.offsetY,
		OffsetX: t.offsetX,
		Scale:   t.scale,
		Line:    t.line,
	}
}

func (t *Task) OffsetX() int {
	return t.offsetX
}

func (t *Task) SetOffsetX(x int) {
	t.offsetX = x
}

func (t *Task) OffsetY() int {
	return t.offsetY
}

func (t *Task) SetOffsetY(y int) {
	t.offsetY = y
}

func (t *Task) Scale() float64 {
	return t.scale
}

func (t *Task) SetScale(scale float64) {
	t.scale = scale
}

func (t *Task) Line() string {
	return t.line
}

// PID returns the process identifier (PID) of the task.
func (t *Task) PID() int {
	return t.pid
}

// GetCommand returns the ICommand instance associated with the Task.
func (t *Task) GetCommand() interfaces.ICommand {
	return t.cmd
}

// SetContext sets the context of the Task to the provided value.
func (t *Task) SetContext(ctx interface{}) {
	t.context = ctx
}

// GetContext retrieves the current context associated with the Task instance.
// It returns the context as an interface{}.
func (t *Task) GetContext() interface{} {
	return t.context
}

// CreateTimer initializes a timer with specified delay, interval, and repetition count for the current task. Returns success status.
func (t *Task) CreateTimer(first int, interval int, count int) bool {
	return t.tasks.CreateTimer(t.pid, first, interval, count)
}

// StopTimer stops a timer with the given timer ID (tid) associated with the task and returns true if successful.
func (t *Task) StopTimer(tid int) bool {
	return t.tasks.StopTimer(t.pid, tid)
}

// IsActive checks if the task with the given process ID (pid) is currently active. Returns true if active, false otherwise.
func (t *Task) IsActive(pid int) bool {
	return t.tasks.IsActive(pid)
}

// Deactivate terminates the task identified by the given pid and returns true if successfully terminated, false otherwise.
func (t *Task) Deactivate(pid int) bool {
	return t.tasks.Kill(pid)
}

// DeactivateAll deactivates all tasks matching the provided name and returns the count of tasks that were deactivated.
func (t *Task) DeactivateAll(name string) int {
	return t.tasks.KillAll(name)
}

// SaveTasks saves the current tasks to a storage medium using the provided name and returns true if successful.
func (t *Task) SaveTasks(name string) bool {
	return t.tasks.SaveTasks(name)
}

// RestoreTasks restores previously saved tasks identified by the provided name. Returns true if successful, false otherwise.
func (t *Task) RestoreTasks(name string) bool {
	return t.tasks.RestoreTasks(name)
}

// ListTasks retrieves the list of all task names currently managed by the TaskManager. It returns a slice of strings.
func (t *Task) ListTasks() []string {
	return t.tasks.ListTasks()
}

// SetCaption sets the caption for the task and returns true if the operation was successful, otherwise false.
func (t *Task) SetCaption(caption string) bool {
	return t.tasks.SetCaption(t.pid, caption)
}

// PaintRequest forwards a paint request to the TaskManager and returns a boolean indicating success.
func (t *Task) PaintRequest() bool {
	return t.tasks.PaintRequest()
}

// GetScreenSize returns the width and height of the screen as two integer values.
func (t *Task) GetScreenSize() (int, int) {
	return t.tasks.GetScreenSize()
}

// CWD returns the current working directory command interface associated with the task.
func (t *Task) CWD() interfaces.ICommand {
	return t.tasks.CWD()
}

// CWDSet sets the current working directory for the task to the specified path and returns true if successful.
func (t *Task) CWDSet(arg string) bool {
	return t.tasks.CWDSet(arg)
}

// CWDGet retrieves the current working directory as a string for the task.
func (t *Task) CWDGet() string {
	return t.tasks.CWDGet()
}

// CWDPath retrieves the current working directory path as a slice of strings. It delegates the call to the underlying TaskManager.
func (t *Task) CWDPath() []string {
	return t.tasks.CWDPath()
}

// CWDChilds retrieves the names of child directories or files under the current working directory.
func (t *Task) CWDChilds() []string {
	return t.tasks.CWDChilds()
}

func (t *Task) Help(arg string) (string, error) {
	return t.tasks.Help(arg)
}

// SetSelectionMode sets the selection mode for a task identified by the given pid.
func (t *Task) SetSelectionMode(pid int) {
	t.tasks.SetSelectionMode(pid)
}

// SetSelectionModeNext moves the task selection mode to the next item via the task manager.
func (t *Task) SetSelectionModeNext() {
	t.tasks.SetSelectionModePrevious()
}

// SetSelectionModePrevious sets the selection mode to the previous task in the task manager.
func (t *Task) SetSelectionModePrevious() {
	t.tasks.SetSelectionModePrevious()
}

// SetSelectionOptions configures selection mode options using the given option rune and value. Returns true if successful.
func (t *Task) SetSelectionOptions(option rune, value float64) bool {
	return t.tasks.SetSelectionOptions(option, value)
}

// SetId sets the task's process ID (PID) to the provided integer value.
func (t *Task) SetId(id int) {
	t.pid = id
}

// Paint renders the task's visual representation on the provided surface by applying transformations and calling the paint function.
func (t *Task) Paint(surface *Surface) {
	fn := t.cmd.PaintEvent()
	if fn == nil {
		return
	}
	caption := strconv.Itoa(t.pid)
	if len(t.caption) > 0 {
		caption += " - " + t.caption
	}
	surface.SetOffsetX(t.offsetX)
	surface.SetOffsetY(t.offsetY)
	surface.SetScale(t.scale)
	surface.SetCaption(caption)
	surface.Begin()
	fn(t, surface)
	surface.End()
}

// GetWriter returns an io.Writer instance associated with the current task's context.
func (t *Task) GetWriter() io.Writer {
	return t.ctx.GetWriter()
}

// SetFg sets the specified process ID (pid) as the foreground task and returns true if successful.
func (t *Task) SetFg(pid int) bool {
	return t.tasks.SetFg(pid)
}

// TaskList returns a string representation of all tasks managed by the TaskManager.
func (t *Task) TaskList() string {
	return t.tasks.List()
}

// Write sends the provided string `data` to the underlying context for writing.
func (t *Task) Write(data string) {
	t.ctx.Write(data)
}

// WriteLn writes a line of text (data) to the task's associated context, followed by a newline.
func (t *Task) WriteLn(data string) {
	t.ctx.WriteLn(data)
}

// WriteColor outputs the given text with specified foreground and background colors and the specified color mode.
func (t *Task) WriteColor(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	t.ctx.WriteColor(data, fg, bg, mode)
}

// WriteColorLn writes the provided string to the output with the specified foreground and background colors in the given mode.
func (t *Task) WriteColorLn(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	t.ctx.WriteColorLn(data, fg, bg, mode)
}

// ClearScreen clears the screen of the current context associated with the task.
func (t *Task) ClearScreen() {
	t.ctx.ClearScreen()
}

// SetExit signals the context to exit, marking the task for termination or indicating it should stop execution immediately.
func (t *Task) SetExit() {
	t.ctx.SetExit()
}

// History performs a history-related action specified by the provided verb and index on the task's context.
func (t *Task) History(verb interfaces.HistoryAction, idx int) {
	t.ctx.History(verb, idx)
}
