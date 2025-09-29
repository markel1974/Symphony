package handler

import (
	"fmt"
	"io"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

const version = "0.1"

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
	gk        objects.IGateKeeper
	seq       ISequencer
	op        opcodes.IOpcodes
	bc        *bytecode.Bytecode
	imports   *Imports
	constants *Constants
	globals   *Globals
	retValues bool
	error     *Error
	shutdown  bool
	mainCore  *Core
	cores     []*Core
}

func NewVM(gk objects.IGateKeeper, seq ISequencer, op opcodes.IOpcodes) *VM {
	v := &VM{
		seq:       seq,
		gk:        gk,
		op:        op,
		error:     nil,
		retValues: false,
	}
	v.imports = NewImports(gk)
	v.constants = NewConstants(gk)
	v.globals = NewGlobals(gk)
	return v
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
	v.mainCore = NewCore(v.gk, 0, v.coreShutdown, v.coreCreate)
	if err = v.mainCore.Setup(v.imports, v.constants, v.globals, v.seq); err != nil {
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

// Statistics returns three uint64 values: start, allocated objects, and a counter from the Core instance.
func (v *VM) Statistics() (uint64, uint64, uint64, uint64) {
	return v.mainCore.Statistics()
}

// EnableRetValues sets the flag to enable or disable returning multiple values from the virtual machine's execution.
func (v *VM) EnableRetValues(retValues bool) {
	v.retValues = retValues
}

// Print prints the current state of the virtual machine's stack to the console.'
func (v *VM) Print(writer io.Writer) {
	v.mainCore.Print(writer)
}

// GetReturnValue returns the value from the top of the stack as an interface value.
func (v *VM) GetReturnValue(idx int) interface{} {
	return v.mainCore.GetReturnValue(idx)
}

// GetReturnValues returns the values from the top of the stack as an array of interface values.
func (v *VM) GetReturnValues() []interface{} {
	return v.mainCore.GetReturnValues()
}

// Run executes the main function identified by mainId with the provided arguments in the virtual machine context.
func (v *VM) Run(mainId uint, args ...interface{}) ([]interface{}, error) {
	obj, err := v.constants.Retrieve(mainId)
	if err != nil {
		return nil, err
	}
	mainFn, ok := obj.(*objects.Func)
	if !ok {
		return nil, fmt.Errorf("entry point not found: %d", mainId)
	}
	return v.exec(mainFn, v.retValues, args...)
}

// Run executes the virtual machine's bytecode, managing the stack, frames, and instruction pointer state.
func (v *VM) exec(mainFn *objects.Func, ret bool, args ...interface{}) ([]interface{}, error) {
	v.shutdown = false
	v.error = nil
	v.cores = []*Core{v.mainCore}
	argsObj := v.gk.FromArrayInterfaces(objects.FrameStatic, args)
	if err := v.mainCore.Initialize(mainFn, argsObj); err != nil {
		return nil, err
	}
	for v.shutdown == false {
		for _, core := range v.cores {
			core.Execute()
		}
	}
	for _, core := range v.cores {
		core.Finalize()
	}
	if v.error != nil {
		if v.error.id >= uint(len(v.cores)) {
			return nil, fmt.Errorf("invalid core index: %d", v.error.id)
		}
		core := v.cores[v.error.id]
		filePos, _ := v.bc.Position(core.SourcePos())
		err := fmt.Errorf("%w at %s", v.error, filePos)
		for _, frame := range core.FramesUnroll() {
			filePos, _ = v.bc.Position(frame.SourcePos(int(frame.SavedIP()) - 1))
			err = fmt.Errorf("%w at %s", err, filePos)
		}
		return nil, err
	}
	if ret {
		return v.GetReturnValues(), nil
	}
	return nil, nil
}

// coreCreate initializes and adds a new Core instance to the VM, setting it up and handling errors during initialization.
func (v *VM) coreCreate(_ uint, callee *objects.Func, args []objects.IObject) {
	core := NewCore(v.gk, uint(len(v.cores)), v.coreShutdown, v.coreCreate)
	v.cores = append(v.cores, core)
	if err := core.Setup(v.imports, v.constants, v.globals, v.seq); err != nil {
		v.coreShutdown(core.Id(), err)
		return
	}
	if err := core.Initialize(callee, args); err != nil {
		v.coreShutdown(core.Id(), err)
		return
	}
}

// Shutdown sets the error state and marks the virtual machine as shut down.
func (v *VM) coreShutdown(id uint, err error) {
	if err != nil {
		v.shutdown = true
		v.error = NewError(err, id)
		return
	}
	if id == 0 {
		v.shutdown = true
		return
	}
	if id >= uint(len(v.cores)) {
		v.shutdown = true
		v.error = NewError(fmt.Errorf("invalid core index: %d", id), id)
		return
	}
	core := v.cores[id]
	core.Finalize()
	v.cores = append(v.cores[:id], v.cores[id+1:]...)
}
