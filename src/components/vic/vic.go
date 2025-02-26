package mos6569

import (
	"github.com/markel1974/c64emu/src/bits"
	"github.com/markel1974/c64emu/src/config"
	"log"
)

// https://www.cebix.net/VIC-Article.txt
// https://www.oxyron.de/html/registers_vic2.html

// irqRasterBit represents the interrupt request bit for raster events.
// irqSpriteToGraphicBit represents the interrupt request bit for sprite-to-graphic collisions.
// irqSpriteToSpriteBit represents the interrupt request bit for sprite-to-sprite collisions.
// irqLightPenBit represents the interrupt request bit for light pen events.
// irqMasterBit represents the master interrupt request bit.
const (
	irqRasterBit          = uint8(0x01)
	irqSpriteToGraphicBit = uint8(0x02)
	irqSpriteToSpriteBit  = uint8(0x04)
	irqLightPenBit        = uint8(0x08)
	irqMasterBit          = uint8(0x80)
)

// irqUnsetMasterBit is the bitwise negation of irqMasterBit, used to clear the master IRQ bit in irqLatch.
const (
	irqUnsetMasterBit = ^irqMasterBit
)

//https://dustlayer.com/c64-architecture

// cycleData represents a node in a cyclic linked list to manage cycles and associated operations.
type cycleData struct {
	fn          func(vic *VIC)
	next        *cycleData
	cycle       uint8
	cycleBorder uint8
}

// VIC represents the core component for emulating a video interface chip, managing graphics rendering and sprites.
// It maintains configuration, collision detection, registers, and graphical memory base addresses.
// This structure includes VIC-specific registers, raster positions, display modes, and control logic for video operations.
// The structure integrates functionality for sprite handling, display timing, bad line conditions, and IRQ management.
// Graphics systems and border handling are also supported using their specific components within the VIC structure.
type VIC struct {
	id              string
	cfg             *config.Config
	collisions      *Collisions
	sprites         *Sprites
	graphics        *Graphics
	borders         *Borders
	lineStart       int
	drawLine        bool
	vBlankNextCycle bool
	curr            *cycleData

	socket           ISocket
	banks            IBanks
	mXx              []uint16 // VIC registers [m0x - m1x - m2x - m3x - m4x - m5x - m6x - m7x]
	mXy              []uint8  // VIC registers [m0y - m1y - m2y - m3y - m4y - m5y - m6y - m7y]
	mXc              []uint8  // VIC registers [m0c - m1c - m2c - m3c - m4c - m5c - m6c - m7c]
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
	vaBase           uint8    // vaBase
	ciaVaBase        uint16   // CIA VA14/15 video base
	matrixBase       uint16   // Video matrix base
	charBase         uint16   // Character generator base
	bitmapBase       uint16   // Bitmap base
	xScroll          uint16   // X scroll value
	yScroll          uint16   // Y scroll value
	irqLatch         uint8    //
	irqMask          uint8    //
	irqRaster        uint16   // Interrupt raster line
	sprExpY          uint8    // 8 sprite y expansion FlipFlops
	sprBgrClx        uint8    // Sprite to background collision
	sprSprClx        uint8    // Sprite to sprite collision
	rasterX          uint16   // Current raster x position
	rasterY          uint16   // Current raster line
	dyTop            uint16   // Comparison values for borders logic
	dyBottom         uint16   // Comparison values for borders logic
	displayMode      int      // Index of current display mode
	lpTriggered      bool     // LightPen was triggered in this frame
	badLineEnabler   bool     // Bad Lines enabled for this frame
	badLineCondition bool     // Current line is bad line
	baLow            bool     // BA Line
	aecLow           bool     // AEC Line
	aecLowNextCycle  uint64   //
	lastByte         uint8    // Last byte read by VIC
	refreshCounter   uint8    //
	den              bool     //
	bmm              bool     //
	ecm              bool     //
	columnSel        bool     //
}

