package core

import (
	"fmt"
	"log"

	"github.com/markel1974/c64emu/src/kernel/adaptiveticker"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/messages"
	"github.com/markel1974/c64emu/src/kernel/process"
	"github.com/markel1974/c64emu/src/kernel/process_factory"
)

// Kernel represents the core component responsible for managing rendering, input/output, process execution, and timers.
type Kernel struct {
	user         string
	ticker       *adaptiveticker.AdaptiveTicker
	inputDriver  interfaces.IKeyboardDriver
	foreground   interfaces.IProcess
	pidGenerator *adaptiveticker.Ids
	running      map[int]*KernelProcess
	shellPath    string
	messageChan  chan interfaces.IMessage
	pf           *process_factory.ProcessFactory
	timersChan   chan *adaptiveticker.TimerHandler
	servers      []interfaces.IServer
	exit         bool
	routingTable map[interfaces.MessageType]func(interfaces.IMessage)
	process      *KernelProcess
}

// NewKernel creates and returns a new Kernel instance, initializing its dependencies and internal fields.
func NewKernel(user string, ticker *adaptiveticker.AdaptiveTicker, inputDriver interfaces.IKeyboardDriver, shellPath string) *Kernel {
	t := &Kernel{
		user:         user,
		ticker:       ticker,
		inputDriver:  inputDriver,
		foreground:   nil,
		pidGenerator: adaptiveticker.NewIds(1024),
		messageChan:  make(chan interfaces.IMessage, contextMaQueueLen),
		timersChan:   make(chan *adaptiveticker.TimerHandler, contextMaQueueLen),
		exit:         false,
		shellPath:    shellPath,
		running:      make(map[int]*KernelProcess),
		routingTable: make(map[interfaces.MessageType]func(interfaces.IMessage)),
	}
	t.pf = process_factory.NewProcessFactory(t)

	return t
}

// PID returns the current process ID (pid) of the Kernel instance.
func (c *Kernel) PID() int {
	return c.process.PID()
}

func (c *Kernel) createKernelProcess(user string, commandLine string, parent *KernelProcess, cmd interfaces.ICommand, protected bool) (*KernelProcess, error) {
	pid := NewPID()
	if _, ok := c.pidGenerator.Set(pid); !ok {
		return nil, fmt.Errorf("error creating task: pid generator overflow")
	}
	userProcess := c.pf.Create(pid.GetId(), user, cmd)
	kernelProcess := NewKernelProcess(user, commandLine, cmd.Name(), parent, pid, protected, userProcess)
	routingTable := make(map[interfaces.MessageType]func(interfaces.IMessage))
	for k, v := range c.routingTable {
		routingTable[k] = v
	}
	kernelProcess.SetRoutingTable(routingTable)
	c.running[kernelProcess.PID()] = kernelProcess
	return kernelProcess, nil
}

func (c *Kernel) Setup(servers []interfaces.IServer) error {
	//routing table
	c.routingTable[interfaces.MessageTypeRead] = c.handleReadEvent
	c.routingTable[interfaces.MessageTypeTimer] = c.handleTimerEvent
	c.routingTable[interfaces.MessageTypeQuit] = c.handleQuitEvent
	c.routingTable[interfaces.MessageTypeTimedMessage] = c.handleTimedMessage
	c.routingTable[interfaces.MessageTypeProcessExit] = c.handleProcessExit
	c.routingTable[interfaces.MessageTypeProcessExec] = c.handleProcessExec
	c.routingTable[interfaces.MessageTypeProcessSetForeground] = c.handleProcessSetForeground
	c.routingTable[interfaces.MessageTypeProcessKill] = c.handleProcessKill
	c.routingTable[interfaces.MessageTypeProcessKillAll] = c.handleProcessKillAll
	c.routingTable[interfaces.MessageTypeProcessKillForeground] = c.handleProcessKillForeground
	c.routingTable[interfaces.MessageTypeTimerCreate] = c.handleTimerCreate
	c.routingTable[interfaces.MessageTypeTimerStop] = c.handleTimerStop
	c.routingTable[interfaces.MessageTypeProcessList] = c.handleProcessList
	c.routingTable[interfaces.MessageTypeFileSystemFindResponse] = c.handleFileSystemFindResponse
	for _, server := range servers {
		c.servers = append(c.servers, server)
		for _, h := range server.Register(c) {
			c.routingTable[h] = server.PostMessage
		}
	}
	//kernel process
	onCreate := func(process interfaces.IProcess, args []string) error { return nil }
	kernelCmd := process.NewCommand("kernel", interfaces.CommandTypeFile, nil, true, onCreate)
	kernelProcess, err := c.createKernelProcess(c.user, kernelCmd.Name(), nil, kernelCmd, true)
	if err != nil {
		return err
	}
	c.process = kernelProcess
	//server process
	for _, server := range c.servers {
		serverCmd := process.NewCommand(server.Name(), interfaces.CommandTypeFile, nil, true, func(process interfaces.IProcess, args []string) error { return nil })
		serverProcess, err := c.createKernelProcess(c.user, server.Name(), c.process, serverCmd, true)
		if err != nil {
			return err
		}
		server.SetProcess(serverProcess)
		server.Start()
	}
	return nil
}

