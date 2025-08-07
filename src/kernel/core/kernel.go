package core

import (
	"fmt"
	"log"

	"github.com/markel1974/c64emu/src/kernel/adaptiveticker"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/messages"
	"github.com/markel1974/c64emu/src/kernel/process_factory"
)

// Kernel represents the core component responsible for managing rendering, input/output, process execution, and timers.
type Kernel struct {
	user         string
	ticker       *adaptiveticker.AdaptiveTicker
	inputDriver  interfaces.IKeyboardDriver
	renderServer interfaces.IRender
	foreground   interfaces.IProcess
	pidGenerator *adaptiveticker.Ids
	running      map[int]*KernelProcess
	shellPath    string
	messageChan  chan interfaces.IMessage
	pf           *process_factory.ProcessFactory
	timersChan   chan *adaptiveticker.TimerHandler
	servers      []interfaces.IServer
	exit         bool
	handlers     map[interfaces.MessageType]func(interfaces.IMessage)
}

// NewKernel creates and returns a new Kernel instance, initializing its dependencies and internal fields.
func NewKernel(user string, ticker *adaptiveticker.AdaptiveTicker, inputDriver interfaces.IKeyboardDriver, renderServer interfaces.IRender, shellPath string) *Kernel {
	t := &Kernel{
		user:         user,
		ticker:       ticker,
		inputDriver:  inputDriver,
		renderServer: renderServer,
		foreground:   nil,
		pidGenerator: adaptiveticker.NewIds(1024),
		messageChan:  make(chan interfaces.IMessage, contextMaQueueLen),
		timersChan:   make(chan *adaptiveticker.TimerHandler, contextMaQueueLen),
		exit:         false,
		shellPath:    shellPath,
		running:      make(map[int]*KernelProcess),
		handlers:     make(map[interfaces.MessageType]func(interfaces.IMessage)),
	}
	t.pf = process_factory.NewProcessFactory(t)

	t.handlers[interfaces.MessageTypeRead] = t.handleReadEvent
	t.handlers[interfaces.MessageTypeTimer] = t.handleTimerEvent
	t.handlers[interfaces.MessageTypeQuit] = t.handleQuitEvent
	t.handlers[interfaces.MessageTypeTimedMessage] = t.handleTimedMessage
	t.handlers[interfaces.MessageTypeProcessExit] = t.handleProcessExit
	t.handlers[interfaces.MessageTypeProcessExec] = t.handleProcessExec
	t.handlers[interfaces.MessageTypeProcessSetForeground] = t.handleProcessSetForeground
	t.handlers[interfaces.MessageTypeProcessKill] = t.handleProcessKill
	t.handlers[interfaces.MessageTypeProcessKillAll] = t.handleProcessKillAll
	t.handlers[interfaces.MessageTypeProcessKillForeground] = t.handleProcessKillForeground
	t.handlers[interfaces.MessageTypeTimerCreate] = t.handleTimerCreate
	t.handlers[interfaces.MessageTypeTimerStop] = t.handleTimerStop
	t.handlers[interfaces.MessageTypeProcessList] = t.handleProcessList
	t.handlers[interfaces.MessageTypeFileSystemFindResponse] = t.handleFileSystemFindResponse
	return t
}

// SetScreenSize adjusts the screen dimensions to the specified width and height values.
func (c *Kernel) SetScreenSize(w int, h int) {
	c.messageChan <- messages.NewMessageSetScreenSize(c, w, h)
}

// Process returns the current foreground process.
func (c *Kernel) Process() interfaces.IProcess {
	return nil
}

// PID retrieves the process identifier (PID) of the kernel instance.
func (c *Kernel) PID() int {
	return 0
}

// User returns the name of the user associated with the Kernel instance.
func (c *Kernel) User() string {
	return c.user
}

