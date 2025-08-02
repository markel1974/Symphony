package core

import (
	"fmt"
	"github.com/markel1974/c64emu/src/kernel/adaptiveticker"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/messages"
	"github.com/markel1974/c64emu/src/kernel/process_factory"
	"io"
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
	shellPath    string
	shell        interfaces.IProcess
	messageChan  chan messages.IMessage
	pf           *process_factory.ProcessFactory
	timersChan   chan *adaptiveticker.TimerHandler
	exit         bool
}

// NewKernel creates and returns a new Kernel instance, initializing its dependencies and internal fields.
func NewKernel(ticker *adaptiveticker.AdaptiveTicker, timersChan chan *adaptiveticker.TimerHandler, inputDriver io.Reader, outputDriver io.Writer, render interfaces.IRender, fs interfaces.IFileSystem, shellPath string) *Kernel {
	t := &Kernel{
		ticker:       ticker,
		inputDriver:  inputDriver,
		outputDriver: outputDriver,
		render:       render,
		fs:           fs,
		foreground:   nil,
		selector:     NewProcessSelector(),
		ids:          adaptiveticker.NewIds(1024),
		messageChan:  make(chan messages.IMessage, contextMaQueueLen),
		timersChan:   timersChan,
		exit:         false,
		shellPath:    shellPath,
	}
	t.pf = process_factory.NewProcessFactory(t)
	return t
}

// CallProcessList returns a formatted string containing process IDs and their respective command names managed by the Kernel.
func (c *Kernel) CallProcessList() []*interfaces.ProcessDescription {
	return c.doProcessList()
}

