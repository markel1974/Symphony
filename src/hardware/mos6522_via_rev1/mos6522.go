package mos6522

import (
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

// see drive/iecieee/via2d.c [store_pra - store_prb]
// 1541, 1541II, 1571 and 2031
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
	socket        references.IMos6522Socket
}

// NewVIA initializes and returns a new instance of the VIA type, associating it with a parent component and factory.
// It sets the VIA's internal registers to their default values and registers the component within its hierarchy.
func NewVIA(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *VIA {
	v := &VIA{
		BaseComponent: component.NewBaseComponent(),
		pra:           0,
		ddra:          0,
		prb:           0,
		ddrb:          0,
	}
	v.BaseComponent.Register(factory, parent, Identifier(), v, references.IdIMos6522(v, label, instance))
	v.timer0 = NewTimer(v, v.GetFactory(), label, 0)
	v.timer1 = NewTimer(v, v.GetFactory(), label, 1)
	v.shiftRegister = NewShiftRegister(v, v.GetFactory(), label, 0)
	return v
}

func (v *VIA) Setup() error {
	return nil
}

func (v *VIA) Bind(socket references.IMos6522Socket) error {
	v.socket = socket
	v.shiftRegister.Initialize(v.socket.ReadCB2, v.socket.WriteCB2)
	return nil
}

// Connect establishes a connection or initializes state for the VIA component, returning an error if any issue occurs.
func (v *VIA) Connect() error {
	return nil
}

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
	reg := addr & 0x0f
	switch reg {
	case 0x0: //0x1800 | 0x1c00
		return v.socket.ReadPRB(v.prb, v.ddrb)
	case 0x1: //0x1801 | 0x1c01
		return v.socket.ReadPRA(v.pra, v.ddra)
	case 0x2: //0x1802 | 0x1c02
		return v.ddrb
	case 0x3: //0x1803 | 0x1c03
		return v.ddra
	case 0x4: //0x1804 | 0x1c04
		v.ifr &= 0xbf
		v.socket.IRQClearTrigger()
		return v.timer0.CounterLow()
	case 0x5: // 0x1805 | 0x1c05
		return v.timer0.CounterHigh() //uint8(v.timer0.counter >> 8)
	case 0x6: // 0x1806 | 0x1c06
		return v.timer0.LatchLow() //uint8(v.timer0.latch)
	case 0x7: // 0x1807 | 0x1c07
		return v.timer0.LatchHigh() //uint8(v.timer0.latch >> 8)
	case 0x8: // 0x1808 | 0x1c08
		v.ifr &= 0xdf
		return v.timer1.CounterLow() //uint8(v.t2c)
	case 0x9: // 0x1809 | 0x1c09
		return v.timer1.CounterHigh() //uint8(v.t2c >> 8)
	case 0xa: // 0x180a | 0x1c0a
		return v.shiftRegister.Get()
	case 0xb: // 0x180b | 0x1c0b
		return v.acr
	case 0xc: // 0x180c | 0x1c0c
		return v.pcr
	case 0xd: // 0x180d | 0x1c0d
		if (v.ifr & v.ier) != 0 {
			return v.ifr | 0x80
		}
		return v.ifr
	case 0xe: // 0x180e | 0x1c0e
		return v.ier | 0x80
	case 0xf: // 0x180f | 0x1c0f
		//TODO Implement not ok!
		return v.socket.ReadPRA(v.pra, v.ddra)
	default:
		fmt.Printf("%s READ UNKNOWN - %x\n", v.GetId(), addr)
		return 0
	}
}

// WriteByte writes a byte of data to the specified register address in the VIA and updates the internal state accordingly.
func (v *VIA) WriteByte(addr uint16, data uint8) {
	reg := addr & 0x0f
	switch reg {
	case 0x0: //0x1800:
		v.prb = data
		v.socket.WritePRB(v.prb, v.ddrb)
	case 0x1: //0x1801:
		v.pra = data
		v.socket.WritePRA(v.pra, v.ddra)
	case 0x2: //0x1802:
		v.ddrb = data
		v.socket.WriteDDRB(v.prb, v.ddrb)
	case 0x3: //0x1803:
		v.ddra = data
		v.socket.WriteDDRA(v.pra, v.ddra)
	case 0x4: //0x1804:
		v.timer0.SetLatchLow(data) //v.t1l = (v.t1l & 0xff00) | uint16(data)
	case 0x5: //0x1805:
		v.timer0.SetLatchHigh(data) //v.t1l = (v.t1l & 0xff) | (uint16(data) << 8)
		v.ifr &= 0xbf
		v.timer0.Load() //v.t1c = v.t1l
	case 0x6: //0x1806:
		v.timer0.SetLatchLow(data) //v.t1l = (v.t1l & 0xff00) | uint16(data)
	case 0x7: //0x1807:
		v.timer0.SetLatchHigh(data) //v.t1l = (v.t1l & 0xff) | (uint16(data) << 8)
	case 0x8: //0x1808:
		v.timer1.SetLatchLow(data) //v.t2l = (v.t2l & 0xff00) | uint16(data)
	case 0x9: //0x1809:
		v.timer1.SetLatchHigh(data) // v.t2l = (v.t2l & 0xff) | (uint16(data) << 8)
		v.ifr &= 0xdf
		v.timer1.Load() //v.t2c = v.t2l
	case 0xa: //0x180a:
		v.shiftRegister.Set(data)
	case 0xb: //0x180b:
		v.acr = data
	case 0xc: //0x180c || 0x1c0c
		//Bits #1-#3: %111 = Attach Byte Ready line to overflow processor flag.
		//Whenever a data byte has been successfully read from or written to disk, V flag is set to 1.
		//Bits #5-#7: Head control; %111 = Read (0xE0); %110 = Write (0xC0).
		v.pcr = data
		headControl := v.pcr & 0xE0
		if headControl == 0xC0 {
			v.socket.WriteCA2(true)
		} else {
			v.socket.WriteCA2(false)
		}
		//NOT CONNECTED
		//if (v.pcr & 0x0E) == 0x0C {
		//	v.socket.WriteCB2(false)
		//} else if (v.pcr & 0x0E) == 0x0E {
		//	v.socket.WriteCB2(true)
		//}
	case 0xd: //0x180d:
		v.ifr &= ^data
	case 0xe: //0x180e:
		if (data & 0x80) != 0 {
			v.ier |= data & 0x7f
		} else {
			v.ier &= ^data
		}
	case 0xf: //0x180f:
		v.pra = data
		v.socket.WritePRA(v.pra, v.ddra)
	default:
		fmt.Printf("%s WRITE UNKNOWN - %x\n", v.GetId(), addr)
	}
}

