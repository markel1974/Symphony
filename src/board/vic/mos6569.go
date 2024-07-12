package vic

import (
	"github.com/markel1974/c64emu/src/board/iboard"
	"github.com/markel1974/c64emu/src/preferences"
)

//https://dustlayer.com/c64-architecture

type MOS6569 struct {
	board          iboard.IBoard
	prefs          *preferences.Prefs
	cycle          int
	lastByte       uint8 // Last byte read by VIC
	displayBuffer  []uint8
	lineStart      int     // Offset from current line in bitmap buffer
	lineOffset     int     // Offset from chunky bitmap buffer
	foreMaskBuf    []uint8 // Foreground mask for sprite-graphics collisions and priorities
	foreMaskOffset int     // Offset from foreMaskBuf

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

	ecColor  uint8    // Index ec Color Mapping
	b0cColor uint8    // Index b0c Color Mapping
	b1cColor uint8    // Index b1c Color Mapping
	b2cColor uint8    // Index b2c Color Mapping
	b3cColor uint8    // Index b3c Color Mapping
	mm0Color uint8    // Index mm0 Color Mapping
	mm1Color uint8    // Index mm1 Color Mapping
	mXcColor [8]uint8 // Indices for m0c - m1c - m2c - m3c - m4c - m5c - m6c - m7c Color Mapping

	irqFlag   uint8
	irqMask   uint8
	irqRaster uint16 // Interrupt raster line

	colors     [256]uint8 // Indices of the 16 colors (16 times mirrored to avoid "& 0x0f")
	matrixLine [40]uint8  // Buffer for video line, read in Bad Lines
	colorLine  [40]uint8  // Buffer for color line, read in Bad Lines

	rasterY          uint16 // Current raster line
	dyStart          uint16 // Comparison values for border logic
	dyStop           uint16 // Comparison values for border logic
	rowCounter       uint16 // Row counter
	videoCounter     uint16 // Video counter
	videoCounterBase uint16 // Video counter base
	xScroll          uint16 // X scroll value
	yScroll          uint16 // Y scroll value

	vaBase     uint8
	ciaVaBase  uint16 // CIA VA14/15 video base
	matrixBase uint16 // Video matrix base
	charBase   uint16 // Character generator base
	bitmapBase uint16 // Bitmap base

	borderColorSample [DisplayXFill + 1]uint8 // Samples of border color at each "displayed" cycle
	borderOn          bool                    // Upper/lower border on (Main border FlipFlop)
	borderOnUL        bool                    // Upper/lower border on
	borderOnSample    [5]bool                 // Samples of border state at different cycles (1, 17, 18, 56, 57)

	displayIdx      int  // Index of current display mode
	displayOn       bool // Display state
	badLinesEnabled bool // Bad Lines enabled for this frame
	lpTriggered     bool // LightPen was triggered in this frame
	isBadLine       bool // Current line is bad line
	drawThisLine    bool // This line is drawn

	refreshCounter uint8 // Refresh counter

	sprClxBgr    uint8     // Sprite to background collision
	sprClx       uint8     // Sprite to sprite collision
	sprCollBuf   []uint8   // Buffer for sprite-sprite collisions and priorities
	sprPtr       [8]uint16 // Sprite data pointers
	sprMC        [8]uint16 // Sprite data counters
	sprMCBase    [8]uint16 // Sprite data counter bases
	sprData      [][]uint8 // Sprite data read
	sprDrawData  [][]uint8 // Sprite data for drawing
	sprExpY      uint8     // 8 sprite y expansion flip flops
	sprDmaOn     uint8     // 8 flags: Sprite DMA active
	sprDisplayOn uint8     // 8 flags: Sprite display active
	sprDraw      uint8     // 8 flags: Draw sprite in this line

	rasterX      uint16 // Current raster x position
	mlIndex      int    // Index in matrix/colorLine[]
	gfxData      uint8
	charData     uint8
	charDataLast uint8
	colorData    uint8
	firstBaCycle uint64
	vBlanking    bool // Flag: VBlank in next cycle
	baLow        bool
	ready        bool // VIC Initialization Complete
}

