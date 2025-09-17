package objects

import (
	"bytes"
	"encoding/gob"
)

func init() {
	gob.Register(&FuncJit{})
}

// FuncJit is a callable object type that encapsulates a function and provides execution context information.
type FuncJit struct {
	IAllocator
	name string
	data []byte
}

// NewFuncImport creates a new FuncImport instance with the specified Id and callable function.
func newFuncJit(allocator IAllocator, name string, fn []byte) IObject {
	return &FuncJit{
		IAllocator: allocator,
		name:       name,
		data:       fn,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *FuncJit) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *FuncJit) AsInterface() interface{} {
	return nil
}

// AsBool returns a boolean representation of the object, always returning false for FuncJit.
func (o *FuncJit) AsBool() bool {
	return false
}

// AsInt64 returns the len of the array as an int64 Code.
func (o *FuncJit) AsInt64() int64 {
	return 0
}

// AsFloat64 returns the len of the array as an int64 Code.
func (o *FuncJit) AsFloat64() float64 {
	return 0
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *FuncJit) AsBytes() []byte {
	return nil
}

// AsString returns the string representation of a FuncImport object.
func (o *FuncJit) AsString() string {
	return "<FuncJit>"
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *FuncJit) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// Nil checks if the object is nil and always returns false.
func (o *FuncJit) Nil() bool {
	return false
}

// LogicalOp performs a logical operation between the current object and a provided IObject using the specified operator.
// Always returns nil and ErrInvalidOperator as logical operations are not supported for this object type.
func (o *FuncJit) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.GateKeeper(), op)
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation using the specified operator and operand, returning the result or an error.
func (o *FuncJit) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Falsy returns false for all objects.
func (o *FuncJit) Falsy() bool {
	return false
}

// IndexGet attempts to retrieve a Code at the given index and returns an error if the object is not indexable.
func (o *FuncJit) IndexGet(_ int, _ IObject) (IObject, error) {
	return o.GateKeeper().UndefinedValue(), ErrIndexNotIndexable
}

// IndexSet attempts to assign a Code to an index in the object but always returns ErrIndexUnsupported,
// as this operation is unsupported.
func (o *FuncJit) IndexSet(_, _ IObject) error {
	return ErrIndexUnsupported
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *FuncJit) Iterate(_ int) IIterator {
	return nil
}

// Iterable determines if the object can be iterated over and returns false for this implementation.
func (o *FuncJit) Iterable() bool {
	return false
}

// Length returns the len of the Int object.
func (o *FuncJit) Length() int {
	return 0
}

// Name returns the name of the FuncImport as a string.
func (o *FuncJit) Name() string {
	return o.name
}

// TypeName returns the type name of the FuncImport as a string.
func (o *FuncJit) TypeName() string {
	return "FuncJit:" + o.name
}

// Copy creates and returns a new FuncImport instance with the same Value field as the original object.
func (o *FuncJit) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewFuncJit(frame, o.name, o.data)
}

// Equals checks whether the current FuncImport is equal to another object of type IObject. Always returns false.
func (o *FuncJit) Equals(_ IObject) bool {
	return false
}

// Call invokes the function encapsulated within the FuncImport with the provided arguments and returns the result or an error.
func (o *FuncJit) Call(_ int, _ ...IObject) (uint, IObject, error) {
	return 0, nil, ErrUnimplemented
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *FuncJit) Count() int {
	return 1
}

// GobEncode serializes the FuncJit's data into a byte slice using gob encoding and returns the result or an error.
func (o *FuncJit) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(o.name); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode decodes the provided byte slice into the FuncJit's data field using the gob package.
func (o *FuncJit) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&o.name); err != nil {
		return err
	}
	if err := decoder.Decode(&o.data); err != nil {
		return err
	}
	return nil
}

/*
func (*FuncJit) create(machineCode []byte) error {
	//machineCode := []byte{0x48, 0xc7, 0xc0, 0x2a, 0x00, 0x00, 0x00, 0xc3}
	memory := make([]byte, len(machineCode))
	copy(memory, machineCode)
	memory, err := unix.Mmap(
		-1,   // fd (file descriptor), -1 for anonymous memory
		0,    // offset
		4096, // len (es. 1 un page)
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
