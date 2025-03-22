package board

import (
	"fmt"
	"github.com/markel1974/c64emu/src/common/signals"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/hardware/c64/prg"
	"github.com/markel1974/c64emu/src/hardware/vic"
	"github.com/markel1974/c64emu/src/references"
	"os"
)

// intrIrqVicBit represents the interrupt IRQ bit for the VIC.
// intrIrqCia1Bit represents the interrupt IRQ bit for CIA1.
// intrIrqCia2Bit represents the interrupt IRQ bit for CIA2.
// intrIrqExpansionBit represents the interrupt IRQ bit for Expansion.
const (
	intrIrqVicBit       = 0
	intrIrqCia1Bit      = 1
	intrIrqCia2Bit      = 2
	intrIrqExpansionBit = 3
)

// Board represents a complex hardware configuration comprising various component sockets and related system functionalities.
type Board struct {
	*component.BaseComponent
	db              references.IDisplayBuffer
	keysSocket      *KeyboardSocket
	joySocket1      *JoystickSocket
	joySocket2      *JoystickSocket
	quartzSocket    *QuartzSocket
	romSocket       *RomLoaderSocket
	iecSocket       *IECSocket
	picSocket       *PICSocket
	cia1Socket      *CIA1Socket
	cia2Socket      *CIA2Socket
	vicSocket       *VICSocket
	sidSocket       *SIDSocket
	cpuSocket       *CPUSocket
	plaSocket       *PLASocket
	cartSocket      *CartridgeSocket
	throttleSocket  *ThrottleSocket
	expansionSocket *ExpansionSocket
	cfg             *config.Config
	prg             *prg.PRG
	joySwap         bool
	dmaLow          bool
	vBlankSignal    *signals.Signal
	ledSignal       *signals.SignalUint32
	//expansionIrqTrigger *signals.SignalUint32
	//expansionIrqClear   *signals.SignalUint32
}

// NewBoard initializes and returns a new Board instance with specific sockets and properties configured.
// It registers the created instance as a child of the provided parent IComponent.
func NewBoard(parent references.IComponent, factory references.IComponentFactory, instance int) *Board {
	b := &Board{
		BaseComponent:   component.NewBaseComponent(),
		db:              nil,
		keysSocket:      NewKeyboardSocket(),
		joySocket1:      NewJoystickSocket(),
		joySocket2:      NewJoystickSocket(),
		quartzSocket:    NewQuartzSocket(),
		romSocket:       NewRomLoaderSocket(),
		picSocket:       NewPICSocket(),
		cpuSocket:       NewCPUSocket(),
		iecSocket:       NewIECSocket(),
		vicSocket:       NewVICSocket(),
		sidSocket:       NewSIDSocket(),
		cia1Socket:      NewCIA1Socket(),
		cia2Socket:      NewCIA2Socket(),
		plaSocket:       NewPLASocket(),
		cartSocket:      NewCartridgeSocket(),
		throttleSocket:  NewThrottleSocket(),
		expansionSocket: NewExpansionSocket(),
		dmaLow:          false,
		prg:             nil,
		joySwap:         true,
		vBlankSignal:    signals.NewSignal(),
		ledSignal:       signals.NewSignalUint32(),
	}
	b.BaseComponent.Register(factory, parent, Identifier(), instance, b, references.IdIBoardC64(b))
	return b
}

func (s *Board) VBlankSignal() *signals.Signal {
	return s.vBlankSignal
}

func (s *Board) LEDSignal() *signals.SignalUint32 {
	return s.ledSignal
}

