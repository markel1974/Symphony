package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// IVIA1SocketConnections defines an interface for managing socket-specific connections for a VIA1 adapter.
// GetDeviceNumber retrieves the device number associated with the socket.
// IRQClearTrigger clears the specified interrupt request (IRQ) using the provided bitmask.
// IRQTrigger triggers the specified interrupt request (IRQ) using the provided bitmask.
type IVIA1SocketConnections interface {
	GetDeviceNumber() uint8

	IRQClearTrigger(uint32)

	IRQTrigger(uint32)
}

// VIA1Socket represents a VIA (Versatile Interface Adapter) socket implementation with specific connection and signaling logic.
type VIA1Socket struct {
	references.IMos6522
	label       string
	parent      references.IComponent
	component   references.IComponent
	connections IVIA1SocketConnections
	iec         references.IIec
	intrId      uint32
	dipSwitch   uint8
	prbFilter   uint8
	hwId        string
}

// NewVIA1Socket initializes and returns a new instance of VIA1Socket with the provided IVIA1SocketConnections implementation.
func NewVIA1Socket(parent references.IComponent, label string, connections IVIA1SocketConnections, iec references.IIec) *VIA1Socket {
	v := &VIA1Socket{
		IMos6522:    nil,
		parent:      parent,
		label:       label,
		iec:         iec,
		connections: connections,
		intrId:      intrIrqVIA1Bit,
		prbFilter:   0,
	}
	v.hwId = references.IdIMos6522(v.IMos6522, v.label, 0)
	return v
}

func (v *VIA1Socket) HardwareId() string {
	return v.hwId
}

// Wire initializes the VIA1Socket by setting up dependencies and configurations provided in the input parameters.
func (v *VIA1Socket) Wire() error {
	var err error
	v.component = v.parent.GetChildByHardwareId(v.HardwareId())
	if v.IMos6522, err = references.ComponentToIMos6522(v.component); err != nil {
		return err
	}

	v.prbFilter = v.createPRBFilter()
	v.dipSwitch = v.createDipSwitch(v.connections.GetDeviceNumber())
	if err = v.IMos6522.Bind(v); err != nil {
		return err
	}
	return nil
}

// GetDeviceNumber retrieves the device number associated with the current VIA1 socket connection.
func (v *VIA1Socket) GetDeviceNumber() uint8 {
	return v.connections.GetDeviceNumber()
}

// IRQClearTrigger clears the interrupt request (IRQ) associated with the VIA1 socket by invoking the IRQClearTrigger method on its connections.
func (v *VIA1Socket) IRQClearTrigger() {
	v.connections.IRQClearTrigger(v.intrId)
}

// IRQTrigger triggers an interrupt request (IRQ) signal using the associated interrupt ID of the VIA1Socket connection.
func (v *VIA1Socket) IRQTrigger() {
	v.connections.IRQTrigger(v.intrId)
}

// ReadPortA reads from Port A of the VIA1 socket and always returns the value 0xff, indicating response for track 0 sensor.
func (v *VIA1Socket) ReadPortA() uint8 {
	//We ha respond 0xff for track 0 sensor
	return 0xff
}

// ReadPortB computes the value to return based on the current PRB state, dip switch configuration, and PRB filter rules.
func (v *VIA1Socket) ReadPortB() uint8 {
	prb := v.ReadPRB()
	//bit 0 - 2 - 7 = 0x85
	const bits = uint8((1 << 0) | (1 << 2) | (1 << 7))
	data := v.iec.PeripheralRead()
	p := (prb | v.dipSwitch) & v.prbFilter
	ret := (p | data) ^ bits //0x85
	return ret
}

// SignalPRA writes data to Peripheral Register A (PRA), performing any necessary updates or signaling based on input values.
func (v *VIA1Socket) SignalPRA(_ uint8) {
}

// SignalPRB writes the given values to the peripheral using prb and ddrb parameters for data and direction configuration.
func (v *VIA1Socket) SignalPRB(prb uint8) {
	ddrb := v.ReadDDRB()
	v.peripheralWrite(prb, ddrb)
}

// SignalDDRA updates the Data Direction Register A (DDRA) with the provided values.
func (v *VIA1Socket) SignalDDRA(_ uint8) {
}

// SignalDDRB performs a write operation to the Data Direction Register B (DDRB) and invokes the peripheral write method.
func (v *VIA1Socket) SignalDDRB(ddrb uint8) {
	prb := v.ReadPRB()
	v.peripheralWrite(prb, ddrb)
}

