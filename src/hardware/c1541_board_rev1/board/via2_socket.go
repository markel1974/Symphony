package board

import (
	"github.com/markel1974/c64emu/src/hardware/c1541_board_rev1/mechanic"
	"github.com/markel1974/c64emu/src/references"
)

// headControl defines bitmask [0,1] used for head step direction management; controls upward or downward movement.
const headControl = uint8(0x3)

// motorControl is a constant used to represent motor control functionality. Bit [2] toggles motor state: 0 = Off; 1 = On.
const motorControl = uint8(0x4)

// ledControl represents the bit mask for controlling the LED; 0 means Off and 1 means On.
const ledControl = uint8(0x8)

// photocellControl represents the control code for managing photocell-related operations in the system.
const photocellControl = uint8(0x10)

// densityControl represents a constant value of type uint8 used for defining density settings or controls.
const densityControl = uint8(0x60)

// dataArrivedControl defines a bitmask used to indicate the arrival of data in specific operations or controls.
const dataArrivedControl = uint8(0x80)

// noPhotocellControl is the bitwise complement of photocellControl, used to disable photocell control functionality.
const noPhotocellControl = ^photocellControl

// syncArrivedControl is the bitwise NOT of dataArrivedControl, used as a mask to determine sync state in data operations.
const syncArrivedControl = ^dataArrivedControl

// IVIA2SocketConnections represents an interface for managing socket communication specific to VIA2 interactions.
// LEDTrigger is responsible for controlling LED signals, accepting an 8-bit unsigned integer as an argument.
// IRQClearTrigger clears the IRQ (Interrupt Request) signal, taking a 32-bit unsigned integer as input.
// IRQTrigger triggers the IRQ signal, taking a 32-bit unsigned integer as an input parameter.
type IVIA2SocketConnections interface {
	LedActivity(bool2 bool)

	IRQClearTrigger(uint32)

	IRQTrigger(uint32)
}

// VIA2Socket represents a data and signal interface between a VIA chip and a socket, utilizing a mechanic component for operations.
type VIA2Socket struct {
	references.IMos6522
	label       string
	parent      references.IComponent
	component   references.IComponent
	mec         mechanic.IMechanic
	connections IVIA2SocketConnections
	intrId      uint32
	prbPrev     uint8
	hwId        string
}

// NewVIA2Socket initializes a new VIA2Socket with the provided connections and mechanic, configuring IRQ and initial state.
func NewVIA2Socket(parent references.IComponent, label string, connections IVIA2SocketConnections, mec mechanic.IMechanic) *VIA2Socket {
	v := &VIA2Socket{
		IMos6522:    nil,
		parent:      parent,
		label:       label,
		mec:         mec,
		connections: connections,
		intrId:      intrIrqVIA2Bit,
		prbPrev:     0,
	}
	v.hwId = references.IdIMos6522(v.IMos6522, v.label, 1)
	return v
}

func (v *VIA2Socket) HardwareId() string {
	return v.hwId
}

// Wire initializes the VIA2Socket by configuring its IMos6522 component and applying its configuration settings.
func (v *VIA2Socket) Wire() error {
	var err error
	v.component = v.parent.GetChildByHardwareId(v.HardwareId())
	if v.IMos6522, err = references.ComponentToIMos6522(v.component); err != nil {
		return err
	}
	if err = v.IMos6522.Bind(v); err != nil {
		return err
	}
	return nil
}

// Reset reinitializes the VIA2Socket's internal state by setting prbPrev to 0 and invoking the Reset method of IMos6522.
func (v *VIA2Socket) Reset() {
	v.prbPrev = 0
	v.IMos6522.Reset()
}

// LedActivity triggers an LED indication based on the provided byte data, utilizing the established socket connections.
func (v *VIA2Socket) LedActivity(led bool) {
	v.connections.LedActivity(led)
}

// IRQClearTrigger clears the interrupt request (IRQ) signal associated with the current VIA2Socket instance.
func (v *VIA2Socket) IRQClearTrigger() {
	v.connections.IRQClearTrigger(v.intrId)
}

// IRQTrigger triggers an interrupt request (IRQ) for the associated connection using the stored interrupt ID.
func (v *VIA2Socket) IRQTrigger() {
	v.connections.IRQTrigger(v.intrId)
}

