package vic

import (
	"github.com/markel1974/c64emu/src/board/iboard"
	"github.com/markel1974/c64emu/src/preferences"
)

//https://dustlayer.com/c64-architecture

var _emptyForeMaskBuffer = make([]uint8, DisplayXDiv8)

type MOS6569 struct {
	*Core
	board            iboard.IBoard
	prefs            *preferences.Prefs
	sprites          *Sprites
	graphics         *Graphics
	cycle            int      // Cycle
	lastByte         uint8    // Last byte read by VIC
	lineStart        int      // Offset from current line in bitmap buffer
	lineOffset       int      // Offset from chunky bitmap buffer
	matrixLine       []uint8  // Buffer for video line, read in Bad Lines
	colorLine        []uint8  // Buffer for color line, read in Bad Lines
	rasterY          uint16   // Current raster line
	dyStart          uint16   // Comparison values for border logic
	dyStop           uint16   // Comparison values for border logic
	rowCounter       uint16   // Row counter
	videoCounter     uint16   // Video counter
	videoCounterBase uint16   // Video counter base
	borderOnUL       bool     // Upper/lower border on
	displayOn        bool     // Display state
	badLinesEnabled  bool     // Bad Lines enabled for this frame
	lpTriggered      bool     // LightPen was triggered in this frame
	isBadLine        bool     // Current line is bad line
	drawThisLine     bool     // This line is drawn
	sprDmaOn         uint8    // 8 flags: Sprite DMA active
	sprMC            []uint16 // Sprite data counters
	sprMCBase        []uint16 // Sprite data counter bases
	sprDisplayOn     uint8    // 8 flags: Sprite display active
	refreshCounter   uint8    // Refresh counter
	rasterX          uint16   // Current raster x position
	mlIndex          int      // Index in matrix/colorLine[]
	firstBaCycle     uint64   //
	vBlanking        bool     // VBlank in next cycle
	baLow            bool     // BALine
	ready            bool     // VIC Initialization Complete
}

func NewMOS6569() *MOS6569 {
	core := NewCore()
	vic := &MOS6569{
		Core:             core,
		graphics:         NewGraphics(core),
		sprites:          NewSprites(core),
		ready:            false,
		lineOffset:       0,
		isBadLine:        false,
		baLow:            false,
		badLinesEnabled:  false,
		rasterY:          TotalRasters - 1,
		rowCounter:       7,
		videoCounter:     0,
		videoCounterBase: 0,
		dyStart:          Row24YStart,
		dyStop:           Row24YStop,
		mlIndex:          0,
		cycle:            1,
		displayOn:        false,
		vBlanking:        false,
		lpTriggered:      false,
		drawThisLine:     false,
		sprDisplayOn:     0,
		sprDmaOn:         0,
		matrixLine:       make([]uint8, 40),
		colorLine:        make([]uint8, 40),
		sprMC:            make([]uint16, SpriteNumber),
		sprMCBase:        make([]uint16, SpriteNumber),
	}
	for i := range vic.sprMC {
		vic.sprMC[i] = 63
	}
	return vic
}

func (vic *MOS6569) Setup(board iboard.IBoard, prefs *preferences.Prefs) {
	vic.board = board
	vic.prefs = prefs
	vic.graphics.Setup(board)
	vic.sprites.Setup(board)
}

func (vic *MOS6569) GetDisplayBuffer() []uint8 {
	return vic.displayBuffer
}

func (vic *MOS6569) Reset() {
	vic.ready = false
}

func (vic *MOS6569) GetLastByte() uint8 {
	return vic.lastByte
}

func (vic *MOS6569) GetBALow() bool {
	return vic.baLow
}

func (vic *MOS6569) setBALow() {
	if !vic.baLow {
		vic.baLow = true
		vic.firstBaCycle = vic.board.Cycle()
	}
}

func (vic *MOS6569) clearBALow() {
	vic.baLow = false
}

func (vic *MOS6569) refreshAccess() {
	_ = vic.readByte(0x3f00 | uint16(vic.refreshCounter))
	vic.refreshCounter--
}

