package iec

import (
	"fmt"
	"github.com/markel1974/c64emu/src/common/conversion"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
	"log"
)

//serial-iec-device.c
//static void serial_iec_device_exec_main(unsigned int devnr, CLOCK clk_value)

// IECBUS_DEVICE_READ_ATN represents the signal for device read attention (ATN_IN).
// IECBUS_DEVICE_READ_CLK represents the signal for device read clock (CLK_IN).
// IECBUS_DEVICE_READ_DATA represents the signal for device read data (DATA_IN).
// IECBUS_DEVICE_WRITE_DATA represents the signal for device write data (DATA_OUT).
// IECBUS_DEVICE_WRITE_CLK represents the signal for device write clock (CLK_OUT).
// IECBUS_DEVICE_ATNA represents the signal for attention acknowledge (ATN_A).
const (
	/* ATN_IN */ IECBUS_DEVICE_READ_ATN = 0x80
	/*  CLK_IN */ IECBUS_DEVICE_READ_CLK = 0x04
	/* DATA_IN */ IECBUS_DEVICE_READ_DATA = 0x01

	/* DATA_OUT*/
	IECBUS_DEVICE_WRITE_DATA = 0x02
	/* CLK_OUT*/ IECBUS_DEVICE_WRITE_CLK = 0x08
	///* ATN_A */ IECBUS_DEVICE_ATNA = 0x10
)

/*
const (
	IECBUS_DEVICE_READ_DATA = uint8(0x01)
	IECBUS_DEVICE_READ_CLK  = uint8(0x04)
	IECBUS_DEVICE_READ_ATN  = uint8(0x80)

	IECBUS_DEVICE_ATNA = uint8(0x10)

	IECBUS_DEVICE_WRITE_CLK  = uint8(0x40)
	IECBUS_DEVICE_WRITE_DATA = uint8(0x80)
)

*/

// P_PRE0 to P_FRAMEERR1 are constants representing various protocol states or phases in an iota enumeration.
const (
	P_PRE0 = 0 + iota
	P_PRE1
	P_PRE2
	P_READY
	P_EOI
	P_EOIw
	P_BIT0
	P_BIT0w
	P_BIT1
	P_BIT1w
	P_BIT2
	P_BIT2w
	P_BIT3
	P_BIT3w
	P_BIT4
	P_BIT4w
	P_BIT5
	P_BIT5w
	P_BIT6
	P_BIT6w
	P_BIT7
	P_BIT7w
	P_DONE0
	P_DONE1
	P_FRAMEERR0
	P_FRAMEERR1
)

// P_TALKING represents a state or flag with the value 0x20.
// P_LISTENING represents a state or flag with the value 0x40.
// P_ATN represents a state or flag with the value 0x80.
const (
	P_TALKING   = uint8(0x20)
	P_LISTENING = uint8(0x40)
	P_ATN       = uint8(0x80)
)

const stateLast = 0xf

// Protocol represents a structure for managing IEC protocol-based communication and its state machine.
type Protocol struct {
	*component.BaseComponent
	iec           references.IIec
	device        references.IIecProtocolDevice
	quartz        references.IQuartzSocket
	gs            *Global
	cfg           *config.Config
	deviceNumber  uint8
	stateMachine  uint8
	flags         uint8
	primary       uint8
	secondaryPrev uint8
	secondary     uint8
	timeout       uint64
	byte          uint8
	state         [stateLast + 1]uint8
}

// NewProtocol creates a new Protocol instance, initializes it with the provided parameters, and registers it with the parent.
func NewProtocol(parent references.IComponent, suffix string, q references.IQuartzSocket, deviceNumber uint8, device references.IIecProtocolDevice) *Protocol {
	p := &Protocol{
		BaseComponent: component.NewBaseComponent("iec_protocol", suffix),
		quartz:        q,
		gs:            _gs,
		iec:           nil,
		device:        device,
		deviceNumber:  deviceNumber,
	}
	component.Register(parent, p)
	return p
}

// Setup initializes the Protocol with the provided IEC interface and configuration settings.
func (v *Protocol) Setup(iec references.IIec, cfg *config.Config) {
	v.iec = iec
	v.cfg = cfg
}

