package vic

import (
	"github.com/markel1974/c64emu/src/board/quartz"
	"github.com/markel1974/c64emu/src/config"
)

//https://dustlayer.com/c64-architecture

type MOS6569 struct {
	core           *Core
	cfg            *config.Config
	sprites        *Sprites
	graphics       *Graphics
	foreMask       *ForeMask
	cycle          int   // Cycle
	lineStart      int   // Offset from current line in bitmap buffer
	drawThisLine   bool  // This line is drawn
	refreshCounter uint8 // Refresh counter
	vBlanking      bool  // VBlank in next cycle
}

func NewMOS6569(db IDisplayBuffer) *MOS6569 {
	core := NewCore()
	foreMask := NewForeMask()
	vic := &MOS6569{
		core:         core,
		foreMask:     foreMask,
		graphics:     NewGraphics(core, foreMask, db),
		sprites:      NewSprites(core, foreMask, db),
		cycle:        1,
		vBlanking:    false,
		drawThisLine: false,
		cfg:          nil,
	}
	return vic
}

func (vic *MOS6569) Setup(quartz *quartz.Quartz, banks IBanks, cfg *config.Config) {
	//vic.board = board
	vic.cfg = cfg
	vic.cfg.Bind(vic.configChanged)
	vic.core.Setup(quartz, banks)
	vic.graphics.Setup()
	vic.sprites.Setup()
}

func (vic *MOS6569) Reset() {
	vic.core.ready = false
}

func (vic *MOS6569) GetLastByte() uint8 {
	return vic.core.lastByte
}

/*
func (vic *MOS6569) GetBALow() bool {
	return vic.core.baLow
}

func (vic *MOS6569) GetAECLow() bool {
	return vic.core.aecLow
}*/

func (vic *MOS6569) SignalBALowBind(fn func(bool)) {
	vic.core.signalBALow.Bind(fn)
}

func (vic *MOS6569) SignalAECLowBind(fn func(bool)) {
	vic.core.signalAECLow.Bind(fn)
}

func (vic *MOS6569) SignalReadyBind(fn func()) {
	vic.core.signalReady.Bind(fn)
}

func (vic *MOS6569) SignalTriggerIRQBind(fn func(uint32)) {
	vic.core.signalIRQTrigger.Bind(fn)
}

func (vic *MOS6569) SignalClearIRQBind(fn func(uint32)) {
	vic.core.signalIRQClear.Bind(fn)
}

//func (vic *MOS6569) GetBALowSignal() *signals.Signal1[bool] {
//	return vic.core.baLowSignal
//}

func (vic *MOS6569) configChanged() {
	//vic.skipFrames = vic.cfg.SkipFrames()
}

func (vic *MOS6569) ReadRegister(addr uint16) uint8 {
	return vic.core.ReadRegister(addr)
}

func (vic *MOS6569) WriteRegister(addr uint16, data uint8) {
	vic.core.WriteRegister(addr, data)
}

func (vic *MOS6569) ChangedVA(va uint8) {
	vic.core.ChangedVA(va)
}

func (vic *MOS6569) LightPenTrigger() {
	vic.core.LightPenTrigger()
}

