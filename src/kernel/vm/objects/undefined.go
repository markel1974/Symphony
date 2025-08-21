package objects

const (
	UndefinedType  = "undefined"
	UndefinedLabel = "<" + UndefinedType + ">"
)

// Undefined represents an undefined values.
type Undefined struct {
	*Object
}

func newUndefined(factory *GateKeeper, frame int) *Undefined {
	return &Undefined{
		Object: factory.NewObject(frame),
	}
}

// TypeName returns the name of the type.
func (o *Undefined) TypeName() string {
	return UndefinedType
}

func (o *Undefined) String() string {
	return UndefinedLabel
}

// Copy returns a copy of the type.
func (o *Undefined) Copy(frame int, _ int) IObject {
	return o.GateKeeper().NewUndefined(frame)
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
func (o *Undefined) IndexGet(_ int, _ IObject) (IObject, error) {
	return o.GateKeeper().UndefinedValue(), nil
}

// Iterate creates a map iterator.
func (o *Undefined) Iterate(_ int) IIterator {
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
func (o *Undefined) Key(_ int) IObject {
	return o
}

// Value returns the values of the current element.
func (o *Undefined) Value(_ int) IObject {
	return o
}
