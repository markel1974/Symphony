package core

import (
	"fmt"
	"io"
	"time"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
)

const Version = "0.1"

// resetIp is the instruction pointer value used to reset the VM's instruction pointer to the beginning of the main function.'
const (
	resetIp = -1
)

// stackSize specifies the size limit of the stack for function execution.
// maxFrames indicates the maximum number of call frames allowed.
const (
	stackSize = 2048
	maxFrames = 1024
)

// VM represents a virtual machine that executes bytecode instructions, handles stack, and manages execution frames.
type VM struct {
	gk                objects.IGateKeeper
	op                *bytecode.Opcodes
	bc                *bytecode.Bytecode
	stack             *Stack
	frames            *Frames
	currFrame         *Frame
	ip                int
	shutdown          bool
	err               error
	sequencer         []*Decoder
	sequencerMask     int
	imports           *Imports
	constants         *Constants
	globals           *Globals
	seq               ISequencer
	counterStart      uint64
	counterIterations uint64
	retValues         bool
}

// New initializes and returns a new virtual machine instance configured with the provided components and settings.
func New(gk objects.IGateKeeper, seq ISequencer, op *bytecode.Opcodes) *VM {
	v := &VM{
		gk:        gk,
		op:        op,
		ip:        resetIp,
		imports:   nil,
		retValues: false,
	}
	v.constants = NewConstants(gk, v.SetError)
	v.imports = NewImports(gk, v.SetError)
	v.globals = NewGlobals(gk, v.SetError)
	v.stack = NewStack(gk, stackSize, v.SetError)
	v.frames = NewFrames(gk, maxFrames, v.SetError)
	v.seq = seq
	return v
}

// Setup initializes the virtual machine with the provided bytecode and loader components.
func (v *VM) Setup(loader bytecode.ILoader, codes ...*bytecode.Bytecode) (map[string]uint, error) {
	sequencer, err := v.seq.Create(v)
	if err != nil {
		return nil, err
	}
	v.sequencerMask = len(sequencer) - 1
	v.sequencer = make([]*Decoder, len(sequencer))
	for i, s := range sequencer {
		if v.sequencer[i], err = NewDecoder(s); err != nil {
			return nil, err
		}
	}
	switch len(codes) {
	case 1:
		v.bc = codes[0]
	default:
		relocator := bytecode.NewRelocator(v.gk, loader, v.op, bytecode.PreInitFunction, bytecode.InitFunction)
		if v.bc, err = relocator.Relocate(codes); err != nil {
			return nil, err
		}
	}
	if v.bc == nil {
		return nil, fmt.Errorf("no bytecode provided")
	}

	if err = v.imports.Setup(loader, v.bc.Imports()); err != nil {
		return nil, err
	}
	if err = v.globals.Setup(v.bc.Globals()); err != nil {
		return nil, err
	}
	entryPoints, err := v.constants.Setup(v.bc.Constants(), bytecode.PreInitFunction, bytecode.InitFunction)
	if err != nil {
		return nil, err
	}
	for _, fn := range v.constants.PreInitFuncs() {
		if _, err = v.exec(fn, false); err != nil {
			return nil, err
		}
	}
	for _, fn := range v.constants.InitFuncs() {
		if _, err = v.exec(fn, false); err != nil {
			return nil, err
		}
	}
	return entryPoints, nil
}

// Version returns the version of the virtual machine.
func (v *VM) Version() string {
	return Version
}

// Statistics returns three uint64 values: start, allocated objects, and a counter from the VM instance.
func (v *VM) Statistics() (uint64, uint64, uint64, uint64) {
	return v.counterStart, v.gk.AllocatedObjects(), v.counterIterations, uint64(v.frames.Max())
}

// EnableRetValues sets the flag to enable or disable returning multiple values from the virtual machine's execution.
func (v *VM) EnableRetValues(retValues bool) {
	v.retValues = retValues
}

