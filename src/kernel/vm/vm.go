package vm

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecodes"
	"github.com/markel1974/c64emu/src/kernel/vm/errors"
	"github.com/markel1974/c64emu/src/kernel/vm/modules"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
	"github.com/markel1974/c64emu/src/kernel/vm/opcodes"
	"github.com/markel1974/c64emu/src/kernel/vm/tokens"
)

const (
	sequenceLen  = 1 << 8
	sequenceMask = sequenceLen - 1
)

type ISequencer interface {
	Create(vm *VM) []func()
}

type VM struct {
	constants       []objects.Object
	stack           [objects.StackSize]objects.Object
	sp              int
	globals         []objects.Object
	fileSet         *objects.SourceFileSet
	frames          [objects.MaxFrames]frame
	framesIndex     int
	curFrame        *frame
	curInstructions []byte
	ip              int
	abort           bool
	suspend         bool
	maxAllocations  int64
	allocations     int64
	err             error
	sequencer       []func()
}

func NewVM(sequencer ISequencer, bytecode *bytecodes.Bytecode, globals []objects.Object, maxAllocations int64) *VM {
	if globals == nil {
		globals = make([]objects.Object, objects.GlobalsSize)
	}
	v := &VM{
		constants:      bytecode.Constants,
		sp:             0,
		globals:        globals,
		fileSet:        bytecode.FileSet,
		framesIndex:    1,
		ip:             -1,
		maxAllocations: maxAllocations,
		suspend:        false,
	}
	v.frames[0].fn = bytecode.MainFunction
	v.frames[0].ip = -1
	v.curFrame = &v.frames[0]
	v.curInstructions = v.curFrame.fn.Instructions
	v.sequencer = sequencer.Create(v)
	return v
}

func (v *VM) Abort() {
	v.abort = true
}

func (v *VM) Run() error {
	v.sp = 0
	v.curFrame = &(v.frames[0])
	v.curInstructions = v.curFrame.fn.Instructions
	v.framesIndex = 1
	v.ip = -1
	v.allocations = v.maxAllocations + 1

	v.run()
	v.abort = false
	v.suspend = false
	if v.err != nil {
		filePos := v.fileSet.Position(v.curFrame.fn.SourcePos(v.ip - 1))
		err := fmt.Errorf("runtime error %w at %s", v.err, filePos)
		for v.framesIndex > 1 {
			v.framesIndex--
			v.curFrame = &v.frames[v.framesIndex-1]
			filePos = v.fileSet.Position(v.curFrame.fn.SourcePos(v.curFrame.ip - 1))
			err = fmt.Errorf("%w at %s", err, filePos)
		}
		return err
	}
	return nil
}

func (v *VM) run() {
	for !v.abort || !v.suspend || v.err == nil {
		v.ip++
		v.sequencer[v.ip&sequenceMask]()
	}
}

// IsStackEmpty tests if the stack is empty or not.
func (v *VM) IsStackEmpty() bool {
	return v.sp == 0
}

func (v *VM) indexAssign(dst, src objects.Object, selectors []objects.Object) error {
	numSel := len(selectors)
	for sIdx := numSel - 1; sIdx > 0; sIdx-- {
		next, err := dst.IndexGet(selectors[sIdx])
		if err != nil {
			if errors.Is(err, errors.ErrNotIndexable) {
				return fmt.Errorf("not indexable: %s", dst.TypeName())
			}
			if errors.Is(err, errors.ErrInvalidIndexType) {
				return fmt.Errorf("invalid index type: %s",
					selectors[sIdx].TypeName())
			}
			return err
		}
		dst = next
	}

	if err := dst.IndexSet(selectors[0], src); err != nil {
		if errors.Is(err, errors.ErrNotIndexAssignable) {
			return fmt.Errorf("not index-assignable: %s", dst.TypeName())
		}
		if errors.Is(err, errors.ErrInvalidIndexValueType) {
			return fmt.Errorf("invaid index value type: %s", src.TypeName())
		}
		return err
	}
	return nil
}

