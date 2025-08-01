package core

import (
	"encoding/json"
	"fmt"
	"github.com/markel1974/c64emu/src/kernel/adaptiveticker"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/messages"
	"github.com/markel1974/c64emu/src/kernel/process_factory"
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
	foreground   interfaces.IProcess
	selector     *ProcessSelector
	ids          *adaptiveticker.Ids
	fs           interfaces.IFileSystem
	sh           *shell.Shell
	messageChan  chan messages.IMessage
	pf           *process_factory.ProcessFactory
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
		selector:     NewProcessSelector(),
		ids:          adaptiveticker.NewIds(1024),
		messageChan:  make(chan messages.IMessage, contextMaQueueLen),
		timersChan:   timersChan,
		exit:         false,
	}
	t.pf = process_factory.NewProcessFactory(t)
	return t
}

// CallTaskExec executes a command by parsing the input line, creating a task, and managing its lifecycle and state.
// Returns true and error if the task was created but execution failed, or true and nil if execution succeeded.
func (c *Kernel) CallTaskExec(line string, options *interfaces.ProcessOptions) (bool, error) {
	cmd, args, err := c.fs.Find(line)
	if err != nil {
		return false, fmt.Errorf("error creating task: invalid command '%s'", line)
	}
	task := c.pf.Create(cmd, line, options)
	if !c.ids.Set(task) {
		return false, fmt.Errorf("error creating task: can't set pid")
	}
	if err = cmd.Execute(task, args); err != nil {
		c.CallTaskKill(task.PID())
		return true, err
	}
	if !cmd.Daemon() {
		c.CallTaskKill(task.PID())
		return true, nil
	}
	task.SetState(interfaces.ProcessStateRunning)
	if !cmd.Background() {
		c.foreground = task
	}
	return true, nil
}

// CallTaskKill terminates and removes a task by its process ID (pid). Returns true if successful, false if the pid is not found.
func (c *Kernel) CallTaskKill(pid int) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}
	task, ok := t.(interfaces.IProcess)
	if !ok {
		return false
	}
	if len(task.Timers()) > 0 {
		c.ticker.RemoveEntries(task.Timers())
	}
	if c.foreground != nil {
		if c.foreground.PID() == pid {
			c.foreground = nil
		}
	}
	c.ids.Unset(pid)
	return true
}

// CallTaskKillAll terminates all tasks matching the specified name. Returns the number of tasks successfully terminated.
func (c *Kernel) CallTaskKillAll(name string) int {
	count := 0
	var tasks []interfaces.IProcess
	c.ids.Range(func(item adaptiveticker.IIds) bool {
		task, ok := item.(interfaces.IProcess)
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
			if task.GetCommand().Name() == name {
				deactivate = true
			}
		}
		if deactivate {
			if ok := c.CallTaskKill(task.PID()); ok {
				count++
			}
		}
	}
	return count
}

// CallTaskGetForegroundName retrieves the process ID and command name of the current foreground process.
// Returns a known invalid ID and empty string if no foreground process is set.
func (c *Kernel) CallTaskGetForegroundName() (int, string) {
	if c.foreground == nil {
		return adaptiveticker.UnknownId, ""
	}
	return c.foreground.PID(), c.foreground.GetCommand().Name()
}

// CallTaskSetBackground sets the foreground object to nil, effectively resetting it, and returns true if it was not already nil.
func (c *Kernel) CallTaskSetBackground() bool {
	if c.foreground == nil {
		return false
	}
	c.foreground = nil
	return true
}

// CallTaskKillForeground terminates the process currently set as the foreground process in the kernel, if one exists.
func (c *Kernel) CallTaskKillForeground() {
	if c.foreground == nil {
		return
	}
	c.CallTaskKill(c.foreground.PID())
}

