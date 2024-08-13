package iec

import (
	"fmt"
	"github.com/markel1974/c64emu/src/board/iec/drives/c1541"
	"github.com/markel1974/c64emu/src/board/iec/drives/fsdrive"
	"github.com/markel1974/c64emu/src/board/iec/virtualdrive"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/signals"
	"strings"
)

/*
Diamo un'occhiata alla sequenza in cui un carattere che sta per essere trasmesso. In questo momento, sia la clock line
che lq data line vengono mantenute down [true]. Il talker gestisce la clock line e il listener la data line. Potrebbe esserci
più di listener, nel qual caso tutti gestiscono la data line
*/
const (
	BusNum       = 32
	MaxDriveSize = 4
)

/*
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
*/

type IEC struct {
	cfg             *config.Config
	atnState        uint8
	cpuPort         uint8
	cpuData         uint8
	cpuBus          uint8
	peripheralsPort uint8
	peripheralsData []uint8
	virtualDrives   []virtualdrive.IVirtualDrive
	ledSignal       *signals.Signal2[int, uint8]

	//openData                []byte
	//listener                virtualdrive.IVirtualDrive // Pointer to active listener
	//talker                  virtualdrive.IVirtualDrive // Pointer to active talker
	//listenerActive          bool                       // Listener selected, listener_data is valid
	//talkerActive            bool                       // Talker selected, talker_data is valid
	//listening               bool                       // Last ATN was listen (to decide between sec_listen/sec_talk)
	//receivedCmd             uint8                      // Received command code ($x0)
	//secAddr                 uint8                      // Received secondary address ($0x)
	//emu1541 bool
}

func NewIEC() *IEC {
	c := &IEC{
		peripheralsData: make([]uint8, BusNum),
		virtualDrives:   nil,
		ledSignal:       signals.NewSignal2[int, uint8](),
	}
	return c
}

func (c *IEC) AddPeripheral(peripheral *c1541.Board) {
	//if c.peripheralsCount >= BusNum {
	//	return
	//}
	//for i := uint8(0); i < c.peripheralsCount; i++ {
	//	if c.peripheralStorage[i] == peripheral {
	//		return
	//	}
	//}
	//c.peripheralStorage[c.peripheralsCount] = peripheral
	//c.peripheralsCount++
	//c.rebuildPeripherals()
	//TODO
	//peripheral->LedStateChangedEvent.Bind(new SignalExecutor2<IECBus, int, uint8>(this, &IECBus::ledStateChangedEventHandler));
}

func (c *IEC) RemovePeripheral(peripheral *c1541.Board) {
	//found := false
	//for i := uint8(0); i < c.peripheralsCount; i++ {
	//	if c.peripheralStorage[i] == peripheral {
	//		c.peripheralsCount--
	//		c.peripheralStorage[i] = nil
	//		found = true
	//		break
	//	}
	//}
	//if found {
	//	for i := uint8(0); i < c.peripheralsCount; i++ {
	//		for j := i + 1; j < c.peripheralsCount; j++ {
	//			if c.peripheralStorage[i].GetDeviceNumber() < c.peripheralStorage[j].GetDeviceNumber() {
	//				tmp := c.peripheralStorage[i]
	//				c.peripheralStorage[i] = c.peripheralStorage[j]
	//				c.peripheralStorage[j] = tmp
	//			}
	//		}
	//	}
	//}
	//c.rebuildPeripherals()
}

func (c *IEC) Setup(cfg *config.Config) {
	c.cfg = cfg
	c.cfg.Bind(c.configChanged)

	//TODO ATTIVARE PER TEST
	vd8 := c.createVirtualDrive(8, 1)
	c.virtualDrives = append(c.virtualDrives, vd8)
	//vd9 := c.createVirtualDrive(9, 1)
	//c.virtualDrives = append(c.virtualDrives, vd9)

	//c.rebuildPeripherals()
}

func (c *IEC) configChanged() {
	//TODO IMPLEMENT
}

func (c *IEC) Emulate() {
	if len(c.virtualDrives) == 0 {
		return
	}
	if len(c.virtualDrives) == 1 {
		c.virtualDrives[0].Emulate()
		return
	}
	for _, vd := range c.virtualDrives {
		vd.Emulate()
	}
}

func (c *IEC) Reset() {
	//c.listener = nil
	//c.talker = nil
	//c.listenerActive = false
	//c.talkerActive = false
	//c.listening = false
	//c.receivedCmd = 0
	//c.secAddr = 0
	//c.openData = nil
	for _, vd := range c.virtualDrives {
		if vd.Ready() {
			vd.Reset()
		}
	}
}

