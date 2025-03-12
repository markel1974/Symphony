package iec

import (
	"github.com/markel1974/c64emu/src/components/iec/virtualdrive"
)

const (
	CmdData  = 0x60 // Data transfer
	CmdClose = 0xe0 // Close channel
	CmdOpen  = 0xf0 // Open channel
)

const (
	AtnListen   = 0x20
	AtnUnlisten = 0x30
	AtnTalk     = 0x40
	AtnUntalk   = 0x5
)

type IDrive interface {
	Close(v uint8) uint8
	Open(v uint8, data []uint8) uint8
	Write(v uint8, d uint8, b bool) uint8
	Read(v uint8) (uint8, uint8)
	Ready() bool
}

type VirtualBus struct {
	openData       []byte
	listener       IDrive // Pointer to active listener
	talker         IDrive // Pointer to active talker
	listenerActive bool   // Listener selected, listener_data is valid
	talkerActive   bool   // Talker selected, talker_data is valid
	listening      bool   // Last ATN was listen (to decide between sec_listen/sec_talk)
	receivedCmd    uint8  // Received command code ($x0)
	secAddr        uint8  // Received secondary address ($0x)
	virtualDrives  []IDrive
}

func NewEmulator() *VirtualBus {
	return &VirtualBus{}
}

func (vb *VirtualBus) Reset() {
	vb.listener = nil
	vb.talker = nil
	vb.listenerActive = false
	vb.talkerActive = false
	vb.listening = false
	vb.receivedCmd = 0
	vb.secAddr = 0
	vb.openData = nil
}

func (vb *VirtualBus) Out(data uint8, eoi bool) uint8 {
	if vb.listenerActive {
		if vb.receivedCmd == CmdOpen {
			return vb.openOut(data, eoi)
		}
		if vb.receivedCmd == CmdData {
			return vb.dataOut(data, eoi)
		}
		return virtualdrive.StTimeout
	} else {
		return virtualdrive.StTimeout
	}
}

func (vb *VirtualBus) OutATN(data uint8) uint8 {
	vb.receivedCmd = 0 // Command is sent with secondary address
	vb.secAddr = 0     // Command is sent with secondary address
	d := data & 0xf0
	switch d {
	case AtnListen:
		vb.listening = true
		return vb.listen(d)
	case AtnUnlisten:
		vb.listening = false
		return vb.unListen()
	case AtnTalk:
		vb.listening = false
		return vb.talk(d)
	case AtnUntalk:
		vb.listening = false
		return vb.unTalk()
	}
	return virtualdrive.StTimeout
}

func (vb *VirtualBus) OutSec(data uint8) uint8 {
	if vb.listening {
		if vb.listenerActive {
			vb.secAddr = data & 0x0f
			vb.receivedCmd = data & 0xf0
			return vb.secListen()
		}
	} else {
		if vb.talkerActive {
			vb.secAddr = data & 0x0f
			vb.receivedCmd = data & 0xf0
			return vb.secTalk()
		}
	}
	return virtualdrive.StTimeout
}

func (vb *VirtualBus) In() (uint8, uint8) {
	if vb.talkerActive && (vb.receivedCmd == CmdData) {
		return vb.dataIn()
	}
	return virtualdrive.StTimeout, 0
}

func (vb *VirtualBus) SetATN() {
	// Only needed for real Dispatcher
}

func (vb *VirtualBus) RelATN() {
	// Only needed for real Dispatcher
}

func (vb *VirtualBus) Turnaround() {
	// Only needed for real Dispatcher
}

func (vb *VirtualBus) Release() {
	// Only needed for real Dispatcher
}

func (vb *VirtualBus) listen(device uint8) uint8 {
	vb.listenerActive = false
	if device < 8 || device > 11 {
		return virtualdrive.StNotPresent
	}
	if vb.listener = vb.virtualDrives[device-8]; vb.listener == nil {
		return virtualdrive.StNotPresent
	}
	if !vb.listener.Ready() {
		return virtualdrive.StNotPresent
	}
	vb.listenerActive = true
	return virtualdrive.StOk
}

func (vb *VirtualBus) talk(device uint8) uint8 {
	vb.talkerActive = false
	if device < 8 || device > 11 {
		return virtualdrive.StNotPresent
	}
	if vb.talker = vb.virtualDrives[device-8]; vb.talker == nil {
		return virtualdrive.StNotPresent
	}
	if !vb.talker.Ready() {
		return virtualdrive.StNotPresent
	}
	vb.talkerActive = true
	return virtualdrive.StOk
}

func (vb *VirtualBus) unListen() uint8 {
	vb.listenerActive = false
	return virtualdrive.StOk
}

func (vb *VirtualBus) unTalk() uint8 {
	vb.talkerActive = false
	return virtualdrive.StOk
}