func (v *VM) doOpConstant() {
	v.ip += 2
	cIdx := int(v.curInstructions[v.ip]) | int(v.curInstructions[v.ip-1])<<8
	v.stack[v.sp] = v.constants[cIdx]
	v.sp++
}

func (v *VM) doOpNull() {
	v.stack[v.sp] = objects.UndefinedValue
	v.sp++
}

func (v *VM) doOpBinary() {
	v.ip++
	right := v.stack[v.sp-1]
	left := v.stack[v.sp-2]
	tok := tokens.Token(v.curInstructions[v.ip])
	res, e := left.BinaryOp(tok, right)
	if e != nil {
		v.sp -= 2
		if errors.Is(e, errors.ErrInvalidOperator) {
			v.err = fmt.Errorf("invalid operation: %s %s %s", left.TypeName(), tok.String(), right.TypeName())
		}
		v.err = e
		return
	}
	v.allocations--
	if v.allocations == 0 {
		v.err = errors.ErrObjectAllocLimit
		return
	}
	v.stack[v.sp-2] = res
	v.sp--
}

func (v *VM) doOpEqual() {
	right := v.stack[v.sp-1]
	left := v.stack[v.sp-2]
	v.sp -= 2
	if left.Equals(right) {
		v.stack[v.sp] = objects.TrueValue
	} else {
		v.stack[v.sp] = objects.FalseValue
	}
	v.sp++
}

func (v *VM) doOpNotEqual() {
	right := v.stack[v.sp-1]
	left := v.stack[v.sp-2]
	v.sp -= 2
	if left.Equals(right) {
		v.stack[v.sp] = objects.FalseValue
	} else {
		v.stack[v.sp] = objects.TrueValue
	}
	v.sp++
}
func (v *VM) doOpPop() {
	v.sp--
}

func (v *VM) doOpTrue() {
	v.stack[v.sp] = objects.TrueValue
	v.sp++
}

func (v *VM) doOpFalse() {
	v.stack[v.sp] = objects.FalseValue
	v.sp++
}

func (v *VM) doOpLNot() {
	operand := v.stack[v.sp-1]
	v.sp--
	if operand.IsFalsy() {
		v.stack[v.sp] = objects.TrueValue
	} else {
		v.stack[v.sp] = objects.FalseValue
	}
	v.sp++
}

func (v *VM) doOpBComplement() {
	operand := v.stack[v.sp-1]
	v.sp--
	switch x := operand.(type) {
	case *objects.Int:
		var res objects.Object = &objects.Int{Value: ^x.Value}
		v.allocations--
		if v.allocations == 0 {
			v.err = errors.ErrObjectAllocLimit
			return
		}
		v.stack[v.sp] = res
		v.sp++
	default:
		v.err = fmt.Errorf("invalid operation: ^%s", operand.TypeName())
		return
	}
}

func (v *VM) doOpMinus() {
	operand := v.stack[v.sp-1]
	v.sp--

	switch x := operand.(type) {
	case *objects.Int:
		var res objects.Object = &objects.Int{Value: -x.Value}
		v.allocations--
		if v.allocations == 0 {
			v.err = errors.ErrObjectAllocLimit
			return
		}
		v.stack[v.sp] = res
		v.sp++
	case *objects.Float:
		var res objects.Object = &objects.Float{Value: -x.Value}
		v.allocations--
		if v.allocations == 0 {
			v.err = errors.ErrObjectAllocLimit
			return
		}
		v.stack[v.sp] = res
		v.sp++
	default:
		v.err = fmt.Errorf("invalid operation: -%s", operand.TypeName())
		return
	}
}

func (v *VM) doOpJumpFalsy() {
	v.ip += 2
	v.sp--
	if v.stack[v.sp].IsFalsy() {
		pos := int(v.curInstructions[v.ip]) | int(v.curInstructions[v.ip-1])<<8
		v.ip = pos - 1
	}
}

func (v *VM) doOpAndJump() {
	v.ip += 2
	if v.stack[v.sp-1].IsFalsy() {
		pos := int(v.curInstructions[v.ip]) | int(v.curInstructions[v.ip-1])<<8
		v.ip = pos - 1
	} else {
		v.sp--
	}
}

