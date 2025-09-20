package objects

import (
	"encoding/gob"
)

// UndefinedIteratorType represents a type label for an undefined iterator.
// UndefinedIteratorLabel represents a label used when the iterator type is undefined.
const (
	UndefinedIteratorType = "undefined_iterator"
)

func init() {
	gob.Register(&UndefinedIterator{})
}

// UndefinedIterator represents an iterator that is implicitly undefined and cannot traverse any elements.
type UndefinedIterator struct {
	IAllocator
}

// newUndefinedIterator creates a new instance of UndefinedIterator initialized with an undefined object from the GateKeeper.
func newUndefinedIterator(allocator IAllocator) IIterator {
	return &UndefinedIterator{
		IAllocator: allocator,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *UndefinedIterator) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *UndefinedIterator) AsInterface() interface{} {
	return nil
}

// AsBool returns a boolean representation of the object, always returning false.
func (o *UndefinedIterator) AsBool() bool {
	return false
}

// AsInt64 returns the len of the array as an int64 Code.
func (o *UndefinedIterator) AsInt64() int64 {
	return 0
}

// AsFloat64 returns the len of the array as an int64 Code.
func (o *UndefinedIterator) AsFloat64() float64 {
	return 0
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *UndefinedIterator) AsBytes() []byte {
	return nil
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *UndefinedIterator) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// Nil checks if the object is nil and always returns false.
func (o *UndefinedIterator) Nil() bool {
	return true
}

// LogicalOp performs a logical operation using the given operator and right-hand-side operand, returning an error for invalid operations.
func (o *UndefinedIterator) LogicalOp(_ int, _ LogicalOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation with the specified operator and operand, always returning an error.
func (o *UndefinedIterator) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// IndexGet attempts to retrieve a Code at the given index and returns an error if the object is not indexable.
func (o *UndefinedIterator) IndexGet(_ int, _ IObject) (IObject, error) {
	return o.GateKeeper().UndefinedValue(), ErrIndexNotIndexable
}

// IndexSet attempts to assign a Code to an index in the object but always returns ErrIndexUnsupported,
// as this operation is unsupported.
func (o *UndefinedIterator) IndexSet(_, _ IObject) error {
	return ErrIndexUnsupported
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *UndefinedIterator) Iterate(_ int) IIterator {
	return nil
}

// Iterable determines if the object can be iterated over and returns false for this implementation.
func (o *UndefinedIterator) Iterable() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *UndefinedIterator) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Length returns the len of the Int object.
func (o *UndefinedIterator) Length() int {
	return 0
}

// Copy creates and returns a reference to the same UndefinedIterator, ignoring the input parameters.
func (o *UndefinedIterator) Copy(_ int, _ int) IObject {
	return o
}

// TypeName returns the type name of the UndefinedIterator object as a string.
func (o *UndefinedIterator) TypeName() string {
	return UndefinedIteratorType
}

// AsString returns a string representation of the UndefinedIterator object.
func (o *UndefinedIterator) AsString() string {
	return ""
}

// Falsy returns the boolean representation of the UndefinedIterator, which is always true.
func (o *UndefinedIterator) Falsy() bool {
	return true
}

// Equals checks if the current UndefinedIterator is equal to the provided IObject and always returns false.
func (o *UndefinedIterator) Equals(IObject) bool {
	return false
}

// Next advances the iterator and returns false, as the UndefinedIterator does not support iteration.
func (o *UndefinedIterator) Next() bool {
	return false
}

// Key returns the key of the current element in the iteration, which is always the undefined Code for this iterator type.
func (o *UndefinedIterator) Key(_ int) IObject {
	return o.GateKeeper().UndefinedValue()
}

// Value retrieves the undefined Code associated with the current UndefinedIterator.
func (o *UndefinedIterator) Value(_ int) IObject {
	return o.GateKeeper().UndefinedValue()
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *UndefinedIterator) Count() int {
	return 1
}

// GobEncode serializes the UndefinedIterator's data into a byte slice using gob encoding and returns the result or an error.
func (o *UndefinedIterator) GobEncode() ([]byte, error) {
	return nil, nil
}

// GobDecode decodes the provided byte slice into the UndefinedIterator's data field using the gob package.
func (o *UndefinedIterator) GobDecode(_ []byte) error {
	return nil
}