// NewVIC creates and returns a pointer to a newly initialized VIC instance with default values and given id.
func NewVIC(id string) *VIC {
	vic := &VIC{
		id:               id,
		banks:            nil,
		mXx:              make([]uint16, SpriteNumber),
		mXy:              make([]uint8, SpriteNumber),
		mXc:              make([]uint8, SpriteNumber),
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
		matrixBase:       0,
		charBase:         0,
		bitmapBase:       0,
		vaBase:           0,
		ciaVaBase:        0,
		xScroll:          0,
		yScroll:          0,
		irqRaster:        0,
		irqLatch:         0,
		irqMask:          0,
		sprExpY:          0,
		sprSprClx:        0,
		sprBgrClx:        0,
		rasterX:          0,
		rasterY:          TotalRasters - 1,
		dyTop:            Row24YStart,
		dyBottom:         Row24YStop,
		displayMode:      0,
		lpTriggered:      false,
		badLineCondition: false,
		badLineEnabler:   false,
		baLow:            false,
		aecLowNextCycle:  0,
		aecLow:           false,
		lastByte:         0,
		refreshCounter:   0,
		den:              false,
		bmm:              false,
		ecm:              false,
		columnSel:        false,
	}
	return vic
}

// Setup initializes the VIC component with the provided socket and configuration.
func (vic *VIC) Setup(socket ISocket, cfg *config.Config) {
	vic.cfg = cfg
	db := socket.GetDisplayBuffer()
	vic.socket = socket
	vic.banks = vic.socket.GetBanks()
	vic.collisions = NewCollisions(vic)
	vic.graphics = NewGraphics(vic, vic.collisions, db)
	vic.sprites = NewSprites(vic, vic.collisions, db)
	vic.borders = NewBorder(vic, db)
	vic.vBlankNextCycle = false
	vic.drawLine = false
	vic.cfg.Bind(vic.configChanged)
	vic.graphics.Setup()
	vic.sprites.Setup()
	vic.curr = _pal
}

// Reset reinitializes the VIC instance to its default state by resetting its internal readiness flag.
func (vic *VIC) Reset() {
	//vic.core.ready = false
}

// GetText retrieves the current text data as a byte slice from the VIC's graphics system.
func (vic *VIC) GetText() []byte {
	return vic.graphics.GetText()
}

// GetLastByte returns the last byte stored in the VIC instance.
func (vic *VIC) GetLastByte() uint8 {
	return vic.lastByte
}

// configChanged handles updates to the VIC configuration and applies necessary changes to its state.
func (vic *VIC) configChanged() {
	//vic.skipFrames = vic.cfg.SkipFrames()
}

// Emulate executes one cycle of the VIC, processing the current function and updating the raster position.
func (vic *VIC) Emulate() {
	vic.TryAcquireAEC()
	vic.curr.fn(vic)
	vic.curr = vic.curr.next
	vic.UpdateRasterX()
}

// GetRasterY returns the current vertical raster position as a uint16.
func (vic *VIC) GetRasterY() uint16 {
	return vic.rasterY
}

// ResetRasterX sets the rasterX field to its default initial value, typically used to reset position or state.
func (vic *VIC) ResetRasterX() {
	vic.rasterX = 0xfffc
}

// UpdateRasterX increments the rasterX property of the VIC object by 8 each time it is called.
func (vic *VIC) UpdateRasterX() {
	vic.rasterX += 8
}

// TryBALowIfBadLine checks if the bad line condition is met and sets the BA line to low if true.
func (vic *VIC) TryBALowIfBadLine() {
	if vic.badLineCondition {
		vic.SetBALow()
	}
}

// GetBALow returns the state of the baLow variable, indicating whether the BA low condition is active.
func (vic *VIC) GetBALow() bool {
	return vic.baLow
}

// SetBALow sets the BA (bus available) signal to low and schedules the AEC signal to be low after 3 cycles if not already set.
func (vic *VIC) SetBALow() {
	if vic.baLow {
		return
	}
	vic.baLow = true
	vic.aecLowNextCycle = vic.socket.Cycle() + 3
	vic.socket.BALow(true)
}

