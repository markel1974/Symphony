package mechanic

import (
	"github.com/markel1974/c64emu/src/c1541/disk/gcr"
	"github.com/markel1974/c64emu/src/c1541/disk/void"
)

type IDisk interface {
	Read() uint8
	Write(uint8)
	GetTracksNumber() uint8
	SetHeadTrack(uint8)
	Rotate()
	Usable() bool
}

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) Create(image []byte) (IDisk, error) {
	if image == nil {
		return void.NewDisk(), nil
	}
	g, err := gcr.NewDisk(image)
	return g, err
}
