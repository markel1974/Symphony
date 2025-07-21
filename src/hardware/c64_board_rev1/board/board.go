package board

import (
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/hardware/c64_board_rev1/prg"
	"github.com/markel1974/c64emu/src/hardware/mos6569_vic_rev1"
	"github.com/markel1974/c64emu/src/references"
	"os"
)

// intrIrqVicBit represents the bit position for the VIC IRQ interrupt.
// intrIrqCia1Bit represents the bit position for the CIA1 IRQ interrupt.
// intrIrqCia2Bit represents the bit position for the CIA2 IRQ interrupt.
// intrIrqExpansionBit represents the bit position for the Expansion IRQ interrupt.
const (
	intrIrqVicBit       = 0
	intrIrqCia1Bit      = 1
	intrIrqCia2Bit      = 2
	intrIrqExpansionBit = 3
)

// Board represents the main structure of a system board with various hardware component sockets and configurations.
type Board struct {
	*component.BaseComponent
	keysSocket     *KeyboardSocket
	joySocket1     *JoystickSocket
	joySocket2     *JoystickSocket
	quartzSocket   *QuartzSocket
	ramSocket      *RamSocket
	colorRamSocket *ColorRamSocket
	romSocket      *RomSocket
	iecSocket      *IECSocket
	//picSocket       *PICSocket
	cia1Socket      *CIA1Socket
	cia2Socket      *CIA2Socket
	vicSocket       *VICSocket
	sidSocket       *SIDSocket
	cpuSocket       *CPUSocket
	plaSocket       *PLASocket
	throttleSocket  *ThrottleSocket
	cartridgeSocket *CartridgeManagerSocket
	prg             *prg.PRG
	joySwap         bool
	dmaLow          bool
	emulation       []func()
	sockets         []references.ISocket
	label           string
	connections     references.IC64BoardConnections
}

