package mos6569

import (
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

// VIC represents a versatile interface controller for managing video output and graphical resources in a system.
// It encapsulates configurations, graphics components, collision detection, and rendering capabilities for the display.
type VIC struct {
	*component.BaseComponent
	cfg        *config.Config
	collisions *Collisions
	sprites    *SpriteHandler
	graphics   *Graphics
	borders    *Borders
	memory     *Memory

	label     string
	reads     [RegisterCount]func() uint8
	writes    [RegisterCount]func(uint8)
	sequencer *Sequencer

	socketCycle           func() uint64
	socketBALow           func(bool)
	socketAECLow          func(bool)
	socketIRQTrigger      func()
	socketIRQClearTrigger func()
	socketLastCycle       func()
	socketVBlank          func()

	cr1    uint8 // VIC register
	cr2    uint8 // VIC register
	lpx    uint8 // VIC register
	lpy    uint8 // VIC register
	denBit bool

	irqLatch  uint8  // irqLatch holds an 8-bit value that latches the IRQ (Interrupt Request) configuration.
	irqMask   uint8  // irqMask represents an 8-bit mask used for interrupt request (IRQ) management.
	irqRaster uint16 // Interrupt raster line

	rasterX          uint16 // Current raster x position
	rasterY          uint16 // Current raster line
	lpTriggered      bool   // LightPen was triggered in this frame
	badLineEnabler   bool   // Bad Lines enabled for this frame
	badLineCondition bool   // Current line is bad line
	baLow            bool   // BA Line
	aecLow           bool   // AEC Line
	aecLowNextCycle  uint64 // aecLowNextCycle represents the counter for the next cycle in the AEC low-level operation.
	vBlankNextCycle  bool   // vBlankNextCycle indicates whether the next cycle will trigger a vertical blanking interval (vBlank) in the display.
	lineStart        int
	drawLine         bool
}

// NewVIC creates and initializes a new VIC instance with default configuration and registers it with the parent component.
func NewVIC(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *VIC {
	vic := &VIC{
		BaseComponent:    component.NewBaseComponent(),
		cr1:              0,
		cr2:              0,
		lpx:              0,
		lpy:              0,
		irqRaster:        0,
		irqLatch:         0,
		irqMask:          0,
		rasterX:          0,
		rasterY:          0,
		lpTriggered:      false,
		badLineCondition: false,
		badLineEnabler:   false,
		baLow:            false,
		aecLowNextCycle:  0,
		aecLow:           false,
		denBit:           false,
		label:            label,
	}
	vic.BaseComponent.Register(factory, parent, Identifier(), vic, references.IdIMos6569(vic, label, instance))
	return vic
}

// Setup initializes the VIC instance by retrieving and assigning its configuration from the factory.
func (vic *VIC) Setup() error {
	vic.cfg = vic.GetFactory().GetConfig()
	return nil
}

// Bind initializes the VIC chip by connecting it to the provided MOS 6569 socket and setting up its components.
// It configures internal structures like sequencers, memory, collisions, graphics, sprites, and borders.
// Returns an error if any component setup fails during initialization.
func (vic *VIC) Bind(socket references.IMos6569Socket) error {
	displayBuffer := vic.GetFactory().GetIDisplayBuffer()

	vic.socketIRQTrigger = socket.IRQTrigger
	vic.socketIRQClearTrigger = socket.IRQClearTrigger
	vic.socketCycle = socket.Cycle
	vic.socketLastCycle = socket.LastCycle
	vic.socketBALow = socket.BALow
	vic.socketAECLow = socket.AECLow
	vic.socketVBlank = socket.VBlank

	vic.sequencer = NewSequencerPal()
	vic.rasterY = vic.sequencer.rasterYMax

	vic.memory = NewMemory(vic, vic.GetFactory(), vic.label, 0, socket.ReadRam, socket.ReadColorRam, socket.ReadCharRom)
	vic.collisions = NewCollisions(vic, vic.GetFactory(), vic.label, 0, vic.irqEmit, vic.sequencer.width)
	vic.graphics = NewGraphics(vic, vic.GetFactory(), vic.label, 0, vic.memory, vic.collisions, displayBuffer, vic.sequencer.rasterYMax)
	vic.sprites = NewSprites(vic, vic.GetFactory(), vic.label, 0, vic.memory, vic.collisions, displayBuffer)
	vic.borders = NewBorder(vic, vic.GetFactory(), vic.label, 0, displayBuffer, vic.sequencer.width)

	vic.borders.SetDYTop(vic.sequencer.row24YStart)
	vic.borders.SetDYBottom(vic.sequencer.row24YStop)
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
	if err := vic.memory.Setup(); err != nil {
		return err
	}

	vic.reads = vic.createReadRegister()
	vic.writes = vic.createWriteRegister()

	return nil
}

// Connect establishes a connection using the VIC instance and returns an error if the connection fails.
func (vic *VIC) Connect() error {
	return nil
}

// Internal determines internal functionality status and returns a boolean indicating its state.
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
	return vic.memory.GetLastByte()
}

// configChanged handles updates to the VIC configuration and applies necessary changes to its state.
func (vic *VIC) configChanged() {
}

// Emulate executes one cycle of the VIC, processing the current function and updating the raster position.
//
//go:nosplit
func (vic *VIC) Emulate() {
	vic.TryAcquireAEC()
	vic.sequencer.Sequence(vic)
	vic.UpdateRasterX()
}

// EmulationRequired returns true if emulation is required for the current VIC (Video Interface Controller) state.
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

// SetBALow sets the BA (bus-available) signal to low and schedules the AEC signal to be low after 3 cycles if not already set.
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

// badLineUpdate updates the bad line condition based on the current raster position, DEN bit, and YSCROLL value.
// The bad line condition occurs when specific raster and scroll conditions are met, enabling certain VIC behavior.
func (vic *VIC) badLineUpdate() {
	// Bad Line Condition is given at any arbitrary clock cycle, if at the
	// negative edge of ø0 at the beginning of the cycle RASTER >= $30 and RASTER <= $f7
	// and the lower three bits of RASTER are equal to YSCROLL
	// and if the DEN bit has been set for at least one cycle somewhere in raster line $30
	// So clearing the DEN bit will normally prevent Bad Lines

	if (vic.rasterY >= vic.sequencer.firstDmaLine) && (vic.rasterY <= vic.sequencer.lastDmaLine) {
		if vic.rasterY == vic.sequencer.firstDmaLine && vic.denBit {
			//If YSCROLL=0, a Bad Line Condition occurs in raster line $30 as soon as the DEN bit
			vic.badLineEnabler = true
			if vic.graphics.GetYScroll() == 0 {
				vic.badLineCondition = true
				return
			}
		}
		if vic.badLineEnabler {
			vic.badLineCondition = vic.graphics.GetYScroll() == (vic.rasterY & 7)
		}
	} else {
		vic.badLineEnabler = false
		vic.badLineCondition = false
	}
}

// ChangedVA updates the VIC's virtual address base and triggers the memory pointer update process.
func (vic *VIC) ChangedVA(newVA uint8) {
	vic.memory.SetCIAVABase(newVA)
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
	vic.memory.ResetRefreshCounter()
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

func (vic *VIC) irqSetLatch(data uint8) {
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

// irqSetMask updates the IRQ mask register and triggers verification, potentially emitting an interrupt.
func (vic *VIC) irqSetMask(data uint8) {
	vic.irqMask = data & 0xf
	vic.irqVerify()
}

// irqSetRasterHigh sets the high byte of the raster interrupt compare value and updates the internal raster counter.
func (vic *VIC) irqSetRasterHigh(data uint8) {
	irqRaster := (vic.irqRaster & 0xff00) | uint16(data)
	vic.irqRasterSet(irqRaster)
}

// irqSetRasterLow sets the lower byte of the IRQ raster value and updates the raster, potentially triggering an IRQ.
func (vic *VIC) irqSetRasterLow(data uint16) {
	irqRaster := (vic.irqRaster & 0xff) | data
	vic.irqRasterSet(irqRaster)
}

// rasterUpdate updates the VIC raster interrupt value and triggers an interrupt if the raster line matches the new value.
func (vic *VIC) irqRasterSet(irqRaster uint16) {
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

// setCR1 updates the control register CR1 and adjusts various graphical and border settings based on the given data.
func (vic *VIC) setCR1(data uint8) {
	vic.cr1 = data
	vic.graphics.SetYScroll(uint16(vic.cr1) & 7)
	if rowSel := (vic.cr1 & 0x8) != 0; rowSel {
		vic.borders.SetDYTop(vic.sequencer.row25YStart)
		vic.borders.SetDYBottom(vic.sequencer.row25YStop)
	} else {
		vic.borders.SetDYTop(vic.sequencer.row24YStart)
		vic.borders.SetDYBottom(vic.sequencer.row24YStop)
	}
	vic.denBit = (vic.cr1 & 0x10) != 0
	vic.graphics.SetBmm((vic.cr1 & 0x20) != 0)
	vic.graphics.SetEcm((vic.cr1 & 0x40) != 0)
	//rst8 := (vic.cr1 & 0x80) != 0
	displayMode := ((int(vic.cr1) & 0x60) | (int(vic.cr2) & 0x10)) >> 4 //cr1 bit 5-6 (BMM|ECM)| cr2 bit 4 (MCM)
	vic.graphics.SetDisplayMode(displayMode)
	vic.irqSetRasterLow((uint16(vic.cr1) & 0x80) << 1)
	vic.badLineUpdate()
}

// setCR2 updates the CR2 register and adjusts associated graphics settings like XScroll, column selection, and display mode.
func (vic *VIC) setCR2(data uint8) {
	vic.cr2 = data
	vic.graphics.SetXScroll(uint16(vic.cr2) & 7)
	vic.borders.SetColumnSel((vic.cr2 & 0x8) != 0)
	displayMode := ((int(vic.cr1) & 0x60) | (int(vic.cr2) & 0x10)) >> 4 //cr1 bit 5-6 (BMM|ECM)| cr2 bit 4 (MCM)
	vic.graphics.SetDisplayMode(displayMode)
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
	reads[0x00] = vic.sprites.ReadMXx0
	reads[0x01] = vic.sprites.ReadMXy0
	reads[0x02] = vic.sprites.ReadMXx1
	reads[0x03] = vic.sprites.ReadMXy1
	reads[0x04] = vic.sprites.ReadMXx2
	reads[0x05] = vic.sprites.ReadMXy2
	reads[0x06] = vic.sprites.ReadMXx3
	reads[0x07] = vic.sprites.ReadMXy3
	reads[0x08] = vic.sprites.ReadMXx4
	reads[0x09] = vic.sprites.ReadMXy4
	reads[0x0a] = vic.sprites.ReadMXx5
	reads[0x0b] = vic.sprites.ReadMXy5
	reads[0x0c] = vic.sprites.ReadMXx6
	reads[0x0d] = vic.sprites.ReadMXy6
	reads[0x0e] = vic.sprites.ReadMXx7
	reads[0x0f] = vic.sprites.ReadMXy7
	reads[0x10] = vic.sprites.ReadMX8
	reads[0x11] = func() uint8 { return uint8((uint16(vic.cr1) & 0x7f) | ((vic.rasterY & 0x100) >> 1)) } // Control register 1
	reads[0x12] = func() uint8 { return uint8(vic.rasterY) }                                             // Raster counter
	reads[0x13] = func() uint8 { return vic.lpx }                                                        // Light pen X
	reads[0x14] = func() uint8 { return vic.lpy }                                                        // Light pen Y
	reads[0x15] = vic.sprites.ReadMe                                                                     // Sprite enabled
	reads[0x16] = func() uint8 { return vic.cr2 | 0xc0 }                                                 // Control register 2
	reads[0x17] = vic.sprites.ReadMYe                                                                    // Sprite Y expansion
	reads[0x18] = vic.memory.GetVABase                                                                   // Memory pointers
	reads[0x19] = func() uint8 { return vic.irqLatch | 0x70 }                                            // IRQ latch
	reads[0x1a] = func() uint8 { return vic.irqMask | 0xf0 }                                             // IRQ mask
	reads[0x1b] = vic.sprites.ReadMDp                                                                    // Sprite data priority
	reads[0x1c] = vic.sprites.ReadMMc                                                                    // Sprite multicolor
	reads[0x1d] = vic.sprites.ReadMXe                                                                    // Sprite X expansion
	reads[0x1e] = vic.collisions.RetrieveSprite2Sprite
	reads[0x1f] = vic.collisions.RetrieveSprite2Background
	reads[0x20] = vic.borders.ReadEc
	reads[0x21] = vic.graphics.ReadB0c
	reads[0x22] = vic.graphics.ReadB1c
	reads[0x23] = vic.graphics.ReadB2c
	reads[0x24] = vic.graphics.ReadB3c
	reads[0x25] = vic.sprites.ReadMM0
	reads[0x26] = vic.sprites.ReadMM1
	reads[0x27] = vic.sprites.ReadMXc0
	reads[0x28] = vic.sprites.ReadMXc1
	reads[0x29] = vic.sprites.ReadMXc2
	reads[0x2a] = vic.sprites.ReadMXc3
	reads[0x2b] = vic.sprites.ReadMXc4
	reads[0x2c] = vic.sprites.ReadMXc5
	reads[0x2d] = vic.sprites.ReadMXc6
	reads[0x2e] = vic.sprites.ReadMXc7
	return reads
}

// createWriteRegister initializes an array of functions for writing data to various VIC-II registers.
func (vic *VIC) createWriteRegister() [RegisterCount]func(uint8) {
	var writes [RegisterCount]func(uint8)
	var unconnected = func(uint8) {}
	for idx := range writes {
		writes[idx] = unconnected
	}
	writes[0x00] = vic.sprites.WriteMXx0
	writes[0x01] = vic.sprites.WriteMXy0
	writes[0x02] = vic.sprites.WriteMXx1
	writes[0x03] = vic.sprites.WriteMXy1
	writes[0x04] = vic.sprites.WriteMXx2
	writes[0x05] = vic.sprites.WriteMXy2
	writes[0x06] = vic.sprites.WriteMXx3
	writes[0x07] = vic.sprites.WriteMXy3
	writes[0x08] = vic.sprites.WriteMXx4
	writes[0x09] = vic.sprites.WriteMXy4
	writes[0x0a] = vic.sprites.WriteMXx5
	writes[0x0b] = vic.sprites.WriteMXy5
	writes[0x0c] = vic.sprites.WriteMXx6
	writes[0x0d] = vic.sprites.WriteMXy6
	writes[0x0e] = vic.sprites.WriteMXx7
	writes[0x0f] = vic.sprites.WriteMXy7
	writes[0x10] = vic.sprites.WriteMX8
	writes[0x11] = vic.setCR1                          // Control register 1
	writes[0x12] = vic.irqSetRasterHigh                // Raster counter
	writes[0x13] = func(data uint8) { vic.lpx = data } // Light pen X
	writes[0x14] = func(data uint8) { vic.lpy = data } // Light pen Y
	writes[0x15] = vic.sprites.WriteMe                 // Sprite enabled
	writes[0x16] = vic.setCR2                          // Control register 2
	writes[0x17] = vic.sprites.WriteMYe                // Sprite Y expansion
	writes[0x18] = vic.memory.SetVABase                // Memory pointers
	writes[0x19] = vic.irqSetLatch                     // IRQ Latch
	writes[0x1a] = vic.irqSetMask                      // IRQ mask
	writes[0x1b] = vic.sprites.WriteMDp                // Sprite data priority
	writes[0x1c] = vic.sprites.WriteMMc                // Sprite Color
	writes[0x1d] = vic.sprites.WriteMXe                // Sprite X expansion
	writes[0x1e] = vic.collisions.SetSprite            // Sprite-sprite collision
	writes[0x1f] = vic.collisions.SetBackground        // Sprite-background collision
	writes[0x20] = vic.borders.WriteEc
	writes[0x21] = vic.graphics.WriteB0c
	writes[0x22] = vic.graphics.WriteB1c
	writes[0x23] = vic.graphics.WriteB2c
	writes[0x24] = vic.graphics.WriteB3c
	writes[0x25] = vic.sprites.WriteMM0
	writes[0x26] = vic.sprites.WriteMM1
	writes[0x27] = vic.sprites.WriteMXc0
	writes[0x28] = vic.sprites.WriteMXc1
	writes[0x29] = vic.sprites.WriteMXc2
	writes[0x2a] = vic.sprites.WriteMXc3
	writes[0x2b] = vic.sprites.WriteMXc4
	writes[0x2c] = vic.sprites.WriteMXc5
	writes[0x2d] = vic.sprites.WriteMXc6
	writes[0x2e] = vic.sprites.WriteMXc7
	return writes
}
