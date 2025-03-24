package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// IVIA1SocketConnections defines an interface for managing IRQ operations between connected components.
// IRQClear clears the specified interrupt request by its identifier.
// IRQTrigger triggers the specified interrupt request by its identifier.
type IVIA1SocketConnections interface {
	GetDeviceNumber() uint8
	IRQClear(uint32)
	IRQTrigger(uint32)
}

// VIA1Socket represents a socket interface for a VIA1 device, managing communication, signaling, and device configuration.
type VIA1Socket struct {
	references.IVIA
	connections IVIA1SocketConnections
	iec         references.IIec
	intrId      uint32
	dipSwitch   uint8
	prbFilter   uint8
}

// NewVIA1Socket creates and initializes a new instance of VIA1Socket with default values.
func NewVIA1Socket(connections IVIA1SocketConnections) *VIA1Socket {
	return &VIA1Socket{
		IVIA:        nil,
		connections: connections,
		iec:         nil,
		intrId:      intrIrqVIA1Bit,
		prbFilter:   0,
	}
}

// Setup initializes the VIA1Socket with the provided VIA interface, socket connections, IEC interface, and device number.
// It configures internal filters, DIP switch settings, and invokes the Setup method of the VIA interface.
func (v *VIA1Socket) Setup(c map[string]references.IComponent, _ *config.Config) error {
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
	if err = v.IVIA.Setup(v); err != nil {
		return err
	}
	return nil
}

func (v *VIA1Socket) Connect() error {
	return nil
}

// GetDeviceNumber retrieves the device number associated with the VIA1Socket instance.
func (v *VIA1Socket) GetDeviceNumber() uint8 {
	return v.connections.GetDeviceNumber()
}

// IRQClear clears the interrupt request (IRQ) signal by invoking the IRQClear method on the associated connections instance.
func (v *VIA1Socket) IRQClear() {
	v.connections.IRQClear(v.intrId)
}

// IRQTrigger triggers an interrupt request (IRQ) signal using the associated interrupt ID of the VIA1 socket.
func (v *VIA1Socket) IRQTrigger() {
	v.connections.IRQTrigger(v.intrId)
}

// setFilters configures the prbFilter register using specific bitmask settings to manage data flow and device behaviors.
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

// setDipSwitch configures the dipSwitch property based on the given deviceNumber, adjusting its bitwise values accordingly.
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

// ReadPRA reads and returns the value from Peripheral Register A (PRA) to maintain compatibility with 1541C ROMs.
func (v *VIA1Socket) ReadPRA(_ uint8, _ uint8) uint8 {
	// Keep 1541C ROMs happy (track 0 sensor)
	return 0xff
}

// ReadPRB reads and processes the Peripheral Register B value by combining it with the dip switch and PRB filter settings.
func (v *VIA1Socket) ReadPRB(prb uint8, _ uint8) uint8 {
	data := v.iec.PeripheralRead()
	p := (prb | v.dipSwitch) & v.prbFilter
	//bit 0 - 2 - 7 = 0x85
	ret := (p | data) ^ 0x85
	return ret
}

// WritePRA writes data to the Peripheral Register A (PRA) with specified mask and shift values.
func (v *VIA1Socket) WritePRA(_ uint8, _ uint8) {
}

// WritePRB writes data to the Peripheral Register B (PRB) and updates its data-direction register (DDRB).
func (v *VIA1Socket) WritePRB(prb uint8, ddrb uint8) {
	v.peripheralWrite(prb, ddrb)
}

// WriteDDRA writes a value to the Data Direction Register A (DDRA), controlling input/output modes for PRA pins.
func (v *VIA1Socket) WriteDDRA(_ uint8, _ uint8) {

}

// WriteDDRB updates the data direction register 'B' and writes to peripherals based on the provided PRB and DDRB values.
func (v *VIA1Socket) WriteDDRB(prb uint8, ddrb uint8) {
	v.peripheralWrite(prb, ddrb)
}

// peripheralWrite manages writing operations to an IEC peripheral by combining PRB, DDRB, and DIP switch configurations.
func (v *VIA1Socket) peripheralWrite(prb uint8, ddrb uint8) {
	p := prb | v.dipSwitch
	wd := (^p) & ddrb
	v.iec.PeripheralWrite(v.connections.GetDeviceNumber(), wd)
}