func (vic *MOS6569) sampleBorder() {
	vic.graphics.SetBorderColorSample(vic.cycle)
	vic.lineOffset += 8
	vic.foreMaskOffset++
}

func (vic *MOS6569) NewPrefs(_ *preferences.Prefs) {
	//vic.skipFrames = prefs.SkipFrames()
}

func (vic *MOS6569) rasterIrq() {
	vic.irqFlag |= 0x01
	if (vic.irqMask & 0x01) != 0 {
		vic.irqFlag |= 0x80
		vic.board.VICTriggerIRQ()
	}
}

func (vic *MOS6569) readByte(addr uint16) uint8 {
	va := addr | vic.ciaVaBase
	if (va & 0x7000) == 0x1000 {
		vic.lastByte = vic.board.CharRomRead(va & 0x0fff)
		return vic.lastByte
	}
	vic.lastByte = vic.board.RamRead(va)
	return vic.lastByte
}

func (vic *MOS6569) checkSpriteDMA() {
	rasterY := vic.rasterY & 0xff
	for i, mask := 0, uint8(1); i < 8; i, mask = i+1, mask<<1 {
		if (vic.me&mask) != 0 && rasterY == uint16(vic.mXy[i]) {
			vic.sprDmaOn |= mask
			vic.sprMCBase[i] = 0
			if (vic.mye & mask) != 0 {
				vic.sprExpY &= ^mask
			}
		}
	}
}

func (vic *MOS6569) fetchSpriteDataPtr(num int) {
	addr := vic.matrixBase | 0x03f8 | uint16(num)
	data := vic.readByte(addr)
	ptr := uint16(data) << 6
	vic.sprites.SetSpritePtr(num, ptr)
}

func (vic *MOS6569) fetchSpriteData(num int, byteNum int) {
	if (vic.sprDmaOn & (1 << num)) != 0 {
		ptr := vic.sprites.GetSpritePtr(num)
		addr := (vic.sprMC[num] & 0x3f) | ptr
		data := vic.readByte(addr)
		vic.sprites.SetSpriteData(num, byteNum, data)
		vic.sprMC[num]++
	} else if byteNum == 1 {
		//idleAccess
		vic.readByte(0x3fff)
	}
}

