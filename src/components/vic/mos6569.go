package vic

import (
	"github.com/markel1974/c64emu/src/components/quartz"
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

	switch vic.cycle {
	case 1:
		if rasterY := vic.core.GetRasterY(); rasterY == RasterYMax {
			vic.vBlanking = true
		} else {
			vic.core.IncrementCounters()
			vic.drawThisLine = (rasterY >= FirstDisplayedLine) && (rasterY <= LastDisplayedLine)
		}
		vic.graphics.SetBorderOnSample(0)
		vic.sprites.FetchDataPtr(3)
		vic.sprites.FetchData(3, 0)
		vic.graphics.TryDisplayOn()
		if (vic.sprites.GetDMAFlags() & 0x18) == 0 {
			vic.core.ClearBALow()
		}

	case 2:
		if vic.vBlanking {
			vBlank = true
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
		vic.graphics.TryDisplayOn()
		if (vic.sprites.GetDMAFlags() & 0x20) != 0 {
			vic.core.SetBALow()
		}

	case 3:
		vic.sprites.FetchDataPtr(4)
		vic.sprites.FetchData(4, 0)
		vic.graphics.TryDisplayOn()
		if (vic.sprites.GetDMAFlags() & 0x30) == 0 {
			vic.core.ClearBALow()
		}

	case 4:
		vic.sprites.FetchData(4, 1)
		vic.sprites.FetchData(4, 2)
		vic.graphics.TryDisplayOn()
		if (vic.sprites.GetDMAFlags() & 0x40) != 0 {
			vic.core.SetBALow()
		}

	case 5:
		vic.sprites.FetchDataPtr(5)
		vic.sprites.FetchData(5, 0)
		vic.graphics.TryDisplayOn()
		if (vic.sprites.GetDMAFlags() & 0x60) == 0 {
			vic.core.ClearBALow()
		}

	case 6:
		vic.sprites.FetchData(5, 1)
		vic.sprites.FetchData(5, 2)
		vic.graphics.TryDisplayOn()
		if (vic.sprites.GetDMAFlags() & 0x80) != 0 {
			vic.core.SetBALow()
		}

	case 7:
		vic.sprites.FetchDataPtr(6)
		vic.sprites.FetchData(6, 0)
		vic.graphics.TryDisplayOn()
		if (vic.sprites.GetDMAFlags() & 0xc0) == 0 {
			vic.core.ClearBALow()
		}

	case 8:
		vic.sprites.FetchData(6, 1)
		vic.sprites.FetchData(6, 2)
		vic.graphics.TryDisplayOn()

	case 9:
		vic.sprites.FetchDataPtr(7)
		vic.sprites.FetchData(7, 0)
		vic.graphics.TryDisplayOn()
		if (vic.sprites.GetDMAFlags() & 0x80) == 0 {
			vic.core.ClearBALow()
		}

	case 10:
		vic.sprites.FetchData(7, 1)
		vic.sprites.FetchData(7, 2)
		vic.graphics.TryDisplayOn()

	case 11:
		vic.accessRefresh()
		vic.graphics.TryDisplayOn()
		vic.core.ClearBALow()

	case 12:
		vic.accessRefresh()
		vic.graphics.TryDisplayOn()
		vic.core.TryBALow()

	case 13:
		if vic.drawThisLine {
			vic.graphics.DrawBackground()
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.accessRefresh()
		vic.graphics.TryDisplayOn()
		vic.core.TryBALow()
		vic.core.ResetRasterX()

	case 14:
		if vic.drawThisLine {
			vic.graphics.DrawBackground()
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.accessRefresh()
		vic.graphics.TryDisplayOn()
		vic.graphics.TryResetRowCounter()
		vic.core.TryBALow()
		vic.graphics.UpdateVideoCounter()

	case 15:
		if vic.drawThisLine {
			vic.graphics.DrawBackground()
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.accessRefresh()
		vic.graphics.TryDisplayOn()
		vic.core.TryBALow()
		vic.sprites.UpdateCounterBase()
		vic.graphics.ResetMatrixLineIndex()
		vic.graphics.TryMatrixAccess()

	case 16:
		if vic.drawThisLine {
			vic.graphics.DrawBackground()
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.graphics.GraphicsAccess()
		vic.graphics.TryDisplayOn()
		vic.core.TryBALow()
		vic.sprites.UpdateDMACounterBase()
		vic.graphics.TryMatrixAccess()

	case 17:
		if vic.core.ModeColumn40() {
			vic.graphics.BorderUpdate()
		}
		vic.graphics.SetBorderOnSample(1)
		if vic.drawThisLine {
			vic.graphics.Draw()
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.graphics.GraphicsAccess()
		vic.graphics.TryDisplayOn()
		vic.core.TryBALow()
		vic.graphics.TryMatrixAccess()

	case 18:
		if vic.core.ModeColumn38() {
			vic.graphics.BorderUpdate()
		}
		vic.graphics.SetBorderOnSample(2)
		if vic.drawThisLine {
			vic.graphics.Draw()
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.graphics.GraphicsAccess()
		vic.graphics.TryDisplayOn()
		vic.core.TryBALow()
		vic.graphics.TryMatrixAccess()
		vic.graphics.UpdateLastCharData()

	case 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54:
		if vic.drawThisLine {
			vic.graphics.Draw()
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.graphics.GraphicsAccess()
		vic.graphics.TryDisplayOn()
		vic.core.TryBALow()
		vic.graphics.TryMatrixAccess()
		vic.graphics.UpdateLastCharData()

	case 55:
		if vic.drawThisLine {
			vic.graphics.Draw()
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.graphics.GraphicsAccess()
		vic.graphics.TryDisplayOn()
		vic.core.FlipFlopMYE()
		vic.sprites.UpdateDMA()
		if (vic.sprites.GetDMAFlags() & 0x01) != 0 {
			vic.core.SetBALow()
		} else {
			vic.core.ClearBALow()
		}

	case 56:
		if vic.core.ModeColumn38() {
			vic.graphics.SetBorderOn()
		}
		vic.graphics.SetBorderOnSample(3)
		if vic.drawThisLine {
			vic.graphics.Draw()
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.idleAccess()
		vic.graphics.TryDisplayOn()
		vic.sprites.UpdateDMA()
		if (vic.sprites.GetDMAFlags() & 0x01) != 0 {
			vic.core.SetBALow()
		}

	case 57:
		if vic.core.ModeColumn40() {
			vic.graphics.SetBorderOn()
		}
		vic.graphics.SetBorderOnSample(4)
		vic.sprites.ApplyDisplayFlags()
		vic.sprites.UpdateDisplayFlags()
		if vic.drawThisLine {
			vic.graphics.DrawBackground()
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.idleAccess()
		vic.graphics.TryDisplayOn()
		if (vic.sprites.GetDMAFlags() & 0x02) != 0 {
			vic.core.SetBALow()
		}

	case 58:
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
		if vic.drawThisLine {
			vic.graphics.DrawBackground()
			vic.graphics.SampleBorder(vic.cycle)
			vic.graphics.IncrementOffset()
		}
		vic.sprites.FetchData(0, 1)
		vic.sprites.FetchData(0, 2)
		vic.graphics.TryDisplayOn()
		if (vic.sprites.GetDMAFlags() & 0x04) != 0 {
			vic.core.SetBALow()
		}

	case 60:
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
		vic.graphics.TryDisplayOn()
		if (vic.sprites.GetDMAFlags() & 0x06) == 0 {
			vic.core.ClearBALow()
		}

	case 61:
		vic.sprites.FetchData(1, 1)
		vic.sprites.FetchData(1, 2)
		vic.graphics.TryDisplayOn()
		if (vic.sprites.GetDMAFlags() & 0x08) != 0 {
			vic.core.SetBALow()
		}

	case 62:
		vic.sprites.FetchDataPtr(2)
		vic.sprites.FetchData(2, 0)
		vic.graphics.TryDisplayOn()
		if (vic.sprites.GetDMAFlags() & 0x0c) == 0 {
			vic.core.ClearBALow()
		}

	case 63:
		vic.sprites.FetchData(2, 1)
		vic.sprites.FetchData(2, 2)
		vic.graphics.TryDisplayOn()
		vic.graphics.UpdateBorderUpperLower()
		if (vic.sprites.GetDMAFlags() & 0x10) != 0 {
			vic.core.SetBALow()
		}
		lastCycle = true
	}
	vic.core.UpdateRasterX()

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

func (vic *MOS6569) idleAccess() {
	_ = vic.core.ReadByte(0x3fff)
}
