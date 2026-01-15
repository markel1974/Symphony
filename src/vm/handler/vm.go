package handler

import (
	"fmt"
	"io"
	"time"

	"github.com/markel1974/symphony/src/vm/bytecode"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

const version = "0.1"

// stackSize specifies the size limit of the stack for function execution.
// maxFrames indicates the maximum number of call frames allowed.
const (
	stackSize = 1024
	maxFrames = 512
	maxCores  = 8
)

type Error struct {
	error
	id uint
}

func NewError(err error, id uint) *Error {
	return &Error{
		error: err,
		id:    id,
	}
}

type VM struct {
	gk                objects.IGateKeeper
	seq               ISequencer
	op                opcodes.IOpcodes
	bc                *bytecode.Bytecode
	imports           *Imports
	constants         *Constants
	globals           *Globals
	retValues         bool
	error             *Error
	shutdown          bool
	coreMainId        int
	cores             []*Core
	coresRunning      []*Core
	coresFree         []*Core
	counterIterations uint64
	counterStart      uint64
	maxCores          int
	maxFrames         int
	stackSize         int
}

func NewVM(gk objects.IGateKeeper, seq ISequencer, op opcodes.IOpcodes) *VM {
	v := &VM{
		seq:       seq,
		gk:        gk,
		op:        op,
		error:     nil,
		retValues: false,
		maxCores:  maxCores,
		maxFrames: maxFrames,
		stackSize: stackSize,
	}
	v.imports = NewImports(gk)
	v.constants = NewConstants(gk)
	v.globals = NewGlobals(gk)
	return v
}

// SetMaxCores sets the maximum number of cores that the virtual machine can utilize.
func (v *VM) SetMaxCores(c int) {
	v.maxCores = c
}

// SetMaxFrames sets the maximum number of stack frames allowed in the virtual machine.
func (v *VM) SetMaxFrames(f int) {
	v.maxFrames = f
}

// SetStackSize sets the maximum size of the stack for the virtual machine.
func (v *VM) SetStackSize(s int) {
	v.stackSize = s
}

// Version returns the version of the virtual machine.
func (v *VM) Version() string {
	return version
}

// Setup initializes the virtual machine with the provided bytecode and loader components.
func (v *VM) Setup(loader bytecode.ILoader, codes ...*bytecode.Bytecode) (map[string]uint, error) {
	//var test IVMFullAccess = v.core
	//fmt.Println("Core.Setup", test)
	var err error
	switch len(codes) {
	case 1:
		v.bc = codes[0]
	default:
		relocator := bytecode.NewRelocator(v.gk, loader, v.op, bytecode.PreInitFunction, bytecode.InitFunction)
		for _, code := range codes {
			relocator.Add(code)
		}
		if v.bc, err = relocator.Relocate(); err != nil {
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
	v.cores = make([]*Core, v.maxCores)
	for idx := range v.cores {
		core := NewCore(v.gk, v.maxFrames, v.stackSize, uint(idx), v.coreShutdown, v.coreCreate)
		if err = core.Setup(v.imports, v.constants, v.globals, v.seq); err != nil {
			return nil, err
		}
		v.cores[idx] = core
	}
	for _, preInit := range v.constants.PreInitFuncs() {
		if _, err = v.exec(preInit, false); err != nil {
			return nil, err
		}
	}
	for _, init := range v.constants.InitFuncs() {
		if _, err = v.exec(init, false); err != nil {
			return nil, err
		}
	}
	return entryPoints, nil
}

// Statistics returns three uint64 values: start, allocated objects, and a counter from the Core instance.
func (v *VM) Statistics() (uint64, uint64, uint64, uint64) {
	framesMax := uint64(0)
	for _, c := range v.cores {
		framesMax += c.FramesMax()
	}
	return v.counterStart, v.gk.AllocatedObjects(), v.counterIterations, framesMax
}

// EnableRetValues sets the flag to enable or disable returning multiple values from the virtual machine's execution.
func (v *VM) EnableRetValues(retValues bool) {
	v.retValues = retValues
}

// Print prints the current state of the virtual machine's stack to the console.'
func (v *VM) Print(writer io.Writer) {
	for _, core := range v.cores {
		core.Print(writer)
	}
}

// GetReturnValue returns the value from the top of the stack as an interface value.
func (v *VM) GetReturnValue(idx int) interface{} {
	if v.coreMainId < 0 || v.coreMainId > len(v.cores) {
		return nil
	}
	return v.cores[v.coreMainId].GetReturnValue(idx)
}

// GetReturnValues returns the values from the top of the stack as an array of interface values.
func (v *VM) GetReturnValues() []interface{} {
	if v.coreMainId < 0 || v.coreMainId > len(v.cores) {
		return nil
	}
	return v.cores[v.coreMainId].GetReturnValues()
}

// Run executes the main function identified by mainId with the provided arguments in the virtual machine context.
func (v *VM) Run(mainId uint, args ...interface{}) ([]interface{}, error) {
	entryFn, err := v.constants.RetrieveFunc(mainId)
	if err != nil {
		return nil, err
	}
	return v.exec(entryFn, v.retValues, args...)
}

// Run executes the virtual machine's bytecode, managing the stack, frames, and instruction pointer state.
func (v *VM) exec(mainFn *objects.Func, ret bool, args ...interface{}) ([]interface{}, error) {
	v.shutdown = false
	v.error = nil
	v.counterIterations = 0
	v.counterStart = uint64(time.Now().UnixMilli())

	v.gk.Reset()

	v.coreMainId = -1
	v.coresRunning = make([]*Core, 0, len(v.cores))
	v.coresFree = make([]*Core, len(v.cores))
	for idx, core := range v.cores {
		core.Reset()
		v.coresFree[idx] = core
	}

	v.coreCreate(0, mainFn, v.gk.FromArrayInterfaces(objects.FrameStatic, args))

	for v.shutdown == false {
		for _, core := range v.coresRunning {
			v.counterIterations++
			core.Execute()
		}
	}

	for _, core := range v.coresRunning {
		core.Finalize()
	}

	if v.error != nil {
		return nil, v.error
	}

	if ret {
		return v.GetReturnValues(), nil
	}

	return nil, nil
}

// coreCreate allocates a new core for execution, initializes it with the given function and arguments, and handles errors.
func (v *VM) coreCreate(_ uint, callee *objects.Func, args []objects.IObject) {
	if len(v.coresFree) == 0 {
		v.shutdown = true
		v.error = NewError(fmt.Errorf("no cores available"), 0)
		return
	}
	coreId := len(v.coresFree) - 1
	if v.coreMainId < 0 {
		v.coreMainId = coreId
	}
	core := v.coresFree[coreId]
	v.coresFree = v.coresFree[:coreId]

	runningIndex := len(v.coresRunning)
	v.coresRunning = append(v.coresRunning, core)
	if err := core.Initialize(runningIndex, callee, args); err != nil {
		v.coreShutdown(core.Id(), err)
	}
}

// coreShutdown handles shutting down a core by finalizing its state, removing it from the running list, and marking it as free.
// If an error is provided, the VM is set to shutdown mode with the error stored for diagnostics.
// If the ID is 0, the VM is also shut down to prevent invalid operations.
// The method ensures proper reallocation of resources while maintaining consistency in the running cores list.
func (v *VM) coreShutdown(id uint, err error) {
	core := v.cores[id]
	if err != nil {
		filePos, _ := v.bc.Position(core.SourcePos())
		err = fmt.Errorf("%w at %s", err, filePos)
		for _, frame := range core.FramesUnroll() {
			filePos, _ = v.bc.Position(frame.SourcePos(int(frame.SavedIP()) - 1))
			err = fmt.Errorf("%w at %s", err, filePos)
		}
		v.shutdown = true
		v.error = NewError(err, id)
		return
	}
	if int(id) == v.coreMainId {
		v.shutdown = true
		return
	}

	runningIndex := core.Finalize()
	if runningIndex < 0 {
		return
	}
	lastIndex := len(v.coresRunning) - 1
	lastCore := v.coresRunning[lastIndex]
	v.coresRunning[runningIndex] = lastCore

	lastCore.Update(runningIndex)
	v.coresRunning = v.coresRunning[:lastIndex]

	v.coresFree = append(v.coresFree, core)
}
