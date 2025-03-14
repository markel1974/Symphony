package board

import (
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

// CIA2Socket represents a connection interface to the CIA2 chip on a hardware board.
// It contains a reference to the board and an interrupt identifier.
type CIA2Socket struct {
	cia2   references.ICia
	board  *Board
	intrId uint32
}

// NewCIA2Socket creates and returns a new instance of CIA2Socket.
func NewCIA2Socket() *CIA2Socket {
	c := &CIA2Socket{
		cia2:   nil,
		board:  nil,
		intrId: intrIrqCia2Bit,
	}
	return c
}

// Setup initializes the CIA2Socket with the provided board reference and interrupt ID.
func (w *CIA2Socket) Setup(board *Board, c2 component.IComponent) error {
	cia2, ok := c2.(references.ICia)
	if !ok {
		return fmt.Errorf("unknown cia2 interface")
	}
	w.board = board
	w.cia2 = cia2
	w.cia2.Setup(w)
	return nil
}

func (w *CIA2Socket) Emulate() {
	w.cia2.Emulate()
}

// Reset resets the connected CIA2 component to its default state by invoking its Reset method.
func (w *CIA2Socket) Reset() {
	w.cia2.Reset()
}

// Update invokes the Update method on the CIA2 instance associated with the CIA2Socket.
func (w *CIA2Socket) Update() {
	w.cia2.Update()
}

// ReadPortA reads data from Port A considering the data direction register and input from the CPU's IEC interface.
func (w *CIA2Socket) ReadPortA(prA uint8, ddrA uint8, _ uint8, _ uint8) uint8 {
	data := w.board.iec.CpuRead()
	ret := ((prA | (^ddrA)) & 0x3f) | data
	return ret
}

// ReadPortB reads the value of Port B by combining the port register (prB) with the inverted data direction register (ddrB).
func (w *CIA2Socket) ReadPortB(_ uint8, _ uint8, prB uint8, ddrB uint8) uint8 {
	ret := prB | (^ddrB)
	return ret
}

// WritePortA writes data to Port A by updating its virtual address and signaling the connected IEC interface.
func (w *CIA2Socket) WritePortA(prA uint8, ddrA uint8, _ uint8, _ uint8) {
	w.updateVA(prA, ddrA)
	w.board.iec.CpuWrite(prA)
}

// WritePortB writes data to Port B with the given input parameters, modifying internal state as required.
func (w *CIA2Socket) WritePortB(_ uint8, _ uint8, _ uint8, _ uint8) {
}

// WriteDdrA updates the DDR register A and its associated state by calling the updateVA method.
func (w *CIA2Socket) WriteDdrA(prA uint8, ddrA uint8, _ uint8, _ uint8) {
	w.updateVA(prA, ddrA)
}

// WriteDdrB updates the Data Direction Register (DDR) for Port B with specified parameters, typically controlling pin direction settings.
func (w *CIA2Socket) WriteDdrB(_ uint8, _ uint8, _ uint8, _ uint8) {
}

// updateVA updates the VIC-memory bank selection based on the given port A and data direction register A values.
func (w *CIA2Socket) updateVA(prA uint8, ddrA uint8) {
	//Bit 0..1: Select the position of the VIC-memory
	//Bit 2: RS-232: TXD Output, userPort: Data PA 2 (pin M)
	//Bit 3..5: serial bus Output (0=High/Inactive, 1=Low/Active)
	//Bit 6..7: serial bus Input (0=Low/Active, 1=High/Inactive)

	//%00, 0: Bank 3: $C000-$FFFF, 49152-65535
	//%01, 1: Bank 2: $8000-$BFFF, 32768-49151
	//%10, 2: Bank 1: $4000-$7FFF, 16384-32767
	//%11, 3: Bank 0: $0000-$3FFF, 0-16383 (standard)
	va := (^(prA | (^ddrA))) & 3
	w.board.vicSocket.ChangedVA(va)
}

// IRQTrigger sends an interrupt request to the connected board's NMI trigger slot.
func (w *CIA2Socket) IRQTrigger() {
	w.board.nmiTriggerSlot()
}

// IRQClear clears the currently triggered IRQ (Interrupt Request) on the associated board's NMI slot.
func (w *CIA2Socket) IRQClear() {
	w.board.nmiClearSlot()
}
