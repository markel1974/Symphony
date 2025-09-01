package objects

import "encoding/gob"

func init() {
	gob.Register(&FuncExternal{})
}

// FuncCallable is a type alias for a function that takes a variadic list of IObject arguments and returns an IObject and an error.
type FuncCallable = func(gk IGateKeeper, frame int, args ...IObject) (ret IObject, err error)

// FuncExternal is a callable object type that encapsulates a function and provides execution context information.
type FuncExternal struct {
	gk    IGateKeeper
	frame int
	name  string
	value FuncCallable
}

// NewFuncExternal creates a new FuncExternal instance with the specified Id and callable function.
func newFuncExternal(factory IGateKeeper, frame int, name string, fn FuncCallable) IObject {
	return &FuncExternal{
		gk:    factory,
		frame: frame,
		name:  name,
		value: fn,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *FuncExternal) GateKeeper() IGateKeeper {
	return o.gk
}

// AsBool returns a boolean representation of the object, always returning false.
func (o *FuncExternal) AsBool() bool {
	return false
}

// AsInt64 returns the length of the array as an int64 value.
func (o *FuncExternal) AsInt64() int64 {
	return 0
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *FuncExternal) AsFloat64() float64 {
	return 0
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *FuncExternal) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// Nil checks if the object is nil and always returns false.
func (o *FuncExternal) Nil() bool {
	return false
}

// SetStatic sets the frame to FrameStatic, marking it with a static execution context.
func (o *FuncExternal) SetStatic() {
	o.frame = FrameStatic
}

// Frame returns the current frame value of the Object.
func (o *FuncExternal) Frame() int {
	return o.frame
}

// LogicalOp performs a logical operation between the object and a right-hand operand using the specified operator.
// Returns the resulting object or an ErrInvalidOperator error if the operation is not valid.
func (o *FuncExternal) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.gk, op)
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation using the specified operator and operand, returning an error for invalid operations.
func (o *FuncExternal) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Falsy returns false for all objects.
func (o *FuncExternal) Falsy() bool {
	return false
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *FuncExternal) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *FuncExternal) IndexSet(_, _ IObject) (err error) {
	return ErrUnsupportedIndex
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *FuncExternal) Iterate(_ int) IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *FuncExternal) CanIterate() bool {
	return false
}

// Length returns the length of the Int object.
func (o *FuncExternal) Length() int {
	return 0
}

// Name returns the name of the FuncExternal as a string.
func (o *FuncExternal) Name() string {
	return o.name
}

// TypeName returns the type name of the FuncExternal as a string.
func (o *FuncExternal) TypeName() string {
	return "FuncExternal:" + o.name
}

// AsString returns the string representation of a FuncExternal object.
func (o *FuncExternal) AsString() string {
	return "<FuncExternal>"
}

// Copy creates and returns a new FuncExternal instance with the same Value field as the original object.
func (o *FuncExternal) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewFuncExternal(frame, o.name, o.value)
}

// Equals checks whether the current FuncExternal is equal to another object of type IObject. Always returns false.
func (o *FuncExternal) Equals(_ IObject) bool {
	return false
}

// Call invokes the function encapsulated within the FuncExternal with the provided arguments and returns the result or an error.
func (o *FuncExternal) Call(frame int, args ...IObject) (IObject, error) {
	return o.value(o.gk, frame, args...)
}

// CanCall checks whether the FuncExternal instance can be invoked as a callable function. Always returns true.
func (o *FuncExternal) CanCall() bool {
	return true
}