// Run executes the main function identified by mainId with the provided arguments in the virtual machine context.
func (v *VM) Run(mainId uint, args ...interface{}) ([]interface{}, error) {
	obj := v.constants.Get(mainId)
	mainFn, ok := obj.(*objects.FuncCompiled)
	if !ok {
		return nil, fmt.Errorf("entry point not found: %d", mainId)
	}
	return v.exec(mainFn, v.retValues, args...)
}

// Stack returns the current stack instance associated with the VM.
func (v *VM) Stack() *Stack {
	return v.stack
}

// Constants returns a pointer to the Constants associated with the VM instance.
func (v *VM) Constants() *Constants {
	return v.constants
}

// Globals returns the global variables associated with the VM instance.
func (v *VM) Globals() *Globals {
	return v.globals
}

// Imports return a pointer to the Imports object associated with the VM instance.
func (v *VM) Imports() *Imports {
	return v.imports
}

// Factory returns the IGateKeeper instance associated with the VM.
func (v *VM) Factory() objects.IGateKeeper {
	return v.gk
}

// Frame returns the current frame instance associated with the VM.
func (v *VM) Frame() *Frame {
	return v.currFrame
}

// SetIp sets the virtual machine's instruction pointer to the specified value.
func (v *VM) SetIp(ip int) {
	v.ip = ip
}

// GetIp retrieves the current instruction pointer value from the virtual machine.
func (v *VM) GetIp() int {
	return v.ip
}

// Call executes a function or method with the specified number of arguments and handles variadic functions if applicable.
func (v *VM) Call(value objects.IObject, spread bool, numArgs int) {
	if numArgs < 0 {
		v.SetError(fmt.Errorf("invalid number of arguments: %d", numArgs))
		return
	}

	if spread {
		var args []objects.IObject
		obj := v.Stack().Pop()
		switch z := obj.(type) {
		case *objects.Array:
			args = z.Values()
		default:
			v.SetError(fmt.Errorf("unexpected type (array required): %s", obj.TypeName()))
			return
		}
		if argsLen := len(args); argsLen > 0 {
			for _, item := range args {
				v.Stack().Push(item)
			}
			numArgs += argsLen - 1
		}
	}

	switch ce := value.(type) {
	case *objects.FuncCompiled:
		numParams := ce.NumParameters()
		if ce.VarArgs() {
			if numParams > 0 {
				numParams--
			}
			v.Stack().PushVarArgs(v.currFrame.Id(), numArgs, numParams)
			numArgs = ce.NumParameters()
		} else {
			if numArgs != numParams {
				v.SetError(fmt.Errorf("%s invalid arguments [vargs: %v] : want>=%d, got=%d", ce.Name(), ce.VarArgs(), numParams, numArgs))
				return
			}
		}
		v.prepareForCall(ce, numArgs)
	default:
		var args []objects.IObject
		args = append(args, v.Stack().PeekArrayObject(numArgs)...)
		v.CallObject(value, numArgs, args...)
	}
}

// CallObject invokes a callable object with the given arguments and handles stack cleanup and error management.
func (v *VM) CallObject(value objects.IObject, numArgs int, args ...objects.IObject) {
	retCount, ret, err := value.Call(v.currFrame.Id(), args...)
	v.stack.DecrementCount(numArgs + 1)
	if err != nil {
		v.SetError(objects.ComputeCallError(err, value.TypeName()))
		return
	}
	switch retCount {
	case 0:
		v.stack.Push(v.gk.UndefinedValue())
		return
	case 1:
		if ret != nil {
			v.stack.Push(ret)
		} else {
			v.stack.Push(v.gk.UndefinedValue())
		}
		return
	default:
		container, ok := ret.(*objects.Array)
		if !ok {
			v.SetError(fmt.Errorf("invalid return count: %d", retCount))
			return
		}
		for _, item := range container.Values() {
			v.stack.Push(item)
		}
		return
	}
}

