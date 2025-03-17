package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// CPUSocket represents a connection hub for the 6510 CPU, PIC, and PLA components to integrate and interact cohesively.
type CPUSocket struct {
	references.I6510
	pic references.IPIC6510
	pla references.IPlaC64
}

// NewCPUSocket creates and returns a new instance of CPUSocket with its internal references uninitialized.
func NewCPUSocket() *CPUSocket {
	c := &CPUSocket{
		I6510: nil,
		pic:   nil,
		pla:   nil,
	}
	return c
}

// Connect initializes the CPUSocket with the provided CPU, PIC, and PLA, and sets up the CPU for interaction.
func (w *CPUSocket) Connect(cpu references.I6510, pic references.IPIC6510, pla references.IPlaC64) error {
	w.I6510 = cpu
	w.pic = pic
	w.pla = pla
	if err := w.I6510.Setup(w); err != nil {
		return err
	}
	return nil
}

// GetPic returns the programmable interrupt controller (PIC) associated with the CPUSocket.
func (w *CPUSocket) GetPic() references.IPIC6510 {
	return w.pic
}

// GetBanks retrieves the memory bank interface used for managing and accessing read/write operations within the CPU socket.
func (w *CPUSocket) GetBanks() references.I6510Banks {
	return w.pla
}
