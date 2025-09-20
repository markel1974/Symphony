package sequencers

import (
	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/sequencers/native"
)

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
