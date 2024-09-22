package board

import (
	"github.com/markel1974/c64emu/src/c64/banks"
	"github.com/markel1974/c64emu/src/c64/cartridges"
	"github.com/markel1974/c64emu/src/c64/iec"
	"github.com/markel1974/c64emu/src/c64/inputs"
	"github.com/markel1974/c64emu/src/c64/prg"
	"github.com/markel1974/c64emu/src/components/6510"
	"github.com/markel1974/c64emu/src/components/cia"
	"github.com/markel1974/c64emu/src/components/quartz"
	"github.com/markel1974/c64emu/src/components/sid"
	"github.com/markel1974/c64emu/src/components/vic"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/signals"
	"golang.design/x/clipboard"
	"log"
	"os"
)

const (
	intrIrqVicBit       = 0
	intrIrqCia1Bit      = 1
	intrIrqCia2Bit      = 2
	intrIrqExpansionBit = 3
)

const baseId = "c64"

type Board struct {
	db                  mos6569.IDisplayBuffer
	player              mos6581.IPlayer
	quartz              *quartz.Quartz
	cpu                 *mos6510.CPU
	vic                 *mos6569.VIC
	sid                 *mos6581.SID
	cia1                *mos6526.CIA
	cia2                *mos6526.CIA
	cia1Socket          *CIA1Socket
	cia2Socket          *CIA2Socket
	vicSocket           *VicSocket
	sidSocket           *SidSocket
	cpuSocket           *CPUSocket
	expansion           *Expansion
	pic                 *mos6510.Pic
	iec                 *iec.IEC
	keys                *inputs.Keyboard
	joy1                *inputs.Joystick
	joy2                *inputs.Joystick
	joySwap             bool
	cfg                 *config.Config
	hasClipboard        bool
	cartMan             *cartridges.Manager
	banks               *banks.Banks
	expansionIrqTrigger *signals.SignalUint32
	expansionIrqClear   *signals.SignalUint32
	vBlank              bool
	lastVicCycle        bool
	dmaLow              bool
	prg                 *prg.PRG
}

