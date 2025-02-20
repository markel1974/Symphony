package executor

import "github.com/go-gl/gl/v3.3-core/gl"

// Binder is a utility type that manages OpenGL object bindings and restoration to previous states.
type Binder struct {
	restoreLoc uint32
	bindFunc   func(uint32)
	obj        uint32
	prev       []uint32
}

// Bind binds the current object if it is not already bound and pushes the previously bound object onto a stack.
func (b *Binder) Bind() *Binder {
	var prev int32
	gl.GetIntegerv(b.restoreLoc, &prev)
	b.prev = append(b.prev, uint32(prev))
	if b.prev[len(b.prev)-1] != b.obj {
		b.bindFunc(b.obj)
	}
	return b
}

// Restore reverts the binder to the previously bound object and removes the last state from the stack. Returns the binder.
func (b *Binder) Restore() *Binder {
	if b.prev[len(b.prev)-1] != b.obj {
		b.bindFunc(b.prev[len(b.prev)-1])
	}
	b.prev = b.prev[:len(b.prev)-1]
	return b
}