func NewMOS6569() *MOS6569 {
	vic := &MOS6569{}
	vic.ready = false
	vic.displayBuffer = make([]uint8, DisplaySize)
	vic.lineOffset = 0
	vic.isBadLine = false
	vic.sprExpY = 0
	for i := 0; i < 8; i++ {
		vic.sprMCBase[i] = 0
	}
	vic.sprData = make([][]uint8, 8)
	for i := range vic.sprData {
		vic.sprData[i] = make([]uint8, 4)
	}
	vic.sprDrawData = make([][]uint8, 8)
	for i := range vic.sprDrawData {
		vic.sprDrawData[i] = make([]uint8, 4)
	}
	// Set pointers
	vic.matrixBase = 0
	vic.charBase = 0
	vic.bitmapBase = 0
	// Initialize VIC registers
	vic.mx8 = 0
	vic.cr1 = 0
	vic.cr2 = 0
	vic.lpx = 0
	vic.lpy = 0
	vic.me = 0
	vic.mxe = 0
	vic.mye = 0
	vic.mdp = 0
	vic.mmc = 0
	vic.ec = 0
	vic.b0c = 0
	vic.b1c = 0
	vic.b2c = 0
	vic.b3c = 0
	vic.mm0 = 0
	vic.mm1 = 0
	for i := 0; i < 8; i++ {
		vic.mXx[i] = 0
		vic.mXy[i] = 0
		vic.mXc[i] = 0
	}
	vic.vaBase = 0
	vic.ciaVaBase = 0
	vic.irqFlag = 0
	vic.irqMask = 0
	vic.sprClx = 0
	vic.sprClxBgr = 0

	// Initialize other variables
	vic.rasterY = TotalRasters - 1
	vic.rowCounter = 7
	vic.irqRaster = 0
	vic.videoCounter = 0
	vic.videoCounterBase = 0
	vic.xScroll = 0
	vic.yScroll = 0
	vic.dyStart = Row24YStart
	vic.dyStop = Row24YStop
	vic.mlIndex = 0
	vic.cycle = 1
	vic.displayIdx = 0
	vic.displayOn = false
	vic.borderOn = false
	vic.borderOnUL = false
	vic.vBlanking = false
	vic.lpTriggered = false
	vic.drawThisLine = false
	vic.sprDmaOn = 0
	vic.sprDisplayOn = 0
	for i := 0; i < 8; i++ {
		vic.sprMC[i] = 63
		vic.sprPtr[i] = 0
	}
	vic.sprCollBuf = make([]uint8, DisplayXFillMax)
	copy(vic.sprCollBuf, _displayXFillMax)
	vic.foreMaskBuf = make([]uint8, DisplayXFill+1)
	copy(vic.foreMaskBuf, _displayXDiv8)
	for i := 0; i < 256; i++ {
		vic.colors[i] = (uint8)(i & 0x0f)
	}
	// Preset colors to black
	vic.ecColor = vic.colors[0]
	vic.b0cColor = vic.colors[0]
	vic.b1cColor = vic.colors[0]
	vic.b2cColor = vic.colors[0]
	vic.b3cColor = vic.colors[0]
	vic.mm0Color = vic.colors[0]
	vic.mm1Color = vic.colors[0]
	for i := 0; i < 8; i++ {
		vic.mXcColor[i] = vic.colors[0]
	}
	vic.baLow = false
	vic.badLinesEnabled = false
	return vic
}

func (vic *MOS6569) Setup(board iboard.IBoard, prefs *preferences.Prefs) {
	vic.board = board
	vic.prefs = prefs
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
	vic.readByte(0x3f00 | uint16(vic.refreshCounter))
	vic.refreshCounter--
}

