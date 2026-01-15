package disk

import (
	"github.com/markel1974/symphony/src/hardware/c1541_board_rev1/disk/gcr"
)

// Factory represents an abstract factory used to create and manage disk instances based on provided image data.
type Factory struct {
}

// NewFactory initializes and returns a new instance of Factory.
func NewFactory() *Factory {
	return &Factory{}
}

// Create initializes and returns an IDisk instance based on the provided image data or an empty disk if the image is nil.
func (f *Factory) Create(image []byte, wp bool) (IDisk, error) {
	g := gcr.NewDisk(wp)
	if err := g.Load(image); err != nil {
		return nil, err
	}
	return g, nil
}
