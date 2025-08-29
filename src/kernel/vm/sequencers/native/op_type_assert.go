// File: vm/sequencers/native/op_type_assert.go

package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// init initializes the system by registering the NewOpTypeAssert operation with the sequencer via SequencerRegister.
func init() {
	SequencerRegister(NewOpTypeAssert)
}

// OpTypeAssert represents a bytecode operation for performing type assertions in a virtual machine.
// It embeds bytecode.Opcode to utilize opcode-related functionalities and operates on a core.IVMFullAccess instance.
type OpTypeAssert struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpTypeAssert creates a new instance of OpTypeAssert, ensuring the provided IVM implements IVMFullAccess.
func NewOpTypeAssert(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpTypeAssert{
		Opcode: op.Opcode(bytecode.OpTypeAssert),
		vm:     vmT,
	}, nil
}

// Execute processes the type assertion operation. It attempts to assert if the top stack object matches the desired type.
// It decodes the target type from the constant pool using the provided decoder and validates against the interface object.
// On success, the concrete value and a boolean 'true' are pushed onto the stack. On failure, undefined and 'false' are pushed.
func (op *OpTypeAssert) Execute(decoder *core.Decoder) {
	// The operand is the index of the target type name in the constants table.
	typeNameIndex := decoder.Read(0)

	// The interface object is at the top of the stack.
	interfaceObj := op.vm.Stack().Pop()

	targetTypeObj := op.vm.Constants().Get(uint(typeNameIndex))
	targetTypeName, ok := targetTypeObj.(*objects.String)
	if !ok {
		op.vm.SetError(fmt.Errorf("constant for type assertion is not a string"))
		return
	}

	io, isInterface := interfaceObj.(*objects.Interface)
	if !isInterface {
		// If not an interface, assertion always fails.
		op.vm.Stack().Push(op.vm.Factory().UndefinedValue())
		op.vm.Stack().Push(op.vm.Factory().FalseValue())
		return
	}

	concreteValue := io.Value()
	if concreteValue.TypeName() == targetTypeName.Value() {
		// Success!
		op.vm.Stack().Push(concreteValue)               // Push the "unwrapped" concrete value.
		op.vm.Stack().Push(op.vm.Factory().TrueValue()) // Push the 'ok' boolean.
	} else {
		// Failure.
		op.vm.Stack().Push(op.vm.Factory().UndefinedValue()) // Push a null value.
		op.vm.Stack().Push(op.vm.Factory().FalseValue())     // Push the 'ok' boolean.
	}
}
