package iec

import (
	"github.com/markel1974/c64emu/src/board/iboard"
	"github.com/markel1974/c64emu/src/board/iec/filedrive"
	"github.com/markel1974/c64emu/src/board/iec/virtualdrive"
	"github.com/markel1974/c64emu/src/preferences"
)

const (
	BusNum       = 32
	MaxDriveSize = 4
)

const (
	StOk          = 0    // No error
	StReadTimeout = 0x02 // Timeout on reading
	StTimeout     = 0x03 // Timeout
	StEof         = 0x40 // End of file
	StNotPresent  = 0x80 // Device not present
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

type IEC struct {
	cpuPort uint8
	cpuData uint8
	cpuBus  uint8
	drvPort uint8
	oldAtn  uint8
	drvBus  []uint8
	drvData []uint8

	peripheralStorage       []*C1541Model
	peripheralStorageActive []*C1541Model
	peripheralsCount        uint8
	peripheralsActiveCount  uint8
	virtualDrives           []virtualdrive.IVirtualDrive
	openData                []byte
	listener                virtualdrive.IVirtualDrive // Pointer to active listener
	talker                  virtualdrive.IVirtualDrive // Pointer to active talker
	listenerActive          bool                       // Listener selected, listener_data is valid
	talkerActive            bool                       // Talker selected, talker_data is valid
	listening               bool                       // Last ATN was listen (to decide between sec_listen/sec_talk)
	receivedCmd             uint8                      // Received command code ($x0)
	secAddr                 uint8                      // Received secondary address ($0x)
	emu1541                 bool
}

func NewIEC() *IEC {
	c := &IEC{
		//stubIdx:               0,
		//stub:                  []uint8{208, 144},
		drvBus:                  make([]uint8, BusNum),
		drvData:                 make([]uint8, BusNum),
		peripheralStorage:       make([]*C1541Model, BusNum),
		peripheralStorageActive: make([]*C1541Model, BusNum),
		virtualDrives:           make([]virtualdrive.IVirtualDrive, MaxDriveSize),
	}
	return c
}

func (c *IEC) AddPeripheral(peripheral *C1541Model) {
	if c.peripheralsCount >= BusNum {
		return
	}
	for i := uint8(0); i < c.peripheralsCount; i++ {
		if c.peripheralStorage[i] == peripheral {
			return
		}
	}
	c.peripheralStorage[c.peripheralsCount] = peripheral
	c.peripheralsCount++

	c.rebuildPeripherals()
	//TODO
	//peripheral->LedStateChangedEvent.Bind(new SignalExecutor2<IECBus, int, uint8>(this, &IECBus::ledStateChangedEventHandler));
}

func (c *IEC) RemovePeripheral(peripheral *C1541Model) {
	found := false
	for i := uint8(0); i < c.peripheralsCount; i++ {
		if c.peripheralStorage[i] == peripheral {
			c.peripheralsCount--
			c.peripheralStorage[i] = nil
			found = true
			break
		}
	}
	if found {
		for i := uint8(0); i < c.peripheralsCount; i++ {
			for j := i + 1; j < c.peripheralsCount; j++ {
				if c.peripheralStorage[i].GetDeviceNumber() < c.peripheralStorage[j].GetDeviceNumber() {
					tmp := c.peripheralStorage[i]
					c.peripheralStorage[i] = c.peripheralStorage[j]
					c.peripheralStorage[j] = tmp
				}
			}
		}
	}
	c.rebuildPeripherals()
}

func (c *IEC) Setup(board iboard.IBoard, prefs *preferences.Prefs) {
	for i := 0; i < MaxDriveSize; i++ {
		c.virtualDrives[i] = c.createVirtualDrive(prefs.Emul1541Proc(), 8+i, prefs.GetDrivePath(i))
	}
	for i := uint8(0); i < c.peripheralsCount; i++ {
		c.peripheralStorage[i].NewPrefs(prefs)
	}
	c.rebuildPeripherals()
	for i := uint8(0); i < MaxDriveSize; i++ {
		oldPath := ""
		if c.virtualDrives[i] != nil {
			oldPath = c.virtualDrives[i].GetPath()
		}
		newPath := prefs.GetDrivePath(int(i))
		if (oldPath != newPath) || c.emu1541 != prefs.Emul1541Proc() {
			c.destroyVirtualDrive(c.virtualDrives[i])
			c.virtualDrives[i] = c.createVirtualDrive(prefs.Emul1541Proc(), int(8+i), newPath)
		}
	}
	c.emu1541 = prefs.Emul1541Proc()
}

func (c *IEC) NewPrefs(prefs *preferences.Prefs) {
	//TODO IMPLEMENT
}

func (c *IEC) Reset() {
	c.listener = nil
	c.talker = nil
	c.listenerActive = false
	c.talkerActive = false
	c.listening = false
	c.receivedCmd = 0
	c.secAddr = 0
	c.openData = nil

	for i := 0; i < MaxDriveSize; i++ {
		if c.virtualDrives[i] != nil && c.virtualDrives[i].Ready() {
			c.virtualDrives[i].Reset()
		}
	}
}

func (c *IEC) CpuWrite(data uint8) {
	if c.emu1541 {
		c.updateCpuBus(^data)
		c.updateDrvBus()
		c.updatePorts()
		c.dispatchCpuWrite()
	}
}

func (c *IEC) CpuRead() uint8 {
	//TODO IMPLEMENT
	return StNotPresent
	//return c.cpuPort
}

func (c *IEC) PeripheralAtnResponse(data uint8, deviceNumber uint8) {
	c.PeripheralWrite(deviceNumber, data)
}

func (c *IEC) PeripheralRead(deviceNumber uint8) uint8 {
	return c.drvPort
}

func (c *IEC) PeripheralWrite(deviceNumber uint8, data uint8) {
	c.drvBus[deviceNumber] = ((data << 3) & 0x40) | ((data << 6) & (((^data) ^ c.cpuBus) << 3) & 0x80)
	c.drvData[deviceNumber] = data
	c.updatePorts()
}

func (c *IEC) Out(data uint8, eoi bool) uint8 {
	if c.listenerActive {
		if c.receivedCmd == CmdOpen {
			return c.openOut(data, eoi)
		}
		if c.receivedCmd == CmdData {
			return c.dataOut(data, eoi)
		}
		return StTimeout
	} else {
		return StTimeout
	}

}

func (c *IEC) OutATN(data uint8) uint8 {
	c.receivedCmd = 0 // Command is sent with secondary address
	c.secAddr = 0     // Command is sent with secondary address
	switch data & 0xf0 {
	case AtnListen:
		c.listening = true
		return c.listen(data & 0x0f)
	case AtnUnlisten:
		c.listening = false
		return c.unListen()
	case AtnTalk:
		c.listening = false
		return c.talk(data & 0x0f)
	case AtnUntalk:
		c.listening = false
		return c.unTalk()
	}
	return StTimeout
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
	return StTimeout
}

func (c *IEC) In(data *uint8) uint8 {
	if c.talkerActive && (c.receivedCmd == CmdData) {
		return c.dataIn(data)
	}
	*data = 0
	return StTimeout
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

func (c *IEC) createVirtualDrive(emul1541 bool, deviceNumber int, newPath string) virtualdrive.IVirtualDrive {
	if emul1541 {
		return nil
	}
	if len(newPath) == 0 {
		return nil
	}
	vd := filedrive.CreateDrive(uint8(deviceNumber), newPath)
	if vd != nil {
		//vd->LedStateChangedEvent.Bind(new SignalExecutor2<IECBus, int, uint8>(this, &IECBus::ledStateChangedEventHandler));
	}
	return vd
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

func (c *IEC) updateCpuBus(data uint8) {
	c.cpuBus = ((data << 2) & 0x80) | ((data << 2) & 0x40) | ((data << 1) & 0x10)
}

func (c *IEC) updateDrvBus() {
	for x := uint8(0); x < c.peripheralsActiveCount; x++ {
		p := c.peripheralStorageActive[x]
		//if (p->IsActive()) {
		unit := p.GetDeviceNumber()
		data := c.drvData[unit]
		c.drvBus[unit] = ((data << 3) & 0x40) | ((data << 6) & (((^data) ^ c.cpuBus) << 3) & 0x80)
		//drvBus[unit] = (((drvData[unit] << 3) & 0x40) | ((drvData[unit] << 6) & ((~drvData[unit] ^ cpuBus) << 3) & 0x80));
		//}
	}
}

func (c *IEC) updatePorts() {
	c.cpuPort = c.cpuBus
	for x := uint8(0); x < c.peripheralsActiveCount; x++ {
		p := c.peripheralStorageActive[x]
		//if (p->IsActive()) {
		unit := p.GetDeviceNumber()
		data := c.drvBus[unit]
		c.cpuPort &= data
		//}
	}
	c.drvPort = ((c.cpuPort >> 4) & 0x4) | (c.cpuPort >> 7) | ((c.cpuBus << 3) & 0x80)
}

func (c *IEC) dispatchCpuWrite() {
	newAtn := c.cpuBus & 0x10
	if c.oldAtn != newAtn {
		for x := uint8(0); x < c.peripheralsActiveCount; x++ {
			p := c.peripheralStorageActive[x]
			//if (p->IsActive()) {
			p.AtnStateChanged((c.oldAtn) != 0)
			//}
		}
		c.oldAtn = newAtn
	}
}

func (c *IEC) listen(device uint8) uint8 {
	c.listenerActive = false
	if device < 8 || device > 11 {
		return StNotPresent
	}
	if c.listener = c.virtualDrives[device-8]; c.listener == nil {
		return StNotPresent
	}
	if !c.listener.Ready() {
		return StNotPresent
	}
	c.listenerActive = true
	return StOk
}

func (c *IEC) talk(device uint8) uint8 {
	c.talkerActive = false
	if device < 8 || device > 11 {
		return StNotPresent
	}
	if c.talker = c.virtualDrives[device-8]; c.talker == nil {
		return StNotPresent
	}
	if !c.talker.Ready() {
		return StNotPresent
	}
	c.talkerActive = true
	return StOk
}

func (c *IEC) unListen() uint8 {
	c.listenerActive = false
	return StOk
}

func (c *IEC) unTalk() uint8 {
	c.talkerActive = false
	return StOk
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
	return StOk
}

func (c *IEC) secTalk() uint8 {
	return StOk
}

func (c *IEC) openOut(data uint8, eoi bool) uint8 {
	c.openData = append(c.openData, data)
	if eoi {
		if c.listener != nil {
			return c.listener.Open(c.secAddr, c.openData)
		}
	}
	return StOk
}

func (c *IEC) dataOut(data uint8, eoi bool) uint8 {
	if c.listener != nil {
		return c.listener.Write(c.secAddr, data, eoi)
	}
	return StOk
}

func (c *IEC) dataIn(data *uint8) uint8 {
	if c.talker != nil {
		return c.talker.Read(c.secAddr, data)
	}
	return 0
}

//void IECBus::ledStateChangedEventHandler(int  deviceNumber, uint8 state) {
//LedStateChangedEvent.Emit(deviceNumber, state);
//}

func (c *IEC) rebuildPeripherals() {
	c.peripheralsActiveCount = 0
	for driveId := uint8(0); driveId < c.peripheralsCount; driveId++ {
		if c1541 := c.peripheralStorage[driveId]; c1541 != nil && c1541.IsActive() {
			c.peripheralStorageActive[c.peripheralsActiveCount] = c1541
			c.peripheralsActiveCount++
		}
	}
}

func (c *IEC) ledStateChangedEventHandler(deviceNumber int, state uint8) {

}

//SignalBinder2<int, uint8> LedStateChangedEvent;
