package iec

import (
	"log"
	"strconv"
	"strings"

	"github.com/markel1974/c64emu/src/common/signals"
	"github.com/markel1974/c64emu/src/components/board"
	"github.com/markel1974/c64emu/src/components/iec/fsdrive"
	"github.com/markel1974/c64emu/src/components/iec/virtualdrive"
	"github.com/markel1974/c64emu/src/config"

	c1541board "github.com/markel1974/c64emu/src/c1541/board"
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


*/

type Dispatcher struct {
	*board.BaseComponent
	cfg             *config.Config
	atnState        uint8
	cpuPort         uint8
	cpuData         uint8
	cpuBus          uint8
	peripheralsPort uint8
	peripheralsData []uint8
	virtualDrives   []virtualdrive.IVirtualDrive
	ledSignal       *signals.Signal2[int, uint8]
}

func NewDispatcher(parent board.IComponent, suffix string) *Dispatcher {
	c := &Dispatcher{
		BaseComponent:   board.NewBaseComponent("iec", suffix),
		peripheralsData: make([]uint8, BusNum),
		virtualDrives:   nil,
		ledSignal:       signals.NewSignal2[int, uint8](),
	}
	board.Register(parent, c)
	return c
}

func (c *Dispatcher) AddPeripheral(peripheral *c1541board.Board) {
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

func (c *Dispatcher) RemovePeripheral(peripheral *c1541board.Board) {
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

func (c *Dispatcher) Setup(cfg *config.Config) {
	c.cfg = cfg
	c.cfg.Bind(c.configChanged)

	for idx, d := range cfg.GetDrives() {
		vd := c.createVirtualDrive(d.Kind, d.Opts, uint8(idx))
		c.virtualDrives = append(c.virtualDrives, vd)
	}

	//vd9 := c.createVirtualDrive(9, 1)
	//c.virtualDrives = append(c.virtualDrives, vd9)

	//c.rebuildPeripherals()
}

func (c *Dispatcher) configChanged() {
	//TODO IMPLEMENT
}

func (c *Dispatcher) Emulate() {
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

func (c *Dispatcher) Reset() {
	for _, vd := range c.virtualDrives {
		if vd.Ready() {
			vd.Reset()
		}
	}
}

func (c *Dispatcher) buildCpuBus(data uint8) uint8 {
	b6 := (data << 2) & 0x80
	b5 := (data << 2) & 0x40
	b4 := (data << 1) & 0x10
	value := b6 | b5 | b4
	return value
}

func (c *Dispatcher) buildPeripheralBus(cpuBus uint8, data uint8) uint8 {
	nData := ^data
	bBus := ((nData ^ cpuBus) << 3) & 0x80
	p1 := (data << 3) & 0x40
	p2 := (data << 6) & bBus
	value := p1 | p2
	return value
}

func (c *Dispatcher) updatePorts() {
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

func (c *Dispatcher) CpuWrite(data uint8) {
	c.cpuBus = c.buildCpuBus(^data)
	//c.debugCpuWrite(^c.cpuBus)
	c.updatePorts()
	c.notifyCpuWrite()
	//_board->GetRam()[0x90] |= _board->GetBus()->Out(_board->GetRam()[0x95], _board->GetRam()[0xa3] & 0x80);
	//_board->GetRam()[0x90] |= _board->GetBus()->OutATN(_board->GetRam()[0x95]);
	//_board->GetRam()[0x90] |= _board->GetBus()->OutSec(_board->GetRam()[0x95]);
	//_board->GetRam()[0x90] |= _board->GetBus()->In(_a);
	//_board->GetBus()->SetATN();
	//_board->GetBus()->RelATN();
	//_board->GetBus()->Turnaround();
	//_board->GetBus()->Release();
}

func (c *Dispatcher) CpuRead() uint8 {
	return c.cpuPort
}

func (c *Dispatcher) PeripheralRead() uint8 {
	return c.peripheralsPort
}

func (c *Dispatcher) PeripheralWrite(deviceNumber uint8, data uint8) {
	c.peripheralsData[deviceNumber] = data
	//c.debugPeripheralWrite(c.peripheralBus[deviceNumber])
	c.updatePorts()
}

func (c *Dispatcher) createVirtualDrive(kind string, opts string, deviceId uint8) virtualdrive.IVirtualDrive {
	deviceNumber := deviceId + 8
	var vd virtualdrive.IVirtualDrive
	switch kind {
	case "C1541":
		vd = c1541board.New(c, strconv.Itoa(int(deviceNumber)), c, deviceId, deviceNumber, opts)
	case "FSDRIVE":
		vd = fsdrive.New(c, deviceId, deviceNumber, opts)
	default:
		vd = c1541board.New(c, strconv.Itoa(int(deviceNumber)), c, deviceId, deviceNumber, opts)
	}
	vd.Setup(c.cfg)
	return vd
}

func (c *Dispatcher) notifyCpuWrite() {
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

func (c *Dispatcher) ledStateChangedEventHandler(deviceNumber int, state uint8) {
	c.ledSignal.Emit(deviceNumber, state)
}

func (c *Dispatcher) debugCpuWrite(data uint8) {
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
	log.Printf("CPU SEND: [%x] [%08b] %s\n", value, value, strings.Join(message, " "))
}

func (c *Dispatcher) debugCpuRead(data uint8) {
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
	log.Printf("CPU SEND: [%x] [%08b] %s\n", value, value, strings.Join(message, " "))
}

/*
// virtualbus.go
func (vb *VirtualBus) OutATN(data uint8) uint8 {
    //Invia il comando
    return vb.iec.CpuWrite(data)
}

func (vb *VirtualBus) OutSec(data uint8) uint8 {
     //Invia il comando
     return vb.iec.CpuWrite(data)
}

func (vb *VirtualBus) Out(data uint8, eoi bool) uint8 {
     //Invia il comando
     return vb.iec.CpuWrite(data)
}

func (vb *VirtualBus) In() (uint8, uint8) {
    return vb.iec.CpuRead()
}

func (vb *VirtualBus) SetATN() {
    // Non fare nulla (gestito internamente da Dispatcher e dal 1541 emulato)
}

func (vb *VirtualBus) RelATN() {
    // Non fare nulla (gestito internamente da Dispatcher e dal 1541 emulato)
}

func (vb *VirtualBus) Turnaround() {
    // Non fare nulla (gestito internamente da Dispatcher e dal 1541 emulato)
}

func (vb *VirtualBus) Release() {
    // Non fare nulla (gestito internamente da Dispatcher e dal 1541 emulato)
}
*/
