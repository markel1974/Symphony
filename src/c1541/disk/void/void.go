package void

import "fmt"

// Void is a placeholder struct that implements methods with default or no functionality.
type Void struct{}

// NewDisk initializes and returns an empty Void disk instance.
func NewDisk() *Void {
	return &Void{}
}

// Load loads the given image byte slice into the Void instance. Currently, this method is not implemented.
func (e *Void) Load(image []byte) error {
	return fmt.Errorf("not implemented")
}

// TrackLen returns the length of the track, which is always 0 for the Void type.
func (e *Void) TrackLen() int {
	return 0
}

// TrackSectors returns the number of sectors in the current track as an unsigned 8-bit integer.
func (e *Void) TrackSectors() uint8 {
	return 0
}

// Read returns a fixed value of 0 and does not perform any dynamic operations.
func (e *Void) Read() uint8 {
	return 0
}

// Write writes a byte to the Void type. This implementation does nothing as the Void type is not functional.
func (e *Void) Write(_ uint8) {
}

// Next returns the next value, always 0 for the Void type implementation.
func (e *Void) Next() uint8 {
	return 0
}

// SetHeadHalfTrack sets the current head position to a specified half-track.
// The input is a uint8 representing the desired half-track position.
// Returns an integer indicating the result of the operation.
func (e *Void) SetHeadHalfTrack(uint8) int {
	return 0
}

// MicroSecPerByte returns the number of microseconds required to process a single byte in the Void implementation.
func (e *Void) MicroSecPerByte() uint8 {
	return 0
}

// Rotate simulates rotating the drive to the next position. Currently, this method does not perform any operation.
func (e *Void) Rotate() {
}

// Usable returns false, indicating that the Void type is not usable.
func (e *Void) Usable() bool {
	return false
}
