package vm

import (
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Frames is a structure that manages function call frames in a virtual machine execution context.
// It maintains a stack of frames and tracks the current frame index for managing execution state.
type Frames struct {
	frames      []*Frame
	framesIndex uint
	errSignal   func(err error)
}

// NewFrames initializes and returns a new Frames instance with the specified main function and maximum frame count.
func NewFrames(maxFrames int, errSignal func(err error)) *Frames {
	f := &Frames{
		frames:      make([]*Frame, maxFrames),
		framesIndex: 1,
		errSignal:   errSignal,
	}
	for i := range f.frames {
		f.frames[i] = NewFunctionCallFrame(errSignal)
	}
	return f
}

// Reset resets the frame index to its initial position, effectively clearing the current frame context.
func (f *Frames) Reset() {
	f.framesIndex = 1
}

// Index returns the current index of the frame stack.
func (f *Frames) Index() uint {
	return f.framesIndex
}

// Head returns the first frame in the frames slice.
func (f *Frames) Head() *Frame {
	return f.frames[0]
}

// Get retrieves the current frame from the `frames` slice based on the `framesIndex`.
func (f *Frames) Get() *Frame {
	return f.frames[f.framesIndex]
}

// GetPrev returns the frame at the position immediately before the current frame based on the framesIndex.
func (f *Frames) GetPrev() *Frame {
	return f.frames[f.framesIndex-1]
}

// Next advances the frame index by one. Returns ErrStackOverflow if the index exceeds the bounds of the frames.
func (f *Frames) Next() {
	f.framesIndex++
	if f.framesIndex >= uint(len(f.frames)) {
		f.errSignal(objects.ErrStackOverflow)
		return
	}
}

// Previous decrements the framesIndex and returns an error if it goes below zero, indicating a stack underflow.
func (f *Frames) Previous() {
	if f.framesIndex <= 1 {
		f.errSignal(objects.ErrStackOverflow)
		return
	}
	f.framesIndex--
}

// Unroll iterates through the frames in reverse order, appending each one to a slice, and decreases the frame index.
func (f *Frames) Unroll() []*Frame {
	var currFrame []*Frame
	for f.framesIndex > 1 {
		f.framesIndex--
		currFrame = append(currFrame, f.frames[f.framesIndex-1])
	}
	return currFrame
}
