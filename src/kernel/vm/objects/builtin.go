package objects

import (
	"fmt"
)

// BuiltinType is a constant string identifying the type name of builtin objects.
const (
	BuiltinType = "builtin"
)

// Builtin defines a struct representing a built-in object with a name and an integer value.
type Builtin struct {
	*Object
	name  string
	index int
}

// NewBuiltin creates a new instance of Builtin with the specified name and value.
func newBuiltin(factory *GateKeeper, frame int, name string, index int) *Builtin {
	return &Builtin{
		Object: factory.NewObject(frame),
		name:   name,
		index:  index,
	}
}

// Name returns the name of the Builtin object.
func (o *Builtin) Name() string {
	return o.name
}

// Value returns the integer value associated with the Builtin object.
func (o *Builtin) Value() int {
	return o.index
}

// String returns a string representation of the Builtin object, formatted as "name: value".
func (o *Builtin) String() string {
	return fmt.Sprintf("%s: %d", o.name, o.index)
}

// TypeName returns the type name of the Builtin object as a string.
func (o *Builtin) TypeName() string {
	return BuiltinType
}

// Copy creates and returns a new Builtin object with the same name and value as the current instance.
func (o *Builtin) Copy(frame int) IObject {
	return o.GateKeeper().NewBuiltin(frame, o.name, o.index)
}

// Boolean returns false, indicating the boolean representation of the Builtin object is always false.
func (o *Builtin) Boolean() bool {
	return false
}

// Equals checks if the current Builtin object is equal to another IObject based on its name and value.
func (o *Builtin) Equals(x IObject) bool {
	t, ok := x.(*Builtin)
	if !ok {
		return false
	}
	return o.name == t.name && o.index == t.index
}
