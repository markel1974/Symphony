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

// tasksFileExtension defines the file extension used for storing task-related files.
const tasksFileExtension = ".task"

// commandActivate represents the string command for activating a feature or functionality.
// commandTask represents the string command for initiating or managing a task.
const (
	commandActivate = "activate"
	commandTask     = "task"
)

// TaskManager handles task scheduling, execution, and management within a context.
type TaskManager struct {
	ctx        *Context
	ticker     *adaptiveticker.AdaptiveTicker
	foreground *Task
	selector   *TaskSelector
	dirty      bool
	width      int
	height     int
	fullPaint  bool
	timersChan chan *adaptiveticker.TimerHandler
	ids        *adaptiveticker.Ids
	interactor *CommandInteractor
}

// NewTaskManager initializes and returns a new instance of TaskManager with the provided context, ticker, timers channel, and commands.
func NewTaskManager(ctx *Context, ticker *adaptiveticker.AdaptiveTicker, timersChannel chan *adaptiveticker.TimerHandler, system []interfaces.ICommand, commands interfaces.ICommand) *TaskManager {
	t := &TaskManager{
		ctx:        ctx,
		ticker:     ticker,
		foreground: nil,
		selector:   NewTaskSelector(),
		timersChan: timersChannel,
		dirty:      false,
		fullPaint:  true,
		width:      80,
		height:     24,
		ids:        adaptiveticker.NewIds(1024),
		interactor: NewCommandInteractor(commands, system),
	}
	return t
}

// Execute parses and executes a given command line string, associating it with a task, and manages its lifecycle.
func (c *TaskManager) Execute(line string, options *TaskOptions) (bool, error) {
	cmd, args, err := c.interactor.Find(line)
	if err != nil {
		return false, fmt.Errorf("error creating task: invalid command '%s'", line)
	}
	task := NewTask(c, c.ctx, cmd, line)
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

// SetScreenSize sets the screen dimensions for the TaskManager to the specified width and height, triggering a full repaint.
func (c *TaskManager) SetScreenSize(width int, height int) {
	c.width = width
	c.height = height
	c.fullPaint = true

	//if fgPid := c.GetForegroundPid(); fgPid > unknownId {
	//	c.PaintRequest()
	//}
}

// GetScreenSize retrieves the current screen width and height as two integer values.
func (c *TaskManager) GetScreenSize() (int, int) {
	return c.width, c.height
}

// SetSelectionMode updates the current selection state based on the requested PID or defaults to the first available task.
// If no task matching the requested PID is found, the first available task is selected.
// Triggers a repaint request after updating the selection.
func (c *TaskManager) SetSelectionMode(requestedPid int) {
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

	c.PaintRequest()
}

// SetSelectionModeNext increments the selection index to the next in the available list, wrapping to the start if necessary.
func (c *TaskManager) SetSelectionModeNext() {
	if !c.selector.Next() {
		return
	}
	c.PaintRequest()
}

// SetSelectionModePrevious sets the selection mode to the previous available task in the TaskManager.
func (c *TaskManager) SetSelectionModePrevious() {
	if !c.selector.Prev() {
		return
	}
	c.PaintRequest()
}

// SetSelectionDisabled resets the selection state by clearing index, process ID, and available options, and requests a repaint.
func (c *TaskManager) SetSelectionDisabled() {
	c.selector.Clear()
	c.PaintRequest()
}

// CWDChilds retrieves the names of all child files or directories under the current working directory as a string slice.
func (c *TaskManager) CWDChilds() []string {
	var out []string
	for _, z := range c.interactor.CWD().Childs() {
		out = append(out, z.Name())
	}
	return out
}

// CWD retrieves the current working directory command interface.
func (c *TaskManager) CWD() interfaces.ICommand {
	return c.interactor.CWD()
}

// CWDGet retrieves the current working directory path as a string using the TaskManager's command interface.
func (c *TaskManager) CWDGet() string {
	return c.interactor.CWD().CommandPath()
}

// CWDPath returns the current working directory path as a slice of strings by delegating to the TaskManager's implementation.
func (c *TaskManager) CWDPath() []string {
	return c.interactor.CWD().Path()
}

// CWDSet updates the current working directory to the specified path and returns true if the operation is successful.
func (c *TaskManager) CWDSet(arg string) bool {
	return c.interactor.CWDSet(arg)
}

// Help updates the current working directory to the specified path and returns true if the operation is successful.
func (c *TaskManager) Help(arg string) (string, error) {
	return c.interactor.Help(arg)
}

// SetSelectionOptions updates selection attributes for a task identified by the current process ID based on the given option and value.
// Returns true if the operation succeeds, otherwise false.
// Valid options are 'y', 'x', and 'z' for adjusting vertical offset, horizontal offset, and scale respectively.
func (c *TaskManager) SetSelectionOptions(option rune, value float64) bool {
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
	c.fullPaint = true
	c.PaintRequest()
	return true
}

// SetCaption updates the caption of a task identified by the given pid.
// Returns true if the caption was successfully updated, false if no task with the given pid exists.
func (c *TaskManager) SetCaption(pid int, caption string) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}
	task := t.(*Task)
	task.caption = caption
	return true
}

// SetFg sets the task with the given process ID as the foreground task and returns true if successful, false otherwise.
func (c *TaskManager) SetFg(pid int) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}
	task := t.(*Task)
	c.foreground = task
	return true
}

