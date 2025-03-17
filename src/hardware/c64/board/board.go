package board

import (
	"github.com/markel1974/c64emu/src/common/signals"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/hardware/c64/prg"
	"github.com/markel1974/c64emu/src/hardware/vic"
	"github.com/markel1974/c64emu/src/references"
	"golang.design/x/clipboard"
	"log"
	"os"
)

// Constants representing interrupt request bits for various hardware components.
const (
	intrIrqVicBit       = 0
	intrIrqCia1Bit      = 1
	intrIrqCia2Bit      = 2
	intrIrqExpansionBit = 3
)

// Board represents the central processing and coordination system for handling hardware components and peripherals.
type Board struct {
	*component.BaseComponent
	db                  references.IDisplayBuffer
	player              references.IPlayer
	quartz              references.IQuartz
	pic                 references.IPIC6510
	iec                 references.IIec
	keys                references.IKeyboard
	joy1                references.IJoystick
	joy2                references.IJoystick
	throttle            references.IThrottle
	factory             references.IComponentFactory
	cia1Socket          *CIA1Socket
	cia2Socket          *CIA2Socket
	vicSocket           *VICSocket
	sidSocket           *SIDSocket
	cpuSocket           *CPUSocket
	plaSocket           *PLASocket
	cartSocket          *CartridgeSocket
	joySwap             bool
	cfg                 *config.Config
	hasClipboard        bool
	expansionIrqTrigger *signals.SignalUint32
	expansionIrqClear   *signals.SignalUint32
	vBlank              bool
	lastVicCycle        bool
	dmaLow              bool
	prg                 *prg.PRG
}

// NewBoard initializes and returns a pointer to a new Board instance with default settings and dependencies.
func NewBoard(parent references.IComponent, factory references.IComponentFactory, suffix string) *Board {
	b := &Board{
		BaseComponent:       component.NewBaseComponent(componentId, suffix),
		factory:             factory,
		db:                  nil,
		player:              nil,
		quartz:              nil,
		iec:                 nil,
		pic:                 nil,
		keys:                nil,
		joy1:                nil,
		joy2:                nil,
		hasClipboard:        false,
		cartSocket:          nil,
		plaSocket:           nil,
		expansionIrqTrigger: nil,
		expansionIrqClear:   nil,
		vBlank:              false,
		lastVicCycle:        false,
		dmaLow:              false,
		prg:                 nil,
		joySwap:             true,
		throttle:            nil,
	}
	component.Register(parent, b)
	return b
}

func NewBoardComponent(parent references.IComponent, factory references.IComponentFactory, suffix string) references.IComponent {
	return NewBoard(parent, factory, suffix)
}

