package void

import "fmt"

// Disk is a placeholder struct that implements methods with default or no functionality.
type Disk struct {
}

// NewVoidDisk initializes and returns an empty Void disk instance.
func NewVoidDisk() *Disk {
	return &Disk{}
}

func (e *Disk) WriteProtected() bool {
	return true
}

// Load loads the given image byte slice into the Void instance. Currently, this method is not implemented.
func (e *Disk) Load(image []byte) error {
	return fmt.Errorf("not implemented")
}

// TrackLen returns the length of the track, which is always 0 for the Void type.
func (e *Disk) TrackLen() int {
	return 0
}

// TrackSectors returns the number of sectors in the current track as an unsigned 8-bit integer.
func (e *Disk) TrackSectors() uint8 {
	return 0
}

// Read returns a fixed value of 0 and does not perform any dynamic operations.
func (e *Disk) Read() uint8 {
	return 0
}

// Write writes a byte to the Void type. This implementation does nothing as the Void type is not functional.
func (e *Disk) Write(_ uint8) {
}

// Next returns the next value, always 0 for the Void type implementation.
func (e *Disk) Next() uint8 {
	return 0
}

// SetHeadHalfTrack sets the current head position to a specified half-track.
// The input is a uint8 representing the desired half-track position.
// Returns an integer indicating the result of the operation.
func (e *Disk) SetHeadHalfTrack(uint8) bool {
	return false
}

// MicroSecPerByte returns the number of microseconds required to process a single byte in the Void implementation.
func (e *Disk) MicroSecPerByte() int {
	return 0
}

// Rotate simulates rotating the drive to the next position. Currently, this method does not perform any operation.
func (e *Disk) Rotate() {
}

// Usable returns false, indicating that the Void type is not usable.
func (e *Disk) Usable() bool {
	return false
}
