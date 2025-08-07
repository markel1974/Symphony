package core

import "github.com/markel1974/c64emu/src/kernel/interfaces"

// PID represents a unique process identifier encapsulating an ID value.
type PID struct {
	id int
}

// NewPID creates and returns a new instance of the PID structure with default values initialized.
func NewPID() *PID {
	return &PID{}
}

// SetId assigns a new integer value to the id field of the PID object.
func (p *PID) SetId(i int) {
	p.id = i
}

// GetId retrieves the identifier value stored in the PID instance.
func (p *PID) GetId() int {
	return p.id
}

// KernelProcess represents a kernel-level process with parent linkage, protection status, and associated timer IDs.
type KernelProcess struct {
	interfaces.IProcess
	user      string
	line      string
	name      string
	pid       *PID
	parent    interfaces.IProcess
	protected bool
	timers    []int
}

// NewKernelProcess creates a new KernelProcess instance with a parent process, protection flag, and assigned process.
func NewKernelProcess(user string, line, name string, parent interfaces.IProcess, pid *PID, protected bool, process interfaces.IProcess) *KernelProcess {
	return &KernelProcess{
		user:      user,
		line:      line,
		name:      name,
		parent:    parent,
		protected: protected,
		IProcess:  process,
		pid:       pid,
	}
}

// Description provides a brief summary of the process including its name, PID, and line information.
func (kp *KernelProcess) Description() *interfaces.ProcessDescription {
	return interfaces.NewProcessDescription(kp.name, kp.pid.GetId(), kp.line)
}

// AddTimer adds a timer ID to the KernelProcess's list of timers.
func (kp *KernelProcess) AddTimer(tid int) {
	kp.timers = append(kp.timers, tid)
}

// Parent returns the parent process of the current KernelProcess, implementing the IProcess interface.
func (kp *KernelProcess) Parent() interfaces.IProcess {
	return kp.parent
}

// Protected returns true if the KernelProcess is marked as protected, otherwise false.
func (kp *KernelProcess) Protected() bool {
	return kp.protected
}

// Timers returns the list of timer IDs associated with the KernelProcess.
func (kp *KernelProcess) Timers() []int {
	return kp.timers
}

// TimersIterator iterates over registered timer IDs in the process, calling the provided callback for each timer ID.
// If the callback returns true, the iteration stops immediately.
func (kp *KernelProcess) TimersIterator(callback func(tid int) bool) {
	for _, tid := range kp.timers {
		if callback(tid) {
			break
		}
	}
}
