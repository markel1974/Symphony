package vm

import (
	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	_nativeSequencer "github.com/markel1974/c64emu/src/vm/sequencers/native"
)

// NewVM creates a new instance of core.VM with a sequencer determined by the provided sequencerId and given dependencies.
func NewVM(factory objects.IGateKeeper, op *bytecode.Opcodes, sequencerId string) (*core.VM, error) {
	var seq core.ISequencer
	switch sequencerId {
	case "native":
		seq = _nativeSequencer.NewSequencer(op)
	default:
		seq = _nativeSequencer.NewSequencer(op)
	}
	return core.New(factory, seq, op), nil
}
