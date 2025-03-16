package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// CPUSocket represents the CPU's interface to the board, providing access to the interrupt controller and memory pla.
type CPUSocket struct {
	board *Board
	cpu   references.I6510
	pic   references.IPic6510
	banks references.I6510Banks
}

// NewCPUSocket initializes and returns a new instance of CPUSocket with default nil components.
func NewCPUSocket() *CPUSocket {
	c := &CPUSocket{
		board: nil,
		cpu:   nil,
		pic:   nil,
		banks: nil,
	}
	return c
}

// Setup initializes the CPUSocket by associating it with the provided Board and configuring related components.
func (w *CPUSocket) Setup(board *Board, cpu references.I6510) {
	w.board = board
	w.pic = board.pic
	w.banks = board.pla
	w.cpu = cpu
	w.cpu.Setup(w)
	w.cpu.SetOverflowBranch(w.board.via2Socket.ByteReady())
}

// Reset reinitializes the CPU within the board to its default state.
func (w *CPUSocket) Reset() {
	w.cpu.Reset()
}

func (w *CPUSocket) Emulate() {
	w.cpu.Emulate()
}

// GetPic retrieves the programmable interrupt controller (PIC) associated with the CPUSocket instance.
func (w *CPUSocket) GetPic() references.IPic6510 {
	return w.pic
}

// GetBanks returns the memory bank interface associated with the current CPUSocket.
func (w *CPUSocket) GetBanks() references.I6510Banks {
	return w.banks
}
