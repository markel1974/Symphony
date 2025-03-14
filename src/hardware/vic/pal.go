package mos6569

// _pal is a pointer to the starting cycleData node, representing the first cycle in the PAL video cycle sequence.
var _pal *cycleData

// sprite0 to sprite7 represent bitmask values for individual sprite identifiers.
const (
	sprite0 = 0x1
	sprite1 = 0x2
	sprite2 = 0x4
	sprite3 = 0x8
	sprite4 = 0x10
	sprite5 = 0x20
	sprite6 = 0x40
	sprite7 = 0x80
)

// init initializes the palette cycle data, assigns cycle borders, and connects cycles to create a linked sequence.
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

// palCycle1 handles the raster line processing logic for the VIC, including display updates and sprite DMA handling.
func palCycle1(vic *VIC) {
	if rasterY := vic.GetRasterY(); rasterY == RasterYMax {
		vic.vBlankNextCycle = true
	} else {
		vic.IncrementRasterY()
		vic.drawLine = (rasterY >= FirstDisplayedLine) && (rasterY <= LastDisplayedLine)
	}
	vic.borders.Reset()
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite3) != 0 {
		vic.sprites.FetchPtr(3)
		vic.sprites.Fetch(3, 0)
	}
	if vic.sprites.GetDMAFlag(sprite3|sprite4) == 0 {
		vic.ClearBALow()
	}
}

// palCycle2 handles the video interface chip's state transitions at the start of a PAL screen cycle.
func palCycle2(vic *VIC) {
	if vic.vBlankNextCycle {
		vic.vBlankNextCycle = false
		vic.lineStart = 0
		vic.graphics.ResetVideoCounterLatch()
		vic.ResetRasterY()
		vic.socket.VBlank()
	}
	vic.collisions.ClearGraphics()
	vic.graphics.SetOffset(vic.lineStart)
	vic.sprites.SetOffset(vic.lineStart)
	vic.borders.SetOffset(vic.lineStart)
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite3) != 0 {
		vic.sprites.Fetch(3, 1)
		vic.sprites.Fetch(3, 2)
	} else {
		vic.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(sprite5) != 0 {
		vic.SetBALow()
	}
}

// palCycle3 handles sprite DMA access and display updates for the VIC, modifying display parameters accordingly.
func palCycle3(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite4) != 0 {
		vic.sprites.FetchPtr(4)
		vic.sprites.Fetch(4, 0)
	}
	if vic.sprites.GetDMAFlag(sprite4|sprite5) == 0 {
		vic.ClearBALow()
	}
}

// palCycle4 handles sprite DMA operations for sprite 4 and modifies display access based on VIC state.
// It checks if the DMA flag for sprite 4 is set, performs fetch operations if true, otherwise idles the VIC.
// Additionally, it sets the BA low if the DMA flag for sprite 6 is set.
func palCycle4(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite4) != 0 {
		vic.sprites.Fetch(4, 1)
		vic.sprites.Fetch(4, 2)
	} else {
		vic.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(sprite6) != 0 {
		vic.SetBALow()
	}
}

// palCycle5 handles sprite DMA logic specifically for sprite 5 and updates display or bus access accordingly.
func palCycle5(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite5) != 0 {
		vic.sprites.FetchPtr(5)
		vic.sprites.Fetch(5, 0)
	}
	if vic.sprites.GetDMAFlag(sprite5|sprite6) == 0 {
		vic.ClearBALow()
	}
}

// palCycle6 performs DMA and idle access handling for sprites and sets the BA low flag in VIC graphics operations.
func palCycle6(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite5) != 0 {
		vic.sprites.Fetch(5, 1)
		vic.sprites.Fetch(5, 2)
	} else {
		vic.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(sprite7) != 0 {
		vic.SetBALow()
	}
}

// palCycle7 manages sprite DMA operations and display access for sprites 6 and 7 within the VIC graphics system.
// It acquires display access, checks DMA flags, fetches sprite data, and adjusts bus access levels accordingly.
func palCycle7(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite6) != 0 {
		vic.sprites.FetchPtr(6)
		vic.sprites.Fetch(6, 0)
	}
	if vic.sprites.GetDMAFlag(sprite6|sprite7) == 0 {
		vic.ClearBALow()
	}
}

