package objects

import (
	"encoding/gob"
)

func init() {
	gob.Register(&FuncPackage{})
}

// FuncJit is a callable object type that encapsulates a function and provides execution context information.
type FuncJit struct {
	gk    IGateKeeper
	frame int
	kind  string
	name  string
	value []byte
}

// NewFuncPackage creates a new FuncPackage instance with the specified Id and callable function.
func newFuncJit(factory IGateKeeper, frame int, kind string, name string, fn []byte) IObject {
	return &FuncJit{
		gk:    factory,
		frame: frame,
		kind:  kind,
		name:  name,
		value: fn,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *FuncJit) GateKeeper() IGateKeeper {
	return o.gk
}

// Frame returns the current frame value of the Object.
func (o *FuncJit) Frame() int {
	return o.frame
}

// BinaryOp performs a binary operation on the current object and another object using the specified operator.
// Returns the result of the operation or an error if the operation is not supported.
func (o *FuncJit) BinaryOp(_ int, _ Operator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Falsy returns false for all objects.
func (o *FuncJit) Falsy() bool {
	return false
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *FuncJit) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *FuncJit) IndexSet(_, _ IObject) (err error) {
	return ErrUnsupportedIndex
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *FuncJit) Iterate(_ int) IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *FuncJit) CanIterate() bool {
	return false
}

// Length returns the length of the Int object.
func (o *FuncJit) Length() int {
	return 0
}

// Name returns the name of the FuncPackage as a string.
func (o *FuncJit) Name() string {
	return o.name
}

// TypeName returns the type name of the FuncPackage as a string.
func (o *FuncJit) TypeName() string {
	return o.kind + ":" + o.name
}

// String returns the string representation of a FuncPackage object.
func (o *FuncJit) String() string {
	return "<" + o.kind + ">"
}

// Copy creates and returns a new FuncPackage instance with the same Value field as the original object.
func (o *FuncJit) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewFuncJitFrame(frame, o.kind, o.name, o.value)
}

// Equals checks whether the current FuncPackage is equal to another object of type IObject. Always returns false.
func (o *FuncJit) Equals(_ IObject) bool {
	return false
}

// Call invokes the function encapsulated within the FuncPackage with the provided arguments and returns the result or an error.
func (o *FuncJit) Call(_ int, _ ...IObject) (IObject, error) {
	return nil, ErrUnimplemented
}

// CanCall checks whether the FuncPackage instance can be invoked as a callable function. Always returns true.
func (o *FuncJit) CanCall() bool {
	return true
}

/*
func (*FuncJit) create(machineCode []byte) error {
	//machineCode := []byte{0x48, 0xc7, 0xc0, 0x2a, 0x00, 0x00, 0x00, 0xc3}
	memory := make([]byte, len(machineCode))
	copy(memory, machineCode)
	memory, err := unix.Mmap(
		-1,   // fd (file descriptor), -1 for anonymous memory
		0,    // offset
		4096, // length (es. 1 un page)
		unix.PROT_READ|unix.PROT_WRITE|unix.PROT_EXEC,
		unix.MAP_ANON|unix.MAP_PRIVATE,
	)
	if err != nil {
		return err
	}
	ptrToMemory := &memory[0]
	type nativeFunc func() int
	funcPtr := (nativeFunc)(unsafe.Pointer(ptrToMemory))
	result := funcPtr()
}
*/
