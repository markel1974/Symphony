package iec_rev1

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/hardware/quartz_rev1"
	"github.com/markel1974/c64emu/src/references"
	"log"
)

// defaultDDRBMask is a bitmask defining active bits for DDRB configuration, enabling bits 1, 3, and 4 in an 8-bit register.
const defaultDDRBMask = uint8((1 << 1) | (1 << 3) | (1 << 4))

// DeviceReadData represents the input data signal (DATA_IN).
// DeviceReadClk represents the input clock signal (CLK_IN).
// DeviceReadAtn represents the input attention signal (ATN_IN).
// DeviceWriteData represents the output data signal (DATA_OUT).
// DeviceWriteClk represents the output clock signal (CLK_OUT).
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

// pPre0 to pPre2 represent initial states in a sequence based on the iota increment.
// pReady represents the ready state in the sequence.
// pEOI and pEOIw represent end-of-input states with and without additional conditions.
// pBit0 to pBit7 signify bit states in a sequence, with their corresponding writable states suffixed with 'w'.
// pDone0 and pDone1 denote completion states in the process.
// pFrameError0 and PFrameError1 signify frame error states in the sequence.
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

// pTalking represents the state or flag for talking with a hexadecimal value of 0x20.
// pListening represents the state or flag for listening with a hexadecimal value of 0x40.
// pAtn represents the state or flag for attention with a hexadecimal value of 0x80.
// pRequestListen represents the request flag for listening with a hexadecimal value of 0x20.
// pRequestTalking represents the request flag for talking with a hexadecimal value of 0x40.
// pUnlisten represents the flag to clear the listening state with a hexadecimal value of 0x3f.
// pUntalk represents the flag to clear the talking state with a hexadecimal value of 0x5f.
const (
	pTalking   = uint8(0x20)
	pListening = uint8(0x40)
	pAtn       = uint8(0x80)

	pRequestListen  = uint8(0x20)
	pRequestTalking = uint8(0x40)

	pUnlisten = uint8(0x3f)
	pUntalk   = uint8(0x5f)
)

// stateLast defines the maximum valid index for state-related operations, used to mask and bound state array access.
const stateLast = 0xf

// _pBits is a precomputed lookup table used to store bit masks for quick access in bitwise operations.
var _pBits [0xff]uint8

// init initializes the _pBits array by mapping specific bit positions (pBit0 to pBit7) to their corresponding binary values.
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

// Protocol represents a structure for handling IEC protocol communication and related state management.
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

// NewProtocol initializes a new Protocol instance with the given factory, parent component, label, and instance number.
// It sets up the base component, protocol state, device, and quartz instance to manage protocol functionality.
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

// SetDevice assigns a high-level device implementation to the protocol instance, enabling specific device handling logic.
func (v *Protocol) SetDevice(device references.IIecProtocolDevice) {
	// Associates a high-level logic implementation (e.g., a file system handler)
	// with this protocol instance. The protocol will handle the low-level (bit-by-bit)
	// communication, while the 'device' object will handle data interpretation (e.g., opening files, reading blocks).
	v.device = device
}

// Setup prepares the Protocol component for operation.
// It initializes configuration, dispatch tables, and manages timing setup via the quartz component.
func (v *Protocol) Setup() error {
	// This method prepares the component for use.
	// It initializes the pointer to the global configuration, builds the dispatch tables
	// for the state machine (one for 'Listen/ATN' mode and one for 'Talk' mode),
	// and initializes the quartz component for timing management.
	v.cfg = v.GetFactory().GetConfig()
	v.atnOrListenFn = v.setupAtnOrListen()
	v.talkFn = v.setupTalk()
	if err := v.quartz.Setup(); err != nil {
		return err
	}
	return nil
}

