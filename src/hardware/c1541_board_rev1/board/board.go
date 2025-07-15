// Package board
//
//	This package contains the implementation of the C1541 board.
//	It is responsible for the emulation of the C1541 floppy disk drive.
//	The C1541 is a peripheral device for the Commodore 64 home computer.
//	It was the standard floppy disk drive for the C64.
package board

import (
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/hardware/c1541_board_rev1/disk"
	"github.com/markel1974/c64emu/src/hardware/c1541_board_rev1/mechanic"
	"github.com/markel1974/c64emu/src/references"
	"log"
)

var _c1541hardware = []struct {
	id       string
	instance int
}{
	{"c1541_roms", 0},
	{"c1541_ram", 0},
	{"quartz", 0},
	{"c1541_pla", 0},
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
	iec        references.IIec
	cpuSocket  *CPUSocket
	via1Socket *VIA1Socket
	via2Socket *VIA2Socket
	//picSocket    *PICSocket
	plaSocket    *PLASocket
	ram          *RamSocket
	romSocket    *RomLoaderSocket
	quartzSocket *QuartzSocket
	mechanics    *mechanic.Factory
	mec          mechanic.IMechanic
	disks        *disk.Factory
	deviceId     uint8
	deviceNumber uint8
	diskId       string
	cfg          *config.Config
	emulation    []func()
	sockets      []references.ISocket
	label        string
}

// NewBoard creates and initializes a new Board with the specified IEC interface, device ID, device number, and options string.
func NewBoard(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Board {
	m := &Board{
		BaseComponent: component.NewBaseComponent(),
		deviceId:      0,
		diskId:        "",
		cfg:           nil,
		disks:         disk.NewFactory(),
		mechanics:     mechanic.NewFactory(),
		emulation:     []func(){},
		label:         label,
	}

	m.BaseComponent.Register(factory, parent, Identifier(), m, references.IdIIecDevice(m, label, instance))

	return m
}

func (m *Board) Setup() error {
	m.cfg = m.GetFactory().GetConfig()
	m.cfg.Bind(m.configChanged)
	return nil
}

func (m *Board) Bind(_ references.IIecDeviceSocket, deviceId uint8, deviceNumber uint8) error {
	iec, err := references.ComponentToIEC(m.Parent())
	if err != nil {
		return err
	}
	m.iec = iec
	m.deviceId = deviceId
	m.deviceNumber = deviceNumber
	m.mec = m.mechanics.Create("sync")

	m.romSocket = NewRomLoaderSocket(m, m.label)
	m.ram = NewRamSocket(m, m.label)
	m.quartzSocket = NewQuartzSocket(m, m.label)
	m.via1Socket = NewVIA1Socket(m, m.label, m, iec)
	m.via2Socket = NewVIA2Socket(m, m.label, m, m.mec)
	m.cpuSocket = NewCPUSocket(m, m.label, m.via2Socket)
	//m.picSocket = NewPICSocket(m, m.label)
	m.plaSocket = NewPLASocket(m, m.label)

	m.sockets = append(m.sockets, m.romSocket)
	m.sockets = append(m.sockets, m.ram)
	m.sockets = append(m.sockets, m.quartzSocket)
	m.sockets = append(m.sockets, m.cpuSocket)
	m.sockets = append(m.sockets, m.via1Socket)
	m.sockets = append(m.sockets, m.via2Socket)
	//m.sockets = append(m.sockets, m.picSocket)
	m.sockets = append(m.sockets, m.plaSocket)

	return nil
}

// Connect initializes the Board instance by configuring its components and setting up the necessary connections using the given config.
func (m *Board) Connect() error {
	var err error
	m.diskId = ""
	if err = m.mec.Setup(); err != nil {
		return err
	}

	//TODO REMOVE WHEN THREE IS READY... BEGIN
	var hardware []references.IComponent
	for _, hw := range _c1541hardware {
		comp, err := m.GetFactory().Create(m, m.label, hw.id, hw.instance)
		if err != nil {
			return err
		}
		hardware = append(hardware, comp)
	}
	for _, comp := range hardware {
		if err = comp.Setup(); err != nil {
			return err
		}
	}
	//TODO REMOVE WHEN THREE IS READY... END

	for _, c := range m.sockets {
		if err = c.Wire(); err != nil {
			return err
		}
	}

	if m.emulation, err = m.rebuildEmulation(); err != nil {
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
	//m.picSocket.Reset()
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

func (m *Board) LedActivity(led bool) {
	m.iec.LedActivity(m.deviceNumber, led)
}

// IRQClearTrigger clears the specified interrupt request (IRQ) in the programmable interrupt controller (PIC) associated with the board.
func (m *Board) IRQClearTrigger(intr uint32) {
	m.cpuSocket.ClearIRQ(intr)
}

// IRQTrigger triggers an IRQ by setting the specified interrupt bit in the programmable interrupt controller (PIC).
func (m *Board) IRQTrigger(intr uint32) {
	m.cpuSocket.TriggerIRQ(intr)
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

// rebuildEmulation constructs a sequence of emulation functions based on the given components and hardware sequence.
// Returns the constructed sequence of emulation functions or an error if the sequence is incomplete.
func (m *Board) rebuildEmulation() ([]func(), error) {
	hardwareSequence := []string{
		m.via1Socket.HardwareId(),
		m.via2Socket.HardwareId(),
		m.cpuSocket.HardwareId(),
		m.quartzSocket.HardwareId(),
	}
	var emulation []func()
	components := make(map[string]references.IComponent)
	for _, v := range m.GetChildren() {
		components[v.HardwareId()] = v
	}
	for _, x := range hardwareSequence {
		if comp, ok := components[x]; ok {
			if comp.EmulationRequired() {
				emulation = append(emulation, comp.Emulate)
			}
		}
	}
	if len(emulation) != len(hardwareSequence) {
		return nil, fmt.Errorf("emulation sequence is not complete")
	}

	if m.mec.EmulationRequired() {
		emulation = append(emulation, m.mec.Emulate)
	}

	return emulation, nil
}
