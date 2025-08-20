package objects

const (
	MapIteratorType  = "map_iterator"
	MapIteratorLabel = "<" + MapIteratorType + ">"
)

// MapIterator is a type used for iterating over key-values pairs in a map-like structure.
// It implements the IIterator interface and provides methods for traversal and element access.
// The embedded Object provides default implementations for methods from the IObject interface.
// The internal state includes the map (values), keys (keys), current position index (index), and total keys length (length).
type MapIterator struct {
	*Object
	values map[string]IObject
	keys   []string
	index  int
	length int
}

// NewMapIterator creates and returns a new instance of MapIterator.
func _newMapIterator(factory *Factory, frame int, v map[string]IObject) *MapIterator {
	var keys []string
	for k := range v {
		keys = append(keys, k)
	}
	return &MapIterator{
		Object: factory.NewObject(frame),
		values: v,
		keys:   keys,
		length: len(keys),
		index:  0,
	}
}

// Copy creates and returns a new instance of MapIterator, duplicating its current state.
func (i *MapIterator) Copy(frame int) IObject {
	ret := i.Factory().NewMapIterator(frame, i.values)
	ret.index = i.index
	return ret
}

// TypeName returns the type name of the MapIterator as a string.
func (i *MapIterator) TypeName() string {
	return MapIteratorType
}

// String returns the string representation of the MapIterator.
func (i *MapIterator) String() string {
	return MapIteratorLabel
}

// Boolean returns true, indicating the MapIterator is considered falsy in a boolean context.
func (i *MapIterator) Boolean() bool {
	return true
}

// Equals determine if the current MapIterator is equal to another IObject. Returns false by default.
func (i *MapIterator) Equals(IObject) bool {
	return false
}

// Next advances the iterator to the next element and returns true if there are more elements to iterate over.
func (i *MapIterator) Next() bool {
	i.index++
	return i.index <= i.length
}

// Key retrieves the key of the current element in the iteration as an IObject.
func (i *MapIterator) Key(frame int) IObject {
	k := i.keys[i.index-1]
	return i.Factory().NewStringNoSize(frame, k)
}

// Value retrieves the values of the current element in the iteration based on the iterator's current position.
func (i *MapIterator) Value(_ int) IObject {
	k := i.keys[i.index-1]
	return i.values[k]
}
