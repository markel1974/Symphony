package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// VIA1Socket represents the interface between the VIA1 chip and the board, managing IRQs, filters, and dip switches.
type VIA1Socket struct {
	via1      references.IVIA
	board     *Board
	intrId    uint32
	dipSwitch uint8
	prbFilter uint8
}

// NewVIA1Socket creates and returns a new instance of VIA1Socket with default, uninitialized state values.
func NewVIA1Socket() *VIA1Socket {
	return &VIA1Socket{
		via1:      nil,
		board:     nil,
		intrId:    intrIrqVIA1Bit,
		prbFilter: 0,
	}
}

// Connect initializes the VIA1Socket by assigning the board and interrupt ID, setting filters, and configuring the dip switch.
func (v *VIA1Socket) Connect(board *Board, via1 references.IVIA) error {
	v.board = board
	v.via1 = via1
	v.setFilters()
	v.setDipSwitch(board.deviceNumber)
	v.via1.Setup(v)
	return nil
}

func (v *VIA1Socket) Emulate() {
	v.via1.Emulate()
}

// Reset reinitializes the VIA1 socket by invoking the Reset method of the associated VIA instance.
func (v *VIA1Socket) Reset() {
	v.via1.Reset()
}

func (v *VIA1Socket) SignalPRB() {
	v.via1.SignalPRB()
}

// IRQClear clears the interrupt request (IRQ) for the current VIA1Socket using the associated board and interrupt ID.
func (v *VIA1Socket) IRQClear() {
	v.board.IRQClear(v.intrId)
}

// IRQTrigger triggers an interrupt request using the associated board and interrupt ID of the VIA1Socket.
func (v *VIA1Socket) IRQTrigger() {
	v.board.IRQTrigger(v.intrId)
}

// setFilters configures the `prbFilter` property by applying bitwise operations to set specific filter flags.
func (v *VIA1Socket) setFilters() {
	v.prbFilter |= 0 << 0 //Bit #0: DATA IN; 0 = Low; 1 = High.
	v.prbFilter |= 1 << 1 //Bit #1: DATA OUT; 0 = Low; 1 = High.
	v.prbFilter |= 0 << 2 //Bit #2: CLOCK IN; 0 = Low; 1 = High.
	v.prbFilter |= 1 << 3 //Bit #3: CLOCK OUT; 0 = Low; 1 = High..
	v.prbFilter |= 1 << 4 //Bit #4: ATNA OUT; 1 = Enable device presence detection by automatically acknowledging ATN IN signals on DATA OUT.
	v.prbFilter |= 1 << 5 //Bits #5 - #6: Device number, set with jumper, minus 8; % 00 = 8; % 01 = 9; % 10 = 10; % 11 = 11. Default: % 00, 8.
	v.prbFilter |= 1 << 6
	v.prbFilter |= 0 << 7 //Bit #7: ATN IN; 0 = Low; 1 = High.
}

// setDipSwitch configures the dip switch settings based on the provided device number.
func (v *VIA1Socket) setDipSwitch(deviceNumber uint8) {
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

// ReadPRA returns a fixed value of 0xff to ensure compatibility with 1541C ROMs (track 0 sensor).
func (v *VIA1Socket) ReadPRA(_ uint8, _ uint8) uint8 {
	// Keep 1541C ROMs happy (track 0 sensor)
	return 0xff
}

// ReadPRB reads from the Peripheral Register B, applying filters, dip switch values, and a specific bitwise operation.
func (v *VIA1Socket) ReadPRB(prb uint8, _ uint8) uint8 {
	data := v.board.iec.PeripheralRead()
	p := (prb | v.dipSwitch) & v.prbFilter
	//bit 0 - 2 - 7 = 0x85
	ret := (p | data) ^ 0x85
	return ret
}

// WritePRA performs a write operation on Port A of the VIA1 socket but does not modify any state in this implementation.
func (v *VIA1Socket) WritePRA(_ uint8, _ uint8) {
}

// WritePRB writes data to the peripheral bus using the specified PRB and DDRB values.
func (v *VIA1Socket) WritePRB(prb uint8, ddrb uint8) {
	v.peripheralWrite(prb, ddrb)
}

// WriteDDRA writes data direction settings to the DDRA register, specifying input or output configuration for PORT A.
func (v *VIA1Socket) WriteDDRA(_ uint8, _ uint8) {

}

// WriteDDRB writes data direction values to the PRB register and updates the peripheral accordingly.
func (v *VIA1Socket) WriteDDRB(prb uint8, ddrb uint8) {
	v.peripheralWrite(prb, ddrb)
}

// peripheralWrite performs a write operation to an IEC peripheral device using the given port B value and data direction.
func (v *VIA1Socket) peripheralWrite(prb uint8, ddrb uint8) {
	p := prb | v.dipSwitch
	wd := (^p) & ddrb
	v.board.iec.PeripheralWrite(v.board.deviceNumber, wd)
}