func (vic *MOS6569) Emulate() (bool, bool) {
	vBlank := false
	lastCycle := false

	switch vic.cycle {
	case 1:
		// Fetch sprite pointer 3, increment raster counter, trigger raster IRQ,
		// test for Bad Line, reset BA if sprites 3 and 4 off, read data of sprite 3
		if vic.rasterY == RasterYMax {
			// Trigger VBlank in cycle 2
			vic.vBlanking = true
		} else {
			// Increment raster counter
			vic.rasterY++
			// Trigger raster IRQ if IRQ line reached
			if vic.rasterY == vic.irqRaster {
				vic.rasterIrq()
			}
			// In line $30, the DEN bit controls if Bad Lines can occur
			if vic.rasterY == 0x30 {
				vic.badLinesEnabled = (vic.cr1 & 0x10) != 0
			}
			// Bad Line condition?
			vic.isBadLine = (vic.rasterY >= FirstDmaLine) && (vic.rasterY <= LastDmaLine) && ((vic.rasterY & 7) == vic.yScroll) && (vic.badLinesEnabled)
			// Don't draw all lines, hide some at the top and bottom
			vic.drawThisLine = (vic.rasterY >= FirstDisplayedLine) && (vic.rasterY <= LastDisplayedLine)
		}
		// First sample of border state
		vic.graphics.SetBorderOnSample(0)
		vic.fetchSpriteDataPtr(3)
		vic.fetchSpriteData(3, 0)
		if vic.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDmaOn & 0x18) == 0 {
			vic.clearBALow()
		}

	case 2:
		// Set BA for sprite 5, read data of sprite 3
		if vic.vBlanking {
			vBlank = true
			// Vertical blank, reset counters
			vic.videoCounterBase = 0
			vic.rasterY = 0
			vic.refreshCounter = 0xff
			vic.vBlanking = false
			vic.lpTriggered = false
			vic.lineStart = 0
			if vic.irqRaster == 0 {
				// Trigger raster IRQ if IRQ in line 0
				vic.rasterIrq()
			}
		}
		// Our output goes here
		vic.lineOffset = vic.lineStart
		// Clear foreground mask
		copy(vic.foreMaskBuf, _emptyForeMaskBuffer)

		vic.foreMaskOffset = 0
		vic.fetchSpriteData(3, 1)
		vic.fetchSpriteData(3, 2)
		if vic.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDmaOn & 0x20) != 0 {
			vic.setBALow()
		}

	case 3:
		// Fetch sprite pointer 4, reset BA is sprite 4 and 5 off
		vic.fetchSpriteDataPtr(4)
		vic.fetchSpriteData(4, 0)
		if vic.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDmaOn & 0x30) == 0 {
			vic.clearBALow()
		}

	case 4:
		// Set BA for sprite 6, read data of sprite 4
		vic.fetchSpriteData(4, 1)
		vic.fetchSpriteData(4, 2)
		if vic.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDmaOn & 0x40) != 0 {
			vic.setBALow()
		}

	case 5:
		// Fetch sprite pointer 5, reset BA if sprite 5 and 6 off
		vic.fetchSpriteDataPtr(5)
		vic.fetchSpriteData(5, 0)
		if vic.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDmaOn & 0x60) == 0 {
			vic.clearBALow()
		}

	case 6:
		// Set BA for sprite 7, read data of sprite 5
		vic.fetchSpriteData(5, 1)
		vic.fetchSpriteData(5, 2)
		if vic.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDmaOn & 0x80) != 0 {
			vic.setBALow()
		}

	case 7:
		// Fetch sprite pointer 6, reset BA if sprite 6 and 7 off
		vic.fetchSpriteDataPtr(6)
		vic.fetchSpriteData(6, 0)
		if vic.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDmaOn & 0xc0) == 0 {
			vic.clearBALow()
		}

	case 8:
		// Read data of sprite 6
		vic.fetchSpriteData(6, 1)
		vic.fetchSpriteData(6, 2)
		if vic.isBadLine {
			vic.displayOn = true
		}

	case 9:
		// Fetch sprite pointer 7, reset BA if sprite 7 off
		vic.fetchSpriteDataPtr(7)
		vic.fetchSpriteData(7, 0)
		if vic.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDmaOn & 0x80) == 0 {
			vic.clearBALow()
		}

	case 10:
		// Read data of sprite 7
		vic.fetchSpriteData(7, 1)
		vic.fetchSpriteData(7, 2)
		if vic.isBadLine {
			vic.displayOn = true
		}

	case 11:
		// Refresh, reset BA
		vic.refreshAccess()
		if vic.isBadLine {
			vic.displayOn = true
		}
		vic.clearBALow()

	case 12:
		// Refresh, turn on matrix access if Bad Line
		vic.refreshAccess()
		// Turn on display and matrix access if Bad Line
		if vic.isBadLine {
			vic.displayOn = true
			vic.setBALow()
		}

	case 13:
		// Refresh, turn on matrix access if Bad Line, reset rasterX, graphics display starts here
		if vic.drawThisLine {
			vic.graphics.DrawBackground(vic.lineOffset)
			vic.sampleBorder()
		}
		vic.refreshAccess()
		// Turn on display and matrix access if Bad Line
		if vic.isBadLine {
			vic.displayOn = true
			vic.setBALow()
		}
		vic.rasterX = 0xfffc

	case 14:
		// Refresh, videoCounter -> videoCounterBase, turn on matrix access and reset RC if Bad Line
		if vic.drawThisLine {
			vic.graphics.DrawBackground(vic.lineOffset)
			vic.sampleBorder()
		}
		vic.refreshAccess()
		// Turn on display and matrix access and reset RC if Bad Line
		if vic.isBadLine {
			vic.rowCounter = 0
			vic.displayOn = true
			vic.setBALow()
		}
		vic.videoCounter = vic.videoCounterBase

	case 15:
		// Refresh and matrix access, increment sprMCBase by 2 if y expansion FlipFlop is set
		if vic.drawThisLine {
			vic.graphics.DrawBackground(vic.lineOffset)
			vic.sampleBorder()
		}
		vic.refreshAccess()
		// Turn on display and matrix access if Bad Line
		if vic.isBadLine {
			vic.displayOn = true
			vic.setBALow()
		}
		for idx := 0; idx < 8; idx++ {
			if (vic.sprExpY & (1 << idx)) != 0 {
				vic.sprMCBase[idx] += 2
			}
		}
		vic.mlIndex = 0
		vic.matrixAccess()

	case 16:
		// Graphics and matrix access, increment sprMCBase by 1 if y expansion FlipFlop is set
		// and check if sprite DMA can be turned off
		if vic.drawThisLine {
			vic.graphics.DrawBackground(vic.lineOffset)
			vic.sampleBorder()
		}
		vic.graphicsAccess()
		// Turn on display and matrix access if Bad Line
		if vic.isBadLine {
			vic.displayOn = true
			vic.setBALow()
		}
		for idx, mask := 0, uint8(1); idx < 8; idx, mask = idx+1, mask<<1 {
			if (vic.sprExpY & mask) != 0 {
				vic.sprMCBase[idx]++
			}
			if (vic.sprMCBase[idx] & 0x3f) == 0x3f {
				vic.sprDmaOn &= ^mask
			}
		}
		vic.matrixAccess()

	case 17:
		// Graphics and matrix access, turn off border in 40 column mode, display window starts here
		if (vic.cr2 & 8) != 0 {
			if vic.rasterY == vic.dyStop {
				vic.borderOnUL = true
			} else {
				if (vic.cr1 & 0x10) != 0 {
					if vic.rasterY == vic.dyStart {
						vic.borderOnUL = false
						vic.graphics.SetBorderOn(false)
					} else if !vic.borderOnUL {
						vic.graphics.SetBorderOn(false)
					}
				} else if !vic.borderOnUL {
					vic.graphics.SetBorderOn(false)
				}
			}
		}
		// Second sample of border state
		vic.graphics.SetBorderOnSample(1)
		if vic.drawThisLine {
			if vic.borderOnUL {
				vic.graphics.DrawBackground(vic.lineOffset)
			} else {
				vic.graphics.DrawBackground(vic.lineOffset)
				vic.graphics.DrawGraphics(vic.lineOffset)
			}
			vic.sampleBorder()
		}
		vic.graphicsAccess()
		// Turn on display and matrix access if Bad Line
		if vic.isBadLine {
			vic.displayOn = true
			vic.setBALow()
		}
		vic.matrixAccess()

	case 18:
		// Turn off border in 38 column mode
		if (vic.cr2 & 8) == 0 {
			if vic.rasterY == vic.dyStop {
				vic.borderOnUL = true
			} else {
				if (vic.cr1 & 0x10) != 0 {
					if vic.rasterY == vic.dyStart {
						vic.borderOnUL = false
						vic.graphics.SetBorderOn(false)
					} else if !vic.borderOnUL {
						vic.graphics.SetBorderOn(false)
					}
				} else {
					if !vic.borderOnUL {
						vic.graphics.SetBorderOn(false)
					}
				}
			}
		}

		// Third sample of border state
		vic.graphics.SetBorderOnSample(2)
		if vic.drawThisLine {
			if vic.borderOnUL {
				vic.graphics.DrawBackground(vic.lineOffset)
			} else {
				vic.graphics.DrawGraphics(vic.lineOffset)
			}
			vic.sampleBorder()
		}

		vic.graphicsAccess()
		// Turn on display and matrix access if Bad Line
		if vic.isBadLine {
			vic.displayOn = true
			vic.setBALow()
		}
		vic.matrixAccess()

		vic.graphics.SetCharDataLast()

	case 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54:
		if vic.drawThisLine {
			if vic.borderOnUL {
				vic.graphics.DrawBackground(vic.lineOffset)
			} else {
				vic.graphics.DrawGraphics(vic.lineOffset)
			}
			vic.sampleBorder()
		}
		vic.graphicsAccess()
		// Turn on display and matrix access if Bad Line
		if vic.isBadLine {
			vic.displayOn = true
			vic.setBALow()
		}
		vic.matrixAccess()
		vic.graphics.SetCharDataLast()

	case 55:
		// Last graphics access, turn off matrix access, turn on sprite DMA if Y coordinate is
		// right and sprite is enabled, handle sprite y expansion, set BA for sprite 0
		if vic.drawThisLine {
			if vic.borderOnUL {
				vic.graphics.DrawBackground(vic.lineOffset)
			} else {
				vic.graphics.DrawGraphics(vic.lineOffset)
			}
			vic.sampleBorder()
		}
		vic.graphicsAccess()
		if vic.isBadLine {
			vic.displayOn = true
		}
		// Invert y expansion FlipFlop (if MYE bit is set)
		for idx, mask := 0, uint8(1); idx < 8; idx, mask = idx+1, mask<<1 {
			if (vic.mye & mask) != 0 {
				vic.sprExpY ^= mask
			}
		}
		vic.checkSpriteDMA()
		if (vic.sprDmaOn & 0x01) != 0 {
			vic.setBALow()
		} else {
			vic.clearBALow()
		}

	case 56:
		// Turn on border in 38 column mode, turn on sprite DMA if Y coordinate is right and
		// sprite is enabled, set BA for sprite 0, display window ends here
		if (vic.cr2 & 8) == 0 {
			vic.graphics.SetBorderOn(true)
		}
		// Fourth sample of border state
		vic.graphics.SetBorderOnSample(3)
		if vic.drawThisLine {
			if vic.borderOnUL {
				vic.graphics.DrawBackground(vic.lineOffset)
			} else {
				vic.graphics.DrawGraphics(vic.lineOffset)
			}
			vic.sampleBorder()
		}
		//idleAccess
		vic.readByte(0x3fff)
		if vic.isBadLine {
			vic.displayOn = true
		}
		vic.checkSpriteDMA()
		if (vic.sprDmaOn & 0x01) != 0 {
			vic.setBALow()
		}

	// Turn on border in 40 column mode, set BA for sprite 1, paint sprites
	case 57:
		if (vic.cr2 & 8) != 0 {
			vic.graphics.SetBorderOn(true)
		}
		// Fifth sample of border state
		vic.graphics.SetBorderOnSample(4)
		vic.sprites.SetDisplayOn(vic.sprDisplayOn)
		// Turn off sprite display if DMA is off
		for idx, mask := 0, uint8(1); idx < 8; idx, mask = idx+1, mask<<1 {
			if (vic.sprDisplayOn&mask) != 0 && (vic.sprDmaOn&mask) == 0 {
				vic.sprDisplayOn &= ^mask
			}
		}
		if vic.drawThisLine {
			vic.graphics.DrawBackground(vic.lineOffset)
			vic.sampleBorder()
		}
		//idleAccess
		vic.readByte(0x3fff)
		if vic.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDmaOn & 0x02) != 0 {
			vic.setBALow()
		}
	case 58:
		// Fetch sprite pointer 0, sprMCBase->sprMC, turn on sprite display if necessary,
		// turn off display if RC=7, read data of sprite 0
		if vic.drawThisLine {
			vic.graphics.DrawBackground(vic.lineOffset)
			vic.sampleBorder()
		}
		rasterY := vic.rasterY & 0xff
		for idx, mask := 0, uint8(1); idx < 8; idx, mask = idx+1, mask<<1 {
			vic.sprMC[idx] = vic.sprMCBase[idx]
			if (vic.sprDmaOn&mask) != 0 && (rasterY == uint16(vic.mXy[idx])) {
				vic.sprDisplayOn |= mask
			}
		}
		vic.fetchSpriteDataPtr(0)
		vic.fetchSpriteData(0, 0)
		if vic.rowCounter == 7 {
			vic.videoCounterBase = vic.videoCounter
			vic.displayOn = false
		}
		if vic.isBadLine || vic.displayOn {
			vic.displayOn = true
			vic.rowCounter = (vic.rowCounter + 1) & 7
		}
	case 59:
		// Set BA for sprite 2, read data of sprite 0
		if vic.drawThisLine {
			vic.graphics.DrawBackground(vic.lineOffset)
			vic.sampleBorder()
		}
		vic.fetchSpriteData(0, 1)
		vic.fetchSpriteData(0, 2)
		if vic.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDmaOn & 0x04) != 0 {
			vic.setBALow()
		}

	// Fetch sprite pointer 1, reset BA if sprite 1 and 2 off, graphics display ends here
	case 60:
		if vic.drawThisLine {
			vic.graphics.DrawBackground(vic.lineOffset)
			vic.sampleBorder()
			vic.sprites.Draw(vic.lineStart)
			vic.graphics.DrawBorder(vic.lineStart)
			// Increment pointer in chunky buffer
			vic.lineStart += DisplayX
		}
		vic.fetchSpriteDataPtr(1)
		vic.fetchSpriteData(1, 0)
		if vic.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDmaOn & 0x06) == 0 {
			vic.clearBALow()
		}

	// Set BA for sprite 3, read data of sprite 1
	case 61:
		vic.fetchSpriteData(1, 1)
		vic.fetchSpriteData(1, 2)
		if vic.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDmaOn & 0x08) != 0 {
			vic.setBALow()
		}

	// Read sprite pointer 2, reset BA if sprite 2 and 3 off, read data of sprite 2
	case 62:
		vic.fetchSpriteDataPtr(2)
		vic.fetchSpriteData(2, 0)
		if vic.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDmaOn & 0x0c) == 0 {
			vic.clearBALow()
		}

	// Set BA for sprite 4, read data of sprite 2
	case 63:
		vic.fetchSpriteData(2, 1)
		vic.fetchSpriteData(2, 2)
		if vic.isBadLine {
			vic.displayOn = true
		}
		if vic.rasterY == vic.dyStop {
			vic.borderOnUL = true
		} else if (vic.cr1&0x10) != 0 && vic.rasterY == vic.dyStart {
			vic.borderOnUL = false
		}
		if (vic.sprDmaOn & 0x10) != 0 {
			vic.setBALow()
		}
		lastCycle = true
	}
	vic.rasterX += 8
	if lastCycle {
		vic.cycle = 1
	} else {
		vic.cycle++
	}
	return vBlank, lastCycle
}

