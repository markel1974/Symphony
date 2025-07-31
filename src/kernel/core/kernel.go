package core

import (
	"encoding/json"
	"fmt"
	"github.com/markel1974/c64emu/src/kernel/adaptiveticker"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/messages"
	"github.com/markel1974/c64emu/src/kernel/servers/shell"
	"io"
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
	ticker       *adaptiveticker.AdaptiveTicker
	inputDriver  io.Reader
	outputDriver io.Writer
	render       interfaces.IRender
	foreground   interfaces.ITask
	selector     *TaskSelector
	ids          *adaptiveticker.Ids
	fs           interfaces.IFileSystem
	sh           *shell.Shell
	messageChan  chan messages.IMessage
	timersChan   chan *adaptiveticker.TimerHandler
	exit         bool
}

// NewKernel creates and returns a new Kernel instance, initializing its dependencies and internal fields.
func NewKernel(ticker *adaptiveticker.AdaptiveTicker, timersChan chan *adaptiveticker.TimerHandler, inputDriver io.Reader, outputDriver io.Writer, render interfaces.IRender, fs interfaces.IFileSystem, sh *shell.Shell) *Kernel {
	t := &Kernel{
		ticker:       ticker,
		inputDriver:  inputDriver,
		outputDriver: outputDriver,
		render:       render,
		fs:           fs,
		sh:           sh,
		foreground:   nil,
		selector:     NewTaskSelector(),
		ids:          adaptiveticker.NewIds(1024),
		messageChan:  make(chan messages.IMessage, contextMaQueueLen),
		timersChan:   timersChan,
		exit:         false,
	}
	return t
}

// Start initializes the kernel's event handling loop and begins processing I/O operations asynchronously.
func (c *Kernel) Start() {
	c.render.WriteHighlight("Admin Console Ready")
	c.sh.SetPromptPrefix(c.fs.CWD().Name())
	c.sh.NextLine(true)

	d := make(chan bool)
	go func() {
		d <- true
		readBuffer := make([]byte, 4096)
		for {
			n, err := c.inputDriver.Read(readBuffer)
			if err == nil {
				if n > 0 {
					re := messages.NewMessageRead(readBuffer, n)
					c.messageChan <- re
				}
			} else {
				qe := messages.NewMessageQuit()
				c.messageChan <- qe
				return
			}
		}
	}()
	_ = <-d
	c.eventLoop()
}

// IOWrite writes the provided byte slice to the output driver and returns the number of bytes written and any error encountered.
func (c *Kernel) IOWrite(data []byte) (int, error) {
	return c.outputDriver.Write(data)
}

// IORead reads data from the input driver into the provided byte slice and returns the number of bytes read and any error encountered.
func (c *Kernel) IORead(p []byte) (int, error) {
	return c.inputDriver.Read(p)
}

// IOType processes input events based on their type and key value to handle control, foreground tasks, and system state.
func (c *Kernel) IOType(kind interfaces.KeyType, key rune) {
	if kind == interfaces.KeyTypeCtrl {
		switch key {
		case 3:
			c.selector.Clear()
			c.killForeground()
			c.sh.NextLine(true)
		case 4:
			c.handleExecActivate()
		}
		return
	}
	if fgPid := c.getForegroundPid(); fgPid != adaptiveticker.UnknownId {
		c.handleTaskKeyEvent(fgPid, int(kind), key)
		return
	}
	if quit := c.handleKeyEvent(kind, key); quit {
		c.ExitRequested()
	}
}

// SetScreenSize sets the display dimensions of the screen to the specified width (w) and height (h).
func (c *Kernel) SetScreenSize(w int, h int) {
	c.render.SetScreenSize(w, h)
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
}

// SetSelectionTaskNext advances to the next available task in the selection and triggers a repaint request if successful.
func (c *Kernel) SetSelectionTaskNext() bool {
	return c.selector.Next()
}

// SetSelectionTaskPrevious updates selection to the previous item using the internal selector and requests a repaint if successful.
func (c *Kernel) SetSelectionTaskPrevious() bool {
	return c.selector.Prev()
}

// ExitRequested signals that the kernel should terminate its execution and exit.
func (c *Kernel) ExitRequested() {
	c.exit = true
}

