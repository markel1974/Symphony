package vic

type Graphics struct {
	core              *Core
	foreMask          *ForeMask
	db                IDisplayBuffer
	gfxData           uint8
	colorData         uint8
	charData          uint8
	charDataLast      uint8
	borderOn          bool    // Upper/lower border on (Main border FlipFlop)
	borderOnSample    []bool  // Samples of border state at different cycles (1, 17, 18, 56, 57)
	borderColorSample []uint8 // Samples of border color at each "displayed" cycle
}

func NewGraphics(core *Core, foreMask *ForeMask, db IDisplayBuffer) *Graphics {
	gr := &Graphics{
		core:              core,
		foreMask:          foreMask,
		db:                db,
		gfxData:           0,
		colorData:         0,
		charData:          0,
		charDataLast:      0,
		borderOnSample:    make([]bool, 5),
		borderColorSample: make([]uint8, DisplayXFill+1),
		borderOn:          false,
	}
	return gr
}

func (gr *Graphics) Setup() {
}

func (gr *Graphics) SetGfxData(data uint8) {
	gr.gfxData = data
}

func (gr *Graphics) SetColorData(data uint8) {
	gr.colorData = data
}

func (gr *Graphics) SetCharData(data uint8) {
	gr.charData = data
}

func (gr *Graphics) SetCharDataLast() {
	gr.charDataLast = gr.charData
}

func (gr *Graphics) SetBorderOn(b bool) {
	gr.borderOn = b
}

func (gr *Graphics) SetBorderColorSample(cycle int) {
	if gr.borderOn {
		idx := cycle - 13
		gr.borderColorSample[idx&DisplayXFill] = gr.core.ecColor
	}
}

func (gr *Graphics) SetBorderOnSample(idx int) {
	gr.borderOnSample[idx] = gr.borderOn
}

