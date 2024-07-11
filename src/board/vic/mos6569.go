package vic

import (
	"github.com/markel1974/c64emu/src/board/iboard"
	"github.com/markel1974/c64emu/src/preferences"
)

//https://dustlayer.com/c64-architecture

type MOS6569 struct {
	cycle             int
	lastByte          uint8 // Last byte read by VIC
	board             iboard.IBoard
	prefs             *preferences.Prefs
	displayBuffer     []uint8
	chunkyLineStart   int       // Offset from current line in bitmap buffer
	chunkyPtr         int       // Offset from chunky bitmap buffer
	foreMaskPtr       int       // offset from fmBuf
	mx                [8]uint16 // VIC registers
	my                [8]uint8  // VIC registers
	mx8               uint8     // VIC registers
	ctrl1             uint8
	ctrl2             uint8
	lpx               uint8
	lpy               uint8
	me                uint8
	mxe               uint8
	mye               uint8
	mdp               uint8
	mmc               uint8
	vaBase            uint8
	irqFlag           uint8
	irqMask           uint8
	clxSpr            uint8
	clxBgr            uint8
	ec                uint8
	b0c               uint8
	b1c               uint8
	b2c               uint8
	b3c               uint8
	mm0               uint8
	mm1               uint8
	ecColor           uint8
	b0cColor          uint8 // Indices for exterior/background colors
	b1cColor          uint8 // Indices for exterior/background colors
	b2cColor          uint8 // Indices for exterior/background colors
	b3cColor          uint8 // Indices for exterior/background colors
	mm0Color          uint8 // Indices for MOB multi colors
	mm1Color          uint8 // Indices for MOB multi colors
	sc                [8]uint8
	colors            [256]uint8              // Indices of the 16 colors (16 times mirrored to avoid "& 0x0f")
	sprColor          [8]uint8                // Indices for MOB colors
	matrixLine        [40]uint8               // Buffer for video line, read in Bad Lines
	colorLine         [40]uint8               // Buffer for color line, read in Bad Lines
	rasterY           uint16                  // Current raster line
	irqRaster         uint16                  // Interrupt raster line
	dyStart           uint16                  // Comparison values for border logic
	dyStop            uint16                  // Comparison values for border logic
	rc                uint16                  // Row counter
	vc                uint16                  // Video counter
	vcBase            uint16                  // Video counter base
	xScroll           uint16                  // X scroll value
	yScroll           uint16                  // Y scroll value
	ciaVaBase         uint16                  // CIA VA14/15 video base
	mc                [8]uint16               // Sprite data counters
	sprCollBuf        []uint8                 // REAL DisplayX = 0x180 | Buffer for sprite-sprite collisions and priorities
	fmBuf             []uint8                 // DisplayX / 8 | Foreground mask for sprite-graphics collisions and priorities
	borderColorSample [DisplayXFill + 1]uint8 // DisplayX / 8 | Samples of border color at each "displayed" cycle
	borderOnSample    [5]bool                 // Samples of border state at different cycles (1, 17, 18, 56, 57) //PROTECTION AGAIN BUFFER OVERFLOW !!! borderColorSample FROM OUT OF BUFFER!!!!
	sprPtr            [8]uint16               // Sprite data pointers
	mcBase            [8]uint16               // Sprite data counter bases
	sprData           [][]uint8               // Sprite data read
	sprDrawData       [][]uint8               // Sprite data for drawing
	displayIdx        int                     // Index of current display mode
	displayOn         bool                    // true: Display state, false: Idle state
	borderOn          bool                    // Flag: Upper/lower border on (Main border FlipFlop)
	badLinesEnabled   bool                    // Flag: Bad Lines enabled for this frame
	lpTriggered       bool                    // Flag: LightPen was triggered in this frame
	matrixBase        uint16                  // Video matrix base
	charBase          uint16                  // Character generator base
	bitmapBase        uint16                  // Bitmap base
	isBadLine         bool                    // Flag: Current line is bad line
	drawThisLine      bool                    // Flag: This line is drawn on the _screen
	udBorderOn        bool                    // Flag: Upper/lower border on
	refreshCounter    uint8                   // Refresh counter
	sprExpY           uint8                   // 8 sprite y expansion flip flops
	sprDmaOn          uint8                   // 8 flags: Sprite DMA active
	sprDisplayOn      uint8                   // 8 flags: Sprite display active
	sprDraw           uint8                   // 8 flags: Draw sprite in this line
	rasterX           uint16                  // Current raster x position
	mlIndex           int                     // Index in matrix/colorLine[]
	gfxData           uint8
	charData          uint8
	colorData         uint8
	lastCharData      uint8
	firstBaCycle      uint64
	vBlanking         bool // Flag: VBlank in next cycle
	baLow             bool
	ready             bool // VIC Initialization Complete
	//ecColorLong       uint32 // ecColor expanded to 32 bits
}