// AddServer adds a new server to the kernel, registers its handlers, sets the router, and starts the server.
func (c *Kernel) AddServer(server interfaces.IServer) {
	c.servers = append(c.servers, server)
	handlers := server.Register(c)
	for _, h := range handlers {
		c.handlers[h] = server.PostMessage
	}
	server.Start()
}

// PostMessage sends the provided IMessage to the Kernel's internal message channel for further processing.
func (c *Kernel) PostMessage(msg interfaces.IMessage) {
	c.messageChan <- msg
}

// CallProcessIsActive checks if a process with the given PID is currently active in the Kernel's activeProcess map.
func (c *Kernel) CallProcessIsActive(router interfaces.IRouter, pid int) bool {
	active, _ := c.running[pid]
	return active != nil
}

// CallScreenSize retrieves the screen's width and height as integers from the render instance.
func (c *Kernel) CallScreenSize(router interfaces.IRouter) (int, int) {
	return c.renderServer.CallGetScreenSize(router)
}

// CallCWDPath returns the command path of the current working directory from the file system.
//func (c *Kernel) CallCWDPath(router interfaces.IRouter) string {
//	return c.fsServer.CallCWDPath(router)
//}

// CallCWDDirectoryListing retrieves the directory listing of the current working directory as a slice of strings.
//func (c *Kernel) CallCWDDirectoryListing(router interfaces.IRouter) []string {
//	return c.fsServer.CallCWDDirectoryListing(router)
//}

// CallFileSystemHelp retrieves the help information associated with the given argument and returns it as a string.
// Returns an error if the help information cannot be fetched.
//func (c *Kernel) CallFileSystemHelp(router interfaces.IRouter, arg string) (string, error) {
//	return c.fsServer.CallHelp(router, arg)
//}

// CallExitRequested sets the `exit` flag to true, signaling that an exit has been requested for the kernel.
func (c *Kernel) CallExitRequested(router interfaces.IRouter) {
	c.exit = true
}

// Start initializes the kernel's event handling loop and begins processing I/O operations asynchronously.
func (c *Kernel) Start() {
	c.PostMessage(messages.NewMessageFileSystemFindRequest(c, c, c.shellPath, true))
	d := make(chan bool)
	go func() {
		d <- true
		readBuffer := make([]byte, 4096)
		for {
			k, v, err := c.inputDriver.ScanKey(readBuffer)
			if err == nil {
				if k != interfaces.KeyTypeNone {
					re := messages.NewMessageRead(c, k, v, false)
					c.messageChan <- re
				}
			} else {
				qe := messages.NewMessageQuit(c)
				c.messageChan <- qe
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
			c.handleMessageEvent(m)
		case t := <-c.timersChan:
			c.handleMessageEvent(t.Event.(interfaces.IMessage))
		}
		if c.exit {
			c.doShutdown()
			return
		}
	}
}

// handleMessageEvent processes an incoming IMessage by dispatching it to the appropriate handlers based on its type.
func (c *Kernel) handleMessageEvent(m interfaces.IMessage) {
	if m.Response() {
		m.Router().PostMessage(m)
		return
	}
	id := m.GetType()
	if handler, _ := c.handlers[id]; handler != nil {
		handler(m)
	} else {
		log.Printf("Kernel: unknown message type: %d", id)
	}
}

// handleReadEvent processes input events based on their type and key value to handle control, foreground processes, and system state.
func (c *Kernel) handleReadEvent(m interfaces.IMessage) {
	mm, ok := m.(*messages.MessageRead)
	if !ok {
		return
	}
	for _, process := range c.running {
		if readBroadcastEvent := process.GetCommand().OnReadBroadcast(); readBroadcastEvent != nil {
			process.PostMessage(messages.NewMessageRead(m.Router(), mm.Kind(), mm.Data(), true))
		}
	}
	if c.foreground != nil {
		if readBroadcastEvent := c.foreground.GetCommand().OnRead(); readBroadcastEvent != nil {
			c.foreground.PostMessage(mm)
		}
	}
}

// handleTimerEvent triggers a timer event for a process identified by the given pid and tid, with the specified interval.
// Returns true if the event was successfully triggered, otherwise false.
func (c *Kernel) handleTimerEvent(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageTimer)
	if !ok {
		return
	}
	process, _ := c.running[mt.PID()]
	if process == nil {
		return
	}
	process.PostMessage(mt)
}

