package objects

import "github.com/markel1974/c64emu/src/kernel/compiler"

// CompiledFunction represents a compiled function with bytecode instructions, metadata, and associated free variables.
type CompiledFunction struct {
	ObjectImpl
	instructions  []byte
	numLocals     int
	numParameters int
	varArgs       bool
	sourceMap     map[int]compiler.Pos
	free          []*ObjectPtr
}

// NewCompiledFunction creates a new instance of CompiledFunction with the given instructions, locals, parameters, varArgs, sourceMap, and free vars.
func NewCompiledFunction(instructions []byte, numLocals int, numParameters int, varArgs bool, sourceMap map[int]compiler.Pos, free []*ObjectPtr) *CompiledFunction {
	if sourceMap == nil {
		sourceMap = make(map[int]compiler.Pos)
	}
	return &CompiledFunction{
		instructions:  instructions,
		numLocals:     numLocals,
		numParameters: numParameters,
		varArgs:       varArgs,
		sourceMap:     sourceMap,
		free:          free,
	}
}

// Instructions returns the bytecode instructions of the compiled function.
func (o *CompiledFunction) Instructions() []byte {
	return o.instructions
}

// NumLocals returns the number of local variables required by the compiled function.
func (o *CompiledFunction) NumLocals() int {
	return o.numLocals
}

// NumParameters returns the total number of parameters required by the compiled function.
func (o *CompiledFunction) NumParameters() int {
	return o.numParameters
}

// VarArgs returns true if the function accepts a variable number of arguments, otherwise false.
func (o *CompiledFunction) VarArgs() bool {
	return o.varArgs
}

// Free returns the slice of *ObjectPtr representing the free variables captured by the compiled function.
func (o *CompiledFunction) Free() []*ObjectPtr {
	return o.free
}

// TypeName returns the type name of the object as a string, specifically "compiled-function".
func (o *CompiledFunction) TypeName() string {
	return "compiled-function"
}

// String returns a human-readable string representation of the CompiledFunction object.
func (o *CompiledFunction) String() string {
	return "<compiled-function>"
}

// Copy creates and returns a new instance of CompiledFunction, duplicating its state, except for its variable pointers.
func (o *CompiledFunction) Copy() IObject {
	return &CompiledFunction{
		instructions:  append([]byte{}, o.instructions...),
		numLocals:     o.numLocals,
		numParameters: o.numParameters,
		varArgs:       o.varArgs,
		free:          append([]*ObjectPtr{}, o.free...), // DO NOT Copy() of elements; these are variable pointers
	}
}

// Equals checks if the current CompiledFunction is equal to the provided IObject. Always returns false.
func (o *CompiledFunction) Equals(_ IObject) bool {
	return false
}

// SourcePos retrieves the source position for the instruction pointer (ip) in the compiled function's source map.
// If the ip is not found in the map, it decrements ip until a valid position is found or returns NoPos if none exists.
func (o *CompiledFunction) SourcePos(ip int) compiler.Pos {
	for ip >= 0 {
		if p, ok := o.sourceMap[ip]; ok {
			return p
		}
		ip--
	}
	return compiler.NoPos
}

// CanCall determines if the compiled function object can be invoked as a callable, always returning true.
func (o *CompiledFunction) CanCall() bool {
	return true
}
