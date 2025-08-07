package interfaces

// ProcessState represents the state of a task, defined as an integer-based enumerated type.
type ProcessState int

// ProcessStateSetup represents the initial setup state of a process.
// ProcessStateRunning represents the state where a task is currently running.
const (
	ProcessStateSetup   ProcessState = iota
	ProcessStateRunning ProcessState = iota
)
