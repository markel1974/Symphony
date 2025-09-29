package handler

import (
	"fmt"
	"io"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// Stack represents a data structure that operates on a LIFO (Last In, First Out) principle.
// It manages a slice of objects implementing the IObject interface and tracks the stack pointer.
type Stack struct {
	gk             objects.IGateKeeper
	stack          []objects.IObject
	sp             uint
	shutdownSignal func(err error)
}

// NewStack creates and initializes a new Stack with the specified size and returns a pointer to it.
func NewStack(gk objects.IGateKeeper, size int, shutdownSignal func(err error)) *Stack {
	s := &Stack{
		gk:             gk,
		sp:             0,
		stack:          make([]objects.IObject, size),
		shutdownSignal: shutdownSignal,
	}
	for i := range s.stack {
		s.stack[i] = gk.UndefinedValue()
	}
	return s
}

// StackPointer returns the current stack pointer (sp) indicating the top position of the stack.
func (v *Stack) StackPointer() uint {
	return v.sp
}

// SetStackPointer sets the stack pointer to the specified value.
func (v *Stack) SetStackPointer(sp uint) {
	if sp >= uint(len(v.stack)) {
		v.shutdownSignal(objects.ErrIndexOutOfBounds)
		return
	}
	v.sp = sp
}

// Reset resets the stack pointer to zero, effectively clearing the stack.
func (v *Stack) Reset() {
	v.sp = 0
}

// Decrement decreases the stack pointer (sp) by one, effectively moving the pointer to the previous stack position.
func (v *Stack) Decrement() {
	if v.sp == 0 {
		v.shutdownSignal(objects.ErrIndexOutOfBounds)
		return
	}
	v.sp--
}

// DecrementCount reduces the stack pointer (sp) by the specified count, effectively moving it down the stack.
func (v *Stack) DecrementCount(count uint) {
	if count > v.sp {
		v.shutdownSignal(objects.ErrIndexOutOfBounds)
		return
	}
	v.sp -= count
}

// Set assigns the given object to the position indicated by the current stack pointer minus one.
func (v *Stack) Set(obj objects.IObject) {
	if v.sp == 0 {
		v.shutdownSignal(objects.ErrIndexOutOfBounds)
		return
	}
	if obj == nil {
		obj = v.gk.UndefinedValue()
	}
	v.stack[v.sp-1] = obj
}

// Push adds the provided object to the top of the stack and increments the stack pointer.
func (v *Stack) Push(obj objects.IObject) {
	if v.sp+1 >= uint(len(v.stack)) {
		v.shutdownSignal(objects.ErrIndexOutOfBounds)
		return
	}
	if obj == nil {
		obj = v.gk.UndefinedValue()
	}
	v.stack[v.sp] = obj
	v.sp++
}

// Pop removes and returns the object at the top of the stack.
// It returns UndefinedValue if the stack is empty (stack underflow).
func (v *Stack) Pop() objects.IObject {
	if v.sp == 0 {
		v.shutdownSignal(objects.ErrIndexOutOfBounds)
		return v.gk.UndefinedValue()
	}
	v.sp--
	return v.stack[v.sp]
}

// PopArray removes and returns a specified number of elements from the stack as a slice of IObject.
func (v *Stack) PopArray(numElements uint) []objects.IObject {
	if numElements > v.sp {
		v.shutdownSignal(objects.ErrIndexOutOfBounds)
		return []objects.IObject{}
	}
	elements := make([]objects.IObject, numElements)
	target := v.sp - numElements
	for i := target; i < v.sp; i++ {
		elements[i-target] = v.stack[i]
	}
	v.sp -= numElements
	return elements
}

// PopStruct removes multiple elements from the stack and returns a string representation and a map of objects.
func (v *Stack) PopStruct(numElements uint) (string, map[string]objects.IObject) {
	nameObj := v.Pop()
	s := v.PopMap(numElements)
	return nameObj.AsString(), s
}

// PopMap removes the specified number of key-value pairs from the stack and returns them as a map.
func (v *Stack) PopMap(numElements uint) map[string]objects.IObject {
	if numElements&1 == 1 {
		v.shutdownSignal(objects.ErrIndexOutOfBounds)
		return nil
	}
	kv := make(map[string]objects.IObject, numElements)
	if numElements > v.sp {
		v.shutdownSignal(objects.ErrIndexOutOfBounds)
		return kv
	}
	target := v.sp - numElements
	for i := target; i < v.sp; i += 2 {
		k := v.stack[i]
		value := v.stack[i+1]
		kv[k.AsString()] = value
	}
	v.sp -= numElements
	return kv
}

// PopInterface pops a concrete value and a number of method-function pairs from the stack, returning them as an interface table.
// It constructs a map where method names are keys and their corresponding functions are values, along with the concrete value.
// The parameter numMethods specifies the number of method-function pairs to pop from the stack.
func (v *Stack) PopInterface(numMethods int) (objects.IObject, map[string]objects.IObject) {
	iTable := make(map[string]objects.IObject, numMethods)
	for i := 0; i < numMethods; i++ {
		methodFunc := v.Pop()
		methodNameObj := v.Pop()
		iTable[methodNameObj.AsString()] = methodFunc
	}
	concreteValue := v.Pop()
	return concreteValue, iTable
}

// SetAbsolute assigns the specified object to the stack at the given absolute index.
func (v *Stack) SetAbsolute(absolute uint, obj objects.IObject) {
	if absolute >= v.sp {
		v.shutdownSignal(objects.ErrIndexOutOfBounds)
		return
	}
	if obj == nil {
		obj = v.gk.UndefinedValue()
	}
	v.stack[absolute] = obj
}

// PeekAbsolute retrieves the object at the specified absolute index in the stack without modifying the stack pointer.
func (v *Stack) PeekAbsolute(absolute uint) objects.IObject {
	if absolute >= v.sp {
		v.shutdownSignal(objects.ErrIndexOutOfBounds)
		return v.gk.UndefinedValue()
	}
	return v.stack[absolute]
}

// PeekOffset retrieves the object at the stack pointer minus the specified offset without modifying the stack pointer.
// Returns UndefinedValue if the resolved index is out of stack bounds.
func (v *Stack) PeekOffset(offset uint) objects.IObject {
	if offset == 0 || offset > v.sp {
		v.shutdownSignal(objects.ErrIndexOutOfBounds)
		return v.gk.UndefinedValue()
	}
	ret := v.stack[v.sp-offset]
	return ret
}

// SetOffset assigns the given object to the stack at a position determined by the current stack pointer minus the offset.
// If the resolved index is out of bounds, an error signal is triggered.
func (v *Stack) SetOffset(offset uint, obj objects.IObject) {
	if offset == 0 || offset > v.sp {
		v.shutdownSignal(objects.ErrIndexOutOfBounds)
		return
	}
	v.stack[v.sp-offset] = obj
}

// Peek returns the object at the top of the stack.
func (v *Stack) Peek() objects.IObject {
	if v.sp == 0 {
		v.shutdownSignal(objects.ErrIndexOutOfBounds)
		return v.gk.UndefinedValue()
	}
	ret := v.stack[v.sp-1]
	return ret
}

// PeekInterval returns a slice of objects from the stack within the specified range [start:end). Returns nil for invalid ranges.
func (v *Stack) PeekInterval(start uint, end uint) []objects.IObject {
	if start == end {
		v.shutdownSignal(objects.ErrIndexOutOfBounds)
		return []objects.IObject{}
	}
	if end > v.sp || start > end {
		v.shutdownSignal(objects.ErrIndexOutOfBounds)
		return []objects.IObject{}
	}
	return v.stack[start:end]
}

// PeekArray retrieves a slice of IObject elements from the stack, based on the specified number of arguments.
func (v *Stack) PeekArray(numArgs uint) []objects.IObject {
	if numArgs > v.sp {
		v.shutdownSignal(objects.ErrIndexOutOfBounds)
		return nil
	}
	start := v.sp - numArgs
	z := v.stack[start:v.sp]
	return z
}

// CopyOffset copies a specified number of arguments from the source stack area to the designated frame in the stack.
// It performs bounds checks to prevent stack underflow or overflow and signals an error if any condition is violated.
func (v *Stack) CopyOffset(start uint, count uint) {
	if count == 0 {
		return
	}
	if count > v.sp || start > v.sp-count {
		v.shutdownSignal(objects.ErrIndexOutOfBounds)
		return
	}
	destinationEnd := start + count
	sourceStart := v.sp - count
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
	for x := uint(0); x < v.sp; x++ {
		_, _ = fmt.Fprintf(writer, "%v\n", v.stack[x])
	}
	_, _ = fmt.Fprintln(writer, "--------------------------------")
}
