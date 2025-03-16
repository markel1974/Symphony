package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// CPUSocket represents a CPU interface that connects to a board, picture processing unit, and memory banks for operations.
type CPUSocket struct {
	board *Board
	cpu   references.I6510
}

// NewCPUSocket creates and returns a new instance of CPUSocket with initialized properties.
func NewCPUSocket() *CPUSocket {
	c := &CPUSocket{
		board: nil,
		cpu:   nil,
	}
	return c
}

// Connect initializes the CPUSocket by linking it with a given Board, setting the PIC and memory banks from the board.
func (w *CPUSocket) Connect(board *Board, cpu references.I6510) error {
	w.board = board
	w.cpu = cpu
	w.cpu.Setup(w)
	return nil
}

// GetPic retrieves the programmable interrupt controller (IPic6510) associated with the CPUSocket instance.
func (w *CPUSocket) GetPic() references.IPic6510 {
	return w.board.pic
}

// GetBanks retrieves the memory banks interface associated with the CPUSocket.
func (w *CPUSocket) GetBanks() references.I6510Banks {
	return w.board.plaSocket
}

// Reset invokes the CPU reset functionality through the associated board, resetting the internal CPU state.
func (w *CPUSocket) Reset() {
	w.cpu.Reset()
}

func (w *CPUSocket) Emulate() {
	w.cpu.Emulate()
}

func (w *CPUSocket) SetRDYLow(rdyLow bool) {
	w.cpu.SetRDYLow(rdyLow)
}

func (w *CPUSocket) SetAECLow(aecLow bool) {
	w.cpu.SetAECLow(aecLow)
}
