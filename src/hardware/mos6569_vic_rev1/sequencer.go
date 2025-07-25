package mos6569

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

const drawLoopCycles = 36 // drawLoopCycles defines the number of iterations for the main display loop in video cycle data processing.

// SequencerData represents a chainable unit in a sequencer containing a function, a next node, and cycle information.
type SequencerData struct {
	fn          func(vic *VIC)
	next        *SequencerData
	cycle       uint8
	cycleBorder uint8
}

// NewSequencerData creates and returns a new instance of SequencerData with the provided function initialized.
func NewSequencerData(fn func(vic *VIC)) *SequencerData {
	return &SequencerData{fn: fn}
}

// Sequencer represents a structure to manage video cycle sequencing and raster timing for display rendering.
type Sequencer struct {
	*component.BaseComponent
	width              int
	height             int
	totalRasters       uint16 // Total number of raster lines (PAL)
	firstDisplayedLine uint16 // First displayed line
	lastDisplayedLine  uint16 // Last displayed line
	firstDmaLine       uint16 // First possible line for Bad Lines
	lastDmaLine        uint16 // Last possible line for Bad Lines
	row25YStart        uint16 // Display window coordinates
	row25YStop         uint16
	row24YStart        uint16
	row24YStop         uint16
	rasterYMax         uint16
	displaySize        int
	borderFirstCycle   uint8
	data               []*SequencerData
	curr               *SequencerData
}

// NewSequencer creates a Sequencer based on the total number of raster lines, selecting PAL or NTSC timing accordingly.
func NewSequencer(parent references.IComponent, factory references.IComponentFactory, label string, instance int /* screenFreq */, _ int, totalRaster int) *Sequencer {
	if totalRaster >= 300 {
		return newSequencerPal(parent, factory, label, instance)
	}
	return newSequencerNtsc(parent, factory, label, instance)
}

