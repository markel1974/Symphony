package context

import (
	"github.com/markel1974/c64emu/src/shell/adaptiveticker"
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"io"
	"strconv"
)

type Task struct {
	tasks   *TaskManager
	ctx     *Context
	cmd     interfaces.ICommand
	context interface{}
	timers  []int
	pid     int
	state   taskState
	caption string
	Line    string
	OffsetX int
	OffsetY int
	Scale   float64
}

func NewTask(tasks *TaskManager, ctx *Context, cmd interfaces.ICommand, line string) *Task {
	return &Task{
		tasks:   tasks,
		ctx:     ctx,
		cmd:     cmd,
		context: nil,
		state:   taskStateSetup,
		caption: "",
		Line:    line,
		OffsetX: 0,
		OffsetY: 0,
		Scale:   1.0,
	}
}

func (t *Task) PID() int {
	return t.pid
}

func (t *Task) GetCommand() interfaces.ICommand {
	return t.cmd
}

func (t *Task) SetContext(ctx interface{}) {
	t.context = ctx
}

func (t *Task) GetContext() interface{} {
	return t.context
}

func (t *Task) CreateTimer(first int, interval int, count int) bool {
	return t.tasks.CreateTimer(t.pid, first, interval, count)
}

func (t *Task) StopTimer(tid int) bool {
	return t.tasks.StopTimer(t.pid, tid)
}

func (t *Task) IsActive(pid int) bool {
	return t.tasks.IsActive(pid)
}

func (t *Task) Deactivate(pid int) bool {
	return t.tasks.Kill(pid)
}

func (t *Task) DeactivateAll(name string) int {
	return t.tasks.KillAll(name)
}

func (t *Task) SaveTasks(name string) bool {
	return t.tasks.SaveTasks(name)
}

func (t *Task) RestoreTasks(name string) bool {
	return t.tasks.RestoreTasks(name)
}

func (t *Task) ListTasks() []string {
	return t.tasks.ListTasks()
}
func (t *Task) SetCaption(caption string) bool {
	return t.tasks.SetCaption(t.pid, caption)
}

func (t *Task) PaintRequest() bool {
	return t.tasks.PaintRequest()
}

func (t *Task) GetScreenSize() (int, int) {
	return t.tasks.GetScreenSize()
}

func (t *Task) CWD() interfaces.ICommand {
	return t.tasks.CWD()
}

func (t *Task) CWDSet(arg string) bool {
	return t.tasks.CWDSet(arg)
}

func (t *Task) CWDGet() string {
	return t.tasks.CWDGet()
}

func (t *Task) CWDPath() []string {
	return t.tasks.CWDPath()
}

func (t *Task) CWDChilds() []string {
	return t.tasks.CWDChilds()
}

func (t *Task) SetSelectionMode(pid int) {
	t.tasks.SetSelectionMode(pid)
}

func (t *Task) SetSelectionModeNext() {
	t.tasks.SetSelectionModePrevious()
}

func (t *Task) SetSelectionModePrevious() {
	t.tasks.SetSelectionModePrevious()
}

func (t *Task) SetSelectionOptions(option rune, value float64) bool {
	return t.tasks.SetSelectionOptions(option, value)
}

func (t *Task) SetId(id int) {
	t.pid = id
}

func (t *Task) Unset() {
	t.pid = adaptiveticker.UnknownId
}

func (t *Task) Paint(surface *Surface) {
	fn := t.cmd.PaintEvent()
	if fn == nil {
		return
	}
	caption := strconv.Itoa(t.pid)
	if len(t.caption) > 0 {
		caption += " - " + t.caption
	}
	surface.SetOffsetX(t.OffsetX)
	surface.SetOffsetY(t.OffsetY)
	surface.SetScale(t.Scale)
	surface.SetCaption(caption)
	surface.Begin()
	fn(t, surface)
	surface.End()
}

func (t *Task) GetWriter() io.Writer {
	return t.ctx.GetWriter()
}

func (t *Task) SetFg(pid int) bool {
	return t.tasks.SetFg(pid)
}

func (t *Task) TaskList() string {
	return t.tasks.List()
}

func (t *Task) Write(data string) {
	t.ctx.Write(data)
}

func (t *Task) WriteLn(data string) {
	t.ctx.WriteLn(data)
}

func (t *Task) WriteColor(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	t.ctx.WriteColor(data, fg, bg, mode)
}

func (t *Task) WriteColorLn(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	t.ctx.WriteColorLn(data, fg, bg, mode)
}

func (t *Task) ClearScreen() {
	t.ctx.ClearScreen()
}

func (t *Task) SetExit() {
	t.ctx.SetExit()
}

func (t *Task) History(verb interfaces.HistoryAction, idx int) {
	t.ctx.History(verb, idx)
}
