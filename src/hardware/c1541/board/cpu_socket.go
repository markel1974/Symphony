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

// Setup initializes the CPUSocket by linking required components and configuring dependencies using the provided map and config.
// Returns an error if any component setup or binding fails.
func (w *CPUSocket) Setup(c map[string]references.IComponent, cfg *config.Config) error {
	var err error
	w.pic, err = references.ComponentsToIPIC6510(c, 0)
	if err != nil {
		return err
	}
	w.pla, err = references.ComponentsToIPLAc1541(c, 0)
	if err != nil {
		return err
	}
	w.I6510, err = references.ComponentsToI6510(c, 0)
	if err != nil {
		return err
	}
	w.via2, err = references.ComponentsToIVIA(c, 1)
	if err != nil {
		return err
	}
	if err = w.I6510.Setup(w, cfg); err != nil {
		return err
	}
	w.I6510.SetOverflowBranch(w.via2.ByteReady)
	return nil
}

// Connect establishes the connection for the 6510 CPU interface and returns an error if the operation fails.
func (w *CPUSocket) Connect() error {
	if err := w.I6510.Connect(); err != nil {
		return err
	}
	return nil
}

// GetPic retrieves the programmable interrupt controller (PIC) associated with the CPUSocket instance.
func (w *CPUSocket) GetPic() references.IPIC6510 {
	return w.pic
}

// GetBanks retrieves the interface for managing and accessing memory banks in the CPUSocket.
func (w *CPUSocket) GetBanks() references.I6510Banks {
	return w.pla
}
