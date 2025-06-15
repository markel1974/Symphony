package mos6522

import (
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

// see drive/iecieee/via2d.c [store_pra - store_prb]
// 1541, 1541II, 1571 and 2031
// see https://sta.c64.org/cbm1541mem.html

// defaultViaTimeout is the maximum threshold value used for timer underflow checks during VIA emulation.
const defaultViaTimeout = 0xffff

// VIA represents a versatile interface adapter used for I/O, timing, and control in a system.
type VIA struct {
	*component.BaseComponent
	pra          uint8
	ddra         uint8
	prb          uint8
	ddrb         uint8
	t1c          uint16
	t1l          uint16
	t2c          uint16
	t2l          uint16
	sr           uint8
	acr          uint8
	pcr          uint8
	ifr          uint8
	ier          uint8
	lastCA1      bool  // Stato precedente di CA1 per il rilevamento del fronte
	lastCB1      bool  // Stato precedente di CB1 per il rilevamento del fronte
	lastCB2      bool  // Stato precedente di CB2 per la modalità Shift Register
	shiftCounter uint8 // Contatore per gli 8 bit dello Shift Register
	socket       references.IMos6522Socket
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
	return v
}

func (v *VIA) Setup() error {
	return nil
}

func (v *VIA) Bind(socket references.IMos6522Socket) error {
	v.socket = socket
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
	v.t1c = 0
	v.t1l = 0
	v.t2c = 0
	v.t2l = 0
	v.sr = 0
	v.acr = 0
	v.pcr = 0
	v.ifr = 0
	v.ier = 0
	v.lastCA1 = false
	v.lastCB1 = false
	v.lastCB2 = false
	v.shiftCounter = 0
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
		//v.signalIRQClear.Emit(v.intrId) //intrVIA1Id)
		return uint8(v.t1c)
	case 0x5: // 0x1805 | 0x1c05
		return uint8(v.t1c >> 8)
	case 0x6: // 0x1806 | 0x1c06
		return uint8(v.t1l)
	case 0x7: // 0x1807 | 0x1c07
		return uint8(v.t1l >> 8)
	case 0x8: // 0x1808 | 0x1c08
		v.ifr &= 0xdf
		return uint8(v.t2c)
	case 0x9: // 0x1809 | 0x1c09
		return uint8(v.t2c >> 8)
	case 0xa: // 0x180a | 0x1c0a
		return v.sr
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
		v.t1l = (v.t1l & 0xff00) | uint16(data)
	case 0x5: //0x1805:
		v.t1l = (v.t1l & 0xff) | (uint16(data) << 8)
		v.ifr &= 0xbf
		v.t1c = v.t1l
	case 0x6: //0x1806:
		v.t1l = (v.t1l & 0xff00) | uint16(data)
	case 0x7: //0x1807:
		v.t1l = (v.t1l & 0xff) | (uint16(data) << 8)
	case 0x8: //0x1808:
		v.t2l = (v.t2l & 0xff00) | uint16(data)
	case 0x9: //0x1809:
		v.t2l = (v.t2l & 0xff) | (uint16(data) << 8)
		v.ifr &= 0xdf
		v.t2c = v.t2l
	case 0xa: //0x180a:
		v.sr = data
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
func (v *VIA) Emulate() {
	//v.handleHandshakeInput()
	// Timer 1
	t1c := uint(v.t1c) - 1
	v.t1c = uint16(t1c)
	if t1c > defaultViaTimeout {
		if (v.acr & 0x40) != 0 {
			v.t1c = v.t1l
		}
		v.ifr |= 0x40
		if (v.ier & 0x40) != 0 {
			v.socket.IRQTrigger()
		}
	}
	// Timer 2
	//t2Underflow := false
	if (v.acr & 0x20) == 0 {
		t2c := uint(v.t2c) - 1
		v.t2c = uint16(t2c)
		if t2c > defaultViaTimeout {
			v.ifr |= 0x20
			//t2Underflow = true
			if (v.ier & 0x20) != 0 {
				v.socket.IRQTrigger()
			}
		}
	}
	//v.handleShiftRegister(t2Underflow)
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

/*
// handleShiftRegister gestisce la logica completa dello Shift Register.
func (v *VIA) handleShiftRegister(t2Clock bool) {
	srMode := (v.acr >> 2) & 0x07
	if srMode == 0 {
		return
	}
	shiftTrigger := false
	switch srMode {
	case 1, 4, 5: // Modalità con clock da Timer 2
		shiftTrigger = t2Clock
	case 2, 6: // Modalità con clock da Phase 2
		shiftTrigger = true
	case 3, 7: // Modalità con clock esterno da CB1
		cb1 := v.socket.ReadCB1()
		if cb1 && !v.lastCB1 {
			shiftTrigger = true
		}
		v.lastCB1 = cb1
	}
	if !shiftTrigger {
		return
	}
	isShiftIn := (srMode & 0x04) == 0
	if v.shiftCounter < 8 {
		if isShiftIn {
			inBit := uint8(0)
			if v.socket.ReadCB2() {
				inBit = 1
			}
			v.sr = (v.sr << 1) | inBit
		} else {
			outBit := (v.sr & 0x80) != 0
			v.socket.WriteCB2(outBit)
			v.sr <<= 1
		}
		v.shiftCounter++
	}
	if v.shiftCounter == 8 {
		v.shiftCounter = 9
		v.ifr |= 0x04
		if (v.ier & 0x04) != 0 {
			v.socket.IRQTrigger()
		}
	}
}
*/