// Reset resets the internal state of the Protocol to its initial configuration.
func (v *Protocol) Reset() {
	v.flags = 0
	v.timeout = 0
	for i := 0; i < len(v.state); i++ {
		v.state[i] = 0
	}
	v.transmitData(IECBUS_DEVICE_WRITE_CLK | IECBUS_DEVICE_WRITE_DATA)
}

// Ready checks if the protocol is ready for operation and returns a boolean value indicating readiness.
func (v *Protocol) Ready() bool {
	return true
}

// GetDeviceNumber returns the device number associated with the Protocol instance.
func (v *Protocol) GetDeviceNumber() uint8 {
	return v.deviceNumber
}

// AtnStateChanged toggles the ATN (Attention) state and prints the current state as "ON" or "OFF".
func (v *Protocol) AtnStateChanged(atn bool) {
	if !atn {
		fmt.Println("ATN STATE CHANGED", "ON")
	} else {
		fmt.Println("ATN STATE CHANGED", "OFF")
	}
}

// BusStateChanged handles changes in the communication bus state and updates the Protocol's behavior accordingly.
func (v *Protocol) BusStateChanged(uint8) {
	//TODO REMOVE
}

// Emulate handles the state transitions and bus communication logic for the Protocol according to the IEC device interface.
func (v *Protocol) Emulate() {
	bus := v.iec.PeripheralRead()

	//log.Printf("serial_iec_device_exec_main(%u, %u) F=%i, S=%i, ATN=%i CLK=%i DTA=%i", deviceNumber, clkValue, device.flags, device.state, (bus & IECBUS_DEVICE_READ_ATN) ? 1 : 0, (bus & IECBUS_DEVICE_READ_CLK) ? 1 : 0, (bus & IECBUS_DEVICE_READ_DATA) ? 1 : 0)
	if !conversion.Uint8ToBool(v.getFlags()&P_ATN) && !conversion.Uint8ToBool(bus&IECBUS_DEVICE_READ_ATN) {
		//Falling flank on ATN (bus master addressing all devices) */
		v.setStateMachine(P_PRE0)
		v.setFlags(v.getFlags() | P_ATN)
		v.setPrimary(0)
		v.setSecondaryPrev(v.getSecondary())
		v.setSecondary(0)
		v.setTimeout(100)
		//P_PRE0_100_time = time.Now()

		//Set DATA=0("I am here").If nobody on the bus does this within 1 ms, bus-master will assume that "DeviceAdapter not present"
		v.transmitData(IECBUS_DEVICE_WRITE_CLK)
	} else if conversion.Uint8ToBool(v.getFlags()&P_ATN) && conversion.Uint8ToBool(bus&IECBUS_DEVICE_READ_ATN) {
		//Rising flank on ATN (bus master finished addressing all devices) */
		v.setFlags(v.getFlags() & ^P_ATN)

		if (v.getPrimary() == 0x20+v.deviceNumber) || (v.getPrimary() == 0x40+v.deviceNumber) {
			if (v.getSecondary() & 0xf0) == 0x60 {
				switch v.getPrimary() & 0xf0 {
				case 0x20:
					v.device.Listen(v.getSecondary())
					v.print("after device.Listen", bus)
				case 0x40:
					v.device.Talk(v.getSecondary())
				}
			} else if (v.getSecondary() & 0xf0) == 0xe0 {
				v.gs.SetState(0)
				state := v.device.Close(v.getSecondary())
				v.setState(v.getSecondary(), state)
				//device.setState(device.getSecondary(), v.gs.getState())
			} else if (v.getSecondary() & 0xf0) == 0xf0 {
				//device.open() will not actually open the file (since we don't have a filename yet) but just set things up so that
				//the characters passed to device.
				//write() before the next call to device.unlisten() will be interpreted as the filename.
				//The file will actually be opened during the next call to device.unlisten()
				v.gs.SetState(0)
				state := v.device.Open(v.getSecondary())
				v.print("after device.Open", bus)
				v.setState(v.getSecondary(), state)
				//device.setState(device.getSecondary(), v.gs.getState())
			}

			if v.getPrimary() == 0x20+v.deviceNumber {
				//We were told to listen
				v.setFlags(v.getFlags() & ^P_TALKING)
				//st!=0 means that the previous OPEN command failed, i.e. we could not open a file for writing.
				//In that case, ignore the "LISTEN" request which will signal the error to the sender
				if v.getState(v.getSecondary()) == 0 {
					v.setFlags(v.getFlags() | P_LISTENING)
					v.setStateMachine(P_PRE1)
					log.Printf("device %d start listening", v.deviceNumber)
				}
				//set DATA=0 ("I am here")
				v.transmitData(IECBUS_DEVICE_WRITE_CLK)
			} else if v.getPrimary() == 0x40+v.deviceNumber {
				//We were told to talk
				v.setFlags(v.getFlags() & ^P_LISTENING)
				v.setFlags(v.getFlags() | P_TALKING)
				v.setStateMachine(P_PRE0)
				log.Printf("device %d start talking", v.deviceNumber)
			}
		} else if (v.getPrimary() == 0x3f) && conversion.Uint8ToBool(v.getFlags()&P_LISTENING) {
			//All devices were told to stop listening
			v.setFlags(v.getFlags() & ^P_LISTENING)
			log.Printf("device %d stop listening", v.deviceNumber)

			//If this is an UNLISTEN that followed an OPEN (0x2_ 0xf_), then
			//device.unlisten will try to open the file with the filename that
			//was received in between the OPEN and now.
			//If the file cannot be opened, it will set st != 0.
			v.gs.SetState(v.getState(v.getSecondaryPrev()))
			v.device.Unlisten(v.getSecondaryPrev())
			v.setState(v.getSecondaryPrev(), v.gs.GetState())
		} else if (v.getPrimary() == 0x5f) && conversion.Uint8ToBool(v.getFlags()&P_TALKING) {
			//All devices were told to stop talking
			v.device.Untalk(v.getSecondaryPrev())
			v.setFlags(v.getFlags() & ^P_TALKING)
			log.Printf("device %d stop talking", v.deviceNumber)
		}

		if !conversion.Uint8ToBool(v.getFlags() & (P_LISTENING | P_TALKING)) {
			//We're neither listening nor talking => make sure we're not holding DATA  or CLOCK line to 0 )
			v.transmitData(IECBUS_DEVICE_WRITE_CLK | IECBUS_DEVICE_WRITE_DATA)
		}
	}

	if conversion.Uint8ToBool(v.getFlags() & (P_ATN | P_LISTENING)) {
		v.doListen(bus)
	} else if conversion.Uint8ToBool(v.getFlags() & P_TALKING) {
		v.doTalk(bus)
	}
}