func NewBoard(db mos6569.IDisplayBuffer, player mos6581.IPlayer) *Board {
	b := &Board{
		db:                  db,
		player:              player,
		quartz:              quartz.NewQuartz(),
		iec:                 nil,
		cpu:                 nil,
		vic:                 nil,
		sid:                 nil,
		cia1:                nil,
		cia2:                nil,
		pic:                 nil,
		keys:                nil,
		joy1:                nil,
		joy2:                nil,
		hasClipboard:        false,
		cartMan:             cartridges.NewManager(),
		banks:               nil,
		expansionIrqTrigger: nil,
		expansionIrqClear:   nil,
		vBlank:              false,
		lastVicCycle:        false,
		dmaLow:              false,
		prg:                 nil,
		joySwap:             true,
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

	s.cpuSocket = NewCPUSocket()
	s.vicSocket = NewVicSocket()
	s.sidSocket = NewSidSocket()
	s.cia1Socket = NewCIA1Socket()
	s.cia2Socket = NewCIA2Socket()

	s.pic = mos6510.NewPic()
	s.iec = iec.NewIEC()
	s.cpu = mos6510.NewCPU(baseId + "_cpu")
	s.vic = mos6569.NewVIC(baseId + "_vic")
	s.sid = mos6581.NewSID(baseId + "_sid")
	s.cia1 = mos6526.NewCIA(baseId + "_cia1")
	s.cia2 = mos6526.NewCIA(baseId + "_cia2")
	s.keys = inputs.NewKeyboard()
	s.joy1 = inputs.NewJoystick()
	s.joy2 = inputs.NewJoystick()
	s.banks = banks.NewBanks()
	s.expansion = NewExpansion(s)

	s.pic.Setup(s.quartz)
	s.iec.Setup(cfg)
	s.cpuSocket.Setup(s)
	s.cpu.Setup(s.cpuSocket)
	s.vicSocket.Setup(s, intrIrqVicBit)
	s.vic.Setup(s.vicSocket, cfg)
	s.sidSocket.Setup(s)
	s.sid.Setup(s.sidSocket, cfg, mos6569.ScreenFreq, mos6569.TotalRasters)
	s.cia1Socket.Setup(s, intrIrqCia1Bit)
	s.cia1.Setup(s.cia1Socket)
	s.cia2Socket.Setup(s, intrIrqCia2Bit)
	s.cia2.Setup(s.cia2Socket)
	s.cartMan.Setup(s.expansion, cfg)
	s.banks.Setup(s.vic, s.sid, s.cia1, s.cia2, s.cartMan, cfg)

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

	if prgPath := s.cfg.GetPrg(); len(prgPath) > 0 {
		s.prg = prg.NewPRG(banks.NewObserver(s.banks), s.keys)
		if err := s.prg.Load(prgPath); err != nil {
			log.Printf("can't load prg: %s", err.Error())
			s.prg = nil
		}
	}

	s.reset()

	return nil
}

func (s *Board) reset() {
	s.pic.Reset()
	s.banks.Reset()
	s.cpuSocket.Reset()
	s.cartMan.Reset()
	s.sidSocket.Reset()
	s.cia1Socket.Reset()
	s.cia2Socket.Reset()
	s.iec.Reset()
	s.expansion.Reset()
}

func (s *Board) AsyncReset() {
	s.banks.Reset()
	s.pic.TriggerReset()
	//s.cpuSocket.AsyncReset()
	s.vicSocket.Reset()
	s.sidSocket.Reset()
	s.cia1Socket.Reset()
	s.cia2Socket.Reset()
	s.iec.Reset()
	s.expansion.Reset()

	s.dmaLow = false
}

func (s *Board) configChanged() {
}

func (s *Board) Emulate() bool {
	s.vBlank = false

	//PHI1
	s.vic.Emulate()

	//PHI2
	s.cia1.Emulate()
	s.cia2.Emulate()
	s.cartMan.Emulate()
	s.iec.Emulate()
	s.cpu.Emulate()

	s.quartz.AddCycle()

	return s.vBlank
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
	s.keys.SetCommand(string(data))
}

func (s *Board) KeyboardNumLockToggle() {
	s.keys.NumLockToggle()
}

func (s *Board) KeyboardCapitalToggle() {
	s.keys.CapitalToggle()
}

func (s *Board) SetMouse(x uint8, y uint8) {
	s.sidSocket.SetPotXY(x, y)
}

func (s *Board) KeyboardSetVirtualKey(pressed bool, vKey int) {
	s.keys.SetVirtualKey(pressed, vKey)
}

func (s *Board) JoystickSetKey(pressed bool, vKey int) {
	if s.joySwap {
		s.joy2.SetKey(pressed, vKey)
	} else {
		s.joy1.SetKey(pressed, vKey)
	}
}

func (s *Board) JoystickSwap(pressed bool) {
	if !pressed {
		return
	}
	s.joySwap = !s.joySwap
	s.joy1.Reset()
	s.joy2.Reset()
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

func (s *Board) dmaLowSlot(v bool) {
	//if _DMA=Low the CPU can be requested to release the bus.
	//It will stop after the next read cycle and all bus lines will go to high resistance state.
	//So other units can use the computer hardware. At _DMA=High the CPU continues to work.
	//The DMA line on the expansion port gets pulled low, or the VIC-II's BA line goes low.
	//The DMA line is used to put the CPU in a wait state.
	//The DMA line also forces the CPU's AEC line low so while it's waiting its R/W, address bus and data bus lines are put in HighZ,
	//so they don't have any influence over the buses.
	//This allows a device on the expansion port, such as an REU, to perform direct memory accesses (DMA) to the main RAM.
	s.dmaLow = v
	s.cpu.SetRDYLow(s.dmaLow || s.vic.GetBALow())
	s.cpu.SetAECLow(s.dmaLow || s.vic.GetAECLow())
}

func (s *Board) rdyLowSlot(v bool) {
	//The RDY signal the result of logical AND between BA and DMA produced by the chip U27
	s.cpu.SetRDYLow(v || s.dmaLow)
	//TODO SIGNAL
}

func (s *Board) aecLowSlot(v bool) {
	s.cpu.SetAECLow(v || s.dmaLow)
	//TODO SIGNAL
}

func (s *Board) irqTriggerSlot(i uint32) {
	s.pic.TriggerIRQ(i)
	if s.expansionIrqTrigger != nil {
		s.expansionIrqTrigger.Emit(i)
	}
}

func (s *Board) irqClearSlot(i uint32) {
	s.pic.ClearIRQ(i)
	if s.expansionIrqClear != nil {
		s.expansionIrqClear.Emit(i)
	}
}

func (s *Board) nmiTriggerSlot() {
	s.pic.TriggerNMI()
}

func (s *Board) nmiClearSlot() {
	s.pic.ClearNMI()
}

func (s *Board) vicLastCycleSLot() {
	s.sidSocket.Update()
}

func (s *Board) vicVBlankSlot() {
	s.vBlank = true

	//TODO MOVE
	s.sidSocket.Render()

	s.cia1Socket.Update()
	s.cia2Socket.Update()

	if s.prg != nil {
		if s.prg.Inject(s.vic.GetText()) {
			s.prg = nil
		}
	}
}
