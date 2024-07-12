package vic

type Core struct {
	mXx [8]uint16 // VIC registers [m0x - m1x - m2x - m3x - m4x - m5x - m6x - m7x]
	mXy [8]uint8  // VIC registers [m0y - m1y - m2y - m3y - m4y - m5y - m6y - m7y]
	mx8 uint8     // VIC register
	cr1 uint8     // VIC register
	cr2 uint8     // VIC register
	lpx uint8     // VIC register
	lpy uint8     // VIC register
	me  uint8     // VIC register
	mxe uint8     // VIC register
	mye uint8     // VIC register
	mdp uint8     // VIC register
	mmc uint8     // VIC register
	ec  uint8     // VIC register
	b0c uint8     // VIC register
	b1c uint8     // VIC register
	b2c uint8     // VIC register
	b3c uint8     // VIC register
	mm0 uint8     // VIC register
	mm1 uint8     // VIC register
	mXc [8]uint8  // VIC registers [m0c - m1c - m2c - m3c - m4c - m5c - m6c - m7c]
}

func NewCore() *Core {
	c := &Core{
		mXx: [8]uint16{},
		mXy: [8]uint8{},
		mx8: 0,
		cr1: 0,
		cr2: 0,
		lpx: 0,
		lpy: 0,
		me:  0,
		mxe: 0,
		mye: 0,
		mdp: 0,
		mmc: 0,
		ec:  0,
		b0c: 0,
		b1c: 0,
		b2c: 0,
		b3c: 0,
		mm0: 0,
		mm1: 0,
		mXc: [8]uint8{},
	}
	for i := 0; i < 8; i++ {
		c.mXx[i] = 0
		c.mXy[i] = 0
		c.mXc[i] = 0
	}
	return c
}