// NewBoard initializes and returns a new Board instance configured with various hardware sockets and components.
func NewBoard(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Board {
	s := &Board{
		BaseComponent: component.NewBaseComponent(),
		dmaLow:        false,
		prg:           nil,
		joySwap:       true,
		emulation:     []func(){},
		label:         label,
	}
	s.BaseComponent.Register(factory, parent, Identifier(), s, references.IdIC64Board(s, label, instance))

	s.keysSocket = NewKeyboardSocket(s, s.label)
	s.joySocket1 = NewJoystickSocket(s, s.label, 0)
	s.joySocket2 = NewJoystickSocket(s, s.label, 1)
	s.ramSocket = NewRamSocket(s, s.label)
	s.colorRamSocket = NewColorRamSocket(s, s.label)
	s.romSocket = NewRomSocket(s, s.label)
	s.quartzSocket = NewQuartzSocket(s, s.label)
	s.iecSocket = NewIECSocket(s, s.label, s)
	s.cartridgeSocket = NewExpansionSocket(s, s.label, s)
	//s.picSocket = NewPICSocket(s, s.label)
	s.cpuSocket = NewCPUSocket(s, s.label)
	s.vicSocket = NewVICSocket(s, s.label, s)
	s.sidSocket = NewSIDSocket(s, s.label, mos6569.ScreenFreq, mos6569.TotalRasters)
	s.cia1Socket = NewCIA1Socket(s, s.label, s)
	s.cia2Socket = NewCIA2Socket(s, s.label, s)
	s.plaSocket = NewPLASocket(s, s.label)
	s.throttleSocket = NewThrottleSocket(s, s.label, 1000/mos6569.ScreenFreq)

	s.sockets = append(s.sockets, s.romSocket)
	s.sockets = append(s.sockets, s.ramSocket)
	s.sockets = append(s.sockets, s.colorRamSocket)
	s.sockets = append(s.sockets, s.quartzSocket)
	s.sockets = append(s.sockets, s.keysSocket)
	s.sockets = append(s.sockets, s.joySocket1)
	s.sockets = append(s.sockets, s.joySocket2)
	s.sockets = append(s.sockets, s.iecSocket)
	s.sockets = append(s.sockets, s.cartridgeSocket)
	//s.sockets = append(s.sockets, s.picSocket)
	s.sockets = append(s.sockets, s.cpuSocket)
	s.sockets = append(s.sockets, s.vicSocket)
	s.sockets = append(s.sockets, s.sidSocket)
	s.sockets = append(s.sockets, s.cia1Socket)
	s.sockets = append(s.sockets, s.cia2Socket)
	s.sockets = append(s.sockets, s.plaSocket)
	s.sockets = append(s.sockets, s.throttleSocket)

	return s
}

// Setup initializes the Board by assigning its configuration using the associated factory method. Returns an error if any issue occurs.
func (s *Board) Setup() error {
	return nil
}

// Connect establishes connections to all sockets on the board and returns an error if any socket fails to mount.
func (s *Board) Connect() error {
	for _, c := range s.sockets {
		if err := c.Wire(); err != nil {
			return err
		}
	}
	return nil
}

// Wire establishes the provided IC64BoardConnections interface to the Board, enabling VBlank and LED activity management.
// It returns an error if the operation fails.
func (s *Board) Wire(conn references.IC64BoardConnections) error {
	s.connections = conn
	return nil
}

// Start initializes the board by setting up emulation, peripherals, cartridges, loading PRG if specified, and resetting the board.
// Returns an error if any step in the initialization process fails.
func (s *Board) Start() error {
	cfg := s.GetFactory().GetConfig()
	var err error
	if err = s.iecSocket.CreatePeripherals(); err != nil {
		return err
	}
	if err = s.cartridgeSocket.CreateCartridges(); err != nil {
		return err
	}
	if prgData := cfg.Prg(); len(prgData) > 0 {
		if err = s.startPRG(prgData); err != nil {
			return err
		}
	}
	if s.emulation, err = s.rebuildEmulation(); err != nil {
		return err
	}
	s.reset()
	s.Print(os.Stdout, " ", true)
	return nil
}

// Internal determines if the current state or configuration of the Board is considered internal or not. Returns a boolean.
func (s *Board) Internal() bool {
	return false
}

// Reset clears the current state of the Board and initializes it to its default state.
func (s *Board) Reset() {
	//nothing to do
}

// reset reinitializes all sockets of the Board, setting them to their default states.
func (s *Board) reset() {
	//s.picSocket.Reset()
	s.plaSocket.Reset()
	s.cpuSocket.Reset()
	s.cartridgeSocket.Reset()
	s.sidSocket.Reset()
	s.cia1Socket.Reset()
	s.cia2Socket.Reset()
	s.iecSocket.Reset()
}

// AsyncReset performs an asynchronous reset on various components of the board, initializing them to their default state.
func (s *Board) AsyncReset() {
	s.plaSocket.Reset()
	s.cpuSocket.TriggerReset()
	//s.cpuSocket.AsyncReset()
	s.vicSocket.Reset()
	s.sidSocket.Reset()
	s.cia1Socket.Reset()
	s.cia2Socket.Reset()
	s.iecSocket.Reset()

	s.dmaLow = false
}

// Emulate executes all functions stored in the emulation slice of the Board in sequence.
//
//go:nosplit
func (s *Board) Emulate() {
	for _, fn := range s.emulation {
		fn()
	}
}

// EmulationRequired checks if emulation is needed for the current Board instance. Returns true if emulation is required.
func (s *Board) EmulationRequired() bool {
	return true
}

// AECLow checks if the AEC (Automatic Exposure Control) is in the low state via the vicSocket connection.
func (s *Board) AECLow() bool {
	return s.vicSocket.GetAECLow()
}

// BALow checks if the 'BALow' status is set via the vicSocket and returns a boolean indicating the result.
func (s *Board) BALow() bool {
	return s.vicSocket.GetBALow()
}

// DMALowTrigger sets the DMA line to low or high, controlling CPU bus access and enabling other units to utilize the hardware.
// When set to low, the CPU halts after the next read cycle, and all bus lines shift to high impedance state.
// This state allows devices such as those on the expansion port to perform direct memory access (DMA) to the main RAM.
// Additionally, the CPU's AEC line is forced low, effectively placing the CPU in a wait state.
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

// RDYLowTrigger controls the RDY signal based on the logical conditions of BA and DMA as set by chip U27.
func (s *Board) RDYLowTrigger(v bool) {
	//The RDY signal the result of logical AND between BA and DMA produced by the chip U27
	s.cpuSocket.SetRDYLow(v || s.dmaLow)
}

// AECLowTrigger sets the AEC (Address Enable Control) low signal based on the given value and DMA low state.
func (s *Board) AECLowTrigger(v bool) {
	s.cpuSocket.SetAECLow(v || s.dmaLow)
}

// IRQTrigger triggers an interrupt request (IRQ) using the associated interrupt ID managed by the connection interface.
func (s *Board) IRQTrigger(d uint32) {
	s.cpuSocket.TriggerIRQ(d)
	s.cartridgeSocket.IRQ(d)
}

// IRQClearTrigger clears the interrupt request (IRQ) trigger for the provided device identifier.
func (s *Board) IRQClearTrigger(d uint32) {
	s.cpuSocket.ClearIRQ(d)
	s.cartridgeSocket.IRQClear(d)
}

// NMITrigger initiates a Non-Maskable Interrupt (NMI) on the board through the connected PIC socket.
func (s *Board) NMITrigger() {
	s.cpuSocket.TriggerNMI()
}

// NMIClearTrigger clears the NMI (Non-Maskable Interrupt) trigger by invoking the ClearNMI method on the PIC socket.
func (s *Board) NMIClearTrigger() {
	s.cpuSocket.ClearNMI()
}

// RSTTrigger triggers a reset signal using the PIC socket implementation within the Board structure.
func (s *Board) RSTTrigger() {
	s.cpuSocket.TriggerReset()
}

// LastCycleTrigger triggers the preparation of the socket in the last cycle of board operations.
func (s *Board) LastCycleTrigger() {
	s.sidSocket.Prepare()
}

// VBlankTrigger handles the vertical blank interrupt by emitting signals and updating connected hardware components.
func (s *Board) VBlankTrigger() {
	s.connections.VBlank()

	s.sidSocket.Update()
	//s.cia1Socket.PollInputs()
	//s.cia2Socket.Update()

	if s.prg != nil {
		if s.prg.Inject(s.vicSocket.GetText()) {
			s.prg = nil
		}
	}
	s.throttleSocket.Update()
}

// LedActivityTrigger emits a signal to change the state of the LED to the specified value.
func (s *Board) LedActivityTrigger(deviceNumber uint8, led bool) {
	s.connections.LedActivity(deviceNumber, led)
}

// GetText retrieves a byte slice containing the text data from the underlying vicSocket of the Board.
func (s *Board) GetText() []byte {
	return s.vicSocket.GetText()
}

// HardwareButton handles the state of a hardware button, passing the pressed state and a value to the cartridge socket.
func (s *Board) HardwareButton(pressed bool, val uint8) {
	s.cartridgeSocket.HardwareButton(pressed, val)
}

// KeyboardSetCommand sets the given command to the keysSocket for the Board.
func (s *Board) KeyboardSetCommand(cmd string) {
	s.keysSocket.SetCommand(cmd)
}

// KeyboardNumLockToggle toggles the Num Lock state on the keyboard by interacting with the keysSocket component.
func (s *Board) KeyboardNumLockToggle() {
	s.keysSocket.NumLockToggle()
}

// KeyboardCapitalToggle toggles the capitalization state of the keyboard keys via the keysSocket interface.
func (s *Board) KeyboardCapitalToggle() {
	s.keysSocket.CapitalToggle()
}

// SetMouse sets the mouse pointer coordinates using the specified x and y values.
func (s *Board) SetMouse(x uint8, y uint8) {
	s.sidSocket.SetPotX(x)
	s.sidSocket.SetPotY(y)
}

// KeyboardSetKey updates the state of a virtual key on the keyboard, specifying whether it is pressed or not.
func (s *Board) KeyboardSetKey(pressed bool, vKey int) {
	s.keysSocket.SetKey(pressed, vKey)
}

// Joy1SetKey updates the key state for joystick 1 based on whether it is pressed and the virtual key identifier.
func (s *Board) Joy1SetKey(pressed bool, vKey int) {
	if s.joySwap {
		s.joySocket2.SetKey(pressed, vKey)
	} else {
		s.joySocket1.SetKey(pressed, vKey)
	}
}

// Joy2SetKey sets the key state for the joystick, delegating to the appropriate joystick socket based on joySwap status.
func (s *Board) Joy2SetKey(pressed bool, vKey int) {
	if s.joySwap {
		s.joySocket1.SetKey(pressed, vKey)
	} else {
		s.joySocket2.SetKey(pressed, vKey)
	}
}

// Joystick1Move processes joystick movement and button press inputs for the first joystick based on the current socket configuration.
func (s *Board) Joystick1Move(x uint, y uint, buttons uint) {
	if s.joySwap {
		s.joySocket2.Move(x, y, buttons)
	} else {
		s.joySocket1.Move(x, y, buttons)
	}
}

// Joystick2Move processes joystick movement input and manages state based on the `joySwap` configuration of the board.
func (s *Board) Joystick2Move(x uint, y uint, buttons uint) {
	if s.joySwap {
		s.joySocket1.Move(x, y, buttons)
	} else {
		s.joySocket2.Move(x, y, buttons)
	}
}

// JoySwap toggles the joySwap state and resets both joySocket1 and joySocket2.
func (s *Board) JoySwap() {
	s.joySwap = !s.joySwap
	s.joySocket1.Reset()
	s.joySocket2.Reset()
}

// ExtRamWrite writes a byte of data to an external RAM address based on the specified memory configuration.
func (s *Board) ExtRamWrite(memConfig int, addr uint16, data uint8) {
	s.plaSocket.ExtWrite(memConfig, addr, data)
}

// ExtRamRead reads a byte from external RAM at the specified address using the given memory configuration.
func (s *Board) ExtRamRead(memConfig int, addr uint16) uint8 {
	return s.plaSocket.ExtRead(memConfig, addr)
}

// startPRG initializes and loads a PRG from the specified file path. It returns an error if loading the PRG fails.
func (s *Board) startPRG(prgPath []byte) error {
	s.prg = prg.NewPRG(s.plaSocket, s.keysSocket)
	if err := s.prg.Load(prgPath); err != nil {
		s.prg = nil
		return fmt.Errorf("can't load prg: %s", err.Error())
	}
	return nil
}

// rebuildEmulation constructs a sequence of emulation functions for hardware components in a specific order.
// It checks if emulation is required for each component and ensures all components in the sequence are accounted for.
// Returns an ordered slice of emulation functions or an error if the sequence is incomplete.
func (s *Board) rebuildEmulation() ([]func(), error) {
	hardwareSequence := []string{
		s.vicSocket.HardwareId(),
		s.cia1Socket.HardwareId(),
		s.cia2Socket.HardwareId(),
		s.cartridgeSocket.HardwareId(),
		s.iecSocket.HardwareId(),
		s.cpuSocket.HardwareId(),
		s.quartzSocket.HardwareId(),
	}

	var emulation []func()
	components := make(map[string]references.IComponent)
	for _, v := range s.GetChildren() {
		components[v.HardwareId()] = v
	}
	emulationCounter := 0
	for _, x := range hardwareSequence {
		if comp, ok := components[x]; ok {
			emulationCounter++
			if comp.EmulationRequired() {
				emulation = append(emulation, comp.Emulate)
			}
		}
	}
	if emulationCounter != len(hardwareSequence) {
		return nil, fmt.Errorf("emulation sequence is not complete")
	}
	return emulation, nil
}
