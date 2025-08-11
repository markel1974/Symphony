package objects

import (
	"fmt"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
)

// Map represents a map of objects.
type Map struct {
	ObjectImpl
	Value map[string]Object
}

// TypeName returns the name of the type.
func (o *Map) TypeName() string {
	return "map"
}

func (o *Map) String() string {
	var pairs []string
	for k, v := range o.Value {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

// Copy returns a copy of the type.
func (o *Map) Copy() Object {
	c := make(map[string]Object)
	for k, v := range o.Value {
		c[k] = v.Copy()
	}
	return &Map{Value: c}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Map) IsFalsy() bool {
	return len(o.Value) == 0
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Map) Equals(x Object) bool {
	var xVal map[string]Object
	switch x := x.(type) {
	case *Map:
		xVal = x.Value
	case *ImmutableMap:
		xVal = x.Value
	default:
		return false
	}
	if len(o.Value) != len(xVal) {
		return false
	}
	for k, v := range o.Value {
		tv := xVal[k]
		if !v.Equals(tv) {
			return false
		}
	}
	return true
}

// IndexGet returns the value for the given key.
func (o *Map) IndexGet(index Object) (res Object, err error) {
	strIdx, ok := ToString(index)
	if !ok {
		err = errors.ErrInvalidIndexType
		return
	}
	res, ok = o.Value[strIdx]
	if !ok {
		res = UndefinedValue
	}
	return
}

// IndexSet sets the value for the given key.
func (o *Map) IndexSet(index, value Object) (err error) {
	strIdx, ok := ToString(index)
	if !ok {
		err = errors.ErrInvalidIndexType
		return
	}
	o.Value[strIdx] = value
	return nil
}

// Iterate creates a map iterator.
func (o *Map) Iterate() Iterator {
	var keys []string
	for k := range o.Value {
		keys = append(keys, k)
	}
	return &MapIterator{
		v: o.Value,
		k: keys,
		l: len(keys),
	}
}

// CanIterate returns whether the Object can be Iterated.
func (o *Map) CanIterate() bool {
	return true
}

// MapIterator represents an iterator for the map.
type MapIterator struct {
	ObjectImpl
	v map[string]Object
	k []string
	i int
	l int
}

// TypeName returns the name of the type.
func (i *MapIterator) TypeName() string {
	return "map-iterator"
}

func (i *MapIterator) String() string {
	return "<map-iterator>"
}

// IsFalsy returns true if the value of the type is falsy.
func (i *MapIterator) IsFalsy() bool {
	return true
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (i *MapIterator) Equals(Object) bool {
	return false
}

// Copy returns a copy of the type.
func (i *MapIterator) Copy() Object {
	return &MapIterator{v: i.v, k: i.k, i: i.i, l: i.l}
}

// Next returns true if there are more elements to iterate.
func (i *MapIterator) Next() bool {
	i.i++
	return i.i <= i.l
}

// Key returns the key or index value of the current element.
func (i *MapIterator) Key() Object {
	k := i.k[i.i-1]
	return &String{Value: k}
}

// Value returns the value of the current element.
func (i *MapIterator) Value() Object {
	k := i.k[i.i-1]
	return i.v[k]
}

func ToMap(o Object) (res map[string]interface{}) {
	switch o := o.(type) {
	case *Map:
		res = make(map[string]interface{})
		for key, v := range o.Value {
			res[key] = ToInterface(v)
		}
	}
	return
}
