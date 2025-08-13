package vm

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Stack represents a data structure that operates on a LIFO (Last In, First Out) principle.
// It manages a slice of objects implementing the IObject interface and tracks the stack pointer.
type Stack struct {
	stack []objects.IObject
	sp    int
}

// NewStack creates and initializes a new Stack with the specified size and returns a pointer to it.
func NewStack(size int) *Stack {
	return &Stack{
		sp:    0,
		stack: make([]objects.IObject, size),
	}
}

// StackPointer returns the current stack pointer (sp) indicating the top position of the stack.
func (v *Stack) StackPointer() int {
	return v.sp
}

// SetStackPointer sets the stack pointer to the specified value.
func (v *Stack) SetStackPointer(sp int) {
	v.sp = sp
}

// Clear resets the stack pointer to zero, effectively clearing the stack.
func (v *Stack) Clear() {
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
	v.stack[absolute] = obj
}

// Set assigns the given object to the position indicated by the current stack pointer minus one.
func (v *Stack) Set(obj objects.IObject) {
	v.stack[v.sp-1] = obj
}

// Push adds the provided object to the top of the stack and increments the stack pointer.
func (v *Stack) Push(obj objects.IObject) {
	v.stack[v.sp] = obj
	v.sp++
}

// PushVarArgs processes a variadic argument list, grouping extra arguments into an array and updating the stack pointer.
func (v *Stack) PushVarArgs(numArgs int, realArgs int) {
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
	v.stack[spStart] = objects.NewArray(args)
	v.sp = spStart + 1
}

// Pop removes and returns the object at the top of the stack.
// It returns UndefinedValue if the stack is empty (stack underflow).
func (v *Stack) Pop() objects.IObject {
	if v.sp == 0 {
		return objects.UndefinedValue
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
		key := v.stack[i]
		value := v.stack[i+1]
		kv[key.(*objects.String).Value()] = value
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
	ret := v.stack[sp]
	return ret
}

// Peek returns the object at the top of the stack.
func (v *Stack) Peek() objects.IObject {
	if v.sp == 0 {
		return objects.UndefinedValue
	}
	ret := v.stack[v.sp-1]
	return ret
}

// PeekArrayObject retrieves a slice of IObject elements from the stack, based on the specified number of arguments.
func (v *Stack) PeekArrayObject(numArgs int) []objects.IObject {
	z := v.stack[v.sp-numArgs : v.sp]
	return z
}

// Print outputs each element in the stack from the bottom to the current stack pointer.
func (v *Stack) Print() {
	for x := 0; x < v.sp; x++ {
		fmt.Printf("%v\n", v.stack[x])
	}
}
