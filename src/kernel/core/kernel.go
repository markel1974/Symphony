package core

// Package core implements the kernel subsystem responsible for orchestrating process lifecycle,
// message routing, timer management, and I/O operations within the system environment.

import (
	"fmt"
	"log"

	"github.com/markel1974/c64emu/src/kernel/adaptiveticker"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/messages"
	"github.com/markel1974/c64emu/src/kernel/process"
	"github.com/markel1974/c64emu/src/kernel/process_factory"
)

const (
	// kernelQueueLen defines the buffer capacity for the kernel's message channel
	kernelQueueLen = 8192
	// kernelQueueMax represents the maximum safe queue length to prevent overflow
	kernelQueueMax = kernelQueueLen - 1
)

// Kernel serves as the central orchestrator for the system, managing process execution,
// message routing, timer operations, and I/O handling. It maintains the system's state
// and coordinates communication between all active components.
type Kernel struct {
	user         string                                               // Current user context
	ticker       *adaptiveticker.AdaptiveTicker                       // System timer manager
	inputDriver  interfaces.IKeyboardDriver                           // Keyboard input handler
	foreground   interfaces.IUserProcess                              // Currently active foreground process
	pidGenerator *adaptiveticker.Ids                                  // Process ID allocation manager
	running      map[int]*KernelProcess                               // Registry of active kernel processes
	shellPath    string                                               // Default shell executable path
	messageChan  chan interfaces.IMessage                             // Main message processing queue
	pf           *process_factory.ProcessFactory                      // Factory for creating new processes
	timersChan   chan *adaptiveticker.TimerHandler                    // Timer event delivery channel
	servers      []interfaces.IServer                                 // Registered system servers
	exit         bool                                                 // System shutdown flag
	routingTable map[interfaces.MessageType]func(interfaces.IMessage) // Message dispatch routing table
	process      *KernelProcess                                       // The kernel's own process representation
}

// NewKernel constructs a new kernel instance with the specified user context, timer system,
// input driver, and shell path. It initializes all internal data structures and prepares
// the kernel for operation without starting the main event loop.
func NewKernel(user string, ticker *adaptiveticker.AdaptiveTicker, inputDriver interfaces.IKeyboardDriver, shellPath string) *Kernel {
	t := &Kernel{
		user:         user,
		ticker:       ticker,
		inputDriver:  inputDriver,
		foreground:   nil,
		pidGenerator: adaptiveticker.NewIds(1024),
		messageChan:  make(chan interfaces.IMessage, kernelQueueLen),
		timersChan:   make(chan *adaptiveticker.TimerHandler, kernelQueueLen),
		exit:         false,
		shellPath:    shellPath,
		running:      make(map[int]*KernelProcess),
		routingTable: make(map[interfaces.MessageType]func(interfaces.IMessage)),
		pf:           process_factory.NewProcessFactory(),
	}
	return t
}

// PID returns the kernel's own process identifier, establishing it as a first-class
// participant in the system's process hierarchy.
func (c *Kernel) PID() int {
	return c.process.PID()
}

// Setup initializes the kernel's operational infrastructure by configuring message routing,
// registering system servers, and creating the kernel's own process context. This method
// establishes the foundation for all subsequent kernel operations.
func (c *Kernel) Setup(servers []interfaces.IServer) error {
	// Configure core message routing table
	c.routingTable[interfaces.MessageTypeRead] = c.handleReadEvent
	c.routingTable[interfaces.MessageTypeTimer] = c.handleTimerEvent
	c.routingTable[interfaces.MessageTypeQuit] = c.handleQuitEvent
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
	c.routingTable[interfaces.MessageTypeProcessIsRunning] = c.handleProcessIsRunning
	c.routingTable[interfaces.MessageTypeExitRequested] = c.handleExitRequested

	// Register system servers and their message handlers
	for _, server := range servers {
		c.servers = append(c.servers, server)
		for _, h := range server.Register() {
			c.routingTable[h] = server.PostUserMessage
		}
	}

	// Initialize kernel's own process context
	var err error
	kCmd := c.doCreateKernelCommand("kernel")
	c.process, err = c.doCreateKernelProcess(c.user, kCmd.Name(), nil, kCmd, true)
	if err != nil {
		return err
	}

	// Initialize server processes
	for _, server := range c.servers {
		sCmd := c.doCreateKernelCommand(server.Name())
		serverProcess, err := c.doCreateKernelProcess(c.user, sCmd.Name(), c.process, sCmd, true)
		if err != nil {
			return err
		}
		if err = server.Setup(serverProcess, serverProcess); err != nil {
			return err
		}
	}
	return nil
}

