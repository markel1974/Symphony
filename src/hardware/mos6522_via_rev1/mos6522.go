package mos6522

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

const (
	bit0 = 0x01
	bit2 = 0x04
	bit1 = 0x02
	bit4 = 0x10
	bit5 = 0x20
	bit6 = 0x40
	bit7 = 0x80
)

const (
	RegisterSize  = 0xf
	RegisterCount = RegisterSize + 1
)

// see https://sta.c64.org/cbm1541mem.html

// VIA represents a versatile interface adapter used for I/O, timing, and control in a system.
type VIA struct {
	*component.BaseComponent
	pra           uint8
	ddra          uint8
	prb           uint8
	ddrb          uint8
	timer0        *Timer
	timer1        *Timer
	shiftRegister *ShiftRegister
	acr           uint8
	pcr           uint8
	ifr           uint8
	ier           uint8
	lastCA1       bool
	lastCB2       bool
	lastPB6       bool
	lastCB1       bool
	//socket        references.IMos6522Socket
	reads  [RegisterCount]func() uint8
	writes [RegisterCount]func(uint8)

	socketReadPortA       func() uint8
	socketReadPortB       func() uint8
	socketReadPB6         func() bool
	socketReadCA1         func() bool
	socketReadCB1         func() bool
	socketSignalPRA       func(uint8)
	socketSignalPRB       func(uint8)
	socketSignalDDRA      func(uint8)
	socketSignalDDRB      func(uint8)
	socketSignalPCR       func(uint8)
	socketIRQTrigger      func()
	socketIRQClearTrigger func()
}

