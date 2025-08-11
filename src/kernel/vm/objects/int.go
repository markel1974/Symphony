package objects

import (
	"strconv"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
	"github.com/markel1974/c64emu/src/kernel/vm/tokens"
)

// Int represents an integer value.
type Int struct {
	ObjectImpl
	Value int64
}

func (o *Int) String() string {
	return strconv.FormatInt(o.Value, 10)
}

// TypeName returns the name of the type.
func (o *Int) TypeName() string {
	return "int"
}

// BinaryOp returns another object that is the result of a given binary operator and a right-hand side object.
func (o *Int) BinaryOp(op tokens.Token, rhs Object) (Object, error) {
	switch rhs := rhs.(type) {
	case *Int:
		switch op {
		case tokens.Add:
			r := o.Value + rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case tokens.Sub:
			r := o.Value - rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case tokens.Mul:
			r := o.Value * rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case tokens.Quo:
			r := o.Value / rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case tokens.Rem:
			r := o.Value % rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case tokens.And:
			r := o.Value & rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case tokens.Or:
			r := o.Value | rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case tokens.Xor:
			r := o.Value ^ rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case tokens.AndNot:
			r := o.Value &^ rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case tokens.Shl:
			r := o.Value << uint64(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case tokens.Shr:
			r := o.Value >> uint64(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
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
	case *Float:
		switch op {
		case tokens.Add:
			return &Float{Value: float64(o.Value) + rhs.Value}, nil
		case tokens.Sub:
			return &Float{Value: float64(o.Value) - rhs.Value}, nil
		case tokens.Mul:
			return &Float{Value: float64(o.Value) * rhs.Value}, nil
		case tokens.Quo:
			return &Float{Value: float64(o.Value) / rhs.Value}, nil
		case tokens.Less:
			if float64(o.Value) < rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case tokens.Greater:
			if float64(o.Value) > rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case tokens.LessEq:
			if float64(o.Value) <= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case tokens.GreaterEq:
			if float64(o.Value) >= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		default:
			return nil, errors.ErrInvalidOperator
		}
	case *Char:
		switch op {
		case tokens.Add:
			return &Char{Value: rune(o.Value) + rhs.Value}, nil
		case tokens.Sub:
			return &Char{Value: rune(o.Value) - rhs.Value}, nil
		case tokens.Less:
			if o.Value < int64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case tokens.Greater:
			if o.Value > int64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case tokens.LessEq:
			if o.Value <= int64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case tokens.GreaterEq:
			if o.Value >= int64(rhs.Value) {
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
func (o *Int) Copy() Object {
	return &Int{Value: o.Value}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Int) IsFalsy() bool {
	return o.Value == 0
}

// Equals returns true if the value of the type is equal to the value of another object.
func (o *Int) Equals(x Object) bool {
	t, ok := x.(*Int)
	if !ok {
		return false
	}
	return o.Value == t.Value
}

// ToInt64 will try to convert object o to int64 value.
func ToInt64(o Object) (v int64, ok bool) {
	switch o := o.(type) {
	case *Int:
		v = o.Value
		ok = true
	case *Float:
		v = int64(o.Value)
		ok = true
	case *Char:
		v = int64(o.Value)
		ok = true
	case *Bool:
		if o == TrueValue {
			v = 1
		}
		ok = true
	case *String:
		c, err := strconv.ParseInt(o.Value, 10, 64)
		if err == nil {
			v = c
			ok = true
		}
	}
	return
}

// ToInt will try to convert object o to int value.
func ToInt(o Object) (v int, ok bool) {
	switch o := o.(type) {
	case *Int:
		v = int(o.Value)
		ok = true
	case *Float:
		v = int(o.Value)
		ok = true
	case *Char:
		v = int(o.Value)
		ok = true
	case *Bool:
		if o == TrueValue {
			v = 1
		}
		ok = true
	case *String:
		c, err := strconv.ParseInt(o.Value, 10, 64)
		if err == nil {
			v = int(c)
			ok = true
		}
	}
	return
}
