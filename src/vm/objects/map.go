package objects

import (
	"encoding/gob"
	"fmt"
	"strings"
)

const (
	MapType = "map"
)

func init() {
	gob.Register(&Map{})
}

// Map represents a collection of key-values pairs where keys are strings and values implement the IObject interface.
type Map struct {
	Allocator
	values map[string]IObject
}

// NewMap creates and returns a new instance of Map initialized with the provided map of string keys to IObject values.
func newMap(factory IGateKeeper, frame int, value map[string]IObject) IObject {
	if len(value) > maxMapLen {
		nv := make(map[string]IObject)
		for k, v := range value {
			if len(nv) > maxMapLen {
				break
			}
			nv[k] = v
		}
		value = nv
	}
	return &Map{
		Allocator: Allocator{gk: factory, frame: frame},
		values:    value,
	}
}

// AsBool converts the String object to a boolean. Returns true if the string has non-zero length, otherwise false.
func (o *Map) AsBool() bool {
	return len(o.values) > 0
}

// AsInt64 returns the length of the array as an int64 value.
func (o *Map) AsInt64() int64 {
	return int64(len(o.values))
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *Map) AsFloat64() float64 {
	return float64(len(o.values))
}

// AsString returns the string representation of the Map object in the format of key-values pairs enclosed in braces.
func (o *Map) AsString() string {
	var pairs []string
	for k, v := range o.values {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.AsString()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

// AssignValue assigns the elements of another Map to the current Map if the input is of type *Map, otherwise returns an error.
func (o *Map) AssignValue(v IObject) error {
	target, ok := v.(*Map)
	if !ok {
		return ErrNotAssignable
	}
	o.values = target.values
	return nil
}

// Nil checks if the object is nil and always returns false.
func (o *Map) Nil() bool {
	return false
}

// LogicalOp performs a logical operation using the given operator and operand, returning an error for unsupported operators.
func (o *Map) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.gk, op)
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation on the Map using the specified operator and operand, returning an error.
func (o *Map) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Map) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *Map) CanCall() bool {
	return false
}

// Get retrieves the values associated with the specified key from the map. If the key is not found, it returns nil.
func (o *Map) Get(key string) IObject {
	return o.values[key]
}

// Set assigns the specified values to the given key in the Map. Overrides the values if the key already exists.
func (o *Map) Set(key string, value IObject) {
	if len(o.values) > maxMapLen {
		return
	}
	o.values[key] = value
}

// Delete removes the entry associated with the specified key from the map.
func (o *Map) Delete(key string) {
	delete(o.values, key)
}

// Has checks if the specified key exists in the Map and returns true if found, otherwise false.
func (o *Map) Has(key string) bool {
	_, ok := o.values[key]
	return ok
}

// Values return the internal map of key-values pairs stored in the Map object.
func (o *Map) Values() map[string]IObject {
	return o.values
}

// Length returns the number of key-values pairs in the Map object.
func (o *Map) Length() int {
	return len(o.values)
}

// TypeName returns the string "map", which represents the type name of the Map object.
func (o *Map) TypeName() string {
	return MapType
}

// Copy creates and returns a deep copy of the Map object, duplicating all key-values pairs recursively.
func (o *Map) Copy(frame int, depth int) IObject {
	c := make(map[string]IObject)
	for k, v := range o.values {
		if depth > maxDepth {
			break
		}
		c[k] = v.Copy(frame, depth+1)
	}
	return o.GateKeeper().NewMap(frame, c)
}

// Falsy returns true if the map contains no key-values pairs, indicating it is empty.
func (o *Map) Falsy() bool {
	return len(o.values) == 0
}

// Equals checks if the Map is equal to another IObject by comparing their key-values pairs. Returns true if equal.
func (o *Map) Equals(in IObject) bool {
	var xVal map[string]IObject
	switch x := in.(type) {
	case *Map:
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

// IndexGet retrieves the values associated with the given index in the map. Returns UndefinedValue if the index does not exist.
// An error is returned if the index type is invalid.
func (o *Map) IndexGet(_ int, index IObject) (res IObject, err error) {
	strIdx, ok := o.GateKeeper().ToString(index)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	res, ok = o.values[strIdx]
	if !ok {
		res = o.GateKeeper().UndefinedValue()
	}
	return
}

// IndexSet sets the specified values at the given string-convertible index in the Map. Returns an error for invalid index types.
func (o *Map) IndexSet(index, value IObject) (err error) {
	strIdx, ok := o.GateKeeper().ToString(index)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	o.values[strIdx] = value
	return nil
}

// Iterate creates and returns an iterator for the Map, allowing iteration over its keys and associated values.
func (o *Map) Iterate(frame int) IIterator {
	return o.GateKeeper().NewMapIterator(frame, o.values, 0)
}

// CanIterate returns true, indicating that the Map object supports iteration over its elements.
func (o *Map) CanIterate() bool {
	return true
}
