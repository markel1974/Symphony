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
	"github.com/markel1974/c64emu/src/hardware/6510"
	"github.com/markel1974/c64emu/src/hardware/c1541/mechanic"
	"github.com/markel1974/c64emu/src/hardware/pic_6510"
	"github.com/markel1974/c64emu/src/hardware/pla_c1541"
	"github.com/markel1974/c64emu/src/hardware/roms_c1541"
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
	factory        references.IComponentFactory
	pic            references.IPic6510
	iec            references.IIec
	externalQuartz references.IQuartzSocket
	cpuSocket      *CPUSocket
	via1Socket     *Via1Socket
	via2Socket     *Via2Socket
	pla            references.IPlaC1541
	mec            *mechanic.Mechanic
	deviceId       uint8
	deviceNumber   uint8
	filePath       string
	cfg            *config.Config
	ledChanged     *signals.SignalUint32
}

func NewBoardComponent(parent references.IComponent, factory references.IComponentFactory, suffix string) references.IComponent {
	return NewBoard(parent, factory, suffix)
}

// NewBoard creates and initializes a new Board with the specified IEC interface, device ID, device number, and options string.
func NewBoard(parent references.IComponent, factory references.IComponentFactory, suffix string) *Board {
	b := &Board{
		BaseComponent:  component.NewBaseComponent("c1541", suffix),
		factory:        factory,
		externalQuartz: nil,
		iec:            nil,
		deviceId:       0,
		filePath:       "",
		deviceNumber:   0,
		ledChanged:     signals.NewSignalUint32(),
		cpuSocket:      nil,
		via1Socket:     nil,
		via2Socket:     nil,
		pic:            nil,
		pla:            nil,
		cfg:            nil,
	}
	component.Register(parent, b)
	return b
}

// Setup initializes the Board instance by configuring its components and setting up the necessary connections using the given config.
func (m *Board) Setup(iec references.IIec, quartz references.IQuartzSocket, deviceId uint8, deviceNumber uint8, opts string, cfg *config.Config) error {
	var err error
	m.iec = iec
	m.cfg = cfg
	m.deviceId = deviceId
	m.deviceNumber = deviceNumber
	m.filePath = opts
	m.externalQuartz = quartz
	m.cfg.Bind(m.configChanged)
	m.pla = pla_c1541.NewPLA(m, m.factory, "")
	//m.quartz = quartz.NewQuartz(m, "")
	m.pic = pic_6510.NewPic(m, "")
	cpu := mos6510.NewCPU(m, m.factory, "")
	m.mec = mechanic.NewMechanic()
	m.mec.Setup(m.filePath)
	var v1 references.IComponent
	if v1, err = m.factory.Create(m, "mos6522", "1"); err != nil {
		return err
	}
	via1, ok := v1.(references.IVia)
	if ok {
		return fmt.Errorf("via1 is not a Via")
	}
	var v2 references.IComponent
	if v2, err = m.factory.Create(m, "mos6522", "2"); err != nil {
		return err
	}
	via2, ok := v2.(references.IVia)
	if ok {
		return fmt.Errorf("via2 is not a Via")
	}
	m.cpuSocket = NewCPUSocket()
	m.via1Socket = NewVia1Socket()
	m.via2Socket = NewVia2Socket()
	if err = m.via1Socket.Setup(m, via1); err != nil {
		return err
	}
	if err = m.via2Socket.Setup(m, via2); err != nil {
		return err
	}
	loader := roms_c1541.NewRomLoader(cfg)
	if err = m.pla.Setup(via1, via2, loader, cfg); err != nil {
		return err
	}
	m.pic.Setup(m.externalQuartz)
	m.cpuSocket.Setup(m, cpu)
	return nil
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
	//m.quartz.AddCycle()
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
// - newAtn: Indicates whether the ATN signal is active.
func (m *Board) AtnStateChanged(newAtn bool) {
	m.via1Socket.SignalPRB()
	if !newAtn {
		//Interrupt by negative edge of ATN on IEC bus
		//https://sta.c64.org/cbm1541mem.html
		m.pla.Write(0x7c, 1)
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