func (s *Board) Setup(db references.IDisplayBuffer, player references.IAudioRender, cfg *config.Config) error {
	s.db = db
	s.cfg = cfg
	s.cfg.Bind(s.configChanged)

	quartz, err := references.ComponentToIQuartz(s.GetFactory().Create(s, "quartz", 0))
	if err != nil {
		return err
	}
	throttle, err := references.ComponentToIThrottle(s.GetFactory().Create(s, "dynamic_throttle", 0))
	if err != nil {
		return err
	}
	pic, err := references.ComponentToIPIC6510(s.GetFactory().Create(s, "pic_6510", 0))
	if err != nil {
		return err
	}
	cpu, err := references.ComponentToI6510(s.GetFactory().Create(s, "mos6510", 0))
	if err != nil {
		return err
	}
	iec, err := references.ComponentToIEC(s.GetFactory().Create(s, "iec", 0))
	if err != nil {
		return err
	}
	vic, err := references.ComponentToIVIC(s.GetFactory().Create(s, "mos6569", 0))
	if err != nil {
		return err
	}
	sid, err := references.ComponentToISID(s.GetFactory().Create(s, "mos6581", 0))
	if err != nil {
		return err
	}
	cart, err := references.ComponentToICartridgeManagerC64(s.GetFactory().Create(s, "cartridges_c64", 0))
	if err != nil {
		return err
	}
	rom, err := references.ComponentToIROMLoaderC64(s.GetFactory().Create(s, "roms_c64", 0))
	if err != nil {
		return err
	}
	cia1, err := references.ComponentToICIA(s.GetFactory().Create(s, "mos6526", 0))
	if err != nil {
		return err
	}
	cia2, err := references.ComponentToICIA(s.GetFactory().Create(s, "mos6526", 1))
	if err != nil {
		return err
	}
	pla, err := references.ComponentToIPLAc64(s.GetFactory().Create(s, "pla_c64", 0))
	if err != nil {
		return err
	}
	keys, err := references.ComponentToIKeyboard(s.GetFactory().Create(s, "keyboard_c64", 0))
	if err != nil {
		return err
	}
	joy1, err := references.ComponentToIJoystick(s.GetFactory().Create(s, "joystick_c64", 0))
	if err != nil {
		return err
	}
	joy2, err := references.ComponentToIJoystick(s.GetFactory().Create(s, "joystick_c64", 1))
	if err != nil {
		return err
	}

	if err = s.quartzSocket.Connect(quartz); err != nil {
		return err
	}
	if err = s.expansionSocket.Connect(s, pic, pla, vic, quartz); err != nil {
		return err
	}
	if err = s.throttleSocket.Connect(throttle, mos6569.FrameInterval); err != nil {
		return err
	}
	if err = s.romSocket.Connect(rom, cfg); err != nil {
		return err
	}
	if err = s.picSocket.Connect(pic, quartz); err != nil {
		return err
	}
	if err = s.iecSocket.Connect(iec, s, quartz, cfg); err != nil {
		return err
	}
	if err = s.cpuSocket.Connect(cpu, pic, pla); err != nil {
		return err
	}
	if err = s.vicSocket.Connect(vic, s, db, pic, pla, quartz, cfg); err != nil {
		return err
	}
	if err = s.sidSocket.Connect(sid, player, mos6569.ScreenFreq, mos6569.TotalRasters, cfg); err != nil {
		return err
	}
	if err = s.cia1Socket.Connect(cia1, pic, vic, keys, joy1, joy2); err != nil {
		return err
	}
	if err = s.cia2Socket.Connect(cia2, pic, vic, iec); err != nil {
		return err
	}
	if err = s.cartSocket.Connect(cart, s.expansionSocket, cfg); err != nil {
		return err
	}
	if err = s.plaSocket.Connect(pla, vic, sid, cia1, cia2, cart, rom, cfg); err != nil {
		return err
	}
	if err = s.keysSocket.Connect(keys); err != nil {
		return err
	}
	if err = s.joySocket1.Connect(joy1); err != nil {
		return err
	}
	if err = s.joySocket2.Connect(joy2); err != nil {
		return err
	}

	for _, crt := range cfg.Cartridges() {
		if _, err = s.cartSocket.Add(crt.GetKind(), crt.GetName(), crt.GetData()); err != nil {
			return err
		}
	}

	if prgPath := s.cfg.Prg(); len(prgPath) > 0 {
		if err = s.startPRG(prgPath); err != nil {
			return err
		}
	}

	s.reset()

	s.Print(os.Stdout, " ", true)

	return nil
}

// Reset clears the state of the Board, restoring it back to its initial configuration.
func (s *Board) Reset() {
	//nothing to do
}

// reset resets all the sockets of the Board to their initial state.
func (s *Board) reset() {
	s.picSocket.Reset()
	s.plaSocket.Reset()
	s.cpuSocket.Reset()
	s.cartSocket.Reset()
	s.sidSocket.Reset()
	s.cia1Socket.Reset()
	s.cia2Socket.Reset()
	s.iecSocket.Reset()
	//s.expansion.Reset()
}

