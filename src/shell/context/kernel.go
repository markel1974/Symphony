/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package context

import (
	"encoding/json"
	"fmt"
	"github.com/markel1974/c64emu/src/shell/adaptiveticker"
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"log"
	"os"
	"strings"
)

// tasksFileExtension defines the standard file extension used for task-related files in the application.
const tasksFileExtension = ".task"

// commandActivate represents the string command for "activate" functionality.
// commandTask represents the string command for "task" functionality.
const (
	commandActivate = "activate"
	commandTask     = "task"
)

// Kernel represents the core scheduler and manager for tasks and system resources, facilitating interaction and control.
type Kernel struct {
	ticker      *adaptiveticker.AdaptiveTicker
	render      interfaces.IRender
	io          interfaces.IInputOutput
	foreground  interfaces.ITask
	selector    *TaskSelector
	timersChan  chan *adaptiveticker.TimerHandler
	ids         *adaptiveticker.Ids
	fs          interfaces.IFileSystem
	messageChan chan iMessage
	exit        bool
}

// NewKernel initializes and returns a new Kernel instance with provided ticker, render, I/O, and filesystem components.
func NewKernel(ticker *adaptiveticker.AdaptiveTicker, render interfaces.IRender, io interfaces.IInputOutput, fs interfaces.IFileSystem) *Kernel {
	t := &Kernel{
		ticker:      ticker,
		render:      render,
		io:          io,
		fs:          fs,
		foreground:  nil,
		selector:    NewTaskSelector(),
		ids:         adaptiveticker.NewIds(1024),
		messageChan: make(chan iMessage, contextMaQueueLen),
		timersChan:  make(chan *adaptiveticker.TimerHandler, contextMaQueueLen),
		exit:        false,
	}
	return t
}

// SetScreenSize resizes the screen dimensions by setting the width and height parameters for the rendering engine.
func (c *Kernel) SetScreenSize(w int, h int) {
	c.render.SetScreenSize(w, h)
}

// GetScreenSize returns the screen's width and height as integers.
func (c *Kernel) GetScreenSize() (int, int) {
	return c.render.GetScreenSize()
}

// Write writes the provided data string to the kernel's render for processing or output.
func (c *Kernel) Write(data string) {
	c.render.Write(data)
}

// WriteLn writes a string followed by a newline to the kernel's output stream.
func (c *Kernel) WriteLn(data string) {
	c.render.WriteLn(data)
}

// WriteColor outputs text with specified foreground and background colors and applies the given color rendering mode.
func (c *Kernel) WriteColor(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	c.render.WriteColor(data, fg, bg, mode)
}

// WriteColorLn writes a line of text with specified foreground color, background color, and color mode to the output.
func (c *Kernel) WriteColorLn(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	c.render.WriteColorLn(data, fg, bg, mode)
}

// ClearScreen clears the output screen by invoking the render's ClearScreen method.
func (c *Kernel) ClearScreen() {
	c.render.ClearScreen()
}

// ExecActivate attempts to activate a foreground process by sending an activate command to its associated identifier.
// Returns false if the process ID or name is invalid, or if the process is already in an active state.
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

// ExecCommand executes a command from the given input line and applies the provided task options for configuration.
// Returns a boolean indicating task creation success and an error if the command execution fails.
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

// SetSelectionMode updates the task selection based on a requested process ID and resets selection if needed.
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

// PaintRequest determines whether a paint operation is required by invoking the underlying kernel logic.
func (c *Kernel) PaintRequest() bool {
	return c.doPaintRequest(false)
}

// doPaintRequest attempts to initiate a paint request. If successful, it triggers a timer for repaint and returns true.
func (c *Kernel) doPaintRequest(full bool) bool {
	if c.render.PaintRequest(full) {
		c.ticker.Create(c.timersChan, newMessagePaint(), -1, -1, 1)
		return true
	}
	return false
}

// SetSelectionModeNext advances the task selection to the next available task and triggers a paint request if successful.
func (c *Kernel) SetSelectionModeNext() {
	if !c.selector.Next() {
		return
	}
	c.doPaintRequest(false)
}

// SetSelectionModePrevious switches the selection mode to the previous item in the task selector and triggers a repaint.
func (c *Kernel) SetSelectionModePrevious() {
	if !c.selector.Prev() {
		return
	}
	c.doPaintRequest(false)
}

// SetSelectionDisabled clears the task selector state and disables ongoing selection, triggering a paint request.
func (c *Kernel) SetSelectionDisabled() {
	c.selector.Clear()
	c.doPaintRequest(false)
}

// ExitRequested sets the kernel's exit flag to true, signaling a request to terminate its execution.
func (c *Kernel) ExitRequested() {
	c.exit = true
}

// CWDChilds retrieves the names of child items in the current working directory of the kernel's filesystem.
func (c *Kernel) CWDChilds() []string {
	var out []string
	for _, z := range c.fs.CWD().Childs() {
		out = append(out, z.Name())
	}
	return out
}

// CWD retrieves the current working directory command associated with the kernel, returning an instance of ICommand.
func (c *Kernel) CWD() interfaces.ICommand {
	return c.fs.CWD()
}

