package handler

import (
	"github.com/markel1974/c64emu/src/vm/objects"
)

const (
	RootFrame = 1
)

// Frames is a structure that manages function call frames in a virtual machine execution context.
// It maintains a stack of frames and tracks the current frame index for managing execution state.
type Frames struct {
	gk             objects.IGateKeeper
	frames         []*Frame
	framesIndex    uint
	frameMax       uint
	shutdownSignal func(err error)
}

// NewFrames initializes and returns a new Frames instance with the specified main function and maximum frame count.
func NewFrames(gk objects.IGateKeeper, maxFrames int, startInterval uint, shutdownSignal func(err error)) *Frames {
	f := &Frames{
		gk:             gk,
		frames:         make([]*Frame, maxFrames),
		framesIndex:    RootFrame,
		shutdownSignal: shutdownSignal,
		frameMax:       0,
	}
	for i := range f.frames {
		frameId := startInterval + uint(i)
		f.frames[i] = NewFrame(gk, frameId, shutdownSignal)
	}
	return f
}

// Reset resets the frame index to its initial position, effectively clearing the current frame context.
func (f *Frames) Reset() {
	f.framesIndex = RootFrame
	f.frameMax = 0
}

// Max returns the maximum frame index accessed during execution. It is updated as new frames are added to the stack.
func (f *Frames) Max() uint {
	return f.frameMax
}

// Head returns the first frame in the frames slice.
func (f *Frames) Head() *Frame {
	return f.frames[0]
}

// Current retrieves the current frame from the `frames` slice based on the `framesIndex`.
func (f *Frames) Current() *Frame {
	return f.frames[f.framesIndex]
}

// Previous returns the frame at the position immediately before the current frame based on the framesIndex.
func (f *Frames) Previous() *Frame {
	return f.frames[f.framesIndex-1]
}

// MoveNext advances the frame index by one. Returns ErrStackOverflow if the index exceeds the bounds of the frames.
func (f *Frames) MoveNext() {
	f.framesIndex++
	if f.framesIndex >= uint(len(f.frames)) {
		f.shutdownSignal(objects.ErrStackOverflow)
		return
	}
	if f.framesIndex > f.frameMax {
		f.frameMax = f.framesIndex
	}
}

// MovePrevious decrements the framesIndex and returns an error if it goes below zero, indicating a stack underflow.
func (f *Frames) MovePrevious() {
	if f.framesIndex <= RootFrame {
		f.shutdownSignal(objects.ErrStackOverflow)
		return
	}
	f.framesIndex--
}

// CanMovePrevious checks if the frame index can be decremented without going below the root frame index.
func (f *Frames) CanMovePrevious() bool {
	return f.framesIndex > RootFrame
}

// Unroll iterates through the frames in reverse order, appending each one to a slice, and decreases the frame index.
func (f *Frames) Unroll() []*Frame {
	var currFrame []*Frame
	for f.framesIndex > RootFrame {
		f.framesIndex--
		currFrame = append(currFrame, f.frames[f.framesIndex-1])
	}
	return currFrame
}
