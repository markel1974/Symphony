package board

import (
	"github.com/markel1974/c64emu/src/board/banks"
	"github.com/markel1974/c64emu/src/board/cartridges"
	"github.com/markel1974/c64emu/src/board/iec"
	"github.com/markel1974/c64emu/src/board/keyboard"
	"github.com/markel1974/c64emu/src/board/prg"
	cpu2 "github.com/markel1974/c64emu/src/components/6510"
	cia3 "github.com/markel1974/c64emu/src/components/cia"
	quartz2 "github.com/markel1974/c64emu/src/components/quartz"
	"github.com/markel1974/c64emu/src/components/sid"
	vic2 "github.com/markel1974/c64emu/src/components/vic"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/signals"
	"golang.design/x/clipboard"
	"log"
	"os"
)

//type PhiMode int

const intrExpansionId = 0x10

const (
// PhiIdle = PhiMode(0)
// Phi1    = PhiMode(1)
// Phi2    = PhiMode(2)
)

type Board struct {
	db           vic2.IDisplayBuffer
	quartz       *quartz2.Quartz
	cpu          *cpu2.MOS6510
	vic          *vic2.MOS6569
	sid          *sid.MOS6581
	cia1         *cia3.MOS6526A
	cia2         *cia3.MOS6526B
	pic          *cpu2.Pic
	iec          *iec.IEC
	keys         *keyboard.Keyboard
	cfg          *config.Config
	hasClipboard bool
	cartMan      *cartridges.Manager
	banks        *banks.Banks
	dma          int
	dmaLow       bool
	//baLow        bool
	//aecLow       bool

	irqTrigger *signals.SignalUint32
	irqClear   *signals.SignalUint32
	//phiMode    PhiMode
}

func NewBoard(db vic2.IDisplayBuffer) *Board {
	b := &Board{
		db:           db,
		quartz:       quartz2.NewQuartz(),
		iec:          nil,
		cpu:          nil,
		vic:          nil,
		sid:          nil,
		cia1:         nil,
		cia2:         nil,
		pic:          nil,
		keys:         nil,
		hasClipboard: false,
		cartMan:      cartridges.NewManager(),
		dma:          0,
		dmaLow:       false,
		//baLow:        false,
		//aecLow:       false,
		banks:      nil,
		irqTrigger: nil,
		irqClear:   nil,
		//phiMode:    PhiIdle,
	}
	return b
}

func (s *Board) Setup(cfg *config.Config) error {
	if err := clipboard.Init(); err != nil {
		log.Printf("can't init clipboard: %s", err)
	} else {
		s.hasClipboard = true
	}
	s.cfg = cfg
	s.cfg.Bind(s.configChanged)

	s.pic = cpu2.NewPic()
	s.iec = iec.NewIEC()
	s.cpu = cpu2.NewMOS6510("c64")
	s.vic = vic2.NewMOS6569(s.db)
	s.sid = sid.NewMOS6581()
	s.cia1 = cia3.NewMOS6526A()
	s.cia2 = cia3.NewMOS6526B()
	s.keys = keyboard.NewKeyboard()
	s.banks = banks.NewBanks()

	s.pic.Setup(s.quartz)
	s.iec.Setup(cfg)
	s.cpu.Setup(s.pic, s.banks)

	s.vic.Setup(s.quartz, s.banks, cfg)
	s.vic.SignalReadyBind(s.ReadySlot)
	s.vic.SignalTriggerIRQBind(s.irqTriggerSlot)
	s.vic.SignalClearIRQBind(s.irqClearSlot)
	s.vic.SignalBALowBind(s.cpu.SetRDYLow)
	s.vic.SignalAECLowBind(s.cpu.SetAECLow)

	s.sid.Setup(cfg)

	s.cia1.Setup(cfg)
	s.cia1.SignalTriggerIRQBind(s.irqTriggerSlot)
	s.cia1.SignalClearIRQBind(s.irqClearSlot)
	s.cia1.SignalLightPenTriggerBind(s.vic.LightPenTrigger)

	s.cia2.Setup(s.iec, cfg)
	s.cia2.SignalTriggerNMIBind(s.pic.TriggerNMI)
	s.cia2.SignalClearNMIBind(s.pic.ClearNMI)
	s.cia2.SignalChangedVABind(s.vic.ChangedVA)

	s.cartMan.Setup(s, cfg)
	s.banks.Setup(s.vic, s.sid, s.cia1, s.cia2, s.cartMan, cfg)

	//TODO TUTTE LE CONNESSIONI CON PIN DEVONO ESSERE EFFETTUATE TRAMITE SIGNAL - SLOT
	//test := signals.NewSignalByte()
	//test.Bind(s.sid.SetPotXSlot)

	if !s.cfg.GetDisableCartridgeAutostart() {
		for _, cartName := range s.cfg.GetCartridges() {
			var data []uint8
			var err error
			if len(cartName.Path) > 0 {
				data, err = os.ReadFile(cartName.Path)
				if err != nil {
					log.Printf("can't add cartridge: %s", err.Error())
					continue
				}
			}
			if cartId, err := s.cartMan.Add(cartName.Kind, cartName.Path, data); err != nil {
				log.Printf("can't add cartridge: %s", err.Error())
			} else {
				log.Printf("cartridge: %s [%s] successfully added", cartName, cartId)
			}
		}
	}

	s.Reset()

	return nil
}

