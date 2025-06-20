package iec_rev1

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/hardware/quartz_rev1"
	"github.com/markel1974/c64emu/src/references"
	"log"
)

//serial-iec-device.c
//static void serial_iec_device_exec_main(unsigned int devnr, CLOCK clk_value)

// defaultDDRBMask defines the default Data Direction Register B bitmask enabling bits 1, 3, and 4.
const defaultDDRBMask = uint8((1 << 1) | (1 << 3) | (1 << 4))

// DeviceReadData represents the input data signal (DATA_IN).
// DeviceReadClk represents the input clock signal (CLK_IN).
// DeviceReadAtn represents the input attention signal (ATN_IN).
// DeviceWriteData represents the output data signal (DATA_OUT).
// DeviceWriteClk represents the output clock signal (CLK_OUT).
// DeviceWriteAtn represents the output attention signal (ATN_A).
const (
	DeviceReadData  = uint8(0x01) // DATA_IN
	DeviceReadClk   = uint8(0x04) // CLK_IN
	DeviceReadAtn   = uint8(0x80) // ATN_IN
	DeviceWriteData = uint8(0x02) // DATA_OUT
	DeviceWriteClk  = uint8(0x08) // CLK_OUT
	//DeviceWriteAtn  = uint8(0x10) // ATN_A
	//DeviceWriteClk  = uint8(0x40)
	//DeviceWriteData = uint8(0x80)
)

// Constants representing various states or phases of a process using iota for sequential enumeration.
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

// pTalking represents the bitmask for the talking state in communication.
// pListening represents the bitmask for the listening state in communication.
// pAtn represents the bitmask for the attention state in communication.
// pRequestListen represents the bitmask for requesting the listening state.
// pRequestTalking represents the bitmask for requesting the talking state.
// pUnlisten represents the bitmask for clearing the listening state.
// pUntalk represents the bitmask for clearing the talking state.
const (
	pTalking   = uint8(0x20)
	pListening = uint8(0x40)
	pAtn       = uint8(0x80)

	pRequestListen  = uint8(0x20)
	pRequestTalking = uint8(0x40)

	pUnlisten = uint8(0x3f)
	pUntalk   = uint8(0x5f)
)

// stateLast defines the maximum valid index for state-related operations, ensuring indices are within bounds.
const stateLast = 0xf

// _pBits is a precomputed lookup table mapping state machine values to bitmask representations.
var _pBits [0xff]uint8

// init initializes the `_pBits` array by setting specific indices to corresponding bit values.
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

// Protocol represents a structure defining the protocol logic for managing and interacting with IEC devices.
type Protocol struct {
	*component.BaseComponent
	iec           references.IIec
	device        references.IIecProtocolDevice
	quartz        references.IQuartz
	ps            *ProtocolState
	cfg           *config.Config
	deviceNumber  uint8
	debug         bool
	atnOrListenFn []func(uint8, bool, bool, bool)
	talkFn        []func(uint8, bool, bool, bool)
}

// NewProtocol initializes and returns a new Protocol instance with the provided factory, parent, label, and instance number.
// It registers the Protocol as a component and sets up its internal state and dependencies.
func NewProtocol(factory references.IComponentFactory, parent references.IComponent, label string, instance int) *Protocol {
	p := &Protocol{
		BaseComponent: component.NewBaseComponent(),
		ps:            NewProtocolState(),
		iec:           nil,
		device:        nil,
		quartz:        nil,
		debug:         false,
	}
	p.BaseComponent.Register(factory, parent, "iec_device_protocol", p, references.IdIIecDevice(p, label, instance))
	p.quartz = quartz_rev1.NewQuartz(p, factory, label, 0)
	return p
}

// SetDevice assigns an IIecProtocolDevice to the Protocol for handling device-specific operations.
func (v *Protocol) SetDevice(device references.IIecProtocolDevice) {
	v.device = device
}

// Setup initializes the Protocol by configuring it and setting up its quartz component. Returns an error if setup fails.
func (v *Protocol) Setup() error {
	v.cfg = v.GetFactory().GetConfig()
	v.atnOrListenFn = v.setupAtnOrListen()
	v.talkFn = v.setupTalk()
	if err := v.quartz.Setup(); err != nil {
		return err
	}
	return nil
}

