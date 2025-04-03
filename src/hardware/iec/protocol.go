package iec

import (
	"github.com/markel1974/c64emu/src/common/signals"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/hardware/quartz"
	"github.com/markel1974/c64emu/src/references"
	"log"
)

//serial-iec-device.c
//static void serial_iec_device_exec_main(unsigned int devnr, CLOCK clk_value)

// defaultDDRBMask defines a bitmask for configuring data direction on DDRB, enabling specific pins for peripheral control.
const defaultDDRBMask = uint8((1 << 1) | (1 << 3) | (1 << 4))

// DeviceReadData represents the data input signal for a device.
// DeviceReadClk represents the clock input signal for a device.
// DeviceReadAtn represents the attention input signal for a device.
// DeviceWriteData represents the data output signal for a device.
// DeviceWriteClk represents the clock output signal for a device.
// DeviceWriteAtn represents the attention output signal for a device.
const (
	DeviceReadData = uint8(0x01) // DATA_IN
	DeviceReadClk  = uint8(0x04) // CLK_IN
	DeviceReadAtn  = uint8(0x80) // ATN_IN

	DeviceWriteData = uint8(0x02) // DATA_OUT
	DeviceWriteClk  = uint8(0x08) // CLK_OUT
	DeviceWriteAtn  = uint8(0x10) // ATN_A

	//DeviceWriteClk  = uint8(0x40)
	//DeviceWriteData = uint8(0x80)
)

// pPre0 to pPre2 represent pre-process states.
// pReady denotes the ready state.
// pEOI and pEOIw signify end-of-input and its waiting state.
// pBit0 to pBit7 indicate bit processing states.
// pBit0w to pBit7w represent waiting states for corresponding bit processing.
// pDone0 and pDone1 symbolize completion states.
// pFrameError0 and PFrameError1 represent error states in a frame.
const (
	pPre0 = 0 + iota
	pPre1
	pPre2
	pReady
	pEOI
	pEOIw
	pBit0
	pBit0w
	pBit1
	pBit1w
	pBit2
	pBit2w
	pBit3
	pBit3w
	pBit4
	pBit4w
	pBit5
	pBit5w
	pBit6
	pBit6w
	pBit7
	pBit7w
	pDone0
	pDone1
	pFrameError0
	PFrameError1
)

// pTalking represents a state indicating talking in progress.
// pListening represents a state indicating listening in progress.
// pAtn represents a state indicating attention state.
// pRequestListen represents a request to transition to listening state.
// pRequestTalking represents a request to transition to talking state.
// pUnlisten represents a mask to clear the listening state.
// pUntalk represents a mask to clear the talking state.
const (
	pTalking   = uint8(0x20)
	pListening = uint8(0x40)
	pAtn       = uint8(0x80)

	pRequestListen  = uint8(0x20)
	pRequestTalking = uint8(0x40)

	pUnlisten = uint8(0x3f)
	pUntalk   = uint8(0x5f)
)

// stateLast defines the last valid index for the state array, used as a mask for ensuring proper index bounds.
const stateLast = 0xf

// _pBits is a lookup table mapping protocol state machine states to corresponding bit masks for an 8-bit data line.
var _pBits [0xff]uint8

// init initializes the _pBits array with bit masks corresponding to constants pBit0 through pBit7.
func init() {
	_pBits[pBit0] = uint8(0x01)
	_pBits[pBit1] = uint8(0x02)
	_pBits[pBit2] = uint8(0x04)
	_pBits[pBit3] = uint8(0x08)
	_pBits[pBit4] = uint8(0x10)
	_pBits[pBit5] = uint8(0x20)
	_pBits[pBit6] = uint8(0x40)
	_pBits[pBit7] = uint8(0x80)
}

// Protocol represents a communication protocol structure encapsulating configurations, state, and linked components.
type Protocol struct {
	*component.BaseComponent
	iec    references.IIec
	device references.IIecProtocolDevice
	quartz references.IQuartz
	//gs            *Global
	cfg           *config.Config
	deviceNumber  uint8
	stateMachine  uint8
	flags         uint8
	primary       uint8
	secondaryPrev uint8
	secondary     uint8
	timeout       uint64
	data          uint8
	state         [stateLast + 1]uint8
	ledSignal     *signals.SignalUint32
}