var _colorMultiplier [][]uint8
var _displayXDiv8 []uint8
var _displayXFillMax []uint8

//var _displayXPlusOne []uint8

func init() {
	_colorMultiplier = make([][]uint8, 0xff)
	for x := uint8(0); x < 0xff; x++ {
		_colorMultiplier[x] = []uint8{x, x, x, x, x, x, x, x}
	}
	_displayXDiv8 = make([]uint8, DisplayXDiv8)
	for x := 0; x < DisplayXDiv8; x++ {
		_displayXDiv8[x] = uint8(0)
	}
	_displayXFillMax = make([]uint8, DisplayXFillMax)
	for x := 0; x < DisplayXFillMax; x++ {
		_displayXFillMax[x] = 0
	}
	//_displayXPlusOne = make([]uint8, DisplayXFill+1)
	//for x := 0; x < DisplayXFill+1; x++ {
	//	_displayXPlusOne[x] = 0
	//}
}

func NewMOS6569() *MOS6569 {
	vic := &MOS6569{}
	vic.sprData = make([][]uint8, 8)
	for i := range vic.sprData {
		vic.sprData[i] = make([]uint8, 4)
	}
	vic.sprDrawData = make([][]uint8, 8)
	for i := range vic.sprDrawData {
		vic.sprDrawData[i] = make([]uint8, 4)
	}
	vic.ready = false
	vic.displayBuffer = make([]uint8, DisplaySize) //(color_t*)malloc(sizeof(color_t) * DisplaySize)
	vic.chunkyPtr = 0                              //&vic.displayBuffer[0]
	vic.isBadLine = false
	vic.sprExpY = 0
	for i := 0; i < 8; i++ {
		vic.mcBase[i] = 0
	}
	// Set pointers
	vic.matrixBase = 0
	vic.charBase = 0
	vic.bitmapBase = 0
	// Initialize VIC registers
	vic.mx8 = 0
	vic.ctrl1 = 0
	vic.ctrl2 = 0
	vic.lpx = 0
	vic.lpy = 0
	vic.me = 0
	vic.mxe = 0
	vic.mye = 0
	vic.mdp = 0
	vic.mmc = 0
	vic.vaBase = 0
	vic.irqFlag = 0
	vic.irqMask = 0
	vic.clxSpr = 0
	vic.clxBgr = 0
	vic.ciaVaBase = 0
	vic.ec = 0
	vic.b0c = 0
	vic.b1c = 0
	vic.b2c = 0
	vic.b3c = 0
	vic.mm0 = 0
	vic.mm1 = 0
	for i := 0; i < 8; i++ {
		vic.mx[i] = 0
		vic.my[i] = 0
		vic.sc[i] = 0
	}
	// Initialize other variables
	vic.rasterY = TotalRasters - 1
	vic.rc = 7
	vic.irqRaster = 0
	vic.vc = 0
	vic.vcBase = 0
	vic.xScroll = 0
	vic.yScroll = 0
	vic.dyStart = Row24YStart
	vic.dyStop = Row24YStop
	vic.mlIndex = 0
	vic.cycle = 1
	vic.displayIdx = 0
	vic.displayOn = false
	vic.borderOn = false
	vic.udBorderOn = false
	vic.vBlanking = false
	vic.lpTriggered = false
	vic.drawThisLine = false
	vic.sprDmaOn = 0
	vic.sprDisplayOn = 0
	for i := 0; i < 8; i++ {
		vic.mc[i] = 63
		vic.sprPtr[i] = 0
	}
	vic.sprCollBuf = make([]uint8, DisplayXFillMax)
	copy(vic.sprCollBuf, _displayXFillMax)
	vic.fmBuf = make([]uint8, DisplayXFill+1)
	copy(vic.fmBuf, _displayXDiv8)
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
		vic.sprColor[i] = vic.colors[0]
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
		if (vic.me&mask) != 0 && (vic.rasterY&0xff) == uint16(vic.my[i]) {
			vic.sprDmaOn |= mask
			vic.mcBase[i] = 0
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
		vic.sprData[num][byteNum] = vic.readByte((vic.mc[num] & 0x3f) | vic.sprPtr[num])
		vic.mc[num]++
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
	vic.chunkyPtr += 8
	vic.foreMaskPtr++
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
		return uint8(vic.mx[addr>>1])
	case 0x01, 0x03, 0x05, 0x07, 0x09, 0x0b, 0x0d, 0x0f:
		return vic.my[addr>>1]
	// Sprite X position MSB
	case 0x10:
		return vic.mx8
	// Control register 1
	case 0x11:
		return uint8((uint16(vic.ctrl1) & 0x7f) | ((vic.rasterY & 0x100) >> 1))
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
		return vic.ctrl2 | 0xc0
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
		ret := vic.clxSpr
		vic.clxSpr = 0 // Read and clear
		return ret
	// Sprite-background collision
	case 0x1f:
		ret := vic.clxBgr
		vic.clxBgr = 0 // Read and clear
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
		return vic.sc[addr-0x27] | 0xf0
	default:
		return 0xff
	}
}

