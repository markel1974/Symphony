package mechanic

import (
	"github.com/markel1974/c64emu/src/c1541/disk/gcr"
	"github.com/markel1974/c64emu/src/c1541/disk/void"
)

type IDisk interface {
	Load(image []byte) error
	Read() uint8
	Write(uint8)
	Next() uint8
	SetHeadHalfTrack(uint8) int
	TrackLen() int
	TrackSectors() uint8
	MicroSecPerByte() uint8
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
	g := gcr.NewDisk()
	if len(image) > 0 {
		if err := g.Load(image); err != nil {
			return nil, err
		}
	}
	return g, nil
}