// NewProtocolDevice creates a new IIecDevice instance using the provided component factory, parent, label, and instance parameters.
func NewProtocolDevice(factory references.IComponentFactory, parent references.IComponent, label string, instance int) references.IIecDevice {
	return NewProtocol(factory, parent, label, instance)
}

// NewProtocol creates a new instance of the Protocol component, initializes its fields, and registers it within the component hierarchy.
func NewProtocol(factory references.IComponentFactory, parent references.IComponent, label string, instance int) *Protocol {
	p := &Protocol{
		BaseComponent: component.NewBaseComponent(),
		//gs:            _gs,
		iec:       nil,
		device:    nil,
		ledSignal: signals.NewSignalUint32(),
		quartz:    nil,
		flags:     0,
	}
	p.BaseComponent.Register(factory, parent, "iec_device_protocol", p, references.IdIIecDevice(p, label, instance))

	p.quartz = quartz.NewQuartz(p, factory, label, 0)

	return p
}

// SetDevice assigns the specified IIecProtocolDevice to the Protocol instance for handling device-specific operations.
func (v *Protocol) SetDevice(device references.IIecProtocolDevice) {
	v.device = device
}

// Setup initializes the Protocol by configuring its settings and setting up the quartz component.
func (v *Protocol) Setup() error {
	v.cfg = v.GetFactory().GetConfig()
	if err := v.quartz.Setup(); err != nil {
		return err
	}
	return nil
}

// Bind associates a Protocol instance with a device socket, device ID, and device number, initializing the quartz and IEC components.
func (v *Protocol) Bind(_ references.IIecDeviceSocket, deviceId uint8, deviceNumber uint8) error {
	v.deviceNumber = deviceNumber
	var err error
	if err = v.quartz.Bind(v, references.IQuartz1Mhz); err != nil {
		return err
	}
	v.iec, err = references.ComponentToIEC(v.Parent())
	if err != nil {
		return err
	}
	return nil
}

// Connect establishes a connection using the protocol. Returns an error if the connection fails.
func (v *Protocol) Connect() error {
	return nil
}

// Internal checks if the protocol operates purely in internal mode without external interactions.
func (v *Protocol) Internal() bool {
	return false
}

// Shutdown gracefully stops the Protocol's operations, ensuring all resources are properly released and cleaned up.
func (v *Protocol) Shutdown() {
	//
}

// Reset resets the state of the Protocol by clearing flags, timeout, and state, and performs a peripheral write operation.
func (v *Protocol) Reset() {
	v.flags = 0
	v.timeout = 0
	for i := 0; i < len(v.state); i++ {
		v.state[i] = 0
	}
	v.peripheralWrite(false, DeviceWriteClk|DeviceWriteData)
}

// Ready returns a boolean indicating whether the protocol is ready to operate.
func (v *Protocol) Ready() bool {
	return true
}

// LEDSignal returns a signal of type SignalUint32 used to manage LED-related events or data changes.
func (v *Protocol) LEDSignal() *signals.SignalUint32 {
	return v.ledSignal
}

// GetDeviceNumber returns the device number associated with the Protocol instance. It is stored as an unsigned 8-bit integer.
func (v *Protocol) GetDeviceNumber() uint8 {
	return v.deviceNumber
}

// AtnStateChanged handles changes in the ATN (Attention) state by performing specific actions based on the boolean input.
func (v *Protocol) AtnStateChanged(atn bool) {
	//if !atn {
	//	log.Println("ATN STATE CHANGED", "ON")
	//} else {
	//	log.Println("ATN STATE CHANGED", "OFF")
	//}
}

// EmulationRequired determines if emulation is required by returning a boolean value.
func (v *Protocol) EmulationRequired() bool {
	return true
}

// Emulate processes the current state of the protocol by reading from the peripheral bus and updating the state machine.
func (v *Protocol) Emulate() {
	b := v.iec.PeripheralRead()
	busReadAtn := (b & DeviceReadAtn) != 0
	//log.Printf("protocol.Emulate(%d, %d) F=%d, S=%d, ATN=%v CLK=%v DTA=%v", v.deviceNumber, v.quartz.Cycle(), v.flags, v.state, busReadAtn, busReadClk, busReadData)
	if !v.flagGet(pAtn) && !busReadAtn {
		v.doAtnFallingFlank(busReadAtn)
	} else if v.flagGet(pAtn) && busReadAtn {
		v.doAtnRisingFlank(busReadAtn)
	}
	if v.flagGet(pAtn | pListening) {
		v.doAtnOrListen(busReadAtn, (b&DeviceReadClk) != 0, (b&DeviceReadData) != 0)
	} else if v.flagGet(pTalking) {
		v.doTalk(busReadAtn, (b&DeviceReadClk) != 0, (b&DeviceReadData) != 0)
	}
	v.quartz.Emulate()
}

