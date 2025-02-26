package mos6526

// ISocket defines an interface for interacting with external device ports and managing IRQ operations.
// ReadPortA reads data from Port A using the provided peripheral and data direction registers.
// ReadPortB reads data from Port B using the provided peripheral and data direction registers.
// WritePortA writes data to Port A using the provided peripheral and data direction registers.
// WritePortB writes data to Port B using the provided peripheral and data direction registers.
// WriteDdrA modifies the Data Direction Register for Port A using the provided register values.
// WriteDdrB modifies the Data Direction Register for Port B using the provided register values.
// IRQTrigger signals an interrupt request activation.
// IRQClear signals an interrupt request deactivation.
type ISocket interface {
	ReadPortA(prA uint8, ddrA uint8, prB uint8, ddrB uint8) uint8
	ReadPortB(prA uint8, ddrA uint8, prB uint8, ddrB uint8) uint8
	WritePortA(prA uint8, ddrA uint8, prB uint8, ddrB uint8)
	WritePortB(prA uint8, ddrA uint8, prB uint8, ddrB uint8)
	WriteDdrA(prA uint8, ddrA uint8, prB uint8, ddrB uint8)
	WriteDdrB(prA uint8, ddrA uint8, prB uint8, ddrB uint8)

	IRQTrigger()
	IRQClear()
}