func (v *VM) doOpOrJump() {
	v.ip += 2
	if v.stack[v.sp-1].IsFalsy() {
		v.sp--
	} else {
		pos := int(v.curInstructions[v.ip]) | int(v.curInstructions[v.ip-1])<<8
		v.ip = pos - 1
	}
}

func (v *VM) doOpJump() {
	pos := int(v.curInstructions[v.ip+2]) | int(v.curInstructions[v.ip+1])<<8
	v.ip = pos - 1
}

func (v *VM) doOpSetGlobal() {
	v.ip += 2
	v.sp--
	globalIndex := int(v.curInstructions[v.ip]) | int(v.curInstructions[v.ip-1])<<8
	v.globals[globalIndex] = v.stack[v.sp]
}

func (v *VM) doOpSetSelGlobal() {
	v.ip += 3
	globalIndex := int(v.curInstructions[v.ip-1]) | int(v.curInstructions[v.ip-2])<<8
	numSelectors := int(v.curInstructions[v.ip])

	// selectors and RHS value
	selectors := make([]objects.Object, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = v.stack[v.sp-numSelectors+i]
	}
	val := v.stack[v.sp-numSelectors-1]
	v.sp -= numSelectors + 1
	e := v.indexAssign(v.globals[globalIndex], val, selectors)
	if e != nil {
		v.err = e
		return
	}
}

func (v *VM) doOpGetGlobal() {
	v.ip += 2
	globalIndex := int(v.curInstructions[v.ip]) | int(v.curInstructions[v.ip-1])<<8
	val := v.globals[globalIndex]
	v.stack[v.sp] = val
	v.sp++
}

func (v *VM) doOpArray() {
	v.ip += 2
	numElements := int(v.curInstructions[v.ip]) | int(v.curInstructions[v.ip-1])<<8
	var elements []objects.Object
	for i := v.sp - numElements; i < v.sp; i++ {
		elements = append(elements, v.stack[i])
	}
	v.sp -= numElements
	arr := &objects.Array{Value: elements}
	v.allocations--
	if v.allocations == 0 {
		v.err = errors.ErrObjectAllocLimit
		return
	}
	v.stack[v.sp] = arr
	v.sp++
}

func (v *VM) doOpMap() {
	v.ip += 2
	numElements := int(v.curInstructions[v.ip]) | int(v.curInstructions[v.ip-1])<<8
	kv := make(map[string]objects.Object, numElements)
	for i := v.sp - numElements; i < v.sp; i += 2 {
		key := v.stack[i]
		value := v.stack[i+1]
		kv[key.(*objects.String).Value] = value
	}
	v.sp -= numElements
	m := &objects.Map{Value: kv}
	v.allocations--
	if v.allocations == 0 {
		v.err = errors.ErrObjectAllocLimit
		return
	}
	v.stack[v.sp] = m
	v.sp++
}

func (v *VM) doOpError() {
	value := v.stack[v.sp-1]
	var e objects.Object = &objects.Error{
		Value: value,
	}
	v.allocations--
	if v.allocations == 0 {
		v.err = errors.ErrObjectAllocLimit
		return
	}
	v.stack[v.sp-1] = e
}

func (v *VM) doOpImmutable() {
	value := v.stack[v.sp-1]
	switch value := value.(type) {
	case *objects.Array:
		var immutableArray objects.Object = &objects.ImmutableArray{
			Value: value.Value,
		}
		v.allocations--
		if v.allocations == 0 {
			v.err = errors.ErrObjectAllocLimit
			return
		}
		v.stack[v.sp-1] = immutableArray
	case *objects.Map:
		var immutableMap objects.Object = &objects.ImmutableMap{
			Value: value.Value,
		}
		v.allocations--
		if v.allocations == 0 {
			v.err = errors.ErrObjectAllocLimit
			return
		}
		v.stack[v.sp-1] = immutableMap
	}
}

