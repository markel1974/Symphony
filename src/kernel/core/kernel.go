package core

import (
	"encoding/json"
	"fmt"
	"github.com/markel1974/c64emu/src/kernel/adaptiveticker"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/shell"
	"log"
	"os"
	"strings"
)

// tasksFileExtension specifies the file extension used for task-related data files.
const tasksFileExtension = ".task"

// commandActivate represents the command string for activation operations.
// commandTask represents the command string for task-related operations.
const (
	commandActivate = "activate"
	commandTask     = "task"
)

// Kernel represents the core component responsible for managing rendering, input/output, task execution, and timers.
type Kernel struct {
	ticker      *adaptiveticker.AdaptiveTicker
	io          interfaces.IInputOutput
	foreground  interfaces.ITask
	selector    *TaskSelector
	timersChan  chan *adaptiveticker.TimerHandler
	ids         *adaptiveticker.Ids
	fs          interfaces.IFileSystem
	sh          *shell.Shell
	messageChan chan iMessage
	exit        bool
}

// NewKernel creates and returns a new Kernel instance, initializing its dependencies and internal fields.
func NewKernel(ticker *adaptiveticker.AdaptiveTicker, io interfaces.IInputOutput, fs interfaces.IFileSystem, sh *shell.Shell) *Kernel {
	t := &Kernel{
		ticker:      ticker,
		io:          io,
		fs:          fs,
		sh:          sh,
		foreground:  nil,
		selector:    NewTaskSelector(),
		ids:         adaptiveticker.NewIds(1024),
		messageChan: make(chan iMessage, contextMaQueueLen),
		timersChan:  make(chan *adaptiveticker.TimerHandler, contextMaQueueLen),
		exit:        false,
	}
	return t
}

// NextLine advances to the next line in the shell output, optionally determining if an end-of-line character is added.
// TODO REMOVE
func (c *Kernel) NextLine(eol bool) {
	c.sh.NextLine(eol)
}

// KeyEvent processes a keyboard event of a given type and key and returns true if the event was handled successfully.
func (c *Kernel) KeyEvent(kind interfaces.KeyType, key rune) bool {
	return c.sh.KeyEvent(kind, key)
}

// Redraw refreshes the current state based on the provided input string, invoking the underlying shell's Redraw function.
func (c *Kernel) Redraw(l string) {
	c.sh.Redraw(l)
}

// HistoryApply performs actions on the command history based on the specified verb (list, clear, or execute at the given index).
func (c *Kernel) HistoryApply(verb interfaces.HistoryAction, idx int) {
	switch verb {
	case interfaces.HistoryActionClear:
		c.sh.ClearHistory()
	case interfaces.HistoryActionExec:
		if arg, found := c.sh.GetHistoryAtPos(idx); found {
			_, _ = c.ExecCommand(arg, nil)
		}
	case interfaces.HistoryActionList:
		c.sh.Write(c.sh.GetHistory())
	}
}

// SetHistoryDefault sets the default history data in the kernel's state handler using the provided string input.
func (c *Kernel) SetHistoryDefault(data string) {
	c.sh.SetHistoryDefault(data)
}

// SetScreenSize sets the display dimensions of the screen to the specified width (w) and height (h).
func (c *Kernel) SetScreenSize(w int, h int) {
	c.sh.SetScreenSize(w, h)
}

// GetScreenSize returns the width and height of the screen in pixels as two integer values.
func (c *Kernel) GetScreenSize() (int, int) {
	return c.sh.GetScreenSize()
}

// Write sends the specified string data to the kernel's renderer for processing or output.
func (c *Kernel) Write(data string) {
	c.sh.Write(data)
}

// WriteLn writes the provided string followed by a newline to the output stream.
func (c *Kernel) WriteLn(data string) {
	c.sh.WriteLn(data)
}

// WriteColor writes text with specified foreground and background colors using the provided color rendering mode.
func (c *Kernel) WriteColor(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	c.sh.WriteColor(data, fg, bg, mode)
}

// WriteColorLn writes a line of text with specified foreground color, background color, and color mode.
func (c *Kernel) WriteColorLn(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	c.sh.WriteColorLn(data, fg, bg, mode)
}

