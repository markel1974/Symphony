package context

import "github.com/markel1974/c64emu/src/shell/adaptiveticker"

// TaskSelector tracks the state for task selection, including process ID, available tasks, and current index.
type TaskSelector struct {
	pid       int
	available []int
	idx       int
}

// NewTaskSelector initializes and returns a new instance of TaskSelector with default values.
func NewTaskSelector() *TaskSelector {
	return &TaskSelector{
		pid:       adaptiveticker.UnknownId,
		available: nil,
		idx:       0,
	}
}

func (ts *TaskSelector) Clear() {
	ts.idx = 0
	ts.pid = adaptiveticker.UnknownId
	ts.available = nil
}

func (ts *TaskSelector) AddAvailable(pid int) {
	ts.available = append(ts.available, pid)
}

func (ts *TaskSelector) Set(pid int, idx int) {
	ts.pid = pid
	ts.idx = idx
}

func (ts *TaskSelector) Next() bool {
	if len(ts.available) == 0 {
		return false
	}
	next := ts.idx + 1
	if next >= len(ts.available) {
		next = 0
	}
	ts.idx = next
	ts.pid = ts.available[next]
	return true
}

func (ts *TaskSelector) Prev() bool {
	if len(ts.available) == 0 {
		return false
	}
	prev := ts.idx - 1
	if prev < 0 {
		prev = len(ts.available) - 1
	}
	ts.idx = prev
	ts.pid = ts.available[prev]
	return true
}

func (ts *TaskSelector) PID() int {
	return ts.pid
}
