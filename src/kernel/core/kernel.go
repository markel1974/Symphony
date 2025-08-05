package core

import (
	"fmt"
	"github.com/markel1974/c64emu/src/kernel/adaptiveticker"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/messages"
	"github.com/markel1974/c64emu/src/kernel/process_factory"
	"log"
)

// Kernel represents the core component responsible for managing rendering, input/output, task execution, and timers.
type Kernel struct {
	user         string
	ticker       *adaptiveticker.AdaptiveTicker
	inputDriver  interfaces.IKeyboardDriver
	renderServer interfaces.IRender
	foreground   interfaces.IProcess
	pidGenerator *adaptiveticker.Ids
	running      map[int]interfaces.IProcess
	fsServer     interfaces.IFileSystem
	shellPath    string
	shell        interfaces.IProcess
	messageChan  chan interfaces.IMessage
	pf           *process_factory.ProcessFactory
	timersChan   chan *adaptiveticker.TimerHandler
	servers      []interfaces.IServer
	exit         bool
	handlers     map[interfaces.MessageType]func(interfaces.IMessage)
}

// NewKernel creates and returns a new Kernel instance, initializing its dependencies and internal fields.
func NewKernel(user string, ticker *adaptiveticker.AdaptiveTicker, timersChan chan *adaptiveticker.TimerHandler, inputDriver interfaces.IKeyboardDriver, renderServer interfaces.IRender, fsServer interfaces.IFileSystem, shellPath string) *Kernel {
	t := &Kernel{
		user:         user,
		ticker:       ticker,
		inputDriver:  inputDriver,
		renderServer: renderServer,
		fsServer:     fsServer,
		foreground:   nil,
		pidGenerator: adaptiveticker.NewIds(1024),
		messageChan:  make(chan interfaces.IMessage, contextMaQueueLen),
		timersChan:   timersChan,
		exit:         false,
		shellPath:    shellPath,
		running:      make(map[int]interfaces.IProcess),
		handlers:     make(map[interfaces.MessageType]func(interfaces.IMessage)),
	}
	t.pf = process_factory.NewProcessFactory(t)

	t.handlers[interfaces.MessageTypeRead] = t.handleReadEvent
	t.handlers[interfaces.MessageTypeTimer] = t.handleTimerEvent
	t.handlers[interfaces.MessageTypeQuit] = t.handleQuitEvent
	t.handlers[interfaces.MessageTypeTimedMessage] = t.handleTimedMessage
	return t
}

func (c *Kernel) PID() int {
	return 0
}

func (c *Kernel) User() string {
	return c.user
}

func (c *Kernel) AddServer(server interfaces.IServer) {
	c.servers = append(c.servers, server)
	for _, r := range server.Register() {
		c.handlers[r] = server.PostMessage
	}
	server.SetRouter(c)
	server.Start()
}

// PostMessage sends the provided IMessage to the Kernel's internal message channel for further processing.
func (c *Kernel) PostMessage(msg interfaces.IMessage) {
	c.messageChan <- msg
}

// PostTimedMessage schedules a message for execution based on specified timing parameters and count.
//func (c *Kernel) PostTimedMessage(msg interfaces.IMessage, first int64, interval int64, count int64) {
//	_ = c.ticker.Create(c.timersChan, msg, first, interval, count)
//}

// CallProcessList returns a formatted string containing process IDs and their respective command names managed by the Kernel.
func (c *Kernel) CallProcessList() []*interfaces.ProcessDescription {
	return c.doProcessList()
}

