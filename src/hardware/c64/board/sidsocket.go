package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// SIDSocket represents a socket connected to a Board for managing or interacting with its state or functionality.
type SIDSocket struct {
	sid   references.ISID
	board *Board
}

// NewSIDSocket creates and returns a new instance of SIDSocket with default initialization.
func NewSIDSocket() *SIDSocket {
	c := &SIDSocket{}
	return c
}

// Connect initializes the SIDSocket with the provided Board instance, assigning it to the internal board field.
func (w *SIDSocket) Connect(board *Board, sid references.ISID, fragFreq int, rasters int, cfg *config.Config) error {
	w.board = board
	w.sid = sid
	w.sid.Setup(w, fragFreq, rasters, cfg)
	return nil
}

// Reset clears the state of the SID chip by resetting all registers and internal components to their default values.
func (w *SIDSocket) Reset() {
	w.sid.Reset()
}

// SetPotXY sets the POT X and POT Y registers to the specified x and y values for the SID chip.
func (w *SIDSocket) SetPotXY(x uint8, y uint8) {
	w.sid.SetPotX(x)
	w.sid.SetPotY(y)
}

// Prepare initializes or readies the attached SID instance for audio processing by loading necessary register values.
func (w *SIDSocket) Prepare() {
	w.sid.Prepare()
}

// Update triggers the internal update process of the SID chip associated with the SIDSocket.
func (w *SIDSocket) Update() {
	w.sid.Update()
}

// GetPlayer returns the instance of IPlayer associated with the SIDSocket's board.
func (w *SIDSocket) GetPlayer() references.IPlayer {
	return w.board.player
}
