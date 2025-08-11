package objects

import (
	"github.com/markel1974/c64emu/src/kernel/vm/errors"
	"github.com/markel1974/c64emu/src/kernel/vm/tokens"
)

// Char represents a character value.
type Char struct {
	ObjectImpl
	Value rune
}

func (o *Char) String() string {
	return string(o.Value)
}

// TypeName returns the name of the type.
func (o *Char) TypeName() string {
	return "char"
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *Char) BinaryOp(op tokens.Token, rhs Object) (Object, error) {
	switch rhs := rhs.(type) {
	case *Char:
		switch op {
		case tokens.Add:
			r := o.Value + rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Char{Value: r}, nil
		case tokens.Sub:
			r := o.Value - rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Char{Value: r}, nil
		case tokens.Less:
			if o.Value < rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case tokens.Greater:
			if o.Value > rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case tokens.LessEq:
			if o.Value <= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case tokens.GreaterEq:
			if o.Value >= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		default:
			return nil, errors.ErrInvalidOperator
		}
	case *Int:
		switch op {
		case tokens.Add:
			r := o.Value + rune(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Char{Value: r}, nil
		case tokens.Sub:
			r := o.Value - rune(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Char{Value: r}, nil
		case tokens.Less:
			if int64(o.Value) < rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case tokens.Greater:
			if int64(o.Value) > rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case tokens.LessEq:
			if int64(o.Value) <= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case tokens.GreaterEq:
			if int64(o.Value) >= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		default:
			return nil, errors.ErrInvalidOperator
		}
	}
	return nil, errors.ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *Char) Copy() Object {
	return &Char{Value: o.Value}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Char) IsFalsy() bool {
	return o.Value == 0
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Char) Equals(x Object) bool {
	t, ok := x.(*Char)
	if !ok {
		return false
	}
	return o.Value == t.Value
}
