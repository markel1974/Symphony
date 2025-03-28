package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// CPUSocket represents a connection hub for the 6510 CPU, PIC, and PLA components to integrate and interact cohesively.
type CPUSocket struct {
	references.I6510
	pic   references.IPIC6510
	banks references.IPlaC64
}

// NewCPUSocket creates and returns a new instance of CPUSocket with its internal references uninitialized.
func NewCPUSocket() *CPUSocket {
	c := &CPUSocket{
		I6510: nil,
	}
	return c
}

// Mount initializes the CPUSocket with the provided CPU, PIC, and PLA, and sets up the CPU for interaction.
func (w *CPUSocket) Mount(cc map[string]references.IComponent, _ *config.Config, label string) error {
	var err error
	if w.I6510, err = references.ComponentsToI6510(cc, label, 0); err != nil {
		return err
	}
	if w.pic, err = references.ComponentsToIPIC6510(cc, label, 0); err != nil {
		return err
	}
	if w.banks, err = references.ComponentsToIPLAc64(cc, label, 0); err != nil {
		return err
	}
	if err = w.I6510.Bind(w, w.pic, w.banks); err != nil {
		return err
	}
	return nil
}
