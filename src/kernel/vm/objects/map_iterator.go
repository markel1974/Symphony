package objects

// MapIterator is a type used for iterating over key-values pairs in a map-like structure.
// It implements the IIterator interface and provides methods for traversal and element access.
// The embedded ObjectImpl provides default implementations for methods from the IObject interface.
// The internal state includes the map (v), keys (k), current position index (i), and total keys length (l).
type MapIterator struct {
	ObjectImpl
	v map[string]IObject
	k []string
	i int
	l int
}

func NewMapIterator(v map[string]IObject, keys []string) *MapIterator {
	return &MapIterator{v: v, k: keys, l: len(keys)}
}

// TypeName returns the type name of the MapIterator as a string.
func (i *MapIterator) TypeName() string {
	return "map-iterator"
}

// String returns the string representation of the MapIterator.
func (i *MapIterator) String() string {
	return "<map-iterator>"
}

// Falsy returns true, indicating the MapIterator is considered falsy in a boolean context.
func (i *MapIterator) Falsy() bool {
	return true
}

// Equals determines if the current MapIterator is equal to another IObject. Returns false by default.
func (i *MapIterator) Equals(IObject) bool {
	return false
}

// Copy creates and returns a new instance of MapIterator, duplicating its current state.
func (i *MapIterator) Copy() IObject {
	return &MapIterator{v: i.v, k: i.k, i: i.i, l: i.l}
}

// Next advances the iterator to the next element and returns true if there are more elements to iterate over.
func (i *MapIterator) Next() bool {
	i.i++
	return i.i <= i.l
}

// Key retrieves the key of the current element in the iteration as an IObject.
func (i *MapIterator) Key() IObject {
	k := i.k[i.i-1]
	return NewStringNoSize(k)
}

// Value retrieves the values of the current element in the iteration based on the iterator's current position.
func (i *MapIterator) Value() IObject {
	k := i.k[i.i-1]
	return i.v[k]
}

// ToMap converts an IObject to a map[string]interface{} if the object is a *Map, recursively applying ToInterface.
func ToMap(o IObject) (res map[string]interface{}) {
	switch o := o.(type) {
	case *Map:
		res = make(map[string]interface{})
		for key, v := range o.values {
			res[key] = ToInterface(v)
		}
	}
	return
}
