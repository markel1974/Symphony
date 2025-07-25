package mos6569

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// https://www.cebix.net/VIC-Article.txt
// https://www.oxyron.de/html/registers_vic2.html

const (
	RegisterSize  = 0x3f
	RegisterCount = RegisterSize + 1
)

//https://dustlayer.com/c64-architecture

type ISequencer interface {
	Setup() error

	Sequence(vic *VIC)

	GetRasterYMax() uint16

	GetWidth() int

	GetRow24YStart() uint16

	GetRow24YStop() uint16

	GetRow25YStart() uint16

	GetRow25YStop() uint16

	GetFirstDmaLine() uint16

	GetLastDmaLine() uint16
}

// VIC represents a versatile interface controller for managing video output and graphical resources in a system.
// It encapsulates configurations, graphics components, collision detection, and rendering capabilities for the display.
type VIC struct {
	*component.BaseComponent
	cfg        *config.Config
	memory     *MemoryUnit
	interrupts *Interrupts
	collisions *CollisionsUnit
	sprites    *SpritesUnit
	graphics   *GraphicsUnit
	borders    *BordersUnit
	lightPen   *LightPen

	sequencer ISequencer
	label     string
	reads     [RegisterCount]func() uint8
	writes    [RegisterCount]func(uint8)

	sequencerSequence func(vic *VIC)

	socketCycle     func() uint64
	socketBALow     func(bool)
	socketAECLow    func(bool)
	socketLastCycle func()
	socketVBlank    func()

	sequencerRow25YStart  uint16
	sequencerRow25YStop   uint16
	sequencerRow24YStart  uint16
	sequencerRow24YStop   uint16
	sequencerFirstDmaLine uint16
	sequencerLastDmaLine  uint16

	cr1 uint8 // VIC register
	cr2 uint8 // VIC register

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
		badLineCondition: false,
		badLineEnabler:   false,
		baLow:            false,
		aecLowNextCycle:  0,
		aecLow:           false,
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

	vic.socketCycle = socket.Cycle
	vic.socketLastCycle = socket.LastCycle
	vic.socketBALow = socket.BALow
	vic.socketAECLow = socket.AECLow
	vic.socketVBlank = socket.VBlank

	if socket.TotalRaster() > 300 {
		vic.sequencer = NewSequencerPal(vic, vic.GetFactory(), vic.label, 0)
	} else {
		vic.sequencer = NewSequencerNtsc(vic, vic.GetFactory(), vic.label, 0)
	}
	vic.sequencerSequence = vic.sequencer.Sequence
	vic.sequencerRow25YStart = vic.sequencer.GetRow25YStart()
	vic.sequencerRow25YStop = vic.sequencer.GetRow25YStop()
	vic.sequencerRow24YStart = vic.sequencer.GetRow24YStart()
	vic.sequencerRow24YStop = vic.sequencer.GetRow24YStop()
	vic.sequencerFirstDmaLine = vic.sequencer.GetFirstDmaLine()
	vic.sequencerLastDmaLine = vic.sequencer.GetLastDmaLine()

	vic.memory = NewMemory(vic, vic.GetFactory(), vic.label, 0, socket.ReadRam, socket.ReadColorRam, socket.ReadCharRom)
	vic.interrupts = NewInterrupts(vic, vic.GetFactory(), vic.label, 0, socket.IRQTrigger, socket.IRQClearTrigger, vic.sequencer.GetRasterYMax())
	vic.collisions = NewCollisions(vic, vic.GetFactory(), vic.label, 0, vic.interrupts.Emit, vic.sequencer.GetWidth())
	vic.graphics = NewGraphics(vic, vic.GetFactory(), vic.label, 0, vic.memory, vic.collisions, displayBuffer, vic.sequencer.GetRasterYMax())
	vic.sprites = NewSprites(vic, vic.GetFactory(), vic.label, 0, vic.memory, vic.collisions, displayBuffer)
	vic.borders = NewBorder(vic, vic.GetFactory(), vic.label, 0, displayBuffer, vic.sequencer.GetWidth())
	vic.lightPen = NewLightPen(vic, vic.GetFactory(), vic.label, 0, vic.interrupts.Emit)
	vic.borders.SetDYTop(vic.sequencer.GetRow24YStart())
	vic.borders.SetDYBottom(vic.sequencer.GetRow24YStop())
	vic.vBlankNextCycle = false
	vic.drawLine = false
	vic.cfg.Bind(vic.configChanged)

	if err := vic.sequencer.Setup(); err != nil {
		return err
	}
	if err := vic.memory.Setup(); err != nil {
		return err
	}
	if err := vic.interrupts.Setup(); err != nil {
		return err
	}
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
	if err := vic.lightPen.Setup(); err != nil {
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
	vic.sequencerSequence(vic)
	vic.interrupts.RasterXIncrement()
}

// EmulationRequired returns true if emulation is required for the current VIC (Video Interface Controller) state.
func (vic *VIC) EmulationRequired() bool {
	return true
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
//
//go:nosplit
func (vic *VIC) TryAcquireAEC() {
	if vic.baLow && !vic.aecLow {
		if vic.socketCycle() >= vic.aecLowNextCycle {
			vic.aecLow = true
			vic.socketAECLow(true)
		}
	}
}

// BadLineUpdate updates the bad line condition based on the current raster position, DEN bit, and YSCROLL value.
// The bad line condition occurs when specific raster and scroll conditions are met, enabling certain VIC behavior.
func (vic *VIC) BadLineUpdate() {
	// Bad Line Condition is given at any arbitrary clock cycle, if at the
	// negative edge of ø0 at the beginning of the cycle RASTER >= $30 and RASTER <= $f7
	// and the lower three bits of RASTER are equal to YSCROLL
	// and if the DEN bit has been set for at least one cycle somewhere in raster line $30
	// So clearing the DEN bit will normally prevent Bad Lines
	rasterY := vic.interrupts.RasterY()

	if (rasterY >= vic.sequencerFirstDmaLine) && (rasterY <= vic.sequencerLastDmaLine) {
		if rasterY == vic.sequencerFirstDmaLine && vic.borders.GetDen() {
			//If YSCROLL=0, a Bad Line Condition occurs in raster line $30 as soon as the DEN bit
			vic.badLineEnabler = true
			if vic.graphics.GetYScroll() == 0 {
				vic.badLineCondition = true
				return
			}
		}
		if vic.badLineEnabler {
			vic.badLineCondition = vic.graphics.GetYScroll() == (rasterY & 7)
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

// LightPenTrigger handles triggering the light pen functionality based on the current raster position.
func (vic *VIC) LightPenTrigger() {
	vic.lightPen.Trigger(vic.interrupts.RasterY(), vic.interrupts.RasterY())
}

// WriteCR1 updates the control register CR1 and adjusts various graphical and border settings based on the given data.
func (vic *VIC) WriteCR1(data uint8) {
	vic.cr1 = data
	vic.graphics.SetYScroll(uint16(vic.cr1) & 7)
	if rowSel := (vic.cr1 & 0x8) != 0; rowSel {
		vic.borders.SetDYTop(vic.sequencerRow25YStart)
		vic.borders.SetDYBottom(vic.sequencerRow25YStop)
	} else {
		vic.borders.SetDYTop(vic.sequencerRow24YStart)
		vic.borders.SetDYBottom(vic.sequencerRow24YStop)
	}
	vic.borders.SetDen((vic.cr1 & 0x10) != 0)
	vic.graphics.SetBmm((vic.cr1 & 0x20) != 0)
	vic.graphics.SetEcm((vic.cr1 & 0x40) != 0)
	//rst8 := (vic.cr1 & 0x80) != 0
	displayMode := ((int(vic.cr1) & 0x60) | (int(vic.cr2) & 0x10)) >> 4 //cr1 bit 5-6 (BMM|ECM)| cr2 bit 4 (MCM)
	vic.graphics.SetDisplayMode(displayMode)
	vic.interrupts.WriteRasterLow((uint16(vic.cr1) & 0x80) << 1)

	vic.BadLineUpdate()
}

// ReadCR1 reads the CR1 register, combining specific raster and control bits, and returns the resulting 8-bit value.
func (vic *VIC) ReadCR1() uint8 {
	return uint8((uint16(vic.cr1) & 0x7f) | ((vic.interrupts.RasterY() & 0x100) >> 1))
}

// WriteCR2 updates the CR2 register and adjusts associated graphics settings like XScroll, column selection, and display mode.
func (vic *VIC) WriteCR2(data uint8) {
	vic.cr2 = data
	vic.graphics.SetXScroll(uint16(vic.cr2) & 7)
	vic.borders.SetColumnSel((vic.cr2 & 0x8) != 0)
	displayMode := ((int(vic.cr1) & 0x60) | (int(vic.cr2) & 0x10)) >> 4 //cr1 bit 5-6 (BMM|ECM)| cr2 bit 4 (MCM)
	vic.graphics.SetDisplayMode(displayMode)
}

// ReadCR2 reads the CR2 register value, applies a bitwise OR with 0xC0, and returns the result as an 8-bit unsigned integer.
func (vic *VIC) ReadCR2() uint8 {
	return vic.cr2 | 0xc0
}

// ReadRegister reads a register at the given address and returns the corresponding 8-bit value.
//
//go:nosplit
func (vic *VIC) ReadRegister(addr uint16) uint8 {
	return vic.reads[addr&RegisterSize]()
}

// WriteRegister writes data to a register at the specified address, handling various control and memory settings.
//
//go:nosplit
func (vic *VIC) WriteRegister(addr uint16, data uint8) {
	vic.writes[addr&RegisterSize](data)
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
	reads[0x11] = vic.ReadCR1                // Control register 1
	reads[0x12] = vic.interrupts.ReadRasterY // Raster counter
	reads[0x13] = vic.lightPen.ReadX         // Light pen X
	reads[0x14] = vic.lightPen.ReadY         // Light pen Y
	reads[0x15] = vic.sprites.ReadMe         // Sprite enabled
	reads[0x16] = vic.ReadCR2                // Control register 2
	reads[0x17] = vic.sprites.ReadMYe        // Sprite Y expansion
	reads[0x18] = vic.memory.GetVABase       // MemoryUnit pointers
	reads[0x19] = vic.interrupts.ReadLatch   // IRQ latch
	reads[0x1a] = vic.interrupts.ReadMask    // IRQ mask
	reads[0x1b] = vic.sprites.ReadMDp        // Sprite data priority
	reads[0x1c] = vic.sprites.ReadMMc        // Sprite multicolor
	reads[0x1d] = vic.sprites.ReadMXe        // Sprite X expansion
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
	writes[0x11] = vic.WriteCR1                   // Control register 1
	writes[0x12] = vic.interrupts.WriteRasterHigh // Raster counter
	writes[0x13] = vic.lightPen.WriteX            // Light pen X
	writes[0x14] = vic.lightPen.WriteY            // Light pen Y
	writes[0x15] = vic.sprites.WriteMe            // Sprite enabled
	writes[0x16] = vic.WriteCR2                   // Control register 2
	writes[0x17] = vic.sprites.WriteMYe           // Sprite Y expansion
	writes[0x18] = vic.memory.SetVABase           // MemoryUnit pointers
	writes[0x19] = vic.interrupts.WriteLatch      // IRQ Latch
	writes[0x1a] = vic.interrupts.WriteMask       // IRQ mask
	writes[0x1b] = vic.sprites.WriteMDp           // Sprite data priority
	writes[0x1c] = vic.sprites.WriteMMc           // Sprite Color
	writes[0x1d] = vic.sprites.WriteMXe           // Sprite X expansion
	writes[0x1e] = vic.collisions.SetSprite       // Sprite-sprite collision
	writes[0x1f] = vic.collisions.SetBackground   // Sprite-background collision
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
