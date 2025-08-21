package objects

const (
	FuncCompiledDef   = "func_compiled"
	FuncCompiledLabel = "<" + FuncCompiledDef + ">"
)

// FuncCompiled represents a compiled function with bytecode instructions, metadata, and associated free variables.
type FuncCompiled struct {
	IObject
	name          string
	instructions  *Instructions
	numLocals     int
	numParameters int
	varArgs       bool
	sourceMap     map[int]int
	free          []*ObjectPointer
}

// NewFunctionCompiled creates a new instance of FuncCompiled with the given instructions, locals, parameters, varArgs, sourceMap, and free vars.
func newFuncCompiled(factory *GateKeeper, frame int, name string, instructions []byte, numLocals int, numParameters int, varArgs bool, sourceMap map[int]int, free []*ObjectPointer) *FuncCompiled {
	if sourceMap == nil {
		sourceMap = make(map[int]int)
	}
	return &FuncCompiled{
		IObject:       factory.newObject(frame),
		name:          name,
		instructions:  _newInstructions(instructions),
		numLocals:     numLocals,
		numParameters: numParameters,
		varArgs:       varArgs,
		sourceMap:     sourceMap,
		free:          free,
	}
}

// Name returns the name of the compiled function.
func (o *FuncCompiled) Name() string {
	return o.name
}

// Data returns the bytecode instructions of the compiled function.
func (o *FuncCompiled) Data() []byte {
	return o.instructions.Data()
}

// Instructions returns the bytecode instructions associated with the compiled function.
func (o *FuncCompiled) Instructions() *Instructions {
	return o.instructions
}

// NumLocals returns the number of local variables required by the compiled function.
func (o *FuncCompiled) NumLocals() int {
	return o.numLocals
}

// NumParameters returns the total number of parameters required by the compiled function.
func (o *FuncCompiled) NumParameters() int {
	return o.numParameters
}

// VarArgs returns true if the function accepts a variable number of arguments, otherwise false.
func (o *FuncCompiled) VarArgs() bool {
	return o.varArgs
}

// Free returns the slice of *ObjectPointer representing the free variables captured by the compiled function.
func (o *FuncCompiled) Free() []*ObjectPointer {
	return o.free
}

// TypeName returns the type name of the object as a string, specifically "compiled-function".
func (o *FuncCompiled) TypeName() string {
	return FuncCompiledDef
}

// String returns a human-readable string representation of the FuncCompiled object.
func (o *FuncCompiled) String() string {
	return FuncCompiledLabel
}

// Copy creates and returns a new instance of FuncCompiled, duplicating its state, except for its variable pointers.
func (o *FuncCompiled) Copy(frame int, _ int) IObject {
	obj := o.GateKeeper().NewFuncCompiled(frame, o.name, nil, o.numLocals, o.numParameters, o.varArgs, nil, nil)
	ret, ok := obj.(*FuncCompiled)
	if !ok {
		return o.GateKeeper().UndefinedValue()
	}
	ret.instructions = o.instructions.Copy()
	ret.free = append([]*ObjectPointer{}, o.free...)
	return ret
}

// Equals checks if the current FuncCompiled is equal to the provided IObject. Always returns false.
func (o *FuncCompiled) Equals(_ IObject) bool {
	return false
}

// SourcePos retrieves the source position for the instruction pointer (ip) in the compiled function's source map.
// If the ip is not found in the map, it decrements ip until a valid position is found or returns NoPos if none exists.
func (o *FuncCompiled) SourcePos(ip int) int {
	for ip >= 0 {
		if p, ok := o.sourceMap[ip]; ok {
			return p
		}
		ip--
	}
	return 0
}

// CanCall determines if the compiled function object can be invoked as a callable, always returning true.
func (o *FuncCompiled) CanCall() bool {
	return true
}
