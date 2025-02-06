package executor

import "github.com/go-gl/gl/v3.3-core/gl"

type Binder struct {
	restoreLoc uint32
	bindFunc   func(uint32)
	obj        uint32
	prev       []uint32
}

func (b *Binder) Bind() *Binder {
	var prev int32
	gl.GetIntegerv(b.restoreLoc, &prev)
	b.prev = append(b.prev, uint32(prev))
	if b.prev[len(b.prev)-1] != b.obj {
		b.bindFunc(b.obj)
	}
	return b
}

func (b *Binder) Restore() *Binder {
	if b.prev[len(b.prev)-1] != b.obj {
		b.bindFunc(b.prev[len(b.prev)-1])
	}
	b.prev = b.prev[:len(b.prev)-1]
	return b
}
