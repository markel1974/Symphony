package mos6569

// _pal is a pointer to the starting cycleData node, representing the first cycle in the PAL video cycle sequence.
var _pal *cycleData

// init initializes the PAL video timing cycle data. It constructs a circular linked list of 63 cycleData nodes,
// where each node represents one CPU clock cycle of a single PAL scanline. It pre-calculates border-related
// values for each cycle and links them in sequence to form the complete 63-cycle timeline.
func init() {
	const palBorderFirstCycle uint8 = 13

	var palCycles []*cycleData
	palCycles = append(palCycles, &cycleData{fn: palCycle1})
	palCycles = append(palCycles, &cycleData{fn: palCycle2})
	palCycles = append(palCycles, &cycleData{fn: palCycle3})
	palCycles = append(palCycles, &cycleData{fn: palCycle4})
	palCycles = append(palCycles, &cycleData{fn: palCycle5})
	palCycles = append(palCycles, &cycleData{fn: palCycle6})
	palCycles = append(palCycles, &cycleData{fn: palCycle7})
	palCycles = append(palCycles, &cycleData{fn: palCycle8})
	palCycles = append(palCycles, &cycleData{fn: palCycle9})
	palCycles = append(palCycles, &cycleData{fn: palCycle10})
	palCycles = append(palCycles, &cycleData{fn: palCycle11})
	palCycles = append(palCycles, &cycleData{fn: palCycle12})
	palCycles = append(palCycles, &cycleData{fn: palCycle13})
	palCycles = append(palCycles, &cycleData{fn: palCycle14})
	palCycles = append(palCycles, &cycleData{fn: palCycle15})
	palCycles = append(palCycles, &cycleData{fn: palCycle16})
	palCycles = append(palCycles, &cycleData{fn: palCycle17})
	palCycles = append(palCycles, &cycleData{fn: palCycle18})
	for x := 19; x <= 54; x++ {
		palCycles = append(palCycles, &cycleData{fn: palCycle19to54})
	}
	palCycles = append(palCycles, &cycleData{fn: palCycle55})
	palCycles = append(palCycles, &cycleData{fn: palCycle56})
	palCycles = append(palCycles, &cycleData{fn: palCycle57})
	palCycles = append(palCycles, &cycleData{fn: palCycle58})
	palCycles = append(palCycles, &cycleData{fn: palCycle59})
	palCycles = append(palCycles, &cycleData{fn: palCycle60})
	palCycles = append(palCycles, &cycleData{fn: palCycle61})
	palCycles = append(palCycles, &cycleData{fn: palCycle62})
	palCycles = append(palCycles, &cycleData{fn: palCycle63})

	last := len(palCycles) - 1
	for idx := 0; idx < len(palCycles); idx++ {
		palCycles[idx].cycleBorder = 0xff
		palCycles[idx].cycle = uint8(idx) + 1
		if palCycles[idx].cycle >= palBorderFirstCycle {
			palCycles[idx].cycleBorder = palCycles[idx].cycle - palBorderFirstCycle
		}
		if idx == last {
			palCycles[idx].next = palCycles[0]
		} else {
			palCycles[idx].next = palCycles[idx+1]
		}
	}
	_pal = palCycles[0]
}