func (c *IEC) buildCpuBus(data uint8) uint8 {
	b6 := (data << 2) & 0x80
	b5 := (data << 2) & 0x40
	b4 := (data << 1) & 0x10
	value := b6 | b5 | b4
	return value
}

func (c *IEC) buildPeripheralBus(cpuBus uint8, data uint8) uint8 {
	nData := ^data
	bBus := ((nData ^ cpuBus) << 3) & 0x80
	p1 := (data << 3) & 0x40
	p2 := (data << 6) & bBus
	value := p1 | p2
	return value
}

func (c *IEC) updatePorts() {
	c.cpuPort = c.cpuBus
	for _, vd := range c.virtualDrives {
		unit := vd.GetDeviceNumber()
		pData := c.peripheralsData[unit]
		pBus := c.buildPeripheralBus(c.cpuBus, pData)
		c.cpuPort &= pBus
	}
	bp7 := (c.cpuPort >> 4) & 0x04
	bp8 := c.cpuPort >> 7
	bb5 := (c.cpuBus << 3) & 0x80
	value := bp7 | bp8 | bb5
	c.peripheralsPort = value
}

func (c *IEC) CpuWrite(data uint8) {
	c.cpuBus = c.buildCpuBus(^data)
	//c.debugCpuWrite(^c.cpuBus)
	c.updatePorts()
	c.notifyCpuWrite()
	/*
		_board->GetRam()[0x90] |= _board->GetBus()->Out(_board->GetRam()[0x95], _board->GetRam()[0xa3] & 0x80);
		_board->GetRam()[0x90] |= _board->GetBus()->OutATN(_board->GetRam()[0x95]);
		_board->GetRam()[0x90] |= _board->GetBus()->OutSec(_board->GetRam()[0x95]);
		_board->GetRam()[0x90] |= _board->GetBus()->In(_a);
		_board->GetBus()->SetATN();
		_board->GetBus()->RelATN();
		_board->GetBus()->Turnaround();
		_board->GetBus()->Release();
	*/
}

func (c *IEC) CpuRead() uint8 {
	return c.cpuPort
}

func (c *IEC) PeripheralRead() uint8 {
	return c.peripheralsPort
}

func (c *IEC) PeripheralWrite(deviceNumber uint8, data uint8) {
	c.peripheralsData[deviceNumber] = data
	//c.debugPeripheralWrite(c.peripheralBus[deviceNumber])
	c.updatePorts()
}

func (c *IEC) createVirtualDrive(deviceNumber uint8, kind int) virtualdrive.IVirtualDrive {
	switch kind {
	case 1:
		vd := c1541.New(c, deviceNumber)
		vd.Setup(c.cfg)
		return vd
	case 2:
		vd := fsdrive.New(c, deviceNumber)
		vd.Setup(c.cfg)
		//vd->LedStateChangedEvent.Bind(new SignalExecutor2<IECBus, int, uint8>(this, &IECBus::ledStateChangedEventHandler));
		return vd
	}
	return nil
}

func (c *IEC) notifyCpuWrite() {
	newAtnState := c.cpuBus & 0x10
	if c.atnState == newAtnState {
		for _, vd := range c.virtualDrives {
			vd.BusStateChanged(c.peripheralsPort)
		}
		return
	}
	for _, vd := range c.virtualDrives {
		vd.AtnStateChanged(c.atnState != 0, newAtnState != 0)
	}
	c.atnState = newAtnState
}

//void IECBus::ledStateChangedEventHandler(int  deviceNumber, uint8 state) {
//LedStateChangedEvent.Emit(deviceNumber, state);
//}

func (c *IEC) ledStateChangedEventHandler(deviceNumber int, state uint8) {
	c.ledSignal.Emit(deviceNumber, state)
}

func (c *IEC) debugCpuWrite(data uint8) {
	//value := ^data
	value := data
	var message []string
	if value&0x20 != 0 {
		message = append(message, "[DATA_OUT]")
	}
	if value&0x10 != 0 {
		message = append(message, "[CLK_OUT]")
	}
	if value&0x08 != 0 {
		message = append(message, "[ATN_OUT]")
	}
	fmt.Printf("CPU SEND: [%x] [%08b] %s\n", value, value, strings.Join(message, " "))
}

func (c *IEC) debugCpuRead(data uint8) {
	value := data
	var message []string
	if value&0x80 != 0 {
		message = append(message, "[CLK_IN]")
	}
	if value&0x40 != 0 {
		message = append(message, "[DATA_IN]")
	}
	if value&0x20 != 0 {
		message = append(message, "[DATA_OUT]")
	}
	if value&0x10 != 0 {
		message = append(message, "[CLK_OUT]")
	}
	if value&0x08 != 0 {
		message = append(message, "[ATN_OUT]")
	}
	fmt.Printf("CPU SEND: [%x] [%08b] %s\n", value, value, strings.Join(message, " "))
}

