package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// CIA2Socket represents a specialized implementation of ICIA, handling interactions with CIA2 and emulation components.
// The connections field manages CIA2-specific port and communication events.
// The iec field facilitates communication over the IEC serial bus.
// The intrId field defines the unique interrupt identifier for CIA2 within the system.
type CIA2Socket struct {
	references.ICIA
	pic    references.IPIC6510
	vic    references.IVIC
	iec    references.IIec
	intrId uint32
}

// NewCIA2Socket creates and returns a pointer to a new instance of CIA2Socket with default uninitialized fields.
func NewCIA2Socket() *CIA2Socket {
	c := &CIA2Socket{
		ICIA:   nil,
		pic:    nil,
		vic:    nil,
		iec:    nil,
		intrId: intrIrqCia2Bit,
	}
	return c
}

// Setup initializes the CIA2Socket instance with the provided CIA, connections, and IEC interface, and sets up the CIA.
func (w *CIA2Socket) Setup(c map[string]references.IComponent, cfg *config.Config) error {
	var err error
	w.ICIA, err = references.ComponentsToICIA(c, 1)
	if err != nil {
		return err
	}
	w.pic, err = references.ComponentsToIPIC6510(c, 0)
	if err != nil {
		return err
	}
	w.vic, err = references.ComponentsToIVIC(c, 0)
	if err != nil {
		return err
	}
	w.iec, err = references.ComponentsToIEC(c, 0)
	if err = w.ICIA.Setup(w, cfg); err != nil {
		return err
	}
	return nil
}

// Connect initializes and establishes a connection using the ICIA interface. Returns an error if the operation fails.
func (w *CIA2Socket) Connect() error {
	if err := w.ICIA.Connect(); err != nil {
		return err
	}
	return nil
}

// ReadPortA reads the value of port A by combining peripheral data, data direction bits, and the IEC CpuRead result.
func (w *CIA2Socket) ReadPortA(prA uint8, ddrA uint8, _ uint8, _ uint8) uint8 {
	data := w.iec.CpuRead()
	ret := ((prA | (^ddrA)) & 0x3f) | data
	return ret
}

// ReadPortB reads the current state of Port B by combining the peripheral register and inverted direction register.
func (w *CIA2Socket) ReadPortB(_ uint8, _ uint8, prB uint8, ddrB uint8) uint8 {
	ret := prB | (^ddrB)
	return ret
}

// WritePortA writes the state of Port A using the given peripheral and direction register values.
func (w *CIA2Socket) WritePortA(prA uint8, ddrA uint8, _ uint8, _ uint8) {
	w.updateVA(prA, ddrA)
	w.iec.CpuWrite(prA)
}

// WritePortB updates the state of port B using the provided peripheral and direction registers.
func (w *CIA2Socket) WritePortB(_ uint8, _ uint8, _ uint8, _ uint8) {
}

// WriteDdrA updates the data direction register for port A and triggers actions based on the updated state.
func (w *CIA2Socket) WriteDdrA(prA uint8, ddrA uint8, _ uint8, _ uint8) {
	w.updateVA(prA, ddrA)
}

// WriteDdrB updates the data direction register for port B with the provided parameters.
func (w *CIA2Socket) WriteDdrB(_ uint8, _ uint8, _ uint8, _ uint8) {
}

// updateVA updates the VIC-memory bank based on the current states of prA and ddrA and triggers a corresponding VA change event.
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
	w.vic.ChangedVA(va)
}

// IRQTrigger triggers a Non-Maskable Interrupt (NMI) by invoking the NMITrigger method on the connected socket.
func (w *CIA2Socket) IRQTrigger() {
	w.pic.TriggerNMI()
}

// IRQClear clears the NMI (Non-Maskable Interrupt) request by invoking the NMIClear method on the connections object.
func (w *CIA2Socket) IRQClear() {
	w.pic.ClearNMI()
}
