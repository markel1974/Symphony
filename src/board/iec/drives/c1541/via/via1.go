package via

import (
	"github.com/markel1974/c64emu/src/board/iec/drives/c1541/cpu"
	"github.com/markel1974/c64emu/src/board/iec/virtualdrive"
)

type Via1 struct {
	*Core
	iec          virtualdrive.IIec
	intr         *cpu.Interrupts
	deviceNumber uint8
	_dipSwitch   uint8
	_prbFilter   uint8
}

func NewVia1(iec virtualdrive.IIec, intr *cpu.Interrupts, deviceNumber uint8) *Via1 {
	v := &Via1{
		Core:         NewCore(),
		intr:         intr,
		iec:          iec,
		_prbFilter:   0,
		deviceNumber: deviceNumber,
	}
	v._prbFilter |= 0 << 0 //Bit #0: DATA IN; 0 = Low; 1 = High.
	v._prbFilter |= 1 << 1 //Bit #1: DATA OUT; 0 = Low; 1 = High.
	v._prbFilter |= 0 << 2 //Bit #2: CLOCK IN; 0 = Low; 1 = High.
	v._prbFilter |= 1 << 3 //Bit #3: CLOCK OUT; 0 = Low; 1 = High..
	v._prbFilter |= 1 << 4 //Bit #4: ATNA OUT; 1 = Enable device presence detection by automatically acknowledging ATN IN signals on DATA OUT.
	v._prbFilter |= 1 << 5 //Bits #5 - #6: Device number, set with jumper, minus 8; % 00 = 8; % 01 = 9; % 10 = 10; % 11 = 11. Default: % 00, 8.
	v._prbFilter |= 1 << 6
	v._prbFilter |= 0 << 7 //Bit #7: ATN IN; 0 = Low; 1 = High.
	v.setDipSwitch()
	return v
}

func (v *Via1) Reset() {
	v.Core.Reset()
}

func (v *Via1) Setup() {

}

func (v *Via1) ReadByte(addr uint16) uint8 {
	switch addr {
	case 0x1800:
		data := v.iec.PeripheralRead(v.deviceNumber)
		return (v._prb&v._prbFilter | data) ^ 0x85
	case 0x1801:
		return 0xff // Keep 1541C ROMs happy (track 0 sensor)
	case 0x1802:
		return v._ddrb
	case 0x1803:
		return v._ddra
	case 0x1804:
		v._ifr &= 0xbf
		//TODO TEST
		v.intr.ClearVIA1IRQ()
		return uint8(v._t1c)
	case 0x1805:
		return uint8(v._t1c >> 8)
	case 0x1806:
		return uint8(v._t1l)
	case 0x1807:
		return uint8(v._t1l >> 8)
	case 0x1808:
		v._ifr &= 0xdf
		return uint8(v._t2c)
	case 0x1809:
		return uint8(v._t2c >> 8)
	case 0x180a:
		return v._sr
	case 0x180b:
		return v._acr
	case 0x180c:
		return v._pcr
	case 0x180d:
		if (v._ifr & v._ier) != 0 {
			return v._ifr | 0x80
		}
		return v._ifr
	case 0x180e:
		return v._ier | 0x80
	case 0x180f:
		return 0xff // Keep 1541C ROMs happy (track 0 sensor)
	default:
		return 0
	}
}

func (v *Via1) WriteByte(addr uint16, data uint8) {
	switch addr {
	case 0x1800:
		v._prb = data | v._dipSwitch
		data = (^v._prb) & v._ddrb
		v.iec.PeripheralWrite(v.deviceNumber, data)
	case 0x1801:
		v._pra = data
	case 0x1802:
		v._ddrb = data
		data &= ^v._prb
		v.iec.PeripheralWrite(v.deviceNumber, data)
	case 0x1803:
		v._ddra = data
	case 0x1804:
		v._t1l = (v._t1l & 0xff00) | uint16(data)
	case 0x1805:
		v._t1l = (v._t1l & 0xff) | (uint16(data) << 8)
		v._ifr &= 0xbf
		v._t1c = v._t1l
	case 0x1806:
		v._t1l = (v._t1l & 0xff00) | uint16(data)
	case 0x1807:
		v._t1l = (v._t1l & 0xff) | (uint16(data) << 8)
		break
	case 0x1808:
		v._t2l = (v._t2l & 0xff00) | uint16(data)
	case 0x1809:
		v._t2l = (v._t2l & 0xff) | (uint16(data) << 8)
		v._ifr &= 0xdf
		v._t2c = v._t2l
	case 0x180a:
		v._sr = data
	case 0x180b:
		v._acr = data
	case 0x180c:
		v._pcr = data
	case 0x180d:
		v._ifr &= ^data
	case 0x180e:
		if data&0x80 != 0 {
			v._ier |= data & 0x7f
		} else {
			v._ier &= ^data
		}
	case 0x180f:
		v._pra = data
	}
}

func (v *Via1) CountTimers() {
	tmp := uint(v._t1c) - 1
	v._t1c = uint16(tmp)

	if tmp > DEFAULT_VIA_TIMEOUT {
		if v._acr&0x40 != 0 {
			// Reload from latch in free-run mode
			v._t1c = v._t1l
		}
		v._ifr |= 0x40
		//TODO TEST
		if v._ier&0x40 != 0 {
			v.intr.TriggerVIA1()
		}
	}

	if v._acr&0x20 == 0 {
		// Only count in one-shot mode
		tmp = uint(v._t2c) - 1
		v._t2c = uint16(tmp)
		if tmp > DEFAULT_VIA_TIMEOUT {
			v._ifr |= 0x20
		}
	}
}

func (v *Via1) AtnStateChanged(state bool) {
	data := (^v._prb) & v._ddrb
	v.iec.PeripheralAtnResponse(v.deviceNumber, data)
}

func (v *Via1) setDipSwitch() {
	switch v.deviceNumber - 8 {
	case 0:
		v._dipSwitch |= 0 << 5
		v._dipSwitch |= 0 << 6
	case 1:
		v._dipSwitch |= 1 << 5
		v._dipSwitch |= 0 << 6
	case 2:
		v._dipSwitch |= 0 << 5
		v._dipSwitch |= 1 << 6
	case 3:
		v._dipSwitch |= 1 << 5
		v._dipSwitch |= 1 << 6
	default:
		v._dipSwitch |= 0 << 5
		v._dipSwitch |= 0 << 6
	}
}
