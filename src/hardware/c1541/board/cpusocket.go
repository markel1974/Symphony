package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// CPUSocket represents the CPU's interface to the board, providing access to the interrupt controller and memory pla.
type CPUSocket struct {
	references.I6510
	pic   references.IPIC6510
	banks references.I6510Banks
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

// Connect initializes the CPUSocket by associating it with the provided Board and configuring related components.
func (w *CPUSocket) Connect(cpu references.I6510, pic references.IPIC6510, pla references.IPLAc1541, via2 references.IVIA) error {
	w.pic = pic
	w.banks = pla
	w.I6510 = cpu
	if err := w.I6510.Setup(w); err != nil {
		return err
	}
	w.I6510.SetOverflowBranch(via2.ByteReady)
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