// NewVIA initializes and returns a new instance of the VIA type, associating it with a parent component and factory.
// It sets the VIA's internal registers to their default values and registers the component within its hierarchy.
func NewVIA(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *VIA {
	v := &VIA{
		BaseComponent: component.NewBaseComponent(),
		pra:           0,
		prb:           0,
		ddra:          0,
		ddrb:          0,
	}
	v.BaseComponent.Register(factory, parent, Identifier(), v, references.IdIMos6522(v, label, instance))
	v.timer0 = NewTimer(v, v.GetFactory(), label, 0)
	v.timer1 = NewTimer(v, v.GetFactory(), label, 1)
	v.shiftRegister = NewShiftRegister(v, v.GetFactory(), label, 0)
	return v
}

// Setup initializes the VIA component, preparing it for operation and ensuring all dependencies are properly configured.
func (v *VIA) Setup() error {
	return nil
}

// Bind assigns an IMos6522Socket to the VIA and initializes its shift register with the socket's CB2 read/write callbacks.
func (v *VIA) Bind(socket references.IMos6522Socket) error {
	v.socketIRQTrigger = socket.IRQTrigger
	v.socketIRQClearTrigger = socket.IRQClearTrigger
	v.socketReadPortA = socket.ReadPortA
	v.socketReadPortB = socket.ReadPortB
	v.socketReadPB6 = socket.ReadPB6
	v.socketReadCA1 = socket.ReadCA1
	v.socketReadCB1 = socket.ReadCB1
	v.socketSignalPRA = socket.SignalPRA
	v.socketSignalPRB = socket.SignalPRB
	v.socketSignalDDRA = socket.SignalDDRA
	v.socketSignalDDRB = socket.SignalDDRB
	v.socketSignalPCR = socket.SignalPCR
	v.shiftRegister.Initialize(socket.ReadCB2, socket.WriteCB2)
	v.reads = v.createReadRegister()
	v.writes = v.createWriteRegister()
	return nil
}

// Connect establishes a connection or initializes state for the VIA component, returning an error if any issue occurs.
func (v *VIA) Connect() error {
	return nil
}

// Internal determines if the VIA component should operate in internal mode, returning false in its current implementation.
func (v *VIA) Internal() bool {
	return false
}

// Reset initializes all VIA internal registers to their default state of zero.
func (v *VIA) Reset() {
	v.pra = 0
	v.ddra = 0
	v.prb = 0
	v.ddrb = 0
	v.timer0.Reset()
	v.timer1.Reset()
	v.shiftRegister.Reset()
	v.acr = 0
	v.pcr = 0
	v.ifr = 0
	v.ier = 0
	v.lastCA1 = false
	v.lastCB1 = false
	v.lastCB2 = false
	v.lastPB6 = false
}

// ReadByte reads a byte from the specified VIA register address and returns the corresponding value based on its state.
func (v *VIA) ReadByte(addr uint16) uint8 {
	//0x1800 | 0x1c00
	reg := addr & RegisterSize
	return v.reads[reg]()
}

// WriteByte writes a byte of data to the specified register address in the VIA and updates the internal state accordingly.
func (v *VIA) WriteByte(addr uint16, data uint8) {
	//0x1800 | 0x1c00
	reg := addr & RegisterSize
	v.writes[reg](data)
}

// Emulate executes a single emulation cycle for VIA, decrementing timers and handling interrupts based on current settings.
//
//go:nosplit
func (v *VIA) Emulate() {
	v.handleHandshakeInput()

	t2ClockPulse := false

	if (v.acr & bit5) == 0 {
		t2ClockPulse = true
	} else {
		currentPB6 := v.socketReadPB6()
		if v.lastPB6 && !currentPB6 {
			t2ClockPulse = true
		}
		v.lastPB6 = currentPB6
	}
	v.timer1.SetClockPulse(t2ClockPulse)

	v.timer0.Emulate()
	v.timer1.Emulate()

	if v.timer0.Underflow() {
		v.ifr |= bit6
		if (v.acr & bit6) != 0 {
			v.timer0.Load()
		}
		if (v.ier & bit6) != 0 {
			v.socketIRQTrigger()
		}
	}
	t2Underflow := v.timer1.Underflow()
	if t2Underflow {
		v.ifr |= bit5
		if (v.ier & bit5) != 0 {
			v.socketIRQTrigger()
		}
	}
	if srMode := (v.acr >> 2) & 0x07; srMode != 0 {
		v.handleShiftRegister(srMode, t2Underflow)
	}
}

// EmulationRequired indicates whether emulation of the VIA functionality is required, always returning true.
func (v *VIA) EmulationRequired() bool {
	return true
}

// ReadDDRA returns the current value of the Data Direction Register A (DDRA) of the VIA.
func (v *VIA) ReadDDRA() uint8 {
	return v.ddra
}

// ReadDDRB returns the current value of the Data Direction Register A (DDRB) of the VIA.
func (v *VIA) ReadDDRB() uint8 {
	return v.ddrb
}

// ReadPRA returns the current value of the Peripheral Register A (PRA) from the VIA.
func (v *VIA) ReadPRA() uint8 {
	return v.pra
}

// ReadPRB returns the current value of the Peripheral Register B (PRB) from the VIA.
func (v *VIA) ReadPRB() uint8 {
	return v.prb
}

// ReadACR returns the current value of the Auxiliary Control Register (ACR) from the VIA.
func (v *VIA) ReadACR() uint8 {
	return v.acr
}

// ReadPCR retrieves the current value of the Peripheral Control Register (PCR) from the VIA.
func (v *VIA) ReadPCR() uint8 {
	return v.pcr
}

func (v *VIA) createWriteRegister() [RegisterCount]func(uint8) {
	var writes [RegisterCount]func(uint8)
	writes[0x0] = func(data uint8) {
		v.prb = data
		v.socketSignalPRB(v.prb)
	}
	writes[0x1] = func(data uint8) {
		v.pra = data
		v.socketSignalPRA(v.pra)
	}
	writes[0x2] = func(data uint8) {
		v.ddrb = data
		v.socketSignalDDRB(v.ddrb)
	}
	writes[0x3] = func(data uint8) {
		v.ddra = data
		v.socketSignalDDRA(v.ddra)
	}
	writes[0x4] = v.timer0.SetLatchLow
	writes[0x5] = func(data uint8) {
		v.timer0.SetLatchHigh(data)
		v.ifr &= 0xbf
		v.timer0.Load()
	}
	writes[0x6] = v.timer0.SetLatchLow
	writes[0x7] = v.timer0.SetLatchHigh
	writes[0x8] = v.timer1.SetLatchLow
	writes[0x9] = func(data uint8) {
		v.timer1.SetLatchHigh(data)
		v.ifr &= 0xdf
		v.timer1.Load()
	}
	writes[0xa] = v.shiftRegister.Set
	writes[0xb] = func(data uint8) {
		v.acr = data
	}
	writes[0xc] = func(data uint8) {
		v.pcr = data
		v.socketSignalPCR(v.pcr)
	}
	writes[0xd] = func(data uint8) {
		v.ifr &= ^data
	}
	writes[0xe] = func(data uint8) {
		if (data & bit7) != 0 {
			v.ier |= data & 0x7f
		} else {
			v.ier &= ^data
		}
	}
	writes[0xf] = func(data uint8) {
		v.pra = data
		v.socketSignalPRA(v.pra)
	}
	return writes
}

func (v *VIA) createReadRegister() [RegisterCount]func() uint8 {
	var reads [RegisterCount]func() uint8
	reads[0x0] = v.socketReadPortB
	reads[0x1] = v.socketReadPortA
	reads[0x2] = v.ReadDDRB
	reads[0x3] = v.ReadDDRA
	reads[0x4] = func() uint8 {
		v.ifr &= 0xbf
		v.socketIRQClearTrigger()
		return v.timer0.CounterLow()
	}
	reads[0x5] = v.timer0.CounterHigh
	reads[0x6] = v.timer0.LatchLow
	reads[0x7] = v.timer0.LatchHigh
	reads[0x8] = func() uint8 {
		v.ifr &= 0xdf
		return v.timer1.CounterLow()
	}
	reads[0x9] = v.timer1.CounterHigh
	reads[0xa] = v.shiftRegister.Get
	reads[0xb] = v.ReadACR
	reads[0xc] = v.ReadPCR
	reads[0xd] = func() uint8 {
		if (v.ifr & v.ier) != 0 {
			return v.ifr | bit7
		}
		return v.ifr
	}
	reads[0xe] = func() uint8 {
		return v.ier | bit7
	}
	reads[0xf] = func() uint8 {
		return v.socketReadPortA()
	}
	return reads
}

// handleShiftRegister processes the serial shift register based on the mode and trigger conditions, updating the interrupt flags if necessary.
func (v *VIA) handleShiftRegister(srMode uint8, t2Underflow bool) {
	enabled := false
	switch srMode {
	case 1, 4, 5: // Clock mode with Timer 2
		enabled = t2Underflow
	case 2, 6: // Clock mode with Phase 2
		enabled = true
	case 3, 7: // Clock mode with external clock from CB1
		cb1 := v.socketReadCB1()
		if cb1 && !v.lastCB1 {
			enabled = true
		}
		v.lastCB1 = cb1
	}
	if enabled {
		shiftIn := (srMode & bit2) == 0
		if v.shiftRegister.Handle(shiftIn) {
			v.ifr |= bit2
			if (v.ier & bit2) != 0 {
				v.socketIRQTrigger()
			}
		}
	}
}

// handleHandshakeInput manages edge detection for control pins CA1 and CB1, updates interrupt flags, and handles port latching.
func (v *VIA) handleHandshakeInput() {
	// Port A
	// Reads the configuration from PCR: true for falling edge (bit 0 = 0), false for rising edge (bit 0 = 1).
	ca1DetectsFallingEdge := (v.pcr & bit0) == 0
	currentCA1 := v.socketReadCA1()

	// Detect the edge by comparing the current state with the previous cycle state.
	ca1EdgeDetected := (ca1DetectsFallingEdge && !currentCA1 && v.lastCA1) ||
		(!ca1DetectsFallingEdge && currentCA1 && !v.lastCA1)

	if ca1EdgeDetected {
		// Set the interrupt flag for CA1 (bit 1 of IFR).
		v.ifr |= bit1
		// If latching on Port A is enabled (bit 0 of ACR = 0), "freeze" the port value.
		if (v.acr & bit0) == 0 {
			v.pra = v.socketReadPortA()
		}
		// If the CA1 interrupt is enabled (bit 1 of IER), trigger the IRQ.
		if (v.ier & bit1) != 0 {
			v.socketIRQTrigger()
		}
	}
	v.lastCA1 = currentCA1

	// Port B
	cb1DetectsFallingEdge := (v.pcr & bit4) == 0
	currentCB1 := v.socketReadCB1()

	// Detect the edge by comparing the current state with the previous cycle state.
	cb1EdgeDetected := (cb1DetectsFallingEdge && !currentCB1 && v.lastCB1) ||
		(!cb1DetectsFallingEdge && currentCB1 && !v.lastCB1)

	if cb1EdgeDetected {
		// Sets the interrupt flag for CB1 (bit 4 of IFR).
		v.ifr |= bit4
		// If Port B latching is enabled (bit 4 of ACR = 0), "freeze" the port value.
		if (v.acr & bit4) == 0 {
			v.prb = v.socketReadPortB()
		}
		// If CB1 interrupt is enabled (bit 4 of IER), trigger the IRQ.
		if (v.ier & bit4) != 0 {
			v.socketIRQTrigger()
		}
	}
	// Update CB1 state for the next cycle.
	v.lastCB1 = currentCB1
}
