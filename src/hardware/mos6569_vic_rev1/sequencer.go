package mos6569

// SequencerData represents a node in a cyclic linked list to manage cycles and associated operations.
type SequencerData struct {
	fn          func(vic *VIC)
	next        *SequencerData
	cycle       uint8
	cycleBorder uint8
}

// CreatePalSequencer initializes the PAL video timing cycle data. It constructs a circular linked list of 63 cycleData nodes,
// where each node represents one CPU clock cycle of a single PAL scanline. It pre-calculates border-related
// values for each cycle and links them in sequence to form the complete 63-cycle sequencer.
func CreatePalSequencer() []*SequencerData {
	const palBorderFirstCycle uint8 = 13

	var sequencer []*SequencerData
	sequencer = append(sequencer, &SequencerData{fn: phaseInitAndSprite3DMAPtrData0})
	sequencer = append(sequencer, &SequencerData{fn: phaseVBlankAndSprite3DMAData1Data2})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite4DMAPtrData0})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite4DMAData1Data2})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite5DMAPtrData0})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite5DMAData1Data2})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite6DMAPtrData0})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite6DMAData1Data2})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite7DMAPtrData0})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite7DMAData1Data2})
	sequencer = append(sequencer, &SequencerData{fn: phaseRefresh})
	sequencer = append(sequencer, &SequencerData{fn: phaseSetupBadLineCheck})
	sequencer = append(sequencer, &SequencerData{fn: phaseSetupRasterXReset})
	sequencer = append(sequencer, &SequencerData{fn: phaseSetupVCounterLoad})
	sequencer = append(sequencer, &SequencerData{fn: phaseSetupRCounterCheckAndSpritePipe1})
	sequencer = append(sequencer, &SequencerData{fn: phaseDisplayFirstFetchAndSpritePipe2})
	sequencer = append(sequencer, &SequencerData{fn: phaseDisplayMainFetchC40})
	sequencer = append(sequencer, &SequencerData{fn: phaseDisplayMainFetchC38})
	for x := 19; x <= 54; x++ {
		sequencer = append(sequencer, &SequencerData{fn: phaseDisplayMainFetch})
	}
	sequencer = append(sequencer, &SequencerData{fn: phaseTeardownLastFetchAndDMASetup})
	sequencer = append(sequencer, &SequencerData{fn: phaseTeardownIdle})
	sequencer = append(sequencer, &SequencerData{fn: phaseTeardownCommitSpriteFlags})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite0DMAPtrData0})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite0DMAData1Data2})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite1DMAPtrData0})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite1DMAData1Data2})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite2DMAPtrData0})
	sequencer = append(sequencer, &SequencerData{fn: phaseTeardownFinalSprite2DMA})

	last := len(sequencer) - 1
	for idx := 0; idx < len(sequencer); idx++ {
		sequencer[idx].cycleBorder = 0xff
		sequencer[idx].cycle = uint8(idx) + 1
		if sequencer[idx].cycle >= palBorderFirstCycle {
			sequencer[idx].cycleBorder = sequencer[idx].cycle - palBorderFirstCycle
		}
		if idx == last {
			sequencer[idx].next = sequencer[0]
		} else {
			sequencer[idx].next = sequencer[idx+1]
		}
	}
	return sequencer
}