// Return concludes the execution of the current frame and handles return values, including handling deferred calls.
// If deferred calls are present, prepares the frame for execution of the first deferred call without immediate execution.
// Saves return values in the parent frame during 'defer' execution chains and recursively processes parent returns.
// Defaults to standard frame return if no deferred calls or chaining is present.
func (v *VM) Return(returnValues []objects.IObject) {
	// CASE 1: We are in a frame that is ending AND has pending 'defer' calls.
	if v.currFrame.HasDeferredCalls() {
		if deferredCall := v.currFrame.DeferredPop(); deferredCall != nil {
			// Prepare a frame for 'defer' call but don't execute it yet.
			deferredFrame := v.prepareForCall(deferredCall, 0)
			// Save return values of the current frame (parent) in a new 'defer' call frame. This creates the chain link.
			deferredFrame.SaveParentReturnValues(returnValues)
			// Break function. VM will execute the 'defer' call in the next cycle.
			return
		}
	}

	// CASE 2: We are in a frame that just finished (e.g. 'defer' closure) and must restore its parent's return flow.
	if v.currFrame.HasParentReturnValues() {
		savedReturnValues := v.currFrame.PopParentReturnValues()
		// A. Finalize the return of the current frame ('defer' frame).
		v.returnApply(returnValues)
		// B. Now that we're back in the parent frame, continue its return process.
		//	This creates recursion: calling 'Return' with parent values.
		v.Return(savedReturnValues)
	} else {
		// CASE 3: Standard return. No, pending 'defer' calls, and we're not inside a 'defer' chain.
		v.returnApply(returnValues)
	}
}

func (v *VM) returnApply(returnValues []objects.IObject) {
	shutdown := false
	prevIp := v.currFrame.SavedIP()
	leavingFrameBasePointer := v.currFrame.BasePointer()
	leavingFrameId := v.currFrame.Id()

	// Function arguments belong to the caller and should not be released when the function ends.
	numArgs := v.currFrame.NumParameters()
	totalArgs := numArgs + len(returnValues)
	var objectsToPreserve []objects.IObject

	if totalArgs > 0 {
		// Create a combined list, pre-allocating capacity for efficiency
		objectsToPreserve = make([]objects.IObject, 0, totalArgs)
		for _, obj := range returnValues {
			if obj.Frame() != objects.FrameStatic {
				objectsToPreserve = append(objectsToPreserve, obj)
			}
		}
		// add arguments which are located at the start of the current frame
		for i := 0; i < numArgs; i++ {
			obj := v.stack.PeekAbsolute(leavingFrameBasePointer + i)
			if obj.Frame() != objects.FrameStatic {
				objectsToPreserve = append(objectsToPreserve, obj)
			}
		}
	}

	v.stack.ReleaseObjects(leavingFrameId, leavingFrameBasePointer, v.stack.StackPointer(), objectsToPreserve)

	if v.frames.CanMovePrevious() {
		v.frames.MovePrevious()
		v.currFrame = v.frames.Previous()
		v.SetIp(prevIp)
	} else {
		shutdown = true
	}

	v.stack.SetStackPointer(leavingFrameBasePointer)

	if lRet := len(returnValues); lRet > 0 {
		for i := lRet - 1; i >= 0; i-- {
			v.stack.Push(returnValues[i])
		}
	} else {
		v.stack.Push(v.gk.UndefinedValue())
	}

	if shutdown {
		v.Shutdown()
	}
}

// ReseIp resets the instruction pointer of the virtual machine to its initial reset state defined by `resetIp`.
func (v *VM) ReseIp() {
	v.ip = resetIp
}

// Shutdown gracefully shuts down the virtual machine by setting its internal state to signify termination.
func (v *VM) Shutdown() {
	v.shutdown = true
}

// Print prints the current state of the virtual machine's stack to the console.'
func (v *VM) Print(writer io.Writer) {
	v.stack.Print(writer)
}

// GetReturnValue returns the value from the top of the stack as an interface value.
func (v *VM) GetReturnValue(idx int) interface{} {
	obj := v.stack.PeekAbsolute(idx)
	if obj == nil {
		return nil
	}
	return v.gk.ToInterface(obj)
}

// GetReturnValues returns the values from the top of the stack as an array of interface values.
func (v *VM) GetReturnValues() []interface{} {
	values := v.stack.StackPointer()
	if values == 0 {
		return nil
	}
	out := make([]interface{}, values)
	for x := 0; x < values; x++ {
		out[x] = v.GetReturnValue(x)
	}
	return out
}