// Setup initializes the Board with provided display buffer, player, and configuration settings.
func (s *Board) Setup(db references.IDisplayBuffer, player references.IPlayer, cfg *config.Config) error {
	var err error
	s.db = db
	s.player = player
	if err = clipboard.Init(); err != nil {
		log.Printf("can't init clipboard: %s", err)
	} else {
		s.hasClipboard = true
	}
	s.cfg = cfg
	s.cfg.Bind(s.configChanged)

	s.cpuSocket = NewCPUSocket()
	s.vicSocket = NewVICSocket()
	s.sidSocket = NewSIDSocket()
	s.cia1Socket = NewCIA1Socket()
	s.cia2Socket = NewCIA2Socket()
	s.plaSocket = NewPLASocket()
	s.cartSocket = NewCartridgeSocket()

	var cpu references.I6510
	var vic references.IVIC
	var sid references.ISID
	var cia1 references.ICIA
	var cia2 references.ICIA
	var plaC references.IPlaC64
	var cart references.IExpansionSocketC64
	var rl references.IROMLoaderC64

	if _, s.quartz, err = s.factory.CreateIQuartz(s, "quartz", ""); err != nil {
		return err
	}
	if _, s.throttle, err = s.factory.CreateIThrottle(s, "dynamic_throttle", ""); err != nil {
		return err
	}
	if _, s.pic, err = s.factory.CreateIPIC6510(s, "pic_6510", ""); err != nil {
		return err
	}
	if _, cpu, err = s.factory.CreateI6510(s, "mos6510", ""); err != nil {
		return err
	}
	if _, s.iec, err = s.factory.CreateIEC(s, "iec", ""); err != nil {
		return err
	}
	if _, vic, err = s.factory.CreateIVIC(s, "mos6569", ""); err != nil {
		return err
	}
	if _, sid, err = s.factory.CreateISID(s, "mos6581", ""); err != nil {
		return err
	}
	if _, cart, err = s.factory.CreateIExpansionSocketC64(s, "cartridges_c64", ""); err != nil {
		return err
	}
	if _, rl, err = s.factory.CreateIROMLoaderC64(s, "roms_c64", ""); err != nil {
		return err
	}
	if _, cia1, err = s.factory.CreateICIA(s, "mos6526", "1"); err != nil {
		return err
	}
	if _, cia2, err = s.factory.CreateICIA(s, "mos6526", "2"); err != nil {
		return err
	}
	if _, plaC, err = s.factory.CreateIPLAc64(s, "pla_c64", ""); err != nil {
		return err
	}
	if _, s.keys, err = s.factory.CreateIKeyboard(s, "keyboard_c64", ""); err != nil {
		return err
	}
	if _, s.joy1, err = s.factory.CreateIJoystick(s, "joystick_c64", "1"); err != nil {
		return err
	}
	if _, s.joy2, err = s.factory.CreateIJoystick(s, "joystick_c64", "2"); err != nil {
		return err
	}
	expansion := NewExpansion(s, "")
	if err = expansion.Setup(s); err != nil {
		return err
	}

	s.throttle.SetInterval(mos6569.FrameInterval)

	if err = rl.Setup(cfg); err != nil {
		return err
	}
	if err = s.pic.Setup(s.quartz); err != nil {
		return err
	}
	if err = s.iec.Setup(s.quartz, cfg); err != nil {
		return err
	}
	if err = s.cpuSocket.Connect(s, cpu); err != nil {
		return err
	}
	if err = s.vicSocket.Connect(s, vic); err != nil {
		return err
	}
	if err = s.sidSocket.Connect(s, sid, mos6569.ScreenFreq, mos6569.TotalRasters, cfg); err != nil {
		return err
	}
	if err = s.cia1Socket.Connect(s, cia1); err != nil {
		return err
	}
	if err = s.cia2Socket.Connect(s, cia2); err != nil {
		return err
	}
	if err = s.cartSocket.Connect(s, cart, expansion, cfg); err != nil {
		return err
	}
	if err = s.plaSocket.Connect(s, plaC, vic, sid, cia1, cia2, s.cartSocket, rl); err != nil {
		return err
	}

	for _, cartName := range s.cfg.GetCartridges() {
		var data []uint8
		if len(cartName.Path) > 0 {
			if data, err = os.ReadFile(cartName.Path); err != nil {
				log.Printf("can't add cartridge: %s", err.Error())
				continue
			}
		}
		if cartId, err := s.cartSocket.Add(cartName.Kind, cartName.Path, data); err != nil {
			log.Printf("can't add cartridge: %s", err.Error())
		} else {
			log.Printf("cartridge: %s [%s] successfully added", cartName, cartId)
		}
	}

	if prgPath := s.cfg.GetPrg(); len(prgPath) > 0 {
		s.prg = prg.NewPRG(plaC, s.keys)
		if err := s.prg.Load(prgPath); err != nil {
			log.Printf("can't load prg: %s", err.Error())
			s.prg = nil
		}
	}

	s.reset()

	s.Print(os.Stdout, " ", true)
	//state, _ := s.DumpAll()
	//buf, _ := json.MarshalIndent(state, "", " ")
	//fmt.Println(string(buf))
	//if err := s.RestoreAll(state); err != nil {
	//	fmt.Println(err)
	//}

	return nil
}

func (s *Board) Reset() {
	//nothing to do
}

// reset resets the internal components of the Board to their initial states.
func (s *Board) reset() {
	s.pic.Reset()
	s.plaSocket.Reset()
	s.cpuSocket.Reset()
	s.cartSocket.Reset()
	s.sidSocket.Reset()
	s.cia1Socket.Reset()
	s.cia2Socket.Reset()
	s.iec.Reset()
	//s.expansion.Reset()
}

