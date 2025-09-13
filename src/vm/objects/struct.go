package objects

import (
	"encoding/gob"
	"fmt"
	"strings"
)

func init() {
	gob.Register(&Struct{})
}

// Struct is a composite object that implements the IObject interface and stores a collection of key-value pairs.
type Struct struct {
	Allocator
	typeName string
	values   map[string]IObject
}

// NewStruct creates a new instance of MapImmutable with the provided map of string keys and IObject values.
func newStruct(factory IGateKeeper, frame int, value map[string]IObject) IObject {
	if len(value) > maxStructLen {
		nv := make(map[string]IObject)
		for k, v := range value {
			if len(nv) > maxStructLen {
				break
			}
			nv[k] = v
		}
		value = nv
	}
	return &Struct{
		Allocator: Allocator{gk: factory, frame: frame},
		typeName:  "",
		values:    value,
	}
}

// AsBool converts the Struct to a boolean, returning true if the Struct contains at least one key-value pair; otherwise false.
func (o *Struct) AsBool() bool {
	return len(o.values) > 0
}

// AsInt64 returns the length of the array as an int64 value.
func (o *Struct) AsInt64() int64 {
	return int64(len(o.values))
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *Struct) AsFloat64() float64 {
	return float64(len(o.values))
}

// AsString returns a string representation of the Struct object, formatting its key-value pairs into a map-like structure.
func (o *Struct) AsString() string {
	var pairs []string
	for k, v := range o.values {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.AsString()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

// AssignValue assigns the elements of another Struct to the current Struct if the input is of type *Struct, otherwise returns an error.
func (o *Struct) AssignValue(v IObject) error {
	switch v := v.(type) {
	case *Struct:
		o.values = v.values
		return nil
	case *Map:
		o.values = v.values
	default:
		return ErrNotAssignable
	}
	return nil
}

// Nil checks if the object is nil and always returns false.
func (o *Struct) Nil() bool {
	return false
}

// LogicalOp performs a logical operation using the specified operator and operand, returning the result or an error.
func (o *Struct) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.gk, op)
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation using the given operator and operand, returning the result or an error.
func (o *Struct) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Struct) Call(_ int, _ ...IObject) (retCount uint, ret IObject, err error) {
	return 0, nil, nil
}

// Values returns the underlying map of string keys to IObject values contained within the Struct.
func (o *Struct) Values() map[string]IObject {
	return o.values
}

// Length returns the number of key-value pairs stored in the Struct.
func (o *Struct) Length() int {
	return len(o.values)
}

// SetValue sets the specified key to the given value in the values map of the Struct.
func (o *Struct) SetValue(k string, v IObject) {
	if len(o.values) > maxStructLen {
		return
	}
	o.values[k] = v
}

// GetValue retrieves the value associated with the given key in the values map and a boolean indicating its presence.
func (o *Struct) GetValue(k string) (IObject, bool) {
	v, ok := o.values[k]
	return v, ok
}

// TypeName returns the type name of the object as a string.
func (o *Struct) TypeName() string {
	return o.typeName
}

// Copy creates and returns a new IObject by duplicating the internal state of the Struct instance.
func (o *Struct) Copy(frame int, depth int) IObject {
	c := make(map[string]IObject)
	for k, v := range o.values {
		if depth >= maxDepth {
			break
		}
		c[k] = v.Copy(frame, depth+1)
	}
	return o.GateKeeper().NewStruct(frame, o.typeName, c)
}

// Falsy returns true if the Struct contains no values, otherwise false.
func (o *Struct) Falsy() bool {
	return len(o.values) == 0
}

// IndexGet retrieves the value associated with the given index within the Struct. Returns an error for invalid index types.
func (o *Struct) IndexGet(_ int, index IObject) (IObject, error) {
	strIdx, ok := o.GateKeeper().ToString(index)
	if !ok {
		return nil, ErrInvalidIndexType
	}
	res, ok := o.values[strIdx]
	if !ok {
		res = o.GateKeeper().UndefinedValue()
	}
	return res, nil
}

// IndexSet updates or assigns a value to the specified index within the Struct. Returns an error for invalid index types.
func (o *Struct) IndexSet(index, value IObject) error {
	strIdx, ok := o.GateKeeper().ToString(index)
	if !ok {
		return ErrInvalidIndexType
	}
	o.values[strIdx] = value
	return nil
}

// Equals checks if the current Struct is equal to another IObject by comparing their key-value pairs and lengths.
func (o *Struct) Equals(in IObject) bool {
	var xVal map[string]IObject
	switch x := in.(type) {
	case *Struct:
		xVal = x.values
	default:
		return false
	}
	if len(o.values) != len(xVal) {
		return false
	}
	for k, v := range o.values {
		tv := xVal[k]
		if !v.Equals(tv) {
			return false
		}
	}
	return true
}

// Iterate returns an IIterator for traversing the key-value pairs in the Struct's internal map.
func (o *Struct) Iterate(frame int) IIterator {
	return o.GateKeeper().NewStructIterator(frame, o.values, 0)
}

// Iterable checks if the object can be iterated over. Always returns true for this implementation.
func (o *Struct) Iterable() bool {
	return true
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *Struct) Count() int {
	counter := 0
	for _, v := range o.values {
		counter += v.Count()
	}
	return counter
}
