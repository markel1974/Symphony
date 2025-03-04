package board

import mos6510 "github.com/markel1974/c64emu/src/components/6510"

// CPUSocket represents the CPU's interface to the board, providing access to the interrupt controller and memory banks.
type CPUSocket struct {
	board *Board
	pic   mos6510.IPic
	banks mos6510.IBanks
}

// NewCPUSocket initializes and returns a new instance of CPUSocket with default nil components.
func NewCPUSocket() *CPUSocket {
	c := &CPUSocket{
		board: nil,
		pic:   nil,
		banks: nil,
	}
	return c
}

// Setup initializes the CPUSocket by associating it with the provided Board and configuring related components.
func (w *CPUSocket) Setup(board *Board) {
	w.board = board
	w.pic = board.pic
	w.banks = board.banks
}

// Reset reinitializes the CPU within the board to its default state.
func (w *CPUSocket) Reset() {
	w.board.cpu.Reset()
}

// GetPic retrieves the programmable interrupt controller (PIC) associated with the CPUSocket instance.
func (w *CPUSocket) GetPic() mos6510.IPic {
	return w.pic
}

// GetBanks returns the memory bank interface associated with the current CPUSocket.
func (w *CPUSocket) GetBanks() mos6510.IBanks {
	return w.banks
}
