package mechanic

import (
	"github.com/markel1974/c64emu/src/hardware/c1541/disk/gcr"
	"github.com/markel1974/c64emu/src/hardware/c1541/disk/void"
)

// IDisk is an interface representing a disk with operations to load, read, write, rotate, and manage tracks and sectors.
// It provides methods to interact with the disk's data, manage the read/write head, and determine the disk's usability.
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

// Factory represents an abstract factory used to create and manage disk instances based on provided image data.
type Factory struct {
}

// NewFactory initializes and returns a new instance of Factory.
func NewFactory() *Factory {
	return &Factory{}
}

// Create initializes and returns an IDisk instance based on the provided image data or an empty disk if the image is nil.
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
