package vm

import (
	"fmt"
	"io"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Stack represents a data structure that operates on a LIFO (Last In, First Out) principle.
// It manages a slice of objects implementing the IObject interface and tracks the stack pointer.
type Stack struct {
	factory        *objects.Factory
	stack          []objects.IObject
	sp             int
	allocations    int64
	maxAllocations int64
	errSignal      func(err error)
}

// NewStack creates and initializes a new Stack with the specified size and returns a pointer to it.
func NewStack(factory *objects.Factory, size int, maxAllocations int64, errSignal func(err error)) *Stack {
	s := &Stack{
		factory:        factory,
		sp:             0,
		stack:          make([]objects.IObject, size),
		maxAllocations: maxAllocations,
		allocations:    maxAllocations + 1,
		errSignal:      errSignal,
	}
	for i := range s.stack {
		s.stack[i] = factory.UndefinedValue()
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
	v.allocations = v.maxAllocations + 1
	v.sp = 0
}

// Decrement decreases the stack pointer (sp) by one, effectively moving the pointer to the previous stack position.
func (v *Stack) Decrement() {
	v.sp--
}

func (v *Stack) DecrementCount(count int) {
	v.sp -= count
}

// SetAbsolute assigns the specified object to the stack at the given absolute index.
func (v *Stack) SetAbsolute(absolute int, obj objects.IObject) {
	if v.allocations--; v.allocations == 0 {
		v.errSignal(objects.ErrObjectAllocLimit)
		return
	}
	if absolute < 0 || absolute >= len(v.stack) {
		v.errSignal(objects.ErrIndexOutOfBounds)
		return
	}
	v.stack[absolute] = obj
}

// Set assigns the given object to the position indicated by the current stack pointer minus one.
func (v *Stack) Set(obj objects.IObject) {
	if v.allocations--; v.allocations == 0 {
		v.errSignal(objects.ErrObjectAllocLimit)
		return
	}
	sp := v.sp - 1
	if sp < 0 || sp >= len(v.stack) {
		v.errSignal(objects.ErrIndexOutOfBounds)
		return
	}
	v.stack[sp] = obj
}

// Push adds the provided object to the top of the stack and increments the stack pointer.
func (v *Stack) Push(obj objects.IObject) {
	if v.allocations--; v.allocations == 0 {
		v.errSignal(objects.ErrObjectAllocLimit)
		return
	}
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
	if v.allocations--; v.allocations <= 0 {
		v.errSignal(objects.ErrObjectAllocLimit)
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
	v.stack[spStart] = v.factory.NewArray(frame, args)
	v.sp = spStart + 1
}

// Pop removes and returns the object at the top of the stack.
// It returns UndefinedValue if the stack is empty (stack underflow).
func (v *Stack) Pop() objects.IObject {
	if v.sp == 0 {
		return v.factory.UndefinedValue()
	}
	v.sp--
	return v.stack[v.sp]
}

// PopArrayElements removes and returns a specified number of elements from the stack as a slice of IObject.
func (v *Stack) PopArrayElements(numElements int) []objects.IObject {
	var elements []objects.IObject
	for i := v.sp - numElements; i < v.sp; i++ {
		elements = append(elements, v.stack[i])
	}
	v.sp -= numElements
	return elements
}

// PopMapElements removes the specified number of key-value pairs from the stack and returns them as a map.
func (v *Stack) PopMapElements(numElements int) map[string]objects.IObject {
	kv := make(map[string]objects.IObject, numElements)
	for i := v.sp - numElements; i < v.sp; i += 2 {
		k := v.stack[i]
		value := v.stack[i+1]
		key, ok := k.(*objects.String)
		if !ok {
			v.errSignal(fmt.Errorf("expected key to be of type String, got %T (value is %v)", k, value))
			return nil
		}
		kv[key.Value()] = value
	}
	v.sp -= numElements
	return kv
}

// PeekAbsolute retrieves the object at the specified absolute index in the stack without modifying the stack pointer.
func (v *Stack) PeekAbsolute(absolute int) objects.IObject {
	return v.stack[absolute]
}

// PeekOffset returns the stack object at the specified offset relative to the current stack pointer.
func (v *Stack) PeekOffset(offset int) objects.IObject {
	sp := v.sp + offset
	if sp < 0 || sp >= len(v.stack) {
		return v.factory.UndefinedValue()
	}
	ret := v.stack[sp]
	return ret
}

// Peek returns the object at the top of the stack.
func (v *Stack) Peek() objects.IObject {
	sp := v.sp - 1
	if sp < 0 || sp >= len(v.stack) {
		return v.factory.UndefinedValue()
	}
	ret := v.stack[v.sp-1]
	return ret
}

// PeekArrayObject retrieves a slice of IObject elements from the stack, based on the specified number of arguments.
func (v *Stack) PeekArrayObject(numArgs int) []objects.IObject {
	start := v.sp - numArgs
	if start < 0 || start >= len(v.stack) {
		return nil
	}
	end := v.sp
	if end < 0 || end > len(v.stack) {
		return nil
	}
	z := v.stack[start:v.sp]
	return z
}

func (v *Stack) ReleaseObjects(start, end int) {
	//TODO IMPLMENT!
	/*
		for i := start; i < end; i++ {
			obj := v.stack[i]
			switch o := obj.(type) {
			case *objects.Int:
				v.factory.ReleaseInt(o)
			case *objects.Float:
				v.factory.ReleaseFloat(o)
			case *objects.String:
				v.factory.ReleaseString(o)
			case *objects.Array:
				v.factory.ReleaseArray(o)
				// Oggetti come Bool e Undefined sono singleton, non vanno rilasciati.
				// Altri oggetti complessi potrebbero non avere un pool, quindi non fare nulla.
			}
			v.stack[i] = v.factory.UndefinedValue()
		}

	*/
}

// Print outputs each element in the stack from the bottom to the current stack pointer.
func (v *Stack) Print(writer io.Writer) {
	_, _ = fmt.Fprintln(writer, "------------- Stack -------------")
	for x := 0; x < v.sp; x++ {
		_, _ = fmt.Fprintf(writer, "%v\n", v.stack[x])
	}
	_, _ = fmt.Fprintln(writer, "--------------------------------")
}
