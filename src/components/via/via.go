package mos6522

import (
	"fmt"
	"github.com/markel1974/c64emu/src/signals"
)

// see drive/iecieee/via2d.c [store_pra - store_prb]
// 1541, 1541II, 1571 and 2031
// see https://sta.c64.org/cbm1541mem.html

const defaultViaTimeout = 0xffff

type Via struct {
	pra              uint8
	ddra             uint8
	prb              uint8
	ddrb             uint8
	t1c              uint16
	t1l              uint16
	t2c              uint16
	t2l              uint16
	sr               uint8
	acr              uint8
	pcr              uint8
	ifr              uint8
	ier              uint8
	id               string
	intrId           uint32
	wiring           IWiring
	signalIRQTrigger *signals.SignalUint32
	signalIRQClear   *signals.SignalUint32
}

func NewVia(id string, intrId uint32) *Via {
	v := &Via{
		id:               id,
		intrId:           intrId,
		signalIRQTrigger: signals.NewSignalUint32(),
		signalIRQClear:   signals.NewSignalUint32(),
	}
	return v
}

func (v *Via) Setup(wiring IWiring) {
	v.wiring = wiring
}

func (v *Via) Reset() {
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

	v.wiring.Reset()
}

func (v *Via) SignalTriggerIRQBind(fn func(uint32)) {
	v.signalIRQTrigger.Bind(fn)
}

func (v *Via) SignalClearIRQBind(fn func(uint32)) {
	v.signalIRQClear.Bind(fn)
}

func (v *Via) ReadByte(addr uint16) uint8 {
	reg := addr & 0x0f
	switch reg {
	case 0x0: //0x1800:
		return v.wiring.ReadPRB(v.prb, v.ddrb)
	case 0x1: //0x1801:
		return v.wiring.ReadPRA(v.pra, v.ddra)
	case 0x2: //0x1802:
		return v.ddrb
	case 0x3: //0x1803:
		return v.ddra
	case 0x4: //x1804:
		v.ifr &= 0xbf
		v.signalIRQClear.Emit(v.intrId) //intrVIA1Id)
		return uint8(v.t1c)
	case 0x5: //0x1805:
		return uint8(v.t1c >> 8)
	case 0x6: //0x1806:
		return uint8(v.t1l)
	case 0x7: //0x1807:
		return uint8(v.t1l >> 8)
	case 0x8: //0x1808:
		v.ifr &= 0xdf
		return uint8(v.t2c)
	case 0x9: //0x1809:
		return uint8(v.t2c >> 8)
	case 0xa: //0x180a:
		return v.sr
	case 0xb: ////0x180b:
		return v.acr
	case 0xc: //0x180c:
		return v.pcr
	case 0xd: //0x180d:
		if (v.ifr & v.ier) != 0 {
			return v.ifr | 0x80
		}
		return v.ifr
	case 0xe: //0x180e:
		return v.ier | 0x80
	case 0xf: //0x180f:
		return v.wiring.ReadPRA(v.pra, v.ddra)
	default:
		fmt.Printf("%s READ UNKNOWN - %x\n", v.id, addr)
		return 0
	}
}

func (v *Via) WriteByte(addr uint16, data uint8) {
	reg := addr & 0x0f
	switch reg {
	case 0x0: //0x1800:
		v.prb = data
		v.wiring.WritePRB(v.prb, v.ddrb)
	case 0x1: //0x1801:
		v.pra = data
		v.wiring.WritePRA(v.pra, v.ddra)
	case 0x2: //0x1802:
		v.ddrb = data
		v.wiring.WriteDDRB(v.prb, v.ddrb)
	case 0x3: //0x1803:
		v.ddra = data
		v.wiring.WriteDDRA(v.pra, v.ddra)
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
	case 0xc: //0x180c:
		v.pcr = data
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
		v.wiring.WritePRA(v.pra, v.ddra)
	default:
		fmt.Printf("%s WRITE UNKNOWN - %x\n", v.id, addr)
	}
}

func (v *Via) Emulate() {
	t1c := uint(v.t1c) - 1
	v.t1c = uint16(t1c)

	if t1c > defaultViaTimeout {
		if (v.acr & 0x40) != 0 {
			// Reload from latch in free-run mode
			v.t1c = v.t1l
		}
		v.ifr |= 0x40
		if (v.ier & 0x40) != 0 {
			v.signalIRQTrigger.Emit(v.intrId) //intrVIA1Id)
		}
	}

	if (v.acr & 0x20) == 0 {
		// Only count in one-shot mode
		t2c := uint(v.t2c) - 1
		v.t2c = uint16(t2c)
		if t2c > defaultViaTimeout {
			v.ifr |= 0x20
		}
	}
}

func (v *Via) SignalPRA() {
	v.wiring.WritePRA(v.pra, v.ddra)
}

func (v *Via) SignalPRB() {
	v.wiring.WritePRB(v.prb, v.ddrb)
}

func (v *Via) ByteReady() bool {
	if (v.pcr & 0x0e) == 0x0e {
		return true
	}
	return false
}