func (v *VM) doOpIndex() {
	index := v.stack[v.sp-1]
	left := v.stack[v.sp-2]
	v.sp -= 2

	val, err := left.IndexGet(index)
	if err != nil {
		if err == errors.ErrNotIndexable {
			v.err = fmt.Errorf("not indexable: %s", index.TypeName())
			return
		}
		if err == errors.ErrInvalidIndexType {
			v.err = fmt.Errorf("invalid index type: %s",
				index.TypeName())
			return
		}
		v.err = err
		return
	}
	if val == nil {
		val = objects.UndefinedValue
	}
	v.stack[v.sp] = val
	v.sp++
}

func (v *VM) doOpSliceIndex() {
	highStack := v.stack[v.sp-1]
	lowStack := v.stack[v.sp-2]
	leftStack := v.stack[v.sp-3]
	v.sp -= 3

	var lowIdx int64
	if lowStack != objects.UndefinedValue {
		if low, ok := lowStack.(*objects.Int); ok {
			lowIdx = low.Value
		} else {
			v.err = fmt.Errorf("invalid slice index type: %s",
				low.TypeName())
			return
		}
	}

	switch left := leftStack.(type) {
	case *objects.Array:
		numElements := int64(len(left.Value))
		var highIdx int64
		if highStack == objects.UndefinedValue {
			highIdx = numElements
		} else if high, ok := highStack.(*objects.Int); ok {
			highIdx = high.Value
		} else {
			v.err = fmt.Errorf("invalid slice index type: %s",
				high.TypeName())
			return
		}
		if lowIdx > highIdx {
			v.err = fmt.Errorf("invalid slice index: %d > %d",
				lowIdx, highIdx)
			return
		}
		if lowIdx < 0 {
			lowIdx = 0
		} else if lowIdx > numElements {
			lowIdx = numElements
		}
		if highIdx < 0 {
			highIdx = 0
		} else if highIdx > numElements {
			highIdx = numElements
		}
		var val objects.Object = &objects.Array{
			Value: left.Value[lowIdx:highIdx],
		}
		v.allocations--
		if v.allocations == 0 {
			v.err = errors.ErrObjectAllocLimit
			return
		}
		v.stack[v.sp] = val
		v.sp++

	case *objects.ImmutableArray:
		numElements := int64(len(left.Value))
		var highIdx int64
		if highStack == objects.UndefinedValue {
			highIdx = numElements
		} else if high, ok := highStack.(*objects.Int); ok {
			highIdx = high.Value
		} else {
			v.err = fmt.Errorf("invalid slice index type: %s",
				high.TypeName())
			return
		}
		if lowIdx > highIdx {
			v.err = fmt.Errorf("invalid slice index: %d > %d",
				lowIdx, highIdx)
			return
		}
		if lowIdx < 0 {
			lowIdx = 0
		} else if lowIdx > numElements {
			lowIdx = numElements
		}
		if highIdx < 0 {
			highIdx = 0
		} else if highIdx > numElements {
			highIdx = numElements
		}
		var val objects.Object = &objects.Array{
			Value: left.Value[lowIdx:highIdx],
		}
		v.allocations--
		if v.allocations == 0 {
			v.err = errors.ErrObjectAllocLimit
			return
		}
		v.stack[v.sp] = val
		v.sp++
	case *objects.String:
		numElements := int64(len(left.Value))
		var highIdx int64
		if highStack == objects.UndefinedValue {
			highIdx = numElements
		} else if high, ok := highStack.(*objects.Int); ok {
			highIdx = high.Value
		} else {
			v.err = fmt.Errorf("invalid slice index type: %s",
				high.TypeName())
			return
		}
		if lowIdx > highIdx {
			v.err = fmt.Errorf("invalid slice index: %d > %d",
				lowIdx, highIdx)
			return
		}
		if lowIdx < 0 {
			lowIdx = 0
		} else if lowIdx > numElements {
			lowIdx = numElements
		}
		if highIdx < 0 {
			highIdx = 0
		} else if highIdx > numElements {
			highIdx = numElements
		}
		var val objects.Object = &objects.String{
			Value: left.Value[lowIdx:highIdx],
		}
		v.allocations--
		if v.allocations == 0 {
			v.err = errors.ErrObjectAllocLimit
			return
		}
		v.stack[v.sp] = val
		v.sp++
	case *objects.Bytes:
		numElements := int64(len(left.Value))
		var highIdx int64
		if highStack == objects.UndefinedValue {
			highIdx = numElements
		} else if high, ok := highStack.(*objects.Int); ok {
			highIdx = high.Value
		} else {
			v.err = fmt.Errorf("invalid slice index type: %s",
				high.TypeName())
			return
		}
		if lowIdx > highIdx {
			v.err = fmt.Errorf("invalid slice index: %d > %d",
				lowIdx, highIdx)
			return
		}
		if lowIdx < 0 {
			lowIdx = 0
		} else if lowIdx > numElements {
			lowIdx = numElements
		}
		if highIdx < 0 {
			highIdx = 0
		} else if highIdx > numElements {
			highIdx = numElements
		}
		var val objects.Object = &objects.Bytes{
			Value: left.Value[lowIdx:highIdx],
		}
		v.allocations--
		if v.allocations == 0 {
			v.err = errors.ErrObjectAllocLimit
			return
		}
		v.stack[v.sp] = val
		v.sp++
	}
}

