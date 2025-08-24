package objects

const (
	CharType = "char"
)

// Char represents a character type, encapsulating a single rune values and inheriting behavior from Object.
type Char struct {
	gk    IGateKeeper
	frame int
	value rune
}

// NewChar creates and returns a new Char object with the specified rune values.
func newChar(factory IGateKeeper, frame int, value rune) IObject {
	return &Char{
		gk:    factory,
		frame: frame,
		value: value,
	}
}

// GateKeeper returns a reference to the GateKeeper associated with the Object.
func (o *Char) GateKeeper() IGateKeeper {
	return o.gk
}

// Frame returns the current frame value of the Object.
func (o *Char) Frame() int {
	return o.frame
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *Char) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrNotIndexAssignable,
// as this operation is unsupported.
func (o *Char) IndexSet(_, _ IObject) (err error) {
	return ErrNotIndexAssignable
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *Char) Iterate(_ int) IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *Char) CanIterate() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Char) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *Char) CanCall() bool {
	return false
}

// Length returns the length of the Int object.
func (o *Char) Length() int {
	return 0
}

// Value returns the rune values stored in the Char object.
func (o *Char) Value() rune {
	return o.value
}

// String returns the string representation of the Char object's values.
func (o *Char) String() string {
	return string(o.value)
}

// TypeName returns the name of the type as a string.
func (o *Char) TypeName() string {
	return CharType
}

// BinaryOp performs a binary operation between the Char object and another IObject using the specified operator.
func (o *Char) BinaryOp(frame int, op Operator, in IObject) (IObject, error) {
	switch rhs := in.(type) {
	case *Char:
		switch op {
		case OperatorAdd:
			r := o.value + rhs.value
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewChar(frame, r), nil
		case OperatorSub:
			r := o.value - rhs.value
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewChar(frame, r), nil
		case OperatorLess:
			if o.value < rhs.value {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorGreater:
			if o.value > rhs.value {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorLessEq:
			if o.value <= rhs.value {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorGreaterEq:
			if o.value >= rhs.value {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		default:
			return nil, ErrInvalidOperator
		}
	case *Int:
		switch op {
		case OperatorAdd:
			r := o.value + rune(rhs.value)
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewChar(frame, r), nil
		case OperatorSub:
			r := o.value - rune(rhs.value)
			if r == o.value {
				return o, nil
			}
			return o.GateKeeper().NewChar(frame, r), nil
		case OperatorLess:
			if int64(o.Value()) < rhs.value {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorGreater:
			if int64(o.Value()) > rhs.value {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorLessEq:
			if int64(o.Value()) <= rhs.value {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		case OperatorGreaterEq:
			if int64(o.Value()) >= rhs.value {
				return o.GateKeeper().TrueValue(), nil
			}
			return o.GateKeeper().FalseValue(), nil
		default:
			return nil, ErrInvalidOperator
		}
	}
	return nil, ErrInvalidOperator
}

// Copy creates and returns a new instance of the Char object with the same values.
func (o *Char) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewChar(frame, o.value)
}

// Boolean checks whether the Char object represents a falsy state, returning true if the underlying values is 0.
func (o *Char) Boolean() bool {
	return o.value == 0
}

// Equals checks if the current Char object is equal to another IObject. Returns true if both objects are equal.
func (o *Char) Equals(x IObject) bool {
	t, ok := x.(*Char)
	if !ok {
		return false
	}
	return o.value == t.value
}
