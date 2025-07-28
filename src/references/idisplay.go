package references

import "github.com/markel1974/c64emu/src/config"

// IDisplayBuffer defines an interface for managing and setting data in a display buffer.
// SetArray sets a section of the buffer at a specified index with provided data and properties like width.
type IDisplayBuffer interface {
	SetArray(idx int, data *[]uint8, width int)
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
