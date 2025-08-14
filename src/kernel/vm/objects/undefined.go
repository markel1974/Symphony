package objects

const (
	UndefinedType  = "undefined"
	UndefinedLabel = "<" + UndefinedType + ">"
)

var (
	// UndefinedValue represents an undefined values.
	UndefinedValue IObject = &Undefined{}
)

// Undefined represents an undefined values.
type Undefined struct {
	Object
}

// TypeName returns the name of the type.
func (o *Undefined) TypeName() string {
	return UndefinedType
}

func (o *Undefined) String() string {
	return UndefinedLabel
}

// Copy returns a copy of the type.
func (o *Undefined) Copy() IObject {
	return o
}

// Boolean returns true.
func (o *Undefined) Boolean() bool {
	return true
}

// Equals returns true if the values of the type are equal to the values of
// another object.
func (o *Undefined) Equals(x IObject) bool {
	return o == x
}

// IndexGet returns an element at a given index.
func (o *Undefined) IndexGet(_ IObject) (IObject, error) {
	return UndefinedValue, nil
}

// Iterate creates a map iterator.
func (o *Undefined) Iterate() IIterator {
	return o
}

// CanIterate returns whether the IObject can be Iterated.
func (o *Undefined) CanIterate() bool {
	return true
}

// Next returns true if there are more elements to iterate.
func (o *Undefined) Next() bool {
	return false
}

// Key returns the key or index values of the current element.
func (o *Undefined) Key() IObject {
	return o
}

// Value returns the values of the current element.
func (o *Undefined) Value() IObject {
	return o
}
