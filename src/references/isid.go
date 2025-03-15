package references

import "github.com/markel1974/c64emu/src/config"

// ISid is an interface representing a Sound Interface Device (SID) for audio synthesis and emulation in a system.
// Setup initializes the interface with a socket, frequency, rasters, and a configuration.
// Reset resets the internal state of the SID to default.
// SetPotX sets the X-axis of the paddle controller input.
// SetPotY sets the Y-axis of the paddle controller input.
// Prepare prepares the SID for the next operational update.
// Update processes the SID's internal state and audio generation for the current cycle.
// WriteRegister writes data to a specific register address of the SID.
// ReadRegister retrieves data from a specified register address of the SID.
type ISid interface {
	Setup(socket ISidSocket, fragFreq int, rasters int, cfg *config.Config)

	Reset()

	SetPotX(x uint8)

	SetPotY(y uint8)

	Prepare()

	Update()

	WriteRegister(addr uint16, data uint8)

	ReadRegister(addr uint16) uint8
}

// ISidSocket is an interface representing a socket for SID integration and player management functionality.
// GetPlayer retrieves the IPlayer instance associated with the socket, enabling player-related operations.
type ISidSocket interface {
	GetPlayer() IPlayer
}