func (gr *Graphics) DrawBorder(lineStart int) {
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

func (gr *Graphics) DrawBackground(lineOffset int) {
	var c uint8
	switch gr.core.displayIdx {
	case 0, 1, 3: // Standard text, Multicolor text, Multicolor bitmap
		c = gr.core.b0cColor
	case 2: // Standard bitmap
		c = gr.core.colors[gr.charDataLast]
	case 4: // ECM text
		if (gr.charDataLast & 0x80) != 0 {
			if (gr.charDataLast & 0x40) != 0 {
				c = gr.core.b3cColor
			} else {
				c = gr.core.b2cColor
			}
		} else {
			if (gr.charDataLast & 0x40) != 0 {
				c = gr.core.b1cColor
			} else {
				c = gr.core.b0cColor
			}
		}
	default:
		c = gr.core.colors[0]
	}
	gr.db.SetMulti8(lineOffset, c)
}

func (gr *Graphics) DrawGraphics(lineOffset int) {
	offset := lineOffset + int(gr.core.xScroll)
	switch gr.core.displayIdx {
	case 0: // Standard text
		gr.drawGraphicStandard(offset, gr.core.b0cColor, gr.core.colors[gr.colorData])
	case 1: // Multicolor text
		if (gr.colorData & 8) != 0 {
			gr.drawGraphicMulticolor(offset, gr.core.b0cColor, gr.core.b1cColor, gr.core.b2cColor, gr.core.colors[gr.colorData&7])
		} else {
			gr.drawGraphicStandard(offset, gr.core.b0cColor, gr.core.colors[gr.colorData])
		}
	case 2: // Standard bitmap
		gr.drawGraphicStandard(offset, gr.core.colors[gr.charData], gr.core.colors[gr.charData>>4])
	case 3: // Multicolor bitmap
		gr.drawGraphicMulticolor(offset, gr.core.b0cColor, gr.core.colors[gr.charData>>4], gr.core.colors[gr.charData], gr.core.colors[gr.colorData])
	case 4: // ECM text
		if (gr.charData & 0x80) != 0 {
			if (gr.charData & 0x40) != 0 {
				gr.drawGraphicStandard(offset, gr.core.b3cColor, gr.core.colors[gr.colorData])
			} else {
				gr.drawGraphicStandard(offset, gr.core.b2cColor, gr.core.colors[gr.colorData])
			}
		} else if (gr.charData & 0x40) != 0 {
			gr.drawGraphicStandard(offset, gr.core.b1cColor, gr.core.colors[gr.colorData])
		} else {
			gr.drawGraphicStandard(offset, gr.core.b0cColor, gr.core.colors[gr.colorData])
		}
	case 5: //Invalid multicolor text
		if (gr.colorData & 8) != 0 {
			gr.drawGraphicsInvalidMulticolor(offset, gr.core.colors[0])
		} else {
			gr.drawGraphicsInvalidStandard(offset, gr.core.colors[0])
		}
	case 6: //Invalid standard bitmap
		gr.drawGraphicsInvalidStandard(offset, gr.core.colors[0])
	case 7: // Invalid multicolor bitmap
		gr.drawGraphicsInvalidMulticolor(offset, gr.core.colors[0])
	}
}

func (gr *Graphics) drawGraphicsInvalidStandard(offset int, a uint8) {
	gr.db.SetMulti8(offset, a)
	p1 := gr.gfxData >> gr.core.xScroll
	p2 := gr.gfxData << (7 - gr.core.xScroll)
	gr.foreMask.Update(p1, p2)
}

func (gr *Graphics) drawGraphicsInvalidMulticolor(offset int, a uint8) {
	gr.db.SetMulti8(offset, a)
	p1 := ((gr.gfxData & 0xaa) | (gr.gfxData&0xaa)>>1) >> gr.core.xScroll
	p2 := ((gr.gfxData & 0xaa) | (gr.gfxData&0xaa)>>1) << (8 - gr.core.xScroll)
	gr.foreMask.Update(p1, p2)
}

func (gr *Graphics) drawGraphicStandard(offset int, a uint8, b uint8) {
	p1 := gr.gfxData >> gr.core.xScroll
	p2 := gr.gfxData << (7 - gr.core.xScroll)
	gr.foreMask.Update(p1, p2)
	colorBuffer := [4]uint8{a, b, 0, 0}
	data := gr.gfxData
	gr.db.Set(offset+7, colorBuffer[data&1])
	data >>= 1
	gr.db.Set(offset+6, colorBuffer[data&1])
	data >>= 1
	gr.db.Set(offset+5, colorBuffer[data&1])
	data >>= 1
	gr.db.Set(offset+4, colorBuffer[data&1])
	data >>= 1
	gr.db.Set(offset+3, colorBuffer[data&1])
	data >>= 1
	gr.db.Set(offset+2, colorBuffer[data&1])
	data >>= 1
	gr.db.Set(offset+1, colorBuffer[data&1])
	data >>= 1
	gr.db.Set(offset, colorBuffer[data])
}

func (gr *Graphics) drawGraphicMulticolor(offset int, a uint8, b uint8, c uint8, d uint8) {
	p1 := ((gr.gfxData & 0xaa) | (gr.gfxData&0xaa)>>1) >> gr.core.xScroll
	p2 := ((gr.gfxData & 0xaa) | (gr.gfxData&0xaa)>>1) << (8 - gr.core.xScroll)
	gr.foreMask.Update(p1, p2)
	colorBuffer := [4]uint8{a, b, c, d}
	data := gr.gfxData
	gr.db.Set(offset+7, colorBuffer[data&3])
	gr.db.Set(offset+6, colorBuffer[data&3])
	data >>= 2
	gr.db.Set(offset+5, colorBuffer[data&3])
	gr.db.Set(offset+4, colorBuffer[data&3])
	data >>= 2
	gr.db.Set(offset+3, colorBuffer[data&3])
	gr.db.Set(offset+2, colorBuffer[data&3])
	data >>= 2
	gr.db.Set(offset+1, colorBuffer[data])
	gr.db.Set(offset, colorBuffer[data])
}
