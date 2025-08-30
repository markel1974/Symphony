package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

// registerFunc is a function type that registers an operation with the provided bytecode.Opcodes and returns a core.IOpExecutor.
type registerFunc = func(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error)

// _registerContainer holds a list of functions that register operational executors for bytecode instructions.
var _registerContainer []registerFunc

// SequencerRegister appends a registerFunc to the internal _registerContainer for further use by the sequencer system.
func SequencerRegister(fn registerFunc) {
	_registerContainer = append(_registerContainer, fn)
}

// Sequencer defines a container for managing and initializing opcode executors for a virtual machine's instruction set.
type Sequencer struct {
	op *bytecode.Opcodes
}

// NewSequencer initializes and returns a new Sequencer instance with the provided Opcodes configuration.
func NewSequencer(op *bytecode.Opcodes) core.ISequencer {
	return &Sequencer{
		op: op,
	}
}

// Create initializes and returns a slice of IOpExecutor with default OpUnknown executors, supplemented by static mappings.
func (ds *Sequencer) Create(vm *core.VM) []core.IOpExecutor {
	v, _ := ds.createRegistered(vm)
	return v
}

// createRegistered constructs and registers custom IOpExecutors using functions from _registerContainer and updates the sequence.
func (ds *Sequencer) createRegistered(vmIn *core.VM) ([]core.IOpExecutor, error) {
	fullAccess := core.IVMFullAccess(vmIn)
	container := make([]core.IOpExecutor, bytecode.OpcodesLen)
	for idx := range container {
		var err error
		if container[idx], err = NewOpUnknown(fullAccess, ds.op); err != nil {
			return nil, err
		}
	}
	for _, fn := range _registerContainer {
		seq, err := fn(fullAccess, ds.op)
		if err != nil {
			return nil, err
		}
		data := container[seq.OpcodeId()&bytecode.OpcodesMask]
		if data.OpcodeId() != bytecode.OpUnknown {
			return nil, fmt.Errorf("opcode %d already registered: %s", seq.OpcodeId(), data.Name())
		}
		container[seq.OpcodeId()&bytecode.OpcodesMask] = seq
	}
	return container, nil
}

// createStatic initializes and assigns specific op executors to the sequencer's container in a pre-defined sequence.
func (ds *Sequencer) createStatic(vmIn *core.VM) ([]core.IOpExecutor, error) {
	fullAccess := core.IVMFullAccess(vmIn)
	container := make([]core.IOpExecutor, bytecode.OpcodesLen)
	for idx := range container {
		var err error
		if container[idx], err = NewOpUnknown(fullAccess, ds.op); err != nil {
			return nil, err
		}
	}

	var z []registerFunc

	z = append(z, NewOpConstant)
	z = append(z, NewOpNull)
	z = append(z, NewOpBinary)
	z = append(z, NewOpReferences)
	z = append(z, NewOpEqual)
	z = append(z, NewOpNotEqual)
	z = append(z, NewOpPop)
	z = append(z, NewOpTrue)
	z = append(z, NewOpFalse)
	z = append(z, NewOpNotLogical)
	z = append(z, NewOpBitwiseComplement)
	z = append(z, NewOpMinus)
	z = append(z, NewOpJumpFalsy)
	z = append(z, NewOpJumpAnd)
	z = append(z, NewOpJumpOr)
	z = append(z, NewOpJump)
	z = append(z, NewOpGlobalSet)
	z = append(z, NewOpGlobalSelSet)
	z = append(z, NewOpGetGlobal)
	z = append(z, NewOpArray)
	z = append(z, NewOpMap)
	z = append(z, NewOpStruct)
	z = append(z, NewOpError)
	z = append(z, NewOpImmutable)
	z = append(z, NewOpIndexGet)
	z = append(z, NewOpIndexSet)
	z = append(z, NewOpIndexSlice)
	z = append(z, NewOpCall)
	z = append(z, NewOpReturn)
	z = append(z, NewOpLocalDefine)
	z = append(z, NewOpLocalSet)
	z = append(z, NewOpLocalSelSet)
	z = append(z, NewOpLocalGet)
	z = append(z, NewOpClosure)
	z = append(z, NewOpFreeGetPtr)
	z = append(z, NewOpFreeGet)
	z = append(z, NewOpFreeSet)
	z = append(z, NewOpLocalPtrGet)
	z = append(z, NewOpFreeSelSet)
	z = append(z, NewOpIteratorInit)
	z = append(z, NewOpIteratorNext)
	z = append(z, NewOpIteratorKey)
	z = append(z, NewOpIteratorValue)
	z = append(z, NewOpIntOp)
	z = append(z, NewOpDeref)
	z = append(z, NewOpNoOp)
	z = append(z, NewOpSuspend)

	for _, fn := range z {
		seq, err := fn(fullAccess, ds.op)
		if err != nil {
			return nil, err
		}
		container[seq.OpcodeId()&bytecode.OpcodesMask] = seq
	}
	return container, nil
}

// facadeForOpcode returns a facade for the provided opcodeId.
func (ds *Sequencer) facadeForOpcode(opcodeId bytecode.OpcodeId, vm *core.VM) interface{} {
	switch opcodeId {
	// Category: Stack & Constants (Read-Only)
	case bytecode.OpConstant, bytecode.OpGlobalGet, bytecode.OpReferences, bytecode.OpFreeGet, bytecode.OpLocalGet:
		return core.IVMReadOnly(vm)
	// Category: Control Flow
	case bytecode.OpJump, bytecode.OpJumpFalsy, bytecode.OpJumpAnd, bytecode.OpJumpOr:
		return core.IVMControlFlow(vm)
	// Category: Simple Stack
	case bytecode.OpPop, bytecode.OpTrue, bytecode.OpFalse, bytecode.OpNull, bytecode.OpNotLogical, bytecode.OpMinus:
		// These instructions only manipulate the top of the stack.
		return core.IVMStackOnly(vm)
	// Category: Full Access (use with caution)
	case bytecode.OpCall, bytecode.OpReturn, bytecode.OpCallMethod, bytecode.OpClosure:
		// These complex operations need wider access.
		return core.IVMFullAccess(vm)
	// Default Case
	default:
		// Basic read/write access.
		return core.IVMReadWrite(vm)
	}
}