// AsyncReset performs an asynchronous reset on various hardware components connected to the board.
func (s *Board) AsyncReset() {
	s.plaSocket.Reset()
	s.picSocket.TriggerReset()
	//s.cpuSocket.AsyncReset()
	s.vicSocket.Reset()
	s.sidSocket.Reset()
	s.cia1Socket.Reset()
	s.cia2Socket.Reset()
	s.iecSocket.Reset()
	//s.expansion.Reset()

	s.dmaLow = false
}

// configChanged is triggered whenever the board's configuration is updated to apply necessary changes.
func (s *Board) configChanged() {
}

// Emulate executes a single emulation cycle, triggering connected sockets and updating the system state.
// Returns the current v
func (s *Board) Emulate() {
	//PHI1
	s.vicSocket.Emulate()

	//PHI2
	s.cia1Socket.Emulate()
	s.cia2Socket.Emulate()
	s.cartSocket.Emulate()
	s.iecSocket.Emulate()
	s.cpuSocket.Emulate()

	s.quartzSocket.Emulate()
}

// Throttle returns the IThrottle instance associated with the board to provide control over throttling behavior.
func (s *Board) Throttle() references.IThrottle {
	return s.throttleSocket
}

// GetText retrieves the textual representation of the board's content as a byte slice from the vicSocket.
func (s *Board) GetText() []byte {
	return s.vicSocket.GetText()
}

// GetScreenSize returns the width and height of the screen as two integers: width (X) and height (Y
func (s *Board) GetScreenSize() (int, int) {
	return mos6569.DisplayX, mos6569.DisplayY
}

// KeyboardSetCommand assigns the specified command to the keyboard's key socket for execution.
func (s *Board) KeyboardSetCommand(cmd string) {
	s.keysSocket.SetCommand(cmd)
}

// KeyboardNumLockToggle toggles the Num Lock state on the keyboard via the associated keys' socket.
func (s *Board) KeyboardNumLockToggle() {
	s.keysSocket.NumLockToggle()
}

// KeyboardCapitalToggle toggles the capitalization state of the keyboard keys.
func (s *Board) KeyboardCapitalToggle() {
	s.keysSocket.CapitalToggle()
}

// SetMouse sets the mouse's position on the board by updating the potentiometer X and Y coordinates.
func (s *Board) SetMouse(x uint8, y uint8) {
	s.sidSocket.SetPotX(x)
	s.sidSocket.SetPotY(y)
}

// KeyboardSetKey sets the state of a virtual key, allowing simulation of key presses
func (s *Board) KeyboardSetKey(pressed bool, vKey int) {
	s.keysSocket.SetKey(pressed, vKey)
}

// Joy1SetKey sets the
func (s *Board) Joy1SetKey(pressed bool, vKey int) {
	if s.joySwap {
		s.joySocket2.SetKey(pressed, vKey)
	} else {
		s.joySocket1.SetKey(pressed, vKey)
	}
}

// Joy2SetKey handles key press events for joystick input, delegating the action to the appropriate joystick socket based on joySwap.
func (s *Board) Joy2SetKey(pressed bool, vKey int) {
	if s.joySwap {
		s.joySocket1.SetKey(pressed, vKey)
	} else {
		s.joySocket2.SetKey(pressed, vKey)
	}
}

// Joystick1Move processes joystick movement and button states for the first joystick or swaps to the second based on joySwap.
func (s *Board) Joystick1Move(x uint, y uint, buttons uint) {
	if s.joySwap {
		s.joySocket2.Move(x, y, buttons)
	} else {
		s.joySocket1.Move(x, y, buttons)
	}
}

// Joystick2Move controls the movement of the joystick depending on the joySwap flag, updating either joySocket1 or joySocket2.
func (s *Board) Joystick2Move(x uint, y uint, buttons uint) {
	if s.joySwap {
		s.joySocket1.Move(x, y, buttons)
	} else {
		s.joySocket2.Move(x, y, buttons)
	}
}

// JoySwap toggles the joystick swap state and resets both joystick sockets if the pressed parameter is true.
func (s *Board) JoySwap() {
	s.joySwap = !s.joySwap
	s.joySocket1.Reset()
	s.joySocket2.Reset()
}

