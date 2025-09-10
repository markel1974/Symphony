package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// registerFunc is a function type that registers an operation with the provided bytecode.Opcodes and returns a core.IOpExecutor.
type registerFunc = func(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error)

// _registerContainer holds a list of functions that register operational executors for bytecode instructions.
var _registerContainer []registerFunc

// SequencerRegister appends a registerFunc to the internal _registerContainer for further use by the sequencer system.
func SequencerRegister(fn registerFunc) {
	_registerContainer = append(_registerContainer, fn)
}

// Sequencer defines a container for managing and initializing opcode executors for a virtual machine's instruction set.
type Sequencer struct {
	op *opcodes.Opcodes
}

// NewSequencer initializes and returns a new Sequencer instance with the provided Opcodes configuration.
func NewSequencer(op *opcodes.Opcodes) core.ISequencer {
	return &Sequencer{
		op: op,
	}
}

// Create initializes and returns a slice of IOpExecutor with default OpUnknown executors, supplemented by static mappings.
func (ds *Sequencer) Create(vm *core.VM) ([]core.IOpExecutor, error) {
	v, err := ds.createRegistered(vm)
	return v, err
}

// createRegistered constructs and registers custom IOpExecutors using functions from _registerContainer and updates the sequence.
func (ds *Sequencer) createRegistered(vmIn *core.VM) ([]core.IOpExecutor, error) {
	fullAccess := core.IVMFullAccess(vmIn)
	container := make([]core.IOpExecutor, ds.op.Len())
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
		opId := seq.OpcodeId()
		if opId < 0 || int(opId) >= len(container) {
			return nil, fmt.Errorf("opcode %d out of range", opId)
		}
		if registered := container[opId]; registered.OpcodeId() != opcodes.OpUnknown {
			return nil, fmt.Errorf("opcode %d already registered: %s", opId, registered.Name())
		}
		container[opId] = seq
	}
	return container, nil
}

// createStatic initializes and assigns specific op executors to the sequencer's container in a pre-defined sequence.
func (ds *Sequencer) createStatic(vmIn *core.VM) ([]core.IOpExecutor, error) {
	fullAccess := core.IVMFullAccess(vmIn)
	container := make([]core.IOpExecutor, ds.op.Len())
	for idx := range container {
		var err error
		if container[idx], err = NewOpUnknown(fullAccess, ds.op); err != nil {
			return nil, err
		}
	}

	var z []registerFunc

	z = append(z, NewOpConstant)
	z = append(z, NewOpNull)
	z = append(z, NewOpLogical)
	z = append(z, NewOpArithmetic)
	z = append(z, NewOpImport)
	//z = append(z, NewOpFuncInternal)
	z = append(z, NewOpPop)
	z = append(z, NewOpTrue)
	z = append(z, NewOpFalse)
	z = append(z, NewOpUnaryNot)
	z = append(z, NewOpUnaryBitwiseComplement)
	z = append(z, NewOpUnarySub)
	z = append(z, NewOpJumpFalsy)
	z = append(z, NewOpJumpTruthy)
	z = append(z, NewOpJumpAnd)
	z = append(z, NewOpJumpOr)
	z = append(z, NewOpJump)
	z = append(z, NewOpJumpIndirect)
	z = append(z, NewOpGlobalSet)
	z = append(z, NewOpGlobalIndex)
	z = append(z, NewOpGlobalGet)
	z = append(z, NewOpGlobalCopy)
	z = append(z, NewOpCreateArray)
	z = append(z, NewOpCreateMap)
	z = append(z, NewOpCreateStruct)
	z = append(z, NewOpCreateError)
	z = append(z, NewOpIndexGet)
	z = append(z, NewOpIndexSet)
	z = append(z, NewOpIndexSlice)
	z = append(z, NewOpCall)
	z = append(z, NewOpReturn)
	z = append(z, NewOpLocalDefine)
	z = append(z, NewOpLocalSet)
	z = append(z, NewOpLocalIndex)
	z = append(z, NewOpLocalGet)
	z = append(z, NewOpCreateClosure)
	z = append(z, NewOpFreeGetPtr)
	z = append(z, NewOpFreeGet)
	z = append(z, NewOpFreeSet)
	z = append(z, NewOpLocalPtrGet)
	z = append(z, NewOpIteratorInit)
	z = append(z, NewOpIteratorNext)
	z = append(z, NewOpIteratorKey)
	z = append(z, NewOpIteratorValue)
	z = append(z, NewOpIntLogical)
	z = append(z, NewOpIntArithmetic)
	z = append(z, NewOpDerefGet)
	z = append(z, NewOpNoOp)
	z = append(z, NewOpSuspend)

	for _, fn := range z {
		seq, err := fn(fullAccess, ds.op)
		if err != nil {
			return nil, err
		}
		opId := seq.OpcodeId()
		if opId < 0 || int(opId) >= len(container) {
			return nil, fmt.Errorf("opcode %d out of range", opId)
		}
		container[opId] = seq
	}
	return container, nil
}

// facadeForOpcode returns a facade for the provided opcodeId.
func (ds *Sequencer) facadeForOpcode(opcodeId opcodes.OpcodeId, vm *core.VM) interface{} {
	switch opcodeId {
	// Category: Stack & Constants (Read-Only)
	case opcodes.OpConstant, opcodes.OpGlobalGet, opcodes.OpImport, opcodes.OpFreeGet, opcodes.OpLocalGet:
		return core.IVMReadOnly(vm)
	// Category: Control Flow
	case opcodes.OpJump, opcodes.OpJumpFalsy, opcodes.OpJumpAnd, opcodes.OpJumpOr:
		return core.IVMControlFlow(vm)
	// Category: Simple Stack
	case opcodes.OpPop, opcodes.OpTrue, opcodes.OpFalse, opcodes.OpNull, opcodes.OpUnaryNot, opcodes.OpUnarySub:
		// These instructions only manipulate the top of the stack.
		return core.IVMStackOnly(vm)
	// Category: Full Access (use with caution)
	case opcodes.OpCall, opcodes.OpReturn, opcodes.OpCallMethod, opcodes.OpCreateClosure:
		// These complex operations need wider access.
		return core.IVMFullAccess(vm)
	// Default Case
	default:
		// Basic read/write access.
		return core.IVMReadWrite(vm)
	}
}