func (v *VM) doOpCall() {
	numArgs := int(v.curInstructions[v.ip+1])
	spread := int(v.curInstructions[v.ip+2])
	v.ip += 2

	value := v.stack[v.sp-1-numArgs]
	if !value.CanCall() {
		v.err = fmt.Errorf("not callable: %s", value.TypeName())
		return
	}

	if spread == 1 {
		v.sp--
		switch arr := v.stack[v.sp].(type) {
		case *objects.Array:
			for _, item := range arr.Value {
				v.stack[v.sp] = item
				v.sp++
			}
			numArgs += len(arr.Value) - 1
		case *objects.ImmutableArray:
			for _, item := range arr.Value {
				v.stack[v.sp] = item
				v.sp++
			}
			numArgs += len(arr.Value) - 1
		default:
			v.err = fmt.Errorf("not an array: %s", arr.TypeName())
			return
		}
	}

	if callee, ok := value.(*objects.CompiledFunction); ok {
		if callee.VarArgs {
			// if the closure is variadic, roll up all variadic parameters into an array
			realArgs := callee.NumParameters - 1
			varArgs := numArgs - realArgs
			if varArgs >= 0 {
				numArgs = realArgs + 1
				args := make([]objects.Object, varArgs)
				spStart := v.sp - varArgs
				for i := spStart; i < v.sp; i++ {
					args[i-spStart] = v.stack[i]
				}
				v.stack[spStart] = &objects.Array{Value: args}
				v.sp = spStart + 1
			}
		}
		if numArgs != callee.NumParameters {
			numParams := callee.NumParameters
			if callee.VarArgs {
				numParams = callee.NumParameters - 1
			}
			v.err = fmt.Errorf("wrong number of arguments: want>=%d, got=%d", numParams, numArgs)
			return
		}

		// test if it's tail-call
		if callee == v.curFrame.fn { // recursion
			nextOp := v.curInstructions[v.ip+1]
			if nextOp == opcodes.OpReturn ||
				(nextOp == opcodes.OpPop &&
					opcodes.OpReturn == v.curInstructions[v.ip+2]) {
				for p := 0; p < numArgs; p++ {
					v.stack[v.curFrame.basePointer+p] =
						v.stack[v.sp-numArgs+p]
				}
				v.sp -= numArgs + 1
				v.ip = -1 // reset IP to beginning of the frame
				return
				//continue
			}
		}
		if v.framesIndex >= objects.MaxFrames {
			v.err = errors.ErrStackOverflow
			return
		}

		// update call frame
		v.curFrame.ip = v.ip // store current ip before call
		v.curFrame = &(v.frames[v.framesIndex])
		v.curFrame.fn = callee
		v.curFrame.freeVars = callee.Free
		v.curFrame.basePointer = v.sp - numArgs
		v.curInstructions = callee.Instructions
		v.ip = -1
		v.framesIndex++
		v.sp = v.sp - numArgs + callee.NumLocals
	} else {
		var args []objects.Object
		args = append(args, v.stack[v.sp-numArgs:v.sp]...)
		ret, e := value.Call(args...)
		v.sp -= numArgs + 1

		// runtime error
		if e != nil {
			if e == errors.ErrWrongNumArguments {
				v.err = fmt.Errorf("wrong number of arguments in call to '%s'", value.TypeName())
				return
			}
			if e, ok := e.(errors.ErrInvalidArgumentType); ok {
				v.err = fmt.Errorf("invalid type for argument '%s' in call to '%s': "+"expected %s, found %s", e.Name, value.TypeName(), e.Expected, e.Found)
				return
			}
			v.err = e
			return
		}

		// nil return -> undefined
		if ret == nil {
			ret = objects.UndefinedValue
		}
		v.allocations--
		if v.allocations == 0 {
			v.err = errors.ErrObjectAllocLimit
			return
		}
		v.stack[v.sp] = ret
		v.sp++
	}
}
func (v *VM) doOpReturn() {
	v.ip++
	var retVal objects.Object
	if int(v.curInstructions[v.ip]) == 1 {
		retVal = v.stack[v.sp-1]
	} else {
		retVal = objects.UndefinedValue
	}
	//v.sp--
	v.framesIndex--
	v.curFrame = &v.frames[v.framesIndex-1]
	v.curInstructions = v.curFrame.fn.Instructions
	v.ip = v.curFrame.ip
	//v.sp = lastFrame.basePointer - 1
	v.sp = v.frames[v.framesIndex].basePointer
	// skip stack overflow check because (newSP) <= (oldSP)
	v.stack[v.sp-1] = retVal
	//v.sp++
}

