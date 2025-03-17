package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// PLASocket represents a socket that integrates a programmable logic array (PLA) into the board for emulation purposes.
// It combines board interactions and the IPLAc1541 interface for managing the logic needed in emulation scenarios.
type PLASocket struct {
	references.IPLAc1541
}

// NewPLASocket creates and returns a new instance of PLASocket with its fields initialized to nil.
func NewPLASocket() *PLASocket {
	c := &PLASocket{
		IPLAc1541: nil,
	}
	return c
}

// Connect establishes links between the PLASocket and specified components, initializing them using the provided configuration.
// It sets up the PLA logic by associating it with VIA components, the ROM loader, and the configuration data.
func (w *PLASocket) Connect(pla references.IPLAc1541, via1 references.IVIA, via2 references.IVIA, roms references.IROMLoaderC1541, cfg *config.Config) error {
	w.IPLAc1541 = pla
	if err := w.IPLAc1541.Setup(via1, via2, roms, cfg); err != nil {
		return err
	}
	return nil
}
