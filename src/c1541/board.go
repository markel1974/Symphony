package c1541

import (
	"fmt"
	"github.com/markel1974/c64emu/src/c1541/banks"
	"github.com/markel1974/c64emu/src/c1541/mechanics"
	"github.com/markel1974/c64emu/src/c64/iec/virtualdrive"
	"github.com/markel1974/c64emu/src/components/6510"
	"github.com/markel1974/c64emu/src/components/quartz"
	"github.com/markel1974/c64emu/src/components/via"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/signals"
)

const (
	intrRstBit = 0
	intrNmiBit = 1
	intrIrqBit = 2
)

const (
	intrIrqVIA1Bit = 0
	intrIrqVIA2Bit = 1
)

type Board struct {
	pic          *mos6510.Pic
	cpu          *mos6510.MOS6510
	iec          virtualdrive.IIec
	quartz       *quartz.Quartz
	via1         *mos6522.Via
	via2         *mos6522.Via
	via1Wiring   *Via1Wiring
	via2Wiring   *Via2Wiring
	banks        *banks.Banks
	mec          *mechanics.Mechanics
	deviceId     uint8
	deviceNumber uint8
	filePath     string
	cfg          *config.Config
	ledChanged   *signals.SignalUint32
}

func New(iec virtualdrive.IIec, deviceId uint8, deviceNumber uint8, opts string) *Board {
	return &Board{
		iec:          iec,
		deviceId:     deviceId,
		filePath:     opts,
		deviceNumber: deviceNumber,
		ledChanged:   signals.NewSignalUint32(),
		via1Wiring:   nil,
		via2Wiring:   nil,
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
	m.pic = mos6510.NewPic(intrRstBit, intrNmiBit, intrIrqBit)
	m.cpu = mos6510.NewMOS6510("c1541")
	m.mec = mechanics.NewMechanics(m.banks, m.deviceNumber)

	m.mec.Setup(m.filePath)

	m.via1Wiring = NewVia1Wiring(m.iec, m.deviceNumber)
	m.via2Wiring = NewVia2Wiring(m.iec, m.mec, m.deviceNumber)
	m.via2Wiring.SignalLedBind(m.ledChangedSlot)

	m.via1 = mos6522.NewVia("VIA1", intrIrqVIA1Bit)
	m.via2 = mos6522.NewVia("VIA2", intrIrqVIA2Bit)

	m.via1.Setup(m.via1Wiring)
	m.via1.SignalTriggerIRQBind(m.pic.TriggerIRQ)
	m.via1.SignalClearIRQBind(m.pic.ClearIRQ)

	m.via2.Setup(m.via2Wiring)
	m.via2.SignalTriggerIRQBind(m.pic.TriggerIRQ)
	m.via2.SignalClearIRQBind(m.pic.ClearIRQ)

	m.banks.Setup(m.via1, m.via2, cfg)

	m.pic.Setup(m.quartz)
	m.cpu.Setup(m.pic, m.banks)
	m.cpu.SetOverflowBranch(m.via2.ByteReady)
}

func (m *Board) Reset() {
	m.pic.Reset()
	m.cpu.Reset()
	m.via1.Reset()
	m.via2.Reset()
}

func (m *Board) Emulate() {
	//m.mec.Emulate()
	m.via1.Emulate()
	m.via2.Emulate()
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
	m.via1.SignalPRB()
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
	if opt, ok := m.cfg.GetDrivesOpt(m.deviceId); ok {
		if opt != m.filePath {
			m.filePath = opt
			m.Reset()
			m.mec.Setup(m.filePath)
		}
	}
}

func (m *Board) ledChangedSlot(d byte) {
	fmt.Println("LED", m.deviceNumber, d)
	m.ledChanged.Emit(uint32(d)<<8 | uint32(m.deviceNumber))
}
