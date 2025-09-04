package objects

import "encoding/gob"

// BoolType defines the string representation of the boolean type. It is used as the type name for boolean objects.
const (
	BoolType = "bool"
)

// init registers the Bool type with the gob package for encoding and decoding.
func init() {
	gob.Register(&Bool{})
}

// Bool represents a boolean type with frame-specific and gatekeeper-managed value assignments.
type Bool struct {
	Allocator
	value bool
}

// newBool creates and returns a new instance of Bool, initializing it with the given frame and boolean value.
func newBool(factory IGateKeeper, frame int, value bool) IObject {
	return &Bool{
		Allocator: Allocator{gk: factory, frame: frame},
		value:     value,
	}
}

func (o *Bool) AsBool() bool {
	return o.value
}

// AsInt64 returns the length of the array as an int64 value.
func (o *Bool) AsInt64() int64 {
	if o.value {
		return 1
	}
	return 0
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *Bool) AsFloat64() float64 {
	if o.value {
		return 1
	}
	return 0
}

// AsString returns the string representation of the Bool object, either "true" or "false" based on its boolean value.
func (o *Bool) AsString() string {
	if o.value {
		return "true"
	}
	return "false"
}

// AssignValue assigns the value of another Bool instance to the current instance. Returns an error if the input is not a Bool.
func (o *Bool) AssignValue(v IObject) error {
	target, ok := v.(*Bool)
	if !ok {
		return ErrNotAssignable
	}
	o.value = target.value
	return nil
}

// Nil checks if the object is nil and always returns false.
func (o *Bool) Nil() bool {
	return false
}

// LogicalOp performs a logical operation on the Bool instance using the provided operator and right-hand operand.
// Returns the result of the operation as an IObject or an error if the operation is invalid.
func (o *Bool) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.gk, op)
	}
	lhsValue := int64(0)
	if o.value {
		lhsValue = 1
	}
	ret, err := logicalOpInt64(lhsValue, op, rhsIn.AsInt64())
	if err != nil {
		return nil, err
	}
	if ret {
		return o.gk.TrueValue(), nil
	}
	return o.gk.FalseValue(), nil
}

// ArithmeticOp performs an arithmetic operation between the Bool object and a given IObject using the specified operator.
// Returns the result as an IObject and an error if the operation is not valid or executable.
func (o *Bool) ArithmeticOp(_ int, op ArithmeticOperator, rhsIn IObject) (IObject, error) {
	lhsValue := int64(0)
	if o.value {
		lhsValue = 1
	}
	ret, err := arithmeticOpInt64(lhsValue, op, rhsIn.AsInt64())
	if err != nil {
		return nil, err
	}
	if ret != 0 {
		return o.gk.TrueValue(), nil
	}
	return o.gk.FalseValue(), nil
}

// IndexGet retrieves the value at a given index from the Bool object, but always returns an error as Bool is not indexable.
func (o *Bool) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to set an index on the Bool object but always returns ErrUnsupportedIndex as Bool is not indexable.
func (o *Bool) IndexSet(_, _ IObject) (err error) {
	return ErrUnsupportedIndex
}

// Iterate returns nil as Bool does not support iteration.
func (o *Bool) Iterate(_ int) IIterator {
	return nil
}

// Iterable indicates whether the Bool object can be iterated. Always returns false.
func (o *Bool) Iterable() bool {
	return false
}

// Call invokes the Bool object as a callable function with the provided arguments, returning nil and no error.
func (o *Bool) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Length returns the length of the Bool object, which is always 0.
func (o *Bool) Length() int {
	return 0
}

// TypeName returns the type name of the Bool object as a string.
func (o *Bool) TypeName() string {
	return BoolType
}

// Copy creates and returns a new Bool instance with the same value and the specified execution frame.
func (o *Bool) Copy(frame int, _ int) IObject {
	return o.gk.NewBool(frame, o.value)
}

// Falsy determines if the Bool's value is logically false by returning the negation of its `value` field.
func (o *Bool) Falsy() bool {
	return !o.value
}

// Equals checks if the current Bool object is equal to the provided IObject.
func (o *Bool) Equals(x IObject) bool {
	return o == x
}

// GobDecode decodes the Bool object from a byte slice encoded using Gob into its internal value.
func (o *Bool) GobDecode(b []byte) (err error) {
	o.value = b[0] == 1
	return
}

// GobEncode encodes the Bool instance into a byte slice representation. Returns the byte slice and any encoding error.
func (o *Bool) GobEncode() (b []byte, err error) {
	if o.value {
		b = []byte{1}
	} else {
		b = []byte{0}
	}
	return
}
