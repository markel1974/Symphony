package handler

import (
	"github.com/markel1974/symphony/src/vm/objects"
)

// IVM defines an interface for managing a virtual machine with functionalities for version retrieval.
type IVM interface {
	Version() int
	CoreId() uint
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

	CreateError(obj objects.IObject) objects.IObject
	CreateClosure(fn *objects.Func, objRequired []objects.IObject) objects.IObject
	CreateObjectPointer(obj objects.IObject) objects.IObject
	CreateSlice(highIdx int, lowIdx int, targetObj objects.IObject) objects.IObject

	Factory() objects.IGateKeeper
	Shutdown(err error)
}

// IVMReadOnly is an interface providing read-only access to Core components like constants, globals, imports, and frame.
type IVMReadOnly interface {
	IVM
	IVMStackOnly
	GlobalsGet(idx uint) (objects.IObject, error)
	GlobalsSet(idx uint, obj objects.IObject) error
	ImportsGet(idx uint) (objects.IObject, error)
	ConstantsGet(idx uint) (objects.IObject, error)
}

// IVMReadWrite represents an interface extending IVMReadOnly with additional write capabilities for Core global state.
type IVMReadWrite interface {
	IVMReadOnly
	GlobalsGet(idx uint) (objects.IObject, error)
	GlobalsSet(idx uint, obj objects.IObject) error
}

// IVMControlFlow defines an interface for managing instruction pointers and control flow within a virtual machine.
// It extends the IVMStackOnly interface to include stack-dependent operations related to reading conditions.
type IVMControlFlow interface {
	IVMStackOnly
	IVMFrameOnly
	SetIp(ip uint)
	GetIp() uint
}

// IVMFullAccess is an interface that merges IVMReadWrite and IVMControlFlow functionalities with additional execution controls.
// It enables invoking function-like objects and managing returned values within the Virtual Machine environment.
type IVMFullAccess interface {
	IVMReadWrite
	IVMControlFlow

	Call(value objects.IObject, async bool, spread bool, numArgs int)
	CallObject(value objects.IObject, numArgs int, args ...objects.IObject)
	Return(returnValues []objects.IObject)
	Shutdown(err error)
	BlockCurrentCore()
	WakeCore(id uint)
}
