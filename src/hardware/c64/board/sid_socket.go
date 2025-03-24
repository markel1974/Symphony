package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// SIDSocket represents a socket connected to a Board for managing or interacting with its state or functionality.
type SIDSocket struct {
	references.ISID
	fragFreq int
	rasters  int
	player   references.IAudioRender
}

// NewSIDSocket creates and returns a new instance of SIDSocket with default initialization.
func NewSIDSocket(fragFreq int, rasters int) *SIDSocket {
	c := &SIDSocket{
		ISID:     nil,
		fragFreq: fragFreq,
		rasters:  rasters,
		player:   nil,
	}
	return c
}

func (w *SIDSocket) SetPlayer(player references.IAudioRender) {
	w.player = player
}

// Setup initializes the SIDSocket with the provided Board instance, assigning it to the internal board field.
func (w *SIDSocket) Setup(c map[string]references.IComponent, cfg *config.Config) error {
	var err error
	if w.ISID, err = references.ComponentsToISID(c, 0); err != nil {
		return err
	}
	if err = w.ISID.Setup(w, w.player, w.fragFreq, w.rasters, cfg); err != nil {
		return err
	}
	return nil
}

func (s *SIDSocket) Connect() error {
	return nil
}