// doListen handles the state transitions for the device during the listening phase on the IEC bus based on the current clock and data signals.
// It ensures proper synchronization, processes incoming data or commands, and acknowledges frames as needed or signals errors.
func (v *Protocol) doListen(bus uint8) {
	//We are either under ATN or in "listening" mode
	//fmt.Println("State:", clkValue, device.getStateMachine())
	switch v.getStateMachine() {
	case P_PRE0:
		//Ignore anything that happens during the first 100 us after falling
		//flank on ATN (other devices may have been sending and need some time to set CLK=1)
		if v.timeoutExpired() {
			//test := time.Since(P_PRE0_100_time)
			//fmt.Println("P_PRE0_100_time", test)
			v.setStateMachine(P_PRE1)
		}
	case P_PRE1:
		//Make sure CLK=0 so we actually detect a rising flank instate P_PRE2
		if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
			v.setStateMachine(P_PRE2)
		}
	case P_PRE2:
		// wait for rising flank on CLK ("ready-to-send")
		if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
			//React by setting DATA=1 ("ready-for-data")
			v.transmitData(IECBUS_DEVICE_WRITE_CLK | IECBUS_DEVICE_WRITE_DATA)
			v.setTimeout(200)
			v.setStateMachine(P_READY)
		}
	case P_READY:
		if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
			//Sender set CLK=0, is about to send first bit
			v.setStateMachine(P_BIT0)
		} else if !conversion.Uint8ToBool(v.getFlags()&P_ATN) && v.timeoutExpired() {
			//Sender did not set CLK=0 within 200 us after we set DATA=1 => it is signaling EOI
			//(not so if we are under ATN) acknowledge we received it by setting DATA=0 for 60us
			log.Printf("device %d got EOI on channel %d", v.deviceNumber, v.getSecondary()&0x0f)
			v.transmitData(IECBUS_DEVICE_WRITE_CLK)
			v.setStateMachine(P_EOI)
			v.setTimeout(60)
		}
	case P_EOI:
		if v.timeoutExpired() {
			//Set DATA back to 1 and wait for sender to set CLK=0
			v.transmitData(IECBUS_DEVICE_WRITE_CLK | IECBUS_DEVICE_WRITE_DATA)
			v.setStateMachine(P_EOIw)
		}
	case P_EOIw:
		if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
			//Sender set CLK=0, is about to send first bit
			v.setStateMachine(P_BIT0)
		}
	case P_BIT0, P_BIT1, P_BIT2, P_BIT3, P_BIT4, P_BIT5, P_BIT6, P_BIT7:
		if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
			//Sender set CLK=1, signaling that the DATA line represents a valid bit
			bit := uint8(1 << ((int(v.getStateMachine()) - P_BIT0) / 2))
			p1 := v.getByte() & ^bit
			p2 := uint8(0)
			if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
				p2 = bit
			}
			v.setByte(p1 | p2)
			//Go to associated P_BIT(n)w state, waiting for sender to set CLK=0
			v.setStateMachineNext()
		}
	case P_BIT0w, P_BIT1w, P_BIT2w, P_BIT3w, P_BIT4w, P_BIT5w, P_BIT6w:
		if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
			//Sender set CLK=0. go to P_BIT(n+1) state to receive the next bit
			v.setStateMachineNext()
		}
	case P_BIT7w:
		if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
			//Sender set CLK=0 and this was the last bit
			log.Printf("device %d received : 0x%02x (%c)", v.deviceNumber, v.getByte(), v.getByte())
			if conversion.Uint8ToBool(v.getFlags() & P_ATN) {
				//We are currently receiving under ATN. Store the first two bytes received (contain primary and secondary address)
				if v.getPrimary() == 0 {
					v.setPrimary(v.getByte())
				} else if v.getSecondary() == 0 {
					v.setSecondary(v.getByte())
				}
				if v.getPrimary() != 0x3f && v.getPrimary() != 0x5f && ((uint8(v.getPrimary()) & 0x1f) != v.deviceNumber) {
					//This is NOT a UNLISTEN (0x3f) or UNTALK (0x5f) command and the primary address is not ours =>
					//Don't acknowledge the frame and stop listening. If all devices on the bus do this, the bus-master knows that "DeviceAdapter not present"
					v.setStateMachine(P_DONE0)
				} else {
					//Acknowledge frame by setting DATA=0
					v.transmitData(IECBUS_DEVICE_WRITE_CLK)
					//repeat from P_PRE2 (we know that CLK=0 so no need to go to P_PRE1)
					v.setStateMachine(P_PRE2)
				}
			} else if conversion.Uint8ToBool(v.getFlags() & P_LISTENING) {
				//We are currently listening for data pass received byte on to the upper level
				log.Printf("device %d received 0x%02x (%c) on channel %d", v.deviceNumber, v.getByte(), v.getByte(), v.getSecondary())
				v.gs.SetState(v.getState(v.getSecondary()))
				state := v.device.Write(v.getSecondary(), v.getByte())
				v.setState(v.getSecondary(), state)
				//device.setState(device.getSecondary(), v.gs.getState())

				if v.getState(v.getSecondary()) != 0 {
					//There was an error during iec_bus_write => stop listening. This will signal an error condition to the sender
					v.setStateMachine(P_DONE0)
				} else {
					//Acknowledge frame by setting DATA=0
					v.transmitData(IECBUS_DEVICE_WRITE_CLK)
					//repeat from P_PRE2 (we know that CLK=0 so no need to go to P_PRE1)
					v.setStateMachine(P_PRE2)
				}
			}
		}
	case P_DONE0:
		fmt.Println("We're just waiting for the bus-master to set ATN back to 1")
		//We're just waiting for the bus-master to set ATN back to 1
	default:
		panic("unhandled default case")
	}
}

