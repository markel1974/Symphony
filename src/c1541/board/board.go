package board

import (
	"fmt"
	"github.com/markel1974/c64emu/src/c1541/mechanic"
	"github.com/markel1974/c64emu/src/c1541/pla"
	"github.com/markel1974/c64emu/src/common/signals"
	"github.com/markel1974/c64emu/src/components/6510"
	"github.com/markel1974/c64emu/src/components/iec/virtualdrive"
	"github.com/markel1974/c64emu/src/components/quartz"
	"github.com/markel1974/c64emu/src/components/via"
	"github.com/markel1974/c64emu/src/config"
)

const (
	intrIrqVIA1Bit = 0
	intrIrqVIA2Bit = 1
)

const baseId = "c1541"

type Board struct {
	pic          *mos6510.Pic
	cpu          *mos6510.CPU
	iec          virtualdrive.IIec
	quartz       *quartz.Quartz
	via1         *mos6522.Via
	via2         *mos6522.Via
	cpuSocket    *CPUSocket
	via1Socket   *Via1Socket
	via2Socket   *Via2Socket
	banks        *pla.PLA
	mec          *mechanic.Mechanic
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
		cpuSocket:    nil,
		via1Socket:   nil,
		via2Socket:   nil,
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
	m.banks = pla.New()
	m.quartz = quartz.NewQuartz()
	m.pic = mos6510.NewPic()
	m.cpu = mos6510.NewCPU(baseId)
	m.mec = mechanic.NewMechanic()
	m.mec.Setup(m.filePath)
	m.via1 = mos6522.NewVia(baseId + "_via1")
	m.via2 = mos6522.NewVia(baseId + "_via2")
	m.cpuSocket = NewCPUSocket()
	m.via1Socket = NewVia1Socket()
	m.via2Socket = NewVia2Socket()
	m.via1Socket.Setup(m, intrIrqVIA1Bit)
	m.via2Socket.Setup(m, intrIrqVIA2Bit)
	m.via1.Setup(m.via1Socket)
	m.via2.Setup(m.via2Socket)
	m.banks.Setup(m.via1, m.via2, cfg)
	m.pic.Setup(m.quartz)
	m.cpuSocket.Setup(m)
	m.cpu.Setup(m.cpuSocket)
	m.cpu.SetOverflowBranch(m.via2.ByteReady)
}

func (m *Board) Reset() {
	m.pic.Reset()
	m.cpuSocket.Reset()
	m.via1Socket.Reset()
	m.via2Socket.Reset()
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

func (m *Board) LedChanged(d byte) {
	fmt.Println("LED", m.deviceNumber, d)
	m.ledChanged.Emit(uint32(d)<<8 | uint32(m.deviceNumber))
}

func (m *Board) IRQClear(intr uint32) {
	m.pic.ClearIRQ(intr)
}

func (m *Board) IRQTrigger(intr uint32) {
	m.pic.TriggerIRQ(intr)
}
