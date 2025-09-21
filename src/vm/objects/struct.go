package objects

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"
)

func init() {
	gob.Register(&Struct{})
}

// Struct is a composite object that implements the IObject interface and stores a collection of key-Code pairs.
type Struct struct {
	IAllocator
	typeName string
	data     map[string]IObject
}

// NewStruct creates a new instance of MapImmutable with the provided map of string keys and IObject Code.
func newStruct(allocator IAllocator, value map[string]IObject) IObject {
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
		IAllocator: allocator,
		typeName:   "",
		data:       value,
	}
}

// setAllocator sets the allocator for the instance, defining its memory management and lifecycle behavior.
func (o *Struct) setAllocator(allocator IAllocator) {
	o.IAllocator = allocator
}

// AsInterface converts the object into a generic interface{} type and returns the underlying data.
func (o *Struct) AsInterface() interface{} {
	res := make(map[string]interface{})
	for key, v := range o.data {
		res[key] = v.AsInterface()
	}
	return res
}

// AsBool converts the Struct to a boolean, returning true if the Struct contains at least one key-Code pair; otherwise false.
func (o *Struct) AsBool() bool {
	return len(o.data) > 0
}

// AsInt64 returns the len of the array as an int64 Code.
func (o *Struct) AsInt64() int64 {
	return int64(len(o.data))
}

// AsFloat64 returns the len of the array as an int64 Code.
func (o *Struct) AsFloat64() float64 {
	return float64(len(o.data))
}

// AsBytes converts the object elements into a single concatenated slice of bytes by calling AsBytes on each element.
func (o *Struct) AsBytes() []byte {
	var res []byte
	for _, v := range o.data {
		res = append(res, v.AsBytes()...)
	}
	return res
}

// AsString returns a string representation of the Struct object, formatting its key-Code pairs into a map-like structure.
func (o *Struct) AsString() string {
	var pairs []string
	for k, v := range o.data {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.AsString()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

// AssignValue assigns the elements of another Struct to the current Struct if the input is of type *Struct, otherwise returns an error.
func (o *Struct) AssignValue(v IObject) error {
	switch v := v.(type) {
	case *Struct:
		o.data = v.data
		return nil
	case *Map:
		o.data = v.data
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
		return logicalOpNil(o.GateKeeper(), op)
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

// Values returns the underlying map of string keys to IObject Code contained within the Struct.
func (o *Struct) Values() map[string]IObject {
	return o.data
}

// Length returns the number of key-Code pairs stored in the Struct.
func (o *Struct) Length() int {
	return len(o.data)
}

// SetValue sets the specified key to the given Code in the Code map of the Struct.
func (o *Struct) SetValue(k string, v IObject) {
	if len(o.data) > maxStructLen {
		return
	}
	o.data[k] = v
}

// GetValue retrieves the Code associated with the given key in the Code map and a boolean indicating its presence.
func (o *Struct) GetValue(k string) (IObject, bool) {
	v, ok := o.data[k]
	return v, ok
}

// TypeName returns the type name of the object as a string.
func (o *Struct) TypeName() string {
	return o.typeName
}

// Copy creates and returns a new IObject by duplicating the internal state of the Struct instance.
func (o *Struct) Copy(frame int, depth int) IObject {
	c := make(map[string]IObject)
	for k, v := range o.data {
		if depth >= maxDepth {
			break
		}
		c[k] = v.Copy(frame, depth+1)
	}
	return o.GateKeeper().NewStruct(frame, o.typeName, c)
}

// Falsy returns true if the Struct contains no Code, otherwise false.
func (o *Struct) Falsy() bool {
	return len(o.data) == 0
}

// IndexGet retrieves the Code associated with the given index within the Struct. Returns an error for invalid index types.
func (o *Struct) IndexGet(_ int, index IObject) (IObject, error) {
	strIdx := index.AsString()
	res, ok := o.data[strIdx]
	if !ok {
		res = o.GateKeeper().UndefinedValue()
	}
	return res, nil
}

// IndexSet updates or assigns a Code to the specified index within the Struct. Returns an error for invalid index types.
func (o *Struct) IndexSet(index, value IObject) error {
	strIdx := index.AsString()
	o.data[strIdx] = value
	return nil
}

// Equals checks if the current Struct is equal to another IObject by comparing their key-Code pairs and lengths.
func (o *Struct) Equals(in IObject) bool {
	var xVal map[string]IObject
	switch x := in.(type) {
	case *Struct:
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

// Iterate returns an IIterator for traversing the key-Code pairs in the Struct's internal map.
func (o *Struct) Iterate(frame int) IIterator {
	return o.GateKeeper().NewStructIterator(frame, o.data, 0)
}

// Iterable checks if the object can be iterated over. Always returns true for this implementation.
func (o *Struct) Iterable() bool {
	return true
}

// Count returns the total number of elements in the instance and its sub-elements.
func (o *Struct) Count() int {
	counter := 0
	for _, v := range o.data {
		counter += v.Count()
	}
	return counter
}

// GobEncode serializes the Struct's data into a byte slice using gob encoding and returns the result or an error.
func (o *Struct) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(o.typeName); err != nil {
		return nil, err
	}
	if err := encoder.Encode(o.data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode decodes the provided byte slice into the Struct's data field using the gob package.
func (o *Struct) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&o.typeName); err != nil {
		return err
	}
	if err := decoder.Decode(&o.data); err != nil {
		return err
	}
	return nil
}