// doAtnFallingFlank handles the logic triggered on the falling edge of the ATN signal during bus communication.
// It sets the protocol state machine, updates flags and primary/secondary addresses, initializes a timeout,
// and performs a peripheral write signaling presence on the bus.
func (v *Protocol) doAtnFallingFlank(busReadAtn bool) {
	//bus master addressing all devices
	v.setStateMachine(pPre0)
	v.flagsSet(pAtn)
	v.setPrimary(0)
	v.setSecondaryPrev(v.getSecondary())
	v.setSecondary(0)
	v.setTimeout(100)
	//Set DATA=0("I am here").If nobody on the bus does this within 1 ms,
	//bus-master will assume that "DeviceAdapter not present"
	v.peripheralWrite(busReadAtn, DeviceWriteClk)
}

// doAtnRisingFlank handles the state transitions of a device when the ATN signal's rising edge is detected on the bus.
// It determines the device's action (e.g., listening, talking, unlistening, or untalking) based on the received commands.
// This method updates the device's flags and state machine based on the addressing and operation requests.
// It ensures proper handling of bus communication lines for the device's operational state.
func (v *Protocol) doAtnRisingFlank(busReadAtn bool) {
	//bus master finished addressing all devices
	v.flagsRemove(pAtn)

	if (v.getPrimary() == (v.deviceNumber + pRequestListen)) || (v.getPrimary() == (v.deviceNumber + pRequestTalking)) {
		if (v.getSecondary() & 0xf0) == pTalking|pListening {
			switch v.getPrimary() & 0xf0 {
			case pRequestListen:
				v.device.Listen(v.getSecondary())
			case pRequestTalking:
				v.device.Talk(v.getSecondary())
			}
		} else if (v.getSecondary() & 0xf0) == pTalking|pListening|pAtn {
			//v.gs.SetState(0)
			state := v.device.Close(v.getSecondary())
			v.setState(v.getSecondary(), state)
			//device.setState(device.getSecondary(), v.gs.getState())
		} else if (v.getSecondary() & 0xf0) == 0xf0 {
			//v.device.Open() will not actually open the file (since we don't have a filename yet) but just set things up so that
			//the characters passed to device.
			//v.device.Write() before the next call to device.unlisten() will be interpreted as the filename.
			//The file will actually be opened during the next call to device.unlisten()
			//v.gs.SetState(0)
			state := v.device.Open(v.getSecondary())
			v.setState(v.getSecondary(), state)
			//device.setState(device.getSecondary(), v.gs.getState())
		}

		if v.getPrimary() == (v.deviceNumber + pRequestListen) {
			//We were told to listen
			v.flagsRemove(pTalking)
			//The state !=0 means that the previous OPEN command failed, i.e. we could not open a file for writing.
			//In that case, ignore the "LISTEN" request which will signal the error to the sender
			if v.getState(v.getSecondary()) == 0 {
				v.flagsSet(pListening)
				v.setStateMachine(pPre1)
				log.Printf("device %d start listening", v.deviceNumber)
			}
			//set DATA=0 ("I am here")
			v.peripheralWrite(busReadAtn, DeviceWriteClk)
		} else if v.getPrimary() == (v.deviceNumber + pRequestTalking) {
			//We were told to talk
			v.flagsRemove(pListening)
			v.flagsSet(pTalking)
			v.setStateMachine(pPre0)
			log.Printf("device %d start talking", v.deviceNumber)
		}
	} else if (v.getPrimary() == 0x3f) && v.flagGet(pListening) {
		//All devices were told to stop listening
		v.flagsRemove(pListening)
		log.Printf("device %d stop listening", v.deviceNumber)

		//If this is an UNLISTEN that followed an OPEN (0x2_ 0xf_), then
		//device.unlisten will try to open the file with the filename that
		//was received in between the OPEN and now.
		//If the file cannot be opened, it will set st != 0.
		//v.gs.SetState(v.getState(v.getSecondaryPrev()))
		v.device.Unlisten(v.getSecondaryPrev())
		//v.setState(v.getSecondaryPrev(), v.gs.GetState())
	} else if (v.getPrimary() == 0x5f) && v.flagGet(pTalking) {
		//All devices were told to stop talking
		v.device.Untalk(v.getSecondaryPrev())
		v.flagsRemove(pTalking)
		log.Printf("device %d stop talking", v.deviceNumber)
	}

	if !v.flagGet(pListening | pTalking) {
		//We're neither listening nor talking => make sure we're not holding DATA  or CLOCK line to 0 )
		v.peripheralWrite(busReadAtn, DeviceWriteClk|DeviceWriteData)
	}
}