// ClearBALow resets the BA low and AEC low flags in the VIC instance and updates the corresponding socket states.
func (vic *VIC) ClearBALow() {
	if vic.baLow {
		vic.baLow = false
		vic.socket.BALow(false)
	}
	if vic.aecLow {
		vic.aecLow = false
		vic.socket.AECLow(false)
	}
}

// GetAECLow retrieves the current state of the AEC low flag for the VIC instance. Returns true if enabled, false otherwise.
func (vic *VIC) GetAECLow() bool {
	return vic.aecLow
}

// TryAcquireAEC attempts to acquire the AEC (Address Enable Control) signal if BA is low and AEC is not already low.
// It ensures the AEC signal is acquired only when the current cycle meets the required condition.
// This method controls the AEC line state by interacting with the VIC's socket mechanism.
func (vic *VIC) TryAcquireAEC() {
	if vic.baLow && !vic.aecLow {
		if vic.socket.Cycle() >= vic.aecLowNextCycle {
			vic.aecLow = true
			vic.socket.AECLow(true)
		}
	}
}

// UpdateSpriteExpY adjusts the sprite's vertical expansion state based on the MYE register using an inversion technique.
func (vic *VIC) UpdateSpriteExpY() {
	// Invert y expansion FlipFlop (if MYE bit is set)
	for idx, mask := 0, uint8(1); idx < SpriteNumber; idx, mask = idx+1, mask<<1 {
		if (vic.mye & mask) != 0 {
			vic.sprExpY ^= mask
		}
	}
}

// badLineUpdate updates the bad line condition based on the current raster position, DEN bit, and YSCROLL value.
// The bad line condition occurs when specific raster and scroll conditions are met, enabling certain VIC behavior.
func (vic *VIC) badLineUpdate() {
	// Bad Line Condition is given at any arbitrary clock cycle, if at the
	// negative edge of ø0 at the beginning of the cycle RASTER >= $30 and RASTER <= $f7
	// and the lower three bits of RASTER are equal to YSCROLL
	// and if the DEN bit has been set for at least one cycle somewhere in raster line $30
	// So clearing the DEN bit will normally prevent Bad Lines

	if (vic.rasterY >= FirstDmaLine) && (vic.rasterY <= LastDmaLine) {
		if vic.rasterY == FirstDmaLine && vic.den {
			//If YSCROLL=0, a Bad Line Condition occurs in raster line $30 as soon as the DEN bit
			vic.badLineEnabler = true
			if vic.yScroll == 0 {
				vic.badLineCondition = true
				return
			}
		}
		if vic.badLineEnabler {
			vic.badLineCondition = vic.yScroll == (vic.rasterY & 7)
		}
	} else {
		vic.badLineEnabler = false
		vic.badLineCondition = false
	}
}

// ChangedVA updates the VIC's virtual address base and triggers the memory pointer update process.
func (vic *VIC) ChangedVA(newVA uint8) {
	vic.ciaVaBase = uint16(newVA) << 14
	vic.memoryPointerUpdate()
}

// LightPenTrigger triggers the light pen interrupt and updates the light pen coordinates if not already triggered.
func (vic *VIC) LightPenTrigger() {
	if !vic.lpTriggered {
		vic.lpTriggered = true
		vic.lpx = uint8(vic.rasterX >> 1)
		vic.lpy = uint8(vic.rasterY)
		vic.irqEmit(irqLightPenBit)
	}
}

// ResetRasterY resets the VIC's raster Y position and refresh counter, and handles IRQ emission for raster line 0.
func (vic *VIC) ResetRasterY() {
	vic.rasterY = 0
	vic.refreshCounter = 0xff
	vic.lpTriggered = false
	if vic.irqRaster == 0 {
		vic.irqEmit(irqRasterBit)
	}
}

// IncrementRasterY increments the rasterY field, triggers an IRQ if it matches irqRaster, and updates bad lines.
func (vic *VIC) IncrementRasterY() {
	vic.rasterY++
	if vic.rasterY == vic.irqRaster {
		vic.irqEmit(irqRasterBit)
	}
	vic.badLineUpdate()
}

// AccessIdle triggers a read operation on address 0x3fff to access the idle state of the VIC component.
func (vic *VIC) AccessIdle() {
	_ = vic.ReadByte(0x3fff)
}

