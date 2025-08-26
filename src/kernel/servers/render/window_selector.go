package render

import (
	"github.com/markel1974/c64emu/src/kernel/adaptiveticker"
)

// WindowSelector is a type that manages available task identifiers and tracks the current one for selection and navigation.
type WindowSelector struct {
	pid       int
	available []int
	idx       int
}

// NewWindowSelector creates and returns a new WindowSelector initialized with default values.
func NewWindowSelector() *WindowSelector {
	return &WindowSelector{
		pid:       adaptiveticker.UnknownId,
		available: nil,
		idx:       0,
	}
}

// Clear resets the WindowSelector's state by clearing the index, setting the process Id to UnknownId, and nullifying the available list.
func (ts *WindowSelector) Clear() {
	ts.idx = 0
	ts.pid = adaptiveticker.UnknownId
	ts.available = nil
}

// Len returns the length of the available list.
func (ts *WindowSelector) Len() int {
	return len(ts.available)
}

// AddAvailable adds the given process Id to the available pool for selection.
func (ts *WindowSelector) AddAvailable(pid int) {
	ts.available = append(ts.available, pid)
}

// Set updates the current process Id and index in the WindowSelector.
func (ts *WindowSelector) Set(pid int, idx int) {
	ts.pid = pid
	ts.idx = idx
}

// Get retrieves the process Id at the specified index in the available list.
func (ts *WindowSelector) Get(idx int) (int, bool) {
	if idx < 0 || idx >= len(ts.available) {
		return adaptiveticker.UnknownId, false
	}
	return ts.available[idx], true
}

// Next advances the selection index to the next available task in the list and updates the current task Id. Returns false if the list is empty.
func (ts *WindowSelector) Next() bool {
	if len(ts.available) == 0 {
		return false
	}
	next := ts.idx + 1
	if next < 0 {
		next = len(ts.available) - 1
	} else if next >= len(ts.available) {
		next = 0
	}
	ts.idx = next
	ts.pid = ts.available[next]
	return true
}

// Prev moves the selection to the previous item in the available list, wrapping around if at the beginning. Returns true if successful.
func (ts *WindowSelector) Prev() bool {
	if len(ts.available) == 0 {
		return false
	}
	prev := ts.idx - 1
	if prev < 0 {
		prev = len(ts.available) - 1
	} else if prev >= len(ts.available) {
		prev = 0
	}

	ts.idx = prev
	ts.pid = ts.available[prev]
	return true
}

// PID returns the currently selected process Id (pid) stored in the WindowSelector instance.
func (ts *WindowSelector) PID() int {
	return ts.pid
}
