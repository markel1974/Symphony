package mos6569

// https://www.cebix.net/VIC-Article.txt
// https://www.oxyron.de/html/registers_vic2.html

const (
	flagRasterBit   = 0x01
	flagLightPenBit = 0x08
	flagMasterBit   = 0x80
)

type Core struct {
	socket           ISocket
	banks            IBanks
	mXx              []uint16 // VIC registers [m0x - m1x - m2x - m3x - m4x - m5x - m6x - m7x]
	mXy              []uint8  // VIC registers [m0y - m1y - m2y - m3y - m4y - m5y - m6y - m7y]
	mx8              uint8    // VIC register
	cr1              uint8    // VIC register
	cr2              uint8    // VIC register
	lpx              uint8    // VIC register
	lpy              uint8    // VIC register
	me               uint8    // VIC register
	mxe              uint8    // VIC register
	mye              uint8    // VIC register
	mdp              uint8    // VIC register
	mmc              uint8    // VIC register
	ec               uint8    // VIC register
	b0c              uint8    // VIC register
	b1c              uint8    // VIC register
	b2c              uint8    // VIC register
	b3c              uint8    // VIC register
	mm0              uint8    // VIC register
	mm1              uint8    // VIC register
	mXc              []uint8  // VIC registers [m0c - m1c - m2c - m3c - m4c - m5c - m6c - m7c]
	ecColor          uint8    // Index ec Color Mapping
	b0cColor         uint8    // Index b0c Color Mapping
	b1cColor         uint8    // Index b1c Color Mapping
	b2cColor         uint8    // Index b2c Color Mapping
	b3cColor         uint8    // Index b3c Color Mapping
	mm0Color         uint8    // Index mm0 Color Mapping
	mm1Color         uint8    // Index mm1 Color Mapping
	mXcColor         []uint8  // Color Index Mapping for m0c - m1c - m2c - m3c - m4c - m5c - m6c - m7c
	vaBase           uint8    // vaBase
	ciaVaBase        uint16   // CIA VA14/15 video base
	matrixBase       uint16   // Video matrix base
	charBase         uint16   // Character generator base
	bitmapBase       uint16   // Bitmap base
	xScroll          uint16   // X scroll value
	yScroll          uint16   // Y scroll value
	irqFlag          uint8    //
	irqMask          uint8    //
	irqRaster        uint16   // Interrupt raster line
	sprExpY          uint8    // 8 sprite y expansion FlipFlops
	sprClxBgr        uint8    // Sprite to background collision
	sprClx           uint8    // Sprite to sprite collision
	rasterX          uint16   // Current raster x position
	rasterY          uint16   // Current raster line
	dyStart          uint16   // Comparison values for borders logic
	dyStop           uint16   // Comparison values for borders logic
	colors           []uint8  // Indices of the 16 colors (16 times mirrored to avoid "& 0x0f")
	displayIdx       int      // Index of current display mode
	lpTriggered      bool     // LightPen was triggered in this frame
	badLinesEnabled  bool     // Bad Lines enabled for this frame
	badLineCondition bool     // Current line is bad line
	ready            bool     // VIC Initialization Complete
	baLow            bool     // BA Line
	aecLow           bool     // AEC Line
	baLowFirstCycle  uint64   //
	//aecLowClearNextCycle bool     //
	lastByte       uint8 // Last byte read by VIC
	refreshCounter uint8
}

func NewCore(socket ISocket) *Core {
	colors := make([]uint8, 256)
	for i := range colors {
		colors[i] = (uint8)(i & 0x0f)
	}
	c := &Core{
		socket:           socket,
		banks:            nil,
		mXx:              make([]uint16, SpriteNumber),
		mXy:              make([]uint8, SpriteNumber),
		mx8:              0,
		cr1:              0,
		cr2:              0,
		lpx:              0,
		lpy:              0,
		me:               0,
		mxe:              0,
		mye:              0,
		mdp:              0,
		mmc:              0,
		ec:               0,
		b0c:              0,
		b1c:              0,
		b2c:              0,
		b3c:              0,
		mm0:              0,
		mm1:              0,
		mXc:              make([]uint8, SpriteNumber),
		mXcColor:         make([]uint8, SpriteNumber),
		matrixBase:       0,
		charBase:         0,
		bitmapBase:       0,
		vaBase:           0,
		ciaVaBase:        0,
		xScroll:          0,
		yScroll:          0,
		irqRaster:        0,
		irqFlag:          0,
		irqMask:          0,
		sprExpY:          0,
		sprClx:           0,
		sprClxBgr:        0,
		ecColor:          colors[0], // Preset colors to black
		b0cColor:         colors[0], // Preset colors to black
		b1cColor:         colors[0], // Preset colors to black
		b2cColor:         colors[0], // Preset colors to black
		b3cColor:         colors[0], // Preset colors to black
		mm0Color:         colors[0], // Preset colors to black
		mm1Color:         colors[0], // Preset colors to black
		colors:           colors,
		rasterX:          0,
		rasterY:          TotalRasters - 1,
		dyStart:          Row24YStart,
		dyStop:           Row24YStop,
		displayIdx:       0,
		lpTriggered:      false,
		badLineCondition: false,
		badLinesEnabled:  false,
		baLow:            false,
		baLowFirstCycle:  0,
		aecLow:           false,
		ready:            false,
		lastByte:         0,
		refreshCounter:   0,
	}
	// Preset colors to black
	for i := range c.mXcColor {
		c.mXcColor[i] = c.colors[0]
	}
	return c
}

