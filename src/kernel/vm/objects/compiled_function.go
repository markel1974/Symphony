package objects

// CompiledFunction represents a compiled function.
type CompiledFunction struct {
	ObjectImpl
	Instructions  []byte
	NumLocals     int // number of local variables (including function parameters)
	NumParameters int
	VarArgs       bool
	SourceMap     map[int]Pos
	Free          []*ObjectPtr
}

// TypeName returns the name of the type.
func (o *CompiledFunction) TypeName() string {
	return "compiled-function"
}

func (o *CompiledFunction) String() string {
	return "<compiled-function>"
}

// Copy returns a copy of the type.
func (o *CompiledFunction) Copy() Object {
	return &CompiledFunction{
		Instructions:  append([]byte{}, o.Instructions...),
		NumLocals:     o.NumLocals,
		NumParameters: o.NumParameters,
		VarArgs:       o.VarArgs,
		Free:          append([]*ObjectPtr{}, o.Free...), // DO NOT Copy() of elements; these are variable pointers
	}
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *CompiledFunction) Equals(_ Object) bool {
	return false
}

// SourcePos returns the source position of the instruction at ip.
func (o *CompiledFunction) SourcePos(ip int) Pos {
	for ip >= 0 {
		if p, ok := o.SourceMap[ip]; ok {
			return p
		}
		ip--
	}
	return NoPos
}

// CanCall returns whether the Object can be Called.
func (o *CompiledFunction) CanCall() bool {
	return true
}