// palCycle1: This cycle marks the beginning of the horizontal blanking period. The raster line counter (rasterY)
// is checked against the maximum value. If it matches, a V-Blank is scheduled for the next cycle. Otherwise,
// rasterY is incremented for the new scanline. Sprite 3 DMA for the upcoming line begins if enabled, fetching
// the sprite pointer (phi1) and the first byte of sprite data (phi2).
//
//go:nosplit
func palCycle1(vic *VIC) {
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

// palCycle2: The V-Blank is triggered if scheduled in the previous cycle, resetting raster counters and firing
// the V-Blank interrupt. Internal pointers for graphics, sprites, and borders are reset. Sprite 3 DMA
// continues, fetching the second and third bytes of sprite data. The BA (Bus-Available) signal is asserted
// (pulled low) to prepare for sprite 5 DMA if it is enabled for the upcoming line.
//
//go:nosplit
func palCycle2(vic *VIC) {
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

// palCycle3: Sprite 4 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 4 and 5.
//
//go:nosplit
func palCycle3(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite4) != 0 {
		vic.sprites.FetchPtr(4)     //phi1
		vic.sprites.FetchData(4, 0) //phi2
	}
	if vic.sprites.GetDMAFlag(bitSprite4|bitSprite5) == 0 {
		vic.ClearBALow()
	}
}

// palCycle4: Sprite 4 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 6 DMA if it is enabled.
//
//go:nosplit
func palCycle4(vic *VIC) {
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

// palCycle5: Sprite 5 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 5 and 6.
//
//go:nosplit
func palCycle5(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite5) != 0 {
		vic.sprites.FetchPtr(5)     //phi1
		vic.sprites.FetchData(5, 0) //phi2
	}
	if vic.sprites.GetDMAFlag(bitSprite5|bitSprite6) == 0 {
		vic.ClearBALow()
	}
}

// palCycle6: Sprite 5 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 7 DMA if it is enabled.
//
//go:nosplit
func palCycle6(vic *VIC) {
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

// palCycle7: Sprite 6 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 6 and 7.
//
//go:nosplit
func palCycle7(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite6) != 0 {
		vic.sprites.FetchPtr(6)     //phi1
		vic.sprites.FetchData(6, 0) //phi12
	}
	if vic.sprites.GetDMAFlag(bitSprite6|bitSprite7) == 0 {
		vic.ClearBALow()
	}
}

// palCycle8: Sprite 6 DMA continues, fetching its second and third data bytes.
//
//go:nosplit
func palCycle8(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite6) != 0 {
		vic.sprites.FetchData(6, 1) //phi1
		vic.sprites.FetchData(6, 2) //phi2
	} else {
		vic.AccessIdle()
	}
}

// palCycle9: Sprite 7 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprite 7.
//
//go:nosplit
func palCycle9(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite7) != 0 {
		vic.sprites.FetchPtr(7)     //phi1
		vic.sprites.FetchData(7, 0) //phi2
	}
	if vic.sprites.GetDMAFlag(bitSprite7) == 0 {
		vic.ClearBALow()
	}
}

// palCycle10: Sprite 7 DMA continues, fetching its second and third data bytes. This concludes the main
// sprite DMA phase for the upcoming line.
//
//go:nosplit
func palCycle10(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite7) != 0 {
		vic.sprites.FetchData(7, 1) //phi1
		vic.sprites.FetchData(7, 2) //phi2
	} else {
		vic.AccessIdle()
	}
}

// palCycle11: This is a refresh cycle. The VIC performs a DRAM refresh operation by accessing an address
// in the range $3C00-$3FFF. The address bus is released for the CPU, and the BA signal is cleared.
//
//go:nosplit
func palCycle11(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	vic.ClearBALow()
	vic.AccessRefresh()
}

// palCycle12: This is a refresh cycle. The VIC checks for a "bad line" condition, which occurs if the
// DEN bit is set and the lower 3 bits of the raster counter match the lower 3 bits of the YSCROLL register.
// If it is a bad line, the VIC prepares to halt the CPU by asserting the BA signal.
//
//go:nosplit
func palCycle12(vic *VIC) {
	vic.AccessRefresh()
	vic.graphics.TryAcquireDisplayAccess()
	vic.TryBALowIfBadLine()
}

// palCycle13: This is a refresh cycle. The horizontal raster counter (rasterX) is reset to 0. The VIC is now
// in the left border area. The BA signal is asserted if the bad line condition was met in the previous cycle.
//
//go:nosplit
func palCycle13(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.AccessRefresh()
	vic.graphics.TryAcquireDisplayAccess()
	vic.TryBALowIfBadLine()
	vic.ResetRasterX()
}