func (vic *Core) Setup() {
	vic.banks = vic.socket.GetBanks()
}

func (vic *Core) GetRasterY() uint16 {
	return vic.rasterY
}

func (vic *Core) ResetRasterX() {
	vic.rasterX = 0xfffc
}

func (vic *Core) UpdateRasterX() {
	vic.rasterX += 8
}

func (vic *Core) ModeColumn38() bool {
	return (vic.cr2 & 8) == 0
}

func (vic *Core) ModeColumn40() bool {
	return (vic.cr2 & 8) != 0
}

func (vic *Core) TryBALowIfBadLine() {
	if vic.badLineCondition {
		vic.SetBALow()
	}
}

func (vic *Core) GetBALow() bool {
	return vic.baLow
}

func (vic *Core) SetBALow() {
	if vic.baLow {
		return
	}
	vic.baLow = true
	vic.baLowFirstCycle = vic.socket.Cycle()
	vic.socket.BALow(true)
}

func (vic *Core) ClearBALow() {
	if vic.baLow {
		vic.baLow = false
		vic.baLowFirstCycle = 0
		vic.socket.BALow(vic.baLow)
	}
	if vic.aecLow {
		//vic.aecLowClearNextCycle = true
		vic.aecLow = false
		vic.socket.AECLow(vic.aecLow)
	}
}

func (vic *Core) GetAECLow() bool {
	return vic.aecLow
}

func (vic *Core) TryAcquireAEC() {
	if vic.baLow && !vic.aecLow {
		if dist := vic.socket.Cycle() - vic.baLowFirstCycle; dist >= 3 {
			vic.aecLow = true
			vic.socket.AECLow(vic.aecLow)
		}
	}
}

func (vic *Core) FlipFlopMYE() {
	// Invert y expansion FlipFlop (if MYE bit is set)
	for idx, mask := 0, uint8(1); idx < SpriteNumber; idx, mask = idx+1, mask<<1 {
		if (vic.mye & mask) != 0 {
			vic.sprExpY ^= mask
		}
	}
}

func (vic *Core) badLineUpdate() {
	// Bad Line Condition is given at any arbitrary clock cycle, if at the
	// negative edge of ø0 at the beginning of the cycle RASTER >= $30 and RASTER <= $f7
	// and the lower three bits of RASTER are equal to YSCROLL
	// and if the DEN bit was set during an arbitrary cycle of raster line $30.

	//if vic.rasterY == 0x30 {
	//	vic.badLinesEnabled = (vic.cr1 & 0x10) != 0 //denBit
	//}

	vic.badLineCondition = (vic.rasterY >= FirstDmaLine) && (vic.rasterY <= LastDmaLine) && ((vic.rasterY & 7) == vic.yScroll) && vic.badLinesEnabled
}

func (vic *Core) ChangedVA(newVA uint8) {
	vic.ciaVaBase = uint16(newVA) << 14
	// Force update memory pointers
	vic.WriteRegister(0x18, vic.vaBase)
}

func (vic *Core) LightPenTrigger() {
	if !vic.lpTriggered {
		vic.lpTriggered = true
		// Latch current coordinates
		vic.lpx = uint8(vic.rasterX >> 1)
		vic.lpy = uint8(vic.rasterY)
		vic.irqFlag |= flagLightPenBit
		if (vic.irqMask & flagLightPenBit) != 0 {
			vic.irqFlag |= flagMasterBit
			vic.socket.IRQTrigger()
		}
	}
}

func (vic *Core) ResetCounters() {
	vic.refreshCounter = 0xff
	vic.rasterY = 0
	vic.lpTriggered = false
	if vic.irqRaster == 0 {
		vic.rasterIrq()
	}
}

func (vic *Core) IncrementCounters() {
	vic.rasterY++
	if vic.rasterY == vic.irqRaster {
		vic.rasterIrq()
	}
	if vic.rasterY == 0x30 {
		vic.badLinesEnabled = (vic.cr1 & 0x10) != 0 //denBit
	}
	vic.badLineUpdate()
}

func (vic *Core) AccessIdle() {
	_ = vic.ReadByte(0x3fff)
}

