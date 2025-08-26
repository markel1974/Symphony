package vm

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
	"github.com/markel1974/c64emu/src/kernel/vm/sequencers/native"
)

// NewVM creates a new instance of core.VM with a sequencer determined by the provided sequencerId and given dependencies.
func NewVM(sequencerId string, factory objects.IGateKeeper, op *bytecode.Opcodes) *core.VM {
	switch sequencerId {
	case "native":
		seq := native.NewSequencer(op)
		return core.New(factory, seq)
	default:
		seq := native.NewSequencer(op)
		return core.New(factory, seq)
	}
}
