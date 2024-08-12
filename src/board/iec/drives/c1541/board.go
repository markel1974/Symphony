package c1541

import (
	"github.com/markel1974/c64emu/src/board/cpu"
	"github.com/markel1974/c64emu/src/board/iec/drives/c1541/banks"
	"github.com/markel1974/c64emu/src/board/iec/drives/c1541/mechanics"
	"github.com/markel1974/c64emu/src/board/iec/drives/c1541/via"
	"github.com/markel1974/c64emu/src/board/iec/virtualdrive"
	"github.com/markel1974/c64emu/src/board/quartz"
	"github.com/markel1974/c64emu/src/config"
)

type Board struct {
	pic          *cpu.Pic
	iec          virtualdrive.IIec
	quartz       *quartz.Quartz
	cpu          *cpu.MOS6510
	via1         *via.Via1
	via2         *via.Via2
	banks        *banks.Banks
	mec          *mechanics.Mechanics
	deviceNumber uint8
	cfg          *config.Config
}

func New(quartz *quartz.Quartz, iec virtualdrive.IIec, deviceNumber uint8) *Board {
	return &Board{
		iec:          iec,
		quartz:       quartz,
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
	m.pic = cpu.NewPic()
	m.cpu = cpu.NewMOS6510()
	m.mec = mechanics.NewMechanics(m.banks, m.deviceNumber)
	m.via1 = via.NewVia1(m.iec, m.deviceNumber)
	m.via2 = via.NewVia2(m.iec, m.mec)

	m.banks.Setup(m.via1, m.via2, cfg)
	m.pic.Setup(m.quartz)
	m.cpu.Setup(m.pic, m.banks, cfg)

	m.mec.Setup(cfg)

	m.via1.Setup()
	m.via1.SignalTriggerIRQBind(m.pic.TriggerIRQ)
	m.via1.SignalClearIRQBind(m.pic.ClearIRQ)

	m.via2.Setup()
	m.via2.SignalTriggerIRQBind(m.pic.TriggerIRQ)
	m.via2.SignalClearIRQBind(m.pic.ClearIRQ)
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
		m.banks.AtnWakeUp()
	}
}

func (m *Board) BusStateChanged(u uint8) {
	//nothing to do
}

func (m *Board) configChanged() {
}
