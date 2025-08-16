package objects

const (
	BoolType = "bool"
)

// TrueValue is a predefined constant representing the boolean true value as an IObject.
// FalseValue is a predefined constant representing the boolean false value as an IObject.
// UndefinedValue is a predefined constant representing an undefined value as an IObject.
var (
	// TrueValue represents a true values.
	TrueValue IObject = &Bool{value: true}

	// FalseValue represents a false values.
	FalseValue IObject = &Bool{value: false}
)

// Bool is a custom type representing a boolean values, implementing IObject interface and encapsulating a boolean value.
type Bool struct {
	Object
	value bool
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
func (o *Bool) Copy() IObject {
	return o
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

// ToBool converts the given IObject to a bool based on its Boolean() method and returns the result along with a success flag.
func ToBool(o IObject) (v bool, ok bool) {
	ok = true
	v = !o.Boolean()
	return
}

// FromBool converts a boolean values into its corresponding IObject representation, returning TrueValue or FalseValue.
func FromBool(v bool) IObject {
	if v {
		return TrueValue
	}
	return FalseValue
}

// ToBoolArg converts the given IObject to a boolean if possible or returns an error indicating an invalid argument type.
func ToBoolArg(index int, o IObject) (bool, error) {
	b1, ok := o.(*Bool)
	if !ok {
		return false, NewInvalidArgumentError(index, "bool(compatible)", o.TypeName())
	}
	return b1.value, nil
}