func (vic *MOS6569) WriteRegister(addr uint16, data uint8) {
	addr = addr & 0x3f
	switch addr {
	case 0x00, 0x02, 0x04, 0x06, 0x08, 0x0a, 0x0c, 0x0e:
		target := addr >> 1
		vic.mx[target] = (vic.mx[target] & 0xff00) | uint16(data)
	case 0x10:
		vic.mx8 = data
		for i, j := 0, uint8(1); i < 8; i, j = i+1, j<<1 {
			if (vic.mx8 & j) != 0 {
				vic.mx[i] |= 0x100
			} else {
				vic.mx[i] &= 0xff
			}
		}
	case 0x01, 0x03, 0x05, 0x07, 0x09, 0x0b, 0x0d, 0x0f:
		vic.my[addr>>1] = data
	case 0x11: // Control register 1
		vic.ctrl1 = data
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
		vic.displayIdx = ((int(vic.ctrl1) & 0x60) | (int(vic.ctrl2) & 0x10)) >> 4
	case 0x12: // Raster counter
		newIRQRaster := (vic.irqRaster & 0xff00) | uint16(data)
		if vic.irqRaster != newIRQRaster && vic.rasterY == newIRQRaster {
			vic.rasterIrq()
		}
		vic.irqRaster = newIRQRaster
	case 0x15: // Sprite enable
		vic.me = data
	case 0x16: // Control register 2
		vic.ctrl2 = data
		vic.xScroll = uint16(data) & 7
		vic.displayIdx = ((int(vic.ctrl1) & 0x60) | (int(vic.ctrl2) & 0x10)) >> 4
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
		vic.sc[target] = data
		vic.sprColor[target] = vic.colors[data]
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
				vic.badLinesEnabled = (vic.ctrl1 & 0x10) != 0
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
			vic.vcBase = 0
			vic.rasterY = 0
			vic.refreshCounter = 0xff
			vic.vBlanking = false
			vic.lpTriggered = false
			vic.chunkyLineStart = 0
			if vic.irqRaster == 0 {
				// Trigger raster IRQ if IRQ in line 0
				vic.rasterIrq()
			}
		}
		// Our output goes here
		vic.chunkyPtr = vic.chunkyLineStart
		// Clear foreground mask
		copy(vic.fmBuf, _displayXDiv8)

		vic.foreMaskPtr = 0
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
		// Refresh, vc -> vcBase, turn on matrix access and reset RC if Bad Line
		if vic.drawThisLine {
			vic.drawBackground()
			vic.sampleBorder()
		}
		vic.refreshAccess()
		// Turn on display and matrix access and reset RC if Bad Line
		if vic.isBadLine {
			vic.rc = 0
			vic.displayOn = true
			vic.setBALow()
		}
		vic.vc = vic.vcBase

	case 15:
		// Refresh and matrix access, increment mcBase by 2 if y expansion flip flop is set
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
				vic.mcBase[idx] += 2
			}
		}
		vic.mlIndex = 0
		vic.matrixAccess()

	case 16:
		// Graphics and matrix access, increment mcBase by 1 if y expansion flip flop is set
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
				vic.mcBase[idx]++
			}
			if (vic.mcBase[idx] & 0x3f) == 0x3f {
				vic.sprDmaOn &= ^mask
			}
		}
		vic.matrixAccess()

	case 17:
		// Graphics and matrix access, turn off border in 40 column mode, display window starts here
		if (vic.ctrl2 & 8) != 0 {
			if vic.rasterY == vic.dyStop {
				vic.udBorderOn = true
			} else {
				if (vic.ctrl1 & 0x10) != 0 {
					if vic.rasterY == vic.dyStart {
						vic.udBorderOn = false
						vic.borderOn = false
					} else if !vic.udBorderOn {
						vic.borderOn = false
					}
				} else if !vic.udBorderOn {
					vic.borderOn = false
				}
			}
		}
		// Second sample of border state
		vic.borderOnSample[1] = vic.borderOn
		if vic.drawThisLine {
			if vic.udBorderOn {
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
		if (vic.ctrl2 & 8) == 0 {
			if vic.rasterY == vic.dyStop {
				vic.udBorderOn = true
			} else {
				if (vic.ctrl1 & 0x10) != 0 {
					if vic.rasterY == vic.dyStart {
						vic.udBorderOn = false
						vic.borderOn = false
					} else if !vic.udBorderOn {
						vic.borderOn = false
					}
				} else {
					if !vic.udBorderOn {
						vic.borderOn = false
					}
				}
			}
		}

		// Third sample of border state
		vic.borderOnSample[2] = vic.borderOn
		if vic.drawThisLine {
			if vic.udBorderOn {
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
		vic.lastCharData = vic.charData

	case 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54:
		if vic.drawThisLine {
			if vic.udBorderOn {
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
		vic.lastCharData = vic.charData

	case 55:
		// Last graphics access, turn off matrix access, turn on sprite DMA if Y coordinate is
		// right and sprite is enabled, handle sprite y expansion, set BA for sprite 0
		if vic.drawThisLine {
			if vic.udBorderOn {
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
		if (vic.ctrl2 & 8) == 0 {
			vic.borderOn = true
		}
		// Fourth sample of border state
		vic.borderOnSample[3] = vic.borderOn
		if vic.drawThisLine {
			if vic.udBorderOn {
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
		if (vic.ctrl2 & 8) != 0 {
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
		// Fetch sprite pointer 0, mcBase->mc, turn on sprite display if necessary,
		// turn off display if RC=7, read data of sprite 0
		if vic.drawThisLine {
			vic.drawBackground()
			vic.sampleBorder()
		}
		rasterY := vic.rasterY & 0xff
		for idx, mask := 0, uint8(1); idx < 8; idx, mask = idx+1, mask<<1 {
			vic.mc[idx] = vic.mcBase[idx]
			if (vic.sprDmaOn&mask) != 0 && (rasterY == uint16(vic.my[idx])) {
				vic.sprDisplayOn |= mask
			}
		}
		vic.fetchSpriteDataPtr(0)
		vic.fetchSpriteData(0, 0)
		if vic.rc == 7 {
			vic.vcBase = vic.vc
			vic.displayOn = false
		}
		if vic.isBadLine || vic.displayOn {
			vic.displayOn = true
			vic.rc = (vic.rc + 1) & 7
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
			displayPtr := vic.chunkyLineStart
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
			vic.chunkyLineStart += DisplayX
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
			vic.udBorderOn = true
		} else if (vic.ctrl1&0x10) != 0 && vic.rasterY == vic.dyStart {
			vic.udBorderOn = false
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
			addr := (vic.vc & 0x03ff) | vic.matrixBase
			vic.matrixLine[vic.mlIndex] = vic.readByte(addr)
			vic.colorLine[vic.mlIndex] = vic.board.ColorRead(addr & 0x03ff)
		}
	}
}

func (vic *MOS6569) graphicsAccess() {
	if vic.displayOn {
		var addr uint16
		if (vic.ctrl1 & 0x20) != 0 {
			addr = ((vic.vc & 0x03ff) << 3) | vic.bitmapBase | vic.rc // Bitmap
		} else {
			addr = (uint16(vic.matrixLine[vic.mlIndex]) << 3) | vic.charBase | vic.rc // Text
		}
		if (vic.ctrl1 & 0x40) != 0 {
			addr &= 0xf9ff // ECM
		}
		vic.gfxData = vic.readByte(addr)
		vic.charData = vic.matrixLine[vic.mlIndex]
		vic.colorData = vic.colorLine[vic.mlIndex]
		vic.mlIndex++
		vic.vc++
	} else {
		if (vic.ctrl1 & 0x40) != 0 {
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
		c = vic.colors[vic.lastCharData]
	case 4: // ECM text
		if (vic.lastCharData & 0x80) != 0 {
			if (vic.lastCharData & 0x40) != 0 {
				c = vic.b3cColor
			} else {
				c = vic.b2cColor
			}
		} else {
			if (vic.lastCharData & 0x40) != 0 {
				c = vic.b1cColor
			} else {
				c = vic.b0cColor
			}
		}
	default:
		c = vic.colors[0]
	}
	copy(vic.displayBuffer[vic.chunkyPtr:], _colorMultiplier[c])
}

func (vic *MOS6569) drawGraphics() {
	offset := vic.chunkyPtr + int(vic.xScroll)
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
	vic.fmBuf[vic.foreMaskPtr+0] |= vic.gfxData >> vic.xScroll
	vic.fmBuf[vic.foreMaskPtr+1] |= vic.gfxData << (7 - vic.xScroll)
}

func (vic *MOS6569) drawGraphicsInvalidMulticolor(offset int, a uint8) {
	copy(vic.displayBuffer[offset:], _colorMultiplier[a])
	vic.fmBuf[vic.foreMaskPtr+0] |= ((vic.gfxData & 0xaa) | (vic.gfxData&0xaa)>>1) >> vic.xScroll
	vic.fmBuf[vic.foreMaskPtr+1] |= ((vic.gfxData & 0xaa) | (vic.gfxData&0xaa)>>1) << (8 - vic.xScroll)
}

func (vic *MOS6569) drawGraphicStandard(offset int, a uint8, b uint8) {
	vic.fmBuf[vic.foreMaskPtr+0] |= vic.gfxData >> vic.xScroll
	vic.fmBuf[vic.foreMaskPtr+1] |= vic.gfxData << (7 - vic.xScroll)
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
	vic.fmBuf[vic.foreMaskPtr+0] |= ((vic.gfxData & 0xaa) | (vic.gfxData&0xaa)>>1) >> vic.xScroll
	vic.fmBuf[vic.foreMaskPtr+1] |= ((vic.gfxData & 0xaa) | (vic.gfxData&0xaa)>>1) << (8 - vic.xScroll)
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
	if vic.clxSpr != 0 {
		vic.clxSpr |= sprColl
	} else {
		vic.clxSpr |= sprColl
		vic.irqFlag |= 0x04
		if vic.irqMask&0x04 != 0 {
			vic.irqFlag |= 0x80
			vic.board.VICTriggerIRQ()
		}
	}
	// sprite-background collisions
	if vic.clxBgr != 0 {
		vic.clxBgr |= gfxColl
	} else {
		vic.clxBgr |= gfxColl
		vic.irqFlag |= 0x02
		if vic.irqMask&0x02 != 0 {
			vic.irqFlag |= 0x80
			vic.board.VICTriggerIRQ()
		}
	}
}

func (vic *MOS6569) drawSpriteExpandedMulticolor(sNum uint8, sBit uint8, gfxColl *uint8, sprColl *uint8) {
	q := int(vic.mx[sNum]) + 8
	displayPtr := vic.chunkyLineStart + q
	color := vic.sprColor[sNum]
	m := q / 8
	s := q & 7
	foreMask := (((uint32(vic.fmBuf[m]) << 24) | (uint32(vic.fmBuf[m+1]) << 16) | (uint32(vic.fmBuf[m+2]) << 8) | (uint32(vic.fmBuf[m+3]))) << s) | (uint32(vic.fmBuf[m+4]) >> (8 - s))
	foreMaskR := (((uint32(vic.fmBuf[m+4]) << 24) | (uint32(vic.fmBuf[m+5]) << 16) | (uint32(vic.fmBuf[m+6]) << 8) | (uint32(vic.fmBuf[m+7]))) << s) | (uint32(vic.fmBuf[m+8]) >> (8 - s))
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
	q := int(vic.mx[sNum]) + 8
	displayPtr := vic.chunkyLineStart + q
	color := vic.sprColor[sNum]
	m := q / 8
	s := q & 7
	foreMask := (((uint32(vic.fmBuf[m]) << 24) | (uint32(vic.fmBuf[m+1]) << 16) | (uint32(vic.fmBuf[m+2]) << 8) | (uint32(vic.fmBuf[m+3]))) << s) | (uint32(vic.fmBuf[m+4]) >> (8 - s))
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
	q := int(vic.mx[sNum]) + 8
	displayPtr := vic.chunkyLineStart + q
	color := vic.sprColor[sNum]
	m := q / 8
	s := q & 7
	foreMask := (((uint32(vic.fmBuf[m+0]) << 24) | (uint32(vic.fmBuf[m+1]) << 16) | (uint32(vic.fmBuf[m+2]) << 8) | (uint32(vic.fmBuf[m+3]))) << s) | (uint32(vic.fmBuf[m+4]) >> (8 - s))
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
	q := int(vic.mx[sNum]) + 8
	displayPtr := vic.chunkyLineStart + q
	color := vic.sprColor[sNum]
	m := q / 8
	s := q & 7
	foreMask := (((uint32(vic.fmBuf[m]) << 24) | (uint32(vic.fmBuf[m+1]) << 16) | (uint32(vic.fmBuf[m+2]) << 8) | (uint32(vic.fmBuf[m+3]))) << s) | (uint32(vic.fmBuf[m+4]) >> (8 - s))
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
