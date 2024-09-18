package mos6569

var _pal *cycleData

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

// Check BA for matrix fetch
// Check BA for Sprite Phi2 fetch

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

func palCycle1(vic *VIC) {
	if rasterY := vic.core.GetRasterY(); rasterY == RasterYMax {
		vic.vBlankNextCycle = true
	} else {
		vic.core.IncrementRasterY()
		vic.drawLine = (rasterY >= FirstDisplayedLine) && (rasterY <= LastDisplayedLine)
	}
	vic.borders.Reset()
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite3) != 0 {
		vic.sprites.FetchPtr(3)
		vic.sprites.Fetch(3, 0)
	}
	if vic.sprites.GetDMAFlag(sprite3|sprite4) == 0 {
		vic.core.ClearBALow()
	}
}

func palCycle2(vic *VIC) {
	if vic.vBlankNextCycle {
		vic.vBlankNextCycle = false
		vic.lineStart = 0
		vic.graphics.ResetVideoCounterBase()
		vic.core.ResetRasterY()
		vic.core.socket.VBlank()
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
		vic.core.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(sprite5) != 0 {
		vic.core.SetBALow()
	}
}

func palCycle3(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite4) != 0 {
		vic.sprites.FetchPtr(4)
		vic.sprites.Fetch(4, 0)
	}
	if vic.sprites.GetDMAFlag(sprite4|sprite5) == 0 {
		vic.core.ClearBALow()
	}
}

func palCycle4(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite4) != 0 {
		vic.sprites.Fetch(4, 1)
		vic.sprites.Fetch(4, 2)
	} else {
		vic.core.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(sprite6) != 0 {
		vic.core.SetBALow()
	}
}

func palCycle5(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite5) != 0 {
		vic.sprites.FetchPtr(5)
		vic.sprites.Fetch(5, 0)
	}
	if vic.sprites.GetDMAFlag(sprite5|sprite6) == 0 {
		vic.core.ClearBALow()
	}
}

func palCycle6(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite5) != 0 {
		vic.sprites.Fetch(5, 1)
		vic.sprites.Fetch(5, 2)
	} else {
		vic.core.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(sprite7) != 0 {
		vic.core.SetBALow()
	}
}

func palCycle7(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite6) != 0 {
		vic.sprites.FetchPtr(6)
		vic.sprites.Fetch(6, 0)
	}
	if vic.sprites.GetDMAFlag(sprite6|sprite7) == 0 {
		vic.core.ClearBALow()
	}
}

func palCycle8(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite6) != 0 {
		vic.sprites.Fetch(6, 1)
		vic.sprites.Fetch(6, 2)
	} else {
		vic.core.AccessIdle()
	}
}

func palCycle9(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite7) != 0 {
		vic.sprites.FetchPtr(7)
		vic.sprites.Fetch(7, 0)
	}
	if vic.sprites.GetDMAFlag(sprite7) == 0 {
		vic.core.ClearBALow()
	}
}

func palCycle10(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite7) != 0 {
		vic.sprites.Fetch(7, 1)
		vic.sprites.Fetch(7, 2)
	} else {
		vic.core.AccessIdle()
	}
}

func palCycle11(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	vic.core.ClearBALow()
	vic.core.AccessRefresh()
}

func palCycle12(vic *VIC) {
	vic.core.AccessRefresh()
	vic.graphics.TryAcquireDisplayAccess()
	vic.core.TryBALowIfBadLine()
}

func palCycle13(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.core.AccessRefresh()
	vic.graphics.TryAcquireDisplayAccess()
	vic.core.TryBALowIfBadLine()
	vic.core.ResetRasterX()
}

func palCycle14(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.core.AccessRefresh()
	vic.graphics.TryAcquireDisplayAccess()
	vic.graphics.TryResetRowCounter()
	vic.core.TryBALowIfBadLine()
	vic.graphics.UpdateVideoCounter()
}

