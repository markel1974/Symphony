package board

import (
	"github.com/markel1974/c64emu/src/board/banks"
	"github.com/markel1974/c64emu/src/board/cartridges"
	"github.com/markel1974/c64emu/src/board/cia"
	"github.com/markel1974/c64emu/src/board/cpu"
	"github.com/markel1974/c64emu/src/board/iec"
	"github.com/markel1974/c64emu/src/board/keyboard"
	"github.com/markel1974/c64emu/src/board/quartz"
	"github.com/markel1974/c64emu/src/board/sid"
	"github.com/markel1974/c64emu/src/board/vic"
	"github.com/markel1974/c64emu/src/preferences"
	"golang.design/x/clipboard"
	"log"
	"os"
)

type PhiMode int

const (
	PhiIdle = PhiMode(0)
	Phi1    = PhiMode(1)
	Phi2    = PhiMode(2)
)

type Board struct {
	db           vic.IDisplayBuffer
	quartz       *quartz.Quartz
	cpu          *cpu.MOS6510
	vic          *vic.MOS6569
	sid          *sid.MOS6581
	cia1         *cia.MOS6526_1
	cia2         *cia.MOS6526_2
	interrupts   *cpu.Interrupts
	iec          *iec.IEC
	keys         *keyboard.Keyboard
	prefs        *preferences.Prefs
	hasClipboard bool
	phiMode      PhiMode
	cartMan      *cartridges.Manager
	banks        *banks.Banks
}

func NewBoard(db vic.IDisplayBuffer) *Board {
	b := &Board{
		db:           db,
		quartz:       quartz.NewQuartz(),
		iec:          nil,
		cpu:          nil,
		vic:          nil,
		sid:          nil,
		cia1:         nil,
		cia2:         nil,
		interrupts:   nil,
		keys:         nil,
		hasClipboard: false,
		cartMan:      cartridges.NewManager(),
		phiMode:      PhiIdle,
		banks:        nil,
	}
	return b
}

func (s *Board) Setup(prefs *preferences.Prefs) error {
	if err := clipboard.Init(); err != nil {
		log.Printf("can't init clipboard: %s", err)
	} else {
		s.hasClipboard = true
	}

	s.prefs = prefs

	s.iec = iec.NewIEC()
	s.cpu = cpu.NewMOS6510()
	s.vic = vic.NewMOS6569(s.db)
	s.sid = sid.NewMOS6581()
	s.cia1 = cia.NewMOS6526_1()
	s.cia2 = cia.NewMOS6526_2()
	s.keys = keyboard.NewKeyboard()
	s.banks = banks.NewBanks()

	s.iec.Setup(s.quartz, prefs)
	s.cpu.Setup(s.quartz, s.banks, prefs)
	s.interrupts = s.cpu.GetInterrupts()
	s.vic.Setup(s.quartz, s.interrupts, s.banks, prefs)
	s.vic.GetReadySignal().Bind(s.ReadyEvent)
	//s.vic.GetBALowSignal().Bind(baLow)
	s.sid.Setup(prefs)
	s.cia1.Setup(s.interrupts, s.vic, prefs)
	s.cia2.Setup(s.interrupts, s.vic, s.iec, prefs)
	s.cartMan.Setup(s, prefs)
	s.banks.Setup(s.vic, s.sid, s.cia1, s.cia2, s.cartMan, prefs)

	if !s.prefs.GetDisableCartridgeAutostart() {
		for _, cartName := range s.prefs.GetCartridges() {
			data, err := os.ReadFile(cartName)
			if err != nil {
				log.Printf("can't add cartridge: %s", err.Error())
				continue
			}
			if cartId, err := s.cartMan.Add(cartName, data); err != nil {
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
	s.banks.Reset()
	s.cpu.Reset()
	s.sid.Reset()
	s.cia1.Reset()
	s.cia2.Reset()
	s.iec.Reset()
}

func (s *Board) AsyncReset() {
	s.keys.Reset()
	s.banks.AsyncReset()
	s.cpu.AsyncReset()
	s.vic.Reset()
	s.sid.Reset()
	s.cia1.Reset()
	s.cia2.Reset()
	s.iec.Reset()
}

func (s *Board) NewPrefs(prefs *preferences.Prefs) {
	s.prefs = prefs
	s.iec.NewPrefs(prefs)
	s.sid.NewPrefs(prefs)
	s.vic.NewPrefs(prefs)
}

func (s *Board) Emulate() bool {
	s.phiMode = Phi1
	vBlank, lastVicCycle := s.vic.Emulate()
	s.phiMode = Phi2
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
	s.cpu.Emulate(s.vic.GetBALow())
	s.iec.Emulate()
	s.cartMan.Emulate()

	s.quartz.AddCycle()
	s.phiMode = PhiIdle
	return vBlank
}

func (s *Board) GameExRomConfigChanged() {
	s.banks.RebuildMemoryConfig()
}

func (s *Board) CreateAlarm(name string, callback quartz.AlarmCallback) *quartz.Alarm {
	return s.quartz.NewAlarm(name, callback)
}

func (s *Board) ReadyEvent() {
	s.keys.SetReady()
}

func (s *Board) LedStateChangedEvent(deviceNumber int, state uint8) {
	//TODO IMPLEMENT
	//deviceId := deviceNumber - 8
	//if deviceId < 0 || deviceId >= MAX_DRIVE_COUNT {
	//	return
	//}
	//k.leds[deviceId] = state
	//k.updateLedState()
	//s.keys.InputReady(_ledActivities == 0)
}

func (s *Board) RmwFlags() uint8 {
	//TODO IMPLEMENT cpu rmw flags
	return 0
}

func (s *Board) Cycle() uint64 {
	return s.quartz.Cycle()
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
	s.keys.SetCapital()
}

func (s *Board) KeyboardSetMenu(menu bool) {
	s.keys.SetMenu(menu)
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

func (s *Board) RamRead(addr uint16) uint8 {
	return s.banks.Read(addr)
}

func (s *Board) RamWrite(addr uint16, data uint8) {
	s.banks.Write(addr, data)
}

func (s *Board) NMI() {
	s.interrupts.TriggerNMI()
}

func (s *Board) DMA(l bool) {
	//TODO IMPLEMENT
	//if _DMA=Low the CPU can be requested to release the bus.
	//It will stop after the next read cycle and all bus lines will go to high resistance state.
	//So other units can use the computer hardware. At _DMA=High the CPU continues to work.
}

func (s *Board) IRQIn() {
	//TODO IMPLEMENT
	//As an output, reflects the status of the IRQ line
}

func (s *Board) IRQOut() {
	//As an output, reflects the status of the IRQ line
}

func (s *Board) BusAvailable() bool {
	return s.vic.GetBALow()
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
