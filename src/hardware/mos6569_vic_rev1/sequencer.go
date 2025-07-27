package mos6569

const drawLoopCycles = 36 // drawLoopCycles defines the number of iterations for the main display loop in video cycle data processing.

// SequencerData represents a chainable unit in a sequencer containing a function, a next node, and cycle information.
type SequencerData struct {
	fn          func()
	next        *SequencerData
	cycle       uint8
	cycleBorder uint8
}

// NewSequencerData creates and returns a new instance of SequencerData with the provided function initialized.
func NewSequencerData(fn func()) *SequencerData {
	return &SequencerData{fn: fn}
}

// Sequencer represents a structure to manage video cycle sequencing and raster timing for display rendering.
type Sequencer struct {
	lineWidth          int
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
	//curr               *SequencerData
}

// NewSequencerPal initializes the PAL video timing cycle data.
// It constructs a circular linked list of 63 cycleData nodes,
// where each node represents one CPU clock cycle of a single PAL scanline. It pre-calculates border-related
// values for each cycle and links them in sequence to form the complete 63-cycle sequencer.
func NewSequencerPal(vic *VIC) *Sequencer {
	const palWidth = 384
	const palHeight = 272
	const palTotalRasters = 312

	seq := &Sequencer{
		lineWidth:          palWidth,
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

	seq.add(vic.phaseSprite3DMAPhase1AndInit)
	seq.add(vic.phaseSprite3DMAPhase2AndVBlank)
	seq.add(vic.phaseSprite4DMAPhase1)
	seq.add(vic.phaseSprite4DMAPhase2)
	seq.add(vic.phaseSprite5DMAPhase1)
	seq.add(vic.phaseSprite5DMAPhase2)
	seq.add(vic.phaseSprite6DMAPhase1)
	seq.add(vic.phaseSprite6DMAPhase2)
	seq.add(vic.phaseSprite7DMAPhase1)
	seq.add(vic.phaseSprite7DMAPhase2)
	seq.add(vic.phaseRefresh)
	seq.add(vic.phaseSetupBadLineCheck)
	seq.add(vic.phaseSetupRasterXReset)
	seq.add(vic.phaseSetupVCounterLoad)
	seq.add(vic.phaseSetupRCounterCheckAndSpritePipe1)
	seq.add(vic.phaseDisplayFirstFetchAndSpritePipe2)
	seq.add(vic.phaseDisplayMainFetchC40)
	seq.add(vic.phaseDisplayMainFetchC38)
	for i := 0; i < drawLoopCycles; i++ {
		seq.add(vic.phaseDisplayMainFetch)
	}
	seq.add(vic.phaseDMASetupAndTeardownLastFetch)
	seq.add(vic.phaseTeardownIdle)
	seq.add(vic.phaseTeardownCommitSpriteFlags)
	seq.add(vic.phaseSprite0DMAPhase1AndScanLineEnd)
	seq.add(vic.phaseSprite0DMAPhase2)
	seq.add(vic.phaseSprite1DMAPhase1)
	seq.add(vic.phaseSprite1DMAPhase2)
	seq.add(vic.phaseSprite2DMAPhase1)
	seq.add(vic.phaseSprite2DMAPhase2AndTeardownFinal)

	seq.finalize()

	return seq
}

// NewSequencerNtsc initializes the NTSC video timing cycle data with 63 cycleData nodes for a single scanline.
// Constructs a sequencer containing pre-calculated raster and border values for 263 total raster lines.
// Links the nodes sequentially to form a complete NTSC video frame.
func NewSequencerNtsc(vic *VIC) *Sequencer {
	const ntscWidth = 384
	const ntscHeight = 240
	const ntscTotalRasters = 263

	seq := &Sequencer{
		lineWidth:          ntscWidth,
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

	seq.add(vic.phaseSprite3DMAPhase1AndInit)
	seq.add(vic.phaseSprite3DMAPhase2AndVBlank)
	seq.add(vic.phaseSprite4DMAPhase1)
	seq.add(vic.phaseSprite4DMAPhase2)
	seq.add(vic.phaseSprite5DMAPhase1)
	seq.add(vic.phaseSprite5DMAPhase2)
	seq.add(vic.phaseSprite6DMAPhase1)
	seq.add(vic.phaseSprite6DMAPhase2)
	seq.add(vic.phaseSprite7DMAPhase1)
	seq.add(vic.phaseSprite7DMAPhase2)
	seq.add(vic.phaseRefresh)
	seq.add(vic.phaseSetupBadLineCheck)
	seq.add(vic.phaseSetupRasterXReset)
	seq.add(vic.phaseSetupVCounterLoad)
	seq.add(vic.phaseSetupRCounterCheckAndSpritePipe1)
	seq.add(vic.phaseDisplayFirstFetchAndSpritePipe2)
	seq.add(vic.phaseDisplayMainFetchC40)
	seq.add(vic.phaseDisplayMainFetchC38)
	for i := 0; i < drawLoopCycles; i++ {
		seq.add(vic.phaseDisplayMainFetch)
	}
	seq.add(vic.phaseDMASetupAndTeardownLastFetch)
	// The two extra NTSC cycles are refresh cycles.
	seq.add(vic.phaseRefresh)
	seq.add(vic.phaseRefresh)
	seq.add(vic.phaseTeardownIdle)
	seq.add(vic.phaseTeardownCommitSpriteFlags)
	seq.add(vic.phaseSprite0DMAPhase1AndScanLineEnd)
	seq.add(vic.phaseSprite0DMAPhase2)
	seq.add(vic.phaseSprite1DMAPhase1)
	seq.add(vic.phaseSprite1DMAPhase2)
	seq.add(vic.phaseSprite2DMAPhase1)
	seq.add(vic.phaseSprite2DMAPhase2AndTeardownFinal)

	seq.finalize()

	return seq
}

func (seq *Sequencer) Start() *SequencerData {
	return seq.data[0]
}

// GetRasterYMax returns the maximum Y-coordinate value of the raster as uint16.
func (seq *Sequencer) GetRasterYMax() uint16 {
	return seq.rasterYMax
}

// GetLineWidth returns the width of the Sequencer instance.
func (seq *Sequencer) GetLineWidth() int {
	return seq.lineWidth
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

// GetFirstDisplayedLine retrieves the first raster line that is displayed on the screen from the Sequencer instance.
func (seq *Sequencer) GetFirstDisplayedLine() uint16 {
	return seq.firstDisplayedLine
}

// GetLastDisplayedLine retrieves the last displayed raster line value from the Sequencer instance.
func (seq *Sequencer) GetLastDisplayedLine() uint16 {
	return seq.lastDisplayedLine
}

// prepare initializes the Sequencer by creating an empty slice of SequencerData pointers to be used for cycle management.
func (seq *Sequencer) prepare() {
	seq.row24YStart = seq.row25YStart + 4
	seq.row24YStop = seq.row25YStop - 4
	seq.data = make([]*SequencerData, 0)
}

// add appends a new SequencerData instance created with the provided function to the Sequencer's data slice.
func (seq *Sequencer) add(fn func()) {
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
}

// phaseSprite3DMAPhase1AndInit: This cycle marks the beginning of the horizontal blanking period. The raster line counter (sequencerRasterY)
// is checked against the maximum value. If it matches, a V-Blank is scheduled for the next cycle. Otherwise,
// sequencerRasterY is incremented for the new scanline. Sprite 3 DMA for the upcoming line begins if enabled, fetching
// the sprite pointer (phi1) and the first byte of sprite data (phi2).
//
//go:nosplit
func (vic *VIC) phaseSprite3DMAPhase1AndInit() {
	vic.sprites.Prepare()
	if vic.rasterY >= vic.rasterYMax {
		vic.vBlankNextCycle = true
	} else {
		vic.rasterY++
		vic.interrupts.VerifyRasterY(vic.rasterY)
		vic.graphics.BadLineVerify(vic.rasterY, vic.borders.GetDen())
		vic.drawLine = (vic.rasterY >= vic.firstDisplayedLine) && (vic.rasterY <= vic.lastDisplayedLine)
	}
	vic.borders.ColumnInitialize()
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.sprites.GetDMAFlag(bitSprite3) != 0 {
		vic.sprites.ActivateSprite(3)
		vic.sprites.LatchAttributes()
		vic.sprites.ReadPtr(vic.memory.GetMatrixBase()) //ph1
		vic.sprites.ReadData(0)                         //ph2
	}
	if vic.sprites.GetDMAFlag(bitSprite3|bitSprite4) == 0 {
		vic.ClearBALow()
	}
	vic.rasterX += 8
}

// phaseSprite3DMAPhase2AndVBlank: The V-Blank is triggered if scheduled in the previous cycle, resetting raster counters and firing
// the V-Blank interrupt. Internal pointers for graphics, sprites, and borders are reset. Sprite 3 DMA
// continues, fetching the second and third bytes of sprite data. The BA (Bus-Available) signal is asserted
// (pulled low) to prepare for sprite 5 DMA if it is enabled for the upcoming line.
//
//go:nosplit
func (vic *VIC) phaseSprite3DMAPhase2AndVBlank() {
	if vic.vBlankNextCycle {
		vic.vBlankNextCycle = false
		vic.beam.ResetLineOffset()
		vic.graphics.ResetVideoCounterLatch()
		vic.memory.ResetRefreshCounter()
		vic.rasterY = 0
		vic.interrupts.VerifyRasterY(vic.rasterY)
		vic.lightPen.TriggerClear()
		vic.socketVBlank()
	}
	vic.collisions.ClearGraphics()
	vic.graphics.ResetOffset()
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.sprites.GetDMAFlag(bitSprite3) != 0 {
		vic.sprites.ActivateSprite(3)
		vic.sprites.ReadData(1) //ph1
		vic.sprites.ReadData(2) //ph1
	} else {
		vic.memory.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite5) != 0 {
		vic.SetBALow()
	}
	vic.rasterX += 8
}

// phaseSprite4DMAPhase1: Sprite 4 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 4 and 5.
//
//go:nosplit
func (vic *VIC) phaseSprite4DMAPhase1() {
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.sprites.GetDMAFlag(bitSprite4) != 0 {
		vic.sprites.ActivateSprite(4)
		vic.sprites.LatchAttributes()
		vic.sprites.ReadPtr(vic.memory.GetMatrixBase()) //ph1
		vic.sprites.ReadData(0)                         //ph2
	}
	if vic.sprites.GetDMAFlag(bitSprite4|bitSprite5) == 0 {
		vic.ClearBALow()
	}
	vic.rasterX += 8
}

// phaseSprite4DMAPhase2: Sprite 4 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 6 DMA if it is enabled.
//
//go:nosplit
func (vic *VIC) phaseSprite4DMAPhase2() {
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.sprites.GetDMAFlag(bitSprite4) != 0 {
		vic.sprites.ActivateSprite(4)
		vic.sprites.ReadData(1) //ph1
		vic.sprites.ReadData(2) //ph1
	} else {
		vic.memory.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite6) != 0 {
		vic.SetBALow()
	}
	vic.rasterX += 8
}

// phaseSprite5DMAPhase1: Sprite 5 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 5 and 6.
//
//go:nosplit
func (vic *VIC) phaseSprite5DMAPhase1() {
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.sprites.GetDMAFlag(bitSprite5) != 0 {
		vic.sprites.ActivateSprite(5)
		vic.sprites.LatchAttributes()
		vic.sprites.ReadPtr(vic.memory.GetMatrixBase()) //ph1
		vic.sprites.ReadData(0)                         //ph2
	}
	if vic.sprites.GetDMAFlag(bitSprite5|bitSprite6) == 0 {
		vic.ClearBALow()
	}
	vic.rasterX += 8
}

// phaseSprite5DMAPhase2: Sprite 5 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 7 DMA if it is enabled.
//
//go:nosplit
func (vic *VIC) phaseSprite5DMAPhase2() {
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.sprites.GetDMAFlag(bitSprite5) != 0 {
		vic.sprites.ActivateSprite(5)
		vic.sprites.ReadData(1) //ph1
		vic.sprites.ReadData(2) //ph1
	} else {
		vic.memory.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite7) != 0 {
		vic.SetBALow()
	}
	vic.rasterX += 8
}

// phaseSprite6DMAPhase1: Sprite 6 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 6 and 7.
//
//go:nosplit
func (vic *VIC) phaseSprite6DMAPhase1() {
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.sprites.GetDMAFlag(bitSprite6) != 0 {
		vic.sprites.ActivateSprite(6)
		vic.sprites.LatchAttributes()
		vic.sprites.ReadPtr(vic.memory.GetMatrixBase()) //ph1
		vic.sprites.ReadData(0)                         //ph2
	}
	if vic.sprites.GetDMAFlag(bitSprite6|bitSprite7) == 0 {
		vic.ClearBALow()
	}
	vic.rasterX += 8
}

// phaseSprite6DMAData1Data2: Sprite 6 DMA continues, fetching its second and third data bytes.
//
//go:nosplit
func (vic *VIC) phaseSprite6DMAPhase2() {
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.sprites.GetDMAFlag(bitSprite6) != 0 {
		vic.sprites.ActivateSprite(6)
		vic.sprites.ReadData(1) //ph1
		vic.sprites.ReadData(2) //ph1
	} else {
		vic.memory.AccessIdle()
	}
	vic.rasterX += 8
}

// phaseSprite7DMAPhase1: Sprite 7 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprite 7.
//
//go:nosplit
func (vic *VIC) phaseSprite7DMAPhase1() {
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.sprites.GetDMAFlag(bitSprite7) != 0 {
		vic.sprites.ActivateSprite(7)
		vic.sprites.LatchAttributes()
		vic.sprites.ReadPtr(vic.memory.GetMatrixBase()) //ph1
		vic.sprites.ReadData(0)                         //ph2
	}
	if vic.sprites.GetDMAFlag(bitSprite7) == 0 {
		vic.ClearBALow()
	}
	vic.rasterX += 8
}

// phaseSprite7DMAPhase2: Sprite 7 DMA continues, fetching its second and third data bytes. This concludes the main
// sprite DMA phase for the upcoming line.
//
//go:nosplit
func (vic *VIC) phaseSprite7DMAPhase2() {
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.sprites.GetDMAFlag(bitSprite7) != 0 {
		vic.sprites.ActivateSprite(7)
		vic.sprites.ReadData(1) //ph1
		vic.sprites.ReadData(2) //ph1
	} else {
		vic.memory.AccessIdle()
	}
	vic.rasterX += 8
}

// phaseRefresh: This is a refresh cycle. The VIC performs a DRAM refresh operation by accessing an address
// in the range $3C00-$3FFF. The address bus is released for the CPU, and the BA signal is cleared.
//
//go:nosplit
func (vic *VIC) phaseRefresh() {
	vic.graphics.AcquireDisplayAccessIfBadLine()
	vic.ClearBALow()
	vic.memory.AccessRefresh()
	vic.rasterX += 8
}

// phaseSetupBadLineCheck: This is a refresh cycle. The VIC checks for a "bad line" condition, which occurs if the
// DEN bit is set and the lower 3 bits of the raster counter match the lower 3 bits of the YSCROLL register.
// If it is a bad line, the VIC prepares to halt the CPU by asserting the BA signal.
//
//go:nosplit
func (vic *VIC) phaseSetupBadLineCheck() {
	vic.memory.AccessRefresh()
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.graphics.BadLine() {
		vic.SetBALow()
	}
	vic.rasterX += 8
}

// phaseSetupRasterXReset: This is a refresh cycle. The horizontal raster counter (rasterX) is reset to 0. The VIC is now
// in the left border area. The BA signal is asserted if the bad line condition was met in the previous cycle.
//
//go:nosplit
func (vic *VIC) phaseSetupRasterXReset() {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.memory.AccessRefresh()
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.graphics.BadLine() {
		vic.SetBALow()
	}
	vic.rasterX = 0xfffc
	vic.rasterX += 8
}

// phaseSetupVCounterLoad: This is a refresh cycle. The Video Counter (VC) is loaded from the Video Counter Base (VCBASE),
// pointing to the current character row in screen memory. The Row Counter (RC) is reset to 0 if this is the
// first scanline of a character row (i.e., sequencerRasterY & 7 == 0).
//
//go:nosplit
func (vic *VIC) phaseSetupVCounterLoad() {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.memory.AccessRefresh()
	vic.graphics.AcquireDisplayAccessIfBadLine()
	vic.graphics.ResetRowCounterIfBadLine()
	if vic.graphics.BadLine() {
		vic.SetBALow()
	}
	vic.graphics.UpdateVideoCounter()
	vic.rasterX += 8
}

// phaseSetupRCounterCheckAndSpritePipe1: This is the critical "bad line" decision cycle. If it's a bad line, the VIC takes full control
// of the bus for the next 40 cycles. The graphics pipeline begins its first access, fetching a character code
// from screen RAM using the Video Counter (VC). Sprite Y-expansion counters for the *next* scanline are checked.
//
//go:nosplit
func (vic *VIC) phaseSetupRCounterCheckAndSpritePipe1() {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.memory.AccessRefresh()
	vic.graphics.AcquireDisplayAccessIfBadLine()
	vic.sprites.TryIncrementCounterBase()
	vic.graphics.ResetLineIndex()
	if vic.graphics.BadLine() {
		vic.SetBALow()
	}
	if vic.baLow { //phi2
		if vic.aecLow {
			vic.graphics.FetchData(vic.memory.GetMatrixBase())
			vic.graphics.FetchColor(vic.memory.GetMatrixBase())
		} else {
			vic.graphics.FetchDataFake()
			vic.graphics.FetchColorFake()
		}
	}
	vic.rasterX += 8
}

// phaseDisplayFirstFetchAndSpritePipe2: First graphics data fetch cycle. The VIC fetches the character's bitmap data from Character
// ROM or RAM, using the character code fetched in the previous cycle and the current Row Counter (RC).
// Sprite Y-expansion counters are committed.
//
//go:nosplit
func (vic *VIC) phaseDisplayFirstFetchAndSpritePipe2() {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.graphics.FetchMemory(vic.rasterY)
	vic.graphics.AcquireDisplayAccessIfBadLine()
	vic.sprites.TryIncrementCounterBase()
	vic.sprites.CommitIncrementCounterBase()
	if vic.graphics.BadLine() {
		vic.SetBALow()
	}
	if vic.baLow { //phi2
		if vic.aecLow {
			vic.graphics.FetchData(vic.memory.GetMatrixBase())
			vic.graphics.FetchColor(vic.memory.GetMatrixBase())
		} else {
			vic.graphics.FetchDataFake()
			vic.graphics.FetchColorFake()
		}
	}
	vic.rasterX += 8
}

// phaseDisplayMainFetchC40: Second graphics data fetch cycle. The VIC fetches the next character code from screen RAM.
// The first pixels of the 40-column window (or the side border in 38-column mode) are drawn. The side
// border logic for 40-column mode is triggered.
//
//go:nosplit
func (vic *VIC) phaseDisplayMainFetchC40() {
	vic.borders.Column40Update(vic.rasterY)

	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		if vic.borders.VerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
	}
	vic.graphics.FetchMemory(vic.rasterY)
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.graphics.BadLine() {
		vic.SetBALow()
	}
	if vic.baLow { //phi2
		if vic.aecLow {
			vic.graphics.FetchData(vic.memory.GetMatrixBase())
			vic.graphics.FetchColor(vic.memory.GetMatrixBase())
		} else {
			vic.graphics.FetchDataFake()
			vic.graphics.FetchColorFake()
		}
	}
	vic.rasterX += 8
}

// phaseDisplayMainFetchC38: The main display window begins. The VIC fetches the next character's bitmap data. The fetched
// bitmap data is loaded into the graphics pipeline's internal shift register. The side border logic for
// 38-column mode is triggered. Pixels are drawn.
//
//go:nosplit
func (vic *VIC) phaseDisplayMainFetchC38() {
	vic.borders.Column38Update(vic.rasterY)

	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		if vic.borders.VerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
	}
	vic.graphics.FetchMemory(vic.rasterY)
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.graphics.BadLine() {
		vic.SetBALow()
	}
	if vic.baLow { //phi2
		if vic.aecLow {
			vic.graphics.FetchData(vic.memory.GetMatrixBase())
			vic.graphics.FetchColor(vic.memory.GetMatrixBase())
		} else {
			vic.graphics.FetchDataFake()
			vic.graphics.FetchColorFake()
		}
	}
	vic.graphics.CommitCharData()
	vic.rasterX += 8
}

// phaseDisplayMainFetch: These 36 cycles form the core of the visible display area. In each pair of cycles, the VIC
// fetches a character code from screen RAM and its corresponding bitmap data. The graphics pipeline
// continuously shifts out 8 pixels per character, drawing either foreground or background pixels. On a "bad line",
// the CPU is halted throughout this entire phase.
//
//go:nosplit
func (vic *VIC) phaseDisplayMainFetch() {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		if vic.borders.VerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
	}
	// Pipelining
	vic.graphics.FetchMemory(vic.rasterY) //phi1
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.graphics.BadLine() {
		vic.SetBALow()
	}
	if vic.baLow { //phi2
		if vic.aecLow {
			vic.graphics.FetchData(vic.memory.GetMatrixBase())
			vic.graphics.FetchColor(vic.memory.GetMatrixBase())
		} else {
			vic.graphics.FetchDataFake()
			vic.graphics.FetchColorFake()
		}
	}
	vic.graphics.CommitCharData()
	vic.rasterX += 8
}

