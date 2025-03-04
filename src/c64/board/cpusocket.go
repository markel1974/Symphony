package board

import mos6510 "github.com/markel1974/c64emu/src/components/6510"

// CPUSocket represents a CPU interface that connects to a board, picture processing unit, and memory banks for operations.
type CPUSocket struct {
	board *Board
	pic   mos6510.IPic
	banks mos6510.IBanks
}

// NewCPUSocket creates and returns a new instance of CPUSocket with initialized properties.
func NewCPUSocket() *CPUSocket {
	c := &CPUSocket{}
	return c
}

// Setup initializes the CPUSocket by linking it with a given Board, setting the PIC and memory banks from the board.
func (w *CPUSocket) Setup(board *Board) {
	w.board = board
	w.pic = board.pic
	w.banks = board.pla
}

// Reset invokes the CPU reset functionality through the associated board, resetting the internal CPU state.
func (w *CPUSocket) Reset() {
	w.board.cpu.Reset()
}

// GetPic retrieves the programmable interrupt controller (IPic) associated with the CPUSocket instance.
func (w *CPUSocket) GetPic() mos6510.IPic {
	return w.pic
}

// GetBanks retrieves the memory banks interface associated with the CPUSocket.
func (w *CPUSocket) GetBanks() mos6510.IBanks {
	return w.banks
}
