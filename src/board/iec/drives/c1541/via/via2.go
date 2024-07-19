package via

import (
	"github.com/markel1974/c64emu/src/board/iec/drives/c1541/cpu"
	"github.com/markel1974/c64emu/src/board/iec/drives/c1541/mechanics"
	"github.com/markel1974/c64emu/src/board/iec/virtualdrive"
)

type Via2 struct {
	*Core
	iec  virtualdrive.IIec
	job  *mechanics.Job
	intr *cpu.Interrupts
}

func NewVia2(iec virtualdrive.IIec, intr *cpu.Interrupts, job *mechanics.Job) *Via2 {
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
			return (v._prb & 0x7f) | ps
		} else {
			return (v._prb | 0x80) | ps
		}
	case 0x1c01:
		return v.job.ReadGCRByte()
	case 0x1c02:
		return v._ddrb
	case 0x1c03:
		return v._ddra
	case 0x1c04:
		v._ifr &= 0xbf
		v.intr.ClearVIA2IRQ()
		//_cpu1541->_interrupt.intr[INT_VIA2IRQ] = false;
		return uint8(v._t1c)
	case 0x1c05:
		return uint8(v._t1c >> 8)
	case 0x1c06:
		return uint8(v._t1l)
	case 0x1c07:
		return uint8(v._t1l >> 8)
	case 0x1c08:
		v._ifr &= 0xdf
		return uint8(v._t2c)
	case 0x1c09:
		return uint8(v._t2c >> 8)
	case 0x1c0a:
		return v._sr
	case 0x1c0b:
		return v._acr
	case 0x1c0c:
		return v._pcr
	case 0x1c0d:
		if v._ifr&v._ier != 0 {
			return v._ifr | 0x80
		}
		return v._ifr
	case 0x1c0e:
		return v._ier | 0x80
	case 0x1c0f:
		return v.job.ReadGCRByte()
	default:
		return 0
	}
}

func (v *Via2) WriteByte(addr uint16, data uint8) {
	switch addr {
	case 0x1c00:
		if ((v._prb ^ data) & 8) != 0 {
			l := 0
			if (data & 8) != 0 {
				l = 1
			}
			v.job.UpdateLEDs(l) // Bit 3: VirtualDrive LED
		}

		if ((v._prb ^ data) & 3) != 0 {
			/* Bits 0/1: Stepper motor */
			if (v._prb & 3) == ((data + 1) & 3) {
				v.job.MoveHeadOut()
			} else if (v._prb & 3) == ((data - 1) & 3) {
				v.job.MoveHeadIn()
			}
		}
		v._prb = data & 0xef
	case 0x1c01:
		v._pra = data
	case 0x1c02:
		v._ddrb = data
	case 0x1c03:
		v._ddra = data
	case 0x1c04:
		v._t1l = (v._t1l & 0xff00) | uint16(data)
	case 0x1c05:
		v._t1l = (v._t1l & 0xff) | (uint16(data) << 8)
		v._ifr &= 0xbf
		v._t1c = v._t1l
	case 0x1c06:
		v._t1l = (v._t1l & 0xff00) | uint16(data)
	case 0x1c07:
		v._t1l = (v._t1l & 0xff) | (uint16(data) << 8)
	case 0x1c08:
		v._t2l = (v._t2l & 0xff00) | uint16(data)
	case 0x1c09:
		v._t2l = (v._t2l & 0xff) | (uint16(data) << 8)
		v._ifr &= 0xdf
		v._t2c = v._t2l
	case 0x1c0a:
		v._sr = data
	case 0x1c0b:
		v._acr = data
	case 0x1c0c:
		v._pcr = data
	case 0x1c0d:
		v._ifr &= ^data
	case 0x1c0e:
		if data&0x80 != 0 {
			v._ier |= data & 0x7f
		} else {
			v._ier &= ^data
		}
	case 0x1c0f:
		v._pra = data
	}
}

func (v *Via2) CountTimers() {
	tmp := uint(v._t1c) - 1

	v._t1c = uint16(tmp)
	if tmp > DEFAULT_VIA_TIMEOUT {
		// Reload from latch in free-run mode
		if (v._acr & 0x40) != 0 {
			v._t1c = v._t1l
		}
		v._ifr |= 0x40
		if (v._ier & 0x40) != 0 {
			v.intr.TriggerVIA2()
		}
	}

	if (v._acr & 0x20) == 0 {
		// Only count in one-shot mode
		tmp = uint(v._t2c) - 1
		v._t2c = uint16(tmp)
		if tmp > DEFAULT_VIA_TIMEOUT {
			v._ifr |= 0x20
		}
	}
}

func (v *Via2) WriteSector() {
	v.job.WriteSector()
}

func (v *Via2) FormatTrack() {
	v.job.FormatTrack()
}
