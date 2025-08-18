package objects

// IIterator defines an interface for iterating over a collection of elements.
// Next determines if more elements are available in the iteration.
// Key retrieves the key or index of the current element in the iteration.
// Value returns the values of the current element in the iteration.
type IIterator interface {
	IObject

	Next() bool

	Key() IObject

	Value() IObject
}
