package objects

import (
	"bytes"
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

// Map represents a collection of key-Code pairs where keys are strings and Code implement the IObject interface.
type Map struct {
	IAllocator
	data map[string]IObject
}

// NewMap creates and returns a new instance of Map initialized with the provided map of string keys to IObject Code.
func newMap(allocator IAllocator, value map[string]IObject) IObject {
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
		IAllocator: allocator,
		data:       value,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *Map) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *Map) AsInterface() interface{} {
	res := make(map[string]interface{})
	for key, v := range o.data {
		res[key] = v.AsInterface()
	}
	return res
}

// AsBool converts the String object to a boolean. Returns true if the string has non-zero len, otherwise false.
func (o *Map) AsBool() bool {
	return len(o.data) > 0
}

// AsInt64 returns the len of the array as an int64 Code.
func (o *Map) AsInt64() int64 {
	return int64(len(o.data))
}

// AsFloat64 returns the len of the array as an int64 Code.
func (o *Map) AsFloat64() float64 {
	return float64(len(o.data))
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *Map) AsBytes() []byte {
	var res []byte
	for _, v := range o.data {
		res = append(res, v.AsBytes()...)
	}
	return res
}

// AsString returns the string representation of the Map object in the format of key-Code pairs enclosed in braces.
func (o *Map) AsString() string {
	var pairs []string
	for k, v := range o.data {
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
	o.data = target.data
	return nil
}

// Nil checks if the object is nil and always returns false.
func (o *Map) Nil() bool {
	return false
}

// LogicalOp performs a logical operation using the given operator and operand, returning an error for unsupported operators.
func (o *Map) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		ret, err := logicalOpNil(op)
		return o.GateKeeper().FromBoolError(ret, err)
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation on the Map using the specified operator and operand, returning an error.
func (o *Map) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// UnaryOp performs a unary operation on the Map object and always returns ErrInvalidOperator.
func (o *Map) UnaryOp(_ int, _ UnaryOperator) (IObject, error) {
	return nil, ErrInvalidOperator
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *Map) Call(_ int, _ ...IObject) (uint, IObject, error) {
	return 0, nil, nil
}

// Get retrieves the Code associated with the specified key from the map. If the key is not found, it returns nil.
func (o *Map) Get(key string) IObject {
	return o.data[key]
}

// Set assigns the specified Code to the given key in the Map. Overrides the Code if the key already exists.
func (o *Map) Set(key string, value IObject) {
	if len(o.data) > maxMapLen {
		return
	}
	o.data[key] = value
}

// Delete removes the entry associated with the specified key from the map.
func (o *Map) Delete(key string) {
	delete(o.data, key)
}

// Has checks if the specified key exists in the Map and returns true if found, otherwise false.
func (o *Map) Has(key string) bool {
	_, ok := o.data[key]
	return ok
}

// Values return the internal map of key-Code pairs stored in the Map object.
func (o *Map) Values() map[string]IObject {
	return o.data
}

// Length returns the number of key-Code pairs in the Map object.
func (o *Map) Length() int {
	return len(o.data)
}

// TypeName returns the string "map", which represents the type name of the Map object.
func (o *Map) TypeName() string {
	return MapType
}

// Copy creates and returns a deep copy of the Map object, duplicating all key-Code pairs recursively.
func (o *Map) Copy(frame int, depth int) IObject {
	c := make(map[string]IObject)
	for k, v := range o.data {
		if depth > maxDepth {
			break
		}
		c[k] = v.Copy(frame, depth+1)
	}
	return o.GateKeeper().NewMap(frame, c)
}

// Falsy returns true if the map contains no key-Code pairs, indicating it is empty.
func (o *Map) Falsy() bool {
	return len(o.data) == 0
}

// Equals checks if the Map is equal to another IObject by comparing their key-Code pairs. Returns true if equal.
func (o *Map) Equals(in IObject) bool {
	var xVal map[string]IObject
	switch x := in.(type) {
	case *Map:
		xVal = x.data
	default:
		return false
	}
	if len(o.data) != len(xVal) {
		return false
	}
	for k, v := range o.data {
		tv := xVal[k]
		if !v.Equals(tv) {
			return false
		}
	}
	return true
}

// IndexGet retrieves the Code associated with the given index in the map. Returns UndefinedValue if the index does not exist.
// An error is returned if the index type is invalid.
func (o *Map) IndexGet(_ int, index IObject) (IObject, error) {
	strIdx := index.AsString()
	res, ok := o.data[strIdx]
	if !ok {
		return o.GateKeeper().UndefinedValue(), nil
	}
	return res, nil
}

// IndexSet sets the specified Code at the given string-convertible index in the Map. Returns an error for invalid index types.
func (o *Map) IndexSet(index, value IObject) error {
	strIdx := index.AsString()
	o.data[strIdx] = value
	return nil
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *Map) Count() int {
	counter := 0
	for _, v := range o.data {
		counter += v.Count()
	}
	return counter
}

// Iterate creates and returns an iterator for the Map, allowing iteration over its keys and associated Code.
func (o *Map) Iterate(frame int) IIterator {
	return o.GateKeeper().NewMapIterator(frame, o.data, 0)
}

// Iterable returns true, indicating that the Map object supports iteration over its elements.
func (o *Map) Iterable() bool {
	return true
}

// GobEncode serializes the Map's data into a byte slice using gob encoding and returns the result or an error.
func (o *Map) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(o.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode decodes the provided byte slice into the Map's data field using the gob package.
func (o *Map) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&o.data); err != nil {
		return err
	}
	return nil
}
