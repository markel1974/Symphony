package core

import (
	"fmt"
	"io"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// resetIp is the instruction pointer value used to reset the VM's instruction pointer to the beginning of the main function.'
const (
	resetIp = -1
)

// globalsSize defines the maximum size of the global variables space.
// stackSize specifies the size limit of the stack for function execution.
// maxFrames indicates the maximum number of call frames allowed.
const (
	stackSize = 2048
	maxFrames = 1024
)

// VM represents a virtual machine that executes bytecode instructions, handles stack, and manages execution frames.
type VM struct {
	factory     objects.IGateKeeper
	sourceFiles *bytecode.Files
	stack       *Stack
	frames      *Frames
	currFrame   *Frame
	ip          int
	shutdown    bool
	err         error
	sequencer   []*Decoder
	references  *References
	constants   *Constants
	globals     *Globals
	entryPoints map[string]*objects.FuncCompiled
}

// New initializes and returns a new virtual machine instance configured with the provided components and settings.
func New(factory objects.IGateKeeper, sequencer ISequencer) *VM {
	v := &VM{
		factory:     factory,
		ip:          resetIp,
		sourceFiles: nil,
		references:  nil,
		entryPoints: make(map[string]*objects.FuncCompiled),
	}
	v.constants = NewConstants(factory, v.SetError)
	v.references = NewReferences(factory, v.SetError)
	v.globals = NewGlobals(factory, v.SetError)
	v.stack = NewStack(factory, stackSize, v.SetError)
	v.frames = NewFrames(factory, maxFrames, v.SetError)
	seq := sequencer.Create()
	v.sequencer = make([]*Decoder, len(seq))
	for i, s := range seq {
		v.sequencer[i] = NewDecoder(s.Execute, s.Operands())
	}
	return v
}

// Setup initializes the virtual machine with the provided bytecode and loader components.
func (v *VM) Setup(loader bytecode.ILoader, bc *bytecode.Bytecode) error {
	if err := v.references.Setup(loader, bc.References()); err != nil {
		return err
	}
	if err := v.constants.Setup(bc.Constants()); err != nil {
		return err
	}
	if err := v.globals.Setup(bc.Global()); err != nil {
		return err
	}
	v.sourceFiles = bc.SourceFiles()
	for _, global := range bc.Global() {
		switch c := global.(type) {
		case *objects.FuncCompiled:
			v.entryPoints[c.Name()] = c
		}
	}
	return nil
}

// Reset reinitializes the virtual machine's state, clears the stack and frames, and resets execution-related variables.
func (v *VM) Reset() {
	v.ip = resetIp
	v.factory.Reset()
	v.stack.Reset()
	v.frames.Reset()
	v.err = nil
	v.shutdown = false
}

// Run executes the virtual machine's bytecode, managing the stack, frames, and instruction pointer state.
func (v *VM) Run(mainId string, args ...interface{}) error {
	v.Reset()

	mainFn, ok := v.entryPoints[mainId]
	if !ok {
		return fmt.Errorf("entry point not found: %s", mainId)
	}

	v.currFrame = v.frames.Head()
	v.currFrame.Bind(v.ip, mainFn, 0)
	v.stack.SetStackPointer(v.currFrame.NumLocals())
	if v.currFrame.NumParameters() != len(args) {
		return fmt.Errorf("[%s] wrong number of arguments provided: want=%d, got=%d", mainId, v.currFrame.NumParameters(), len(args))
	}

	for idx, arg := range args {
		argObj := v.factory.FromInterface(objects.FrameStatic, arg)
		v.stack.SetAbsolute(idx, argObj)
	}

	v.loop()

	if v.err != nil {
		filePos, _ := v.sourceFiles.Position(v.currFrame.SourcePos(v.ip - 1))
		err := fmt.Errorf("runtime error %w at %s", v.err, filePos)
		for _, frame := range v.frames.Unroll() {
			filePos, _ = v.sourceFiles.Position(frame.SourcePos(frame.SavedIP() - 1))
			err = fmt.Errorf("%w at %s", err, filePos)
		}
		return err
	}
	return nil
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

// References return a pointer to the References object associated with the VM instance.
func (v *VM) References() *References {
	return v.references
}

// SetIp sets the virtual machine's instruction pointer to the specified value.
func (v *VM) SetIp(ip int) {
	v.ip = ip
}

// GetIp retrieves the current instruction pointer value from the virtual machine.
func (v *VM) GetIp() int {
	return v.ip
}

// FunctionLibraryCall invokes a callable object with the given arguments and handles stack cleanup and error management.
func (v *VM) FunctionLibraryCall(value objects.IObject, args []objects.IObject, numArgs int) {
	ret, err := value.Call(v.FrameID(), args...)
	// Cleans the stack from the function and its arguments
	v.stack.DecrementCount(numArgs + 1)
	if err != nil {
		if objects.Is(err, objects.ErrWrongNumArguments) {
			v.SetError(fmt.Errorf("wrong number of arguments in call to '%s'", value.TypeName()))
		} else {
			v.SetError(err)
		}
		return
	}
	if ret == nil {
		v.stack.Push(v.factory.UndefinedValue())
	} else {
		v.stack.Push(ret)
	}
}

// FunctionCompiledCall sets up a new execution frame for a compiled function and manages stack allocation for local variables.
// callee specifies the compiled function to be executed, and numArgs determines the number of arguments passed.
// It reserves stack space for all local variables and adjusts the instruction pointer accordingly.
func (v *VM) FunctionCompiledCall(callee *objects.FuncCompiled, numArgs int) {
	// Frame setup
	v.currFrame = v.frames.Get()
	v.frames.Next()
	bp := v.stack.StackPointer() - numArgs
	v.currFrame.Bind(v.GetIp(), callee, bp)
	// Reserve space for *all* local variables of the new function
	// by simply advancing the stack pointer.
	// This ensures that space for temporary calculations starts *after*
	// the space reserved for local variables, avoiding collisions.
	v.stack.SetStackPointer(v.stack.StackPointer() + callee.NumLocals())
	v.ReseIp()
}

// FunctionCompiledReturn handles the return operation by unwinding the current call frame, restoring the previous frame, and managing the stack.
func (v *VM) FunctionCompiledReturn(returnValues []objects.IObject) {
	shutdown := false
	prevIp := v.currFrame.SavedIP()
	leavingFrameBasePointer := v.BasePointer()
	v.stack.ReleaseObjects(leavingFrameBasePointer, v.stack.StackPointer())
	if v.frames.Index() > 1 {
		v.frames.Previous()
		v.currFrame = v.frames.GetPrev()
		v.SetIp(prevIp)
	} else {
		shutdown = true
	}
	v.stack.SetStackPointer(leavingFrameBasePointer)
	// push return values onto the new stack (caller's stack).
	if lRet := len(returnValues); lRet > 0 {
		// iterate over the slice in reverse to restore the original order.
		for i := lRet - 1; i >= 0; i-- {
			v.stack.Push(returnValues[i])
		}
	} else {
		v.stack.Push(v.factory.UndefinedValue())
	}
	if shutdown {
		v.Shutdown()
	}
}

// FreeVarsIndex retrieves the pointer to a free variable at the specified index from the current frame.
func (v *VM) FreeVarsIndex(idx int) *objects.ObjectPointer {
	return v.currFrame.FreeVarsIndex(idx)
}

// FrameID returns the ID of the current execution frame within the virtual machine.
func (v *VM) FrameID() int {
	return v.currFrame.ID()
}

// BasePointer returns the base pointer of the current frame in the virtual machine.
func (v *VM) BasePointer() int {
	return v.currFrame.BasePointer()
}

// StackPointer returns the current position of the stack pointer within the virtual machine's stack.
func (v *VM) StackPointer() int {
	return v.stack.StackPointer()
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
	return v.factory.ToInterface(obj)
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

// loop executes the main instruction loop for the virtual machine, updating the instruction pointer and processing opcodes.
func (v *VM) loop() {
	var opcode byte
	var decoder *Decoder
	for {
		v.ip++
		opcode = v.currFrame.Get8(v.ip)
		decoder = v.sequencer[opcode]
		v.ip = decoder.Decode(v.currFrame, v.ip)
		//log.Println("Executing instruction ", opcode, bytecode.OpcodeNames(opcode))
		decoder.Execute(v)
		if v.shutdown {
			break
		}
	}
}