// ClearScreen clears the screen by invoking the render system's ClearScreen operation.
func (c *Kernel) ClearScreen() {
	c.sh.ClearScreen()
}

// ExecActivate attempts to activate a foreground process by executing the associated command and returns its success status.
func (c *Kernel) ExecActivate() bool {
	pid, name := c.GetForegroundName()
	if pid == adaptiveticker.UnknownId {
		return false
	}
	if name == commandActivate {
		return false
	}
	c.SetBackground()
	_, _ = c.ExecCommand(fmt.Sprint(commandActivate, " ", pid), nil)
	return false
}

// ExecCommand executes a command by parsing the input line, creating a task, and managing its lifecycle and state.
// Returns true and error if the task was created but execution failed, or true and nil if execution succeeded.
func (c *Kernel) ExecCommand(line string, options *TaskOptions) (bool, error) {
	cmd, args, err := c.fs.Find(line)
	if err != nil {
		return false, fmt.Errorf("error creating task: invalid command '%s'", line)
	}
	task := NewTask(c, cmd, line)
	if !c.ids.Set(task) {
		return false, fmt.Errorf("error creating task: can't set pid")
	}
	task.SetOptions(options)
	if err = cmd.Execute(task, args); err != nil {
		c.Kill(task.pid)
		return true, err
	}
	if !cmd.Daemon() {
		c.Kill(task.pid)
		return true, nil
	}
	task.state = TaskStateRunning
	if !cmd.Background() {
		c.foreground = task
	}
	return true, nil
}

// SetSelectionMode sets the current selection mode for tasks based on a requested process ID. Defaults to the first task if the requested ID is unavailable.
func (c *Kernel) SetSelectionMode(requestedPid int) {
	var idx = 0
	var firstPid = adaptiveticker.UnknownId
	var firstIdx = 0

	c.selector.Clear()

	c.ids.Range(func(item adaptiveticker.IIds) bool {
		task, ok := item.(*Task)
		if ok && task != nil {
			if task.cmd.PaintEvent() != nil {
				c.selector.AddAvailable(task.pid)
				if firstPid == adaptiveticker.UnknownId {
					firstPid = task.PID()
					firstIdx = idx
				}
				if task.PID() == requestedPid {
					c.selector.Set(requestedPid, idx)
				}
				idx++
			}
		}
		return true
	})
	if c.selector.PID() == adaptiveticker.UnknownId {
		if firstPid == adaptiveticker.UnknownId {
			return
		}
		c.selector.Set(firstPid, firstIdx)
	}
	c.doPaintRequest(false)
}

// PaintRequest determines if a paint operation is required by invoking an underlying kernel method.
func (c *Kernel) PaintRequest() bool {
	return c.doPaintRequest(false)
}

// doPaintRequest handles paint requests by invoking the render's PaintRequest method and potentially scheduling a paint event.
// Returns true if a paint event was successfully scheduled, otherwise returns false.
// The full parameter specifies whether the entire view should be repainted.
func (c *Kernel) doPaintRequest(full bool) bool {
	if c.sh.PaintRequest(full) {
		c.ticker.Create(c.timersChan, newMessagePaint(), -1, -1, 1)
		return true
	}
	return false
}

// SetSelectionModeNext advances to the next available task in the selection and triggers a repaint request if successful.
func (c *Kernel) SetSelectionModeNext() {
	if !c.selector.Next() {
		return
	}
	c.doPaintRequest(false)
}

// SetSelectionModePrevious updates selection to the previous item using the internal selector and requests a repaint if successful.
func (c *Kernel) SetSelectionModePrevious() {
	if !c.selector.Prev() {
		return
	}
	c.doPaintRequest(false)
}

// SetSelectionDisabled disables task selection by clearing the selector's state and cancels any pending paint requests.
func (c *Kernel) SetSelectionDisabled() {
	c.selector.Clear()
	c.doPaintRequest(false)
}

// ExitRequested signals that the kernel should terminate its execution and exit.
func (c *Kernel) ExitRequested() {
	c.exit = true
}

