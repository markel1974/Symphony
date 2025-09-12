package core

import (
	"fmt"
	"io"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// Stack represents a data structure that operates on a LIFO (Last In, First Out) principle.
// It manages a slice of objects implementing the IObject interface and tracks the stack pointer.
type Stack struct {
	gk        objects.IGateKeeper
	stack     []objects.IObject
	sp        int
	errSignal func(err error)
}

// NewStack creates and initializes a new Stack with the specified size and returns a pointer to it.
func NewStack(gk objects.IGateKeeper, size int, errSignal func(err error)) *Stack {
	s := &Stack{
		gk:        gk,
		sp:        0,
		stack:     make([]objects.IObject, size),
		errSignal: errSignal,
	}
	for i := range s.stack {
		s.stack[i] = gk.UndefinedValue()
	}
	return s
}

// StackPointer returns the current stack pointer (sp) indicating the top position of the stack.
func (v *Stack) StackPointer() int {
	return v.sp
}

// SetStackPointer sets the stack pointer to the specified value.
func (v *Stack) SetStackPointer(sp int) {
	v.sp = sp
}

// Reset resets the stack pointer to zero, effectively clearing the stack.
func (v *Stack) Reset() {
	v.sp = 0
}

// Decrement decreases the stack pointer (sp) by one, effectively moving the pointer to the previous stack position.
func (v *Stack) Decrement() {
	if v.sp == 0 {
		v.errSignal(objects.ErrIndexOutOfBounds)
		return
	}
	v.sp--
}

// DecrementCount reduces the stack pointer (sp) by the specified count, effectively moving it down the stack.
func (v *Stack) DecrementCount(count int) {
	if v.sp < count {
		v.errSignal(objects.ErrIndexOutOfBounds)
		return
	}
	v.sp -= count
}

// SetAbsolute assigns the specified object to the stack at the given absolute index.
func (v *Stack) SetAbsolute(absolute int, obj objects.IObject) {
	if absolute < 0 || absolute >= len(v.stack) {
		v.errSignal(objects.ErrIndexOutOfBounds)
		return
	}
	v.stack[absolute] = obj
}

// SetOffset assigns the given object to the stack at a position determined by the current stack pointer minus the offset.
// If the resolved index is out of bounds, an error signal is triggered.
func (v *Stack) SetOffset(offset int, obj objects.IObject) {
	sp := v.sp - offset
	if sp < 0 || sp >= len(v.stack) {
		v.errSignal(objects.ErrIndexOutOfBounds)
		return
	}
	v.stack[sp] = obj
}

// Set assigns the given object to the position indicated by the current stack pointer minus one.
func (v *Stack) Set(obj objects.IObject) {
	sp := v.sp - 1
	if sp < 0 || sp >= len(v.stack) {
		v.errSignal(objects.ErrIndexOutOfBounds)
		return
	}
	v.stack[sp] = obj
}

// Push adds the provided object to the top of the stack and increments the stack pointer.
func (v *Stack) Push(obj objects.IObject) {
	if v.sp < 0 || v.sp >= len(v.stack) {
		v.errSignal(objects.ErrIndexOutOfBounds)
		return
	}
	v.stack[v.sp] = obj
	v.sp++
}

// PushVarArgs processes a variadic argument list, grouping extra arguments into an array and updating the stack pointer.
func (v *Stack) PushVarArgs(frame int, numArgs int, realArgs int) {
	varArgs := numArgs - realArgs
	if varArgs < 0 {
		return
	}
	numArgs = realArgs + 1
	args := make([]objects.IObject, varArgs)
	spStart := v.sp - varArgs
	for i := spStart; i < v.sp; i++ {
		args[i-spStart] = v.stack[i]
	}
	if spStart < 0 || spStart >= len(v.stack) {
		v.errSignal(objects.ErrIndexOutOfBounds)
		return
	}
	v.stack[spStart] = v.gk.NewArray(frame, args)
	v.sp = spStart + 1
}

// Pop removes and returns the object at the top of the stack.
// It returns UndefinedValue if the stack is empty (stack underflow).
func (v *Stack) Pop() objects.IObject {
	if v.sp == 0 {
		v.errSignal(objects.ErrIndexOutOfBounds)
		return v.gk.UndefinedValue()
	}
	v.sp--
	return v.stack[v.sp]
}

// PopArray removes and returns a specified number of elements from the stack as a slice of IObject.
func (v *Stack) PopArray(numElements int) []objects.IObject {
	if numElements > v.sp {
		v.errSignal(objects.ErrIndexOutOfBounds)
		return []objects.IObject{}
	}
	var elements []objects.IObject
	for i := v.sp - numElements; i < v.sp; i++ {
		elements = append(elements, v.stack[i])
	}
	v.sp -= numElements
	return elements
}

// PopMap removes the specified number of key-value pairs from the stack and returns them as a map.
func (v *Stack) PopMap(numElements int) map[string]objects.IObject {
	kv := make(map[string]objects.IObject, numElements)
	if numElements > v.sp {
		v.errSignal(objects.ErrIndexOutOfBounds)
		return kv
	}
	for i := v.sp - numElements; i < v.sp; i += 2 {
		k := v.stack[i]
		value := v.stack[i+1]
		key, ok := k.(*objects.String)
		if !ok {
			v.errSignal(fmt.Errorf("expected key to be of type AsString, got %T (value is %v)", k, value))
			return nil
		}
		kv[key.Value()] = value
	}
	v.sp -= numElements
	return kv
}

// PeekAbsolute retrieves the object at the specified absolute index in the stack without modifying the stack pointer.
func (v *Stack) PeekAbsolute(absolute int) objects.IObject {
	sp := absolute
	if sp < 0 || sp >= len(v.stack) {
		v.errSignal(objects.ErrIndexOutOfBounds)
		return v.gk.UndefinedValue()
	}
	return v.stack[absolute]
}

// PeekOffset retrieves the object at the stack pointer minus the specified offset without modifying the stack pointer.
// Returns UndefinedValue if the resolved index is out of stack bounds.
func (v *Stack) PeekOffset(offset int) objects.IObject {
	sp := v.sp - offset
	if sp < 0 || sp >= len(v.stack) {
		v.errSignal(objects.ErrIndexOutOfBounds)
		return v.gk.UndefinedValue()
	}
	ret := v.stack[sp]
	return ret
}

// Peek returns the object at the top of the stack.
func (v *Stack) Peek() objects.IObject {
	sp := v.sp - 1
	if sp < 0 || sp >= len(v.stack) {
		v.errSignal(objects.ErrIndexOutOfBounds)
		return v.gk.UndefinedValue()
	}
	ret := v.stack[sp]
	return ret
}

// PeekInterval returns a slice of objects from the stack within the specified range [start:end). Returns nil for invalid ranges.
func (v *Stack) PeekInterval(start int, end int) []objects.IObject {
	if start == end {
		v.errSignal(objects.ErrIndexOutOfBounds)
		return []objects.IObject{}
	}
	if start < 0 || end > len(v.stack) || start > end {
		v.errSignal(objects.ErrIndexOutOfBounds)
		return []objects.IObject{}
	}
	return v.stack[start:end]
}

// PeekArray retrieves a slice of IObject elements from the stack, based on the specified number of arguments.
func (v *Stack) PeekArray(numArgs int) []objects.IObject {
	start := v.sp - numArgs
	if start < 0 || start >= len(v.stack) {
		v.errSignal(objects.ErrIndexOutOfBounds)
		return nil
	}
	end := v.sp
	if end < 0 || end > len(v.stack) {
		return nil
	}
	z := v.stack[start:v.sp]
	return z
}

// CopyOffset copies a specified number of arguments from the source stack area to the designated frame in the stack.
// It performs bounds checks to prevent stack underflow or overflow and signals an error if any condition is violated.
func (v *Stack) CopyOffset(start int, count int) {
	if count == 0 {
		return
	}
	sourceStart := v.sp - count
	destinationEnd := start + count
	if sourceStart < 0 {
		v.errSignal(fmt.Errorf("stack underflow during argument copy: sourceStart is %d", sourceStart))
		return
	}
	if destinationEnd > len(v.stack) {
		v.errSignal(fmt.Errorf("stack overflow during argument copy: destinationEnd is %d", destinationEnd))
		return
	}
	copy(v.stack[start:destinationEnd], v.stack[sourceStart:v.sp])
}

// ReleaseAll releases all managed resources through the gatekeeper (gk) and ensures proper cleanup.
func (v *Stack) ReleaseAll() {
	v.gk.ReleaseAll()
}

// ReleaseObjects releases objects from the stack within the range [start:end), excluding those in the preserve list.
// Objects are managed by a gatekeeper (gk) for proper disposal. It handles optimized cases for single-item preservation.
// The function ensures safe operations with bounds checks and no operation for invalid ranges.
func (v *Stack) ReleaseObjects(frame int, start int, end int, preserve []objects.IObject) {
	if start >= end {
		return
	}
	if start < 0 || end > len(v.stack) {
		return
	}
	switch len(preserve) {
	case 0:
		v.gk.ReleaseObjects(frame, v.stack[start:end])
		return
	case 1:
		// Optimized case: single value, avoiding map usage.
		preserveObj := preserve[0]
		batchList := make([]objects.IObject, 0, end-start)
		for i := start; i < end; i++ {
			if obj := v.stack[i]; obj != preserveObj {
				batchList = append(batchList, obj)
			}
		}
		if len(batchList) > 0 {
			v.gk.ReleaseObjects(frame, batchList)
		}
		return
	default:
		preserveHelper := make(map[objects.IObject]bool, len(preserve))
		for _, obj := range preserve {
			preserveHelper[obj] = true
		}
		batchList := make([]objects.IObject, 0, end-start)
		for i := start; i < end; i++ {
			if obj := v.stack[i]; !preserveHelper[obj] {
				batchList = append(batchList, obj)
			}
		}
		if len(batchList) > 0 {
			v.gk.ReleaseObjects(frame, batchList)
		}
		return
	}
}

// Print outputs each element in the stack from the bottom to the current stack pointer.
func (v *Stack) Print(writer io.Writer) {
	_, _ = fmt.Fprintln(writer, "------------- Stack -------------")
	for x := 0; x < v.sp; x++ {
		_, _ = fmt.Fprintf(writer, "%v\n", v.stack[x])
	}
	_, _ = fmt.Fprintln(writer, "--------------------------------")
}
