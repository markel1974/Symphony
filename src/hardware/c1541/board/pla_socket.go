package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// PLASocket represents a socket that integrates a programmable logic array (PLA) into the board for emulation purposes.
// It combines board interactions and the IPLAc1541 interface for managing the logic needed in emulation scenarios.
type PLASocket struct {
	references.IPLAc1541
	via1 references.IVIA
	via2 references.IVIA
	roms references.IROMLoaderC1541
	cfg  *config.Config
}

// NewPLASocket creates and returns a new instance of PLASocket with its fields initialized to nil.
func NewPLASocket() *PLASocket {
	c := &PLASocket{
		IPLAc1541: nil,
	}
	return c
}

// Setup establishes links between the PLASocket and specified components, initializing them using the provided configuration.
// It sets up the PLA logic by associating it with VIA components, the ROM loader, and the configuration data.
func (w *PLASocket) Setup(c map[string]references.IComponent, cfg *config.Config) error {
	var err error
	w.IPLAc1541, err = references.ComponentsToIPLAc1541(c, 0)
	if err != nil {
		return err
	}
	w.via1, err = references.ComponentsToIVIA(c, 0)
	if err != nil {
		return err
	}
	w.via2, err = references.ComponentsToIVIA(c, 1)
	if err != nil {
		return err
	}
	w.roms, err = references.ComponentsToIROMLoaderC1541(c, 0)
	if err != nil {
		return err
	}
	if err = w.IPLAc1541.Setup(w, w.cfg); err != nil {
		return err
	}
	return nil
}

func (w *PLASocket) Connect() error {
	if err := w.IPLAc1541.Connect(w.via1, w.via2, w.roms); err != nil {
		return err
	}
	return nil
}