// CWDDirectoryListing retrieves the names of the child elements in the current working directory.
func (c *Kernel) CWDDirectoryListing() []string {
	var out []string

	cwd := c.fs.CWD()
	for _, z := range cwd.DirectoryListing() {
		out = append(out, z) // z.Name())
	}
	return out
}

// CWD returns the current working directory command from the kernel's file system interface.
func (c *Kernel) CWD() interfaces.ICommand {
	return c.fs.CWD()
}

// CWDGet returns the command path of the current working directory from the kernel's file system.
func (c *Kernel) CWDGet() string {
	return c.fs.CWD().CommandPath()
}

// CWDPath returns the current working directory path as a slice of strings.
func (c *Kernel) CWDPath() []string {
	return c.fs.CWD().Path()
}

// CWDSet sets the current working directory to the specified path and returns true if the operation is successful.
func (c *Kernel) CWDSet(arg string) bool {
	return c.fs.CWDSet(arg)
}

// Help retrieves the help documentation associated with the provided argument from the kernel's filesystem.
func (c *Kernel) Help(arg string) (string, error) {
	return c.fs.Help(arg)
}

// SetSelectionOptions modifies selection parameters like offsets or scale based on the provided option and value.
// Returns true if the operation is applied successfully; otherwise, returns false.
func (c *Kernel) SetSelectionOptions(option rune, value float64) bool {
	t, ok := c.ids.Get(c.selector.PID())
	if !ok {
		return false
	}
	task := t.(*Task)
	switch option {
	case 'y':
		task.SetOffsetY(task.OffsetY() + int(value))
	case 'x':
		task.SetOffsetX(task.OffsetX() + int(value))
	case 'z':
		if scale := task.Scale() + value; scale >= 0.2 && scale <= 1 {
			task.SetScale(scale)
		}
	}
	c.doPaintRequest(true)
	return true
}

// SetFg sets the task with the given pid as the foreground task. Returns true if successful, false if the pid is invalid.
func (c *Kernel) SetFg(pid int) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}
	task := t.(*Task)
	c.foreground = task
	return true
}

// GetSuggestion provides auto-complete suggestions based on the given input and cursor position, returning a prefix, suggestions, and a success flag.
func (c *Kernel) GetSuggestion(in string, cursor int) (string, []string, bool) {
	prefix, suggestions, found := c.fs.Suggestion(in, cursor)
	return prefix, suggestions, found
}

// CreateTimer initializes a timer for a process based on its ID with specified delay, interval, and count settings.
// Returns true if the timer was created successfully, false otherwise.
func (c *Kernel) CreateTimer(pid int, first int, interval int, count int) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}
	task := t.(*Task)
	if task.cmd.TimerEvent == nil {
		return false
	}
	m := newMessageTimer(pid, interval)
	m.tid = c.ticker.Create(c.timersChan, m, int64(first), int64(interval), int64(count))
	if m.tid > -1 {
		task.timers = append(task.timers, m.tid)
	}
	return true
}

// StopTimer stops a timer identified by the task ID (tid) within the process ID (pid). Returns true if successful, false otherwise.
func (c *Kernel) StopTimer(pid int, tid int) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}
	task := t.(*Task)
	return c.closeTimer(task, tid)
}

// IsActive checks if the specified process ID (pid) is currently active in the kernel and returns true if found.
func (c *Kernel) IsActive(pid int) bool {
	_, ret := c.ids.Get(pid)
	return ret
}

// GetForegroundPid retrieves the process ID of the currently active foreground process, or UnknownId if none is active.
func (c *Kernel) GetForegroundPid() int {
	if c.foreground == nil {
		return adaptiveticker.UnknownId
	}
	return c.foreground.PID()
}

// GetForegroundName retrieves the process ID and command name of the current foreground process.
// Returns a known invalid ID and empty string if no foreground process is set.
func (c *Kernel) GetForegroundName() (int, string) {
	if c.foreground == nil {
		return adaptiveticker.UnknownId, ""
	}
	return c.foreground.PID(), c.foreground.GetCommand().Name()
}

