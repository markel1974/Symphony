package vic

var _colorMultiplier = make([][]uint8, 0xff)

func init() {
	for idx := range _colorMultiplier {
		x := uint8(idx)
		_colorMultiplier[idx] = []uint8{x, x, x, x, x, x, x, x}
	}
}

type Graphics struct {
	*Core
	gfxData           uint8
	colorData         uint8
	charData          uint8
	charDataLast      uint8
	borderOn          bool    // Upper/lower border on (Main border FlipFlop)
	borderOnSample    []bool  // Samples of border state at different cycles (1, 17, 18, 56, 57)
	borderColorSample []uint8 // Samples of border color at each "displayed" cycle

}

func NewGraphics(core *Core) *Graphics {
	gr := &Graphics{
		Core:              core,
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
		gr.borderColorSample[idx&DisplayXFill] = gr.ecColor
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
			copy(gr.displayBuffer[lineStart+(idx*8):], _colorMultiplier[gr.borderColorSample[idx]])
		}
	}
	if gr.borderOnSample[1] {
		//32 = 4*8
		copy(gr.displayBuffer[lineStart+(32):], _colorMultiplier[gr.borderColorSample[4]])
	}
	if gr.borderOnSample[2] {
		for idx := 5; idx < BorderS; idx++ {
			copy(gr.displayBuffer[lineStart+(idx*8):], _colorMultiplier[gr.borderColorSample[idx]])
		}
	}
	if gr.borderOnSample[3] {
		copy(gr.displayBuffer[lineStart+(BorderOffset):], _colorMultiplier[gr.borderColorSample[BorderS]])
	}
	if gr.borderOnSample[4] {
		for idx := 44; idx < DisplayXDiv8; idx++ {
			copy(gr.displayBuffer[lineStart+(idx*8):], _colorMultiplier[gr.borderColorSample[idx]])
		}
	}
}

func (gr *Graphics) DrawBackground(lineOffset int) {
	var c uint8
	switch gr.displayIdx {
	case 0, 1, 3: // Standard text, Multicolor text, Multicolor bitmap
		c = gr.b0cColor
	case 2: // Standard bitmap
		c = gr.colors[gr.charDataLast]
	case 4: // ECM text
		if (gr.charDataLast & 0x80) != 0 {
			if (gr.charDataLast & 0x40) != 0 {
				c = gr.b3cColor
			} else {
				c = gr.b2cColor
			}
		} else {
			if (gr.charDataLast & 0x40) != 0 {
				c = gr.b1cColor
			} else {
				c = gr.b0cColor
			}
		}
	default:
		c = gr.colors[0]
	}
	copy(gr.displayBuffer[lineOffset:], _colorMultiplier[c])
}

func (gr *Graphics) DrawGraphics(lineOffset int) {
	offset := lineOffset + int(gr.xScroll)
	switch gr.displayIdx {
	case 0: // Standard text
		gr.drawGraphicStandard(offset, gr.b0cColor, gr.colors[gr.colorData])
	case 1: // Multicolor text
		if (gr.colorData & 8) != 0 {
			gr.drawGraphicMulticolor(offset, gr.b0cColor, gr.b1cColor, gr.b2cColor, gr.colors[gr.colorData&7])
		} else {
			gr.drawGraphicStandard(offset, gr.b0cColor, gr.colors[gr.colorData])
		}
	case 2: // Standard bitmap
		gr.drawGraphicStandard(offset, gr.colors[gr.charData], gr.colors[gr.charData>>4])
	case 3: // Multicolor bitmap
		gr.drawGraphicMulticolor(offset, gr.b0cColor, gr.colors[gr.charData>>4], gr.colors[gr.charData], gr.colors[gr.colorData])
	case 4: // ECM text
		if (gr.charData & 0x80) != 0 {
			if (gr.charData & 0x40) != 0 {
				gr.drawGraphicStandard(offset, gr.b3cColor, gr.colors[gr.colorData])
			} else {
				gr.drawGraphicStandard(offset, gr.b2cColor, gr.colors[gr.colorData])
			}
		} else if (gr.charData & 0x40) != 0 {
			gr.drawGraphicStandard(offset, gr.b1cColor, gr.colors[gr.colorData])
		} else {
			gr.drawGraphicStandard(offset, gr.b0cColor, gr.colors[gr.colorData])
		}
	case 5: //Invalid multicolor text
		if (gr.colorData & 8) != 0 {
			gr.drawGraphicsInvalidMulticolor(offset, gr.colors[0])
		} else {
			gr.drawGraphicsInvalidStandard(offset, gr.colors[0])
		}
	case 6: //Invalid standard bitmap
		gr.drawGraphicsInvalidStandard(offset, gr.colors[0])
	case 7: // Invalid multicolor bitmap
		gr.drawGraphicsInvalidMulticolor(offset, gr.colors[0])
	}
}

