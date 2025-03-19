// Package board
//
//	This package contains the implementation of the C1541 board.
//	It is responsible for the emulation of the C1541 floppy disk drive.
//	The C1541 is a peripheral device for the Commodore 64 home computer.
//	It was the standard floppy disk drive for the C64.
package board

import (
	"fmt"
	"github.com/markel1974/c64emu/src/common/signals"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/hardware/c1541/mechanic"
	"github.com/markel1974/c64emu/src/references"
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
	*component.BaseComponent
	factory    references.IComponentFactory
	iec        references.IIec
	cpuSocket  *CPUSocket
	via1Socket *VIA1Socket
	via2Socket *VIA2Socket
	picSocket  *PICSocket
	plaSocket  *PLASocket
	rlSocket   *RomLoaderSocket
	mec        *mechanic.Mechanic
	deviceId   uint8
	filePath   string
	cfg        *config.Config
	ledChanged *signals.SignalUint32
}

// NewBoard creates and initializes a new Board with the specified IEC interface, device ID, device number, and options string.
func NewBoard(parent references.IComponent, factory references.IComponentFactory, label int) *Board {
	b := &Board{
		BaseComponent: component.NewBaseComponent(componentId, label, references.IdIBoardC1541),
		factory:       factory,
		iec:           nil,
		deviceId:      0,
		filePath:      "",
		ledChanged:    signals.NewSignalUint32(),
		cpuSocket:     NewCPUSocket(),
		via1Socket:    NewVIA1Socket(),
		via2Socket:    NewVIA2Socket(),
		picSocket:     NewPICSocket(),
		plaSocket:     NewPLASocket(),
		rlSocket:      NewRomLoaderSocket(),
		cfg:           nil,
	}
	component.Register(parent, b)
	return b
}

// Setup initializes the Board instance by configuring its components and setting up the necessary connections using the given config.
func (m *Board) Setup(iec references.IIec, quartz references.IQuartz, deviceId uint8, deviceNumber uint8, opts string, cfg *config.Config) error {
	var err error
	m.iec = iec
	m.cfg = cfg
	m.deviceId = deviceId
	m.filePath = opts
	//quartz := quartz.NewQuartz(m, "")
	m.cfg.Bind(m.configChanged)

	_, rl, err := m.factory.CreateIROMLoaderC1541(m, "roms_c1541", 0)
	if err != nil {
		return err
	}
	_, pla, err := m.factory.CreateIPLAc1541(m, "pla_c1541", 0)
	if err != nil {
		return err
	}
	_, pic, err := m.factory.CreateIPIC6510(m, "pic_6510", 0)
	if err != nil {
		return err
	}
	_, cpu, err := m.factory.CreateI6510(m, "mos6510", 0)
	if err != nil {
		return err
	}
	_, via1, err := m.factory.CreateIVIA(m, "mos6522", 0)
	if err != nil {
		return err
	}
	_, via2, err := m.factory.CreateIVIA(m, "mos6522", 1)
	if err != nil {
		return err
	}

	m.mec = mechanic.NewMechanic()
	m.mec.Setup(m.filePath)

	if err = m.rlSocket.Connect(rl, cfg); err != nil {
		return err
	}
	if err = m.plaSocket.Connect(pla, via1, via2, rl, cfg); err != nil {
		return err
	}
	if err = m.picSocket.Connect(pic, quartz); err != nil {
		return err
	}
	if err = m.via1Socket.Connect(via1, m, iec, deviceNumber); err != nil {
		return err
	}
	if err = m.via2Socket.Connect(via2, m, m.mec); err != nil {
		return err
	}
	if err = m.cpuSocket.Connect(cpu, pic, pla, via2); err != nil {
		return err
	}
	return nil
}

// Reset reinitializes the Board's internal components to their default states by calling their respective Reset methods.
func (m *Board) Reset() {
	m.picSocket.Reset()
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
	//m.quartz.AddCycle()
}

// Ready checks if the Board's internal state is prepared for operations and returns true if it is ready.
func (m *Board) Ready() bool {
	//TODO
	return true
}

// GetDeviceNumber retrieves the device number associated with the Board instance.
func (m *Board) GetDeviceNumber() uint8 {
	return m.via1Socket.GetDeviceNumber()
}

// AtnStateChanged handles changes in the ATN signal on the IEC bus and updates internal state accordingly.
// It triggers the PRB signal on VIA1 and writes to the bank memory if the ATN signal is active.
// Parameters:
// - newAtn: Indicates whether the ATN signal is active.
func (m *Board) AtnStateChanged(newAtn bool) {
	m.via1Socket.SignalPRB()
	if !newAtn {
		//Interrupt by negative edge of ATN on IEC bus
		//https://sta.c64.org/cbm1541mem.html
		m.plaSocket.Write(0x7c, 1)
		//fmt.Println("ATN", b, "RECEIVED - WAKE UP")
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
	fmt.Println("LED", m.via1Socket.GetDeviceNumber(), d)
	m.ledChanged.Emit(uint32(d)<<8 | uint32(m.via1Socket.GetDeviceNumber()))
}

// IRQClear clears the specified interrupt request (IRQ) in the programmable interrupt controller (PIC) associated with the board.
func (m *Board) IRQClear(intr uint32) {
	m.picSocket.ClearIRQ(intr)
}

// IRQTrigger triggers an IRQ by setting the specified interrupt bit in the programmable interrupt controller (PIC).
func (m *Board) IRQTrigger(intr uint32) {
	m.picSocket.TriggerIRQ(intr)
}
