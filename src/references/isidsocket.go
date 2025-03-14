package references

// ISidSocket represents a socket interface primarily used to interact with an IPlayer within a connected system.
// GetPlayer retrieves the associated IPlayer instance for managing audio or game-player-related operations.
type ISidSocket interface {
	GetPlayer() IPlayer
}