// Emulate executes a single emulation cycle for VIA, decrementing timers and handling interrupts based on current settings.
//
//go:nosplit
func (v *VIA) Emulate() {
	v.handleHandshakeInput()

	t2ClockPulse := false

	if (v.acr & 0x20) == 0 {
		t2ClockPulse = true
	} else {
		currentPB6 := v.socket.ReadPB6()
		if v.lastPB6 && !currentPB6 {
			t2ClockPulse = true
		}
		v.lastPB6 = currentPB6
	}
	v.timer1.SetClockPulse(t2ClockPulse)

	v.timer0.Emulate()
	v.timer1.Emulate()

	if v.timer0.Underflow() {
		v.ifr |= 0x40
		if (v.acr & 0x40) != 0 {
			v.timer0.Load()
		}
		if (v.ier & 0x40) != 0 {
			v.socket.IRQTrigger()
		}
	}
	t2Underflow := v.timer1.Underflow()
	if t2Underflow {
		v.ifr |= 0x20
		if (v.ier & 0x20) != 0 {
			v.socket.IRQTrigger()
		}
	}
	if srMode := (v.acr >> 2) & 0x07; srMode != 0 {
		shiftTrigger := false
		switch srMode {
		case 1, 4, 5: // Clock mode with Timer 2
			shiftTrigger = t2Underflow
		case 2, 6: // Clock mode with Phase 2
			shiftTrigger = true
		case 3, 7: // Clock mode with external clock from CB1
			cb1 := v.socket.ReadCB1()
			if cb1 && !v.lastCB1 {
				shiftTrigger = true
			}
			v.lastCB1 = cb1
		}
		if shiftTrigger {
			isShiftIn := (srMode & 0x04) == 0
			if v.shiftRegister.Handle(isShiftIn) {
				v.ifr |= 0x04
				if (v.ier & 0x04) != 0 {
					v.socket.IRQTrigger()
				}
			}
		}
	}
}

// EmulationRequired indicates whether emulation of the VIA functionality is required, always returning true.
func (v *VIA) EmulationRequired() bool {
	return true
}

// SignalPRA writes the current values of PRA and DDRA to the connected socket via the WritePRA method.
func (v *VIA) SignalPRA() {
	v.socket.WritePRA(v.pra, v.ddra)
}

// SignalPRB sends the current contents of the PRB and DDRB registers to the connected IMos6522Socket.
func (v *VIA) SignalPRB() {
	v.socket.WritePRB(v.prb, v.ddrb)
}

// handleHandshakeInput manages edge detection for control pins CA1 and CB1, updates interrupt flags, and handles port latching.
func (v *VIA) handleHandshakeInput() {
	// Port A
	// Reads the configuration from PCR: true for falling edge (bit 0 = 0), false for rising edge (bit 0 = 1).
	ca1DetectsFallingEdge := (v.pcr & 0x01) == 0
	currentCA1 := v.socket.ReadCA1()

	// Detect the edge by comparing the current state with the previous cycle state.
	ca1EdgeDetected := (ca1DetectsFallingEdge && !currentCA1 && v.lastCA1) ||
		(!ca1DetectsFallingEdge && currentCA1 && !v.lastCA1)

	if ca1EdgeDetected {
		// Set the interrupt flag for CA1 (bit 1 of IFR).
		v.ifr |= 0x02
		// If latching on Port A is enabled (bit 0 of ACR = 0), "freeze" the port value.
		if (v.acr & 0x01) == 0 {
			v.pra = v.socket.ReadPRA(v.pra, v.ddra)
		}
		// If the CA1 interrupt is enabled (bit 1 of IER), trigger the IRQ.
		if (v.ier & 0x02) != 0 {
			v.socket.IRQTrigger()
		}
	}
	v.lastCA1 = currentCA1

	// Port B
	cb1DetectsFallingEdge := (v.pcr & 0x10) == 0
	currentCB1 := v.socket.ReadCB1()

	// Detect the edge by comparing the current state with the previous cycle state.
	cb1EdgeDetected := (cb1DetectsFallingEdge && !currentCB1 && v.lastCB1) ||
		(!cb1DetectsFallingEdge && currentCB1 && !v.lastCB1)

	if cb1EdgeDetected {
		// Sets the interrupt flag for CB1 (bit 4 of IFR).
		v.ifr |= 0x10
		// If Port B latching is enabled (bit 4 of ACR = 0), "freeze" the port value.
		if (v.acr & 0x10) == 0 {
			v.prb = v.socket.ReadPRB(v.prb, v.ddrb)
		}
		// If CB1 interrupt is enabled (bit 4 of IER), trigger the IRQ.
		if (v.ier & 0x10) != 0 {
			v.socket.IRQTrigger()
		}
	}
	// Update CB1 state for the next cycle.
	v.lastCB1 = currentCB1
}
