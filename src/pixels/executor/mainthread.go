package executor

import (
	"errors"
	"runtime"
	"sync"
)

const (
	callQueueCap = 16
)

var Thread *MainThread

func init() {
	runtime.LockOSThread()
	Thread = NewMainThread()
}

type MainThread struct {
	callQueue chan func()
	respMutex sync.Mutex
	respChan  chan interface{}
}

func NewMainThread() *MainThread {
	return &MainThread{
		callQueue: make(chan func(), callQueueCap),
		respChan:  make(chan interface{}),
	}
}

func (m *MainThread) Run(fn func()) {
	done := make(chan bool)
	go func() {
		fn()
		done <- true
	}()

	for {
		select {
		case fnCq := <-m.callQueue:
			fnCq()
		case <-done:
			return
		}
	}
}

func (m *MainThread) Post(f func()) {
	m.callQueue <- f
}

func (m *MainThread) Call(f func()) {
	m.respMutex.Lock()
	m.callQueue <- func() {
		f()
		m.respChan <- true
	}
	<-m.respChan
	m.respMutex.Unlock()
}

func (m *MainThread) CallErr(f func() error) error {
	m.respMutex.Lock()
	m.callQueue <- func() {
		m.respChan <- f()
	}
	resp := <-m.respChan
	m.respMutex.Unlock()
	if resp == nil {
		return nil
	}
	if err, ok := resp.(error); ok {
		return err
	}
	return errors.New("invalid response")
}

func (m *MainThread) CallVal(f func() interface{}) interface{} {
	m.respMutex.Lock()
	m.callQueue <- func() {
		m.respChan <- f()
	}
	val := <-m.respChan
	m.respMutex.Unlock()
	return val
}
