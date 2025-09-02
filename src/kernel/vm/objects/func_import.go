package objects

import "encoding/gob"

func init() {
	gob.Register(&FuncImport{})
}

// FuncCallable is a type alias for a function that takes a variadic list of IObject arguments and returns an IObject and an error.
type FuncCallable = func(gk IGateKeeper, frame int, args ...IObject) (ret IObject, err error)

// FuncImport is a callable object type that encapsulates a function and provides execution context information.
type FuncImport struct {
	gk    IGateKeeper
	frame int
	name  string
	value FuncCallable
}

// NewFuncImport creates a new FuncImport instance with the specified Id and callable function.
func newFuncImport(factory IGateKeeper, frame int, name string, fn FuncCallable) IObject {
	return &FuncImport{
		gk:    factory,
		frame: frame,
		name:  name,
		value: fn,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *FuncImport) GateKeeper() IGateKeeper {
	return o.gk
}

// AsBool returns a boolean representation of the object, always returning false.
func (o *FuncImport) AsBool() bool {
	return false
}

// AsInt64 returns the length of the array as an int64 value.
func (o *FuncImport) AsInt64() int64 {
	return 0
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *FuncImport) AsFloat64() float64 {
	return 0
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *FuncImport) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// Nil checks if the object is nil and always returns false.
func (o *FuncImport) Nil() bool {
	return false
}

// SetStatic sets the frame to FrameStatic, marking it with a static execution context.
func (o *FuncImport) SetStatic() {
	o.frame = FrameStatic
}

// Frame returns the current frame value of the Object.
func (o *FuncImport) Frame() int {
	return o.frame
}

// LogicalOp performs a logical operation between the object and a right-hand operand using the specified operator.
// Returns the resulting object or an ErrInvalidOperator error if the operation is not valid.
func (o *FuncImport) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.gk, op)
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

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *FuncImport) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *FuncImport) IndexSet(_, _ IObject) (err error) {
	return ErrUnsupportedIndex
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *FuncImport) Iterate(_ int) IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *FuncImport) CanIterate() bool {
	return false
}

// Length returns the length of the Int object.
func (o *FuncImport) Length() int {
	return 0
}

// Name returns the name of the FuncImport as a string.
func (o *FuncImport) Name() string {
	return o.name
}

// TypeName returns the type name of the FuncImport as a string.
func (o *FuncImport) TypeName() string {
	return "FuncImport:" + o.name
}

// AsString returns the string representation of a FuncImport object.
func (o *FuncImport) AsString() string {
	return "<FuncImport>"
}

// Copy creates and returns a new FuncImport instance with the same Value field as the original object.
func (o *FuncImport) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewFuncImport(frame, o.name, o.value)
}

// Equals checks whether the current FuncImport is equal to another object of type IObject. Always returns false.
func (o *FuncImport) Equals(_ IObject) bool {
	return false
}

// Call invokes the function encapsulated within the FuncImport with the provided arguments and returns the result or an error.
func (o *FuncImport) Call(frame int, args ...IObject) (IObject, error) {
	return o.value(o.gk, frame, args...)
}

// CanCall checks whether the FuncImport instance can be invoked as a callable function. Always returns true.
func (o *FuncImport) CanCall() bool {
	return true
}