// CreateNtscSequencer constructs a sequencer of cycles for NTSC display operation, including phases and DMA management.
// It defines the NTSC-specific scanline timing logic based on internal cycles and boundary conditions.
// Returns a slice of SequencerData nodes configured in a cyclic linked list structure.
func CreateNtscSequencer() []*SequencerData {
	// The border logic in NTSC starts slightly later
	const ntscBorderFirstCycle uint8 = 15
	var sequencer []*SequencerData

	// Initial Phases (H-Blank and Sprite DMA 3-7)
	// This sequence is identical to PAL. The hardware performs the same operations
	// at the start of a scan line.
	sequencer = append(sequencer, &SequencerData{fn: phaseInitAndSprite3DMAPtrData0})
	sequencer = append(sequencer, &SequencerData{fn: phaseVBlankAndSprite3DMAData1Data2})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite4DMAPtrData0})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite4DMAData1Data2})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite5DMAPtrData0})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite5DMAData1Data2})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite6DMAPtrData0})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite6DMAData1Data2})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite7DMAPtrData0})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite7DMAData1Data2})

	// Setup Drawing Window (identical to PAL)
	sequencer = append(sequencer, &SequencerData{fn: phaseRefresh})
	sequencer = append(sequencer, &SequencerData{fn: phaseSetupBadLineCheck})
	sequencer = append(sequencer, &SequencerData{fn: phaseSetupRasterXReset})
	sequencer = append(sequencer, &SequencerData{fn: phaseSetupVCounterLoad})
	sequencer = append(sequencer, &SequencerData{fn: phaseSetupRCounterCheckAndSpritePipe1})
	sequencer = append(sequencer, &SequencerData{fn: phaseDisplayFirstFetchAndSpritePipe2})

	// Main Drawing Window (Extended for NTSC)
	sequencer = append(sequencer, &SequencerData{fn: phaseDisplayMainFetchC40})
	sequencer = append(sequencer, &SequencerData{fn: phaseDisplayMainFetchC38})

	// Here the difference: The loop lasts 2 cycles more (up to 56 instead of 54)
	for x := 19; x <= 56; x++ {
		sequencer = append(sequencer, &SequencerData{fn: phaseDisplayMainFetch})
	}

	// DMA Cleanup and Preparation (identical to PAL, but shifted)
	sequencer = append(sequencer, &SequencerData{fn: phaseTeardownLastFetchAndDMASetup})
	sequencer = append(sequencer, &SequencerData{fn: phaseTeardownIdle})
	sequencer = append(sequencer, &SequencerData{fn: phaseTeardownCommitSpriteFlags})

	// DMA Sprite 0-2 and Final Cycle (identical to PAL, but shifted)
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite0DMAPtrData0})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite0DMAData1Data2})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite1DMAPtrData0})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite1DMAData1Data2})
	sequencer = append(sequencer, &SequencerData{fn: phaseSprite2DMAPtrData0})
	sequencer = append(sequencer, &SequencerData{fn: phaseTeardownFinalSprite2DMA})

	last := len(sequencer) - 1
	for idx := 0; idx < len(sequencer); idx++ {
		sequencer[idx].cycleBorder = 0xff
		sequencer[idx].cycle = uint8(idx) + 1
		if sequencer[idx].cycle >= ntscBorderFirstCycle {
			sequencer[idx].cycleBorder = sequencer[idx].cycle - ntscBorderFirstCycle
		}
		if idx == last {
			sequencer[idx].next = sequencer[0]
		} else {
			sequencer[idx].next = sequencer[idx+1]
		}
	}
	return sequencer
}

// phaseInitAndSprite3DMAPtrData0: This cycle marks the beginning of the horizontal blanking period. The raster line counter (rasterY)
// is checked against the maximum value. If it matches, a V-Blank is scheduled for the next cycle. Otherwise,
// rasterY is incremented for the new scanline. Sprite 3 DMA for the upcoming line begins if enabled, fetching
// the sprite pointer (phi1) and the first byte of sprite data (phi2).
//
//go:nosplit
func phaseInitAndSprite3DMAPtrData0(vic *VIC) {
	if rasterY := vic.GetRasterY(); rasterY == RasterYMax {
		vic.vBlankNextCycle = true
	} else {
		vic.IncrementRasterY()
		vic.drawLine = (rasterY >= FirstDisplayedLine) && (rasterY <= LastDisplayedLine)
	}
	vic.borders.ColumnInitialize()
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite3) != 0 {
		vic.sprites.FetchPtr(3)     //phi1
		vic.sprites.FetchData(3, 0) //phi2
	}
	if vic.sprites.GetDMAFlag(bitSprite3|bitSprite4) == 0 {
		vic.ClearBALow()
	}
}

// phaseVBlankAndSprite3DMAData1Data2: The V-Blank is triggered if scheduled in the previous cycle, resetting raster counters and firing
// the V-Blank interrupt. Internal pointers for graphics, sprites, and borders are reset. Sprite 3 DMA
// continues, fetching the second and third bytes of sprite data. The BA (Bus-Available) signal is asserted
// (pulled low) to prepare for sprite 5 DMA if it is enabled for the upcoming line.
//
//go:nosplit
func phaseVBlankAndSprite3DMAData1Data2(vic *VIC) {
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
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite3) != 0 {
		vic.sprites.FetchData(3, 1) //phi1
		vic.sprites.FetchData(3, 2) //phi2
	} else {
		vic.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite5) != 0 {
		vic.SetBALow()
	}
}