// Start activates the kernel's main execution loop, beginning asynchronous I/O processing
// and entering the primary event handling cycle. This method transitions the kernel
// from initialization to active operation state.
func (c *Kernel) Start() error {
	// Request shell initialization
	c.postSelfMessage(messages.NewMessageFileSystemFindRequest(c.PID(), c.shellPath, true))

	// Launch asynchronous input processing
	d := make(chan bool)
	go func() {
		d <- true
		readBuffer := make([]byte, 4096)
		for {
			k, v, err := c.inputDriver.ScanKey(readBuffer)
			if err == nil {
				if k != interfaces.KeyTypeNone {
					re := messages.NewMessageRead(k, v, false)
					c.postSelfMessage(re)
				}
			} else {
				qe := messages.NewMessageQuit()
				c.postSelfMessage(qe)
				return
			}
		}
	}()
	_ = <-d
	c.eventLoop()
	return nil
}

// SetScreenSize notifies the system of display dimension changes, allowing processes
// to adapt their output accordingly.
func (c *Kernel) SetScreenSize(w int, h int) {
	msg := messages.NewMessageSetScreenSize(w, h)
	c.postSelfMessage(msg)
}

// Process provides access to the kernel's process representation, enabling it to
// participate in standard process operations and communication patterns.
func (c *Kernel) Process() interfaces.IUserProcess {
	return c.process
}

// User returns the current user context under which the kernel operates,
// establishing the security and permission context for all operations.
func (c *Kernel) User() string {
	return c.user
}

// postSelfMessage queues a message for processing by the kernel's event loop, implementing
// flow control to prevent queue overflow and maintain system stability.
func (c *Kernel) postSelfMessage(msg interfaces.IMessage) {
	if len(c.messageChan) >= kernelQueueMax {
		log.Printf("Kernel: message queue full, dropping message: %d", msg.GetType())
		return
	}
	msg.SetPID(c.PID())
	c.messageChan <- msg
}

// handleProcessIsRunning verifies the active status of a specified process and returns
// the result to the requesting process, supporting process management and coordination.
func (c *Kernel) handleProcessIsRunning(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageProcessIsRunning)
	if !ok {
		return
	}
	kProc, _ := c.running[mt.PID()]
	if kProc == nil {
		return
	}
	kActive, _ := c.running[mt.VerifyPID()]
	mt.SetResult(kActive != nil)
	kProc.PostUserMessage(mt)
}

// handleExitRequested processes system shutdown requests by setting the exit flag,
// initiating the kernel's graceful shutdown sequence.
func (c *Kernel) handleExitRequested(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageExitRequested)
	if !ok {
		return
	}
	kProc, _ := c.running[mt.PID()]
	if kProc == nil {
		return
	}
	c.exit = true
}

// eventLoop implements the kernel's main execution cycle, processing messages and timer
// events while monitoring for shutdown conditions. This is the heart of kernel operation.
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

// handleMessageEvent dispatches incoming messages to their appropriate handlers based
// on message type and process routing tables, forming the core of the kernel's
// message processing architecture.
func (c *Kernel) handleMessageEvent(m interfaces.IMessage) {
	log.Printf("Kernel: message received: %d - %d", m.PID(), m.GetType())

	kProc, _ := c.running[m.PID()]
	if kProc == nil {
		log.Printf("Kernel: unknown process: %d - type %d", m.PID(), m.GetType())
		return
	}

	//if m.PID() < 0 {
	//	m.SetPID(kProc.PID())
	//}

	// Handle response messages directly
	if m.Response() {
		kProc.PostUserMessage(m)
		return
	}

	// Route messages based on process-specific or kernel routing tables
	id := m.GetType()
	if route := kProc.GetRoute(id); route != nil {
		route(m)
		return
	}
	log.Printf("Kernel: unknown message type: %d", id)
}

// handleReadEvent processes keyboard input by distributing events to interested processes,
// supporting both broadcast and targeted delivery mechanisms for flexible I/O handling.
func (c *Kernel) handleReadEvent(m interfaces.IMessage) {
	mm, ok := m.(*messages.MessageRead)
	if !ok {
		return
	}

	// Broadcast to processes requesting global input notifications
	//sentForeground := false
	for _, kProc := range c.running {
		if readBroadcastEvent := kProc.GetCommand().OnReadBroadcast(); readBroadcastEvent != nil {
			kProc.PostUserMessage(messages.NewMessageRead(mm.Kind(), mm.Data(), true))
			//if c.foreground != nil && c.foreground == kProc{
			//sentForeground = true
			//}
		}
	}

	// Send to foreground process if it accepts input
	if c.foreground != nil {
		if readEvent := c.foreground.GetCommand().OnRead(); readEvent != nil {
			c.foreground.PostUserMessage(mm)
		}
	}
}

