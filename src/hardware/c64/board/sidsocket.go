package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// SIDSocket represents a socket connected to a Board for managing or interacting with its state or functionality.
type SIDSocket struct {
	references.ISID
	player references.IPlayer
}

// NewSIDSocket creates and returns a new instance of SIDSocket with default initialization.
func NewSIDSocket() *SIDSocket {
	c := &SIDSocket{
		ISID:   nil,
		player: nil,
	}
	return c
}

// Connect initializes the SIDSocket with the provided Board instance, assigning it to the internal board field.
func (w *SIDSocket) Connect(sid references.ISID, player references.IPlayer, fragFreq int, rasters int, cfg *config.Config) error {
	w.ISID = sid
	w.player = player
	w.ISID.Setup(w, fragFreq, rasters, cfg)
	return nil
}

// GetPlayer returns the instance of IPlayer associated with the SIDSocket's board.
func (w *SIDSocket) GetPlayer() references.IPlayer {
	return w.player
}