// CallProcessExec executes a command by parsing the input line, creating a process, and managing its lifecycle and state.
// Returns true and error if the process was created but execution failed, or true and nil if execution succeeded.
func (c *Kernel) CallProcessExec(line string, options *interfaces.ProcessOptions) (bool, error) {
	_, err := c.doProcessExec(line, options)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CallProcessKill terminates and removes a task by its process ID (pid). Returns true if successful, false if the pid is not found.
func (c *Kernel) CallProcessKill(pid int) bool {
	return c.doProcessKill(pid)
}

// CallProcessKillAll terminates all tasks matching the specified name. Returns the number of tasks successfully terminated.
func (c *Kernel) CallProcessKillAll(name string) int {
	return c.doProcessKillAll(name)
}

// CallProcessSetBackground sets the foreground object to nil, effectively resetting it, and returns true if it was not already nil.
func (c *Kernel) CallProcessSetBackground() bool {
	return c.doProcessSetBackground()
}

// CallProcessSelection updates the selection mode for a specific process and triggers a repaint without requesting a redraw.
func (c *Kernel) CallProcessSelection(pid int) {
	c.doProcessSetSelectionMode(pid)
	c.render.PaintRequest(false)
}

// CallProcessSelectionPrevious moves the task selection to the previous task and triggers a render update if successful.
func (c *Kernel) CallProcessSelectionPrevious() {
	if c.selector.Prev() {
		c.render.PaintRequest(false)
	}
}

// CallProcessSelectionNext advances the task selector to the next task and triggers a repaint if the task selection changes.
func (c *Kernel) CallProcessSelectionNext() {
	if c.selector.Next() {
		c.render.PaintRequest(false)
	}
}

// CallProcessSelectionOptions updates the selected task's option with the given rune and value, then triggers a repaint request.
// Returns true on successful task retrieval and option update, otherwise returns false.
func (c *Kernel) CallProcessSelectionOptions(option rune, value float64) bool {
	return c.doProcessSelectionOptions(option, value)
}

// CallPaintRequest triggers a paint request via the render component and returns true if successful.
func (c *Kernel) CallPaintRequest() bool {
	return c.render.PaintRequest(false)
}

// CallWritePromptEOL writes the specified prompt followed by an end-of-line based on the eol flag using the render instance.
func (c *Kernel) CallWritePromptEOL(prompt string, eol bool) {
	c.render.WritePromptEOL(prompt, eol)
}

// CallWritePromptLine sends a formatted prompt and line to the renderer for output using the WritePromptLine method.
func (c *Kernel) CallWritePromptLine(prompt string, line string) {
	c.render.WritePromptLine(prompt, line)
}

// CallWrite sends the provided string data to the kernel's rendering writer for output.
func (c *Kernel) CallWrite(data string) {
	c.render.Write(data)
}

// CallWriteNormal writes the provided string data to the render instance using the WriteNormal method.
func (c *Kernel) CallWriteNormal(data string) {
	c.render.WriteNormal(data)
}

// CallWriteHighlights writes syntax-highlighted content to the render component using the provided data string.
func (c *Kernel) CallWriteHighlights(data string) {
	c.render.WriteHighlight(data)
}

// CallWriteCritical writes critical data to the render component of the Kernel instance.
func (c *Kernel) CallWriteCritical(data string) {
	c.render.WriteCritical(data)
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

// CallMoveCursorLeft moves the cursor one position to the left within the render context.
func (c *Kernel) CallMoveCursorLeft() {
	c.render.MoveCursorLeft()
}

// CallMoveCursorRight moves the cursor one position to the right by invoking the render's MoveCursorRight method.
func (c *Kernel) CallMoveCursorRight() {
	c.render.MoveCursorRight()
}

// CallSaveCursor saves the current cursor state by invoking the SaveCursor method on the associated renderer.
func (c *Kernel) CallSaveCursor() {
	c.render.SaveCursor()
}

// CallRestoreCursor restores the cursor to its previous position using the render instance of the Kernel.
func (c *Kernel) CallRestoreCursor() {
	c.render.RestoreCursor()
}

// CallCWDSet sets the current working directory to the specified path and updates the shell prompt accordingly.
func (c *Kernel) CallCWDSet(arg string) bool {
	return c.fs.CWDSet(arg)
}

// CallCWDGet returns the command path of the current working directory from the file system.
func (c *Kernel) CallCWDGet() string {
	return c.fs.CWDCommandPath()
}

// CallCWDPath retrieves the current working directory's path as a slice of strings from the filesystem instance.
func (c *Kernel) CallCWDPath() []string {
	return c.fs.CWDPath()
}

// CallCWDName returns the name of the current working directory as a string.
func (c *Kernel) CallCWDName() string {
	return c.fs.CWDName()
}

// CallCWDDirectoryListing retrieves the directory listing of the current working directory as a slice of strings.
func (c *Kernel) CallCWDDirectoryListing() []string {
	return c.fs.CWDDirectoryListing()
}

// CallHistory applies a history action to the shell and invokes task execution if arguments are produced.
func (c *Kernel) CallHistory(_ interfaces.HistoryAction, _ int) {
	//TOD IMPLEMENT HISTORY!!!!
	//if arg := c.shOld.HistoryApply(verb, idx); len(arg) > 0 {
	//	_, _ = c.CallTaskExec(arg, nil)
	//}
}

// CallSuggestion provides autocomplete suggestions and context for a given input string at a specified cursor position.
func (c *Kernel) CallSuggestion(in string, cursor int) (string, []string, bool) {
	return c.fs.Suggestion(in, cursor)
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
	return c.doProcessSetFg(pid)
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
	//TODO REIMPLEMENT!!!!
	const commandActivate = "activate"

	if kind == interfaces.KeyTypeCtrl {
		switch key {
		case 3:
			c.selector.Clear()
			if c.foreground != nil {
				c.doProcessKill(c.foreground.PID())
			}
		case 4:
			pid, name := c.doProcessGetForegroundName()
			if pid != adaptiveticker.UnknownId && name != commandActivate {
				c.doProcessSetBackground()
				_, _ = c.doProcessExec(fmt.Sprint(commandActivate, " ", pid), nil)
			}
		}
		return
	}

	if c.foreground != nil {
		if process := c.pid2Process(c.foreground.PID()); process != nil {
			if readEvent := process.GetCommand().ReadEvent(); readEvent != nil {
				readEvent(process, int(kind), key)
			}
		}
		return
	}
}

// Start initializes the kernel's event handling loop and begins processing I/O operations asynchronously.
func (c *Kernel) Start() {
	c.shell, _ = c.doProcessExec(c.shellPath, nil)
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

// taskExecutor executes a command by parsing the input line, creating a task, and managing its lifecycle and state.
// Returns true and error if the task was created but execution failed, or true and nil if execution succeeded.
func (c *Kernel) doProcessExec(line string, options *interfaces.ProcessOptions) (interfaces.IProcess, error) {
	cmd, args, err := c.fs.Find(line)
	if err != nil {
		return nil, fmt.Errorf("error creating task: invalid command '%s'", line)
	}
	process := c.pf.Create(cmd, line, options)
	if !c.ids.Set(process) {
		return nil, fmt.Errorf("error creating task: can't set pid")
	}
	if err = cmd.Execute(process, args); err != nil {
		c.doProcessKill(process.PID())
		return nil, err
	}
	if !cmd.Daemon() {
		c.doProcessKill(process.PID())
		return nil, nil
	}
	process.SetState(interfaces.ProcessStateRunning)
	if !cmd.Background() {
		c.foreground = process
	}
	return process, nil
}

// doProcessKill terminates and removes a task by its process ID (pid). Returns true if successful, false if the pid is not found.
func (c *Kernel) doProcessKill(pid int) bool {
	t, ok := c.ids.Get(pid)
	if !ok {
		return false
	}
	task, ok := t.(interfaces.IProcess)
	if !ok {
		return false
	}
	if task == c.shell {
		return false
	}
	if len(task.Timers()) > 0 {
		c.ticker.RemoveEntries(task.Timers())
	}
	if c.foreground != nil {
		if c.foreground.PID() == pid {
			c.foreground = c.shell //nil
		}
	}
	c.ids.Unset(pid)
	return true
}

// CallTaskKillAll terminates all tasks matching the specified name. Returns the number of tasks successfully terminated.
func (c *Kernel) doProcessKillAll(name string) int {
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
			if ok := c.doProcessKill(task.PID()); ok {
				count++
			}
		}
	}
	return count
}

// doProcessSelectionOptions updates the current process's option and triggers a repaint request. Returns false if process not found.
func (c *Kernel) doProcessSelectionOptions(option rune, value float64) bool {
	t, ok := c.ids.Get(c.selector.PID())
	if !ok {
		return false
	}
	process, ok := t.(interfaces.IProcess)
	if !ok {
		return false
	}
	process.SetOption(option, value)
	c.render.PaintRequest(true)
	return true
}

// doProcessSetBackground sets the current foreground task to the shell and returns true, or returns false if none exists.
func (c *Kernel) doProcessSetBackground() bool {
	if c.foreground == nil {
		return false
	}
	c.foreground = c.shell //nil
	return true
}

// doProcessSetSelectionMode sets the current selection mode for tasks based on a requested process ID. Defaults to the first task if the requested ID is unavailable.
func (c *Kernel) doProcessSetSelectionMode(requestedPid int) {
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

// doProcessSetFg sets the foreground task to the one associated with the given PID. Returns true if successful, false otherwise.
func (c *Kernel) doProcessSetFg(pid int) bool {
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

// doProcessGetForegroundName retrieves the PID and command name of the currently foregrounded process.
// If no foreground process exists, it returns UnknownId and an empty string.
func (c *Kernel) doProcessGetForegroundName() (int, string) {
	if c.foreground == nil {
		return adaptiveticker.UnknownId, ""
	}
	return c.foreground.PID(), c.foreground.GetCommand().Name()
}

// doProcessList retrieves a list of process descriptions by iterating through all stored processes in the Kernel.
func (c *Kernel) doProcessList() []*interfaces.ProcessDescription {
	var out []*interfaces.ProcessDescription
	c.ids.Range(func(item adaptiveticker.IIds) bool {
		process, ok := item.(interfaces.IProcess)
		if ok && process != nil {
			out = append(out, process.Description())
		}
		return true
	})
	return out
}

// pid2Process retrieves a process implementing the interfaces.IProcess interface by its PID. Returns nil if not found.
func (c *Kernel) pid2Process(pid int) interfaces.IProcess {
	p, ok := c.ids.Get(pid)
	if !ok {
		return nil
	}
	process, _ := p.(interfaces.IProcess)
	return process
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
	c.doProcessKillAll("")
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

// handleTimerEvent triggers a timer event for a task identified by the given pid and tid, with the specified interval.
// Returns true if the event was successfully triggered, otherwise false.
func (c *Kernel) handleTimerEvent(pid int, tid int, interval int) bool {
	if process := c.pid2Process(pid); process != nil {
		if timerEvent := process.GetCommand().TimerEvent(); timerEvent != nil {
			timerEvent(process, tid, interval)
			return true
		}
	}
	return false
}

// handlePaintEvent executes a rendering operation if the surface is marked as dirty, processing selected and other tasks.
// Returns true if the rendering process is executed, false otherwise.
func (c *Kernel) handlePaintEvent() bool {
	if !c.render.IsDirty() {
		return false
	}
	var selectedProcess interfaces.IProcess = nil
	var tasks []interfaces.IProcess
	c.ids.Range(func(item adaptiveticker.IIds) bool {
		process, ok := item.(interfaces.IProcess)
		if ok {
			if process.PID() == c.selector.PID() {
				selectedProcess = process
			} else {
				tasks = append(tasks, process)
			}
		}
		return true
	})
	return c.render.ExecPaint(selectedProcess, tasks)
}

/*

// tasksFileExtension specifies the file extension used for task-related data files.
const tasksFileExtension = ".task"

// commandActivate represents the command string for activation operations.
// commandTask represents the command string for task-related operations.
const (
	commandActivate = "activate"
	commandTask     = "task"
)


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

*/