func (vic *Core) AccessRefresh() {
	_ = vic.ReadByte(0x3f00 | uint16(vic.refreshCounter))
	vic.refreshCounter--
}

func (vic *Core) ReadByte(addr uint16) uint8 {
	va := addr | vic.ciaVaBase
	if (va & 0x7000) == 0x1000 {
		vic.lastByte = vic.banks.ReadCharRom(va & 0x0fff)
		return vic.lastByte
	}
	vic.lastByte = vic.banks.ReadDirect(va)
	return vic.lastByte
}

func (vic *Core) rasterIrq() {
	vic.irqFlag |= flagRasterBit
	if (vic.irqMask & flagRasterBit) != 0 {
		vic.irqFlag |= flagMasterBit
		vic.socket.IRQTrigger()
	}
}

func (vic *Core) ReadRegister(addr uint16) uint8 {
	addr = addr & 0x3f
	switch addr {
	case 0x00, 0x02, 0x04, 0x06, 0x08, 0x0a, 0x0c, 0x0e:
		return uint8(vic.mXx[addr>>1])
	case 0x01, 0x03, 0x05, 0x07, 0x09, 0x0b, 0x0d, 0x0f:
		return vic.mXy[addr>>1]

	case 0x10: // Sprite X position MSB
		return vic.mx8

	case 0x11: // Control register 1
		return uint8((uint16(vic.cr1) & 0x7f) | ((vic.rasterY & 0x100) >> 1))

	case 0x12: // Raster counter
		return uint8(vic.rasterY)

	case 0x13: // Light pen X
		return vic.lpx

	case 0x14: // Light pen Y
		return vic.lpy

	case 0x15: // Sprite enable
		return vic.me

	case 0x16: // Control register 2
		return vic.cr2 | 0xc0

	case 0x17: // Sprite Y expansion
		return vic.mye

	case 0x18: // Memory pointers
		return vic.vaBase | 0x01

	case 0x19: // IRQ spriteFlags
		if !vic.ready {
			vic.ready = true
			vic.socket.Ready()
		}
		return vic.irqFlag | 0x70

	case 0x1a: // IRQ mask
		return vic.irqMask | 0xf0

	case 0x1b: // Sprite data priority
		return vic.mdp

	case 0x1c: // Sprite multicolor
		return vic.mmc

	case 0x1d: // Sprite X expansion
		return vic.mxe

	case 0x1e: // Sprite-sprite collision
		ret := vic.sprClx
		vic.sprClx = 0 // Read and clear
		return ret

	case 0x1f: // Sprite-background collision
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

func (vic *Core) WriteRegister(addr uint16, data uint8) {
	addr = addr & 0x3f
	switch addr {
	case 0x00, 0x02, 0x04, 0x06, 0x08, 0x0a, 0x0c, 0x0e:
		target := addr >> 1
		vic.mXx[target] = (vic.mXx[target] & 0xff00) | uint16(data)
	case 0x10:
		vic.mx8 = data
		for i, j := 0, uint8(1); i < SpriteNumber; i, j = i+1, j<<1 {
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
		irqRaster := (vic.irqRaster & 0xff) | ((uint16(data) & 0x80) << 1)
		if vic.irqRaster != irqRaster && vic.rasterY == irqRaster {
			vic.rasterIrq()
		}
		vic.irqRaster = irqRaster
		if (data & 8) != 0 {
			vic.dyStart = Row25YStart
			vic.dyStop = Row25YStop
		} else {
			vic.dyStart = Row24YStart
			vic.dyStop = Row24YStop
		}
		if (vic.rasterY == 0x30) && ((vic.cr1 & 0x10) != 0) {
			vic.badLinesEnabled = true
		}
		vic.badLineUpdate()
		vic.displayIdx = ((int(vic.cr1) & 0x60) | (int(vic.cr2) & 0x10)) >> 4
	case 0x12: // Raster counter
		irqRaster := (vic.irqRaster & 0xff00) | uint16(data)
		if vic.irqRaster != irqRaster && vic.rasterY == irqRaster {
			vic.rasterIrq()
		}
		vic.irqRaster = irqRaster
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
	case 0x19: // IRQ spriteFlags
		vic.irqFlag = vic.irqFlag & (^data & 0x0f)
		if (vic.irqFlag & vic.irqMask) != 0 {
			// Set master bit if allowed interrupt still pending
			vic.irqFlag |= flagMasterBit
		} else {
			vic.socket.IRQClear()
		}
	case 0x1a: // IRQ mask
		vic.irqMask = data & 0x0f
		if (vic.irqFlag & vic.irqMask) != 0 {
			// Trigger interrupt if pending (now allowed)
			vic.irqFlag |= flagMasterBit
			vic.socket.IRQTrigger()
		} else {
			vic.irqFlag &= 0x7f
			vic.socket.IRQClear()
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