// phaseSprite4DMAPtrData0: Sprite 4 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 4 and 5.
//
//go:nosplit
func phaseSprite4DMAPtrData0(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite4) != 0 {
		vic.sprites.FetchPtr(4)     //phi1
		vic.sprites.FetchData(4, 0) //phi2
	}
	if vic.sprites.GetDMAFlag(bitSprite4|bitSprite5) == 0 {
		vic.ClearBALow()
	}
}

// phaseSprite4DMAData1Data2: Sprite 4 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 6 DMA if it is enabled.
//
//go:nosplit
func phaseSprite4DMAData1Data2(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite4) != 0 {
		vic.sprites.FetchData(4, 1) //phi1
		vic.sprites.FetchData(4, 2) //phi2
	} else {
		vic.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite6) != 0 {
		vic.SetBALow()
	}
}

// phaseSprite5DMAPtrData0: Sprite 5 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 5 and 6.
//
//go:nosplit
func phaseSprite5DMAPtrData0(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite5) != 0 {
		vic.sprites.FetchPtr(5)     //phi1
		vic.sprites.FetchData(5, 0) //phi2
	}
	if vic.sprites.GetDMAFlag(bitSprite5|bitSprite6) == 0 {
		vic.ClearBALow()
	}
}

// phaseSprite5DMAData1Data2: Sprite 5 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 7 DMA if it is enabled.
//
//go:nosplit
func phaseSprite5DMAData1Data2(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite5) != 0 {
		vic.sprites.FetchData(5, 1) //phi1
		vic.sprites.FetchData(5, 2) //phi2
	} else {
		vic.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite7) != 0 {
		vic.SetBALow()
	}
}

// phaseSprite6DMAPtrData0: Sprite 6 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 6 and 7.
//
//go:nosplit
func phaseSprite6DMAPtrData0(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite6) != 0 {
		vic.sprites.FetchPtr(6)     //phi1
		vic.sprites.FetchData(6, 0) //phi12
	}
	if vic.sprites.GetDMAFlag(bitSprite6|bitSprite7) == 0 {
		vic.ClearBALow()
	}
}

// phaseSprite6DMAData1Data2: Sprite 6 DMA continues, fetching its second and third data bytes.
//
//go:nosplit
func phaseSprite6DMAData1Data2(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite6) != 0 {
		vic.sprites.FetchData(6, 1) //phi1
		vic.sprites.FetchData(6, 2) //phi2
	} else {
		vic.AccessIdle()
	}
}

// phaseSprite7DMAPtrData0: Sprite 7 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprite 7.
//
//go:nosplit
func phaseSprite7DMAPtrData0(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite7) != 0 {
		vic.sprites.FetchPtr(7)     //phi1
		vic.sprites.FetchData(7, 0) //phi2
	}
	if vic.sprites.GetDMAFlag(bitSprite7) == 0 {
		vic.ClearBALow()
	}
}

// phaseSprite7DMAData1Data2: Sprite 7 DMA continues, fetching its second and third data bytes. This concludes the main
// sprite DMA phase for the upcoming line.
//
//go:nosplit
func phaseSprite7DMAData1Data2(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite7) != 0 {
		vic.sprites.FetchData(7, 1) //phi1
		vic.sprites.FetchData(7, 2) //phi2
	} else {
		vic.AccessIdle()
	}
}

// phaseRefresh: This is a refresh cycle. The VIC performs a DRAM refresh operation by accessing an address
// in the range $3C00-$3FFF. The address bus is released for the CPU, and the BA signal is cleared.
//
//go:nosplit
func phaseRefresh(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	vic.ClearBALow()
	vic.AccessRefresh()
}

// phaseSetupBadLineCheck: This is a refresh cycle. The VIC checks for a "bad line" condition, which occurs if the
// DEN bit is set and the lower 3 bits of the raster counter match the lower 3 bits of the YSCROLL register.
// If it is a bad line, the VIC prepares to halt the CPU by asserting the BA signal.
//
//go:nosplit
func phaseSetupBadLineCheck(vic *VIC) {
	vic.AccessRefresh()
	vic.graphics.TryAcquireDisplayAccess()
	vic.TryBALowIfBadLine()
}

// phaseSetupRasterXReset: This is a refresh cycle. The horizontal raster counter (rasterX) is reset to 0. The VIC is now
// in the left border area. The BA signal is asserted if the bad line condition was met in the previous cycle.
//
//go:nosplit
func phaseSetupRasterXReset(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.AccessRefresh()
	vic.graphics.TryAcquireDisplayAccess()
	vic.TryBALowIfBadLine()
	vic.ResetRasterX()
}