func (v *VM) doOpDefineLocal() {
	v.ip++
	localIndex := int(v.curInstructions[v.ip])
	sp := v.curFrame.basePointer + localIndex

	// local variables can be mutated by other actions
	// so always store the copy of popped value
	val := v.stack[v.sp-1]
	v.sp--
	v.stack[sp] = val
}

func (v *VM) doOpSetLocal() {
	localIndex := int(v.curInstructions[v.ip+1])
	v.ip++
	sp := v.curFrame.basePointer + localIndex

	// update pointee of v.stack[sp] instead of replacing the pointer
	// itself. this is needed because there can be free variables
	// referencing the same local variables.
	val := v.stack[v.sp-1]
	v.sp--
	if obj, ok := v.stack[sp].(*objects.ObjectPtr); ok {
		*obj.Value = val
		val = obj
	}
	v.stack[sp] = val // also use a copy of popped value
}

func (v *VM) doOpSetSelLocal() {
	localIndex := int(v.curInstructions[v.ip+1])
	numSelectors := int(v.curInstructions[v.ip+2])
	v.ip += 2

	// selectors and RHS value
	selectors := make([]objects.Object, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = v.stack[v.sp-numSelectors+i]
	}
	val := v.stack[v.sp-numSelectors-1]
	v.sp -= numSelectors + 1
	dst := v.stack[v.curFrame.basePointer+localIndex]
	if obj, ok := dst.(*objects.ObjectPtr); ok {
		dst = *obj.Value
	}
	if e := v.indexAssign(dst, val, selectors); e != nil {
		v.err = e
		return
	}
}

func (v *VM) doOpGetLocal() {
	v.ip++
	localIndex := int(v.curInstructions[v.ip])
	val := v.stack[v.curFrame.basePointer+localIndex]
	if obj, ok := val.(*objects.ObjectPtr); ok {
		val = *obj.Value
	}
	v.stack[v.sp] = val
	v.sp++
}

func (v *VM) doOpGetBuiltin() {
	v.ip++
	builtinIndex := int(v.curInstructions[v.ip])
	v.stack[v.sp] = modules.GetBuiltin(builtinIndex)
	v.sp++
}

