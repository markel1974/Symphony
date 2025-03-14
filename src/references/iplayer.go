package references

// IPlayer is an interface for managing player-related operations in a game or multimedia context.
// GetCurrentPosition returns the current position of the player.
// Write writes audio or data buffer with specified parameters.
type IPlayer interface {
	GetCurrentPosition() int
	Write([]uint32, int, int)
}
