package banks

//from c64pla.c

type Ports struct {
	dataOut  uint8
	dir      uint8
	data     uint8
	dataRead uint8
	//dirRead         uint8
	oldDataOut      uint8 /* Tape motor status.  */
	oldWriteBit     uint8 /* Tape write line status.  */
	oldSenseOut     uint8 /* Tape sense line out status. */
	dataSetBit6     uint8
	dataSetBit7     uint8
	dataFalloffBit6 uint8
	dataFalloffBit7 uint8
	tapeSense       int
	tapeWriteIn     int
	tapeMotorIn     int
	capsSense       int
	pullUp          uint8
}

func NewPorts() *Ports {
	return &Ports{
		capsSense: 1,
		pullUp:    0x17,
		dataOut:   0,
		dir:       0,
		data:      0,
		dataRead:  0,
		//dirRead:     0,
		oldDataOut:  0xff,
		oldWriteBit: 0xff,
		oldSenseOut: 0xff,
		tapeSense:   0,
		tapeWriteIn: 0,
		tapeMotorIn: 0,
	}
}

// Reset
// both DDR and DATA are 0 after poweron/reset, this means all pins are input,
// and the pullups at charen/hiram/loram/sense will pull up the respective lines
// so the kernal will be banked in and i/o is active.
// c64pla_pport_reset
func (p *Ports) Reset() {
	p.data = 0
	p.dataOut = 0
	p.dataRead = 0
	p.dir = 0
	//p.dirRead = 0
	p.dataSetBit6 = 0
	p.dataSetBit7 = 0
	p.dataFalloffBit6 = 0
	p.dataFalloffBit7 = 0
}

func (p *Ports) SetDir(data uint8) {
	p.dir = data
}

func (p *Ports) SetData(data uint8) {
	p.data = data
}

func (p *Ports) GetDirection() uint8 {
	return p.dir
}

func (p *Ports) GetDataRead() uint8 {
	return p.dataRead
}

//func (p *Ports) GetDataOut() uint8 {
//	return p.dataOut
//}

// SetTape - Tape sense status: 1 = some button pressed, 0 = no buttons pressed
func (p *Ports) SetTape(tapeSense int, tapeWriteIn int, tapeMotorIn int) {
	p.tapeSense = tapeSense
	p.tapeWriteIn = tapeWriteIn
	p.tapeMotorIn = tapeMotorIn
}

func (p *Ports) GetMemoryConfig(exRom uint8, game uint8) uint8 {
	c := ((^p.dir | p.data) & 0x7) | (exRom << 3) | (game << 4)
	return c
}

// Update - c64pla_config_changed
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
		//tapeport_set_motor(TAPEPORT_PORT_1, !p.oldDataOut)
	}
	if ((^p.dir | p.data) & 0x8) != p.oldWriteBit {
		p.oldWriteBit = (^p.dir | p.data) & 0x8
		//TODO IMPLEMENT
		//tapeport_toggle_write_bit(TAPEPORT_PORT_1, (^p.dir | p.data) & 0x8)
	}
	if ((p.dir & p.data) & 0x10) != p.oldSenseOut {
		p.oldSenseOut = (p.dir & p.data) & 0x10
		//TODO IMPLEMENT
		//tapeport_set_sense_out(TAPEPORT_PORT_1, !p.oldSenseOut)
	}
	//p.dirRead = p.dir
}
