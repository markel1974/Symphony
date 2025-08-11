package objects

// TrueValue represents the boolean true values as an IObject.
// FalseValue represents the boolean false values as an IObject.
// UndefinedValue represents an undefined values as an IObject.
var (
	// TrueValue represents a true values.
	TrueValue IObject = &Bool{value: true}

	// FalseValue represents a false values.
	FalseValue IObject = &Bool{value: false}

	// UndefinedValue represents an undefined values.
	UndefinedValue IObject = &Undefined{}
)

// Bool represents a boolean type with true or false values.
// It embeds ObjectImpl and provides methods for common object operations.
// The values field indicates the boolean state: `true` or `false`.
type Bool struct {
	ObjectImpl
	value bool
}

// String returns the string representation of the Bool object, "true" for true values and "false" for false values.
func (o *Bool) String() string {
	if o.value {
		return "true"
	}
	return "false"
}

// TypeName returns the string "bool", representing the name of the type of the Bool object.
func (o *Bool) TypeName() string {
	return "bool"
}

// Copy returns the Bool object itself as it is considered immutable.
func (o *Bool) Copy() IObject {
	return o
}

// Falsy determines if the Bool object's values should be considered falsy (returns true if the values is false).
func (o *Bool) Falsy() bool {
	return !o.value
}

// Equals checks whether the Bool object is equal to another IObject based on reference comparison.
func (o *Bool) Equals(x IObject) bool {
	return o == x
}

// GobDecode implements the gob.GobDecoder interface, decoding a byte slice to set the Bool's values.
func (o *Bool) GobDecode(b []byte) (err error) {
	o.value = b[0] == 1
	return
}

// GobEncode encodes the Bool object into a byte slice for serialization, representing true as 1 and false as 0.
func (o *Bool) GobEncode() (b []byte, err error) {
	if o.value {
		b = []byte{1}
	} else {
		b = []byte{0}
	}
	return
}

// ToBool converts an IObject to a boolean by checking its falsy state. Returns the boolean values and a success flag.
func ToBool(o IObject) (v bool, ok bool) {
	ok = true
	v = !o.Falsy()
	return
}

// FromBool converts a boolean values into an IObject, returning TrueValue for true and FalseValue for false.
func FromBool(v bool) IObject {
	if v {
		return TrueValue
	}
	return FalseValue
}