// SendSelectionTaskOptions modifies selection parameters like offsets or scale based on the provided option and value.
// Returns true if the operation is applied successfully; otherwise, returns false.
func (c *Kernel) SendSelectionTaskOptions(option rune, value float64) bool {
	t, ok := c.ids.Get(c.selector.PID())
	if !ok {
		return false
	}
	task := t.(*Task)
	task.SetOption(option, value)
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
	m := messages.NewMessageTimer(pid, interval)
	m.SetTID(c.ticker.Create(c.timersChan, m, int64(first), int64(interval), int64(count)))
	if m.TID() > -1 {
		task.timers = append(task.timers, m.TID())
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

// getForegroundPid retrieves the process ID of the currently active foreground process, or UnknownId if none is active.
func (c *Kernel) getForegroundPid() int {
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

// killForeground terminates the process currently set as the foreground process in the kernel, if one exists.
func (c *Kernel) killForeground() {
	if c.foreground == nil {
		return
	}
	c.Kill(c.foreground.PID())
}

// ExecCommand executes a command by parsing the input line, creating a task, and managing its lifecycle and state.
// Returns true and error if the task was created but execution failed, or true and nil if execution succeeded.
func (c *Kernel) ExecCommand(line string, options *TaskOptions) (bool, error) {
	cmd, args, err := c.fs.Find(line)
	if err != nil {
		return false, fmt.Errorf("error creating task: invalid command '%s'", line)
	}
	task := NewTask(c, c.fs, c.render, c.sh, cmd, line)
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

// handleExecActivate attempts to activate a foreground process by executing the associated command and returns its success status.
func (c *Kernel) handleExecActivate() bool {
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

// eventLoop is the main execution loop handling incoming messages and timers, and initiates shutdown when needed.
func (c *Kernel) eventLoop() {
	for {
		select {
		case m := <-c.messageChan:
			c.handleMessageEvent(m)
		case t := <-c.timersChan:
			c.handleMessageEvent(t.Event.(messages.IMessage))
		}
		if c.exit {
			c.shutdown()
			return
		}
	}
}

// handleMessageEvent processes incoming messages based on their type and performs associated actions within the kernel.
// Handles MessageTypeRead, MessageTypeTimer, MessageTypePaint, and MessageTypeQuit to execute corresponding logic.
func (c *Kernel) handleMessageEvent(m messages.IMessage) {
	if m != nil {
		switch m.GetType() {
		case messages.MessageTypeRead:
			if mm, ok := m.(*messages.MessageRead); ok {
				c.render.Scan(mm.Data())
			}
		case messages.MessageTypeTimer:
			if mt, ok := m.(*messages.MessageTimer); ok {
				c.handleTimerEvent(mt.PID(), mt.TID(), mt.Interval())
			}
		case messages.MessageTypePaint:
			if _, ok := m.(*messages.MessagePaint); ok {
				c.handlePaintEvent()
			}
		case messages.MessageTypeQuit:
			if _, ok := m.(*messages.MessageQuit); ok {
				c.exit = true
			}
		}
	}
}

// ExecRead invokes the ReadEvent function for a task identified by pid. Returns true if execution is successful.
func (c *Kernel) handleTaskKeyEvent(pid int, code int, buffer rune) bool {
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

// handleKeyEvent processes a keyboard event of a given type and key and returns true if the event was handled successfully.
func (c *Kernel) handleKeyEvent(kind interfaces.KeyType, key rune) bool {
	c.sh.KeyHandler(kind, key)
	if !c.sh.Authenticated() {
		return true
	}
	if kind == interfaces.KeyTypeEnter {
		if buffer := c.sh.InputBuffer(); len(buffer) > 0 {
			c.render.WriteLn("")
			_, _ = c.ExecCommand(buffer, nil)
			c.sh.NextLine(false)
		} else {
			c.sh.NextLine(true)
		}
	} else if kind == interfaces.KeyTypeTab {
		if c.sh.TabFound() {
			tabData, cursor, tabCount := c.sh.TabData()
			data, suggestions, found := c.fs.Suggestion(tabData, cursor)
			if found && len(suggestions) > 0 {
				sLen := len(suggestions)
				if idx := tabCount % sLen; idx < sLen {
					if complete := suggestions[idx]; len(complete) > len(data) {
						tabLine := complete
						c.sh.Redraw(tabLine)
						c.sh.SetHistoryDefault(tabLine)
						if sLen == 1 {
							c.sh.TabReset()
						}
					}
				}
			}
		}
	}
	return false
}

// ExecPaint executes a rendering operation if the surface is marked as dirty, processing selected and other tasks.
// Returns true if the rendering process is executed, false otherwise.
func (c *Kernel) handlePaintEvent() bool {
	if !c.render.IsDirty() {
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
	return c.render.ExecPaint(selectedTask, tasks)
}

// ExecTimer triggers a timer event for a task identified by the given pid and tid, with the specified interval.
// Returns true if the event was successfully triggered, otherwise false.
func (c *Kernel) handleTimerEvent(pid int, tid int, interval int) bool {
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
