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
	"github.com/markel1974/c64emu/src/hardware/c1541/disk"
	"github.com/markel1974/c64emu/src/hardware/c1541/mechanic"
	"github.com/markel1974/c64emu/src/references"
	"log"
)

var _c1541hardware = []struct {
	id       string
	instance int
}{
	{"roms_c1541", 0},
	{"quartz", 0},
	{"pla_c1541", 0},
	{"pic_6510", 0},
	{"mos6510", 0},
	{"mos6522", 0},
	{"mos6522", 1},
}

// intrIrqVIA1Bit represents the interrupt request bit for VIA1.
// intrIrqVIA2Bit represents the interrupt request bit for VIA2.
const (
	intrIrqVIA1Bit = 0
	intrIrqVIA2Bit = 1
)

// Board represents the main hardware abstraction, containing critical components like CPU, memory, and IO devices.
type Board struct {
	*component.BaseComponent
	cpuSocket        *CPUSocket
	via1Socket       *VIA1Socket
	via2Socket       *VIA2Socket
	picSocket        *PICSocket
	plaSocket        *PLASocket
	romSocket        *RomLoaderSocket
	quartzSocket     *QuartzSocket
	mec              *mechanic.Mechanic
	disks            *disk.Factory
	deviceId         uint8
	deviceNumber     uint8
	diskId           string
	cfg              *config.Config
	ledSignal        *signals.SignalUint32
	emulation        []func()
	sockets          []references.ISocket
	hardwareSequence []string
	label            string
	//cc               map[string]references.IComponent
}

// NewBoard creates and initializes a new Board with the specified IEC interface, device ID, device number, and options string.
func NewBoard(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Board {
	m := &Board{
		BaseComponent: component.NewBaseComponent(),
		deviceId:      0,
		diskId:        "",
		ledSignal:     signals.NewSignalUint32(),
		cfg:           nil,
		disks:         disk.NewFactory(),
		mec:           mechanic.NewMechanic(),
		emulation:     []func(){},
		label:         label,
	}

	m.BaseComponent.Register(factory, parent, Identifier(), m, references.IdIIecDevice(m, label, instance))

	m.hardwareSequence = []string{
		references.IdIVIA(nil, label, 0),
		references.IdIVIA(nil, label, 1),
		references.IdI6510(nil, label, 0),
		references.IdIQuartz(nil, label, 0),
	}

	return m
}

func (m *Board) Setup(cc map[string]references.IComponent, cfg *config.Config) error {
	//m.cc = cc
	m.cfg = cfg
	m.cfg.Bind(m.configChanged)

	return nil
}

func (m *Board) Bind(_ references.IIecDeviceSocket, iecC references.IComponent, deviceId uint8, deviceNumber uint8) error {
	iec, ok := iecC.(references.IIec)
	if !ok {
		return fmt.Errorf("invalid IEC interface")
	}
	//q, ok := quartzC.(references.IQuartz)
	//if !ok {
	//	return fmt.Errorf("invalid IEC interface")
	//}
	m.deviceId = deviceId
	m.deviceNumber = deviceNumber

	m.romSocket = NewRomLoaderSocket()
	m.quartzSocket = NewQuartzSocket()
	m.cpuSocket = NewCPUSocket()
	m.via1Socket = NewVIA1Socket(m, iec)
	m.via2Socket = NewVIA2Socket(m, m.mec)
	m.picSocket = NewPICSocket()
	m.plaSocket = NewPLASocket()

	m.sockets = append(m.sockets, m.romSocket)
	m.sockets = append(m.sockets, m.quartzSocket)
	m.sockets = append(m.sockets, m.cpuSocket)
	m.sockets = append(m.sockets, m.via1Socket)
	m.sockets = append(m.sockets, m.via2Socket)
	m.sockets = append(m.sockets, m.picSocket)
	m.sockets = append(m.sockets, m.plaSocket)

	return nil
}

// Connect initializes the Board instance by configuring its components and setting up the necessary connections using the given config.
func (m *Board) Connect() error {
	var err error
	m.diskId = ""
	//quartz := quartz.NewQuartz(m, "")

	if err = m.mec.Setup(); err != nil {
		return err
	}

	//TODO REMOVE WHEN THREE IS READY...
	components := make(map[string]references.IComponent)
	var hardware []references.IComponent
	for _, hw := range _c1541hardware {
		comp, err := m.GetFactory().Create(m, m.label, hw.id, hw.instance)
		if err != nil {
			return err
		}
		components[comp.HardwareId()] = comp
		hardware = append(hardware, comp)
	}
	for _, comp := range hardware {
		if err = comp.Setup(components, m.cfg); err != nil {
			return err
		}
	}

	//components := m.cc

	for _, c := range m.sockets {
		if err = c.Mount(components, m.cfg, m.label); err != nil {
			return err
		}
	}

	if m.emulation, err = m.rebuildEmulation(components); err != nil {
		return err
	}
	if err = m.mountConfigDisk(m.cfg); err != nil {
		return err
	}
	return nil
}

func (m *Board) Internal() bool {
	return false
}

// Shutdown gracefully shuts down the Board's components, ensuring proper cleanup and resource deallocation.
func (m *Board) Shutdown() {
	//
}

// Reset reinitializes the Board's internal components to their default states by calling their respective Reset methods.
func (m *Board) Reset() {
	m.picSocket.Reset()
	m.cpuSocket.Reset()
	m.via1Socket.Reset()
	m.via2Socket.Reset()
}

func (m *Board) EmulationRequired() bool {
	return true
}

// Emulate executes one emulation cycle for the Board by simulating the VIA1, VIA2, CPU, and incrementing the quartz clock cycle.
func (m *Board) Emulate() {
	for _, fn := range m.emulation {
		fn()
	}
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
		m.plaSocket.Write(0x7c, 1)
		//fmt.Println("ATN", "RECEIVED - WAKE UP")
	}
}