func (vic *MOS6569) TriggerLightPen() {
	// LightPen triggers only once per frame
	if !vic.lpTriggered {
		vic.lpTriggered = true
		vic.lpx = uint8(vic.rasterX >> 1) // Latch current coordinates
		vic.lpy = uint8(vic.rasterY)
		vic.irqFlag |= 0x08 // Trigger IRQ
		if (vic.irqMask & 0x08) != 0 {
			vic.irqFlag |= 0x80
			vic.board.VICTriggerIRQ()
		}
	}
}

func (vic *MOS6569) ChangedVA(newVA uint8) {
	vic.ciaVaBase = uint16(newVA) << 14
	vic.WriteRegister(0x18, vic.vaBase) // Force update of memory pointers
}

func (vic *MOS6569) matrixAccess() {
	if vic.baLow {
		if vic.board.Cycle()-vic.firstBaCycle < 3 {
			vic.colorLine[vic.mlIndex] = 0xff
			vic.matrixLine[vic.mlIndex] = 0xff
		} else {
			addr := (vic.videoCounter & 0x03ff) | vic.matrixBase
			vic.matrixLine[vic.mlIndex] = vic.readByte(addr)
			vic.colorLine[vic.mlIndex] = vic.board.ColorRead(addr & 0x03ff)
		}
	}
}