// doAtnOrListen handles the protocol state transitions based on the bus signals during ATN or listening mode.
// It processes conditions such as clock state changes, timeouts, and received bits to update the state machine.
func (v *Protocol) doAtnOrListen(busReadAtn bool, busReadClk bool, busReadData bool) {
	//We are either under ATN or in "listening" mode
	//fmt.Println("State:", clkValue, device.getStateMachine())
	sm := v.getStateMachine()
	switch sm {
	case pPre0:
		//Ignore anything that happens during the first 100 us after falling
		//flank on ATN (other devices may have been sending and need some time to set CLK=1)
		if v.timeoutExpired() {
			//test := time.Since(P_PRE0_100_time)
			//fmt.Println("P_PRE0_100_time", test)
			v.setStateMachine(pPre1)
		}
	case pPre1:
		//Make sure CLK=0 so we actually detect a rising flank instate pPre2
		if !busReadClk {
			v.setStateMachine(pPre2)
		}
	case pPre2:
		// wait for rising flank on CLK ("ready-to-send")
		if busReadClk {
			//React by setting DATA=1 ("ready-for-data")
			v.peripheralWrite(busReadAtn, DeviceWriteClk|DeviceWriteData)
			v.setTimeout(200)
			v.setStateMachine(pReady)
		}
	case pReady:
		if !busReadClk {
			//Sender set CLK=0, is about to send first bit
			v.setStateMachine(pBit0)
		} else if !v.flagGet(pAtn) && v.timeoutExpired() {
			//Sender did not set CLK=0 within 200 us after we set DATA=1 => it is signaling EOI
			//(not so if we are under ATN) acknowledge we received it by setting DATA=0 for 60us
			log.Printf("device %d got EOI on channel %d", v.deviceNumber, v.getSecondary()&0x0f)
			v.peripheralWrite(busReadAtn, DeviceWriteClk)
			v.setStateMachine(pEOI)
			v.setTimeout(60)
		}
	case pEOI:
		if v.timeoutExpired() {
			//Set DATA back to 1 and wait for sender to set CLK=0
			v.peripheralWrite(busReadAtn, DeviceWriteClk|DeviceWriteData)
			v.setStateMachine(pEOIw)
		}
	case pEOIw:
		if !busReadClk {
			//Sender set CLK=0, is about to send first bit
			v.setStateMachine(pBit0)
		}
	case pBit0, pBit1, pBit2, pBit3, pBit4, pBit5, pBit6, pBit7:
		if busReadClk {
			//Sender set CLK=1, signaling that the DATA line represents a valid bit
			bit := _pBits[sm]
			if busReadData {
				v.dataSetBit(bit)
			} else {
				v.dataClearBit(bit)
			}
			//Go to associated P_BIT(n)w state, waiting for sender to set CLK=0
			v.setStateMachineNext()
		}
	case pBit0w, pBit1w, pBit2w, pBit3w, pBit4w, pBit5w, pBit6w:
		if !busReadClk {
			//Sender set CLK=0. go to P_BIT(n+1) state to receive the next bit
			v.setStateMachineNext()
		}
	case pBit7w:
		if !busReadClk {
			//Sender set CLK=0 and this was the last bit
			log.Printf("device %d received : 0x%02x (%c)", v.deviceNumber, v.dataGetByte(), v.dataGetByte())
			if v.flagGet(pAtn) {
				//We are currently receiving under ATN. Store the first two bytes received (contain primary and secondary address)
				if v.getPrimary() == 0 {
					v.setPrimary(v.dataGetByte())
				} else if v.getSecondary() == 0 {
					v.setSecondary(v.dataGetByte())
				}
				if (v.getPrimary() != pUnlisten) && (v.getPrimary() != pUntalk) && ((v.getPrimary() & 0x1f) != v.deviceNumber) {
					//This is NOT a UNLISTEN (0x3f) or UNTALK (0x5f) command and the primary address is not ours =>
					//Don't acknowledge the frame and stop listening. If all devices on the bus do this, the bus-master knows that "DeviceAdapter not present"
					v.setStateMachine(pDone0)
				} else {
					//Acknowledge frame by setting DATA=0
					v.peripheralWrite(busReadAtn, DeviceWriteClk)
					//repeat from pPre2 (we know that CLK=0 so no need to go to pPre1)
					v.setStateMachine(pPre2)
				}
			} else if v.flagGet(pListening) {
				//We are currently listening for data pass received byte on to the upper level
				log.Printf("device %d received 0x%02x (%c) on channel %d", v.deviceNumber, v.dataGetByte(), v.dataGetByte(), v.getSecondary())
				//v.gs.SetState(v.getState(v.getSecondary()))
				state := v.device.Write(v.getSecondary(), v.dataGetByte())
				v.setState(v.getSecondary(), state)
				//device.setState(device.getSecondary(), v.gs.getState())

				if v.getState(v.getSecondary()) != 0 {
					//There was an error during iec_bus_write => stop listening. This will signal an error condition to the sender
					v.setStateMachine(pDone0)
				} else {
					//Acknowledge frame by setting DATA=0
					v.peripheralWrite(busReadAtn, DeviceWriteClk)
					//repeat from pPre2 (we know that CLK=0 so no need to go to pPre1)
					v.setStateMachine(pPre2)
				}
			}
		}
	case pDone0:
		//fmt.Println("We're just waiting for the bus-master to set ATN back to 1")
		//We're just waiting for the bus-master to set ATN back to 1
	default:
		panic("unhandled default case")
	}
}

