package vm

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

// ISequencer defines an interface to generate a sequence of functions for a given Virtual Machine instance.
type ISequencer interface {
	Create() []IOpExecutor
}

type SequencerData struct {
	execute   func(vm *VM, operands *[]int)
	operands  []func(*Frame, int) (int, int)
	fullWidth int
}

func NewSequencerData(execute func(vm *VM, operands *[]int), operands []int) *SequencerData {
	sd := &SequencerData{
		execute:   execute,
		operands:  make([]func(*Frame, int) (int, int), len(operands)),
		fullWidth: 0,
	}
	for i, width := range operands {
		switch width {
		case 1:
			sd.operands[i] = func(frame *Frame, ip int) (int, int) { return int(frame.Get8(ip)), 1 }
		case 2:
			sd.operands[i] = func(frame *Frame, ip int) (int, int) { return int(frame.Get16(ip)), 2 }
		}
		sd.fullWidth += width
	}
	return sd
}

// VM represents a virtual machine that executes bytecode instructions, handles stack, and manages execution frames.
type VM struct {
	factory     *objects.GateKeeper
	sourceFiles *bytecode.Files
	stack       *Stack
	frames      *Frames
	currFrame   *Frame
	ip          int
	shutdown    bool
	err         error
	sequencer   []*SequencerData
	references  []objects.IObject
	loader      bytecode.ILoader
	constants   *Constants
}

// New initializes and returns a new virtual machine instance configured with the provided components and settings.
func New(factory *objects.GateKeeper, op *bytecode.Opcodes, sequencer ISequencer) *VM {
	v := &VM{
		factory:     factory,
		ip:          resetIp,
		loader:      nil,
		sourceFiles: nil,
		references:  nil,
	}
	v.constants = NewConstants(factory, v.SetError)
	v.stack = NewStack(factory, stackSize, v.SetError)
	v.frames = NewFrames(factory, maxFrames, v.SetError)
	if sequencer == nil {
		sequencer = NewSequencer(op)
	}
	seq := sequencer.Create()
	v.sequencer = make([]*SequencerData, len(seq))
	for i, s := range seq {
		v.sequencer[i] = NewSequencerData(s.Execute, s.Operands())
	}
	v.Reset()
	return v
}

// Shutdown gracefully shuts down the virtual machine by setting its internal state to signify termination.
func (v *VM) Shutdown() {
	v.shutdown = true
}

// Print prints the current state of the virtual machine's stack to the console.'
func (v *VM) Print(writer io.Writer) {
	v.stack.Print(writer)
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

// SetIp sets the virtual machine's instruction pointer to the specified value.
func (v *VM) SetIp(ip int) {
	v.ip = ip
}

// GetIp retrieves the current instruction pointer value from the virtual machine.
func (v *VM) GetIp() int {
	return v.ip
}

// ReseIp resets the instruction pointer of the virtual machine to its initial reset state defined by `resetIp`.
func (v *VM) ReseIp() {
	v.ip = resetIp
}

// Run executes the virtual machine's bytecode, managing the stack, frames, and instruction pointer state.
func (v *VM) Run(loader bytecode.ILoader, bc *bytecode.Bytecode, mainId string, args ...interface{}) error {
	references, err := loader.ResolveSymbols(bc.References())
	if err != nil {
		return err
	}
	var mainFn *objects.FuncCompiled
	constants := make([]objects.IObject, len(bc.Constants()))
	for idx, constant := range bc.Constants() {
		constants[idx] = constant
		switch c := constant.(type) {
		case *objects.Builtin:
			symbol := loader.BuiltinResolve(idx)
			if symbol == nil {
				return fmt.Errorf("builtin symbol not found: %s", c.Name())
			}
			constants[idx] = symbol
		case *objects.FuncCompiled:
			if mainId == c.Name() {
				mainFn = c
			}
		}
	}
	if mainFn == nil {
		return fmt.Errorf("main function not found")
	}
	v.loader = loader
	v.sourceFiles = bc.SourceFiles()
	v.references = references
	v.constants.SetConstants(constants)
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
		err = fmt.Errorf("runtime error %w at %s", v.err, filePos)
		for _, frame := range v.frames.Unroll() {
			filePos, _ = v.sourceFiles.Position(frame.SourcePos(frame.SavedIP() - 1))
			err = fmt.Errorf("%w at %s", err, filePos)
		}
		return err
	}
	return nil
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

// BoundsCheck validates and adjusts slice bounds using provided low and high indices, ensuring they are within valid range.
func (v *VM) BoundsCheck(lowStack objects.IObject, highStack objects.IObject, numElements int64) (int64, int64, error) {
	var lowIdx int64
	if lowStack != v.factory.UndefinedValue() {
		if low, ok := lowStack.(*objects.Int); ok {
			lowIdx = low.Value()
		} else {
			return 0, 0, fmt.Errorf("invalid slice index type: %s", low.TypeName())
		}
	}
	var highIdx int64
	if highStack == v.factory.UndefinedValue() {
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

// IndexAssign assigns a value to a nested structure, using selectors to determine the target location.
// It navigates through the provided selectors and performs an assignment on the target object at the final index.
// Returns an error if any selector is invalid, the object is not indexable, or the assignment fails.
func (v *VM) IndexAssign(frame int, dst objects.IObject, src objects.IObject, selectors []objects.IObject) error {
	numSel := len(selectors)
	for sIdx := numSel - 1; sIdx > 0; sIdx-- {
		next, err := dst.IndexGet(frame, selectors[sIdx])
		if err != nil {
			if objects.Is(err, objects.ErrNotIndexable) {
				return fmt.Errorf("not indexable: %s", dst.TypeName())
			}
			if objects.Is(err, objects.ErrInvalidIndexType) {
				return fmt.Errorf("invalid index type: %s",
					selectors[sIdx].TypeName())
			}
			return err
		}
		dst = next
	}
	if err := dst.IndexSet(selectors[0], src); err != nil {
		if objects.Is(err, objects.ErrNotIndexAssignable) {
			return fmt.Errorf("not index-assignable: %s", dst.TypeName())
		}
		if objects.Is(err, objects.ErrInvalidIndexValueType) {
			return fmt.Errorf("invaid index values type: %s", src.TypeName())
		}
		return err
	}
	return nil
}

// loop executes the main instruction loop for the virtual machine, updating the instruction pointer and processing opcodes.
func (v *VM) loop() {
	cOperands := make([]int, 16)
	cOperandsPtr := &cOperands
	for {
		v.ip++
		inst := v.currFrame.Get8(v.ip)
		opcode := inst & bytecode.OpcodesMask
		data := v.sequencer[opcode]
		if data.fullWidth > 0 {
			v.ip += data.fullWidth
			readOffset := v.ip
			for idx, fn := range data.operands {
				val, width := fn(v.currFrame, readOffset)
				cOperands[idx] = val
				readOffset -= width
			}
		}
		//log.Println("Executing instruction ", opcode, bytecode.OpcodeNames(opcode))
		data.execute(v, cOperandsPtr)
		if v.shutdown {
			break
		}
	}
}
