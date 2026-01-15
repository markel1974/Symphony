package references

import "github.com/markel1974/symphony/src/config"

// IAudioRender is an interface for managing player-related operations in a game or multimedia context.
// GetCurrentPosition returns the current position of the player.
// Write writes audio or data buffer with specified parameters.
type IAudioRender interface {
	Setup(cfg *config.Config) error

	Write(*[]float32, int)
}
