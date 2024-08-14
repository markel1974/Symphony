package c1541

import (
	"github.com/markel1974/c64emu/src/board/iec/drives/c1541/banks"
	"github.com/markel1974/c64emu/src/board/iec/drives/c1541/mechanics"
	"github.com/markel1974/c64emu/src/board/iec/virtualdrive"
	cpu2 "github.com/markel1974/c64emu/src/components/6510"
	"github.com/markel1974/c64emu/src/components/quartz"
	via3 "github.com/markel1974/c64emu/src/components/via"
	"github.com/markel1974/c64emu/src/config"
)

type Board struct {
	pic          *cpu2.Pic
	iec          virtualdrive.IIec
	quartz       *quartz.Quartz
	cpu          *cpu2.MOS6510
	via1         *via3.Via1
	via2         *via3.Via2
	banks        *banks.Banks
	mec          *mechanics.Mechanics
	deviceNumber uint8
	filePath     string
	cfg          *config.Config
}

func New(iec virtualdrive.IIec, deviceNumber uint8, opts string) *Board {
	return &Board{
		iec:          iec,
		filePath:     opts,
		deviceNumber: deviceNumber,
		pic:          nil,
		via1:         nil,
		via2:         nil,
		cpu:          nil,
		banks:        nil,
		cfg:          nil,
	}
}

func (m *Board) Setup(cfg *config.Config) {
	m.cfg = cfg
	m.cfg.Bind(m.configChanged)

	m.banks = banks.New()
	m.quartz = quartz.NewQuartz()
	m.pic = cpu2.NewPic()
	m.cpu = cpu2.NewMOS6510("c1541")
	m.mec = mechanics.NewMechanics(m.banks, m.deviceNumber)
	m.via1 = via3.NewVia1(m.iec, m.deviceNumber)
	m.via2 = via3.NewVia2(m.iec, m.mec)

	m.banks.Setup(m.via1, m.via2, cfg)
	m.pic.Setup(m.quartz)
	m.cpu.Setup(m.pic, m.banks, cfg)

	m.mec.Setup(m.filePath)

	m.via1.Setup()
	m.via1.SignalTriggerIRQBind(m.pic.TriggerIRQ)
	m.via1.SignalClearIRQBind(m.pic.ClearIRQ)

	m.via2.Setup()
	m.via2.SignalTriggerIRQBind(m.pic.TriggerIRQ)
	m.via2.SignalClearIRQBind(m.pic.ClearIRQ)
	m.cpu.SetOverflowBranch(m.via2.ByteReady)
}

func (m *Board) Reset() {
	m.pic.Reset()
	m.cpu.Reset()
	m.via1.Reset()
	m.via2.Reset()
}

func (m *Board) Emulate() {
	m.via1.CountTimers()
	m.via2.CountTimers()
	m.cpu.Emulate()
	m.quartz.AddCycle()
}

func (m *Board) Ready() bool {
	//TODO
	return true
}

func (m *Board) GetDeviceNumber() uint8 {
	return m.deviceNumber
}

func (m *Board) AtnStateChanged(b bool, b2 bool) {
	m.via1.AtnStateChanged()
	if b {
		//fmt.Println("ATN", b, "RECEIVED - WAKE UP")
		//https://sta.c64.org/cbm1541mem.html
		//Interrupt by negative edge of ATN on IEC bus
		m.banks.Write(0x7c, 1)
	}
}

func (m *Board) BusStateChanged(u uint8) {
	//nothing to do
}

func (m *Board) configChanged() {
}