// handleTimedMessage processes a timed message by extracting its properties and scheduling it via the ticker.
func (c *Kernel) handleTimedMessage(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageTimedMessage)
	if !ok {
		return
	}
	_ = c.ticker.Create(c.timersChan, mt.Message(), mt.First(), mt.Interval(), mt.Count())
}

// handleProcessSetForeground handles a process set foreground message by setting the foreground process to the specified process.
func (c *Kernel) handleProcessExit(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageProcessExit)
	if !ok {
		return
	}
	if process, _ := c.running[mt.Router().PID()]; process != nil {
		c.doProcessExit(process)
		return
	}
}

// handleProcessExec handles process execution by validating the message type and invoking the process execution logic.
func (c *Kernel) handleProcessExec(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageProcessExec)
	if !ok {
		return
	}
	c.PostMessage(messages.NewMessageFileSystemFindRequest(c, mt.Router(), mt.Line(), false))
}

// handleProcessSetForeground handles a process set foreground message by setting the foreground process to the specified process.
func (c *Kernel) handleProcessSetForeground(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageProcessSetForeground)
	if !ok {
		return
	}
	if process, _ := c.running[mt.PID()]; process != nil {
		c.doProcessSetForeground(mt.Router(), process)
	}
}

// handleProcessKill terminates and removes a process by its process ID (pid). Returns true if successful, false if the pid is not found.
func (c *Kernel) handleProcessKill(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageProcessKill)
	if !ok {
		return
	}
	process, _ := c.running[mt.PID]
	if process == nil {
		return
	}
	c.doProcessExit(process)
}

// handleProcessKillAll handles the termination of all processes except the sender's and optionally filters by process name.
func (c *Kernel) handleProcessKillAll(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageProcessKillAll)
	if !ok {
		return
	}
	var processes []*KernelProcess
	for _, process := range c.running {
		if len(mt.Name()) != 0 {
			if process.GetCommand().Name() != mt.Name() {
				continue
			}
		}
		processes = append(processes, process)
	}
	for _, process := range processes {
		c.doProcessExit(process)
	}
}

// handleProcessKillForeground handles the termination of the foreground process.
func (c *Kernel) handleProcessKillForeground(m interfaces.IMessage) {
	_, ok := m.(*messages.MessageProcessKillForeground)
	if !ok {
		return
	}
	process, _ := c.running[c.foreground.PID()]
	if process == nil {
		return
	}
	c.doProcessExit(process)
}

// handleTimerCreate processes a timer creation request, initializes the timer, and sends a response or error message.
func (c *Kernel) handleTimerCreate(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageTimerCreate)
	if !ok {
		return
	}
	process, _ := c.running[m.Router().PID()]
	if process == nil {
		return
	}
	msgTimer := messages.NewMessageTimer(mt.Router(), mt.Router().PID(), mt.Interval())
	msgTimer.SetTID(c.ticker.Create(c.timersChan, msgTimer, int64(mt.First()), int64(mt.Interval()), int64(mt.Count())))
	if msgTimer.TID() < 0 {
		m.Router().PostMessage(messages.NewMessageError(m.Router(), fmt.Errorf("error creating timer")))
		return
	}
	process.AddTimer(msgTimer.TID())
	m.Router().PostMessage(messages.NewMessageTimerCreated(m.Router(), msgTimer.TID()))
}

func (c *Kernel) handleTimerStop(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageTimerStop)
	if !ok {
		return
	}
	process, _ := c.running[m.Router().PID()]
	if process == nil {
		m.Router().PostMessage(messages.NewMessageError(m.Router(), fmt.Errorf("error creating timer")))
		return
	}
	c.doCloseTimer(process, mt.TID())
}

