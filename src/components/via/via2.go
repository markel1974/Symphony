package via

import (
	"github.com/markel1974/c64emu/src/board/iec/drives/c1541/mechanics"
	"github.com/markel1974/c64emu/src/board/iec/virtualdrive"
	"github.com/markel1974/c64emu/src/signals"
)

// see drive/iecieee/via2d.c [store_pra - store_prb]
// 1541, 1541II, 1571 and 2031
// see https://sta.c64.org/cbm1541mem.html

type Via2 struct {
	*Core
	iec              virtualdrive.IIec
	mec              *mechanics.Mechanics
	signalIRQTrigger *signals.SignalUint32
	signalIRQClear   *signals.SignalUint32
}

func NewVia2(iec virtualdrive.IIec, mec *mechanics.Mechanics) *Via2 {
	v := &Via2{
		Core:             NewCore(),
		iec:              iec,
		mec:              mec,
		signalIRQTrigger: signals.NewSignalUint32(),
		signalIRQClear:   signals.NewSignalUint32(),
	}
	return v
}

func (v *Via2) Reset() {
	v.Core.Reset()
}

func (v *Via2) Setup() {
}

func (v *Via2) SignalTriggerIRQBind(fn func(uint32)) {
	v.signalIRQTrigger.Bind(fn)
}

func (v *Via2) SignalClearIRQBind(fn func(uint32)) {
	v.signalIRQClear.Bind(fn)
}

func (v *Via2) ReadByte(addr uint16) uint8 {
	switch addr {
	case 0x1c00:
		wps := v.mec.WriteProtectionState()
		if v.mec.SyncFound() {
			return (v.prb & 0x7f) | wps
		}
		return (v.prb | 0x80) | wps
	case 0x1c01:
		d := v.mec.ReadGCRByte()
		v.mec.RotateDisk()
		return d
	case 0x1c02:
		return v.ddrb
	case 0x1c03:
		return v.ddra
	case 0x1c04:
		v.ifr &= 0xbf
		v.signalIRQClear.Emit(intrVIA2Id)
		return uint8(v.t1c)
	case 0x1c05:
		return uint8(v.t1c >> 8)
	case 0x1c06:
		return uint8(v.t1l)
	case 0x1c07:
		return uint8(v.t1l >> 8)
	case 0x1c08:
		v.ifr &= 0xdf
		return uint8(v.t2c)
	case 0x1c09:
		return uint8(v.t2c >> 8)
	case 0x1c0a:
		return v.sr
	case 0x1c0b:
		return v.acr
	case 0x1c0c:
		return v.pcr
	case 0x1c0d:
		if v.ifr&v.ier != 0 {
			return v.ifr | 0x80
		}
		return v.ifr
	case 0x1c0e:
		return v.ier | 0x80
	case 0x1c0f:
		d := v.mec.ReadGCRByte()
		v.mec.RotateDisk()
		return d
	default:
		return 0
	}
}

func (v *Via2) WriteByte(addr uint16, data uint8) {
	switch addr {
	case 0x1c00:
		v.mec.UpdatePRB(v.prb, data)
		v.prb = data & 0xef
	case 0x1c01:
		v.mec.WriteGCRByte(data)
		v.mec.RotateDisk()
		v.pra = data
	case 0x1c02:
		v.ddrb = data
	case 0x1c03:
		v.ddra = data
	case 0x1c04:
		v.t1l = (v.t1l & 0xff00) | uint16(data)
	case 0x1c05:
		v.t1l = (v.t1l & 0xff) | (uint16(data) << 8)
		v.ifr &= 0xbf
		v.t1c = v.t1l
	case 0x1c06:
		v.t1l = (v.t1l & 0xff00) | uint16(data)
	case 0x1c07:
		v.t1l = (v.t1l & 0xff) | (uint16(data) << 8)
	case 0x1c08:
		v.t2l = (v.t2l & 0xff00) | uint16(data)
	case 0x1c09:
		v.t2l = (v.t2l & 0xff) | (uint16(data) << 8)
		v.ifr &= 0xdf
		v.t2c = v.t2l
	case 0x1c0a:
		v.sr = data
	case 0x1c0b:
		v.acr = data
	case 0x1c0c:
		v.pcr = data
	case 0x1c0d:
		v.ifr &= ^data
	case 0x1c0e:
		if data&0x80 != 0 {
			v.ier |= data & 0x7f
		} else {
			v.ier &= ^data
		}
	case 0x1c0f:
		v.mec.WriteGCRByte(data)
		v.mec.RotateDisk()
		v.pra = data
	}
}

func (v *Via2) CountTimers() {
	t1c := uint(v.t1c) - 1
	v.t1c = uint16(t1c)
	if t1c > defaultViaTimeout {
		// Reload from latch in free-run mode
		if (v.acr & 0x40) != 0 {
			v.t1c = v.t1l
		}
		v.ifr |= 0x40
		if (v.ier & 0x40) != 0 {
			v.signalIRQTrigger.Emit(intrVIA2Id)
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

func (v *Via2) ByteReady() bool {
	if v.pcr&0x0e == 0x0e {
		return true
	}
	return false
}

func (v *Via2) WriteSector() {
	v.mec.WriteSector()
}

func (v *Via2) FormatTrack() {
	v.mec.FormatTrack()
}