// phaseDMASetupAndTeardownLastFetch: This is the last graphics data fetch cycle for the line. The VIC finalizes the graphics pipeline
// and prepares for the *next* scanline by updating sprite Y-expansion flags and calculating which sprites will
// be active (setting up their DMA flags). The BA signal is asserted to prepare for sprite 0 DMA.
//
//go:nosplit
func (vic *VIC) phaseDMASetupAndTeardownLastFetch() {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		if vic.borders.VerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
	}
	vic.graphics.FetchMemory(vic.rasterY)
	vic.graphics.AcquireDisplayAccessIfBadLine()
	vic.sprites.UpdateYExpansion()
	vic.sprites.UpdateDMA(vic.rasterY)
	if vic.sprites.GetDMAFlag(bitSprite0) != 0 {
		vic.SetBALow()
	} else {
		vic.ClearBALow()
	}
	vic.rasterX += 8
}

// phaseTeardownIdle: Idle cycle. The main graphics fetch is complete. The VIC is now in the right border area.
// The 38-column side border logic is applied. The BA signal is asserted for sprite 0 DMA if needed.
//
//go:nosplit
func (vic *VIC) phaseTeardownIdle() {
	vic.borders.Column38Apply()
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		if vic.borders.VerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
	}
	vic.memory.AccessIdle()
	vic.graphics.AcquireDisplayAccessIfBadLine()
	vic.sprites.UpdateDMA(vic.rasterY)
	if vic.sprites.GetDMAFlag(bitSprite0) != 0 {
		vic.SetBALow()
	}
	vic.rasterX += 8
}

