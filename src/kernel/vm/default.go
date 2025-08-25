package vm

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/executors"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func NewVM(factory objects.IGateKeeper, op *bytecode.Opcodes) *core.VM {
	seq := executors.NewSequencer(op)
	vm := core.New(factory, seq)
	return vm
}