// SetBackground sets the foreground object to nil, effectively resetting it, and returns true if it was not already nil.
func (c *Kernel) SetBackground() bool {
	if c.foreground == nil {
		return false
	}
	c.foreground = nil
	return true
}

// KillForeground terminates the process currently set as the foreground process in the kernel, if one exists.
func (c *Kernel) KillForeground() {
	if c.foreground == nil {
		return
	}
	c.Kill(c.foreground.PID())
}

// Kill terminates and removes a task by its process ID (pid). Returns true if successful, false if the pid is not found.
func (c *Kernel) Kill(pid int) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}
	task := t.(*Task)
	if len(task.timers) > 0 {
		c.ticker.Remove(task.timers)
	}
	if c.foreground != nil {
		if c.foreground.PID() == pid {
			c.foreground = nil
		}
	}
	c.ids.Unset(pid)
	return true
}

// KillAll terminates all tasks matching the specified name. Returns the number of tasks successfully terminated.
func (c *Kernel) KillAll(name string) int {
	count := 0
	var tasks []*Task

	c.ids.Range(func(item adaptiveticker.IIds) bool {
		task, ok := item.(*Task)
		if ok && task != nil {
			tasks = append(tasks, task)
		}
		return true
	})

	for _, task := range tasks {
		deactivate := false
		if len(name) == 0 {
			deactivate = true
		} else {
			if task.cmd.Name() == name {
				deactivate = true
			}
		}

		if deactivate {
			if ok := c.Kill(task.pid); ok {
				count++
			}
		}
	}
	return count
}

// List returns a formatted string containing task process IDs and their respective command names managed by the Kernel.
func (c *Kernel) List() string {
	out := "\r\nPid: Task"
	c.ids.Range(func(item adaptiveticker.IIds) bool {
		task, ok := item.(*Task)
		if ok && task != nil {
			out += fmt.Sprintf("\r\n%d: %s", task.pid, task.cmd.Name())
		}
		return true
	})
	return out
}

// ExecTimer triggers a timer event for a task identified by the given pid and tid, with the specified interval.
// Returns true if the event was successfully triggered, otherwise false.
func (c *Kernel) ExecTimer(pid int, tid int, interval int) bool {
	ret := false
	if t, ok := c.ids.Get(pid); ok {
		task := t.(*Task)
		if fn := task.cmd.TimerEvent(); fn != nil {
			fn(task, tid, interval)
			ret = true
		}
	}
	return ret
}

// ExecRead invokes the ReadEvent function for a task identified by pid. Returns true if execution is successful.
func (c *Kernel) ExecRead(pid int, code int, buffer rune) bool {
	ret := false
	if t, ok := c.ids.Get(pid); ok {
		task := t.(*Task)
		if fn := task.cmd.ReadEvent(); fn != nil {
			fn(task, code, buffer)
			ret = true
		}
	}
	return ret
}

// ExecPaint executes a rendering operation if the surface is marked as dirty, processing selected and other tasks.
// Returns true if the rendering process is executed, false otherwise.
func (c *Kernel) ExecPaint() bool {
	if !c.sh.IsDirty() {
		return false
	}
	var selectedTask interfaces.ITask = nil
	var tasks []interfaces.ITask
	c.ids.Range(func(item adaptiveticker.IIds) bool {
		task, ok := item.(*Task)
		if ok && task != nil {
			if task.PID() == c.selector.PID() {
				selectedTask = task
			} else {
				tasks = append(tasks, task)
			}
		}
		return true
	})
	return c.sh.ExecPaint(selectedTask, tasks)
}

// ListTasks retrieves a list of task names by scanning files in the current directory with a specific extension.
func (c *Kernel) ListTasks() []string {
	var out []string
	dir := "./"
	if files, err := os.ReadDir(dir); err == nil {
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			file := f.Name()
			pos := strings.LastIndex(file, tasksFileExtension)
			if pos < 0 {
				continue
			}
			out = append(out, file[:pos])
		}
	}
	return out
}