// SetError sets the internal error state of the VM and marks it for shutdown.
func (v *VM) SetError(err error) {
	v.err = err
	v.shutdown = true
}

// Rewrite updates the global function mapping with the given JIT-compiled function for the specified identifier.
func (v *VM) Rewrite(id uint, jit *objects.FuncJit) {
	v.globals.Set(id, jit)
}

// Reset reinitializes the virtual machine's state, clears the stack and frames, and resets execution-related variables.
func (v *VM) prepare() {
	v.ip = resetIp
	v.gk.Reset()
	v.stack.Reset()
	v.frames.Reset()
	v.err = nil
	v.shutdown = false
}

// Run executes the virtual machine's bytecode, managing the stack, frames, and instruction pointer state.
func (v *VM) exec(mainFn *objects.FuncCompiled, ret bool, args ...interface{}) ([]interface{}, error) {
	v.prepare()
	defer func() {
		v.stack.ReleaseAll()
	}()
	v.currFrame = v.frames.Head()
	v.currFrame.Bind(v.ip, mainFn, 0)
	v.stack.SetStackPointer(v.currFrame.NumLocals())
	if v.currFrame.NumParameters() != len(args) {
		return nil, fmt.Errorf("[%s] wrong number of arguments provided: want=%d, got=%d", mainFn.Name(), v.currFrame.NumParameters(), len(args))
	}
	for idx, arg := range args {
		argObj := v.gk.FromInterface(objects.FrameStatic, arg)
		v.stack.SetAbsolute(idx, argObj)
	}

	v.loop()

	if v.err != nil {
		filePos, _ := v.bc.Position(v.currFrame.SourcePos(v.ip - 1))
		err := fmt.Errorf("%w at %s", v.err, filePos)
		for _, frame := range v.frames.Unroll() {
			filePos, _ = v.bc.Position(frame.SourcePos(frame.SavedIP() - 1))
			err = fmt.Errorf("%w at %s", err, filePos)
		}
		return nil, err
	}
	if ret {
		return v.GetReturnValues(), nil
	}
	return nil, nil
}

// loop executes the main instruction loop for the virtual machine, updating the instruction pointer and processing opcodes.
func (v *VM) loop() {
	//log.Printf("starting......")
	var opcode bytecode.OpcodeId
	var decoder *Decoder
	v.counterIterations = 0
	v.counterStart = uint64(time.Now().UnixMilli())
	for {
		v.counterIterations++
		v.ip, opcode = v.currFrame.Fetch(v.ip)
		decoder = v.sequencer[int(opcode)&v.sequencerMask]
		v.ip = decoder.Decode(v.currFrame, v.ip)
		//log.Printf("Executing instruction opcode: %d name: %s ip: %d decoded: %v", opcode, decoder.Name(), v.ip, decoder.decodedOperands[:decoder.fullWidth])
		decoder.Execute()
		if v.shutdown {
			break
		}
	}
}

// callCompiled sets up a new execution frame for a compiled function and manages stack allocation for local variables.
// Callee specifies the compiled function to be executed, and numArgs determines the number of arguments passed.
// It reserves stack space for all local variables and adjusts the instruction pointer accordingly.
func (v *VM) prepareForCall(callee *objects.FuncCompiled, numArgs int) *Frame {
	// 1. Calculate the new basePointer safely, anchoring it to the caller's frame.
	//	The new "floor" begins exactly where the caller's local variable space ends.
	bp := v.currFrame.BasePointer() + v.currFrame.NumLocals()

	v.stack.CopyOffset(bp, numArgs)

	// 3. Advance to the next frame
	v.currFrame = v.frames.Current()
	v.frames.MoveNext()

	// 4. Bind the new frame to function and correct basePointer
	v.currFrame.Bind(v.GetIp(), callee, bp)

	// 5. Set stack pointer to include arguments and space for new locals
	v.stack.SetStackPointer(bp + numArgs + callee.NumLocals())
	v.ReseIp()

	return v.currFrame
}
