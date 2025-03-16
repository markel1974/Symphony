package iec

import (
	"fmt"
	"github.com/markel1974/c64emu/src/common/signals"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
	"strconv"
)

// BusNum defines the number of buses available in the system.
// MaxDriveSize specifies the maximum size of a drive supported.
const (
	BusNum       = 32
	MaxDriveSize = 4
)

// Dispatcher is a structure responsible for managing CPU interactions, peripherals, virtual drives, and LED signals.
type Dispatcher struct {
	*component.BaseComponent
	factory         references.IComponentFactory
	atn             bool
	cpuPort         uint8
	cpuData         uint8
	cpuBus          uint8
	peripheralsPort uint8
	peripheralsData []uint8
	virtualDrives   []references.IIecDevice
	ledSignal       *signals.Signal2[int, uint8]
}

func NewDispatcherComponent(parent references.IComponent, factory references.IComponentFactory, suffix string) references.IComponent {
	return NewDispatcher(parent, factory, suffix)
}

// NewDispatcher creates and initializes a new Dispatcher instance with the given parent component and suffix.
func NewDispatcher(parent references.IComponent, factory references.IComponentFactory, suffix string) *Dispatcher {
	c := &Dispatcher{
		BaseComponent:   component.NewBaseComponent(componentId, suffix),
		factory:         factory,
		peripheralsData: make([]uint8, BusNum),
		virtualDrives:   nil,
		ledSignal:       signals.NewSignal2[int, uint8](),
	}
	component.Register(parent, c)
	return c
}

// Setup initializes the Dispatcher by configuring the DriveFactory and setting up virtual drives based on the provided config.
func (c *Dispatcher) Setup(q references.IQuartz, cfg *config.Config) error {
	for deviceId, d := range cfg.GetDrives() {
		deviceNumber := deviceId + 8
		suffix := strconv.Itoa(deviceNumber)
		kind := "c1541"
		if len(d.Kind) > 0 {
			kind = d.Kind
		}
		device, err := c.factory.Create(c, kind, suffix)
		if err != nil {
			return err
		}
		vd, ok := device.(references.IIecDevice)
		if !ok {
			return fmt.Errorf("device %s is not an IEC device", kind)
		}
		if err = vd.Setup(c, q, uint8(deviceId), uint8(deviceNumber), d.Opts, cfg); err != nil {
			return err
		}
		c.virtualDrives = append(c.virtualDrives, vd)
	}
	return nil
}

// Emulate manages the emulation logic for all virtual drives linked to the Dispatcher, invoking their Emulate method sequentially.
func (c *Dispatcher) Emulate() {
	if len(c.virtualDrives) == 0 {
		return
	}
	if len(c.virtualDrives) == 1 {
		c.virtualDrives[0].Emulate()
		return
	}
	for _, vd := range c.virtualDrives {
		vd.Emulate()
	}
}

// Reset reinitializes all ready virtual drives managed by the Dispatcher by invoking their individual Reset methods.
func (c *Dispatcher) Reset() {
	for _, vd := range c.virtualDrives {
		if vd.Ready() {
			vd.Reset()
		}
	}
}

// buildCpuBus transforms the input data into a CPU bus value by shifting and masking specific bits.
func (c *Dispatcher) buildCpuBus(data uint8) uint8 {
	b6 := (data << 2) & 0x80
	b5 := (data << 2) & 0x40
	b4 := (data << 1) & 0x10
	value := b6 | b5 | b4
	return value
}

// buildPeripheralBus calculates the peripheral bus value based on the CPU bus and data using bitwise operations.
func (c *Dispatcher) buildPeripheralBus(cpuBus uint8, data uint8) uint8 {
	nData := ^data
	bBus := ((nData ^ cpuBus) << 3) & 0x80
	p1 := (data << 3) & 0x40
	p2 := (data << 6) & bBus
	value := p1 | p2
	return value
}

