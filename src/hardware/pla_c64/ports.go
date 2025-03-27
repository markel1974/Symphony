package pla_c64

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

//from c64pla.c

// Ports represents the state and management of hardware I/O ports used for device communication and control.
type Ports struct {
	*component.BaseComponent
	dataOut         uint8
	dir             uint8
	data            uint8
	dataRead        uint8
	oldDataOut      uint8 // Tape motor status
	oldWriteBit     uint8 // Tape write line status
	oldSenseOut     uint8 // Tape sense line out status
	dataSetBit6     uint8
	dataSetBit7     uint8
	dataFalloffBit6 uint8
	dataFalloffBit7 uint8
	tapeSense       int
	tapeWriteIn     int
	tapeMotorIn     int
	capsSense       int
	pullUp          uint8
	//dirRead         uint8
}

// NewPorts initializes and returns a new instance of the Ports struct with default values set.
func NewPorts(factory references.IComponentFactory, parent references.IComponent, instance int) *Ports {
	p := &Ports{
		BaseComponent: component.NewBaseComponent(),
		capsSense:     1,
		pullUp:        0x17,
		dataOut:       0,
		dir:           0,
		data:          0,
		dataRead:      0,
		oldDataOut:    0xff,
		oldWriteBit:   0xff,
		oldSenseOut:   0xff,
		tapeSense:     0,
		tapeWriteIn:   0,
		tapeMotorIn:   0,
		//dirRead:     0,
	}
	p.BaseComponent.Register(factory, parent, "ports", p, references.IdInternalComponent("Ports", instance))
	return p
}

func (p *Ports) Connect() error {
	return nil
}

func (p *Ports) Internal() bool {
	return true
}

// Reset resets all port-related variables in the Ports struct to their default state.
func (p *Ports) Reset() {
	p.data = 0
	p.dataOut = 0
	p.dataRead = 0
	p.dir = 0
	p.dataSetBit6 = 0
	p.dataSetBit7 = 0
	p.dataFalloffBit6 = 0
	p.dataFalloffBit7 = 0
	//p.dirRead = 0
}

func (p *Ports) Emulate() {
	//
}

func (p *Ports) EmulationRequired() bool {
	return false
}

// SetDir sets the direction register of the Ports to the specified value.
func (p *Ports) SetDir(data uint8) {
	p.dir = data
}

// SetData sets the data property of the Ports instance to the specified value.
func (p *Ports) SetData(data uint8) {
	p.data = data
}

// GetDirection returns the current direction configuration of the port as an 8-bit unsigned integer.
func (p *Ports) GetDirection() uint8 {
	return p.dir
}

// GetDataRead retrieves the current value of the `dataRead` field from the Ports structure.
func (p *Ports) GetDataRead() uint8 {
	return p.dataRead
}

//func (p *Ports) GetDataOut() uint8 {
//	return p.dataOut
//}

// SetTape sets the tape sense, tape write input, and tape motor input values for the Ports instance.
func (p *Ports) SetTape(tapeSense int, tapeWriteIn int, tapeMotorIn int) {
	p.tapeSense = tapeSense
	p.tapeWriteIn = tapeWriteIn
	p.tapeMotorIn = tapeMotorIn
}

// GetMemoryConfig calculates the memory configuration index based on port direction, data values, and cartridge signals.
func (p *Ports) GetMemoryConfig(exRom uint8, game uint8) uint8 {
	c := ((^p.dir | p.data) & 0x7) | (exRom << 3) | (game << 4)
	return c
}

// Update recalculates the state of the Ports structure based on current register values, adjusting outputs and data logic.
func (p *Ports) Update() {
	//6 Bits - (on cpu are P0 - P1 - P2 - P3 - P4 - P5 - P6)
	//Bit 3: Datasette output signal level.
	//Bit 4: Datasette button status; 0 = One or more of PLAY, RECORD, F.FWD or REW pressed; 1 = No button is pressed.
	//Bit 5: Datasette motor control; 0 = On; 1 = Off.
	p.dataOut = (p.dataOut & ^p.dir) | (p.data & p.dir)
	p.dataRead = (p.data | ^p.dir) & (p.dataOut | p.pullUp)
	if (p.pullUp&0x40) != 0 && (p.capsSense == 0) {
		p.dataRead &= 0xbf
	}
	if p.dir&0x20 == 0 {
		p.dataRead &= 0xdf
	}
	if p.tapeSense != 0 && ((p.dir & 0x10) == 0) {
		p.dataRead &= 0xef
	}
	if p.tapeWriteIn != 0 && ((p.dir & 0x08) == 0) {
		p.dataRead &= 0xf7
	}
	if p.tapeMotorIn != 0 && ((p.dir & 0x20) == 0) {
		p.dataRead &= 0xdf
	}
	if ((p.dir & p.data) & 0x20) != p.oldDataOut {
		p.oldDataOut = (p.dir & p.data) & 0x20
		//TODO IMPLEMENT
		//tapePort.setMotor(TAPEPORT_PORT_1, !p.oldDataOut)
	}
	if ((^p.dir | p.data) & 0x8) != p.oldWriteBit {
		p.oldWriteBit = (^p.dir | p.data) & 0x8
		//TODO IMPLEMENT
		//tapePort.toggleWriteBit(TAPEPORT_PORT_1, (^p.dir | p.data) & 0x8)
	}
	if ((p.dir & p.data) & 0x10) != p.oldSenseOut {
		p.oldSenseOut = (p.dir & p.data) & 0x10
		//TODO IMPLEMENT
		//tapePort.setSenseOut(TAPEPORT_PORT_1, !p.oldSenseOut)
	}
	//p.dirRead = p.dir
}
