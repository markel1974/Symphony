package board

const headControl = uint8(0x3)
const motorControl = uint8(0x4)
const ledControl = uint8(0x8)
const photocellControl = uint8(0x10)
const densityControl = uint8(0x60)
const syncControl = uint8(0x80)

const noPhotocellControl = ^photocellControl
const noSyncControl = ^syncControl

type Via2Socket struct {
	board   *Board
	intrId  uint32
	prbPrev uint8
}

func NewVia2Socket(board *Board, intrId uint32) *Via2Socket {
	return &Via2Socket{
		board:   board,
		intrId:  intrId,
		prbPrev: 0,
	}
}

func (v *Via2Socket) Reset() {
	v.prbPrev = 0
	v.board.via2.Reset()
}

func (v *Via2Socket) LedChanged(data byte) {
	v.board.LedChanged(data)
}

func (v *Via2Socket) IRQClear() {
	v.board.IRQClear(v.intrId)
}

func (v *Via2Socket) IRQTrigger() {
	v.board.IRQTrigger(v.intrId)
}

func (v *Via2Socket) ReadPRA(_ uint8, _ uint8) uint8 {
	d := v.board.mec.ReadByte()
	v.board.mec.RotateDisk()
	return d
}

func (v *Via2Socket) ReadPRB(prb uint8, _ uint8) uint8 {
	p := prb & noPhotocellControl
	photocellState := v.board.mec.WriteProtectionState()
	if v.board.mec.SyncFound() {
		return (p & noSyncControl) | photocellState
	} else {
		v.board.mec.RotateDisk()
		return (p | syncControl) | photocellState
	}
}

func (v *Via2Socket) WritePRA(pra uint8, _ uint8) {
	v.board.mec.WriteByte(pra)
	v.board.mec.RotateDisk()
}

func (v *Via2Socket) WritePRB(prb uint8, _ uint8) {
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
	if (m & photocellControl) != 0 {
		//photocell := (data & photocellControl) != 0
		//fmt.Println("TODO - PHOTOCELL", photocell)
	}
	//bit [5-6]:
	//Data density; %00 = Lowest; %11 = Highest.
	if (m & densityControl) != 0 {
		//density := (prb & densityControl) >> 5
		//fmt.Printf("TODO - DENSITY %2b\n", density)
	}
	//Bit [7]
	//0 = SYNC marks are being currently read from disk; 1 = Data bytes are being read.
	if (m & syncControl) != 0 {
		//sync := (prb & syncControl) != 0
		//fmt.Println("TODO - SYNC", !sync)
	}
}

func (v *Via2Socket) WriteDDRA(_ uint8, _ uint8) {

}

func (v *Via2Socket) WriteDDRB(_ uint8, _ uint8) {

}

/*
func (v *Via2) WriteSector() {
	v.mec.WriteSector()
}

func (v *Via2) FormatTrack() {
	v.mec.FormatTrack()
}
*/
