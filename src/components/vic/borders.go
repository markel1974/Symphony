package vic

type Borders struct {
	core              *Core
	db                IDisplayBuffer
	mainFlipFlop      bool    // Main borders FlipFlop
	verticalFlipFlop  bool    // Vertical Border FlipFlop
	borderOnSample    []bool  // Samples of borders state at different cycles (1, 17, 18, 56, 57)
	borderColorSample []uint8 // Samples of borders color at each "displayed" cycle
}

func NewBorder(core *Core, db IDisplayBuffer) *Borders {
	gr := &Borders{
		db:                db,
		core:              core,
		borderOnSample:    make([]bool, 5),
		borderColorSample: make([]uint8, DisplayXFill+1),
		mainFlipFlop:      false,
		verticalFlipFlop:  false,
	}
	return gr
}

func (b *Borders) Sample(cycle int) {
	if b.mainFlipFlop {
		idx := cycle - 13
		b.borderColorSample[idx&DisplayXFill] = b.core.ecColor
	}
}

func (b *Borders) SetBorderOn() {
	b.mainFlipFlop = true
}

func (b *Borders) SetBorderOnSample(idx int) {
	b.borderOnSample[idx] = b.mainFlipFlop
}

func (b *Borders) GetVerticalFlipFlop() bool {
	return b.verticalFlipFlop
}

func (b *Borders) UpdateVerticalFlipFlop() {
	if b.core.rasterY == b.core.dyStop {
		b.verticalFlipFlop = true
	} else if (b.core.cr1&0x10) != 0 && b.core.rasterY == b.core.dyStart {
		b.verticalFlipFlop = false
	}
}

func (b *Borders) Update() {
	if b.core.rasterY == b.core.dyStop {
		b.verticalFlipFlop = true
	} else {
		if (b.core.cr1 & 0x10) != 0 {
			if b.core.rasterY == b.core.dyStart {
				b.verticalFlipFlop = false
				b.mainFlipFlop = false
			} else if !b.verticalFlipFlop {
				b.mainFlipFlop = false
			}
		} else if !b.verticalFlipFlop {
			b.mainFlipFlop = false
		}
	}
}

func (b *Borders) Draw(lineStart int) {
	const BorderS = 43
	const BorderOffset = BorderS * 8
	if b.borderOnSample[0] {
		for idx := 0; idx < 4; idx++ {
			b.db.SetMulti8(lineStart+(idx*8), b.borderColorSample[idx])
		}
	}
	if b.borderOnSample[1] {
		//32 = 4*8
		b.db.SetMulti8(lineStart+(32), b.borderColorSample[4])
	}
	if b.borderOnSample[2] {
		//TODO VERIFICA
		for idx := 5; idx < BorderS; idx++ {
			b.db.SetMulti8(lineStart+(idx*8), b.borderColorSample[idx])
		}
	}
	if b.borderOnSample[3] {
		b.db.SetMulti8(lineStart+(BorderOffset), b.borderColorSample[BorderS])
	}
	if b.borderOnSample[4] {
		for idx := 44; idx < DisplayXDiv8; idx++ {
			b.db.SetMulti8(lineStart+(idx*8), b.borderColorSample[idx])
		}
	}
}
