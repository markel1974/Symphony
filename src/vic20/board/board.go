package board

import (
	"github.com/markel1974/c64emu/src/c64/cartridges"
	"github.com/markel1974/c64emu/src/c64/inputs"
	"github.com/markel1974/c64emu/src/c64/pla"
	"github.com/markel1974/c64emu/src/c64/prg"
	"github.com/markel1974/c64emu/src/common/signals"
	mos6510 "github.com/markel1974/c64emu/src/components/6510"
	"github.com/markel1974/c64emu/src/components/board"
	"github.com/markel1974/c64emu/src/components/iec"
	"github.com/markel1974/c64emu/src/components/throttling"
	mos6569 "github.com/markel1974/c64emu/src/components/vic"
	"github.com/markel1974/c64emu/src/config"
	"golang.design/x/clipboard"
	"log"
)

const (
	intrIrqVicBit       = 0
	intrIrqCia1Bit      = 1
	intrIrqCia2Bit      = 2
	intrIrqExpansionBit = 3
)

type Board struct {
	cia1Socket          *CIA1Socket
	cia2Socket          *CIA2Socket
	vicSocket           *VicSocket
	cpuSocket           *CPUSocket
	expansion           *Expansion
	db                  board.IDisplayBuffer
	p                   board.IPlayer
	pic                 *mos6510.Pic
	iec                 *iec.IEC
	keys                *inputs.Keyboard
	joy1                *inputs.Joystick
	joy2                *inputs.Joystick
	joySwap             bool
	cfg                 *config.Config
	hasClipboard        bool
	cartMan             *cartridges.Manager
	pla                 *pla.PLA
	expansionIrqTrigger *signals.SignalUint32
	expansionIrqClear   *signals.SignalUint32
	vBlank              bool
	lastVicCycle        bool
	dmaLow              bool
	prg                 *prg.PRG
	dt                  board.IThrottling
}

func NewBoard() *Board {
	b := &Board{
		iec:                 nil,
		pic:                 nil,
		keys:                nil,
		joy1:                nil,
		joy2:                nil,
		hasClipboard:        false,
		cartMan:             cartridges.NewManager(),
		pla:                 nil,
		expansionIrqTrigger: nil,
		expansionIrqClear:   nil,
		vBlank:              false,
		lastVicCycle:        false,
		dmaLow:              false,
		prg:                 nil,
		joySwap:             true,
		dt:                  nil,
	}
	return b
}

func (s *Board) Setup(db board.IDisplayBuffer, p board.IPlayer, cfg *config.Config) error {
	s.db = db
	s.p = p
	if err := clipboard.Init(); err != nil {
		log.Printf("can't init clipboard: %s", err)
	} else {
		s.hasClipboard = true
	}
	s.cfg = cfg
	s.cfg.Bind(s.configChanged)
	s.dt = throttling.NewDynamicThrottling(mos6569.FrameInterval)

	s.cpuSocket = NewCPUSocket()
	s.vicSocket = NewVicSocket()
	//s.sidSocket = NewSidSocket()
	s.cia1Socket = NewCIA1Socket()
	s.cia2Socket = NewCIA2Socket()

	const componentId = "vic20"
	s.pic = mos6510.NewPic(componentId, "")
	s.iec = iec.NewIEC(componentId, "")
	s.keys = inputs.NewKeyboard()
	s.joy1 = inputs.NewJoystick()
	s.joy2 = inputs.NewJoystick()
	s.pla = pla.NewPLA(componentId, "")
	s.expansion = NewExpansion(s)

	s.iec.Setup(cfg)
	s.cpuSocket.Setup(s)
	s.vicSocket.Setup(s, intrIrqVicBit)
	//s.sidSocket.Setup(s)
	s.cia1Socket.Setup(s, intrIrqCia1Bit)
	s.cia2Socket.Setup(s, intrIrqCia2Bit)
	s.cartMan.Setup(s.expansion, cfg)

	s.reset()

	return nil
}

func (s *Board) reset() {
	s.pic.Reset()
	s.pla.Reset()
	s.cpuSocket.Reset()
	s.cartMan.Reset()
	//s.sidSocket.Reset()
	s.cia1Socket.Reset()
	s.cia2Socket.Reset()
	s.iec.Reset()
	s.expansion.Reset()
}

func (s *Board) AsyncReset() {
	s.pla.Reset()
	s.pic.TriggerReset()
	//s.cpuSocket.AsyncReset()
	s.vicSocket.Reset()
	//s.sidSocket.Reset()
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

	s.cartMan.Emulate()
	s.iec.Emulate()

	//s.quartz.AddCycle()

	return s.vBlank
}

func (s *Board) GetText() []byte {
	//return s.vic.GetText()
	return nil
}

func (s *Board) GetScreenSize() (int, int) {
	return 320, 200
}

func (s *Board) DiskChange() {
	//s.cfg.GetDrives()
	//aaa

	s.cfg.SwitchDisk()
	s.cfg.SetDriveOpt("", 8)
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

func (s *Board) KeyboardSetCommand(cmd string) {
	s.keys.SetCommand(cmd)
}

func (s *Board) KeyboardNumLockToggle() {
	s.keys.NumLockToggle()
}

func (s *Board) KeyboardCapitalToggle() {
	s.keys.CapitalToggle()
}

func (s *Board) SetMouse(x uint8, y uint8) {
	//s.sidSocket.SetPotXY(x, y)
}

func (s *Board) KeyboardSetVirtualKey(pressed bool, vKey int) {
	s.keys.SetVirtualKey(pressed, vKey)
}

func (s *Board) Joy1SetKey(pressed bool, vKey int) {
	if s.joySwap {
		s.joy2.SetKey(pressed, vKey)
	} else {
		s.joy1.SetKey(pressed, vKey)
	}
}

func (s *Board) Joy2SetKey(pressed bool, vKey int) {
	if s.joySwap {
		s.joy1.SetKey(pressed, vKey)
	} else {
		s.joy2.SetKey(pressed, vKey)
	}
}

func (s *Board) Joystick1Move(x uint, y uint, buttons uint) {
	if s.joySwap {
		s.joy2.Move(x, y, buttons)
	} else {
		s.joy1.Move(x, y, buttons)
	}
}

func (s *Board) Joystick2Move(x uint, y uint, buttons uint) {
	if s.joySwap {
		s.joy1.Move(x, y, buttons)
	} else {
		s.joy2.Move(x, y, buttons)
	}
}

func (s *Board) JoySwap(pressed bool) {
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
		prev = s.pla.GetMemoryConfig()
		s.pla.SetMemoryEntry(uint8(memConfig))
	}
	s.pla.Write(addr, data)
	if prev != nil {
		s.pla.SetMemoryConfig(prev)
	}
}

func (s *Board) ExtRamRead(memConfig int, addr uint16) uint8 {
	var prev []uint8 = nil
	if memConfig >= 0 {
		prev = s.pla.GetMemoryConfig()
		s.pla.SetMemoryEntry(uint8(memConfig))
	}
	rb := s.pla.Read(addr)
	if prev != nil {
		s.pla.SetMemoryConfig(prev)
	}
	return rb
}

func (s *Board) Throttle() board.IThrottling {
	return s.dt
}
