package core

import (
	"github.com/markel1974/c64emu/src/vm/objects"
)

// IVM defines an interface for managing a virtual machine with functionalities for version retrieval.
type IVM interface {
	Version() string
}

type IVMFrameOnly interface {
	IVM
	FrameId() int
	FrameDeferredAdd(obj objects.IObject)
	FrameFreeVarsIndex(index uint) *objects.ObjectPointer
}

// IVMStackOnly defines an interface for managing a virtual machine stack with functionalities for stack retrieval,
// gatekeeper factory access, and error handling.
type IVMStackOnly interface {
	IVM
	StackDecrement()
	StackDecrementCount(count uint)

	StackPeek() objects.IObject
	StackPop() objects.IObject
	StackPush(obj objects.IObject)
	StackSet(objects.IObject)
	StackPeekBP(offset uint) objects.IObject
	StackSetBP(offset uint, obj objects.IObject)
	StackPeekSP(offset uint) objects.IObject
	StackSetSP(offset uint, obj objects.IObject)

	StackPopArray(numElem uint) objects.IObject
	StackPopMap(numElem uint) objects.IObject
	StackPopStruct(numElem uint) objects.IObject
	StackPopInterface(numElem int) objects.IObject

	CreateClosure(fnObj objects.IObject) objects.IObject
	CreateObjectPointer(obj objects.IObject) (*objects.ObjectPointer, error)

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
}

// IVMReadWrite represents an interface extending IVMReadOnly with additional write capabilities for VM global state.
type IVMReadWrite interface {
	IVMReadOnly
	Globals() *Globals
}

// IVMControlFlow defines an interface for managing instruction pointers and control flow within a virtual machine.
// It extends the IVMStackOnly interface to include stack-dependent operations related to reading conditions.
type IVMControlFlow interface {
	IVMStackOnly
	IVMFrameOnly
	SetIp(ip int)
	GetIp() int
}

// IVMFullAccess is an interface that merges IVMReadWrite and IVMControlFlow functionalities with additional execution controls.
// It enables invoking function-like objects and managing returned values within the Virtual Machine environment.
type IVMFullAccess interface {
	IVMReadWrite
	IVMControlFlow

	Call(value objects.IObject, spread bool, numArgs int)
	CallObject(value objects.IObject, numArgs int, args ...objects.IObject)
	Return(returnValues []objects.IObject)
	Shutdown()
}