// doTalk handles the communication protocol state machine for a device, managing data transmission and error handling.
func (v *Protocol) doTalk(busReadAtn bool, busReadClk bool, busReadData bool) {
	sm := v.getStateMachine()
	switch sm {
	case pPre0:
		if busReadClk {
			//Bus-master set CLK=1 (and before that should have set DATA=0)
			//we are getting ready for role reversal.Set CLK=0,DATA=1
			v.peripheralWrite(busReadAtn, DeviceWriteData)
			v.setStateMachine(pPre1)
			v.setTimeout(80)
		}
	case pPre1:
		if v.timeoutExpired() {
			//Signal "ready-to-send" (CLK=1)
			v.peripheralWrite(busReadAtn, DeviceWriteClk|DeviceWriteData)
			v.setStateMachine(pReady)
		}
	case pReady:
		if busReadData {
			//Receiver signaled "ready-for-data" (DATA=1)
			//v.gs.SetState(v.getState(v.getSecondary()))
			b, state := v.device.Read(v.getSecondary() & 0xf)
			v.dataSetByte(b)
			v.setState(v.getSecondary(), state)
			//device.gs.setState(device.getSecondary(), v.gs.getState())
			if v.getState(v.getSecondary()) == 0 {
				//At least two bytes left to send. Go on to send the first bit.
				v.setStateMachine(pBit0)
				//no need to wait before sending the first bit
				v.setTimeout(0)
			} else if v.getState(v.getSecondary()) == pRequestTalking {
				//Only this byte left to send => signal EOI by keeping CLK=1
				log.Printf("device %d signaling EOI on channel %d", v.deviceNumber, v.getSecondary())
				v.setStateMachine(pEOI)
			} else {
				//There was some kind of error; we have nothing to send.
				//Just stop talking and wait for ATN (This will produce a "File not found" when loading)
				v.flagsRemove(pTalking)
			}
		}
	case pEOI:
		if !busReadData {
			//Receiver set DATA=0, first part of acknowledging the EOI
			v.setStateMachine(pEOIw)
		}
	case pEOIw:
		if busReadData {
			//Receiver set DATA=1, final part of acknowledging the EOI. Go on to send first bit
			v.setStateMachine(pBit0)
			//no need to wait before sending the first bit
			v.setTimeout(0)
		}
	case pBit0, pBit1, pBit2, pBit3, pBit4, pBit5, pBit6, pBit7:
		if v.timeoutExpired() {
			//60 us have passed since we set CLK=1 to signal "data valid" for the previous bit.
			//Pull CLK=0 and put the next bit out of DATA.
			bit := _pBits[sm]
			if v.dataHasBit(bit) {
				v.peripheralWrite(busReadAtn, DeviceWriteData)
			} else {
				v.peripheralWrite(busReadAtn, 0)
			}
			//Go to associated P_BIT(n)w state
			v.setTimeout(90) //orig 60
			v.setStateMachineNext()
		}
	case pBit0w, pBit1w, pBit2w, pBit3w, pBit4w, pBit5w, pBit6w, pBit7w:
		if v.timeoutExpired() {
			//60 us have passed since we pulled CLK=0 and put the current bit on DATA.
			//set CLK=1, keeping data as it is (this signals "data valid" to the receiver)
			if busReadData {
				v.peripheralWrite(busReadAtn, DeviceWriteClk|DeviceWriteData)
			} else {
				v.peripheralWrite(busReadAtn, DeviceWriteClk)
			}
			//Go to associated P_BIT(n+1) state to send the next bit.
			//If this was the final bit, then the next state is pDone0
			v.setTimeout(90) //orig 60
			v.setStateMachineNext()
		}
	case pDone0:
		if v.timeoutExpired() {
			//60 us have passed since we set CLK=1 to signal "data valid" for the final bit.
			//Pull CLK=0 and set DATA=1.This prepares for the receiver acknowledgement.
			v.peripheralWrite(busReadAtn, DeviceWriteData)
			v.setTimeout(1000)
			v.setStateMachine(pDone1)
		}
	case pDone1:
		if !busReadData {
			//Receiver set DATA=0, acknowledging the frame
			log.Printf("device %d sent 0x%02x (%c) on channel %d", v.deviceNumber, v.dataGetByte(), v.dataGetByte(), v.getSecondary())
			if v.getState(v.getSecondary()) == pRequestTalking {
				//This was the last byte => stop talking.This leaves us waiting for ATN.
				v.flagsRemove(pTalking)
				v.setState(v.getSecondary(), 0)
				//Release the CLOCK line to 1
				v.peripheralWrite(busReadAtn, DeviceWriteClk|DeviceWriteData)
			} else {
				//There is at least one more byte to send Start over from pPre1
				v.setTimeout(0)
				v.setStateMachine(pPre1)
			}
		} else if v.timeoutExpired() {
			//We didn't receive an acknowledgement within 1 ms.Set CLOCK=0 and after 100 us back to CLOCK=1
			log.Printf("device %d got NACK on channel %d", v.deviceNumber, v.getSecondary())
			v.peripheralWrite(busReadAtn, DeviceWriteClk|DeviceWriteData)
			v.setTimeout(100)
			v.setStateMachine(pFrameError0)
		}
	case pFrameError0:
		if v.timeoutExpired() {
			//Finished 1-0-1 sequence of CLOCK signal
			//to acknowledge the frame-error.Now wait for sender to set DATA=0 so we can continue.
			v.peripheralWrite(busReadAtn, DeviceWriteData)
			v.setStateMachine(PFrameError1)
		}
	case PFrameError1:
		if !busReadData {
			// sender set DATA=0, we can retry to send the byte
			v.setTimeout(0)
			v.setStateMachine(pPre1)
		}
	default:
		panic("unhandled default case")
	}
}