// palCycle8 manages DMA access and fetch operations for sprite 6 in the VIC graphics system if the DMA flag is set.
func palCycle8(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite6) != 0 {
		vic.sprites.Fetch(6, 1)
		vic.sprites.Fetch(6, 2)
	} else {
		vic.AccessIdle()
	}
}

// palCycle9 processes VIC sprite DMA flag and display access for sprite7, managing sprite fetching and BA signal control.
func palCycle9(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite7) != 0 {
		vic.sprites.FetchPtr(7)
		vic.sprites.Fetch(7, 0)
	}
	if vic.sprites.GetDMAFlag(sprite7) == 0 {
		vic.ClearBALow()
	}
}

// palCycle10 manages display access and sprite data fetch operation for sprite 7 based on its DMA flag status.
func palCycle10(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite7) != 0 {
		vic.sprites.Fetch(7, 1)
		vic.sprites.Fetch(7, 2)
	} else {
		vic.AccessIdle()
	}
}

// palCycle11 manipulates the VIC object to acquire display access, clear the BA low signal, and perform a refresh action.
func palCycle11(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	vic.ClearBALow()
	vic.AccessRefresh()
}

// palCycle12 handles the VIC's PAL refresh cycle by accessing refresh operations and managing graphics and display access.
// It also triggers BA low if required during a bad line condition.
func palCycle12(vic *VIC) {
	vic.AccessRefresh()
	vic.graphics.TryAcquireDisplayAccess()
	vic.TryBALowIfBadLine()
}

// palCycle13 processes a single PAL (Phase Alternating Line) cycle, managing VIC-II border colors, graphics, and raster refresh.
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

// palCycle14 executes a sequence of operations for handling VIC graphics, borders, and video counter updates for a PAL cycle 14.
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

// palCycle15 processes a single VIC-II cycle, updating visual elements, sprite counters, and handling display access logic.
func palCycle15(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.AccessRefresh()
	vic.graphics.TryAcquireDisplayAccess()
	vic.sprites.UpdateCounterBase()
	vic.graphics.ResetLineIndex()
	vic.TryBALowIfBadLine()
	vic.graphics.TryPhi2Access()
}

// palCycle16 executes a set of operations per cycle for the VIC chip in PAL mode, handling graphics, borders, and sprite updates.
func palCycle16(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.graphics.TryGraphicsAccess()
	vic.graphics.TryAcquireDisplayAccess()
	vic.sprites.UpdateCounterBaseDMA()
	vic.TryBALowIfBadLine()
	vic.graphics.TryPhi2Access()
}