// phaseSetupVCounterLoad: This is a refresh cycle. The Video Counter (VC) is loaded from the Video Counter Base (VCBASE),
// pointing to the current character row in screen memory. The Row Counter (RC) is reset to 0 if this is the
// first scanline of a character row (i.e., rasterY & 7 == 0).
//
//go:nosplit
func phaseSetupVCounterLoad(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.AccessRefresh()
	vic.graphics.TryAcquireDisplayAccess()
	vic.graphics.TryResetRowCounter()
	vic.TryBALowIfBadLine()
	vic.graphics.UpdateVideoCounter()
}

// phaseSetupRCounterCheckAndSpritePipe1: This is the critical "bad line" decision cycle. If it's a bad line, the VIC takes full control
// of the bus for the next 40 cycles. The graphics pipeline begins its first access, fetching a character code
// from screen RAM using the Video Counter (VC). Sprite Y-expansion counters for the *next* scanline are checked.
//
//go:nosplit
func phaseSetupRCounterCheckAndSpritePipe1(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.AccessRefresh()
	vic.graphics.TryAcquireDisplayAccess()
	vic.sprites.TryIncrementCounterBase()
	vic.graphics.ResetLineIndex()
	vic.TryBALowIfBadLine()
	vic.graphics.TryPhi2Access()
}

// phaseDisplayFirstFetchAndSpritePipe2: First graphics data fetch cycle. The VIC fetches the character's bitmap data from Character
// ROM or RAM, using the character code fetched in the previous cycle and the current Row Counter (RC).
// Sprite Y-expansion counters are committed.
//
//go:nosplit
func phaseDisplayFirstFetchAndSpritePipe2(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.graphics.TryGraphicsAccess()
	vic.graphics.TryAcquireDisplayAccess()
	vic.sprites.TryIncrementCounterBase()
	vic.sprites.CommitIncrementCounterBase()
	vic.TryBALowIfBadLine()
	vic.graphics.TryPhi2Access()
}

// phaseDisplayMainFetchC40: Second graphics data fetch cycle. The VIC fetches the next character code from screen RAM.
// The first pixels of the 40-column window (or the side border in 38-column mode) are drawn. The side
// border logic for 40-column mode is triggered.
//
//go:nosplit
func phaseDisplayMainFetchC40(vic *VIC) {
	vic.borders.Column40Update()
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		if vic.borders.VerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
	}
	vic.graphics.TryGraphicsAccess()
	vic.graphics.TryAcquireDisplayAccess()
	vic.TryBALowIfBadLine()
	vic.graphics.TryPhi2Access()
}

// phaseDisplayMainFetchC38: The main display window begins. The VIC fetches the next character's bitmap data. The fetched
// bitmap data is loaded into the graphics pipeline's internal shift register. The side border logic for
// 38-column mode is triggered. Pixels are drawn.
//
//go:nosplit
func phaseDisplayMainFetchC38(vic *VIC) {
	vic.borders.Column38Update()
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		if vic.borders.VerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
	}
	vic.graphics.TryGraphicsAccess()
	vic.graphics.TryAcquireDisplayAccess()
	vic.TryBALowIfBadLine()
	vic.graphics.TryPhi2Access()
	vic.graphics.UpdateCharDataLast()
}

// phaseDisplayMainFetch: These 36 cycles form the core of the visible display area. In each pair of cycles, the VIC
// fetches a character code from screen RAM and its corresponding bitmap data. The graphics pipeline
// continuously shifts out 8 pixels per character, drawing either foreground or background pixels. On a "bad line",
// the CPU is halted throughout this entire phase.
//
//go:nosplit
func phaseDisplayMainFetch(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		if vic.borders.VerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
	}
	vic.graphics.TryGraphicsAccess()
	vic.graphics.TryAcquireDisplayAccess()
	vic.TryBALowIfBadLine()
	vic.graphics.TryPhi2Access()
	vic.graphics.UpdateCharDataLast()
}

