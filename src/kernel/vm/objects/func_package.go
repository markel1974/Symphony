package objects

import "encoding/gob"

const (
	FuncBuiltinDef = "func_builtin"
	FuncPackageDef = "func_package"
)

func init() {
	gob.Register(&FuncPackage{})
}

// FuncCallable is a type alias for a function that takes a variadic list of IObject arguments and returns an IObject and an error.
type FuncCallable = func(frame int, args ...IObject) (ret IObject, err error)

// FuncPackage is a callable object type that encapsulates a function and provides execution context information.
type FuncPackage struct {
	gk    IGateKeeper
	frame int
	kind  string
	name  string
	value FuncCallable
}

// NewFuncPackage creates a new FuncPackage instance with the specified Id and callable function.
func newFuncPackage(factory IGateKeeper, frame int, kind string, name string, fn FuncCallable) IObject {
	return &FuncPackage{
		gk:    factory,
		frame: frame,
		kind:  kind,
		name:  name,
		value: fn,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *FuncPackage) GateKeeper() IGateKeeper {
	return o.gk
}

// AsBool returns a boolean representation of the object, always returning false.
func (o *FuncPackage) AsBool() bool {
	return false
}

// AsInt64 returns the length of the array as an int64 value.
func (o *FuncPackage) AsInt64() int64 {
	return 0
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *FuncPackage) AsFloat64() float64 {
	return 0
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *FuncPackage) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// SetStatic sets the frame to FrameStatic, marking it with a static execution context.
func (o *FuncPackage) SetStatic() {
	o.frame = FrameStatic
}

// Frame returns the current frame value of the Object.
func (o *FuncPackage) Frame() int {
	return o.frame
}

// LogicalOp performs a logical operation between the object and a right-hand operand using the specified operator.
// Returns the resulting object or an ErrInvalidOperator error if the operation is not valid.
func (o *FuncPackage) LogicalOp(_ int, _ LogicalOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation using the specified operator and operand, returning an error for invalid operations.
func (o *FuncPackage) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Falsy returns false for all objects.
func (o *FuncPackage) Falsy() bool {
	return false
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *FuncPackage) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *FuncPackage) IndexSet(_, _ IObject) (err error) {
	return ErrUnsupportedIndex
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *FuncPackage) Iterate(_ int) IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *FuncPackage) CanIterate() bool {
	return false
}

// Length returns the length of the Int object.
func (o *FuncPackage) Length() int {
	return 0
}

// Name returns the name of the FuncPackage as a string.
func (o *FuncPackage) Name() string {
	return o.name
}

// TypeName returns the type name of the FuncPackage as a string.
func (o *FuncPackage) TypeName() string {
	return o.kind + ":" + o.name
}

// AsString returns the string representation of a FuncPackage object.
func (o *FuncPackage) AsString() string {
	return "<" + o.kind + ">"
}

// Copy creates and returns a new FuncPackage instance with the same Value field as the original object.
func (o *FuncPackage) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewFuncPackageFrame(frame, o.kind, o.name, o.value)
}

// Equals checks whether the current FuncPackage is equal to another object of type IObject. Always returns false.
func (o *FuncPackage) Equals(_ IObject) bool {
	return false
}

// Call invokes the function encapsulated within the FuncPackage with the provided arguments and returns the result or an error.
func (o *FuncPackage) Call(frame int, args ...IObject) (IObject, error) {
	return o.value(frame, args...)
}

// CanCall checks whether the FuncPackage instance can be invoked as a callable function. Always returns true.
func (o *FuncPackage) CanCall() bool {
	return true
}
