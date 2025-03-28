package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// CPUSocket represents a CPU socket managing the integration of a 6510 CPU, programmable interrupt controller, memory banks, and VIA.
type CPUSocket struct {
	references.I6510
	pic  references.IPIC6510
	pla  references.IPLAc1541
	via2 references.IVIA
}

// NewCPUSocket creates and initializes a new CPUSocket instance with default nil values for its fields.
func NewCPUSocket() *CPUSocket {
	c := &CPUSocket{
		I6510: nil,
		pic:   nil,
		pla:   nil,
	}
	return c
}

// Mount initializes the CPUSocket by linking required components and configuring dependencies using the provided map and config.
// Returns an error if any component setup or binding fails.
func (w *CPUSocket) Mount(c map[string]references.IComponent, _ *config.Config, label string) error {
	var err error
	if w.pic, err = references.ComponentsToIPIC6510(c, label, 0); err != nil {
		return err
	}
	if w.pla, err = references.ComponentsToIPLAc1541(c, label, 0); err != nil {
		return err
	}
	if w.I6510, err = references.ComponentsToI6510(c, label, 0); err != nil {
		return err
	}
	if w.via2, err = references.ComponentsToIVIA(c, label, 1); err != nil {
		return err
	}
	if err = w.I6510.Bind(w, w.pic, w.pla); err != nil {
		return err
	}
	w.I6510.SetOverflowBranch(w.via2.ByteReady)
	return nil
}
