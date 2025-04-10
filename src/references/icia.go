package references

import (
	"fmt"
)

// ICIASocket defines an interface for interacting with CIA I/O ports and triggering IRQ events.
// ReadPortA reads the current value of port A based on the peripheral and direction registers.
// ReadPortB reads the current value of port B based on the peripheral and direction registers.
// WritePortA updates the state of port A based on the peripheral and direction registers.
// WritePortB updates the state of port B based on the peripheral and direction registers.
// WriteDdrA updates the data direction register for port A.
// WriteDdrB updates the data direction register for port B.
// IRQTrigger triggers an interrupt request on the IRQ line.
// IRQClear clears the interrupt request on the IRQ line.
type ICIASocket interface {
	ReadPortA(prA uint8, ddrA uint8, prB uint8, ddrB uint8) uint8

	ReadPortB(prA uint8, ddrA uint8, prB uint8, ddrB uint8) uint8

	WritePortA(prA uint8, ddrA uint8, prB uint8, ddrB uint8)

	WritePortB(prA uint8, ddrA uint8, prB uint8, ddrB uint8)

	WriteDdrA(prA uint8, ddrA uint8, prB uint8, ddrB uint8)

	WriteDdrB(prA uint8, ddrA uint8, prB uint8, ddrB uint8)

	IRQTrigger()

	IRQClear()
}

func IdICIA(_ ICIA, label string, instance int) string {
	return IdInternalComponent(label, instance, "ICIA")
}

// ICIA defines the interface for a CIA (Complex Interface Adapter) component in a computing or emulation context.
// Setup initializes the CIA with the provided socket connection.
// Reset resets the CIA to its default state.
// Emulate executes the necessary operations for the current emulation cycle.
// Update updates the CIA's state and handles any required operations for the current frame or step.
// WriteRegister writes a byte of data to the CIA at the specified register address.
// ReadRegister reads a byte of data from the CIA at the specified register address.
type ICIA interface {
	Setup() error

	Bind(conn ICIASocket) error

	Connect() error

	Reset()

	Emulate()

	Update()

	WriteRegister(addr uint16, data uint8)

	ReadRegister(addr uint16) uint8
}

func ComponentToICIA(component IComponent) (ICIA, error) {
	if component == nil {
		return nil, fmt.Errorf("component ICIA is nil")
	}
	v, ok := component.(ICIA)
	if !ok {
		return nil, fmt.Errorf("component is not a ICIA")
	}
	return v, nil
}