// Start initializes the kernel's event handling loop and begins processing I/O operations asynchronously.
func (c *Kernel) Start() error {
	c.PostMessage(messages.NewMessageFileSystemFindRequest(c.PID(), c.PID(), c.shellPath, true))

	d := make(chan bool)
	go func() {
		d <- true
		readBuffer := make([]byte, 4096)
		for {
			k, v, err := c.inputDriver.ScanKey(readBuffer)
			if err == nil {
				if k != interfaces.KeyTypeNone {
					re := messages.NewMessageRead(c.PID(), k, v, false)
					c.messageChan <- re
				}
			} else {
				qe := messages.NewMessageQuit(c.PID())
				c.messageChan <- qe
				return
			}
		}
	}()
	_ = <-d
	c.eventLoop()
	return nil
}

// SetScreenSize adjusts the screen dimensions to the specified width and height values.
func (c *Kernel) SetScreenSize(w int, h int) {
	c.messageChan <- messages.NewMessageSetScreenSize(c.PID(), w, h)
}

// Process returns the KernelProcess instance associated with the Kernel instance.
func (c *Kernel) Process() interfaces.IProcess {
	return c.process
}

// User returns the name of the user associated with the Kernel instance.
func (c *Kernel) User() string {
	return c.user
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

// CallExitRequested sets the `exit` flag to true, signaling that an exit has been requested for the kernel.
func (c *Kernel) CallExitRequested(router interfaces.IRouter) {
	c.exit = true
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

// handleMessageEvent processes an incoming IMessage by dispatching it to the appropriate routingTable based on its type.
func (c *Kernel) handleMessageEvent(m interfaces.IMessage) {
	kernelProcess, _ := c.running[m.PID()]
	if kernelProcess == nil {
		log.Printf("Kernel: unknown process: %d - type %d", m.PID(), m.GetType())
		return
	}
	//TODO IMPLEMENTARE TUTTI I CONTROLLI E I LOG
	//fmt.Println("Kernel: dispatching", m.GetType(), "")
	if m.Response() {
		kernelProcess.PostMessage(m)
		return
	}
	id := m.GetType()
	if route := kernelProcess.GetRoute(id); route != nil {
		route(m)
		return
	}
	log.Printf("Kernel: unknown message type: %d", id)
}

// handleReadEvent processes input events based on their type and key value to handle control, foreground processes, and system state.
func (c *Kernel) handleReadEvent(m interfaces.IMessage) {
	mm, ok := m.(*messages.MessageRead)
	if !ok {
		return
	}
	for _, proc := range c.running {
		if readBroadcastEvent := proc.GetCommand().OnReadBroadcast(); readBroadcastEvent != nil {
			proc.PostMessage(messages.NewMessageRead(m.PID(), mm.Kind(), mm.Data(), true))
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
	proc, _ := c.running[mt.PID()]
	if proc == nil {
		return
	}
	proc.PostMessage(mt)
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
		log.Printf("Kernel [handleProcessExit]: unknown type: %d - process %d", mt.GetType(), mt.PID())
		return
	}
	kernelProcess, _ := c.running[mt.PID()]
	if kernelProcess == nil {
		log.Printf("Kernel [handleProcessExit]: unknown process: %d - type %d", mt.PID(), mt.GetType())
		return
	}
	c.doProcessExit(kernelProcess)
}

// handleProcessExec handles process execution by validating the message type and invoking the process execution logic.
func (c *Kernel) handleProcessExec(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageProcessExec)
	if !ok {
		return
	}
	c.PostMessage(messages.NewMessageFileSystemFindRequest(c.PID(), mt.PID(), mt.Line(), false))
}

// handleProcessSetForeground handles a process set foreground message by setting the foreground process to the specified process.
func (c *Kernel) handleProcessSetForeground(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageProcessSetForeground)
	if !ok {
		return
	}
	if proc, _ := c.running[mt.PID()]; proc != nil {
		c.doProcessSetForeground(mt.PID(), proc)
	}
}

// handleProcessKill terminates and removes a process by its process ID (pid). Returns true if successful, false if the pid is not found.
func (c *Kernel) handleProcessKill(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageProcessKill)
	if !ok {
		return
	}
	proc, _ := c.running[mt.PID()]
	if proc == nil {
		return
	}
	c.doProcessExit(proc)
}

