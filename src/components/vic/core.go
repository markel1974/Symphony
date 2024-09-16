package mos6569

import (
	"log"
)

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
	dyTop            uint16   // Comparison values for borders logic
	dyBottom         uint16   // Comparison values for borders logic
	colors           []uint8  // Indices of the 16 colors (16 times mirrored to avoid "& 0x0f")
	displayIdx       int      // Index of current display mode
	lpTriggered      bool     // LightPen was triggered in this frame
	badLinesEnabled  bool     // Bad Lines enabled for this frame
	badLineCondition bool     // Current line is bad line
	baLow            bool     // BA Line
	aecLow           bool     // AEC Line
	aecLowNextCycle  uint64   //
	lastByte         uint8    // Last byte read by VIC
	refreshCounter   uint8    //
	//ready            bool   // VIC Initialization Complete
}

func NewCore(socket ISocket) *Core {
	colors := make([]uint8, 256)
	for i := range colors {
		colors[i] = (uint8)(i & 0x0f)
	}
	defaultColor := colors[0] //black
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
		ecColor:          defaultColor,
		b0cColor:         defaultColor,
		b1cColor:         defaultColor,
		b2cColor:         defaultColor,
		b3cColor:         defaultColor,
		mm0Color:         defaultColor,
		mm1Color:         defaultColor,
		colors:           colors,
		rasterX:          0,
		rasterY:          TotalRasters - 1,
		dyTop:            Row24YStart,
		dyBottom:         Row24YStop,
		displayIdx:       0,
		lpTriggered:      false,
		badLineCondition: false,
		badLinesEnabled:  false,
		baLow:            false,
		aecLowNextCycle:  0,
		aecLow:           false,
		lastByte:         0,
		refreshCounter:   0,
		//ready:          false,
	}
	for i := range c.mXcColor {
		c.mXcColor[i] = defaultColor
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

/*
func (vic *Core) ModeColumn38() bool {
	return (vic.cr2 & 8) == 0
}

func (vic *Core) ModeColumn40() bool {
	return (vic.cr2 & 8) != 0
}
*/

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
	vic.aecLowNextCycle = vic.socket.Cycle() + 3
	vic.socket.BALow(true)
}

func (vic *Core) ClearBALow() {
	if vic.baLow {
		vic.baLow = false
		vic.socket.BALow(false)
	}
	if vic.aecLow {
		vic.aecLow = false
		vic.socket.AECLow(false)
	}
}

func (vic *Core) GetAECLow() bool {
	return vic.aecLow
}