// handleProcessList processes a message requesting a list of running processes and sends the response with process details.
func (c *Kernel) handleProcessList(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageProcessList)
	if !ok {
		return
	}
	var out []*interfaces.ProcessDescription
	for _, process := range c.running {
		out = append(out, process.Description())
	}
	mt.SetProcesses(out)
	mt.Router().PostMessage(mt)
}

// handleQuitEvent handles a quit message by verifying its type and setting the kernel's exit flag to true.
func (c *Kernel) handleQuitEvent(m interfaces.IMessage) {
	_, ok := m.(*messages.MessageQuit)
	if !ok {
		return
	}
	c.exit = true
}

// doProcessExec executes a process by creating it, assigning a pid, and starting it with the given user and input line.
// Configures the command arguments, initializes the process, and notifies servers about its creation.
func (c *Kernel) handleFileSystemFindResponse(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageFileSystemFindResponse)
	if !ok {
		return
	}
	originator := mt.Parent()
	if originator == nil {
		log.Printf("error creating task: invalid originator")
		return
	}
	cmd, args, err := mt.GetResult()
	if err != nil {
		originator.PostMessage(messages.NewMessageError(originator, fmt.Errorf("error creating task: invalid command '%s'", mt.Line())))
		return
	}
	pid := NewPID()
	_, ok = c.pidGenerator.Set(pid)
	if !ok {
		originator.PostMessage(messages.NewMessageError(originator, fmt.Errorf("error creating task: can't set pid")))
		return
	}
	parent, _ := c.running[originator.PID()]
	kernelProcess := NewKernelProcess(originator.User(), mt.Line(), cmd.Name(), parent, pid, mt.Protected(), c.pf.Create(pid.GetId(), originator.User(), cmd))
	c.running[kernelProcess.PID()] = kernelProcess
	kernelProcess.Setup()
	for _, server := range c.servers {
		server.NotifyProcessCreation(kernelProcess.PID(), kernelProcess.GetCommand().Name())
	}
	kernelProcess.PostMessage(messages.NewMessageProcessStart(originator, args))
}

// doProcessSetForeground sets the specified process as the foreground process and sends activation messages if needed.
func (c *Kernel) doProcessSetForeground(router interfaces.IRouter, process interfaces.IProcess) {
	for _, s := range c.servers {
		s.NotifyProcessForeground(process.PID())
	}
	if c.foreground != process {
		c.foreground = process
		c.foreground.PostMessage(messages.NewMessageProcessActivate(router))
	}
}

// doProcessExit handles the termination process of a given IProcess, ensuring cleanup of resources and notifying observers.
func (c *Kernel) doProcessExit(process *KernelProcess) {
	if process.Protected() {
		return
	}
	if len(process.Timers()) > 0 {
		c.ticker.RemoveEntries(process.Timers())
	}
	for _, server := range c.servers {
		server.NotifyProcessTermination(process.PID())
	}
	if c.foreground != nil {
		if c.foreground.PID() == process.PID() {
			if parent := process.Parent(); parent != nil {
				c.doProcessSetForeground(c, parent)
			} else {
				log.Printf("Fatal Error: foreground process is nil")
			}
		}
	}
	delete(c.running, process.PID())
	c.pidGenerator.Unset(process.PID())
	process.PostMessage(messages.NewMessageQuit(process))
}

// doCloseTimer removes a timer with the specified ID from the task and ticker, returning true if the timer is successfully removed.
func (c *Kernel) doCloseTimer(process *KernelProcess, tid int) bool {
	ret := false
	if process != nil {
		process.TimersIterator(func(timerId int) bool {
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
func (c *Kernel) doShutdown() {
	var processes []*KernelProcess
	for _, process := range c.running {
		processes = append(processes, process)
	}
	for _, process := range processes {
		c.doProcessExit(process)
	}
}
