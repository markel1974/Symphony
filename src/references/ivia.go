package references

// IVIASocket provides methods for interacting with a VIA socket, including reading, writing, and signaling operations.
// ReadPRA reads a value from Peripheral Register A (PRA) using a specified mask and shift.
// ReadPRB reads a value from Peripheral Register B (PRB) using a specified mask and shift.
// WritePRA writes a value to Peripheral Register A (PRA) using a specified mask and shift.
// WritePRB writes a value to Peripheral Register B (PRB) using a specified mask and shift.
// WriteDDRA writes a value to the Data Direction Register A (DDRA) using a specified mask and shift.
// WriteDDRB writes a value to the Data Direction Register B (DDRB) using a specified mask and shift.
// IRQClear clears the interrupt request (IRQ) signal.
// IRQTrigger triggers an interrupt request (IRQ) signal.
type IVIASocket interface {
	ReadPRA(uint8, uint8) uint8

	ReadPRB(uint8, uint8) uint8

	WritePRA(uint8, uint8)

	WritePRB(uint8, uint8)

	WriteDDRA(uint8, uint8)

	WriteDDRB(uint8, uint8)

	IRQClear()

	IRQTrigger()
}

const IdIVIA = "VIA"

// IVIA represents an interface for a VIA (Versatile Interface Adapter) component, managing communication and signaling.
// Setup initializes the IVia instance by associating it with a provided socket.
// Reset reinitializes the state of the IVia to its default operational state.
// Emulate performs an emulation cycle for the IVia.
// ReadByte retrieves a byte of data from the specified memory address.
// WriteByte writes a byte of data to the specified memory address.
// SignalPRA triggers the VIA PRA (Peripheral Register A) signal.
// SignalPRB triggers the VIA PRB (Peripheral Register B) signal.
// ByteReady checks if the VIA is ready to handle a new byte of data and returns true if ready.
type IVIA interface {
	Setup(conn IVIASocket) error

	Reset()

	Emulate()

	ReadByte(addr uint16) uint8

	WriteByte(addr uint16, data uint8)

	SignalPRA()

	SignalPRB()

	ByteReady() bool
}
