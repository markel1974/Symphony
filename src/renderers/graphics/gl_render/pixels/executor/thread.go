package executor

import (
	"runtime"
)

const (
	queueCap = 32
)

var GraphicThread *Thread

func init() {
	runtime.LockOSThread()
	GraphicThread = NewThread()
}

type Thread struct {
	done          chan bool
	postQueue     chan func()
	callQueue     chan func()
	callWaitQueue chan bool
}

func NewThread() *Thread {
	return &Thread{
		done:          make(chan bool),
		postQueue:     make(chan func(), queueCap),
		callQueue:     make(chan func(), queueCap),
		callWaitQueue: make(chan bool),
	}
}

func (m *Thread) Post(f func()) {
	m.postQueue <- f
}

func (m *Thread) Call(f func()) {
	m.callQueue <- f
	<-m.callWaitQueue
}

func (m *Thread) CallVal(f func() interface{}) interface{} {
	var val interface{}
	fnCq := func() {
		val = f()
	}
	m.callQueue <- fnCq
	<-m.callWaitQueue
	return val
}

func (m *Thread) CallErr(f func() error) error {
	var err error
	fnCq := func() {
		err = f()
	}
	m.callQueue <- fnCq
	<-m.callWaitQueue
	return err
}

func (m *Thread) Run(main func()) error {
	m.mainLoop(main)
	m.eventLoop()
	return nil
}

func (m *Thread) mainLoop(main func()) {
	go func() {
		main()
		m.done <- true
	}()
}

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
