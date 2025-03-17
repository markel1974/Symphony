package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// CPUSocket represents a CPU interface that connects to a board, picture processing unit, and memory banks for operations.
type CPUSocket struct {
	board *Board
	references.I6510
}

// NewCPUSocket creates and returns a new instance of CPUSocket with initialized properties.
func NewCPUSocket() *CPUSocket {
	c := &CPUSocket{
		board: nil,
		I6510: nil,
	}
	return c
}

// Connect initializes the CPUSocket by linking it with a given Board, setting the PIC and memory banks from the board.
func (w *CPUSocket) Connect(board *Board, cpu references.I6510) error {
	w.board = board
	w.I6510 = cpu
	if err := w.I6510.Setup(w); err != nil {
		return err
	}
	return nil
}

// GetPic retrieves the programmable interrupt controller (IPIC6510) associated with the CPUSocket instance.
func (w *CPUSocket) GetPic() references.IPIC6510 {
	return w.board.picSocket
}

// GetBanks retrieves the memory banks interface associated with the CPUSocket.
func (w *CPUSocket) GetBanks() references.I6510Banks {
	return w.board.plaSocket
}
