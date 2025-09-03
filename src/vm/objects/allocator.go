package objects

import "encoding/gob"

func init() {
	gob.Register(&Allocator{})
}

// Allocator manages memory allocation, references, and execution frame tracking using an IGateKeeper implementation.
type Allocator struct {
	gk         IGateKeeper
	frame      int
	references int
}

// GateKeeper returns the IGateKeeper instance associated with the Allocator.
func (a *Allocator) GateKeeper() IGateKeeper {
	return a.gk
}

// AddRef increments the reference count for the Allocator instance.
func (a *Allocator) AddRef() int {
	a.references++
	return a.references
}

// ReleaseRef decrements the reference count for the Allocator, ensuring it doesn't drop below zero.
func (a *Allocator) ReleaseRef() int {
	if a.references > 0 {
		a.references--
	}
	return a.references
}

// RefCount returns the current reference count for the Allocator instance.
func (a *Allocator) RefCount() int {
	return a.references
}

// Frame returns the current execution frame value managed by the Allocator instance.
func (a *Allocator) Frame() int {
	return a.frame
}

// SetStatic sets the Allocator's frame to a static state by assigning it the FrameStatic constant value.
func (a *Allocator) SetStatic() {
	a.frame = FrameStatic
}