// handleTimerEvent delivers timer notifications to the target process, enabling
// time-based operations and periodic task execution within the system.
func (c *Kernel) handleTimerEvent(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageTimer)
	if !ok {
		return
	}
	kProc, _ := c.running[mt.PID()]
	if kProc == nil {
		return
	}
	kProc.PostUserMessage(mt)
}

// handleProcessExit manages process termination by performing cleanup operations
// and updating system state to reflect the process's removal from active execution.
func (c *Kernel) handleProcessExit(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageProcessExit)
	if !ok {
		return
	}
	kProc, _ := c.running[mt.PID()]
	if kProc == nil {
		log.Printf("Kernel [handleProcessExit]: unknown process: %d - type %d", mt.PID(), mt.GetType())
		return
	}
	c.doProcessExit(kProc)
}

// handleProcessExec initiates process execution by resolving the command through the
// filesystem service, preparing the execution environment for the new process.
func (c *Kernel) handleProcessExec(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageProcessExec)
	if !ok {
		return
	}
	c.postSelfMessage(messages.NewMessageFileSystemFindRequest(mt.PID(), mt.Line(), false))
}

// handleProcessSetForeground manages foreground process transitions, ensuring proper
// focus management and user interface responsiveness within the system.
func (c *Kernel) handleProcessSetForeground(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageProcessSetForeground)
	if !ok {
		return
	}
	kProc, _ := c.running[mt.PID()]
	if kProc == nil {
		return
	}
	c.doProcessSetForeground(mt.PID(), kProc)
}

// handleProcessKill terminates a specific process by its identifier, providing controlled
// process lifecycle management and resource cleanup capabilities.
func (c *Kernel) handleProcessKill(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageProcessKill)
	if !ok {
		return
	}
	kProc, _ := c.running[mt.PID()]
	if kProc == nil {
		return
	}
	c.doProcessExit(kProc)
}

// handleProcessKillAll terminates multiple processes based on optional name filtering,
// supporting bulk process management operations for system maintenance tasks.
func (c *Kernel) handleProcessKillAll(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageProcessKillAll)
	if !ok {
		return
	}
	var processes []*KernelProcess
	for _, kProc := range c.running {
		if len(mt.Name()) != 0 {
			if kProc.GetCommand().Name() != mt.Name() {
				continue
			}
		}
		processes = append(processes, kProc)
	}
	for _, proc := range processes {
		c.doProcessExit(proc)
	}
}

// handleProcessKillForeground terminates the currently active foreground process,
// typically used for interrupt handling and user-initiated task cancellation.
func (c *Kernel) handleProcessKillForeground(m interfaces.IMessage) {
	_, ok := m.(*messages.MessageProcessKillForeground)
	if !ok {
		return
	}
	kProc, _ := c.running[c.foreground.PID()]
	if kProc == nil {
		return
	}
	c.doProcessExit(kProc)
}

// handleTimerCreate processes timer creation requests, establishing periodic or delayed
// execution contexts and associating them with the requesting process.
func (c *Kernel) handleTimerCreate(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageTimerCreate)
	if !ok {
		return
	}
	kProc, _ := c.running[m.PID()]
	if kProc == nil {
		return
	}
	msgTimer := messages.NewMessageTimer(mt.PID(), mt.Interval())
	msgTimer.SetTID(c.ticker.Create(c.timersChan, msgTimer, int64(mt.First()), int64(mt.Interval()), int64(mt.Count())))
	if msgTimer.TID() < 0 {
		kProc.PostUserMessage(messages.NewMessageError(fmt.Errorf("error creating timer")))
		return
	}
	kProc.AddTimer(msgTimer.TID())
	kProc.PostUserMessage(messages.NewMessageTimerCreated(msgTimer.TID()))
}

// handleTimerStop processes timer cancellation requests, removing active timers
// and cleaning up associated resources to prevent memory leaks and unnecessary processing.
func (c *Kernel) handleTimerStop(m interfaces.IMessage) {
	mt, ok := m.(*messages.MessageTimerStop)
	if !ok {
		return
	}
	kProc, _ := c.running[m.PID()]
	if kProc == nil {
		log.Printf("error stopping timer: invalid originator %d", m.PID())
		return
	}
	c.doCloseTimer(kProc, mt.TID())
}

// handleProcessList generates and returns a comprehensive list of active processes,
// supporting system monitoring and process management functionality.
func (c *Kernel) handleProcessList(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageProcessList)
	if !ok {
		return
	}
	kProc, _ := c.running[mt.PID()]
	if kProc == nil {
		log.Printf("error in handleProcessList: invalid originator %d", mt.PID())
		return
	}
	var out []*interfaces.ProcessDescription
	for _, proc := range c.running {
		out = append(out, proc.Description())
	}
	mt.SetResult(out)
	kProc.PostUserMessage(mt)
}

