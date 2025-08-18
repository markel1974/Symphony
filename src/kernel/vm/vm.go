package vm

import (
	"fmt"
	"log"

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

// ISequencer defines an interface to generate a sequence of functions for a given Virtual Machine instance.
type ISequencer interface {
	Create() []IOpExecutor
}

// VM represents a virtual machine that executes bytecode instructions, handles stack, and manages execution frames.
type VM struct {
	sourceFiles *bytecode.Files
	stack       *Stack
	frames      *Frames
	currFrame   *Frame
	ip          int
	shutdown    bool
	err         error
	sequencer   []func(vm *VM)
	references  []objects.IObject
	loader      bytecode.ILoader
	constants   *Constants
	entryPoints map[string]*objects.FunctionCompiled
}

// NewVM initializes and returns a new virtual machine instance configured with the provided components and settings.
func NewVM(loader bytecode.ILoader, sequencer ISequencer, bc *bytecode.Bytecode, maxAllocations int64) (*VM, error) {
	if bc == nil {
		return nil, fmt.Errorf("bytecode is nil")
	}
	if maxAllocations < 1 {
		return nil, fmt.Errorf("max allocations must be greater than 0")
	}
	references, err := loader.ResolveSymbols(bc.References())
	if err != nil {
		return nil, err
	}
	constants := make([]objects.IObject, len(bc.Constants()))
	entryPoints := make(map[string]*objects.FunctionCompiled)
	for idx, constant := range bc.Constants() {
		constants[idx] = constant
		switch c := constant.(type) {
		case *objects.Builtin:
			fn := loader.BuiltinResolve(idx)
			if fn == nil {
				return nil, fmt.Errorf("builtin function not found: %s", c.Name())
			}
			constants[idx] = fn
		case *objects.FunctionCompiled:
			entryPoints[c.Name()] = c
		}
	}
	v := &VM{
		sourceFiles: bc.SourceFiles(),
		ip:          resetIp,
		loader:      loader,
		references:  references,
		entryPoints: entryPoints,
	}
	v.stack = NewStack(stackSize, maxAllocations, v.setError)
	v.frames = NewFrames(maxFrames, v.setError)
	v.constants = NewConstants(constants, v.setError)
	if sequencer == nil {
		sequencer = NewSequencer()
	}
	seq := sequencer.Create()
	v.sequencer = make([]func(vm *VM), len(seq))
	for i, s := range seq {
		v.sequencer[i] = s.Execute
	}

	v.Reset()
	return v, nil
}

// Shutdown gracefully shuts down the virtual machine by setting its internal state to signify termination.
func (v *VM) Shutdown() {
	v.shutdown = true
}

// Reset reinitializes the virtual machine's state, clears the stack and frames, and resets execution-related variables.
func (v *VM) Reset() {
	v.ip = resetIp
	v.stack.Reset()
	v.frames.Reset()
	v.err = nil
	v.shutdown = false
}

// Run executes the virtual machine's bytecode, managing the stack, frames, and instruction pointer state.
func (v *VM) Run(main string) error {
	mainFn, _ := v.entryPoints[main]
	if mainFn == nil {
		return fmt.Errorf("main function not found")
	}
	v.currFrame = v.frames.Head()
	v.currFrame.Bind(v.ip, mainFn, 0)
	v.stack.SetStackPointer(v.currFrame.NumLocals())

	v.run()

	if v.err != nil {
		filePos, _ := v.sourceFiles.Position(v.currFrame.SourcePos(v.ip - 1))
		err := fmt.Errorf("runtime error %w at %s", v.err, filePos)
		for _, frame := range v.frames.Unroll() {
			filePos, _ = v.sourceFiles.Position(frame.SourcePos(frame.StartIP() - 1))
			err = fmt.Errorf("%w at %s", err, filePos)
		}
		return err
	}
	return nil
}

// run executes the main instruction loop for the virtual machine, updating the instruction pointer and processing opcodes.
func (v *VM) run() {
	for {
		v.ip++
		inst := v.currFrame.Get(v.ip)
		opcode := bytecode.Opcode(inst & bytecode.OpcodesMask)
		log.Println("Executing instruction ", opcode, bytecode.OpcodeNames(opcode))
		v.sequencer[opcode](v)
		if v.shutdown {
			break
		}
	}
}

// setError sets the internal error state of the VM and marks it for shutdown.
func (v *VM) setError(err error) {
	v.err = err
	v.shutdown = true
}

// checkBounds validates and adjusts slice bounds using provided low and high indices, ensuring they are within valid range.
func (v *VM) checkBounds(lowStack objects.IObject, highStack objects.IObject, numElements int64) (int64, int64, error) {
	var lowIdx int64
	if lowStack != objects.UndefinedValue {
		if low, ok := lowStack.(*objects.Int); ok {
			lowIdx = low.Value()
		} else {
			return 0, 0, fmt.Errorf("invalid slice index type: %s", low.TypeName())
		}
	}
	var highIdx int64
	if highStack == objects.UndefinedValue {
		highIdx = numElements
	} else if high, ok := highStack.(*objects.Int); ok {
		highIdx = high.Value()
	} else {
		return 0, 0, fmt.Errorf("invalid slice index type: %s", high.TypeName())
	}
	if lowIdx > highIdx {
		return 0, 0, fmt.Errorf("invalid slice index: %d > %d", lowIdx, highIdx)
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
	return lowIdx, highIdx, nil
}