// AccessRefresh performs a memory read using the refresh counter and decrements the counter afterward.
func (vic *VIC) AccessRefresh() {
	_ = vic.ReadByte(0x3f00 | uint16(vic.refreshCounter))
	vic.refreshCounter--
}

// ReadByte reads a byte from the given address after applying the VIC's address translation logic.
func (vic *VIC) ReadByte(addr uint16) uint8 {
	va := addr | vic.ciaVaBase
	if (va & 0x7000) == 0x1000 {
		vic.lastByte = vic.banks.ReadCharRom(va & 0x0fff)
		return vic.lastByte
	}
	vic.lastByte = vic.banks.ReadDirect(va)
	return vic.lastByte
}

// CollisionApply processes sprite-to-sprite and sprite-to-background collisions and emits appropriate IRQ signals.
func (vic *VIC) CollisionApply(sprites uint8, graphics uint8) {
	if vic.sprSprClx != 0 {
		vic.sprSprClx |= sprites
	} else {
		vic.sprSprClx |= sprites
		vic.irqEmit(irqSpriteToSpriteBit)
	}
	if vic.sprBgrClx != 0 {
		vic.sprBgrClx |= graphics
	} else {
		vic.sprBgrClx |= graphics
		vic.irqEmit(irqSpriteToGraphicBit)
	}
}

// ReadRegister reads a register at the given address and returns the corresponding 8-bit value.
func (vic *VIC) ReadRegister(addr uint16) uint8 {
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
	case 0x19: // IRQ latch
		return vic.irqLatch | 0x70
	case 0x1a: // IRQ mask
		return vic.irqMask | 0xf0
	case 0x1b: // Sprite data priority
		return vic.mdp
	case 0x1c: // Sprite multicolor
		return vic.mmc
	case 0x1d: // Sprite X expansion
		return vic.mxe
	case 0x1e: // Sprite-sprite collision
		ret := vic.sprSprClx
		vic.sprSprClx = 0 // Read and clear
		return ret
	case 0x1f: // Sprite-background collision
		ret := vic.sprBgrClx
		vic.sprBgrClx = 0 // Read and clear
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

// WriteRegister writes data to a register at the specified address, handling various control and memory settings.
func (vic *VIC) WriteRegister(addr2 uint16, data uint8) {
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
			if (data & bits.Uint8s[i]) != 0 {
				vic.mXx[i] |= 0x100
			} else {
				vic.mXx[i] &= 0xff
			}
		}
	case 0x11: // Control register 1
		vic.cr1 = data
		vic.yScroll = uint16(vic.cr1) & 7
		if rowSel := (vic.cr1 & 0x8) != 0; rowSel {
			vic.dyTop = Row25YStart
			vic.dyBottom = Row25YStop
		} else {
			vic.dyTop = Row24YStart
			vic.dyBottom = Row24YStop
		}
		vic.den = (vic.cr1 & 0x10) != 0
		vic.bmm = (vic.cr1 & 0x20) != 0
		vic.ecm = (vic.cr1 & 0x40) != 0
		//rst8 := (vic.cr1 & 0x80) != 0
		vic.displayMode = ((int(vic.cr1) & 0x60) | (int(vic.cr2) & 0x10)) >> 4 //cr1 bit 5-6 (BMM|ECM)| cr2 bit 4 (MCM)
		irqRaster := (vic.irqRaster & 0xff) | ((uint16(vic.cr1) & 0x80) << 1)
		vic.rasterUpdate(irqRaster) //can emit irq
		vic.badLineUpdate()
	case 0x12: // Raster counter
		irqRaster := (vic.irqRaster & 0xff00) | uint16(data)
		vic.rasterUpdate(irqRaster) //can emit irq
	case 0x13: // Light pen X
		vic.lpx = data
	case 0x14: // Light pen Y
		vic.lpy = data
	case 0x15: // Sprite enable
		vic.me = data
	case 0x16: // Control register 2
		vic.cr2 = data
		vic.xScroll = uint16(vic.cr2) & 7
		vic.columnSel = (vic.cr2 & 0x8) != 0
		vic.displayMode = ((int(vic.cr1) & 0x60) | (int(vic.cr2) & 0x10)) >> 4 //cr1 bit 5-6 (BMM|ECM)| cr2 bit 4 (MCM)
	case 0x17: // Sprite Y expansion
		vic.mye = data
		vic.sprExpY |= ^data
	case 0x18: // Memory pointers
		vic.vaBase = data
		vic.memoryPointerUpdate()
	case 0x19: // IRQ Latch
		// TODO VERIFICA IMPLEMENTAZIONE
		vic.irqLatch &= ^((data & 0xf) | irqMasterBit)
		vic.irqVerify() //can emit irq
		//old
		//vic.irqLatch &= ^(data & 0xf)
		//if (vic.irqLatch & vic.irqMask) != 0 {
		//	vic.irqLatch |= irqMasterBit // Set master bit if allowed interrupt still pending
		//} else {
		//	vic.socket.IRQClear()
		//}
	case 0x1a: // IRQ mask
		vic.irqMask = data & 0xf
		vic.irqVerify() //can emit irq
	case 0x1b: // Sprite data priority
		vic.mdp = data
	case 0x1c: // Sprite multicolor
		vic.mmc = data
	case 0x1d: // Sprite X expansion
		vic.mxe = data
	case 0x1e: // Sprite-sprite collision
		vic.sprSprClx = data
	case 0x1f: // Sprite-background collision
		vic.sprBgrClx = data
	case 0x20:
		vic.ec = data
	case 0x21:
		vic.b0c = data
	case 0x22:
		vic.b1c = data
	case 0x23:
		vic.b2c = data
	case 0x24:
		vic.b3c = data
	case 0x25:
		vic.mm0 = data
	case 0x26:
		vic.mm1 = data
	case 0x27:
		vic.mXc[0] = data
	case 0x28:
		vic.mXc[1] = data
	case 0x29:
		vic.mXc[2] = data
	case 0x2a:
		vic.mXc[3] = data
	case 0x2b:
		vic.mXc[4] = data
	case 0x2c:
		vic.mXc[5] = data
	case 0x2d:
		vic.mXc[6] = data
	case 0x2e:
		vic.mXc[7] = data
	case 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f: //unconnected
	default:
		log.Printf("WriteRegister: unknown reg 0x%x", reg)
	}
}

