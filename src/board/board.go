package board

import (
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
)

type PhiMode int

const (
	PhiIdle = PhiMode(0)
	Phi1    = PhiMode(1)
	Phi2    = PhiMode(2)
)

type Board struct {
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
	cart         cartridges.ICartridge
	cartFactory  *cartridges.Factory
	banks        *Banks
}

func NewBoard() *Board {
	b := &Board{
		quartz:       quartz.NewQuartz(),
		iec:          nil,
		cpu:          nil,
		vic:          nil,
		sid:          nil,
		cia1:         nil,
		cia2:         nil,
		interrupts:   nil,
		keys:         nil,
		cart:         nil,
		hasClipboard: false,
		cartFactory:  cartridges.NewFactory(),
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
	s.vic = vic.NewMOS6569()
	s.sid = sid.NewMOS6581()
	s.cia1 = cia.NewMOS6526_1()
	s.cia2 = cia.NewMOS6526_2()
	s.keys = keyboard.NewKeyboard()
	s.banks = NewBanks()

	s.cpu.Setup(s, prefs)
	s.interrupts = s.cpu.GetInterrupts()
	s.vic.Setup(s, prefs)
	s.sid.Setup(s, prefs)
	s.cia1.Setup(s, prefs)
	s.cia2.Setup(s, prefs)
	s.iec.Setup(s, prefs)
	s.cartFactory.Setup(s, prefs)

	s.banks.Setup(s, prefs)

	if !s.prefs.GetDisableCartridgeAutostart() {
		if cartFile := s.prefs.GetCartridge(); len(cartFile) > 0 {
			if cart, err := s.cartFactory.Load(cartFile); err == nil {
				s.cart = cart
				//s.portsUpdate()
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
}

func (s *Board) AsyncReset() {
	s.keys.Reset()
	s.banks.AsyncReset()
	s.cpu.AsyncReset()
	s.vic.Reset()
	s.sid.Reset()
	s.cia1.Reset()
	s.cia2.Reset()
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

	s.quartz.AddCycle()
	s.phiMode = PhiIdle
	return vBlank
}

func (s *Board) RmwFlags() uint8 {
	//TODO IMPLEMENT cpu rmw flags
	return 0
}

func (s *Board) Cycle() uint64 {
	return s.quartz.Cycle()
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

func (s *Board) KeyboardSetExt(ext bool) {
	s.keys.SetExt(ext)
}

func (s *Board) KeyboardSetNumLock(numLock bool) {
	s.keys.SetNumLock(numLock)
}

func (s *Board) KeyboardSetCapital(capital bool) {
	s.keys.SetCapital(capital)
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

func (s *Board) VICTriggerIRQ() {
	s.interrupts.TriggerVICIRQ()
}

func (s *Board) VICClearIRQ() {
	s.interrupts.ClearVICIRQ()
}

func (s *Board) VICChangedVA(d uint8) {
	s.vic.ChangedVA(d)
}

func (s *Board) VICLightPenTrigger() {
	s.vic.TriggerLightPen()
}

func (s *Board) GetDisplayBuffer() []byte {
	return s.vic.GetDisplayBuffer()
}

func (s *Board) CpuRamRead(addr uint16) uint8 {
	return s.banks.Read(addr)
}

func (s *Board) CpuRamWrite(addr uint16, data uint8) {
	s.banks.Write(addr, data)
}

func (s *Board) ColorRead(addr uint16) uint8 {
	return s.banks.ReadColor(addr)
}

//func (s *Board) ColorWrite(addr uint16, data uint8) {
//	s.banks.WriteColor(addr, data)
//}

func (s *Board) RamRead(addr uint16) uint8 {
	return s.banks.ReadDirect(addr)
}

func (s *Board) RamWrite(addr uint16, data uint8) {
	s.banks.WriteDirect(addr, data)
}

//func (s *Board) BasicRomRead(addr uint16) uint8 {
//	return s.banks.ReadBasicRom(addr)
//}

func (s *Board) CharRomRead(addr uint16) uint8 {
	return s.banks.ReadCharRom(addr)
}

//func (s *Board) KernalRomRead(addr uint16) uint8 {
//	return s.banks.ReadKernalRom(addr)
//}

func (s *Board) NMITrigger() {
	s.interrupts.TriggerNMI()
}

func (s *Board) NMIClear() {
	s.interrupts.ClearNMI()
}

func (s *Board) CIATriggerIRQ() {
	s.interrupts.TriggerCIAIRQ()
}

func (s *Board) CIAClearIRQ() {
	s.interrupts.ClearCIAIRQ()
}

func (s *Board) BusCpuWrite(d uint8) {
	s.iec.CpuWrite(d)
}

func (s *Board) BusCpuRead() uint8 {
	return s.iec.CpuRead()
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

func (s *Board) ExtRamWrite(memConfig int, addr uint16, data uint8) {
	var prev []uint8 = nil
	if memConfig >= 0 {
		prev = s.banks.GetMemoryConfig()
		s.banks.SetMemoryEntry(uint8(memConfig))
	}
	s.CpuRamWrite(addr, data)
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
	rb := s.CpuRamRead(addr)
	if prev != nil {
		s.banks.SetMemoryConfig(prev)
	}
	return rb
}