func (s *Board) Reset() {
	s.pic.Reset()
	s.banks.Reset()
	s.cpu.Reset()
	s.cartMan.Reset()
	s.sid.Reset()
	s.cia1.Reset()
	s.cia2.Reset()
	s.iec.Reset()

	s.dma = 0
	s.dmaLow = false
	//s.baLow = false
	//s.aecLow = false
}

func (s *Board) AsyncReset() {
	s.keys.Reset()
	s.banks.Reset()
	s.pic.TriggerReset()
	//s.cpu.AsyncReset()
	s.vic.Reset()
	s.sid.Reset()
	s.cia1.Reset()
	s.cia2.Reset()
	s.iec.Reset()

	s.dma = 0
	s.dmaLow = false
	//s.baLow = false
	//s.aecLow = false
}

func (s *Board) configChanged() {
}

//var testCounter = 0

func (s *Board) Emulate() bool {
	//s.phiMode = Phi1
	vBlank, lastVicCycle := s.vic.Emulate()
	//s.phiMode = Phi2
	if vBlank {
		//sidCounter := s.sid.ResetHistoryCounter()
		//TODO
		_ = s.sid.ResetHistoryCounter()
		s.cia1.CountTOD()
		s.cia2.CountTOD()
		s.updateKeyboard()
	}
	if lastVicCycle {
		s.sid.Emulate()
	}
	s.cia1.CheckIRQs()
	s.cia2.CheckIRQs()
	s.cia1.Emulate()
	s.cia2.Emulate()

	//TODO SLOW CPU BUG (1/32)
	//if s.speed > 0 {
	//	if testCounter == s.speed {
	//		s.cpu.Emulate()
	//		testCounter = 0
	//	} else {
	//		testCounter++
	//	}
	//}

	s.cpu.Emulate()
	s.cartMan.Emulate()
	s.iec.Emulate()
	s.quartz.AddCycle()
	//s.phiMode = PhiIdle
	return vBlank
}

func (s *Board) ReadySlot() {
	s.keys.SetReady()
}

// LedStateChangedEvent (deviceNumber int, state uint8)
func (s *Board) LedStateChangedEvent(_ int, _ uint8) {
	//TODO IMPLEMENT
	//deviceId := deviceNumber - 8
	//if deviceId < 0 || deviceId >= MAX_DRIVE_COUNT {
	//	return
	//}
	//k.leds[deviceId] = state
	//k.updateLedState()
	//s.keys.InputReady(_ledActivities == 0)
}

func (s *Board) KeyboardPaste(pressed bool) {
	if !pressed {
		return
	}
	if !s.hasClipboard {
		return
	}
	data := clipboard.Read(clipboard.FmtText)
	s.keys.SetCommand(string(data), "")
}

func (s *Board) KeyboardSetExt(pressed bool) {
	if !pressed {
		return
	}
	s.keys.SetExt()
}

func (s *Board) KeyboardSetNumLock(pressed bool) {
	if !pressed {
		return
	}
	s.keys.SetNumLock()
}

func (s *Board) KeyboardSetCapital(pressed bool) {
	if !pressed {
		return
	}
	//TODO RIMUOVERE
	//s.loadPRG(s.cfg.GetPrg())
	//return

	//TODO RIMUOVERE
	//s.speed++
	//fmt.Println("SPEED", s.speed)
	//return

	s.keys.SetCapital()
}

func (s *Board) KeyboardSetMenu(pressed bool) {
	if !pressed {
		return
	}
	s.keys.SetMenu()
}

func (s *Board) KeyboardSetVirtualKey(pressed bool, vKey int) {
	s.keys.SetVirtualKey(pressed, vKey)
}

func (s *Board) KeyboardSwapJoystick(pressed bool) {
	if !pressed {
		return
	}
	s.keys.SwapJoystick()
}

func (s *Board) updateKeyboard() {
	if c64Byte, c64Bit, pressed, joyKey, shifted, ok := s.keys.PollKeyboard(); ok {
		joyKey1 := joyKey
		joyKey2 := uint8(0xff)
		if s.keys.HasJoystickSwap() {
			joyKey2 = joyKey
			joyKey1 = uint8(0xff)
		}
		if pressed {
			s.cia1.SetKeyDown(c64Byte, c64Bit, shifted, joyKey1, joyKey2)
		} else {
			s.cia1.SetKeyUp(c64Byte, c64Bit, shifted, joyKey1, joyKey2)
		}
	}
}