// CallTaskSaveAll saves task configurations to a JSON file, returning true if successful, otherwise false.
func (c *Kernel) CallTaskSaveAll(name string) bool {
	options := make(map[int]*interfaces.ProcessOptions)
	c.ids.Range(func(item adaptiveticker.IIds) bool {
		task, ok := item.(interfaces.IProcess)
		if ok && task != nil {
			if !strings.HasPrefix(task.Line(), commandTask) {
				options[task.PID()] = task.Options()
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

// CallTaskRestoreAll attempts to restore tasks from a file by name, executing commands and reactivating the task environment.
func (c *Kernel) CallTaskRestoreAll(name string) bool {
	var tasks map[int]*interfaces.ProcessOptions
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
		_, _ = c.CallTaskExec(task.Line, task)
	}
	_, _ = c.CallTaskExec(commandActivate, nil)
	return true
}

// CallTaskList returns a formatted string containing task process IDs and their respective command names managed by the Kernel.
func (c *Kernel) CallTaskList() string {
	out := "\r\nPid: Process"
	c.ids.Range(func(item adaptiveticker.IIds) bool {
		task, ok := item.(interfaces.IProcess)
		if ok && task != nil {
			out += fmt.Sprintf("\r\n%d: %s", task.PID(), task.GetCommand().Name())
		}
		return true
	})
	return out
}

// CallTaskSavedList retrieves a list of task names by scanning files in the current directory with a specific extension.
func (c *Kernel) CallTaskSavedList() []string {
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

// CallTaskSelection updates the selection mode for a specific process and triggers a repaint without requesting a redraw.
func (c *Kernel) CallTaskSelection(pid int) {
	c.setSelectionMode(pid)
	c.render.PaintRequest(false)
}

// CallTaskSelectionPrevious moves the task selection to the previous task and triggers a render update if successful.
func (c *Kernel) CallTaskSelectionPrevious() {
	if c.selector.Prev() {
		c.render.PaintRequest(false)
	}
}

// CallTaskSelectionNext advances the task selector to the next task and triggers a repaint if the task selection changes.
func (c *Kernel) CallTaskSelectionNext() {
	if c.selector.Next() {
		c.render.PaintRequest(false)
	}
}

// CallTaskSelectionOptions updates the selected task's option with the given rune and value, then triggers a repaint request.
// Returns true on successful task retrieval and option update, otherwise returns false.
func (c *Kernel) CallTaskSelectionOptions(option rune, value float64) bool {
	t, ok := c.ids.Get(c.selector.PID())
	if !ok {
		return false
	}
	task, ok := t.(interfaces.IProcess)
	if !ok {
		return false
	}
	task.SetOption(option, value)
	c.render.PaintRequest(true)
	return true
}

// CallPaintRequest triggers a paint request via the render component and returns true if successful.
func (c *Kernel) CallPaintRequest() bool {
	return c.render.PaintRequest(false)
}

// CallWrite sends the provided string data to the kernel's rendering writer for output.
func (c *Kernel) CallWrite(data string) {
	c.render.Write(data)
}

// CallWriteLn writes the provided string followed by a new line to the kernel's output stream.
func (c *Kernel) CallWriteLn(data string) {
	c.render.WriteLn(data)
}

// CallWriteColor writes a string to the output with specified foreground color, background color, and color mode.
func (c *Kernel) CallWriteColor(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	c.render.WriteColor(data, fg, bg, mode)
}

// CallWriteColorLn writes a line of text with specified foreground and background colors and a given color mode.
func (c *Kernel) CallWriteColorLn(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	c.render.WriteColorLn(data, fg, bg, mode)
}

// CallClearScreen clears the screen by invoking the associated renderer's ClearScreen method.
func (c *Kernel) CallClearScreen() {
	c.render.ClearScreen()
}

// CallScreenSize retrieves the screen's width and height as integers from the render instance.
func (c *Kernel) CallScreenSize() (int, int) {
	return c.render.GetScreenSize()
}

// CallCWDSet sets the current working directory to the specified path and updates the shell prompt accordingly.
func (c *Kernel) CallCWDSet(arg string) bool {
	if ok := c.fs.CWDSet(arg); ok {
		c.sh.SetPromptPrefix(c.fs.CWDName())
		return true
	}
	return false
}

// CallCWDGet returns the command path of the current working directory from the file system.
func (c *Kernel) CallCWDGet() string {
	return c.fs.CWDCommandPath()
}

// CallCWDPath retrieves the current working directory's path as a slice of strings from the filesystem instance.
func (c *Kernel) CallCWDPath() []string {
	return c.fs.CWDPath()
}

// CallCWDDirectoryListing retrieves the directory listing of the current working directory as a slice of strings.
func (c *Kernel) CallCWDDirectoryListing() []string {
	return c.fs.CWDDirectoryListing()
}

// CallHistory applies a history action to the shell and invokes task execution if arguments are produced.
func (c *Kernel) CallHistory(verb interfaces.HistoryAction, idx int) {
	if arg := c.sh.HistoryApply(verb, idx); len(arg) > 0 {
		_, _ = c.CallTaskExec(arg, nil)
	}
}

// CallHelp retrieves the help information associated with the given argument and returns it as a string.
// Returns an error if the help information cannot be fetched.
func (c *Kernel) CallHelp(arg string) (string, error) {
	return c.fs.Help(arg)
}

// CallSetScreenSize adjusts the screen dimensions to the specified width and height values.
func (c *Kernel) CallSetScreenSize(w int, h int) {
	c.render.SetScreenSize(w, h)
}

// CallExitRequested sets the `exit` flag to true, signaling that an exit has been requested for the kernel.
func (c *Kernel) CallExitRequested() {
	c.exit = true
}

// CallSetFg sets the foreground task to the one associated with the given PID. Returns true if successful, false otherwise.
func (c *Kernel) CallSetFg(pid int) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}
	task, ok := t.(interfaces.IProcess)
	if !ok {
		return false
	}
	c.foreground = task
	return true
}

// CallCreateTimer creates a timer for the specified process ID with given start time, interval, and execution count.
// Returns true if the timer is successfully created, otherwise false.
// It requires the process to have a valid TimerEvent handler.
// Adds the timer ID to the process's list of active timers.
func (c *Kernel) CallCreateTimer(pid int, first int, interval int, count int) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}
	task, ok := t.(interfaces.IProcess)
	if !ok {
		return false
	}
	if task.GetCommand().TimerEvent == nil {
		return false
	}
	m := messages.NewMessageTimer(pid, interval)
	m.SetTID(c.ticker.Create(c.timersChan, m, int64(first), int64(interval), int64(count)))
	if m.TID() > -1 {
		task.AddTimer(m.TID())
	}
	return true
}

func (c *Kernel) CallStopTimer(pid int, tid int) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}
	task, ok := t.(interfaces.IProcess)
	if !ok {
		return false
	}
	return c.closeTimer(task, tid)
}

