package core

import (
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// IVM defines an interface for managing a virtual machine with functionalities for version retrieval.
type IVM interface {
	Version() string
}

// IVMStackOnly defines an interface for managing a virtual machine stack with functionalities for stack retrieval,
// gatekeeper factory access, and error handling.
type IVMStackOnly interface {
	IVM
	Stack() *Stack
	Factory() objects.IGateKeeper
	SetError(err error)
}

// IVMReadOnly is an interface providing read-only access to VM components like constants, globals, imports, and frame.
type IVMReadOnly interface {
	IVM
	IVMStackOnly
	Constants() *Constants
	Globals() *Globals
	Imports() *Imports
	Internals() *Internals
	Frame() *Frame
}

// IVMReadWrite represents an interface extending IVMReadOnly with additional write capabilities for VM global state.
type IVMReadWrite interface {
	IVMReadOnly
	Globals() *Globals // Ridefiniamo per chiarezza, anche se già presente
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
	Shutdown()
}
