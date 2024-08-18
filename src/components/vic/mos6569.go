package vic

import (
	"github.com/markel1974/c64emu/src/components/quartz"
	"github.com/markel1974/c64emu/src/config"
)

//https://dustlayer.com/c64-architecture

const (
	cycleFirst = 1
	cycleLast  = 63
)

type MOS6569 struct {
	core            *Core
	cfg             *config.Config
	sprites         *Sprites
	graphics        *Graphics
	borders         *Borders
	foreMask        *ForeMask
	cycle           int   // Cycle
	lineStart       int   // Offset from current line in bitmap buffer
	drawThisLine    bool  // This line is drawn
	refreshCounter  uint8 // Refresh counter
	vBlank          bool  // VBlank in current cycle
	vBlankNextCycle bool  // VBlank in next cycle
	cycles          []func()
}

func NewMOS6569(db IDisplayBuffer) *MOS6569 {
	core := NewCore()
	foreMask := NewForeMask()
	vic := &MOS6569{
		core:            core,
		foreMask:        foreMask,
		graphics:        NewGraphics(core, foreMask, db),
		sprites:         NewSprites(core, foreMask, db),
		borders:         NewBorder(core, db),
		cycle:           cycleFirst,
		vBlank:          false,
		vBlankNextCycle: false,
		drawThisLine:    false,
		cfg:             nil,
		cycles:          make([]func(), cycleLast+1),
	}

	//TODO NTSC / PAL
	vic.cycles[0] = nil
	vic.cycles[1] = vic.cycle1
	vic.cycles[2] = vic.cycle2
	vic.cycles[3] = vic.cycle3
	vic.cycles[4] = vic.cycle4
	vic.cycles[5] = vic.cycle5
	vic.cycles[6] = vic.cycle6
	vic.cycles[7] = vic.cycle7
	vic.cycles[8] = vic.cycle8
	vic.cycles[9] = vic.cycle9
	vic.cycles[10] = vic.cycle10
	vic.cycles[11] = vic.cycle11
	vic.cycles[12] = vic.cycle12
	vic.cycles[13] = vic.cycle13
	vic.cycles[14] = vic.cycle14
	vic.cycles[15] = vic.cycle15
	vic.cycles[16] = vic.cycle16
	vic.cycles[17] = vic.cycle17
	vic.cycles[18] = vic.cycle18
	vic.cycles[19] = vic.cycle19to54
	vic.cycles[20] = vic.cycle19to54
	vic.cycles[21] = vic.cycle19to54
	vic.cycles[22] = vic.cycle19to54
	vic.cycles[23] = vic.cycle19to54
	vic.cycles[24] = vic.cycle19to54
	vic.cycles[25] = vic.cycle19to54
	vic.cycles[26] = vic.cycle19to54
	vic.cycles[27] = vic.cycle19to54
	vic.cycles[28] = vic.cycle19to54
	vic.cycles[29] = vic.cycle19to54
	vic.cycles[30] = vic.cycle19to54
	vic.cycles[31] = vic.cycle19to54
	vic.cycles[32] = vic.cycle19to54
	vic.cycles[33] = vic.cycle19to54
	vic.cycles[34] = vic.cycle19to54
	vic.cycles[35] = vic.cycle19to54
	vic.cycles[36] = vic.cycle19to54
	vic.cycles[37] = vic.cycle19to54
	vic.cycles[38] = vic.cycle19to54
	vic.cycles[39] = vic.cycle19to54
	vic.cycles[40] = vic.cycle19to54
	vic.cycles[41] = vic.cycle19to54
	vic.cycles[42] = vic.cycle19to54
	vic.cycles[43] = vic.cycle19to54
	vic.cycles[44] = vic.cycle19to54
	vic.cycles[45] = vic.cycle19to54
	vic.cycles[46] = vic.cycle19to54
	vic.cycles[47] = vic.cycle19to54
	vic.cycles[48] = vic.cycle19to54
	vic.cycles[49] = vic.cycle19to54
	vic.cycles[50] = vic.cycle19to54
	vic.cycles[51] = vic.cycle19to54
	vic.cycles[52] = vic.cycle19to54
	vic.cycles[53] = vic.cycle19to54
	vic.cycles[54] = vic.cycle19to54
	vic.cycles[55] = vic.cycle55
	vic.cycles[56] = vic.cycle56
	vic.cycles[57] = vic.cycle57
	vic.cycles[58] = vic.cycle58
	vic.cycles[59] = vic.cycle59
	vic.cycles[60] = vic.cycle60
	vic.cycles[61] = vic.cycle61
	vic.cycles[62] = vic.cycle62
	vic.cycles[63] = vic.cycle63
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

func (vic *MOS6569) GetBALow() bool {
	return vic.core.GetBALow()
}

func (vic *MOS6569) GetAECLow() bool {
	return vic.core.GetAECLow()
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

func (vic *MOS6569) accessRefresh() {
	_ = vic.core.ReadByte(0x3f00 | uint16(vic.refreshCounter))
	vic.refreshCounter--
}

func (vic *MOS6569) idleAccess() {
	_ = vic.core.ReadByte(0x3fff)
}

func (vic *MOS6569) Emulate() (bool, bool) {
	vic.cycles[vic.cycle]()
	vic.core.UpdateRasterX()
	vBlank := false
	if vic.vBlank {
		vBlank = true
		vic.vBlank = false
	}
	lastCycle := vic.cycle == cycleLast
	if lastCycle {
		vic.cycle = cycleFirst
	} else {
		vic.cycle++
	}
	return vBlank, lastCycle
}

func (vic *MOS6569) cycle1() {
	if rasterY := vic.core.GetRasterY(); rasterY == RasterYMax {
		vic.vBlankNextCycle = true
	} else {
		vic.core.IncrementCounters()
		vic.drawThisLine = (rasterY >= FirstDisplayedLine) && (rasterY <= LastDisplayedLine)
	}
	vic.borders.SetBorderOnSample(0)
	vic.sprites.FetchPtr(3)
	vic.sprites.Fetch(3, 0)
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(0x18) == 0 {
		vic.core.ClearBALow()
	}
}

func (vic *MOS6569) cycle2() {
	if vic.vBlankNextCycle {
		vic.vBlank = true
		vic.graphics.ResetVideoCounterBase()
		vic.refreshCounter = 0xff
		vic.vBlankNextCycle = false
		vic.lineStart = 0
		vic.core.ResetCounters()
	}
	vic.graphics.SetLineOffset(vic.lineStart)
	vic.foreMask.Clear()
	vic.sprites.Fetch(3, 1)
	vic.sprites.Fetch(3, 2)
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(0x20) != 0 {
		vic.core.SetBALow()
	}
}

func (vic *MOS6569) cycle3() {
	vic.sprites.FetchPtr(4)
	vic.sprites.Fetch(4, 0)
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(0x30) == 0 {
		vic.core.ClearBALow()
	}
}

func (vic *MOS6569) cycle4() {
	vic.sprites.Fetch(4, 1)
	vic.sprites.Fetch(4, 2)
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(0x40) != 0 {
		vic.core.SetBALow()
	}
}

func (vic *MOS6569) cycle5() {
	vic.sprites.FetchPtr(5)
	vic.sprites.Fetch(5, 0)
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(0x60) == 0 {
		vic.core.ClearBALow()
	}
}

func (vic *MOS6569) cycle6() {
	vic.sprites.Fetch(5, 1)
	vic.sprites.Fetch(5, 2)
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(0x80) != 0 {
		vic.core.SetBALow()
	}
}

func (vic *MOS6569) cycle7() {
	vic.sprites.FetchPtr(6)
	vic.sprites.Fetch(6, 0)
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(0xc0) == 0 {
		vic.core.ClearBALow()
	}
}

func (vic *MOS6569) cycle8() {
	vic.sprites.Fetch(6, 1)
	vic.sprites.Fetch(6, 2)
	vic.graphics.TryAcquireDisplayAccess()
}

func (vic *MOS6569) cycle9() {
	vic.sprites.FetchPtr(7)
	vic.sprites.Fetch(7, 0)
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(0x80) == 0 {
		vic.core.ClearBALow()
	}
}

func (vic *MOS6569) cycle10() {
	vic.sprites.Fetch(7, 1)
	vic.sprites.Fetch(7, 2)
	vic.graphics.TryAcquireDisplayAccess()
}

func (vic *MOS6569) cycle11() {
	vic.accessRefresh()
	vic.graphics.TryAcquireDisplayAccess()
	vic.core.ClearBALow()
}

func (vic *MOS6569) cycle12() {
	vic.accessRefresh()
	vic.graphics.TryAcquireDisplayAccess()
	vic.core.TryAcquireBA()
}

func (vic *MOS6569) cycle13() {
	if vic.drawThisLine {
		vic.graphics.DrawBackground()
		vic.borders.Sample(vic.cycle)
	}
	vic.accessRefresh()
	vic.graphics.TryAcquireDisplayAccess()
	vic.core.TryAcquireBA()
	vic.core.ResetRasterX()
}

func (vic *MOS6569) cycle14() {
	if vic.drawThisLine {
		vic.graphics.DrawBackground()
		vic.borders.Sample(vic.cycle)
	}
	vic.accessRefresh()
	vic.graphics.TryAcquireDisplayAccess()
	vic.graphics.TryResetRowCounter()
	vic.core.TryAcquireBA()
	vic.graphics.UpdateVideoCounter()
}

func (vic *MOS6569) cycle15() {
	if vic.drawThisLine {
		vic.graphics.DrawBackground()
		vic.borders.Sample(vic.cycle)
	}
	vic.accessRefresh()
	vic.graphics.TryAcquireDisplayAccess()
	vic.sprites.UpdateCounterBase()
	vic.graphics.ResetMatrixLineIndex()
	vic.core.TryAcquireBA()
	vic.core.TryAcquireAEC()
	vic.graphics.TryMatrixAccess()
}

func (vic *MOS6569) cycle16() {
	if vic.drawThisLine {
		vic.graphics.DrawBackground()
		vic.borders.Sample(vic.cycle)
	}
	vic.graphics.TryGraphicsAccess()
	vic.graphics.TryAcquireDisplayAccess()
	vic.sprites.UpdateDMACounterBase()
	vic.core.TryAcquireBA()
	vic.core.TryAcquireAEC()
	vic.graphics.TryMatrixAccess()
}

func (vic *MOS6569) cycle17() {
	if vic.core.ModeColumn40() {
		vic.borders.Update()
	}
	vic.borders.SetBorderOnSample(1)
	if vic.drawThisLine {
		if vic.borders.GetVerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
		vic.borders.Sample(vic.cycle)
	}
	vic.graphics.TryGraphicsAccess()
	vic.graphics.TryAcquireDisplayAccess()
	vic.core.TryAcquireBA()
	vic.core.TryAcquireAEC()
	vic.graphics.TryMatrixAccess()
}

func (vic *MOS6569) cycle18() {
	if vic.core.ModeColumn38() {
		vic.borders.Update()
	}
	vic.borders.SetBorderOnSample(2)
	if vic.drawThisLine {
		if vic.borders.GetVerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
		vic.borders.Sample(vic.cycle)
	}
	vic.graphics.TryGraphicsAccess()
	vic.graphics.TryAcquireDisplayAccess()
	vic.core.TryAcquireBA()
	vic.core.TryAcquireAEC()
	vic.graphics.TryMatrixAccess()
	vic.graphics.UpdateLastCharData()
}

func (vic *MOS6569) cycle19to54() {
	if vic.drawThisLine {
		if vic.borders.GetVerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
		vic.borders.Sample(vic.cycle)
	}
	vic.graphics.TryGraphicsAccess()
	vic.graphics.TryAcquireDisplayAccess()
	vic.core.TryAcquireBA()
	vic.core.TryAcquireAEC()
	vic.graphics.TryMatrixAccess()
	vic.graphics.UpdateLastCharData()
}

func (vic *MOS6569) cycle55() {
	if vic.drawThisLine {
		if vic.borders.GetVerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
		vic.borders.Sample(vic.cycle)
	}
	vic.graphics.TryGraphicsAccess()
	vic.graphics.TryAcquireDisplayAccess()
	vic.core.FlipFlopMYE()
	vic.sprites.UpdateDMA()
	if vic.sprites.GetDMAFlag(0x01) != 0 {
		vic.core.SetBALow()
	} else {
		vic.core.ClearBALow()
	}
}
func (vic *MOS6569) cycle56() {
	if vic.core.ModeColumn38() {
		vic.borders.SetBorderOn()
	}
	vic.borders.SetBorderOnSample(3)
	if vic.drawThisLine {
		if vic.borders.GetVerticalFlipFlop() {
			vic.graphics.DrawBackground()
		} else {
			vic.graphics.DrawForeground()
		}
		vic.borders.Sample(vic.cycle)
	}
	vic.idleAccess()
	vic.graphics.TryAcquireDisplayAccess()
	vic.sprites.UpdateDMA()
	if vic.sprites.GetDMAFlag(0x01) != 0 {
		vic.core.SetBALow()
	}
}

func (vic *MOS6569) cycle57() {
	if vic.core.ModeColumn40() {
		vic.borders.SetBorderOn()
	}
	vic.borders.SetBorderOnSample(4)
	vic.sprites.UpdateDisplayFlags()
	if vic.drawThisLine {
		vic.graphics.DrawBackground()
		vic.borders.Sample(vic.cycle)
	}
	vic.idleAccess()
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(0x02) != 0 {
		vic.core.SetBALow()
	}
}

func (vic *MOS6569) cycle58() {
	if vic.drawThisLine {
		vic.graphics.DrawBackground()
		vic.borders.Sample(vic.cycle)
	}
	vic.sprites.UpdateRasterYDisplayFlags()
	vic.sprites.FetchPtr(0)
	vic.sprites.Fetch(0, 0)
	vic.graphics.UpdateDisplayAccess()
}

func (vic *MOS6569) cycle59() {
	if vic.drawThisLine {
		vic.graphics.DrawBackground()
		vic.borders.Sample(vic.cycle)
	}
	vic.sprites.Fetch(0, 1)
	vic.sprites.Fetch(0, 2)
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(0x04) != 0 {
		vic.core.SetBALow()
	}
}

func (vic *MOS6569) cycle60() {
	if vic.drawThisLine {
		vic.graphics.DrawBackground()
		vic.borders.Sample(vic.cycle)
		vic.sprites.Draw(vic.lineStart)
		vic.borders.Draw(vic.lineStart)
		vic.lineStart += DisplayX
	}
	vic.sprites.FetchPtr(1)
	vic.sprites.Fetch(1, 0)
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(0x06) == 0 {
		vic.core.ClearBALow()
	}
}

func (vic *MOS6569) cycle61() {
	vic.sprites.Fetch(1, 1)
	vic.sprites.Fetch(1, 2)
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(0x08) != 0 {
		vic.core.SetBALow()
	}
}

func (vic *MOS6569) cycle62() {
	vic.sprites.FetchPtr(2)
	vic.sprites.Fetch(2, 0)
	vic.graphics.TryAcquireDisplayAccess()
	if vic.sprites.GetDMAFlag(0x0c) == 0 {
		vic.core.ClearBALow()
	}
}

func (vic *MOS6569) cycle63() {
	vic.sprites.Fetch(2, 1)
	vic.sprites.Fetch(2, 2)
	vic.graphics.TryAcquireDisplayAccess()
	vic.borders.UpdateVerticalFlipFlop()
	if vic.sprites.GetDMAFlag(0x10) != 0 {
		vic.core.SetBALow()
	}
}
