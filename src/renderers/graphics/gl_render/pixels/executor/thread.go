package executor

import (
	"runtime"
)

// queueCap defines the capacity for channel buffers used to manage thread-safe function execution queues.
const (
	queueCap = 32
)

// GraphicThread is a global instance of *Thread designed to handle graphics-related tasks on a dedicated thread.
var GraphicThread *Thread

// init initializes the main graphic thread by locking it to the OS thread and creating an instance of GraphicThread.
func init() {
	runtime.LockOSThread()
	GraphicThread = NewThread()
}

// Thread represents a struct designed to manage concurrent execution of tasks using separate post and call queues.
// It provides a mechanism to post and synchronously or asynchronously execute functions in a controlled environment.
type Thread struct {
	done          chan bool
	postQueue     chan func()
	callQueue     chan func()
	callWaitQueue chan bool
}

// NewThread creates a new Thread instance with initialized channels for communication and queue handling.
func NewThread() *Thread {
	return &Thread{
		done:          make(chan bool),
		postQueue:     make(chan func(), queueCap),
		callQueue:     make(chan func(), queueCap),
		callWaitQueue: make(chan bool),
	}
}

// Post schedules the provided function to be executed within the Thread's context.
func (m *Thread) Post(f func()) {
	m.postQueue <- f
}

// Call schedules a function to be executed on the thread and blocks until the function has finished execution.
func (m *Thread) Call(f func()) {
	m.callQueue <- f
	<-m.callWaitQueue
}

// CallVal executes a function synchronously on the thread and returns its resulting value.
func (m *Thread) CallVal(f func() interface{}) interface{} {
	var val interface{}
	fnCq := func() {
		val = f()
	}
	m.callQueue <- fnCq
	<-m.callWaitQueue
	return val
}

// CallErr queues a function returning an error in the callQueue, waits for its execution, and returns the resulting error.
func (m *Thread) CallErr(f func() error) error {
	var err error
	fnCq := func() {
		err = f()
	}
	m.callQueue <- fnCq
	<-m.callWaitQueue
	return err
}

// Run starts the thread's primary loop, executing the provided main function and managing event handling. Returns error if any.
func (m *Thread) Run(main func()) error {
	m.mainLoop(main)
	m.eventLoop()
	return nil
}

// mainLoop initializes the main thread function and executes it in a separate goroutine.
// Signals thread completion via the done channel after the main function execution finishes.
func (m *Thread) mainLoop(main func()) {
	go func() {
		main()
		m.done <- true
	}()
}

// eventLoop orchestrates and executes functions from the postQueue and callQueue, managing the thread lifecycle until done.
func (m *Thread) eventLoop() {
	defer close(m.done)
	defer close(m.postQueue)
	defer close(m.callQueue)
	defer close(m.callWaitQueue)

	for {
		select {
		case fnCq := <-m.postQueue:
			fnCq()
		case fnCq := <-m.callQueue:
			fnCq()
			m.callWaitQueue <- true
		case <-m.done:
			return
		}
	}
}