func (vic *MOS6569) Emulate() (bool, bool) {
	vBlank := false
	lastCycle := false

	//vic.core.CheckAEC()

	switch vic.cycle {
	case 1:
		// Fetch sprite pointer 3, increment raster counter, trigger raster IRQ,
		// test for Bad Line, reset BA if sprites 3 and 4 off, read data of sprite 3
		if vic.core.rasterY == RasterYMax {
			// Trigger VBlank in cycle 2
			vic.vBlanking = true
		} else {
			// Increment raster counter
			vic.core.UpdateRasterY()
			// Bad Line condition?
			vic.core.BadLineUpdate()
			// Don't draw all lines, hide some at the top and bottom
			vic.drawThisLine = (vic.core.rasterY >= FirstDisplayedLine) && (vic.core.rasterY <= LastDisplayedLine)
		}
		// First sample of border state
		vic.graphics.SetBorderOnSample(0)
		vic.sprites.FetchDataPtr(3)
		vic.sprites.FetchData(3, 0)
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
		}
		if (vic.sprites.GetDMAFlags() & 0x18) == 0 {
			vic.core.ClearBALow()
		}
	case 2:
		// Set BA for sprite 5, read data of sprite 3
		if vic.vBlanking {
			vBlank = true
			// Vertical blank, reset counters
			vic.graphics.ResetVideoCounterBase()
			vic.refreshCounter = 0xff
			vic.vBlanking = false
			vic.lineStart = 0
			vic.core.ResetCounters()
		}
		vic.graphics.SetLineOffset(vic.lineStart)
		vic.foreMask.Clear()
		vic.sprites.FetchData(3, 1)
		vic.sprites.FetchData(3, 2)
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
		}
		if (vic.sprites.GetDMAFlags() & 0x20) != 0 {
			vic.core.SetBALow()
		}
	case 3:
		// Fetch sprite pointer 4, reset BA is sprite 4 and 5 off
		vic.sprites.FetchDataPtr(4)
		vic.sprites.FetchData(4, 0)
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
		}
		if (vic.sprites.GetDMAFlags() & 0x30) == 0 {
			vic.core.ClearBALow()
		}
	case 4:
		// Set BA for sprite 6, read data of sprite 4
		vic.sprites.FetchData(4, 1)
		vic.sprites.FetchData(4, 2)
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
		}
		if (vic.sprites.GetDMAFlags() & 0x40) != 0 {
			vic.core.SetBALow()
		}
	case 5:
		// Fetch sprite pointer 5, reset BA if sprite 5 and 6 off
		vic.sprites.FetchDataPtr(5)
		vic.sprites.FetchData(5, 0)
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
		}
		if (vic.sprites.GetDMAFlags() & 0x60) == 0 {
			vic.core.ClearBALow()
		}
	case 6:
		// Set BA for sprite 7, read data of sprite 5
		vic.sprites.FetchData(5, 1)
		vic.sprites.FetchData(5, 2)
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
		}
		if (vic.sprites.GetDMAFlags() & 0x80) != 0 {
			vic.core.SetBALow()
		}
	case 7:
		// Fetch sprite pointer 6, reset BA if sprite 6 and 7 off
		vic.sprites.FetchDataPtr(6)
		vic.sprites.FetchData(6, 0)
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
		}
		if (vic.sprites.GetDMAFlags() & 0xc0) == 0 {
			vic.core.ClearBALow()
		}
	case 8:
		// Read data of sprite 6
		vic.sprites.FetchData(6, 1)
		vic.sprites.FetchData(6, 2)
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
		}
	case 9:
		// Fetch sprite pointer 7, reset BA if sprite 7 off
		vic.sprites.FetchDataPtr(7)
		vic.sprites.FetchData(7, 0)
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
		}
		if (vic.sprites.GetDMAFlags() & 0x80) == 0 {
			vic.core.ClearBALow()
		}
	case 10:
		// Read data of sprite 7
		vic.sprites.FetchData(7, 1)
		vic.sprites.FetchData(7, 2)
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
		}
	case 11:
		// Refresh, reset BA
		vic.accessRefresh()
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
		}
		vic.core.ClearBALow()
	case 12:
		// Refresh, turn on matrix access if Bad Line
		vic.accessRefresh()
		// Turn on display and matrix access if Bad Line
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
			vic.core.SetBALow()
		}
	case 13:
		// Refresh, turn on matrix access if Bad Line, reset rasterX, graphics display starts here
		if vic.drawThisLine {
			vic.graphics.DrawBackground()
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.accessRefresh()
		// Turn on display and matrix access if Bad Line
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
			vic.core.SetBALow()
		}
		vic.core.rasterX = 0xfffc
	case 14:
		// Refresh, videoCounter -> videoCounterBase, turn on matrix access and reset RC if Bad Line
		if vic.drawThisLine {
			vic.graphics.DrawBackground()
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.accessRefresh()
		// Turn on display and matrix access and reset RowCounter if Bad Line
		if vic.core.isBadLine {
			vic.graphics.ResetRowCounter()
			vic.graphics.SetDisplayOn()
			vic.core.SetBALow()
		}
		vic.graphics.UpdateVideoCounter()
	case 15:
		// Refresh and matrix access, increment dataCounterBase by 2 if y expansion FlipFlop is set
		if vic.drawThisLine {
			vic.graphics.DrawBackground()
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.accessRefresh()
		// Turn on display and matrix access if Bad Line
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
			vic.core.SetBALow()
		}
		vic.sprites.UpdateCounterBase()
		vic.graphics.ResetMatrixLineIndex()
		vic.core.TryAcquireAEC()
		vic.graphics.MatrixAccess()
	case 16:
		// Graphics and matrix access, increment dataCounterBase by 1 if y expansion FlipFlop is set
		// and check if sprite DMA can be turned off
		if vic.drawThisLine {
			vic.graphics.DrawBackground()
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.graphics.GraphicsAccess()
		// Turn on display and matrix access if Bad Line
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
			vic.core.SetBALow()
		}
		vic.sprites.UpdateDMACounterBase()
		vic.core.TryAcquireAEC()
		vic.graphics.MatrixAccess()
	case 17:
		// Graphics and matrix access, turn off border in 40 column mode, display window starts here
		vic.graphics.BorderUpdate()
		// Second sample of border state
		vic.graphics.SetBorderOnSample(1)
		if vic.drawThisLine {
			vic.graphics.Draw(true)
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.graphics.GraphicsAccess()
		// Turn on display and matrix access if Bad Line
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
			vic.core.SetBALow()
		}
		vic.core.TryAcquireAEC()
		vic.graphics.MatrixAccess()
	case 18:
		// Turn off border in 38 column mode
		vic.graphics.BorderUpdate2()
		// Third sample of border state
		vic.graphics.SetBorderOnSample(2)
		if vic.drawThisLine {
			vic.graphics.Draw(false)
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.graphics.GraphicsAccess()
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
			vic.core.SetBALow()
		}
		vic.core.TryAcquireAEC()
		vic.graphics.MatrixAccess()
		vic.graphics.SetCharDataLast()

	case 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54:
		if vic.drawThisLine {
			vic.graphics.Draw(false)
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.graphics.GraphicsAccess()
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
			vic.core.SetBALow()
		}
		vic.core.TryAcquireAEC()
		vic.graphics.MatrixAccess()
		vic.graphics.SetCharDataLast()
	case 55:
		// Last graphics access, turn off matrix access, turn on sprite DMA if Y coordinate is
		// right and sprite is enabled, handle sprite y expansion, set BA for sprite 0
		if vic.drawThisLine {
			vic.graphics.Draw(false)
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.graphics.GraphicsAccess()
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
		}
		vic.core.FlipFlopMYE()
		vic.sprites.UpdateDMA()
		if (vic.sprites.GetDMAFlags() & 0x01) != 0 {
			vic.core.SetBALow()
		} else {
			vic.core.ClearBALow()
		}
	case 56:
		// Turn on border in 38 column mode, turn on sprite DMA if Y coordinate is right and
		// sprite is enabled, set BA for sprite 0, display window ends here
		if (vic.core.cr2 & 8) == 0 {
			vic.graphics.SetBorderOn()
		}
		// Fourth sample of border state
		vic.graphics.SetBorderOnSample(3)
		if vic.drawThisLine {
			vic.graphics.Draw(false)
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.accessIdle()
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
		}
		vic.sprites.UpdateDMA()
		if (vic.sprites.GetDMAFlags() & 0x01) != 0 {
			vic.core.SetBALow()
		}
	case 57:
		// Turn on border in 40 column mode, set BA for sprite 1, paint sprites
		if (vic.core.cr2 & 8) != 0 {
			vic.graphics.SetBorderOn()
		}
		// Fifth sample of border state
		vic.graphics.SetBorderOnSample(4)
		vic.sprites.ApplyDisplayFlags()
		// Turn off sprite display if DMA is off
		vic.sprites.UpdateDisplayFlags()
		if vic.drawThisLine {
			vic.graphics.DrawBackground()
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.accessIdle()
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
		}
		if (vic.sprites.GetDMAFlags() & 0x02) != 0 {
			vic.core.SetBALow()
		}
	case 58:
		// Fetch sprite pointer 0, dataCounterBase->dataCounter, turn on sprite display if necessary,
		// turn off display if RC=7, read data of sprite 0
		if vic.drawThisLine {
			vic.graphics.DrawBackground()
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.sprites.UpdateRasterYDisplayFlags()
		vic.sprites.FetchDataPtr(0)
		vic.sprites.FetchData(0, 0)
		vic.graphics.UpdateDisplayOn()
	case 59:
		// Set BA for sprite 2, read data of sprite 0
		if vic.drawThisLine {
			vic.graphics.DrawBackground()
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.sprites.FetchData(0, 1)
		vic.sprites.FetchData(0, 2)
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
		}
		if (vic.sprites.GetDMAFlags() & 0x04) != 0 {
			vic.core.SetBALow()
		}
	case 60:
		// Fetch sprite pointer 1, reset BA if sprite 1 and 2 off, graphics display ends here
		if vic.drawThisLine {
			vic.graphics.DrawBackground()
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()

			vic.sprites.Draw(vic.lineStart)
			vic.graphics.DrawBorder(vic.lineStart)
			vic.lineStart += DisplayX
		}
		vic.sprites.FetchDataPtr(1)
		vic.sprites.FetchData(1, 0)
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
		}
		if (vic.sprites.GetDMAFlags() & 0x06) == 0 {
			vic.core.ClearBALow()
		}
	case 61:
		// Set BA for sprite 3, read data of sprite 1
		vic.sprites.FetchData(1, 1)
		vic.sprites.FetchData(1, 2)
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
		}
		if (vic.sprites.GetDMAFlags() & 0x08) != 0 {
			vic.core.SetBALow()
		}
	case 62:
		// Read sprite pointer 2, reset BA if sprite 2 and 3 off, read data of sprite 2
		vic.sprites.FetchDataPtr(2)
		vic.sprites.FetchData(2, 0)
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
		}
		if (vic.sprites.GetDMAFlags() & 0x0c) == 0 {
			vic.core.ClearBALow()
		}
	case 63:
		// Set BA for sprite 4, read data of sprite 2
		vic.sprites.FetchData(2, 1)
		vic.sprites.FetchData(2, 2)
		if vic.core.isBadLine {
			vic.graphics.SetDisplayOn()
		}
		vic.graphics.UpdateBorderUpperLower()
		if (vic.sprites.GetDMAFlags() & 0x10) != 0 {
			vic.core.SetBALow()
		}
		lastCycle = true
	}
	vic.core.rasterX += 8
	if lastCycle {
		vic.cycle = 1
	} else {
		vic.cycle++
	}
	return vBlank, lastCycle
}

func (vic *MOS6569) accessRefresh() {
	_ = vic.core.ReadByte(0x3f00 | uint16(vic.refreshCounter))
	vic.refreshCounter--
}

func (vic *MOS6569) accessIdle() {
	_ = vic.core.ReadByte(0x3fff)
}