// updatePorts recalculates the CPU and peripheral ports based on the current state of the CPU bus and peripheral data.
func (c *Dispatcher) updatePorts() {
	c.cpuPort = c.cpuBus
	for _, vd := range c.virtualDrives {
		unit := vd.GetDeviceNumber()
		pData := c.peripheralsData[unit]
		pBus := c.buildPeripheralBus(c.cpuBus, pData)
		c.cpuPort &= pBus
	}
	bp7 := (c.cpuPort >> 4) & 0x04
	bp8 := c.cpuPort >> 7
	bb5 := (c.cpuBus << 3) & 0x80
	value := bp7 | bp8 | bb5
	c.peripheralsPort = value
}

// CpuWrite updates the CPU bus with the inverted data, adjusts port states, and triggers CPU write notifications.
func (c *Dispatcher) CpuWrite(data uint8) {
	//fmt.Printf("CpuWrite 0x%.2x\n", data)
	c.cpuBus = c.buildCpuBus(^data)
	//c.debugCpuWrite(^c.cpuBus)
	c.updatePorts()
	c.notifyCpuWrite()
}

// CpuRead retrieves the current value of the CPU data port from the Dispatcher instance.
func (c *Dispatcher) CpuRead() uint8 {
	return c.cpuPort
}

// PeripheralRead retrieves the current state of the peripherals' port value.
func (c *Dispatcher) PeripheralRead() uint8 {
	return c.peripheralsPort
}

func (c *Dispatcher) PeripheralWrite(deviceNumber uint8, data uint8) {
	c.peripheralsData[deviceNumber] = data
	//fmt.Printf("PeripheralWrite 0x%.2x\n", data)
	//c.debugPeripheralWrite(c.peripheralBus[deviceNumber])
	c.updatePorts()
}

// notifyCpuWrite adjusts the state of the CPU bus, notifying virtual drives of changes in ATN or bus states as necessary.
func (c *Dispatcher) notifyCpuWrite() {
	newAtn := (c.cpuBus & 0x10) != 0
	if newAtn != c.atn {
		for _, vd := range c.virtualDrives {
			vd.AtnStateChanged(newAtn)
		}
		c.atn = newAtn
	} else {
		for _, vd := range c.virtualDrives {
			vd.BusStateChanged(c.peripheralsPort)
		}
	}
}

// ledStateChangedEventHandler handles LED state change events and emits the updated state via the associated signal.
func (c *Dispatcher) ledStateChangedEventHandler(deviceNumber int, state uint8) {
	c.ledSignal.Emit(deviceNumber, state)
}

// AddPeripheral registers a peripheral device of a specified kind and options to the dispatcher using the given device ID.
func (c *Dispatcher) AddPeripheral(kind string, opts string, deviceId uint8) {
	//if c.peripheralsCount >= BusNum {
	//	return
	//}
	//for i := uint8(0); i < c.peripheralsCount; i++ {
	//	if c.peripheralStorage[i] == peripheral {
	//		return
	//	}
	//}
	//c.peripheralStorage[c.peripheralsCount] = peripheral
	//c.peripheralsCount++
	//c.rebuildPeripherals()
	//TODO
	//peripheral->LedStateChangedEvent.Bind(new SignalExecutor2<IECBus, int, uint8>(this, &IECBus::ledStateChangedEventHandler));
}

// RemovePeripheral removes a peripheral device from the dispatcher's peripheral list based on the provided device ID.
func (c *Dispatcher) RemovePeripheral(deviceId uint8) {
	//found := false
	//for i := uint8(0); i < c.peripheralsCount; i++ {
	//	if c.peripheralStorage[i] == peripheral {
	//		c.peripheralsCount--
	//		c.peripheralStorage[i] = nil
	//		found = true
	//		break
	//	}
	//}
	//if found {
	//	for i := uint8(0); i < c.peripheralsCount; i++ {
	//		for j := i + 1; j < c.peripheralsCount; j++ {
	//			if c.peripheralStorage[i].GetDeviceNumber() < c.peripheralStorage[j].GetDeviceNumber() {
	//				tmp := c.peripheralStorage[i]
	//				c.peripheralStorage[i] = c.peripheralStorage[j]
	//				c.peripheralStorage[j] = tmp
	//			}
	//		}
	//	}
	//}
	//c.rebuildPeripherals()
}