func (gr *Graphics) drawGraphicsInvalidStandard(offset int, a uint8) {
	copy(gr.displayBuffer[offset:], _colorMultiplier[a])
	gr.foreMaskBuf[gr.foreMaskOffset+0] |= gr.gfxData >> gr.xScroll
	gr.foreMaskBuf[gr.foreMaskOffset+1] |= gr.gfxData << (7 - gr.xScroll)
}

func (gr *Graphics) drawGraphicsInvalidMulticolor(offset int, a uint8) {
	copy(gr.displayBuffer[offset:], _colorMultiplier[a])
	gr.foreMaskBuf[gr.foreMaskOffset+0] |= ((gr.gfxData & 0xaa) | (gr.gfxData&0xaa)>>1) >> gr.xScroll
	gr.foreMaskBuf[gr.foreMaskOffset+1] |= ((gr.gfxData & 0xaa) | (gr.gfxData&0xaa)>>1) << (8 - gr.xScroll)
}

func (gr *Graphics) drawGraphicStandard(offset int, a uint8, b uint8) {
	gr.foreMaskBuf[gr.foreMaskOffset+0] |= gr.gfxData >> gr.xScroll
	gr.foreMaskBuf[gr.foreMaskOffset+1] |= gr.gfxData << (7 - gr.xScroll)
	colorBuffer := [4]uint8{a, b, 0, 0}
	data := gr.gfxData
	gr.displayBuffer[offset+7] = colorBuffer[data&1]
	data >>= 1
	gr.displayBuffer[offset+6] = colorBuffer[data&1]
	data >>= 1
	gr.displayBuffer[offset+5] = colorBuffer[data&1]
	data >>= 1
	gr.displayBuffer[offset+4] = colorBuffer[data&1]
	data >>= 1
	gr.displayBuffer[offset+3] = colorBuffer[data&1]
	data >>= 1
	gr.displayBuffer[offset+2] = colorBuffer[data&1]
	data >>= 1
	gr.displayBuffer[offset+1] = colorBuffer[data&1]
	data >>= 1
	gr.displayBuffer[offset] = colorBuffer[data]
}

func (gr *Graphics) drawGraphicMulticolor(offset int, a uint8, b uint8, c uint8, d uint8) {
	gr.foreMaskBuf[gr.foreMaskOffset+0] |= ((gr.gfxData & 0xaa) | (gr.gfxData&0xaa)>>1) >> gr.xScroll
	gr.foreMaskBuf[gr.foreMaskOffset+1] |= ((gr.gfxData & 0xaa) | (gr.gfxData&0xaa)>>1) << (8 - gr.xScroll)
	colorBuffer := [4]uint8{a, b, c, d}
	data := gr.gfxData
	gr.displayBuffer[offset+7] = colorBuffer[data&3]
	gr.displayBuffer[offset+6] = colorBuffer[data&3]
	data >>= 2
	gr.displayBuffer[offset+5] = colorBuffer[data&3]
	gr.displayBuffer[offset+4] = colorBuffer[data&3]
	data >>= 2
	gr.displayBuffer[offset+3] = colorBuffer[data&3]
	gr.displayBuffer[offset+2] = colorBuffer[data&3]
	data >>= 2
	gr.displayBuffer[offset+1] = colorBuffer[data]
	gr.displayBuffer[offset] = colorBuffer[data]
}