func (c *Kernel) CallIsActive(pid int) bool {
	_, ret := c.ids.Get(pid)
	return ret
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
			c.CallTaskKillForeground()
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
		c.CallExitRequested()
	}
}

// Start initializes the kernel's event handling loop and begins processing I/O operations asynchronously.
func (c *Kernel) Start() {
	c.render.WriteHighlight("Admin Console Ready")
	c.sh.SetPromptPrefix(c.fs.CWDName())
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

// SetSelectionMode sets the current selection mode for tasks based on a requested process ID. Defaults to the first task if the requested ID is unavailable.
func (c *Kernel) setSelectionMode(requestedPid int) {
	var idx = 0
	var firstPid = adaptiveticker.UnknownId
	var firstIdx = 0

	c.selector.Clear()

	c.ids.Range(func(item adaptiveticker.IIds) bool {
		task, ok := item.(interfaces.IProcess)
		if ok && task != nil {
			if task.GetCommand().PaintEvent() != nil {
				c.selector.AddAvailable(task.PID())
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

// getForegroundPid retrieves the process ID of the currently active foreground process, or UnknownId if none is active.
func (c *Kernel) getForegroundPid() int {
	if c.foreground == nil {
		return adaptiveticker.UnknownId
	}
	return c.foreground.PID()
}

// handleExecActivate attempts to activate a foreground process by executing the associated command and returns its success status.
func (c *Kernel) handleExecActivate() bool {
	pid, name := c.CallTaskGetForegroundName()
	if pid == adaptiveticker.UnknownId {
		return false
	}
	if name == commandActivate {
		return false
	}
	c.CallTaskSetBackground()
	_, _ = c.CallTaskExec(fmt.Sprint(commandActivate, " ", pid), nil)
	return false
}

// closeTimer removes a timer with the specified ID from the task and ticker, returning true if the timer is successfully removed.
func (c *Kernel) closeTimer(task interfaces.IProcess, tid int) bool {
	ret := false
	if task != nil {
		task.TimersIterator(func(timerId int) bool {
			if timerId == tid {
				ret = c.ticker.RemoveEntries([]int{timerId})
				return true
			}
			return false
		})
	}
	return ret
}

// shutdown stops all processes and cleans up resources managed by the Kernel instance.
func (c *Kernel) shutdown() {
	c.CallTaskKillAll("")
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
		task := t.(interfaces.IProcess)
		if fn := task.GetCommand().ReadEvent(); fn != nil {
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
			_, _ = c.CallTaskExec(buffer, nil)
			c.sh.NextLine(false)
		} else {
			c.sh.NextLine(true)
		}
	} else if kind == interfaces.KeyTypeTab {
		if tabData, cursor, ok := c.sh.TabData(); ok {
			data, suggestions, found := c.fs.Suggestion(tabData, cursor)
			c.sh.HistorySuggest(data, suggestions, found)
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
	var selectedTask interfaces.IProcess = nil
	var tasks []interfaces.IProcess
	c.ids.Range(func(item adaptiveticker.IIds) bool {
		task, ok := item.(interfaces.IProcess)
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
		task, ok := t.(interfaces.IProcess)
		if ok && task != nil {
			if fn := task.GetCommand().TimerEvent(); fn != nil {
				fn(task, tid, interval)
				ret = true
			}
		}
	}
	return ret
}