// phaseTeardownCommitSpriteFlags: Idle cycle. The 40-column side border logic is applied. The sprite DMA flags calculated in
// cycle 55 are now committed for use in the upcoming DMA phase. The BA signal is asserted to prepare for
// sprite 1 DMA.
//
//go:nosplit
func (vic *VIC) phaseTeardownCommitSpriteFlags() {
	vic.borders.Column40Apply()
	vic.sprites.CommitSpriteFlags()
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.memory.AccessIdle()
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.sprites.GetDMAFlag(bitSprite1) != 0 {
		vic.SetBALow()
	}
	vic.rasterX += 8
}

// phaseSprite0DMAPhase1: Sprite 0 DMA for the upcoming line begins if enabled, fetching its pointer and first data byte.
// Sprite flags are prepared for the line *after* the next one.
//
//go:nosplit
func (vic *VIC) phaseSprite0DMAPhase1AndScanLineEnd() {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.sprites.PrepareSpriteFlags(vic.rasterY)
	if vic.sprites.GetDMAFlag(bitSprite0) != 0 {
		vic.sprites.ActivateSprite(0)
		vic.sprites.LatchAttributes()
		vic.sprites.ReadPtr(vic.memory.GetMatrixBase()) //ph1
		vic.sprites.ReadData(0)                         //ph2
	}
	vic.graphics.TryAcquireDisplayAccessOnScanlineEnd()
	vic.rasterX += 8
}