// palCycle17 executes the PAL-based VIC-II emulation cycle 17 logic, updating borders, graphics, and access conditions.
func palCycle17(vic *VIC) {
	vic.borders.UpdateColumn40()
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		if vic.borders.GetVerticalFlipFlop() {
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

// palCycle18 executes a sequence of operations for the VIC, updating borders, graphics, and access cycles efficiently.
// It determines whether to draw background or foreground based on the vertical flip-flop state and current line settings.
// The function manages access attempts to graphics and display while handling character data and bad line conditions.
func palCycle18(vic *VIC) {
	vic.borders.UpdateColumn38()
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		if vic.borders.GetVerticalFlipFlop() {
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

// palCycle19to54 handles VIC-II operations between cycle 19 to 54 for rendering, including border and graphics updates.
// Uses state flags to determine whether to update background or foreground graphics during the specified cycle range.
// Tries memory and display accesses required for graphical operations and updates character data accordingly.
func palCycle19to54(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		if vic.borders.GetVerticalFlipFlop() {
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

// palCycle55 processes a line cycle for the VIC-II chip, handling border colors, graphics rendering, sprite updates, and DMA logic.
func palCycle55(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		if vic.borders.GetVerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
	}
	vic.graphics.TryGraphicsAccess()
	vic.graphics.TryAcquireDisplayAccess()
	vic.UpdateSpriteExpY()
	vic.sprites.UpdateDMA()
	if vic.sprites.GetDMAFlag(sprite0) != 0 {
		vic.SetBALow() //BALow for Sprite 0 [cycle 58 = 55 + 3]
	} else {
		vic.ClearBALow() //Clear BALow for Sprite 0
	}
}

// palCycle56 handles cycle 56 operations on the VIC, including border drawing, DMA updates, and sprite BALow settings.
func palCycle56(vic *VIC) {
	vic.borders.EnableColumn38()
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		if vic.borders.GetVerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
	}
	vic.AccessIdle()
	vic.graphics.TryAcquireDisplayAccess()
	vic.sprites.UpdateDMA()
	if vic.sprites.GetDMAFlag(sprite0) != 0 {
		//Wrong cycle 59
		vic.SetBALow() //BALow for Sprite 0 [cycle 59 = 56 + 3]
	}
}

// palCycle57 performs a sequence of operations for VIC, enabling 40-column mode, updating sprites, and managing display.
// It handles tasks such as drawing the background, acquiring colors, and setting sprite-specific DMA flags.
func palCycle57(vic *VIC) {
	vic.borders.EnableColumn40()
	vic.sprites.UpdateDisplayFlags()
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.AccessIdle()
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite1) != 0 {
		vic.SetBALow() //BALow for Sprite 1 [cycle 60 = 57 + 3]
	}
}

// palCycle58 executes the 58th PAL cycle operations for the VIC, processing borders, sprites, and graphics updates.
func palCycle58(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.sprites.UpdateDisplayYFlags()
	if vic.sprites.GetDMAFlag(sprite0) != 0 {
		vic.sprites.FetchPtr(0) //phi1
		vic.sprites.Fetch(0, 0) //phi2
	}
	vic.graphics.UpdateDisplayAccess()
}

// palCycle59 executes the 59th cycle of the PAL video timing, managing sprite DMA fetches and access states.
// It handles background graphic drawing, border color acquisition, and memory access for sprites.
// The function ensures proper DMA operations for sprites 0 and 2 and adjusts the BA (Bus Available) signal as needed.
func palCycle59(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite0) != 0 {
		vic.sprites.Fetch(0, 1) //phi1
		vic.sprites.Fetch(0, 2) //phi2
	} else {
		vic.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(sprite2) != 0 {
		vic.SetBALow() //BALow for Sprite 2 [cycle 62 = 59 + 3]
	}
}

// palCycle60 manages the rendering of VIC-II components for a single PAL clock cycle, ensuring proper sprite and border updates.
func palCycle60(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
		vic.sprites.Draw()
		vic.borders.Draw()
		vic.lineStart += DisplayX
	}
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite1) != 0 {
		vic.sprites.FetchPtr(1) //phi1
		vic.sprites.Fetch(1, 0) //phi2
	} else {
		vic.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(sprite1|sprite2) == 0 {
		vic.ClearBALow() //Clear BALow for Sprite 1 - 2
	}
}

// palCycle61 handles graphics display access and sprite DMA operations for VIC during the specific cycle 61 in PAL mode.
// It performs sprite fetch operations or switches to idle access if no DMA flag is detected for certain sprites.
// Additionally, it sets the BA low signal for sprite 3 based on its DMA flag status in the cycle.
func palCycle61(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite1) != 0 {
		vic.sprites.Fetch(1, 1) //phi1
		vic.sprites.Fetch(1, 2) //phi2
	} else {
		vic.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(sprite3) != 0 {
		vic.SetBALow() //BALow for Sprite 3 [cycle 1 = 61 + 3]
	}
}

// palCycle62 handles the PAL (Phase Alternating Line) cycle 62 operations for sprite processing and display access in the VIC-II.
// It attempts to acquire display access, fetches sprite pointers, and performs operations based on DMA flags for sprites 2 and 3.
// The method clears certain flags (e.g., BALow) when specific DMA conditions are not met for the given sprites.
func palCycle62(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite2) != 0 {
		vic.sprites.FetchPtr(2) //phi1
		vic.sprites.Fetch(2, 0) //phi2
	}
	if vic.sprites.GetDMAFlag(sprite2|sprite3) == 0 {
		vic.ClearBALow() //Clear BALow for Sprite 2 - 3
	}
}

// palCycle63 handles the operations and state updates for cycle 63 in the VIC chip's PAL display timing.
// It manages sprite DMA access, updates display flip-flop state, and interacts with video memory.
func palCycle63(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	vic.borders.UpdateVerticalFlipFlop()
	if vic.sprites.GetDMAFlag(sprite2) != 0 {
		vic.sprites.Fetch(2, 1) //phi1
		vic.sprites.Fetch(2, 2) //phi2
	} else {
		vic.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(sprite4) != 0 {
		vic.SetBALow() //BALow for Sprite 4 [cycle 3 = 63 + 3]
	}
	vic.socket.LastCycle()
}
