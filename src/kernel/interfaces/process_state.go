package interfaces

// ProcessState represents the state of a task, defined as an integer-based enumerated type.
type ProcessState int

// TaskStateSetup represents the initial setup state of a task.
// TaskStateRunning represents the state where a task is currently running.
const (
	ProcessStateSetup   ProcessState = iota
	ProcessStateRunning ProcessState = iota
)
