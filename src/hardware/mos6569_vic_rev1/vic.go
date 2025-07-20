package mos6569

import (
	"github.com/markel1974/c64emu/src/common/bits"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
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

const (
	ioAndCharRomArea = 0x7000
	charRomOffset    = 0x1000
)

// irqUnsetMasterBit is the bitwise negation of irqMasterBit, used to clear the master IRQ bit in irqLatch.
const (
	irqUnsetMasterBit = ^irqMasterBit
)

const (
	RegisterSize  = 0x3f
	RegisterCount = RegisterSize + 1
)

//https://dustlayer.com/c64-architecture

// cycleData represents a node in a cyclic linked list to manage cycles and associated operations.
type cycleData struct {
	fn          func(vic *VIC)
	next        *cycleData
	cycle       uint8
	cycleBorder uint8
}

// VIC represents a versatile interface controller for managing video output and graphical resources in a system.
// It encapsulates configurations, graphics components, collision detection, and rendering capabilities for the display.
type VIC struct {
	*component.BaseComponent
	cfg                   *config.Config
	collisions            *Collisions
	sprites               *SpriteHandler
	graphics              *Graphics
	borders               *Borders
	readRam               func(addr uint16) uint8
	readColorRam          func(addr uint16) uint8
	readCharRom           func(addr uint16) uint8
	socketCycle           func() uint64
	socketBALow           func(bool)
	socketAECLow          func(bool)
	socketIRQTrigger      func()
	socketIRQClearTrigger func()
	socketLastCycle       func()
	socketVBlank          func()
	curr                  *cycleData
	lineStart             int
	drawLine              bool
	vBlankNextCycle       bool
	mXx                   []uint16 // VIC registers [m0x - m1x - m2x - m3x - m4x - m5x - m6x - m7x]
	mXy                   []uint8  // VIC registers [m0y - m1y - m2y - m3y - m4y - m5y - m6y - m7y]
	mXc                   []uint8  // VIC registers [m0c - m1c - m2c - m3c - m4c - m5c - m6c - m7c]
	mx8                   uint8    // VIC register
	cr1                   uint8    // VIC register
	cr2                   uint8    // VIC register
	lpx                   uint8    // VIC register
	lpy                   uint8    // VIC register
	me                    uint8    // VIC register
	mxe                   uint8    // VIC register
	mye                   uint8    // VIC register
	mdp                   uint8    // VIC register
	mmc                   uint8    // VIC register
	ec                    uint8    // VIC register
	b0c                   uint8    // VIC register
	b1c                   uint8    // VIC register
	b2c                   uint8    // VIC register
	b3c                   uint8    // VIC register
	mm0                   uint8    // VIC register
	mm1                   uint8    // VIC register
	vaBase                uint8    // vaBase
	ciaVaBase             uint16   // CIA VA14/15 video base
	matrixBase            uint16   // Video matrix base
	charBase              uint16   // Character generator base
	bitmapBase            uint16   // Bitmap base
	xScroll               uint16   // X scroll value
	yScroll               uint16   // Y scroll value
	irqLatch              uint8    // irqLatch holds an 8-bit value that latches the IRQ (Interrupt Request) configuration.
	irqMask               uint8    // irqMask represents an 8-bit mask used for interrupt request (IRQ) management.
	irqRaster             uint16   // Interrupt raster line
	sprExpY               uint8    // 8 sprite y expansion FlipFlops
	sprBgrClx             uint8    // Sprite to background collision
	sprSprClx             uint8    // Sprite to sprite collision
	rasterX               uint16   // Current raster x position
	rasterY               uint16   // Current raster line
	dyTop                 uint16   // Comparison values for borders logic
	dyBottom              uint16   // Comparison values for borders logic
	displayMode           int      // Index of current display mode
	lpTriggered           bool     // LightPen was triggered in this frame
	badLineEnabler        bool     // Bad Lines enabled for this frame
	badLineCondition      bool     // Current line is bad line
	baLow                 bool     // BA Line
	aecLow                bool     // AEC Line
	aecLowNextCycle       uint64   // aecLowNextCycle represents the counter for the next cycle in the AEC low-level operation.
	lastByte              uint8    // Last byte read by VIC
	refreshCounter        uint8    // refreshCounter tracks the number of times a refresh operation has been performed.
	den                   bool     // den indicates a boolean value typically used as a flag or condition.
	bmm                   bool     // bmm indicates a boolean value used for specific conditional checks or state representation.
	ecm                   bool     // ecm indicates whether the ECM is active or not.
	columnSel             bool     // columnSel indicates whether column selection mode is enabled.
	label                 string
	reads                 [RegisterCount]func() uint8
	writes                [RegisterCount]func(uint8)
}

// NewVIC creates and initializes a new VIC instance with default configuration and registers it with the parent component.
func NewVIC(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *VIC {
	vic := &VIC{
		BaseComponent:    component.NewBaseComponent(),
		readRam:          nil,
		readColorRam:     nil,
		readCharRom:      nil,
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
		label:            label,
	}
	vic.BaseComponent.Register(factory, parent, Identifier(), vic, references.IdIMos6569(vic, label, instance))
	return vic
}

func (vic *VIC) Setup() error {
	vic.cfg = vic.GetFactory().GetConfig()
	return nil
}

func (vic *VIC) Bind(socket references.IMos6569Socket) error {
	displayBuffer := vic.GetFactory().GetIDisplayBuffer()

	vic.socketIRQTrigger = socket.IRQTrigger
	vic.socketIRQClearTrigger = socket.IRQClearTrigger
	vic.socketCycle = socket.Cycle
	vic.socketLastCycle = socket.LastCycle
	vic.socketBALow = socket.BALow
	vic.socketAECLow = socket.AECLow
	vic.socketVBlank = socket.VBlank
	vic.readRam = socket.ReadRam
	vic.readColorRam = socket.ReadColorRam
	vic.readCharRom = socket.ReadCharRom

	vic.collisions = NewCollisions(vic, vic.GetFactory(), vic.label, 0, vic)
	vic.graphics = NewGraphics(vic, vic.GetFactory(), vic.label, 0, vic, vic.collisions, displayBuffer)
	vic.sprites = NewSprites(vic, vic.GetFactory(), vic.label, 0, vic, vic.collisions, displayBuffer)
	vic.borders = NewBorder(vic, vic.GetFactory(), vic.label, 0, vic, displayBuffer)
	vic.vBlankNextCycle = false
	vic.drawLine = false
	vic.cfg.Bind(vic.configChanged)
	if err := vic.collisions.Setup(); err != nil {
		return err
	}
	if err := vic.graphics.Setup(); err != nil {
		return err
	}
	if err := vic.sprites.Setup(); err != nil {
		return err
	}
	if err := vic.borders.Setup(); err != nil {
		return err
	}
	vic.curr = _pal
	vic.reads = vic.createReadRegister()
	vic.writes = vic.createWriteRegister()

	return nil
}

func (vic *VIC) Connect() error {
	return nil
}

func (vic *VIC) Internal() bool {
	return false
}

// Reset reinitializes the VIC instance to its default state by resetting its internal readiness flag.
func (vic *VIC) Reset() {
	//vic.core.ready = false
}

// GetText retrieves the current text data as a byte slice from the VIC's graphics system.
func (vic *VIC) GetText() []byte {
	return vic.graphics.GetText()
}

// GetVASignal returns the last byte stored in the VIC instance.
func (vic *VIC) GetVASignal() uint8 {
	return vic.lastByte
}

// configChanged handles updates to the VIC configuration and applies necessary changes to its state.
func (vic *VIC) configChanged() {
}

// Emulate executes one cycle of the VIC, processing the current function and updating the raster position.
//
//go:nosplit
func (vic *VIC) Emulate() {
	vic.TryAcquireAEC()
	vic.curr.fn(vic)
	vic.curr = vic.curr.next
	vic.UpdateRasterX()
}

func (vic *VIC) EmulationRequired() bool {
	return true
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
	vic.aecLowNextCycle = vic.socketCycle() + 3
	vic.socketBALow(true)
}

// ClearBALow resets the BA low and AEC low flags in the VIC instance and updates the corresponding socket states.
func (vic *VIC) ClearBALow() {
	if vic.baLow {
		vic.baLow = false
		vic.socketBALow(false)
	}
	if vic.aecLow {
		vic.aecLow = false
		vic.socketAECLow(false)
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
		if vic.socketCycle() >= vic.aecLowNextCycle {
			vic.aecLow = true
			vic.socketAECLow(true)
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
	if (va & ioAndCharRomArea) == charRomOffset {
		vic.lastByte = vic.readCharRom(va)
		return vic.lastByte
	}
	vic.lastByte = vic.readRam(va)
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
	reg := addr & RegisterSize
	return vic.reads[reg]()
}

// WriteRegister writes data to a register at the specified address, handling various control and memory settings.
func (vic *VIC) WriteRegister(addr uint16, data uint8) {
	reg := addr & RegisterSize
	vic.writes[reg](data)
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
		vic.socketIRQTrigger()
	}
}

// irqVerify checks the IRQ latch and mask, sets or clears the master IRQ bit, and triggers or clears the IRQ signal.
func (vic *VIC) irqVerify() {
	if (vic.irqLatch & vic.irqMask) != 0 {
		vic.irqLatch |= irqMasterBit
		vic.socketIRQTrigger() // Trigger interrupt if pending (now allowed)
	} else {
		vic.irqLatch &= irqUnsetMasterBit
		vic.socketIRQClearTrigger()
	}
}

// createReadRegister initializes an array of functions for reading VIC-II registers based on their respective indices.
// Each register is mapped to a specific read function, or defaults to returning 0xff if unconnected.
func (vic *VIC) createReadRegister() [RegisterCount]func() uint8 {
	var reads [RegisterCount]func() uint8
	var unconnected = func() uint8 {
		return 0xff
	}
	for idx := range reads {
		reads[idx] = unconnected
	}
	reads[0x00] = func() uint8 {
		return uint8(vic.mXx[0])
	}
	reads[0x01] = func() uint8 {
		return vic.mXy[0]
	}
	reads[0x02] = func() uint8 {
		return uint8(vic.mXx[1])
	}
	reads[0x03] = func() uint8 {
		return vic.mXy[1]
	}
	reads[0x04] = func() uint8 {
		return uint8(vic.mXx[2])
	}
	reads[0x05] = func() uint8 {
		return vic.mXy[2]
	}
	reads[0x06] = func() uint8 {
		return uint8(vic.mXx[3])
	}
	reads[0x07] = func() uint8 {
		return vic.mXy[3]
	}
	reads[0x08] = func() uint8 {
		return uint8(vic.mXx[4])
	}
	reads[0x09] = func() uint8 {
		return vic.mXy[4]
	}
	reads[0x0a] = func() uint8 {
		return uint8(vic.mXx[5])
	}
	reads[0x0b] = func() uint8 {
		return vic.mXy[5]
	}
	reads[0x0c] = func() uint8 {
		return uint8(vic.mXx[6])
	}
	reads[0x0d] = func() uint8 {
		return vic.mXy[6]
	}
	reads[0x0e] = func() uint8 {
		return uint8(vic.mXx[7])
	}
	reads[0x0f] = func() uint8 {
		return vic.mXy[7]
	}
	reads[0x10] = func() uint8 {
		// Sprite X position MSB
		return vic.mx8
	}
	reads[0x11] = func() uint8 {
		// Control register 1
		return uint8((uint16(vic.cr1) & 0x7f) | ((vic.rasterY & 0x100) >> 1))
	}
	reads[0x12] = func() uint8 {
		// Raster counter
		return uint8(vic.rasterY)
	}
	reads[0x13] = func() uint8 {
		// Light pen X
		return vic.lpx
	}
	reads[0x14] = func() uint8 {
		// Light pen Y
		return vic.lpy
	}
	reads[0x15] = func() uint8 {
		// Sprite enable
		return vic.me
	}
	reads[0x16] = func() uint8 {
		// Control register 2
		return vic.cr2 | 0xc0
	}
	reads[0x17] = func() uint8 {
		// Sprite Y expansion
		return vic.mye
	}
	reads[0x18] = func() uint8 {
		// Memory pointers
		return vic.vaBase | 0x01
	}
	reads[0x19] = func() uint8 {
		// IRQ latch
		return vic.irqLatch | 0x70
	}
	reads[0x1a] = func() uint8 {
		// IRQ mask
		return vic.irqMask | 0xf0
	}
	reads[0x1b] = func() uint8 {
		// Sprite data priority
		return vic.mdp
	}
	reads[0x1c] = func() uint8 {
		// Sprite multicolor
		return vic.mmc
	}
	reads[0x1d] = func() uint8 {
		// Sprite X expansion
		return vic.mxe
	}
	reads[0x1e] = func() uint8 {
		// Sprite-sprite collision
		ret := vic.sprSprClx
		vic.sprSprClx = 0 // Read and clear
		return ret
	}
	reads[0x1f] = func() uint8 {
		// Sprite-background collision
		ret := vic.sprBgrClx
		vic.sprBgrClx = 0 // Read and clear
		return ret
	}
	reads[0x20] = func() uint8 {
		return vic.ec | 0xf0
	}
	reads[0x21] = func() uint8 {
		return vic.b0c | 0xf0
	}
	reads[0x22] = func() uint8 {
		return vic.b1c | 0xf0
	}
	reads[0x23] = func() uint8 {
		return vic.b2c | 0xf0
	}
	reads[0x24] = func() uint8 {
		return vic.b3c | 0xf0
	}
	reads[0x25] = func() uint8 {
		return vic.mm0 | 0xf0
	}
	reads[0x26] = func() uint8 {
		return vic.mm1 | 0xf0
	}
	reads[0x27] = func() uint8 {
		return vic.mXc[0] | 0xf0
	}
	reads[0x28] = func() uint8 {
		return vic.mXc[1] | 0xf0
	}
	reads[0x29] = func() uint8 {
		return vic.mXc[2] | 0xf0
	}
	reads[0x2a] = func() uint8 {
		return vic.mXc[3] | 0xf0
	}
	reads[0x2b] = func() uint8 {
		return vic.mXc[4] | 0xf0
	}
	reads[0x2c] = func() uint8 {
		return vic.mXc[5] | 0xf0
	}
	reads[0x2d] = func() uint8 {
		return vic.mXc[6] | 0xf0
	}
	reads[0x2e] = func() uint8 {
		return vic.mXc[7] | 0xf0
	}
	return reads
}

// createWriteRegister initializes an array of functions for writing data to various VIC-II registers.
func (vic *VIC) createWriteRegister() [RegisterCount]func(uint8) {
	var writes [RegisterCount]func(uint8)
	var unconnected = func(uint8) {
	}
	for idx := range writes {
		writes[idx] = unconnected
	}
	writes[0x00] = func(data uint8) {
		vic.mXx[0] = (vic.mXx[0] & 0xff00) | uint16(data)
	}
	writes[0x01] = func(data uint8) {
		vic.mXy[0] = data
	}
	writes[0x02] = func(data uint8) {
		vic.mXx[1] = (vic.mXx[1] & 0xff00) | uint16(data)
	}
	writes[0x03] = func(data uint8) {
		vic.mXy[1] = data
	}
	writes[0x04] = func(data uint8) {
		vic.mXx[2] = (vic.mXx[2] & 0xff00) | uint16(data)
	}
	writes[0x05] = func(data uint8) {
		vic.mXy[2] = data
	}
	writes[0x06] = func(data uint8) {
		vic.mXx[3] = (vic.mXx[3] & 0xff00) | uint16(data)
	}
	writes[0x07] = func(data uint8) {
		vic.mXy[3] = data
	}
	writes[0x08] = func(data uint8) {
		vic.mXx[4] = (vic.mXx[4] & 0xff00) | uint16(data)
	}
	writes[0x09] = func(data uint8) {
		vic.mXy[4] = data
	}
	writes[0x0a] = func(data uint8) {
		vic.mXx[5] = (vic.mXx[5] & 0xff00) | uint16(data)
	}
	writes[0x0b] = func(data uint8) {
		vic.mXy[5] = data
	}
	writes[0x0c] = func(data uint8) {
		vic.mXx[6] = (vic.mXx[6] & 0xff00) | uint16(data)
	}
	writes[0x0d] = func(data uint8) {
		vic.mXy[6] = data
	}
	writes[0x0e] = func(data uint8) {
		vic.mXx[7] = (vic.mXx[7] & 0xff00) | uint16(data)
	}
	writes[0x0f] = func(data uint8) {
		vic.mXy[7] = data
	}
	writes[0x10] = func(data uint8) { //MSBs of X coordinates
		vic.mx8 = data
		for i := 0; i < SpriteNumber; i++ {
			if (data & bits.Uint8s[i]) != 0 {
				vic.mXx[i] |= 0x100
			} else {
				vic.mXx[i] &= 0xff
			}
		}
	}
	writes[0x11] = func(data uint8) { // Control register 1
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
	}
	writes[0x12] = func(data uint8) { // Raster counter
		irqRaster := (vic.irqRaster & 0xff00) | uint16(data)
		vic.rasterUpdate(irqRaster) //can emit irq
	}
	writes[0x13] = func(data uint8) { // Light pen X
		vic.lpx = data
	}
	writes[0x14] = func(data uint8) { // Light pen Y
		vic.lpy = data
	}
	writes[0x15] = func(data uint8) { // Sprite enable
		vic.me = data
	}
	writes[0x16] = func(data uint8) { // Control register 2
		vic.cr2 = data
		vic.xScroll = uint16(vic.cr2) & 7
		vic.columnSel = (vic.cr2 & 0x8) != 0
		vic.displayMode = ((int(vic.cr1) & 0x60) | (int(vic.cr2) & 0x10)) >> 4 //cr1 bit 5-6 (BMM|ECM)| cr2 bit 4 (MCM)
	}
	writes[0x17] = func(data uint8) { // Sprite Y expansion
		vic.mye = data
		vic.sprExpY |= ^data
	}
	writes[0x18] = func(data uint8) { // Memory pointers
		vic.vaBase = data
		vic.memoryPointerUpdate()
	}
	writes[0x19] = func(data uint8) { // IRQ Latch
		// Verify implementation
		vic.irqLatch &= ^((data & 0xf) | irqMasterBit)
		vic.irqVerify() //can emit irq
		//old
		//vic.irqLatch &= ^(data & 0xf)
		//if (vic.irqLatch & vic.irqMask) != 0 {
		//	vic.irqLatch |= irqMasterBit // Set master bit if allowed interrupt still pending
		//} else {
		//	vic.socket.IRQClearTrigger()
		//}
	}
	writes[0x1a] = func(data uint8) { // IRQ mask
		vic.irqMask = data & 0xf
		vic.irqVerify() //can emit irq
	}
	writes[0x1b] = func(data uint8) { // Sprite data priority
		vic.mdp = data
	}
	writes[0x1c] = func(data uint8) { // Sprite multicolor
		vic.mmc = data
		vic.sprites.ModeUpdate()
	}
	writes[0x1d] = func(data uint8) { // Sprite X expansion
		vic.mxe = data
		vic.sprites.ModeUpdate()
	}
	writes[0x1e] = func(data uint8) { // Sprite-sprite collision
		vic.sprSprClx = data
	}
	writes[0x1f] = func(data uint8) { // Sprite-background collision
		vic.sprBgrClx = data
	}
	writes[0x20] = func(data uint8) {
		vic.ec = data
	}
	writes[0x21] = func(data uint8) {
		vic.b0c = data
	}
	writes[0x22] = func(data uint8) {
		vic.b1c = data
	}
	writes[0x23] = func(data uint8) {
		vic.b2c = data
	}
	writes[0x24] = func(data uint8) {
		vic.b3c = data
	}
	writes[0x25] = func(data uint8) {
		vic.mm0 = data
	}
	writes[0x26] = func(data uint8) {
		vic.mm1 = data
	}
	writes[0x27] = func(data uint8) {
		vic.mXc[0] = data
	}
	writes[0x28] = func(data uint8) {
		vic.mXc[1] = data
	}
	writes[0x29] = func(data uint8) {
		vic.mXc[2] = data
	}
	writes[0x2a] = func(data uint8) {
		vic.mXc[3] = data
	}
	writes[0x2b] = func(data uint8) {
		vic.mXc[4] = data
	}
	writes[0x2c] = func(data uint8) {
		vic.mXc[5] = data
	}
	writes[0x2d] = func(data uint8) {
		vic.mXc[6] = data
	}
	writes[0x2e] = func(data uint8) {
		vic.mXc[7] = data
	}
	return writes
}