// ExtRamWrite writes a single byte of data to the external RAM at the specified address using the provided memory configuration.
func (s *Board) ExtRamWrite(memConfig int, addr uint16, data uint8) {
	var prev []uint8 = nil
	if memConfig >= 0 {
		prev = s.plaSocket.GetMemoryConfig()
		s.plaSocket.SetMemoryEntry(uint8(memConfig))
	}
	s.plaSocket.Write(addr, data)
	if prev != nil {
		s.plaSocket.SetMemoryConfig(prev)
	}
}

// ExtRamRead reads a byte from external RAM at the specified address using the given memory configuration.
func (s *Board) ExtRamRead(memConfig int, addr uint16) uint8 {
	var prev []uint8 = nil
	if memConfig >= 0 {
		prev = s.plaSocket.GetMemoryConfig()
		s.plaSocket.SetMemoryEntry(uint8(memConfig))
	}
	rb := s.plaSocket.Read(addr)
	if prev != nil {
		s.plaSocket.SetMemoryConfig(prev)
	}
	return rb
}

// DMALowTrigger sets the DMA low trigger state, enabling or disabling the CPU bus release for direct memory access.
// If set to true, the CPU enters a wait state, and its bus lines go to high resistance, allowing external hardware access.
// The method adjusts the DMA state and configures the CPU's RDY and AEC lines based on the VIC-II chip's signals.
func (s *Board) DMALowTrigger(v bool) {
	//If _DMA=Low the CPU can be requested to release the bus.
	//It will stop after the next read cycle, and all bus lines will go to high resistance state.
	//So other units can use the computer hardware. At _DMA=High the CPU continues to work.
	//The DMA line on the expansion port gets pulled low, or the VIC-II's BA line goes low.
	//The DMA line is used to put the CPU in a wait state.
	//The DMA line also forces the CPU's AEC line low, so while it's waiting, its R/W, address bus and data bus lines are put in HighZ,
	//so they don't have any influence over the buses.
	//This allows a device on the expansion port, such as an REU, to perform direct memory accesses (DMA) to the main RAM.
	s.dmaLow = v
	s.cpuSocket.SetRDYLow(s.dmaLow || s.vicSocket.GetBALow())
	s.cpuSocket.SetAECLow(s.dmaLow || s.vicSocket.GetAECLow())
}

// RDYLowTrigger sets the RDY signal based on the logical AND between BA and DMA produced by chip U27.
// It also updates the RDY state taking into consideration the provided value and the dmaLow state.
func (s *Board) RDYLowTrigger(v bool) {
	//The RDY signal the result of logical AND between BA and DMA produced by the chip U27
	s.cpuSocket.SetRDYLow(v || s.dmaLow)
	//TODO SIGNAL
}

// AECLowTrigger sets the AEC (Asynchronous Enable Control) low signal for the CPU socket based on the given value.
func (s *Board) AECLowTrigger(v bool) {
	s.cpuSocket.SetAECLow(v || s.dmaLow)
	//TODO SIGNAL
}

// LastCycleTrigger performs the final preparation step by invoking the Prepare method on the sidSocket instance.
func (s *Board) LastCycleTrigger() {
	s.sidSocket.Prepare()
}

// VBlankTrigger sets the vBlank flag to true and updates connected sockets and programmable logic if applicable.
func (s *Board) VBlankTrigger() {
	s.vBlankSignal.Emit()

	s.sidSocket.Update()
	s.cia1Socket.Update()
	s.cia2Socket.Update()
	if s.prg != nil {
		if s.prg.Inject(s.vicSocket.GetText()) {
			s.prg = nil
		}
	}
}

// LedTrigger handles LED state change events for a specified device and updates the LED state accordingly.
func (s *Board) LedTrigger(state uint32) {
	s.ledSignal.Emit(state)
}

// startPR
func (s *Board) startPRG(prgPath string) error {
	s.prg = prg.NewPRG(s.plaSocket, s.keysSocket)
	if err := s.prg.Load(prgPath); err != nil {
		s.prg = nil
		return fmt.Errorf("can't load prg: %s", err.Error())
	}
	return nil
}
