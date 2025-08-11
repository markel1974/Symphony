package objects

// CompiledFunction represents a compiled sequence of instructions with metadata for execution in a virtual machine.
type CompiledFunction struct {
	ObjectImpl
	Instructions  []byte
	NumLocals     int
	NumParameters int
	VarArgs       bool
	SourceMap     map[int]Pos
	Free          []*ObjectPtr
}

// TypeName returns the string "compiled-function" representing the object's type.
func (o *CompiledFunction) TypeName() string {
	return "compiled-function"
}

// String returns a fixed string representation "<compiled-function>" for the CompiledFunction object.
func (o *CompiledFunction) String() string {
	return "<compiled-function>"
}

// Copy creates a new instance of CompiledFunction with duplicated Instructions and Free slices, preserving the original values.
func (o *CompiledFunction) Copy() IObject {
	return &CompiledFunction{
		Instructions:  append([]byte{}, o.Instructions...),
		NumLocals:     o.NumLocals,
		NumParameters: o.NumParameters,
		VarArgs:       o.VarArgs,
		Free:          append([]*ObjectPtr{}, o.Free...), // DO NOT Copy() of elements; these are variable pointers
	}
}

// Equals checks whether the current object is equal to another IObject and returns false in this implementation.
func (o *CompiledFunction) Equals(_ IObject) bool {
	return false
}

// SourcePos returns the source position associated with the given instruction pointer (ip), or NoPos if not found.
func (o *CompiledFunction) SourcePos(ip int) Pos {
	for ip >= 0 {
		if p, ok := o.SourceMap[ip]; ok {
			return p
		}
		ip--
	}
	return NoPos
}

// CanCall checks if the CompiledFunction object can be called as a function, always returning true.
func (o *CompiledFunction) CanCall() bool {
	return true
}
