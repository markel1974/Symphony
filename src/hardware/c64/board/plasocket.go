package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// PLASocket represents a socket encapsulating a PLA and its associated board.
type PLASocket struct {
	references.IPlaC64
}

// NewPLASocket creates a new instance of PLASocket with default uninitialized board and PLA components.
func NewPLASocket() *PLASocket {
	c := &PLASocket{
		IPlaC64: nil,
	}
	return c
}

// Connect initializes the PLASocket with the provided board, PLA, and associated components like VIC, SID, CIAs, and cartridge manager.
func (w *PLASocket) Connect(p references.IPlaC64, vic references.IVIC, sid references.ISID, cia1 references.ICIA, cia2 references.ICIA, cartMan references.ICartridgeManagerC64, roms references.IROMLoaderC64, cfg *config.Config) error {
	w.IPlaC64 = p
	w.IPlaC64.Setup(vic, sid, cia1, cia2, cartMan, roms, cfg)
	return nil
}