// handleProcessKillAll handles the termination of all processes except the sender's and optionally filters by process name.
func (c *Kernel) handleProcessKillAll(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageProcessKillAll)
	if !ok {
		return
	}
	var processes []*KernelProcess
	for _, proc := range c.running {
		if len(mt.Name()) != 0 {
			if proc.GetCommand().Name() != mt.Name() {
				continue
			}
		}
		processes = append(processes, proc)
	}
	for _, proc := range processes {
		c.doProcessExit(proc)
	}
}

// handleProcessKillForeground handles the termination of the foreground process.
func (c *Kernel) handleProcessKillForeground(m interfaces.IMessage) {
	_, ok := m.(*messages.MessageProcessKillForeground)
	if !ok {
		return
	}
	proc, _ := c.running[c.foreground.PID()]
	if proc == nil {
		return
	}
	c.doProcessExit(proc)
}

// handleTimerCreate processes a timer creation request, initializes the timer, and sends a response or error message.
func (c *Kernel) handleTimerCreate(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageTimerCreate)
	if !ok {
		return
	}
	proc, _ := c.running[m.PID()]
	if proc == nil {
		return
	}
	msgTimer := messages.NewMessageTimer(mt.PID(), mt.PID(), mt.Interval())
	msgTimer.SetTID(c.ticker.Create(c.timersChan, msgTimer, int64(mt.First()), int64(mt.Interval()), int64(mt.Count())))
	if msgTimer.TID() < 0 {
		proc.PostMessage(messages.NewMessageError(m.PID(), fmt.Errorf("error creating timer")))
		return
	}
	proc.AddTimer(msgTimer.TID())
	proc.PostMessage(messages.NewMessageTimerCreated(m.PID(), msgTimer.TID()))
}

func (c *Kernel) handleTimerStop(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageTimerStop)
	if !ok {
		return
	}
	proc, _ := c.running[m.PID()]
	if proc == nil {
		log.Printf("error stopping time: invalid originator %d", m.PID())
		return
	}
	c.doCloseTimer(proc, mt.TID())
}

// handleProcessList processes a message requesting a list of running processes and sends the response with process details.
func (c *Kernel) handleProcessList(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageProcessList)
	if !ok {
		return
	}
	kernelProc, _ := c.running[mt.PID()]
	if kernelProc == nil {
		log.Printf("error in handleProcessList: invalid originator %d", mt.PID())
		return
	}
	var out []*interfaces.ProcessDescription
	for _, proc := range c.running {
		out = append(out, proc.Description())
	}
	mt.SetProcesses(out)
	kernelProc.PostMessage(mt)
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
	originator, _ := c.running[mt.RequestorPID()]
	if originator == nil {
		log.Printf("error creating task: invalid originator")
		return
	}
	cmd, args, err := mt.GetResult()
	if err != nil {
		originator.PostMessage(messages.NewMessageError(originator.PID(), fmt.Errorf("error creating task: invalid command '%s'", mt.Line())))
		return
	}
	parent, _ := c.running[originator.PID()]
	kernelProcess, err := c.createKernelProcess(originator.User(), mt.Line(), parent, cmd, mt.Protected())
	if err != nil {
		originator.PostMessage(messages.NewMessageError(originator.PID(), fmt.Errorf("error creating task: %s", err.Error())))
		return
	}
	kernelProcess.Setup()
	for _, server := range c.servers {
		server.NotifyProcessCreation(kernelProcess.PID(), kernelProcess.GetCommand().Name())
	}
	kernelProcess.PostMessage(messages.NewMessageProcessStart(originator.PID(), args))
}

// doProcessSetForeground sets the specified process as the foreground process and sends activation messages if needed.
func (c *Kernel) doProcessSetForeground(requestorPID int, process interfaces.IProcess) {
	for _, s := range c.servers {
		s.NotifyProcessForeground(process.PID())
	}
	if c.foreground != process {
		c.foreground = process
		c.foreground.PostMessage(messages.NewMessageProcessActivate(requestorPID))
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
				c.doProcessSetForeground(c.PID(), parent)
			} else {
				log.Printf("Fatal Error: foreground process is nil")
			}
		}
	}
	delete(c.running, process.PID())
	c.pidGenerator.Unset(process.PID())
	process.PostMessage(messages.NewMessageQuit(process.PID()))
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
	for _, proc := range c.running {
		processes = append(processes, proc)
	}
	for _, proc := range processes {
		c.doProcessExit(proc)
	}
}
