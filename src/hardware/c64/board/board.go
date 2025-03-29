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
	keysSocket       *KeyboardSocket
	joySocket1       *JoystickSocket
	joySocket2       *JoystickSocket
	quartzSocket     *QuartzSocket
	romSocket        *RomLoaderSocket
	iecSocket        *IECSocket
	picSocket        *PICSocket
	cia1Socket       *CIA1Socket
	cia2Socket       *CIA2Socket
	vicSocket        *VICSocket
	sidSocket        *SIDSocket
	cpuSocket        *CPUSocket
	plaSocket        *PLASocket
	cartSocket       *CartridgeSocket
	throttleSocket   *ThrottleSocket
	expansionSocket  *ExpansionSocket
	cfg              *config.Config
	prg              *prg.PRG
	joySwap          bool
	dmaLow           bool
	vBlankSignal     *signals.Signal
	ledSignal        *signals.SignalUint32
	emulation        []func()
	sockets          []references.ISocket
	components       map[string]references.IComponent
	hardwareSequence []string
	label            string
	//expansionIrqTrigger *signals.SignalUint32
	//expansionIrqClear   *signals.SignalUint32
}

// NewBoard initializes and returns a new Board instance with specific sockets and properties configured.
// It registers the created instance as a child of the provided parent IComponent.
func NewBoard(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Board {
	s := &Board{
		BaseComponent: component.NewBaseComponent(),
		dmaLow:        false,
		prg:           nil,
		joySwap:       true,
		vBlankSignal:  signals.NewSignal(),
		ledSignal:     signals.NewSignalUint32(),
		emulation:     []func(){},
		label:         label,
	}
	s.BaseComponent.Register(factory, parent, Identifier(), s, references.IdIBoardC64(s, label, instance))

	s.hardwareSequence = []string{
		references.IdIVIC(nil, label, 0),
		references.IdICIA(nil, label, 0),
		references.IdICIA(nil, label, 1),
		references.IdICartridgeManagerC64(nil, label, 0),
		references.IdIIec(nil, label, 0),
		references.IdI6510(nil, label, 0),
		references.IdIQuartz(nil, label, 0),
	}

	s.keysSocket = NewKeyboardSocket(s, s.label)
	s.joySocket1 = NewJoystickSocket(s, s.label, 0)
	s.joySocket2 = NewJoystickSocket(s, s.label, 1)
	s.romSocket = NewRomLoaderSocket(s, s.label)
	s.quartzSocket = NewQuartzSocket(s, s.label)
	s.iecSocket = NewIECSocket(s, s.label, s)
	s.expansionSocket = NewExpansionSocket(s, s.label, s)
	s.cartSocket = NewCartridgeSocket(s, s.label, s.expansionSocket)
	s.picSocket = NewPICSocket(s, s.label)
	s.cpuSocket = NewCPUSocket(s, s.label)
	s.vicSocket = NewVICSocket(s, s.label, s)
	s.sidSocket = NewSIDSocket(s, s.label, mos6569.ScreenFreq, mos6569.TotalRasters)
	s.cia1Socket = NewCIA1Socket(s, s.label)
	s.cia2Socket = NewCIA2Socket(s, s.label)
	s.plaSocket = NewPLASocket(s, s.label)
	s.throttleSocket = NewThrottleSocket(s, s.label, mos6569.FrameInterval)

	s.sockets = append(s.sockets, s.romSocket)
	s.sockets = append(s.sockets, s.quartzSocket)
	s.sockets = append(s.sockets, s.keysSocket)
	s.sockets = append(s.sockets, s.joySocket1)
	s.sockets = append(s.sockets, s.joySocket2)
	s.sockets = append(s.sockets, s.iecSocket)
	s.sockets = append(s.sockets, s.expansionSocket)
	s.sockets = append(s.sockets, s.cartSocket)
	s.sockets = append(s.sockets, s.picSocket)
	s.sockets = append(s.sockets, s.cpuSocket)
	s.sockets = append(s.sockets, s.vicSocket)
	s.sockets = append(s.sockets, s.sidSocket)
	s.sockets = append(s.sockets, s.cia1Socket)
	s.sockets = append(s.sockets, s.cia2Socket)
	s.sockets = append(s.sockets, s.plaSocket)
	s.sockets = append(s.sockets, s.throttleSocket)

	return s
}

func (s *Board) VBlankSignal() *signals.Signal {
	return s.vBlankSignal
}

func (s *Board) LEDSignal() *signals.SignalUint32 {
	return s.ledSignal
}

func (s *Board) Setup(components map[string]references.IComponent, cfg *config.Config) error {
	s.cfg = cfg
	s.components = components
	return nil
}

func (s *Board) Connect() error {
	for _, c := range s.sockets {
		if err := c.Mount(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Board) Start() error {
	var err error
	if s.emulation, err = s.rebuildEmulation(s.hardwareSequence, s.components); err != nil {
		return err
	}
	if err = s.iecSocket.CreatePeripherals(); err != nil {
		return err
	}
	if err = s.cartSocket.CreateCartridges(); err != nil {
		return err
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

func (s *Board) Internal() bool {
	return false
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

// Emulate executes all functions in the Board's emulation sequence in order.
func (s *Board) Emulate() {
	for _, fn := range s.emulation {
		fn()
	}
}

func (s *Board) EmulationRequired() bool {
	return true
}

// GetText retrieves the textual representation of the board's content as a byte slice from the vicSocket.
func (s *Board) GetText() []byte {
	return s.vicSocket.GetText()
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
	s.throttleSocket.Update()
}

// LedTrigger handles LED state change events for a specified device and updates the LED state accordingly.
func (s *Board) LedTrigger(state uint32) {
	s.ledSignal.Emit(state)
}

// startPRG initializes and loads a PRG file into the Board using the specified path, returning an error if loading fails.
func (s *Board) startPRG(prgPath string) error {
	s.prg = prg.NewPRG(s.plaSocket, s.keysSocket)
	if err := s.prg.Load(prgPath); err != nil {
		s.prg = nil
		return fmt.Errorf("can't load prg: %s", err.Error())
	}
	return nil
}

// rebuildEmulation constructs a sequence of emulation functions based on the given components and hardware sequence.
// Returns the constructed sequence of emulation functions or an error if the sequence is incomplete.
func (s *Board) rebuildEmulation(seq []string, components map[string]references.IComponent) ([]func(), error) {
	var emulation []func()
	for _, x := range seq {
		if comp, ok := components[x]; ok {
			if comp.EmulationRequired() {
				emulation = append(emulation, comp.Emulate)
			}
		}
	}
	if len(emulation) != len(seq) {
		return nil, fmt.Errorf("emulation sequence is not complete")
	}
	return emulation, nil
}
