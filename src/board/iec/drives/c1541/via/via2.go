package via

import (
	"github.com/markel1974/c64emu/src/board/iec/drives/c1541/cpu"
	"github.com/markel1974/c64emu/src/board/iec/drives/c1541/mechanics"
	"github.com/markel1974/c64emu/src/board/iec/virtualdrive"
)

type Via2 struct {
	*Core
	iec  virtualdrive.IIec
	job  *mechanics.Mechanics
	intr *cpu.Interrupts
}

func NewVia2(iec virtualdrive.IIec, intr *cpu.Interrupts, job *mechanics.Mechanics) *Via2 {
	v := &Via2{
		Core: NewCore(),
		iec:  iec,
		intr: intr,
		job:  job,
	}
	return v
}

func (v *Via2) Reset() {
	v.Core.Reset()
}

func (v *Via2) Setup() {

}

func (v *Via2) ReadByte(addr uint16) uint8 {
	switch addr {
	case 0x1c00:
		ps := v.job.WriteProtectionState()
		if v.job.SyncFound() {
			return (v.prb & 0x7f) | ps
		} else {
			return (v.prb | 0x80) | ps
		}
	case 0x1c01:
		return v.job.ReadGCRByte()
	case 0x1c02:
		return v.ddrb
	case 0x1c03:
		return v.ddra
	case 0x1c04:
		v.ifr &= 0xbf
		v.intr.ClearVIA2IRQ()
		//_cpu1541->_interrupt.intr[INT_VIA2IRQ] = false;
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
		return v.job.ReadGCRByte()
	default:
		return 0
	}
}

func (v *Via2) WriteByte(addr uint16, data uint8) {
	switch addr {
	case 0x1c00:
		if ((v.prb ^ data) & 8) != 0 {
			l := 0
			if (data & 8) != 0 {
				l = 1
			}
			v.job.UpdateLEDs(l) // Bit 3: VirtualDrive LED
		}

		if ((v.prb ^ data) & 3) != 0 {
			/* Bits 0/1: Stepper motor */
			if (v.prb & 3) == ((data + 1) & 3) {
				v.job.MoveHeadOut()
			} else if (v.prb & 3) == ((data - 1) & 3) {
				v.job.MoveHeadIn()
			}
		}
		v.prb = data & 0xef
	case 0x1c01:
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
		v.pra = data
	}
}

func (v *Via2) CountTimers() {
	tmp := uint(v.t1c) - 1

	v.t1c = uint16(tmp)
	if tmp > defaultViaTimeout {
		// Reload from latch in free-run mode
		if (v.acr & 0x40) != 0 {
			v.t1c = v.t1l
		}
		v.ifr |= 0x40
		if (v.ier & 0x40) != 0 {
			v.intr.TriggerVIA2()
		}
	}

	if (v.acr & 0x20) == 0 {
		// Only count in one-shot mode
		tmp = uint(v.t2c) - 1
		v.t2c = uint16(tmp)
		if tmp > defaultViaTimeout {
			v.ifr |= 0x20
		}
	}
}

func (v *Via2) WriteSector() {
	v.job.WriteSector()
}

func (v *Via2) FormatTrack() {
	v.job.FormatTrack()
}