// palCycle14: This is a refresh cycle. The Video Counter (VC) is loaded from the Video Counter Base (VCBASE),
// pointing to the current character row in screen memory. The Row Counter (RC) is reset to 0 if this is the
// first scanline of a character row (i.e., rasterY & 7 == 0).
//
//go:nosplit
func palCycle14(vic *VIC) {
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

// palCycle15: This is the critical "bad line" decision cycle. If it's a bad line, the VIC takes full control
// of the bus for the next 40 cycles. The graphics pipeline begins its first access, fetching a character code
// from screen RAM using the Video Counter (VC). Sprite Y-expansion counters for the *next* scanline are checked.
//
//go:nosplit
func palCycle15(vic *VIC) {
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

// palCycle16: First graphics data fetch cycle. The VIC fetches the character's bitmap data from Character
// ROM or RAM, using the character code fetched in the previous cycle and the current Row Counter (RC).
// Sprite Y-expansion counters are committed.
//
//go:nosplit
func palCycle16(vic *VIC) {
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

// palCycle17: Second graphics data fetch cycle. The VIC fetches the next character code from screen RAM.
// The first pixels of the 40-column window (or the side border in 38-column mode) are drawn. The side
// border logic for 40-column mode is triggered.
//
//go:nosplit
func palCycle17(vic *VIC) {
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

// palCycle18: The main display window begins. The VIC fetches the next character's bitmap data. The fetched
// bitmap data is loaded into the graphics pipeline's internal shift register. The side border logic for
// 38-column mode is triggered. Pixels are drawn.
//
//go:nosplit
func palCycle18(vic *VIC) {
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

// palCycle19to54: These 36 cycles form the core of the visible display area. In each pair of cycles, the VIC
// fetches a character code from screen RAM and its corresponding bitmap data. The graphics pipeline
// continuously shifts out 8 pixels per character, drawing either foreground or background pixels. On a "bad line",
// the CPU is halted throughout this entire phase.
//
//go:nosplit
func palCycle19to54(vic *VIC) {
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

// palCycle55: This is the last graphics data fetch cycle for the line. The VIC finalizes the graphics pipeline
// and prepares for the *next* scanline by updating sprite Y-expansion flags and calculating which sprites will
// be active (setting up their DMA flags). The BA signal is asserted to prepare for sprite 0 DMA.
//
//go:nosplit
func palCycle55(vic *VIC) {
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

// palCycle56: Idle cycle. The main graphics fetch is complete. The VIC is now in the right border area.
// The 38-column side border logic is applied. The BA signal is asserted for sprite 0 DMA if needed.
//
//go:nosplit
func palCycle56(vic *VIC) {
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

// palCycle57: Idle cycle. The 40-column side border logic is applied. The sprite DMA flags calculated in
// cycle 55 are now committed for use in the upcoming DMA phase. The BA signal is asserted to prepare for
// sprite 1 DMA.
//
//go:nosplit
func palCycle57(vic *VIC) {
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

// palCycle58: Sprite 0 DMA for the upcoming line begins if enabled, fetching its pointer and first data byte.
// Sprite flags are prepared for the line *after* the next one.
//
//go:nosplit
func palCycle58(vic *VIC) {
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

// palCycle59: Sprite 0 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 2 DMA.
//
//go:nosplit
func palCycle59(vic *VIC) {
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

// palCycle60: Sprite 1 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 1 and 2.
//
//go:nosplit
func palCycle60(vic *VIC) {
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

// palCycle61: Sprite 1 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 3 DMA.
//
//go:nosplit
func palCycle61(vic *VIC) {
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

// palCycle62: Sprite 2 DMA begins if enabled, fetching its pointer and first data byte. The BA signal is
// managed based on the DMA status of sprites 2 and 3.
//
//go:nosplit
func palCycle62(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(bitSprite2) != 0 {
		vic.sprites.FetchPtr(2)     //phi1
		vic.sprites.FetchData(2, 0) //phi2
	}
	if vic.sprites.GetDMAFlag(bitSprite2|bitSprite3) == 0 {
		vic.ClearBALow()
	}
}

// palCycle61: Sprite 2 DMA continues, fetching its second and third data bytes. The BA signal is asserted
// to prepare for sprite 4 DMA.
//
//go:nosplit
func palCycle63(vic *VIC) {
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
