package disk

import (
	"github.com/markel1974/c64emu/src/hardware/c1541_board_rev1/disk/gcr"
)

// IDisk represents an interface for disk operations, allowing reading, writing, loading images, and managing track state.
// WriteProtected checks if the disk is write-protected and returns a boolean indicating its status.
// Load initializes the disk by loading image data and setting up tracks and sectors.
// Read retrieves the current byte from the disk at the cursor position of the active track.
// Write writes a byte value to the current position on the active track of the disk.
// Next retrieves the byte value at the next position without advancing the track cursor.
// SetHeadHalfTrack updates the disk head to a specified half-track position.
// TrackLen provides the length of the current track's data in bytes.
// TrackSectors returns the total number of sectors present in the current disk track.
// MicroSecPerByte calculates and returns the time required to process one byte on the current track in microseconds.
// Rotate simulates disk rotation by advancing the track's cursor to the next position.
// Usable checks if the disk is functional and returns its usability status as a boolean.
type IDisk interface {
	WriteProtected() bool

	Load(image []byte) error

	Read() uint8

	Write(uint8)

	Next() uint8

	SetHeadHalfTrack(uint8) bool

	TrackLen() int

	TrackSectors() uint8

	MicroSecPerByte() int

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
func (f *Factory) Create(image []byte, wp bool) (IDisk, error) {
	g := gcr.NewDisk(wp)
	if err := g.Load(image); err != nil {
		return nil, err
	}
	return g, nil
}