// CWDGet returns the current working directory path as a string from the kernel's filesystem.
func (c *Kernel) CWDGet() string {
	return c.fs.CWD().CommandPath()
}

// CWDPath returns the current working directory path as a slice of strings.
func (c *Kernel) CWDPath() []string {
	return c.fs.CWD().Path()
}

// CWDSet sets the current working directory for the kernel to the provided path and returns true if successful.
func (c *Kernel) CWDSet(arg string) bool {
	return c.fs.CWDSet(arg)
}

// Help retrieves help text associated with the provided argument and returns it along with any potential error.
func (c *Kernel) Help(arg string) (string, error) {
	return c.fs.Help(arg)
}

// SetSelectionOptions adjusts specific task attributes based on the provided option and value, triggering a paint update.
// Returns true if successful; otherwise, false if the task is not found.
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

// SetFg sets the task with the specified PID as the foreground task if it exists and returns true; otherwise, returns false.
func (c *Kernel) SetFg(pid int) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}
	task := t.(*Task)
	c.foreground = task
	return true
}

// GetSuggestion generates command-line suggestions based on the input and current cursor position.
// It returns the shared prefix, a list of suggestions, and a boolean indicating if suggestions were found.
func (c *Kernel) GetSuggestion(in string, cursor int) (string, []string, bool) {
	prefix, suggestions, found := c.fs.Suggestion(in, cursor)
	return prefix, suggestions, found
}

// CreateTimer initializes a timer for a specific process and adds it to the task's timer list.
// Returns true if the timer is successfully created; otherwise, returns false.
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

// StopTimer stops the timer associated with the given process ID (pid) and timer ID (tid). Returns true if successful.
func (c *Kernel) StopTimer(pid int, tid int) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}
	task := t.(*Task)
	return c.closeTimer(task, tid)
}

// IsActive checks if a process with the given pid is currently active in the kernel and returns true if so.
func (c *Kernel) IsActive(pid int) bool {
	_, ret := c.ids.Get(pid)
	return ret
}

// GetForegroundPid returns the process ID (PID) of the current foreground process or a default invalid ID if none exists.
func (c *Kernel) GetForegroundPid() int {
	if c.foreground == nil {
		return adaptiveticker.UnknownId
	}
	return c.foreground.PID()
}

// GetForegroundName returns the PID and command name of the currently foregrounded process, or a default if none exists.
func (c *Kernel) GetForegroundName() (int, string) {
	if c.foreground == nil {
		return adaptiveticker.UnknownId, ""
	}
	return c.foreground.PID(), c.foreground.GetCommand().Name()
}

// SetBackground resets the foreground process of the kernel to nil.
// Returns true if the operation is successful; false otherwise.
func (c *Kernel) SetBackground() bool {
	if c.foreground == nil {
		return false
	}
	c.foreground = nil
	return true
}

// KillForeground terminates the process currently running in the foreground if one exists.
func (c *Kernel) KillForeground() {
	if c.foreground == nil {
		return
	}
	c.Kill(c.foreground.PID())
}

// Kill removes the task identified by the given PID from the kernel and deallocates its resources. Returns true if successful.
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

// KillAll terminates all tasks with a specified name or all tasks if the name is empty, returning the count of terminated tasks.
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

// List returns a formatted string that lists all tasks, including their process ID and command name.
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

// ExecTimer triggers a timer event for a task identified by pid and tid with the specified interval.
// It returns true if the event is successfully executed; otherwise, it returns false.
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

// ExecRead attempts to execute a read operation for the specified process ID using the given code and buffer values.
// It fetches the task by pid, invokes its registered read event if available, and returns true upon success.
// Returns false if the process ID is not found or no read event is registered.
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

// ExecPaint processes and renders tasks, prioritizing the task matching the current selector's PID, if any.
func (c *Kernel) ExecPaint() bool {
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

// ListTasks retrieves a list of task names by scanning files in the current directory with the specific task file extension.
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

// SaveTasks saves the current task options to a file with the specified name in JSON format. Returns true if successful.
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

// RestoreTasks restores a set of tasks from a file identified by name, initializing and activating them. Returns true on success.
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

// History performs a history operation using the specified action and index in the kernel context.
func (c *Kernel) History(verb interfaces.HistoryAction, idx int) {
	//TODO IMPLEMENT
	c.io.History(verb, idx)
}

// closeTimer removes a timer identified by tid from the task's list of timers and returns true if successful, false otherwise.
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

// shutdown terminates all running processes associated with the Kernel by invoking KillAll with an empty string argument.
func (c *Kernel) shutdown() {
	c.KillAll("")
}

// Start initializes and begins the main event loop of the Kernel, handling I/O operations and event dispatching.
func (c *Kernel) Start() {
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

// eventLoop continuously processes messages and timers, delegating events to the appropriate handler until exit is triggered.
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

// messageEventHandler processes different message types and performs appropriate actions based on the message type.
func (c *Kernel) messageEventHandler(m iMessage) {
	if m != nil {
		switch m.getType() {
		case MessageTypeRead:
			if mm, ok := m.(*MessageRead); ok {
				c.render.Scan(mm.data)
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