func (vic *Core) TryAcquireAEC() {
	if vic.baLow && !vic.aecLow {
		if vic.socket.Cycle() >= vic.aecLowNextCycle {
			vic.aecLow = true
			vic.socket.AECLow(true)
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
	//vic.WriteRegister(0x18, vic.vaBase)
	vic.matrixBase = (uint16(vic.vaBase) & 0xf0) << 6
	vic.charBase = (uint16(vic.vaBase) & 0x0e) << 10
	vic.bitmapBase = (uint16(vic.vaBase) & 0x08) << 10
}

func (vic *Core) LightPenTrigger() {
	if !vic.lpTriggered {
		vic.lpTriggered = true
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
	reg := addr & 0x3f
	switch reg {
	case 0x00:
		return uint8(vic.mXx[0])
	case 0x01:
		return vic.mXy[0]
	case 0x02:
		return uint8(vic.mXx[1])
	case 0x03:
		return vic.mXy[1]
	case 0x04:
		return uint8(vic.mXx[2])
	case 0x05:
		return vic.mXy[2]
	case 0x06:
		return uint8(vic.mXx[3])
	case 0x07:
		return vic.mXy[3]
	case 0x08:
		return uint8(vic.mXx[4])
	case 0x09:
		return vic.mXy[4]
	case 0x0a:
		return uint8(vic.mXx[5])
	case 0x0b:
		return vic.mXy[5]
	case 0x0c:
		return uint8(vic.mXx[6])
	case 0x0d:
		return vic.mXy[6]
	case 0x0e:
		return uint8(vic.mXx[7])
	case 0x0f:
		return vic.mXy[7]
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
	case 0x27:
		return vic.mXc[0] | 0xf0
	case 0x28:
		return vic.mXc[1] | 0xf0
	case 0x29:
		return vic.mXc[2] | 0xf0
	case 0x2a:
		return vic.mXc[3] | 0xf0
	case 0x2b:
		return vic.mXc[4] | 0xf0
	case 0x2c:
		return vic.mXc[5] | 0xf0
	case 0x2d:
		return vic.mXc[6] | 0xf0
	case 0x2e:
		return vic.mXc[7] | 0xf0
	case 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f: //unconnected
		return 0xff
	default:
		log.Printf("ReadRegister: unknown reg 0x%x", reg)
		return 0xff
	}
}

func (vic *Core) WriteRegister(addr2 uint16, data uint8) {
	reg := addr2 & 0x3f
	switch reg {
	case 0x00:
		vic.mXx[0] = (vic.mXx[0] & 0xff00) | uint16(data)
	case 0x01:
		vic.mXy[0] = data
	case 0x02:
		vic.mXx[1] = (vic.mXx[1] & 0xff00) | uint16(data)
	case 0x03:
		vic.mXy[1] = data
	case 0x04:
		vic.mXx[2] = (vic.mXx[2] & 0xff00) | uint16(data)
	case 0x05:
		vic.mXy[2] = data
	case 0x06:
		vic.mXx[3] = (vic.mXx[3] & 0xff00) | uint16(data)
	case 0x07:
		vic.mXy[3] = data
	case 0x08:
		vic.mXx[4] = (vic.mXx[4] & 0xff00) | uint16(data)
	case 0x09:
		vic.mXy[4] = data
	case 0x0a:
		vic.mXx[5] = (vic.mXx[5] & 0xff00) | uint16(data)
	case 0x0b:
		vic.mXy[5] = data
	case 0x0c:
		vic.mXx[6] = (vic.mXx[6] & 0xff00) | uint16(data)
	case 0x0d:
		vic.mXy[6] = data
	case 0x0e:
		vic.mXx[7] = (vic.mXx[7] & 0xff00) | uint16(data)
	case 0x0f:
		vic.mXy[7] = data
	case 0x10: //MSBs of X coordinates
		vic.mx8 = data
		for i := 0; i < SpriteNumber; i++ {
			if (data & _spriteBit[i]) != 0 {
				vic.mXx[i] |= 0x100
			} else {
				vic.mXx[i] &= 0xff
			}
		}

		/*
			for i, j := 0, uint8(1); i < SpriteNumber; i, j = i+1, j<<1 {
				//fmt.Println(j)
				if (vic.mx8 & j) != 0 {
					vic.mXx[i] |= 0x100
				} else {
					vic.mXx[i] &= 0xff
				}
			}
		*/
	case 0x11: // Control register 1
		vic.cr1 = data
		vic.yScroll = uint16(data) & 7
		irqRaster := (vic.irqRaster & 0xff) | ((uint16(data) & 0x80) << 1)
		if vic.irqRaster != irqRaster && vic.rasterY == irqRaster {
			vic.rasterIrq()
		}
		vic.irqRaster = irqRaster
		if (data & 8) != 0 {
			vic.dyTop = Row25YStart
			vic.dyBottom = Row25YStop
		} else {
			vic.dyTop = Row24YStart
			vic.dyBottom = Row24YStop
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
	case 0x13: // Light pen X
		vic.lpx = data
	case 0x14: // Light pen Y
		vic.lpy = data
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
	case 0x1e: // Sprite-sprite collision
		//log.Printf("Write Sprite-sprite collision %d", data)
		vic.sprClx = data
	case 0x1f: // Sprite-background collision
		//log.Printf("Write Sprite-background collision %d", data)
		vic.sprClxBgr = data
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
	case 0x27:
		vic.mXc[0] = data
		vic.mXcColor[0] = vic.colors[data]
	case 0x28:
		vic.mXc[1] = data
		vic.mXcColor[1] = vic.colors[data]
	case 0x29:
		vic.mXc[2] = data
		vic.mXcColor[2] = vic.colors[data]
	case 0x2a:
		vic.mXc[3] = data
		vic.mXcColor[3] = vic.colors[data]
	case 0x2b:
		vic.mXc[4] = data
		vic.mXcColor[4] = vic.colors[data]
	case 0x2c:
		vic.mXc[5] = data
		vic.mXcColor[5] = vic.colors[data]
	case 0x2d:
		vic.mXc[6] = data
		vic.mXcColor[6] = vic.colors[data]
	case 0x2e:
		vic.mXc[7] = data
		vic.mXcColor[7] = vic.colors[data]
	case 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f: //unconnected
	default:
		log.Printf("WriteRegister: unknown reg 0x%x", reg)
	}
}
