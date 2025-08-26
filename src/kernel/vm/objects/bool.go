package objects

import "encoding/gob"

const (
	BoolType = "bool"
)

func init() {
	gob.Register(&Bool{})
}

// Bool is a custom type representing a boolean values, implementing IObject interface and encapsulating a boolean value.
type Bool struct {
	gk    IGateKeeper
	frame int
	value bool
}

// NewBool creates and returns a new Bool object with the specified boolean value.
func newBool(factory IGateKeeper, frame int, value bool) IObject {
	return &Bool{
		gk:    factory,
		frame: frame,
		value: value,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *Bool) GateKeeper() IGateKeeper {
	return o.gk
}

// Frame returns the current frame value of the Object.
func (o *Bool) Frame() int {
	return o.frame
}

// BinaryOp performs a binary operation on the current object and another object using the specified operator.
// Returns the result of the operation or an error if the operation is not supported.
func (o *Bool) BinaryOp(_ int, _ Operator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *Bool) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrNotIndexAssignable,
// as this operation is unsupported.
func (o *Bool) IndexSet(_, _ IObject) (err error) {
	return ErrNotIndexAssignable
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *Bool) Iterate(_ int) IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *Bool) CanIterate() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Bool) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *Bool) CanCall() bool {
	return false
}

// Length returns the length of the Int object.
func (o *Bool) Length() int {
	return 0
}

// String returns the string representation of the Bool object, index.e., "true" if the value is true, otherwise "false".
func (o *Bool) String() string {
	if o.value {
		return "true"
	}
	return "false"
}

// TypeName returns the name of the type as a string, specifically "bool" for the Bool type.
func (o *Bool) TypeName() string {
	return BoolType
}

// Copy creates and returns a reference to the current Bool object.
func (o *Bool) Copy(frame int, _ int) IObject {
	return o.gk.NewBool(frame, o.value)
}

// Boolean returns true if the Bool value is false, otherwise returns false.
func (o *Bool) Boolean() bool {
	return !o.value
}

// Equals compares the Bool object with another IObject and returns true if both are equal, otherwise false.
func (o *Bool) Equals(x IObject) bool {
	return o == x
}

// GobDecode decodes a byte slice into the Bool value, setting the Bool's value to true if the first byte equals 1.
func (o *Bool) GobDecode(b []byte) (err error) {
	o.value = b[0] == 1
	return
}

// GobEncode serializes the Bool object into a byte slice based on its boolean value. Returns the serialized data and error, if any.
func (o *Bool) GobEncode() (b []byte, err error) {
	if o.value {
		b = []byte{1}
	} else {
		b = []byte{0}
	}
	return
}