// memoryPointerUpdate updates the memory pointers for various VIC-II display components based on the current vaBase value.
// It recalculates matrixBase, charBase, and bitmapBase by applying bitwise operations on the vaBase property.
func (vic *VIC) memoryPointerUpdate() {
	vic.matrixBase = (uint16(vic.vaBase) & 0xf0) << 6
	vic.charBase = (uint16(vic.vaBase) & 0x0e) << 10
	vic.bitmapBase = (uint16(vic.vaBase) & 0x08) << 10
}

// rasterUpdate updates the VIC raster interrupt value and triggers an interrupt if the raster line matches the new value.
func (vic *VIC) rasterUpdate(irqRaster uint16) {
	if irqRaster != vic.irqRaster {
		if vic.rasterY == irqRaster {
			vic.irqEmit(irqRasterBit)
		}
		vic.irqRaster = irqRaster
	}
}

// irqEmit sets the given IRQ bit in irqLatch and triggers IRQ if it matches the irqMask.
func (vic *VIC) irqEmit(irq uint8) {
	vic.irqLatch |= irq
	if (vic.irqMask & irq) != 0 {
		vic.irqLatch |= irqMasterBit
		vic.socket.IRQTrigger()
	}
}

// irqVerify checks the IRQ latch and mask, sets or clears the master IRQ bit, and triggers or clears the IRQ signal.
func (vic *VIC) irqVerify() {
	if (vic.irqLatch & vic.irqMask) != 0 {
		vic.irqLatch |= irqMasterBit
		vic.socket.IRQTrigger() // Trigger interrupt if pending (now allowed)
	} else {
		vic.irqLatch &= irqUnsetMasterBit
		vic.socket.IRQClear()
	}
}