func (v *VM) doOpClosure() {
	v.ip += 3
	constIndex := int(v.curInstructions[v.ip-1]) | int(v.curInstructions[v.ip-2])<<8
	numFree := int(v.curInstructions[v.ip])
	fn, ok := v.constants[constIndex].(*objects.CompiledFunction)
	if !ok {
		v.err = fmt.Errorf("not function: %s", fn.TypeName())
		return
	}
	free := make([]*objects.ObjectPtr, numFree)
	for i := 0; i < numFree; i++ {
		switch freeVar := (v.stack[v.sp-numFree+i]).(type) {
		case *objects.ObjectPtr:
			free[i] = freeVar
		default:
			free[i] = &objects.ObjectPtr{
				Value: &v.stack[v.sp-numFree+i],
			}
		}
	}
	v.sp -= numFree
	cl := &objects.CompiledFunction{
		Instructions:  fn.Instructions,
		NumLocals:     fn.NumLocals,
		NumParameters: fn.NumParameters,
		VarArgs:       fn.VarArgs,
		Free:          free,
	}
	v.allocations--
	if v.allocations == 0 {
		v.err = errors.ErrObjectAllocLimit
		return
	}
	v.stack[v.sp] = cl
	v.sp++
}

func (v *VM) doOpGetFreePtr() {
	v.ip++
	freeIndex := int(v.curInstructions[v.ip])
	val := v.curFrame.freeVars[freeIndex]
	v.stack[v.sp] = val
	v.sp++
}

func (v *VM) doOpGetFree() {
	v.ip++
	freeIndex := int(v.curInstructions[v.ip])
	val := *v.curFrame.freeVars[freeIndex].Value
	v.stack[v.sp] = val
	v.sp++
}

func (v *VM) doOpSetFree() {
	v.ip++
	freeIndex := int(v.curInstructions[v.ip])
	*v.curFrame.freeVars[freeIndex].Value = v.stack[v.sp-1]
	v.sp--
}

func (v *VM) doOpGetLocalPtr() {
	v.ip++
	localIndex := int(v.curInstructions[v.ip])
	sp := v.curFrame.basePointer + localIndex
	val := v.stack[sp]
	var freeVar *objects.ObjectPtr
	if obj, ok := val.(*objects.ObjectPtr); ok {
		freeVar = obj
	} else {
		freeVar = &objects.ObjectPtr{Value: &val}
		v.stack[sp] = freeVar
	}
	v.stack[v.sp] = freeVar
	v.sp++
}

func (v *VM) doOpSetSelFree() {
	v.ip += 2
	freeIndex := int(v.curInstructions[v.ip-1])
	numSelectors := int(v.curInstructions[v.ip])

	// selectors and RHS value
	selectors := make([]objects.Object, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = v.stack[v.sp-numSelectors+i]
	}
	val := v.stack[v.sp-numSelectors-1]
	v.sp -= numSelectors + 1
	e := v.indexAssign(*v.curFrame.freeVars[freeIndex].Value,
		val, selectors)
	if e != nil {
		v.err = e
		return
	}
}

func (v *VM) doOpIteratorInit() {
	var iterator objects.Object
	dst := v.stack[v.sp-1]
	v.sp--
	if !dst.CanIterate() {
		v.err = fmt.Errorf("not iterable: %s", dst.TypeName())
		return
	}
	iterator = dst.Iterate()
	v.allocations--
	if v.allocations == 0 {
		v.err = errors.ErrObjectAllocLimit
		return
	}
	v.stack[v.sp] = iterator
	v.sp++
}

func (v *VM) doOpIteratorNext() {
	iterator := v.stack[v.sp-1]
	v.sp--
	hasMore := iterator.(objects.Iterator).Next()
	if hasMore {
		v.stack[v.sp] = objects.TrueValue
	} else {
		v.stack[v.sp] = objects.FalseValue
	}
	v.sp++
}

func (v *VM) doOpIteratorKey() {
	iterator := v.stack[v.sp-1]
	v.sp--
	val := iterator.(objects.Iterator).Key()
	v.stack[v.sp] = val
	v.sp++
}

func (v *VM) doOpIteratorValue() {
	iterator := v.stack[v.sp-1]
	v.sp--
	val := iterator.(objects.Iterator).Value()
	v.stack[v.sp] = val
	v.sp++
}

func (v *VM) doOpSuspend() {
	v.suspend = true
}

func (v *VM) doOpUnknown() {
	v.err = fmt.Errorf("unknown opcode: %d", v.curInstructions[v.ip])
}
