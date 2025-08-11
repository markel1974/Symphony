package objects

// IIterator defines an interface for iterating over a collection of elements.
// Next determines if more elements are available in the iteration.
// Key retrieves the key or index of the current element in the iteration.
// Value returns the values of the current element in the iteration.
type IIterator interface {
	IObject

	// Next returns true if there are more elements to iterate.
	Next() bool

	// Key returns the key or index values of the current element.
	Key() IObject

	// Value returns the values of the current element.
	Value() IObject
}