func (vic *MOS6569) graphicsAccess() {
	if vic.displayOn {
		var addr uint16
		if (vic.cr1 & 0x20) != 0 {
			addr = ((vic.videoCounter & 0x03ff) << 3) | vic.bitmapBase | vic.rowCounter // Bitmap
		} else {
			addr = (uint16(vic.matrixLine[vic.mlIndex]) << 3) | vic.charBase | vic.rowCounter // Text
		}
		if (vic.cr1 & 0x40) != 0 {
			addr &= 0xf9ff // ECM
		}
		vic.graphics.SetGfxData(vic.readByte(addr))
		vic.graphics.SetCharData(vic.matrixLine[vic.mlIndex])
		vic.graphics.SetColorData(vic.colorLine[vic.mlIndex])
		vic.mlIndex++
		vic.videoCounter++
	} else {
		if (vic.cr1 & 0x40) != 0 {
			vic.graphics.SetGfxData(vic.readByte(0x39ff))
		} else {
			vic.graphics.SetGfxData(vic.readByte(0x3fff))
		}
		vic.graphics.SetColorData(0)
		vic.graphics.SetCharData(0)
	}
}

func (vic *MOS6569) ReadRegister(addr uint16) uint8 {
	addr = addr & 0x3f
	switch addr {
	case 0x00, 0x02, 0x04, 0x06, 0x08, 0x0a, 0x0c, 0x0e:
		return uint8(vic.mXx[addr>>1])
	case 0x01, 0x03, 0x05, 0x07, 0x09, 0x0b, 0x0d, 0x0f:
		return vic.mXy[addr>>1]
	// Sprite X position MSB
	case 0x10:
		return vic.mx8
	// Control register 1
	case 0x11:
		return uint8((uint16(vic.cr1) & 0x7f) | ((vic.rasterY & 0x100) >> 1))
	// Raster counter
	case 0x12:
		return uint8(vic.rasterY)
	// Light pen X
	case 0x13:
		return vic.lpx
	// Light pen Y
	case 0x14:
		return vic.lpy
	// Sprite enable
	case 0x15:
		return vic.me
	// Control register 2
	case 0x16:
		return vic.cr2 | 0xc0
	// Sprite Y expansion
	case 0x17:
		return vic.mye
	// Memory pointers
	case 0x18:
		return vic.vaBase | 0x01
	// IRQ flags
	case 0x19:
		if !vic.ready {
			vic.ready = true
			vic.board.ReadyEvent()
		}
		return vic.irqFlag | 0x70
	// IRQ mask
	case 0x1a:
		return vic.irqMask | 0xf0
	// Sprite data priority
	case 0x1b:
		return vic.mdp
	// Sprite multicolor
	case 0x1c:
		return vic.mmc
	// Sprite X expansion
	case 0x1d:
		return vic.mxe
	// Sprite-sprite collision
	case 0x1e:
		ret := vic.sprClx
		vic.sprClx = 0 // Read and clear
		return ret
	// Sprite-background collision
	case 0x1f:
		ret := vic.sprClxBgr
		vic.sprClxBgr = 0 // Read and clear
		return ret
	case 0x20:
		return vic.ec | 0xf0
	case 0x21:
		return vic.b0c | 0xf0
	case 0x22:
		return vic.b1c | 0xf0
	case 0x23:
		return vic.b2c | 0xf0
	case 0x24:
		return vic.b3c | 0xf0
	case 0x25:
		return vic.mm0 | 0xf0
	case 0x26:
		return vic.mm1 | 0xf0
	case 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e:
		return vic.mXc[addr-0x27] | 0xf0
	default:
		return 0xff
	}
}

