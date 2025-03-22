package iec

import (
	"fmt"
	"github.com/markel1974/c64emu/src/common/signals"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// cpuClkIn represents the IEC signal for CPU clock input (Bit 7).
// cpuDataIn represents the IEC signal for CPU data input (Bit 6).
// cpuDataOut represents the IEC signal for CPU data output (Bit 5).
// cpuClkOut represents the IEC signal for CPU clock output (Bit 4).
// cpuAtnOut represents the IEC signal for CPU attention output (Bit 3).
const (
	// IEC Signals  (CPU - input).
	cpuClkIn  = 0x80 // Bit 7: CLK (input)
	cpuDataIn = 0x40 // Bit 6: DATA (input)

	//  IEC Signals (CPU - output).
	cpuDataOut = 0x20 // Bit 5: DATA (output)
	cpuClkOut  = 0x10 // Bit 4: CLK (output)
	cpuAtnOut  = 0x08 // Bit 3: ATN (output)
)

// BusNum defines the maximum number of peripherals or devices that can be supported by the dispatcher bus system.
const (
	BusNum = 32
)

// Dispatcher represents a hardware communication interface for managing CPU and peripheral interactions.
type Dispatcher struct {
	*component.BaseComponent
	atn             bool
	cpuPort         uint8
	cpuBus          uint8
	peripheralsPort uint8
	peripheralsData []uint8
	virtualDrives   []references.IIecDevice
	ledSignal       *signals.SignalUint32 //*signals.Signal2[int, uint8]
}

// NewDispatcher initializes a new Dispatcher instance with the given parent component, factory, and instance number.
// It sets up base component properties, allocates memory for peripheral data, and initializes the LED signal.
func NewDispatcher(parent references.IComponent, factory references.IComponentFactory, instance int) *Dispatcher {
	c := &Dispatcher{
		BaseComponent:   component.NewBaseComponent(),
		peripheralsData: make([]uint8, BusNum),
		virtualDrives:   nil,
		ledSignal:       signals.NewSignalUint32(), //ledSignal:       signals.NewSignal2[int, uint8](),
	}
	c.BaseComponent.Register(factory, parent, Identifier(), instance, c, references.IdIIec(c))
	return c
}

// Setup initializes the Dispatcher component, configures drives based on the provided configuration, and prepares devices.
func (c *Dispatcher) Setup(q references.IQuartz, cfg *config.Config) error {
	for deviceId, d := range cfg.GetDrives() {
		kind := "c1541"
		if len(d.Kind) > 0 {
			kind = d.Kind
		}
		if err := c.AddPeripheral(q, cfg, kind, d.Opts, uint8(deviceId)); err != nil {
			return err
		}
		/*
			deviceNumber := deviceId + 8

			device, err := c.GetFactory().Create(c, kind, deviceNumber)
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
			vd.LEDSignal().Bind(func(state uint32) {
				c.ledSignal.Emit(state)
			})
			c.virtualDrives = append(c.virtualDrives, vd)

		*/
	}
	return nil
}

// Emulate processes the emulation logic for all virtual drives managed by the Dispatcher, calling each drive's Emulate method.
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

// Reset iterates through all virtual drives and resets those that are in a ready state.
func (c *Dispatcher) Reset() {
	for _, vd := range c.virtualDrives {
		if vd.Ready() {
			vd.Reset()
		}
	}
}

// AddPeripheral adds a new peripheral to the dispatcher with the given kind, options, and device ID.
func (c *Dispatcher) AddPeripheral(q references.IQuartz, cfg *config.Config, kind string, opts string, deviceId uint8) error {
	deviceNumber := deviceId + 8
	device, err := c.GetFactory().Create(c, kind, int(deviceNumber))
	if err != nil {
		return err
	}
	vd, ok := device.(references.IIecDevice)
	if !ok {
		return fmt.Errorf("device %s is not an IEC device", kind)
	}
	if err = vd.Setup(c, q, deviceId, deviceNumber, opts, cfg); err != nil {
		return err
	}
	vd.LEDSignal().Bind(func(state uint32) {
		c.ledSignal.Emit(state)
	})
	c.virtualDrives = append(c.virtualDrives, vd)
	//c.updatePorts()

	return nil
}

// RemovePeripheral removes a peripheral identified by the given device ID from the dispatcher.
func (c *Dispatcher) RemovePeripheral(deviceId uint8) {
	deviceNumber := deviceId + 8
	for x, v := range c.virtualDrives {
		if v.GetDeviceNumber() == deviceNumber {
			v.Shutdown()
			c.virtualDrives = append(c.virtualDrives[:x], c.virtualDrives[x+1:]...)
			//c.updatePorts()
		}
	}
}

func (c *Dispatcher) LEDSignal() *signals.SignalUint32 {
	return c.ledSignal
}

// CpuWrite updates the CPU bus and synchronizes ports and peripherals based on the provided data.
func (c *Dispatcher) CpuWrite(data uint8) {
	c.cpuBus = c.buildCpuBus(^data)
	//DebugCpuWrite(^c.cpuBus)
	c.updatePorts()
	c.notifyCpuWrite()
}

// CpuRead returns the current value of the CPU port.
func (c *Dispatcher) CpuRead() uint8 {
	return c.cpuPort
}

// PeripheralRead returns the current value from the peripherals port as a uint8.
func (c *Dispatcher) PeripheralRead() uint8 {
	return c.peripheralsPort
}

// PeripheralWrite updates the data for a specific peripheral device and triggers the port update mechanism.
func (c *Dispatcher) PeripheralWrite(deviceNumber uint8, data uint8) {
	c.peripheralsData[deviceNumber] = data
	//DebugPeripheralWrite(c.peripheralBus[deviceNumber])
	c.updatePorts()
}

// buildCpuBus processes the given data and constructs a CPU bus value by combining specific signal bits.
func (c *Dispatcher) buildCpuBus(data uint8) uint8 {
	b6 := (data << 2) & cpuClkIn
	b5 := (data << 2) & cpuDataIn
	b4 := (data << 1) & cpuClkOut
	value := b6 | b5 | b4
	return value
}

// buildPeripheralBus computes the peripheral bus signal based on the given cpuBus and data values.
// It applies bitwise operations to manipulate the input signals to generate the resulting bus configuration.
// The function utilizes specific input flags (cpuClkIn, cpuDataIn) to correctly modify and combine signals.
// Returns the newly computed peripheral bus value as uint8.
func (c *Dispatcher) buildPeripheralBus(cpuBus uint8, data uint8) uint8 {
	nData := ^data
	bBus := ((nData ^ cpuBus) << 3) & cpuClkIn
	p1 := (data << 3) & cpuDataIn
	p2 := (data << 6) & bBus
	value := p1 | p2
	return value
}

// updatePorts updates the state of the CPU and peripherals ports based on the current bus and peripheral data values.
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
	bb5 := (c.cpuBus << 3) & cpuClkIn
	value := bp7 | bp8 | bb5
	c.peripheralsPort = value
}

// notifyCpuWrite checks if the attention (ATN) signal has changed and notifies all virtual drives of the new ATN state.
func (c *Dispatcher) notifyCpuWrite() {
	if newAtn := (c.cpuBus & cpuClkOut) != 0; newAtn != c.atn {
		for _, vd := range c.virtualDrives {
			vd.AtnStateChanged(newAtn)
		}
		c.atn = newAtn
	}
}