/*
func (c *IEC) debugPeripheralWrite(data uint8) {
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

func (c *IEC) debugPeripheralRead(data uint8) {
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
*/
/*
void debug_iec_drv_write(unsigned int data)
{
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


*/

/*
func (c *IEC) Out(data uint8, eoi bool) uint8 {
	if c.listenerActive {
		if c.receivedCmd == CmdOpen {
			return c.openOut(data, eoi)
		}
		if c.receivedCmd == CmdData {
			return c.dataOut(data, eoi)
		}
		return virtualdrive.StTimeout
	} else {
		return virtualdrive.StTimeout
	}
}

func (c *IEC) OutATN(data uint8) uint8 {
	c.receivedCmd = 0 // Command is sent with secondary address
	c.secAddr = 0     // Command is sent with secondary address
	d := data & 0xf0
	switch d {
	case AtnListen:
		c.listening = true
		return c.listen(d)
	case AtnUnlisten:
		c.listening = false
		return c.unListen()
	case AtnTalk:
		c.listening = false
		return c.talk(d)
	case AtnUntalk:
		c.listening = false
		return c.unTalk()
	}
	return virtualdrive.StTimeout
}

func (c *IEC) OutSec(data uint8) uint8 {
	if c.listening {
		if c.listenerActive {
			c.secAddr = data & 0x0f
			c.receivedCmd = data & 0xf0
			return c.secListen()
		}
	} else {
		if c.talkerActive {
			c.secAddr = data & 0x0f
			c.receivedCmd = data & 0xf0
			return c.secTalk()
		}
	}
	return virtualdrive.StTimeout
}

func (c *IEC) In() (uint8, uint8) {
	if c.talkerActive && (c.receivedCmd == CmdData) {
		return c.dataIn()
	}
	return virtualdrive.StTimeout, 0
}

func (c *IEC) SetATN() {
	// Only needed for real IEC
}

func (c *IEC) RelATN() {
	// Only needed for real IEC
}

func (c *IEC) Turnaround() {
	// Only needed for real IEC
}

func (c *IEC) Release() {
	// Only needed for real IEC
}

func (c *IEC) listen(device uint8) uint8 {
	c.listenerActive = false
	if device < 8 || device > 11 {
		return virtualdrive.StNotPresent
	}
	if c.listener = c.virtualDrives[device-8]; c.listener == nil {
		return virtualdrive.StNotPresent
	}
	if !c.listener.Ready() {
		return virtualdrive.StNotPresent
	}
	c.listenerActive = true
	return virtualdrive.StOk
}

func (c *IEC) talk(device uint8) uint8 {
	c.talkerActive = false
	if device < 8 || device > 11 {
		return virtualdrive.StNotPresent
	}
	if c.talker = c.virtualDrives[device-8]; c.talker == nil {
		return virtualdrive.StNotPresent
	}
	if !c.talker.Ready() {
		return virtualdrive.StNotPresent
	}
	c.talkerActive = true
	return virtualdrive.StOk
}

func (c *IEC) unListen() uint8 {
	c.listenerActive = false
	return virtualdrive.StOk
}

func (c *IEC) unTalk() uint8 {
	c.talkerActive = false
	return virtualdrive.StOk
}

func (c *IEC) secListen() uint8 {
	switch c.receivedCmd {
	case CmdOpen:
		c.openData = nil
	case CmdClose:
		if c.listener != nil {
			return c.listener.Close(c.secAddr)
		}
	}
	return virtualdrive.StOk
}

func (c *IEC) secTalk() uint8 {
	return virtualdrive.StOk
}

func (c *IEC) openOut(data uint8, eoi bool) uint8 {
	c.openData = append(c.openData, data)
	if eoi {
		if c.listener != nil {
			return c.listener.Open(c.secAddr, c.openData)
		}
	}
	return virtualdrive.StOk
}

func (c *IEC) dataOut(data uint8, eoi bool) uint8 {
	if c.listener != nil {
		return c.listener.Write(c.secAddr, data, eoi)
	}
	return virtualdrive.StOk
}

func (c *IEC) dataIn() (uint8, uint8) {
	if c.talker != nil {
		return c.talker.Read(c.secAddr)
	}
	return 0, 0
}

func (c *IEC) destroyVirtualDrive(vd virtualdrive.IVirtualDrive) {
	if vd == nil {
		return
	}
	if c.listener == vd {
		c.listener = nil
		c.listenerActive = false
	}
	if c.talker == vd {
		c.talker = nil
		c.talkerActive = false
	}
}
*/
