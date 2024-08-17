package _510

type Core struct {
	banks   IBanks
	pic     IPic
	nFlag   uint8  // Negative flag
	zFlag   uint8  // Zero flag
	vFlag   uint8  // Overflow flag
	dFlag   uint8  // Decimal mode flag
	iFlag   uint8  // Interrupt disable flag
	cFlag   uint8  // Carry flag
	a       uint8  // Register
	x       uint8  // Register
	y       uint8  // Register
	sp      uint8  // Stack pointer
	pc      uint16 // Program counter
	op      uint8  // Current opcode
	ar      uint16 // Address register
	ar2     uint16 // Address register 2
	rmw     uint8  // Data buffer for RMW instructions
	state   uint8  // Current state
	opFlags uint8  //
	stop    bool   //
	rdyLow  bool   // current RDY state
	aecLow  bool   // current AEC state
}

func NewCore(pic IPic) *Core {
	regs := &Core{
		banks:   nil,
		pic:     pic,
		a:       0,
		x:       0,
		y:       0,
		sp:      0xff,
		nFlag:   0,
		zFlag:   0,
		vFlag:   0,
		dFlag:   0,
		cFlag:   0,
		iFlag:   1,
		opFlags: 0,
		ar:      0,
		ar2:     0,
		rdyLow:  false,
		aecLow:  false,
		stop:    false,
	}
	return regs
}
