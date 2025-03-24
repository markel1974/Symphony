package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/hardware/c1541/mechanic"
	"github.com/markel1974/c64emu/src/references"
)

// headControl is a bitmask (0x3) used for controlling the head step direction during mechanical operations.
const headControl = uint8(0x3)

// motorControl represents the motor control bit in the VIA2 PRB register. 0 = Motor off; 1 = Motor on.
const motorControl = uint8(0x4)

// ledControl defines the bit mask for controlling the LED state; 0 represents Off, and 1 represents On.
const ledControl = uint8(0x8)

// photocellControl represents the control bitmask for managing photocell-related functionalities in a device.
const photocellControl = uint8(0x10)

// densityControl defines a constant with a uint8 value of 0x60, likely representing a control code for density settings.
const densityControl = uint8(0x60)

// dataArrivedControl is a constant used to indicate the presence of newly received data in certain control operations.
const dataArrivedControl = uint8(0x80)

// noPhotocellControl is a constant used to disable the photocell control by flipping all bits of `photocellControl`.
const noPhotocellControl = ^photocellControl

// syncArrivedControl is the bitwise complement of dataArrivedControl, typically used to denote a sync state in operations.
const syncArrivedControl = ^dataArrivedControl

// IVIA2SocketConnections defines an interface for managing connections and signal events between VIA2 and a socket.
// LedChanged signals a change in the LED state identified by the provided uint8 value.
// IRQClear clears an interrupt request by using the specified interrupt ID as a uint32.
// IRQTrigger triggers an interrupt using the specified interrupt ID as a uint32.
type IVIA2SocketConnections interface {
	LEDTrigger(uint8)

	IRQClear(uint32)

	IRQTrigger(uint32)
}

// VIA2Socket represents a VIA (Versatile Interface Adapter) connected to a socket for communication and control.
// It integrates a VIA interface, mechanic, and socket connection handlers for IRQ and LED signal changes.
// intrId defines the unique interrupt ID for handling IRQs through the associated connections.
// prbPrev tracks the previous state of the PRB (Peripheral Register B) for state change detection.
type VIA2Socket struct {
	references.IVIA
	mec         *mechanic.Mechanic
	connections IVIA2SocketConnections
	intrId      uint32
	prbPrev     uint8
}

// NewVIA2Socket creates and returns a pointer to a new VIA2Socket instance with default initializations.
func NewVIA2Socket(connections IVIA2SocketConnections, mec *mechanic.Mechanic) *VIA2Socket {
	return &VIA2Socket{
		IVIA:        nil,
		mec:         mec,
		connections: connections,
		intrId:      intrIrqVIA2Bit,
		prbPrev:     0,
	}
}

// Setup initializes the VIA2Socket instance with the provided IVIA, connections, and Mechanic, and sets up the IVIA.
func (v *VIA2Socket) Setup(c map[string]references.IComponent, _ *config.Config) error {
	via2, err := references.ComponentsToIVIA(c, 1)
	if err != nil {
		return err
	}
	v.IVIA = via2
	if err = v.IVIA.Setup(v); err != nil {
		return err
	}
	return nil
}

func (v *VIA2Socket) Connect() error {
	return nil
}

// Reset reinitializes the state of the VIA2Socket and its associated VIA by resetting their internal configurations.
func (v *VIA2Socket) Reset() {
	v.prbPrev = 0
	v.IVIA.Reset()
}

// LEDTrigger notifies the connected system that the LED state has changed, passing the updated data byte as a parameter.
func (v *VIA2Socket) LEDTrigger(data byte) {
	v.connections.LEDTrigger(data)
}

// IRQClear clears the interrupt request associated with the intrId of the VIA2Socket instance by invoking the connections' IRQClear method.
func (v *VIA2Socket) IRQClear() {
	v.connections.IRQClear(v.intrId)
}

// IRQTrigger sends an interrupt request trigger signal through the associated connection using the stored interrupt ID.
func (v *VIA2Socket) IRQTrigger() {
	v.connections.IRQTrigger(v.intrId)
}

// ReadPRA reads a byte of data from the Mechanic instance associated with the VIA2Socket and returns the value.
func (v *VIA2Socket) ReadPRA(_ uint8, _ uint8) uint8 {
	d := v.mec.ReadByte()
	return d
}

// ReadPRB reads the PRB register, applies controls based on the photocell state and sync status, and returns the result.
func (v *VIA2Socket) ReadPRB(prb uint8, _ uint8) uint8 {
	p := prb & noPhotocellControl
	photocellState := v.mec.WriteProtectionState()
	if v.mec.SyncFound() {
		return (p & syncArrivedControl) | photocellState
	} else {
		return (p | dataArrivedControl) | photocellState
	}
}

// WritePRA writes the specified PRA (Peripheral Register A) value to the Mechanic via `WriteByte`.
func (v *VIA2Socket) WritePRA(pra uint8, _ uint8) {
	v.mec.WriteByte(pra)
}

// WritePRB processes a value written to the PRB register, updating hardware state based on changed bits in the input.
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
		led := uint8(0)
		if (prb & ledControl) != 0 {
			led = 1
		}
		v.connections.LEDTrigger(led)
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

// WriteDDRA handles writing to the Data Direction Register A (DDRA) by accepting the current and next values, performing no action.
func (v *VIA2Socket) WriteDDRA(_ uint8, _ uint8) {

}

// WriteDDRB handles the write operation for DDRB (Data Direction Register B) with given input parameters.
func (v *VIA2Socket) WriteDDRB(_ uint8, _ uint8) {

}
