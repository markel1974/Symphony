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
	core          *Core
	gk            objects.IGateKeeper
	seq           ISequencer
	op            opcodes.IOpcodes
	bc            *bytecode.Bytecode
	sequencer     []*Decoder
	sequencerMask int
	imports       *Imports
	constants     *Constants
	globals       *Globals
	retValues     bool
	err           error
	shutdown      bool
}

func NewVM(gk objects.IGateKeeper, seq ISequencer, op opcodes.IOpcodes) *VM {
	v := &VM{
		seq:       seq,
		gk:        gk,
		op:        op,
		err:       nil,
		retValues: false,
	}
	v.imports = NewImports(gk, v.SetError)
	v.constants = NewConstants(gk, v.SetError)
	v.globals = NewGlobals(gk, v.SetError)
	v.core = NewCore(gk, v.SetError)
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
	if err := v.seq.Bind(v.core); err != nil {
		return nil, err
	}
	sequencer := v.seq.Executors()
	var err error
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
		for _, code := range codes {
			relocator.Add(code)
		}
		v.bc, err = relocator.Relocate()
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
	if err = v.core.Setup(v.imports, v.constants, v.globals, v.sequencer, v.sequencerMask); err != nil {
		return nil, err
	}
	for _, fn := range v.constants.PreInitFuncs() {
		if _, err = v.exec(v.core, fn, false); err != nil {
			return nil, err
		}
	}
	for _, fn := range v.constants.InitFuncs() {
		if _, err = v.exec(v.core, fn, false); err != nil {
			return nil, err
		}
	}
	return entryPoints, nil
}

// Statistics returns three uint64 values: start, allocated objects, and a counter from the Core instance.
func (v *VM) Statistics() (uint64, uint64, uint64, uint64) {
	return v.core.Statistics()
}

// EnableRetValues sets the flag to enable or disable returning multiple values from the virtual machine's execution.
func (v *VM) EnableRetValues(retValues bool) {
	v.retValues = retValues
}

// Print prints the current state of the virtual machine's stack to the console.'
func (v *VM) Print(writer io.Writer) {
	v.core.Print(writer)
}

// GetReturnValue returns the value from the top of the stack as an interface value.
func (v *VM) GetReturnValue(idx int) interface{} {
	return v.core.GetReturnValue(idx)
}

// GetReturnValues returns the values from the top of the stack as an array of interface values.
func (v *VM) GetReturnValues() []interface{} {
	return v.core.GetReturnValues()
}

// SetError sets the internal error state of the Core and marks it for shutdown.
func (v *VM) SetError(err error) {
	v.err = err
	v.shutdown = true
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
	return v.exec(v.core, mainFn, v.retValues, args...)
}

// Run executes the virtual machine's bytecode, managing the stack, frames, and instruction pointer state.
func (v *VM) exec(core *Core, mainFn *objects.Func, ret bool, args ...interface{}) ([]interface{}, error) {
	v.shutdown = false
	v.err = nil
	if err := core.Initialize(mainFn, args...); err != nil {
		return nil, err
	}
	for v.shutdown == false {
		core.Execute()
	}
	if err := core.Finalize(); err != nil {
		return nil, err
	}
	if v.err != nil {
		filePos, _ := v.bc.Position(core.SourcePos())
		err := fmt.Errorf("%w at %s", v.err, filePos)
		for _, frame := range v.core.FramesUnroll() {
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