func (s *Board) GameExRomConfigChanged() {
	s.banks.RebuildMemoryConfig()
}

func (s *Board) Read(addr uint16) uint8 {
	return s.banks.Read(addr)
}

func (s *Board) Write(addr uint16, data uint8) {
	s.banks.Write(addr, data)
}

func (s *Board) NMITrigger() {
	s.pic.TriggerNMI()
}

func (s *Board) SetDMALow(v bool) {
	//TODO IMPLEMENT
	//if _DMA=Low the CPU can be requested to release the bus.
	//It will stop after the next read cycle and all bus lines will go to high resistance state.
	//So other units can use the computer hardware. At _DMA=High the CPU continues to work.
	if v {
		s.dma++
		s.dmaLow = true
	} else {
		if s.dma > 0 {
			s.dma--
			if s.dma == 0 {
				s.dmaLow = false
			}
		}
	}
	//s.updateCpuRdy()
}

func (s *Board) ResetTrigger() {
	s.pic.TriggerReset()
}

func (s *Board) IRQTrigger() {
	s.pic.TriggerIRQ(intrExpansionId)
}

func (s *Board) IRQClear() {
	s.pic.ClearIRQ(intrExpansionId)
}

func (s *Board) IRQTriggerBind(fn func(uint32)) {
	if s.irqTrigger == nil {
		s.irqTrigger = signals.NewSignalUint32()
	}
	s.irqTrigger.Bind(fn)
}

func (s *Board) IRQClearBind(fn func(uint32)) {
	if s.irqClear == nil {
		s.irqClear = signals.NewSignalUint32()
	}
	s.irqClear.Bind(fn)
}

func (s *Board) BusAvailable() bool {
	return !s.vic.GetBALow()
}

func (s *Board) AECAvailable() bool {
	return !s.vic.GetAECLow()
}

func (s *Board) Cycle() uint64 {
	return s.quartz.Cycle()
}

func (s *Board) CycleAlarm(id string, callback quartz2.AlarmCallback) *quartz2.Alarm {
	return s.quartz.NewAlarm(id, callback)
}

func (s *Board) RamSetWriteTrigger(addr uint16, fn func(uint16, uint8)) int {
	return s.banks.SetWriteTrigger(addr, fn)
}

func (s *Board) RamRemoveWriteTrigger(addr uint16, id int) {
	s.banks.RemoveRamTrigger(addr, id)
}

func (s *Board) RmwFlags() uint8 {
	//TODO IMPLEMENT cpu rmw flags
	return 0
}

func (s *Board) ExtRamWrite(memConfig int, addr uint16, data uint8) {
	var prev []uint8 = nil
	if memConfig >= 0 {
		prev = s.banks.GetMemoryConfig()
		s.banks.SetMemoryEntry(uint8(memConfig))
	}
	s.banks.Write(addr, data)
	if prev != nil {
		s.banks.SetMemoryConfig(prev)
	}
}

func (s *Board) ExtRamRead(memConfig int, addr uint16) uint8 {
	var prev []uint8 = nil
	if memConfig >= 0 {
		prev = s.banks.GetMemoryConfig()
		s.banks.SetMemoryEntry(uint8(memConfig))
	}
	rb := s.banks.Read(addr)
	if prev != nil {
		s.banks.SetMemoryConfig(prev)
	}
	return rb
}

func (s *Board) irqTriggerSlot(i uint32) {
	s.pic.TriggerIRQ(i)
	if s.irqTrigger != nil {
		s.irqTrigger.Emit(i)
	}
}

func (s *Board) irqClearSlot(i uint32) {
	s.pic.ClearIRQ(i)
	if s.irqClear != nil {
		s.irqClear.Emit(i)
	}
}

/*
func (s *Board) baLowSlot(baLow bool) {
	s.baLow = baLow
	//s.updateCpuRdy()
}

func (s *Board) aecLowSlot(aecLow bool) {
	s.aecLow = aecLow
	s.updateCpuRdy()
}

func (s *Board) updateCpuRdy() {
	//s.cpu.SetRDYLow((s.baLow) || s.dmaLow)
	//s.cpu.SetRDYLow((s.baLow && s.aecLow) || s.dmaLow)
	s.cpu.SetRDYLow(s.aecLow || s.dmaLow)
}
*/

func (s *Board) loadPRG(prgFile string) {
	//TODO TEST - WE HAVE TO WAIT READY
	//s.loadPRG(s.cfg.GetPrg())
	//return
	if len(prgFile) == 0 {
		return
	}
	p := prg.NewPRG()
	if err := p.Load(prgFile); err != nil {
		log.Printf("can't set load prg: %s", err.Error())
		return
	}
	inject := banks.NewObserver(s.banks)
	if err := inject.Inject(false, p.GetStartAddress(), p.GetData()); err != nil {
		log.Printf("can't set prg: %s", err.Error())
		return
	}
}
