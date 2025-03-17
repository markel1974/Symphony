package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// CPUSocket represents the CPU's interface to the board, providing access to the interrupt controller and memory pla.
type CPUSocket struct {
	board *Board
	references.I6510
	pic   references.IPIC6510
	banks references.I6510Banks
}

// NewCPUSocket initializes and returns a new instance of CPUSocket with default nil components.
func NewCPUSocket() *CPUSocket {
	c := &CPUSocket{
		board: nil,
		I6510: nil,
		pic:   nil,
		banks: nil,
	}
	return c
}

// Connect initializes the CPUSocket by associating it with the provided Board and configuring related components.
func (w *CPUSocket) Connect(board *Board, cpu references.I6510) {
	w.board = board
	w.pic = board.pic
	w.banks = board.pla
	w.I6510 = cpu
	w.I6510.Setup(w)
	w.I6510.SetOverflowBranch(w.board.via2Socket.ByteReady())
}

// GetPic retrieves the programmable interrupt controller (PIC) associated with the CPUSocket instance.
func (w *CPUSocket) GetPic() references.IPIC6510 {
	return w.pic
}

// GetBanks returns the memory bank interface associated with the current CPUSocket.
func (w *CPUSocket) GetBanks() references.I6510Banks {
	return w.banks
}