// phaseSprite0DMAPhase2: Sprite 0 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 2 DMA.
//
//go:nosplit
func (vic *VIC) phaseSprite0DMAPhase2() {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.sprites.GetDMAFlag(bitSprite0) != 0 {
		vic.sprites.ActivateSprite(0)
		vic.sprites.ReadData(1) //ph1
		vic.sprites.ReadData(2) //ph1
	} else {
		vic.memory.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite2) != 0 {
		vic.SetBALow()
	}
	vic.rasterX += 8
}

// phaseSprite1DMAPhase1: Sprite 1 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 1 and 2.
//
//go:nosplit
func (vic *VIC) phaseSprite1DMAPhase1() {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.sprites.GetDMAFlag(bitSprite1) != 0 {
		vic.sprites.ActivateSprite(1)
		vic.sprites.LatchAttributes()
		vic.sprites.ReadPtr(vic.memory.GetMatrixBase()) //ph1
		vic.sprites.ReadData(0)                         //ph2
	} else {
		vic.memory.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite1|bitSprite2) == 0 {
		vic.ClearBALow()
	}
	vic.rasterX += 8
}

// phaseSprite1DMAPhase2: Sprite 1 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 3 DMA.
//
//go:nosplit
func (vic *VIC) phaseSprite1DMAPhase2() {
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.sprites.GetDMAFlag(bitSprite1) != 0 {
		vic.sprites.ActivateSprite(1)
		vic.sprites.ReadData(1) //ph1
		vic.sprites.ReadData(2) //ph1
	} else {
		vic.memory.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite3) != 0 {
		vic.SetBALow()
	}
	vic.rasterX += 8
}

