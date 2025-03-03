package externalcpu

import (
	mos6510 "github.com/markel1974/c64emu/src/components/6510"
)

// CPUSocket represents a connector that integrates programmable interrupt controller and memory banks for a CPU system.
type CPUSocket struct {
	pic   mos6510.IPic
	banks mos6510.IBanks
}

// NewCPUSocket creates and returns a new instance of CPUSocket.
func NewCPUSocket() *CPUSocket {
	c := &CPUSocket{}
	return c
}

// Setup initializes the CPUSocket by setting its programmable interrupt controller and memory banks from the provided board.
func (w *CPUSocket) Setup(board *ExternalCPU) {
	w.pic = board.pic
	w.banks = board.board
}

// GetPic returns the current programmable interrupt controller (IPic) associated with the CPUSocket instance.
func (w *CPUSocket) GetPic() mos6510.IPic {
	return w.pic
}

// GetBanks returns the IBanks interface used for managing memory banks in the CPUSocket.
func (w *CPUSocket) GetBanks() mos6510.IBanks {
	return w.banks
}