// SaveTasks saves task configurations to a JSON file, returning true if successful, otherwise false.
func (c *Kernel) SaveTasks(name string) bool {
	options := make(map[int]*TaskOptions)
	c.ids.Range(func(item adaptiveticker.IIds) bool {
		task, ok := item.(*Task)
		if ok && task != nil {
			if !strings.HasPrefix(task.Line(), commandTask) {
				options[task.pid] = task.Options()
			}
		}
		return true
	})
	data, err := json.Marshal(options)
	if err != nil {
		log.Println("Error marshalling task file ", name, ": ", err.Error())
		return false
	}
	if pos := strings.LastIndex(name, string(os.PathSeparator)); pos > -1 {
		name = name[pos+1:]
	}
	name += tasksFileExtension
	if err = os.WriteFile(name, data, 0644); err != nil {
		log.Println("Error writing task file ", name, ": ", err.Error())
		return false
	}
	return true
}

// RestoreTasks attempts to restore tasks from a file by name, executing commands and reactivating the task environment.
func (c *Kernel) RestoreTasks(name string) bool {
	var tasks map[int]*TaskOptions
	if pos := strings.LastIndex(name, string(os.PathSeparator)); pos > -1 {
		name = name[pos+1:]
	}
	name += tasksFileExtension
	data, err := os.ReadFile(name)
	if err != nil {
		return false
	}
	if err = json.Unmarshal(data, &tasks); err != nil {
		return false
	}
	for _, task := range tasks {
		if strings.HasPrefix(task.Line, commandTask) {
			continue
		}
		_, _ = c.ExecCommand(task.Line, task)
	}
	_, _ = c.ExecCommand(commandActivate, nil)
	return true
}

// closeTimer removes a timer with the specified ID from the task and ticker, returning true if the timer is successfully removed.
func (c *Kernel) closeTimer(task *Task, tid int) bool {
	ret := false
	if task != nil {
		for _, timer := range task.timers {
			if timer == tid {
				ret = c.ticker.Remove([]int{timer})
				break
			}
		}
	}
	return ret
}

// shutdown stops all processes and cleans up resources managed by the Kernel instance.
func (c *Kernel) shutdown() {
	c.KillAll("")
}

// Start initializes the kernel's event handling loop and begins processing I/O operations asynchronously.
func (c *Kernel) Start() {
	c.sh.WriteColor("Admin Console Ready", interfaces.ColorBlueDef, interfaces.ColorRedDef, interfaces.ModeNormal)
	c.sh.NextLine(true)

	d := make(chan bool)
	go func() {
		d <- true
		readBuffer := make([]byte, 1024)
		for {
			n, err := c.io.Read(readBuffer)
			if err == nil {
				if n > 0 {
					re := newMessageRead(readBuffer, n)
					re.postEvent(c.messageChan)
				}
			} else {
				qe := newMessageQuit()
				qe.postEvent(c.messageChan)
				return
			}
		}
	}()
	_ = <-d
	c.eventLoop()
}

// eventLoop is the main execution loop handling incoming messages and timers, and initiates shutdown when needed.
func (c *Kernel) eventLoop() {
	for {
		select {
		case m := <-c.messageChan:
			c.messageEventHandler(m)
		case t := <-c.timersChan:
			c.messageEventHandler(t.Event.(iMessage))
		}
		if c.exit {
			c.shutdown()
			return
		}
	}
}

// messageEventHandler processes incoming messages based on their type and performs associated actions within the kernel.
// Handles MessageTypeRead, MessageTypeTimer, MessageTypePaint, and MessageTypeQuit to execute corresponding logic.
func (c *Kernel) messageEventHandler(m iMessage) {
	if m != nil {
		switch m.getType() {
		case MessageTypeRead:
			if mm, ok := m.(*MessageRead); ok {
				c.sh.Scan(mm.data)
			}

		case MessageTypeTimer:
			if mt, ok := m.(*MessageTimer); ok {
				c.ExecTimer(mt.pid, mt.tid, mt.interval)
			}

		case MessageTypePaint:
			if _, ok := m.(*MessagePaint); ok {
				c.ExecPaint()
			}

		case MessageTypeQuit:
			if _, ok := m.(*MessageQuit); ok {
				c.exit = true
			}
		}
	}
}
