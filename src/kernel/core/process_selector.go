package core

import "github.com/markel1974/c64emu/src/kernel/adaptiveticker"

// ProcessSelector is a type that manages available task identifiers and tracks the current one for selection and navigation.
type ProcessSelector struct {
	pid       int
	available []int
	idx       int
}

// NewProcessSelector creates and returns a new ProcessSelector initialized with default values.
func NewProcessSelector() *ProcessSelector {
	return &ProcessSelector{
		pid:       adaptiveticker.UnknownId,
		available: nil,
		idx:       0,
	}
}

// Clear resets the ProcessSelector's state by clearing the index, setting the process ID to UnknownId, and nullifying the available list.
func (ts *ProcessSelector) Clear() {
	ts.idx = 0
	ts.pid = adaptiveticker.UnknownId
	ts.available = nil
}

// AddAvailable adds the given process ID to the available pool for selection.
func (ts *ProcessSelector) AddAvailable(pid int) {
	ts.available = append(ts.available, pid)
}

// Set updates the current process ID and index in the ProcessSelector.
func (ts *ProcessSelector) Set(pid int, idx int) {
	ts.pid = pid
	ts.idx = idx
}

// Next advances the selection index to the next available task in the list and updates the current task ID. Returns false if the list is empty.
func (ts *ProcessSelector) Next() bool {
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
func (ts *ProcessSelector) Prev() bool {
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

// PID returns the currently selected process ID (pid) stored in the ProcessSelector instance.
func (ts *ProcessSelector) PID() int {
	return ts.pid
}