// doTalk handles the communication protocol for transmitting data between devices using a state machine.
// It manages the transition between various states, timing, and signaling during the data transmission process.
func (v *Protocol) doTalk(bus uint8) {
	switch v.getStateMachine() {
	case P_PRE0:
		if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
			//Bus-master set CLK=1 (and before that should have set DATA=0)
			//we are getting ready for role reversal.Set CLK=0,DATA=1
			v.transmitData(IECBUS_DEVICE_WRITE_DATA)
			v.setStateMachine(P_PRE1)
			v.setTimeout(80)
		}
	case P_PRE1:
		if v.timeoutExpired() {
			//Signal "ready-to-send" (CLK=1)
			v.transmitData(IECBUS_DEVICE_WRITE_CLK | IECBUS_DEVICE_WRITE_DATA)
			v.setStateMachine(P_READY)
		}
	case P_READY:
		if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
			//Receiver signaled "ready-for-data" (DATA=1)
			v.gs.SetState(v.getState(v.getSecondary()))
			b, state := v.device.Read(v.getSecondary())
			v.setByte(b)
			v.setState(v.getSecondary(), state)
			//device.gs.setState(device.getSecondary(), v.gs.getState())
			if v.getState(v.getSecondary()) == 0 {
				//At least two bytes left to send. Go on to send the first bit.
				v.setStateMachine(P_BIT0)
				//no need to wait before sending the first bit
				v.setTimeout(0)
			} else if v.getState(v.getSecondary()) == 0x40 {
				//Only this byte left to send => signal EOI by keeping CLK=1
				log.Printf("device %d signaling EOI on channel %d", v.deviceNumber, v.getSecondary())
				v.setStateMachine(P_EOI)
			} else {
				//There was some kind of error; we have nothing to send.
				//Just stop talking and wait for ATN (This will produce a "File not found" when loading)
				v.setFlags(v.getFlags() & ^P_TALKING)
			}
		}
	case P_EOI:
		if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
			//Receiver set DATA=0, first part of acknowledging the EOI
			v.setStateMachine(P_EOIw)
		}
	case P_EOIw:
		if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
			//Receiver set DATA=1, final part of acknowledging the EOI. Go on to send first bit
			v.setStateMachine(P_BIT0)
			//no need to wait before sending the first bit
			v.setTimeout(0)
		}
	case P_BIT0, P_BIT1, P_BIT2, P_BIT3, P_BIT4, P_BIT5, P_BIT6, P_BIT7:
		if v.timeoutExpired() {
			//60 us have passed since we set CLK=1 to signal "data valid" for the previous bit.
			//Pull CLK=0 and put the next bit out of DATA.
			bit := uint8(1 << ((int(v.getStateMachine()) - P_BIT0) / 2))
			res := uint8(0)
			if conversion.Uint8ToBool(v.getByte() & bit) {
				res = IECBUS_DEVICE_WRITE_DATA
			}
			v.transmitData(res)
			//Go to associated P_BIT(n)w state
			v.setTimeout(60)

			v.setStateMachineNext()
		}
	case P_BIT0w, P_BIT1w, P_BIT2w, P_BIT3w, P_BIT4w, P_BIT5w, P_BIT6w, P_BIT7w:
		if v.timeoutExpired() {
			//60 us have passed since we pulled CLK=0 and put the current bit on DATA.
			//set CLK=1, keeping data as it is (this signals "data valid" to the receiver)
			if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
				v.transmitData(IECBUS_DEVICE_WRITE_CLK | IECBUS_DEVICE_WRITE_DATA)
			} else {
				v.transmitData(IECBUS_DEVICE_WRITE_CLK)
			}
			//Go to associated P_BIT(n+1) state to send the next bit.
			//If this was the final bit, then the next state is P_DONE0
			v.setTimeout(60)
			v.setStateMachineNext()
		}
	case P_DONE0:
		if v.timeoutExpired() {
			//60 us have passed since we set CLK=1 to signal "data valid" for the final bit.
			//Pull CLK=0 and set DATA=1.This prepares for the receiver acknowledgement.
			v.transmitData(IECBUS_DEVICE_WRITE_DATA)
			v.setTimeout(1000)
			v.setStateMachine(P_DONE1)
		}
	case P_DONE1:
		if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
			//Receiver set DATA=0, acknowledging the frame
			log.Printf("device %d sent 0x%02x (%c) on channel %d", v.deviceNumber, v.getByte(), v.getByte(), v.getSecondary())
			if v.getState(v.getSecondary()) == 0x40 {
				//This was the last byte => stop talking.This leaves us waiting for ATN.
				v.setFlags(v.getFlags() & ^P_TALKING)
				v.setState(v.getSecondary(), 0)
				//Release the CLOCK line to 1
				v.transmitData(IECBUS_DEVICE_WRITE_CLK | IECBUS_DEVICE_WRITE_DATA)
			} else {
				//There is at least one more byte to send Start over from P_PRE1
				v.setTimeout(0)
				v.setStateMachine(P_PRE1)
			}
		} else if v.timeoutExpired() {
			//We didn't receive an acknowledgement within 1 ms.Set CLOCK=0 and after 100 us back to CLOCK=1
			log.Printf("device %d got NACK on channel %d", v.deviceNumber, v.getSecondary())
			v.transmitData(IECBUS_DEVICE_WRITE_CLK | IECBUS_DEVICE_WRITE_DATA)
			v.setTimeout(100)
			v.setStateMachine(P_FRAMEERR0)
		}
	case P_FRAMEERR0:
		if v.timeoutExpired() {
			//Finished 1-0-1 sequence of CLOCK signal
			//to acknowledge the frame-error.Now wait for sender to set DATA=0 so we can continue.
			v.transmitData(IECBUS_DEVICE_WRITE_DATA)
			v.setStateMachine(P_FRAMEERR1)
		}
	case P_FRAMEERR1:
		if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
			// sender set DATA=0, we can retry to send the byte
			v.setTimeout(0)
			v.setStateMachine(P_PRE1)
		}
	default:
		panic("unhandled default case")
	}
}

