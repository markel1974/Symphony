package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// headControl represents the bit mask for controlling the head movement direction in the PRB register of VIA2.
const headControl = uint8(0x3)

// motorControl represents the bit mask for controlling the motor state in PRB. 0 indicates off; 1 indicates on.
const motorControl = uint8(0x4)

// ledControl is a constant representing the LED control bit (bit [3]) in a control register. 0 = Off; 1 = On.
const ledControl = uint8(0x8)

// photocellControl represents a uint8 configuration flag for enabling or indicating a photocell-related control feature.
const photocellControl = uint8(0x10)

// densityControl is a constant representing a fixed uint8 value of 0x60, typically used for density configuration or identification.
const densityControl = uint8(0x60)

// dataArrivedControl is a constant that signals the presence of new data by its associated bit set in an 8-bit value.
const dataArrivedControl = uint8(0x80)

// noPhotocellControl inverts all bits of photocellControl, disabling its effect when applied in certain operations.
const noPhotocellControl = ^photocellControl

// syncArrivedControl is the bitwise complement of dataArrivedControl, used to manage synchronization states in the system.
const syncArrivedControl = ^dataArrivedControl

// VIA2Socket represents a socket interface for interacting with the VIA2 (Versatile Interface Adapter) component on the board.
type VIA2Socket struct {
	via2    references.IVIA
	board   *Board
	intrId  uint32
	prbPrev uint8
}

// NewVIA2Socket creates and returns a new instance of VIA2Socket with default initialized fields.
func NewVIA2Socket() *VIA2Socket {
	return &VIA2Socket{
		via2:    nil,
		board:   nil,
		intrId:  intrIrqVIA2Bit,
		prbPrev: 0,
	}
}

// Connect initializes the VIA2Socket by associating it with a Board instance and configuring the interrupt ID.
func (v *VIA2Socket) Connect(board *Board, via2 references.IVIA) error {
	v.board = board
	v.via2 = via2
	v.via2.Setup(v)
	return nil
}

func (v *VIA2Socket) Emulate() {
	v.via2.Emulate()
}

// Reset reinitializes the VIA2Socket to its default state by clearing prbPrev and invoking the Reset method on board.via2.
func (v *VIA2Socket) Reset() {
	v.prbPrev = 0
	v.via2.Reset()
}

func (v *VIA2Socket) ByteReady() func() bool {
	return v.via2.ByteReady
}

// LedChanged updates the LED state by forwarding the provided data to the board's LED change handler.
func (v *VIA2Socket) LedChanged(data byte) {
	v.board.LedChanged(data)
}

// IRQClear clears the interrupt request associated with this VIA2Socket instance by delegating to the board's IRQClear method.
func (v *VIA2Socket) IRQClear() {
	v.board.IRQClear(v.intrId)
}

// IRQTrigger triggers an interrupt request (IRQ) on the board using the interrupt ID associated with the VIA2Socket instance.
func (v *VIA2Socket) IRQTrigger() {
	v.board.IRQTrigger(v.intrId)
}

// ReadPRA reads a byte from the board's `Mechanic` and returns the value.
func (v *VIA2Socket) ReadPRA(_ uint8, _ uint8) uint8 {
	d := v.board.mec.ReadByte()
	return d
}

// ReadPRB processes the PRB value, combines it with the write protection state, and adjusts based on synchronization status.
func (v *VIA2Socket) ReadPRB(prb uint8, _ uint8) uint8 {
	p := prb & noPhotocellControl
	photocellState := v.board.mec.WriteProtectionState()
	if v.board.mec.SyncFound() {
		return (p & syncArrivedControl) | photocellState
	} else {
		return (p | dataArrivedControl) | photocellState
	}
}

// WritePRA writes the given PRA value to the Mechanic's disk via the WriteByte method.
func (v *VIA2Socket) WritePRA(pra uint8, _ uint8) {
	v.board.mec.WriteByte(pra)
}

// WritePRB updates the state and behavior of the `VIA2Socket` based on the given PRB byte input.
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
			v.board.mec.MoveHeadOut()
		} else if (prevPrb & headControl) == ((prb - 1) & headControl) {
			v.board.mec.MoveHeadIn()
		}
	}
	//bit [2]
	//Motor control; 0 = Off; 1 = On.
	if (m & motorControl) != 0 {
		motorOn := (prb & motorControl) != 0
		v.board.mec.SetMotor(motorOn)
		//fmt.Println("TODO - MOTOR", motorOn)
	}
	//bit [3]
	//LED control; 0 = Off; 1 = On.
	if (m & ledControl) != 0 {
		led := uint8(0)
		if (prb & ledControl) != 0 {
			led = 1
		}
		v.board.LedChanged(led)
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

// WriteDDRA handles writing operations to the data direction register A (DDRA) for the VIA2 component. This method is unimplemented.
func (v *VIA2Socket) WriteDDRA(_ uint8, _ uint8) {

}

// WriteDDRB sets the Data Direction Register B for VIA2 but currently contains no implementation.
func (v *VIA2Socket) WriteDDRB(_ uint8, _ uint8) {

}
