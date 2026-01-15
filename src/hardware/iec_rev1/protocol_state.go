package iec_rev1

import (
	"github.com/markel1974/symphony/src/references"
	"log"
)

type ProtocolState struct {
	flags         uint8
	primary       uint8
	secondaryPrev uint8
	secondary     uint8
	data          uint8
	stateMachine  uint8
	state         [stateLast + 1]uint8
	timeout       uint64
}

func NewProtocolState() *ProtocolState {
	return &ProtocolState{
		flags:         0,
		primary:       0,
		secondaryPrev: 0,
		secondary:     0,
		data:          0,
		stateMachine:  0,
		timeout:       0,
	}
}

func (v *ProtocolState) Reset() {
	v.flags = 0
	v.timeout = 0
	for i := 0; i < len(v.state); i++ {
		v.state[i] = 0
	}
}

// FlagsSet sets the specified flags in the Protocol's internal flags field using a bitwise OR operation.
func (v *ProtocolState) FlagsSet(f uint8) {
	v.flags |= f
}

// FlagsRemove removes specific flags from the Protocol's flags field using bitwise operations.
func (v *ProtocolState) FlagsRemove(f uint8) {
	v.flags &= ^f
}

// FlagGet checks if the specified flag is set in the Protocol's internal flags field. Returns true if set, false otherwise.
func (v *ProtocolState) FlagGet(f uint8) bool {
	return (v.flags & f) != 0
}

// DataHasBit checks if the specified bit is set in the data field of the Protocol instance. Returns true if the bit is set.
func (v *ProtocolState) DataHasBit(bit uint8) bool {
	return (v.data & bit) != 0
}

// DataSetBit sets the specified bit(s) in the `data` field to 1 without altering other bits.
func (v *ProtocolState) DataSetBit(m uint8) {
	//m := uint8(1 << pos)
	v.data |= m
}

// DataClearBit clears specific bits in the `data` field of the Protocol struct, based on the provided mask `m`.
func (v *ProtocolState) DataClearBit(m uint8) {
	//m := uint8(^(1 << pos))
	n := ^m
	v.data &= n
}

// DataSetByte sets the Protocol's `data` field to the specified byte value `b`.
func (v *ProtocolState) DataSetByte(b uint8) {
	v.data = b
}

// DataGetByte returns the current value of the `data` field in the Protocol structure.
func (v *ProtocolState) DataGetByte() uint8 {
	return v.data
}

// StateMachineSet sets the state machine to the specified state value.
func (v *ProtocolState) StateMachineSet(m uint8) {
	v.stateMachine = m
}

// StateMachineGet returns the current state of the protocol's state machine as an unsigned 8-bit integer.
func (v *ProtocolState) StateMachineGet() uint8 {
	return v.stateMachine
}

// StateMachineAdvance increments the stateMachine variable by one, advancing the state machine to the next state.
func (v *ProtocolState) StateMachineAdvance() {
	v.stateMachine++
}

// PrimarySet sets the primary device address to the specified value.
func (v *ProtocolState) PrimarySet(p uint8) {
	v.primary = p
}

// PrimaryGet retrieves the current primary address of the Protocol. Returns the value as an unsigned 8-bit integer.
func (v *ProtocolState) PrimaryGet() uint8 {
	return v.primary
}

// SecondarySet sets the value of the secondary address/state in the Protocol instance.
func (v *ProtocolState) SecondarySet(s uint8) {
	v.secondary = s
}

// SecondaryGet retrieves the current value of the secondary byte in the Protocol instance.
func (v *ProtocolState) SecondaryGet() uint8 {
	return v.secondary
}

// SecondaryPrevSet sets the secondary previous address value in the Protocol instance.
func (v *ProtocolState) SecondaryPrevSet() {
	s := v.secondary
	v.secondaryPrev = s
}

// SecondaryPrevGet retrieves the previous value of the secondary byte within the Protocol structure.
func (v *ProtocolState) SecondaryPrevGet() uint8 {
	return v.secondaryPrev
}

// StateSet updates the state at the given index after masking the index with stateLast.
func (v *ProtocolState) StateSet(idx uint8, s uint8) {
	x := idx & stateLast
	v.state[x] = s
}

// StateGet retrieves the state value at the given index from the state array after masking the index with stateLast.
func (v *ProtocolState) StateGet(idx uint8) uint8 {
	x := idx & stateLast
	return v.state[x]
}

// TimeoutSet sets a timeout in microseconds by calculating the required number of cycles and updating the timeout property.
func (v *ProtocolState) TimeoutSet(q references.IQuartz, uSec uint64) {
	cycles := q.USecToCycleRounded(uSec)
	v.timeout = q.Cycle() + cycles
}

// TimeoutExpired checks if the current cycle exceeds or equals the timeout value, indicating that the timeout has expired.
func (v *ProtocolState) TimeoutExpired(q references.IQuartz) bool {
	if b := q.Cycle(); b >= v.timeout {
		return true
	}
	return false
}

// Print logs the current state of the Protocol instance along with an identifier and bus number.
func (v *ProtocolState) Print(id string, bus uint8) {
	log.Printf("%s -> bus: %d, stateMachine: %d, flags: %d, primary: %d, secondary: %d, secondaryPrev: %d\n", id, bus, v.stateMachine, v.flags, v.primary, v.secondary, v.secondaryPrev)
}