// flagsSet sets the specified flags in the Protocol's internal flags field using a bitwise OR operation.
func (v *Protocol) flagsSet(f uint8) {
	v.flags |= f
}

// flagsRemove removes specific flags from the Protocol's flags field using bitwise operations.
func (v *Protocol) flagsRemove(f uint8) {
	v.flags &= ^f
}

// flagGet checks if the specified flag is set in the Protocol's internal flags field. Returns true if set, false otherwise.
func (v *Protocol) flagGet(f uint8) bool {
	return (v.flags & f) != 0
}

// dataHasBit checks if the specified bit is set in the data field of the Protocol instance. Returns true if the bit is set.
func (v *Protocol) dataHasBit(bit uint8) bool {
	return (v.data & bit) != 0
}

// dataSetBit sets the specified bit(s) in the `data` field to 1 without altering other bits.
func (v *Protocol) dataSetBit(m uint8) {
	//m := uint8(1 << pos)
	v.data |= m
}

// dataClearBit clears specific bits in the `data` field of the Protocol struct, based on the provided mask `m`.
func (v *Protocol) dataClearBit(m uint8) {
	//m := uint8(^(1 << pos))
	n := ^m
	v.data &= n
}

// dataSetByte sets the Protocol's `data` field to the specified byte value `b`.
func (v *Protocol) dataSetByte(b uint8) {
	v.data = b
}

