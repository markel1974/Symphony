package vic

type Borders struct {
	core              *Core
	db                IDisplayBuffer
	borderOn          bool    // Upper-lower borders on (Main borders FlipFlop)
	borderOnSample    []bool  // Samples of borders state at different cycles (1, 17, 18, 56, 57)
	borderColorSample []uint8 // Samples of borders color at each "displayed" cycle
	borderULOn        bool    // Upper-lower borders on
}

func NewBorder(core *Core, db IDisplayBuffer) *Borders {
	gr := &Borders{
		db:                db,
		core:              core,
		borderOnSample:    make([]bool, 5),
		borderColorSample: make([]uint8, DisplayXFill+1),
		borderOn:          false,
		borderULOn:        false,
	}
	return gr
}

func (gr *Borders) Sample(cycle int) {
	if gr.borderOn {
		idx := cycle - 13
		gr.borderColorSample[idx&DisplayXFill] = gr.core.ecColor
	}
}

func (gr *Borders) SetBorderOn() {
	gr.borderOn = true
}

func (gr *Borders) SetBorderOnSample(idx int) {
	gr.borderOnSample[idx] = gr.borderOn
}

func (gr *Borders) GetBorderULOn() bool {
	return gr.borderULOn
}

func (gr *Borders) UpdateBorderUpperLower() {
	if gr.core.rasterY == gr.core.dyStop {
		gr.borderULOn = true
	} else if (gr.core.cr1&0x10) != 0 && gr.core.rasterY == gr.core.dyStart {
		gr.borderULOn = false
	}
}

func (gr *Borders) UpdateBorder() {
	if gr.core.rasterY == gr.core.dyStop {
		gr.borderULOn = true
	} else {
		if (gr.core.cr1 & 0x10) != 0 {
			if gr.core.rasterY == gr.core.dyStart {
				gr.borderULOn = false
				gr.borderOn = false
			} else if !gr.borderULOn {
				gr.borderOn = false
			}
		} else if !gr.borderULOn {
			gr.borderOn = false
		}
	}
}

func (gr *Borders) Draw(lineStart int) {
	const BorderS = 43
	const BorderOffset = BorderS * 8
	if gr.borderOnSample[0] {
		for idx := 0; idx < 4; idx++ {
			gr.db.SetMulti8(lineStart+(idx*8), gr.borderColorSample[idx])
		}
	}
	if gr.borderOnSample[1] {
		//32 = 4*8
		gr.db.SetMulti8(lineStart+(32), gr.borderColorSample[4])
	}
	if gr.borderOnSample[2] {
		for idx := 5; idx < BorderS; idx++ {
			gr.db.SetMulti8(lineStart+(idx*8), gr.borderColorSample[idx])
		}
	}
	if gr.borderOnSample[3] {
		gr.db.SetMulti8(lineStart+(BorderOffset), gr.borderColorSample[BorderS])
	}
	if gr.borderOnSample[4] {
		for idx := 44; idx < DisplayXDiv8; idx++ {
			gr.db.SetMulti8(lineStart+(idx*8), gr.borderColorSample[idx])
		}
	}
}
