package board

import (
	"github.com/markel1974/c64emu/src/components/board"
)

// SidSocket represents a socket connected to a Board for managing or interacting with its state or functionality.
type SidSocket struct {
	board *Board
}

// NewSidSocket creates and returns a new instance of SidSocket with default initialization.
func NewSidSocket() *SidSocket {
	c := &SidSocket{}
	return c
}

// Setup initializes the SidSocket with the provided Board instance, assigning it to the internal board field.
func (w *SidSocket) Setup(board *Board) {
	w.board = board
}

// Reset clears the state of the SID chip by resetting all registers and internal components to their default values.
func (w *SidSocket) Reset() {
	w.board.sid.Reset()
}

// SetPotXY sets the POT X and POT Y registers to the specified x and y values for the SID chip.
func (w *SidSocket) SetPotXY(x uint8, y uint8) {
	w.board.sid.SetPotX(x)
	w.board.sid.SetPotY(y)
}

// Prepare initializes or readies the attached SID instance for audio processing by loading necessary register values.
func (w *SidSocket) Prepare() {
	w.board.sid.Prepare()
}

// Update triggers the internal update process of the SID chip associated with the SidSocket.
func (w *SidSocket) Update() {
	w.board.sid.Update()
}

// GetPlayer returns the instance of IPlayer associated with the SidSocket's board.
func (w *SidSocket) GetPlayer() board.IPlayer {
	return w.board.player
}
