package objects

// FunctionCompiled represents a compiled function with bytecode instructions, metadata, and associated free variables.
type FunctionCompiled struct {
	ObjectImpl
	instructions  *Instructions
	numLocals     int
	numParameters int
	varArgs       bool
	sourceMap     map[int]int
	free          []*ObjectPointer
}

// NewFunctionCompiled creates a new instance of FunctionCompiled with the given instructions, locals, parameters, varArgs, sourceMap, and free vars.
func NewFunctionCompiled(instructions []byte, numLocals int, numParameters int, varArgs bool, sourceMap map[int]int, free []*ObjectPointer) *FunctionCompiled {
	if sourceMap == nil {
		sourceMap = make(map[int]int)
	}
	return &FunctionCompiled{
		instructions:  NewInstructions(instructions),
		numLocals:     numLocals,
		numParameters: numParameters,
		varArgs:       varArgs,
		sourceMap:     sourceMap,
		free:          free,
	}
}

// Data returns the bytecode instructions of the compiled function.
func (o *FunctionCompiled) Data() []byte {
	return o.instructions.Data()
}

// Instructions returns the bytecode instructions associated with the compiled function.
func (o *FunctionCompiled) Instructions() *Instructions {
	return o.instructions
}

// NumLocals returns the number of local variables required by the compiled function.
func (o *FunctionCompiled) NumLocals() int {
	return o.numLocals
}

// NumParameters returns the total number of parameters required by the compiled function.
func (o *FunctionCompiled) NumParameters() int {
	return o.numParameters
}

// VarArgs returns true if the function accepts a variable number of arguments, otherwise false.
func (o *FunctionCompiled) VarArgs() bool {
	return o.varArgs
}

// Free returns the slice of *ObjectPointer representing the free variables captured by the compiled function.
func (o *FunctionCompiled) Free() []*ObjectPointer {
	return o.free
}

// TypeName returns the type name of the object as a string, specifically "compiled-function".
func (o *FunctionCompiled) TypeName() string {
	return "compiled-function"
}

// String returns a human-readable string representation of the FunctionCompiled object.
func (o *FunctionCompiled) String() string {
	return "<compiled-function>"
}

// Copy creates and returns a new instance of FunctionCompiled, duplicating its state, except for its variable pointers.
func (o *FunctionCompiled) Copy() IObject {
	return &FunctionCompiled{
		instructions:  o.instructions.Copy(), //append([]byte{}, o.instructions...),
		numLocals:     o.numLocals,
		numParameters: o.numParameters,
		varArgs:       o.varArgs,
		free:          append([]*ObjectPointer{}, o.free...), // DO NOT Copy() of elements; these are variable pointers
	}
}

// Equals checks if the current FunctionCompiled is equal to the provided IObject. Always returns false.
func (o *FunctionCompiled) Equals(_ IObject) bool {
	return false
}

// SourcePos retrieves the source position for the instruction pointer (ip) in the compiled function's source map.
// If the ip is not found in the map, it decrements ip until a valid position is found or returns NoPos if none exists.
func (o *FunctionCompiled) SourcePos(ip int) int {
	for ip >= 0 {
		if p, ok := o.sourceMap[ip]; ok {
			return p
		}
		ip--
	}
	return 0
}

// CanCall determines if the compiled function object can be invoked as a callable, always returning true.
func (o *FunctionCompiled) CanCall() bool {
	return true
}
