package iecprotocol

import (
	"github.com/markel1974/c64emu/src/common/conversion"
	"github.com/markel1974/c64emu/src/components/board"
	"github.com/markel1974/c64emu/src/components/iec/iecdevice"
	"github.com/markel1974/c64emu/src/config"
)

//serial-iec-device.c
//static void serial_iec_device_exec_main(unsigned int devnr, CLOCK clk_value)

const (
	IECBUS_DEVICE_READ_DATA = uint8(0x01)
	IECBUS_DEVICE_READ_CLK  = uint8(0x04)
	IECBUS_DEVICE_READ_ATN  = uint8(0x80)

	IECBUS_DEVICE_ATNA = uint8(0x10)

	IECBUS_DEVICE_WRITE_CLK  = uint8(0x40)
	IECBUS_DEVICE_WRITE_DATA = uint8(0x80)
)

const (
	P_PRE0 = iota
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

const (
	P_TALKING   = uint8(0x20)
	P_LISTENING = uint8(0x40)
	P_ATN       = uint8(0x80)
)

type Protocol struct {
	*board.BaseComponent
	iec    iecdevice.IIec
	gs     *Global
	device *DeviceAdapter
	cfg    *config.Config
}

func NewProtocol(parent board.IComponent, suffix string, deviceNumber uint8, device iecdevice.IIecProtocolDevice) *Protocol {
	p := &Protocol{
		BaseComponent: board.NewBaseComponent("iec_protocol", suffix),
		gs:            _gs,
		iec:           nil,
		device:        nil,
	}
	board.Register(parent, p)
	p.device = NewDeviceAdapter(p, "", deviceNumber, p.gs, device)
	return p
}

func (v *Protocol) Setup(iec iecdevice.IIec, cfg *config.Config) {
	v.iec = iec
	v.cfg = cfg
}

func (v *Protocol) Reset() {
	//TODO
}

func (v *Protocol) Ready() bool {
	return true
}

func (v *Protocol) GetDeviceNumber() uint8 {
	return v.device.GetDeviceNumber()
}

func (v *Protocol) AtnStateChanged(bool) {
	//Nothing TO DO
}

func (v *Protocol) BusStateChanged(uint8) {
	//TODO REMOVE
}

func (v *Protocol) Emulate(clkValue uint64) {
	//Read bus
	bus := v.iec.PeripheralRead()
	device := v.device
	deviceNumber := v.device.GetDeviceNumber()

	//log.Printf("serial_iec_device_exec_main(%u, %u) F=%i, S=%i, ATN=%i CLK=%i DTA=%i", deviceNumber, clkValue, device.flags, device.state, (bus & IECBUS_DEVICE_READ_ATN) ? 1 : 0, (bus & IECBUS_DEVICE_READ_CLK) ? 1 : 0, (bus & IECBUS_DEVICE_READ_DATA) ? 1 : 0)

	if !conversion.Uint8ToBool(device.flags&P_ATN) && !conversion.Uint8ToBool(bus&IECBUS_DEVICE_READ_ATN) {
		//Falling flank on ATN (bus master addressing all devices) */
		device.SetStateMachine(P_PRE0)
		device.flags |= P_ATN
		device.SetPrimary(0)
		device.SetSecondaryPrev(device.GetSecondary())
		device.SetSecondary(0)
		device.SetTimeout(clkValue + US2CYCLES(100))

		//Set DATA=0("I am here").If nobody on the bus does this within 1 ms, bus-master will assume that "DeviceAdapter not present"
		v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK)
	} else if conversion.Uint8ToBool(device.flags&P_ATN) && conversion.Uint8ToBool(bus&IECBUS_DEVICE_READ_ATN) {
		//Rising flank on ATN (bus master finished addressing all devices) */
		device.flags &= ^P_ATN

		if (device.GetPrimary() == 0x20+deviceNumber) || (device.GetPrimary() == 0x40+deviceNumber) {
			if (device.GetSecondary() & 0xf0) == 0x60 {
				switch device.GetPrimary() & 0xf0 {
				case 0x20:
					device.Listen(device.GetSecondary())
				case 0x40:
					device.Talk(device.GetSecondary())
				}
			} else if (device.GetSecondary() & 0xf0) == 0xe0 {
				v.gs.SetState(0)
				device.Close(device.GetSecondary())
				device.SetState(device.GetSecondary(), v.gs.GetState())
			} else if (device.GetSecondary() & 0xf0) == 0xf0 {
				//device.open() will not actually open the file (since we don't have a filename yet) but just set things up so that
				//the characters passed to device.
				//write() before the next call to device.unlisten() will be interpreted as the filename.
				//The file will actually be opened during the next call to device.unlisten()
				v.gs.SetState(0)
				device.Open(device.GetSecondary())
				device.SetState(device.GetSecondary(), v.gs.GetState())
			}

			if device.GetPrimary() == 0x20+deviceNumber {
				//We were told to listen
				device.flags &= ^P_TALKING
				//st!=0 means that the previous OPEN command failed, i.e. we could not open a file for writing.
				//In that case, ignore the "LISTEN" request which will signal the error to the sender
				if device.GetState(device.GetSecondary()) == 0 {
					device.flags |= P_LISTENING
					device.SetStateMachine(P_PRE1)
					//log.Printf("device %i start listening", deviceNumber)
				}
				//set DATA=0 ("I am here")
				v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK)
			} else if device.GetPrimary() == 0x40+deviceNumber {
				//We were told to talk
				device.flags &= ^P_LISTENING
				device.flags |= P_TALKING
				device.SetStateMachine(P_PRE0)
				//log.Printf("device %i start talking", deviceNumber)
			}
		} else if (device.GetPrimary() == 0x3f) && conversion.Uint8ToBool(device.flags&P_LISTENING) {
			//All devices were told to stop listening
			device.flags &= ^P_LISTENING
			//log.Printf("device %i stop listening", deviceNumber)

			//If this is an UNLISTEN that followed an OPEN (0x2_ 0xf_), then
			//device.unlisten will try to open the file with the filename that
			//was received in between the OPEN and now.
			//If the file cannot be opened, it will set st != 0.
			v.gs.SetState(device.GetState(device.GetSecondaryPrev()))
			device.Unlisten(device.GetSecondaryPrev())
			device.SetState(device.GetSecondaryPrev(), v.gs.GetState())
		} else if (device.GetPrimary() == 0x5f) && conversion.Uint8ToBool(device.flags&P_TALKING) {
			//All devices were told to stop talking
			device.Untalk(device.GetSecondaryPrev())
			device.flags &= ^P_TALKING
			//log.Printf("device %i stop talking", deviceNumber)
		}

		if !conversion.Uint8ToBool(device.flags & (P_LISTENING | P_TALKING)) {
			//We're neither listening nor talking => make sure we're not holding DATA  or CLOCK line to 0
			v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK|IECBUS_DEVICE_WRITE_DATA)
		}
	}

	if conversion.Uint8ToBool(device.flags & (P_ATN | P_LISTENING)) {
		//We are either under ATN or in "listening" mode
		switch device.GetStateMachine() {
		case P_PRE0:
			//Ignore anything that happens during the first 100 us after falling
			//flank on ATN (other devices may have been sending and need some time to set CLK=1)
			if clkValue >= device.GetTimeout() {
				device.SetStateMachine(P_PRE1)
			}
		case P_PRE1:
			//Make sure CLK=0 so we actually detect a rising flank instate P_PRE2
			if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
				device.SetStateMachine(P_PRE2)
			}
		case P_PRE2:
			// wait for rising flank on CLK ("ready-to-send")
			if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
				//React by setting DATA=1 ("ready-for-data")
				v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK|IECBUS_DEVICE_WRITE_DATA)
				device.SetTimeout(clkValue + US2CYCLES(200))
				device.SetStateMachine(P_READY)
			}
		case P_READY:
			if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
				//Sender set CLK=0, is about to send first bit
				device.SetStateMachine(P_BIT0)
			} else if !conversion.Uint8ToBool(device.flags&P_ATN) && (clkValue >= device.GetTimeout()) {
				//Sender did not set CLK=0 within 200 us after we set DATA=1 => it is signaling EOI (not so if we are under ATN) acknowledge we received it by setting DATA=0 for 60us
				//log.Printf("device %i got EOI on channel %i", deviceNumber, device.secondary & 0x0f)
				v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK)
				device.SetStateMachine(P_EOI)
				device.SetTimeout(clkValue + US2CYCLES(60))
			}
		case P_EOI:
			if clkValue >= device.GetTimeout() {
				//Set DATA back to 1 and wait for sender to set CLK=0
				v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK|IECBUS_DEVICE_WRITE_DATA)
				device.SetStateMachine(P_EOIw)
			}
		case P_EOIw:
			if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
				//Sender set CLK=0, is about to send first bit
				device.SetStateMachine(P_BIT0)
			}
		case P_BIT0:
		case P_BIT1:
		case P_BIT2:
		case P_BIT3:
		case P_BIT4:
		case P_BIT5:
		case P_BIT6:
		case P_BIT7:
			if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
				//Sender set CLK=1, signaling that the DATA line represents a valid bit
				bit := uint8(1 << (uint8(device.GetStateMachine()-P_BIT0) / 2))
				p := uint8(0)
				if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
					p = 1
				}
				device.SetByte((device.GetByte() & ^bit) | p)
				//Go to associated P_BIT(n)w state, waiting for sender to set CLK=0
				device.SetStateMachineNext()
			}
		case P_BIT0w:
		case P_BIT1w:
		case P_BIT2w:
		case P_BIT3w:
		case P_BIT4w:
		case P_BIT5w:
		case P_BIT6w:
			if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
				//Sender set CLK=0. go to P_BIT(n+1) state to receive the next bit
				device.SetStateMachineNext()
			}
		case P_BIT7w:
			if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
				//Sender set CLK=0 and this was the last bit
				//log.Printf("device %i received : 0x%02x (%c)", deviceNumber, device.byte, device.byte)
				if conversion.Uint8ToBool(device.flags & P_ATN) {
					//We are currently receiving under ATN. Store the first two bytes received (contain primary and secondary address)
					if device.GetPrimary() == 0 {
						device.SetPrimary(device.GetByte())
					} else if device.GetSecondary() == 0 {
						device.SetSecondary(device.GetByte())
					}
					if device.GetPrimary() != 0x3f && device.GetPrimary() != 0x5f && ((uint8(device.GetPrimary()) & 0x1f) != deviceNumber) {
						//This is NOT a UNLISTEN (0x3f) or UNTALK (0x5f) command and the primary address is not ours =>
						//Don't acknowledge the frame and stop listening. If all devices on the bus do this, the bus-master knows that "DeviceAdapter not present"
						device.SetStateMachine(P_DONE0)
					} else {
						//Acknowledge frame by setting DATA=0
						v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK)
						//repeat from P_PRE2 (we know that CLK=0 so no need to go to P_PRE1)
						device.SetStateMachine(P_PRE2)
					}
				} else if conversion.Uint8ToBool(device.flags & P_LISTENING) {
					//We are currently listening for data pass received byte on to the upper level
					//log.Printf("device %i received 0x%02x (%c) on channel %i", deviceNumber, device.byte, isprint((unsigned char)device.byte) ? device.byte : '.', device.secondary & 0x0f)
					v.gs.SetState(device.GetState(device.GetSecondary()))
					device.Write(device.GetSecondary(), device.GetByte())
					device.SetState(device.GetSecondary(), v.gs.GetState())

					if device.GetState(device.GetSecondary()) != 0 {
						//There was an error during iec_bus_write => stop listening. This will signal an error condition to the sender
						device.SetStateMachine(P_DONE0)
					} else {
						//Acknowledge frame by setting DATA=0
						v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK)
						//repeat from P_PRE2 (we know that CLK=0 so no need to go to P_PRE1)
						device.SetStateMachine(P_PRE2)
					}
				}
			}
		case P_DONE0:
			//We're just waiting for the bus-master to set ATN back to 1
		default:
			panic("unhandled default case")
		}
	} else if conversion.Uint8ToBool(device.flags & P_TALKING) {
		//We are in "talking" mode
		switch device.GetStateMachine() {
		case P_PRE0:
			if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
				//Bus-master set CLK=1 (and before that should have set DATA=0)
				//we are getting ready for role reversal.Set CLK=0,DATA=1
				v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_DATA)
				device.SetStateMachine(P_PRE1)
				device.SetTimeout(clkValue + US2CYCLES(80))
			}
		case P_PRE1:
			if clkValue >= device.GetTimeout() {
				//Signal "ready-to-send" (CLK=1)
				v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK|IECBUS_DEVICE_WRITE_DATA)
				device.SetStateMachine(P_READY)
			}
		case P_READY:
			if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
				//Receiver signaled "ready-for-data" (DATA=1)
				v.gs.SetState(device.GetState(device.GetSecondary()))
				b := device.Read(device.GetSecondary())
				device.SetByte(b)
				device.SetState(device.GetSecondary(), v.gs.GetState())
				if device.GetState(device.GetSecondary()) == 0 {
					//At least two bytes left to send. Go on to send the first bit.
					device.SetStateMachine(P_BIT0)
					//no need to wait before sending the first bit
					device.SetTimeout(clkValue)
				} else if device.GetState(device.GetSecondary()) == 0x40 {
					//Only this byte left to send => signal EOI by keeping CLK=1
					//log.Printf(serial_iec_device_log,"device %i signaling EOI on channel %i", deviceNumber, device.secondary & 0x0f)
					device.SetStateMachine(P_EOI)
				} else {
					//There was some kind of error; we have nothing to send.
					//Just stop talking and wait for ATN (This will produce a "File not found" when loading)
					device.flags &= ^P_TALKING
				}
			}
		case P_EOI:
			if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
				//Receiver set DATA=0, first part of acknowledging the EOI
				device.SetStateMachine(P_EOIw)
			}
		case P_EOIw:
			if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
				//Receiver set DATA=1, final part of acknowledging the EOI. Go on to send first bit
				device.SetStateMachine(P_BIT0)
				//no need to wait before sending the first bit
				device.SetTimeout(clkValue)
			}
		case P_BIT0:
		case P_BIT1:
		case P_BIT2:
		case P_BIT3:
		case P_BIT4:
		case P_BIT5:
		case P_BIT6:
		case P_BIT7:
			if clkValue >= device.GetTimeout() {
				//60 us have passed since we set CLK=1 to signal "data valid" for the previous bit.
				//Pull CLK=0 and put the next bit out of DATA.
				bit := 1 << ((int(device.GetStateMachine()) - P_BIT0) / 2)
				res := uint8(0)
				if conversion.Uint8ToBool(device.GetByte() & uint8(bit)) {
					res = IECBUS_DEVICE_WRITE_DATA
				}
				v.iec.PeripheralWrite(deviceNumber, res)
				//v.device.PeripheralWrite(deviceNumber, (uint8_t)((device.byte & bit) ? IECBUS_DEVICE_WRITE_DATA : 0))

				//Go to associated P_BIT(n)w state
				device.SetTimeout(clkValue + US2CYCLES(60))
				device.SetStateMachineNext()
			}
		case P_BIT0w:
		case P_BIT1w:
		case P_BIT2w:
		case P_BIT3w:
		case P_BIT4w:
		case P_BIT5w:
		case P_BIT6w:
		case P_BIT7w:
			if clkValue >= device.GetTimeout() {
				//60 us have passed since we pulled CLK=0 and put the current bit on DATA.
				//set CLK=1, keeping data as it is (this signals "data valid" to the receiver)
				if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
					v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK|IECBUS_DEVICE_WRITE_DATA)
				} else {
					v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK)
				}
				//Go to associated P_BIT(n+1) state to send the next bit.
				//If this was the final bit, then the next state is P_DONE0
				device.SetTimeout(clkValue + US2CYCLES(60))
				device.SetStateMachineNext()
			}
		case P_DONE0:
			if clkValue >= device.GetTimeout() {
				//60 us have passed since we set CLK=1 to signal "data valid" for the final bit.
				//Pull CLK=0 and set DATA=1.This prepares for the receiver acknowledgement.
				v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_DATA)
				device.SetTimeout(clkValue + US2CYCLES(1000))
				device.SetStateMachine(P_DONE1)
			}
		case P_DONE1:
			if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
				//Receiver set DATA=0, acknowledging the frame
				//log.Printf("device %i sent 0x%02x (%c) on channel %i", deviceNumber, device.byte, device.byte, device.secondary & 0x0f)
				if device.GetState(device.GetSecondary()) == 0x40 {
					//This was the last byte => stop talking.This leaves us waiting for ATN.
					device.flags &= ^P_TALKING
					device.SetState(device.GetSecondary(), 0)
					//Release the CLOCK line to 1
					v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK|IECBUS_DEVICE_WRITE_DATA)
				} else {
					//There is at least one more byte to send Start over from P_PRE1
					device.SetTimeout(clkValue)
					device.SetStateMachine(P_PRE1)
				}
			} else if clkValue >= device.GetTimeout() {
				//We didn't receive an acknowledgement within 1 ms.Set CLOCK=0 and after 100 us back to CLOCK=1
				//log.Printf("device %i got NACK on channel %i", deviceNumber, device.secondary & 0x0f)
				v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK|IECBUS_DEVICE_WRITE_DATA)
				device.SetTimeout(clkValue + US2CYCLES(100))
				device.SetStateMachine(P_FRAMEERR0)
			}
		case P_FRAMEERR0:
			if clkValue >= device.GetTimeout() {
				//Finished 1-0-1 sequence of CLOCK signal
				//to acknowledge the frame-error.Now wait for sender to set DATA=0 so we can continue.
				v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_DATA)
				device.SetStateMachine(P_FRAMEERR1)
			}
		case P_FRAMEERR1:
			if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
				// sender set DATA=0, we can retry to send the byte
				device.SetTimeout(clkValue)
				device.SetStateMachine(P_PRE1)
			}
		default:
			panic("unhandled default case")
		}
	}
}

func US2CYCLES(v uint64) uint64 {
	return v
}