// SignalCA2 sets the state of the CA2 pin based on the provided boolean value.
func (v *VIA1Socket) SignalCA2(w bool) {
}

// SignalCB2 sets or clears the CB2 control line based on the provided boolean value.
func (v *VIA1Socket) SignalCB2(bool) {

}

// ReadCA1 checks and returns the current state of the CA1 input signal for the VIA1 socket as a boolean value.
func (v *VIA1Socket) ReadCA1() bool {
	return false
}

// ReadCB1 checks the current state of the CB1 line on the VIA1Socket and returns its value as a boolean.
func (v *VIA1Socket) ReadCB1() bool {
	return false
}

// ReadCB2 checks the state of the CB2 signal line on the VIA1 socket and returns its boolean status.
func (v *VIA1Socket) ReadCB2() bool {
	return false
}

// ReadPB6 reads the state of the PB6 pin from the VIA2Socket and returns it as a boolean value.
func (v *VIA1Socket) ReadPB6() bool {
	return false
}

// EmitPRB reads the current PRB value and signals it, ensuring proper communication with the connected peripheral.
func (v *VIA1Socket) EmitPRB() {
	prb := v.ReadPRB()
	v.SignalPRB(prb)
}

// peripheralWrite computes the write data for the peripheral bus and sends it to the corresponding device via IIec.
func (v *VIA1Socket) peripheralWrite(prb uint8, ddrb uint8) {
	p := prb | v.dipSwitch
	np := ^p
	wd := np & ddrb

	//const DeviceWriteData = 0x02 // DATA_OUT
	//const DeviceWriteClk = 0x08
	//const DeviceWriteAtn = 0x10
	//d := DeviceWriteData & wd
	//c := DeviceWriteClk & wd
	//a := DeviceWriteAtn & wd
	//fmt.Printf("c1541 transmitting 0x%x - atn: 0x%x, clock: 0x%x, data: 0x%x\n", wd, a, c, d)
	//0x1 use sidecar | use external atn a | external atn a
	const nAtnBit = ^(references.IECAtnABit)
	sidecarData := references.IECSidecarEnabled | references.IECSidecarAtnAEnabled | (uint16(wd&references.IECAtnABit) << 8)
	wd &= nAtnBit
	data := sidecarData | uint16(wd)

	v.iec.PeripheralWrite(v.connections.GetDeviceNumber(), data)
}

// createFilter configures the appropriate bit values for the `prbFilter` field based on specific hardware function requirements.
func (v *VIA1Socket) createPRBFilter() uint8 {
	prbFilter := uint8(0)
	prbFilter |= 0 << 0 //Bit #0: DATA IN; 0 = Low; 1 = High.
	prbFilter |= 1 << 1 //Bit #1: DATA OUT; 0 = Low; 1 = High.
	prbFilter |= 0 << 2 //Bit #2: CLOCK IN; 0 = Low; 1 = High.
	prbFilter |= 1 << 3 //Bit #3: CLOCK OUT; 0 = Low; 1 = High..
	prbFilter |= 1 << 4 //Bit #4: ATNA OUT; 1 = Enable device presence detection by automatically acknowledging ATN IN signals on DATA OUT.
	prbFilter |= 1 << 5 //Bits #5 - #6: Device number, set with jumper, minus 8; % 00 = 8; % 01 = 9; % 10 = 10; % 11 = 11. Default: % 00, 8.
	prbFilter |= 1 << 6
	prbFilter |= 0 << 7 //Bit #7: ATN IN; 0 = Low; 1 = High.
	return prbFilter
}

// createDipSwitch configures the dipSwitch field based on the input device number by adjusting specific bits in the field.
func (v *VIA1Socket) createDipSwitch(deviceNumber uint8) uint8 {
	dipSwitch := uint8(0)
	switch deviceNumber - 8 {
	case 0:
		dipSwitch |= 0 << 5
		dipSwitch |= 0 << 6
	case 1:
		dipSwitch |= 1 << 5
		dipSwitch |= 0 << 6
	case 2:
		dipSwitch |= 0 << 5
		dipSwitch |= 1 << 6
	case 3:
		dipSwitch |= 1 << 5
		dipSwitch |= 1 << 6
	default:
		dipSwitch |= 0 << 5
		dipSwitch |= 0 << 6
	}
	return dipSwitch
}