// setFlags sets the flags property in the Protocol struct to the provided uint8 value f.
func (v *Protocol) setFlags(f uint8) {
	v.flags = f
}

// getFlags retrieves the current flags value of the Protocol instance.
func (v *Protocol) getFlags() uint8 {
	return v.flags
}

// setByte sets the value of the byte field in the Protocol struct.
func (v *Protocol) setByte(b uint8) {
	v.byte = b
}

// getByte retrieves the current byte value stored in the Protocol struct.
func (v *Protocol) getByte() uint8 {
	return v.byte
}

// setStateMachine sets the state machine value to the specified uint8 value.
func (v *Protocol) setStateMachine(m uint8) {
	v.stateMachine = m
}

// getStateMachine retrieves the current value of the stateMachine field in the Protocol instance.
func (v *Protocol) getStateMachine() uint8 {
	return v.stateMachine
}

// setStateMachineNext increments the state machine counter by 1.
func (v *Protocol) setStateMachineNext() {
	v.stateMachine++
}

// setPrimary sets the primary address of the protocol to the specified value.
func (v *Protocol) setPrimary(p uint8) {
	v.primary = p
}

// getPrimary returns the primary address of the Protocol instance as an unsigned 8-bit integer.
func (v *Protocol) getPrimary() uint8 {
	return v.primary
}