// dataGetByte returns the current value of the `data` field in the Protocol structure.
func (v *Protocol) dataGetByte() uint8 {
	return v.data
}

// setStateMachine sets the state machine to the specified state value.
func (v *Protocol) setStateMachine(m uint8) {
	v.stateMachine = m
}

// getStateMachine returns the current state of the protocol's state machine as an unsigned 8-bit integer.
func (v *Protocol) getStateMachine() uint8 {
	return v.stateMachine
}

// setStateMachineNext increments the stateMachine variable by one, advancing the state machine to the next state.
func (v *Protocol) setStateMachineNext() {
	v.stateMachine++
}

// setPrimary sets the primary device address to the specified value.
func (v *Protocol) setPrimary(p uint8) {
	v.primary = p
}

// getPrimary retrieves the current primary address of the Protocol. Returns the value as an unsigned 8-bit integer.
func (v *Protocol) getPrimary() uint8 {
	return v.primary
}

// setSecondary sets the value of the secondary address/state in the Protocol instance.
func (v *Protocol) setSecondary(s uint8) {
	v.secondary = s
}

// getSecondary retrieves the current value of the secondary byte in the Protocol instance.
func (v *Protocol) getSecondary() uint8 {
	return v.secondary
}

// setSecondaryPrev sets the secondary previous address value in the Protocol instance.
func (v *Protocol) setSecondaryPrev(s uint8) {
	v.secondaryPrev = s
}

// getSecondaryPrev retrieves the previous value of the secondary byte within the Protocol structure.
func (v *Protocol) getSecondaryPrev() uint8 {
	return v.secondaryPrev
}

// setTimeout sets a timeout in microseconds by calculating the required number of cycles and updating the timeout property.
func (v *Protocol) setTimeout(uSec uint64) {
	cycles := v.quartz.USecToCycleRounded(uSec)
	v.timeout = v.quartz.Cycle() + cycles
}

// timeoutExpired checks if the current cycle exceeds or equals the timeout value, indicating that the timeout has expired.
func (v *Protocol) timeoutExpired() bool {
	if b := v.quartz.Cycle(); b >= v.timeout {
		return true
	}
	return false
}

// setState updates the state at the given index after masking the index with stateLast.
func (v *Protocol) setState(idx uint8, s uint8) {
	x := idx & stateLast
	v.state[x] = s
}

// getState retrieves the state value at the given index from the state array after masking the index with stateLast.
func (v *Protocol) getState(idx uint8) uint8 {
	x := idx & stateLast
	return v.state[x]
}

// peripheralWrite writes data to a peripheral device, with optional ATN (Attention) line assertion based on busReadAtn.
func (v *Protocol) peripheralWrite(busReadAtn bool, data uint8) {
	sidecarData := references.IECSidecarEnabled | references.IECSidecarAtnAEnabled
	if busReadAtn {
		sidecarData |= uint16(references.IECAtnABit) << 8
	}
	out := sidecarData | uint16(data&defaultDDRBMask)
	v.iec.PeripheralWrite(v.deviceNumber, out)
}

// print logs the current state of the Protocol instance along with an identifier and bus number.
func (v *Protocol) print(id string, bus uint8) {
	log.Printf("%s -> bus: %d, stateMachine: %d, flags: %d, primary: %d, secondary: %d, secondaryPrev: %d\n", id, bus, v.stateMachine, v.flags, v.primary, v.secondary, v.secondaryPrev)
}

/*
// IGlobal represents an interface for managing and retrieving a state value as an unsigned 8-bit integer.
type IGlobal interface {
	GetState() uint8
	SetState() uint8
}

// Global represents a structure encapsulating the state of a global entity as an unsigned 8-bit integer.
type Global struct {
	state uint8
}

// NewGlobalState initializes and returns a pointer to a new Global instance with its state set to 0.
func NewGlobalState() *Global {
	return &Global{
		state: 0,
	}
}

// SetState updates the internal state of the Global instance with the given value.
func (v *Global) SetState(state uint8) {
	v.state = state
}

// GetState retrieves the current state value of the Global instance.
func (v *Global) GetState() uint8 {
	return v.state
}

// _gs is a singleton instance of Global, initialized using NewGlobalState, and shared across various components.
var _gs = NewGlobalState()
*/