// handleQuitEvent processes system quit requests by setting the exit flag,
// initiating the kernel's graceful shutdown sequence and resource cleanup.
func (c *Kernel) handleQuitEvent(m interfaces.IMessage) {
	_, ok := m.(*messages.MessageQuit)
	if !ok {
		return
	}
	c.exit = true
}

// handleFileSystemFindResponse completes the process execution cycle by creating
// and initializing new processes based on filesystem command resolution results.
func (c *Kernel) handleFileSystemFindResponse(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageFileSystemFindResponse)
	if !ok {
		return
	}
	kRequestor, _ := c.running[mt.RequestorPID()]
	if kRequestor == nil {
		log.Printf("error creating task: invalid originator")
		return
	}
	cmd, args, err := mt.GetResult()
	if err != nil {
		kRequestor.PostUserMessage(messages.NewMessageError(fmt.Errorf("error creating task: invalid command '%s'", mt.Line())))
		return
	}
	parent, _ := c.running[kRequestor.PID()]
	kProc, err := c.doCreateKernelProcess(kRequestor.User(), mt.Line(), parent, cmd, mt.Protected())
	if err != nil {
		kRequestor.PostUserMessage(messages.NewMessageError(fmt.Errorf("error creating task: %s", err.Error())))
		return
	}
	kProc.Start()
	for _, server := range c.servers {
		server.PostUserMessage(messages.NewMessageNotifyProcessCreate(kProc.PID(), kProc.GetCommand().Name()))
	}
	kProc.PostUserMessage(messages.NewMessageProcessStart(args))
}

// doProcessSetForeground updates the system's foreground process state and notifies
// relevant system components of the focus change, maintaining UI consistency.
func (c *Kernel) doProcessSetForeground(requestorPID int, process interfaces.IUserProcess) {
	for _, server := range c.servers {
		server.PostUserMessage(messages.NewMessageNotifyProcessForeground(process.PID()))
	}
	if c.foreground != process {
		c.foreground = process
		c.foreground.PostUserMessage(messages.NewMessageProcessActivate())
	}
}

// doProcessExit orchestrates complete process termination including resource cleanup,
// timer cancellation, server notification, and foreground process management.
func (c *Kernel) doProcessExit(process *KernelProcess) {
	if process.Protected() {
		return
	}
	if len(process.Timers()) > 0 {
		c.ticker.RemoveEntries(process.Timers())
	}
	for _, server := range c.servers {
		server.PostUserMessage(messages.NewMessageMessageNotifyProcessTerminate(process.PID()))
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
	process.PostUserMessage(messages.NewMessageQuit())
}

// doCloseTimer removes a specific timer from both the process context and the system
// timer manager, ensuring complete timer lifecycle management and resource cleanup.
func (c *Kernel) doCloseTimer(kProc *KernelProcess, tid int) bool {
	ret := false
	if kProc != nil {
		kProc.TimersIterator(func(timerId int) bool {
			if timerId == tid {
				ret = c.ticker.RemoveEntries([]int{timerId})
				return true
			}
			return false
		})
	}
	return ret
}

// doShutdown performs comprehensive system shutdown by terminating all active processes
// and cleaning up system resources, ensuring graceful system state transition.
func (c *Kernel) doShutdown() {
	var processes []*KernelProcess
	for _, kProc := range c.running {
		processes = append(processes, kProc)
	}
	for _, kProc := range processes {
		c.doProcessExit(kProc)
	}
}

// doCreateKernelCommand constructs basic kernel command objects for internal system
// processes, providing the minimal command interface required for kernel operation.
func (c *Kernel) doCreateKernelCommand(name string) interfaces.ICommand {
	onCreate := func(process interfaces.IUserProcess, args []string) error { return nil }
	kCmd := process.NewCommand(name, interfaces.CommandTypeFile, nil, true, onCreate)
	return kCmd
}

// doCreateKernelProcess creates fully configured kernel processes with proper PID
// allocation, routing table setup, and registration in the kernel's process registry.
func (c *Kernel) doCreateKernelProcess(user string, commandLine string, parent *KernelProcess, cmd interfaces.ICommand, protected bool) (*KernelProcess, error) {
	pid := NewPID()
	if _, ok := c.pidGenerator.Set(pid); !ok {
		return nil, fmt.Errorf("error creating task: pid generator overflow")
	}
	routingTable := make(map[interfaces.MessageType]func(interfaces.IMessage))
	for k, v := range c.routingTable {
		routingTable[k] = v
	}
	uProc := c.pf.Create(pid.GetId(), user, cmd)
	kProc := NewKernelProcess(uProc, user, commandLine, cmd.Name(), parent, pid, protected, c.messageChan)
	kProc.SetRoutingTable(routingTable)
	c.running[kProc.PID()] = kProc
	return kProc, nil
}