// Bind connects the Protocol instance to the IEC system, initializing required dependencies like quartz and IEC bus handler.
// It assigns the specified device number and verifies successful integration of components.
func (v *Protocol) Bind(_ references.IIecDeviceSocket, _ uint8, deviceNumber uint8) error {
	// This method "wires" the component to the rest of the system after its creation.
	// It sets the device number (e.g., 8, 9) and gets a reference
	// to the main IEC bus handler to be able to read and write to the bus lines.
	// It also initializes the quartz to synchronize timed operations.
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

// Connect initializes or establishes a connection for the protocol and returns an error if any issue occurs.
func (v *Protocol) Connect() error {
	// Method part of the Symphony framework lifecycle.
	// Currently does not perform specific operations for this component,
	// but is present for architectural consistency.
	return nil
}

// Internal determines if the protocol component is a top-level entity rather than being internal to another component.
func (v *Protocol) Internal() bool {
	// Method part of the Symphony framework lifecycle.
	// Indicates that this component is not "internal" to another, but a top-level entity.
	return false
}

// Shutdown terminates the Protocol component gracefully without requiring additional cleanup operations.
func (v *Protocol) Shutdown() {
	// Method for the clean termination of the component.
	// No specific cleanup operations are currently needed.
}

// Reset returns the component to its initial state.
// Resets the internal state machine and ensures the device releases the bus lines (CLK and DATA to high-impedance state).
func (v *Protocol) Reset() {
	// Returns the component to its initial state.
	// Resets the internal state machine and ensures the device
	// releases the bus lines (CLK and DATA to 1, high-impedance state).
	v.ps.Reset()
	v.peripheralWrite(false, DeviceWriteClk|DeviceWriteData)
}

// Ready checks if the Protocol is ready to operate and always returns true after initialization.
func (v *Protocol) Ready() bool {
	// Method to check if the component is ready to operate.
	// For this component, it is always considered ready after initialization.
	return true
}

// LedActivity controls the LED state by sending a request to the IEC bus manager to update the emulator's frontend.
func (v *Protocol) LedActivity(led bool) {
	// Propagates the request to turn the LED on/off to the IEC bus manager,
	// which in turn will communicate it to the emulator's frontend.
	v.iec.LedActivity(v.deviceNumber, led)
}

// GetDeviceNumber returns the configured device number of the component as an 8-bit unsigned integer.
func (v *Protocol) GetDeviceNumber() uint8 {
	// Returns the device number (e.g., 8, 9) with which this component was configured.
	return v.deviceNumber
}

// AtnStateChanged is invoked when the ATN (Attention) line state changes in the IEC bus, typically for debugging purposes.
// The method provides traceability of command phase transitions by logging state changes when debugging is enabled.
func (v *Protocol) AtnStateChanged(atn bool) {
	// Method called by the IEC bus manager when the ATN (Attention) line changes state.
	// Useful for debugging to trace when the computer (master) starts or ends a command phase.
	if v.debug {
		if !atn {
			log.Println("ATN STATE CHANGED", "ON")
		} else {
			log.Println("ATN STATE CHANGED", "OFF")
		}
	}
}

// EmulationRequired returns true to signal that the component requires the Emulate() method to be invoked every cycle.
func (v *Protocol) EmulationRequired() bool {
	// Indicates to the Symphony framework that this component needs to have its Emulate() method
	// called on every emulation cycle, as it manages an active state machine.
	return true
}

// Emulate executes the main logic cycle for the protocol, managing state transitions and bus communication processes.
func (v *Protocol) Emulate() {
	// This is the heartbeat of the component, executed on every cycle.
	// 1. Reads the current state of the IEC bus lines (CLK, DATA, ATN).
	// 2. Detects the rising and falling edges of the ATN line to switch between command mode and data transfer mode.
	// 3. Based on the current flags (pAtn, pListening, pTalking), it selects the correct dispatch table.
	// 4. Executes the specific handler function for the current state of the state machine, passing the bus status.
	// 5. Advances the internal time for timeout management.
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

// doAtnFallingFlank handles the transition when ATN line is pulled low, signaling the start of a command sequence on the bus.
// Resets the state machine, sets the ATN flag, clears device addresses, sets a timeout, and asserts presence on the DATA line.
func (v *Protocol) doAtnFallingFlank(busReadAtn bool) {
	// This function is executed when the computer (master) pulls the ATN line low,
	// indicating the start of a command sequence.
	// 1. The state machine is reset to the initial state (pPre0).
	// 2. The 'pAtn' flag is set to indicate that we are in command mode.
	// 3. The primary and secondary addresses are cleared.
	// 4. A timeout is set for bus stabilization.
	// 5. The device pulls the DATA line low to signal its presence ("I am here").
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

// doAtnRisingFlank is triggered when the master releases the ATN line, signaling the end of a command sequence.
// It clears the pAtn flag, checks the message target, and handles LISTEN or TALK requests for the device.
// Depending on the command and device state, it transitions the device to listening, talking, or standby mode.
// It also processes global UNLISTEN or UNTALK commands to deactivate communication and potentially open files.
// If the device is neither listening nor talking, it releases the DATA and CLOCK lines.
func (v *Protocol) doAtnRisingFlank(busReadAtn bool) {
	// This function is executed when the computer (master) releases the ATN line.
	// This means the command sequence has finished, and the device must interpret it.
	// 1. The 'pAtn' flag is removed.
	// 2. It checks if the command was addressed to this device (LISTEN or TALK).
	// 3. If so, it performs the requested action (Open, Close, Listen, Talk) by calling the high-level 'device'.
	// 4. It sets the 'pListening' or 'pTalking' flags to switch to data transfer mode.
	// 5. It handles the global UNLISTEN and UNTALK commands to deactivate communication.
	// 6. If it is neither listening nor talking, it releases the bus lines.
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
			//set DATA=0 -> sending on the bus "present"
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

// doAtnOrListenPre0 handles the P_PRE0 state where the protocol waits for bus stabilization after ATN drops to 0.
// It ensures a delay of 100µs to allow devices to release CLK and DATA lines before transitioning to the next state.
func (v *Protocol) doAtnOrListenPre0(_ uint8, _ bool, _ bool, _ bool) {
	// State P_PRE0: Waiting for bus stabilization.
	// After ATN goes to 0, wait for a short period (100µs) to allow
	// all devices to release the CLK and DATA lines.
	// Only after the timeout has expired, proceed to the next state (pPre1).
	//Ignore anything that happens during the first 100 us after falling
	//flank on ATN (other devices may have been sending and need some time to set CLK=1)
	if v.ps.TimeoutExpired(v.quartz) {
		//test := time.Since(P_PRE0_100_time)
		//fmt.Println("P_PRE0_100_time", test)
		v.ps.StateMachineSet(pPre1)
	}
}

// doAtnOrListenPre1 transitions the protocol state to pPre2 if the CLK line is low, ensuring readiness for the next phase.
func (v *Protocol) doAtnOrListenPre1(_ uint8, _ bool, busReadClk bool, _ bool) {
	// State P_PRE1: Waiting for the sender (the C64) to release the CLK line.
	// The protocol dictates that before starting a transmission, the sender
	// must pull CLK to 0. This state waits for that to happen before proceeding.
	//Make sure CLK=0 so we actually detect a rising flank instate pPre2
	if !busReadClk {
		v.ps.StateMachineSet(pPre2)
	}
}

// doAtnOrListenPre2 processes the state P_PRE2, handling "ready-to-send" and "ready-to-receive" handshakes in communication.
func (v *Protocol) doAtnOrListenPre2(_ uint8, busReadAtn bool, busReadClk bool, _ bool) {
	// State P_PRE2: Waiting for the "ready to send" signal from the sender.
	// The sender (C64) pulls the CLK line to 1.
	// The receiver (this device) responds by pulling the DATA line to 1,
	// signaling "ready to receive data".
	// wait for rising flank on CLK ("ready-to-send")
	if busReadClk {
		//React by setting DATA=1 ("ready-for-data")
		v.peripheralWrite(busReadAtn, DeviceWriteClk|DeviceWriteData)
		v.ps.TimeoutSet(v.quartz, 200)
		v.ps.StateMachineSet(pReady)
	}
}

// doAtnOrListenReady handles the state where the device is ready to receive, reacting to CLK changes or timeout conditions.
// If CLK is pulled low, it transitions to receiving the first bit. If a timeout occurs without ATN, it signals EOI.
func (v *Protocol) doAtnOrListenReady(_ uint8, busReadAtn bool, busReadClk bool, _ bool) {
	// State P_READY: The device is ready to receive.
	// It waits for the sender to pull CLK to 0, which is the signal that the first bit of data is about to be placed on the DATA line.
	// If the sender doesn't pull CLK low within a timeout and we are not in ATN mode,
	// it means the sender is signaling EOI (End-Of-Identify), the end of the transmission.
	if !busReadClk {
		//Sender set CLK=0, is about to send first bit
		v.ps.StateMachineSet(pBit0)
	} else if !v.ps.FlagGet(pAtn) && v.ps.TimeoutExpired(v.quartz) {
		//Sender did not set CLK=0 within 200 us after we set DATA=1 => it is signaling EOI
		//(not so if we are under ATN) acknowledge we received it by setting DATA=0 for 60us
		if v.debug {
			log.Printf("device %d got EOI on channel %d", v.deviceNumber, v.ps.SecondaryGet()&0x0f)
		}
		v.device.EOI(v.ps.SecondaryGet())
		v.peripheralWrite(busReadAtn, DeviceWriteClk)
		v.ps.StateMachineSet(pEOI)
		v.ps.TimeoutSet(v.quartz, 60)
	}
}

// doAtnOrListenEOI manages the transition during the End-Or-Identify (EOI) phase in the protocol state machine.
// Handles DATA line behavior and state advancement based on timeout conditions and sender actions.
func (v *Protocol) doAtnOrListenEOI(_ uint8, busReadAtn bool, _ bool, _ bool) {
	// State P_EOI: Acknowledging the EOI signal.
	// The device holds the DATA line low for 60µs.
	// After the timeout, it releases DATA back to 1 and waits for the sender to continue.
	if v.ps.TimeoutExpired(v.quartz) {
		//Set DATA back to 1 and wait for sender to set CLK=0
		v.peripheralWrite(busReadAtn, DeviceWriteClk|DeviceWriteData)
		v.ps.StateMachineSet(pEOIw)
	}
}

// doAtnOrListenEOIw transitions the protocol's state machine upon detecting a clock signal change during the EOI wait state.
// If the CLK signal is pulled low, it indicates the sender is ready to transmit, and the state advances to `pBit0`.
func (v *Protocol) doAtnOrListenEOIw(_ uint8, _ bool, busReadClk bool, _ bool) {
	// State P_EOIW: Waiting state after acknowledging EOI.
	// The device waits for the sender to pull CLK to 0 to resume data transfer (if any).
	if !busReadClk {
		//Sender set CLK=0, is about to send first bit
		v.ps.StateMachineSet(pBit0)
	}
}

// doAtnOrListenBit reads a data bit during P_BIT(n) state when CLK is high, updates the data buffer, and advances the state machine.
func (v *Protocol) doAtnOrListenBit(sm uint8, _ bool, busReadClk bool, busReadData bool) {
	// State P_BIT(n): Receiving a data bit.
	// When the sender pulls CLK high, the bit on the DATA line is valid.
	// This function reads the DATA line and sets or clears the corresponding bit
	// in the internal byte buffer (`v.ps.data`). It then advances to the corresponding "wait" state.
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

// doAtnOrListenBitW handles the state P_BIT(n)W, waiting for the sender to pull CLK low before advancing state.
// If the busReadClk signal is low, it triggers an advance in the state machine to prepare for the next bit.
func (v *Protocol) doAtnOrListenBitW(_ uint8, _ bool, busReadClk bool, _ bool) {
	// State P_BIT(n)W: Waiting after receiving a bit.
	// The device waits for the sender to pull CLK low again, which signals
	// that the sender is preparing the next bit. Once CLK is low, the state machine advances.
	if !busReadClk {
		//Sender set CLK=0. go to P_BIT(n+1) state to receive the next bit
		v.ps.StateMachineAdvance()
	}
}

// doAtnOrListenBit7W handles the protocol's logic upon receiving the 8th bit of a byte in ATN or LISTEN mode.
// In ATN mode, it checks and processes command bytes, ensuring proper device addressing and acknowledgment.
// In LISTEN mode, it processes received data, forwarding it to high-level logic and signaling errors if necessary.
func (v *Protocol) doAtnOrListenBit7W(_ uint8, busReadAtn bool, busReadClk bool, _ bool) {
	// State P_BIT7W: The final state after receiving the 8th bit of a byte.
	if !busReadClk {
		// The sender has pulled CLK low, signaling the end of the byte transfer.
		if v.debug {
			log.Printf("device %d received : 0x%02x (%c)", v.deviceNumber, v.ps.DataGetByte(), v.ps.DataGetByte())
		}
		if v.ps.FlagGet(pAtn) {
			// If in ATN mode, the received byte is part of a command (primary or secondary address).
			// It checks if the command is for this device. If not, it stops listening.
			// If it is, it acknowledges the byte by pulling DATA low and gets ready for the next command byte.
			if v.ps.PrimaryGet() == 0 {
				v.ps.PrimarySet(v.ps.DataGetByte())
			} else if v.ps.SecondaryGet() == 0 {
				v.ps.SecondarySet(v.ps.DataGetByte())
			}
			if (v.ps.PrimaryGet() != pUnlisten) && (v.ps.PrimaryGet() != pUntalk) && ((v.ps.PrimaryGet() & 0x1f) != v.deviceNumber) {
				v.ps.StateMachineSet(pDone0)
			} else {
				v.peripheralWrite(busReadAtn, DeviceWriteClk)
				v.ps.StateMachineSet(pPre2)
			}
		} else if v.ps.FlagGet(pListening) {
			// If in LISTEN mode, the byte is data. Pass it to the high-level device logic.
			// If the device logic reports an error, stop listening to signal the error.
			// Otherwise, acknowledge the byte and get ready for the next one.
			if v.debug {
				log.Printf("device %d received 0x%02x (%c) on channel %d", v.deviceNumber, v.ps.DataGetByte(), v.ps.DataGetByte(), v.ps.SecondaryGet())
			}
			v.ps.StateSet(v.ps.SecondaryGet(), v.device.Write(v.ps.SecondaryGet(), v.ps.DataGetByte()))
			if v.ps.StateGet(v.ps.SecondaryGet()) != 0 {
				v.ps.StateMachineSet(pDone0)
			} else {
				v.peripheralWrite(busReadAtn, DeviceWriteClk)
				v.ps.StateMachineSet(pPre2)
			}
		}
	}
}

// doAtnOrListenDone is invoked in the P_DONE0 state, a final waiting state when an error occurs or the command isn't applicable.
// It waits for the bus-master to pull ATN high, signaling the end of the command sequence.
func (v *Protocol) doAtnOrListenDone(_ uint8, _ bool, _ bool, _ bool) {
	// State P_DONE0: A final waiting state.
	// This state is entered if an error occurred or if a command was not for this device.
	// The device simply waits for the master to end the command sequence by pulling ATN high.
	//We're just waiting for the bus-master to set ATN back to 1
	//fmt.Println("We're just waiting for the bus-master to set ATN back to 1")
}

// setupAtnOrListen initializes and returns a state machine table with handlers for "Listen" and "ATN" modes.
func (v *Protocol) setupAtnOrListen() []func(uint8, bool, bool, bool) {
	// This function builds the dispatch table for the "Listen" and "ATN" modes.
	// It creates a slice of functions and assigns the correct handler function
	// to the index corresponding to each state constant (pPre0, pBit0, etc.).
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

// doTalkPre0 transitions the device to the Talker role, preparing to send data by taking control of the clock line.
func (v *Protocol) doTalkPre0(_ uint8, busReadAtn bool, busReadClk bool, _ bool) {
	// State P_PRE0 (Talk): The first step for a device becoming a sender (Talker).
	// It waits for the master (C64, now the receiver) to signal it's ready by pulling CLK high.
	// The device then pulls CLK low to take control of the clock line and prepares for transmission.
	if busReadClk {
		//Bus-master set CLK=1 (and before that should have set DATA=0)
		//we are getting ready for role reversal.Set CLK=0,DATA=1
		v.peripheralWrite(busReadAtn, DeviceWriteData)
		v.ps.StateMachineSet(pPre1)
		v.ps.TimeoutSet(v.quartz, 80)
	}
}

// doTalkPre1 handles the P_PRE1 (Talk) state in the protocol, initiating a wait period before signaling readiness to send data.
func (v *Protocol) doTalkPre1(_ uint8, busReadAtn bool, _ bool, _ bool) {
	// State P_PRE1 (Talk): A brief waiting period.
	// After taking control of the clock line, the device waits for a short timeout
	// and then pulls CLK high to signal "I'm ready to send data".
	if v.ps.TimeoutExpired(v.quartz) {
		//Signal "ready-to-send" (CLK=1)
		v.peripheralWrite(busReadAtn, DeviceWriteClk|DeviceWriteData)
		v.ps.StateMachineSet(pReady)
	}
}

// doTalkReady handles the P_READY state, managing the device's readiness to send data when the receiver signals readiness.
// If the receiver signals "ready for data," it fetches the next byte to send, updating the state machine accordingly.
// Depending on the state, it moves to send the first bit (if more bytes remain) or signals EOI (if it's the last byte).
// Handles edge cases like errors where no data is available to send, stopping communication appropriately.
func (v *Protocol) doTalkReady(_ uint8, _ bool, _ bool, busReadData bool) {
	// State P_READY (Talk): The device is ready to send.
	// It waits for the receiver to signal "ready for data" by pulling the DATA line high.
	// Once this happens, it fetches the next byte to send from the high-level device logic.
	// Depending on whether it's the last byte or not, it transitions to send the first bit (pBit0) or to signal EOI.
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

// doTalkEOI handles the "Talk" state where the device signals the end of transmission (EOI) and waits for acknowledgment.
func (v *Protocol) doTalkEOI(_ uint8, _ bool, _ bool, busReadData bool) {
	// State P_EOI (Talk): Signaling the end of transmission.
	// The device holds CLK high and waits for the receiver to acknowledge the EOI
	// by pulling the DATA line low.
	if !busReadData {
		//Receiver set DATA=0, first part of acknowledging the EOI
		v.ps.StateMachineSet(pEOIw)
	}
}

// doTalkEOIw handles the protocol state when waiting for the receiver to acknowledge EOI and release the DATA line.
// Transitions the state machine to send the first bit once DATA is released by the receiver.
func (v *Protocol) doTalkEOIw(_ uint8, _ bool, _ bool, busReadData bool) {
	// State P_EOIW (Talk): Waiting after the receiver acknowledged EOI.
	// The device waits for the receiver to release the DATA line (pull it high).
	// It then proceeds to send the final byte.
	if busReadData {
		//Receiver set DATA=1, final part of acknowledging the EOI. Go on to send first bit
		v.ps.StateMachineSet(pBit0)
		//no need to wait before sending the first bit
		v.ps.TimeoutSet(v.quartz, 0)
	}
}

// doTalkBit handles the P_BIT(n) (Talk) state, transmitting a data bit by controlling CLK and DATA lines, and advancing state.
func (v *Protocol) doTalkBit(sm uint8, busReadAtn bool, _ bool, _ bool) {
	// State P_BIT(n) (Talk): Sending a data bit.
	// The device (sender) pulls CLK low and places the value of the current bit
	// onto the DATA line. It then sets a timeout and advances to the corresponding "wait" state.
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

// doTalkBitW handles the state P_BIT(n)W where the device ensures a bit on the DATA line is valid by signaling via CLK.
// Ensures a 60 µs timeout before signaling "data valid" to the receiver by pulling the CLK line high.
// Advances the state machine to the next bit's state or final state after signaling the receiver.
func (v *Protocol) doTalkBitW(_ uint8, busReadAtn bool, _ bool, busReadData bool) {
	// State P_BIT(n)W (Talk): Waiting after placing a bit on the DATA line.
	// The device waits for a short timeout, then pulls the CLK line high to signal
	// to the receiver that the bit on the DATA line is now valid and can be read.
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

// doTalkDone0 transitions to the final stage of the talk operation by pulling CLK low and starting acknowledgment timing.
// A longer timeout of 1ms is set, awaiting acknowledgment from the receiver before advancing the state machine.
func (v *Protocol) doTalkDone0(_ uint8, busReadAtn bool, _ bool, _ bool) {
	// State P_DONE0 (Talk): The final bit has been sent.
	// The device pulls CLK low and prepares to receive acknowledgement from the receiver.
	// It sets a longer timeout (1ms) to wait for this acknowledgement.
	if v.ps.TimeoutExpired(v.quartz) {
		//60 us have passed since we set CLK=1 to signal "data valid" for the final bit.
		//Pull CLK=0 and set DATA=1.This prepares for the receiver acknowledgement.
		v.peripheralWrite(busReadAtn, DeviceWriteData)
		v.ps.TimeoutSet(v.quartz, 1000)
		v.ps.StateMachineSet(pDone1)
	}
}

// doTalkDone1 handles the P_DONE1 state in the Talk process, verifying byte transmission and managing acknowledgments or timeouts.
func (v *Protocol) doTalkDone1(_ uint8, busReadAtn bool, _ bool, busReadData bool) {
	// State P_DONE1 (Talk): Waiting for byte acknowledgement.
	// If the receiver pulls DATA low, the byte was successfully received. The device then checks
	// if there is more data to send and either starts over or stops talking.
	// If the receiver doesn't acknowledge within the timeout, a frame error is triggered.
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

// doTalkFrameError0 handles a frame error scenario in the protocol's talk state by signaling and waiting for acknowledgment.
// It toggles the clock line, sends the necessary data, and transitions to the next error state after acknowledgment.
func (v *Protocol) doTalkFrameError0(_ uint8, busReadAtn bool, _ bool, _ bool) {
	// State P_FRAMEERROR0 (Talk): Handling a frame error (no acknowledgement).
	// The device signals the error by toggling the CLK line and then waits
	// for the receiver to respond to the error signal.
	if v.ps.TimeoutExpired(v.quartz) {
		//Finished 1-0-1 sequence of CLOCK signal
		//to acknowledge the frame-error.Now wait for sender to set DATA=0 so we can continue.
		v.peripheralWrite(busReadAtn, DeviceWriteData)
		v.ps.StateMachineSet(PFrameError1)
	}
}

// doTalkFrameError1 handles the P_FRAMEERROR1 state where the device waits for an acknowledgment of a frame error.
// If the busReadData signal is low (DATA=0), the device retries sending the byte by transitioning to the pPre1 state.
func (v *Protocol) doTalkFrameError1(_ uint8, _ bool, _ bool, busReadData bool) {
	// State P_FRAMEERROR1 (Talk): Waiting after signaling a frame error.
	// The device waits for the receiver to pull DATA low to acknowledge the error condition.
	// Once acknowledged, it will attempt to resend the byte by going back to the pPre1 state.
	if !busReadData {
		// sender set DATA=0, we can retry to send the byte
		v.ps.TimeoutSet(v.quartz, 0)
		v.ps.StateMachineSet(pPre1)
	}
}

// setupTalk initializes and returns a dispatch table for handling different "Talk" mode states in a protocol.
// It assigns specific state-handling functions to corresponding indices in the table.
// Defaults to an unsupported function for unhandled states to ensure errors are reported.
func (v *Protocol) setupTalk() []func(uint8, bool, bool, bool) {
	// This function builds the dispatch table for the "Talk" mode.
	// It creates a slice of functions and assigns the correct handler function
	// to the index corresponding to each state constant (pPre0, pBit0, etc.).
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

// peripheralWrite sends a combined data and control signal to the IEC bus for communication with peripheral devices.
// Depending on the ATN line state, it adjusts the sidecar bits to include ATN information along with the provided data.
func (v *Protocol) peripheralWrite(busReadAtn bool, data uint8) {
	// Helper function to construct the full 16-bit value to be sent to the bus dispatcher.
	// It combines the actual data to be written on the bus lines with sidecar information,
	// such as the state of the ATN line, before sending it to the main IEC bus component.
	sidecarData := references.IECSidecarEnabled | references.IECSidecarAtnAEnabled
	if busReadAtn {
		sidecarData |= uint16(references.IECAtnABit) << 8
	}
	out := sidecarData | uint16(data&defaultDDRBMask)
	v.iec.PeripheralWrite(v.deviceNumber, out)
}