// Bind associates a device socket with the protocol and initializes the device number and IEC communication interface.
func (v *Protocol) Bind(_ references.IIecDeviceSocket, _ uint8, deviceNumber uint8) error {
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

// Connect establishes a connection using the configured protocol and returns an error if the connection fails.
func (v *Protocol) Connect() error {
	return nil
}

// Internal determines whether the protocol is configured for internal use and returns a boolean value.
func (v *Protocol) Internal() bool {
	return false
}

// Shutdown gracefully terminates the protocol processing and releases associated resources.
func (v *Protocol) Shutdown() {
	//
}

// Reset reinitializes the Protocol to its default state by resetting internal components and clearing communication signals.
func (v *Protocol) Reset() {
	v.ps.Reset()
	v.peripheralWrite(false, DeviceWriteClk|DeviceWriteData)
}

// Ready determines if the Protocol instance is fully initialized and prepared for execution.
func (v *Protocol) Ready() bool {
	return true
}

// LedActivity sets the LED state for a specific device based on the provided boolean value.
func (v *Protocol) LedActivity(led bool) {
	v.iec.LedActivity(v.deviceNumber, led)
}

// GetDeviceNumber retrieves the device number associated with the Protocol instance.
func (v *Protocol) GetDeviceNumber() uint8 {
	return v.deviceNumber
}

// AtnStateChanged handles changes in the ATN state by logging state transitions if debug mode is enabled.
func (v *Protocol) AtnStateChanged(atn bool) {
	if v.debug {
		if !atn {
			log.Println("ATN STATE CHANGED", "ON")
		} else {
			log.Println("ATN STATE CHANGED", "OFF")
		}
	}
}

// EmulationRequired indicates whether emulation is required for the specified protocol.
func (v *Protocol) EmulationRequired() bool {
	return true
}

// Emulate executes the core behavior of the protocol, handling state transitions, ATN flag changes, and communication logic.
func (v *Protocol) Emulate() {
	b := v.iec.PeripheralRead()
	busReadAtn := (b & DeviceReadAtn) != 0
	//log.Printf("protocol.Emulate(%d, %d) F=%d, S=%d, ATN=%v CLK=%v DTA=%v", v.deviceNumber, v.quartz.Cycle(), v.flags, v.state, busReadAtn, busReadClk, busReadData)
	if !v.ps.FlagGet(pAtn) && !busReadAtn {
		v.doAtnFallingFlank(busReadAtn)
	} else if v.ps.FlagGet(pAtn) && busReadAtn {
		v.doAtnRisingFlank(busReadAtn)
	}
	if v.ps.FlagGet(pAtn | pListening) {
		sm := v.ps.StateMachineGet()
		v.atnOrListenFn[sm](sm, busReadAtn, (b&DeviceReadClk) != 0, (b&DeviceReadData) != 0)
	} else if v.ps.FlagGet(pTalking) {
		sm := v.ps.StateMachineGet()
		v.talkFn[sm](sm, busReadAtn, (b&DeviceReadClk) != 0, (b&DeviceReadData) != 0)
	}
	v.quartz.Emulate()
}

// doAtnFallingFlank configures the state machine and flags during an ATN line falling edge condition on the bus.
// It attempts to signal device presence and updates timeout for response validation.
func (v *Protocol) doAtnFallingFlank(busReadAtn bool) {
	//bus master addressing all devices
	v.ps.StateMachineSet(pPre0)
	v.ps.FlagsSet(pAtn)
	v.ps.PrimarySet(0)
	v.ps.SecondaryPrevSet()
	v.ps.SecondarySet(0)
	v.ps.TimeoutSet(v.quartz, 100)

	//Set DATA=0("I am here").If nobody on the bus does this within 1 ms,
	//bus-master will assume that "DeviceAdapter not present"
	v.peripheralWrite(busReadAtn, DeviceWriteClk)
}

// doAtnRisingFlank handles the transition on the rising edge of the Attention (ATN) signal in a device communication system.
// It determines the subsequent operations based on the device state, addressing commands, and active flags.
// The method manages device states for listening, talking, and unlistening/untalking processes, while also updating protocol flags.
// It ensures proper handling of data and clock lines depending on the current device mode and communication status.
func (v *Protocol) doAtnRisingFlank(busReadAtn bool) {
	//bus master finished addressing all devices
	v.ps.FlagsRemove(pAtn)

	if (v.ps.PrimaryGet() == (v.deviceNumber + pRequestListen)) || (v.ps.PrimaryGet() == (v.deviceNumber + pRequestTalking)) {
		if (v.ps.SecondaryGet() & 0xf0) == pTalking|pListening {
			switch v.ps.PrimaryGet() & 0xf0 {
			case pRequestListen:
				v.ps.StateSet(v.ps.SecondaryGet(), v.device.Listen(v.ps.SecondaryGet()))
			case pRequestTalking:
				v.ps.StateSet(v.ps.SecondaryGet(), v.device.Talk(v.ps.SecondaryGet()))
			}
		} else if (v.ps.SecondaryGet() & 0xf0) == pTalking|pListening|pAtn {
			v.ps.StateSet(v.ps.SecondaryGet(), v.device.Close(v.ps.SecondaryGet()))
		} else if (v.ps.SecondaryGet() & 0xf0) == 0xf0 {
			//v.device.Open() will not actually open the file (since we don't have a filename yet) but just set things up so that
			//the characters passed to device.
			//v.device.Write() before the next call to device.unlisten() will be interpreted as the filename.
			//The file will actually be opened during the next call to device.unlisten()
			v.ps.StateSet(v.ps.SecondaryGet(), v.device.Open(v.ps.SecondaryGet()))
		}

		if v.ps.PrimaryGet() == (v.deviceNumber + pRequestListen) {
			//We were told to listen
			v.ps.FlagsRemove(pTalking)
			//The state !=0 means that the previous OPEN command failed, i.e. we could not open a file for writing.
			//In that case, ignore the "LISTEN" request which will signal the error to the sender
			if v.ps.StateGet(v.ps.SecondaryGet()) == 0 {
				v.ps.FlagsSet(pListening)
				v.ps.StateMachineSet(pPre1)
				if v.debug {
					log.Printf("device %d start listening", v.deviceNumber)
				}
			}
			//set DATA=0 ("I am here")
			v.peripheralWrite(busReadAtn, DeviceWriteClk)
		} else if v.ps.PrimaryGet() == (v.deviceNumber + pRequestTalking) {
			//We were told to talk
			v.ps.FlagsRemove(pListening)
			v.ps.FlagsSet(pTalking)
			v.ps.StateMachineSet(pPre0)
			if v.debug {
				log.Printf("device %d start talking", v.deviceNumber)
			}
		}
	} else if (v.ps.PrimaryGet() == 0x3f) && v.ps.FlagGet(pListening) {
		//All devices were told to stop listening
		v.ps.FlagsRemove(pListening)
		if v.debug {
			log.Printf("device %d stop listening", v.deviceNumber)
		}
		//If this is an UNLISTEN that followed an OPEN (0x2_ 0xf_), then
		//device.unlisten will try to open the file with the filename that
		//was received in between the OPEN and now.
		//If the file cannot be opened, it will set st != 0.
		v.ps.StateSet(v.ps.SecondaryPrevGet(), v.device.Unlisten(v.ps.SecondaryPrevGet()))
	} else if (v.ps.PrimaryGet() == 0x5f) && v.ps.FlagGet(pTalking) {
		//All devices were told to stop talking
		v.ps.StateSet(v.ps.SecondaryPrevGet(), v.device.Untalk(v.ps.SecondaryPrevGet()))
		v.ps.FlagsRemove(pTalking)
		if v.debug {
			log.Printf("device %d stop talking", v.deviceNumber)
		}
	}

	if !v.ps.FlagGet(pListening | pTalking) {
		//We're neither listening nor talking => make sure we're not holding DATA  or CLOCK line to 0 )
		v.peripheralWrite(busReadAtn, DeviceWriteClk|DeviceWriteData)
	}
}

// doAtnOrListenPre0 handles the ATN signal during the initial state, ensuring proper timing by ignoring early events.
// If the timeout has expired, it transitions the state machine to pPre1.
func (v *Protocol) doAtnOrListenPre0(_ uint8, _ bool, _ bool, _ bool) {
	//Ignore anything that happens during the first 100 us after falling
	//flank on ATN (other devices may have been sending and need some time to set CLK=1)
	if v.ps.TimeoutExpired(v.quartz) {
		//test := time.Since(P_PRE0_100_time)
		//fmt.Println("P_PRE0_100_time", test)
		v.ps.StateMachineSet(pPre1)
	}
}

// doAtnOrListenPre1 advances the state machine to pPre2 if the bus clock (CLK) is low, detecting a rising edge condition.
func (v *Protocol) doAtnOrListenPre1(_ uint8, _ bool, busReadClk bool, _ bool) {
	//Make sure CLK=0 so we actually detect a rising flank instate pPre2
	if !busReadClk {
		v.ps.StateMachineSet(pPre2)
	}
}

// doAtnOrListenPre2 reacts to changes in the CLK signal, setting DATA high and updating the state machine accordingly.
func (v *Protocol) doAtnOrListenPre2(_ uint8, busReadAtn bool, busReadClk bool, _ bool) {
	// wait for rising flank on CLK ("ready-to-send")
	if busReadClk {
		//React by setting DATA=1 ("ready-for-data")
		v.peripheralWrite(busReadAtn, DeviceWriteClk|DeviceWriteData)
		v.ps.TimeoutSet(v.quartz, 200)
		v.ps.StateMachineSet(pReady)
	}
}

// doAtnOrListenReady manages protocol state transitions during ATN or listen ready phase based on bus signals and timeouts.
// sm indicates current state; busReadAtn, busReadClk, and busReadData represent bus signal states.
// Sets state to pBit0 if CLK=0; handles EOI signal acknowledgment with timeout if ATN is not active.
// Updates the state machine and logs debug info if enabled, addressing channel EOI handling.
func (v *Protocol) doAtnOrListenReady(_ uint8, busReadAtn bool, busReadClk bool, _ bool) {
	if !busReadClk {
		//Sender set CLK=0, is about to send first bit
		v.ps.StateMachineSet(pBit0)
	} else if !v.ps.FlagGet(pAtn) && v.ps.TimeoutExpired(v.quartz) {
		//Sender did not set CLK=0 within 200 us after we set DATA=1 => it is signaling EOI
		//(not so if we are under ATN) acknowledge we received it by setting DATA=0 for 60us
		if v.debug {
			log.Printf("device %d got EOI on channel %d", v.deviceNumber, v.ps.SecondaryGet()&0x0f)
		}
		v.peripheralWrite(busReadAtn, DeviceWriteClk)
		v.ps.StateMachineSet(pEOI)
		v.ps.TimeoutSet(v.quartz, 60)
	}
}

// doAtnOrListenEOI manages the protocol state when attention or EOI (End-Or-Identify) conditions are encountered.
// It sets the DATA line high and transitions the state machine when the timeout has expired.
func (v *Protocol) doAtnOrListenEOI(_ uint8, busReadAtn bool, _ bool, _ bool) {
	if v.ps.TimeoutExpired(v.quartz) {
		//Set DATA back to 1 and wait for sender to set CLK=0
		v.peripheralWrite(busReadAtn, DeviceWriteClk|DeviceWriteData)
		v.ps.StateMachineSet(pEOIw)
	}
}

// doAtnOrListenEOIw handles the state transition to pBit0 if the clock signal (busReadClk) is low.
// This is used when the sender is preparing to send the first bit.
func (v *Protocol) doAtnOrListenEOIw(_ uint8, _ bool, busReadClk bool, _ bool) {
	if !busReadClk {
		//Sender set CLK=0, is about to send first bit
		v.ps.StateMachineSet(pBit0)
	}
}

// doAtnOrListenBit processes attention or listen bit state transitions based on clock (CLK) and data line signals.
// If the CLK line is high, reads the DATA line to set or clear a specific bit and advances the state machine.
func (v *Protocol) doAtnOrListenBit(sm uint8, _ bool, busReadClk bool, busReadData bool) {
	if busReadClk {
		//Sender set CLK=1, signaling that the DATA line represents a valid bit
		bit := _pBits[sm]
		if busReadData {
			v.ps.DataSetBit(bit)
		} else {
			v.ps.DataClearBit(bit)
		}
		//Go to associated P_BIT(n)w state, waiting for sender to set CLK=0
		v.ps.StateMachineAdvance()
	}
}

// doAtnOrListenBitW processes a single protocol state machine bit in ATN or Listen mode when the clock signal is low.
func (v *Protocol) doAtnOrListenBitW(_ uint8, _ bool, busReadClk bool, _ bool) {
	if !busReadClk {
		//Sender set CLK=0. go to P_BIT(n+1) state to receive the next bit
		v.ps.StateMachineAdvance()
	}
}

// doAtnOrListenBit7W handles ATN or listening operations on bit 7 based on the current state of the IEC bus signals.
// It checks the clock and data lines to determine whether to store addresses, acknowledge frames, or stop listening.
// If operating under ATN, it manages primary/secondary address detection and determines if the device should respond.
// When listening for data, it passes received bytes to the upper level and evaluates error conditions during writes.
// The state machine transitions are updated according to the bus operation context and acknowledgment needs.
func (v *Protocol) doAtnOrListenBit7W(_ uint8, busReadAtn bool, busReadClk bool, _ bool) {
	if !busReadClk {
		//Sender set CLK=0 and this was the last bit
		if v.debug {
			log.Printf("device %d received : 0x%02x (%c)", v.deviceNumber, v.ps.DataGetByte(), v.ps.DataGetByte())
		}
		if v.ps.FlagGet(pAtn) {
			//We are currently receiving under ATN. Store the first two bytes received (contain primary and secondary address)
			if v.ps.PrimaryGet() == 0 {
				v.ps.PrimarySet(v.ps.DataGetByte())
			} else if v.ps.SecondaryGet() == 0 {
				v.ps.SecondarySet(v.ps.DataGetByte())
			}
			if (v.ps.PrimaryGet() != pUnlisten) && (v.ps.PrimaryGet() != pUntalk) && ((v.ps.PrimaryGet() & 0x1f) != v.deviceNumber) {
				//This is NOT a UNLISTEN (0x3f) or UNTALK (0x5f) command and the primary address is not ours =>
				//Don't acknowledge the frame and stop listening. If all devices on the bus do this, the bus-master knows that "DeviceAdapter not present"
				v.ps.StateMachineSet(pDone0)
			} else {
				//Acknowledge frame by setting DATA=0
				v.peripheralWrite(busReadAtn, DeviceWriteClk)
				//repeat from pPre2 (we know that CLK=0 so no need to go to pPre1)
				v.ps.StateMachineSet(pPre2)
			}
		} else if v.ps.FlagGet(pListening) {
			//We are currently listening for data pass received byte on to the upper level
			if v.debug {
				log.Printf("device %d received 0x%02x (%c) on channel %d", v.deviceNumber, v.ps.DataGetByte(), v.ps.DataGetByte(), v.ps.SecondaryGet())
			}
			v.ps.StateSet(v.ps.SecondaryGet(), v.device.Write(v.ps.SecondaryGet(), v.ps.DataGetByte()))
			if v.ps.StateGet(v.ps.SecondaryGet()) != 0 {
				//There was an error during iec_bus_write => stop listening. This will signal an error condition to the sender
				v.ps.StateMachineSet(pDone0)
			} else {
				//Acknowledge frame by setting DATA=0
				v.peripheralWrite(busReadAtn, DeviceWriteClk)
				//repeat from pPre2 (we know that CLK=0 so no need to go to pPre1)
				v.ps.StateMachineSet(pPre2)
			}
		}
	}
}

// doAtnOrListenDone is executed when the bus-master must set ATN back to 1. It finalizes the state machine's operation.
func (v *Protocol) doAtnOrListenDone(_ uint8, _ bool, _ bool, _ bool) {
	//We're just waiting for the bus-master to set ATN back to 1
	//fmt.Println("We're just waiting for the bus-master to set ATN back to 1")
}

// setupAtnOrListen initializes a state machine for handling ATN or Listen operations and returns a slice of handler functions.
// Each state in the machine corresponds to a function that accepts parameters controlling its behavior.
func (v *Protocol) setupAtnOrListen() []func(uint8, bool, bool, bool) {
	unsupportedFn := func(sm uint8, _ bool, _ bool, _ bool) {
		log.Fatal("unsupported function setupAtnOrListen called with state", sm)
	}
	atnOrListenSM := make([]func(uint8, bool, bool, bool), 0xff)
	for idx := range atnOrListenSM {
		atnOrListenSM[idx] = unsupportedFn
	}
	atnOrListenSM[pPre0] = v.doAtnOrListenPre0
	atnOrListenSM[pPre1] = v.doAtnOrListenPre1
	atnOrListenSM[pPre2] = v.doAtnOrListenPre2
	atnOrListenSM[pReady] = v.doAtnOrListenReady
	atnOrListenSM[pEOI] = v.doAtnOrListenEOI
	atnOrListenSM[pEOIw] = v.doAtnOrListenEOIw
	atnOrListenSM[pBit0] = v.doAtnOrListenBit
	atnOrListenSM[pBit1] = v.doAtnOrListenBit
	atnOrListenSM[pBit2] = v.doAtnOrListenBit
	atnOrListenSM[pBit3] = v.doAtnOrListenBit
	atnOrListenSM[pBit4] = v.doAtnOrListenBit
	atnOrListenSM[pBit5] = v.doAtnOrListenBit
	atnOrListenSM[pBit6] = v.doAtnOrListenBit
	atnOrListenSM[pBit7] = v.doAtnOrListenBit
	atnOrListenSM[pBit0w] = v.doAtnOrListenBitW
	atnOrListenSM[pBit1w] = v.doAtnOrListenBitW
	atnOrListenSM[pBit2w] = v.doAtnOrListenBitW
	atnOrListenSM[pBit3w] = v.doAtnOrListenBitW
	atnOrListenSM[pBit4w] = v.doAtnOrListenBitW
	atnOrListenSM[pBit5w] = v.doAtnOrListenBitW
	atnOrListenSM[pBit6w] = v.doAtnOrListenBitW
	atnOrListenSM[pBit7w] = v.doAtnOrListenBit7W
	atnOrListenSM[pDone0] = v.doAtnOrListenDone
	return atnOrListenSM
}

// doTalkPre0 handles the transition to pPre1 state during role reversal by setting CLK to 0 and DATA to 1 if CLK is high.
func (v *Protocol) doTalkPre0(_ uint8, busReadAtn bool, busReadClk bool, _ bool) {
	if busReadClk {
		//Bus-master set CLK=1 (and before that should have set DATA=0)
		//we are getting ready for role reversal.Set CLK=0,DATA=1
		v.peripheralWrite(busReadAtn, DeviceWriteData)
		v.ps.StateMachineSet(pPre1)
		v.ps.TimeoutSet(v.quartz, 80)
	}
}

// doTalkPre1 performs a pre-phase talk check by verifying timeout expiration and signaling readiness if applicable.
func (v *Protocol) doTalkPre1(_ uint8, busReadAtn bool, _ bool, _ bool) {
	if v.ps.TimeoutExpired(v.quartz) {
		//Signal "ready-to-send" (CLK=1)
		v.peripheralWrite(busReadAtn, DeviceWriteClk|DeviceWriteData)
		v.ps.StateMachineSet(pReady)
	}
}

// doTalkReady processes the "ready-for-data" signal and advances or terminates the data transmission state machine.
func (v *Protocol) doTalkReady(_ uint8, _ bool, _ bool, busReadData bool) {
	if busReadData {
		//Receiver signaled "ready-for-data" (DATA=1)
		b, state := v.device.Read(v.ps.SecondaryGet())
		v.ps.StateSet(v.ps.SecondaryGet(), state)
		v.ps.DataSetByte(b)
		if v.ps.StateGet(v.ps.SecondaryGet()) == 0 {
			//At least two bytes left to send. Go on to send the first bit.
			v.ps.StateMachineSet(pBit0)
			//no need to wait before sending the first bit
			v.ps.TimeoutSet(v.quartz, 0)
		} else if v.ps.StateGet(v.ps.SecondaryGet()) == pRequestTalking {
			//Only this byte left to send => signal EOI by keeping CLK=1
			if v.debug {
				log.Printf("device %d signaling EOI on channel %d", v.deviceNumber, v.ps.SecondaryGet())
			}
			v.ps.StateMachineSet(pEOI)
		} else {
			//There was some kind of error; we have nothing to send.
			//Just stop talking and wait for ATN (This will produce a "File not found" when loading)
			v.ps.FlagsRemove(pTalking)
		}
	}
}

// doTalkEOI handles the End-Or-Identify (EOI) operation in the protocol, transitioning to the pEOIw state if busReadData is low.
func (v *Protocol) doTalkEOI(_ uint8, _ bool, _ bool, busReadData bool) {
	if !busReadData {
		//Receiver set DATA=0, first part of acknowledging the EOI
		v.ps.StateMachineSet(pEOIw)
	}
}

// doTalkEOIw transitions the protocol state machine based on busReadData, setting it to pBit0 and configuring timeout.
func (v *Protocol) doTalkEOIw(_ uint8, _ bool, _ bool, busReadData bool) {
	if busReadData {
		//Receiver set DATA=1, final part of acknowledging the EOI. Go on to send first bit
		v.ps.StateMachineSet(pBit0)
		//no need to wait before sending the first bit
		v.ps.TimeoutSet(v.quartz, 0)
	}
}

// doTalkBit manages the protocol timing and signals for data transmission by monitoring timeout and state transitions.
func (v *Protocol) doTalkBit(sm uint8, busReadAtn bool, _ bool, _ bool) {
	if v.ps.TimeoutExpired(v.quartz) {
		//60 us have passed since we set CLK=1 to signal "data valid" for the previous bit.
		//Pull CLK=0 and put the next bit out of DATA.
		bit := _pBits[sm]
		if v.ps.DataHasBit(bit) {
			v.peripheralWrite(busReadAtn, DeviceWriteData)
		} else {
			v.peripheralWrite(busReadAtn, 0)
		}
		//Go to associated P_BIT(n)w state
		v.ps.TimeoutSet(v.quartz, 90) //orig 60
		v.ps.StateMachineAdvance()
	}
}

// doTalkBitW handles the signaling of a data bit over the protocol, including timeout checks and state progression.
// sm represents the state machine identifier.
// busReadAtn indicates the current attention signal status on the bus.
// busReadClk represents the clock signal read status from the bus.
// busReadData specifies the current data signal status on the bus.
func (v *Protocol) doTalkBitW(_ uint8, busReadAtn bool, _ bool, busReadData bool) {
	if v.ps.TimeoutExpired(v.quartz) {
		//60 us have passed since we pulled CLK=0 and put the current bit on DATA.
		//set CLK=1, keeping data as it is (this signals "data valid" to the receiver)
		if busReadData {
			v.peripheralWrite(busReadAtn, DeviceWriteClk|DeviceWriteData)
		} else {
			v.peripheralWrite(busReadAtn, DeviceWriteClk)
		}
		//Go to associated P_BIT(n+1) state to send the next bit.
		//If this was the final bit, then the next state is pDone0
		v.ps.TimeoutSet(v.quartz, 90) //orig 60
		v.ps.StateMachineAdvance()
	}
}

// doTalkDone0 manages the final steps of a communication cycle, preparing the bus and state for acknowledgement.
func (v *Protocol) doTalkDone0(_ uint8, busReadAtn bool, _ bool, _ bool) {
	if v.ps.TimeoutExpired(v.quartz) {
		//60 us have passed since we set CLK=1 to signal "data valid" for the final bit.
		//Pull CLK=0 and set DATA=1.This prepares for the receiver acknowledgement.
		v.peripheralWrite(busReadAtn, DeviceWriteData)
		v.ps.TimeoutSet(v.quartz, 1000)
		v.ps.StateMachineSet(pDone1)
	}
}

// doTalkDone1 manages the process of sending or receiving data, acknowledges frames, and handles timeouts or errors.
func (v *Protocol) doTalkDone1(_ uint8, busReadAtn bool, _ bool, busReadData bool) {
	if !busReadData {
		//Receiver set DATA=0, acknowledging the frame
		if v.debug {
			log.Printf("device %d sent 0x%02x (%c) on channel %d", v.deviceNumber, v.ps.DataGetByte(), v.ps.DataGetByte(), v.ps.SecondaryGet())
		}
		if v.ps.StateGet(v.ps.SecondaryGet()) == pRequestTalking {
			//This was the last byte => stop talking.This leaves us waiting for ATN.
			v.ps.FlagsRemove(pTalking)
			v.ps.StateSet(v.ps.SecondaryGet(), 0)
			//Release the CLOCK line to 1
			v.peripheralWrite(busReadAtn, DeviceWriteClk|DeviceWriteData)
		} else {
			//There is at least one more byte to send Start over from pPre1
			v.ps.TimeoutSet(v.quartz, 0)
			v.ps.StateMachineSet(pPre1)
		}
	} else if v.ps.TimeoutExpired(v.quartz) {
		//We didn't receive an acknowledgement within 1 ms.Set CLOCK=0 and after 100 us back to CLOCK=1
		if v.debug {
			log.Printf("device %d got NACK on channel %d", v.deviceNumber, v.ps.SecondaryGet())
		}
		v.peripheralWrite(busReadAtn, DeviceWriteClk|DeviceWriteData)
		v.ps.TimeoutSet(v.quartz, 100)
		v.ps.StateMachineSet(pFrameError0)
	}
}

// doTalkFrameError0 handles the 1-0-1 CLOCK signal sequence to acknowledge a frame error and transitions the state machine.
func (v *Protocol) doTalkFrameError0(_ uint8, busReadAtn bool, _ bool, _ bool) {
	if v.ps.TimeoutExpired(v.quartz) {
		//Finished 1-0-1 sequence of CLOCK signal
		//to acknowledge the frame-error.Now wait for sender to set DATA=0 so we can continue.
		v.peripheralWrite(busReadAtn, DeviceWriteData)
		v.ps.StateMachineSet(PFrameError1)
	}
}

// doTalkFrameError1 handles error recovery when a communication issue occurs during a talk frame operation.
// If busReadData is false, it resets the timeout and sets the state machine to pPre1 for retry.
func (v *Protocol) doTalkFrameError1(_ uint8, _ bool, _ bool, busReadData bool) {
	if !busReadData {
		// sender set DATA=0, we can retry to send the byte
		v.ps.TimeoutSet(v.quartz, 0)
		v.ps.StateMachineSet(pPre1)
	}
}

// setupTalk initializes an array of state machine functions for talk handling and returns the configured function slice.
// Each state is associated with a specific handler, defaulting to an unsupported state function when unconfigured.
func (v *Protocol) setupTalk() []func(uint8, bool, bool, bool) {
	unsupportedFn := func(sm uint8, _ bool, _ bool, _ bool) {
		log.Fatal("unsupported talk method called with state", sm)
	}
	talkSM := make([]func(uint8, bool, bool, bool), 0xff)
	for idx := range talkSM {
		talkSM[idx] = unsupportedFn
	}
	talkSM[pPre0] = v.doTalkPre0
	talkSM[pPre1] = v.doTalkPre1
	talkSM[pReady] = v.doTalkReady
	talkSM[pEOI] = v.doTalkEOI
	talkSM[pEOIw] = v.doTalkEOIw
	talkSM[pBit0] = v.doTalkBit
	talkSM[pBit1] = v.doTalkBit
	talkSM[pBit2] = v.doTalkBit
	talkSM[pBit3] = v.doTalkBit
	talkSM[pBit4] = v.doTalkBit
	talkSM[pBit5] = v.doTalkBit
	talkSM[pBit6] = v.doTalkBit
	talkSM[pBit7] = v.doTalkBit
	talkSM[pBit0w] = v.doTalkBitW
	talkSM[pBit1w] = v.doTalkBitW
	talkSM[pBit2w] = v.doTalkBitW
	talkSM[pBit3w] = v.doTalkBitW
	talkSM[pBit4w] = v.doTalkBitW
	talkSM[pBit5w] = v.doTalkBitW
	talkSM[pBit6w] = v.doTalkBitW
	talkSM[pBit7w] = v.doTalkBitW
	talkSM[pDone0] = v.doTalkDone0
	talkSM[pDone1] = v.doTalkDone1
	talkSM[pFrameError0] = v.doTalkFrameError0
	talkSM[PFrameError1] = v.doTalkFrameError1
	return talkSM
}

// peripheralWrite performs a write operation to the peripheral using the provided data and bus attention state.
// busReadAtn determines if the ATN (Attention) line is active during the write operation.
// data is an 8-bit value to be sent to the peripheral.
func (v *Protocol) peripheralWrite(busReadAtn bool, data uint8) {
	sidecarData := references.IECSidecarEnabled | references.IECSidecarAtnAEnabled
	if busReadAtn {
		sidecarData |= uint16(references.IECAtnABit) << 8
	}
	out := sidecarData | uint16(data&defaultDDRBMask)
	v.iec.PeripheralWrite(v.deviceNumber, out)
}
