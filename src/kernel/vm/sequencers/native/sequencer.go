package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

// registerFunc is a function type that registers an operation with the provided bytecode.Opcodes and returns a core.IOpExecutor.
type registerFunc = func(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor

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
func (ds *Sequencer) createRegistered(vm *core.VM) ([]core.IOpExecutor, error) {
	container := make([]core.IOpExecutor, bytecode.OpcodesLen)
	for idx := range container {
		container[idx] = NewOpUnknown(vm, ds.op)
	}
	for _, fn := range _registerContainer {
		seq := fn(vm, ds.op)
		data := container[seq.OpcodeId()&bytecode.OpcodesMask]
		if data.OpcodeId() != bytecode.OpUnknown {
			return nil, fmt.Errorf("opcode %d already registered: %s", seq.OpcodeId(), data.Name())
		}
		container[seq.OpcodeId()&bytecode.OpcodesMask] = seq
	}
	return container, nil
}

// createStatic initializes and assigns specific op executors to the sequencer's container in a pre-defined sequence.
func (ds *Sequencer) createStatic(vm *core.VM) []core.IOpExecutor {
	container := make([]core.IOpExecutor, bytecode.OpcodesLen)
	for idx := range container {
		container[idx] = NewOpUnknown(vm, ds.op)
	}
	ds.setSequence(container[:], NewOpConstant(vm, ds.op))
	ds.setSequence(container[:], NewOpNull(vm, ds.op))
	ds.setSequence(container[:], NewOpBinary(vm, ds.op))
	ds.setSequence(container[:], NewOpReferences(vm, ds.op))
	ds.setSequence(container[:], NewOpEqual(vm, ds.op))
	ds.setSequence(container[:], NewOpNotEqual(vm, ds.op))
	ds.setSequence(container[:], NewOpPop(vm, ds.op))
	ds.setSequence(container[:], NewOpTrue(vm, ds.op))
	ds.setSequence(container[:], NewOpFalse(vm, ds.op))
	ds.setSequence(container[:], NewOpNotLogical(vm, ds.op))
	ds.setSequence(container[:], NewOpBitwiseComplement(vm, ds.op))
	ds.setSequence(container[:], NewOpMinus(vm, ds.op))
	ds.setSequence(container[:], NewOpJumpFalsy(vm, ds.op))
	ds.setSequence(container[:], NewOpJumpAnd(vm, ds.op))
	ds.setSequence(container[:], NewOpJumpOr(vm, ds.op))
	ds.setSequence(container[:], NewOpJump(vm, ds.op))
	ds.setSequence(container[:], NewOpGlobalSet(vm, ds.op))
	ds.setSequence(container[:], NewOpGlobalSelSet(vm, ds.op))
	ds.setSequence(container[:], NewOpGetGlobal(vm, ds.op))
	ds.setSequence(container[:], NewOpArray(vm, ds.op))
	ds.setSequence(container[:], NewOpMap(vm, ds.op))
	ds.setSequence(container[:], NewOpStruct(vm, ds.op))
	ds.setSequence(container[:], NewOpError(vm, ds.op))
	ds.setSequence(container[:], NewOpImmutable(vm, ds.op))
	ds.setSequence(container[:], NewOpIndex(vm, ds.op))
	ds.setSequence(container[:], NewOpIndexSlice(vm, ds.op))
	ds.setSequence(container[:], NewOpCall(vm, ds.op))
	ds.setSequence(container[:], NewOpReturn(vm, ds.op))
	ds.setSequence(container[:], NewOpLocalDefine(vm, ds.op))
	ds.setSequence(container[:], NewOpLocalSet(vm, ds.op))
	ds.setSequence(container[:], NewOpLocalSelSet(vm, ds.op))
	ds.setSequence(container[:], NewOpLocalGet(vm, ds.op))
	ds.setSequence(container[:], NewOpClosure(vm, ds.op))
	ds.setSequence(container[:], NewOpFreeGetPtr(vm, ds.op))
	ds.setSequence(container[:], NewOpFreeGet(vm, ds.op))
	ds.setSequence(container[:], NewOpFreeSet(vm, ds.op))
	ds.setSequence(container[:], NewOpLocalPtrGet(vm, ds.op))
	ds.setSequence(container[:], NewOpFreeSelSet(vm, ds.op))
	ds.setSequence(container[:], NewOpIteratorInit(vm, ds.op))
	ds.setSequence(container[:], NewOpIteratorNext(vm, ds.op))
	ds.setSequence(container[:], NewOpIteratorKey(vm, ds.op))
	ds.setSequence(container[:], NewOpIteratorValue(vm, ds.op))
	ds.setSequence(container[:], NewOpIntOp(vm, ds.op))
	ds.setSequence(container[:], NewOpDeref(vm, ds.op))
	ds.setSequence(container[:], NewOpNoOp(vm, ds.op))
	ds.setSequence(container[:], NewOpSuspend(vm, ds.op))
	return container
}

// setSequence assigns the given IOpExecutor to the Sequencer's container, using the bit-masked opcode as the index.
func (ds *Sequencer) setSequence(container []core.IOpExecutor, seq core.IOpExecutor) {
	container[seq.OpcodeId()&bytecode.OpcodesMask] = seq
}

func (ds *Sequencer) facadeForOpcode(opcodeId bytecode.OpcodeId, vm *core.VM) interface{} {
	switch opcodeId {
	// Category: Stack & Constants (Read-Only)
	case bytecode.OpConstant, bytecode.OpGlobalGet, bytecode.OpReferences, bytecode.OpFreeGet, bytecode.OpLocalGet:
		return IVMReadOnly(vm)
	// Category: Control Flow
	case bytecode.OpJump, bytecode.OpJumpFalsy, bytecode.OpJumpAnd, bytecode.OpJumpOr:
		return IVMControlFlow(vm)
	// Category: Simple Stack
	case bytecode.OpPop, bytecode.OpTrue, bytecode.OpFalse, bytecode.OpNull, bytecode.OpNotLogical, bytecode.OpMinus:
		// These instructions only manipulate the top of the stack.
		return IVMStackOnly(vm)
	// Category: Full Access (use with caution)
	case bytecode.OpCall, bytecode.OpReturn, bytecode.OpCallMethod, bytecode.OpClosure:
		// These complex operations need wider access.
		return IVMFullAccess(vm)
	// Default Case
	default:
		// Basic read/write access.
		return IVMReadWrite(vm)
	}
}