// phaseTeardownLastFetchAndDMASetup: This is the last graphics data fetch cycle for the line. The VIC finalizes the graphics pipeline
// and prepares for the *next* scanline by updating sprite Y-expansion flags and calculating which sprites will
// be active (setting up their DMA flags). The BA signal is asserted to prepare for sprite 0 DMA.
//
//go:nosplit
func phaseTeardownLastFetchAndDMASetup(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		if vic.borders.VerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
	}
	vic.graphics.TryGraphicsAccess()
	vic.graphics.TryAcquireDisplayAccess()
	vic.UpdateSpriteExpY()
	vic.sprites.UpdateDMA()
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
func phaseTeardownIdle(vic *VIC) {
	vic.borders.Column38Apply()
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		if vic.borders.VerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
	}
	vic.AccessIdle()
	vic.graphics.TryAcquireDisplayAccess()
	vic.sprites.UpdateDMA()
	if vic.sprites.GetDMAFlag(bitSprite0) != 0 {
		vic.SetBALow()
	}
}

// phaseTeardownCommitSpriteFlags: Idle cycle. The 40-column side border logic is applied. The sprite DMA flags calculated in
// cycle 55 are now committed for use in the upcoming DMA phase. The BA signal is asserted to prepare for
// sprite 1 DMA.
//
//go:nosplit
func phaseTeardownCommitSpriteFlags(vic *VIC) {
	vic.borders.Column40Apply()
	vic.sprites.CommitSpriteFlags()
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.AccessIdle()
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite1) != 0 {
		vic.SetBALow() //BALow for Sprite 1 [cycle 60 = 57 + 3]
	}
}

// phaseSprite0DMAPtrData0: Sprite 0 DMA for the upcoming line begins if enabled, fetching its pointer and first data byte.
// Sprite flags are prepared for the line *after* the next one.
//
//go:nosplit
func phaseSprite0DMAPtrData0(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.sprites.PrepareSpriteFlags()
	if vic.sprites.GetDMAFlag(bitSprite0) != 0 {
		vic.sprites.FetchPtr(0)     //phi1
		vic.sprites.FetchData(0, 0) //phi2
	}
	vic.graphics.UpdateDisplayAccess()
}

// phaseSprite0DMAData1Data2: Sprite 0 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 2 DMA.
//
//go:nosplit
func phaseSprite0DMAData1Data2(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite0) != 0 {
		vic.sprites.FetchData(0, 1) //phi1
		vic.sprites.FetchData(0, 2) //phi2
	} else {
		vic.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite2) != 0 {
		vic.SetBALow() //BALow for Sprite 2 [cycle 62 = 59 + 3]
	}
}

// phaseSprite1DMAPtrData0: Sprite 1 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 1 and 2.
//
//go:nosplit
func phaseSprite1DMAPtrData0(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	//if vic.drawLine {
	//	vic.graphics.DrawBackground()
	//	vic.sprites.Draw()
	//	vic.borders.Draw()
	//	vic.lineStart += DisplayX
	//}
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite1) != 0 {
		vic.sprites.FetchPtr(1)     //phi1
		vic.sprites.FetchData(1, 0) //phi2
	} else {
		vic.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite1|bitSprite2) == 0 {
		vic.ClearBALow()
	}
}

// phaseSprite1DMAData1Data2: Sprite 1 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 3 DMA.
//
//go:nosplit
func phaseSprite1DMAData1Data2(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite1) != 0 {
		vic.sprites.FetchData(1, 1) //phi1
		vic.sprites.FetchData(1, 2) //phi2
	} else {
		vic.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite3) != 0 {
		vic.SetBALow()
	}
}

// phaseSprite2DMAPtrData0: Sprite 2 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 2 and 3.
//
//go:nosplit
func phaseSprite2DMAPtrData0(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite2) != 0 {
		vic.sprites.FetchPtr(2)     //phi1
		vic.sprites.FetchData(2, 0) //phi2
	}
	if vic.sprites.GetDMAFlag(bitSprite2|bitSprite3) == 0 {
		vic.ClearBALow()
	}
}

// phaseSprite1DMAData1Data2: Sprite 2 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 4 DMA.
//
//go:nosplit
func phaseTeardownFinalSprite2DMA(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	vic.borders.UpdateVerticalFlipFlop()
	if vic.sprites.GetDMAFlag(bitSprite2) != 0 {
		vic.sprites.FetchData(2, 1) //phi1
		vic.sprites.FetchData(2, 2) //phi2
	} else {
		vic.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(bitSprite4) != 0 {
		vic.SetBALow()
	}

	if vic.drawLine {
		vic.sprites.Draw()
		vic.borders.Draw()
		vic.lineStart += DisplayX
	}

	vic.socketLastCycle()
}
