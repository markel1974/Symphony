package context

import "github.com/markel1974/c64emu/src/shell/adaptiveticker"

// TaskSelector is a type that manages available task identifiers and tracks the current one for selection and navigation.
type TaskSelector struct {
	pid       int
	available []int
	idx       int
}

// NewTaskSelector creates and returns a new TaskSelector initialized with default values.
func NewTaskSelector() *TaskSelector {
	return &TaskSelector{
		pid:       adaptiveticker.UnknownId,
		available: nil,
		idx:       0,
	}
}

// Clear resets the TaskSelector's state by clearing the index, setting the process ID to UnknownId, and nullifying the available list.
func (ts *TaskSelector) Clear() {
	ts.idx = 0
	ts.pid = adaptiveticker.UnknownId
	ts.available = nil
}

// AddAvailable adds the given process ID to the available pool for selection.
func (ts *TaskSelector) AddAvailable(pid int) {
	ts.available = append(ts.available, pid)
}

// Set updates the current process ID and index in the TaskSelector.
func (ts *TaskSelector) Set(pid int, idx int) {
	ts.pid = pid
	ts.idx = idx
}

// Next advances the selection index to the next available task in the list and updates the current task ID. Returns false if the list is empty.
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

// Prev moves the selection to the previous item in the available list, wrapping around if at the beginning. Returns true if successful.
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

// PID returns the currently selected process ID (pid) stored in the TaskSelector instance.
func (ts *TaskSelector) PID() int {
	return ts.pid
}
