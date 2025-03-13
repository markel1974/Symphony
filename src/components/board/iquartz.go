package board

// IQuartz represents an interface for components that provide cyclic timing functionality.
// Cycle retrieves the number of cycles elapsed since the start of the quartz-based timing mechanism.
type IQuartz interface {
	Cycle() uint64

	ToUSec(uint64) uint64
}