func (c *TaskManager) GetSuggestion(in string, cursor int) (string, []string, bool) {
	prefix, suggestions, found := c.interactor.Suggestion(in, cursor)
	return prefix, suggestions, found
}

// PaintRequest marks the TaskManager as dirty and initiates a paint request if not already pending.
func (c *TaskManager) PaintRequest() bool {
	ret := false
	if !c.dirty {
		c.dirty = true
		ret = true
		c.ticker.Create(c.timersChan, newMessagePaint(), -1, -1, 1)
	}
	return ret
}

// CreateTimer initializes a timer for a task with a given process ID, start delay, interval, and repeat count.
// Returns true if the timer is successfully created; otherwise, false.
func (c *TaskManager) CreateTimer(pid int, first int, interval int, count int) bool {
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

// StopTimer stops a timer identified by the given task ID (tid) for the process identified by the given process ID (pid).
// Returns true if the timer was successfully stopped; otherwise, returns false.
func (c *TaskManager) StopTimer(pid int, tid int) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}
	task := t.(*Task)
	return c.closeTimer(task, tid)
}

// IsActive checks if a task with the specified process ID (pid) exists and is currently active in the TaskManager.
func (c *TaskManager) IsActive(pid int) bool {
	_, ret := c.ids.Get(pid)
	return ret
}

// GetForegroundPid retrieves the process ID of the currently active foreground task or returns UnknownId if none exists.
func (c *TaskManager) GetForegroundPid() int {
	if c.foreground == nil {
		return adaptiveticker.UnknownId
	}
	return c.foreground.PID()
}

// GetForegroundName returns the PID and name of the current foreground task. If none exists, it returns UnknownId and an empty string.
func (c *TaskManager) GetForegroundName() (int, string) {
	if c.foreground == nil {
		return adaptiveticker.UnknownId, ""
	}
	return c.foreground.PID(), c.foreground.GetCommand().Name()
}

// SetBackground sets the task manager to background mode by clearing the foreground context and returns true if successful.
func (c *TaskManager) SetBackground() bool {
	if c.foreground == nil {
		return false
	}
	c.foreground = nil
	return true
}

// KillForeground terminates the currently active foreground process if one exists.
func (c *TaskManager) KillForeground() {
	if c.foreground == nil {
		return
	}
	c.Kill(c.foreground.pid)
}

// Kill terminates a task with the specified pid. Returns true if the task is successfully found and removed, false otherwise.
func (c *TaskManager) Kill(pid int) bool {
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

// KillAll terminates all tasks managed by TaskManager matching the provided name, or all tasks if name is empty. Returns the count of terminated tasks.
func (c *TaskManager) KillAll(name string) int {
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

// List returns a formatted string listing all tasks with their process IDs and names.
func (c *TaskManager) List() string {
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

// ExecTimer triggers a timer event for a specific task if the task and its timer function exist, returning true if successful.
func (c *TaskManager) ExecTimer(pid int, tid int, interval int) bool {
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

// ExecRead triggers a read event for a task identified by pid with the provided code and buffer values.
// Returns true if the event execution is successful; otherwise, false.
func (c *TaskManager) ExecRead(pid int, code int, buffer rune) bool {
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

// ExecPaint renders the current task manager state onto the terminal by painting tasks and handling selection logic.
func (c *TaskManager) ExecPaint(terminal interfaces.ITerminal) bool {
	if !c.dirty {
		return false
	}

	w, h := c.GetScreenSize()
	surface := newSurface(terminal, h, w)

	if c.fullPaint {
		surface.SetCompletePaint()
		c.fullPaint = false
	}

	var selectedTask *Task = nil

	c.ids.Range(func(item adaptiveticker.IIds) bool {
		task, ok := item.(*Task)
		if ok && task != nil {
			if task.PID() == c.selector.PID() {
				selectedTask = task
			} else {
				surface.SetSelectionMode(false)
				task.Paint(surface)
			}
		}
		return true
	})

	if selectedTask != nil {
		surface.SetSelectionMode(true)
		selectedTask.Paint(surface)
	}

	surface.Render()

	c.dirty = false

	return true
}

// ListTasks retrieves the names of all tasks available in the current directory with the specified file extension.
func (c *TaskManager) ListTasks() []string {
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

// SaveTasks saves the current state of tasks to a file with the specified name and returns true if the operation succeeds.
func (c *TaskManager) SaveTasks(name string) bool {
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

// RestoreTasks restores tasks from a file identified by the given name, reinitializing their state. Returns success status.
func (c *TaskManager) RestoreTasks(name string) bool {
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
		_, _ = c.Execute(task.Line, task)
	}
	_, _ = c.Execute(commandActivate, nil)
	return true
}

// ExecActivate attempts to activate a background process if it is not already active or invalid. Returns false in all cases.
func (c *TaskManager) ExecActivate() bool {
	pid, name := c.GetForegroundName()
	if pid == adaptiveticker.UnknownId {
		return false
	}
	if name == commandActivate {
		return false
	}
	c.SetBackground()
	_, _ = c.Execute(fmt.Sprint(commandActivate, " ", pid), nil)
	return false
}

// closeTimer removes a timer identified by tid from the task's timer list and returns true if the operation is successful.
func (c *TaskManager) closeTimer(task *Task, tid int) bool {
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