// phaseSprite2DMAPhase1: Sprite 2 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 2 and 3.
//
//go:nosplit
func (vic *VIC) phaseSprite2DMAPhase1() {
	vic.graphics.AcquireDisplayAccessIfBadLine()
	if vic.sprites.GetDMAFlag(bitSprite2) != 0 {
		vic.sprites.ActivateSprite(2)
		vic.sprites.LatchAttributes()
		vic.sprites.ReadPtr(vic.memory.GetMatrixBase()) //ph1
		vic.sprites.ReadData(0)                         //ph2
	}
	if vic.sprites.GetDMAFlag(bitSprite2|bitSprite3) == 0 {
		vic.ClearBALow()
	}
	vic.rasterX += 8
}

// phaseSprite1DMAPhase2: Sprite 2 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 4 DMA.
//
//go:nosplit
func (vic *VIC) phaseSprite2DMAPhase2AndTeardownFinal() {
	vic.graphics.AcquireDisplayAccessIfBadLine()
	vic.borders.UpdateVerticalFlipFlop(vic.rasterY)
	if vic.sprites.GetDMAFlag(bitSprite2) != 0 {
		vic.sprites.ActivateSprite(2)
		vic.sprites.ReadData(1) //ph1
		vic.sprites.ReadData(2) //ph1
	} else {
		vic.memory.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite4) != 0 {
		vic.SetBALow()
	}
	if vic.drawLine {
		vic.sprites.Draw()
		vic.borders.Draw()
		vic.beam.Commit()
	}

	vic.socketLastCycle()
	vic.rasterX += 8
}