// AsyncReset resets various components of the board asynchronously by invoking their respective reset methods.
// It clears the low DMA state and initiates resets for the PLA, PIC, VIC, SID, CIA1, CIA2, Dispatcher, and Expansion modules.
func (s *Board) AsyncReset() {
	s.plaSocket.Reset()
	s.pic.TriggerReset()
	//s.cpuSocket.AsyncReset()
	s.vicSocket.Reset()
	s.sidSocket.Reset()
	s.cia1Socket.Reset()
	s.cia2Socket.Reset()
	s.iec.Reset()
	//s.expansion.Reset()

	s.dmaLow = false
}

// configChanged handles updates or changes in the configuration settings for the Board instance.
func (s *Board) configChanged() {
}

// Emulate executes the main emulation loop for the Board, coordinating CPU, VIC, CIA, and other components in sequence.
func (s *Board) Emulate() bool {
	s.vBlank = false

	//PHI1
	s.vicSocket.Emulate()

	//PHI2
	s.cia1Socket.Emulate()
	s.cia2Socket.Emulate()
	s.cartSocket.Emulate()
	s.iec.Emulate()
	s.cpuSocket.Emulate()

	s.quartz.AddCycle()

	return s.vBlank
}

// Throttle returns the throttling service used by the board.
func (s *Board) Throttle() references.IThrottle {
	return s.throttle
}

// GetText retrieves the text data from the VIC component of the Board instance. It returns the text as a byte slice.
func (s *Board) GetText() []byte {
	return s.vicSocket.GetText()
}

// GetScreenSize returns the dimensions of the screen as width and height in pixels.
func (s *Board) GetScreenSize() (int, int) {
	return mos6569.DisplayX, mos6569.DisplayY
}

// DiskChange triggers disk-related operations, such as switching the current disk and setting the drive configuration options.
func (s *Board) DiskChange() {
	//s.cfg.GetDrives()
	//aaa

	s.cfg.SwitchDisk()
	s.cfg.SetDriveOpt("", 8)
}

// KeyboardPaste checks if the paste button is pressed and clipboard data is available, then sets the keyboard command from the clipboard data.
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

// KeyboardSetCommand sets the specified command string to the keyboard input handler.
func (s *Board) KeyboardSetCommand(cmd string) {
	s.keys.SetCommand(cmd)
}

// KeyboardNumLockToggle toggles the state of the Num Lock key on the keyboard.
func (s *Board) KeyboardNumLockToggle() {
	s.keys.NumLockToggle()
}

// KeyboardCapitalToggle toggles the capital (Caps Lock) state of the keyboard.
func (s *Board) KeyboardCapitalToggle() {
	s.keys.CapitalToggle()
}

// SetMouse sets the mouse position by providing x and y coordinates to the SID socket.
func (s *Board) SetMouse(x uint8, y uint8) {
	s.sidSocket.SetPotXY(x, y)
}

// KeyboardSetVirtualKey modifies the state of a virtual key based on whether it is pressed or released.
func (s *Board) KeyboardSetVirtualKey(pressed bool, vKey int) {
	s.keys.SetVirtualKey(pressed, vKey)
}

// Joy1SetKey sets the state of a key input for joystick 1 or swaps to joystick 2 depending on the joySwap flag.
func (s *Board) Joy1SetKey(pressed bool, vKey int) {
	if s.joySwap {
		s.joy2.SetKey(pressed, vKey)
	} else {
		s.joy1.SetKey(pressed, vKey)
	}
}

// Joy2SetKey sets the state of a joystick key (pressed or released) for Joy2, swapping with Joy1 if joySwap is enabled.
func (s *Board) Joy2SetKey(pressed bool, vKey int) {
	if s.joySwap {
		s.joy1.SetKey(pressed, vKey)
	} else {
		s.joy2.SetKey(pressed, vKey)
	}
}

// Joystick1Move updates the position and button states of joystick 1, or joystick 2 if `joySwap` is enabled.
func (s *Board) Joystick1Move(x uint, y uint, buttons uint) {
	if s.joySwap {
		s.joy2.Move(x, y, buttons)
	} else {
		s.joy1.Move(x, y, buttons)
	}
}

