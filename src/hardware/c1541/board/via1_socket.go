package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// IVIA1SocketConnections defines an interface for managing socket-specific connections for a VIA1 adapter.
// GetDeviceNumber retrieves the device number associated with the socket.
// IRQClear clears the specified interrupt request (IRQ) using the provided bitmask.
// IRQTrigger triggers the specified interrupt request (IRQ) using the provided bitmask.
type IVIA1SocketConnections interface {
	GetDeviceNumber() uint8
	IRQClear(uint32)
	IRQTrigger(uint32)
}

// VIA1Socket represents a VIA (Versatile Interface Adapter) socket implementation with specific connection and signaling logic.
type VIA1Socket struct {
	references.IVIA
	connections IVIA1SocketConnections
	iec         references.IIec
	intrId      uint32
	dipSwitch   uint8
	prbFilter   uint8
}

// NewVIA1Socket initializes and returns a new instance of VIA1Socket with the provided IVIA1SocketConnections implementation.
func NewVIA1Socket(connections IVIA1SocketConnections) *VIA1Socket {
	return &VIA1Socket{
		IVIA:        nil,
		connections: connections,
		iec:         nil,
		intrId:      intrIrqVIA1Bit,
		prbFilter:   0,
	}
}

// Setup initializes the VIA1Socket by setting up dependencies and configurations provided in the input parameters.
func (v *VIA1Socket) Setup(c map[string]references.IComponent, cfg *config.Config) error {
	var err error
	v.IVIA, err = references.ComponentsToIVIA(c, 0)
	if err != nil {
		return err
	}
	v.iec, err = references.ComponentsToIEC(c, 0)
	if err != nil {
		return err
	}
	v.setFilters()
	v.setDipSwitch(v.connections.GetDeviceNumber())
	if err = v.IVIA.Setup(v, cfg); err != nil {
		return err
	}
	return nil
}

// Connect establishes a connection for the VIA1Socket by invoking the Connect method of its embedded IVIA component.
func (v *VIA1Socket) Connect() error {
	if err := v.IVIA.Connect(); err != nil {
		return err
	}
	return nil
}

// GetDeviceNumber retrieves the device number associated with the current VIA1 socket connection.
func (v *VIA1Socket) GetDeviceNumber() uint8 {
	return v.connections.GetDeviceNumber()
}

// IRQClear clears the interrupt request (IRQ) associated with the VIA1 socket by invoking the IRQClear method on its connections.
func (v *VIA1Socket) IRQClear() {
	v.connections.IRQClear(v.intrId)
}

// IRQTrigger triggers an interrupt request (IRQ) signal using the associated interrupt ID of the VIA1Socket connection.
func (v *VIA1Socket) IRQTrigger() {
	v.connections.IRQTrigger(v.intrId)
}

// setFilters configures the appropriate bit values for the `prbFilter` field based on specific hardware function requirements.
func (v *VIA1Socket) setFilters() {
	v.prbFilter |= 0 << 0 //Bit #0: DATA IN; 0 = Low; 1 = High.
	v.prbFilter |= 1 << 1 //Bit #1: DATA OUT; 0 = Low; 1 = High.
	v.prbFilter |= 0 << 2 //Bit #2: CLOCK IN; 0 = Low; 1 = High.
	v.prbFilter |= 1 << 3 //Bit #3: CLOCK OUT; 0 = Low; 1 = High..
	v.prbFilter |= 1 << 4 //Bit #4: ATNA OUT; 1 = Enable device presence detection by automatically acknowledging ATN IN signals on DATA OUT.
	v.prbFilter |= 1 << 5 //Bits #5 - #6: Device number, set with jumper, minus 8; % 00 = 8; % 01 = 9; % 10 = 10; % 11 = 11. Default: % 00, 8.
	v.prbFilter |= 1 << 6
	v.prbFilter |= 0 << 7 //Bit #7: ATN IN; 0 = Low; 1 = High.
}

// setDipSwitch configures the dipSwitch field based on the input device number by adjusting specific bits in the field.
func (v *VIA1Socket) setDipSwitch(deviceNumber uint8) {
	switch deviceNumber - 8 {
	case 0:
		v.dipSwitch |= 0 << 5
		v.dipSwitch |= 0 << 6
	case 1:
		v.dipSwitch |= 1 << 5
		v.dipSwitch |= 0 << 6
	case 2:
		v.dipSwitch |= 0 << 5
		v.dipSwitch |= 1 << 6
	case 3:
		v.dipSwitch |= 1 << 5
		v.dipSwitch |= 1 << 6
	default:
		v.dipSwitch |= 0 << 5
		v.dipSwitch |= 0 << 6
	}
}

// ReadPRA reads a value from Peripheral Register A (PRA) and always returns 0xff to maintain compatibility with 1541C ROMs.
func (v *VIA1Socket) ReadPRA(_ uint8, _ uint8) uint8 {
	// Keep 1541C ROMs happy (track 0 sensor)
	return 0xff
}

// ReadPRB reads the Peripheral Register B (PRB) value, applies filters and dipswitch adjustments, then modifies the result.
func (v *VIA1Socket) ReadPRB(prb uint8, _ uint8) uint8 {
	data := v.iec.PeripheralRead()
	p := (prb | v.dipSwitch) & v.prbFilter
	//bit 0 - 2 - 7 = 0x85
	ret := (p | data) ^ 0x85
	return ret
}

// WritePRA writes data to Peripheral Register A (PRA), performing any necessary updates or signaling based on input values.
func (v *VIA1Socket) WritePRA(_ uint8, _ uint8) {
}

// WritePRB writes the given values to the peripheral using prb and ddrb parameters for data and direction configuration.
func (v *VIA1Socket) WritePRB(prb uint8, ddrb uint8) {
	v.peripheralWrite(prb, ddrb)
}

// WriteDDRA updates the Data Direction Register A (DDRA) with the provided values.
func (v *VIA1Socket) WriteDDRA(_ uint8, _ uint8) {

}

// WriteDDRB performs a write operation to the Data Direction Register B (DDRB) and invokes the peripheral write method.
func (v *VIA1Socket) WriteDDRB(prb uint8, ddrb uint8) {
	v.peripheralWrite(prb, ddrb)
}

// peripheralWrite computes the write data for the peripheral bus and sends it to the corresponding device via IIec.
func (v *VIA1Socket) peripheralWrite(prb uint8, ddrb uint8) {
	p := prb | v.dipSwitch
	wd := (^p) & ddrb
	v.iec.PeripheralWrite(v.connections.GetDeviceNumber(), wd)
}