func palCycle15(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.core.AccessRefresh()
	vic.graphics.TryAcquireDisplayAccess()
	vic.sprites.UpdateCounterBase()
	vic.graphics.ResetLineIndex()
	vic.core.TryBALowIfBadLine()
	//vic.core.TryAcquireAEC()
	vic.graphics.TryVideoMatrixAccess()
}

func palCycle16(vic *VIC) {
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.graphics.TryGraphicsAccess()
	vic.graphics.TryAcquireDisplayAccess()
	vic.sprites.UpdateCounterBaseDMA()
	vic.core.TryBALowIfBadLine()
	//vic.core.TryAcquireAEC()
	vic.graphics.TryVideoMatrixAccess()
}

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
	vic.core.TryBALowIfBadLine()
	//vic.core.TryAcquireAEC()
	vic.graphics.TryVideoMatrixAccess()
}

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
	vic.core.TryBALowIfBadLine()
	//vic.core.TryAcquireAEC()
	vic.graphics.TryVideoMatrixAccess()
	vic.graphics.UpdateLastCharData()
}

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
	vic.core.TryBALowIfBadLine()
	//vic.core.TryAcquireAEC()
	vic.graphics.TryVideoMatrixAccess()
	vic.graphics.UpdateLastCharData()
}

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
	vic.core.UpdateSpriteExpY()
	vic.sprites.UpdateDMA()
	if vic.sprites.GetDMAFlag(sprite0) != 0 {
		vic.core.SetBALow() //BALow for Sprite 0 [cycle 58 = 55 + 3]
	} else {
		vic.core.ClearBALow() //Clear BALow for Sprite 0
	}
}

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
	vic.core.AccessIdle()
	vic.graphics.TryAcquireDisplayAccess()
	vic.sprites.UpdateDMA()
	if vic.sprites.GetDMAFlag(sprite0) != 0 {
		//Wrong cycle 59
		vic.core.SetBALow() //BALow for Sprite 0 [cycle 59 = 56 + 3]
	}
}

func palCycle57(vic *VIC) {
	vic.borders.EnableColumn40()
	vic.sprites.UpdateDisplayFlags()
	if vic.drawLine {
		vic.borders.AcquireColor(vic.curr.cycleBorder)
		vic.graphics.DrawBackground()
	}
	vic.core.AccessIdle()
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite1) != 0 {
		vic.core.SetBALow() //BALow for Sprite 1 [cycle 60 = 57 + 3]
	}
}

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
		vic.core.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(sprite2) != 0 {
		vic.core.SetBALow() //BALow for Sprite 2 [cycle 62 = 59 + 3]
	}
}

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
		vic.core.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(sprite1|sprite2) == 0 {
		vic.core.ClearBALow() //Clear BALow for Sprite 1 - 2
	}
}

func palCycle61(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite1) != 0 {
		vic.sprites.Fetch(1, 1) //phi1
		vic.sprites.Fetch(1, 2) //phi2
	} else {
		vic.core.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(sprite3) != 0 {
		vic.core.SetBALow() //BALow for Sprite 3 [cycle 1 = 61 + 3]
	}
}

func palCycle62(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(sprite2) != 0 {
		vic.sprites.FetchPtr(2) //phi1
		vic.sprites.Fetch(2, 0) //phi2
	}
	if vic.sprites.GetDMAFlag(sprite2|sprite3) == 0 {
		vic.core.ClearBALow() //Clear BALow for Sprite 2 - 3
	}
}

func palCycle63(vic *VIC) {
	vic.graphics.TryAcquireDisplayAccess()
	vic.borders.UpdateVerticalFlipFlop()
	if vic.sprites.GetDMAFlag(sprite2) != 0 {
		vic.sprites.Fetch(2, 1) //phi1
		vic.sprites.Fetch(2, 2) //phi2
	} else {
		vic.core.AccessIdle()
	}
	if vic.sprites.GetDMAFlag(sprite4) != 0 {
		vic.core.SetBALow() //BALow for Sprite 4 [cycle 3 = 63 + 3]
	}
	vic.core.socket.LastCycle()
}