// Joystick2Move moves the second joystick based on the given x, y coordinates and button presses.
// If joystick swapping is enabled, the move applies to the first joystick instead.
func (s *Board) Joystick2Move(x uint, y uint, buttons uint) {
	if s.joySwap {
		s.joy1.Move(x, y, buttons)
	} else {
		s.joy2.Move(x, y, buttons)
	}
}

// JoySwap toggles the joystick swap state and resets both joysticks if the passed button state is pressed.
func (s *Board) JoySwap(pressed bool) {
	if !pressed {
		return
	}
	s.joySwap = !s.joySwap
	s.joy1.Reset()
	s.joy2.Reset()
}

// ExtRamWrite writes a byte of data to an external RAM address with an optional memory configuration switch.
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

// ExtRamRead reads a byte from the external RAM at the specified address and memory configuration setting.
// If memConfig is non-negative, it temporarily applies the memory configuration, performs the read operation,
// and restores the previous memory configuration. Returns the retrieved byte value.
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

// dmaLowSlot sets the DMA low slot state for the board by updating internal flags and CPU control lines.
// If DMA is low, the CPU releases the bus for other devices to access memory via direct memory access (DMA).
// The method updates CPU RDY and AEC lines depending on the DMA state and VIC-II bus access status.
func (s *Board) dmaLowSlot(v bool) {
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

// rdyLowSlot updates the RDY signal based on the logical OR of the given value and the dmaLow state.
func (s *Board) rdyLowSlot(v bool) {
	//The RDY signal the result of logical AND between BA and DMA produced by the chip U27
	s.cpuSocket.SetRDYLow(v || s.dmaLow)
	//TODO SIGNAL
}

// aecLowSlot sets the AEC (Address Enable Control) line state based on the provided value and the DMA state.
func (s *Board) aecLowSlot(v bool) {
	s.cpuSocket.SetAECLow(v || s.dmaLow)
	//TODO SIGNAL
}

// irqTriggerSlot triggers an interrupt request (IRQ) specified by the given identifier and emits it to expansion if available.
func (s *Board) irqTriggerSlot(i uint32) {
	s.pic.TriggerIRQ(i)
	if s.expansionIrqTrigger != nil {
		s.expansionIrqTrigger.Emit(i)
	}
}

// irqClearSlot clears the specified IRQ signal on the PIC and emits a clear signal through expansionIrqClear if defined.
func (s *Board) irqClearSlot(i uint32) {
	s.pic.ClearIRQ(i)
	if s.expansionIrqClear != nil {
		s.expansionIrqClear.Emit(i)
	}
}

// nmiTriggerSlot triggers a Non-Maskable Interrupt (NMI) via the board's Programmable Interrupt Controller (PIC).
func (s *Board) nmiTriggerSlot() {
	s.pic.TriggerNMI()
}

// nmiClearSlot clears the Non-Maskable Interrupt (NMI) by utilizing the PIC's `ClearNMI` method.
func (s *Board) nmiClearSlot() {
	s.pic.ClearNMI()
}

// vicLastCycleSLot prepares the SID socket for operations during the last VIC-II cycle.
func (s *Board) vicLastCycleSLot() {
	s.sidSocket.Prepare()
}

// vicVBlankSlot handles the vertical blanking interval updates for connected components and processes injected programs if applicable.
func (s *Board) vicVBlankSlot() {
	s.vBlank = true
	s.sidSocket.Update()
	s.cia1Socket.Update()
	s.cia2Socket.Update()
	if s.prg != nil {
		if s.prg.Inject(s.vicSocket.GetText()) {
			s.prg = nil
		}
	}
}

// ledStateChangedSlot is triggered to handle changes in the LED state for a specified device.
func (s *Board) ledStateChangedSlot(_ int, _ uint8) {
	//TODO IMPLEMENT
	//deviceId := deviceNumber - 8
	//if deviceId < 0 || deviceId >= MAX_DRIVE_COUNT {
	//	return
	//}
	//k.leds[deviceId] = state
	//k.updateLedState()
}
