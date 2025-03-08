// Package board
//
//	This package contains the implementation of the C1541 board.
//	It is responsible for the emulation of the C1541 floppy disk drive.
//	The C1541 is a peripheral device for the Commodore 64 home computer.
//	It was the standard floppy disk drive for the C64.
package board

import (
	"fmt"
	"github.com/markel1974/c64emu/src/c1541/mechanic"
	"github.com/markel1974/c64emu/src/c1541/pla"
	"github.com/markel1974/c64emu/src/common/signals"
	"github.com/markel1974/c64emu/src/components/6510"
	"github.com/markel1974/c64emu/src/components/board"
	"github.com/markel1974/c64emu/src/components/iec/virtualdrive"
	"github.com/markel1974/c64emu/src/components/quartz"
	"github.com/markel1974/c64emu/src/components/via"
	"github.com/markel1974/c64emu/src/config"
)

// intrIrqVIA1Bit represents the interrupt request bit for VIA1.
// intrIrqVIA2Bit represents the interrupt request bit for VIA2.
const (
	intrIrqVIA1Bit = 0
	intrIrqVIA2Bit = 1
)

// baseId is a constant defining the base identifier used for initializing components like CPU and VIA instances.
const baseId = "c1541"

// Board represents the main hardware abstraction, containing critical components like CPU, memory, and IO devices.
type Board struct {
	*board.BaseComponent
	pic          *mos6510.Pic
	iec          virtualdrive.IIec
	quartz       *quartz.Quartz
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

// New creates and initializes a new Board with the specified IEC interface, device ID, device number, and options string.
func New(parentNode *board.Node, suffix string, iec virtualdrive.IIec, deviceId uint8, deviceNumber uint8, opts string) *Board {
	b := &Board{
		BaseComponent: board.NewBaseComponent("c1541", suffix, nil),
		iec:           iec,
		deviceId:      deviceId,
		filePath:      opts,
		deviceNumber:  deviceNumber,
		ledChanged:    signals.NewSignalUint32(),
		cpuSocket:     nil,
		via1Socket:    nil,
		via2Socket:    nil,
		pic:           nil,
		banks:         nil,
		cfg:           nil,
	}
	b.SetNode(board.CreateNode(parentNode, b))
	return b
}

// Setup initializes the Board instance by configuring its components and setting up the necessary connections using the given config.
func (m *Board) Setup(cfg *config.Config) {
	m.cfg = cfg
	m.cfg.Bind(m.configChanged)
	m.banks = pla.New()
	m.quartz = quartz.NewQuartz(m.GetNode(), "")
	m.pic = mos6510.NewPic(m.GetNode(), "")
	cpu := mos6510.NewCPU(m.GetNode(), "")
	m.mec = mechanic.NewMechanic()
	m.mec.Setup(m.filePath)
	via1 := mos6522.NewVia(m.GetNode(), "1")
	via2 := mos6522.NewVia(m.GetNode(), "2")
	m.cpuSocket = NewCPUSocket()
	m.via1Socket = NewVia1Socket()
	m.via2Socket = NewVia2Socket()
	m.via1Socket.Setup(m, via1)
	m.via2Socket.Setup(m, via2)
	m.banks.Setup(via1, via2, cfg)
	m.pic.Setup(m.quartz)
	m.cpuSocket.Setup(m, cpu)

}

// Reset reinitializes the Board's internal components to their default states by calling their respective Reset methods.
func (m *Board) Reset() {
	m.pic.Reset()
	m.cpuSocket.Reset()
	m.via1Socket.Reset()
	m.via2Socket.Reset()
}

// Emulate executes one emulation cycle for the Board by simulating the VIA1, VIA2, CPU, and incrementing the quartz clock cycle.
func (m *Board) Emulate() {
	//m.mec.Emulate()
	m.via1Socket.Emulate()
	m.via2Socket.Emulate()
	m.cpuSocket.Emulate()
	m.quartz.AddCycle()
}

// Ready checks if the Board's internal state is prepared for operations and returns true if it is ready.
func (m *Board) Ready() bool {
	//TODO
	return true
}

// GetDeviceNumber retrieves the device number associated with the Board instance.
func (m *Board) GetDeviceNumber() uint8 {
	return m.deviceNumber
}

// AtnStateChanged handles changes in the ATN signal on the IEC bus and updates internal state accordingly.
// It triggers the PRB signal on VIA1 and writes to the bank memory if the ATN signal is active.
// Parameters:
// - b: Indicates whether the ATN signal is active.
// - b2: Currently unused, reserved for future extension.
func (m *Board) AtnStateChanged(b bool, b2 bool) {
	m.via1Socket.SignalPRB()
	if b {
		//fmt.Println("ATN", b, "RECEIVED - WAKE UP")
		//https://sta.c64.org/cbm1541mem.html
		//Interrupt by negative edge of ATN on IEC bus
		m.banks.Write(0x7c, 1)
	}
}

// BusStateChanged updates the board's state in response to changes on the bus. Currently, no implementation is provided.
func (m *Board) BusStateChanged(u uint8) {
	//nothing to do
}

// configChanged handles configuration changes by checking if the drive options have been updated and applies necessary adjustments.
func (m *Board) configChanged() {
	if opt, ok := m.cfg.GetDrivesOpt(m.deviceId); ok {
		if opt != m.filePath {
			m.filePath = opt
			m.Reset()
			m.mec.Setup(m.filePath)
		}
	}
}

// LedChanged updates the state of the LED for the device and emits the change as a signal with the device identifier and status.
func (m *Board) LedChanged(d byte) {
	fmt.Println("LED", m.deviceNumber, d)
	m.ledChanged.Emit(uint32(d)<<8 | uint32(m.deviceNumber))
}

// IRQClear clears the specified interrupt request (IRQ) in the programmable interrupt controller (PIC) associated with the board.
func (m *Board) IRQClear(intr uint32) {
	m.pic.ClearIRQ(intr)
}

// IRQTrigger triggers an IRQ by setting the specified interrupt bit in the programmable interrupt controller (PIC).
func (m *Board) IRQTrigger(intr uint32) {
	m.pic.TriggerIRQ(intr)
}