// CallProcessExec executes a command by parsing the input line, creating a process, and managing its lifecycle and state.
// Returns true and error if the process was created but execution failed, or true and nil if execution succeeded.
func (c *Kernel) CallProcessExec(user string, line string) (bool, error) {
	_, err := c.doProcessExec(user, line)
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

// CallProcessSetForeground sets the foreground task to the one associated with the given PID. Returns true if successful, false otherwise.
func (c *Kernel) CallProcessSetForeground(pid int) bool {
	return c.doProcessSetForeground(pid)
}

// CallProcessIsActive checks if a process with the given PID is currently active in the Kernel's activeProcess map.
func (c *Kernel) CallProcessIsActive(pid int) bool {
	process, _ := c.running[pid]
	return process != nil
}

// CallWindowsSelectionBegin updates the selection mode for a specific process and triggers a repaint without requesting a redraw.
func (c *Kernel) CallWindowsSelectionBegin() {
	c.renderServer.CallWindowsSelectionBegin()
}

// CallWindowsSelectionOptions updates the selected task's option with the given rune and value, then triggers a repaint request.
// Returns true on successful task retrieval and option update, otherwise returns false.
func (c *Kernel) CallWindowsSelectionOptions(option rune, value float64) {
	c.renderServer.CallWindowsSelectionOptions(option, value)
}

// CallWindowsSelectionPrevious moves the task selection to the previous task and triggers a render update if successful.
func (c *Kernel) CallWindowsSelectionPrevious() {
	c.renderServer.CallWindowsSelectionPrevious()
}

// CallWindowsSelectionNext advances the task windowSelector to the next task and triggers a repaint if the task selection changes.
func (c *Kernel) CallWindowsSelectionNext() {
	c.renderServer.CallWindowsSelectionNext()
}

// CallWindowsSelectionEnd clears the state of the associated WindowSelector instance by resetting its index and available list.
func (c *Kernel) CallWindowsSelectionEnd() {
	c.renderServer.CallWindowsSelectionEnd()
}

// CallPaintRequest triggers a paint request via the render component and returns true if successful.
//func (c *Kernel) CallPaintRequest() {
//	c.renderServer.CallPaintRequest()
//}

// CallWritePromptEOL writes the specified prompt followed by an end-of-line based on the eol flag using the render instance.
func (c *Kernel) CallWritePromptEOL(prompt string, eol bool) {
	c.renderServer.CallWritePromptEOL(prompt, eol)
}

// CallWritePromptLine sends a formatted prompt and line to the renderer for output using the WritePromptLine method.
func (c *Kernel) CallWritePromptLine(prompt string, line string) {
	c.renderServer.CallWritePromptLine(prompt, line)
}

// CallWrite sends the provided string data to the kernel's rendering writer for output.
func (c *Kernel) CallWrite(data string) {
	c.renderServer.CallWrite(data)
}

// CallWriteNormal writes the provided string data to the render instance using the WriteNormal method.
func (c *Kernel) CallWriteNormal(data string) {
	c.renderServer.CallWriteNormal(data)
}

// CallWriteHighlights writes syntax-highlighted content to the render component using the provided data string.
func (c *Kernel) CallWriteHighlights(data string) {
	c.renderServer.CallWriteHighlight(data)
}

// CallWriteCritical writes critical data to the render component of the Kernel instance.
func (c *Kernel) CallWriteCritical(data string) {
	c.renderServer.CallWriteCritical(data)
}

// CallWriteLn writes the provided string followed by a new line to the kernel's output stream.
func (c *Kernel) CallWriteLn(data string) {
	c.renderServer.CallWriteLn(data)
}

// CallWriteColor writes a string to the output with specified foreground color, background color, and color mode.
func (c *Kernel) CallWriteColor(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	c.renderServer.CallWriteColor(data, fg, bg, mode)
}

// CallWriteColorLn writes a line of text with specified foreground and background colors and a given color mode.
func (c *Kernel) CallWriteColorLn(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	c.renderServer.CallWriteColorLn(data, fg, bg, mode)
}

// CallClearScreen clears the screen by invoking the associated renderer's CreateClearScreen method.
func (c *Kernel) CallClearScreen() {
	c.renderServer.CallClearScreen()
}

// CallScreenSize retrieves the screen's width and height as integers from the render instance.
func (c *Kernel) CallScreenSize() (int, int) {
	return c.renderServer.CallGetScreenSize()
}

// CallMoveCursorLeft moves the cursor one position to the left within the render context.
func (c *Kernel) CallMoveCursorLeft() {
	c.renderServer.CallMoveCursorLeft()
}

// CallMoveCursorRight moves the cursor one position to the right by invoking the render's CreateMoveCursorRight method.
func (c *Kernel) CallMoveCursorRight() {
	c.renderServer.CallMoveCursorRight()
}

// CallSaveCursor saves the current cursor state by invoking the CreateSaveCursor method on the associated renderer.
func (c *Kernel) CallSaveCursor() {
	c.renderServer.CallSaveCursor()
}

// CallRestoreCursor restores the cursor to its previous position using the render instance of the Kernel.
func (c *Kernel) CallRestoreCursor() {
	c.renderServer.CallRestoreCursor()
}

// CallCWDSet sets the current working directory to the specified path and updates the shell prompt accordingly.
func (c *Kernel) CallCWDSet(arg string) bool {
	return c.fsServer.CWDSet(arg)
}

// CallCWDGet returns the command path of the current working directory from the file system.
func (c *Kernel) CallCWDGet() string {
	return c.fsServer.CWDCommandPath()
}

// CallCWDPath retrieves the current working directory's path as a slice of strings from the filesystem instance.
func (c *Kernel) CallCWDPath() []string {
	return c.fsServer.CWDPath()
}

// CallCWDName returns the name of the current working directory as a string.
func (c *Kernel) CallCWDName() string {
	return c.fsServer.CWDName()
}

// CallCWDDirectoryListing retrieves the directory listing of the current working directory as a slice of strings.
func (c *Kernel) CallCWDDirectoryListing() []string {
	return c.fsServer.CWDDirectoryListing()
}

// CallFileSystemSuggestion provides autocomplete suggestions and context for a given input string at a specified cursor position.
func (c *Kernel) CallFileSystemSuggestion(in string, cursor int) (string, []string, bool) {
	return c.fsServer.Suggestion(in, cursor)
}

// CallFileSystemHelp retrieves the help information associated with the given argument and returns it as a string.
// Returns an error if the help information cannot be fetched.
func (c *Kernel) CallFileSystemHelp(arg string) (string, error) {
	return c.fsServer.Help(arg)
}

// CallSetScreenSize adjusts the screen dimensions to the specified width and height values.
func (c *Kernel) CallSetScreenSize(w int, h int) {
	c.renderServer.CallSetScreenSize(w, h)
}

// CallExitRequested sets the `exit` flag to true, signaling that an exit has been requested for the kernel.
func (c *Kernel) CallExitRequested() {
	c.exit = true
}

// CallTimerCreate initializes a timer for a process with specified timing parameters if the process and its timer event exist.
// It creates a new message timer, assigns a timer ID, and associates the timer with the process if successful.
func (c *Kernel) CallTimerCreate(pid int, first int, interval int, count int) {
	process, _ := c.running[pid]
	if process == nil {
		return
	}
	if process.GetCommand().TimerEvent == nil {
		return
	}
	m := messages.NewMessageTimer(pid, interval)
	m.SetTID(c.ticker.Create(c.timersChan, m, int64(first), int64(interval), int64(count)))
	if m.TID() > -1 {
		process.AddTimer(m.TID())
	}
	return
}

// CallTimerStop stops a timer associated with a specific process and thread id, returning true if successful or false otherwise.
func (c *Kernel) CallTimerStop(pid int, tid int) {
	process, _ := c.running[pid]
	if process == nil {
		return
	}
	c.closeTimer(process, tid)
}

// Start initializes the kernel's event handling loop and begins processing I/O operations asynchronously.
func (c *Kernel) Start() {
	c.shell, _ = c.doProcessExec(c.user, c.shellPath)
	d := make(chan bool)
	go func() {
		d <- true
		readBuffer := make([]byte, 4096)
		for {
			k, v, err := c.inputDriver.ScanKey(readBuffer)
			if err == nil {
				if k != interfaces.KeyTypeNone {
					re := messages.NewMessageRead(k, v, false)
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
func (c *Kernel) doProcessExec(user string, line string) (interfaces.IProcess, error) {
	cmd, args, err := c.fsServer.Find(line)
	if err != nil {
		return nil, fmt.Errorf("error creating task: invalid command '%s'", line)
	}
	process := c.pf.Create(user, cmd, line)
	if !c.pidGenerator.Set(process) {
		return nil, fmt.Errorf("error creating task: can't set pid")
	}
	c.running[process.PID()] = process
	process.Start()
	for _, server := range c.servers {
		server.NotifyProcessCreation(process.Description())
	}
	if !cmd.Background() {
		c.doProcessSetForeground(process.PID())
	}
	//process.PostMessage(messages.NewMessageProcessStart(args))
	if err = cmd.Execute(process, args); err != nil {
		c.doProcessKill(process.PID())
		return nil, err
	}
	if !cmd.Daemon() {
		c.doProcessKill(process.PID())
		return nil, nil
	}
	return process, nil
}

func (c *Kernel) doProcessSetForeground(pid int) bool {
	process, _ := c.running[pid]
	if process == nil {
		return false
	}
	for _, s := range c.servers {
		s.NotifyProcessForeground(process.Description())
	}
	if c.foreground != process {
		c.foreground = process
		fmt.Println("TODO FOREGROUND MESSAGE!", c.foreground.GetCommand().Name())
	}
	return true
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
	process.PostMessage(messages.NewMessageQuit())
	for _, server := range c.servers {
		server.NotifyProcessTermination(process.Description())
	}
	if c.foreground != nil {
		if c.foreground.PID() == pid {
			c.doProcessSetForeground(c.shell.PID())
		}
	}
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
			c.handleMessageEvent(t.Event.(interfaces.IMessage))
		}
		if c.exit {
			c.shutdown()
			return
		}
	}
}

// handleMessageEvent processes an incoming IMessage by dispatching it to the appropriate handlers based on its type.
func (c *Kernel) handleMessageEvent(m interfaces.IMessage) {
	if handler, ok := c.handlers[m.GetType()]; ok {
		handler(m)
	} else {
		log.Printf("unknown message type: %d", m.GetType())
	}
}

// handleReadEvent processes input events based on their type and key value to handle control, foreground tasks, and system state.
func (c *Kernel) handleReadEvent(m interfaces.IMessage) {
	mm, ok := m.(*messages.MessageRead)
	if !ok {
		return
	}
	for _, process := range c.running {
		if readBroadcastEvent := process.GetCommand().ReadBroadcastEvent(); readBroadcastEvent != nil {
			process.PostMessage(messages.NewMessageRead(mm.Kind(), mm.Data(), true))
		}
	}
	if c.foreground != nil {
		if readBroadcastEvent := c.foreground.GetCommand().ReadEvent(); readBroadcastEvent != nil {
			c.foreground.PostMessage(mm)
		}
	}
}

// handleTimerEvent triggers a timer event for a task identified by the given pid and tid, with the specified interval.
// Returns true if the event was successfully triggered, otherwise false.
func (c *Kernel) handleTimerEvent(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageTimer)
	if !ok {
		return
	}
	if process := c.running[mt.PID()]; process != nil {
		process.PostMessage(mt)
	}
}

func (c *Kernel) handleTimedMessage(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageTimedMessage)
	if !ok {
		return
	}
	_ = c.ticker.Create(c.timersChan, mt.Message(), mt.First(), mt.Interval(), mt.Count())
}

// handleQuitEvent handles a quit message by verifying its type and setting the kernel's exit flag to true.
func (c *Kernel) handleQuitEvent(m interfaces.IMessage) {
	_, ok := m.(*messages.MessageQuit)
	if !ok {
		return
	}
	c.exit = true
}