func (vb *VirtualBus) secListen() uint8 {
	switch vb.receivedCmd {
	case CmdOpen:
		vb.openData = nil
	case CmdClose:
		if vb.listener != nil {
			return vb.listener.Close(vb.secAddr)
		}
	}
	return virtualdrive.StOk
}

func (vb *VirtualBus) secTalk() uint8 {
	return virtualdrive.StOk
}

func (vb *VirtualBus) openOut(data uint8, eoi bool) uint8 {
	vb.openData = append(vb.openData, data)
	if eoi {
		if vb.listener != nil {
			return vb.listener.Open(vb.secAddr, vb.openData)
		}
	}
	return virtualdrive.StOk
}

func (vb *VirtualBus) dataOut(data uint8, eoi bool) uint8 {
	if vb.listener != nil {
		return vb.listener.Write(vb.secAddr, data, eoi)
	}
	return virtualdrive.StOk
}

func (vb *VirtualBus) dataIn() (uint8, uint8) {
	if vb.talker != nil {
		return vb.talker.Read(vb.secAddr)
	}
	return 0, 0
}

func (vb *VirtualBus) destroyVirtualDrive(vd IDrive) {
	if vd == nil {
		return
	}
	if vb.listener == vd {
		vb.listener = nil
		vb.listenerActive = false
	}
	if vb.talker == vd {
		vb.talker = nil
		vb.talkerActive = false
	}
}

/*
func (c *Dispatcher) debugPeripheralWrite(data uint8) {
	value := data
	var message []string
	if value&0x02 != 0 {
		message = append(message, "[DATA_OUT]")
	}
	if value&0x08 != 0 {
		message = append(message, "[CLK_OUT]")
	}
	if value&0x10 != 0 {
		message = append(message, "[ATN_IN]")
	}
	fmt.Printf("DRV RECV: [%x] [%08b] %s\n", value, value, strings.Join(message, " "))
}

func (c *Dispatcher) debugPeripheralRead(data uint8) {
	value := data
	var message []string
	if value&0x01 != 0 {
		message = append(message, "[DATA_IN]")
	}
	if value&0x02 != 0 {
		message = append(message, "[DATA_OUT]")
	}
	if value&0x04 != 0 {
		message = append(message, "[CLK_IN]")
	}
	if value&0x08 != 0 {
		message = append(message, "[CLK_OUT]")
	}
	if value&0x10 != 0 {
		message = append(message, "[ATN_IN]")
	}
	if value&0x80 != 0 {
		message = append(message, "[ATN_OUT]")
	}
	fmt.Printf("DRV RECV: [%x] [%08b] %s\n", value, value, strings.Join(message, " "))
}

void debug_iec_drv_write(unsigned int data) {
    if (debug.iec) {
        uint8_t value = data;
        static uint8_t oldvalue = 0;

        if (value != oldvalue) {
            oldvalue = value;

            log_debug("$1800 store: %s %s %s",
                      value & 0x02 ? "DATA OUT" : "        ",
                      value & 0x08 ? "CLK OUT" : "       ",
                      value & 0x10 ? "ATNA   " : "       "
                      );
        }
    }
}

void debug_iec_drv_read(unsigned int data)
{
    if (debug.iec) {
        uint8_t value = data;
        static uint8_t oldvalue = { 0 };
        const char * data_correct = "";

        if (value != oldvalue) {
            unsigned int atn = value & 0x80 ? 1 : 0;
            unsigned int atna = value & 0x10 ? 1 : 0;
            unsigned int ddata = value & 0x01 ? 1 : 0;

            oldvalue = value;

            if (atn ^ atna) {
                if (!ddata) {
                    data_correct = " ***** ERROR: ATN, ATNA & DATA! *****";
                }
            }

            log_debug("$1800 read:  %s %s %s %s %s %s%s",
                      value & 0x02 ? "DATA OUT" : "        ",
                      value & 0x08 ? "CLK OUT" : "       ",
                      value & 0x10 ? "ATNA   " : "       ",

                      value & 0x01 ? "DATA IN" : "       ",
                      value & 0x04 ? "CLK IN" : "       ",
                      value & 0x80 ? "ATN" : "   ",
                      data_correct
                      );
        }
    }
}

	//_board->GetRam()[0x90] |= _board->GetBus()->Out(_board->GetRam()[0x95], _board->GetRam()[0xa3] & 0x80);
	//_board->GetRam()[0x90] |= _board->GetBus()->OutATN(_board->GetRam()[0x95]);
	//_board->GetRam()[0x90] |= _board->GetBus()->OutSec(_board->GetRam()[0x95]);
	//_board->GetRam()[0x90] |= _board->GetBus()->In(_a);
	//_board->GetBus()->SetATN();
	//_board->GetBus()->RelATN();
	//_board->GetBus()->Turnaround();
	//_board->GetBus()->Release();
*/
