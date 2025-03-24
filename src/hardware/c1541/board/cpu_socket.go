package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// CPUSocket represents the CPU's interface to the board, providing access to the interrupt controller and memory pla.
type CPUSocket struct {
	references.I6510
	pic   references.IPIC6510
	banks references.IPLAc1541
	via2  references.IVIA
}

// NewCPUSocket initializes and returns a new instance of CPUSocket with default nil components.
func NewCPUSocket() *CPUSocket {
	c := &CPUSocket{
		I6510: nil,
		pic:   nil,
		banks: nil,
	}
	return c
}

// Setup initializes the CPUSocket by associating it with the provided Board and configuring related components.
func (w *CPUSocket) Setup(c map[string]references.IComponent, cfg *config.Config) error {
	var err error
	w.pic, err = references.ComponentsToIPIC6510(c, 0)
	if err != nil {
		return err
	}
	w.banks, err = references.ComponentsToIPLAc1541(c, 0)
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

// GetBanks returns the memory bank interface associated with the current CPUSocket.
func (w *CPUSocket) GetBanks() references.I6510Banks {
	return w.banks
}
