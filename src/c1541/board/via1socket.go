package board

type Via1Socket struct {
	board     *Board
	intrId    uint32
	dipSwitch uint8
	prbFilter uint8
}

func NewVia1Socket(board *Board, intrId uint32) *Via1Socket {
	v := &Via1Socket{
		board:     board,
		intrId:    intrId,
		prbFilter: 0,
	}
	v.setFilters()
	v.setDipSwitch(board.deviceNumber)
	return v
}

func (v *Via1Socket) Reset() {
}

func (v *Via1Socket) IRQClear() {
	v.board.IRQClear(v.intrId)
}

func (v *Via1Socket) IRQTrigger() {
	v.board.IRQTrigger(v.intrId)
}

func (v *Via1Socket) setFilters() {
	v.prbFilter |= 0 << 0 //Bit #0: DATA IN; 0 = Low; 1 = High.
	v.prbFilter |= 1 << 1 //Bit #1: DATA OUT; 0 = Low; 1 = High.
	v.prbFilter |= 0 << 2 //Bit #2: CLOCK IN; 0 = Low; 1 = High.
	v.prbFilter |= 1 << 3 //Bit #3: CLOCK OUT; 0 = Low; 1 = High..
	v.prbFilter |= 1 << 4 //Bit #4: ATNA OUT; 1 = Enable device presence detection by automatically acknowledging ATN IN signals on DATA OUT.
	v.prbFilter |= 1 << 5 //Bits #5 - #6: Device number, set with jumper, minus 8; % 00 = 8; % 01 = 9; % 10 = 10; % 11 = 11. Default: % 00, 8.
	v.prbFilter |= 1 << 6
	v.prbFilter |= 0 << 7 //Bit #7: ATN IN; 0 = Low; 1 = High.
}

func (v *Via1Socket) setDipSwitch(deviceNumber uint8) {
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

func (v *Via1Socket) ReadPRA(_ uint8, _ uint8) uint8 {
	// Keep 1541C ROMs happy (track 0 sensor)
	return 0xff
}

func (v *Via1Socket) ReadPRB(prb uint8, _ uint8) uint8 {
	data := v.board.iec.PeripheralRead()
	p := (prb | v.dipSwitch) & v.prbFilter
	//bit 0 - 2 - 7 = 0x85
	ret := (p | data) ^ 0x85
	return ret
}

func (v *Via1Socket) WritePRA(_ uint8, _ uint8) {
}

func (v *Via1Socket) WritePRB(prb uint8, ddrb uint8) {
	v.peripheralWrite(prb, ddrb)
}

func (v *Via1Socket) WriteDDRA(_ uint8, _ uint8) {

}

func (v *Via1Socket) WriteDDRB(prb uint8, ddrb uint8) {
	v.peripheralWrite(prb, ddrb)
}

func (v *Via1Socket) peripheralWrite(prb uint8, ddrb uint8) {
	p := prb | v.dipSwitch
	wd := (^p) & ddrb
	v.board.iec.PeripheralWrite(v.board.deviceNumber, wd)
}
