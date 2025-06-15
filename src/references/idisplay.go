package references

import "github.com/markel1974/c64emu/src/config"

// IDisplayBuffer defines methods for interacting with a display buffer, allowing data manipulation at specified indices.
// Set sets a single byte of data at the given index.
// SetMulti8 sets a single byte of data and applies it across multiple relevant sections.
// Set8 sets an array of 8 bytes of data starting at the given index.
type IDisplayBuffer interface {
	Set(idx int, data uint8)

	SetMulti8(idx int, data uint8)

	Set8(idx int, data [8]uint8)
}

// IDisplayRender provides methods for rendering displays and managing display buffers.
// Setup initializes the renderer with a board and configuration.
// Start begins the rendering process and returns an error if it fails.
// CreateDisplayBuffer creates a new display buffer with the specified width and height.
type IDisplayRender interface {
	Setup(board IC64Board, cfg *config.Config) error

	Start() error

	CreateDisplayBuffer(w int, h int) (IDisplayBuffer, error)
}
