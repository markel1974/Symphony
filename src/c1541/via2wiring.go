package c1541

import (
	"fmt"
	"github.com/markel1974/c64emu/src/c1541/mechanics"
	"github.com/markel1974/c64emu/src/c64/iec/virtualdrive"
	"github.com/markel1974/c64emu/src/signals"
)

//TODO MOVE IN WIRED

const headControl = uint8(0x3)
const motorControl = uint8(0x4)
const ledControl = uint8(0x8)
const photocellControl = uint8(0x10)
const densityControl = uint8(0x60)
const syncControl = uint8(0x80)

const noPhotocellControl = ^photocellControl
const noSyncControl = ^syncControl

type Via2Wiring struct {
	mec          *mechanics.Mechanics
	iec          virtualdrive.IIec
	deviceNumber uint8
	prbPrev      uint8
	signalLed    *signals.SignalByte
}

func (v *Via2Wiring) WriteDDRA(u uint8, u2 uint8) {
}

func (v *Via2Wiring) WriteDDRB(u uint8, u2 uint8) {
}

func NewVia2Wiring(iec virtualdrive.IIec, mec *mechanics.Mechanics, deviceNumber uint8) *Via2Wiring {
	return &Via2Wiring{
		iec:          iec,
		mec:          mec,
		deviceNumber: deviceNumber,
		prbPrev:      0,
		signalLed:    signals.NewSignalByte(),
	}
}

func (v *Via2Wiring) Reset() {
	v.prbPrev = 0
}

func (v *Via2Wiring) SignalLedBind(fn func(byte)) {
	v.signalLed.Bind(fn)
}

func (v *Via2Wiring) ReadPRA(_ uint8, _ uint8) uint8 {
	d := v.mec.ReadByte()
	v.mec.RotateDisk()
	return d
}

func (v *Via2Wiring) ReadPRB(prb uint8, _ uint8) uint8 {
	p := prb & noPhotocellControl
	photocellState := v.mec.WriteProtectionState()
	if v.mec.SyncFound() {
		return (p & noSyncControl) | photocellState
	} else {
		v.mec.RotateDisk()
		return (p | syncControl) | photocellState
	}
}

func (v *Via2Wiring) WritePRA(pra uint8, _ uint8) {
	v.mec.WriteByte(pra)
	v.mec.RotateDisk()
}

func (v *Via2Wiring) WritePRB(prb uint8, _ uint8) {
	prevPrb := v.prbPrev
	v.prbPrev = prb
	m := prevPrb ^ prb

	//bit [0,1]
	//Head step direction.
	//Decrease value (%00-%11-%10-%01-%00...) to move head downwards
	//Increase value (%00-%01-%10-%11-%00...) to move head upwards
	if (m & headControl) != 0 {
		if (prevPrb & headControl) == ((prb + 1) & headControl) {
			v.mec.MoveHeadOut()
		} else if (prevPrb & headControl) == ((prb - 1) & headControl) {
			v.mec.MoveHeadIn()
		}
	}
	//bit [2]
	//Motor control; 0 = Off; 1 = On.
	if (m & motorControl) != 0 {
		motorOn := (prb & motorControl) != 0
		v.mec.SetMotor(motorOn)
		fmt.Println("TODO - MOTOR", motorOn)
	}
	//bit [3]
	//LED control; 0 = Off; 1 = On.
	if (m & ledControl) != 0 {
		led := uint8(0)
		if (prb & ledControl) != 0 {
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
		density := (prb & densityControl) >> 5
		fmt.Printf("TODO - DENSITY %2b\n", density)
	}
	//Bit [7]
	//0 = SYNC marks are being currently read from disk; 1 = Data bytes are being read.
	if (m & syncControl) != 0 {
		sync := (prb & syncControl) != 0
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
