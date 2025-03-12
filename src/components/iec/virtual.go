package iec

import (
	"github.com/markel1974/c64emu/src/common/conversion"
	"github.com/markel1974/c64emu/src/components/iec/virtualdrive"
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

//type IGlobalSt interface {
//	GetSt() uint8
//	SetSt() uint8
//}

type Device struct {
	//globalSt       IGlobalSt
	deviceNumber   uint8
	state          uint8
	flags          uint8
	primary        uint8
	secondary_prev uint8
	secondary      uint8
	timeout        uint64
	byte           uint8
	st             [0xff]uint8
}

func NewState(deviceNumber uint8) *Device {
	v := &Device{
		deviceNumber: deviceNumber,
		//globalSt: State,
	}
	return v
}

func (s *Device) Listen(b uint8) {
}

func (s *Device) Write(b uint8, c uint8, val uint8) uint8 {
	return val
}

func (s *Device) Read(b uint8, val uint8) (uint8, uint8) {
	return 0, val
}

func (s *Device) Talk(b uint8) {
}

func (s *Device) Close(b uint8, val uint8) uint8 {
	return val
}

func (s *Device) Open(b uint8, val uint8) uint8 {
	return val
}

func (s *Device) Unlisten(b uint8, val uint8) uint8 {
	return val
}

func (s *Device) Untalk(sec uint8) {
}

var serial_iec_device_state []Device

type Virtual struct {
	iec                  virtualdrive.IIec
	serial_iec_device_st uint8
}

func NewVirtual(iec virtualdrive.IIec) *Virtual {
	return &Virtual{
		iec:                  iec,
		serial_iec_device_st: 0,
	}
}

func (v *Virtual) Emulate(deviceNumber uint8, clkValue uint64) {
	iec := serial_iec_device_state[deviceNumber]

	//read bus
	bus := v.iec.PeripheralRead()

	//log.Printf("serial_iec_device_exec_main(%u, %u) F=%i, S=%i, ATN=%i CLK=%i DTA=%i", deviceNumber, clkValue, iec.flags, iec.state, (bus & IECBUS_DEVICE_READ_ATN) ? 1 : 0, (bus & IECBUS_DEVICE_READ_CLK) ? 1 : 0, (bus & IECBUS_DEVICE_READ_DATA) ? 1 : 0)

	if !conversion.Uint8ToBool(iec.flags&P_ATN) && !conversion.Uint8ToBool(bus&IECBUS_DEVICE_READ_ATN) {
		// falling flank on ATN (bus master addressing all devices) */
		iec.state = P_PRE0
		iec.flags |= P_ATN
		iec.primary = 0
		iec.secondary_prev = iec.secondary
		iec.secondary = 0
		iec.timeout = clkValue + US2CYCLES(100)

		// set DATA=0 ("I am here"). If nobody on the bus does this within 1ms, bus-master will assume that "Device not present"
		v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK)
	} else if conversion.Uint8ToBool(iec.flags&P_ATN) && conversion.Uint8ToBool(bus&IECBUS_DEVICE_READ_ATN) {
		// rising flank on ATN (bus master finished addressing all devices) */
		iec.flags &= ^P_ATN

		if (iec.primary == 0x20+deviceNumber) || (iec.primary == 0x40+deviceNumber) {
			if (iec.secondary & 0xf0) == 0x60 {
				switch iec.primary & 0xf0 {
				case 0x20:
					iec.Listen(iec.secondary)
				case 0x40:
					iec.Talk(iec.secondary)
				}
			} else if (iec.secondary & 0xf0) == 0xe0 {
				//v.set_st(0)
				val := iec.Close(iec.secondary, uint8(0))
				iec.st[iec.secondary&0x0f] = val
				//iec.st[iec.secondary&0x0f] = v.get_st()
			} else if (iec.secondary & 0xf0) == 0xf0 {
				// iec.open() will not actually open the file (since we don't have a filename yet) but just set things up so that
				// the characters passed to iec.write() before the next call to iec.unlisten() will be interpreted as the filename.
				// The file will actually be opened during the next call to iec.unlisten()
				//v.set_st(0)
				val := iec.Open(iec.secondary, uint8(0))
				//iec.st[iec.secondary&0x0f] = v.get_st()
				iec.st[iec.secondary&0x0f] = val
			}

			if iec.primary == 0x20+deviceNumber {
				// we were told to listen
				iec.flags &= ^P_TALKING
				// st!=0 means that the previous OPEN command failed, i.e. we could not open a file for writing.  In that case, ignore the "LISTEN" request which will signal the error to the sender
				if iec.st[iec.secondary&0x0f] == 0 {
					iec.flags |= P_LISTENING
					iec.state = P_PRE1
					//log.Printf("device %i start listening", deviceNumber)
				}
				//set DATA=0 ("I am here")
				v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK)
			} else if iec.primary == 0x40+deviceNumber {
				// we were told to talk
				iec.flags &= ^P_LISTENING
				iec.flags |= P_TALKING
				iec.state = P_PRE0
				//log.Printf("device %i start talking", deviceNumber)
			}
		} else if (iec.primary == 0x3f) && conversion.Uint8ToBool(iec.flags&P_LISTENING) {
			// all devices were told to stop listening
			iec.flags &= ^P_LISTENING
			//log.Printf("device %i stop listening", deviceNumber)

			//if this is an UNLISTEN that followed an OPEN (0x2_ 0xf_), then
			//iec.unlisten will try to open the file with the filename that
			//was received in between the OPEN and now.  If the file cannot be
			//opened, it will set st != 0.
			//v.set_st(iec.st[iec.secondary_prev&0x0f])
			val := iec.st[iec.secondary_prev&0x0f]
			val = iec.Unlisten(iec.secondary_prev, val)
			iec.st[iec.secondary_prev&0x0f] = val
			//iec.st[iec.secondary_prev&0x0f] = v.get_st()
		} else if (iec.primary == 0x5f) && conversion.Uint8ToBool(iec.flags&P_TALKING) {
			// all devices were told to stop talking
			iec.Untalk(iec.secondary_prev)
			iec.flags &= ^P_TALKING
			//log.Printf("device %i stop talking", deviceNumber)
		}

		if !conversion.Uint8ToBool(iec.flags & (P_LISTENING | P_TALKING)) {
			//we're neither listening nor talking => make sure we're not holding DATA  or CLOCK line to 0
			v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK|IECBUS_DEVICE_WRITE_DATA)
		}
	}

	if conversion.Uint8ToBool(iec.flags & (P_ATN | P_LISTENING)) {
		// we are either under ATN or in "listening" mode
		switch iec.state {
		case P_PRE0:
			//ignore anything that happens during first 100us after falling
			//flank on ATN (other devices may have been sending and need
			//some time to set CLK=1)
			if clkValue >= iec.timeout {
				iec.state = P_PRE1
			}
		case P_PRE1:
			// make sure CLK=0 so we actually detect a rising flank instate P_PRE2
			if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
				iec.state = P_PRE2
			}
		case P_PRE2:
			// wait for rising flank on CLK ("ready-to-send")
			if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
				// react by setting DATA=1 ("ready-for-data")
				v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK|IECBUS_DEVICE_WRITE_DATA)
				iec.timeout = clkValue + US2CYCLES(200)
				iec.state = P_READY
			}
		case P_READY:
			if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
				// sender set CLK=0, is about to send first bit
				iec.state = P_BIT0
			} else if !conversion.Uint8ToBool(iec.flags&P_ATN) && (clkValue >= iec.timeout) {
				// sender did not set CLK=0 within 200us after we set DATA=1 => it is signaling EOI (not so if we are under ATN) acknowledge we received it by setting DATA=0 for 60us
				//log.Printf("device %i got EOI on channel %i", deviceNumber, iec.secondary & 0x0f)
				v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK)
				iec.state = P_EOI
				iec.timeout = clkValue + US2CYCLES(60)
			}
		case P_EOI:
			if clkValue >= iec.timeout {
				// Set DATA back to 1 and wait for sender to set CLK=0 */
				v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK|IECBUS_DEVICE_WRITE_DATA)
				iec.state = P_EOIw
			}
		case P_EOIw:
			if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
				// sender set CLK=0, is about to send first bit
				iec.state = P_BIT0
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
				// sender set CLK=1, signaling that the DATA line represents a valid bit
				bit := uint8(1 << (uint8(iec.state-P_BIT0) / 2))
				p := uint8(0)
				if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
					p = 1
				}
				iec.byte = (iec.byte & ^bit) | p
				// go to associated P_BIT(n)w state, waiting for sender to set CLK=0
				iec.state++
			}
		case P_BIT0w:
		case P_BIT1w:
		case P_BIT2w:
		case P_BIT3w:
		case P_BIT4w:
		case P_BIT5w:
		case P_BIT6w:
			if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
				// sender set CLK=0. go to P_BIT(n+1) state to receive next bit
				iec.state++
			}
		case P_BIT7w:
			if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
				// sender set CLK=0 and this was the last bit
				// log.Printf("device %i received : 0x%02x (%c)", deviceNumber, iec.byte, iec.byte)
				if conversion.Uint8ToBool(iec.flags & P_ATN) {
					// We are currently receiving under ATN.  Store first two bytes received (contain primary and secondary address)
					if iec.primary == 0 {
						iec.primary = iec.byte
					} else if iec.secondary == 0 {
						iec.secondary = iec.byte
					}

					if iec.primary != 0x3f && iec.primary != 0x5f && ((uint8(iec.primary) & 0x1f) != deviceNumber) {
						// This is NOT a UNLISTEN (0x3f) or UNTALK (0x5f) command and the primary address is not ours =>
						// Don't acknowledge the frame and stop listening. If all devices on the bus do this, the bus-master knows that "Device not present"
						iec.state = P_DONE0
					} else {
						// Acknowledge frame by setting DATA=0
						v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK)
						// repeat from P_PRE2 (we know that CLK=0 so no need to go to P_PRE1)
						iec.state = P_PRE2
					}
				} else if conversion.Uint8ToBool(iec.flags & P_LISTENING) {
					// We are currently listening for data pass received byte on to the upper level
					// log.Printf("device %i received 0x%02x (%c) on channel %i", deviceNumber, iec.byte, isprint((unsigned char)iec.byte) ? iec.byte : '.', iec.secondary & 0x0f)
					//v.set_st(iec.st[iec.secondary&0x0f])
					val := iec.st[iec.secondary&0x0f]
					val = iec.Write(iec.secondary, iec.byte, val)
					iec.st[iec.secondary&0x0f] = val
					//ec.st[iec.secondary&0x0f] = v.get_st()

					if iec.st[iec.secondary&0x0f] != 0 {
						// there was an error during iec_bus_write => stop listening. This will signal an error condition to the sender
						iec.state = P_DONE0
					} else {
						// Acknowledge frame by setting DATA=0
						v.iec.PeripheralWrite(deviceNumber, uint8(IECBUS_DEVICE_WRITE_CLK))
						// repeat from P_PRE2 (we know that CLK=0 so no need to go to P_PRE1)
						iec.state = P_PRE2
					}
				}
			}
		case P_DONE0:
			// we're just waiting for the bus-master to set ATN back to 1
		default:
			panic("unhandled default case")
		}
	} else if conversion.Uint8ToBool(iec.flags & P_TALKING) {
		// we are in "talking" mode
		switch iec.state {
		case P_PRE0:
			if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_CLK) {
				// bus-master set CLK=1 (and before that should have set DATA=0) we are getting ready for role reversal. Set CLK=0, DATA=1
				v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_DATA)
				iec.state = P_PRE1
				iec.timeout = clkValue + US2CYCLES(80)
			}
		case P_PRE1:
			if clkValue >= iec.timeout {
				// signal "ready-to-send" (CLK=1)
				v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK|IECBUS_DEVICE_WRITE_DATA)
				iec.state = P_READY
			}
		case P_READY:
			if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
				// receiver signaled "ready-for-data" (DATA=1)
				//v.set_st(iec.st[iec.secondary&0x0f])
				val := iec.st[iec.secondary&0x0f]
				iec.byte, val = iec.Read(iec.secondary, val)
				iec.st[iec.secondary&0x0f] = val
				//iec.st[iec.secondary&0x0f] = v.get_st()
				if iec.st[iec.secondary&0x0f] == 0 {
					// at least two bytes left to send. Go on to send first bit.
					iec.state = P_BIT0
					// no need to wait before sending the first bit
					iec.timeout = clkValue
				} else if iec.st[iec.secondary&0x0f] == 0x40 {
					// only this byte left to send => signal EOI by keeping CLK=1
					//log.Printf(serial_iec_device_log,"device %i signaling EOI on channel %i", deviceNumber, iec.secondary & 0x0f)
					iec.state = P_EOI
				} else {
					// There was some kind of error, we have nothing to send.  Just stop talking and wait for ATN. (This will produce a "File not found" when loading)
					iec.flags &= ^P_TALKING
				}
			}
		case P_EOI:
			if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
				//receiver set DATA=0, first part of acknowledging the EOI
				iec.state = P_EOIw
			}
		case P_EOIw:
			if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
				// receiver set DATA=1, final part of acknowledging the EOI. Go on to send first bit
				iec.state = P_BIT0
				// no need to wait before sending the first bit
				iec.timeout = clkValue
			}
		case P_BIT0:
		case P_BIT1:
		case P_BIT2:
		case P_BIT3:
		case P_BIT4:
		case P_BIT5:
		case P_BIT6:
		case P_BIT7:
			if clkValue >= iec.timeout {
				// 60us have passed since we set CLK=1 to signal "data valid" for the previous bit. Pull CLK=0 and put next bit out on DATA. */
				bit := 1 << ((int(iec.state) - P_BIT0) / 2)
				res := uint8(0)
				if conversion.Uint8ToBool(iec.byte & uint8(bit)) {
					res = IECBUS_DEVICE_WRITE_DATA
				}
				v.iec.PeripheralWrite(deviceNumber, res)
				//v.iec.write(deviceNumber, (uint8_t)((iec.byte & bit) ? IECBUS_DEVICE_WRITE_DATA : 0))
				// go to associated P_BIT(n)w state
				iec.timeout = clkValue + US2CYCLES(60)
				iec.state++
			}
		case P_BIT0w:
		case P_BIT1w:
		case P_BIT2w:
		case P_BIT3w:
		case P_BIT4w:
		case P_BIT5w:
		case P_BIT6w:
		case P_BIT7w:
			if clkValue >= iec.timeout {
				// 60us have passed since we pulled CLK=0 and put the current bit on DATA. set CLK=1, keeping data as it is (this signals "data valid" to the receiver)
				if conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
					v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK|IECBUS_DEVICE_WRITE_DATA)
				} else {
					v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK)
				}
				// go to associated P_BIT(n+1) state to send the next bit. If this was the final bit then next state is P_DONE0
				iec.timeout = clkValue + US2CYCLES(60)
				iec.state++
			}
		case P_DONE0:
			if clkValue >= iec.timeout {
				// 60us have passed since we set CLK=1 to signal "data valid" for the final bit. Pull CLK=0 and set DATA=1. This prepares for the receiver acknowledgement. */
				v.iec.PeripheralWrite(deviceNumber, uint8(IECBUS_DEVICE_WRITE_DATA))
				iec.timeout = clkValue + uint64(US2CYCLES(1000))
				iec.state = P_DONE1
			}
		case P_DONE1:
			if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
				// Receiver set DATA=0, acknowledging the frame
				// log.Printf("device %i sent 0x%02x (%c) on channel %i", deviceNumber, iec.byte, iec.byte, iec.secondary & 0x0f)
				if iec.st[iec.secondary&0x0f] == 0x40 {
					// This was the last byte => stop talking.This leaves us waiting for ATN.
					iec.flags &= ^P_TALKING
					iec.st[iec.secondary&0x0f] = 0
					// Release the CLOCK line to 1
					v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK|IECBUS_DEVICE_WRITE_DATA)
				} else {
					// There is at least one more byte to send Start over from P_PRE1
					iec.timeout = clkValue
					iec.state = P_PRE1
				}
			} else if clkValue >= iec.timeout {
				//We didn't receive an acknowledgement within 1ms. Set CLOCK=0 and after 100us back to CLOCK=1
				//log.Printf("device %i got NACK on channel %i", deviceNumber, iec.secondary & 0x0f)
				v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_CLK|IECBUS_DEVICE_WRITE_DATA)
				iec.timeout = clkValue + US2CYCLES(100)
				iec.state = P_FRAMEERR0
			}
		case P_FRAMEERR0:
			if clkValue >= iec.timeout {
				// finished 1-0-1 sequence of CLOCK signal to acknowledge the frame-error.  Now wait for sender to set DATA=0 so we can continue. */
				v.iec.PeripheralWrite(deviceNumber, IECBUS_DEVICE_WRITE_DATA)
				iec.state = P_FRAMEERR1
			}
		case P_FRAMEERR1:
			if !conversion.Uint8ToBool(bus & IECBUS_DEVICE_READ_DATA) {
				// sender set DATA=0, we can retry to send the byte
				iec.timeout = clkValue
				iec.state = P_PRE1
			}
		default:
			panic("unhandled default case")
		}
	}
}

func US2CYCLES(v uint64) uint64 {
	return v
}

/*
func (v *Virtual) set_st1(b uint8) {
	v.serial_iec_device_st = b
}

func (v *Virtual) get_st1() uint8 {
	return v.serial_iec_device_st
}*/
