package references

// opFlagIrqDisabled represents the flag for interrupt requests being disabled.
// opFlagIrqEnabled represents the flag for interrupt requests being enabled.
// opFlagIntDelayed represents the flag for delayed interrupt requests.
const (
	OpFlagIrqDisabled = 0x01
	OpFlagIrqEnabled  = 0x02
	OpFlagIntDelayed  = 0x04
)

// I6510Pic provides an interface for managing the programmable interrupt controller (PIC) in a 6510 CPU simulation.
// Reset reinitializes the state of the PIC to default values.
// VerifyIrq determines and returns the type of interrupt request (IRQ) based on specified input conditions.
// ClearNMI clears the non-maskable interrupt (NMI) state within the PIC.
// HasNMI checks if a non-maskable interrupt (NMI) has been triggered and returns a boolean result.
type I6510Pic interface {
	Reset()

	VerifyIrq(uint8, uint8) uint8

	ClearNMI()

	HasNMI() bool
}
