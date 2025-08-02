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
	ticker         *adaptiveticker.AdaptiveTicker
	inputDriver    io.Reader
	outputDriver   io.Writer
	render         interfaces.IRender
	foreground     interfaces.IProcess
	windowSelector *WindowSelector
	pidGenerator   *adaptiveticker.Ids
	running        map[int]interfaces.IProcess
	fs             interfaces.IFileSystem
	shellPath      string
	shell          interfaces.IProcess
	messageChan    chan messages.IMessage
	pf             *process_factory.ProcessFactory
	timersChan     chan *adaptiveticker.TimerHandler
	exit           bool
}

// NewKernel creates and returns a new Kernel instance, initializing its dependencies and internal fields.
func NewKernel(ticker *adaptiveticker.AdaptiveTicker, timersChan chan *adaptiveticker.TimerHandler, inputDriver io.Reader, outputDriver io.Writer, render interfaces.IRender, fs interfaces.IFileSystem, shellPath string) *Kernel {
	t := &Kernel{
		ticker:         ticker,
		inputDriver:    inputDriver,
		outputDriver:   outputDriver,
		render:         render,
		fs:             fs,
		foreground:     nil,
		windowSelector: NewWindowSelector(),
		pidGenerator:   adaptiveticker.NewIds(1024),
		messageChan:    make(chan messages.IMessage, contextMaQueueLen),
		timersChan:     timersChan,
		exit:           false,
		shellPath:      shellPath,
		running:        make(map[int]interfaces.IProcess),
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
func (c *Kernel) CallProcessExec(line string, options *interfaces.WindowOptions) (bool, error) {
	_, err := c.doProcessExec(line, options)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CallProcessKillForeground terminates the current foreground process and returns true if successful.
func (c *Kernel) CallProcessKillForeground() bool {
	return c.doProcessKill(c.foreground.PID())
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

// CallProcessIsActive checks if a process with the given PID is currently active in the Kernel's activeProcess map.
func (c *Kernel) CallProcessIsActive(pid int) bool {
	process, _ := c.running[pid]
	return process != nil
}

// CallProcessSetFg sets the foreground task to the one associated with the given PID. Returns true if successful, false otherwise.
func (c *Kernel) CallProcessSetFg(pid int) bool {
	return c.doProcessSetFg(pid)
}

// CallWindowsSelectionBegin updates the selection mode for a specific process and triggers a repaint without requesting a redraw.
func (c *Kernel) CallWindowsSelectionBegin() {
	if c.foreground != nil {
		c.doCreateWindowsSelection(c.foreground.PID())
		c.render.PaintRequest(false)
	}
}

// CallWindowsSelectionEnd clears the state of the associated WindowSelector instance by resetting its index and available list.
func (c *Kernel) CallWindowsSelectionEnd() {
	c.windowSelector.Clear()
}

// CallWindowsSelectionPrevious moves the task selection to the previous task and triggers a render update if successful.
func (c *Kernel) CallWindowsSelectionPrevious() {
	if c.windowSelector.Prev() {
		c.render.PaintRequest(false)
	}
}

// CallWindowsSelectionNext advances the task windowSelector to the next task and triggers a repaint if the task selection changes.
func (c *Kernel) CallWindowsSelectionNext() {
	if c.windowSelector.Next() {
		c.render.PaintRequest(false)
	}
}

// CallWindowsSelectionOptions updates the selected task's option with the given rune and value, then triggers a repaint request.
// Returns true on successful task retrieval and option update, otherwise returns false.
func (c *Kernel) CallWindowsSelectionOptions(option rune, value float64) bool {
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

// CallFileSystemSuggestion provides autocomplete suggestions and context for a given input string at a specified cursor position.
func (c *Kernel) CallFileSystemSuggestion(in string, cursor int) (string, []string, bool) {
	return c.fs.Suggestion(in, cursor)
}

// CallFileSystemHelp retrieves the help information associated with the given argument and returns it as a string.
// Returns an error if the help information cannot be fetched.
func (c *Kernel) CallFileSystemHelp(arg string) (string, error) {
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

// CallTimerCreate creates a timer for the specified process ID with given start time, interval, and execution count.
// Returns true if the timer is successfully created, otherwise false.
// It requires the process to have a valid TimerEvent handler.
// Adds the timer ID to the process's list of active timers.
func (c *Kernel) CallTimerCreate(pid int, first int, interval int, count int) bool {
	process, _ := c.running[pid]
	if process == nil {
		return false
	}
	if process.GetCommand().TimerEvent == nil {
		return false
	}
	m := messages.NewMessageTimer(pid, interval)
	m.SetTID(c.ticker.Create(c.timersChan, m, int64(first), int64(interval), int64(count)))
	if m.TID() > -1 {
		process.AddTimer(m.TID())
	}
	return true
}

// CallTimerStop stops a timer associated with a specific process and thread id, returning true if successful or false otherwise.
func (c *Kernel) CallTimerStop(pid int, tid int) bool {
	process, _ := c.running[pid]
	if process == nil {
		return false
	}
	return c.closeTimer(process, tid)
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
	for _, process := range c.running {
		if readBroadcastEvent := process.GetCommand().ReadBroadcastEvent(); readBroadcastEvent != nil {
			readBroadcastEvent(process, int(kind), key)
		}
	}

	if c.foreground != nil {
		if readEvent := c.foreground.GetCommand().ReadEvent(); readEvent != nil {
			readEvent(c.foreground, int(kind), key)
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
func (c *Kernel) doProcessExec(line string, options *interfaces.WindowOptions) (interfaces.IProcess, error) {
	cmd, args, err := c.fs.Find(line)
	if err != nil {
		return nil, fmt.Errorf("error creating task: invalid command '%s'", line)
	}
	process := c.pf.Create(cmd, line, options)
	if !c.pidGenerator.Set(process) {
		return nil, fmt.Errorf("error creating task: can't set pid")
	}
	c.running[process.PID()] = process
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
	process, _ := c.running[pid]
	if process == nil {
		return false
	}
	if process == c.shell {
		return false
	}
	if len(process.Timers()) > 0 {
		c.ticker.RemoveEntries(process.Timers())
	}
	if c.foreground != nil {
		if c.foreground.PID() == pid {
			c.foreground = c.shell //nil
		}
	}
	c.windowSelector.Clear()
	c.pidGenerator.Unset(pid)
	delete(c.running, pid)
	return true
}

// CallTaskKillAll terminates all tasks matching the specified name. Returns the number of tasks successfully terminated.
func (c *Kernel) doProcessKillAll(name string) int {
	count := 0
	var tasks []interfaces.IProcess
	for _, process := range c.running {
		tasks = append(tasks, process)
	}
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
	process, _ := c.running[c.windowSelector.PID()]
	if process == nil {
		return false
	}
	process.SetWindowOption(option, value)
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

// doCreateWindowsSelection sets the current selection mode for tasks based on a requested process ID. Defaults to the first task if the requested ID is unavailable.
func (c *Kernel) doCreateWindowsSelection(requestedPid int) {
	var idx = 0
	var firstPid = adaptiveticker.UnknownId
	var firstIdx = 0

	c.windowSelector.Clear()

	for _, process := range c.running {
		if process.GetCommand().PaintEvent() != nil {
			c.windowSelector.AddAvailable(process.PID())
			if firstPid == adaptiveticker.UnknownId {
				firstPid = process.PID()
				firstIdx = idx
			}
			if process.PID() == requestedPid {
				c.windowSelector.Set(requestedPid, idx)
			}
			idx++
		}
	}
	if c.windowSelector.PID() == adaptiveticker.UnknownId {
		if firstPid == adaptiveticker.UnknownId {
			return
		}
		c.windowSelector.Set(firstPid, firstIdx)
	}
}

// doProcessSetFg sets the foreground task to the one associated with the given PID. Returns true if successful, false otherwise.
func (c *Kernel) doProcessSetFg(pid int) bool {
	process, _ := c.running[pid]
	if process == nil {
		return false
	}
	c.foreground = process
	return true
}

// doProcessList retrieves a list of process descriptions by iterating through all stored processes in the Kernel.
func (c *Kernel) doProcessList() []*interfaces.ProcessDescription {
	var out []*interfaces.ProcessDescription
	for _, process := range c.running {
		out = append(out, process.Description())
	}
	return out
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
	if process := c.running[pid]; process != nil {
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
	for _, process := range c.running {
		if process.PID() == c.windowSelector.PID() {
			selectedProcess = process
		} else {
			tasks = append(tasks, process)
		}
	}
	return c.render.ExecPaint(selectedProcess, tasks)
}
