package objects

var (
	// TrueValue represents a true value.
	TrueValue Object = &Bool{value: true}

	// FalseValue represents a false value.
	FalseValue Object = &Bool{value: false}

	// UndefinedValue represents an undefined value.
	UndefinedValue Object = &Undefined{}
)

// Bool represents a boolean value.
type Bool struct {
	ObjectImpl
	// this is intentionally non-public to force using objects.TrueValue and FalseValue always
	value bool
}

func (o *Bool) String() string {
	if o.value {
		return "true"
	}

	return "false"
}

// TypeName returns the name of the type.
func (o *Bool) TypeName() string {
	return "bool"
}

// Copy returns a copy of the type.
func (o *Bool) Copy() Object {
	return o
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Bool) IsFalsy() bool {
	return !o.value
}

// Equals returns true if the value of the type is equal to the value of another object.
func (o *Bool) Equals(x Object) bool {
	return o == x
}

// GobDecode decodes bool value from input bytes.
func (o *Bool) GobDecode(b []byte) (err error) {
	o.value = b[0] == 1
	return
}

// GobEncode encodes bool values into bytes.
func (o *Bool) GobEncode() (b []byte, err error) {
	if o.value {
		b = []byte{1}
	} else {
		b = []byte{0}
	}
	return
}

// ToBool will try to convert object o to bool value.
func ToBool(o Object) (v bool, ok bool) {
	ok = true
	v = !o.IsFalsy()
	return
}

func FromBool(v bool) Object {
	if v {
		return TrueValue
	}
	return FalseValue
}