// setSecondary sets the secondary value of the protocol to the provided uint8 value.
func (v *Protocol) setSecondary(s uint8) {
	v.secondary = s
}

// getSecondary retrieves the current value of the `secondary` field in the Protocol instance.
func (v *Protocol) getSecondary() uint8 {
	return v.secondary
}

// setSecondaryPrev sets the value of the secondaryPrev field in the Protocol struct.
func (v *Protocol) setSecondaryPrev(s uint8) {
	v.secondaryPrev = s
}

// getSecondaryPrev retrieves the value of the secondaryPrev field, representing the previous secondary state in the protocol.
func (v *Protocol) getSecondaryPrev() uint8 {
	return v.secondaryPrev
}

// setTimeout sets the timeout value for the protocol instance to the specified value.
func (v *Protocol) setTimeout(offset uint64) {
	val := v.quartz.ToUSec(offset)
	v.timeout = v.quartz.Cycle() + val
}

func (v *Protocol) timeoutExpired() bool {
	b := v.quartz.Cycle() >= v.timeout
	return b
}

// setState updates the protocol state for the given index with the specified state value, using a masked index.
func (v *Protocol) setState(idx uint8, s uint8) {
	x := idx & stateLast
	v.state[x] = s
}

// getState retrieves the state value at the specified index, masked to fit within the bounds of the state array.
func (v *Protocol) getState(idx uint8) uint8 {
	x := idx & stateLast
	return v.state[x]
}

func (v *Protocol) transmitData(data uint8) {
	v.iec.PeripheralWrite(v.deviceNumber, data)
}

func (v *Protocol) print(id string, bus uint8) {
	fmt.Printf("%s -> bus: %d, stateMachine: %d, flags: %d, primary: %d, secondary: %d, secondaryPrev: %d\n", id, bus, v.stateMachine, v.flags, v.primary, v.secondary, v.secondaryPrev)
}
