package sequencers

import (
	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/sequencers/native"
)

// NewSequencers creates a new instance of a sequencer based on the provided id and initializes it.
// Returns an implementation of core.ISequencer or an error if the setup process fails.
func NewSequencers(id string) (core.ISequencer, error) {
	switch id {
	default:
		seq := native.NewSequencer()
		if err := seq.Setup(); err != nil {
			return nil, err
		}
		return seq, nil
	}
}