func (m *Board) LEDSignal() *signals.SignalUint32 {
	return m.ledSignal
}

// IRQClear clears the specified interrupt request (IRQ) in the programmable interrupt controller (PIC) associated with the board.
func (m *Board) IRQClear(intr uint32) {
	m.picSocket.ClearIRQ(intr)
}

// IRQTrigger triggers an IRQ by setting the specified interrupt bit in the programmable interrupt controller (PIC).
func (m *Board) IRQTrigger(intr uint32) {
	m.picSocket.TriggerIRQ(intr)
}

// LEDTrigger updates the state of the LED for the device and emits the change as a signal with the device identifier and status.
func (m *Board) LEDTrigger(d byte) {
	m.ledSignal.Emit(uint32(d)<<8 | uint32(m.via1Socket.GetDeviceNumber()))
}

// InsertDisk inserts a new disk into the board's disk drive using the provided image data and write protection status.
// It resets the board, creates a disk object, and attempts to insert it into the mechanic component.
func (m *Board) InsertDisk(image []byte, writeProtected bool) error {
	m.Reset()
	g, err := m.disks.Create(image, writeProtected)
	if err != nil {
		return err
	}
	if err = m.mec.InsertDisk(g); err != nil {
		return err
	}
	return nil
}

// RemoveDisk removes the currently inserted disk from the mechanic component of the board. It returns an error if the operation fails.
func (m *Board) RemoveDisk() error {
	return m.mec.RemoveDisk()
}

// configChanged handles configuration changes by checking if the drive options have been updated and applies necessary adjustments.
func (m *Board) configChanged() {
	if err := m.mountConfigDisk(m.cfg); err != nil {
		log.Printf("can't mount disk: %s", err.Error())
	}
}

// mountConfigDisk mounts a configuration-specified disk to the Board by fetching its data and applying write protection settings.
func (m *Board) mountConfigDisk(cfg *config.Config) error {
	d := cfg.Drive(m.deviceId)
	if d == nil {
		return nil
	}
	if d.GetId() == m.diskId {
		return nil
	}
	m.diskId = d.GetId()
	if err := m.InsertDisk(d.GetData(), d.IsWriteProtected()); err != nil {
		return err
	}
	return nil
}

// rebuildEmulation constructs a sequence of emulation functions for hardware components requiring emulation.
// Returns an error if the emulation sequence is incomplete.
func (m *Board) rebuildEmulation(components map[string]references.IComponent) ([]func(), error) {
	var emulation []func()
	for _, x := range m.hardwareSequence {
		if comp, ok := components[x]; ok {
			if comp.EmulationRequired() {
				emulation = append(emulation, comp.Emulate)
			}
		}
	}
	if len(emulation) != len(m.hardwareSequence) {
		return nil, fmt.Errorf("emulation sequence is not complete")
	}
	return emulation, nil
}
