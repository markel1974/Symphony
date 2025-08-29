package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// IVMStackOnly defines an interface for managing a virtual machine stack with functionalities for stack retrieval,
// gatekeeper factory access, and error handling.
type IVMStackOnly interface {
	Stack() *core.Stack
	Factory() objects.IGateKeeper
	SetError(err error)
}

// IVMReadOnly is an interface providing read-only access to VM components like constants, globals, references, and frame.
type IVMReadOnly interface {
	IVMStackOnly
	Constants() *core.Constants
	Globals() *core.Globals
	References() *core.References
	Frame() *core.Frame
}

// IVMReadWrite represents an interface extending IVMReadOnly with additional write capabilities for VM global state.
type IVMReadWrite interface {
	IVMReadOnly
	Globals() *core.Globals // Ridefiniamo per chiarezza, anche se già presente
}

// IVMControlFlow defines an interface for managing instruction pointers and control flow within a virtual machine.
// It extends the IVMStackOnly interface to include stack-dependent operations related to reading conditions.
type IVMControlFlow interface {
	IVMStackOnly // Ha bisogno dello stack per leggere le condizioni
	SetIp(ip int)
	GetIp() int
}

// IVMFullAccess is an interface that merges IVMReadWrite and IVMControlFlow functionalities with additional execution controls.
// It enables invoking function-like objects and managing returned values within the Virtual Machine environment.
type IVMFullAccess interface {
	IVMReadWrite
	IVMControlFlow
	Call(value objects.IObject, spread bool, numArgs int)
	Return(returnValues []objects.IObject)
}

// VMFacade serves as a high-level interface to interact with the underlying virtual machine (VM) abstraction.
type VMFacade struct {
	vm *core.VM
}

// NewVMFacade creates and returns a new instance of VMFacade initialized with the provided VM instance.
func NewVMFacade(vm *core.VM) *VMFacade {
	return &VMFacade{vm: vm}
}

// Stack returns the current stack of the virtual machine.
func (f *VMFacade) Stack() *core.Stack { return f.vm.Stack() }

// Factory returns the IGateKeeper instance of the underlying virtual machine.
func (f *VMFacade) Factory() objects.IGateKeeper { return f.vm.Factory() }

// SetError sets the provided error value in the underlying VM instance.
func (f *VMFacade) SetError(err error) { f.vm.SetError(err) }

// Constants returns the core.Constants instance managed by the underlying VM.
func (f *VMFacade) Constants() *core.Constants { return f.vm.Constants() }

// Globals provides access to the shared global state managed by the core VM, including objects and initialization data.
func (f *VMFacade) Globals() *core.Globals { return f.vm.Globals() }

// References returns a pointer to the core.References structure, enabling access to object reference management.
func (f *VMFacade) References() *core.References { return f.vm.References() }

// Frame retrieves the current function call frame from the virtual machine context.
func (f *VMFacade) Frame() *core.Frame { return f.vm.Frame() }

// SetIp sets the instruction pointer of the underlying virtual machine to the specified value.
func (f *VMFacade) SetIp(ip int) { f.vm.SetIp(ip) }

// GetIp retrieves the current instruction pointer (IP) from the underlying VM instance.
func (f *VMFacade) GetIp() int { return f.vm.GetIp() }

// Call invokes a function or callable object represented by 'value' with specified arguments and spread behavior.
func (f *VMFacade) Call(value objects.IObject, spread bool, numArgs int) {
	f.vm.Call(value, spread, numArgs)
}

// Return delegates the return of specified values to the underlying virtual machine instance.
func (f *VMFacade) Return(returnValues []objects.IObject) { f.vm.Return(returnValues) }
