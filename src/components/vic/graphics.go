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
	lineOffset        int     // Offset from chunky bitmap buffer
	matrixLineIndex   int     // Index in matrix/colorLine
	//matrixBufferIndex int     // Index in matrixBuffer
	//matrixBuffer      []uint8
	matrixLine       []uint8 // Buffer for video line, read in Bad Lines
	colorLine        []uint8 // Buffer for color line, read in Bad Lines
	rowCounter       uint16  // Row counter
	videoCounter     uint16  // Video counter
	videoCounterBase uint16  // Video counter base
	borderULOn       bool    // Upper/lower border on
	displayOn        bool    // Display state
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
		lineOffset:        0,
		matrixLineIndex:   0,
		matrixLine:        make([]uint8, 40),
		colorLine:         make([]uint8, 40),
		rowCounter:        7,
		videoCounter:      0,
		videoCounterBase:  0,
		displayOn:         false,
	}
	return gr
}

func (gr *Graphics) Setup() {
}

func (gr *Graphics) SetDisplayOn() {
	gr.displayOn = true
}

func (gr *Graphics) ResetRowCounter() {
	gr.rowCounter = 0
}

func (gr *Graphics) ResetVideoCounterBase() {
	gr.videoCounterBase = 0
}

func (gr *Graphics) ResetMatrixLineIndex() {
	gr.matrixLineIndex = 0
}

func (gr *Graphics) UpdateVideoCounter() {
	gr.videoCounter = gr.videoCounterBase
}

func (gr *Graphics) SetLineOffset(lineStart int) {
	gr.lineOffset = lineStart
}

func (gr *Graphics) UpdateBorderUpperLower() {
	if gr.core.rasterY == gr.core.dyStop {
		gr.borderULOn = true
	} else if (gr.core.cr1&0x10) != 0 && gr.core.rasterY == gr.core.dyStart {
		gr.borderULOn = false
	}
}

func (gr *Graphics) UpdateDisplayOn() {
	if gr.rowCounter == 7 {
		gr.videoCounterBase = gr.videoCounter
		gr.displayOn = false
	}
	if gr.core.isBadLine || gr.displayOn {
		gr.rowCounter = (gr.rowCounter + 1) & 7
		gr.displayOn = true
	}
}

func (gr *Graphics) SetCharDataLast() {
	gr.charDataLast = gr.charData
}

func (gr *Graphics) SetBorderOn() {
	gr.borderOn = true
}

func (gr *Graphics) SetBorderOnSample(idx int) {
	gr.borderOnSample[idx] = gr.borderOn
}

func (gr *Graphics) GraphicsAccess() {
	if gr.displayOn {
		var addr uint16
		if (gr.core.cr1 & 0x20) != 0 {
			addr = ((gr.videoCounter & 0x03ff) << 3) | gr.core.bitmapBase | gr.rowCounter // Bitmap
		} else {
			addr = (uint16(gr.matrixLine[gr.matrixLineIndex]) << 3) | gr.core.charBase | gr.rowCounter // Text
		}
		if (gr.core.cr1 & 0x40) != 0 {
			addr &= 0xf9ff // ECM
		}
		gr.gfxData = gr.core.ReadByte(addr)
		gr.charData = gr.matrixLine[gr.matrixLineIndex]
		gr.colorData = gr.colorLine[gr.matrixLineIndex]
		//gr.matrixBuffer[gr.matrixBufferIndex] = gr.charData + 64
		gr.matrixLineIndex++
		gr.videoCounter++
		//gr.matrixBufferIndex++
	} else {
		if (gr.core.cr1 & 0x40) != 0 {
			gr.gfxData = gr.core.ReadByte(0x39ff)
		} else {
			gr.gfxData = gr.core.ReadByte(0x3fff)
		}
		gr.colorData = 0
		gr.charData = 0
	}
}

func (gr *Graphics) MatrixAccess() {
	if gr.core.baLow {
		if gr.core.aecLow {
			addr := (gr.videoCounter & 0x03ff) | gr.core.matrixBase
			gr.matrixLine[gr.matrixLineIndex] = gr.core.ReadByte(addr)
			//gr.matrixBuffer[gr.matrixBufferIndex] = data + 64
			//TODO screen codes
			//https://sta.c64.org/cbm64scr.html
			//if p := gr.matrixLine[gr.matrixLineIndex]; p != 32 {
			//	fmt.Printf("%s\n", string(p+64))
			//}
			gr.colorLine[gr.matrixLineIndex] = gr.core.banks.ReadColor(addr & 0x03ff)
		} else {
			gr.colorLine[gr.matrixLineIndex] = 0xff
			gr.matrixLine[gr.matrixLineIndex] = 0xff
		}
	}
}

func (gr *Graphics) SampleBorder(cycle int) {
	if gr.borderOn {
		idx := cycle - 13
		gr.borderColorSample[idx&DisplayXFill] = gr.core.ecColor
	}
}

func (gr *Graphics) IncrementOffset() {
	gr.lineOffset += 8
	gr.foreMask.Increment()
}

func (gr *Graphics) BorderUpdate() {
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

func (gr *Graphics) Draw(b bool) {
	if gr.borderULOn {
		gr.DrawBackground()
	} else {
		//if b {
		//gr.DrawBackground()
		//}
		gr.drawGraphics()
	}
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

func (gr *Graphics) DrawBackground() {
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
	gr.db.SetMulti8(gr.lineOffset, c)
}

func (gr *Graphics) drawGraphics() {
	offset := gr.lineOffset + int(gr.core.xScroll)
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
	//foreMask
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
	//foreMask
	p1 := ((gr.gfxData & 0xaa) | (gr.gfxData&0xaa)>>1) >> gr.core.xScroll
	p2 := ((gr.gfxData & 0xaa) | (gr.gfxData&0xaa)>>1) << (8 - gr.core.xScroll)
	gr.foreMask.Update(p1, p2)

	colorBuffer := [4]uint8{a, b, c, d}
	data := gr.gfxData
	color := colorBuffer[data&3]
	gr.db.Set(offset+7, color)
	gr.db.Set(offset+6, color)
	data >>= 2
	color = colorBuffer[data&3]
	gr.db.Set(offset+5, color)
	gr.db.Set(offset+4, color)
	data >>= 2
	color = colorBuffer[data&3]
	gr.db.Set(offset+3, color)
	gr.db.Set(offset+2, color)
	data >>= 2
	color = colorBuffer[data]
	gr.db.Set(offset+1, color)
	gr.db.Set(offset, color)
}