// newSequencerPal initializes the PAL video timing cycle data. It constructs a circular linked list of 63 cycleData nodes,
// where each node represents one CPU clock cycle of a single PAL scanline. It pre-calculates border-related
// values for each cycle and links them in sequence to form the complete 63-cycle sequencer.
func newSequencerPal(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Sequencer {
	const palWidth = 384
	const palHeight = 272
	const palTotalRasters = 312

	seq := &Sequencer{
		BaseComponent:      component.NewBaseComponent(),
		width:              palWidth,
		height:             palHeight,
		totalRasters:       palTotalRasters, // Total number of raster lines
		firstDisplayedLine: 16,              // First displayed line
		lastDisplayedLine:  287,             // Last displayed line
		firstDmaLine:       48,              // First possible line for Bad Lines
		lastDmaLine:        247,             // Last possible line for Bad Lines
		row25YStart:        51,              //
		row25YStop:         251,             //
		rasterYMax:         palTotalRasters - 1,
		displaySize:        (palWidth + 64) * palHeight,
		borderFirstCycle:   13,
	}

	seq.prepare()

	seq.add(seq.phaseSprite3DMAPhase1AndInit)
	seq.add(seq.phaseSprite3DMAPhase2AndVBlank)
	seq.add(seq.phaseSprite4DMAPhase1)
	seq.add(seq.phaseSprite4DMAPhase2)
	seq.add(seq.phaseSprite5DMAPhase1)
	seq.add(seq.phaseSprite5DMAPhase2)
	seq.add(seq.phaseSprite6DMAPhase1)
	seq.add(seq.phaseSprite6DMAPhase2)
	seq.add(seq.phaseSprite7DMAPhase1)
	seq.add(seq.phaseSprite7DMAPhase2)
	seq.add(seq.phaseRefresh)
	seq.add(seq.phaseSetupBadLineCheck)
	seq.add(seq.phaseSetupRasterXReset)
	seq.add(seq.phaseSetupVCounterLoad)
	seq.add(seq.phaseSetupRCounterCheckAndSpritePipe1)
	seq.add(seq.phaseDisplayFirstFetchAndSpritePipe2)
	seq.add(seq.phaseDisplayMainFetchC40)
	seq.add(seq.phaseDisplayMainFetchC38)
	for i := 0; i < drawLoopCycles; i++ {
		seq.add(seq.phaseDisplayMainFetch)
	}
	seq.add(seq.phaseDMASetupAndTeardownLastFetch)
	seq.add(seq.phaseTeardownIdle)
	seq.add(seq.phaseTeardownCommitSpriteFlags)
	seq.add(seq.phaseSprite0DMAPhase1)
	seq.add(seq.phaseSprite0DMAPhase2)
	seq.add(seq.phaseSprite1DMAPhase1)
	seq.add(seq.phaseSprite1DMAPhase2)
	seq.add(seq.phaseSprite2DMAPhase1)
	seq.add(seq.phaseSprite2DMAPhase2AndTeardownFinal)

	seq.finalize()

	seq.BaseComponent.Register(factory, parent, "sequencerPal", seq, references.IdInternalComponent(label, instance, "SequencerPal"))

	return seq
}

// newSequencerNtsc initializes the NTSC video timing cycle data with 63 cycleData nodes for a single scanline.
// Constructs a sequencer containing pre-calculated raster and border values for 263 total raster lines.
// Links the nodes sequentially to form a complete NTSC video frame.
func newSequencerNtsc(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Sequencer {
	const ntscWidth = 384
	const ntscHeight = 240
	const ntscTotalRasters = 263

	seq := &Sequencer{
		BaseComponent:      component.NewBaseComponent(),
		width:              ntscWidth,
		height:             ntscHeight,
		totalRasters:       ntscTotalRasters,
		firstDisplayedLine: 13,  // Lowest value to include top border
		lastDisplayedLine:  261, // Highest value to include bottom border
		firstDmaLine:       48,  //
		lastDmaLine:        247, //
		row25YStart:        51,  //
		row25YStop:         251, //
		rasterYMax:         ntscTotalRasters - 1,
		displaySize:        (ntscWidth + 64) * ntscHeight,
		borderFirstCycle:   15,
	}

	seq.prepare()

	seq.add(seq.phaseSprite3DMAPhase1AndInit)
	seq.add(seq.phaseSprite3DMAPhase2AndVBlank)
	seq.add(seq.phaseSprite4DMAPhase1)
	seq.add(seq.phaseSprite4DMAPhase2)
	seq.add(seq.phaseSprite5DMAPhase1)
	seq.add(seq.phaseSprite5DMAPhase2)
	seq.add(seq.phaseSprite6DMAPhase1)
	seq.add(seq.phaseSprite6DMAPhase2)
	seq.add(seq.phaseSprite7DMAPhase1)
	seq.add(seq.phaseSprite7DMAPhase2)
	seq.add(seq.phaseRefresh)
	seq.add(seq.phaseSetupBadLineCheck)
	seq.add(seq.phaseSetupRasterXReset)
	seq.add(seq.phaseSetupVCounterLoad)
	seq.add(seq.phaseSetupRCounterCheckAndSpritePipe1)
	seq.add(seq.phaseDisplayFirstFetchAndSpritePipe2)
	seq.add(seq.phaseDisplayMainFetchC40)
	seq.add(seq.phaseDisplayMainFetchC38)
	for i := 0; i < drawLoopCycles; i++ {
		seq.add(seq.phaseDisplayMainFetch)
	}
	seq.add(seq.phaseDMASetupAndTeardownLastFetch)
	// The two extra NTSC cycles are refresh cycles.
	seq.add(seq.phaseRefresh)
	seq.add(seq.phaseRefresh)
	seq.add(seq.phaseTeardownIdle)
	seq.add(seq.phaseTeardownCommitSpriteFlags)
	seq.add(seq.phaseSprite0DMAPhase1)
	seq.add(seq.phaseSprite0DMAPhase2)
	seq.add(seq.phaseSprite1DMAPhase1)
	seq.add(seq.phaseSprite1DMAPhase2)
	seq.add(seq.phaseSprite2DMAPhase1)
	seq.add(seq.phaseSprite2DMAPhase2AndTeardownFinal)

	seq.finalize()

	seq.BaseComponent.Register(factory, parent, "sequencerPal", seq, references.IdInternalComponent(label, instance, "SequencerPal"))

	return seq
}

// Setup initializes the Sequencer, preparing it for execution and ensuring all required dependencies are in place.
func (seq *Sequencer) Setup() error {
	return nil
}

// Connect establishes a connection for the Sequencer, returning an error if the connection process fails.
func (seq *Sequencer) Connect() error {
	return nil
}

// EmulationRequired determines if emulation is needed based on the Sequencer's state, returning a boolean value.
func (seq *Sequencer) EmulationRequired() bool {
	return false
}

// Emulate runs the sequencer emulation process, executing its predefined operations and behaviors in sequence.
func (seq *Sequencer) Emulate() {
}

// Internal checks the internal state or logic of the Sequencer and returns a boolean indicating the result.
func (seq *Sequencer) Internal() bool {
	return true
}

// Reset reinitializes the internal state of the Sequencer to its default starting condition.
func (seq *Sequencer) Reset() {
}

// GetRasterYMax returns the maximum Y-coordinate value of the raster as uint16.
func (seq *Sequencer) GetRasterYMax() uint16 {
	return seq.rasterYMax
}

// GetWidth returns the width of the Sequencer instance.
func (seq *Sequencer) GetWidth() int {
	return seq.width
}

// GetFirstDmaLine retrieves the first DMA line of the Sequencer.
func (seq *Sequencer) GetFirstDmaLine() uint16 {
	return seq.firstDmaLine
}

// GetLastDmaLine retrieves the last DMA line value stored in the Sequencer instance.
func (seq *Sequencer) GetLastDmaLine() uint16 {
	return seq.lastDmaLine
}

// GetRow24YStart returns the starting Y-coordinate value for row 24.
func (seq *Sequencer) GetRow24YStart() uint16 {
	return seq.row24YStart
}

// GetRow24YStop retrieves the value of the row24YStop property from the Sequencer.
func (seq *Sequencer) GetRow24YStop() uint16 {
	return seq.row24YStop
}

// GetRow25YStart returns the Y-axis starting position for row 25 within the sequencer.
func (seq *Sequencer) GetRow25YStart() uint16 {
	return seq.row25YStart
}

// GetRow25YStop retrieves the value of row25YStop from the Sequencer instance.
func (seq *Sequencer) GetRow25YStop() uint16 {
	return seq.row25YStop
}

// prepare initializes the Sequencer by creating an empty slice of SequencerData pointers to be used for cycle management.
func (seq *Sequencer) prepare() {
	seq.row24YStart = seq.row25YStart + 4
	seq.row24YStop = seq.row25YStop - 4
	seq.data = make([]*SequencerData, 0)
}

// add appends a new SequencerData instance created with the provided function to the Sequencer's data slice.
func (seq *Sequencer) add(fn func(vic *VIC)) {
	seq.data = append(seq.data, NewSequencerData(fn))
}

// finalize constructs a circular linked list from the sequence data and calculates cycle-specific border values.
func (seq *Sequencer) finalize() {
	last := len(seq.data) - 1
	for idx := 0; idx < len(seq.data); idx++ {
		seq.data[idx].cycleBorder = 0xff
		seq.data[idx].cycle = uint8(idx) + 1
		if seq.data[idx].cycle >= seq.borderFirstCycle {
			seq.data[idx].cycleBorder = seq.data[idx].cycle - seq.borderFirstCycle
		}
		if idx == last {
			seq.data[idx].next = seq.data[0]
		} else {
			seq.data[idx].next = seq.data[idx+1]
		}
	}
	seq.curr = seq.data[0]
}

// Sequence processes the current function in the sequence and advances to the next step in the sequence chain.
//
//go:nosplit
func (seq *Sequencer) Sequence(vic *VIC) {
	seq.curr.fn(vic)
	seq.curr = seq.curr.next
}

// phaseSprite3DMAPhase1AndInit: This cycle marks the beginning of the horizontal blanking period. The raster line counter (sequencerRasterY)
// is checked against the maximum value. If it matches, a V-Blank is scheduled for the next cycle. Otherwise,
// sequencerRasterY is incremented for the new scanline. Sprite 3 DMA for the upcoming line begins if enabled, fetching
// the sprite pointer (phi1) and the first byte of sprite data (phi2).
//
//go:nosplit
func (seq *Sequencer) phaseSprite3DMAPhase1AndInit(vic *VIC) {
	vic.sprites.Prepare()

	if rasterY := vic.GetRasterY(); rasterY == seq.rasterYMax {
		vic.vBlankNextCycle = true
	} else {
		vic.IncrementRasterY()
		vic.drawLine = (rasterY >= seq.firstDisplayedLine) && (rasterY <= seq.lastDisplayedLine)
	}
	vic.borders.ColumnInitialize()
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	if vic.sprites.GetDMAFlag(bitSprite3) != 0 {
		vic.sprites.FetchPhase1(3)
	}
	if vic.sprites.GetDMAFlag(bitSprite3|bitSprite4) == 0 {
		vic.ClearBALow()
	}
}

// phaseSprite3DMAPhase2AndVBlank: The V-Blank is triggered if scheduled in the previous cycle, resetting raster counters and firing
// the V-Blank interrupt. Internal pointers for graphics, sprites, and borders are reset. Sprite 3 DMA
// continues, fetching the second and third bytes of sprite data. The BA (Bus-Available) signal is asserted
// (pulled low) to prepare for sprite 5 DMA if it is enabled for the upcoming line.
//
//go:nosplit
func (seq *Sequencer) phaseSprite3DMAPhase2AndVBlank(vic *VIC) {
	if vic.vBlankNextCycle {
		vic.vBlankNextCycle = false
		vic.lineStart = 0
		vic.graphics.ResetVideoCounterLatch()
		vic.ResetRasterY()
		vic.socketVBlank()
	}
	vic.collisions.ClearGraphics()
	vic.graphics.SetOffset(vic.lineStart)
	vic.sprites.SetOffset(vic.lineStart)
	vic.borders.SetOffset(vic.lineStart)
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	if vic.sprites.GetDMAFlag(bitSprite3) != 0 {
		vic.sprites.FetchPhase2(3)
	} else {
		vic.memory.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite5) != 0 {
		vic.SetBALow()
	}
}

// phaseSprite4DMAPhase1: Sprite 4 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 4 and 5.
//
//go:nosplit
func (seq *Sequencer) phaseSprite4DMAPhase1(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	if vic.sprites.GetDMAFlag(bitSprite4) != 0 {
		vic.sprites.FetchPhase1(4)
	}
	if vic.sprites.GetDMAFlag(bitSprite4|bitSprite5) == 0 {
		vic.ClearBALow()
	}
}

// phaseSprite4DMAPhase2: Sprite 4 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 6 DMA if it is enabled.
//
//go:nosplit
func (seq *Sequencer) phaseSprite4DMAPhase2(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	if vic.sprites.GetDMAFlag(bitSprite4) != 0 {
		vic.sprites.FetchPhase2(4)
		//vic.sprites.FetchData1(seq.spriteOdd, 4) //phi1
		//vic.sprites.FetchData2(seq.spriteOdd, 4) //phi2
	} else {
		vic.memory.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite6) != 0 {
		vic.SetBALow()
	}
}

// phaseSprite5DMAPhase1: Sprite 5 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 5 and 6.
//
//go:nosplit
func (seq *Sequencer) phaseSprite5DMAPhase1(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	if vic.sprites.GetDMAFlag(bitSprite5) != 0 {
		vic.sprites.FetchPhase1(5)
	}
	if vic.sprites.GetDMAFlag(bitSprite5|bitSprite6) == 0 {
		vic.ClearBALow()
	}
}

// phaseSprite5DMAPhase2: Sprite 5 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 7 DMA if it is enabled.
//
//go:nosplit
func (seq *Sequencer) phaseSprite5DMAPhase2(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	if vic.sprites.GetDMAFlag(bitSprite5) != 0 {
		vic.sprites.FetchPhase2(5)
	} else {
		vic.memory.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite7) != 0 {
		vic.SetBALow()
	}
}

// phaseSprite6DMAPhase1: Sprite 6 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 6 and 7.
//
//go:nosplit
func (seq *Sequencer) phaseSprite6DMAPhase1(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	if vic.sprites.GetDMAFlag(bitSprite6) != 0 {
		vic.sprites.FetchPhase1(6)
	}
	if vic.sprites.GetDMAFlag(bitSprite6|bitSprite7) == 0 {
		vic.ClearBALow()
	}
}

// phaseSprite6DMAData1Data2: Sprite 6 DMA continues, fetching its second and third data bytes.
//
//go:nosplit
func (seq *Sequencer) phaseSprite6DMAPhase2(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	if vic.sprites.GetDMAFlag(bitSprite6) != 0 {
		vic.sprites.FetchPhase2(6)
	} else {
		vic.memory.AccessIdle()
	}
}

// phaseSprite7DMAPhase1: Sprite 7 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprite 7.
//
//go:nosplit
func (seq *Sequencer) phaseSprite7DMAPhase1(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	if vic.sprites.GetDMAFlag(bitSprite7) != 0 {
		vic.sprites.FetchPhase1(7)
	}
	if vic.sprites.GetDMAFlag(bitSprite7) == 0 {
		vic.ClearBALow()
	}
}

// phaseSprite7DMAPhase2: Sprite 7 DMA continues, fetching its second and third data bytes. This concludes the main
// sprite DMA phase for the upcoming line.
//
//go:nosplit
func (seq *Sequencer) phaseSprite7DMAPhase2(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	if vic.sprites.GetDMAFlag(bitSprite7) != 0 {
		vic.sprites.FetchPhase2(7)
	} else {
		vic.memory.AccessIdle()
	}
}

// phaseRefresh: This is a refresh cycle. The VIC performs a DRAM refresh operation by accessing an address
// in the range $3C00-$3FFF. The address bus is released for the CPU, and the BA signal is cleared.
//
//go:nosplit
func (seq *Sequencer) phaseRefresh(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	vic.ClearBALow()
	vic.memory.AccessRefresh()
}

// phaseSetupBadLineCheck: This is a refresh cycle. The VIC checks for a "bad line" condition, which occurs if the
// DEN bit is set and the lower 3 bits of the raster counter match the lower 3 bits of the YSCROLL register.
// If it is a bad line, the VIC prepares to halt the CPU by asserting the BA signal.
//
//go:nosplit
func (seq *Sequencer) phaseSetupBadLineCheck(vic *VIC) {
	vic.memory.AccessRefresh()
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	vic.TryBALowIfBadLine()
}

// phaseSetupRasterXReset: This is a refresh cycle. The horizontal raster counter (rasterX) is reset to 0. The VIC is now
// in the left border area. The BA signal is asserted if the bad line condition was met in the previous cycle.
//
//go:nosplit
func (seq *Sequencer) phaseSetupRasterXReset(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(seq.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.memory.AccessRefresh()
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	vic.TryBALowIfBadLine()
	vic.ResetRasterX()
}

// phaseSetupVCounterLoad: This is a refresh cycle. The Video Counter (VC) is loaded from the Video Counter Base (VCBASE),
// pointing to the current character row in screen memory. The Row Counter (RC) is reset to 0 if this is the
// first scanline of a character row (i.e., sequencerRasterY & 7 == 0).
//
//go:nosplit
func (seq *Sequencer) phaseSetupVCounterLoad(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(seq.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.memory.AccessRefresh()
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	vic.graphics.TryResetRowCounter(vic.badLineCondition)
	vic.TryBALowIfBadLine()
	vic.graphics.UpdateVideoCounter()
}

// phaseSetupRCounterCheckAndSpritePipe1: This is the critical "bad line" decision cycle. If it's a bad line, the VIC takes full control
// of the bus for the next 40 cycles. The graphics pipeline begins its first access, fetching a character code
// from screen RAM using the Video Counter (VC). Sprite Y-expansion counters for the *next* scanline are checked.
//
//go:nosplit
func (seq *Sequencer) phaseSetupRCounterCheckAndSpritePipe1(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(seq.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.memory.AccessRefresh()
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	vic.sprites.TryIncrementCounterBase()
	vic.graphics.ResetLineIndex()
	vic.TryBALowIfBadLine()
	vic.graphics.TryPhi2Access(vic.baLow, vic.aecLow)
}

// phaseDisplayFirstFetchAndSpritePipe2: First graphics data fetch cycle. The VIC fetches the character's bitmap data from Character
// ROM or RAM, using the character code fetched in the previous cycle and the current Row Counter (RC).
// Sprite Y-expansion counters are committed.
//
//go:nosplit
func (seq *Sequencer) phaseDisplayFirstFetchAndSpritePipe2(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(seq.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.graphics.TryGraphicsAccess(vic.sequencerRasterY)
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	vic.sprites.TryIncrementCounterBase()
	vic.sprites.CommitIncrementCounterBase()
	vic.TryBALowIfBadLine()
	vic.graphics.TryPhi2Access(vic.baLow, vic.aecLow)
}

// phaseDisplayMainFetchC40: Second graphics data fetch cycle. The VIC fetches the next character code from screen RAM.
// The first pixels of the 40-column window (or the side border in 38-column mode) are drawn. The side
// border logic for 40-column mode is triggered.
//
//go:nosplit
func (seq *Sequencer) phaseDisplayMainFetchC40(vic *VIC) {
	vic.borders.Column40Update(vic.sequencerRasterY, vic.denBit)
	if vic.drawLine {
		vic.borders.AcquireColor(seq.curr.cycleBorder)
		if vic.borders.VerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
	}
	vic.graphics.TryGraphicsAccess(vic.sequencerRasterY)
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	vic.TryBALowIfBadLine()
	vic.graphics.TryPhi2Access(vic.baLow, vic.aecLow)
}

// phaseDisplayMainFetchC38: The main display window begins. The VIC fetches the next character's bitmap data. The fetched
// bitmap data is loaded into the graphics pipeline's internal shift register. The side border logic for
// 38-column mode is triggered. Pixels are drawn.
//
//go:nosplit
func (seq *Sequencer) phaseDisplayMainFetchC38(vic *VIC) {
	vic.borders.Column38Update(vic.sequencerRasterY, vic.denBit)
	if vic.drawLine {
		vic.borders.AcquireColor(seq.curr.cycleBorder)
		if vic.borders.VerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
	}
	vic.graphics.TryGraphicsAccess(vic.sequencerRasterY)
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	vic.TryBALowIfBadLine()
	vic.graphics.TryPhi2Access(vic.baLow, vic.aecLow)
	vic.graphics.UpdateCharDataLast()
}

// phaseDisplayMainFetch: These 36 cycles form the core of the visible display area. In each pair of cycles, the VIC
// fetches a character code from screen RAM and its corresponding bitmap data. The graphics pipeline
// continuously shifts out 8 pixels per character, drawing either foreground or background pixels. On a "bad line",
// the CPU is halted throughout this entire phase.
//
//go:nosplit
func (seq *Sequencer) phaseDisplayMainFetch(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(seq.curr.cycleBorder)
		if vic.borders.VerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
	}
	vic.graphics.TryGraphicsAccess(vic.sequencerRasterY)
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	vic.TryBALowIfBadLine()
	vic.graphics.TryPhi2Access(vic.baLow, vic.aecLow)
	vic.graphics.UpdateCharDataLast()
}

// phaseDMASetupAndTeardownLastFetch: This is the last graphics data fetch cycle for the line. The VIC finalizes the graphics pipeline
// and prepares for the *next* scanline by updating sprite Y-expansion flags and calculating which sprites will
// be active (setting up their DMA flags). The BA signal is asserted to prepare for sprite 0 DMA.
//
//go:nosplit
func (seq *Sequencer) phaseDMASetupAndTeardownLastFetch(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(seq.curr.cycleBorder)
		if vic.borders.VerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
	}
	vic.graphics.TryGraphicsAccess(vic.sequencerRasterY)
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	vic.sprites.UpdateYExpansion()
	vic.sprites.UpdateDMA(vic.sequencerRasterY)
	if vic.sprites.GetDMAFlag(bitSprite0) != 0 {
		vic.SetBALow()
	} else {
		vic.ClearBALow()
	}
}

// phaseTeardownIdle: Idle cycle. The main graphics fetch is complete. The VIC is now in the right border area.
// The 38-column side border logic is applied. The BA signal is asserted for sprite 0 DMA if needed.
//
//go:nosplit
func (seq *Sequencer) phaseTeardownIdle(vic *VIC) {
	vic.borders.Column38Apply()
	if vic.drawLine {
		vic.borders.AcquireColor(seq.curr.cycleBorder)
		if vic.borders.VerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
	}
	vic.memory.AccessIdle()
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	vic.sprites.UpdateDMA(vic.sequencerRasterY)
	if vic.sprites.GetDMAFlag(bitSprite0) != 0 {
		vic.SetBALow()
	}
}

// phaseTeardownCommitSpriteFlags: Idle cycle. The 40-column side border logic is applied. The sprite DMA flags calculated in
// cycle 55 are now committed for use in the upcoming DMA phase. The BA signal is asserted to prepare for
// sprite 1 DMA.
//
//go:nosplit
func (seq *Sequencer) phaseTeardownCommitSpriteFlags(vic *VIC) {
	vic.borders.Column40Apply()
	vic.sprites.CommitSpriteFlags()
	if vic.drawLine {
		vic.borders.AcquireColor(seq.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.memory.AccessIdle()
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	if vic.sprites.GetDMAFlag(bitSprite1) != 0 {
		vic.SetBALow()
	}
}

// phaseSprite0DMAPhase1: Sprite 0 DMA for the upcoming line begins if enabled, fetching its pointer and first data byte.
// Sprite flags are prepared for the line *after* the next one.
//
//go:nosplit
func (seq *Sequencer) phaseSprite0DMAPhase1(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(seq.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.sprites.PrepareSpriteFlags(vic.sequencerRasterY)
	if vic.sprites.GetDMAFlag(bitSprite0) != 0 {
		vic.sprites.FetchPhase1(0)
	}
	vic.graphics.UpdateDisplayAccess(vic.badLineCondition)
}

// phaseSprite0DMAPhase2: Sprite 0 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 2 DMA.
//
//go:nosplit
func (seq *Sequencer) phaseSprite0DMAPhase2(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(seq.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	if vic.sprites.GetDMAFlag(bitSprite0) != 0 {
		vic.sprites.FetchPhase2(0)
	} else {
		vic.memory.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite2) != 0 {
		vic.SetBALow()
	}
}

// phaseSprite1DMAPhase1: Sprite 1 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 1 and 2.
//
//go:nosplit
func (seq *Sequencer) phaseSprite1DMAPhase1(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(seq.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	//if vic.drawLine {
	//	vic.graphics.DrawBackground()
	//	vic.sprites.Draw(seq.spriteOdd)
	//	vic.borders.Draw()
	//	vic.lineStart += seq.width
	//}
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	if vic.sprites.GetDMAFlag(bitSprite1) != 0 {
		vic.sprites.FetchPhase1(1)
	} else {
		vic.memory.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite1|bitSprite2) == 0 {
		vic.ClearBALow()
	}
}

// phaseSprite1DMAPhase2: Sprite 1 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 3 DMA.
//
//go:nosplit
func (seq *Sequencer) phaseSprite1DMAPhase2(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	if vic.sprites.GetDMAFlag(bitSprite1) != 0 {
		vic.sprites.FetchPhase2(1)
	} else {
		vic.memory.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite3) != 0 {
		vic.SetBALow()
	}
}

// phaseSprite2DMAPhase1: Sprite 2 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 2 and 3.
//
//go:nosplit
func (seq *Sequencer) phaseSprite2DMAPhase1(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	if vic.sprites.GetDMAFlag(bitSprite2) != 0 {
		vic.sprites.FetchPhase1(2)
	}
	if vic.sprites.GetDMAFlag(bitSprite2|bitSprite3) == 0 {
		vic.ClearBALow()
	}
}

// phaseSprite1DMAPhase2: Sprite 2 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 4 DMA.
//
//go:nosplit
func (seq *Sequencer) phaseSprite2DMAPhase2AndTeardownFinal(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess(vic.badLineCondition)
	vic.borders.UpdateVerticalFlipFlop(vic.sequencerRasterY, vic.denBit)
	if vic.sprites.GetDMAFlag(bitSprite2) != 0 {
		vic.sprites.FetchPhase2(2)
	} else {
		vic.memory.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite4) != 0 {
		vic.SetBALow()
	}

	if vic.drawLine {
		vic.sprites.Draw()
		vic.borders.Draw()
		vic.lineStart += seq.width
	}

	vic.socketLastCycle()
}
