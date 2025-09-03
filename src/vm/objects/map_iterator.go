package objects

import "encoding/gob"

const (
	MapIteratorType  = "map_iterator"
	MapIteratorLabel = "<" + MapIteratorType + ">"
)

func init() {
	gob.Register(&MapIterator{})
}

// MapIterator is a type used for iterating over key-values pairs in a map-like structure.
// It implements the IIterator interface and provides methods for traversal and element access.
// The embedded Object provides default implementations for methods from the IObject interface.
// The internal state includes the map (values), keys (keys), current position index (index), and total keys length (length).
type MapIterator struct {
	Allocator
	values map[string]IObject
	keys   []string
	index  int
	length int
}

// NewMapIterator creates and returns a new instance of MapIterator.
func newMapIterator(factory IGateKeeper, frame int, v map[string]IObject, index int) IIterator {
	var keys []string
	for k := range v {
		keys = append(keys, k)
	}
	return &MapIterator{
		Allocator: Allocator{gk: factory, frame: frame},
		values:    v,
		keys:      keys,
		length:    len(keys),
		index:     index,
	}
}

// AsBool returns true if the map is not empty, otherwise false.
func (o *MapIterator) AsBool() bool {
	return o.length > 0
}

// AsInt64 returns the length of the array as an int64 value.
func (o *MapIterator) AsInt64() int64 {
	return int64(o.length)
}

// AsFloat64 returns the length of the array as an int64 value.
func (o *MapIterator) AsFloat64() float64 {
	return float64(o.length)
}

// AsString returns the string representation of the MapIterator.
func (o *MapIterator) AsString() string {
	return MapIteratorLabel
}

// AssignValue sets the current object to the provided IObject, returning ErrNotAssignable if the operation is not supported.
func (o *MapIterator) AssignValue(_ IObject) error {
	return ErrNotAssignable
}

// Nil checks if the object is nil and always returns false.
func (o *MapIterator) Nil() bool {
	return false
}

// LogicalOp performs a logical operation on the MapIterator object using the specified operator and operand.
// Returns nil and ErrInvalidOperator as logical operations are unsupported for MapIterator.
func (o *MapIterator) LogicalOp(_ int, op LogicalOperator, rhsIn IObject) (IObject, error) {
	if rhsIn.Nil() {
		return logicalOpNil(o.gk, op)
	}
	return nil, ErrInvalidOperator
}

// ArithmeticOp performs an arithmetic operation specified by the operator on the current object and the given operand.
// Returns an error if the operation is unsupported.
func (o *MapIterator) ArithmeticOp(_ int, _ ArithmeticOperator, _ IObject) (IObject, error) {
	return nil, ErrInvalidOperator
}

// IndexGet attempts to retrieve a value at the given index and returns an error if the object is not indexable.
func (o *MapIterator) IndexGet(_ int, _ IObject) (res IObject, err error) {
	return nil, ErrNotIndexable
}

// IndexSet attempts to assign a value to an index in the object but always returns ErrUnsupportedIndex,
// as this operation is unsupported.
func (o *MapIterator) IndexSet(_, _ IObject) (err error) {
	return ErrUnsupportedIndex
}

// Iterate returns an IIterator to traverse over the elements of the object. If iteration is not supported, it returns nil.
func (o *MapIterator) Iterate(_ int) IIterator {
	return nil
}

// CanIterate determines if the object can be iterated over and returns false for this implementation.
func (o *MapIterator) CanIterate() bool {
	return false
}

// Call invokes the Object with the provided arguments, returning a result object and an error, if any.
func (o *MapIterator) Call(_ int, _ ...IObject) (ret IObject, err error) {
	return nil, nil
}

// CanCall determines if the object can be invoked as a callable. Returns false for non-callable objects.
func (o *MapIterator) CanCall() bool {
	return false
}

// Length returns the length of the Int object.
func (o *MapIterator) Length() int {
	return 0
}

// Copy creates and returns a new instance of MapIterator, duplicating its current state.
func (o *MapIterator) Copy(frame int, _ int) IObject {
	ret := o.GateKeeper().NewMapIterator(frame, o.values, o.index)
	return ret
}

// TypeName returns the type name of the MapIterator as a string.
func (o *MapIterator) TypeName() string {
	return MapIteratorType
}

// Falsy returns true, indicating the MapIterator is considered falsy in a boolean context.
func (o *MapIterator) Falsy() bool {
	return true
}

// Equals determine if the current MapIterator is equal to another IObject. Returns false by default.
func (o *MapIterator) Equals(IObject) bool {
	return false
}

// Next advances the iterator to the next element and returns true if there are more elements to iterate over.
func (o *MapIterator) Next() bool {
	o.index++
	return o.index <= o.length
}

// Key retrieves the key of the current element in the iteration as an IObject.
func (o *MapIterator) Key(frame int) IObject {
	idx := o.index - 1
	if idx < 0 || idx >= o.length {
		return o.GateKeeper().UndefinedValue()
	}
	k := o.keys[idx]
	return o.GateKeeper().NewString(frame, k)
}

// Value retrieves the values of the current element in the iteration based on the iterator's current position.
func (o *MapIterator) Value(_ int) IObject {
	idx := o.index - 1
	if idx < 0 || idx >= o.length {
		return o.GateKeeper().UndefinedValue()
	}
	k := o.keys[idx]
	return o.values[k]
}
