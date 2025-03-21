package iec

// IGlobal represents an interface for managing and retrieving a state value as an unsigned 8-bit integer.
type IGlobal interface {
	GetState() uint8
	SetState() uint8
}

// Global represents a structure encapsulating the state of a global entity as an unsigned 8-bit integer.
type Global struct {
	state uint8
}

// NewGlobalState initializes and returns a pointer to a new Global instance with its state set to 0.
func NewGlobalState() *Global {
	return &Global{
		state: 0,
	}
}

// SetState updates the internal state of the Global instance with the given value.
func (v *Global) SetState(state uint8) {
	v.state = state
}

// GetState retrieves the current state value of the Global instance.
func (v *Global) GetState() uint8 {
	return v.state
}

// _gs is a singleton instance of Global, initialized using NewGlobalState, and shared across various components.
var _gs = NewGlobalState()
