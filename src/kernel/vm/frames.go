package vm

import (
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Frames manages a stack of Frame objects and keeps track of the current frame index.
type Frames struct {
	frames      []*Frame
	framesIndex int
}

// NewFrames initializes and returns a new Frames instance for managing the call stack of a compiled function.
func NewFrames(main *objects.FunctionCompiled, maxFrames int) *Frames {
	f := &Frames{
		frames:      make([]*Frame, maxFrames),
		framesIndex: 1,
	}
	for i := range f.frames {
		f.frames[i] = NewFunctionCallFrame()
	}
	f.frames[0].SetCompiledFunction(main)
	f.frames[0].SetIP(-1)
	return f
}

// Clear resets the framesIndex to 1, effectively clearing the frames state for reuse.
func (f *Frames) Clear() {
	f.framesIndex = 1
}

// Index returns the current index of the frames stack represented by the framesIndex field.
func (f *Frames) Index() int {
	return f.framesIndex
}

// Head returns the first Frame from the frames slice in the Frames struct.
func (f *Frames) Head() *Frame {
	return f.frames[0]
}

// Get retrieves the current frame at the position indicated by the framesIndex in the Frames collection.
func (f *Frames) Get() *Frame {
	return f.frames[f.framesIndex]
}

func (f *Frames) GetPrev() *Frame {
	return f.frames[f.framesIndex-1]
}

// Next advances the frame index to the next frame in the stack, returning an error if the end of the stack is reached.
func (f *Frames) Next() error {
	f.framesIndex++
	if f.framesIndex >= len(f.frames) {
		return objects.ErrStackOverflow
	}
	return nil
}

func (f *Frames) Previous() error {
	f.framesIndex--
	if f.framesIndex < 0 {
		return objects.ErrStackOverflow
	}
	return nil
}

// Unroll returns the current frame and all previous frames in the stack.
func (f *Frames) Unroll() []*Frame {
	var currFrame []*Frame
	for f.framesIndex > 1 {
		f.framesIndex--
		currFrame = append(currFrame, f.frames[f.framesIndex-1])
	}
	return currFrame
}
