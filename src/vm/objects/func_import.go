package objects

import "encoding/gob"

func init() {
	gob.Register(&FuncImport{})
}

// Invocable is a type alias for a function that takes a variadic list of IObject arguments and returns an IObject and an error.
type Invocable = func(gk IGateKeeper, frame int, args ...IObject) (retCount uint, ret IObject, err error)

// FuncImport is a callable object type that encapsulates a function and provides execution context information.
type FuncImport struct {
	IAllocator
	FnName string
	Args   int
	Data   Invocable
}

// NewFuncImport creates a new FuncImport instance with the specified Id and callable function.
func newFuncImport(allocator IAllocator, name string, args int, fn Invocable) IObject {
	return &FuncImport{
		IAllocator: allocator,
		FnName:     name,
		Args:       args,
		Data:       fn,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *FuncImport) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *FuncImport) AsInterface() interface{} {
	return nil
}

// AsBool returns a boolean representation of the object, always returning false.
func (o *FuncImport) AsBool() bool {
	return false
}

// AsInt64 returns the len of the array as an int64 Code.
func (o *FuncImport) AsInt64() int64 {
	return 0
}

// AsFloat64 returns the len of the array as an int64 Code.
func (o *FuncImport) AsFloat64() float64 {
	return 0
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *FuncImport) AsBytes() []byte {
	return nil
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *FuncImport) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// Nil checks if the object is nil and always returns false.
func (o *FuncImport) Nil() bool {
	return false
}

// LogicalOp performs a logical operation between the object and a right-hand operand using the specified operator.
// Returns the resulting object or an ErrInvalidOperator error if the operation is not valid.
func (o *FuncImport) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.GateKeeper(), op)
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation using the specified operator and operand, returning an error for invalid operations.
func (o *FuncImport) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Falsy returns false for all objects.
func (o *FuncImport) Falsy() bool {
	return false
}

// IndexGet attempts to retrieve a Code at the given index and returns an error if the object is not indexable.
func (o *FuncImport) IndexGet(_ int, _ IObject) (IObject, error) {
	return o.GateKeeper().UndefinedValue(), ErrIndexNotIndexable
}

// IndexSet attempts to assign a Code to an index in the object but always returns ErrIndexUnsupported,
// as this operation is unsupported.
func (o *FuncImport) IndexSet(_, _ IObject) error {
	return ErrIndexUnsupported
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *FuncImport) Iterate(_ int) IIterator {
	return nil
}

// Iterable determines if the object can be iterated over and returns false for this implementation.
func (o *FuncImport) Iterable() bool {
	return false
}

// Length returns the len of the Int object.
func (o *FuncImport) Length() int {
	return 0
}

// Name returns the FnName of the FuncImport as a string.
func (o *FuncImport) Name() string {
	return o.FnName
}

// TypeName returns the type FnName of the FuncImport as a string.
func (o *FuncImport) TypeName() string {
	return "FuncImport:" + o.FnName
}

// AsString returns the string representation of a FuncImport object.
func (o *FuncImport) AsString() string {
	return "<FuncImport>"
}

// Copy creates and returns a new FuncImport instance with the same Value field as the original object.
func (o *FuncImport) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewFuncImport(frame, o.FnName, o.Args, o.Data)
}

// Equals checks whether the current FuncImport is equal to another object of type IObject. Always returns false.
func (o *FuncImport) Equals(_ IObject) bool {
	return false
}

// Call invokes the function encapsulated within the FuncImport with the provided arguments and returns the result or an error.
func (o *FuncImport) Call(frame int, args ...IObject) (uint, IObject, error) {
	if o.Args < 0 || len(args) == o.Args {
		return o.Data(o.GateKeeper(), frame, args...)
	}
	return 0, o.GateKeeper().UndefinedValue(), ErrInvalidArgumentsNumber
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *FuncImport) Count() int {
	return 1
}
