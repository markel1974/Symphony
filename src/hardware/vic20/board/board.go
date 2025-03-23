package board

import (
	"github.com/markel1974/c64emu/src/common/signals"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/hardware/c64/prg"
	"github.com/markel1974/c64emu/src/hardware/cartridges_c64"
	"github.com/markel1974/c64emu/src/hardware/iec"
	"github.com/markel1974/c64emu/src/hardware/joystick_c64"
	"github.com/markel1974/c64emu/src/hardware/keyboard_c64"
	"github.com/markel1974/c64emu/src/hardware/pic_6510"
	"github.com/markel1974/c64emu/src/hardware/pla_c64"
	"github.com/markel1974/c64emu/src/references"
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
	*component.BaseComponent
	cia1Socket          *CIA1Socket
	cia2Socket          *CIA2Socket
	vicSocket           *VicSocket
	cpuSocket           *CPUSocket
	expansion           *Expansion
	db                  references.IDisplayBuffer
	p                   references.IAudioRender
	pic                 *pic_6510.Pic
	iec                 *iec.Dispatcher
	keys                *keyboard_c64.Keyboard
	joy1                *joystick_c64.Joystick
	joy2                *joystick_c64.Joystick
	joySwap             bool
	cfg                 *config.Config
	hasClipboard        bool
	cartMan             *cartridges_c64.Manager
	pla                 *pla_c64.PLA
	expansionIrqTrigger *signals.SignalUint32
	expansionIrqClear   *signals.SignalUint32
	vBlankSignal        *signals.Signal
	ledSignal           *signals.SignalUint32
	lastVicCycle        bool
	dmaLow              bool
	prg                 *prg.PRG
	dt                  references.IThrottle
}

func NewBoard(parent references.IComponent, factory references.IComponentFactory, instance int) *Board {
	b := &Board{
		BaseComponent:       component.NewBaseComponent(),
		iec:                 nil,
		pic:                 nil,
		keys:                nil,
		joy1:                nil,
		joy2:                nil,
		hasClipboard:        false,
		cartMan:             nil,
		pla:                 nil,
		expansionIrqTrigger: nil,
		expansionIrqClear:   nil,
		vBlankSignal:        signals.NewSignal(),
		ledSignal:           signals.NewSignalUint32(),
		lastVicCycle:        false,
		dmaLow:              false,
		prg:                 nil,
		joySwap:             true,
		dt:                  nil,
	}
	b.BaseComponent.Register(factory, parent, Identifier(), instance, b, references.IdIBoardVIC20(b))

	return b
}

func (s *Board) VBlankSignal() *signals.Signal {
	return s.vBlankSignal
}

func (s *Board) LEDSignal() *signals.SignalUint32 {
	return s.ledSignal
}

func (s *Board) Setup(db references.IDisplayBuffer, p references.IAudioRender, cfg *config.Config) error {
	s.cartMan = cartridges_c64.NewManager(s, s.GetFactory(), 0)
	s.db = db
	s.p = p
	if err := clipboard.Init(); err != nil {
		log.Printf("can't init clipboard: %s", err)
	} else {
		s.hasClipboard = true
	}
	s.cfg = cfg
	s.cfg.Bind(s.configChanged)
	//s.dt = dynamic_throttle.NewDynamicThrottle(s, s.factory, "")
	//s.dt.SetInterval(mos6569.FrameInterval)

	s.cpuSocket = NewCPUSocket()
	s.vicSocket = NewVicSocket()
	//s.sidSocket = NewSidSocket()
	s.cia1Socket = NewCIA1Socket()
	s.cia2Socket = NewCIA2Socket()

	s.pic = pic_6510.NewPIC(s, s.GetFactory(), 0)
	s.iec = iec.NewDispatcher(s, s.GetFactory(), 0)
	s.keys = keyboard_c64.NewKeyboard(s, s.GetFactory(), 0)
	s.joy1 = joystick_c64.NewJoystick(s, s.GetFactory(), 0)
	s.joy2 = joystick_c64.NewJoystick(s, s.GetFactory(), 1)
	s.pla = pla_c64.NewPLA(s, s.GetFactory(), 0)
	s.expansion = NewExpansion(s)

	//s.iec.Setup(s.quartz, cfg)
	s.cpuSocket.Setup(s)
	s.vicSocket.Setup(s, intrIrqVicBit)
	//s.sidSocket.Setup(s)
	s.cia1Socket.Setup(s, intrIrqCia1Bit)
	s.cia2Socket.Setup(s, intrIrqCia2Bit)
	if err := s.cartMan.Setup(s.expansion, cfg); err != nil {
		return err
	}

	s.reset()

	return nil
}

func (s *Board) Reset() {
	//
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

func (s *Board) Emulate() {
	s.cartMan.Emulate()
	s.iec.Emulate()
	//s.quartz.AddCycle()
}

func (s *Board) EmulationRequired() bool {
	return false
}

func (s *Board) GetText() []byte {
	//return s.vic.GetText()
	return nil
}

func (s *Board) GetScreenSize() (int, int) {
	return 320, 200
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

func (s *Board) KeyboardSetKey(pressed bool, vKey int) {
	s.keys.SetKey(pressed, vKey)
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

func (s *Board) JoySwap() {
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

func (s *Board) Throttle() references.IThrottle {
	return s.dt
}