func (vic *MOS6569) checkSpriteDMA() {
	for i, mask := 0, uint8(1); i < 8; i, mask = i+1, mask<<1 {
		if (vic.me&mask) != 0 && (vic.rasterY&0xff) == uint16(vic.mXy[i]) {
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
	vic.sprPtr[num] = uint16(vic.readByte(addr)) << 6
}

func (vic *MOS6569) fetchSpriteData(num int, byteNum int) {
	if (vic.sprDmaOn & (1 << num)) != 0 {
		vic.sprData[num][byteNum] = vic.readByte((vic.sprMC[num] & 0x3f) | vic.sprPtr[num])
		vic.sprMC[num]++
	} else if byteNum == 1 {
		//idleAccess
		vic.readByte(0x3fff)
	}
}

func (vic *MOS6569) sampleBorder() {
	if vic.borderOn {
		idx := vic.cycle - 13
		vic.borderColorSample[idx&DisplayXFill] = vic.ecColor
	}
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

func (vic *MOS6569) readByte(addr uint16) uint8 {
	va := addr | vic.ciaVaBase
	if (va & 0x7000) == 0x1000 {
		vic.lastByte = vic.board.CharRomRead(va & 0x0fff)
		return vic.lastByte
	}
	vic.lastByte = vic.board.RamRead(va)
	return vic.lastByte
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
		vic.borderOnSample[0] = vic.borderOn
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
		copy(vic.foreMaskBuf, _displayXDiv8)

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
			vic.drawBackground()
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
			vic.drawBackground()
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
		// Refresh and matrix access, increment sprMCBase by 2 if y expansion flip flop is set
		if vic.drawThisLine {
			vic.drawBackground()
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
		// Graphics and matrix access, increment sprMCBase by 1 if y expansion flip flop is set
		// and check if sprite DMA can be turned off
		if vic.drawThisLine {
			vic.drawBackground()
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
						vic.borderOn = false
					} else if !vic.borderOnUL {
						vic.borderOn = false
					}
				} else if !vic.borderOnUL {
					vic.borderOn = false
				}
			}
		}
		// Second sample of border state
		vic.borderOnSample[1] = vic.borderOn
		if vic.drawThisLine {
			if vic.borderOnUL {
				vic.drawBackground()
			} else {
				vic.drawBackground()
				vic.drawGraphics()
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
						vic.borderOn = false
					} else if !vic.borderOnUL {
						vic.borderOn = false
					}
				} else {
					if !vic.borderOnUL {
						vic.borderOn = false
					}
				}
			}
		}

		// Third sample of border state
		vic.borderOnSample[2] = vic.borderOn
		if vic.drawThisLine {
			if vic.borderOnUL {
				vic.drawBackground()
			} else {
				vic.drawGraphics()
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
		vic.charDataLast = vic.charData

	case 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54:
		if vic.drawThisLine {
			if vic.borderOnUL {
				vic.drawBackground()
			} else {
				vic.drawGraphics()
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
		vic.charDataLast = vic.charData

	case 55:
		// Last graphics access, turn off matrix access, turn on sprite DMA if Y coordinate is
		// right and sprite is enabled, handle sprite y expansion, set BA for sprite 0
		if vic.drawThisLine {
			if vic.borderOnUL {
				vic.drawBackground()
			} else {
				vic.drawGraphics()
			}
			vic.sampleBorder()
		}
		vic.graphicsAccess()
		if vic.isBadLine {
			vic.displayOn = true
		}
		// Invert y expansion flip flop if bit in MYE is set
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
			vic.borderOn = true
		}
		// Fourth sample of border state
		vic.borderOnSample[3] = vic.borderOn
		if vic.drawThisLine {
			if vic.borderOnUL {
				vic.drawBackground()
			} else {
				vic.drawGraphics()
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
			vic.borderOn = true
		}
		// Fifth sample of border state
		vic.borderOnSample[4] = vic.borderOn
		// Sample sprDisplayOn and sprData for sprite drawing
		if vic.sprDraw = vic.sprDisplayOn; vic.sprDraw != 0 {
			copy(vic.sprDrawData, vic.sprData)
		}
		// Turn off sprite display if DMA is off
		for idx, mask := 0, uint8(1); idx < 8; idx, mask = idx+1, mask<<1 {
			if (vic.sprDisplayOn&mask) != 0 && (vic.sprDmaOn&mask) == 0 {
				vic.sprDisplayOn &= ^mask
			}
		}
		if vic.drawThisLine {
			vic.drawBackground()
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
			vic.drawBackground()
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
			vic.drawBackground()
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
			const BorderS = 43
			const BorderOffset = BorderS * 8
			vic.drawBackground()
			vic.sampleBorder()
			displayPtr := vic.lineStart
			if vic.sprDraw != 0 {
				vic.drawSprites()
			}
			// Draw border
			if vic.borderOnSample[0] {
				for idx := 0; idx < 4; idx++ {
					copy(vic.displayBuffer[displayPtr+(idx*8):], _colorMultiplier[vic.borderColorSample[idx]])
				}
			}
			if vic.borderOnSample[1] {
				//32 = 4*8
				copy(vic.displayBuffer[displayPtr+(32):], _colorMultiplier[vic.borderColorSample[4]])
			}
			if vic.borderOnSample[2] {
				for idx := 5; idx < BorderS; idx++ {
					copy(vic.displayBuffer[displayPtr+(idx*8):], _colorMultiplier[vic.borderColorSample[idx]])
				}
			}
			if vic.borderOnSample[3] {
				copy(vic.displayBuffer[displayPtr+(BorderOffset):], _colorMultiplier[vic.borderColorSample[BorderS]])
			}
			if vic.borderOnSample[4] {
				for idx := 44; idx < DisplayXDiv8; idx++ {
					copy(vic.displayBuffer[displayPtr+(idx*8):], _colorMultiplier[vic.borderColorSample[idx]])
				}
			}
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
		vic.gfxData = vic.readByte(addr)
		vic.charData = vic.matrixLine[vic.mlIndex]
		vic.colorData = vic.colorLine[vic.mlIndex]
		vic.mlIndex++
		vic.videoCounter++
	} else {
		if (vic.cr1 & 0x40) != 0 {
			vic.gfxData = vic.readByte(0x39ff)
		} else {
			vic.gfxData = vic.readByte(0x3fff)
		}
		vic.colorData = 0
		vic.charData = 0
	}
}

func (vic *MOS6569) drawBackground() {
	var c uint8
	switch vic.displayIdx {
	case 0, 1, 3: // Standard text, Multicolor text, Multicolor bitmap
		c = vic.b0cColor
	case 2: // Standard bitmap
		c = vic.colors[vic.charDataLast]
	case 4: // ECM text
		if (vic.charDataLast & 0x80) != 0 {
			if (vic.charDataLast & 0x40) != 0 {
				c = vic.b3cColor
			} else {
				c = vic.b2cColor
			}
		} else {
			if (vic.charDataLast & 0x40) != 0 {
				c = vic.b1cColor
			} else {
				c = vic.b0cColor
			}
		}
	default:
		c = vic.colors[0]
	}
	copy(vic.displayBuffer[vic.lineOffset:], _colorMultiplier[c])
}

func (vic *MOS6569) drawGraphics() {
	offset := vic.lineOffset + int(vic.xScroll)
	switch vic.displayIdx {
	case 0: // Standard text
		vic.drawGraphicStandard(offset, vic.b0cColor, vic.colors[vic.colorData])
	case 1: // Multicolor text
		if (vic.colorData & 8) != 0 {
			vic.drawGraphicMulticolor(offset, vic.b0cColor, vic.b1cColor, vic.b2cColor, vic.colors[vic.colorData&7])
		} else {
			vic.drawGraphicStandard(offset, vic.b0cColor, vic.colors[vic.colorData])
		}
	case 2: // Standard bitmap
		vic.drawGraphicStandard(offset, vic.colors[vic.charData], vic.colors[vic.charData>>4])
	case 3: // Multicolor bitmap
		vic.drawGraphicMulticolor(offset, vic.b0cColor, vic.colors[vic.charData>>4], vic.colors[vic.charData], vic.colors[vic.colorData])
	case 4: // ECM text
		if (vic.charData & 0x80) != 0 {
			if (vic.charData & 0x40) != 0 {
				vic.drawGraphicStandard(offset, vic.b3cColor, vic.colors[vic.colorData])
			} else {
				vic.drawGraphicStandard(offset, vic.b2cColor, vic.colors[vic.colorData])
			}
		} else if (vic.charData & 0x40) != 0 {
			vic.drawGraphicStandard(offset, vic.b1cColor, vic.colors[vic.colorData])
		} else {
			vic.drawGraphicStandard(offset, vic.b0cColor, vic.colors[vic.colorData])
		}
	case 5: //Invalid multicolor text
		if (vic.colorData & 8) != 0 {
			vic.drawGraphicsInvalidMulticolor(offset, vic.colors[0])
		} else {
			vic.drawGraphicsInvalidStandard(offset, vic.colors[0])
		}
	case 6: //Invalid standard bitmap
		vic.drawGraphicsInvalidStandard(offset, vic.colors[0])
	case 7: // Invalid multicolor bitmap
		vic.drawGraphicsInvalidMulticolor(offset, vic.colors[0])
	}
}

func (vic *MOS6569) drawGraphicsInvalidStandard(offset int, a uint8) {
	copy(vic.displayBuffer[offset:], _colorMultiplier[a])
	vic.foreMaskBuf[vic.foreMaskOffset+0] |= vic.gfxData >> vic.xScroll
	vic.foreMaskBuf[vic.foreMaskOffset+1] |= vic.gfxData << (7 - vic.xScroll)
}

func (vic *MOS6569) drawGraphicsInvalidMulticolor(offset int, a uint8) {
	copy(vic.displayBuffer[offset:], _colorMultiplier[a])
	vic.foreMaskBuf[vic.foreMaskOffset+0] |= ((vic.gfxData & 0xaa) | (vic.gfxData&0xaa)>>1) >> vic.xScroll
	vic.foreMaskBuf[vic.foreMaskOffset+1] |= ((vic.gfxData & 0xaa) | (vic.gfxData&0xaa)>>1) << (8 - vic.xScroll)
}

func (vic *MOS6569) drawGraphicStandard(offset int, a uint8, b uint8) {
	vic.foreMaskBuf[vic.foreMaskOffset+0] |= vic.gfxData >> vic.xScroll
	vic.foreMaskBuf[vic.foreMaskOffset+1] |= vic.gfxData << (7 - vic.xScroll)
	colorBuffer := [4]uint8{a, b, 0, 0}
	data := vic.gfxData
	vic.displayBuffer[offset+7] = colorBuffer[data&1]
	data >>= 1
	vic.displayBuffer[offset+6] = colorBuffer[data&1]
	data >>= 1
	vic.displayBuffer[offset+5] = colorBuffer[data&1]
	data >>= 1
	vic.displayBuffer[offset+4] = colorBuffer[data&1]
	data >>= 1
	vic.displayBuffer[offset+3] = colorBuffer[data&1]
	data >>= 1
	vic.displayBuffer[offset+2] = colorBuffer[data&1]
	data >>= 1
	vic.displayBuffer[offset+1] = colorBuffer[data&1]
	data >>= 1
	vic.displayBuffer[offset] = colorBuffer[data]
}

func (vic *MOS6569) drawGraphicMulticolor(offset int, a uint8, b uint8, c uint8, d uint8) {
	vic.foreMaskBuf[vic.foreMaskOffset+0] |= ((vic.gfxData & 0xaa) | (vic.gfxData&0xaa)>>1) >> vic.xScroll
	vic.foreMaskBuf[vic.foreMaskOffset+1] |= ((vic.gfxData & 0xaa) | (vic.gfxData&0xaa)>>1) << (8 - vic.xScroll)
	colorBuffer := [4]uint8{a, b, c, d}
	data := vic.gfxData
	vic.displayBuffer[offset+7] = colorBuffer[data&3]
	vic.displayBuffer[offset+6] = colorBuffer[data&3]
	data >>= 2
	vic.displayBuffer[offset+5] = colorBuffer[data&3]
	vic.displayBuffer[offset+4] = colorBuffer[data&3]
	data >>= 2
	vic.displayBuffer[offset+3] = colorBuffer[data&3]
	vic.displayBuffer[offset+2] = colorBuffer[data&3]
	data >>= 2
	vic.displayBuffer[offset+1] = colorBuffer[data]
	vic.displayBuffer[offset] = colorBuffer[data]
}

func (vic *MOS6569) drawSprites() {
	//sBit := uint8(1) // bit mask
	sprColl := uint8(0)
	gfxColl := uint8(0)
	copy(vic.sprCollBuf, _displayXFillMax)
	for sNum, sBit := uint8(0), uint8(1); sNum < 8; sNum, sBit = sNum+1, sBit<<1 {
		if vic.sprDraw&sBit != 0 {
			expanded := vic.mxe&sBit != 0
			multiColor := vic.mmc&sBit != 0
			if expanded {
				if multiColor {
					vic.drawSpriteExpandedMulticolor(sNum, sBit, &gfxColl, &sprColl)
				} else {
					vic.drawSpriteExpandedStandard(sNum, sBit, &gfxColl, &sprColl)
				}
			} else {
				if multiColor {
					vic.drawSpriteUnexpandedMulticolor(sNum, sBit, &gfxColl, &sprColl)
				} else {
					vic.drawSpriteUnexpandedStandard(sNum, sBit, &gfxColl, &sprColl)
				}
			}
		}
	}
	// sprite-sprite collisions
	if vic.sprClx != 0 {
		vic.sprClx |= sprColl
	} else {
		vic.sprClx |= sprColl
		vic.irqFlag |= 0x04
		if vic.irqMask&0x04 != 0 {
			vic.irqFlag |= 0x80
			vic.board.VICTriggerIRQ()
		}
	}
	// sprite-background collisions
	if vic.sprClxBgr != 0 {
		vic.sprClxBgr |= gfxColl
	} else {
		vic.sprClxBgr |= gfxColl
		vic.irqFlag |= 0x02
		if vic.irqMask&0x02 != 0 {
			vic.irqFlag |= 0x80
			vic.board.VICTriggerIRQ()
		}
	}
}

func (vic *MOS6569) drawSpriteExpandedMulticolor(sNum uint8, sBit uint8, gfxColl *uint8, sprColl *uint8) {
	q := int(vic.mXx[sNum]) + 8
	displayPtr := vic.lineStart + q
	color := vic.mXcColor[sNum]
	m := q / 8
	s := q & 7
	foreMask := (((uint32(vic.foreMaskBuf[m]) << 24) | (uint32(vic.foreMaskBuf[m+1]) << 16) | (uint32(vic.foreMaskBuf[m+2]) << 8) | (uint32(vic.foreMaskBuf[m+3]))) << s) | (uint32(vic.foreMaskBuf[m+4]) >> (8 - s))
	foreMaskR := (((uint32(vic.foreMaskBuf[m+4]) << 24) | (uint32(vic.foreMaskBuf[m+5]) << 16) | (uint32(vic.foreMaskBuf[m+6]) << 8) | (uint32(vic.foreMaskBuf[m+7]))) << s) | (uint32(vic.foreMaskBuf[m+8]) >> (8 - s))
	dd := vic.sprDrawData[sNum]
	sData := (uint32(dd[0]) << 24) | (uint32(dd[1]) << 16) | (uint32(dd[2]) << 8)
	// Expand sprite data
	sDataL := uint32(_multiExpTable[sData>>24&0xff])<<16 | uint32(_multiExpTable[sData>>16&0xff])
	sDataR := uint32(_multiExpTable[sData>>8&0xff]) << 16
	// Convert sprite chunky pixels to bitPlanes
	plane0L := (sDataL & 0x55555555) | (sDataL&0x55555555)<<1
	plane1L := (sDataL & 0xaaaaaaaa) | (sDataL&0xaaaaaaaa)>>1
	plane0R := (sDataR & 0x55555555) | (sDataR&0x55555555)<<1
	plane1R := (sDataR & 0xaaaaaaaa) | (sDataR&0xaaaaaaaa)>>1
	// Collision with graphics?
	if (foreMask&(plane0L|plane1L)) != 0 || (foreMaskR&(plane0R|plane1R)) != 0 {
		*gfxColl |= sBit
		if vic.mdp&sBit != 0 {
			// Mask sprite if in background
			plane0L &= ^foreMask
			plane1L &= ^foreMask
			plane0R &= ^foreMaskR
			plane1R &= ^foreMaskR
		}
	}
	// Paint sprite
	idx := 0
	for ; idx < 32; idx, plane0L, plane1L = idx+1, plane0L<<1, plane1L<<1 {
		selectedColor := uint8(0)
		if plane1L&0x80000000 != 0 {
			if plane0L&0x80000000 != 0 {
				selectedColor = vic.mm1Color
			} else {
				selectedColor = color
			}
		} else {
			if plane0L&0x80000000 != 0 {
				selectedColor = vic.mm0Color
			} else {
				continue
			}
		}
		if collIdx := q + idx; collIdx < DisplayXFillMax {
			if vic.sprCollBuf[collIdx] != 0 {
				*sprColl |= vic.sprCollBuf[collIdx] | sBit
			} else {
				vic.displayBuffer[displayPtr+idx] = selectedColor
				vic.sprCollBuf[collIdx] = sBit
			}
		}
	}
	for ; idx < 48; idx, plane0R, plane1R = idx+1, plane0R<<1, plane1R<<1 {
		selectedColor := uint8(0)
		if plane1R&0x80000000 != 0 {
			if plane0R&0x80000000 != 0 {
				selectedColor = vic.mm1Color
			} else {
				selectedColor = color
			}
		} else {
			if plane0R&0x80000000 != 0 {
				selectedColor = vic.mm0Color
			} else {
				continue
			}
		}
		if collIdx := q + idx; collIdx < DisplayXFillMax {
			if vic.sprCollBuf[collIdx] != 0 {
				*sprColl |= vic.sprCollBuf[collIdx] | sBit
			} else {
				vic.displayBuffer[displayPtr+idx] = selectedColor
				vic.sprCollBuf[collIdx] = sBit
			}
		}
	}
}

func (vic *MOS6569) drawSpriteExpandedStandard(sNum uint8, sBit uint8, gfxColl *uint8, sprColl *uint8) {
	q := int(vic.mXx[sNum]) + 8
	displayPtr := vic.lineStart + q
	color := vic.mXcColor[sNum]
	m := q / 8
	s := q & 7
	foreMask := (((uint32(vic.foreMaskBuf[m]) << 24) | (uint32(vic.foreMaskBuf[m+1]) << 16) | (uint32(vic.foreMaskBuf[m+2]) << 8) | (uint32(vic.foreMaskBuf[m+3]))) << s) | (uint32(vic.foreMaskBuf[m+4]) >> (8 - s))
	dd := vic.sprDrawData[sNum]
	sData := (uint32(dd[0]) << 24) | (uint32(dd[1]) << 16) | (uint32(dd[2]) << 8)
	// Check graphics collision
	if (foreMask & sData) != 0 {
		*gfxColl |= sBit
		if vic.mdp&sBit != 0 {
			// Mask sprite if in background
			sData &= ^foreMask
		}
	}
	// Paint sprite
	for idx := 0; idx < 24; idx, sData = idx+1, sData<<1 {
		if (sData & 0x80000000) != 0 {
			if collIdx := q + idx; collIdx < DisplayXFillMax {
				if (vic.sprCollBuf[collIdx]) != 0 {
					// Collision with sprite?
					*sprColl |= vic.sprCollBuf[collIdx] | sBit
				} else {
					// Draw pixel if no collision
					vic.displayBuffer[displayPtr+idx] = color
					vic.sprCollBuf[collIdx] = sBit
				}
			}
		}
	}
}

func (vic *MOS6569) drawSpriteUnexpandedMulticolor(sNum uint8, sBit uint8, gfxColl *uint8, sprColl *uint8) {
	q := int(vic.mXx[sNum]) + 8
	displayPtr := vic.lineStart + q
	color := vic.mXcColor[sNum]
	m := q / 8
	s := q & 7
	foreMask := (((uint32(vic.foreMaskBuf[m+0]) << 24) | (uint32(vic.foreMaskBuf[m+1]) << 16) | (uint32(vic.foreMaskBuf[m+2]) << 8) | (uint32(vic.foreMaskBuf[m+3]))) << s) | (uint32(vic.foreMaskBuf[m+4]) >> (8 - s))
	dd := vic.sprDrawData[sNum]
	sData := (uint32(dd[0]) << 24) | (uint32(dd[1]) << 16) | (uint32(dd[2]) << 8)
	// Convert sprite chunky pixels to bitPlanes
	plane0 := (sData & 0x55555555) | (sData&0x55555555)<<1
	plane1 := (sData & 0xaaaaaaaa) | (sData&0xaaaaaaaa)>>1
	// Check graphics collision
	if (foreMask & (plane0 | plane1)) != 0 {
		*gfxColl |= sBit
		if vic.mdp&sBit != 0 {
			// Mask sprite if in background
			plane0 &= ^foreMask
			plane1 &= ^foreMask
		}
	}
	// Paint sprite
	for idx := 0; idx < 24; idx, plane0, plane1 = idx+1, plane0<<1, plane1<<1 {
		var selectedColor uint8
		if (plane1 & 0x80000000) != 0 {
			if (plane0 & 0x80000000) != 0 {
				selectedColor = vic.mm1Color
			} else {
				selectedColor = color
			}
		} else {
			if (plane0 & 0x80000000) != 0 {
				selectedColor = vic.mm0Color
			} else {
				continue
			}
		}
		// Check graphics collision
		if collIdx := q + idx; collIdx < DisplayXFillMax {
			if (vic.sprCollBuf[collIdx]) != 0 {
				// Collision with sprite
				*sprColl |= vic.sprCollBuf[collIdx] | sBit
			} else {
				// Draw pixel if no collision
				vic.displayBuffer[displayPtr+idx] = selectedColor
				vic.sprCollBuf[collIdx] = sBit
			}
		}
	}
}

func (vic *MOS6569) drawSpriteUnexpandedStandard(sNum uint8, sBit uint8, gfxColl *uint8, sprColl *uint8) {
	q := int(vic.mXx[sNum]) + 8
	displayPtr := vic.lineStart + q
	color := vic.mXcColor[sNum]
	m := q / 8
	s := q & 7
	foreMask := (((uint32(vic.foreMaskBuf[m]) << 24) | (uint32(vic.foreMaskBuf[m+1]) << 16) | (uint32(vic.foreMaskBuf[m+2]) << 8) | (uint32(vic.foreMaskBuf[m+3]))) << s) | (uint32(vic.foreMaskBuf[m+4]) >> (8 - s))
	dd := vic.sprDrawData[sNum]
	sData := (uint32(dd[0]) << 24) | (uint32(dd[1]) << 16) | (uint32(dd[2]) << 8)
	// Check graphics collision
	if (foreMask & sData) != 0 {
		*gfxColl |= sBit
		if vic.mdp&sBit != 0 {
			// Mask sprite if in background
			sData &= ^foreMask
		}
	}
	// Paint sprite
	for idx := 0; idx < 24; idx, sData = idx+1, sData<<1 {
		if (sData & 0x80000000) != 0 {
			if collIdx := q + idx; collIdx < DisplayXFillMax {
				if (vic.sprCollBuf[collIdx]) != 0 {
					// Collision with sprite?
					*sprColl |= vic.sprCollBuf[collIdx] | sBit
				} else {
					// Draw pixel if no collision
					vic.displayBuffer[displayPtr+idx] = color
					vic.sprCollBuf[collIdx] = sBit
				}
			}
		}
	}
}
