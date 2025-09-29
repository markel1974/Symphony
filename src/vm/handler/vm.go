package handler

import (
	"fmt"
	"io"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

const version = "0.1"

type VM struct {
	rootCore  *Core
	gk        objects.IGateKeeper
	seq       ISequencer
	op        opcodes.IOpcodes
	bc        *bytecode.Bytecode
	imports   *Imports
	constants *Constants
	globals   *Globals
	retValues bool
	err       error
	shutdown  bool
}

func NewVM(gk objects.IGateKeeper, seq ISequencer, op opcodes.IOpcodes) *VM {
	v := &VM{
		seq:       seq,
		gk:        gk,
		op:        op,
		err:       nil,
		retValues: false,
	}
	v.imports = NewImports(gk, v.shutdownAll)
	v.constants = NewConstants(gk, v.shutdownAll)
	v.globals = NewGlobals(gk, v.shutdownAll)
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
	v.rootCore = NewCore(v.gk, 0, v.shutdownCore)
	if err = v.rootCore.Setup(v.imports, v.constants, v.globals, v.seq); err != nil {
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
	return v.rootCore.Statistics()
}

// EnableRetValues sets the flag to enable or disable returning multiple values from the virtual machine's execution.
func (v *VM) EnableRetValues(retValues bool) {
	v.retValues = retValues
}

// Print prints the current state of the virtual machine's stack to the console.'
func (v *VM) Print(writer io.Writer) {
	v.rootCore.Print(writer)
}

// GetReturnValue returns the value from the top of the stack as an interface value.
func (v *VM) GetReturnValue(idx int) interface{} {
	return v.rootCore.GetReturnValue(idx)
}

// GetReturnValues returns the values from the top of the stack as an array of interface values.
func (v *VM) GetReturnValues() []interface{} {
	return v.rootCore.GetReturnValues()
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
	v.err = nil
	if err := v.rootCore.Initialize(mainFn, args...); err != nil {
		return nil, err
	}
	for v.shutdown == false {
		v.rootCore.Execute()
	}
	v.rootCore.Finalize()
	if v.err != nil {
		filePos, _ := v.bc.Position(v.rootCore.SourcePos())
		err := fmt.Errorf("%w at %s", v.err, filePos)
		for _, frame := range v.rootCore.FramesUnroll() {
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

// ShutdownGlobal sets the virtual machine's error state and initiates a global shutdown sequence.
func (v *VM) shutdownAll(err error) {
	v.err = err
	v.shutdown = true
}

// Shutdown sets the error state and marks the virtual machine as shut down.
func (v *VM) shutdownCore(idx int, err error) {
	if err != nil {
		v.shutdown = true
		v.err = err
		return
	}
	if idx == 0 {
		v.shutdown = true
	}
}