func (vic *MOS6569) WriteRegister(addr uint16, data uint8) {
	addr = addr & 0x3f
	switch addr {
	case 0x00, 0x02, 0x04, 0x06, 0x08, 0x0a, 0x0c, 0x0e:
		target := addr >> 1
		vic.mXx[target] = (vic.mXx[target] & 0xff00) | uint16(data)
	case 0x10:
		vic.mx8 = data
		for i, j := 0, uint8(1); i < 8; i, j = i+1, j<<1 {
			if (vic.mx8 & j) != 0 {
				vic.mXx[i] |= 0x100
			} else {
				vic.mXx[i] &= 0xff
			}
		}
	case 0x01, 0x03, 0x05, 0x07, 0x09, 0x0b, 0x0d, 0x0f:
		vic.mXy[addr>>1] = data
	case 0x11: // Control register 1
		vic.cr1 = data
		vic.yScroll = uint16(data) & 7
		newIRQRaster := (vic.irqRaster & 0xff) | ((uint16(data) & 0x80) << 1)
		if vic.irqRaster != newIRQRaster && vic.rasterY == newIRQRaster {
			vic.rasterIrq()
		}
		vic.irqRaster = newIRQRaster
		if (data & 8) != 0 {
			vic.dyStart = Row25YStart
			vic.dyStop = Row25YStop
		} else {
			vic.dyStart = Row24YStart
			vic.dyStop = Row24YStop
		}
		// In line $30, the DEN bit controls if Bad Lines can occur
		if (vic.rasterY == 0x30) && ((data & 0x10) != 0) {
			vic.badLinesEnabled = true
		}
		// Bad Line condition?
		vic.isBadLine = vic.rasterY >= FirstDmaLine && vic.rasterY <= LastDmaLine && ((vic.rasterY & 7) == vic.yScroll) && vic.badLinesEnabled
		vic.displayIdx = ((int(vic.cr1) & 0x60) | (int(vic.cr2) & 0x10)) >> 4
	case 0x12: // Raster counter
		newIRQRaster := (vic.irqRaster & 0xff00) | uint16(data)
		if vic.irqRaster != newIRQRaster && vic.rasterY == newIRQRaster {
			vic.rasterIrq()
		}
		vic.irqRaster = newIRQRaster
	case 0x15: // Sprite enable
		vic.me = data
	case 0x16: // Control register 2
		vic.cr2 = data
		vic.xScroll = uint16(data) & 7
		vic.displayIdx = ((int(vic.cr1) & 0x60) | (int(vic.cr2) & 0x10)) >> 4
	case 0x17: // Sprite Y expansion
		vic.mye = data
		vic.sprExpY |= ^data
	case 0x18: // Memory pointers
		vic.vaBase = data
		vic.matrixBase = (uint16(data) & 0xf0) << 6
		vic.charBase = (uint16(data) & 0x0e) << 10
		vic.bitmapBase = (uint16(data) & 0x08) << 10
	case 0x19: // IRQ flags
		vic.irqFlag = vic.irqFlag & (^data & 0x0f)
		if (vic.irqFlag & vic.irqMask) != 0 {
			// Set master bit if allowed interrupt still pending
			vic.irqFlag |= 0x80
		} else {
			vic.board.VICClearIRQ()
		}
	case 0x1a: // IRQ mask
		vic.irqMask = data & 0x0f
		if (vic.irqFlag & vic.irqMask) != 0 {
			// Trigger interrupt if pending and now allowed
			vic.irqFlag |= 0x80
			vic.board.VICTriggerIRQ()
		} else {
			vic.irqFlag &= 0x7f
			vic.board.VICClearIRQ()
		}
	case 0x1b: // Sprite data priority
		vic.mdp = data
	case 0x1c: // Sprite multicolor
		vic.mmc = data
	case 0x1d: // Sprite X expansion
		vic.mxe = data
	case 0x20:
		vic.ec = data
		vic.ecColor = vic.colors[data]
	case 0x21:
		vic.b0c = data
		vic.b0cColor = vic.colors[data]
	case 0x22:
		vic.b1c = data
		vic.b1cColor = vic.colors[data]
	case 0x23:
		vic.b2c = data
		vic.b2cColor = vic.colors[data]
	case 0x24:
		vic.b3c = data
		vic.b3cColor = vic.colors[data]
	case 0x25:
		vic.mm0 = data
		vic.mm0Color = vic.colors[data]
	case 0x26:
		vic.mm1 = data
		vic.mm1Color = vic.colors[data]
	case 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e:
		target := addr - 0x27
		vic.mXc[target] = data
		vic.mXcColor[target] = vic.colors[data]
	}
}
