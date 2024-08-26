package via

import (
	"fmt"
	"github.com/markel1974/c64emu/src/c64/iec/virtualdrive"
	"github.com/markel1974/c64emu/src/signals"
)

// see drive/iecieee/via2d.c [store_pra - store_prb]
// 1541, 1541II, 1571 and 2031
// see https://sta.c64.org/cbm1541mem.html

type Via2 struct {
	*Core
	iec              virtualdrive.IIec
	mec              IMechanics
	deviceNumber     uint8
	signalIRQTrigger *signals.SignalUint32
	signalIRQClear   *signals.SignalUint32
	signalLed        *signals.SignalByte
}

func NewVia2(iec virtualdrive.IIec, mec IMechanics, deviceNumber uint8) *Via2 {
	v := &Via2{
		Core:             NewCore(),
		iec:              iec,
		mec:              mec,
		deviceNumber:     deviceNumber,
		signalIRQTrigger: signals.NewSignalUint32(),
		signalIRQClear:   signals.NewSignalUint32(),
		signalLed:        signals.NewSignalByte(),
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

func (v *Via2) SignalLedBind(fn func(byte)) {
	v.signalLed.Bind(fn)
}

func (v *Via2) ReadByte(addr uint16) uint8 {
	switch addr {
	case 0x1c00:
		wps := v.mec.WriteProtectionState()
		if v.mec.SyncFound() {
			return (v.prb & 0x7f) | wps
		} else {
			v.mec.RotateDisk()
			return (v.prb | 0x80) | wps
		}
	case 0x1c01:
		d := v.mec.ReadByte()
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
		d := v.mec.ReadByte()
		v.mec.RotateDisk()
		return d
	default:
		return 0
	}
}

func (v *Via2) WriteByte(addr uint16, data uint8) {
	switch addr {
	case 0x1c00:
		v.updatePRB(v.prb, data)
		v.prb = data & 0xef
	case 0x1c01:
		v.mec.WriteByte(data)
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
		v.mec.WriteByte(data)
		v.mec.RotateDisk()
		v.pra = data
	}
}

func (v *Via2) Emulate() {
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

func (v *Via2) updatePRB(prb uint8, data uint8) {
	const headControl = 0x3
	const motorControl = 0x4
	const ledControl = 0x8
	const photocellControl = 0x10
	const densityControl = 0x60
	const syncControl = 0x80

	m := prb ^ data

	//bit [0,1]
	//Head step direction.
	//Decrease value (%00-%11-%10-%01-%00...) to move head downwards
	//Increase value (%00-%01-%10-%11-%00...) to move head upwards
	if (m & headControl) != 0 {
		if (prb & headControl) == ((data + 1) & headControl) {
			v.mec.MoveHeadOut()
		} else if (prb & headControl) == ((data - 1) & headControl) {
			v.mec.MoveHeadIn()
		}
	}
	//bit [2]
	//Motor control; 0 = Off; 1 = On.
	if (m & motorControl) != 0 {
		motorOn := (data & motorControl) != 0
		v.mec.SetMotor(motorOn)
		fmt.Println("TODO - MOTOR", motorOn)
	}
	//bit [3]
	//LED control; 0 = Off; 1 = On.
	if (m & ledControl) != 0 {
		led := uint8(0)
		if (data & ledControl) != 0 {
			led = 1
		}
		v.signalLed.Emit(led)
	}
	//bit [4]
	//Write protect photocell status; 0 = Write protect tab covered, disk protected; 1 = Tab uncovered, disk not protected.
	if (m & photocellControl) != 0 {
		//photocell := (data & photocellControl) != 0
		//fmt.Println("TODO - PHOTOCELL", photocell)
	}
	//bit [5-6]:
	//Data density; %00 = Lowest; %11 = Highest.
	if (m & densityControl) != 0 {
		density := (data & densityControl) >> 5
		fmt.Printf("TODO - DENSITY %2b\n", density)
	}
	//Bit [7]
	//0 = SYNC marks are being currently read from disk; 1 = Data bytes are being read.
	if (m & syncControl) != 0 {
		sync := (data & syncControl) != 0
		fmt.Println("TODO - SYNC", !sync)
	}
}

/*
func (v *Via2) WriteSector() {
	v.mec.WriteSector()
}

func (v *Via2) FormatTrack() {
	v.mec.FormatTrack()
}
*/
