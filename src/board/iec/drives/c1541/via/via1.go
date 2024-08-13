package via

import (
	"fmt"
	"github.com/markel1974/c64emu/src/board/iec/virtualdrive"
	"github.com/markel1974/c64emu/src/signals"
)

type Via1 struct {
	*Core
	iec              virtualdrive.IIec
	deviceNumber     uint8
	dipSwitch        uint8
	prbFilter        uint8
	signalIRQTrigger *signals.SignalUint32
	signalIRQClear   *signals.SignalUint32
}

func NewVia1(iec virtualdrive.IIec, deviceNumber uint8) *Via1 {
	v := &Via1{
		Core:             NewCore(),
		iec:              iec,
		prbFilter:        0,
		deviceNumber:     deviceNumber,
		signalIRQTrigger: signals.NewSignalUint32(),
		signalIRQClear:   signals.NewSignalUint32(),
	}
	v.prbFilter |= 0 << 0 //Bit #0: DATA IN; 0 = Low; 1 = High.
	v.prbFilter |= 1 << 1 //Bit #1: DATA OUT; 0 = Low; 1 = High.
	v.prbFilter |= 0 << 2 //Bit #2: CLOCK IN; 0 = Low; 1 = High.
	v.prbFilter |= 1 << 3 //Bit #3: CLOCK OUT; 0 = Low; 1 = High..
	v.prbFilter |= 1 << 4 //Bit #4: ATNA OUT; 1 = Enable device presence detection by automatically acknowledging ATN IN signals on DATA OUT.
	v.prbFilter |= 1 << 5 //Bits #5 - #6: Device number, set with jumper, minus 8; % 00 = 8; % 01 = 9; % 10 = 10; % 11 = 11. Default: % 00, 8.
	v.prbFilter |= 1 << 6
	v.prbFilter |= 0 << 7 //Bit #7: ATN IN; 0 = Low; 1 = High.
	v.setDipSwitch(deviceNumber)
	return v
}

func (v *Via1) Reset() {
	v.Core.Reset()
}

func (v *Via1) Setup() {

}

func (v *Via1) SignalTriggerIRQBind(fn func(uint32)) {
	v.signalIRQTrigger.Bind(fn)
}

func (v *Via1) SignalClearIRQBind(fn func(uint32)) {
	v.signalIRQClear.Bind(fn)
}

func (v *Via1) ReadByte(addr uint16) uint8 {
	switch addr {
	case 0x1800:
		data := v.iec.PeripheralRead()
		ret := (v.prb&v.prbFilter | data) ^ 0x85
		fmt.Println("READING FROM IEC", data, ret)
		return ret
	case 0x1801:
		// Keep 1541C ROMs happy (track 0 sensor)
		return 0xff
	case 0x1802:
		return v.ddrb
	case 0x1803:
		return v.ddra
	case 0x1804:
		v.ifr &= 0xbf
		v.signalIRQClear.Emit(intrVIA1Id)
		return uint8(v.t1c)
	case 0x1805:
		return uint8(v.t1c >> 8)
	case 0x1806:
		return uint8(v.t1l)
	case 0x1807:
		return uint8(v.t1l >> 8)
	case 0x1808:
		v.ifr &= 0xdf
		return uint8(v.t2c)
	case 0x1809:
		return uint8(v.t2c >> 8)
	case 0x180a:
		return v.sr
	case 0x180b:
		return v.acr
	case 0x180c:
		return v.pcr
	case 0x180d:
		if (v.ifr & v.ier) != 0 {
			return v.ifr | 0x80
		}
		return v.ifr
	case 0x180e:
		return v.ier | 0x80
	case 0x180f:
		// Keep 1541C ROMs happy (track 0 sensor)
		return 0xff
	default:
		return 0
	}
}

func (v *Via1) WriteByte(addr uint16, data uint8) {
	switch addr {
	case 0x1800:
		v.prb = data | v.dipSwitch
		data = (^v.prb) & v.ddrb
		fmt.Println("0x1800) WRITING TO IEC", data)
		v.iec.PeripheralWrite(v.deviceNumber, data)
	case 0x1801:
		v.pra = data
	case 0x1802:
		v.ddrb = data
		data &= ^v.prb
		fmt.Println("0x1802) WRITING TO IEC", data)
		v.iec.PeripheralWrite(v.deviceNumber, data)
	case 0x1803:
		v.ddra = data
	case 0x1804:
		v.t1l = (v.t1l & 0xff00) | uint16(data)
	case 0x1805:
		v.t1l = (v.t1l & 0xff) | (uint16(data) << 8)
		v.ifr &= 0xbf
		v.t1c = v.t1l
	case 0x1806:
		v.t1l = (v.t1l & 0xff00) | uint16(data)
	case 0x1807:
		v.t1l = (v.t1l & 0xff) | (uint16(data) << 8)
	case 0x1808:
		v.t2l = (v.t2l & 0xff00) | uint16(data)
	case 0x1809:
		v.t2l = (v.t2l & 0xff) | (uint16(data) << 8)
		v.ifr &= 0xdf
		v.t2c = v.t2l
	case 0x180a:
		v.sr = data
	case 0x180b:
		v.acr = data
	case 0x180c:
		v.pcr = data
	case 0x180d:
		v.ifr &= ^data
	case 0x180e:
		if data&0x80 != 0 {
			v.ier |= data & 0x7f
		} else {
			v.ier &= ^data
		}
	case 0x180f:
		v.pra = data
	}
}

func (v *Via1) CountTimers() {
	tmp := uint(v.t1c) - 1
	v.t1c = uint16(tmp)

	if tmp > defaultViaTimeout {
		if v.acr&0x40 != 0 {
			// Reload from latch in free-run mode
			v.t1c = v.t1l
		}
		v.ifr |= 0x40
		//TODO TEST
		if v.ier&0x40 != 0 {
			v.signalIRQTrigger.Emit(intrVIA1Id)
		}
	}

	if v.acr&0x20 == 0 {
		// Only count in one-shot mode
		tmp = uint(v.t2c) - 1
		v.t2c = uint16(tmp)
		if tmp > defaultViaTimeout {
			v.ifr |= 0x20
		}
	}
}

func (v *Via1) AtnStateChanged() {
	data := (^v.prb) & v.ddrb
	//v.iec.PeripheralAtnResponse(v.deviceNumber, data)
	v.iec.PeripheralWrite(v.deviceNumber, data)
}

func (v *Via1) setDipSwitch(deviceNumber uint8) {
	switch deviceNumber - 8 {
	case 0:
		v.dipSwitch |= 0 << 5
		v.dipSwitch |= 0 << 6
	case 1:
		v.dipSwitch |= 1 << 5
		v.dipSwitch |= 0 << 6
	case 2:
		v.dipSwitch |= 0 << 5
		v.dipSwitch |= 1 << 6
	case 3:
		v.dipSwitch |= 1 << 5
		v.dipSwitch |= 1 << 6
	default:
		v.dipSwitch |= 0 << 5
		v.dipSwitch |= 0 << 6
	}
}