// ReadPRA reads a byte from the Mechanic through the VIA2Socket connection and returns the retrieved value.
func (v *VIA2Socket) ReadPRA(_ uint8, _ uint8) uint8 {
	d := v.mec.ReadByte()
	return d
}

// ReadPRB reads the PRB register, evaluates the sync state of the mechanic, and combines it with the photocell control state.
func (v *VIA2Socket) ReadPRB(prb uint8, _ uint8) uint8 {
	p := prb & noPhotocellControl
	photocellState := v.mec.WriteProtectionState()
	if v.mec.SyncFound() {
		return (p & syncArrivedControl) | photocellState
	} else {
		return (p | dataArrivedControl) | photocellState
	}
}

// WritePRA writes the provided value to the Peripheral Register A (PRA) using the mechanic's WriteByte function.
func (v *VIA2Socket) WritePRA(pra uint8, _ uint8) {
	v.mec.WriteByte(pra)
}

// WritePRB updates the PRB register and triggers actions based on the change in its bits compared to the previous value.
// Bits [0,1]: Controls head step direction; increasing or decreasing the value moves the head up or down respectively.
// Bit [2]: Controls motor state; 0 = Off, 1 = On.
// Bit [3]: Controls LED state; 0 = Off, 1 = On.
// Bit [4]: Indicates write protect status; 0 = Disk protected, 1 = Disk not protected (not fully implemented).
// Bits [5-6]: Represents data density levels (not fully implemented).
// Bit [7]: Indicates whether SYNC or data bytes are being read (not fully implemented).
func (v *VIA2Socket) WritePRB(prb uint8, _ uint8) {
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
		//fmt.Println("TODO - MOTOR", motorOn)
	}
	//bit [3]
	//LED control; 0 = Off; 1 = On.
	if (m & ledControl) != 0 {
		led := false
		if (prb & ledControl) != 0 {
			led = true
		}
		v.connections.LedActivity(led)
	}
	//bit [4]
	//Write protect photocell status; 0 = Write protect tab covered, disk protected; 1 = Tab uncovered, disk not protected.
	//if (m & photocellControl) != 0 {
	//photocell := (data & photocellControl) != 0
	//fmt.Println("TODO - PHOTOCELL", photocell)
	//}
	//bit [5-6]:
	//Data density; %00 = Lowest; %11 = Highest.
	//if (m & densityControl) != 0 {
	//density := (prb & densityControl) >> 5
	//fmt.Printf("TODO - DENSITY %2b\n", density)
	//}
	//Bit [7]
	//0 = SYNC marks are being currently read from disk; 1 = Data bytes are being read.
	//if (m & dataArrivedControl) != 0 {
	//sync := (prb & syncControl) != 0
	//fmt.Println("TODO - DATA ARRIVED", !sync)
	//} else {
	//sync := (prb & syncControl) != 0
	//fmt.Println("TODO - SYNC ARRIVED", !sync)
	//}
}

// WriteDDRA updates the Data Direction Register A (DDRA) with specified values, affecting the VIA port configuration.
func (v *VIA2Socket) WriteDDRA(_ uint8, _ uint8) {

}

// WriteDDRB writes data to the Data Direction Register B (DDRB) of the VIA2Socket, configuring the I/O port direction.
func (v *VIA2Socket) WriteDDRB(_ uint8, _ uint8) {

}

func (v *VIA2Socket) WriteCA2(w bool) {
	v.mec.SetWrite(w)
}

func (v *VIA2Socket) WriteCB2(bool) {

}

// ByteReady returns true if the peripheral control register (pcr) is in a ready state for data handling.
// Control SetOverflowBranch on 6502 cpu
func (v *VIA2Socket) ByteReady() bool {
	pcr := v.IMos6522.ReadByte(0xc)
	if (pcr & 0x0e) == 0x0e {
		return v.mec.ByteReady()
	}
	return false
}

func (v *VIA2Socket) ReadCA1() bool {
	return false
}

func (v *VIA2Socket) ReadCB1() bool {
	return false
}

func (v *VIA2Socket) ReadCB2() bool {
	return false
}

func (v *VIA2Socket) ReadPB6() bool {
	return false
}
